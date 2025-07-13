package v1

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/yiran15/api-server/base/apitypes"
	"github.com/yiran15/api-server/base/helper"
	"github.com/yiran15/api-server/base/log"
	"github.com/yiran15/api-server/model"
	"github.com/yiran15/api-server/pkg/jwt"
	"github.com/yiran15/api-server/store"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServicer interface {
	Login(ctx context.Context, req *apitypes.UserLoginRequest) (*apitypes.UserLoginResponse, error)
	CreateUser(ctx context.Context, req *apitypes.UserCreateRequest) error
	UpdateUser(ctx context.Context, req *apitypes.UserUpdateRequest) error
	UpdateRole(ctx context.Context, req *apitypes.UserUpdateRoleRequest) error
	DeleteUser(ctx context.Context, req *apitypes.IDRequest) error
	QueryUser(ctx context.Context, req *apitypes.IDRequest) (*model.User, error)
	ListUser(ctx context.Context, pagination *apitypes.UserListRequest) (*apitypes.UserListResponse, error)
}

type UserService struct {
	userStore  store.UserStorer
	roleStore  store.RoleStorer
	cacheStore store.CacheStorer
	tx         store.TxManagerInterface
	jwt        jwt.JwtInterface
}

func NewUserService(userStore store.UserStorer, roleStore store.RoleStorer, cacheStore store.CacheStorer, tx store.TxManagerInterface, jwt jwt.JwtInterface) UserServicer {
	return &UserService{
		userStore:  userStore,
		roleStore:  roleStore,
		cacheStore: cacheStore,
		tx:         tx,
		jwt:        jwt,
	}
}

func (s *UserService) Login(ctx context.Context, req *apitypes.UserLoginRequest) (*apitypes.UserLoginResponse, error) {
	log.WithBody(ctx, req).Info("login request")
	user, err := s.userStore.Query(ctx, store.Where("email", req.Email), store.Preload(model.PreloadRoles))
	if err != nil {
		return nil, err
	}
	if !s.checkPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid password")
	}
	token, err := s.jwt.GenerateToken(user.ID, user.Name)
	if err != nil {
		return nil, err
	}

	if len(user.Roles) > 0 {
		var roles []any
		for _, role := range user.Roles {
			roles = append(roles, role.Name)
			if err := s.cacheStore.SetSet(ctx, store.RoleType, user.Name, roles, nil); err != nil {
				log.WithRequestID(ctx).Error("set role cache error", zap.String("roleName", role.Name), zap.Any("roles", roles), zap.Error(err))
				return nil, err
			}
		}
	}

	return &apitypes.UserLoginResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *apitypes.UserCreateRequest) error {
	if user, err := s.userStore.Query(ctx, store.Where("email", req.Email)); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if user != nil {
			return fmt.Errorf("user %s already exists", req.Email)
		}
	}

	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return err
	}
	user := &model.User{
		Name:     req.Name,
		NickName: req.NickName,
		Email:    req.Email,
		Password: hashedPassword,
		Avatar:   req.Avatar,
		Mobile:   req.Mobile,
	}
	if err = s.userStore.Create(ctx, user); err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *apitypes.UserUpdateRequest) error {
	log.WithBody(ctx, req).Info("update user request")
	mc, err := s.jwt.GetUser(ctx)
	if err != nil {
		return err
	}
	if mc.UserID != req.ID {
		log.WithRequestID(ctx).Error("no permission", zap.Int64("userId", mc.UserID), zap.Int64("reqId", req.ID))
		return errors.New("no permission")
	}

	user, err := s.userStore.Query(ctx, store.Where("id", req.ID), store.Preload(model.PreloadRoles))
	if err != nil {
		return err
	}

	if req.Password != "" {
		hashedPassword, err := s.hashPassword(req.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}

	user.Name = req.Name
	user.NickName = req.NickName
	user.Email = req.Email
	user.Avatar = req.Avatar
	user.Mobile = req.Mobile
	return s.userStore.Update(ctx, user)
}

func (s *UserService) UpdateRole(ctx context.Context, req *apitypes.UserUpdateRoleRequest) error {
	log.WithBody(ctx, req).Info("update user role request")
	var (
		total int64
		err   error
		roles []*model.Role
	)
	req.RoleIds = helper.RemoveDuplicates(req.RoleIds)
	user, err := s.userStore.Query(ctx, store.Where("id", req.ID), store.Preload(model.PreloadRoles))
	if err != nil {
		return err
	}

	total, roles, err = s.roleStore.List(ctx, 0, 0, "", "", store.In("id", req.RoleIds))
	if err != nil {
		return err
	}

	if err = helper.ValidateRoleIds(req.RoleIds, roles, total); err != nil {
		return err
	}

	if len(user.Roles) == 0 {
		if err := s.userStore.AppendAssociation(ctx, user, model.PreloadRoles, roles); err != nil {
			return err
		}
	} else {
		if err := s.userStore.ReplaceAssociation(ctx, user, model.PreloadRoles, roles); err != nil {
			return err
		}
	}

	// 如果redis缓存中存在该用户的角色，需要删除
	cacheRoles, err := s.cacheStore.GetSet(ctx, store.RoleType, user.Name)
	if err != nil {
		return err
	}

	// 如果未找到缓存，直接返回
	if len(cacheRoles) == 0 {
		return nil
	}

	if err := s.cacheStore.DelKey(ctx, store.RoleType, user.Name); err != nil {
		return err
	}

	roleNames := make([]any, 0, len(roles))
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}

	defer func() {
		time.Sleep(time.Second * 5)
		if err := s.cacheStore.DelKey(ctx, store.RoleType, user.Name); err != nil {
			log.WithRequestID(ctx).Error("del role cache error", zap.String("userName", user.Name), zap.Any("roleNames", roleNames), zap.Error(err))
		}
	}()

	return s.cacheStore.SetSet(ctx, store.RoleType, user.Name, roleNames, nil)
}

