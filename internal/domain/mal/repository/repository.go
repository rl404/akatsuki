package repository

import (
	"context"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/rl404/akatsuki/internal/domain/mal/entity"
)

// Repository contains functions for mal domain.
type Repository interface {
	GetAnimeByID(ctx context.Context, id int) (*mal.Anime, int, error)
	GetUserAnime(ctx context.Context, data entity.GetUserAnimeRequest) ([]mal.UserAnime, int, error)
}
