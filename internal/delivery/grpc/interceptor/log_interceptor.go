package grpc_interceptor

import (
	"context"
	"time"

	shared_context "github.com/FlyKarlik/orderService/pkg/context"
	"google.golang.org/grpc"
)

func (i *GRPCInterceptor) LoggerInterceptor() grpc.UnaryServerInterceptor {
	const layer = "grpc_interceptor"
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		start := time.Now()
		resp, err := handler(ctx, req)

		reqID := ctx.Value(shared_context.ContextKeyEnumXRequestID)
		duration := time.Since(start)

		if err != nil {
			i.logger.Error(layer, info.FullMethod,
				"request failed",
				err,
				"request_id", reqID,
				"duration", duration)
		} else {
			i.logger.Info(layer, info.FullMethod,
				"request completed",
				"request_id", reqID,
				"duration", duration)
		}

		return resp, err
	}
}
