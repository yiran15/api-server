package constant

import (
	"errors"
)

var (
	ErrAuthFailed   = errors.New("auth failed")
	ErrNoPermission = errors.New("access forbidden")
	ErrLoginFailed  = errors.New("incorrect username or password")
)
