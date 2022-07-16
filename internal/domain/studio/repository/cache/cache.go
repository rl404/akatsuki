package cache

import (
	"context"

	"github.com/rl404/akatsuki/internal/domain/studio/entity"
	"github.com/rl404/akatsuki/internal/domain/studio/repository"
	"github.com/rl404/fairy/cache"
)

// Cache contains functions for studio cache.
type Cache struct {
	cacher cache.Cacher
	repo   repository.Repository
}

// New to create new studio cache.
func New(cacher cache.Cacher, repo repository.Repository) *Cache {
	return &Cache{
		cacher: cacher,
		repo:   repo,
	}
}

// BatchUpdate to batch update.
func (c *Cache) BatchUpdate(ctx context.Context, data []entity.Studio) (int, error) {
	return c.repo.BatchUpdate(ctx, data)
}

// GetByIDs to get studio by ids.
func (c *Cache) GetByIDs(ctx context.Context, ids []int64) ([]*entity.Studio, int, error) {
	return c.repo.GetByIDs(ctx, ids)
}
