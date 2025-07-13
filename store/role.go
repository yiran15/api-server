package store

import (
	"context"

	"github.com/yiran15/api-server/model"
)

type RoleStorer interface {
	Create(ctx context.Context, obj *model.Role) error
	CreateBatch(ctx context.Context, objs []*model.Role) error // 批量创建
	Update(ctx context.Context, obj *model.Role, opts ...Option) error
	Delete(ctx context.Context, obj *model.Role, opts ...Option) error // 增加选项，支持where条件删除
	Query(ctx context.Context, opts ...Option) (*model.Role, error)
	List(ctx context.Context, pagination *Pagination, opts ...Option) (total int64, objs []*model.Role, err error)
}

func NewRoleStore(dbProvider DBProviderInterface) RoleStorer {
	return NewRepository[model.Role](dbProvider)
}
