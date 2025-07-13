package service

import (
	"github.com/google/wire"
	v1 "github.com/yiran15/api-server/service/v1"
)

var ServiceProviderSet = wire.NewSet(
	v1.NewUserService,
	v1.NewRoleService,
	v1.NewApiServicer,
)
