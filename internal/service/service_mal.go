package service

import (
	"context"
	"net/http"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/rl404/akatsuki/internal/domain/mal/entity"
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

// GetMalUserAnimeRequest is get mal user anime request model.
type GetMalUserAnimeRequest struct {
	UserName string
	Status   string
	Sort     string
	Limit    int
	Offset   int
}

// GetMalUserAnime to get mal user name.
func (s *service) GetMalUserAnime(ctx context.Context, data GetMalUserAnimeRequest) ([]mal.UserAnime, int, error) {
	anime, code, err := s.mal.GetUserAnime(ctx, entity.GetUserAnimeRequest{
		Username: data.UserName,
		Status:   mal.AnimeStatus(data.Status),
		Sort:     mal.SortAnimeList(data.Sort),
		Limit:    data.Limit,
		Offset:   data.Offset,
	})
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}
	return anime, http.StatusOK, nil
}
