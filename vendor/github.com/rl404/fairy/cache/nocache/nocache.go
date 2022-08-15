// Package nocache is a mock of caching.
package nocache

import (
	"context"
	"errors"
	"time"
)

// Client is nocache client.
type Client struct{}

// ErrNoCache is default Get error return.
var ErrNoCache = errors.New("not using cache")

// New to create fake cache.
func New() (*Client, error) {
	return &Client{}, nil
}

// Set will just return nil.
func (c *Client) Set(_ context.Context, _ string, _ interface{}, _ ...time.Duration) error {
	return nil
}

// Get will just return error to simulate as if data is not
// in cache.
func (c *Client) Get(_ context.Context, _ string, _ interface{}) error {
	return ErrNoCache
}

// Delete will just return nil.
func (c *Client) Delete(_ context.Context, _ string) error {
	return nil
}

// Close will just return nil.
func (c *Client) Close() error {
	return nil
}
