package cache

import (
	"context"
	"net/http"

	"github.com/rl404/akatsuki/internal/domain/empty_id/repository"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/utils"
	"github.com/rl404/fairy/cache"
)

// Cache contains functions for empty_id cache.
type Cache struct {
	cacher cache.Cacher
	repo   repository.Repository
}

// New to create new empty_id cache.
func New(cacher cache.Cacher, repo repository.Repository) *Cache {
	return &Cache{
		cacher: cacher,
		repo:   repo,
	}
}

// IsEmpty to check if id is empty.
func (c *Cache) IsEmpty(ctx context.Context, id int64) (bool, int, error) {
	key := utils.GetKey("empty-id", id)
	var data bool
	if c.cacher.Get(key, &data) == nil {
		return data, http.StatusOK, nil
	}

	isEmpty, code, err := c.repo.IsEmpty(ctx, id)
	if err != nil {
		return true, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(key, isEmpty); err != nil {
		return true, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return isEmpty, code, nil
}

// Create to create empty id.
func (c *Cache) Create(ctx context.Context, id int64) (int, error) {
	key := utils.GetKey("empty-id", id)
	if code, err := c.repo.Create(ctx, id); err != nil {
		return code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(key, true); err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return http.StatusCreated, nil
}

// Delete to delete empty id.
func (c *Cache) Delete(ctx context.Context, id int64) (int, error) {
	key := utils.GetKey("empty-id", id)
	if code, err := c.repo.Delete(ctx, id); err != nil {
		return code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Delete(key); err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return http.StatusOK, nil
}
