package cacheman

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
	ttl    time.Duration
}

// NewRedis creates big cache client
func NewRedis(config *Config) (*RedisClient, error) {
	ttl, e := time.ParseDuration(config.TTL)
	if e != nil {
		ttl, _ = time.ParseDuration(defaultTTL)
	}
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     config.Server,
		Password: config.Password,
		DB:       config.Database.(int),
	})
	return &RedisClient{
		client: client,
		ctx:    ctx,
		ttl:    ttl,
	}, nil
}

func (c *RedisClient) Get(key string) ([]byte, error) {
	result, e := c.client.Get(c.ctx, key).Result()
	if e != nil {
		return nil, e
	}
	return []byte(result), nil
}

func (c *RedisClient) Set(key string, value []byte) error {
	return c.client.Set(c.ctx, key, string(value), c.ttl).Err()
}

func (c *RedisClient) Delete(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

func (c *RedisClient) Reset() error {
	return c.client.FlushAll(c.ctx).Err()
}

func (c *RedisClient) Type() string {
	return fmt.Sprintf("%T", c)
}
