package cache

import (
	"context"
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
	cacher  cache.Cacher
}

// New to create new newrelic plugin for cache.
func New(d string, c cache.Cacher) cache.Cacher {
	return &client{
		dialect: d,
		cacher:  c,
	}
}

// Get to update get metrics.
func (c *client) Get(ctx context.Context, key string, data interface{}) error {
	segment := newrelic.DatastoreSegment{
		StartTime: newrelic.FromContext(ctx).StartSegmentNow(),
		Product:   newrelic.DatastoreProduct(c.dialect),
		Operation: cacheGet,
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
		StartTime: newrelic.FromContext(ctx).StartSegmentNow(),
		Product:   newrelic.DatastoreProduct(c.dialect),
		Operation: cacheSet,
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
		StartTime: newrelic.FromContext(ctx).StartSegmentNow(),
		Product:   newrelic.DatastoreProduct(c.dialect),
		Operation: cacheDelete,
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
