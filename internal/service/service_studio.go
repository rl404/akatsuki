package service

import (
	"context"
	"net/http"

	"github.com/rl404/akatsuki/internal/domain/studio/entity"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/utils"
)

// Studio is studio model.
type Studio struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// GetStudiosRequest is get studio list request model.
type GetStudiosRequest struct {
	Name  string `mod:"lcase,trim"`
	Page  int    `validate:"required,gte=1" mod:"default=1"`
	Limit int    `validate:"required,gte=-1" mod:"default=20"`
}

// GetStudios to get studio list.
func (s *service) GetStudios(ctx context.Context, data GetStudiosRequest) ([]Studio, *Pagination, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, nil, http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	studios, total, code, err := s.studio.Get(ctx, entity.GetRequest{
		Name:  data.Name,
		Page:  data.Page,
		Limit: data.Limit,
	})
	if err != nil {
		return nil, nil, code, errors.Wrap(ctx, err)
	}

	res := make([]Studio, len(studios))
	for i, g := range studios {
		res[i] = Studio{
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

// GetStudioByID to get studio by id.
func (s *service) GetStudioByID(ctx context.Context, id int64) (*Studio, int, error) {
	studio, code, err := s.studio.GetByID(ctx, id)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	return &Studio{
		ID:    studio.ID,
		Name:  studio.Name,
		Count: studio.Count,
	}, http.StatusOK, nil
}
