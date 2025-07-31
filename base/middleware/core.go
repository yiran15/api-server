package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 导入 Gin CORS 中间件

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
			AllowAllOrigins: true,                                                                              // 允许所有来源
			AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},                      // 允许所有常用方法
			AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}, // 允许所有常用请求头
			ExposeHeaders:   []string{"Content-Length"},                                                        // 允许客户端访问的响应头
			MaxAge:          12 * time.Hour,                                                                    // 预检请求的缓存时间
		})

	case CorsSpecificOrigins:
		if len(allowedOrigins) == 0 {
			return cors.New(cors.Config{
				AllowOrigins:  []string{"http://localhost"}, // 默认一个本地开发地址
				AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowHeaders:  []string{"Origin", "Content-Type", "Accept", "Authorization"},
				ExposeHeaders: []string{"Content-Length"},
				MaxAge:        12 * time.Hour,
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
			panic("CorsAllowCredentialsWithSpecificOrigins requires at least one allowed origin")
		}
		return cors.New(cors.Config{
			AllowOrigins:  allowedOrigins, // 必须指定来源，不能是 "*"
			AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:  []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
			ExposeHeaders: []string{"Content-Length"},
			MaxAge:        12 * time.Hour,
		})

	default:
		return cors.New(cors.Config{
			AllowOrigins: []string{"http://localhost:3000"}, // 默认一个常用的开发端口
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
			MaxAge:       5 * time.Minute,
		})
	}
}

// CorssDomainMiddleware 跨域中间件
// func (m *Middleware) Cors(_ CorsOption, _ ...string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if origin := c.Request.Header.Get("Origin"); origin != "" {
// 			c.Header("Access-Control-Allow-Origin", "*")
// 			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
// 			c.Header("Access-Control-Max-Age", "7200")
// 		}

// 		//allows OPTIONS method
// 		if c.Request.Method == http.MethodOptions {
// 			c.AbortWithStatus(http.StatusOK)
// 			return
// 		}
// 		c.Next()
// 	}
// }
