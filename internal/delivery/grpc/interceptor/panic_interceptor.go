package grpc_interceptor

import (
	"context"
	"fmt"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *GRPCInterceptor) UnaryPanicRecoveryInterceptor() grpc.UnaryServerInterceptor {
	const layer = "grpc_interceptor"
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				i.logger.Error(layer, info.FullMethod,
					fmt.Sprintf("panic recovered: %v", r),
					nil,
					"stack", string(stack),
				)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}
