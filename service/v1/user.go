package v1

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/yiran15/api-server/base/apitypes"
	"github.com/yiran15/api-server/base/constant"
	"github.com/yiran15/api-server/base/helper"
	"github.com/yiran15/api-server/base/log"
	"github.com/yiran15/api-server/model"
	"github.com/yiran15/api-server/pkg/jwt"
	"github.com/yiran15/api-server/pkg/oauth"
	"github.com/yiran15/api-server/store"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServicer interface {
	GeneralUserServicer
	OAuthServicer
}

type OAuthServicer interface {
	OAuthLogin(ctx context.Context) (string, error)
	OAuthCallback(ctx context.Context, req *apitypes.OAuthLoginRequest) (*apitypes.UserLoginResponse, error)
}

type GeneralUserServicer interface {
	Login(ctx context.Context, req *apitypes.UserLoginRequest) (*apitypes.UserLoginResponse, error)
	Logout(ctx context.Context) error
	Info(ctx context.Context) (*model.User, error)
	CreateUser(ctx context.Context, req *apitypes.UserCreateRequest) error
	UpdateUserByAdmin(ctx context.Context, req *apitypes.UserUpdateAdminRequest) error
	UpdateUserBySelf(ctx context.Context, req *apitypes.UserUpdateSelfRequest) error
	DeleteUser(ctx context.Context, req *apitypes.IDRequest) error
	QueryUser(ctx context.Context, req *apitypes.IDRequest) (*model.User, error)
	ListUser(ctx context.Context, pagination *apitypes.UserListRequest) (*apitypes.UserListResponse, error)
}

type UserService struct {
	userStore       store.UserStorer
	roleStore       store.RoleStorer
	cacheStore      store.CacheStorer
	tx              store.TxManagerInterface
	jwt             jwt.JwtInterface
	oauth           *oauth.OAuth2
	feishuUserStore store.FeiShuUserStorer
}

func NewUserService(userStore store.UserStorer, roleStore store.RoleStorer, cacheStore store.CacheStorer, tx store.TxManagerInterface, jwt jwt.JwtInterface, feishuOauth *oauth.OAuth2, feishuUserStore store.FeiShuUserStorer) UserServicer {
	return &UserService{
		userStore:       userStore,
		roleStore:       roleStore,
		cacheStore:      cacheStore,
		tx:              tx,
		jwt:             jwt,
		oauth:           feishuOauth,
		feishuUserStore: feishuUserStore,
	}
}

func (s *UserService) Login(ctx context.Context, req *apitypes.UserLoginRequest) (*apitypes.UserLoginResponse, error) {
	user, err := s.userStore.Query(ctx, store.Where("email", req.Email), store.Where("status", 1), store.Preload(model.PreloadRoles))
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, errors.New("user not found")
	}

	if !s.checkPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid password")
	}
	token, err := s.jwt.GenerateToken(user.ID, user.Name)
	if err != nil {
		return nil, err
	}

	if len(user.Roles) == 0 {
		if err := s.cacheStore.SetSet(ctx, store.RoleType, user.ID, []any{constant.EmptyRoleSentinel}, nil); err != nil {
			log.WithRequestID(ctx).Error("login set empty role cache error", zap.Int64("userID", user.ID), zap.Error(err))
		}
	}

	roleNames := make([]any, 0, len(user.Roles))
	for _, role := range user.Roles {
		roleNames = append(roleNames, role.Name)
	}
	if err := s.cacheStore.SetSet(ctx, store.RoleType, user.ID, roleNames, nil); err != nil {
		log.WithRequestID(ctx).Error("login set role cache error", zap.Int64("userID", user.ID), zap.Any("roles", roleNames), zap.Error(err))
	}

	return &apitypes.UserLoginResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *UserService) Logout(ctx context.Context) error {
	mc, err := s.jwt.GetUser(ctx)
	if err != nil {
		return err
	}
	return s.cacheStore.DelKey(ctx, store.RoleType, mc.UserID)
}

