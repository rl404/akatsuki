// Package memcache is a wrapper of the original "github.com/bradfitz/gomemcache/memcache" library.
//
// Only contains basic get, set, delete, and close methods.
// Data will be encoded to JSON before saving to cache.
package memcache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// Client is memcache client.
type Client struct {
	client      *memcache.Client
	expiredTime time.Duration
}

// New to create new cache.
func New(address string, expiredTime time.Duration) (*Client, error) {
	c := memcache.New(address)
	return NewFromGoMemCache(c, expiredTime), c.Ping()
}

// NewFromGoMemCache to create new cache from gomemcache.
func NewFromGoMemCache(client *memcache.Client, expiredTime time.Duration) *Client {
	return &Client{client: client, expiredTime: expiredTime}
}

// Set to save data to cache.
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

	return c.client.Set(&memcache.Item{
		Key:        key,
		Value:      d,
		Expiration: int32(expiredTime.Seconds()),
	})
}

// Get to get data from cache.
func (c *Client) Get(ctx context.Context, key string, data interface{}) error {
	d, err := c.client.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(d.Value, &data)
}

// Delete to delete data from cache.
func (c *Client) Delete(ctx context.Context, key string) error {
	err := c.client.Delete(key)
	if errors.Is(err, memcache.ErrCacheMiss) {
		return nil
	}
	return err
}

// Close to close cache connection.
func (c *Client) Close() error {
	// gomemcache has no function to close client connection. :(
	return nil
}
