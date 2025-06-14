package grpc_sync_handler

import (
	"context"

	"github.com/FlyKarlik/orderService/internal/delivery/grpc/wrapp"
	"github.com/FlyKarlik/orderService/internal/mapper"
	shared_context "github.com/FlyKarlik/orderService/pkg/context"
	"github.com/FlyKarlik/orderService/pkg/validate"
	pb "github.com/FlyKarlik/proto/order_service/gen/order_service/proto"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g *GRPCSyncHandler) CreateOrder(
	ctx context.Context,
	req *pb.CreateOrderRequest,
) (*pb.CreateOrderResponse, error) {
	const layer = "delivery"
	const method = "CreateOrder"

	ctx, span := g.trace.Start(ctx, "GRPCSyncHandler.CreateOrder")
	defer span.End()

	xRequestID := shared_context.XRequestIDFromContext(ctx)

	span.SetAttributes(
		attribute.String("x-request-id", xRequestID),
		attribute.String("rpc.system", "grpc"),
		attribute.String("rpc.service", "OrderSyncService"),
		attribute.String("rpc.method", method),
		attribute.String("order.user_id", req.GetUserId()),
		attribute.String("order.market_id", req.GetMarketId()),
		attribute.String("order.type", req.GetOrderType().String()),
		attribute.String("order.price", req.GetPrice()),
		attribute.Int64("order.quantity", req.GetQuantity()),
	)

	domainReq := mapper.FromProtoCreateOrderRequest(req)

	if err := validate.Validate(domainReq); err != nil {
		g.logger.Error(layer, method, "invalid create order request", err)
		span.RecordError(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := g.usecase.CreateOrder(ctx, domainReq)
	if err != nil {
		code := wrapp.GetStatusCodeFromError(err)
		g.logger.Error(layer, method, "failed to create order", err)
		span.RecordError(err)
		span.SetAttributes(attribute.String("grpc.code", code.String()))
		return nil, status.Error(code, err.Error())
	}

	return mapper.ToProtoCreateOrderResponse(resp), nil
}

func (g *GRPCSyncHandler) GetOrderStatus(
	ctx context.Context,
	req *pb.GetOrderStatusRequest,
) (*pb.GetOrderStatusResponse, error) {
	const layer = "delivery"
	const method = "GetOrderStatus"

	ctx, span := g.trace.Start(ctx, "GRPCSyncHandler.GetOrderStatus")
	defer span.End()

	xRequestID := shared_context.XRequestIDFromContext(ctx)

	span.SetAttributes(
		attribute.String("x-request-id", xRequestID),
		attribute.String("rpc.system", "grpc"),
		attribute.String("rpc.service", "OrderSyncService"),
		attribute.String("rpc.method", method),
		attribute.String("order.id", req.GetOrderId()),
		attribute.String("order.user_id", req.GetUserId()),
	)

	domainReq := mapper.FromProtoGetOrderStatusRequest(req)

	if err := validate.Validate(domainReq); err != nil {
		g.logger.Error(layer, method, "invalid get order status request", err)
		span.RecordError(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := g.usecase.GetOrderStatus(ctx, domainReq)
	if err != nil {
		code := wrapp.GetStatusCodeFromError(err)
		g.logger.Error(layer, method, "failed to get order status", err)
		span.RecordError(err)
		span.SetAttributes(attribute.String("grpc.code", code.String()))
		return nil, status.Error(code, err.Error())
	}

	return mapper.ToProtoGetOrderStatusResponse(resp), nil
}
