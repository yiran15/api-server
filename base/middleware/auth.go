package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yiran15/api-server/base/constant"
)

// Authentication 基于JWT的认证中间件
func (a *Middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			_ = c.Error(errors.New("auth header is empty"))
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			_ = c.Error(errors.New("invalid auth header"))
			return
		}
		mc, err := a.jwtImpl.ParseToken(parts[1])
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.Set(constant.AuthMidwareKey, mc)
		c.Next()
	}
}
