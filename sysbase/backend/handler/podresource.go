package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"sysbase/model"
)

type PodResourceHandler struct{}

func (prh *PodResourceHandler) ListResource(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	idInt, _ := strconv.Atoi(id)

	resp := Response{
		Code: 10000,
		Msg:  "",
	}

	podResource := model.PodResource{}

	r, err := podResource.ListResource(uint(idInt))
	if err == nil {
		resp.Data = r
	} else {
		resp.Code = 10001
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}
