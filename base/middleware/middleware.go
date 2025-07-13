package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yiran15/api-server/pkg/casbin"
	"github.com/yiran15/api-server/pkg/jwt"
)

type MiddlewareInterface interface {
	Auth() gin.HandlerFunc
	Authorization() gin.HandlerFunc
	RequestID() gin.HandlerFunc
	// Logger() gin.HandlerFunc
	Cors(option CorsOption, allowedOrigins ...string) gin.HandlerFunc
}

type Middleware struct {
	jwtImpl   jwt.JwtInterface
	authZImpl casbin.AuthChecker
}

func NewMiddleware(jwtImpl jwt.JwtInterface, authZImpl casbin.AuthChecker) *Middleware {
	return &Middleware{
		jwtImpl:   jwtImpl,
		authZImpl: authZImpl,
	}
}
