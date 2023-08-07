package sql

import (
	"context"
	_errors "errors"
	"net/http"

	"github.com/rl404/akatsuki/internal/domain/genre/entity"
	"github.com/rl404/akatsuki/internal/errors"
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
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return http.StatusOK, nil
}

// GetByIDs to get genre by ids.
func (sql *SQL) GetByIDs(ctx context.Context, ids []int64) ([]*entity.Genre, int, error) {
	var g []Genre
	if err := sql.db.WithContext(ctx).Where("id in ?", ids).Find(&g).Error; err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return sql.toEntities(g), http.StatusOK, nil
}

type genre struct {
	ID    int64
	Name  string
	Count int
}

// Get to get list.
func (sql *SQL) Get(ctx context.Context, data entity.GetRequest) ([]*entity.Genre, int, int, error) {
	subQuery := sql.db.
		Select("genre_id, count(*) as count").
		Table("anime_genre").
		Group("genre_id")

	query := sql.db.
		Select("g.id, g.name, ag.count").
		Table("genre as g").
		Joins("left join (?) ag on ag.genre_id = g.id", subQuery)

	if data.Name != "" {
		query = query.Where("g.name ilike ?", "%"+data.Name+"%")
	}

	var genres []genre
	if err := query.WithContext(ctx).Order("g.name asc").Offset((data.Page - 1) * data.Limit).Limit(data.Limit).Find(&genres).Error; err != nil {
		return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var total int64
	if err := query.WithContext(ctx).Count(&total).Error; err != nil {
		return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	res := make([]*entity.Genre, len(genres))
	for i, g := range genres {
		res[i] = &entity.Genre{
			ID:    g.ID,
			Name:  g.Name,
			Count: g.Count,
		}
	}

	return res, int(total), http.StatusOK, nil
}

// GetByID to get by id.
func (sql *SQL) GetByID(ctx context.Context, id int64) (*entity.Genre, int, error) {
	var g Genre
	if err := sql.db.WithContext(ctx).Where("id = ?", id).First(&g).Error; err != nil {
		if _errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, errors.Wrap(ctx, errors.ErrInvalidGenreID, err)
		}
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return g.toEntity(), http.StatusOK, nil
}
