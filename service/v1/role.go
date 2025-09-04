package v1

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/yiran15/api-server/base/apitypes"
	"github.com/yiran15/api-server/base/helper"
	"github.com/yiran15/api-server/model"
	"github.com/yiran15/api-server/pkg/casbin"
	"github.com/yiran15/api-server/store"
	"gorm.io/gorm"
)

type RoleServicer interface {
	CreateRole(ctx context.Context, req *apitypes.RoleCreateRequest) error
	UpdateRole(ctx context.Context, req *apitypes.RoleUpdateRequest) error
	DeleteRole(ctx context.Context, req *apitypes.IDRequest) error
	QueryRole(ctx context.Context, req *apitypes.IDRequest) (*model.Role, error)
	ListRole(ctx context.Context, pagination *apitypes.RoleListRequest) (*apitypes.RoleListResponse, error)
}

type roleService struct {
	roleRepository store.RoleStorer
	apiRepository  store.ApiStorer
	casbinStore    store.CasbinStorer
	casbinManager  casbin.CasbinManager
	txManager      store.TxManagerInterface
}

func NewRoleService(roleRepository store.RoleStorer, apiRepository store.ApiStorer, casbinStore store.CasbinStorer, casbinManager casbin.CasbinManager, txManager store.TxManagerInterface) RoleServicer {
	return &roleService{
		roleRepository: roleRepository,
		apiRepository:  apiRepository,
		casbinStore:    casbinStore,
		casbinManager:  casbinManager,
		txManager:      txManager,
	}
}

func (receiver *roleService) CreateRole(ctx context.Context, req *apitypes.RoleCreateRequest) error {
	req.Apis = helper.RemoveDuplicates(req.Apis)
	var (
		role  *model.Role
		total int64
		apis  []*model.Api
		err   error
		rules []*model.CasbinRule
	)

	if role, err = receiver.roleRepository.Query(ctx, store.Where("name", req.Name)); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	if role != nil {
		return fmt.Errorf("role %s already exists", req.Name)
	}

	if len(req.Apis) > 0 {
		total, apis, err = receiver.apiRepository.List(ctx, 0, 0, "", "", store.In("id", req.Apis))
		if err != nil {
			return err
		}

		if err := helper.ValidateRoleApis(req.Apis, total, apis); err != nil {
			return err
		}
	}

	for _, api := range apis {
		rules = append(rules, &model.CasbinRule{
			PType: helper.String("p"),
			V0:    helper.String(req.Name),
			V1:    helper.String(api.Path),
			V2:    helper.String(api.Method),
		})
	}

	if err := receiver.txManager.Transaction(ctx, func(ctx context.Context) error {
		if err := receiver.roleRepository.Create(ctx, &model.Role{
			Name:        req.Name,
			Description: req.Description,
			Apis:        apis,
		}); err != nil {
			return err
		}
		if err := receiver.casbinStore.CreateBatch(ctx, rules); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return receiver.casbinManager.LoadPolicy()
}

func (receiver *roleService) UpdateRole(ctx context.Context, req *apitypes.RoleUpdateRequest) error {
	var (
		total int64
		apis  []*model.Api
		err   error
		rules []*model.CasbinRule
	)
	req.Apis = helper.RemoveDuplicates(req.Apis)
	role, err := receiver.roleRepository.Query(ctx, store.Where("id", req.ID))
	if err != nil {
		return err
	}

	role.Description = req.Description
	if len(req.Apis) > 0 {
		total, apis, err = receiver.apiRepository.List(ctx, 0, 0, "", "", store.In("id", req.Apis))
		if err != nil {
			return err
		}

		if err := helper.ValidateRoleApis(req.Apis, total, apis); err != nil {
			return err
		}
	}

	total, casbinRules, err := receiver.casbinStore.List(ctx, 0, 0, "", "", store.Where("v0", role.Name))
	if err != nil {
		return err
	}

	for _, api := range apis {
		rules = append(rules, &model.CasbinRule{
			PType: helper.String("p"),
			V0:    helper.String(role.Name),
			V1:    helper.String(api.Path),
			V2:    helper.String(api.Method),
		})
	}

	if err := receiver.txManager.Transaction(ctx, func(ctx context.Context) error {
		if err := receiver.roleRepository.Update(ctx, role); err != nil {
			return err
		}
		if total > 0 {
			if err := receiver.casbinStore.DeleteBatch(ctx, casbinRules); err != nil {
				return err
			}
		}
		if err := receiver.casbinStore.CreateBatch(ctx, rules); err != nil {
			return err
		}
		if err := receiver.roleRepository.ReplaceAssociation(ctx, role, model.PreloadApis, apis); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return receiver.casbinManager.LoadPolicy()
}

func (receiver *roleService) DeleteRole(ctx context.Context, req *apitypes.IDRequest) error {
	role, err := receiver.roleRepository.Query(ctx, store.Where("id", req.ID), store.Preload(model.PreloadUsers))
	if err != nil {
		return err
	}

	if len(role.Users) > 0 {
		unameArry := make([]string, 0, len(role.Users))
		for _, u := range role.Users {
			unameArry = append(unameArry, u.Name)
		}
		unames := strings.Join(unameArry, ",")
		return fmt.Errorf("the role is being used by the users %s", unames)
	}

	_, casbinRules, err := receiver.casbinStore.List(ctx, 0, 0, "", "", store.Where("v0", role.Name))
	if err != nil {
		return err
	}

	if err := receiver.txManager.Transaction(ctx, func(ctx context.Context) error {
		if err := receiver.roleRepository.Delete(ctx, role); err != nil {
			return err
		}
		if err := receiver.roleRepository.ClearAssociation(ctx, role, model.PreloadApis); err != nil {
			return err
		}
		if err := receiver.casbinStore.DeleteBatch(ctx, casbinRules); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return receiver.casbinManager.LoadPolicy()
}

func (receiver *roleService) QueryRole(ctx context.Context, req *apitypes.IDRequest) (*model.Role, error) {
	return receiver.roleRepository.Query(ctx, store.Where("id", req.ID), store.Preload(model.PreloadApis))
}

func (receiver *roleService) ListRole(ctx context.Context, req *apitypes.RoleListRequest) (*apitypes.RoleListResponse, error) {
	var (
		where store.Option
		colum = "id"
		oder  = "desc"
	)
	if req.Name != "" {
		where = store.Like("name", req.Name+"%")
	}

	if req.Sort != "" && req.Direction != "" {
		colum = req.Sort
		oder = req.Direction
	}

	total, objs, err := receiver.roleRepository.List(ctx, req.Page, req.PageSize, colum, oder, where)
	if err != nil {
		return nil, err
	}
	res := &apitypes.RoleListResponse{
		ListResponse: &apitypes.ListResponse{
			Pagination: &apitypes.Pagination{
				Page:     req.Page,
				PageSize: req.PageSize,
			},
			Total: total,
		},
		List: objs,
	}
	return res, nil
}
