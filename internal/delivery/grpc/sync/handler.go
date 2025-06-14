package grpc_sync_handler

import (
	"github.com/FlyKarlik/orderService/internal/usecase"
	"github.com/FlyKarlik/orderService/pkg/logger"
	pb "github.com/FlyKarlik/proto/order_service/gen/order_service/proto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type GRPCSyncHandler struct {
	logger  logger.Logger
	usecase usecase.Usecase
	trace   trace.Tracer
	pb.UnimplementedOrderSyncServiceServer
}

func New(logger logger.Logger, usecase usecase.Usecase) *GRPCSyncHandler {
	return &GRPCSyncHandler{
		logger:  logger,
		usecase: usecase,
		trace:   otel.Tracer("order-service/grpc-sync-handler"),
	}
}
