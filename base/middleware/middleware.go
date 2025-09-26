package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yiran15/api-server/base/apitypes"
	"github.com/yiran15/api-server/pkg/casbin"
	"github.com/yiran15/api-server/pkg/jwt"
	"github.com/yiran15/api-server/store"
)

type MiddlewareInterface interface {
	Auth() gin.HandlerFunc
	AuthZ() gin.HandlerFunc
	Session() gin.HandlerFunc
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

func (m *Middleware) Abort(c *gin.Context, code int, err error) {
	c.JSON(code, apitypes.NewResponseWithOpts(code, apitypes.WithError(err.Error())))
	c.Error(err)
	c.Abort()
}
