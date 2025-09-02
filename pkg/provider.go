package pkg

import (
	"github.com/google/wire"
	"github.com/yiran15/api-server/pkg/casbin"
	"github.com/yiran15/api-server/pkg/jwt"
	localcache "github.com/yiran15/api-server/pkg/local_cache"
	"github.com/yiran15/api-server/pkg/oauth"
)

var PkgProviderSet = wire.NewSet(
	wire.Bind(new(jwt.JwtInterface), new(*jwt.GenerateToken)),
	jwt.NewGenerateToken,

	casbin.NewEnforcer,
	casbin.NewCasbinManager,
	casbin.NewAuthChecker,
	oauth.NewOAuth2,
	localcache.NewCacher,
)
