package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/yiran15/api-server/base/constant"
	"go.uber.org/zap"
)

// Auth 是一个基于 JWT 的认证中间件
func (m *Middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			zap.L().Error("auth failed, no Authorization header", zap.String("request-id", requestid.Get(c)))
			m.Abort(c, http.StatusUnauthorized, constant.ErrAuthFailed)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			zap.L().Error("auth failed, invalid Authorization header", zap.String("request-id", requestid.Get(c)))
			m.Abort(c, http.StatusUnauthorized, constant.ErrAuthFailed)
			return
		}

		tokenString := parts[1]
		mc, err := m.jwtImpl.ParseToken(tokenString)
		if err != nil {
			zap.L().Error("auth failed, parse token failed", zap.String("request-id", requestid.Get(c)), zap.Error(err))
			m.Abort(c, http.StatusUnauthorized, constant.ErrAuthFailed)
			return
		}
		ctx := context.WithValue(c.Request.Context(), constant.UserContextKey, mc)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
