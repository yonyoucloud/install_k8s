package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"git.yonyou.com/sysbase/backend/model"
)

type ResourceHandler struct{}

func (rh *ResourceHandler) Create(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("Name"))
	category := strings.TrimSpace(c.PostForm("Category"))
	scope := strings.TrimSpace(c.PostForm("Scope"))
	host := strings.TrimSpace(c.PostForm("Host"))
	port := strings.TrimSpace(c.PostForm("Port"))
	user := strings.TrimSpace(c.PostForm("User"))
	password := strings.TrimSpace(c.PostForm("Password"))
	portInt, _ := strconv.Atoi(port)

	resp := Response{
		Code: 10000,
		Msg:  "",
	}

	resource := model.Resource{
		Name:     name,
		Category: category,
		Scope:    scope,
		Host:     host,
		Port:     uint32(portInt),
		User:     user,
		Password: password,
	}

	r, err := resource.Insert()
	if err == nil {
		resp.Data = r
	} else {
		resp.Code = 10001
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}

func (rh *ResourceHandler) List(c *gin.Context) {
	resp := Response{
		Code: 10000,
		Msg:  "",
	}

	resource := model.Resource{}

	r, err := resource.List()
	if err == nil {
		resp.Data = r
	} else {
		resp.Code = 10001
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}

func (rh *ResourceHandler) Delete(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	idInt, _ := strconv.Atoi(id)

	resp := Response{
		Code: 10000,
		Msg:  "",
	}

	resource := model.Resource{
		ID: uint(idInt),
	}

	err := resource.Delete()
	if err != nil {
		resp.Code = 10001
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}

func (rh *ResourceHandler) Edit(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	idInt, _ := strconv.Atoi(id)

	name := strings.TrimSpace(c.PostForm("Name"))
	category := strings.TrimSpace(c.PostForm("Category"))
	scope := strings.TrimSpace(c.PostForm("Scope"))
	host := strings.TrimSpace(c.PostForm("Host"))
	port := strings.TrimSpace(c.PostForm("Port"))
	user := strings.TrimSpace(c.PostForm("User"))
	password := strings.TrimSpace(c.PostForm("Password"))
	portInt, _ := strconv.Atoi(port)

	resp := Response{
		Code: 10000,
		Msg:  "",
	}

	resource := model.Resource{
		ID: uint(idInt),
	}

	resourceData := model.Resource{
		Name:     name,
		Category: category,
		Scope:    scope,
		Host:     host,
		Port:     uint32(portInt),
		User:     user,
		Password: password,
	}

	err := resource.Edit(resourceData)
	if err != nil {
		resp.Code = 10001
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}

func (rh *ResourceHandler) ListK8sCluster(c *gin.Context) {
	resp := Response{
		Code: 10000,
		Msg:  "",
	}

	resource := model.Resource{}

	r, err := resource.ListK8sCluster()
	if err == nil {
		resp.Data = r
	} else {
		resp.Code = 10001
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}

func (rh *ResourceHandler) ListPod(c *gin.Context) {
	resp := Response{
		Code: 10000,
		Msg:  "",
	}

	podID := strings.TrimSpace(c.Query("podID"))
	podIDInt, _ := strconv.Atoi(podID)

	resource := model.Resource{}

	r, err := resource.ListPod(uint(podIDInt))
	if err == nil {
		resp.Data = r
	} else {
		resp.Code = 10001
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}
