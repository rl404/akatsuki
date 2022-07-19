package sql

import (
	"context"
	_errors "errors"
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

// Get to get empty id.
func (sql *SQL) Get(ctx context.Context, id int64) (int64, int, error) {
	var emptyID EmptyID
	if err := sql.db.WithContext(ctx).Where("anime_id = ?", id).First(&emptyID).Error; err != nil {
		if _errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, http.StatusNotFound, errors.Wrap(ctx, errors.ErrAnimeNotFound, err)
		}
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return emptyID.AnimeID, http.StatusOK, nil
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

// GetIDs to get all ids.
func (sql *SQL) GetIDs(ctx context.Context) ([]int64, int, error) {
	var ids []int64
	if err := sql.db.WithContext(ctx).Model(&EmptyID{}).Pluck("anime_id", &ids).Error; err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return ids, http.StatusOK, nil
}
