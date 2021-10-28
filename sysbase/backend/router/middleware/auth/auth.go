package auth

import (
	"github.com/gin-gonic/gin"
)

// 认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()

	}
}
