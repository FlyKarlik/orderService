package grpc_async_handler

import (
	"github.com/FlyKarlik/orderService/internal/usecase"
	"github.com/FlyKarlik/orderService/pkg/logger"
	pb "github.com/FlyKarlik/proto/order_service/gen/order_service/proto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type GRPCAsyncHandler struct {
	logger  logger.Logger
	tracer  trace.Tracer
	usecase usecase.Usecase
	pb.UnimplementedOrderStreamServiceServer
}

func New(logger logger.Logger, usecase usecase.Usecase) *GRPCAsyncHandler {
	return &GRPCAsyncHandler{
		logger:  logger,
		usecase: usecase,
		tracer:  otel.Tracer("order-service/grpc-async-handler"),
	}
}
