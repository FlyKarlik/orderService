package usecase

import (
	"context"
	"time"

	"github.com/FlyKarlik/orderService/internal/domain"
	"github.com/FlyKarlik/orderService/internal/driver"
	"github.com/FlyKarlik/orderService/internal/errs"
	"github.com/FlyKarlik/orderService/internal/repository"
	shared_context "github.com/FlyKarlik/orderService/pkg/context"
	"github.com/FlyKarlik/orderService/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type orderUsecase struct {
	logger logger.Logger
	driver driver.Driver
	repo   repository.Repository
	tracer trace.Tracer
}

func newOrderUsecase(
	logger logger.Logger,
	driver driver.Driver,
	repo repository.Repository) *orderUsecase {
	return &orderUsecase{
		logger: logger,
		driver: driver,
		repo:   repo,
		tracer: otel.Tracer("order-service/usecase"),
	}
}

func (o *orderUsecase) CreateOrder(
	ctx context.Context,
	req domain.CreateOrderRequest,
) (domain.CreateOrderResponse, error) {
	const layer = "usecase"
	const method = "CreateOrder"

	ctx, span := o.tracer.Start(ctx, "OrderUsecase.CreateOrder")
	defer span.End()

	xReqID := shared_context.XRequestIDFromContext(ctx)

	o.logger.Info(layer, method, "creating order",
		"x_request_id", xReqID,
		"user_id", req.UserID.String(),
		"market_id", req.MarketID.String(),
		"user_roles", req.UserRoles,
	)

	span.SetAttributes(
		attribute.String("x-request-id", xReqID),
		attribute.String("user_id", req.UserID.String()),
		attribute.String("market_id", req.MarketID.String()),
	)

	marketsResp, err := o.driver.ViewMarkets(ctx, domain.ViewMarketsRequest{
		UserRoles: req.UserRoles,
	})
	if err != nil {
		o.logger.Error(layer, method, "failed to get markets from SpotInstrumentService", err,
			"x_request_id", xReqID,
		)
		return domain.CreateOrderResponse{}, errs.ErrUnknown
	}

	marketExists := false
	for _, m := range marketsResp.Markets {
		if m.ID.String() == req.MarketID.String() {
			marketExists = true
			break
		}
	}

	if !marketExists {
		o.logger.Warn(layer, method, "market not found or not allowed",
			nil,
			"x_request_id", xReqID,
			"market_id", req.MarketID.String(),
			"user_roles", req.UserRoles,
		)
		return domain.CreateOrderResponse{}, errs.ErrMarketNotFound
	}

	resp, err := o.repo.CreateOrder(ctx, req)
	if err != nil {
		o.logger.Error(layer, method, "failed to create order", err,
			"x_request_id", xReqID,
		)
		return domain.CreateOrderResponse{}, errs.ErrUnknown
	}

	span.SetAttributes(
		attribute.String("order_id", resp.OrderID.String()),
		attribute.String("order_status", string(*resp.OrderStatus)),
	)

	o.logger.Info(layer, method, "order created successfully",
		"x_request_id", xReqID,
		"order_id", resp.OrderID.String(),
		"status", *resp.OrderStatus,
	)

	return resp, nil
}

func (o *orderUsecase) GetOrderStatus(
	ctx context.Context,
	req domain.GetOrderStatusRequest,
) (domain.GetOrderStatusResponse, error) {
	const layer = "usecase"
	const method = "GetOrderStatus"

	ctx, span := o.tracer.Start(ctx, "orderUsecase.GetOrderStatus")
	defer span.End()

	xRequestID := shared_context.XRequestIDFromContext(ctx)

	span.SetAttributes(
		attribute.String("x-request-id", xRequestID),
		attribute.String("order.id", req.OrderID.String()),
		attribute.String("user.id", req.UserID.String()),
	)

	o.logger.Info(layer, method, "fetching order status",
		"x_request_id", xRequestID,
		"order_id", req.OrderID.String(),
		"user_id", req.UserID.String(),
	)

	order, err := o.repo.GetOrderByID(ctx, *req.OrderID)
	if err != nil {
		span.RecordError(err)
		o.logger.Warn(layer, method, "order not found", err,
			"x_request_id", xRequestID,
			"order_id", req.OrderID.String(),
		)
		return domain.GetOrderStatusResponse{}, errs.ErrUnknown
	}

	if *order.UserID != *req.UserID {
		span.RecordError(err)
		span.SetAttributes(
			attribute.String("expected.user_id", order.UserID.String()),
			attribute.String("provided.user_id", req.UserID.String()),
		)

		o.logger.Warn(layer, method, "user ID mismatch", nil,
			"x_request_id", xRequestID,
			"expected_user_id", order.UserID.String(),
			"provided_user_id", req.UserID.String(),
		)
		return domain.GetOrderStatusResponse{}, errs.ErrInvalidUserID
	}

	o.logger.Info(layer, method, "order status retrieved",
		"x_request_id", xRequestID,
		"order_id", req.OrderID.String(),
		"status", *order.Status,
	)

	span.SetAttributes(attribute.String("order.status", string(*order.Status)))

	return domain.GetOrderStatusResponse{Status: order.Status}, nil
}

