package constant

import (
	"errors"
)

var (
	ErrAuthFailed   = errors.New("auth failed")
	ErrNoPermission = errors.New("no permission")
	ErrLoginFailed  = errors.New("incorrect username or password")
)
