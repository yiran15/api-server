package cmd

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yiran15/api-server/base/conf"
	"github.com/yiran15/api-server/base/log"
	"go.uber.org/zap"
)

const (
	FlagConfigPath = "config-path"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "api-server",
		Long:          `api server, feature include user, role, permission`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.NewLogger()
			viper.SetEnvPrefix("QQLX")
			viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
			viper.AutomaticEnv()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runApp(cmd, args)
		},
	}

	cmd.PersistentFlags().StringP(FlagConfigPath, "C", "./config.yaml", "config file path")
	cmd.AddCommand(NewInitCmd())
	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		zap.L().Fatal("unable to bind command line flag: %v", zap.Error(err))
	}
	return cmd
}

func runApp(_ *cobra.Command, _ []string) error {
	cf := viper.GetString(FlagConfigPath)
	if cf == "" {
		return errors.New("config file path is empty")
	}
	err := conf.LoadConfig(cf)
	if err != nil {
		return fmt.Errorf("load config file faild: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.TODO(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	app, cleanup, err := InitApplication()
	if err != nil {
		return fmt.Errorf("init application faild: %w", err)
	}
	defer cleanup()

	if err := app.Run(ctx); err != nil {
		return fmt.Errorf("run application faild: %w", err)
	}
	zap.L().Info("server exiting")
	return nil
}
