package grpc_async_handler

import (
	"context"

	"github.com/FlyKarlik/orderService/internal/delivery/grpc/wrapp"
	"github.com/FlyKarlik/orderService/internal/domain"
	"github.com/FlyKarlik/orderService/internal/mapper"
	shared_context "github.com/FlyKarlik/orderService/pkg/context"
	"github.com/FlyKarlik/orderService/pkg/validate"
	pb "github.com/FlyKarlik/proto/order_service/gen/order_service/proto"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g *GRPCAsyncHandler) StreamOrderUpdates(
	req *pb.StreamOrderUpdatesRequest,
	stream pb.OrderStreamService_StreamOrderUpdatesServer,
) error {
	const layer = "delivery"
	const method = "StreamOrderUpdates"

	ctx, span := g.tracer.Start(stream.Context(), "GRPCAsyncHandler.StreamOrderUpdates")
	defer span.End()

	xRequestID := shared_context.XRequestIDFromContext(ctx)
	span.SetAttributes(
		attribute.String("x-request-id", xRequestID),
		attribute.String("order.id", req.GetOrderId()),
		attribute.String("order.user_id", req.GetUserId()),
	)

	g.logger.Info(layer, method, "stream started",
		"x_request_id", xRequestID,
		"order_id", req.GetOrderId(),
		"user_id", req.GetUserId(),
	)

	domainReq := mapper.FromProtoStreamOrderUpdatesRequest(req)
	if err := validate.Validate(domainReq); err != nil {
		g.logger.Error(layer, method, "invalid stream order updates request", err)
		span.RecordError(err)
		return status.Error(codes.InvalidArgument, err.Error())
	}

	ch, cancel, err := g.usecase.SubscribeToOrderStatus(ctx, domainReq)
	if err != nil {
		g.logger.Error(layer, method, "failed to subscribe to order status", err)
		span.RecordError(err)
		return status.Error(wrapp.GetStatusCodeFromError(err), err.Error())
	}
	defer cancel()

	return g.streamOrderUpdates(ctx, ch, stream, xRequestID)
}

func (g *GRPCAsyncHandler) streamOrderUpdates(
	ctx context.Context,
	ch <-chan domain.StreamOrderUpdatesResponse,
	stream pb.OrderStreamService_StreamOrderUpdatesServer,
	xRequestID string,
) error {
	const layer = "delivery"
	const method = "streamOrderUpdates"

	ctx, span := g.tracer.Start(ctx, "GRPCAsyncHandler.streamOrderUpdates")
	defer span.End()

	for {
		select {
		case <-ctx.Done():
			g.logger.Info(layer, method, "stream context cancelled",
				"x_request_id", xRequestID,
			)
			return nil

		case update, ok := <-ch:
			if !ok {
				g.logger.Info(layer, method, "stream channel closed",
					"x_request_id", xRequestID,
				)
				return nil
			}

			g.logger.Info(layer, method, "sending order status update",
				"x_request_id", xRequestID,
				"order_id", update.OrderID.String(),
				"status", update.OrderStatus.String(),
			)

			if err := stream.Send(mapper.ToProtoStreamOrderUpdatesResponse(update)); err != nil {
				g.logger.Error(layer, method, "failed to send order update", err,
					"x_request_id", xRequestID,
				)
				span.RecordError(err)
				return err
			}
		}
	}
}
