package constant

import "github.com/yiran15/api-server/base/apitypes"

type userContextKey struct{}

var UserContextKey = userContextKey{}

var ApiData apitypes.ServerApiData

const (
	AuthMidwareKey      = "user"
	RequestIDHeader     = "X-Request-Id"
	RequestIDContextKey = "requestID"
	EmptyRoleSentinel   = "__empty__"
)
