package execremote

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	scp "github.com/hnakamur/go-scp"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type (
	// 角色结构定义
	Role struct {
		Name       string
		Parallel   bool
		Hosts      []string
		WaitOutput bool
	}
	// 远程执行结构定义
	ExecRemote struct {
		user      string
		password  string
		timeout   time.Duration
		stdout    chan string
		role      Role
		sshConfig *ssh.ClientConfig
		clients   map[string]*ssh.Client
		cmdReturn []string
	}
)

// 实例化远程执行对象
func New(user, password string, timeout time.Duration, stdout chan string) *ExecRemote {
	er := &ExecRemote{
		user:     user,
		password: password,
		timeout:  timeout,
		stdout:   stdout,
		clients:  make(map[string]*ssh.Client),
	}
	er.insecureClientConfig()

	return er
}

func (er *ExecRemote) SetRole(r ...Role) {
	er.role = Role{}

	for _, role := range r {
		er.role.Name = role.Name
		er.role.Parallel = role.Parallel
		er.role.WaitOutput = role.WaitOutput
		for _, host := range role.Hosts {
			has := false
			for _, h := range er.role.Hosts {
				if h == host {
					has = true
					break
				}
			}
			if !has {
				er.role.Hosts = append(er.role.Hosts, host)
			}
		}
	}

	for _, host := range er.role.Hosts {
		hostArr := strings.Split(host, ":")
		if net.ParseIP(hostArr[0]) == nil {
			fmt.Printf("Host Error:%s\n", host)
			continue
		}

		if _, ok := er.clients[host]; ok {
			continue
		}

		client, err := ssh.Dial("tcp", host, er.sshConfig)
		if err != nil {
			writeChan(er.stdout, fmt.Sprintf("Ssh Dial Error:%s Host:%s\n", err.Error(), host))
			continue
		}
		er.clients[host] = client
	}
}

func (er *ExecRemote) Local(cmd string) ([]string, error) {
	var res []string

	// cmdArgs := strings.Fields(cmd)
	// cmdExec := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)
	cmdExec := exec.Command("bash", "-c", cmd)

	// Get a pipe to read from standard out
	cmdReader, err := cmdExec.StdoutPipe()

	// Use the same pipe for standard error
	cmdExec.Stderr = cmdExec.Stdout

	if err != nil {
		content := fmt.Sprintf("Exec Command StdoutPipe Error:%s\n", err.Error())
		res = append(res, content)
		writeChan(er.stdout, content)
		return res, err
	}

	// Make a new channel which will be used to ensure we get all output
	done := make(chan struct{})

	// Create a scanner which scans r in a line-by-line fashion
	scanner := bufio.NewScanner(cmdReader)
	// 可以在调用Scan之前设置buffer和MaxScanTokenSize的大小 默认值：MaxScanTokenSize = 64 * 1024
	scanner.Buffer([]byte{}, bufio.MaxScanTokenSize*128) // 8M
	scanner.Split(bufio.ScanLines)

	// Use the scanner to scan the output line by line and log it
	// It's running in a goroutine so that it doesn't block
	go func() {
		// Read line by line and process it
		for scanner.Scan() {
			line := scanner.Text()
			res = append(res, line, "\n")
			writeChan(er.stdout, line)
		}
		//  通过 scanner.Err(); 我们可以捕捉到 扫描中的错误信息,这对单行文件超过 MaxScanTokenSize 时特别有用
		if err := scanner.Err(); err != nil {
			content := fmt.Sprintf("Exec Command Reading Standard Input Error:%s\n", err.Error())
			res = append(res, content)
			writeChan(er.stdout, content)
		}

		// We're all done, unblock the channel
		done <- struct{}{}
	}()

	// Start the command and check for errors
	err = cmdExec.Start()
	if err != nil {
		content := fmt.Sprintf("Exec Command Starting Error:%s\n", err.Error())
		res = append(res, content)
		writeChan(er.stdout, content)
		return res, err
	}

	// Wait for all output to be processed
	<-done

	// Wait for the command to finish
	err = cmdExec.Wait()
	if err != nil {
		content := fmt.Sprintf("Exec Command Waiting Error:%s\n", err.Error())
		res = append(res, content)
		writeChan(er.stdout, content)
		return res, err
	}

	return res, nil
}

