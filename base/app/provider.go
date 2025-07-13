package app

import "github.com/google/wire"

var AppProviderSet = wire.NewSet(
	NewApplication,
)
