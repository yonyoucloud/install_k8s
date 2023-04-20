package installk8s

import (
	"fmt"
	"strings"

	"sysbase/tool/execremote"
)

func (ik *InstallK8s) InstallMaster() {
	masterRole, ok := ik.resources["master"]
	if !ok {
		ik.Stdout <- "没有master资源"
		return
	}

	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	etcdLbRole, ok := ik.resources["etcdlb"]
	if !ok {
		ik.Stdout <- "没有etcdlb资源"
		return
	}

	if !ik.createSslMaster() {
		ik.Stdout <- "Master证书创建失败"
		return
	}

	etcdLbHost := strings.Split(etcdLbRole.Hosts[0], ":")[0]

	masterHosts := masterRole.Hosts
	newMasterRole, ok := ik.resources["newmaster"]
	if ik.AddNew && ok {
		masterHosts = append(masterHosts, newMasterRole.Hosts...)
	}
	for _, host := range masterHosts {
		curHost := strings.Split(host, ":")[0]

		ik.er.SetRole(publishRole)
		cmds := []string{
			fmt.Sprintf(`cd %s/master && sed "s#K8S_HOST#%s#g" config.tpl > etc/kubernetes/config`, ik.SourceDir, curHost),
			fmt.Sprintf(`cd %s/master && sed "s#ETCD_LVS_HOST#%s#g" apiserver.tpl > etc/kubernetes/apiserver`, ik.SourceDir, etcdLbHost),
			fmt.Sprintf(`cd %s/master && mkdir -p etc/kubernetes/pki/etcd && chmod 750 etc/kubernetes/pki/etcd`, ik.SourceDir),
			fmt.Sprintf(`/usr/bin/cp -rpf %s/etcd/etc/etcd/ssl/{ca.pem,etcd.pem,etcd-key.pem} %s/master/etc/kubernetes/pki/etcd`, ik.SourceDir, ik.SourceDir),
		}
		ik.er.Run(cmds...)

		isIn := false
		for _, h := range masterRole.Hosts {
			if h == host {
				isIn = true
				break
			}
		}

		cmd := fmt.Sprintf(`cd %s/master && tar zcvf master.gz etc usr`, ik.SourceDir)
		if ik.AddNew && isIn {
			cmd = fmt.Sprintf(`cd %s/master && tar zcvf master.gz etc`, ik.SourceDir)
		}
		ik.er.Run(cmd)

		// 如果是在发布机上运行，此步骤不需要执行
		if !strInArr(strings.Split(publishRole.Hosts[0], ":")[0], localIps) {
			ik.er.Get(fmt.Sprintf("%s/master/master.gz", ik.SourceDir), fmt.Sprintf("%s/master", ik.SourceDir))
		}

		oneHostRole := execremote.Role{
			Hosts: []string{host},
		}
		ik.er.SetRole(oneHostRole)
		ik.er.Put(fmt.Sprintf("%s/master/master.gz", ik.SourceDir), "/tmp")
		cmds = []string{
			`tar zxvf /tmp/master.gz -C / && rm -rf /tmp/master.gz`,
			`systemctl daemon-reload && systemctl enable kube-apiserver && systemctl enable kube-controller-manager && systemctl enable kube-scheduler`,
		}
		ik.er.Run(cmds...)
		ik.er.Local(fmt.Sprintf("rm -rf %s/master/master.gz", ik.SourceDir))
	}
}

