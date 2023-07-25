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
	db           *gorm.DB
	finishedAge  time.Duration
	releasingAge time.Duration
	notYetAge    time.Duration
}

// New to create new anime database.
func New(db *gorm.DB, finishedAge, releasingAge, notYetAge int) *SQL {
	return &SQL{
		db:           db,
		finishedAge:  time.Duration(finishedAge) * 24 * time.Hour,
		releasingAge: time.Duration(releasingAge) * 24 * time.Hour,
		notYetAge:    time.Duration(notYetAge) * 24 * time.Hour,
	}
}

// Get to get anime list.
func (sql *SQL) Get(ctx context.Context, data entity.GetRequest) ([]*entity.Anime, int, int, error) {
	query := sql.db

	if data.Title != "" {
		query = query.Where("title ilike ? or title_synonym ilike ? or title_english ilike ? or title_japanese ilike ?", "%"+data.Title+"%", "%"+data.Title+"%", "%"+data.Title+"%", "%"+data.Title+"%")
	}

	if data.NSFW != nil {
		query = query.Where("nsfw = ?", data.NSFW)
	}

	if data.Type != "" {
		query = query.Where("type = ?", data.Type)
	}

	if data.Status != "" {
		query = query.Where("status = ?", data.Status)
	}

	if data.Season != "" {
		query = query.Where("season = ?", data.Season)
	}

	if data.SeasonYear != 0 {
		query = query.Where("season_year = ?", data.SeasonYear)
	}

	if data.StartMean != 0 {
		query = query.Where("mean >= ?", data.StartMean)
	}

	if data.EndMean != 0 {
		query = query.Where("mean <= ?", data.EndMean)
	}

	var a []Anime
	if err := query.WithContext(ctx).Order(sql.convertSort(data.Sort)).Offset((data.Page - 1) * data.Limit).Limit(data.Limit).Find(&a).Error; err != nil {
		return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var total int64
	if err := query.WithContext(ctx).Model(&Anime{}).Count(&total).Error; err != nil {
		return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	return sql.animeToEntities(a), int(total), http.StatusOK, nil
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
	if err := sql.db.WithContext(ctx).Where("id in ?", ids).Find(&a).Error; err != nil {
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
	if len(data.GenreIDs) > 0 {
		if err := tx.WithContext(ctx).Create(sql.animeGenreFromEntity(data)).Error; err != nil {
			return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
		}
	}

	// Delete existing anime picture.
	if err := tx.WithContext(ctx).Where("anime_id = ?", data.ID).Delete(&AnimePicture{}).Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	// Create new anime picture.
	if len(data.Pictures) > 0 {
		if err := tx.WithContext(ctx).Create(sql.animePictureFromEntity(data)).Error; err != nil {
			return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
		}
	}

	// Delete existing anime related.
	if err := tx.WithContext(ctx).Where("anime_id1 = ?", data.ID).Delete(&AnimeRelated{}).Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	// Create new anime related.
	if len(data.Related) > 0 {
		if err := tx.WithContext(ctx).Create(sql.animeRelatedFromEntity(data)).Error; err != nil {
			return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
		}
	}

	// Delete existing anime studio.
	if err := tx.WithContext(ctx).Where("anime_id = ?", data.ID).Delete(&AnimeStudio{}).Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	// Create new anime studio.
	if len(data.StudioIDs) > 0 {
		if err := tx.WithContext(ctx).Create(sql.animeStudioFromEntity(data)).Error; err != nil {
			return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
		}
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
		Where("id = ? and ((status = ? and updated_at >= ?) or (status = ? and updated_at >= ?) or (status = ? and updated_at >= ?))", id,
			entity.StatusFinished, time.Now().Add(-sql.finishedAge),
			entity.StatusReleasing, time.Now().Add(-sql.releasingAge),
			entity.StatusNotYet, time.Now().Add(-sql.notYetAge)).
		Limit(1).
		Find(&[]Anime{})

	if res.Error != nil {
		return true, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, res.Error)
	}

	return res.RowsAffected == 0, http.StatusOK, nil
}

func (sql *SQL) getOldIDs(ctx context.Context, status entity.Status, age time.Duration) ([]int64, int, error) {
	var ids []int64
	if err := sql.db.WithContext(ctx).Model(&Anime{}).Where("status = ? and updated_at <= ?", status, time.Now().Add(-age)).Pluck("id", &ids).Error; err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return ids, http.StatusOK, nil
}

// GetOldReleasingIDs to get old releasing anime ids.
func (sql *SQL) GetOldReleasingIDs(ctx context.Context) ([]int64, int, error) {
	return sql.getOldIDs(ctx, entity.StatusReleasing, sql.releasingAge)
}

// GetOldFinishedIDs to get old finished anime ids.
func (sql *SQL) GetOldFinishedIDs(ctx context.Context) ([]int64, int, error) {
	return sql.getOldIDs(ctx, entity.StatusFinished, sql.finishedAge)
}

// GetOldNotYetIDs to get old not yet released anime ids.
func (sql *SQL) GetOldNotYetIDs(ctx context.Context) ([]int64, int, error) {
	return sql.getOldIDs(ctx, entity.StatusNotYet, sql.notYetAge)
}

// GetMaxID to get max id.
func (sql *SQL) GetMaxID(ctx context.Context) (int64, int, error) {
	var id int64
	if err := sql.db.WithContext(ctx).Model(&Anime{}).Select("COALESCE(MAX(id), 1)").Row().Scan(&id); err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return id, http.StatusOK, nil
}

// GetIDs to get all anime ids.
func (sql *SQL) GetIDs(ctx context.Context) ([]int64, int, error) {
	var ids []int64
	if err := sql.db.WithContext(ctx).Model(&Anime{}).Pluck("id", &ids).Error; err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return ids, http.StatusOK, nil
}

// GetRelatedByIDs to get related anime by ids.
func (sql *SQL) GetRelatedByIDs(ctx context.Context, ids []int64) ([]*entity.AnimeRelated, int, error) {
	var ar []AnimeRelated
	if err := sql.db.WithContext(ctx).Where("anime_id1 in ?", ids).Find(&ar).Error; err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return sql.animeRelatedToEntities(ar), http.StatusOK, nil
}

// DeleteByID to delete by id.
func (sql *SQL) DeleteByID(ctx context.Context, id int64) (int, error) {
	tx := sql.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, tx.Error)
	}
	defer tx.Rollback()

	if err := tx.WithContext(ctx).Where("id = ?", id).Delete(&Anime{}).Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	if err := tx.WithContext(ctx).Where("anime_id = ?", id).Delete(&AnimeGenre{}).Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	if err := tx.WithContext(ctx).Where("anime_id = ?", id).Delete(&AnimePicture{}).Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	if err := tx.WithContext(ctx).Where("anime_id1 = ? or anime_id2 = ?", id, id).Delete(&AnimeRelated{}).Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	if err := tx.WithContext(ctx).Where("anime_id = ?", id).Delete(&AnimeStudio{}).Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	if err := tx.WithContext(ctx).Where("anime_id = ?", id).Delete(&AnimeStatsHistory{}).Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	if err := tx.Commit().Error; err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	return http.StatusOK, nil
}

// GetHistories to get histories.
func (sql *SQL) GetHistories(ctx context.Context, data entity.GetHistoriesRequest) ([]entity.History, int, error) {
	selects := []string{
		"avg(mean) as mean",
		"floor(avg(rank)) as rank",
		"floor(avg(popularity)) as popularity",
		"floor(avg(member)) as member",
		"floor(avg(voter)) as voter",
		"floor(avg(user_watching)) as user_watching",
		"floor(avg(user_completed)) as user_completed",
		"floor(avg(user_on_hold)) as user_on_hold",
		"floor(avg(user_dropped)) as user_dropped",
		"floor(avg(user_planned)) as user_planned",
	}

	query := sql.db.WithContext(ctx).Model(&animeStatsHistory{}).Where("anime_id = ?", data.AnimeID)

	if data.StartDate != nil {
		query.Where("created_at >= ?", data.StartDate)
	}

	if data.EndDate != nil {
		query.Where("created_at <= ?", data.EndDate)
	}

	switch data.Group {
	case entity.Yearly:
		selects = append(selects, "date_part('year',created_at) as year")
		query.Group("date_part('year',created_at)").Order("year asc")
	case entity.Monthly:
		selects = append(selects, "date_part('year',created_at) as year, date_part('month',created_at) as month")
		query.Group("date_part('year',created_at), date_part('month',created_at)").Order("year asc, month asc")
	case entity.Weekly:
		selects = append(selects, "date_part('year',created_at) as year, date_part('month',created_at) as month, to_char(created_at,'W') as week")
		query.Group("date_part('year',created_at), date_part('month',created_at), to_char(created_at,'W')").Order("year asc, month asc, week asc")
	}

	var histories []animeStatsHistory
	if err := query.Select(selects).Find(&histories).Error; err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	res := make([]entity.History, len(histories))
	for i, h := range histories {
		res[i] = entity.History{
			Year:          h.Year,
			Month:         h.Month,
			Week:          h.Week,
			Mean:          h.Mean,
			Rank:          h.Rank,
			Popularity:    h.Popularity,
			Member:        h.Member,
			Voter:         h.Voter,
			UserWatching:  h.UserWatching,
			UserCompleted: h.UserCompleted,
			UserOnHold:    h.UserOnHold,
			UserDropped:   h.UserDropped,
			UserPlanned:   h.UserPlanned,
		}
	}

	return res, http.StatusOK, nil
}
