package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yiran15/api-server/base/constant"
)

// Auth 是一个基于 JWT 的认证中间件
func (m *Middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			m.Abort(c, http.StatusUnauthorized, constant.ErrNoAuthHeader)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			m.Abort(c, http.StatusUnauthorized, constant.ErrInvalidAuthHeader)
			return
		}

		tokenString := parts[1]
		mc, err := m.jwtImpl.ParseToken(tokenString)
		if err != nil {
			m.Abort(c, http.StatusUnauthorized, err)
			return
		}
		ctx := context.WithValue(c.Request.Context(), constant.UserContextKey, mc)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func (m *Middleware) Abort(c *gin.Context, code int, err error) {
	switch code {
	case http.StatusUnauthorized:
		c.JSON(code, gin.H{
			"code":  code,
			"msg":   "unauthorized",
			"error": err.Error(),
		})
	case http.StatusForbidden:
		c.JSON(code, gin.H{
			"code":  code,
			"msg":   "forbidden",
			"error": err.Error(),
		})
	}
	c.Error(err)
	c.Abort()
}
