package sql

import (
	"context"
	_errors "errors"
	"net/http"

	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/fairy/errors/stack"
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

// Get to get empty id.
func (sql *SQL) Get(ctx context.Context, id int64) (int64, int, error) {
	var emptyID EmptyID
	if err := sql.db.WithContext(ctx).Where("anime_id = ?", id).First(&emptyID).Error; err != nil {
		if _errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, http.StatusNotFound, stack.Wrap(ctx, err, errors.ErrAnimeNotFound)
		}
		return 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return emptyID.AnimeID, http.StatusOK, nil
}

// Create to create empty id.
func (sql *SQL) Create(ctx context.Context, id int64) (int, error) {
	if err := sql.db.WithContext(ctx).Create(&EmptyID{AnimeID: id}).Error; err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return http.StatusCreated, nil
}

// Delete to delete id from empty id.
func (sql *SQL) Delete(ctx context.Context, id int64) (int, error) {
	if err := sql.db.WithContext(ctx).Where("anime_id = ?", id).Delete(&EmptyID{}).Error; err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return http.StatusOK, nil
}

// GetIDs to get all ids.
func (sql *SQL) GetIDs(ctx context.Context) ([]int64, int, error) {
	var ids []int64
	if err := sql.db.WithContext(ctx).Model(&EmptyID{}).Pluck("anime_id", &ids).Error; err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return ids, http.StatusOK, nil
}
