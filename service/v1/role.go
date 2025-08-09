package v1

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/yiran15/api-server/base/apitypes"
	"github.com/yiran15/api-server/base/helper"
	"github.com/yiran15/api-server/base/log"
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

func (s *roleService) CreateRole(ctx context.Context, req *apitypes.RoleCreateRequest) error {
	log.WithBody(ctx, req).Info("create role request")
	req.Apis = helper.RemoveDuplicates(req.Apis)
	var (
		role  *model.Role
		total int64
		apis  []*model.Api
		err   error
		rules []*model.CasbinRule
	)

	if role, err = s.roleRepository.Query(ctx, store.Where("name", req.Name)); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	if role != nil {
		return fmt.Errorf("role %s already exists", req.Name)
	}

	if len(req.Apis) > 0 {
		total, apis, err = s.apiRepository.List(ctx, 0, 0, "", "", store.In("id", req.Apis))
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

	if err := s.txManager.Transaction(ctx, func(ctx context.Context) error {
		if err := s.roleRepository.Create(ctx, &model.Role{
			Name:        req.Name,
			Description: req.Description,
			Apis:        apis,
		}); err != nil {
			return err
		}
		if err := s.casbinStore.CreateBatch(ctx, rules); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return s.casbinManager.LoadPolicy()
}

func (s *roleService) UpdateRole(ctx context.Context, req *apitypes.RoleUpdateRequest) error {
	log.WithBody(ctx, req).Info("update role request")
	var (
		total int64
		apis  []*model.Api
		err   error
		rules []*model.CasbinRule
	)
	req.Apis = helper.RemoveDuplicates(req.Apis)
	role, err := s.roleRepository.Query(ctx, store.Where("id", req.ID))
	if err != nil {
		return err
	}

	role.Description = req.Description
	if len(req.Apis) > 0 {
		total, apis, err = s.apiRepository.List(ctx, 0, 0, "", "", store.In("id", req.Apis))
		if err != nil {
			return err
		}

		if err := helper.ValidateRoleApis(req.Apis, total, apis); err != nil {
			return err
		}
	}

	total, casbinRules, err := s.casbinStore.List(ctx, 0, 0, "", "", store.Where("v0", role.Name))
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

	if err := s.txManager.Transaction(ctx, func(ctx context.Context) error {
		if err := s.roleRepository.Update(ctx, role); err != nil {
			return err
		}
		if total > 0 {
			if err := s.casbinStore.DeleteBatch(ctx, casbinRules); err != nil {
				return err
			}
		}
		if err := s.casbinStore.CreateBatch(ctx, rules); err != nil {
			return err
		}
		if err := s.roleRepository.ReplaceAssociation(ctx, role, model.PreloadApis, apis); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return s.casbinManager.LoadPolicy()
}

func (s *roleService) DeleteRole(ctx context.Context, req *apitypes.IDRequest) error {
	log.WithBody(ctx, req).Info("delete role request")
	role, err := s.roleRepository.Query(ctx, store.Where("id", req.ID), store.Preload(model.PreloadUsers))
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

	_, casbinRules, err := s.casbinStore.List(ctx, 0, 0, "", "", store.Where("v0", role.Name))
	if err != nil {
		return err
	}

	if err := s.txManager.Transaction(ctx, func(ctx context.Context) error {
		if err := s.roleRepository.Delete(ctx, role); err != nil {
			return err
		}
		if err := s.roleRepository.ClearAssociation(ctx, role, model.PreloadApis); err != nil {
			return err
		}
		if err := s.casbinStore.DeleteBatch(ctx, casbinRules); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return s.casbinManager.LoadPolicy()
}

func (s *roleService) QueryRole(ctx context.Context, req *apitypes.IDRequest) (*model.Role, error) {
	log.WithBody(ctx, req).Info("query role request")
	return s.roleRepository.Query(ctx, store.Where("id", req.ID), store.Preload(model.PreloadApis))
}

func (s *roleService) ListRole(ctx context.Context, req *apitypes.RoleListRequest) (*apitypes.RoleListResponse, error) {
	log.WithBody(ctx, req).Info("list role request")
	var (
		where store.Option
		colum string
		oder  string
	)
	if req.Name != "" {
		where = store.Like("name", req.Name+"%")
	}

	if req.Sort != "" && req.Direction != "" {
		colum = req.Sort
		oder = req.Direction
	}

	total, objs, err := s.roleRepository.List(ctx, req.Page, req.PageSize, colum, oder, where)
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
