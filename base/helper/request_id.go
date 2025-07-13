package helper

import "context"

type RequestIDContextKey struct{}

func GetRequestIDFromContext(ctx context.Context) string {
	if reqID, ok := ctx.Value(RequestIDContextKey{}).(string); ok {
		return reqID
	}
	return ""
}
