package redis_helper

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type IRedisCache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (*string, error)
}

type RedisCache struct {
	rdb       *redis.Client
	keyPrefix string
}

func (r RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	value, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.rdb.Set(ctx, fmt.Sprintf("%s:%s", r.keyPrefix, key), value, expiration).Err()
}

func (r RedisCache) Get(ctx context.Context, key string) (*string, error) {
	value, err := r.rdb.Get(ctx, fmt.Sprintf("%s:%s", r.keyPrefix, key)).Result()
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func NewRedisCache(rdb *redis.Client, keyPrefix string) IRedisCache {
	return RedisCache{
		rdb:       rdb,
		keyPrefix: keyPrefix,
	}
}
