package service

import (
	"context"

	"github.com/yiran15/api-server/model"
	"github.com/yiran15/api-server/store"
)

type UserServicer interface {
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, user *model.User) error
	QueryUser(ctx context.Context, user *model.User) (*model.User, error)
	ListUser(ctx context.Context, pagination *store.Pagination, opts ...store.Option) (total int64, objs []*model.User, err error)
}

type UserService struct {
	userStore store.UserStorer
}

func NewUserService(userStore store.UserStorer) UserServicer {
	return &UserService{
		userStore: userStore,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {
	return s.userStore.Create(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, user *model.User) error {
	return s.userStore.Update(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, user *model.User) error {
	return s.userStore.Delete(ctx, user)
}

func (s *UserService) QueryUser(ctx context.Context, user *model.User) (*model.User, error) {
	return s.userStore.Query(ctx, store.Where("id = ?", user.ID))
}

func (s *UserService) ListUser(ctx context.Context, pagination *store.Pagination, opts ...store.Option) (total int64, objs []*model.User, err error) {
	return s.userStore.List(ctx, pagination, opts...)
}
