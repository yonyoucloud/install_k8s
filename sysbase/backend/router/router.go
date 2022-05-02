package router

import (
	"github.com/gin-gonic/gin"

	"git.yonyou.com/sysbase/backend/config"
	"git.yonyou.com/sysbase/backend/handler"
	"git.yonyou.com/sysbase/backend/router/middleware/auth"
	"git.yonyou.com/sysbase/backend/router/middleware/cors"
)

// InitRouter initialize routing information
func InitRouter(c *config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(cors.Cors())

	router.Static("/static", c.Static)

	resourceGroup := router.Group("/api/v1/resource")
	resourceGroup.Use(cors.JsonHeader(), auth.Auth())
	{
		resourceHandler := &handler.ResourceHandler{}

		resourceGroup.POST("/create", resourceHandler.Create)
		resourceGroup.GET("/list", resourceHandler.List)
		resourceGroup.DELETE("/delete/:id", resourceHandler.Delete)
		resourceGroup.POST("/edit/:id", resourceHandler.Edit)
		resourceGroup.GET("/list/k8sCluster", resourceHandler.ListK8sCluster)
		resourceGroup.GET("/list/pod", resourceHandler.ListPod)
	}

	k8sClusterGroup := router.Group("/api/v1/k8sCluster")
	k8sClusterGroup.Use(cors.JsonHeader(), auth.Auth())
	{
		k8sClusterHandler := &handler.K8sClusterHandler{}

		k8sClusterGroup.POST("/create", k8sClusterHandler.Create)
		k8sClusterGroup.GET("/list", k8sClusterHandler.List)
		k8sClusterGroup.DELETE("/delete/:id", k8sClusterHandler.Delete)
		k8sClusterGroup.POST("/edit/:id", k8sClusterHandler.Edit)
		k8sClusterGroup.GET("/get/:id", k8sClusterHandler.Get)
	}

	podGroup := router.Group("/api/v1/pod")
	podGroup.Use(cors.JsonHeader(), auth.Auth())
	{
		podHandler := &handler.PodHandler{}

		podGroup.POST("/create", podHandler.Create)
		podGroup.GET("/list", podHandler.List)
		podGroup.DELETE("/delete/:id", podHandler.Delete)
		podGroup.POST("/edit/:id", podHandler.Edit)
		podGroup.GET("/get/:id", podHandler.Get)
	}

	k8sCluserResourceGroup := router.Group("/api/v1/k8sClusterResource")
	{
		k8sCluserResourceHandler := &handler.K8sClusterResourceHandler{}
		k8sCluserResourceGroup.GET("/listResource/:id", k8sCluserResourceHandler.ListResource)
	}

	podResourceGroup := router.Group("/api/v1/podResource")
	{
		podResourceHandler := &handler.PodResourceHandler{}
		podResourceGroup.GET("/listResource/:id", podResourceHandler.ListResource)
	}

	tenantPodGroup := router.Group("/api/v1/tenantPod")
	{
		tenantPodHandler := &handler.TenantPodHandler{}
		tenantPodGroup.POST("/open", tenantPodHandler.Open)
		tenantPodGroup.GET("/getByTenantID/:tenantID", tenantPodHandler.GetByTenantID)
	}

	installK8sGroup := router.Group("/api/v1/installK8s")
	{
		installK8sHandler := &handler.InstallK8sHandler{
			Config: c.InstallK8s,
		}
		installK8sGroup.GET("/installTest", installK8sHandler.InstallTest)
		installK8sGroup.GET("/installAll", installK8sHandler.InstallAll)
		installK8sGroup.GET("/installBase", installK8sHandler.InstallBase)
		installK8sGroup.GET("/updateKernel", installK8sHandler.UpdateKernel)
		installK8sGroup.GET("/installBin", installK8sHandler.InstallBin)
		installK8sGroup.GET("/installDocker", installK8sHandler.InstallDocker)
		installK8sGroup.GET("/installPriDocker", installK8sHandler.InstallPriDocker)
		installK8sGroup.GET("/installEtcd", installK8sHandler.InstallEtcd)
		installK8sGroup.GET("/installMaster", installK8sHandler.InstallMaster)
		installK8sGroup.GET("/installNode", installK8sHandler.InstallNode)
		installK8sGroup.GET("/installDockerCrt", installK8sHandler.InstallDockerCrt)
		installK8sGroup.GET("/installLvs", installK8sHandler.InstallLvs)
		installK8sGroup.GET("/installDns", installK8sHandler.InstallDns)
		installK8sGroup.GET("/finishInstall", installK8sHandler.FinishInstall)
		installK8sGroup.GET("/servicePublish", installK8sHandler.ServicePublish)
		installK8sGroup.GET("/serviceEtcd", installK8sHandler.ServiceEtcd)
		installK8sGroup.GET("/serviceMaster", installK8sHandler.ServiceMaster)
		installK8sGroup.GET("/serviceNode", installK8sHandler.ServiceNode)
		installK8sGroup.GET("/serviceDns", installK8sHandler.ServiceDns)
		installK8sGroup.GET("/newnodeInstall", installK8sHandler.NewnodeInstall)
		installK8sGroup.GET("/newetcdInstall", installK8sHandler.NewetcdInstall)
		installK8sGroup.GET("/newmasterInstall", installK8sHandler.NewmasterInstall)
		installK8sGroup.GET("/updateSslMaster", installK8sHandler.UpdateSslMaster)
		installK8sGroup.GET("/updateSslEtcd", installK8sHandler.UpdateSslEtcd)
		installK8sGroup.GET("/updateSslNode", installK8sHandler.UpdateSslNode)
	}

	return router
}
