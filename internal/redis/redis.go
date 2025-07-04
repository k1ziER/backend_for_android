package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisClient(addr, password string, db int, tl time.Duration) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &Client{client: rdb, ttl: tl}
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	return c.client.Get(ctx, key).Bytes()
}

func (c *Client) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	if expiration == 0 {
		expiration = c.ttl
	}
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *Client) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *Client) Close() error {
	return c.client.Close()
}
