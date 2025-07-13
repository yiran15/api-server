package constant

type UserContextKey struct{}

const (
	AuthMidwareKey      = "user"
	RequestIDHeader     = "X-Request-Id"
	RequestIDContextKey = "requestID"
)
