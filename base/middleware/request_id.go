package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yiran15/api-server/base/constant"
	"github.com/yiran15/api-server/base/helper"
	"go.uber.org/zap"
)

func (m *Middleware) RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get(constant.RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
			zap.L().Debug("Generated new request ID", zap.String("request_id", requestID), zap.String("path", c.Request.URL.Path))
		}
		c.Set(constant.RequestIDContextKey, requestID)
		c.Header(constant.RequestIDHeader, requestID)
		ctx := context.WithValue(c.Request.Context(), helper.RequestIDContextKey{}, requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
