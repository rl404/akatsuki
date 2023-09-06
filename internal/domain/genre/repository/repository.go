package repository

import (
	"context"

	"github.com/rl404/akatsuki/internal/domain/genre/entity"
)

// Repository contains functions for genre domain.
type Repository interface {
	BatchUpdate(ctx context.Context, data []entity.Genre) (int, error)
	Get(ctx context.Context, data entity.GetRequest) ([]*entity.Genre, int, int, error)
	GetByID(ctx context.Context, id int64) (*entity.Genre, int, error)
	GetByIDs(ctx context.Context, ids []int64) ([]*entity.Genre, int, error)
	GetHistories(ctx context.Context, data entity.GetHistoriesRequest) ([]entity.History, int, error)
}