func (s *UserService) CreateUser(ctx context.Context, req *apitypes.UserCreateRequest) error {
	var (
		user  *model.User
		err   error
		total int64
		roles []*model.Role
	)

	if req.RolesID != nil {
		*req.RolesID = helper.RemoveDuplicates(*req.RolesID)
	}

	if user, err = s.userStore.Query(ctx, store.Where("email", req.Email), store.Where("status", 1)); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	if user != nil {
		return fmt.Errorf("user %s already exists", req.Name)
	}

	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return err
	}

	if req.RolesID != nil {
		total, roles, err = s.roleStore.List(ctx, 0, 0, "", "", store.In("id", *req.RolesID))
		if err != nil {
			return err
		}
		if err = helper.ValidateRoleIds(*req.RolesID, roles, total); err != nil {
			return err
		}
	}

	user = &model.User{
		Name:     req.Name,
		NickName: req.NickName,
		Email:    req.Email,
		Password: hashedPassword,
		Avatar:   req.Avatar,
		Mobile:   req.Mobile,
	}

	return s.tx.Transaction(ctx, func(ctx context.Context) error {
		if err = s.userStore.Create(ctx, user); err != nil {
			return err
		}

		if req.RolesID == nil {
			return nil
		}

		return s.userStore.AppendAssociation(ctx, user, model.PreloadRoles, roles)
	})
}

func (s *UserService) UpdateUserByAdmin(ctx context.Context, req *apitypes.UserUpdateAdminRequest) error {
	if err := s.updateUser(ctx, nil, req); err != nil {
		return err
	}

	if req.RolesID == nil {
		return nil
	}

	return s.updateRole(ctx, &apitypes.UserUpdateRoleRequest{
		ID:      req.ID,
		RolesID: *req.RolesID,
	})
}

func (s *UserService) UpdateUserBySelf(ctx context.Context, req *apitypes.UserUpdateSelfRequest) error {
	mc, err := s.jwt.GetUser(ctx)
	if err != nil {
		return err
	}
	user, err := s.userStore.Query(ctx, store.Where("id", mc.UserID))
	if err != nil {
		return err
	}
	if req.OldPassword == "" {
		return errors.New("old password is required")
	}
	if !s.checkPasswordHash(req.OldPassword, user.Password) {
		return errors.New("invalid old password")
	}
	newReq := new(apitypes.UserUpdateAdminRequest)
	newReq.ID = mc.UserID
	newReq.UserUpdateSelfRequest = req
	return s.updateUser(ctx, user, newReq)
}

func (s *UserService) DeleteUser(ctx context.Context, req *apitypes.IDRequest) error {
	user, err := s.userStore.Query(ctx, store.Where("id", req.ID))
	if err != nil {
		return err
	}
	if err := s.userStore.Delete(ctx, user); err != nil {
		return err
	}

	return s.userStore.ClearAssociation(ctx, user, model.PreloadRoles)
}

func (s *UserService) QueryUser(ctx context.Context, req *apitypes.IDRequest) (*model.User, error) {
	return s.userStore.Query(ctx, store.Where("id", req.ID), store.Preload(model.PreloadRoles))
}

func (s *UserService) Info(ctx context.Context) (*model.User, error) {
	mc, err := s.jwt.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if mc.UserID == 0 {
		log.WithRequestID(ctx).Error("user not found", zap.Int64("userId", mc.UserID), zap.String("userName", mc.UserName))
		return nil, errors.New("user not found")
	}
	return s.userStore.Query(ctx, store.Where("id", mc.UserID), store.Preload(model.PreloadRoles))
}

