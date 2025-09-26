package v1

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/yiran15/api-server/base/apitypes"
	"github.com/yiran15/api-server/base/constant"
	"github.com/yiran15/api-server/base/helper"
	"github.com/yiran15/api-server/base/log"
	"github.com/yiran15/api-server/model"
	"github.com/yiran15/api-server/pkg/jwt"
	localcache "github.com/yiran15/api-server/pkg/local_cache"
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
	OAuth2Provider(ctx context.Context) ([]string, error)
	OAuth2Login(provider, state string) (string, error)
	OAuth2Callback(ctx context.Context, req *apitypes.OAuthLoginRequest) (*apitypes.UserLoginResponse, error)
	OAuth2Activate(ctx context.Context, req *apitypes.OAuthActivateRequest) (*apitypes.UserLoginResponse, error)
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
	localCache      localcache.Cacher
}

func NewUserService(userStore store.UserStorer, roleStore store.RoleStorer, cacheStore store.CacheStorer, tx store.TxManagerInterface, jwt jwt.JwtInterface, feishuOauth *oauth.OAuth2, feishuUserStore store.FeiShuUserStorer, localCache localcache.Cacher) UserServicer {
	return &UserService{
		userStore:       userStore,
		roleStore:       roleStore,
		cacheStore:      cacheStore,
		tx:              tx,
		jwt:             jwt,
		oauth:           feishuOauth,
		feishuUserStore: feishuUserStore,
		localCache:      localCache,
	}
}

func (receiver *UserService) Login(ctx context.Context, req *apitypes.UserLoginRequest) (*apitypes.UserLoginResponse, error) {
	user, err := receiver.userStore.Query(ctx, store.Where("email", req.Email), store.Where("status", 1), store.Preload(model.PreloadRoles))
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		log.WithRequestID(ctx).Error("login failed, user not found", zap.String("email", req.Email))
		return nil, constant.ErrLoginFailed
	}

	if !receiver.checkPasswordHash(req.Password, user.Password) {
		log.WithRequestID(ctx).Error("login failed, invalid password", zap.String("email", req.Email))
		return nil, constant.ErrLoginFailed
	}
	token, err := receiver.jwt.GenerateToken(user.ID, user.Name)
	if err != nil {
		return nil, err
	}

	if len(user.Roles) == 0 {
		if err := receiver.cacheStore.SetSet(ctx, store.RoleType, user.ID, []any{constant.EmptyRoleSentinel}, nil); err != nil {
			log.WithRequestID(ctx).Error("login set empty role cache error", zap.Int64("userID", user.ID), zap.Error(err))
		}
	}

	roleNames := make([]any, 0, len(user.Roles))
	for _, role := range user.Roles {
		roleNames = append(roleNames, role.Name)
	}
	if err := receiver.cacheStore.SetSet(ctx, store.RoleType, user.ID, roleNames, nil); err != nil {
		log.WithRequestID(ctx).Error("login set role cache error", zap.Int64("userID", user.ID), zap.Any("roles", roleNames), zap.Error(err))
	}

	return &apitypes.UserLoginResponse{
		User:  user,
		Token: token,
	}, nil
}

func (receiver *UserService) Logout(ctx context.Context) error {
	mc, err := receiver.jwt.GetUser(ctx)
	if err != nil {
		return err
	}
	return receiver.cacheStore.DelKey(ctx, store.RoleType, mc.UserID)
}

