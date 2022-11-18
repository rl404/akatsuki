package client

import (
	"net/http"
	"time"

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
		Timeout: 10 * time.Second,
		Transport: &clientIDTransport{
			clientID: clientID,
		},
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
	// c.transport = newrelic.NewRoundTripper(c.transport)
	req.Header.Add("X-MAL-CLIENT-ID", c.clientID)
	return c.transport.RoundTrip(req)
}
