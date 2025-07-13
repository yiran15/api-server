package main

import (
	"github.com/yiran15/api-server/cmd"
	"go.uber.org/zap"
)

func main() {
	if err := cmd.NewCmd().Execute(); err != nil {
		zap.L().Fatal("execute command faild", zap.Error(err))
	}
}
