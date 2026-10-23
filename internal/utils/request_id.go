package utils

import "context"

type key string

const RequestIDCtxKey key = "request_id"
const (
	RequestID     = "request_id"
	HTTPRequestID = "X-Request-ID"
)

func GetRequestIDFromCtx(ctx context.Context) string {
	reqID, ok := ctx.Value(RequestIDCtxKey).(string)
	if !ok {
		return ""
	}

	return reqID
}
