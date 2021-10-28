package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"git.yonyou.com/sysbase/backend/model"
)

type TenantPodHandler struct{}

func (tph *TenantPodHandler) Open(c *gin.Context) {
	tenantID := strings.TrimSpace(c.PostForm("tenantID"))
	podID := strings.TrimSpace(c.PostForm("podID"))
	tenantName := strings.TrimSpace(c.PostForm("tenantName"))
	tenantIDInt, _ := strconv.Atoi(tenantID)
	podIDInt, _ := strconv.Atoi(podID)

	resp := Response{
		Code: 10000,
		Msg:  "",
	}

	tenantPod := model.TenantPod{
		TenantID:   uint(tenantIDInt),
		PodID:      uint(podIDInt),
		TenantName: tenantName,
	}

	r, err := tenantPod.Insert()
	if err == nil {
		resp.Data = r
	} else {
		resp.Code = 10001
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}

func (tph *TenantPodHandler) GetByTenantID(c *gin.Context) {
	tenantID := strings.TrimSpace(c.Param("tenantID"))
	tenantIDInt, _ := strconv.Atoi(tenantID)

	resp := Response{
		Code: 10000,
		Msg:  "",
	}

	tenantPod := model.TenantPod{
		TenantID: uint(tenantIDInt),
	}

	r, err := tenantPod.Get()
	if err == nil {
		resp.Data = r
	} else {
		resp.Code = 10001
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}
