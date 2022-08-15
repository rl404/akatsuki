package client

import (
	"net/http"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/nstratos/go-myanimelist/mal"
	"github.com/rl404/fairy/limit"
)

// Client is client for mal.
type Client struct {
	client  *mal.Client
	limiter limit.Limiter
}

// New to create new mal client.
func New(clientID string) *Client {
	limiter, _ := limit.New(limit.Atomic, 1, time.Second)
	return &Client{
		client: mal.NewClient(&http.Client{
			Transport: &clientIDTransport{
				clientID: clientID,
			},
		}),
		limiter: limiter,
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
	c.transport = newrelic.NewRoundTripper(c.transport)
	req.Header.Add("X-MAL-CLIENT-ID", c.clientID)
	return c.transport.RoundTrip(req)
}
