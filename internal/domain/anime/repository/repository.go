package repository

import (
	"context"

	"github.com/rl404/akatsuki/internal/domain/anime/entity"
)

// Repository contains functions for anime domain.
type Repository interface {
	IsDataOld(ctx context.Context, id int64) (bool, int, error)
	GetByID(ctx context.Context, id int64) (*entity.Anime, int, error)
	GetByIDs(ctx context.Context, ids []int64) ([]*entity.Anime, int, error)
	Update(ctx context.Context, data entity.Anime) (int, error)
}
