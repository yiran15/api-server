package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yiran15/api-server/pkg/casbin"
	"github.com/yiran15/api-server/pkg/jwt"
	"github.com/yiran15/api-server/store"
)

type MiddlewareInterface interface {
	Auth() gin.HandlerFunc
	AuthZ() gin.HandlerFunc
	RequestID() gin.HandlerFunc
	ZapLogger() gin.HandlerFunc
	Cors(option CorsOption, allowedOrigins ...string) gin.HandlerFunc
}

type Middleware struct {
	jwtImpl   jwt.JwtInterface
	authZImpl casbin.AuthChecker
	cacheImpl store.CacheStorer
	userStore store.UserStorer
}

func NewMiddleware(jwtImpl jwt.JwtInterface, authZImpl casbin.AuthChecker, cacheImpl store.CacheStorer, userStore store.UserStorer) *Middleware {
	return &Middleware{
		jwtImpl:   jwtImpl,
		authZImpl: authZImpl,
		cacheImpl: cacheImpl,
		userStore: userStore,
	}
}
