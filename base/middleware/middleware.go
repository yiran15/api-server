package middleware

import (
	"net/http"

	"github.com/gin-contrib/requestid"
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
	switch code {
	case http.StatusUnauthorized:
		c.JSON(code, apitypes.NewResponse(code, "", requestid.Get(c), nil, err.Error()))
	case http.StatusForbidden:
		c.JSON(code, apitypes.NewResponse(code, "", requestid.Get(c), nil, err.Error()))
	}
	c.Error(err)
	c.Abort()
}
