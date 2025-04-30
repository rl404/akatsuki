package sql

import (
	"context"
	_errors "errors"
	"net/http"
	"time"

	"github.com/rl404/akatsuki/internal/domain/user_anime/entity"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/fairy/errors/stack"
	"gorm.io/gorm"
)

// SQL contains functions for user sql database.
type SQL struct {
	db  *gorm.DB
	age time.Duration
}

// New to create new user database.
func New(db *gorm.DB, age int) *SQL {
	return &SQL{
		db:  db,
		age: time.Duration(age) * 24 * time.Hour,
	}
}

// Get to get user anime.
func (sql *SQL) Get(ctx context.Context, data entity.GetUserAnimeRequest) ([]*entity.UserAnime, int, int, error) {
	var a []UserAnime
	query := sql.db.WithContext(ctx).Model(&UserAnime{})

	if data.Username != "" {
		query.Where("username = ?", data.Username)
	}

	if err := query.Limit(data.Limit).Offset((data.Page - 1) * data.Limit).Find(&a).Error; err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	var cnt int64
	if err := query.Limit(-1).Count(&cnt).Error; err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	return sql.userAnimeToEntities(a), int(cnt), http.StatusOK, nil
}

// Update to update user anime.
func (sql *SQL) Update(ctx context.Context, data entity.UserAnime) (int, error) {
	var ua UserAnime
	if err := sql.db.WithContext(ctx).Select("id, created_at").Where("username = ? and anime_id = ?", data.Username, data.AnimeID).First(&ua).Error; err != nil {
		if !_errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
		}
	}

	userAnime := sql.userAnimeFromEntity(data)
	userAnime.ID = ua.ID
	userAnime.CreatedAt = ua.CreatedAt
	if err := sql.db.WithContext(ctx).Save(userAnime).Error; err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	return http.StatusOK, nil
}

// IsOld to check if old.
func (sql *SQL) IsOld(ctx context.Context, username string) (bool, int, error) {
	res := sql.db.WithContext(ctx).Where("username = ? and updated_at >= ?", username, time.Now().Add(-sql.age)).Limit(1).Find(&[]UserAnime{})
	if res.Error != nil {
		return true, http.StatusInternalServerError, stack.Wrap(ctx, res.Error, errors.ErrInternalDB)
	}
	return res.RowsAffected == 0, http.StatusOK, nil
}

// GetOldUsernames to get old usernames.
func (sql *SQL) GetOldUsernames(ctx context.Context) ([]string, int, error) {
	var usernames []string
	if err := sql.db.WithContext(ctx).Model(&UserAnime{}).Where("updated_at <= ?", time.Now().Add(-sql.age)).Pluck("distinct(username)", &usernames).Error; err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return usernames, http.StatusOK, nil
}

// DeleteNotInList to delete anime not in list.
func (sql *SQL) DeleteNotInList(ctx context.Context, username string, ids []int64, status entity.Status) (int, error) {
	if err := sql.db.WithContext(ctx).Where("username = ? and anime_id not in ? and status = ?", username, ids, status).Delete(&UserAnime{}).Error; err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return http.StatusOK, nil
}

// DeleteByAnimeID to delete by anime id.
func (sql *SQL) DeleteByAnimeID(ctx context.Context, animeID int64) (int, error) {
	if err := sql.db.WithContext(ctx).Where("anime_id = ?", animeID).Delete(&UserAnime{}).Error; err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return http.StatusOK, nil
}

// DeleteByUsername to delete by username.
func (sql *SQL) DeleteByUsername(ctx context.Context, username string) (int, error) {
	if err := sql.db.WithContext(ctx).Where("username = ?", username).Delete(&UserAnime{}).Error; err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return http.StatusOK, nil
}
