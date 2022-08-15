package sql

import (
	"context"
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
