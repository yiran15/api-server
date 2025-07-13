package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yiran15/api-server/base/constant"
	"github.com/yiran15/api-server/base/log"
	"github.com/yiran15/api-server/pkg/jwt"
	"github.com/yiran15/api-server/store"
	"go.uber.org/zap"
)

func (m *Middleware) AuthZ() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			claims *jwt.JwtClaims
			err    error
			allow  bool
		)
		if claims, err = m.jwtImpl.GetUser(c.Request.Context()); err != nil {
			log.WithRequestID(c.Request.Context()).Error("get jwt claims by ctx failed")
			m.Abort(c, http.StatusForbidden, constant.ErrNoPermission)
			return
		}

		roles, err := m.cacheImpl.GetSet(c.Request.Context(), store.RoleType, claims.UserName)
		if err != nil {
			log.WithRequestID(c.Request.Context()).Error("get role cache failed", zap.Error(err))
			m.Abort(c, http.StatusForbidden, constant.ErrNoPermission)
			return
		}
		if len(roles) == 0 {
			log.WithRequestID(c.Request.Context()).Error("redis cache no role found", zap.String("userName", claims.UserName))
			m.Abort(c, http.StatusForbidden, constant.ErrNoPermission)
			return
		}

		for _, roleName := range roles {
			if allow, err = m.authZImpl.Enforce(roleName, c.Request.URL.Path, c.Request.Method); err != nil {
				m.Abort(c, http.StatusForbidden, err)
				return
			}
			if allow {
				break
			}
		}

		if !allow {
			m.Abort(c, http.StatusForbidden, constant.ErrNoPermission)
			return
		}

		c.Next()
	}
}