func (receiver *UserService) CreateUser(ctx context.Context, req *apitypes.UserCreateRequest) error {
	var (
		user  *model.User
		err   error
		total int64
		roles []*model.Role
	)

	if req.RolesID != nil {
		*req.RolesID = helper.RemoveDuplicates(*req.RolesID)
	}

	if user, err = receiver.userStore.Query(ctx, store.Where("email", req.Email), store.Where("status", 1)); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	if user != nil {
		return fmt.Errorf("user %s already exists", req.Name)
	}

	hashedPassword, err := receiver.hashPassword(req.Password)
	if err != nil {
		return err
	}

	if req.RolesID != nil {
		total, roles, err = receiver.roleStore.List(ctx, 0, 0, "", "", store.In("id", *req.RolesID))
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

	return receiver.tx.Transaction(ctx, func(ctx context.Context) error {
		if err = receiver.userStore.Create(ctx, user); err != nil {
			return err
		}

		if req.RolesID == nil {
			return nil
		}

		return receiver.userStore.AppendAssociation(ctx, user, model.PreloadRoles, roles)
	})
}

func (receiver *UserService) UpdateUserByAdmin(ctx context.Context, req *apitypes.UserUpdateAdminRequest) error {
	if err := receiver.updateUser(ctx, nil, req); err != nil {
		return err
	}

	if req.RolesID == nil {
		return nil
	}

	return receiver.updateRole(ctx, &apitypes.UserUpdateRoleRequest{
		ID:      req.ID,
		RolesID: *req.RolesID,
	})
}

func (receiver *UserService) UpdateUserBySelf(ctx context.Context, req *apitypes.UserUpdateSelfRequest) error {
	mc, err := receiver.jwt.GetUser(ctx)
	if err != nil {
		return err
	}
	user, err := receiver.userStore.Query(ctx, store.Where("id", mc.UserID))
	if err != nil {
		return err
	}
	if req.OldPassword == "" {
		return errors.New("old password is required")
	}
	if !receiver.checkPasswordHash(req.OldPassword, user.Password) {
		return errors.New("invalid old password")
	}
	newReq := new(apitypes.UserUpdateAdminRequest)
	newReq.ID = mc.UserID
	newReq.UserUpdateSelfRequest = req
	return receiver.updateUser(ctx, user, newReq)
}

func (receiver *UserService) DeleteUser(ctx context.Context, req *apitypes.IDRequest) error {
	user, err := receiver.userStore.Query(ctx, store.Where("id", req.ID))
	if err != nil {
		return err
	}
	if err := receiver.userStore.Delete(ctx, user); err != nil {
		return err
	}

	feishuUser, err := receiver.feishuUserStore.Query(ctx, store.Where("user_id", req.ID))
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	if feishuUser != nil {
		if err := receiver.feishuUserStore.Delete(ctx, feishuUser); err != nil {
			return err
		}
	}

	return receiver.userStore.ClearAssociation(ctx, user, model.PreloadRoles)
}

func (receiver *UserService) QueryUser(ctx context.Context, req *apitypes.IDRequest) (*model.User, error) {
	return receiver.userStore.Query(ctx, store.Where("id", req.ID), store.Preload(model.PreloadRoles))
}

func (receiver *UserService) Info(ctx context.Context) (*model.User, error) {
	mc, err := receiver.jwt.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if mc.UserID == 0 {
		log.WithRequestID(ctx).Error("user not found", zap.Int64("userId", mc.UserID), zap.String("userName", mc.UserName))
		return nil, errors.New("user not found")
	}
	return receiver.userStore.Query(ctx, store.Where("id", mc.UserID), store.Preload(model.PreloadRoles))
}

func (receiver *UserService) ListUser(ctx context.Context, req *apitypes.UserListRequest) (*apitypes.UserListResponse, error) {
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

	total, objs, err := receiver.userStore.List(ctx, req.Page, req.PageSize, filed, oder, likeOpt, statusOpt)
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

func (receiver *UserService) updateUser(ctx context.Context, user *model.User, req *apitypes.UserUpdateAdminRequest) error {
	var err error
	if user == nil {
		user, err = receiver.userStore.Query(ctx, store.Where("id", req.ID))
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
			hashedPassword, err := receiver.hashPassword(req.Password)
			if err != nil {
				return err
			}
			user.Password = hashedPassword
		}
	}

	if req.Status != 0 {
		user.Status = &req.Status
	}
	return receiver.userStore.Update(ctx, user)
}

func (receiver *UserService) updateRole(ctx context.Context, req *apitypes.UserUpdateRoleRequest) error {
	var (
		total int64
		err   error
		roles []*model.Role
	)
	req.RolesID = helper.RemoveDuplicates(req.RolesID)
	user, err := receiver.userStore.Query(ctx, store.Where("id", req.ID), store.Preload(model.PreloadRoles))
	if err != nil {
		return err
	}

	total, roles, err = receiver.roleStore.List(ctx, 0, 0, "", "", store.In("id", req.RolesID))
	if err != nil {
		return err
	}

	if err = helper.ValidateRoleIds(req.RolesID, roles, total); err != nil {
		return err
	}

	if err := receiver.userStore.ReplaceAssociation(ctx, user, model.PreloadRoles, roles); err != nil {
		return err
	}

	// 如果redis缓存中存在该用户的角色，需要删除
	cacheRoles, err := receiver.cacheStore.GetSet(ctx, store.RoleType, user.ID)
	if err != nil {
		return err
	}

	// 如果未找到缓存，直接返回
	if len(cacheRoles) == 0 {
		return nil
	}

	if err := receiver.cacheStore.DelKey(ctx, store.RoleType, user.ID); err != nil {
		return err
	}

	roleNames := make([]any, 0, len(roles))
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}

	go func() {
		time.Sleep(time.Second * 5)
		if err := receiver.cacheStore.DelKey(context.TODO(), store.RoleType, user.ID); err != nil {
			log.WithRequestID(ctx).Error("del role cache error", zap.Int64("userID", user.ID), zap.Any("roleNames", roleNames), zap.Error(err))
			return
		}
		log.WithRequestID(ctx).Info("del role cache success", zap.Int64("userID", user.ID), zap.Any("roleNames", roleNames))
	}()

	return receiver.cacheStore.SetSet(ctx, store.RoleType, user.ID, roleNames, nil)
}

