// package middleware
package middleware

import (
	"time"

	"github.com/gin-contrib/cors" // 导入 Gin CORS 中间件
	"github.com/gin-gonic/gin"
)

// CorsOption 定义了 CORS 配置的类型
type CorsOption int

const (
	// CorsAllowAll 允许所有来源、方法和头部的最宽松配置 (开发环境慎用)
	CorsAllowAll CorsOption = iota
	// CorsSpecificOrigins 允许指定来源的配置 (推荐用于生产环境)
	CorsSpecificOrigins
	// CorsAllowCredentialsWithSpecificOrigins 允许指定来源且支持凭证的配置
	CorsAllowCredentialsWithSpecificOrigins
)

// CorsMiddleware 返回一个 Gin CORS 中间件
// 根据传入的 CorsOption 提供不同的 CORS 策略
func (m *Middleware) Cors(option CorsOption, allowedOrigins ...string) gin.HandlerFunc {
	switch option {
	case CorsAllowAll:
		return cors.New(cors.Config{
			AllowAllOrigins:  true,                                                                              // 允许所有来源
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},                      // 允许所有常用方法
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}, // 允许所有常用请求头
			ExposeHeaders:    []string{"Content-Length"},                                                        // 允许客户端访问的响应头
			AllowCredentials: true,                                                                              // 允许发送 Cookie
			MaxAge:           12 * time.Hour,                                                                    // 预检请求的缓存时间
		})

	case CorsSpecificOrigins:
		if len(allowedOrigins) == 0 {
			// 如果没有指定来源，则默认只允许当前域名，或者返回一个错误
			// 实际项目中可以根据需求调整，这里为了示例简单，返回一个限制性配置
			return cors.New(cors.Config{
				AllowOrigins:     []string{"http://localhost"}, // 默认一个本地开发地址
				AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
				ExposeHeaders:    []string{"Content-Length"},
				AllowCredentials: false, // 默认不允许凭证，如果需要请使用 CorsAllowCredentialsWithSpecificOrigins
				MaxAge:           12 * time.Hour,
			})
		}
		return cors.New(cors.Config{
			AllowOrigins:     allowedOrigins, // 允许指定来源
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: false, // 默认不允许凭证
			MaxAge:           12 * time.Hour,
		})

	case CorsAllowCredentialsWithSpecificOrigins:
		if len(allowedOrigins) == 0 {
			// 警告：允许凭证时，AllowOrigins 不能是通配符 "*"。必须明确指定来源。
			// 如果调用者没有提供，这里返回一个错误或者一个非常严格的默认值。
			panic("CorsAllowCredentialsWithSpecificOrigins requires at least one allowed origin")
		}
		return cors.New(cors.Config{
			AllowOrigins:     allowedOrigins, // 必须指定来源，不能是 "*"
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true, // 允许凭证
			MaxAge:           12 * time.Hour,
		})

	default:
		// 默认情况下，返回一个保守的配置，或者直接报错
		return cors.New(cors.Config{
			AllowOrigins: []string{"http://localhost:3000"}, // 默认一个常用的开发端口
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
			MaxAge:       5 * time.Minute,
		})
	}
}
