package constant

type userContextKey struct{}

var UserContextKey = userContextKey{}

const (
	AuthMidwareKey      = "user"
	RequestIDHeader     = "X-Request-Id"
	RequestIDContextKey = "requestID"
	EmptyRoleSentinel   = "__empty__"
)
