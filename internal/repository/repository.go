package repository

import (
	"context"
	"time"

	"github.com/FlyKarlik/orderService/internal/domain"
	redis_cache "github.com/FlyKarlik/orderService/internal/repository/cache"
	in_memory_repo "github.com/FlyKarlik/orderService/internal/repository/in_memory"
	"github.com/FlyKarlik/orderService/pkg/cache"
	"github.com/FlyKarlik/orderService/pkg/logger"
	"github.com/google/uuid"
)

type IOrderRepository interface {
	CreateOrder(ctx context.Context, req domain.CreateOrderRequest) (domain.CreateOrderResponse, error)
	GetOrderByID(ctx context.Context, ID uuid.UUID) (domain.Order, error)
}

type IMarketsCache interface {
	Set(ctx context.Context, key string, value domain.ViewMarketsResponse, ttl time.Duration) error
	Get(ctx context.Context, key string) (domain.ViewMarketsResponse, error)
	Delete(ctx context.Context, key string) error
}

type Repository interface {
	IOrderRepository
	IMarketsCache
}

type repositoryImpl struct {
	IOrderRepository
	IMarketsCache
}

func New(l logger.Logger, redisClient cache.RedisClient) *repositoryImpl {
	return &repositoryImpl{
		IOrderRepository: in_memory_repo.SetupOrderRepo(l),
		IMarketsCache:    redis_cache.NewMarketsCache(l, redisClient),
	}
}
