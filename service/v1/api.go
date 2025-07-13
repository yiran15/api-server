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

func (a *ApiService) CreateApi(ctx context.Context, req *apitypes.ApiCreateRequest) error {
	log.WithBody(ctx, req).Info("create api request")
	if api, err := a.apiStore.Query(ctx, store.Where("name", req.Name)); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if api != nil {
			return fmt.Errorf("api %s already exists", req.Name)
		}
	}

	return a.apiStore.Create(ctx, &model.Api{
		Name:        req.Name,
		Path:        req.Path,
		Method:      req.Method,
		Description: req.Description,
	})
}

func (a *ApiService) UpdateApi(ctx context.Context, req *apitypes.ApiUpdateRequest) error {
	log.WithBody(ctx, req).Info("update api request")
	api, err := a.apiStore.Query(ctx, store.Where("id", req.ID))
	if err != nil {
		return err
	}
	api.Description = req.Description
	return a.apiStore.Update(ctx, api)
}

func (a *ApiService) DeleteApi(ctx context.Context, req *apitypes.IDRequest) error {
	log.WithBody(ctx, req).Info("delete api request")
	api, err := a.apiStore.Query(ctx, store.Where("id", req.ID), store.Preload(model.PreloadRoles))
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

	return a.apiStore.Delete(ctx, api)
}

func (a *ApiService) QueryApi(ctx context.Context, req *apitypes.IDRequest) (*model.Api, error) {
	log.WithBody(ctx, req).Info("query api request")
	return a.apiStore.Query(ctx, store.Where("id", req.ID))
}

func (a *ApiService) ListApi(ctx context.Context, pagination *apitypes.ApiListRequest) (*apitypes.ApiListResponse, error) {
	log.WithBody(ctx, pagination).Info("list api request")
	var (
		where store.Option
		colum string
		oder  string
	)

	if pagination.Name != "" {
		where = store.Like("name", "%"+pagination.Name+"%")
	} else if pagination.Path != "" {
		where = store.Like("path", "%"+pagination.Path+"%")
	} else if pagination.Method != "" {
		where = store.Like("method", "%"+pagination.Method+"%")
	}

	if pagination.SortParam != nil {
		colum = pagination.SortParam.GetApiSortField(pagination.Sort)
		oder = pagination.SortParam.Direction
	}
	total, apis, err := a.apiStore.List(ctx, pagination.Page, pagination.PageSize, colum, oder, where)
	if err != nil {
		return nil, err
	}
	return &apitypes.ApiListResponse{
		ListResponse: &apitypes.ListResponse{
			Pagination: &apitypes.Pagination{
				Page:     pagination.Page,
				PageSize: pagination.PageSize,
			},
			Total: total,
		},
		List: apis,
	}, nil
}
