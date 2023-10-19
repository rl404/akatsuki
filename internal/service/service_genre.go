package service

import (
	"context"
	"net/http"
	"time"

	"github.com/rl404/akatsuki/internal/domain/genre/entity"
	"github.com/rl404/akatsuki/internal/utils"
	"github.com/rl404/fairy/errors/stack"
)

// Genre is genre model.
type Genre struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	Count  int     `json:"count"`
	Mean   float64 `json:"mean"`
	Member int     `json:"member"`
}

// GetGenresRequest is get genre list request model.
type GetGenresRequest struct {
	Name  string      `mod:"lcase,trim"`
	Sort  entity.Sort `validate:"oneof=NAME -NAME COUNT -COUNT MEAN -MEAN MEMBER -MEMBER" mod:"no_space,ucase,default=NAME"`
	Page  int         `validate:"required,gte=1" mod:"default=1"`
	Limit int         `validate:"required,gte=-1" mod:"default=20"`
}

// GetGenres to get genre list.
func (s *service) GetGenres(ctx context.Context, data GetGenresRequest) ([]Genre, *Pagination, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, nil, http.StatusBadRequest, stack.Wrap(ctx, err)
	}

	genres, total, code, err := s.genre.Get(ctx, entity.GetRequest{
		Name:  data.Name,
		Sort:  data.Sort,
		Page:  data.Page,
		Limit: data.Limit,
	})
	if err != nil {
		return nil, nil, code, stack.Wrap(ctx, err)
	}

	res := make([]Genre, len(genres))
	for i, g := range genres {
		res[i] = Genre{
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

// GetGenreByID to get genre by id.
func (s *service) GetGenreByID(ctx context.Context, id int64) (*Genre, int, error) {
	genre, code, err := s.genre.GetByID(ctx, id)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	return &Genre{
		ID:     genre.ID,
		Name:   genre.Name,
		Count:  genre.Count,
		Mean:   genre.Mean,
		Member: genre.Member,
	}, http.StatusOK, nil
}

// GenreHistory is genre stats history.
type GenreHistory struct {
	Year       int     `json:"year"`
	Month      int     `json:"month"`
	Mean       float64 `json:"mean"`
	Rank       int     `json:"rank"`
	Popularity int     `json:"popularity"`
	Member     int     `json:"member"`
	Voter      int     `json:"voter"`
	Count      int     `json:"count"`
}

// GetGenreHistoriesRequest is get genre history request model.
type GetGenreHistoriesRequest struct {
	ID        int64               `validate:"gt=0"`
	StartYear int                 `validate:"gte=0"`
	EndYear   int                 `validate:"gte=0"`
	Group     entity.HistoryGroup `validate:"oneof=MONTHLY YEARLY" mod:"trim,ucase,default=MONTHLY"`
}

// GetGenreHistoriesByID to get anime history by id.
func (s *service) GetGenreHistoriesByID(ctx context.Context, data GetGenreHistoriesRequest) ([]GenreHistory, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, http.StatusBadRequest, stack.Wrap(ctx, err)
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

	histories, code, err := s.genre.GetHistories(ctx, entity.GetHistoriesRequest{
		GenreID:   data.ID,
		StartYear: data.StartYear,
		EndYear:   data.EndYear,
		Group:     data.Group,
	})
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	res := make([]GenreHistory, len(histories))
	for i, h := range histories {
		res[i] = GenreHistory{
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
