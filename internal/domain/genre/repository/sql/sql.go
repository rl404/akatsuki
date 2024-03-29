package sql

import (
	"context"
	_errors "errors"
	"net/http"

	"github.com/rl404/akatsuki/internal/domain/genre/entity"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/fairy/errors/stack"
	"gorm.io/gorm"
)

// SQL contains functions for genre sql database.
type SQL struct {
	db *gorm.DB
}

// New to create new genre database.
func New(db *gorm.DB) *SQL {
	return &SQL{
		db: db,
	}
}

// BatchUpdate to batch update genre.
func (sql *SQL) BatchUpdate(ctx context.Context, data []entity.Genre) (int, error) {
	if err := sql.db.WithContext(ctx).Save(sql.fromEntities(data)).Error; err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return http.StatusOK, nil
}

// GetByIDs to get genre by ids.
func (sql *SQL) GetByIDs(ctx context.Context, ids []int64) ([]*entity.Genre, int, error) {
	var g []Genre
	if err := sql.db.WithContext(ctx).Where("id in ?", ids).Find(&g).Error; err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return sql.toEntities(g), http.StatusOK, nil
}

type genre struct {
	ID     int64
	Name   string
	Count  int
	Mean   float64
	Member int
}

// Get to get list.
func (sql *SQL) Get(ctx context.Context, data entity.GetRequest) ([]*entity.Genre, int, int, error) {
	query := sql.db.
		Select("g.id, g.name, count(*) as count, avg(nullif(a.mean, 0)) as mean, sum(a.member) as member").
		Table("genre as g").
		Joins("left join anime_genre ag on ag.genre_id = g.id").
		Joins("left join anime a on a.id = ag.anime_id").
		Group("g.id, g.name")

	if data.Name != "" {
		query = query.Where("g.name ilike ?", "%"+data.Name+"%")
	}

	var genres []genre
	if err := query.WithContext(ctx).Order(sql.convertSort(data.Sort)).Offset((data.Page - 1) * data.Limit).Limit(data.Limit).Find(&genres).Error; err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	var total int64
	if err := query.WithContext(ctx).Count(&total).Error; err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	res := make([]*entity.Genre, len(genres))
	for i, g := range genres {
		res[i] = &entity.Genre{
			ID:     g.ID,
			Name:   g.Name,
			Count:  g.Count,
			Mean:   g.Mean,
			Member: g.Member,
		}
	}

	return res, int(total), http.StatusOK, nil
}

// GetByID to get by id.
func (sql *SQL) GetByID(ctx context.Context, id int64) (*entity.Genre, int, error) {
	var genre genre
	if err := sql.db.
		Select("g.id, g.name, count(*) as count, avg(nullif(a.mean, 0)) as mean, sum(a.member) as member").
		Table("genre as g").
		Joins("left join anime_genre ag on ag.genre_id = g.id").
		Joins("left join anime a on a.id = ag.anime_id").
		Where("g.id = ?", id).
		Group("g.id, g.name").
		First(&genre).
		Error; err != nil {
		if _errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, stack.Wrap(ctx, err, errors.ErrInvalidGenreID)
		}
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	return &entity.Genre{
		ID:     genre.ID,
		Name:   genre.Name,
		Count:  genre.Count,
		Mean:   genre.Mean,
		Member: genre.Member,
	}, http.StatusOK, nil
}

// GetHistories to get histories.
func (sql *SQL) GetHistories(ctx context.Context, data entity.GetHistoriesRequest) ([]entity.History, int, error) {
	selects := []string{
		"avg(a.mean) as mean",
		"floor(avg(nullif(a.rank, 0))) as rank",
		"floor(avg(nullif(a.popularity, 0))) as popularity",
		"sum(a.member) as member",
		"sum(a.voter) as voter",
		"count(*) as count",
	}

	query := sql.db.WithContext(ctx).
		Table("anime_genre as ag").
		Joins("join anime a on a.id = ag.anime_id").
		Where("ag.genre_id = ?", data.GenreID)

	if data.StartYear > 0 {
		query.Where("a.start_year >= ?", data.StartYear)
	}

	if data.EndYear > 0 {
		query.Where("a.start_year <= ?", data.EndYear)
	}

	switch data.Group {
	case entity.Yearly:
		selects = append(selects, "a.start_year as year")
		query.Group("a.start_year").Order("a.start_year asc")
	case entity.Monthly:
		selects = append(selects, "a.start_year as year, a.start_month as month")
		query.Where("a.start_month != 0")
		query.Group("a.start_year, a.start_month").Order("a.start_year asc, a.start_month asc")
	}

	var histories []genreHistory
	if err := query.Select(selects).Find(&histories).Error; err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	res := make([]entity.History, len(histories))
	for i, h := range histories {
		res[i] = entity.History{
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
