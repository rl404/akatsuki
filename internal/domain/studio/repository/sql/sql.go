package sql

import (
	"context"
	_errors "errors"
	"net/http"

	"github.com/rl404/akatsuki/internal/domain/studio/entity"
	"github.com/rl404/akatsuki/internal/errors"
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
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return http.StatusOK, nil
}

// GetByIDs to get studio by ids.
func (sql *SQL) GetByIDs(ctx context.Context, ids []int64) ([]*entity.Studio, int, error) {
	var s []Studio
	if err := sql.db.WithContext(ctx).Where("id in ?", ids).Find(&s).Error; err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return sql.toEntities(s), http.StatusOK, nil
}

type studio struct {
	ID    int64
	Name  string
	Count int
}

// Get to get list.
func (sql *SQL) Get(ctx context.Context, data entity.GetRequest) ([]*entity.Studio, int, int, error) {
	subQuery := sql.db.
		Select("studio_id, count(*) as count").
		Table("anime_studio").
		Group("studio_id")

	query := sql.db.
		Select("g.id, g.name, ag.count").
		Table("studio as g").
		Joins("left join (?) ag on ag.studio_id = g.id", subQuery)

	if data.Name != "" {
		query = query.Where("g.name ilike ?", "%"+data.Name+"%")
	}

	var studios []studio
	if err := query.WithContext(ctx).Order("lower(g.name) asc").Offset((data.Page - 1) * data.Limit).Limit(data.Limit).Find(&studios).Error; err != nil {
		return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var total int64
	if err := query.WithContext(ctx).Count(&total).Error; err != nil {
		return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	res := make([]*entity.Studio, len(studios))
	for i, g := range studios {
		res[i] = &entity.Studio{
			ID:    g.ID,
			Name:  g.Name,
			Count: g.Count,
		}
	}

	return res, int(total), http.StatusOK, nil
}

// GetByID to get by id.
func (sql *SQL) GetByID(ctx context.Context, id int64) (*entity.Studio, int, error) {
	var g Studio
	if err := sql.db.WithContext(ctx).Where("id = ?", id).First(&g).Error; err != nil {
		if _errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, errors.Wrap(ctx, errors.ErrInvalidStudioID, err)
		}
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return g.toEntity(), http.StatusOK, nil
}
