package cmd

import (
	"context"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yiran15/api-server/base/apitypes"
	"github.com/yiran15/api-server/base/conf"
	"github.com/yiran15/api-server/base/data"
	"github.com/yiran15/api-server/model"
	"github.com/yiran15/api-server/pkg/casbin"
	"github.com/yiran15/api-server/pkg/jwt"
	v1 "github.com/yiran15/api-server/service/v1"
	"github.com/yiran15/api-server/store"
	"go.uber.org/zap"
)

func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "init",
		Long:          `init api server`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			cf := viper.GetString(FlagConfigPath)
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

	userServicer := v1.NewUserService(userRepo, roleRepo, cacheStore, txManager, generateToken)
	roleServicer := v1.NewRoleService(roleRepo, apiRepo, casbinStore, casbinManager, txManager)
	apiServicer := v1.NewApiServicer(apiRepo)
	return &service{
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

	zap.L().Info("create admin api")
	if err = service.apiService.CreateApi(ctx, &apitypes.ApiCreateRequest{
		Name:        "admin",
		Path:        "*",
		Method:      "*",
		Description: "拥有所有接口权限",
	}); err != nil {
		return err
	}
	apis, err := service.apiService.ListApi(ctx, &apitypes.ApiListRequest{
		Pagination: &apitypes.Pagination{
			Page:     1,
			PageSize: 10,
		},
		Name:   "admin",
		Path:   "*",
		Method: "*",
	})
	if err != nil {
		return err
	}
	var adminApi *model.Api
	for _, api := range apis.List {
		if api.Name == "admin" {
			adminApi = api
		}
	}
	if adminApi == nil {
		return errors.New("admin api not found")
	}

	zap.L().Info("create admin role")
	if err = service.roleService.CreateRole(ctx, &apitypes.RoleCreateRequest{
		Name:        "admin",
		Description: "超级管理员",
		Apis: []int64{
			adminApi.ID,
		},
	}); err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)
		if !ok {
			return err
		}
		if mysqlErr.Number != 1062 {
			return err
		}
	}

	adminRoles, err := service.roleService.ListRole(ctx, &apitypes.RoleListRequest{
		Pagination: &apitypes.Pagination{
			Page:     1,
			PageSize: 10,
		},
		Name: "admin",
	})
	if err != nil {
		return err
	}
	var adminRole *model.Role
	for _, v := range adminRoles.List {
		if v.Name == "admin" {
			adminRole = v
		}
	}
	if adminRole == nil {
		return errors.New("admin role not found")
	}

	zap.L().Info("create admin user")
	adminUserReq := &apitypes.UserCreateRequest{
		Name:     "admin",
		NickName: "超级管理员",
		Email:    "admin@qqlx.net",
		Password: "12345678",
		Avatar:   "https://s3-imfile.feishucdn.com/static-resource/v1/v2_79ff6f58-f5c8-41c2-8ffb-8379d4e57acg~?image_size=noop&cut_type=&quality=&format=image&sticker_format=.webp",
	}
	if err = service.userService.CreateUser(ctx, adminUserReq); err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)
		if !ok {
			return err
		}
		if mysqlErr.Number != 1062 {
			return err
		}
	}

	users, err := service.userService.ListUser(ctx, &apitypes.UserListRequest{
		Pagination: &apitypes.Pagination{
			Page:     1,
			PageSize: 10,
		},
		Name: "admin",
	})
	if err != nil {
		return err
	}
	var creaUser bool
	for _, user := range users.List {
		if user.Name == "admin" {
			creaUser = true
			if err := service.userService.UpdateRole(ctx, &apitypes.UserUpdateRoleRequest{
				ID: user.ID,
				RoleIds: []int64{
					adminRole.ID,
				},
			}); err != nil {
				return err
			}
		}
	}
	if !creaUser {
		return errors.New("admin user not found")
	}
	zap.L().Info("init application success")
	return nil
}
