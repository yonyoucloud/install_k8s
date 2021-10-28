package cors

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 跨域Option支持中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With, User-Agent, jweToken")

		if !strings.Contains(c.Request.URL.Path, "/kubeapi/sockjs") {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		origin := c.Request.Header.Get("Origin")
		// if IsOriginOk(origin) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		// }

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}

func JsonHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
	}
}
