package grpc_interceptor

import "github.com/FlyKarlik/orderService/pkg/logger"

type GRPCInterceptor struct {
	logger logger.Logger
}

func New(logger logger.Logger) *GRPCInterceptor {
	return &GRPCInterceptor{
		logger: logger,
	}
}