func (o *orderUsecase) SubscribeToOrderStatus(
	ctx context.Context,
	req domain.StreamOrderUpdatesRequest,
) (<-chan domain.StreamOrderUpdatesResponse, func(), error) {
	const layer = "usecase"
	const method = "SubscribeToOrderStatus"

	ctx, span := o.tracer.Start(ctx, "orderUsecase.SubscribeToOrderStatus")
	defer span.End()

	xRequestID := shared_context.XRequestIDFromContext(ctx)

	orderID := *req.OrderID
	userID := *req.UserID

	order, err := o.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		o.logger.Warn(layer, method, "order not found",
			nil,
			"x_request_id", xRequestID,
			"order_id", orderID.String(),
			"err", err,
		)
		return nil, nil, errs.ErrInvalidOrderID
	}

	if *order.UserID != *req.UserID {
		span.RecordError(err)
		span.SetAttributes(
			attribute.String("expected.user_id", order.UserID.String()),
			attribute.String("provided.user_id", req.UserID.String()),
		)

		o.logger.Warn(layer, method, "user ID mismatch", nil,
			"x_request_id", xRequestID,
			"expected_user_id", order.UserID.String(),
			"provided_user_id", req.UserID.String(),
		)
		return nil, nil, errs.ErrInvalidUserID
	}

	ch := make(chan domain.StreamOrderUpdatesResponse, 10)
	ctx, cancel := context.WithCancel(ctx)

	go o.streamOrderStatusUpdates(ctx, ch, req)

	o.logger.Info(layer, method, "started order status subscription",
		"x_request_id", xRequestID,
		"order_id", orderID.String(),
		"user_id", userID.String(),
	)

	return ch, cancel, nil
}

func (o *orderUsecase) streamOrderStatusUpdates(
	ctx context.Context,
	ch chan<- domain.StreamOrderUpdatesResponse,
	req domain.StreamOrderUpdatesRequest,
) {
	const layer = "usecase"
	const method = "streamOrderStatusUpdates"

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	defer close(ch)

	var lastStatus *domain.OrderStatusEnum

LOOP:
	for {
		select {
		case <-ctx.Done():
			o.logger.Info(layer, method, "subscription cancelled",
				"order_id", req.OrderID.String(),
				"user_id", req.UserID.String(),
			)
			break LOOP

		case <-ticker.C:
			order, err := o.repo.GetOrderByID(ctx, *req.OrderID)
			if err != nil {
				o.logger.Warn(layer, method, "failed to fetch order",
					nil,
					"order_id", req.OrderID.String(),
					"err", err,
				)
				continue
			}

			if *order.UserID != *req.UserID {
				o.logger.Warn(layer, method, "unauthorized update access attempt",
					nil,
					"order_id", req.OrderID.String(),
					"user_id", req.UserID.String(),
				)
				continue
			}

			if order.Status != nil && (lastStatus == nil || *order.Status != *lastStatus) {
				lastStatus = order.Status

				o.logger.Info(layer, method, "order status update streamed",
					"order_id", req.OrderID.String(),
					"user_id", req.UserID.String(),
					"status", *order.Status,
				)

				ch <- domain.StreamOrderUpdatesResponse{
					OrderID:     order.ID,
					OrderStatus: order.Status,
					UpdatedAt:   order.UpdatedAt,
				}
			}
		}
	}
}
