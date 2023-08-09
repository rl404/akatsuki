package repository

import (
	"context"

	"github.com/rl404/akatsuki/internal/domain/studio/entity"
)

// Repository contains functions for studio domain.
type Repository interface {
	BatchUpdate(ctx context.Context, data []entity.Studio) (int, error)
	Get(ctx context.Context, data entity.GetRequest) ([]*entity.Studio, int, int, error)
	GetByID(ctx context.Context, id int64) (*entity.Studio, int, error)
	GetByIDs(ctx context.Context, ids []int64) ([]*entity.Studio, int, error)
	GetHistories(ctx context.Context, data entity.GetHistoriesRequest) ([]entity.History, int, error)
}