func (s *UserService) DeleteUser(ctx context.Context, req *apitypes.IDRequest) error {
	log.WithBody(ctx, req).Info("delete user request")
	mc, err := s.jwt.GetUser(ctx)
	if err != nil {
		return err
	}
	if mc.UserID != req.ID {
		log.WithRequestID(ctx).Error("no permission", zap.Int64("userId", mc.UserID), zap.Int64("reqId", req.ID))
		return errors.New("no permission")
	}
	user, err := s.userStore.Query(ctx, store.Where("id", req.ID))
	if err != nil {
		return err
	}

	return s.userStore.Delete(ctx, user)
}

func (s *UserService) QueryUser(ctx context.Context, req *apitypes.IDRequest) (*model.User, error) {
	log.WithBody(ctx, req).Info("query user request")
	return s.userStore.Query(ctx, store.Where("id", req.ID), store.Preload(model.PreloadRoles))
}

func (s *UserService) ListUser(ctx context.Context, req *apitypes.UserListRequest) (*apitypes.UserListResponse, error) {
	log.WithBody(ctx, req).Info("list user request")
	var (
		where store.Option
		filed string
		oder  string
	)
	if req.Name != "" {
		where = store.Where("name", req.Name)
	} else if req.Email != "" {
		where = store.Where("email", req.Email)
	} else if req.Mobile != "" {
		where = store.Where("mobile", req.Mobile)
	}

	if req.SortParam != nil {
		filed = req.SortParam.GetUserSortField(req.Sort)
		oder = req.SortParam.Direction
	}

	total, objs, err := s.userStore.List(ctx, req.Page, req.PageSize, filed, oder, where, store.Preload(model.PreloadRoles))
	if err != nil {
		return nil, err
	}
	res := &apitypes.UserListResponse{
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

// hashPassword 对密码进行 Bcrypt 哈希
func (s *UserService) hashPassword(password string) (string, error) {
	// bcrypt.DefaultCost 是一个合理的默认值，如果需要更高的安全性可以增加
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash password")
	}
	return string(hashedPassword), nil
}

// checkPasswordHash 验证明文密码是否与哈希密码匹配
func (s *UserService) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // 如果没有错误，则匹配成功
}
