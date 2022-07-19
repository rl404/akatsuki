package sql

import (
	"context"
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
	if err := sql.db.WithContext(ctx).Where("id in (?)", ids).Find(&g).Error; err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return sql.toEntities(g), http.StatusOK, nil
}
