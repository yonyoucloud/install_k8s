package installk8s

import (
	"fmt"
	"strings"

	"sysbase/tool/execremote"
)

func (ik *InstallK8s) InstallNode() {
	nodeRole, ok := ik.resources["node"]
	if !ok {
		ik.Stdout <- "没有node资源"
		return
	}

	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	masterLbRole, ok := ik.resources["masterlb"]
	if !ok {
		ik.Stdout <- "没有masterlb资源"
		return
	}

	if !ik.createSslNode() {
		ik.Stdout <- "Node证书创建失败"
		return
	}

	masterLbHost := strings.Split(masterLbRole.Hosts[0], ":")[0]

	ik.installNode(masterLbHost, publishRole, nodeRole)
}

func (ik *InstallK8s) UpdateSslNode() {
	ik.Stdout <- "[开始]更新node节点证书"
	ik.AddNew = true
	ik.modifyNodeConf()
	ik.Stdout <- "[结束]更新node节点证书"

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

	ik.Stdout <- "[结束]所有节点证书更新完毕[祝您好运！]"
}

func (ik *InstallK8s) installNode(masterLbHost string, publishRole, nodeRole execremote.Role) {
	pridnsRole, ok := ik.resources["pridns"]
	if !ok {
		ik.Stdout <- "没有pridns资源"
		return
	}
	pridnsHost := strings.Split(pridnsRole.Hosts[0], ":")[0]

	for _, host := range nodeRole.Hosts {
		curHost := strings.Split(host, ":")[0]

		ik.er.SetRole(publishRole)
		cmds := []string{
			fmt.Sprintf(`cd %s/node && sed "s#NODE_HOST#%s#g" kubelet-csr.json.tpl > kubelet-csr.json`, ik.SourceDir, curHost),
			fmt.Sprintf(`cd %s/node && sed "s#NODE_HOST#%s#g" kubelet.tpl > etc/kubernetes/kubelet`, ik.SourceDir, curHost),
			fmt.Sprintf(`cd %s/node && sed "s#NODE_HOST#%s#g" kubelet.yaml.tpl > etc/kubernetes/kubelet.yaml`, ik.SourceDir, curHost),
			fmt.Sprintf(`cd %s/node && sed "s#K8S_MASTER_LVS#%s#g" config.tpl > etc/kubernetes/config`, ik.SourceDir, masterLbHost),
			fmt.Sprintf(`cd %s/node && sed -i "s#K8S_MASTER_LVS#%s#g" etc/kubernetes/kubelet`, ik.SourceDir, masterLbHost),
			fmt.Sprintf(`cd %s/node && cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=frognew kubelet-csr.json | cfssljson -bare kubelet`, ik.SourceDir),

			// kubelet.conf
			fmt.Sprintf(`cd %s/node && rm -rf etc/kubernetes/kubelet.conf && kubectl config set-cluster kubernetes --certificate-authority=ca.pem --embed-certs=true --server=https://%s:6443 --kubeconfig=etc/kubernetes/kubelet.conf`, ik.SourceDir, masterLbHost),
			fmt.Sprintf(`cd %s/node && kubectl config set-credentials system:node:%s --client-certificate=kubelet.pem --embed-certs=true --client-key=kubelet-key.pem --kubeconfig=etc/kubernetes/kubelet.conf`, ik.SourceDir, curHost),
			fmt.Sprintf(`cd %s/node && kubectl config set-context system:node:%s@kubernetes --cluster=kubernetes --user=system:node:%s --kubeconfig=etc/kubernetes/kubelet.conf`, ik.SourceDir, curHost, curHost),
			fmt.Sprintf(`cd %s/node && kubectl config use-context system:node:%s@kubernetes --kubeconfig=etc/kubernetes/kubelet.conf`, ik.SourceDir, curHost),

			fmt.Sprintf(`cd %s/node && /usr/bin/cp -rpf *.pem etc/kubernetes/pki`, ik.SourceDir),
		}
		ik.er.Run(cmds...)

		cmd := fmt.Sprintf(`cd %s/node && tar zcvf node.gz etc usr`, ik.SourceDir)
		if ik.OnlyConf {
			cmd = fmt.Sprintf(`cd %s/node && tar zcvf node.gz etc`, ik.SourceDir)
		}
		ik.er.Run(cmd)

		// 如果是在发布机上运行，此步骤不需要执行
		if !strInArr(strings.Split(publishRole.Hosts[0], ":")[0], localIps) {
			ik.er.Get(fmt.Sprintf("%s/node/node.gz", ik.SourceDir), fmt.Sprintf("%s/node", ik.SourceDir))
		}

		oneHostRole := execremote.Role{
			Hosts: []string{host},
		}
		ik.er.SetRole(oneHostRole)
		ik.er.Put(fmt.Sprintf("%s/node/node.gz", ik.SourceDir), "/tmp")
		cmds = []string{
			`tar zxvf /tmp/node.gz -C / && rm -rf /tmp/node.gz`,
			`systemctl daemon-reload && systemctl enable kube-proxy && systemctl enable kubelet && mkdir -p /data/kubelet && chmod 750 /data/kubelet`,
		}
		cmds = getModifyDnsCmds(cmds, pridnsHost)
		ik.er.Run(cmds...)
		ik.er.Local(fmt.Sprintf("rm -rf %s/node/node.gz", ik.SourceDir))
	}
}

