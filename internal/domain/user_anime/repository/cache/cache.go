package cache

import (
	"context"

	"github.com/rl404/akatsuki/internal/domain/user_anime/entity"
	"github.com/rl404/akatsuki/internal/domain/user_anime/repository"
	"github.com/rl404/fairy/cache"
)

// Cache contains functions for user cache.
type Cache struct {
	cacher cache.Cacher
	repo   repository.Repository
}

// New to create new user cache.
func New(cacher cache.Cacher, repo repository.Repository) *Cache {
	return &Cache{
		cacher: cacher,
		repo:   repo,
	}
}

// Get to get user anime.
func (c *Cache) Get(ctx context.Context, data entity.GetUserAnimeRequest) ([]*entity.UserAnime, int, int, error) {
	return c.repo.Get(ctx, data)
}

// Update to update user anime.
func (c *Cache) Update(ctx context.Context, data entity.UserAnime) (int, error) {
	return c.repo.Update(ctx, data)
}

// IsOld to check if old.
func (c *Cache) IsOld(ctx context.Context, username string) (bool, int, error) {
	return c.repo.IsOld(ctx, username)
}

// GetOldUsernames to get old username.
func (c *Cache) GetOldUsernames(ctx context.Context) ([]string, int, error) {
	return c.repo.GetOldUsernames(ctx)
}

// DeleteNotInList to delete anime not in list.
func (c *Cache) DeleteNotInList(ctx context.Context, username string, ids []int64) (int, error) {
	return c.repo.DeleteNotInList(ctx, username, ids)
}

// DeleteByAnimeID to delete by anime id.
func (c *Cache) DeleteByAnimeID(ctx context.Context, animeID int64) (int, error) {
	return c.repo.DeleteByAnimeID(ctx, animeID)
}
