package constant

import "github.com/yiran15/api-server/base/apitypes"

type userContextKey struct{}

type stateContextKey struct{}

type requestIDContextKey struct{}

var UserContextKey = userContextKey{}
var StateContextKey = stateContextKey{}
var RequestIDContextKey = requestIDContextKey{}

var ApiData apitypes.ServerApiData

const (
	FlagConfigPath     = "config-path"
	AuthMidwareKey     = "user"
	RequestIDHeader    = "X-Request-Id"
	EmptyRoleSentinel  = "__empty__"
	OAuth2ProviderList = "oauth2:provider:list"
)
