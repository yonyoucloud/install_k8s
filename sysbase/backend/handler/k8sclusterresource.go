package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"git.yonyou.com/sysbase/backend/model"
)

type K8sClusterResourceHandler struct{}

func (kcrh *K8sClusterResourceHandler) ListResource(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	idInt, _ := strconv.Atoi(id)

	resp := Response{
		Code: 10000,
		Msg:  "",
	}

	k8sClusterResource := model.K8sClusterResource{}

	r, err := k8sClusterResource.ListResource(uint(idInt), []string{})
	if err == nil {
		resp.Data = r
	} else {
		resp.Code = 10001
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}
