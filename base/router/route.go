package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yiran15/api-server/base/middleware"
	"github.com/yiran15/api-server/controller"
)

type RouterInterface interface {
	RegisterRouter(apiGroup *gin.RouterGroup)
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

func (r *Router) RegisterRouter(apiGroup *gin.RouterGroup) {
	apiGroup.Use(r.middleware.ZapLogger(), r.middleware.RequestID())
	v1Group := apiGroup.Group("/v1")

	r.registerUserRouter(v1Group)
	r.registerRoleRouter(v1Group)
	r.registerApiRouter(v1Group)
}

func (r *Router) registerUserRouter(userGroup *gin.RouterGroup) {
	baseGroup := userGroup.Group("/users")
	baseGroup.POST("/login", r.userRouter.UserLogin)
	baseGroup.POST("/register", r.userRouter.UserCreate)
	aGroup := baseGroup.Use(r.middleware.Auth())
	aGroup.GET("/info", r.userRouter.UserInfo)
	aGroup.PUT("/:id", r.userRouter.UserUpdate)

	authGroup := baseGroup.Use(r.middleware.Auth(), r.middleware.AuthZ())
	authGroup.POST("/logout", r.userRouter.UserLogout)

	authGroup.GET("/:id", r.userRouter.UserQuery)
	authGroup.GET("/", r.userRouter.UserList)
	authGroup.DELETE("/:id", r.userRouter.UserDelete)
	authGroup.PUT("/:id/roles", r.userRouter.UserUpdateRole)
}

func (r *Router) registerRoleRouter(roleGroup *gin.RouterGroup) {
	baseGroup := roleGroup.Group("/roles")
	authGroup := baseGroup.Use(r.middleware.Auth(), r.middleware.AuthZ())
	authGroup.POST("/", r.roleRouter.CreateRole)
	authGroup.PUT("/:id", r.roleRouter.UpdateRole)
	authGroup.DELETE("/:id", r.roleRouter.DeleteRole)
	authGroup.GET("/:id", r.roleRouter.QueryRole)
	authGroup.GET("/", r.roleRouter.ListRole)
}

func (r *Router) registerApiRouter(apiGroup *gin.RouterGroup) {
	baseGroup := apiGroup.Group("/apis")
	authGroup := baseGroup.Use(r.middleware.Auth(), r.middleware.AuthZ())
	authGroup.POST("/", r.apiRouter.CreateApi)
	authGroup.PUT("/:id", r.apiRouter.UpdateApi)
	authGroup.DELETE("/:id", r.apiRouter.DeleteApi)
	authGroup.GET("/:id", r.apiRouter.QueryApi)
	authGroup.GET("/", r.apiRouter.ListApi)
}
