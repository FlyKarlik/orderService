package in_memory_repo

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/FlyKarlik/orderService/internal/domain"
	shared_context "github.com/FlyKarlik/orderService/pkg/context"
	"github.com/FlyKarlik/orderService/pkg/logger"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type orderInMemoryRepo struct {
	mu     sync.RWMutex
	logger logger.Logger
	data   map[uuid.UUID]domain.Order
	tracer trace.Tracer
}

func NewInMemoryOrderRepository(l logger.Logger) *orderInMemoryRepo {
	return &orderInMemoryRepo{
		data:   make(map[uuid.UUID]domain.Order),
		logger: l,
		tracer: otel.Tracer("order-service/repo"),
	}
}

func SetupOrderRepo(l logger.Logger) *orderInMemoryRepo {
	orderRepo := NewInMemoryOrderRepository(l)
	orderRepo.StartStatusUpdater(context.Background())
	return orderRepo
}

func (r *orderInMemoryRepo) CreateOrder(
	ctx context.Context, req domain.CreateOrderRequest,
) (domain.CreateOrderResponse, error) {
	const layer = "repo"
	const method = "CreateOrder"

	ctx, span := r.tracer.Start(ctx, "OrderInMemoryRepo.CreateOrder")
	defer span.End()

	xRequestID := shared_context.XRequestIDFromContext(ctx)

	orderID := uuid.New()
	createdAt := time.Now()
	status := domain.OrderStatusEnumCreated

	order := domain.Order{
		ID:        &orderID,
		UserID:    req.UserID,
		MarketID:  req.MarketID,
		OrderType: req.OrderType,
		Price:     req.Price,
		Quantity:  req.Quantity,
		Status:    &status,
		CreatedAt: &createdAt,
	}

	r.mu.Lock()
	r.data[orderID] = order
	r.mu.Unlock()

	span.SetAttributes(
		attribute.String("x-request-id", xRequestID),
		attribute.String("order.id", orderID.String()),
		attribute.String("user.id", req.UserID.String()),
		attribute.String("market.id", req.MarketID.String()),
		attribute.String("order.status", string(*order.Status)),
		attribute.String("order.price", *order.Price),
		attribute.Int64("order.quantity", *req.Quantity),
	)

	r.logger.Info(layer, method, "order created",
		"x_request_id", xRequestID,
		"order_id", orderID.String(),
		"user_id", req.UserID.String(),
		"market_id", req.MarketID.String(),
		"price", req.Price,
		"quantity", req.Quantity,
		"status", status,
		"created_at", createdAt,
	)

	return domain.CreateOrderResponse{
		OrderID:     order.ID,
		OrderStatus: order.Status,
	}, nil
}

func (r *orderInMemoryRepo) GetOrderByID(ctx context.Context, ID uuid.UUID) (domain.Order, error) {
	const layer = "repo"
	const method = "GetOrderByID"

	ctx, span := r.tracer.Start(ctx, "OrderInMemoryRepo.GetOrderByID")
	defer span.End()

	xRequestID := shared_context.XRequestIDFromContext(ctx)

	span.SetAttributes(
		attribute.String("x-request-id", xRequestID),
		attribute.String("order.id", ID.String()),
	)

	r.mu.RLock()
	order, ok := r.data[ID]
	r.mu.RUnlock()

	if !ok {
		err := errors.New("order not found")
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("order.found", false))

		r.logger.Warn(layer, method, "order not found", nil,
			"x_request_id", xRequestID,
			"order_id", ID.String(),
		)
		return domain.Order{}, err
	}

	span.SetAttributes(
		attribute.Bool("order.found", true),
		attribute.String("user.id", order.UserID.String()),
		attribute.String("order.status", string(*order.Status)),
	)

	r.logger.Info(layer, method, "order retrieved",
		"x_request_id", xRequestID,
		"order_id", ID.String(),
		"user_id", order.UserID.String(),
		"status", *order.Status,
	)

	return order, nil
}

// Stub for test solution
func (r *orderInMemoryRepo) StartStatusUpdater(ctx context.Context) {
	const layer = "repo"
	const method = "StatusUpdater"

	ticker := time.NewTicker(30 * time.Second)

	go func() {
		for {
			select {
			case <-ctx.Done():
				r.logger.Info(layer, method, "status updater stopped")
				return
			case <-ticker.C:
				r.mu.Lock()
				for id, order := range r.data {
					if order.Status == nil {
						continue
					}

					if *order.Status == domain.OrderStatusEnumFilled || *order.Status == domain.OrderStatusEnumRejected {
						continue
					}

					var nextStatus domain.OrderStatusEnum
					switch *order.Status {
					case domain.OrderStatusEnumCreated:
						nextStatus = domain.OrderStatusEnumPending
					case domain.OrderStatusEnumPending:
						if rand.Intn(2) == 0 {
							nextStatus = domain.OrderStatusEnumFilled
						} else {
							nextStatus = domain.OrderStatusEnumRejected
						}
					default:
						continue
					}

					order.Status = &nextStatus
					updatedAt := time.Now()
					order.UpdatedAt = &updatedAt
					r.data[id] = order

					r.logger.Info(layer, method, "order status updated",
						"order_id", id.String(),
						"new_status", nextStatus,
					)
				}
				r.mu.Unlock()
			}
		}
	}()
}
