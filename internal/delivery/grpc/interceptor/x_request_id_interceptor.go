package grpc_interceptor

import (
	"context"

	shared_context "github.com/FlyKarlik/orderService/pkg/context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (i *GRPCInterceptor) XRequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		var requestID string
		if ok {
			ids := md.Get(shared_context.ContextKeyEnumXRequestID.String())
			if len(ids) > 0 {
				requestID = ids[0]
			}
		}
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = context.WithValue(ctx, shared_context.ContextKeyEnumXRequestID, requestID)

		return handler(ctx, req)
	}
}

func (i *GRPCInterceptor) XRequestIDUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		xReqID := shared_context.XRequestIDFromContext(ctx)
		if xReqID == "" {
			xReqID = uuid.New().String()
		}

		ctx = metadata.AppendToOutgoingContext(ctx, shared_context.ContextKeyEnumXRequestID.String(), xReqID)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
