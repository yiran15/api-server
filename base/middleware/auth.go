package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/yiran15/api-server/base/apitypes"
	"github.com/yiran15/api-server/base/constant"
	"github.com/yiran15/api-server/base/log"
	"go.uber.org/zap"
)

// Auth 是一个基于 JWT 的认证中间件
func (m *Middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			log.WithRequestID(c.Request.Context()).Error("auth failed, no Authorization header")
			m.Abort(c, http.StatusUnauthorized, constant.ErrAuthFailed)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			log.WithRequestID(c.Request.Context()).Error("auth failed, invalid Authorization header")
			m.Abort(c, http.StatusUnauthorized, constant.ErrAuthFailed)
			return
		}

		tokenString := parts[1]
		mc, err := m.jwtImpl.ParseToken(tokenString)
		if err != nil {
			log.WithRequestID(c.Request.Context()).Error("auth failed, parse token failed", zap.Error(err))
			m.Abort(c, http.StatusUnauthorized, constant.ErrAuthFailed)
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
		c.JSON(code, apitypes.NewResponse(code, fmt.Sprintf("unauthorized, %v", err), requestid.Get(c), nil))
	case http.StatusForbidden:
		c.JSON(code, apitypes.NewResponse(code, fmt.Sprintf("forbidden, %v", err), requestid.Get(c), nil))
	}
	c.Error(err)
	c.Abort()
}
