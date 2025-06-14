package shared_context

import (
	"context"
)

func XRequestIDFromContext(ctx context.Context) string {
	xRequestID, _ := ctx.Value(ContextKeyEnumXRequestID).(string)
	return xRequestID
}