// hashPassword 对密码进行 Bcrypt 哈希
func (receiver *UserService) hashPassword(password string) (string, error) {
	// bcrypt.DefaultCost 是一个合理的默认值，如果需要更高的安全性可以增加
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// checkPasswordHash 验证明文密码是否与哈希密码匹配
func (receiver *UserService) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // 如果没有错误，则匹配成功
}

func (receiver *UserService) OAuth2Login(provider, state string) (string, error) {
	return receiver.oauth.Redirect(state, provider), nil
}

func (receiver *UserService) OAuth2Callback(ctx context.Context, req *apitypes.OAuthLoginRequest) (*apitypes.UserLoginResponse, error) {
	var (
		userID   int64
		userName string
		roles    []*model.Role
		user     *model.User
	)
	provider, ok := ctx.Value(constant.ProviderContextKey).(string)
	if !ok {
		return nil, errors.New("invalid provider")
	}

	oauthToken, err := receiver.oauth.Auth(ctx, req.Code, provider)
	if err != nil {
		return nil, err
	}

	userInfo, err := receiver.oauth.UserInfo(ctx, oauthToken, provider)
	if err != nil {
		return nil, err
	}

	switch v := userInfo.(type) {
	case *model.FeiShuUser:
		feishuUser, err := receiver.feishuLogin(ctx, v)
		if err != nil {
			return nil, err
		}
		if feishuUser == nil || feishuUser.User == nil {
			return nil, errors.New("feishu user not found after login")
		}
		user = feishuUser.User
		userID = user.ID
		userName = user.Name
		roles = user.Roles
		if user.Status != nil && *user.Status != model.UserStatusActive {
			return &apitypes.UserLoginResponse{User: user, Token: ""}, nil
		}

	case *model.KeycloakUser:
		u, err := receiver.genericLogin(ctx, v)
		if err != nil {
			return nil, err
		}
		if u == nil {
			return nil, errors.New("generic user not found after login")
		}
		user = u
		userID = user.ID
		userName = user.Name
		roles = user.Roles
		if user.Status != nil && *user.Status != model.UserStatusActive {
			return &apitypes.UserLoginResponse{User: user, Token: ""}, nil
		}
	default:
		return nil, errors.New("unsupported oauth user type")
	}

	token, err := receiver.jwt.GenerateToken(userID, userName)
	if err != nil {
		return nil, err
	}

	roleNames := make([]any, 0, len(roles))
	if len(roles) > 0 {
		for _, r := range roles {
			if r == nil {
				continue
			}
			roleNames = append(roleNames, r.Name)
		}
	}

	if len(roleNames) > 0 {
		if err := receiver.cacheStore.SetSet(ctx, store.RoleType, userID, roleNames, nil); err != nil {
			log.WithRequestID(ctx).Error("login set role cache error", zap.Int64("userID", userID), zap.Any("roles", roleNames), zap.Error(err))
		}
	} else {
		// set a sentinel so other parts know user has no roles
		if err := receiver.cacheStore.SetSet(ctx, store.RoleType, userID, []any{constant.EmptyRoleSentinel}, nil); err != nil {
			log.WithRequestID(ctx).Error("login set empty role cache error", zap.Int64("userID", userID), zap.Error(err))
		}
	}

	return &apitypes.UserLoginResponse{User: user, Token: token}, nil
}