func (er *ExecRemote) Run(cmds ...string) {
	er.cmdReturn = []string{}
	var wg sync.WaitGroup

	for _, host := range er.role.Hosts {
		client, ok := er.clients[host]
		if !ok {
			continue
		}

		if er.role.WaitOutput {
			content, _ := er.runCmd(client, host, cmds...)
			er.cmdReturn = append(er.cmdReturn, content...)
		} else {
			wg.Add(1)
			if er.role.Parallel {
				go func(client *ssh.Client, host string, cmds ...string) {
					defer wg.Done()
					er.runCmd(client, host, cmds...)
				}(client, host, cmds...)
			} else {
				er.runCmd(client, host, cmds...)
				wg.Done()
			}
		}
	}

	if !er.role.WaitOutput {
		wg.Wait()
	}
}

func (er *ExecRemote) GetCmdReturn() []string {
	var resArr []string
	substr := "-> "
	n := len(er.cmdReturn)
	for i := 0; i < n; i++ {
		if i < 4 || i > n-3 {
			continue
		}
		if strings.Contains(er.cmdReturn[i], substr) {
			resArr = append(resArr, strings.Trim(strings.Split(er.cmdReturn[i], substr)[1], "\r"))
		}
	}
	return resArr
}

func (er *ExecRemote) Put(localPath, remotePath string) {
	var wg sync.WaitGroup
	for _, host := range er.role.Hosts {
		client, ok := er.clients[host]
		if !ok {
			continue
		}

		wg.Add(1)
		if er.role.Parallel {
			go func(client *ssh.Client, host, localPath, remotePath string) {
				defer wg.Done()
				er.runPut(client, host, localPath, remotePath)
			}(client, host, localPath, remotePath)
		} else {
			er.runPut(client, host, localPath, remotePath)
			wg.Done()
		}
	}
	wg.Wait()
}

func (er *ExecRemote) Get(remotePath, localPath string) {
	var wg sync.WaitGroup
	for _, host := range er.role.Hosts {
		client, ok := er.clients[host]
		if !ok {
			continue
		}

		wg.Add(1)
		if er.role.Parallel {
			go func(client *ssh.Client, host, remotePath, localPath string) {
				defer wg.Done()
				er.runGet(client, host, remotePath, localPath)
			}(client, host, remotePath, localPath)
		} else {
			er.runGet(client, host, remotePath, localPath)
			wg.Done()
		}
	}
	wg.Wait()
}

func (er *ExecRemote) Close() {
	for _, client := range er.clients {
		client.Close()
	}
	if !isStringChanClosed(er.stdout) {
		close(er.stdout)
	}
}

