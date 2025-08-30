package cmd

import (
	"context"
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yiran15/api-server/base/apitypes"
	"github.com/yiran15/api-server/base/conf"
	"github.com/yiran15/api-server/base/constant"
	"github.com/yiran15/api-server/base/data"
	"github.com/yiran15/api-server/model"
	"github.com/yiran15/api-server/pkg/casbin"
	"github.com/yiran15/api-server/pkg/jwt"
	v1 "github.com/yiran15/api-server/service/v1"
	"github.com/yiran15/api-server/store"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "init",
		Long:          `init api server`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			cf := viper.GetString(constant.FlagConfigPath)
			if cf == "" {
				zap.L().Fatal("config file path is empty")
			}
			err := conf.LoadConfig(cf)
			if err != nil {
				zap.L().Fatal("load config file faild", zap.String("path", cf), zap.Error(err))
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return initApplication(cmd, args)
		},
	}

	return cmd
}

type service struct {
	db          *gorm.DB
	userService v1.UserServicer
	roleService v1.RoleServicer
	apiService  v1.ApiServicer
}

func getService() (*service, func(), error) {
	db, cleanup1, err := data.NewDB()
	if err != nil {
		return nil, nil, err
	}

	provider := store.NewDBProvider(db)
	userRepo := store.NewUserStore(provider)
	roleRepo := store.NewRoleStore(provider)
	apiRepo := store.NewApiStore(provider)
	casbinStore := store.NewCasbinStore(provider)
	txManager := store.NewTxManager(db)

	redisClient, err := data.NewRDB()
	if err != nil {
		return nil, nil, err
	}
	cacheStore, cleanup2, err := store.NewCacheStore(redisClient)
	if err != nil {
		return nil, nil, err
	}

	generateToken, err := jwt.NewGenerateToken()
	if err != nil {
		return nil, nil, err
	}

	casbinEnforcer, err := casbin.NewEnforcer(db)
	if err != nil {
		return nil, nil, err
	}
	casbinManager := casbin.NewCasbinManager(casbinEnforcer)

	userServicer := v1.NewUserService(userRepo, roleRepo, cacheStore, txManager, generateToken, nil, nil)
	roleServicer := v1.NewRoleService(roleRepo, apiRepo, casbinStore, casbinManager, txManager)
	apiServicer := v1.NewApiServicer(apiRepo)
	return &service{
			db:          db,
			userService: userServicer,
			roleService: roleServicer,
			apiService:  apiServicer,
		}, func() {
			cleanup1()
			cleanup2()
		}, nil
}

func initApplication(_ *cobra.Command, _ []string) error {
	ctx := context.Background()
	service, cleanup, err := getService()
	if err != nil {
		return err
	}
	defer cleanup()

	apis := []apitypes.ApiCreateRequest{
		{
			Name:        "admin",
			Path:        "*",
			Method:      "*",
			Description: "拥有所有接口权限",
		},
		{
			Name:        "readOnly",
			Path:        "*",
			Method:      "GET",
			Description: "只读接口权限",
		},
	}

	var (
		adminApi    model.Api
		readOnlyApi model.Api
	)
	for _, api := range apis {
		var dbApi model.Api
		if err = service.db.Model(&dbApi).Where("name = ?", api.Name).First(&dbApi).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}
		if dbApi.ID == 0 {
			zap.L().Info("create api", zap.String("name", api.Name), zap.String("path", api.Path), zap.String("method", api.Method))
			if err = service.apiService.CreateApi(ctx, &api); err != nil {
				return err
			}
			if err = service.db.Model(&model.Api{}).Where("name = ?", api.Name).First(&dbApi).Error; err != nil {
				return err
			}
		}
		if api.Name == "admin" {
			adminApi = dbApi
		} else {
			readOnlyApi = dbApi
		}
	}

	zap.L().Info("create admin role")
	var adminRole model.Role
	roleCreateRequest := []apitypes.RoleCreateRequest{
		{
			Name:        "admin",
			Description: "所有接口权限",
			Apis: []int64{
				adminApi.ID,
			},
		},
		{
			Name:        "readOnly",
			Description: "只读接口权限",
			Apis: []int64{
				readOnlyApi.ID,
			},
		},
	}
	for _, roleCreateRequest := range roleCreateRequest {
		var dbRole model.Role
		if err = service.db.Model(&model.Role{}).Where("name = ?", roleCreateRequest.Name).First(&dbRole).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}
		if dbRole.ID == 0 {
			if err = service.roleService.CreateRole(ctx, &roleCreateRequest); err != nil {
				return err
			}
			if err = service.db.Model(&model.Role{}).Where("name = ?", roleCreateRequest.Name).First(&dbRole).Error; err != nil {
				return err
			}
		}
		if roleCreateRequest.Name == "admin" {
			adminRole = dbRole
		}
	}

	zap.L().Info("create admin user")
	adminUserReq := &apitypes.UserCreateRequest{
		Name:     "admin",
		NickName: "超级管理员",
		Email:    "admin@qqlx.net",
		Password: "12345678",
		Avatar:   "https://s3-imfile.feishucdn.com/static-resource/v1/v2_79ff6f58-f5c8-41c2-8ffb-8379d4e57acg~?image_size=noop&cut_type=&quality=&format=image&sticker_format=.webp",
	}
	user := model.User{}
	if err = service.db.Model(&model.User{}).Where("name = ?", "admin").First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	if user.ID == 0 {
		if err = service.userService.CreateUser(ctx, adminUserReq); err != nil {
			return err
		}
		if err = service.db.Model(&model.User{}).Where("name = ?", "admin").First(&user).Error; err != nil {
			return err
		}
	}

	userUpdateRequest := &apitypes.UserUpdateAdminRequest{
		ID: user.ID,
		RolesID: &[]int64{
			adminRole.ID,
		},
	}
	if err = service.userService.UpdateUserByAdmin(ctx, userUpdateRequest); err != nil {
		return err
	}
	zap.L().Info("init application success")
	return nil
}
