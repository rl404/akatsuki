package repository

import (
	"context"

	"github.com/rl404/akatsuki/internal/domain/anime/entity"
)

// Repository contains functions for anime domain.
type Repository interface {
	GetByID(ctx context.Context, id int64) (*entity.Anime, int, error)
	GetByIDs(ctx context.Context, ids []int64) ([]*entity.Anime, int, error)
	GetHistories(ctx context.Context, data entity.GetHistoriesRequest) ([]entity.History, int, error)
	Update(ctx context.Context, data entity.Anime) (int, error)
	GetRelatedByIDs(ctx context.Context, ids []int64) ([]*entity.AnimeRelated, int, error)
	DeleteByID(ctx context.Context, id int64) (int, error)

	IsOld(ctx context.Context, id int64) (bool, int, error)
	GetMaxID(ctx context.Context) (int64, int, error)
	GetIDs(ctx context.Context) ([]int64, int, error)
	GetOldFinishedIDs(ctx context.Context) ([]int64, int, error)
	GetOldReleasingIDs(ctx context.Context) ([]int64, int, error)
	GetOldNotYetIDs(ctx context.Context) ([]int64, int, error)
}
