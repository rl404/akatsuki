package repository

import (
	"context"

	"github.com/rl404/akatsuki/internal/domain/user_anime/entity"
)

// Repository contains functions for user domain.
type Repository interface {
	Get(ctx context.Context, data entity.GetUserAnimeRequest) ([]*entity.UserAnime, int, int, error)
	Update(ctx context.Context, data entity.UserAnime) (int, error)
	IsOld(ctx context.Context, username string) (bool, int, error)
}
