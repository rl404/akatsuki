package repository

import (
	"context"

	"github.com/rl404/akatsuki/internal/domain/anime/entity"
)

// Repository contains functions for anime domain.
type Repository interface {
	GetByID(ctx context.Context, id int64) (*entity.Anime, int, error)
	GetByIDs(ctx context.Context, ids []int64) ([]*entity.Anime, int, error)
	Update(ctx context.Context, data entity.Anime) (int, error)

	IsOld(ctx context.Context, id int64) (bool, int, error)
	GetOldFinished(ctx context.Context, limit int) ([]*entity.Anime, int, error)
	GetOldReleasing(ctx context.Context, limit int) ([]*entity.Anime, int, error)
	GetOldNotYet(ctx context.Context, limit int) ([]*entity.Anime, int, error)
	GetMaxID(ctx context.Context) (int64, int, error)
	GetIDs(ctx context.Context) ([]int64, int, error)
}
