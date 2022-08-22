package cache

import (
	"context"
	"strings"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/cache"
)

const (
	cacheGet    = "GET"
	cacheSet    = "SET"
	cacheDelete = "DELETE"
	cacheHit    = "HIT"
	cacheMiss   = "MISS"
)

type client struct {
	dialect string
	host    string
	port    string
	cacher  cache.Cacher
}

// New to create new newrelic plugin for cache.
func New(dialect string, address string, cacher cache.Cacher) cache.Cacher {
	host, port := splitHostPort(address)
	return &client{
		dialect: dialect,
		host:    host,
		port:    port,
		cacher:  cacher,
	}
}

// Get to update get metrics.
func (c *client) Get(ctx context.Context, key string, data interface{}) error {
	segment := newrelic.DatastoreSegment{
		StartTime:    newrelic.FromContext(ctx).StartSegmentNow(),
		Product:      newrelic.DatastoreProduct(c.dialect),
		Operation:    cacheGet,
		Host:         c.host,
		PortPathOrID: c.port,
	}
	defer segment.End()

	if err := c.cacher.Get(ctx, key, data); err != nil {
		segment.ParameterizedQuery = cacheMiss
		return err
	}
	segment.ParameterizedQuery = cacheHit
	return nil
}

// Set to update set metrics.
func (c *client) Set(ctx context.Context, key string, data interface{}, ttl ...time.Duration) error {
	segment := newrelic.DatastoreSegment{
		StartTime:    newrelic.FromContext(ctx).StartSegmentNow(),
		Product:      newrelic.DatastoreProduct(c.dialect),
		Operation:    cacheSet,
		Host:         c.host,
		PortPathOrID: c.port,
	}
	defer segment.End()

	if err := c.cacher.Set(ctx, key, data, ttl...); err != nil {
		segment.ParameterizedQuery = cacheMiss
		return err
	}
	segment.ParameterizedQuery = cacheHit
	return nil
}

// Delete to update delete metrics.
func (c *client) Delete(ctx context.Context, key string) error {
	segment := newrelic.DatastoreSegment{
		StartTime:    newrelic.FromContext(ctx).StartSegmentNow(),
		Product:      newrelic.DatastoreProduct(c.dialect),
		Operation:    cacheDelete,
		Host:         c.host,
		PortPathOrID: c.port,
	}
	defer segment.End()

	if err := c.cacher.Delete(ctx, key); err != nil {
		segment.ParameterizedQuery = cacheMiss
		return err
	}
	segment.ParameterizedQuery = cacheHit
	return nil
}

// Close to close.
func (c *client) Close() error {
	return c.cacher.Close()
}

func splitHostPort(address string) (host string, port string) {
	host = address

	colon := strings.LastIndexByte(host, ':')
	if colon != -1 && validOptionalPort(host[colon:]) {
		host, port = host[:colon], host[colon+1:]
	}

	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = host[1 : len(host)-1]
	}

	return
}

func validOptionalPort(port string) bool {
	if port == "" {
		return true
	}
	if port[0] != ':' {
		return false
	}
	for _, b := range port[1:] {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}
