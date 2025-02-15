package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client  *redis.Client
	options Options
}

func NewRedisCache(opts Options, redisOpts *redis.Options) (*RedisCache, error) {
	client := redis.NewClient(redisOpts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %v", err)
	}

	return &RedisCache{
		client:  client,
		options: opts,
	}, nil
}

func (c *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	key = c.prefixKey(key)

	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, ErrKeyNotFound
	}
	if err != nil {
		return nil, &CacheError{Op: "get", Key: key, Err: err}
	}

	var result interface{}
	if err := c.deserialize([]byte(val), &result); err != nil {
		return nil, &CacheError{Op: "deserialize", Key: key, Err: err}
	}

	return result, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	key = c.prefixKey(key)

	data, err := c.serialize(value)
	if err != nil {
		return &CacheError{Op: "serialize", Key: key, Err: err}
	}

	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return &CacheError{Op: "set", Key: key, Err: err}
	}

	return nil
}
