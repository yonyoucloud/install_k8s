package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"sysbase/model"
)

type PodHandler struct{}

func (ph *PodHandler) Create(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("Name"))
	code := strings.TrimSpace(c.PostForm("Code"))
	k8sClusterID := strings.TrimSpace(c.PostForm("K8sClusterID"))
	domain := strings.TrimSpace(c.PostForm("Domain"))
	cap := strings.TrimSpace(c.PostForm("Cap"))
	iaas := strings.TrimSpace(c.PostForm("Iaas"))
	resourceID := strings.TrimSpace(c.PostForm("ResourceID"))
	k8sClusterIDInt, _ := strconv.Atoi(k8sClusterID)
	capInt, _ := strconv.Atoi(cap)

	resp := Response{
		Code: 10000,
		Msg:  "",
	}

	pod := model.Pod{
		Name:         name,
		Code:         code,
		K8sClusterID: uint(k8sClusterIDInt),
		Domain:       domain,
		Cap:          uint32(capInt),
		Iaas:         iaas,
	}

	r, err := pod.Insert(resourceID)
	if err == nil {
		resp.Data = r
	} else {
		resp.Code = 10001
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}

func (ph *PodHandler) List(c *gin.Context) {
	resp := Response{
		Code: 10000,
		Msg:  "",
	}

	pod := model.Pod{}

	r, err := pod.List()
	if err == nil {
		resp.Data = r
	} else {
		resp.Code = 10001
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}

func (ph *PodHandler) Delete(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	idInt, _ := strconv.Atoi(id)

	resp := Response{
		Code: 10000,
		Msg:  "",
	}

	pod := model.Pod{
		ID: uint(idInt),
	}

	err := pod.Delete()
	if err != nil {
		resp.Code = 10001
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}

func (ph *PodHandler) Edit(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	idInt, _ := strconv.Atoi(id)

	name := strings.TrimSpace(c.PostForm("Name"))
	code := strings.TrimSpace(c.PostForm("Code"))
	k8sClusterID := strings.TrimSpace(c.PostForm("K8sClusterID"))
	domain := strings.TrimSpace(c.PostForm("Domain"))
	cap := strings.TrimSpace(c.PostForm("Cap"))
	iaas := strings.TrimSpace(c.PostForm("Iaas"))
	resourceID := strings.TrimSpace(c.PostForm("ResourceID"))
	k8sClusterIDInt, _ := strconv.Atoi(k8sClusterID)
	capInt, _ := strconv.Atoi(cap)

	resp := Response{
		Code: 10000,
		Msg:  "",
	}

	pod := model.Pod{
		ID: uint(idInt),
	}

	podData := model.Pod{
		Name:         name,
		Code:         code,
		K8sClusterID: uint(k8sClusterIDInt),
		Domain:       domain,
		Cap:          uint32(capInt),
		Iaas:         iaas,
	}

	err := pod.Edit(podData, resourceID)
	if err != nil {
		resp.Code = 10001
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}

func (ph *PodHandler) Get(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	idInt, _ := strconv.Atoi(id)

	resp := Response{
		Code: 10000,
		Msg:  "",
	}

	pod := model.Pod{
		ID: uint(idInt),
	}

	r, err := pod.Get()
	if err == nil {
		resp.Data = r
	} else {
		resp.Code = 10001
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}
