package repository

import (
	"context"

	"github.com/rl404/akatsuki/internal/domain/mal/entity"
	"github.com/rl404/nagato"
)

// Repository contains functions for mal domain.
type Repository interface {
	GetAnimeByID(ctx context.Context, id int) (*nagato.Anime, int, error)
	GetUserAnime(ctx context.Context, data entity.GetUserAnimeRequest) ([]nagato.UserAnime, int, error)
}