func (ik *InstallK8s) UpdateSslMaster() {
	ik.Stdout <- "[开始]关闭master服务服务"
	ik.Params.DoWhat = "stop"
	ik.ServiceMaster()
	ik.Stdout <- "[结束]关闭master服务服务"

	ik.Stdout <- "[开始]更新master、node节点证书"
	ik.AddNew = true
	ik.InstallMaster()
	ik.modifyNodeConf()
	ik.Stdout <- "[结束]更新master、node节点证书"

	ik.Stdout <- "[开始]启动master节点"
	ik.Params.DoWhat = "restart"
	ik.ServiceMaster()
	ik.Stdout <- "[结束]启动master节点"

	ik.Stdout <- "[开始]重启node服务"
	ik.Params.DoWhat = "restart"
	ik.Params.DoContainerd = false
	if nodeRole, ok := ik.resources["node"]; ok {
		ik.serviceNode(nodeRole)
	}
	if newnodeRole, ok := ik.resources["newnode"]; ok {
		ik.serviceNode(newnodeRole)
	}
	ik.Stdout <- "[结束]重启node服务"

	ik.Stdout <- "[开始]更新核心secret和pod，其他secret需要根据实际情况手动更新"
	ik.wait(30, "等待前面服务启动完成...")
	ik.updateSecretAndPod()
	ik.Stdout <- "[结束]更新核心secret和pod，其他secret需要根据实际情况手动更新"

	ik.Stdout <- "[结束]所有节点证书更新完毕[祝您好运！]"
}

func (ik *InstallK8s) createSslMaster() bool {
	masterLbRole, ok := ik.resources["masterlb"]
	if !ok {
		ik.Stdout <- "没有masterlb资源"
		return false
	}

	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return false
	}

	masterRole, ok := ik.resources["master"]
	if !ok {
		ik.Stdout <- "没有master资源"
		return false
	}

	hosts := fmt.Sprintf(`,\n      \"%s\"`, strings.Split(publishRole.Hosts[0], ":")[0])
	masterHosts := masterRole.Hosts
	newMasterRole, ok := ik.resources["newmaster"]
	if ik.AddNew && ok {
		masterHosts = append(masterHosts, newMasterRole.Hosts...)
	}
	for _, host := range masterHosts {
		hosts = fmt.Sprintf(`%s,\n      \"%s\"`, hosts, strings.Split(host, ":")[0])
	}
	masterLbHost := strings.Split(masterLbRole.Hosts[0], ":")[0]

	ik.er.SetRole(publishRole)
	cmds := []string{
		fmt.Sprintf(`cd %s/master && sed "s#LVS#%s#g" admin-csr.json.tpl > admin-csr.json`, ik.SourceDir, masterLbHost),
		fmt.Sprintf(`cd %s/master && sed "s#LVS#%s#g" apiserver-csr.json.tpl > apiserver-csr.json`, ik.SourceDir, masterLbHost),
		fmt.Sprintf(`cd %s/master && sed "s#LVS#%s#g" ca-config.json.tpl > ca-config.json`, ik.SourceDir, masterLbHost),
		fmt.Sprintf(`cd %s/master && sed "s#LVS#%s#g" controller-manager-csr.json.tpl > controller-manager-csr.json`, ik.SourceDir, masterLbHost),
		fmt.Sprintf(`cd %s/master && sed "s#LVS#%s#g" scheduler-csr.json.tpl > scheduler-csr.json`, ik.SourceDir, masterLbHost),
		fmt.Sprintf(`cd %s/master && sed -i "s#HOSTS#%s#g" admin-csr.json apiserver-csr.json ca-config.json controller-manager-csr.json scheduler-csr.json`, ik.SourceDir, hosts),
		fmt.Sprintf(`cd %s/master && ./create_ssl.sh && /usr/bin/cp -rpf *.pem etc/kubernetes/pki`, ik.SourceDir),

		// admin.conf
		fmt.Sprintf(`cd %s/master && rm -rf etc/kubernetes/admin.conf && kubectl config set-cluster kubernetes --certificate-authority=ca.pem --embed-certs=true --server=https://%s:6443 --kubeconfig=etc/kubernetes/admin.conf`, ik.SourceDir, masterLbHost),
		fmt.Sprintf(`cd %s/master && kubectl config set-credentials kubernetes-admin --client-certificate=admin.pem --embed-certs=true --client-key=admin-key.pem --kubeconfig=etc/kubernetes/admin.conf`, ik.SourceDir),
		fmt.Sprintf(`cd %s/master && kubectl config set-context kubernetes-admin@kubernetes --cluster=kubernetes --user=kubernetes-admin --kubeconfig=etc/kubernetes/admin.conf`, ik.SourceDir),
		fmt.Sprintf(`cd %s/master && kubectl config use-context kubernetes-admin@kubernetes --kubeconfig=etc/kubernetes/admin.conf`, ik.SourceDir),
		fmt.Sprintf(`cd %s/master && mkdir -p /root/.kube && chmod 750 /root/.kube && /usr/bin/cp -rpf etc/kubernetes/admin.conf /root/.kube/config`, ik.SourceDir),

		// scheduler.conf
		fmt.Sprintf(`cd %s/master && rm -rf etc/kubernetes/scheduler.conf && kubectl config set-cluster kubernetes --certificate-authority=ca.pem --embed-certs=true --server=https://%s:6443 --kubeconfig=etc/kubernetes/scheduler.conf`, ik.SourceDir, masterLbHost),
		fmt.Sprintf(`cd %s/master && kubectl config set-credentials system:kube-scheduler --client-certificate=scheduler.pem --embed-certs=true --client-key=scheduler-key.pem --kubeconfig=etc/kubernetes/scheduler.conf`, ik.SourceDir),
		fmt.Sprintf(`cd %s/master && kubectl config set-context system:kube-scheduler@kubernetes --cluster=kubernetes --user=system:kube-scheduler --kubeconfig=etc/kubernetes/scheduler.conf`, ik.SourceDir),
		fmt.Sprintf(`cd %s/master && kubectl config use-context system:kube-scheduler@kubernetes --kubeconfig=etc/kubernetes/scheduler.conf`, ik.SourceDir),

		// controller-manager.conf
		fmt.Sprintf(`cd %s/master && rm -rf etc/kubernetes/controller-manager.conf && kubectl config set-cluster kubernetes --certificate-authority=ca.pem --embed-certs=true --server=https://%s:6443 --kubeconfig=etc/kubernetes/controller-manager.conf`, ik.SourceDir, masterLbHost),
		fmt.Sprintf(`cd %s/master && kubectl config set-credentials system:kube-controller-manager --client-certificate=controller-manager.pem --embed-certs=true --client-key=controller-manager-key.pem --kubeconfig=etc/kubernetes/controller-manager.conf`, ik.SourceDir),
		fmt.Sprintf(`cd %s/master && kubectl config set-context system:kube-controller-manager@kubernetes --cluster=kubernetes --user=system:kube-controller-manager --kubeconfig=etc/kubernetes/controller-manager.conf`, ik.SourceDir),
		fmt.Sprintf(`cd %s/master && kubectl config use-context system:kube-controller-manager@kubernetes --kubeconfig=etc/kubernetes/controller-manager.conf`, ik.SourceDir),
	}
	ik.er.Run(cmds...)

	return true
}

