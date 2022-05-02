package installk8s

import (
	"fmt"
	"net"
	"strings"
	"time"

	"git.yonyou.com/sysbase/backend/model"
	"git.yonyou.com/sysbase/backend/tool/execremote"
)

type (
	InstallK8s struct {
		SourceDir string
		Params    Params
		Stdout    chan string
		Defer     func()
		AddNew    bool
		OnlyConf  bool
		UseLvs    bool
		er        *execremote.ExecRemote
		resources map[string]execremote.Role
	}
	Params struct {
		K8sClusterID uint
		DoWhat       string
		DoDocker     bool
	}
)

var localIps []string

func (ik *InstallK8s) Call(funcName string) {
	defer func() {
		ik.Defer()
		ik.er.Close()
	}()

	ik.Stdout <- "开始执行..."

	switch funcName {
	case "InstallTest":
		ik.InstallTest()
		break
	case "InstallAll":
		ik.InstallAll()
		break
	case "InstallBase":
		ik.InstallBase()
		break
	case "UpdateKernel":
		ik.UpdateKernel()
		break
	case "InstallBin":
		ik.InstallBin()
		break
	case "InstallDocker":
		ik.InstallDocker()
		break
	case "InstallPriDocker":
		ik.InstallPriDocker()
		break
	case "InstallEtcd":
		ik.InstallEtcd()
		break
	case "InstallMaster":
		ik.InstallMaster()
		break
	case "InstallNode":
		ik.InstallNode()
		break
	case "InstallDockerCrt":
		ik.InstallDockerCrt()
		break
	case "InstallLvs":
		ik.InstallLvs()
		break
	case "InstallDns":
		ik.InstallDns()
		break
	case "ServicePublish":
		ik.ServicePublish()
		break
	case "ServiceEtcd":
		ik.ServiceEtcd()
		break
	case "ServiceMaster":
		ik.ServiceMaster()
		break
	case "ServiceNode":
		ik.ServiceNode()
		break
	case "ServiceDns":
		ik.ServiceDns()
		break
	case "FinishInstall":
		ik.FinishInstall()
		break
	case "NewnodeInstall":
		ik.NewnodeInstall()
		break
	case "NewetcdInstall":
		ik.NewetcdInstall()
		break
	case "NewmasterInstall":
		ik.NewmasterInstall()
		break
	case "UpdateSslMaster":
		ik.UpdateSslMaster()
		break
	case "UpdateSslEtcd":
		ik.UpdateSslEtcd()
		break
	case "UpdateSslNode":
		ik.UpdateSslNode()
		break
	default:
		ik.Stdout <- fmt.Sprintf("没有%s的Call方法", funcName)
	}

	ik.Stdout <- "结束执行"
	time.Sleep(1 * time.Millisecond)
}

func (ik *InstallK8s) GetResources() {
	localIps, _ = getLocalIp()

	k8sClusterResource := model.K8sClusterResource{}

	resource, err := k8sClusterResource.ListResource(ik.Params.K8sClusterID, []string{})
	if err != nil {
		ik.Stdout <- err.Error()
		return
	}

	var user, password string
	ik.resources = make(map[string]execremote.Role)
	for _, r := range resource {
		user = r.User
		password = r.Password
		role := ik.resources[r.Scope]
		role.Name = r.Scope
		role.Parallel = true
		role.Hosts = append(role.Hosts, fmt.Sprintf("%s:%d", r.Host, r.Port))
		ik.resources[r.Scope] = role
	}

	timeout := 3 * time.Second
	ik.er = execremote.New(user, password, timeout, ik.Stdout)
}

