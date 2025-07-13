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
	middleware middleware.MiddlewareInterface
}

func NewRouter(
	userRouter controller.UserController,
	middleware middleware.MiddlewareInterface) *Router {
	return &Router{
		userRouter: userRouter,
		middleware: middleware,
	}
}

func (r *Router) RegisterRouter(apiGroup *gin.RouterGroup) {
	// TODO 添加路由
}
