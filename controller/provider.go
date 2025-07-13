package controller

import "github.com/google/wire"

var ControllerProviderSet = wire.NewSet(
	NewUserController,
)
