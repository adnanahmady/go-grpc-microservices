package request

import (
	"context"
)

var (
	ctxTraceIDKey = &struct{uint8}{}
)

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, ctxTraceIDKey, traceID)
}

func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(ctxTraceIDKey).(string); ok {
		return traceID
	}
	return "anonymouse"
}