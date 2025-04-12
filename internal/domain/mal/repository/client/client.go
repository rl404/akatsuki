package client

import (
	"net/http"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/limit/atomic"
	"github.com/rl404/nagato"
)

// Client is client for mal.
type Client struct {
	client *nagato.Client
}

// New to create new mal client.
func New(clientID string) *Client {
	c := nagato.New(clientID)
	c.SetLimiter(atomic.New(1, time.Second))
	c.SetHttpClient(&http.Client{
		Timeout: 30 * time.Second,
		Transport: newrelic.NewRoundTripper(&clientIDTransport{
			clientID: clientID,
		}),
	})
	return &Client{
		client: c,
	}
}

type clientIDTransport struct {
	transport http.RoundTripper
	clientID  string
}

func (c *clientIDTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if c.transport == nil {
		c.transport = http.DefaultTransport
	}
	req.Header.Add("X-MAL-CLIENT-ID", c.clientID)
	req.Header.Add("User-Agent", "Akatsuki/0.12.2 (github.com/rl404/akatsuki)")
	return c.transport.RoundTrip(req)
}
