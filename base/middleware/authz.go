package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/yiran15/api-server/pkg/jwt"
)

func (m *Middleware) Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			claims *jwt.JwtClaims
			err    error
			allow  bool
		)
		if claims, err = jwt.GetJwtClaimsByCtx(c); err != nil {
			_ = c.Error(errors.New("no permission"))
			c.Abort()
			return
		}

		// 校验权限
		if allow, err = m.authZImpl.Enforce(claims.UserName, c.Request.URL.Path, c.Request.Method); err != nil {
			_ = c.Error(errors.New("no permission"))
			c.Abort()
			return
		}
		if !allow {
			_ = c.Error(errors.New("no permission"))
			c.Abort()
			return
		}

		c.Next()
	}
}
