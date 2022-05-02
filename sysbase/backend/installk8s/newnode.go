package installk8s

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.yonyou.com/sysbase/backend/model"
	"git.yonyou.com/sysbase/backend/tool/execremote"
)

func (ik *InstallK8s) NewnodeInstall() {
	if !ik.newnodeUpdateKernel() {
		ik.Stdout <- "正在更新系统内核"
		return
	}
	ik.newnodeInstallBase()
	ik.newnodeInstallDocker()
	ik.newnodeInstallDockerCrt()
	ik.newnodeInstallNode()
	ik.newnodeServiceNodeStart()
	ik.newnodeKubeletcniNode()
	ik.newnodeUpdateScope()
	ik.Stdout <- "[结束]所有节点添加完毕[祝您好运！]"
}

func (ik *InstallK8s) newnodeUpdateScope() {
	newnodeRole, ok := ik.resources["newnode"]
	if !ok {
		ik.Stdout <- "没有newnode资源"
		return
	}

	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	publishRole.WaitOutput = true
	ik.er.SetRole(publishRole)
	ik.er.Run(`kubectl get node | grep -v NAME | awk '{print $1}'`)
	hosts := ik.er.GetCmdReturn()

	substr := ":"
	num := 0
	var noAdded []string
	for _, h := range newnodeRole.Hosts {
		if !strings.Contains(h, substr) {
			continue
		}

		contain := false
		for _, hh := range hosts {
			if strings.Split(h, substr)[0] == hh {
				num++
				contain = true
				break
			}
		}

		if !contain {
			noAdded = append(noAdded, h)
		}
	}

	if num == len(newnodeRole.Hosts) {
		k8sClusterResource := model.K8sClusterResource{}

		resource, err := k8sClusterResource.ListResource(ik.Params.K8sClusterID, []string{"newnode"})
		if err != nil {
			ik.Stdout <- err.Error()
			return
		}

		for _, r := range resource {
			model.Resource{
				ID: r.ID,
			}.Edit(model.Resource{
				Scope: "node",
			})
		}
		ik.Stdout <- "所有newnode节点，都被标记为node节点，新节点添加完成"
		return
	}

	ik.Stdout <- fmt.Sprintf("以下这些节点没有添加成功：%v", noAdded)
}

func (ik *InstallK8s) newnodeUpdateKernel() bool {
	newnodeRole, ok := ik.resources["newnode"]
	if !ok {
		ik.Stdout <- "没有newnode资源"
		return false
	}

	hosts := ik.checkHostsKernel(newnodeRole.Hosts)
	if len(hosts) == 0 {
		return true
	}

	role := execremote.Role{
		Hosts: hosts,
	}

	ik.er.SetRole(role)
	ik.updateKernel()

	return false
}

func (ik *InstallK8s) newnodeInstallBase() {
	newnodeRole, ok := ik.resources["newnode"]
	if !ok {
		ik.Stdout <- "没有newnode资源"
		return
	}

	ik.Stdout <- "[开始]安装基础环境"
	ik.er.SetRole(newnodeRole)
	ik.installBase()
	ik.Stdout <- "[结束]安装基础环境"
}

func (ik *InstallK8s) newnodeInstallDocker() {
	newnodeRole, ok := ik.resources["newnode"]
	if !ok {
		ik.Stdout <- "没有newnode资源"
		return
	}

	ik.Stdout <- "[开始]安装docker"
	ik.er.SetRole(newnodeRole)
	ik.installDocker()
	ik.Stdout <- "[结束]安装docker"
}

func (ik *InstallK8s) newnodeInstallDockerCrt() {
	newnodeRole, ok := ik.resources["newnode"]
	if !ok {
		ik.Stdout <- "没有newnode资源"
		return
	}

	pridockerRole, ok := ik.resources["pridocker"]
	if !ok {
		ik.Stdout <- "没有pridocker资源"
		return
	}

	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	ik.Stdout <- "[开始]安装docker证书"
	priDockerHost := strings.Split(pridockerRole.Hosts[0], ":")[0]

	newnodeRole.Parallel = false
	ik.installDockerCrt(priDockerHost, publishRole, newnodeRole)
	ik.Stdout <- "[结束]安装docker证书"
}

func (ik *InstallK8s) newnodeInstallNode() {
	newnodeRole, ok := ik.resources["newnode"]
	if !ok {
		ik.Stdout <- "没有newnode资源"
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

	pridockerRole, ok := ik.resources["pridocker"]
	if !ok {
		ik.Stdout <- "没有pridocker资源"
		return
	}

	masterLbHost := strings.Split(masterLbRole.Hosts[0], ":")[0]
	priDockerHost := strings.Split(pridockerRole.Hosts[0], ":")[0]

	ik.Stdout <- "[开始]安装node节点"
	ik.installNode(masterLbHost, priDockerHost, publishRole, newnodeRole)
	ik.Stdout <- "[结束]安装node节点"
}

func (ik *InstallK8s) newnodeServiceNodeStart() {
	newnodeRole, ok := ik.resources["newnode"]
	if !ok {
		ik.Stdout <- "没有newnode资源"
		return
	}

	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	ik.Stdout <- "[开始]启动新node节点"
	ik.Params.DoWhat = "start"
	ik.Params.DoDocker = true
	ik.serviceNode(newnodeRole)

	publishRole.WaitOutput = true
	ik.er.SetRole(publishRole)
	total := len(newnodeRole.Hosts)
	var hosts []string
	for _, host := range newnodeRole.Hosts {
		hosts = append(hosts, strings.Split(host, ":")[0])
	}
	i := 0
	for {
		i++
		ik.er.Run(fmt.Sprintf(`kubectl get nodes | grep -wE "%s" | grep Ready | wc -l`, strings.Join(hosts, "|")))
		num, _ := strconv.Atoi(ik.er.GetCmdReturn()[0])
		ik.Stdout <- fmt.Sprintf("等待所有节点运行状态变为Ready(%ds)(%d = %d)", i, total, num)
		if num == total {
			break
		}
		time.Sleep(3 * time.Second)
	}

	i = 0
	for {
		i++
		ik.er.Run(fmt.Sprintf(`kubectl get pods -o wide -n kube-system | grep -wE "%s" | grep calico-node | grep Running | wc -l`, strings.Join(hosts, "|")))
		num, _ := strconv.Atoi(ik.er.GetCmdReturn()[0])
		ik.Stdout <- fmt.Sprintf("等待所有节点calico-node容器正常运行(%ds)(%d = %d)", i, total, num)
		if num == total {
			break
		}
		time.Sleep(3 * time.Second)
	}
	ik.Stdout <- "[结束]启动新node节点"
}

func (ik *InstallK8s) newnodeKubeletcniNode() {
	newnodeRole, ok := ik.resources["newnode"]
	if !ok {
		ik.Stdout <- "没有newnode资源"
		return
	}

	ik.Stdout <- "[开始]修改kubelet配置加载cni"
	ik.kubeletcniNode(newnodeRole)
	ik.Stdout <- "[结束]修改kubelet配置加载cni"
}
