package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yiran15/api-server/base/middleware"
	"github.com/yiran15/api-server/controller"
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
	engine.Use(r.middleware.ZapLogger(), r.middleware.Cors(middleware.CorsAllowAll), r.middleware.RequestID())
	apiGroup := engine.Group("/api/v1")
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
		baseGroup.POST("", r.apiRouter.CreateApi)
		baseGroup.PUT("/:id", r.apiRouter.UpdateApi)
		baseGroup.DELETE("/:id", r.apiRouter.DeleteApi)
		baseGroup.GET("/:id", r.apiRouter.QueryApi)
		baseGroup.GET("", r.apiRouter.ListApi)
	}
}