func (receiver *UserService) feishuLogin(ctx context.Context, userInfo *model.FeiShuUser) (*model.FeiShuUser, error) {
	if userInfo.UserID == "" {
		return nil, errors.New("feishu user is empty")
	}

	var email string
	if userInfo.EnterpriseEmail != "" {
		email = userInfo.EnterpriseEmail
	} else if userInfo.Email != "" {
		email = userInfo.Email
	}

	u := &model.User{
		Name:     userInfo.EnName,
		NickName: userInfo.EnName,
		Avatar:   userInfo.AvatarUrl,
		Mobile:   userInfo.Mobile,
		Status:   helper.Int(model.UserStatusInactive),
		Email:    email,
	}

	feishuUser, err := receiver.feishuUserStore.Query(ctx, store.Where("user_id", userInfo.UserID), store.Preload("User.Roles"))
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if feishuUser == nil {
			feishuUser = userInfo
		}
		if feishuUser.User == nil {
			feishuUser.User = u
		}
		if err := receiver.feishuUserStore.Create(ctx, feishuUser); err != nil {
			return nil, err
		}
		return feishuUser, nil
	}

	if feishuUser.User == nil {
		if err := receiver.userStore.Create(ctx, u); err != nil {
			return nil, err
		}
		feishuUser.User = u
	}

	return feishuUser, nil
}

func (receiver *UserService) genericLogin(ctx context.Context, userInfo *model.KeycloakUser) (data *model.User, err error) {
	if userInfo.Sub == "" {
		return nil, errors.New("generic user is empty")
	}

	data, err = receiver.userStore.Query(ctx, store.Where("email", userInfo.Email), store.Preload("Roles"))
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
			_, roles, err := receiver.roleStore.List(ctx, 0, 0, "", "", store.In("name", userInfo.Roles))
			if err != nil {
				return nil, err
			}
			data.Roles = roles
		}
		if err := receiver.userStore.Create(ctx, data); err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (receiver *UserService) OAuth2Provider(_ context.Context) ([]string, error) {
	data, err := receiver.localCache.GetCache(constant.OAuth2ProviderList)
	if err != nil {
		return nil, err
	}
	list, ok := data.([]string)
	if !ok {
		return nil, errors.New("get oauth2 provider list error")
	}
	sort.Strings(list)
	return list, nil
}

func (receiver *UserService) OAuth2Activate(ctx context.Context, req *apitypes.OAuthActivateRequest) (*apitypes.UserLoginResponse, error) {
	if req.Password != req.ConfirmPassword {
		return nil, errors.New("password not match")
	}

	user, err := receiver.userStore.Query(ctx, store.Where("id", req.ID))
	if err != nil {
		return nil, err
	}

	password, err := receiver.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password error: %v", err)
	}
	user.Password = password
	user.Status = helper.Int(model.UserStatusActive)
	if err := receiver.userStore.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update user error: %v", err)
	}
	token, err := receiver.jwt.GenerateToken(user.ID, user.Name)
	if err != nil {
		return nil, err
	}
	return &apitypes.UserLoginResponse{User: user, Token: token}, nil
}
