package v1

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/yiran15/api-server/base/apitypes"
	"github.com/yiran15/api-server/base/log"
	"github.com/yiran15/api-server/model"
	"github.com/yiran15/api-server/store"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ApiServicer interface {
	CreateApi(ctx context.Context, req *apitypes.ApiCreateRequest) error
	UpdateApi(ctx context.Context, req *apitypes.ApiUpdateRequest) error
	DeleteApi(ctx context.Context, req *apitypes.IDRequest) error
	QueryApi(ctx context.Context, req *apitypes.IDRequest) (*model.Api, error)
	ListApi(ctx context.Context, pagination *apitypes.ApiListRequest) (*apitypes.ApiListResponse, error)
}

type ApiService struct {
	apiStore store.ApiStorer
}

func NewApiServicer(apiStore store.ApiStorer) ApiServicer {
	return &ApiService{
		apiStore: apiStore,
	}
}

func (receiver *ApiService) CreateApi(ctx context.Context, req *apitypes.ApiCreateRequest) error {
	if api, err := receiver.apiStore.Query(ctx, store.Where("name", req.Name)); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if api != nil {
			return fmt.Errorf("api %s already exists", req.Name)
		}
	}

	return receiver.apiStore.Create(ctx, &model.Api{
		Name:        req.Name,
		Path:        req.Path,
		Method:      req.Method,
		Description: req.Description,
	})
}

func (receiver *ApiService) UpdateApi(ctx context.Context, req *apitypes.ApiUpdateRequest) error {
	api, err := receiver.apiStore.Query(ctx, store.Where("id", req.ID))
	if err != nil {
		return err
	}
	api.Description = req.Description
	return receiver.apiStore.Update(ctx, api)
}

func (receiver *ApiService) DeleteApi(ctx context.Context, req *apitypes.IDRequest) error {
	api, err := receiver.apiStore.Query(ctx, store.Where("id", req.ID), store.Preload(model.PreloadRoles))
	if err != nil {
		return err
	}

	if len(api.Roles) > 0 {
		roles := make([]string, 0, len(api.Roles))
		for _, role := range api.Roles {
			roles = append(roles, role.Name)
		}
		rolesName := strings.Join(roles, ",")
		log.WithRequestID(ctx).Error("api has roles", zap.String("apiName", api.Name), zap.String("rolesName", rolesName))
		return fmt.Errorf("api %s has roles %s", api.Name, rolesName)
	}

	return receiver.apiStore.Delete(ctx, api)
}

func (receiver *ApiService) QueryApi(ctx context.Context, req *apitypes.IDRequest) (*model.Api, error) {
	return receiver.apiStore.Query(ctx, store.Where("id", req.ID))
}

func (receiver *ApiService) ListApi(ctx context.Context, req *apitypes.ApiListRequest) (*apitypes.ApiListResponse, error) {
	var (
		where store.Option
		colum = "id"
		oder  = "desc"
	)

	if req.Name != "" {
		where = store.Like("name", req.Name+"%")
	} else if req.Path != "" {
		where = store.Like("path", req.Path+"%")
	} else if req.Method != "" {
		where = store.Like("method", req.Method+"%")
	}

	if req.Sort != "" && req.Direction != "" {
		colum = req.Sort
		oder = req.Direction
	}

	total, apis, err := receiver.apiStore.List(ctx, req.Page, req.PageSize, colum, oder, where)
	if err != nil {
		return nil, err
	}
	return &apitypes.ApiListResponse{
		ListResponse: &apitypes.ListResponse{
			Pagination: &apitypes.Pagination{
				Page:     req.Page,
				PageSize: req.PageSize,
			},
			Total: total,
		},
		List: apis,
	}, nil
}
