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

// Update to update data.
func (c *Cache) Update(ctx context.Context, data entity.Anime) (int, error) {
	if code, err := c.repo.Update(ctx, data); err != nil {
		return code, errors.Wrap(ctx, err)
	}

	key := utils.GetKey("anime", data.ID)
	if err := c.cacher.Delete(key); err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return http.StatusOK, nil
}

// IsOld to check if old.
func (c *Cache) IsOld(ctx context.Context, id int64) (bool, int, error) {
	return c.repo.IsOld(ctx, id)
}

// GetOldReleasingIDs to get old releasing anime ids.
func (c *Cache) GetOldReleasingIDs(ctx context.Context) ([]int64, int, error) {
	return c.repo.GetOldReleasingIDs(ctx)
}

// GetOldFinishedIDs to get old finished anime ids.
func (c *Cache) GetOldFinishedIDs(ctx context.Context) ([]int64, int, error) {
	return c.repo.GetOldFinishedIDs(ctx)
}

// GetOldNotYetIDs to get old not yet released anime ids.
func (c *Cache) GetOldNotYetIDs(ctx context.Context) ([]int64, int, error) {
	return c.repo.GetOldNotYetIDs(ctx)
}

// GetMaxID to get max id.
func (c *Cache) GetMaxID(ctx context.Context) (int64, int, error) {
	return c.repo.GetMaxID(ctx)
}

// GetIDs to get all ids.
func (c *Cache) GetIDs(ctx context.Context) ([]int64, int, error) {
	return c.repo.GetIDs(ctx)
}
