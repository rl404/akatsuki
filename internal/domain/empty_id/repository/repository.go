package repository

import "context"

// Repository contains functions for empty_id domain.
type Repository interface {
	IsEmpty(ctx context.Context, id int64) (bool, int, error)
	Create(ctx context.Context, id int64) (int, error)
	Delete(ctx context.Context, id int64) (int, error)
}
