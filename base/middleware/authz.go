package middleware

import (
	"context"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/yiran15/api-server/base/constant"
	"github.com/yiran15/api-server/model"
	"github.com/yiran15/api-server/pkg/jwt"
	"github.com/yiran15/api-server/store"
	"go.uber.org/zap"
)

func (m *Middleware) AuthZ() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := requestid.Get(c)
		claims, err := m.getClaimsFromCtx(c, requestID)
		if err != nil {
			zap.L().Error("get jwt claims failed", zap.String("request-id", requestID), zap.Error(err))
			m.Abort(c, http.StatusForbidden, constant.ErrNoPermission)
			return
		}

		roles, err := m.getRolesByUser(c, claims, requestID)
		if err != nil || len(roles) == 0 {
			if err != nil {
				zap.L().Error("get user roles error", zap.String("request-id", requestID), zap.Error(err))
			}
			if len(roles) == 0 {
				zap.L().Error("user has no roles", zap.String("request-id", requestID), zap.String("userName", claims.UserName))
			}
			m.Abort(c, http.StatusForbidden, constant.ErrNoPermission)
			return
		}

		if !m.checkPermission(c.Request.Context(), roles, c.Request.URL.Path, c.Request.Method, requestID) {
			zap.L().Error("user has no permission", zap.String("request-id", requestID), zap.String("userName", claims.UserName), zap.Strings("roles", roles), zap.String("path", c.Request.URL.Path), zap.String("method", c.Request.Method))
			m.Abort(c, http.StatusForbidden, constant.ErrNoPermission)
			return
		}

		c.Next()
	}
}

// 从上下文获取 JWT claims
func (m *Middleware) getClaimsFromCtx(c *gin.Context, requestID string) (*jwt.JwtClaims, error) {
	claims, err := m.jwtImpl.GetUser(c.Request.Context())
	if err != nil {
		zap.L().Error("get jwt claims failed", zap.String("request-id", requestID), zap.Error(err))
		return nil, err
	}
	return claims, nil
}

// 获取用户角色（缓存优先，缓存 miss 则查询 DB 并回填缓存）
func (m *Middleware) getRolesByUser(c *gin.Context, claims *jwt.JwtClaims, requestID string) ([]string, error) {
	ctx := c.Request.Context()

	roles, err := m.cacheImpl.GetSet(ctx, store.RoleType, claims.UserID)
	if err != nil {
		zap.L().Error("authz get role cache failed", zap.String("request-id", requestID), zap.Error(err))
		return nil, err
	}

	if len(roles) > 0 {
		if len(roles) == 1 && roles[0] == constant.EmptyRoleSentinel {
			return []string{}, nil
		}
		return roles, nil
	}

	user, err := m.userStore.Query(ctx, store.Where("id", claims.UserID), store.Preload(model.PreloadRoles))
	if err != nil {
		zap.L().Error("authz get user by id failed", zap.String("request-id", requestID), zap.Error(err))
		return nil, err
	}

	if len(user.Roles) == 0 {
		// 缓存哨兵值，标记无角色
		if err := m.cacheImpl.SetSet(ctx, store.RoleType, claims.UserID, []any{constant.EmptyRoleSentinel}, nil); err != nil {
			zap.L().Error("authz set empty role cache failed", zap.String("request-id", requestID), zap.Error(err))
		}
		return []string{}, nil
	}

	roles = make([]string, len(user.Roles))
	roleNames := make([]any, len(user.Roles))
	for i, r := range user.Roles {
		roles[i] = r.Name
		roleNames[i] = r.Name
	}

	if err := m.cacheImpl.SetSet(ctx, store.RoleType, claims.UserID, roleNames, nil); err != nil {
		zap.L().Error("authz set role cache failed", zap.String("request-id", requestID), zap.Error(err))
		return nil, err
	}

	return roles, nil
}

// 权限校验
func (m *Middleware) checkPermission(_ context.Context, roles []string, path, method, requestID string) bool {
	for _, role := range roles {
		allow, err := m.authZImpl.Enforce(role, path, method)
		if err != nil {
			zap.L().Error("authz enforce failed", zap.String("request-id", requestID), zap.Error(err), zap.String("role", role), zap.String("path", path), zap.String("method", method))
			return false
		}
		if allow {
			return true
		}
	}
	return false
}