func (s *UserService) ListUser(ctx context.Context, req *apitypes.UserListRequest) (*apitypes.UserListResponse, error) {
	var (
		likeOpt   store.Option
		statusOpt store.Option
		filed     = "id"
		oder      = "desc"
	)

	if req.Name != "" {
		likeOpt = store.Like("name", req.Name+"%")
	} else if req.Email != "" {
		likeOpt = store.Like("email", req.Email+"%")
	} else if req.Mobile != "" {
		likeOpt = store.Like("mobile", req.Mobile+"%")
	} else if req.Department != "" {
		likeOpt = store.Like("department", req.Department+"%")
	}

	if req.Status != 0 {
		statusOpt = store.Where("status", req.Status)
	}

	if req.Sort != "" && req.Direction != "" {
		filed = req.Sort
		oder = req.Direction
	}

	total, objs, err := s.userStore.List(ctx, req.Page, req.PageSize, filed, oder, likeOpt, statusOpt)
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

func (s *UserService) updateUser(ctx context.Context, user *model.User, req *apitypes.UserUpdateAdminRequest) error {
	var err error
	if user == nil {
		user, err = s.userStore.Query(ctx, store.Where("id", req.ID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user %d not found", req.ID)
			}
			return err
		}
	}

	if req.UserUpdateSelfRequest != nil {
		user.Name = req.UserUpdateSelfRequest.Name
		user.NickName = req.UserUpdateSelfRequest.NickName
		user.Email = req.UserUpdateSelfRequest.Email
		user.Avatar = req.UserUpdateSelfRequest.Avatar
		user.Mobile = req.UserUpdateSelfRequest.Mobile
		if req.Password != "" {
			hashedPassword, err := s.hashPassword(req.Password)
			if err != nil {
				return err
			}
			user.Password = hashedPassword
		}
	}

	if req.Status != 0 {
		user.Status = &req.Status
	}
	return s.userStore.Update(ctx, user)
}

func (s *UserService) updateRole(ctx context.Context, req *apitypes.UserUpdateRoleRequest) error {
	var (
		total int64
		err   error
		roles []*model.Role
	)
	req.RolesID = helper.RemoveDuplicates(req.RolesID)
	user, err := s.userStore.Query(ctx, store.Where("id", req.ID), store.Preload(model.PreloadRoles))
	if err != nil {
		return err
	}

	total, roles, err = s.roleStore.List(ctx, 0, 0, "", "", store.In("id", req.RolesID))
	if err != nil {
		return err
	}

	if err = helper.ValidateRoleIds(req.RolesID, roles, total); err != nil {
		return err
	}

	if err := s.userStore.ReplaceAssociation(ctx, user, model.PreloadRoles, roles); err != nil {
		return err
	}

	// 如果redis缓存中存在该用户的角色，需要删除
	cacheRoles, err := s.cacheStore.GetSet(ctx, store.RoleType, user.ID)
	if err != nil {
		return err
	}

	// 如果未找到缓存，直接返回
	if len(cacheRoles) == 0 {
		return nil
	}

	if err := s.cacheStore.DelKey(ctx, store.RoleType, user.ID); err != nil {
		return err
	}

	roleNames := make([]any, 0, len(roles))
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}

	defer func() {
		time.Sleep(time.Second * 5)
		if err := s.cacheStore.DelKey(ctx, store.RoleType, user.ID); err != nil {
			log.WithRequestID(ctx).Error("del role cache error", zap.Int64("userID", user.ID), zap.Any("roleNames", roleNames), zap.Error(err))
		}
	}()

	return s.cacheStore.SetSet(ctx, store.RoleType, user.ID, roleNames, nil)
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

func (s *UserService) OAuthLogin(ctx context.Context) (string, error) {
	state, ok := ctx.Value(constant.StateContextKey).(string)
	if !ok {
		return "", errors.New("state not found")
	}
	return s.oauth.Redirect(state), nil
}

func (s *UserService) OAuthCallback(ctx context.Context, req *apitypes.OAuthLoginRequest) (*apitypes.UserLoginResponse, error) {
	var (
		userID   int64
		userName string
		roles    []*model.Role
		user     *model.User
	)
	oauthToken, err := s.oauth.Auth(ctx, req.State, req.Code)
	if err != nil {
		return nil, err
	}

	userInfo, err := s.oauth.UserInfo(ctx, oauthToken)
	if err != nil {
		return nil, err
	}

	switch v := userInfo.(type) {
	case *model.FeiShuUser:
		feishuUser, err := s.feishuLogin(ctx, v)
		if err != nil {
			return nil, err
		}

		userID = feishuUser.User.ID
		userName = feishuUser.User.Name
		roles = feishuUser.User.Roles
		user = feishuUser.User
		if feishuUser.User == nil || *feishuUser.User.Status == model.UserStatusInactive {
			return &apitypes.UserLoginResponse{
				User:  user,
				Token: "",
			}, nil
		}

	case *model.KeycloakUser:
		u, err := s.genericLogin(ctx, v)
		if err != nil {
			return nil, err
		}
		userID = u.ID
		userName = u.Name
		roles = u.Roles
		user = u
		if *u.Status == model.UserStatusInactive {
			return &apitypes.UserLoginResponse{
				User:  user,
				Token: "",
			}, nil
		}
	}

	token, err := s.jwt.GenerateToken(userID, userName)
	if err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		if err := s.cacheStore.SetSet(ctx, store.RoleType, userID, []any{constant.EmptyRoleSentinel}, nil); err != nil {
			log.WithRequestID(ctx).Error("login set empty role cache error", zap.Int64("userID", userID), zap.Error(err))
		}

		return &apitypes.UserLoginResponse{
			User:  user,
			Token: token,
		}, nil
	}

	roleNames := make([]any, 0, len(roles))
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}
	if err := s.cacheStore.SetSet(ctx, store.RoleType, userID, roleNames, nil); err != nil {
		log.WithRequestID(ctx).Error("login set role cache error", zap.Int64("userID", userID), zap.Any("roles", roleNames), zap.Error(err))
	}

	return &apitypes.UserLoginResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *UserService) feishuLogin(ctx context.Context, userInfo *model.FeiShuUser) (data *model.FeiShuUser, err error) {
	if userInfo.UserID == "" {
		return nil, errors.New("feishu user is empty")
	}
	var email string
	if userInfo.EnterpriseEmail != "" {
		email = userInfo.EnterpriseEmail
	} else if userInfo.Email != "" {
		email = userInfo.Email
	}

	feishuUser, err := s.feishuUserStore.Query(ctx, store.Where("user_id", userInfo.UserID), store.Preload("User.Roles"))
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		userInfo.User = &model.User{
			Name:     userInfo.EnName,
			NickName: userInfo.EnName,
			Avatar:   userInfo.AvatarUrl,
			Mobile:   userInfo.Mobile,
			Status:   helper.Int(model.UserStatusInactive),
			Email:    email,
		}

		if err := s.feishuUserStore.Create(ctx, userInfo); err != nil {
			return nil, err
		}
		feishuUser = userInfo
	}

	return feishuUser, nil
}

func (s *UserService) genericLogin(ctx context.Context, userInfo *model.KeycloakUser) (data *model.User, err error) {
	if userInfo.Sub == "" {
		return nil, errors.New("generic user is empty")
	}

	data, err = s.userStore.Query(ctx, store.Where("email", userInfo.Email), store.Preload("Roles"))
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		data = &model.User{
			Name:       userInfo.PreferredUsername,
			NickName:   userInfo.FamilyName + userInfo.GivenName,
			Email:      userInfo.Email,
			Status:     helper.Int(model.UserStatusInactive),
			Department: strings.Join(userInfo.Group, ","),
		}
		if len(userInfo.Roles) > 0 {
			_, roles, err := s.roleStore.List(ctx, 0, 0, "", "", store.In("name", userInfo.Roles))
			if err != nil {
				return nil, err
			}
			data.Roles = roles
		}
		if err := s.userStore.Create(ctx, data); err != nil {
			return nil, err
		}
	}

	return data, nil
}
