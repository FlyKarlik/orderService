package repository

import (
	"context"

	"github.com/FlyKarlik/orderService/internal/domain"
	in_memory_repo "github.com/FlyKarlik/orderService/internal/repository/in_memory"
	"github.com/FlyKarlik/orderService/pkg/logger"
	"github.com/google/uuid"
)

type IOrderRepository interface {
	CreateOrder(ctx context.Context, req domain.CreateOrderRequest) (domain.CreateOrderResponse, error)
	GetOrderByID(ctx context.Context, ID uuid.UUID) (domain.Order, error)
}

type Repository interface {
	IOrderRepository
}

type repositoryImpl struct {
	IOrderRepository
}

func New(l logger.Logger) *repositoryImpl {
	return &repositoryImpl{
		IOrderRepository: in_memory_repo.SetupOrderRepo(l),
	}
}
