package cache

import (
	"errors"
	"time"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/fairy/cache/inmemory"
	"github.com/rl404/fairy/cache/nop"
	"github.com/rl404/fairy/cache/redis"
)

// CacheType is type for cache.
type CacheType int8

// Available types for cache.
const (
	NOP CacheType = iota
	InMemory
	Redis
)

// ErrInvalidCacheType is error for invalid cache type.
var ErrInvalidCacheType = errors.New("invalid cache type")

// New to create new cache client depends on the type.
func New(cacheType CacheType, address string, password string, expiredTime time.Duration) (cache.Cacher, error) {
	switch cacheType {
	case NOP:
		return nop.New()
	case InMemory:
		return inmemory.New(expiredTime)
	case Redis:
		return redis.New(address, password, expiredTime)
	default:
		return nil, ErrInvalidCacheType
	}
}
