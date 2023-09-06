package service

import (
	"context"
	"net/http"
	"time"

	"github.com/rl404/akatsuki/internal/domain/studio/entity"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/utils"
)

// Studio is studio model.
type Studio struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	Count  int     `json:"count"`
	Mean   float64 `json:"mean"`
	Member int     `json:"member"`
}

// GetStudiosRequest is get studio list request model.
type GetStudiosRequest struct {
	Name  string      `mod:"lcase,trim"`
	Sort  entity.Sort `validate:"oneof=NAME -NAME COUNT -COUNT MEAN -MEAN MEMBER -MEMBER" mod:"no_space,ucase,default=NAME"`
	Page  int         `validate:"required,gte=1" mod:"default=1"`
	Limit int         `validate:"required,gte=-1" mod:"default=20"`
}

// GetStudios to get studio list.
func (s *service) GetStudios(ctx context.Context, data GetStudiosRequest) ([]Studio, *Pagination, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, nil, http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	studios, total, code, err := s.studio.Get(ctx, entity.GetRequest{
		Name:  data.Name,
		Sort:  data.Sort,
		Page:  data.Page,
		Limit: data.Limit,
	})
	if err != nil {
		return nil, nil, code, errors.Wrap(ctx, err)
	}

	res := make([]Studio, len(studios))
	for i, g := range studios {
		res[i] = Studio{
			ID:     g.ID,
			Name:   g.Name,
			Count:  g.Count,
			Mean:   g.Mean,
			Member: g.Member,
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
		ID:     studio.ID,
		Name:   studio.Name,
		Count:  studio.Count,
		Mean:   studio.Mean,
		Member: studio.Member,
	}, http.StatusOK, nil
}

// StudioHistory is studio stats history.
type StudioHistory struct {
	Year       int     `json:"year"`
	Month      int     `json:"month"`
	Mean       float64 `json:"mean"`
	Rank       int     `json:"rank"`
	Popularity int     `json:"popularity"`
	Member     int     `json:"member"`
	Voter      int     `json:"voter"`
	Count      int     `json:"count"`
}

// GetStudioHistoriesRequest is get studio history request model.
type GetStudioHistoriesRequest struct {
	ID        int64               `validate:"gt=0"`
	StartYear int                 `validate:"gte=0"`
	EndYear   int                 `validate:"gte=0"`
	Group     entity.HistoryGroup `validate:"oneof=MONTHLY YEARLY" mod:"trim,ucase,default=MONTHLY"`
}

// GetStudioHistoriesByID to get anime history by id.
func (s *service) GetStudioHistoriesByID(ctx context.Context, data GetStudioHistoriesRequest) ([]StudioHistory, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	if data.StartYear == 0 {
		switch data.Group {
		case entity.Yearly:
			data.StartYear = time.Now().AddDate(-20, 0, 0).Year()
		case entity.Monthly:
			data.StartYear = time.Now().AddDate(-3, 0, 0).Year()
		}
	}

	if data.EndYear == 0 {
		data.EndYear = time.Now().Year()
	}

	histories, code, err := s.studio.GetHistories(ctx, entity.GetHistoriesRequest{
		StudioID:  data.ID,
		StartYear: data.StartYear,
		EndYear:   data.EndYear,
		Group:     data.Group,
	})
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	res := make([]StudioHistory, len(histories))
	for i, h := range histories {
		res[i] = StudioHistory{
			Year:       h.Year,
			Month:      h.Month,
			Mean:       h.Mean,
			Rank:       h.Rank,
			Popularity: h.Popularity,
			Member:     h.Member,
			Voter:      h.Voter,
			Count:      h.Count,
		}
	}

	return res, http.StatusOK, nil
}
