package installk8s

import (
	"fmt"
	"strings"

	"git.yonyou.com/sysbase/backend/tool/execremote"
)

func (ik *InstallK8s) InstallEtcd() {
	etcdRole, ok := ik.resources["etcd"]
	if !ok {
		ik.Stdout <- "没有etcd资源"
		return
	}

	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	if !ik.createSslEtcd() {
		ik.Stdout <- "Etcd证书创建失败"
		return
	}

	clusterHosts := ""
	tmpStr := ""
	etcdHosts := etcdRole.Hosts
	newetcdRole, ok := ik.resources["newetcd"]
	if ik.AddNew && ok {
		etcdHosts = append(etcdHosts, newetcdRole.Hosts...)
	}
	for index, host := range etcdHosts {
		clusterHosts = fmt.Sprintf(`%s%setcd%d=https://%s:2380`, clusterHosts, tmpStr, index+1, strings.Split(host, ":")[0])
		tmpStr = ","
	}

	for index, host := range etcdHosts {
		curHost := strings.Split(host, ":")[0]
		etcdName := fmt.Sprintf("etcd%d", index+1)

		ik.er.SetRole(publishRole)
		cmds := []string{
			fmt.Sprintf(`cd %s/etcd && sed "s#ETCD_INITIAL_CLUSTER_STATE#new#g" etcd.service.tpl > usr/lib/systemd/system/etcd.service`, ik.SourceDir),
			fmt.Sprintf(`cd %s/etcd && sed -i "s#ETCD_HOST#%s#g" usr/lib/systemd/system/etcd.service`, ik.SourceDir, curHost),
			fmt.Sprintf(`cd %s/etcd && sed -i "s#ETCD_NAME#%s#g" usr/lib/systemd/system/etcd.service`, ik.SourceDir, etcdName),
			fmt.Sprintf(`cd %s/etcd && sed -i "s#ETCD_INITIAL_CLUSTER#%s#g" usr/lib/systemd/system/etcd.service`, ik.SourceDir, clusterHosts),
			fmt.Sprintf(`cd %s/etcd && tar zcvf etcd.gz etc usr`, ik.SourceDir),
		}
		ik.er.Run(cmds...)

		// 如果是在发布机上运行，此步骤不需要执行
		if !strInArr(strings.Split(publishRole.Hosts[0], ":")[0], localIps) {
			ik.er.Get(fmt.Sprintf("%s/etcd/etcd.gz", ik.SourceDir), fmt.Sprintf("%s/etcd", ik.SourceDir))
		}

		oneHostRole := execremote.Role{
			Hosts: []string{host},
		}
		ik.er.SetRole(oneHostRole)
		ik.er.Put(fmt.Sprintf("%s/etcd/etcd.gz", ik.SourceDir), "/tmp")
		cmds = []string{
			`id etcd >& /dev/null || useradd -c "etcd user" -s /sbin/nologin -d /var/lib/etcd -r etcd`,
			`tar zxvf /tmp/etcd.gz -C / && rm -rf /tmp/etcd.gz && chown -R etcd:etcd /etc/etcd/ssl && mkdir -p /data/etcd && chown -R etcd:etcd /data/etcd && chmod 750 /data/etcd`,
			`systemctl daemon-reload && systemctl enable etcd`,
		}
		ik.er.Run(cmds...)
		ik.er.Local(fmt.Sprintf("rm -rf %s/etcd/etcd.gz", ik.SourceDir))
	}
}

func (ik *InstallK8s) UpdateSslEtcd() {
	ik.Stdout <- "[开始]关闭master服务及etcd服务"
	ik.Params.DoWhat = "stop"
	ik.ServiceMaster()
	ik.ServiceEtcd()
	ik.Stdout <- "[结束]关闭master服务及etcd服务"

	ik.Stdout <- "[开始]更新etcd节点证书"
	ik.newetcdInstallEtcd()
	ik.updateMasterEtcdSsl()
	ik.Stdout <- "[结束]更新etcd节点证书"

	ik.Stdout <- "[开始]启动etcd节点服务及master服务"
	ik.Params.DoWhat = "restart"
	ik.ServiceEtcd()
	ik.ServiceMaster()
	ik.Stdout <- "[结束]启动etcd节点服务及master服务"

	ik.Stdout <- "[开始]重新配置calico"
	ik.initCalico()
	ik.Stdout <- "[结束]重新配置calico"

	ik.Stdout <- "[结束]所有节点证书更新完毕[祝您好运！]"
}

