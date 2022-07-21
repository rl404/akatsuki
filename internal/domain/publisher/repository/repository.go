package repository

import (
	"context"

	"github.com/rl404/akatsuki/internal/domain/publisher/entity"
)

// Repository contains functions for publisher domain.
type Repository interface {
	PublishParseAnime(ctx context.Context, data entity.ParseAnimeRequest) error
	PublishParseUserAnime(ctx context.Context, data entity.ParseUserAnimeRequest) error
}
