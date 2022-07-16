package service

import (
	"context"
	"net/http"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/rl404/akatsuki/internal/errors"
)

// GetMalAnimeByID to get mal anime by id.
func (s *service) GetMalAnimeByID(ctx context.Context, id int) (*mal.Anime, int, error) {
	anime, code, err := s.mal.GetAnimeByID(ctx, id)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}
	return anime, http.StatusOK, nil
}
