package middleware

import "github.com/google/wire"

var MiddlewareProviderSet = wire.NewSet(
	wire.Bind(new(MiddlewareInterface), new(*Middleware)),
	NewMiddleware,
)
