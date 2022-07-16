package cache

import (
	"context"
	"net/http"

	"github.com/rl404/akatsuki/internal/domain/anime/entity"
	"github.com/rl404/akatsuki/internal/domain/anime/repository"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/utils"
	"github.com/rl404/fairy/cache"
)

// Cache contains functions for anime cache.
type Cache struct {
	cacher cache.Cacher
	repo   repository.Repository
}

// New to create new anime cache.
func New(cacher cache.Cacher, repo repository.Repository) *Cache {
	return &Cache{
		cacher: cacher,
		repo:   repo,
	}
}

// GetByID to get anime by id.
func (c *Cache) GetByID(ctx context.Context, id int64) (data *entity.Anime, code int, err error) {
	key := utils.GetKey("anime", id)
	if c.cacher.Get(key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetByID(ctx, id)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return data, code, nil
}

// GetByIDs to get anime by ids.
func (c *Cache) GetByIDs(ctx context.Context, ids []int64) ([]*entity.Anime, int, error) {
	return c.repo.GetByIDs(ctx, ids)
}

// IsDataOld to check if data is old.
func (c *Cache) IsDataOld(ctx context.Context, id int64) (bool, int, error) {
	return c.repo.IsDataOld(ctx, id)
}

// Update to update data.
func (c *Cache) Update(ctx context.Context, data entity.Anime) (int, error) {
	return c.repo.Update(ctx, data)
}
