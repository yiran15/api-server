package cmd

import (
	"context"

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

func initApplication(_ *cobra.Command, _ []string) error {
	ctx := context.Background()
	db, cleanup1, err := data.NewDB()
	if err != nil {
		return err
	}

	defer func() {
		cleanup1()
	}()

	provider := store.NewDBProvider(db)
	userRepo := store.NewRepository[model.User](provider)
	roleRepo := store.NewRepository[model.Role](provider)
	generateToken, err := jwt.NewGenerateToken()
	if err != nil {
		return err
	}

	enforcer, err := casbin.NewEnforcer(db)
	if err != nil {
		return err
	}
	casbinManager := casbin.NewCasbinManager(enforcer)

	userServicer := v1.NewUserService(userRepo, roleRepo, nil, store.NewTxManager(db), generateToken)

	adminUserReq := &apitypes.UserCreateRequest{
		Name:     "admin",
		NickName: "超级管理员",
		Email:    "admin@qqlx.net",
		Password: "12345678",
		Avatar:   "https://s3-imfile.feishucdn.com/static-resource/v1/v2_79ff6f58-f5c8-41c2-8ffb-8379d4e57acg~?image_size=noop&cut_type=&quality=&format=image&sticker_format=.webp",
	}

	zap.L().Info("create admin user")
	if err = userServicer.CreateUser(ctx, adminUserReq); err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)
		if !ok {
			return err
		}
		if mysqlErr.Number != 1062 {
			return err
		}
	}

	zap.L().Info("create admin role")
	if err = roleRepo.Create(ctx, &model.Role{
		Name: "admin",
	}); err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)
		if !ok {
			return err
		}
		if mysqlErr.Number != 1062 {
			return err
		}
	}

	adminUser, err := userRepo.Query(ctx, store.Where("name", "admin"))
	if err != nil {
		return err
	}

	adminRole, err := roleRepo.Query(ctx, store.Where("name", "admin"))
	if err != nil {
		return err
	}

	zap.L().Info("append admin role to admin user")
	if err = userRepo.AppendAssociation(ctx, adminUser, model.PreloadRoles, adminRole); err != nil {
		return err
	}

	zap.L().Info("create admin role policy")
	_, err = casbinManager.AddRolePolicy("admin_role", &model.Api{
		Path:   "*",
		Method: "*",
	})
	if err != nil {
		return err
	}
	_, err = casbinManager.AddUserToRole("admin", "admin_role")
	if err != nil {
		return err
	}
	zap.L().Info("init application success")
	return nil
}
