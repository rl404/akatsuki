package repository

import (
	"context"

	"github.com/nstratos/go-myanimelist/mal"
)

// Repository contains functions for mal domain.
type Repository interface {
	GetAnimeByID(ctx context.Context, id int) (*mal.Anime, int, error)
}
