package pkg

import (
	"github.com/google/wire"
	"github.com/yiran15/api-server/pkg/casbin"
	"github.com/yiran15/api-server/pkg/jwt"
)

var PkgProviderSet = wire.NewSet(
	wire.Bind(new(jwt.JwtInterface), new(*jwt.GenerateToken)),
	casbin.NewEnforcer,
	casbin.NewCasbinManager,
	casbin.NewAuthChecker,
	jwt.NewGenerateToken,
)
