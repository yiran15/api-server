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
	// 定义测试数据
	testRole := "test_role"
	testAPI1 := &model.Api{
		Path:   "/api/v1/resource1",
		Method: "GET",
	}
	testAPI2 := &model.Api{
		Path:   "/api/v1/resource2/*", // 带通配符的路径
		Method: "POST",
	}
	testAPI3 := &model.Api{
		Path:   "/api/v1/resource3",
		Method: "*", // 动作通配符
	}
	// --- 测试案例 1: 成功添加一个新策略 ---
	t.Run("Add new policy successfully", func(t *testing.T) {
		added, err := casbinManager.AddRolePolicy(testRole, testAPI1)
		if err != nil {
			t.Fatalf("AddRolePolicy failed unexpectedly: %v", err)
		}
		if !added {
			t.Fatal("AddRolePolicy returned false, expected true for new policy")
		}

		// 验证策略是否存在
		hasPolicy, err := enforcer.HasPolicy(testRole, testAPI1.Path, testAPI1.Method)
		if err != nil {
			t.Fatalf("HasPolicy failed unexpectedly: %v", err)
		}
		if !hasPolicy {
			t.Errorf("Policy %s, %s, %s was not found in enforcer", testRole, testAPI1.Path, testAPI1.Method)
		}

		// 验证 GetRolePolicies 是否能获取到
		policies, err := casbinManager.GetRolePolicies(testRole)
		if err != nil {
			t.Fatalf("GetRolePolicies failed unexpectedly: %v", err)
		}
		if len(policies) != 1 {
			t.Errorf("Expected 1 policy for role %s, got %d", testRole, len(policies))
		}
		if policies[0].Path != testAPI1.Path || policies[0].Method != testAPI1.Method {
			t.Errorf("Retrieved policy mismatch. Expected %s %s, got %s %s",
				testAPI1.Path, testAPI1.Method, policies[0].Path, policies[0].Method)
		}
	})

	// --- 测试案例 2: 添加一个已存在的策略 ---
	t.Run("Add existing policy (should return false)", func(t *testing.T) {
		added, err := casbinManager.AddRolePolicy(testRole, testAPI1) // 再次添加相同的策略
		if err != nil {
			t.Fatalf("AddRolePolicy failed unexpectedly for existing policy: %v", err)
		}
		if added { // 如果策略已存在，AddPolicy 会返回 false
			t.Fatal("AddRolePolicy returned true, expected false for existing policy")
		}

		// 验证策略仍然存在
		hasPolicy, err := enforcer.HasPolicy(testRole, testAPI1.Path, testAPI1.Method)
		if err != nil {
			t.Fatalf("HasPolicy failed unexpectedly: %v", err)
		}
		if !hasPolicy {
			t.Errorf("Existing policy %s, %s, %s was not found in enforcer after re-adding", testRole, testAPI1.Path, testAPI1.Method)
		}
	})

	// --- 测试案例 3: 添加带通配符的策略 ---
	t.Run("Add policy with wildcards", func(t *testing.T) {
		added, err := casbinManager.AddRolePolicy(testRole, testAPI2)
		if err != nil {
			t.Fatalf("AddRolePolicy failed unexpectedly for wildcard policy: %v", err)
		}
		if !added {
			t.Fatal("AddRolePolicy returned false, expected true for new wildcard policy")
		}

		// 验证策略是否存在
		hasPolicy, err := enforcer.HasPolicy(testRole, testAPI2.Path, testAPI2.Method)
		if err != nil {
			t.Fatalf("HasPolicy failed unexpectedly for wildcard policy: %v", err)
		}
		if !hasPolicy {
			t.Errorf("Wildcard policy %s, %s, %s was not found in enforcer", testRole, testAPI2.Path, testAPI2.Method)
		}

		// 验证 GetRolePolicies 是否能获取到所有策略
		policies, err := casbinManager.GetRolePolicies(testRole)
		if err != nil {
			t.Fatalf("GetRolePolicies failed unexpectedly after adding wildcard policy: %v", err)
		}
		if len(policies) != 2 { // 现在应该有 2 条策略
			t.Errorf("Expected 2 policies for role %s, got %d", testRole, len(policies))
		}
	})

	// --- 测试案例 4: 添加带动作通配符的策略 ---
	t.Run("Add policy with action wildcard", func(t *testing.T) {
		added, err := casbinManager.AddRolePolicy(testRole, testAPI3)
		if err != nil {
			t.Fatalf("AddRolePolicy failed unexpectedly for action wildcard policy: %v", err)
		}
		if !added {
			t.Fatal("AddRolePolicy returned false, expected true for new action wildcard policy")
		}

		// 验证策略是否存在
		hasPolicy, err := enforcer.HasPolicy(testRole, testAPI3.Path, testAPI3.Method)
		if err != nil {
			t.Fatalf("HasPolicy failed unexpectedly for action wildcard policy: %v", err)
		}
		if !hasPolicy {
			t.Errorf("Action wildcard policy %s, %s, %s was not found in enforcer", testRole, testAPI3.Path, testAPI3.Method)
		}

		// 验证 GetRolePolicies 是否能获取到所有策略
		policies, err := casbinManager.GetRolePolicies(testRole)
		if err != nil {
			t.Fatalf("GetRolePolicies failed unexpectedly after adding action wildcard policy: %v", err)
		}
		if len(policies) != 3 { // 现在应该有 3 条策略
			t.Errorf("Expected 3 policies for role %s, got %d", testRole, len(policies))
		}
	})
}