func (ik *InstallK8s) createSslEtcd() bool {
	etcdLbRole, ok := ik.resources["etcdlb"]
	if !ok {
		ik.Stdout <- "没有etcdlb资源"
		return false
	}

	etcdRole, ok := ik.resources["etcd"]
	if !ok {
		ik.Stdout <- "没有etcd资源"
		return false
	}

	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return false
	}

	hosts := ""
	linesSed := "N;"
	etcdHosts := etcdRole.Hosts

	newEtcdRole, ok := ik.resources["newetcd"]
	if ik.AddNew && ok {
		etcdHosts = append(etcdHosts, newEtcdRole.Hosts...)
	}
	for _, host := range etcdHosts {
		hosts = fmt.Sprintf(`%s,\n      \"%s\"`, hosts, strings.Split(host, ":")[0])
		linesSed = fmt.Sprintf("%sN;", linesSed)
	}

	etcdLbHost := strings.Split(etcdLbRole.Hosts[0], ":")[0]
	ik.er.SetRole(publishRole)
	cmds := []string{
		fmt.Sprintf(`cd %s/etcd && sed "s#LVS#%s#g" etcd-csr.json.tpl > etcd-csr.json`, ik.SourceDir, etcdLbHost),
		fmt.Sprintf(`cd %s/etcd && sed -i "s#HOSTS#%s#g" etcd-csr.json`, ik.SourceDir, hosts),
		fmt.Sprintf("cd %s/etcd && ./create_ssl.sh && /usr/bin/cp -rpf *.pem etc/etcd/ssl", ik.SourceDir),
		// fmt.Sprintf(`cd %s/etcd && sed -i ":label;%ss#' + hosts + '#HOSTS#;b label" etcd-csr.json`, ik.SourceDir, linesSed),
	}
	ik.er.Run(cmds...)

	return true
}

func (ik *InstallK8s) installLvsvipEtcd(etcdNode execremote.Role) {
	if !ik.UseLvs {
		return
	}

	etcdLbRole, ok := ik.resources["etcdlb"]
	if !ok {
		ik.Stdout <- "没有etcdlb资源"
		return
	}

	etcdLbHost := strings.Split(etcdLbRole.Hosts[0], ":")[0]

	ik.er.SetRole(etcdNode)
	cmds := []string{
		fmt.Sprintf(`ifconfig lo:etcd:0 %s broadcast %s netmask 255.255.255.255 up && echo -e "#/bin/sh\n# chkconfig:   2345 90 10\nifconfig lo:etcd:0 %s broadcast %s netmask 255.255.255.255 up" > /etc/rc.d/init.d/vip_route_etcd.sh`, etcdLbHost, etcdLbHost, etcdLbHost, etcdLbHost),
		fmt.Sprintf(`route add -host %s dev lo:etcd:0 && echo "route add -host %s dev lo:etcd:0" >> /etc/rc.d/init.d/vip_route_etcd.sh`, etcdLbHost, etcdLbHost),
		`chmod +x /etc/rc.d/init.d/vip_route_etcd.sh && chkconfig --add vip_route_etcd.sh && chkconfig vip_route_etcd.sh on`,

		`echo "1" > /proc/sys/net/ipv4/conf/lo/arp_ignore && echo "2" > /proc/sys/net/ipv4/conf/lo/arp_announce && echo "1" > /proc/sys/net/ipv4/conf/all/arp_ignore && echo "2" > /proc/sys/net/ipv4/conf/all/arp_announce`,
	}
	ik.er.Run(cmds...)
}
