package constant

import "errors"

var (
	ErrNoAuthHeader      = errors.New("no auth header")
	ErrInvalidAuthHeader = errors.New("invalid auth header")
	ErrNoPermission      = errors.New("no permission")
)