func (er *ExecRemote) runCmd(client *ssh.Client, host string, cmds ...string) ([]string, error) {
	var res []string

	// Create sesssion
	session, err := client.NewSession()

	// Request pseudo terminal
	fileDescriptor := int(os.Stdin.Fd())
	if terminal.IsTerminal(fileDescriptor) {
		originalState, err := terminal.MakeRaw(fileDescriptor)
		if err != nil {
			content := fmt.Sprintf("Terminal MakeRaw Failed:%s Host:%s\n", err.Error(), host)
			res = append(res, content)
			writeChan(er.stdout, content)
		}
		defer terminal.Restore(fileDescriptor, originalState)

		termWidth, termHeight, err := terminal.GetSize(fileDescriptor)
		if err != nil {
			content := fmt.Sprintf("Terminal GetSize Failed:%s Host:%s\n", err.Error(), host)
			res = append(res, content)
			writeChan(er.stdout, content)
		}

		// Set up terminal modes
		modes := ssh.TerminalModes{
			ssh.ECHO:          1, // enable echoing
			ssh.ECHOCTL:       1,
			ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
			ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
		}
		err = session.RequestPty("xterm-256color", termHeight, termWidth, modes)
		if err != nil {
			content := fmt.Sprintf("Request for pseudo terminal failed:%s Host:%s\n", err.Error(), host)
			res = append(res, content)
			writeChan(er.stdout, content)
		}
	}

	if err != nil {
		content := fmt.Sprintf("Client Session Error:%s Host:%s\n", err.Error(), host)
		res = append(res, content)
		writeChan(er.stdout, content)
		return res, err
	}
	defer session.Close()

	// StdinPipee() returns a pipe that will be connected to the remote command's standard input when the command starts.
	// StdoutPipe() returns a pipe that will be connected to the remote command's standard output when the command starts.
	stdin, err := session.StdinPipe()
	if err != nil {
		content := fmt.Sprintf("Session Stdin Error:%s Host:%s\n", err.Error(), host)
		res = append(res, content)
		writeChan(er.stdout, content)
		return res, err
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		content := fmt.Sprintf("Session Stdout Error:%s Host:%s\n", err.Error(), host)
		res = append(res, content)
		writeChan(er.stdout, content)
		return res, err
	}

	// Start remote shell
	err = session.Shell()
	if err != nil {
		content := fmt.Sprintf("Session Shell Error:%s Host:%s\n", err.Error(), host)
		res = append(res, content)
		writeChan(er.stdout, content)
		return res, err
	}

	// Section 3: Send the commands to the remotehost one by one.
	for _, cmd := range cmds {
		n, err := stdin.Write([]byte(cmd + "\r"))
		if err != nil {
			content := fmt.Sprintf("Stdin Write Error:%s Host:%s Cmd:%s\n", err.Error(), host, cmd)
			res = append(res, content)
			writeChan(er.stdout, content)
			continue
		}

		// Error handeling: Check the number of byte is sent
		if n-1 != len(cmd) {
			content := fmt.Sprintf("Stdin Write Error:For the Host %s The command %s is %d byte but %d byte is sent to the device\n", host, cmd, len(cmd), n)
			res = append(res, content)
			writeChan(er.stdout, content)
		}
	}
	// 命令执行完成退出
	stdin.Write([]byte("exit\r"))

	var b bytes.Buffer
	session.Stdout = &b

	scanner := bufio.NewScanner(stdout)
	// 可以在调用Scan之前设置buffer和MaxScanTokenSize的大小 默认值：MaxScanTokenSize = 64 * 1024
	scanner.Buffer([]byte{}, bufio.MaxScanTokenSize*128) // 8M
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		s := scanner.Text()
		if s != "\n" {
			b.WriteString(s)
		} else {
			content := fmt.Sprintf("[%s]-> %s", host, b.String())
			res = append(res, content)
			writeChan(er.stdout, content)
			b.Reset()
		}
		// 通过 scanner.Err(); 我们可以捕捉到 扫描中的错误信息,这对单行文件超过 MaxScanTokenSize 时特别有用
		if err := scanner.Err(); err != nil {
			content := fmt.Sprintf("Scanner reading standard input:%s Host:%s\n", err.Error(), host)
			res = append(res, content)
			writeChan(er.stdout, content)
		}
	}

	if er.role.WaitOutput {
		session.Wait()
	}

	return res, nil
}

func (er *ExecRemote) runPut(client *ssh.Client, host, localPath, remotePath string) {
	acceptFn := func(parentDir string, info os.FileInfo) (bool, error) {
		// str := "处理文件"
		// if info.IsDir() {
		// 	str = "处理目录"
		// }
		// fmt.Println(fmt.Sprintf("%s：%s/%s", str, parentDir, info.Name()))
		return true, nil
	}

	f, err := os.Stat(localPath)
	if err != nil {
		writeChan(er.stdout, fmt.Sprintf("Os Stat Error:%s Host:%s\n", err.Error(), host))
		return
	}
	if f.IsDir() {
		err = scp.NewSCP(client).SendDir(localPath, remotePath, acceptFn)
		if err != nil {
			writeChan(er.stdout, fmt.Sprintf("Scp SendDir Error:%s Host:%s\n", err.Error(), host))
			return
		}
		return
	}
	err = scp.NewSCP(client).SendFile(localPath, remotePath)
	if err != nil {
		writeChan(er.stdout, fmt.Sprintf("Scp SendFile Error:%s Host:%s\n", err.Error(), host))
		return
	}
}

