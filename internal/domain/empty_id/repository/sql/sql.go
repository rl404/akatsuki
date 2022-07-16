package sql

import (
	"context"
	"net/http"

	"github.com/rl404/akatsuki/internal/errors"
	"gorm.io/gorm"
)

// SQL contains functions for empty_id sql database.
type SQL struct {
	db *gorm.DB
}

// New to create new empty_id database.
func New(db *gorm.DB) *SQL {
	return &SQL{
		db: db,
	}
}

// IsEmpty to check if id is empty.
func (sql *SQL) IsEmpty(ctx context.Context, id int64) (bool, int, error) {
	res := sql.db.WithContext(ctx).Where("anime_id = ?", id).Limit(1).Find(&[]EmptyID{})
	if res.Error != nil {
		return true, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, res.Error)
	}
	return res.RowsAffected != 0, http.StatusOK, nil
}

// Create to create empty id.
func (sql *SQL) Create(ctx context.Context, id int64) (int, error) {
	if err := sql.db.WithContext(ctx).Create(&EmptyID{AnimeID: id}).Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return http.StatusCreated, nil
}

// Delete to delete id from empty id.
func (sql *SQL) Delete(ctx context.Context, id int64) (int, error) {
	res := sql.db.WithContext(ctx).Where("anime_id = ?", id).Delete(&EmptyID{})

	if res.Error != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, res.Error)
	}

	if res.RowsAffected == 0 {
		return http.StatusNotFound, errors.Wrap(ctx, errors.ErrAnimeNotFound)
	}

	return http.StatusOK, nil
}