func (ik *InstallK8s) modifyNodeConf() {
	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	masterLbRole, ok := ik.resources["masterlb"]
	if !ok {
		ik.Stdout <- "没有masterlb资源"
		return
	}

	if !ik.createSslNode() {
		ik.Stdout <- "Node证书创建失败"
		return
	}

	masterLbHost := strings.Split(masterLbRole.Hosts[0], ":")[0]

	ik.OnlyConf = true
	if nodeRole, ok := ik.resources["node"]; ok {
		ik.installNode(masterLbHost, publishRole, nodeRole)
		ik.kubeletcniNode(nodeRole)
	}
	if newnodeRole, ok := ik.resources["newnode"]; ok {
		ik.installNode(masterLbHost, publishRole, newnodeRole)
		ik.kubeletcniNode(newnodeRole)
	}
}

func (ik *InstallK8s) createSslNode() bool {
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

	masterLbHost := strings.Split(masterLbRole.Hosts[0], ":")[0]

	ik.er.SetRole(publishRole)
	cmds := []string{
		fmt.Sprintf(`cd %s/node && /usr/bin/cp -rpf ../master/{ca.pem,ca-key.pem,ca-config.json} .`, ik.SourceDir),
		fmt.Sprintf(`cd %s/node && cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=frognew kube-proxy-csr.json | cfssljson -bare kube-proxy`, ik.SourceDir),

		// kube-proxy.conf
		fmt.Sprintf(`cd %s/node && rm -rf etc/kubernetes/kube-proxy.conf && kubectl config set-cluster kubernetes --certificate-authority=ca.pem --embed-certs=true --server=https://%s:6443 --kubeconfig=etc/kubernetes/kube-proxy.conf`, ik.SourceDir, masterLbHost),
		fmt.Sprintf(`cd %s/node && kubectl config set-credentials system:kube-proxy --client-certificate=kube-proxy.pem --embed-certs=true --client-key=kube-proxy-key.pem --kubeconfig=etc/kubernetes/kube-proxy.conf`, ik.SourceDir),
		fmt.Sprintf(`cd %s/node && kubectl config set-context system:kube-proxy@kubernetes --cluster=kubernetes --user=system:kube-proxy --kubeconfig=etc/kubernetes/kube-proxy.conf`, ik.SourceDir),
		fmt.Sprintf(`cd %s/node && kubectl config use-context system:kube-proxy@kubernetes --kubeconfig=etc/kubernetes/kube-proxy.conf`, ik.SourceDir),
	}
	ik.er.Run(cmds...)

	return true
}

func (ik *InstallK8s) kubeletcniNode(nodeRole execremote.Role) {
	ik.er.SetRole(nodeRole)
	ik.er.Run(`sed -i 's#--config=/etc/kubernetes/kubelet.yaml"#--config=/etc/kubernetes/kubelet.yaml --network-plugin=cni --cni-conf-dir=/etc/cni/net.d --cni-bin-dir=/opt/cni/bin"#g' /etc/kubernetes/kubelet`)

	if !ik.OnlyConf {
		ik.er.Run(`systemctl restart kubelet`)
	}
}