func (ik *InstallK8s) updateMasterEtcdSsl() bool {
	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return false
	}

	masterRole, ok := ik.resources["master"]
	if !ok {
		ik.Stdout <- "没有master资源"
		return false
	}

	ik.er.SetRole(publishRole)
	cmds := []string{
		fmt.Sprintf(`/usr/bin/cp -rpf %s/etcd/etc/etcd/ssl/{ca.pem,etcd.pem,etcd-key.pem} %s/master/etc/kubernetes/pki/etcd`, ik.SourceDir, ik.SourceDir),
		fmt.Sprintf(`cd %s/master/etc/kubernetes/pki && tar zcvf etcd.gz etcd`, ik.SourceDir),
	}
	ik.er.Run(cmds...)

	// 如果是在发布机上运行，此步骤不需要执行
	if !strInArr(strings.Split(publishRole.Hosts[0], ":")[0], localIps) {
		ik.er.Get(fmt.Sprintf("%s/master/etc/kubernetes/pki/etcd.gz", ik.SourceDir), fmt.Sprintf("%s/master/etc/kubernetes/pki", ik.SourceDir))
	}

	ik.er.SetRole(masterRole)
	ik.er.Put(fmt.Sprintf("%s/master/etc/kubernetes/pki/etcd.gz", ik.SourceDir), "/etc/kubernetes/pki")
	ik.er.Run(`cd /etc/kubernetes/pki && tar zxvf etcd.gz && rm -rf etcd.gz`)
	ik.er.Local(fmt.Sprintf("rm -rf %s/master/etc/kubernetes/pki/etcd.gz", ik.SourceDir))

	return true
}

