package constant

import "github.com/yiran15/api-server/base/apitypes"

type userContextKey struct{}

type providerContextKey struct{}

type requestIDContextKey struct{}

var UserContextKey = userContextKey{}
var ProviderContextKey = providerContextKey{}
var RequestIDContextKey = requestIDContextKey{}

var ApiData apitypes.ServerApiData

const (
	FlagConfigPath     = "config-path"
	EmptyRoleSentinel  = "__empty__"
	OAuth2ProviderList = "oauth2:provider:list"
)
