package casbin_test

import (
	"testing"

	casbinv2 "github.com/casbin/casbin/v2"
	"github.com/yiran15/api-server/base/conf"
	"github.com/yiran15/api-server/base/data"
	"github.com/yiran15/api-server/model"
	"github.com/yiran15/api-server/pkg/casbin"
)

var (
	casbinManager casbin.CasbinManager
	enforcer      *casbinv2.Enforcer
	// 定义测试数据
	testRole = "test_role"
	testUser = "test_user"
	testApis = []*model.Api{
		{
			Id:          1,
			Name:        "resource1",
			Description: "resource1 GET",
			Path:        "/api/v1/resource1",
			Method:      "GET",
		},
		{
			Id:          2,
			Name:        "resource2",
			Description: "resource2 POST",
			Path:        "/api/v1/resource2/*", // 带通配符的路径
			Method:      "POST",
		},
		{
			Id:          3,
			Name:        "resource3",
			Description: "resource3 *",
			Path:        "/api/v1/resource3",
			Method:      "*", // 动作通配符
		},
	}
)

func init() {
	conf.LoadConfig("../../config.yaml")
	db, _, err := data.NewDB()
	if err != nil {
		panic(err)
	}
	enforcer, err = casbin.NewEnforcer(db)
	if err != nil {
		panic(err)
	}
	casbinManager = casbin.NewCasbinManager(enforcer)
}

func TestAddRole(t *testing.T) {
	for _, api := range testApis {
		_, err := casbinManager.AddRolePolicy(testRole, api)
		if err != nil {
			t.Fatalf("AddRolePolicy failed unexpectedly: %v", err)
		}
	}
}

func TestGetRole(t *testing.T) {
	apis, err := casbinManager.GetRolePolicies(testRole)
	if err != nil {
		t.Fatalf("GetRolePolicies failed unexpectedly: %v", err)
	}
	if len(apis) != len(testApis) {
		t.Fatalf("GetRolePolicies returned unexpected number of policies: expected %d, got %d", len(testApis), len(apis))
	}
}

func TestDeleteRolePolicy(t *testing.T) {
	for _, api := range testApis {
		_, err := casbinManager.DeleteRolePolicy(testRole, api)
		if err != nil {
			t.Fatalf("DeleteRolePolicy failed unexpectedly: %v", err)
		}
	}
}

func TestAddUserToRole(t *testing.T) {
	_, err := casbinManager.AddUserToRole(testUser, testRole)
	if err != nil {
		t.Fatalf("AddUserToRole failed unexpectedly: %v", err)
	}
}

func TestDeleteUserFromRole(t *testing.T) {
	_, err := casbinManager.DeleteUserFromRole(testUser, testRole)
	if err != nil {
		t.Fatalf("DeleteUserFromRole failed unexpectedly: %v", err)
	}
}

func TestGetRolesForUser(t *testing.T) {
	roles, err := casbinManager.GetRolesForUser(testUser)
	if err != nil {
		t.Fatalf("GetRolesForUser failed unexpectedly: %v", err)
	}
	if len(roles) != 1 {
		t.Fatalf("GetRolesForUser returned unexpected number of roles: expected 1, got %d", len(roles))
	}
}
