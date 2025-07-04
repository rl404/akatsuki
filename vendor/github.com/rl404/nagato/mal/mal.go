package mal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rl404/nagato/internal/limiter"
	"golang.org/x/oauth2"
)

const host = "https://api.myanimelist.net/v2"
const version = "0.4.0"

// Client is myanimelist API client.
//
// Override the values if necessary after
// initiated.
type Client struct {
	Host    string
	Http    *http.Client
	Limiter limiter.Limiter
}

// New to create new default client for myanimelist API.
//
// You can't use the client directly since you need to
// provide client id in the header. You need to override
// the http client to follow myanimelist requirement.
//
// Use NewPublic() or NewWithOauth2() for easier to use.
func New() *Client {
	return &Client{
		Host:    host,
		Http:    http.DefaultClient,
		Limiter: limiter.New(1, time.Second),
	}
}

// NewPublic to create new client for myanimelist API
// to retrieve public info only.
func NewPublic(clientID string) *Client {
	return &Client{
		Host: host,
		Http: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &publicTransport{
				clientID: clientID,
			},
		},
		Limiter: limiter.New(1, time.Second),
	}
}

// Oauth2Config is oauth 2 config to
// initiate client with oauth 2.
type Oauth2Config struct {
	ClientID     string
	ClientSecret string
	State        string
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
}

// NewWithOauth2 to create new client for myanimelist API
// with oauth 2. Token should be generated before using this.
//
// See in the `example/mal` for the example of generating the token.
func NewWithOauth2(cfg Oauth2Config) (*Client, error) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &publicTransport{
			clientID: cfg.ClientID,
		},
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)

	oauth2Cfg := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://myanimelist.net/v1/oauth2/authorize",
			TokenURL:  "https://myanimelist.net/v1/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	oauth2Token := oauth2.Token{
		AccessToken:  cfg.AccessToken,
		RefreshToken: cfg.RefreshToken,
		TokenType:    "Bearer",
		Expiry:       cfg.Expiry,
	}

	return &Client{
		Host:    host,
		Http:    oauth2Cfg.Client(ctx, &oauth2Token),
		Limiter: limiter.New(1, time.Second),
	}, nil
}

type publicTransport struct {
	transport http.RoundTripper
	clientID  string
}

// RoundTrip is http round tripper with mal client id.
func (p *publicTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if p.transport == nil {
		p.transport = http.DefaultTransport
	}
	req.Header.Add("X-MAL-CLIENT-ID", p.clientID)
	req.Header.Add("User-Agent", fmt.Sprintf("Nagato/%s (github.com/rl404/nagato)", version))
	return p.transport.RoundTrip(req)
}
