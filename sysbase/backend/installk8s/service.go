package installk8s

import (
	"fmt"

	"git.yonyou.com/sysbase/backend/tool/execremote"
)

func (ik *InstallK8s) ServicePublish() {
	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	if !ik.checkDoWhat() {
		ik.Stdout <- "执行命令非法"
		return
	}

	ik.er.SetRole(publishRole)

	if ik.Params.DoWhat == "start" || ik.Params.DoWhat == "restart" {
		ik.er.Run("iptables -P FORWARD ACCEPT ; systemctl daemon-reload")
	}

	ik.er.Run(fmt.Sprintf("systemctl %s docker", ik.Params.DoWhat))
}

func (ik *InstallK8s) ServiceEtcd() {
	etcdRole, ok := ik.resources["etcd"]
	if !ok {
		ik.Stdout <- "没有etcd资源"
		return
	}

	ik.serviceEtcd(etcdRole)
}

func (ik *InstallK8s) serviceEtcd(etcdRole execremote.Role) {
	if !ik.checkDoWhat() {
		ik.Stdout <- "执行命令非法"
		return
	}

	ik.er.SetRole(etcdRole)

	if ik.Params.DoWhat == "start" || ik.Params.DoWhat == "restart" {
		ik.er.Run("iptables -P FORWARD ACCEPT ; systemctl daemon-reload")
	}

	ik.er.Run(fmt.Sprintf("systemctl %s etcd", ik.Params.DoWhat))
}

func (ik *InstallK8s) ServiceMaster() {
	masterRole, ok := ik.resources["master"]
	if !ok {
		ik.Stdout <- "没有master资源"
		return
	}

	ik.serviceMaster(masterRole)
}

func (ik *InstallK8s) ServiceNewMaster() {
	newmasterRole, ok := ik.resources["newmaster"]
	if !ok {
		ik.Stdout <- "没有newmaster资源"
		return
	}

	ik.serviceMaster(newmasterRole)
}

func (ik *InstallK8s) serviceMaster(masterRole execremote.Role) {
	if !ik.checkDoWhat() {
		ik.Stdout <- "执行命令非法"
		return
	}

	ik.er.SetRole(masterRole)

	if ik.Params.DoWhat == "start" || ik.Params.DoWhat == "restart" {
		ik.er.Run("iptables -P FORWARD ACCEPT ; systemctl daemon-reload")
	}

	cmds := []string{
		fmt.Sprintf("systemctl %s kube-apiserver", ik.Params.DoWhat),
		fmt.Sprintf("systemctl %s kube-controller-manager", ik.Params.DoWhat),
		fmt.Sprintf("systemctl %s kube-scheduler", ik.Params.DoWhat),
	}
	ik.er.Run(cmds...)

	if ik.Params.DoWhat == "stop" {
		cmds := []string{
			`ps aux | grep kube-apiserver | grep -v grep | awk '{if($2 != ""){system("kill -9 "$2)}}'`,
			`ps aux | grep kube-controller-manager | grep -v grep | awk '{if($2 != ""){system("kill -9 "$2)}}'`,
			`ps aux | grep kube-scheduler | grep -v grep | awk '{if($2 != ""){system("kill -9 "$2)}}'`,
		}
		ik.er.Run(cmds...)
	}
}

func (ik *InstallK8s) ServiceNode() {
	nodeRole, ok := ik.resources["node"]
	if !ok {
		ik.Stdout <- "没有node资源"
		return
	}

	ik.serviceNode(nodeRole)
}

func (ik *InstallK8s) serviceNode(nodeRole execremote.Role) {
	if !ik.checkDoWhat() {
		ik.Stdout <- "执行命令非法"
		return
	}

	ik.er.SetRole(nodeRole)

	if ik.Params.DoWhat == "start" || ik.Params.DoWhat == "restart" {
		ik.er.Run("iptables -P FORWARD ACCEPT ; systemctl daemon-reload")
	}

	cmds := []string{
		fmt.Sprintf("systemctl %s kubelet", ik.Params.DoWhat),
		fmt.Sprintf("systemctl %s kube-proxy", ik.Params.DoWhat),
	}
	ik.er.Run(cmds...)

	if ik.Params.DoDocker {
		ik.er.Run(fmt.Sprintf("systemctl %s docker", ik.Params.DoWhat))
	}

	if ik.Params.DoWhat == "stop" {
		cmds := []string{
			`ps aux | grep kubelet | grep -v grep | awk '{if($2 != ""){system("kill -9 "$2)}}'`,
			`ps aux | grep kube-proxy | grep -v grep | awk '{if($2 != ""){system("kill -9 "$2)}}'`,
		}
		ik.er.Run(cmds...)
		if ik.Params.DoDocker {
			ik.er.Run(`ps aux | grep docker | grep -v grep | awk '{if($2 != ""){system("kill -9 "$2)}}'`)
		}
	}
}

func (ik *InstallK8s) ServiceDns() {
	pridnsRole, ok := ik.resources["pridns"]
	if !ok {
		ik.Stdout <- "没有pridns资源"
		return
	}

	if !ik.checkDoWhat() {
		ik.Stdout <- "执行命令非法"
		return
	}

	ik.er.SetRole(pridnsRole)

	if ik.Params.DoWhat == "start" || ik.Params.DoWhat == "restart" {
		ik.er.Run("iptables -P FORWARD ACCEPT ; systemctl daemon-reload")
	}

	ik.er.Run(fmt.Sprintf("systemctl %s named-chroot", ik.Params.DoWhat))
}

func (ik *InstallK8s) checkDoWhat() bool {
	for _, dw := range []string{"start", "restart", "stop"} {
		if dw == ik.Params.DoWhat {
			return true
		}
	}

	return false
}