func (ik *InstallK8s) InstallTest() {
	publishRole, ok := ik.resources["default"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	publishRole.WaitOutput = true
	ik.er.SetRole(publishRole)
	ik.er.Run(`cat /root/1.sh`)
	fmt.Printf("%#v\n", ik.er.GetCmdReturn())

	ik.er.Run(`uname -a | grep 4.19.94 > /dev/null ; echo $?`)
	fmt.Printf("%#v\n", ik.er.GetCmdReturn())
	return

	cmds := []string{
		`ls -la`,
		`date`,
		`for i in {1..10};do echo "你好${i}";sleep 1;done`,
	}
	ik.er.Run(cmds...)
}

func (ik *InstallK8s) InstallAll() {
	var hosts []string
	if r, ok := ik.resources["publish"]; ok {
		hosts = append(hosts, r.Hosts...)
	}
	if r, ok := ik.resources["etcd"]; ok {
		hosts = append(hosts, r.Hosts...)
	}
	if r, ok := ik.resources["master"]; ok {
		hosts = append(hosts, r.Hosts...)
	}
	if r, ok := ik.resources["node"]; ok {
		hosts = append(hosts, r.Hosts...)
	}
	if r, ok := ik.resources["pridocker"]; ok {
		hosts = append(hosts, r.Hosts...)
	}
	if r, ok := ik.resources["lvs"]; ok {
		hosts = append(hosts, r.Hosts...)
	}

	needUpdateHosts := ik.checkHostsKernel(hosts)
	if len(needUpdateHosts) > 0 {
		ik.Stdout <- fmt.Sprintf("请先确保这些机器内核已升级:%v\n", needUpdateHosts)
		ik.UpdateKernel()
		return
	}
	ik.InstallBase()
	ik.InstallBin()
	ik.InstallDocker()
	ik.InstallPriDocker()
	ik.InstallEtcd()
	ik.InstallMaster()
	ik.InstallNode()
	ik.InstallDockerCrt()
	ik.InstallLvs()
	ik.InstallDns()
	ik.FinishInstall()
}

func (ik *InstallK8s) InstallBase() {
	var role []execremote.Role

	if r, ok := ik.resources["publish"]; ok {
		role = append(role, r)
	}
	if r, ok := ik.resources["etcd"]; ok {
		role = append(role, r)
	}
	if r, ok := ik.resources["master"]; ok {
		role = append(role, r)
	}
	if r, ok := ik.resources["node"]; ok {
		role = append(role, r)
	}
	if r, ok := ik.resources["pridocker"]; ok {
		role = append(role, r)
	}
	if r, ok := ik.resources["lvs"]; ok {
		role = append(role, r)
	}

	ik.er.SetRole(role...)
	ik.installBase()
}

func (ik *InstallK8s) installBase() {
	cmds := []string{
		fmt.Sprintf(`chown -R root:root %s`, ik.SourceDir),
		"yum install -y telnet net-tools openssl socat",
		"mkdir /data > /dev/null 2>&1;if [ $? == 0 ];then useradd -d /data/www esn && useradd -d /data/www www && usermod -G esn www && chmod 750 /data/www && mkdir -p /data/log/php && mkdir -p /data/log/nginx && mkdir -p /data/yy_log && chown -R www:www /data/log /data/yy_log && chmod 750 /data/log /data/yy_log;fi",
		"systemctl stop firewalld && systemctl disable firewalld",
		`sed -i "s#SELINUX=enforcing#SELINUX=disabled#g" /etc/selinux/config && setenforce 0`,
		`cat /etc/sysctl.conf | grep net.ipv4.ip_forward > /dev/null 2>&1 ; if [ $? -ne 0 ];then echo "net.ipv4.ip_forward = 1" >> /etc/sysctl.conf && sysctl -p;fi`,
		`cat /etc/sysctl.conf | grep net.ipv4.conf.all.rp_filter > /dev/null 2>&1 ; if [ $? -ne 0 ];then echo "net.ipv4.conf.all.rp_filter = 1" >> /etc/sysctl.conf && sysctl -p;fi`,
	}
	ik.er.Run(cmds...)
}

func (ik *InstallK8s) UpdateKernel() {
	var role []execremote.Role

	if r, ok := ik.resources["etcd"]; ok {
		role = append(role, r)
	}
	if r, ok := ik.resources["master"]; ok {
		role = append(role, r)
	}
	if r, ok := ik.resources["node"]; ok {
		role = append(role, r)
	}
	if r, ok := ik.resources["pridocker"]; ok {
		role = append(role, r)
	}
	if r, ok := ik.resources["lvs"]; ok {
		role = append(role, r)
	}

	ik.er.SetRole(role...)
	ik.updateKernel()
}

func (ik *InstallK8s) checkHostsKernel(hosts []string) []string {
	// 检查每个节点的内核，如果不是最新内核，则升级内核
	var needUpdateHosts []string
	for _, host := range hosts {
		oneHostRole := execremote.Role{
			Hosts:      []string{host},
			WaitOutput: true,
		}
		ik.er.SetRole(oneHostRole)
		ik.er.Run(`uname -a | grep 4.19.94 > /dev/null ; echo $?`)
		if ik.er.GetCmdReturn()[0] == "1" {
			oneHostRole.WaitOutput = false
			needUpdateHosts = append(needUpdateHosts, host)
		}
	}

	return needUpdateHosts
}

func (ik *InstallK8s) updateKernel() {
	ik.er.Put(fmt.Sprintf("%s/kernel-4.19.94.gz", ik.SourceDir), "/tmp/kernel-4.19.94.gz")
	ik.er.Run("cd /tmp && tar zxvf kernel-4.19.94.gz && yum remove -y kernel-headers kernel-tools-libs kernel-tools kernel-ml-tools kernel-ml-tools-libs && yum install -y kernel-4.19.94/* ; rm -rf kernel-4.19.94*")
	ik.er.Run(`num=$(awk -F \' '$1=="menuentry " {print i++ " : " $2}' /etc/grub2.cfg | grep 4.19.94 | awk '{print $1}') && grub2-set-default $num && grub2-mkconfig -o /boot/grub2/grub.cfg ; grub2-editenv list`)
	ik.er.Run("reboot")
}

func (ik *InstallK8s) InstallBin() {
	var role []execremote.Role

	if r, ok := ik.resources["publish"]; ok {
		role = append(role, r)
	}

	ik.er.SetRole(role...)

	ik.er.Put(fmt.Sprintf("%s/bin.gz", ik.SourceDir), "/tmp")
	ik.er.Run("tar zxvf /tmp/bin.gz -C / && rm -rf /tmp/bin.gz")

	if r, ok := ik.resources["etcdlb"]; ok {
		host := strings.Split(r.Hosts[0], ":")
		cmds := []string{
			"mkdir -p /etc/calico",
			fmt.Sprintf(`echo 'apiVersion: projectcalico.org/v3
kind: CalicoAPIConfig
metadata:
spec:
  etcdEndpoints: "https://%s:2379"
  etcdKeyFile: "/etc/cni/net.d/calico-tls/etcd-key"
  etcdCertFile: "/etc/cni/net.d/calico-tls/etcd-cert"
  etcdCACertFile: "/etc/cni/net.d/calico-tls/etcd-ca"' > /etc/calico/calicoctl.cfg`, host[0]),
		}
		ik.er.Run(cmds...)
	}
}

func (ik *InstallK8s) InstallLvs() {
	if !ik.UseLvs {
		return
	}

	etcdRole, ok := ik.resources["etcd"]
	if !ok {
		ik.Stdout <- "没有etcd资源"
		return
	}

	masterRole, ok := ik.resources["master"]
	if !ok {
		ik.Stdout <- "没有master资源"
		return
	}

	ik.installLvs()
	ik.installLvsvipEtcd(etcdRole)
	ik.installLvsvipMaster(masterRole)
}

func (ik *InstallK8s) installLvs() {
	if !ik.UseLvs {
		return
	}

	lvsRole, ok := ik.resources["lvs"]
	if !ok {
		ik.Stdout <- "没有lvs资源"
		return
	}

	etcdLbRole, ok := ik.resources["etcdlb"]
	if !ok {
		ik.Stdout <- "没有etcdlb资源"
		return
	}

	masterLbRole, ok := ik.resources["masterlb"]
	if !ok {
		ik.Stdout <- "没有masterlb资源"
		return
	}

	etcdRole, ok := ik.resources["etcd"]
	if !ok {
		ik.Stdout <- "没有etcd资源"
		return
	}

	masterRole, ok := ik.resources["master"]
	if !ok {
		ik.Stdout <- "没有master资源"
		return
	}

	ik.er.SetRole(lvsRole)

	etcdLbHost := strings.Split(etcdLbRole.Hosts[0], ":")[0]
	masterLbHost := strings.Split(masterLbRole.Hosts[0], ":")[0]

	cmds := []string{
		`yum install -y ipvsadm && systemctl enable ipvsadm`,
		fmt.Sprintf(`ifconfig eth0:lvs:0 %s broadcast %s netmask 255.255.255.255 up && echo -e "#/bin/sh\n# chkconfig:   2345 90 10\nifconfig eth0:lvs:0 %s broadcast %s netmask 255.255.255.255 up" > /etc/rc.d/init.d/vip_route_lvs.sh`, etcdLbHost, etcdLbHost, etcdLbHost, etcdLbHost),
		fmt.Sprintf(`ifconfig eth0:lvs:1 %s broadcast %s netmask 255.255.255.255 up && echo "ifconfig eth0:lvs:1 %s broadcast %s netmask 255.255.255.255 up" >> /etc/rc.d/init.d/vip_route_lvs.sh`, masterLbHost, masterLbHost, masterLbHost, masterLbHost),
		fmt.Sprintf(`route add -host %s dev eth0:lvs:0 ; echo "" > /dev/null && echo "route add -host %s dev eth0:lvs:0 ; echo "" > /dev/null" >> /etc/rc.d/init.d/vip_route_lvs.sh`, etcdLbHost, etcdLbHost),
		fmt.Sprintf(`route add -host %s dev eth0:lvs:1 ; echo "" > /dev/null && echo "route add -host %s dev eth0:lvs:1 ; echo "" > /dev/null" >> /etc/rc.d/init.d/vip_route_lvs.sh`, masterLbHost, masterLbHost),
		`chmod +x /etc/rc.d/init.d/vip_route_lvs.sh && chkconfig --add vip_route_lvs.sh && chkconfig vip_route_lvs.sh on`,
		`echo "1" > /proc/sys/net/ipv4/ip_forward`,
	}
	ik.er.Run(cmds...)

	// etcd
	ipvsadm := fmt.Sprintf(`-A -t %s:2379 -s wrr\n`, etcdLbHost)
	for _, host := range etcdRole.Hosts {
		curHost := strings.Split(host, ":")[0]
		ipvsadm = fmt.Sprintf(`%s-a -t %s:2379 -r %s:2379 -g -w 1\n`, ipvsadm, etcdLbHost, curHost)
	}

	// master
	ipvsadm = fmt.Sprintf(`%s-A -t %s:6443 -s wrr\n`, ipvsadm, masterLbHost)
	for _, host := range masterRole.Hosts {
		curHost := strings.Split(host, ":")[0]
		ipvsadm = fmt.Sprintf(`%s-a -t %s:6443 -r %s:6443 -g -w 1\n`, ipvsadm, masterLbHost, curHost)
	}

	cmds = []string{
		fmt.Sprintf(`echo "%s" > /etc/sysconfig/ipvsadm`, ipvsadm),
		`systemctl restart ipvsadm && ipvsadm -Ln`,
	}

	ik.er.Run(cmds...)
}

func (ik *InstallK8s) installLvsNew() {
	if !ik.UseLvs {
		return
	}

	etcdLbRole, ok := ik.resources["etcdlb"]
	if !ok {
		ik.Stdout <- "没有etcdlb资源"
		return
	}

	masterLbRole, ok := ik.resources["masterlb"]
	if !ok {
		ik.Stdout <- "没有masterlb资源"
		return
	}

	lvsRole, ok := ik.resources["lvs"]
	if !ok {
		ik.Stdout <- "没有lvs资源"
		return
	}

	etcdRole, ok := ik.resources["etcd"]
	if !ok {
		ik.Stdout <- "没有etcd资源"
		return
	}

	masterRole, ok := ik.resources["master"]
	if !ok {
		ik.Stdout <- "没有master资源"
		return
	}

	etcdLbHost := strings.Split(etcdLbRole.Hosts[0], ":")[0]
	masterLbHost := strings.Split(masterLbRole.Hosts[0], ":")[0]

	ik.er.SetRole(lvsRole)
	ik.er.Run("systemctl stop ipvsadm")

	// etcd
	ipvsadm := fmt.Sprintf(`-A -t %s:2379 -s wrr\n`, etcdLbHost)
	var etcdHosts []string
	etcdHosts = append(etcdHosts, etcdRole.Hosts...)
	if newetcdRole, ok := ik.resources["newetcd"]; ok {
		etcdHosts = append(etcdHosts, newetcdRole.Hosts...)
	}
	for _, host := range etcdHosts {
		curHost := strings.Split(host, ":")[0]
		ipvsadm = fmt.Sprintf(`%s-a -t %s:2379 -r %s:2379 -g -w 1\n`, ipvsadm, etcdLbHost, curHost)
	}

	// master
	ipvsadm = fmt.Sprintf(`-A -t %s:6443 -s wrr\n'`, masterLbHost)
	var masterHosts []string
	masterHosts = append(masterHosts, masterRole.Hosts...)
	if newmasterRole, ok := ik.resources["newmaster"]; ok {
		masterHosts = append(masterHosts, newmasterRole.Hosts...)
	}
	for _, host := range masterHosts {
		curHost := strings.Split(host, ":")[0]
		ipvsadm = fmt.Sprintf(`%s-a -t %s:6443 -r %s:6443 -g -w 1\n`, ipvsadm, masterLbHost, curHost)
	}

	cmds := []string{
		fmt.Sprintf(`echo "%s" > /etc/sysconfig/ipvsadm`, ipvsadm),
		`systemctl start ipvsadm && ipvsadm -Ln`,
	}
	ik.er.Run(cmds...)
}

func (ik *InstallK8s) InstallDns() {
	pridnsRole, ok := ik.resources["pridns"]
	if !ok {
		ik.Stdout <- "没有pridns资源"
		return
	}

	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	ik.er.SetRole(pridnsRole)
	ik.er.Run("yum install -y bind-chroot")

	ik.er.SetRole(publishRole)
	ik.er.Run(fmt.Sprintf(`cd %s/bind && tar zcvf bind.gz var etc`, ik.SourceDir))

	// 如果是在发布机上运行，此步骤不需要执行
	if !strInArr(strings.Split(publishRole.Hosts[0], ":")[0], localIps) {
		ik.er.Get(fmt.Sprintf("%s/bind/bind.gz", ik.SourceDir), fmt.Sprintf("%s/bind", ik.SourceDir))
	}

	ik.er.SetRole(pridnsRole)
	ik.er.Put(fmt.Sprintf("%s/bind/bind.gz", ik.SourceDir), "/tmp")
	ik.er.Run("tar zxvf /tmp/bind.gz -C / && rm -rf /tmp/bind.gz && chown -R named:named /var/named/zones && chown root:named /var/named /etc/named.conf /etc/named.rfc1912.zones && systemctl enable named-chroot")
	ik.er.Local(fmt.Sprintf("rm -rf %s/bind/bind.gz", ik.SourceDir))
}

func strInArr(str string, arr []string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

func getLocalIp() ([]string, error) {
	var ips []string

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		// = GET LOCAL IP ADDRESS
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}

	return ips, nil
}
