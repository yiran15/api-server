package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/yiran15/api-server/base/conf"
	"github.com/yiran15/api-server/cmd"
	"go.uber.org/zap"
)

const (
	FlagConfigPath = "config-path"
	ConfigEnv      = "CONFIG_PATH"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "api-server",
		Long: `api server, feature include user, role, permission`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if !cmd.Flags().Changed(FlagConfigPath) {
				envConfigPath := os.Getenv(ConfigEnv)
				if envConfigPath != "" {
					err := cmd.Flags().Set(FlagConfigPath, envConfigPath)
					if err != nil {
						zap.S().Fatalf("set config file path from env %s faild: %v", envConfigPath, err)
						return
					}
				}
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			var (
				cf  string
				err error
			)
			cf, err = cmd.Flags().GetString(FlagConfigPath)
			if err != nil {
				zap.S().Fatalf("get config file path faild: %v", err)
			}
			if cf == "" {
				zap.S().Fatalf("config file path is empty")
			}
			runApp(cf)
		},
	}
	cmd.PersistentFlags().StringP(FlagConfigPath, "C", "./config.yaml", "config file path")
	return cmd
}

func runApp(cf string) {
	err := conf.LoadConfig(cf)
	if err != nil {
		zap.S().Fatalf("load config file %s faild: %v", cf, err)
	}

	ctx, stop := signal.NotifyContext(context.TODO(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	// otelShutdown, err := pkg.SetupOTelSDK(ctx)
	// if err != nil {
	// 	zap.L().Fatal("failed to setup otlp", zap.Error(err))
	// }
	// defer func() {
	// 	errs := errors.Join(err, otelShutdown(ctx))
	// 	if errs != nil {
	// 		zap.L().Error("failed to shutdown otlp", zap.Error(errs))
	// 	}
	// }()

	app, cleanup, err := cmd.InitApplication()
	if err != nil {
		zap.S().Fatalf("init application failed: %v", err)
		return
	}
	defer cleanup()

	if err := app.Run(ctx); err != nil {
		zap.S().Fatalf("run application faild: %v", err)
	}
	zap.S().Infof("server exiting")
}
