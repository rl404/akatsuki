package sql

import (
	"context"
	_errors "errors"
	"net/http"
	"time"

	"github.com/rl404/akatsuki/internal/domain/anime/entity"
	"github.com/rl404/akatsuki/internal/errors"
	"gorm.io/gorm"
)

// SQL contains functions for anime sql database.
type SQL struct {
	db        *gorm.DB
	airingAge time.Duration
	oldAge    time.Duration
}

// New to create new anime database.
func New(db *gorm.DB, airingAge, oldAge int) *SQL {
	return &SQL{
		db:        db,
		airingAge: time.Duration(airingAge) * 24 * time.Hour,
		oldAge:    time.Duration(oldAge) * 24 * time.Hour,
	}
}

// GetByID to get anime by id.
func (sql *SQL) GetByID(ctx context.Context, id int64) (*entity.Anime, int, error) {
	var a Anime
	if err := sql.db.WithContext(ctx).Where("id = ?", id).First(&a).Error; err != nil {
		if _errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, errors.Wrap(ctx, errors.ErrAnimeNotFound, err)
		}
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	anime := a.toEntity()

	// Get genres.
	var animeGenres []AnimeGenre
	if err := sql.db.WithContext(ctx).Where("anime_id = ?", id).Find(&animeGenres).Error; err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	anime.GenreIDs = make([]int64, len(animeGenres))
	for i, g := range animeGenres {
		anime.GenreIDs[i] = g.GenreID
	}

	// Get pictures.
	var animePictures []AnimePicture
	if err := sql.db.WithContext(ctx).Where("anime_id = ?", id).Find(&animePictures).Error; err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	anime.Pictures = make([]string, len(animePictures))
	for i, p := range animePictures {
		anime.Pictures[i] = p.URL
	}

	// Get related.
	var animeRelated []AnimeRelated
	if err := sql.db.WithContext(ctx).Where("anime_id1 = ?", id).Find(&animeRelated).Error; err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	anime.Related = make([]entity.Related, len(animeRelated))
	for i, r := range animeRelated {
		anime.Related[i] = entity.Related{
			ID:       r.AnimeID2,
			Relation: r.Relation,
		}
	}

	// Get studios.
	var animeStudios []AnimeStudio
	if err := sql.db.WithContext(ctx).Where("anime_id = ?", id).Find(&animeStudios).Error; err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	anime.StudioIDs = make([]int64, len(animeStudios))
	for i, s := range animeStudios {
		anime.StudioIDs[i] = s.StudioID
	}

	return anime, http.StatusOK, nil
}

// GetByIDs to get anime by ids.
func (sql *SQL) GetByIDs(ctx context.Context, ids []int64) ([]*entity.Anime, int, error) {
	var a []Anime
	if err := sql.db.WithContext(ctx).Where("id in (?)", ids).Find(&a).Error; err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return sql.animeToEntities(a), http.StatusOK, nil
}

// Update to update anime data.
func (sql *SQL) Update(ctx context.Context, data entity.Anime) (int, error) {
	tx := sql.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, tx.Error)
	}
	defer tx.Rollback()

	// Get existing anime.
	var a Anime
	if err := tx.WithContext(ctx).Select("created_at").Where("id = ?", data.ID).First(&a).Error; err != nil {
		if !_errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
		}
	}

	// Update anime.
	anime := sql.animeFromEntity(data)
	anime.CreatedAt = a.CreatedAt
	if err := tx.WithContext(ctx).Save(anime).Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	// Delete existing anime genre.
	if err := tx.WithContext(ctx).Where("anime_id = ?", data.ID).Delete(&AnimeGenre{}).Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	// Create new anime genre.
	if err := tx.WithContext(ctx).Create(sql.animeGenreFromEntity(data)).Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	// Delete existing anime picture.
	if err := tx.WithContext(ctx).Where("anime_id = ?", data.ID).Delete(&AnimePicture{}).Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	// Create new anime picture.
	if err := tx.WithContext(ctx).Create(sql.animePictureFromEntity(data)).Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	// Delete existing anime related.
	if err := tx.WithContext(ctx).Where("anime_id1 = ?", data.ID).Delete(&AnimeRelated{}).Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	// Create new anime related.
	if err := tx.WithContext(ctx).Create(sql.animeRelatedFromEntity(data)).Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	// Delete existing anime studio.
	if err := tx.WithContext(ctx).Where("anime_id = ?", data.ID).Delete(&AnimeStudio{}).Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	// Create new anime studio.
	if err := tx.WithContext(ctx).Create(sql.animeStudioFromEntity(data)).Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	// Create new anime stats history.
	if err := tx.WithContext(ctx).Create(sql.animeStatsFromEntity(data)).Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	if err := tx.Commit().Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	return http.StatusOK, nil
}

// IsOld to check if old.
func (sql *SQL) IsOld(ctx context.Context, id int64) (bool, int, error) {
	res := sql.db.WithContext(ctx).
		Where("id = ? and ((status = ? and updated_at >= ?) or (status != ? and updated_at >= ?))", id,
			entity.StatusReleasing, time.Now().Add(-sql.airingAge),
			entity.StatusReleasing, time.Now().Add(-sql.oldAge)).
		Limit(1).
		Find(&[]Anime{})

	if res.Error != nil {
		return true, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, res.Error)
	}

	return res.RowsAffected == 0, http.StatusOK, nil
}

// GetOldAiring to get old airing anime.
func (sql *SQL) GetOldAiring(ctx context.Context, limit int) ([]*entity.Anime, int, error) {
	var a []Anime
	if err := sql.db.WithContext(ctx).Where("status = ? and updated_at <= ?", entity.StatusReleasing, time.Now().Add(-sql.airingAge)).Find(&a).Error; err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return sql.animeToEntities(a), http.StatusOK, nil
}

// GetOldData to get old anime.
func (sql *SQL) GetOldData(ctx context.Context, limit int) ([]*entity.Anime, int, error) {
	var a []Anime
	if err := sql.db.WithContext(ctx).Where("updated_at <= ?", time.Now().Add(-sql.oldAge)).Find(&a).Error; err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return sql.animeToEntities(a), http.StatusOK, nil
}
