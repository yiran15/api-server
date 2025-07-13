package store

import (
	"context"

	"github.com/yiran15/api-server/model"
)

type UserStorer interface {
	Create(ctx context.Context, obj *model.User) error
	CreateBatch(ctx context.Context, objs []*model.User) error // 批量创建
	Update(ctx context.Context, obj *model.User, opts ...Option) error
	Delete(ctx context.Context, obj *model.User, opts ...Option) error // 增加选项，支持where条件删除
	Query(ctx context.Context, opts ...Option) (*model.User, error)
	List(ctx context.Context, pagination *Pagination, opts ...Option) (total int64, objs []*model.User, err error)
}

func NewUserStore(dbProvider DBProviderInterface) UserStorer {
	return NewRepository[model.User](dbProvider)
}
