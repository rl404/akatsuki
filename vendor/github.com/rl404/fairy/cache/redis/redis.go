// Package redis is a wrapper of the original "github.com/redis/go-redis" library.
//
// Only contains basic get, set, delete, and close methods.
//
// Data will be encoded to JSON before saving to cache.
package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client is redis client.
type Client struct {
	client      *redis.Client
	expiredTime time.Duration
}

// New to create cache cache with default config.
func New(address, password string, expiredTime time.Duration) (*Client, error) {
	return NewWithConfig(redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	}, expiredTime)
}

// NewWithConfig to create cache from go-redis options.
func NewWithConfig(option redis.Options, expiredTime time.Duration) (*Client, error) {
	client := redis.NewClient(&option)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping test.
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return NewFromGoRedis(client, expiredTime), nil
}

// NewFromGoRedis to create cache from go-redis client.
func NewFromGoRedis(client *redis.Client, expiredTime time.Duration) *Client {
	return &Client{
		client:      client,
		expiredTime: expiredTime,
	}
}

// Set to save data to cache,
func (c *Client) Set(ctx context.Context, key string, data interface{}, ttl ...time.Duration) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Override ttl.
	expiredTime := c.expiredTime
	if len(ttl) > 0 {
		expiredTime = ttl[0]
	}

	return c.client.Set(ctx, key, d, expiredTime).Err()
}

// Get to get data from cache.
func (c *Client) Get(ctx context.Context, key string, data interface{}) error {
	d, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(d), &data)
}

// Delete to delete data from cache.
func (c *Client) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Close to close cache connection.
func (c *Client) Close() error {
	return c.client.Close()
}
