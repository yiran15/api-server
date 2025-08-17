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
	List(ctx context.Context, page, pageSize int, colum, oder string, opts ...Option) (total int64, objs []*model.User, err error)
	AppendAssociation(ctx context.Context, model *model.User, objName string, obj any) error
	ReplaceAssociation(ctx context.Context, model *model.User, objName string, obj any) error
	ClearAssociation(ctx context.Context, model *model.User, objName string) error
}

func NewUserStore(dbProvider DBProviderInterface) UserStorer {
	return NewRepository[model.User](dbProvider)
}

type RoleStorer interface {
	Create(ctx context.Context, obj *model.Role) error
	CreateBatch(ctx context.Context, objs []*model.Role) error // 批量创建
	Update(ctx context.Context, obj *model.Role, opts ...Option) error
	Delete(ctx context.Context, obj *model.Role, opts ...Option) error // 增加选项，支持where条件删除
	Query(ctx context.Context, opts ...Option) (*model.Role, error)
	List(ctx context.Context, page, pageSize int, colum, oder string, opts ...Option) (total int64, objs []*model.Role, err error)
	AppendAssociation(ctx context.Context, model *model.Role, objName string, obj any) error
	ReplaceAssociation(ctx context.Context, model *model.Role, objName string, obj any) error
	ClearAssociation(ctx context.Context, model *model.Role, objName string) error
}

func NewRoleStore(dbProvider DBProviderInterface) RoleStorer {
	return NewRepository[model.Role](dbProvider)
}

type ApiStorer interface {
	Create(ctx context.Context, obj *model.Api) error
	CreateBatch(ctx context.Context, objs []*model.Api) error // 批量创建
	Update(ctx context.Context, obj *model.Api, opts ...Option) error
	Delete(ctx context.Context, obj *model.Api, opts ...Option) error // 增加选项，支持where条件删除
	Query(ctx context.Context, opts ...Option) (*model.Api, error)
	List(ctx context.Context, page, pageSize int, colum, oder string, opts ...Option) (total int64, objs []*model.Api, err error)
}

func NewApiStore(dbProvider DBProviderInterface) ApiStorer {
	return NewRepository[model.Api](dbProvider)
}

type CasbinStorer interface {
	Create(ctx context.Context, obj *model.CasbinRule) error
	CreateBatch(ctx context.Context, objs []*model.CasbinRule) error // 批量创建
	Update(ctx context.Context, obj *model.CasbinRule, opts ...Option) error
	Delete(ctx context.Context, obj *model.CasbinRule, opts ...Option) error         // 增加选项，支持where条件删除
	DeleteBatch(ctx context.Context, objs []*model.CasbinRule, opts ...Option) error // 批量删除
	Query(ctx context.Context, opts ...Option) (*model.CasbinRule, error)
	List(ctx context.Context, page, pageSize int, colum, oder string, opts ...Option) (total int64, objs []*model.CasbinRule, err error)
}

func NewCasbinStore(dbProvider DBProviderInterface) CasbinStorer {
	return NewRepository[model.CasbinRule](dbProvider)
}

type FeiShuUserStorer interface {
	Create(ctx context.Context, obj *model.FeiShuUser) error
	CreateBatch(ctx context.Context, objs []*model.FeiShuUser) error
	Update(ctx context.Context, obj *model.FeiShuUser, opts ...Option) error
	Delete(ctx context.Context, obj *model.FeiShuUser, opts ...Option) error
	Query(ctx context.Context, opts ...Option) (*model.FeiShuUser, error)
	List(ctx context.Context, page, pageSize int, colum, oder string, opts ...Option) (total int64, objs []*model.FeiShuUser, err error)
	AppendAssociation(ctx context.Context, model *model.FeiShuUser, objName string, obj any) error
	ReplaceAssociation(ctx context.Context, model *model.FeiShuUser, objName string, obj any) error
	ClearAssociation(ctx context.Context, model *model.FeiShuUser, objName string) error
}

func NewFeiShuUserStore(dbProvider DBProviderInterface) FeiShuUserStorer {
	return NewRepository[model.FeiShuUser](dbProvider)
}
