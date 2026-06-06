package core_redis

import (
	core_config "api/internal/core/config"
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	rdb *redis.Client
}

func NewRedisClient(
	cfg *core_config.Config,
) *RedisClient {
	return &RedisClient{
		rdb: redis.NewClient(&redis.Options{
			Addr:     cfg.RedisHost,
			Password: cfg.RedisPassword,
			DB:       0,
		}),
	}
}

func (r *RedisClient) Set(
	ctx context.Context,
	key string,
	value interface{},
	ttl time.Duration,
) error {
	return r.rdb.Set(ctx, key, value, ttl).Err()
}

func (r *RedisClient) Get(
	ctx context.Context,
	key string,
) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}

func (r *RedisClient) Del(
	ctx context.Context,
	key string,
) error {
	return r.rdb.Del(ctx, key).Err()
}
