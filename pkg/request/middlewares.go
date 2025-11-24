package request

import (
	"context"
	"time"

	"github.com/adnanahmady/go-grpc-microservices/config"
	"github.com/adnanahmady/go-grpc-microservices/pkg/applog"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type Middlewares struct {
	lgr applog.Logger
	cfg *config.Config
}

func NewMiddlewares(lgr applog.Logger, cfg *config.Config) *Middlewares {
	return &Middlewares{lgr: lgr, cfg: cfg}
}

func (m *Middlewares) UnaryServerLoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now()
		if GetTraceID(ctx) == "anonymouse" {
			ctx = WithTraceID(ctx, uuid.NewString())
			ctx = WithLogger(ctx, m.lgr.NewWith(
				"trace_id", GetTraceID(ctx),
				"method", info.FullMethod,
			))
		}
		reqLgr := GetLogger(ctx)
		reqLgr.Info("received request")

		resp, err := handler(ctx, req)

		duration := time.Since(startTime)
		statusCode := status.Code(err)

		reqLgr.Info("RPC response", "body", resp)
		if err != nil {
			reqLgr.Error(
				"RPC failed",
				err,
				"status_code", statusCode,
				"duration", duration,
			)
		} else {
			reqLgr.Info(
				"RPC completed",
				"status_code", statusCode,
				"duration", duration,
			)
		}

		return resp, err
	}
}
