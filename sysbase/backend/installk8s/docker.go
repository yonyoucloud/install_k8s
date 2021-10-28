package installk8s

import (
	"fmt"
	"strings"

	"git.yonyou.com/sysbase/backend/tool/execremote"
)

func (ik *InstallK8s) InstallDocker() {
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
	if r, ok := ik.resources["pridocker"]; ok {
		role = append(role, r)
	}

	ik.er.SetRole(role...)
	ik.installDocker()
}

func (ik *InstallK8s) installDocker() {
	ik.er.Put(fmt.Sprintf("%s/docker/conf_bin.gz", ik.SourceDir), "/tmp")
	ik.er.Run("yum install -y libseccomp && tar zxvf /tmp/conf_bin.gz -C / && rm -rf /tmp/conf_bin.gz && mkdir -p /data/docker && systemctl daemon-reload && systemctl enable docker")
}

func (ik *InstallK8s) InstallPriDocker() {
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

	priDockerHost := strings.Split(pridockerRole.Hosts[0], ":")[0]

	ik.er.SetRole(publishRole)
	cmds := []string{
		fmt.Sprintf(`cd %s/docker && sed "s#HOST#%s#g" create_ssl.sh.tpl > create_ssl.sh && chmod 750 create_ssl.sh`, ik.SourceDir, priDockerHost),
		fmt.Sprintf(`cd %s/docker && sed "s#HOST#%s#g" start_registry.sh.tpl > start_registry.sh && chmod 750 start_registry.sh`, ik.SourceDir, priDockerHost),
		fmt.Sprintf("cd %s/docker && rm -rf ca.crt", ik.SourceDir),
	}
	ik.er.Run(cmds...)

	// 如果是在发布机上运行，此步骤不需要执行
	if !strInArr(strings.Split(publishRole.Hosts[0], ":")[0], localIps) {
		ik.er.Get(fmt.Sprintf("%s/docker/create_ssl.sh", ik.SourceDir), fmt.Sprintf("%s/docker", ik.SourceDir))
		ik.er.Get(fmt.Sprintf("%s/docker/start_registry.sh", ik.SourceDir), fmt.Sprintf("%s/docker", ik.SourceDir))
	}

	ik.er.SetRole(pridockerRole)
	ik.er.Put(fmt.Sprintf("%s/docker/create_ssl.sh", ik.SourceDir), "/tmp")
	ik.er.Run("sh /tmp/create_ssl.sh && rm -rf /tmp/create_ssl.sh")
	ik.er.Local(fmt.Sprintf("if [ ! -f %s/docker/ca.crt ]; then touch %s/docker/ca.crt; fi", ik.SourceDir, ik.SourceDir))
	ik.er.Get(fmt.Sprintf("/etc/certs/%s.crt", priDockerHost), fmt.Sprintf("%s/docker/ca.crt", ik.SourceDir))

	ik.er.Run("systemctl restart docker")
	ik.er.Put(fmt.Sprintf("%s/images/registry:2.7.1.tar", ik.SourceDir), "/tmp")
	ik.er.Put(fmt.Sprintf("%s/docker/start_registry.sh", ik.SourceDir), "/tmp")
	ik.er.Run(`/tmp/start_registry.sh ; rm -rf /tmp/start_registry.sh`)

	ik.er.SetRole(publishRole)
	ik.er.Local(fmt.Sprintf("chmod 640 %s/docker/ca.crt", ik.SourceDir))
	ik.er.Put(fmt.Sprintf("%s/docker/ca.crt", ik.SourceDir), fmt.Sprintf("%s/docker/ca.crt", ik.SourceDir))
}

func (ik *InstallK8s) InstallDockerCrt() {
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

	nodeRole, ok := ik.resources["node"]
	if !ok {
		ik.Stdout <- "没有node资源"
		return
	}

	priDockerHost := strings.Split(pridockerRole.Hosts[0], ":")[0]

	publishRole.Parallel = false
	nodeRole.Parallel = false

	ik.installDockerCrt(priDockerHost, publishRole, publishRole, nodeRole)
}

func (ik *InstallK8s) installDockerCrt(priDockerHost string, publishRole execremote.Role, role ...execremote.Role) {
	ik.er.SetRole(publishRole)
	cmds := []string{
		fmt.Sprintf(`cd %s/docker && rm -rf etc/docker/certs.d/* && chmod 640 ca.crt && mkdir etc/docker/certs.d/%s:5000 && chmod 750 etc/docker/certs.d/%s:5000 && /usr/bin/cp -rpf ca.crt etc/docker/certs.d/%s:5000`, ik.SourceDir, priDockerHost, priDockerHost, priDockerHost),
		fmt.Sprintf(`cd %s/docker && tar zcvf docker.gz etc`, ik.SourceDir),
	}
	ik.er.Run(cmds...)

	// 如果是在发布机上运行，此步骤不需要执行
	if !strInArr(strings.Split(publishRole.Hosts[0], ":")[0], localIps) {
		ik.er.Get(fmt.Sprintf("%s/docker/docker.gz", ik.SourceDir), fmt.Sprintf("%s/docker", ik.SourceDir))
	}

	ik.er.SetRole(role...)

	ik.er.Put(fmt.Sprintf("%s/docker/docker.gz", ik.SourceDir), "/tmp")
	ik.er.Run("tar zxvf /tmp/docker.gz -C / && rm -rf /tmp/docker.gz")
	ik.er.Local(fmt.Sprintf("rm -rf %s/docker/docker.gz", ik.SourceDir))
}
