package cache

import (
	"context"

	"github.com/rl404/akatsuki/internal/domain/genre/entity"
	"github.com/rl404/akatsuki/internal/domain/genre/repository"
	"github.com/rl404/fairy/cache"
)

// Cache contains functions for genre cache.
type Cache struct {
	cacher cache.Cacher
	repo   repository.Repository
}

// New to create new genre cache.
func New(cacher cache.Cacher, repo repository.Repository) *Cache {
	return &Cache{
		cacher: cacher,
		repo:   repo,
	}
}

// BatchUpdate to batch update.
func (c *Cache) BatchUpdate(ctx context.Context, data []entity.Genre) (int, error) {
	return c.repo.BatchUpdate(ctx, data)
}

// GetByIDs to get genre by ids.
func (c *Cache) GetByIDs(ctx context.Context, ids []int64) ([]*entity.Genre, int, error) {
	return c.repo.GetByIDs(ctx, ids)
}

// Get to get list.
func (c *Cache) Get(ctx context.Context, data entity.GetRequest) ([]*entity.Genre, int, int, error) {
	return c.repo.Get(ctx, data)
}

// GetByID to get by id.
func (c *Cache) GetByID(ctx context.Context, id int64) (*entity.Genre, int, error) {
	return c.repo.GetByID(ctx, id)
}
