package utils_test

import (
	"github.com/rl404/akatsuki/pkg/cache"
	_cache "github.com/rl404/fairy/cache"
)

var cacheType = map[string]cache.CacheType{
	"nocache":  cache.NOP,
	"redis":    cache.Redis,
	"inmemory": cache.InMemory,
}

func GetCache(cfg *config) (_cache.Cacher, error) {
	return cache.New(cacheType[cfg.Cache.Dialect], cfg.Cache.Address, cfg.Cache.Password, cfg.Cache.Time)
}
