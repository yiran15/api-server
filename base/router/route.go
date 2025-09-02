package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yiran15/api-server/base/middleware"
	"github.com/yiran15/api-server/controller"
	_ "github.com/yiran15/api-server/docs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type RouterInterface interface {
	RegisterRouter(engine *gin.Engine)
}

type Router struct {
	userRouter controller.UserController
	roleRouter controller.RoleController
	apiRouter  controller.ApiController
	middleware middleware.MiddlewareInterface
}

func NewRouter(
	userRouter controller.UserController,
	roleRouter controller.RoleController,
	apiRouter controller.ApiController,
	middleware middleware.MiddlewareInterface) *Router {
	return &Router{
		userRouter: userRouter,
		roleRouter: roleRouter,
		apiRouter:  apiRouter,
		middleware: middleware,
	}
}

func (r *Router) RegisterRouter(engine *gin.Engine) {
	engine.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	engine.Use(ginzap.GinzapWithConfig(zap.L(), &ginzap.Config{
		Context: ginzap.Fn(func(c *gin.Context) []zapcore.Field {
			fields := []zapcore.Field{}
			if requestID := requestid.Get(c); requestID != "" {
				fields = append(fields, zap.String("request-id", requestID))
			}
			return fields
		}),
	}))

	engine.Use(ginzap.RecoveryWithZap(zap.L(), true))
	engine.Use(requestid.New())

	apiGroup := engine.Group("/api/v1")
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.registerOAuthRouter(apiGroup)
	r.registerUserRouter(apiGroup)
	r.registerRoleRouter(apiGroup)
	r.registerApiRouter(apiGroup)
}

func (r *Router) registerUserRouter(apiGroup *gin.RouterGroup) {
	userGroup := apiGroup.Group("/user")
	{
		userGroup.POST("/login", r.userRouter.UserLogin)
		userGroup.Use(r.middleware.Auth())
		userGroup.POST("/logout", r.userRouter.UserLogout)
		userGroup.GET("/info", r.userRouter.UserInfo)
		userGroup.PUT("/self", r.userRouter.UserUpdateBySelf)
		userGroup.Use(r.middleware.AuthZ())
		userGroup.POST("/register", r.userRouter.UserCreate)
		userGroup.PUT("/:id", r.userRouter.UserUpdateByAdmin)
		userGroup.GET("/:id", r.userRouter.UserQuery)
		userGroup.GET("", r.userRouter.UserList)
		userGroup.DELETE("/:id", r.userRouter.UserDelete)
	}
}

func (r *Router) registerRoleRouter(apiGroup *gin.RouterGroup) {
	roleGroup := apiGroup.Group("/role")
	{
		roleGroup.Use(r.middleware.Auth(), r.middleware.AuthZ())
		roleGroup.POST("", r.roleRouter.CreateRole)
		roleGroup.PUT("/:id", r.roleRouter.UpdateRole)
		roleGroup.DELETE("/:id", r.roleRouter.DeleteRole)
		roleGroup.GET("/:id", r.roleRouter.QueryRole)
		roleGroup.GET("", r.roleRouter.ListRole)
	}
}

func (r *Router) registerApiRouter(apiGroup *gin.RouterGroup) {
	baseGroup := apiGroup.Group("/api")
	{
		baseGroup.Use(r.middleware.Auth(), r.middleware.AuthZ())
		baseGroup.GET("/serverApi", r.apiRouter.GetServerApi)
		baseGroup.POST("", r.apiRouter.CreateApi)
		baseGroup.PUT("/:id", r.apiRouter.UpdateApi)
		baseGroup.DELETE("/:id", r.apiRouter.DeleteApi)
		baseGroup.GET("/:id", r.apiRouter.QueryApi)
		baseGroup.GET("", r.apiRouter.ListApi)

	}
}

func (r *Router) registerOAuthRouter(apiGroup *gin.RouterGroup) {
	oauthGroup := apiGroup.Group("/oauth2")
	oauthGroup.Use(r.middleware.Session())
	{
		oauthGroup.GET("/provider", r.userRouter.OAuth2Provider)
		oauthGroup.GET("/login", r.userRouter.OAuthLogin)
		oauthGroup.GET("/callback", r.userRouter.OAuthCallback)
	}
}
