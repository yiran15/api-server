package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yiran15/api-server/base/conf"
	"github.com/yiran15/api-server/base/constant"
	baselog "github.com/yiran15/api-server/base/log"
	"go.uber.org/zap"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "api-server",
		Long:          `api server, feature include user, role, permission`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// 设置环境变量前缀和替换规则，这样就可以通过环境变量来配置了
			serverName := os.Getenv("SERVICE_NAME")
			if serverName != "" {
				viper.Set("server.name", serverName)
				viper.SetEnvPrefix(strings.ToUpper(serverName))
			} else {
				viper.Set("server.name", "api-server")
				viper.SetEnvPrefix("API_SERVER")
			}
			viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
			viper.AutomaticEnv()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runApp(cmd, args)
		},
	}

	cmd.PersistentFlags().StringP(constant.FlagConfigPath, "c", "./config.yaml", "config file path")
	cmd.PersistentFlags().StringP("log-level", "l", "info", "log level, enum: debug, info, warn, error")
	cmd.PersistentFlags().StringP("server-bind", "b", ":8080", "server bind address")
	cmd.AddCommand(NewInitCmd())
	// 将命令行参数中的短横线替换为点，例如 --log-level -> log.level
	// 这样 viper 就可以正确解析命令行参数了
	bindAllFlagsWithNormalize(cmd.PersistentFlags())
	return cmd
}

func runApp(_ *cobra.Command, _ []string) error {
	cf := viper.GetString("config.path")
	if cf == "" {
		return errors.New("config file path is empty")
	}
	err := conf.LoadConfig(cf)
	if err != nil {
		return fmt.Errorf("load config file faild: %w", err)
	}
	baselog.NewLogger()
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

func bindAllFlagsWithNormalize(f *pflag.FlagSet) {
	f.VisitAll(func(flag *pflag.Flag) {
		viperKey := strings.ReplaceAll(flag.Name, "-", ".")
		if err := viper.BindPFlag(viperKey, flag); err != nil {
			log.Fatalf("unable to bind flag %s to viper key %s: %v", flag.Name, viperKey, err)
		}
	})
}
