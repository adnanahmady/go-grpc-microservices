package request

import (
	"context"

	"github.com/adnanahmady/go-grpc-microservices/pkg/applog"
)

var (
	ctxTraceIDKey = &struct{ uint8 }{}
	ctxLoggerKey  = &struct{ uint8 }{}
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

func WithLogger(ctx context.Context, lgr applog.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey, lgr)
}

func GetLogger(ctx context.Context) applog.Logger {
	return ctx.Value(ctxLoggerKey).(applog.Logger)
}
