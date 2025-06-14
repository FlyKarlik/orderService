package usecase

import (
	"context"

	"github.com/FlyKarlik/orderService/internal/domain"
	"github.com/FlyKarlik/orderService/internal/driver"
	"github.com/FlyKarlik/orderService/internal/repository"
	"github.com/FlyKarlik/orderService/pkg/logger"
)

type IOrderUsecase interface {
	CreateOrder(ctx context.Context, req domain.CreateOrderRequest) (domain.CreateOrderResponse, error)
	GetOrderStatus(ctx context.Context, req domain.GetOrderStatusRequest) (domain.GetOrderStatusResponse, error)
	SubscribeToOrderStatus(ctx context.Context, req domain.StreamOrderUpdatesRequest) (<-chan domain.StreamOrderUpdatesResponse, func(), error)
}

type Usecase interface {
	IOrderUsecase
}

type usecaseImpl struct {
	IOrderUsecase
}

func New(logger logger.Logger, driver driver.Driver, repo repository.Repository) *usecaseImpl {
	return &usecaseImpl{
		IOrderUsecase: newOrderUsecase(logger, driver, repo),
	}
}
