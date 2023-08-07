package service

import (
	"context"
	"net/http"

	"github.com/rl404/akatsuki/internal/domain/genre/entity"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/utils"
)

// Genre is genre model.
type Genre struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// GetGenresRequest is get genre list request model.
type GetGenresRequest struct {
	Name  string `mod:"lcase,trim"`
	Page  int    `validate:"required,gte=1" mod:"default=1"`
	Limit int    `validate:"required,gte=-1" mod:"default=20"`
}

// GetGenres to get genre list.
func (s *service) GetGenres(ctx context.Context, data GetGenresRequest) ([]Genre, *Pagination, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, nil, http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	genres, total, code, err := s.genre.Get(ctx, entity.GetRequest{
		Name:  data.Name,
		Page:  data.Page,
		Limit: data.Limit,
	})
	if err != nil {
		return nil, nil, code, errors.Wrap(ctx, err)
	}

	res := make([]Genre, len(genres))
	for i, g := range genres {
		res[i] = Genre{
			ID:    g.ID,
			Name:  g.Name,
			Count: g.Count,
		}
	}

	return res, &Pagination{
		Page:  data.Page,
		Limit: data.Limit,
		Total: total,
	}, http.StatusOK, nil
}

// GetGenreByID to get genre by id.
func (s *service) GetGenreByID(ctx context.Context, id int64) (*Genre, int, error) {
	genre, code, err := s.genre.GetByID(ctx, id)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	return &Genre{
		ID:    genre.ID,
		Name:  genre.Name,
		Count: genre.Count,
	}, http.StatusOK, nil
}
