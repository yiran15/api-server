package casbin

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/yiran15/api-server/model"
)

// AuthChecker 授权检查接口
type AuthChecker interface {
	Enforce(sub, obj, act string) (bool, error)
}

// CasbinManager 策略和角色管理接口
type CasbinManager interface {
	// 角色 CRUD
	AddRolePolicy(role string, api *model.Api) (bool, error)
	GetRolePolicies(role string) ([]*model.Api, error)
	UpdateRolePolicy(role string, oldApi *model.Api, newApi *model.Api) (bool, error) // 修正：更新时只需要role, 不需要oldRole
	DeleteRolePolicy(role string, api *model.Api) (bool, error)
	DeleteAllRolePolicies(role string) (bool, error)

	// 角色用户 CRUD
	AddUserToRole(user, role string) (bool, error)
	DeleteUserFromRole(user, role string) (bool, error)
	GetUsersInRole(role string) ([]string, error)
	GetRolesForUser(user string) ([]string, error)
	DeleteUserAllRoles(user string) (bool, error)
}

// casbinManager 实现结构体
type casbinManager struct {
	enforcer *casbin.Enforcer
}

// NewCasbinManager 创建 CasbinManager 实例
func NewCasbinManager(enforcer *casbin.Enforcer) CasbinManager {
	return &casbinManager{
		enforcer: enforcer,
	}
}

// NewAuthChecker 创建 AuthChecker 实例
func NewAuthChecker(enforcer *casbin.Enforcer) AuthChecker {
	return &casbinManager{ // casbinManager 结构体同时实现了 AuthChecker 和 CasbinManager 接口
		enforcer: enforcer,
	}
}

// Enforce 实现 AuthChecker 接口的授权检查方法
func (m *casbinManager) Enforce(sub, obj, act string) (bool, error) {
	ok, err := m.enforcer.Enforce(sub, obj, act)
	if err != nil {
		return false, fmt.Errorf("casbin enforce failed: %w", err)
	}
	return ok, nil
}

// --- CasbinManager 接口方法的具体实现 ---

// AddRolePolicy 为指定角色添加一个 API 权限策略
func (m *casbinManager) AddRolePolicy(role string, api *model.Api) (bool, error) {
	ok, err := m.enforcer.AddPolicy(role, api.Path, api.Method)
	if err != nil {
		return false, fmt.Errorf("failed to add policy for role %s, api %s %s: %w", role, api.Method, api.Path, err)
	}
	return ok, nil
}

// GetRolePolicies 获取指定角色的所有 API 权限策略
func (m *casbinManager) GetRolePolicies(role string) ([]*model.Api, error) {
	policies, err := m.enforcer.GetFilteredPolicy(0, role)
	if err != nil {
		return nil, fmt.Errorf("failed to get policies for role %s: %w", role, err)
	}

	var apis []*model.Api
	for _, p := range policies {
		if len(p) >= 3 {
			apis = append(apis, &model.Api{
				Path:   p[1],
				Method: p[2],
			})
		}
	}
	return apis, nil
}

// UpdateRolePolicy 更新一个角色的特定 API 权限策略
// 修正：UpdateRolePolicy 接收 role 而不是 oldRole，因为我们是在更新某个角色的策略。
// 同时，旧策略的删除和新策略的添加都围绕这个 role。
func (m *casbinManager) UpdateRolePolicy(role string, oldApi *model.Api, newApi *model.Api) (bool, error) {
	deleted, err := m.enforcer.RemovePolicy(role, oldApi.Path, oldApi.Method)
	if err != nil {
		return false, fmt.Errorf("failed to remove old policy for role %s, api %s %s: %w", role, oldApi.Method, oldApi.Path, err)
	}
	if !deleted {
		// 如果旧策略不存在，则直接添加新策略，并返回 false 表示未删除任何策略
		// 但这里我们认为如果旧策略不存在就不是一个真正的“更新”操作
		_, err = m.enforcer.AddPolicy(role, newApi.Path, newApi.Method)
		if err != nil {
			return false, fmt.Errorf("failed to add new policy after old not found for role %s, api %s %s: %w", role, newApi.Method, newApi.Path, err)
		}
		return false, nil // 表示没有旧策略被删除
	}

	added, err := m.enforcer.AddPolicy(role, newApi.Path, newApi.Method)
	if err != nil {
		return false, fmt.Errorf("failed to add new policy for role %s, api %s %s: %w", role, newApi.Method, newApi.Path, err)
	}
	return added, nil
}

// DeleteRolePolicy 删除指定角色的一个 API 权限策略
func (m *casbinManager) DeleteRolePolicy(role string, api *model.Api) (bool, error) {
	ok, err := m.enforcer.RemovePolicy(role, api.Path, api.Method)
	if err != nil {
		return false, fmt.Errorf("failed to delete policy for role %s, api %s %s: %w", role, api.Method, api.Path, err)
	}
	return ok, nil
}

// DeleteAllRolePolicies 删除一个角色的所有权限策略
func (m *casbinManager) DeleteAllRolePolicies(role string) (bool, error) {
	ok, err := m.enforcer.RemoveFilteredPolicy(0, role)
	if err != nil {
		return false, fmt.Errorf("failed to delete all policies for role %s: %w", role, err)
	}
	return ok, nil
}

// AddUserToRole 将用户添加到角色
func (m *casbinManager) AddUserToRole(user, role string) (bool, error) {
	ok, err := m.enforcer.AddGroupingPolicy(user, role)
	if err != nil {
		return false, fmt.Errorf("failed to add user %s to role %s: %w", user, role, err)
	}
	return ok, nil
}

// DeleteUserFromRole 将用户从角色中移除
func (m *casbinManager) DeleteUserFromRole(user, role string) (bool, error) {
	ok, err := m.enforcer.RemoveGroupingPolicy(user, role)
	if err != nil {
		return false, fmt.Errorf("failed to delete user %s from role %s: %w", user, role, err)
	}
	return ok, nil
}

// GetUsersInRole 获取某个角色的所有用户
func (m *casbinManager) GetUsersInRole(role string) ([]string, error) {
	users, err := m.enforcer.GetUsersForRole(role)
	if err != nil {
		return nil, fmt.Errorf("failed to get users for role %s: %w", role, err)
	}
	return users, nil
}

// GetRolesForUser 获取某个用户拥有的所有角色
func (m *casbinManager) GetRolesForUser(user string) ([]string, error) {
	roles, err := m.enforcer.GetRolesForUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles for user %s: %w", user, err)
	}
	return roles, nil
}

// DeleteUserAllRoles 删除用户拥有的所有角色
func (m *casbinManager) DeleteUserAllRoles(user string) (bool, error) {
	ok, err := m.enforcer.RemoveFilteredGroupingPolicy(0, user)
	if err != nil {
		return false, fmt.Errorf("failed to delete all roles for user %s: %w", user, err)
	}
	return ok, nil
}
