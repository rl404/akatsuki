package cache

import (
	"context"
	"errors"
	"time"

	"github.com/rl404/fairy/cache/bigcache"
	"github.com/rl404/fairy/cache/memcache"
	"github.com/rl404/fairy/cache/nocache"
	"github.com/rl404/fairy/cache/redis"
)

// Cacher is caching interface.
//
// See usage example in example folder.
type Cacher interface {
	// Get data from cache. The returned value will be
	// assigned to param `data`. Param `data` should
	// be a pointer just like when using json.Unmarshal.
	Get(ctx context.Context, key string, data interface{}) error
	// Save data to cache. Set and Get should be using
	// the same encoding method. For example, json.Marshal
	// for Set and json.Unmarshal for Get.
	Set(ctx context.Context, key string, data interface{}, ttl ...time.Duration) error
	// Delete data from cache.
	Delete(ctx context.Context, key string) error
	// Close cache connection.
	Close() error
}

// CacheType is type for cache.
type CacheType int8

// Available types for cache.
const (
	NoCache CacheType = iota
	InMemory
	Redis
	Memcache
)

// ErrInvalidCacheType is error for invalid cache type.
var ErrInvalidCacheType = errors.New("invalid cache type")

// New to create new cache client depends on the type.
func New(cacheType CacheType, address string, password string, expiredTime time.Duration) (Cacher, error) {
	switch cacheType {
	case NoCache:
		return nocache.New()
	case InMemory:
		return bigcache.New(expiredTime)
	case Redis:
		return redis.New(address, password, expiredTime)
	case Memcache:
		return memcache.New(address, expiredTime)
	default:
		return nil, ErrInvalidCacheType
	}
}
