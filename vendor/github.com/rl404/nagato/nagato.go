package nagato

import (
	"net/http"

	"github.com/rl404/nagato/internal/limiter"
	"github.com/rl404/nagato/internal/playground"
	"github.com/rl404/nagato/mal"
)

// Client is nagato client.
type Client struct {
	mal       *mal.Client
	validator *playground.Validator
}

// New to create new default nagato client.
func New(clientID string) *Client {
	c := Client{
		mal:       mal.NewPublic(clientID),
		validator: playground.New(true),
	}

	c.initValidator()

	return &c
}

// SetMalClient to override mal client.
func (c *Client) SetMalClient(mal *mal.Client) {
	c.mal = mal
}

// SetHttpClient to override http client used
// by mal client.
func (c *Client) SetHttpClient(http *http.Client) {
	c.mal.Http = http
}

// SetLimiter to override limiter in mal client.
func (c *Client) SetLimiter(limiter limiter.Limiter) {
	c.mal.Limiter = limiter
}
