package cache

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string, target any) error
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

func UseCache[T any](ctx context.Context, cash Cache, key string, ttl time.Duration, callback func() (T, error)) (T, error) {
	var v T
	err := cash.Get(ctx, key, &v)
	if err != cache.ErrCacheMiss {
		return v, err
	}

	v, err = callback()
	if err != nil {
		return v, err
	}

	// fire and forget
	//nolint:errcheck
	cash.Set(ctx, key, v, ttl)
	return v, nil
}

type CacheRedis struct {
	instance *cache.Cache
}

func (c *CacheRedis) Get(ctx context.Context, key string, target any) error {
	return c.instance.Get(ctx, key, target)
}

func (c *CacheRedis) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return c.instance.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   ttl,
	})
}

func (c *CacheRedis) Delete(ctx context.Context, key string) error {
	return c.instance.Delete(ctx, key)
}

func NewCacheRedis(client *redis.Client) (*CacheRedis, error) {
	return &CacheRedis{cache.New(&cache.Options{
		Redis:      client,
		LocalCache: cache.NewTinyLFU(10000, time.Minute),
	})}, nil
}
