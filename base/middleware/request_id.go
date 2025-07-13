package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yiran15/api-server/base/constant"
	"go.uber.org/zap"
)

func (m *Middleware) RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-Id")
		if requestID != "" {
			c.Set(constant.RequestID, requestID)
			c.Header("X-Request-Id", requestID)
		}

		if requestID == "" {
			requestID = uuid.New().String()
			zap.L().Info("generate request id", zap.String("requestID", requestID))
			c.Set(constant.RequestID, requestID)
			c.Header("X-Request-Id", requestID)
		}

		c.Next()
	}
}
