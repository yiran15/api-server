package router

import "github.com/google/wire"

var RouterProviderSet = wire.NewSet(
	wire.Bind(new(RouterInterface), new(*Router)),
	NewRouter,
)
