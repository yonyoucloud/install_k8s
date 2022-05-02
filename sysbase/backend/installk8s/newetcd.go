package installk8s

import (
	"fmt"
	"strings"

	"git.yonyou.com/sysbase/backend/model"
)

func (ik *InstallK8s) NewetcdInstall() {
	newetcdRole, ok := ik.resources["newetcd"]
	if !ok {
		ik.Stdout <- "没有newetcd资源"
		return
	}

	ik.Stdout <- "[开始]关闭master服务及etcd服务"
	ik.Params.DoWhat = "stop"
	ik.ServiceMaster()
	ik.ServiceEtcd()
	ik.Stdout <- "[结束]关闭master服务及etcd服务"

	ik.Stdout <- "[开始]安装etcd节点"
	ik.newetcdInstallBase()
	ik.newetcdInstallEtcd()
	ik.newetcdModifyEtcdConf()
	ik.updateMasterEtcdSsl()
	ik.Stdout <- "[结束]安装etcd节点"

	ik.Stdout <- "[开始]启动新etcd节点服务及master服务"
	ik.Params.DoWhat = "restart"
	ik.ServiceEtcd()
	ik.wait(8, "等待...")
	ik.newetcdClusterAddnew()
	ik.newetcdServiceEtcdStart()
	ik.newetcdClusterAddnewCheck()
	ik.ServiceMaster()
	ik.Stdout <- "[结束]启动新etcd节点服务及master服务"

	ik.Stdout <- "[开始]修改lvs配置"
	ik.installLvsvipEtcd(newetcdRole)
	ik.installLvsNew()
	ik.Stdout <- "[结束]修改lvs配置"

	ik.Stdout <- "[开始]重新配置calico"
	ik.initCalico()
	ik.Stdout <- "[结束]重新配置calico"

	ik.Stdout <- "[开始]修改scope，由newetcd改为etcd"
	ik.newetcdUpdateScope()
	ik.Stdout <- "[结束]修改scope，由newetcd改为etcd"

	ik.Stdout <- "[结束]所有节点添加完毕"
}

func (ik *InstallK8s) newetcdUpdateScope() {
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

	etcdLbHost := strings.Split(etcdLbRole.Hosts[0], ":")[0]

	k8sClusterResource := model.K8sClusterResource{}
	resource, err := k8sClusterResource.ListResource(ik.Params.K8sClusterID, []string{"newetcd"})
	if err != nil {
		ik.Stdout <- err.Error()
		return
	}

	publishRole.WaitOutput = true
	ik.er.SetRole(publishRole)

	for _, r := range resource {
		ik.er.Run(fmt.Sprintf(`etcdctl --cacert=%s/etcd/etc/etcd/ssl/ca.pem --cert=%s/etcd/etc/etcd/ssl/etcd.pem --key=%s/etcd/etc/etcd/ssl/etcd-key.pem --endpoints=https://%s:2379 member list | grep %s > /dev/null ; echo $?`, ik.SourceDir, ik.SourceDir, ik.SourceDir, etcdLbHost, r.Host))
		if ik.er.GetCmdReturn()[0] == "0" {
			model.Resource{
				ID: r.ID,
			}.Edit(model.Resource{
				Scope: "etcd",
			})
		}
	}
}

func (ik *InstallK8s) newetcdInstallBase() {
	newetcdRole, ok := ik.resources["newetcd"]
	if !ok {
		ik.Stdout <- "没有newetcd资源"
		return
	}

	ik.Stdout <- "[开始]安装基础环境"
	ik.er.SetRole(newetcdRole)
	ik.installBase()
	ik.Stdout <- "[结束]安装基础环境"
}

func (ik *InstallK8s) newetcdInstallEtcd() {
	ik.AddNew = true
	ik.InstallEtcd()
}

func (ik *InstallK8s) newetcdModifyEtcdConf() {
	newetcdRole, ok := ik.resources["newetcd"]
	if !ok {
		ik.Stdout <- "没有newetcd资源"
		return
	}

	ik.er.SetRole(newetcdRole)
	ik.er.Run(`rm -rf /data/etcd/* ; sed -i 's#new#existing#g' /usr/lib/systemd/system/etcd.service`)
}

func (ik *InstallK8s) newetcdClusterAddnew() {
	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	newetcdRole, ok := ik.resources["newetcd"]
	if !ok {
		ik.Stdout <- "没有newetcd资源"
		return
	}

	etcdRole, ok := ik.resources["etcd"]
	if !ok {
		ik.Stdout <- "没有etcd资源"
		return
	}

	etcdLbRole, ok := ik.resources["etcdlb"]
	if !ok {
		ik.Stdout <- "没有etcdlb资源"
		return
	}

	etcdIndex := len(etcdRole.Hosts)
	etcdLbHost := strings.Split(etcdLbRole.Hosts[0], ":")[0]
	ik.er.SetRole(publishRole)
	for _, host := range newetcdRole.Hosts {
		etcdIndex++
		cmd := fmt.Sprintf(`etcdctl --cacert=%s/etcd/etc/etcd/ssl/ca.pem --cert=%s/etcd/etc/etcd/ssl/etcd.pem --key=%s/etcd/etc/etcd/ssl/etcd-key.pem --endpoints=https://%s:2379 member add etcd%d --peer-urls=https://%s:2380`, ik.SourceDir, ik.SourceDir, ik.SourceDir, etcdLbHost, etcdIndex, strings.Split(host, ":")[0])
		ik.er.Run(cmd)
	}
}

func (ik *InstallK8s) newetcdServiceEtcdStart() {
	newetcdRole, ok := ik.resources["newetcd"]
	if !ok {
		ik.Stdout <- "没有newetcd资源"
		return
	}

	ik.Params.DoWhat = "start"
	ik.serviceEtcd(newetcdRole)
}

func (ik *InstallK8s) newetcdClusterAddnewCheck() {
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

	etcdLbHost := strings.Split(etcdLbRole.Hosts[0], ":")[0]
	ik.er.SetRole(publishRole)
	cmds := []string{
		fmt.Sprintf(`etcdctl --cacert=%s/etcd/etc/etcd/ssl/ca.pem --cert=%s/etcd/etc/etcd/ssl/etcd.pem --key=%s/etcd/etc/etcd/ssl/etcd-key.pem --endpoints=https://%s:2379 member list`, ik.SourceDir, ik.SourceDir, ik.SourceDir, etcdLbHost),
		fmt.Sprintf(`etcdctl --cacert=%s/etcd/etc/etcd/ssl/ca.pem --cert=%s/etcd/etc/etcd/ssl/etcd.pem --key=%s/etcd/etc/etcd/ssl/etcd-key.pem --endpoints=https://%s:2379 endpoint health`, ik.SourceDir, ik.SourceDir, ik.SourceDir, etcdLbHost),
	}
	ik.er.Run(cmds...)
}
