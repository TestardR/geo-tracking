package redis_cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	redis *redis.Client
}

func NewClient(redis *redis.Client) *Client {
	return &Client{redis: redis}
}

func (c *Client) Set(ctx context.Context, key string, data []byte, expiration time.Duration) error {
	result := c.redis.Set(ctx, key, data, expiration)
	return result.Err()
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	result := c.redis.Get(ctx, key)
	resultBytes, err := result.Bytes()

	return resultBytes, err
}

func (c *Client) Delete(ctx context.Context, key string) error {
	result := c.redis.Del(ctx, key)

	return result.Err()
}
