package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"sysbase/model"
)

type K8sClusterHandler struct{}

func (kch *K8sClusterHandler) Create(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("Name"))
	resourceID := strings.TrimSpace(c.PostForm("ResourceID"))

	resp := Response{
		Code: 10000,
		Msg:  "",
	}

	k8sCluster := model.K8sCluster{
		Name: name,
	}

	r, err := k8sCluster.Insert(resourceID)
	if err == nil {
		resp.Data = r
	} else {
		resp.Code = 10001
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}

func (kch *K8sClusterHandler) List(c *gin.Context) {
	resp := Response{
		Code: 10000,
		Msg:  "",
	}

	k8sCluster := model.K8sCluster{}

	r, err := k8sCluster.List()
	if err == nil {
		resp.Data = r
	} else {
		resp.Code = 10001
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}

func (kch *K8sClusterHandler) Delete(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	idInt, _ := strconv.Atoi(id)

	resp := Response{
		Code: 10000,
		Msg:  "",
	}

	k8sCluster := model.K8sCluster{
		ID: uint(idInt),
	}

	err := k8sCluster.Delete()
	if err != nil {
		resp.Code = 10001
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}

func (kch *K8sClusterHandler) Edit(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	idInt, _ := strconv.Atoi(id)

	name := strings.TrimSpace(c.PostForm("Name"))
	resourceID := strings.TrimSpace(c.PostForm("ResourceID"))

	resp := Response{
		Code: 10000,
		Msg:  "",
	}

	k8sCluster := model.K8sCluster{
		ID: uint(idInt),
	}

	k8sClusterData := model.K8sCluster{
		Name: name,
	}

	err := k8sCluster.Edit(k8sClusterData, resourceID)
	if err != nil {
		resp.Code = 10001
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}

func (kch *K8sClusterHandler) Get(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	idInt, _ := strconv.Atoi(id)

	resp := Response{
		Code: 10000,
		Msg:  "",
	}

	k8sCluster := model.K8sCluster{
		ID: uint(idInt),
	}

	r, err := k8sCluster.Get()
	if err == nil {
		resp.Data = r
	} else {
		resp.Code = 10001
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}
