package installk8s

import (
	"fmt"

	"git.yonyou.com/sysbase/backend/model"
	"git.yonyou.com/sysbase/backend/tool/execremote"
)

func (ik *InstallK8s) NewmasterInstall() {
	newmasterRole, ok := ik.resources["newmaster"]
	if !ok {
		ik.Stdout <- "没有newmaster资源"
		return
	}

	ik.Stdout <- "[开始]关闭master服务服务"
	ik.Params.DoWhat = "stop"
	ik.ServiceMaster()
	ik.Stdout <- "[结束]关闭master服务服务"

	ik.newmasterInstallBase()
	ik.newmasterInstallMaster()

	ik.Stdout <- "[开始]启动新master节点"
	ik.Params.DoWhat = "restart"
	ik.ServiceMaster()
	ik.ServiceNewMaster()
	ik.Stdout <- "[结束]启动新master节点"

	ik.Stdout <- "[开始]修改lvs配置"
	ik.installLvsvipMaster(newmasterRole)
	ik.installLvsNew()
	ik.Stdout <- "[结束]修改lvs配置"

	ik.Stdout <- "[开始]重启node服务"
	ik.Params.DoWhat = "restart"
	ik.Params.DoDocker = false
	if nodeRole, ok := ik.resources["node"]; ok {
		ik.serviceNode(nodeRole)
	}
	if newnodeRole, ok := ik.resources["newnode"]; ok {
		ik.serviceNode(newnodeRole)
	}
	ik.Stdout <- "[结束]重启node服务"

	ik.Stdout <- "[开始]更新核心secret和pod，其他secret需要根据实际情况手动更新"
	ik.wait(8, "等待前面服务启动完成...")
	ik.updateSecretAndPod()
	ik.Stdout <- "[结束]更新核心secret和pod，其他secret需要根据实际情况手动更新"

	ik.Stdout <- "[开始]修改scope，由newmaster改为master"
	ik.newmasterUpdateScope()
	ik.Stdout <- "[结束]修改scope，由newmaster改为master"

	ik.Stdout <- "[结束]所有节点添加完毕[祝您好运！]"
}

func (ik *InstallK8s) newmasterUpdateScope() {
	k8sClusterResource := model.K8sClusterResource{}
	resource, err := k8sClusterResource.ListResource(ik.Params.K8sClusterID, []string{"newmaster"})
	if err != nil {
		ik.Stdout <- err.Error()
		return
	}

	for _, r := range resource {
		oneHostRole := execremote.Role{
			Hosts: []string{
				fmt.Sprintf(`%s:%d`, r.Host, r.Port),
			},
			WaitOutput: true,
		}
		ik.er.SetRole(oneHostRole)

		ik.er.Run(`ps aux | grep kube-apiserver | grep -v grep > /dev/null ; echo $?`)
		if ik.er.GetCmdReturn()[0] == "0" {
			model.Resource{
				ID: r.ID,
			}.Edit(model.Resource{
				Scope: "master",
			})
		}
	}
}

func (ik *InstallK8s) newmasterInstallBase() {
	newmasterRole, ok := ik.resources["newmaster"]
	if !ok {
		ik.Stdout <- "没有newmaster资源"
		return
	}

	ik.Stdout <- "[开始]安装基础环境"
	ik.er.SetRole(newmasterRole)
	ik.installBase()
	ik.Stdout <- "[结束]安装基础环境"
}

func (ik *InstallK8s) newmasterInstallMaster() {
	ik.Stdout <- "[开始]安装master节点"
	ik.AddNew = true
	ik.InstallMaster()
	ik.modifyNodeConf()
	ik.Stdout <- "[结束]安装master节点"
}