func (er *ExecRemote) runGet(client *ssh.Client, host, remotePath, localPath string) {
	acceptFn := func(parentDir string, info os.FileInfo) (bool, error) {
		// str := "处理文件"
		// if info.IsDir() {
		// 	str = "处理目录"
		// }
		// fmt.Println(fmt.Sprintf("%s：%s/%s", str, parentDir, info.Name()))
		return true, nil
	}

	f, err := os.Stat(localPath)
	if err != nil {
		writeChan(er.stdout, fmt.Sprintf("Os Stat Error:%s Host:%s\n", err.Error(), host))
		return
	}
	if f.IsDir() {
		err = scp.NewSCP(client).ReceiveDir(remotePath, localPath, acceptFn)
		if err != nil {
			writeChan(er.stdout, fmt.Sprintf("Scp ReceiveDir Error:%s Host:%s\n", err.Error(), host))
			return
		}
		return
	}
	err = scp.NewSCP(client).ReceiveFile(remotePath, localPath)
	if err != nil {
		writeChan(er.stdout, fmt.Sprintf("Scp ReceiveFile Error:%s Host:%s\n", err.Error(), host))
		return
	}
}

func (er *ExecRemote) runPutScp(client *ssh.Client, host, localPath, remotePath string) {
	// open an SFTP session over an existing ssh connection.
	sftp, err := sftp.NewClient(client)
	if err != nil {
		writeChan(er.stdout, fmt.Sprintf("Sftp NewClient Error:%s Host:%s\n", err.Error(), host))
		return
	}
	defer sftp.Close()

	// Open the source file
	srcFile, err := os.Open(localPath)
	if err != nil {
		writeChan(er.stdout, fmt.Sprintf("Os Open Error:%s Host:%s\n", err.Error(), host))
		return
	}
	defer srcFile.Close()

	finfo, _ := os.Stat(localPath)
	mode := finfo.Mode()

	// Create the destination file
	dstFile, err := sftp.Create(remotePath)
	if err != nil {
		writeChan(er.stdout, fmt.Sprintf("Sftp Create Error:%s Host:%s\n", err.Error(), host))
		return
	}
	defer dstFile.Close()

	// write to file
	if _, err := dstFile.ReadFrom(srcFile); err != nil {
		writeChan(er.stdout, fmt.Sprintf("Sftp ReadFrom Error:%s Host:%s\n", err.Error(), host))
		return
	}
	if err := dstFile.Chmod(mode); err != nil {
		writeChan(er.stdout, fmt.Sprintf("Sftp Chmod Error:%s Host:%s\n", err.Error(), host))
		return
	}
}

func (er *ExecRemote) insecureClientConfig() {
	er.sshConfig = &ssh.ClientConfig{
		User:    er.user,
		Timeout: er.timeout,
		Auth:    []ssh.AuthMethod{ssh.Password(er.password)},

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Config: ssh.Config{
			Ciphers: []string{
				"aes128-ctr", "aes192-ctr", "aes256-ctr",
				"aes128-gcm@openssh.com", "aes256-gcm@openssh.com",
			},
			KeyExchanges: []string{
				"curve25519-sha256", "curve25519-sha256@libssh.org",
				"ecdh-sha2-nistp256", "ecdh-sha2-nistp384", "ecdh-sha2-nistp521",
				"diffie-hellman-group-exchange-sha256",
				"diffie-hellman-group16-sha512",
				"diffie-hellman-group18-sha512",
			},
		},
	}
}

func writeChan(ch chan string, content string) {
	// recover from panic caused by writing to a closed channel
	if r := recover(); r != nil {
		err := fmt.Errorf("%v", r)
		fmt.Fprintf(os.Stderr, "write: error writing %s on channel: %v\n", content, err)
		return
	}
	ch <- content
}

func isStringChanClosed(ch <-chan string) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}
