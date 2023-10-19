package repository

import (
	"context"
)

// Repository contains functions for publisher domain.
type Repository interface {
	PublishParseAnime(ctx context.Context, id int64, forced bool) error
	PublishParseUserAnime(ctx context.Context, username, status string, forced bool) error
}