func (ik *InstallK8s) installLvsvipMaster(masterNode execremote.Role) {
	if !ik.UseLvs {
		return
	}

	masterLbRole, ok := ik.resources["masterlb"]
	if !ok {
		ik.Stdout <- "没有masterlb资源"
		return
	}

	masterLbHost := strings.Split(masterLbRole.Hosts[0], ":")[0]

	ik.er.SetRole(masterNode)
	cmds := []string{
		fmt.Sprintf(`ifconfig lo:master:0 %s broadcast %s netmask 255.255.255.255 up && echo -e "#/bin/sh\n# chkconfig:   2345 90 10\nifconfig lo:master:0 %s broadcast %s netmask 255.255.255.255 up" > /etc/rc.d/init.d/vip_route_master.sh`, masterLbHost, masterLbHost, masterLbHost, masterLbHost),
		fmt.Sprintf(`route add -host %s dev lo:master:0 && echo "route add -host %s dev lo:master:0" >> /etc/rc.d/init.d/vip_route_master.sh`, masterLbHost, masterLbHost),
		`chmod +x /etc/rc.d/init.d/vip_route_master.sh && chkconfig --add vip_route_master.sh && chkconfig vip_route_master.sh on`,

		`echo "1" > /proc/sys/net/ipv4/conf/lo/arp_ignore && echo "2" > /proc/sys/net/ipv4/conf/lo/arp_announce && echo "1" > /proc/sys/net/ipv4/conf/all/arp_ignore && echo "2" > /proc/sys/net/ipv4/conf/all/arp_announce`,
	}
	ik.er.Run(cmds...)
}

func (ik *InstallK8s) updateSecretAndPod() {
	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	ik.er.SetRole(publishRole)
	cmds := []string{
		`updatesecret`, // 更新Secret的Token字段，不需要滚动升级依赖服务了
		// rollUpdateCmd("deployment", "kube-system", "calico-kube-controllers"),
		// rollUpdateCmd("daemonset", "kube-system", "calico-node"),
		// rollUpdateCmd("deployment", "kube-system", "coredns"),
		// rollUpdateCmd("deployment", "kube-system", "metrics-server"),
		// rollUpdateCmd("deployment", "kube-system", "kube-state-metrics"),
		// rollUpdateCmd("statefulset", "kube-system", "prometheus"),
		// rollUpdateCmd("deployment", "kube-system", "kubernetes-dashboard"),
		// rollUpdateCmd("deployment", "istio-system", "istiod"),
		// rollUpdateCmd("deployment", "istio-system", "istio-ingressgateway"),
		// rollUpdateCmd("deployment", "istio-system", "kiali"),
	}
	ik.er.Run(cmds...)
}

func rollUpdateCmd(resource, namespace, name string) string {
	// kubectl -n %s get secret | grep "kubernetes.io/service-account-token" | grep %s-token | awk '{print "kubectl -n %s delete secret "$1}' | sh ;
	return fmt.Sprintf(`num=$(kubectl -n %s get %s %s -o jsonpath='{.spec.template.metadata.annotations.ca-rotation}') ; let num=$num+1 ; kubectl -n %s patch %s %s -p "{\"spec\":{\"template\":{\"metadata\":{\"annotations\":{\"ca-rotation\": \"$num\"}}}}}"`,
		// namespace, name, namespace,
		namespace, resource, name, namespace, resource, name)
}
