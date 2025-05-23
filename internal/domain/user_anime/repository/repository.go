package repository

import (
	"context"

	"github.com/rl404/akatsuki/internal/domain/user_anime/entity"
)

// Repository contains functions for user_anime domain.
type Repository interface {
	Get(ctx context.Context, data entity.GetUserAnimeRequest) ([]*entity.UserAnime, int, int, error)
	Update(ctx context.Context, data entity.UserAnime) (int, error)
	IsOld(ctx context.Context, username string) (bool, int, error)
	GetOldUsernames(ctx context.Context) ([]string, int, error)
	DeleteNotInList(ctx context.Context, username string, ids []int64, status entity.Status) (int, error)
	DeleteByAnimeID(ctx context.Context, animeID int64) (int, error)
	DeleteByUsername(ctx context.Context, username string) (int, error)
}
