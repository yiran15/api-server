package server

import "github.com/google/wire"

var ServerProviderSet = wire.NewSet(
	NewHttpServer,
	wire.Bind(new(ServerInterface), new(*Server)),
	NewServer,
)
