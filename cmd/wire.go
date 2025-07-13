// go:build wireinject
//go:build wireinject
// +build wireinject

package cmd

import (
	"github.com/google/wire"
	"github.com/yiran15/api-server/base/app"
	"github.com/yiran15/api-server/base/data"
	"github.com/yiran15/api-server/base/middleware"
	"github.com/yiran15/api-server/base/router"
	"github.com/yiran15/api-server/base/server"
	"github.com/yiran15/api-server/controller"
	"github.com/yiran15/api-server/pkg"
	"github.com/yiran15/api-server/service"
	"github.com/yiran15/api-server/store"
)

func InitApplication() (*app.Application, func(), error) {
	panic(wire.Build(
		data.DataProviderSet,
		pkg.PkgProviderSet,
		store.StoreProviderSet,
		service.ServiceProviderSet,
		controller.ControllerProviderSet,
		middleware.MiddlewareProviderSet,
		router.RouterProviderSet,
		server.ServerProviderSet,
		app.AppProviderSet,
	))
}
