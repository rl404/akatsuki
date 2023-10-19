package sql

import (
	"context"
	_errors "errors"
	"net/http"

	"github.com/rl404/akatsuki/internal/domain/studio/entity"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/fairy/errors/stack"
	"gorm.io/gorm"
)

// SQL contains functions for studio sql database.
type SQL struct {
	db *gorm.DB
}

// New to create new studio database.
func New(db *gorm.DB) *SQL {
	return &SQL{
		db: db,
	}
}

// BatchUpdate to batch update studio.
func (sql *SQL) BatchUpdate(ctx context.Context, data []entity.Studio) (int, error) {
	if err := sql.db.WithContext(ctx).Save(sql.fromEntities(data)).Error; err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return http.StatusOK, nil
}

// GetByIDs to get studio by ids.
func (sql *SQL) GetByIDs(ctx context.Context, ids []int64) ([]*entity.Studio, int, error) {
	var s []Studio
	if err := sql.db.WithContext(ctx).Where("id in ?", ids).Find(&s).Error; err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return sql.toEntities(s), http.StatusOK, nil
}

type studio struct {
	ID     int64
	Name   string
	Count  int
	Mean   float64
	Member int
}

// Get to get list.
func (sql *SQL) Get(ctx context.Context, data entity.GetRequest) ([]*entity.Studio, int, int, error) {
	query := sql.db.
		Select("g.id, g.name, count(*) as count, avg(nullif(a.mean, 0)) as mean, sum(a.member) as member").
		Table("studio as g").
		Joins("left join anime_studio ag on ag.studio_id = g.id").
		Joins("left join anime a on a.id = ag.anime_id").
		Group("g.id, g.name")

	if data.Name != "" {
		query = query.Where("g.name ilike ?", "%"+data.Name+"%")
	}

	var studios []studio
	if err := query.WithContext(ctx).Order(sql.convertSort(data.Sort)).Offset((data.Page - 1) * data.Limit).Limit(data.Limit).Find(&studios).Error; err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var total int64
	if err := query.WithContext(ctx).Count(&total).Error; err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, errors.ErrInternalDB, err)
	}

	res := make([]*entity.Studio, len(studios))
	for i, g := range studios {
		res[i] = &entity.Studio{
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
func (sql *SQL) GetByID(ctx context.Context, id int64) (*entity.Studio, int, error) {
	var studio studio
	if err := sql.db.
		Select("g.id, g.name, count(*) as count, avg(nullif(a.mean, 0)) as mean, sum(a.member) as member").
		Table("studio as g").
		Joins("left join anime_studio ag on ag.studio_id = g.id").
		Joins("left join anime a on a.id = ag.anime_id").
		Where("g.id = ?", id).
		Group("g.id, g.name").
		First(&studio).
		Error; err != nil {
		if _errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, stack.Wrap(ctx, err, errors.ErrInvalidStudioID)
		}
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, errors.ErrInternalDB, err)
	}

	return &entity.Studio{
		ID:     studio.ID,
		Name:   studio.Name,
		Count:  studio.Count,
		Mean:   studio.Mean,
		Member: studio.Member,
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
		Table("anime_studio as ans").
		Joins("join anime a on a.id = ans.anime_id").
		Where("ans.studio_id = ?", data.StudioID)

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

	var histories []studioHistory
	if err := query.Select(selects).Find(&histories).Error; err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, errors.ErrInternalDB, err)
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
