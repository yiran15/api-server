package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yiran15/api-server/base/constant"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore" // 引入 zapcore
)

// ZapLogger 返回一个 Gin 中间件，用于记录请求日志
func (m *Middleware) ZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		requestID := c.GetString(constant.RequestIDContextKey)
		end := time.Now()
		latency := end.Sub(start)
		latencyMs := latency.Milliseconds()
		statusCode := c.Writer.Status()

		fields := []zapcore.Field{
			zap.Int("status", statusCode),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Int64("latency_ms", latencyMs),
			zap.String("request_id", requestID),
		}

		// 记录 Gin 上下文中的错误
		// c.Errors 是一个 *gin.Error 类型的切片
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("error", c.Errors[0].Error()))
			// 根据状态码和错误情况决定日志级别
			zap.L().Error("request failed", fields...)
		} else {
			zap.L().Info("request success", fields...)

		}
	}
}
