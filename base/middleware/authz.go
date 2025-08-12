package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yiran15/api-server/base/constant"
	"github.com/yiran15/api-server/base/log"
	"github.com/yiran15/api-server/model"
	"github.com/yiran15/api-server/pkg/jwt"
	"github.com/yiran15/api-server/store"
	"go.uber.org/zap"
)

// func (m *Middleware) AuthZ() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var (
// 			claims *jwt.JwtClaims
// 			err    error
// 			roles  []string
// 			// roleNames []any
// 			allow bool
// 		)
// 		if claims, err = m.jwtImpl.GetUser(c.Request.Context()); err != nil {
// 			log.WithRequestID(c.Request.Context()).Error("get jwt claims by ctx failed")
// 			m.Abort(c, http.StatusForbidden, constant.ErrNoPermission)
// 			return
// 		}

// 		roles, err = m.cacheImpl.GetSet(c.Request.Context(), store.RoleType, claims.UserName)
// 		if err != nil {
// 			log.WithRequestID(c.Request.Context()).Error("get role cache failed", zap.Error(err))
// 			m.Abort(c, http.StatusForbidden, constant.ErrNoPermission)
// 			return
// 		}
// 		if len(roles) == 0 {
// 			// user, err := m.userStore.Query(c.Request.Context(), store.Where("id", claims.UserID), store.Preload(model.PreloadRoles))
// 			// if err != nil {
// 			// 	log.WithRequestID(c.Request.Context()).Error("get user by id failed", zap.Error(err))
// 			// 	m.Abort(c, http.StatusForbidden, constant.ErrNoPermission)
// 			// 	return
// 			// }
// 			// if len(user.Roles) != 0 {
// 			// 	roles = make([]string, len(user.Roles))
// 			// 	for i := range user.Roles {
// 			// 		roles[i] = user.Roles[i].Name
// 			// 		roleNames = append(roleNames, user.Roles[i].Name)
// 			// 	}
// 			// }

// 			// if err := m.cacheImpl.SetSet(c.Request.Context(), store.RoleType, claims.UserName, roleNames, nil); err != nil {
// 			// 	log.WithRequestID(c.Request.Context()).Error("set role cache failed", zap.Error(err))
// 			// }
// 			log.WithRequestID(c.Request.Context()).Error("get role cache failed", zap.Error(err))
// 			m.Abort(c, http.StatusForbidden, constant.ErrNoPermission)
// 			return
// 		}

// 		if len(roles) == 0 {
// 			log.WithRequestID(c.Request.Context()).Error("redis cache no role found", zap.String("userName", claims.UserName))
// 			m.Abort(c, http.StatusForbidden, constant.ErrNoPermission)
// 			return
// 		}

// 		for _, roleName := range roles {
// 			if allow, err = m.authZImpl.Enforce(roleName, c.Request.URL.Path, c.Request.Method); err != nil {
// 				m.Abort(c, http.StatusForbidden, err)
// 				return
// 			}
// 			if allow {
// 				break
// 			}
// 		}

// 		if !allow {
// 			m.Abort(c, http.StatusForbidden, constant.ErrNoPermission)
// 			return
// 		}

// 		c.Next()
// 	}
// }

func (m *Middleware) AuthZ() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := m.getClaimsFromCtx(c)
		if err != nil {
			m.Abort(c, http.StatusForbidden, constant.ErrNoPermission)
			return
		}

		roles, err := m.getRolesByUser(c, claims)
		if err != nil || len(roles) == 0 {
			m.Abort(c, http.StatusForbidden, constant.ErrNoPermission)
			return
		}

		if !m.checkPermission(c.Request.Context(), roles, c.Request.URL.Path, c.Request.Method) {
			m.Abort(c, http.StatusForbidden, constant.ErrNoPermission)
			return
		}

		c.Next()
	}
}

// 从上下文获取 JWT claims
func (m *Middleware) getClaimsFromCtx(c *gin.Context) (*jwt.JwtClaims, error) {
	claims, err := m.jwtImpl.GetUser(c.Request.Context())
	if err != nil {
		log.WithRequestID(c.Request.Context()).Error("authz get jwt claims failed", zap.Error(err))
		return nil, err
	}
	return claims, nil
}

// 获取用户角色（缓存优先，缓存 miss 则查询 DB 并回填缓存）
func (m *Middleware) getRolesByUser(c *gin.Context, claims *jwt.JwtClaims) ([]string, error) {
	ctx := c.Request.Context()

	roles, err := m.cacheImpl.GetSet(ctx, store.RoleType, claims.ID)
	if err != nil {
		log.WithRequestID(ctx).Error("authz get role cache failed", zap.Error(err))
		return nil, err
	}

	if len(roles) > 0 {
		return roles, nil
	}

	user, err := m.userStore.Query(ctx, store.Where("id", claims.UserID), store.Preload(model.PreloadRoles))
	if err != nil {
		log.WithRequestID(ctx).Error("authz get user by id failed", zap.Error(err))
		return nil, err
	}

	if len(user.Roles) == 0 {
		return nil, nil
	}

	roles = make([]string, len(user.Roles))
	roleNames := make([]any, len(user.Roles))
	for i, r := range user.Roles {
		roles[i] = r.Name
		roleNames[i] = r.Name
	}

	if err := m.cacheImpl.SetSet(ctx, store.RoleType, claims.ID, roleNames, nil); err != nil {
		log.WithRequestID(ctx).Error("authz set role cache failed", zap.Error(err))
	}

	return roles, nil
}

// 权限校验
func (m *Middleware) checkPermission(ctx context.Context, roles []string, path, method string) bool {
	for _, role := range roles {
		allow, err := m.authZImpl.Enforce(role, path, method)
		if err != nil {
			log.WithRequestID(ctx).Error("authz enforce failed", zap.Error(err), zap.String("role", role), zap.String("path", path), zap.String("method", method))
			return false
		}
		if allow {
			return true
		}
	}
	return false
}
