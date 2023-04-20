package installk8s

import "fmt"

func getModifyDnsCmds(cmds []string, pridnsHost string) []string {
	// 匹配添加私有DNS配置，确保在匹配行的上面只添加一行，并且不重复添加
	cmd1 := fmt.Sprintf(`grep "%s" /etc/resolv.conf > /dev/null 2>&1 || (cp -rp /etc/resolv.conf /tmp/resolv.conf && awk 'BEGIN{d=0};{ if($0 ~ /nameserver / && d==0) {d=1; printf("nameserver %s\n%%s\n", $0)} else {print $0}}' /tmp/resolv.conf > /etc/resolv.conf && rm -rf /tmp/resolv.conf)`, pridnsHost, pridnsHost)

	// 修改网卡配置文件，添加DNS，避免重启后DNS丢失
	cmd2 := fmt.Sprintf(`file=$(grep -rl ^DNS $(ls /etc/sysconfig/network-scripts/ifcfg-*)) && (grep "DNS1=%s" $file > /dev/null 2>&1 || (DNS=$(echo DNS1=%s && cat $file | grep ^DNS | awk -F '=' '{print "DNS"(NR+1)"="$2}' && sed -i '/^DNS.*/d' $file); echo "$DNS" >> $file))`, pridnsHost, pridnsHost)

	return append(cmds, cmd1, cmd2)
}
