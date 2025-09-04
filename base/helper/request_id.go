package helper

import (
	"context"

	"github.com/yiran15/api-server/base/constant"
)

func GetRequestIDFromContext(ctx context.Context) string {
	if reqID, ok := ctx.Value(constant.RequestIDContextKey).(string); ok {
		return reqID
	}
	return ""
}
