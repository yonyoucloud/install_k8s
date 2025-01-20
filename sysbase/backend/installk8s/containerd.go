package installk8s

import (
	"fmt"
	"strings"

	"sysbase/tool/execremote"
)

func (ik *InstallK8s) InstallContainerd() {
	var role []execremote.Role

	if r, ok := ik.resources["publish"]; ok {
		role = append(role, r)
	}
	if r, ok := ik.resources["master"]; ok {
		role = append(role, r)
	}
	if r, ok := ik.resources["node"]; ok {
		role = append(role, r)
	}
	if r, ok := ik.resources["registry"]; ok {
		role = append(role, r)
	}

	ik.er.SetRole(role...)
	ik.installContainerd()
}

func (ik *InstallK8s) installContainerd() {
	ik.er.Put(fmt.Sprintf("%s/containerd/package.gz", ik.SourceDir), "/tmp")
	ik.er.Run("yum install -y libseccomp && tar zxvf /tmp/package.gz -C / && rm -rf /tmp/package.gz && mkdir -p /data/containerd && systemctl daemon-reload && systemctl enable containerd")
}

func (ik *InstallK8s) InstallRegistry() {
	registryRole, ok := ik.resources["registry"]
	if !ok {
		ik.Stdout <- "没有registry资源"
		return
	}

	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	// 获取一下发布机上的镜像文件路径名称
	publishRole.WaitOutput = true
	ik.er.SetRole(publishRole)
	ik.er.Run(fmt.Sprintf("ls %s/images/registry.k8s.io/docker.io/library/registry/* | head -n 1", ik.SourceDir))
	registryImage := ik.er.GetCmdReturn()[0]
	publishRole.WaitOutput = false

	ik.er.SetRole(registryRole)

	// 如果是在发布机上运行，此步骤不需要执行
	if !strInArr(strings.Split(publishRole.Hosts[0], ":")[0], localIps) {
		ik.er.Get(fmt.Sprintf("%s/containerd/create_ssl.sh", ik.SourceDir), fmt.Sprintf("%s/containerd", ik.SourceDir))
		ik.er.Get(fmt.Sprintf("%s/containerd/start_registry.sh", ik.SourceDir), fmt.Sprintf("%s/containerd", ik.SourceDir))
	}

	ik.er.Put(fmt.Sprintf("%s/containerd/create_ssl.sh", ik.SourceDir), "/tmp")
	ik.er.Run("sh /tmp/create_ssl.sh && rm -rf /tmp/create_ssl.sh")
	ik.er.Get("/etc/certs/registry.k8s.io.crt", fmt.Sprintf("%s/containerd", ik.SourceDir))

	ik.er.Run("systemctl restart containerd")
	ik.er.Put(registryImage, "/tmp/registry.image.tar")
	ik.er.Put(fmt.Sprintf("%s/containerd/start_registry.sh", ik.SourceDir), "/tmp")
	cmds := []string{
		`nerdctl -n k8s.io load -i /tmp/registry.image.tar`,
		`/tmp/start_registry.sh ; rm -rf /tmp/registry.image.tar /tmp/start_registry.sh`,
	}
	ik.er.Run(cmds...)

	ik.er.SetRole(publishRole)
	ik.er.Local(fmt.Sprintf("cd %s/containerd && mkdir -p etc/certs && chmod 755 etc etc/certs && chmod 640 registry.k8s.io.crt && mv registry.k8s.io.crt etc/certs && tar zcvf registrycrt.gz etc && rm -rf etc", ik.SourceDir))
	ik.er.Put(fmt.Sprintf("%s/containerd/registrycrt.gz", ik.SourceDir), fmt.Sprintf("%s/containerd/registrycrt.gz", ik.SourceDir))
}

func (ik *InstallK8s) InstallContainerdCrt() {
	publishRole, ok := ik.resources["publish"]
	if !ok {
		ik.Stdout <- "没有publish资源"
		return
	}

	nodeRole, ok := ik.resources["node"]
	if !ok {
		ik.Stdout <- "没有node资源"
		return
	}

	ik.installContainerdCrt(publishRole, nodeRole)
}

func (ik *InstallK8s) installContainerdCrt(publishRole execremote.Role, role ...execremote.Role) {
	// 如果是在发布机上运行，此步骤不需要执行
	if !strInArr(strings.Split(publishRole.Hosts[0], ":")[0], localIps) {
		ik.er.Get(fmt.Sprintf("%s/containerd/registrycrt.gz", ik.SourceDir), fmt.Sprintf("%s/containerd", ik.SourceDir))
	}

	ik.er.SetRole(role...)
	ik.er.Put(fmt.Sprintf("%s/containerd/registrycrt.gz", ik.SourceDir), "/tmp")
	ik.er.Run("tar zxvf /tmp/registrycrt.gz -C / && rm -rf /tmp/registrycrt.gz")
}
