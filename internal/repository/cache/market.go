package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/FlyKarlik/orderService/internal/domain"
	"github.com/FlyKarlik/orderService/pkg/cache"
	"github.com/FlyKarlik/orderService/pkg/logger"
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type redisMarketsCache struct {
	logger logger.Logger
	client cache.RedisClient
	tracer trace.Tracer
}

func NewMarketsCache(logger logger.Logger, client cache.RedisClient) *redisMarketsCache {
	return &redisMarketsCache{
		logger: logger,
		client: client,
		tracer: otel.Tracer("order-service/cache"),
	}
}

func (c *redisMarketsCache) Set(ctx context.Context, key string, value domain.ViewMarketsResponse, ttl time.Duration) error {
	const method = "redisMarketsCache.Set"
	ctx, span := c.tracer.Start(ctx, method)
	defer span.End()

	data, err := json.Marshal(value)
	if err != nil {
		c.logger.Error("cache", method, "failed to marshal ViewMarketsResponse", err, "key", key)
		span.RecordError(err)
		return err
	}

	err = c.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		c.logger.Error("cache", method, "failed to set cache in Redis", err, "key", key)
		span.RecordError(err)
		return err
	}

	c.logger.Debug("cache", method, "cached ViewMarketsResponse successfully", "key", key, "ttl", ttl)
	span.SetAttributes(attribute.String("cache.key", key), attribute.String("cache.status", "set"))
	return nil
}

func (c *redisMarketsCache) Get(ctx context.Context, key string) (domain.ViewMarketsResponse, error) {
	const method = "redisMarketsCache.Get"
	ctx, span := c.tracer.Start(ctx, method)
	defer span.End()

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.logger.Debug("cache", method, "cache miss", "key", key)
			span.SetAttributes(attribute.String("cache.key", key), attribute.String("cache.status", "miss"))
			return domain.ViewMarketsResponse{}, nil
		}
		c.logger.Error("cache", method, "failed to get from Redis", err, "key", key)
		span.RecordError(err)
		return domain.ViewMarketsResponse{}, err
	}

	var result domain.ViewMarketsResponse
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		c.logger.Error("cache", method, "failed to unmarshal cached data", err, "key", key)
		span.RecordError(err)
		return domain.ViewMarketsResponse{}, err
	}

	c.logger.Debug("cache", method, "cache hit", "key", key)
	span.SetAttributes(attribute.String("cache.key", key), attribute.String("cache.status", "hit"))
	return result, nil
}

func (c *redisMarketsCache) Delete(ctx context.Context, key string) error {
	const method = "redisMarketsCache.Delete"
	ctx, span := c.tracer.Start(ctx, method)
	defer span.End()

	err := c.client.Del(ctx, key).Err()
	if err != nil {
		c.logger.Error("cache", method, "failed to delete key from Redis", err, "key", key)
		span.RecordError(err)
		return err
	}

	c.logger.Debug("cache", method, "deleted key from cache", "key", key)
	span.SetAttributes(attribute.String("cache.key", key), attribute.String("cache.status", "deleted"))
	return nil
}
