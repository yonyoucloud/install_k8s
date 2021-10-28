package handler

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"git.yonyou.com/sysbase/backend/config"
	"git.yonyou.com/sysbase/backend/installk8s"
	"git.yonyou.com/sysbase/backend/tool/waitoutput"
)

type InstallK8sHandler struct {
	Config config.InstallK8s
}

var waitOutput *waitoutput.WaitOutput

func init() {
	waitOutput = &waitoutput.WaitOutput{}
}

func (ikh *InstallK8sHandler) InstallAll(c *gin.Context) {
	ikh.call(c, "InstallAll")
}

func (ikh *InstallK8sHandler) InstallBase(c *gin.Context) {
	ikh.call(c, "InstallBase")
}

func (ikh *InstallK8sHandler) UpdateKernel(c *gin.Context) {
	ikh.call(c, "UpdateKernel")
}

func (ikh *InstallK8sHandler) InstallBin(c *gin.Context) {
	ikh.call(c, "InstallBin")
}

func (ikh *InstallK8sHandler) InstallDocker(c *gin.Context) {
	ikh.call(c, "InstallDocker")
}

func (ikh *InstallK8sHandler) InstallPriDocker(c *gin.Context) {
	ikh.call(c, "InstallPriDocker")
}

func (ikh *InstallK8sHandler) InstallEtcd(c *gin.Context) {
	ikh.call(c, "InstallEtcd")
}

func (ikh *InstallK8sHandler) InstallMaster(c *gin.Context) {
	ikh.call(c, "InstallMaster")
}

func (ikh *InstallK8sHandler) InstallNode(c *gin.Context) {
	ikh.call(c, "InstallNode")
}

func (ikh *InstallK8sHandler) InstallDockerCrt(c *gin.Context) {
	ikh.call(c, "InstallDockerCrt")
}

func (ikh *InstallK8sHandler) InstallLvs(c *gin.Context) {
	ikh.call(c, "InstallLvs")
}

func (ikh *InstallK8sHandler) InstallDns(c *gin.Context) {
	ikh.call(c, "InstallDns")
}

func (ikh *InstallK8sHandler) ServicePublish(c *gin.Context) {
	ikh.call(c, "ServicePublish")
}

func (ikh *InstallK8sHandler) ServiceEtcd(c *gin.Context) {
	ikh.call(c, "ServiceEtcd")
}

func (ikh *InstallK8sHandler) ServiceMaster(c *gin.Context) {
	ikh.call(c, "ServiceMaster")
}

func (ikh *InstallK8sHandler) ServiceNode(c *gin.Context) {
	ikh.call(c, "ServiceNode")
}

func (ikh *InstallK8sHandler) ServiceDns(c *gin.Context) {
	ikh.call(c, "ServiceDns")
}

func (ikh *InstallK8sHandler) FinishInstall(c *gin.Context) {
	ikh.call(c, "FinishInstall")
}

func (ikh *InstallK8sHandler) NewnodeInstall(c *gin.Context) {
	ikh.call(c, "NewnodeInstall")
}

func (ikh *InstallK8sHandler) NewetcdInstall(c *gin.Context) {
	ikh.call(c, "NewetcdInstall")
}

func (ikh *InstallK8sHandler) NewmasterInstall(c *gin.Context) {
	ikh.call(c, "NewmasterInstall")
}

func (ikh *InstallK8sHandler) UpdateSslMaster(c *gin.Context) {
	ikh.call(c, "UpdateSslMaster")
}

func (ikh *InstallK8sHandler) UpdateSslEtcd(c *gin.Context) {
	ikh.call(c, "UpdateSslEtcd")
}

func (ikh *InstallK8sHandler) UpdateSslNode(c *gin.Context) {
	ikh.call(c, "UpdateSslNode")
}

func (ikh *InstallK8sHandler) call(c *gin.Context, callFunc string) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")

	id := strings.TrimSpace(c.Query("k8s_cluster_id"))
	doWhat := strings.TrimSpace(c.Query("do_what"))
	doDocker := strings.TrimSpace(c.Query("do_docker"))
	idInt, _ := strconv.Atoi(id)
	doDockerBool := doDocker == "true"

	installCmd := fmt.Sprintf("InstallK8sHandler.%s.%s", callFunc, id)
	running := waitOutput.IsRunning(installCmd)

	if !running {
		waitOutput.SetDataChan(installCmd)
	}

	stdout := waitOutput.GetDataChan(installCmd)
	if stdout == nil {
		c.Stream(func(w io.Writer) bool {
			fmt.Fprint(w, "获取管道错误")
			return false
		})
		return
	}

	if !running {
		ik := &installk8s.InstallK8s{
			SourceDir: ikh.Config.SourceDir,
			Params: installk8s.Params{
				K8sClusterID: uint(idInt),
				DoWhat:       doWhat,
				DoDocker:     doDockerBool,
			},
			Stdout: stdout,
			Defer: func() {
				waitOutput.DeleteByCmd(installCmd)
			},
		}
		ik.GetResources()
		go ik.Call(callFunc)
	}

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-stdout; ok {
			// fmt.Fprint(w, msg)
			c.SSEvent("message", msg)
			return true
		}
		c.SSEvent("close", "")
		return false
	})
}
