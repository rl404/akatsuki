package service

import (
	"context"
	"net/http"
	"time"

	"github.com/rl404/akatsuki/internal/domain/anime/entity"
	publisherEntity "github.com/rl404/akatsuki/internal/domain/publisher/entity"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/utils"
)

// Anime is anime model.
type Anime struct {
	ID                int64            `json:"id"`
	Title             string           `json:"title"`
	AlternativeTitles alternativeTitle `json:"alternative_titles"`
	Picture           string           `json:"picture"`
	StartDate         date             `json:"start_date"`
	EndDate           date             `json:"end_date"`
	Synopsis          string           `json:"synopsis"`
	Background        string           `json:"background"`
	NSFW              bool             `json:"nsfw"`
	Type              entity.Type      `json:"type" swaggertype:"string"`
	Status            entity.Status    `json:"status" swaggertype:"string"`
	Episode           episode          `json:"episode"`
	Season            *season          `json:"season"`
	Broadcast         *broadcast       `json:"broadcast"`
	Source            entity.Source    `json:"source" swaggertype:"string"`
	Rating            entity.Rating    `json:"rating" swaggertype:"string"`
	Mean              float64          `json:"mean"`
	Rank              int              `json:"rank"`
	Popularity        int              `json:"popularity"`
	Member            int              `json:"member"`
	Voter             int              `json:"voter"`
	Stats             stats            `json:"stats"`
	Genres            []genre          `json:"genres"`
	Pictures          []string         `json:"pictures"`
	Related           []related        `json:"related"`
	Studios           []studio         `json:"studios"`
	UpdatedAt         time.Time        `json:"updated_at"`
}

// GetAnimeByID to get anime by id.
func (s *service) GetAnimeByID(ctx context.Context, id int64) (*Anime, int, error) {
	if code, err := s.validateID(ctx, id); err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	// Get anime from db.
	animeDB, code, err := s.anime.GetByID(ctx, id)
	if err != nil {
		if code == http.StatusNotFound {
			// Queue to parse.
			if err := s.publisher.PublishParseAnime(ctx, publisherEntity.ParseAnimeRequest{ID: id}); err != nil {
				return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
			}
			return nil, http.StatusAccepted, nil
		}
		return nil, code, errors.Wrap(ctx, err)
	}

	anime := s.animeFromEntity(animeDB)

	// Get genres.
	if len(animeDB.GenreIDs) > 0 {
		genres, code, err := s.genre.GetByIDs(ctx, animeDB.GenreIDs)
		if err != nil {
			return nil, code, errors.Wrap(ctx, err)
		}

		anime.Genres = make([]genre, len(genres))
		for i, g := range genres {
			anime.Genres[i] = genre{
				ID:   g.ID,
				Name: g.Name,
			}
		}
	}

	// Get related.
	if len(animeDB.Related) > 0 {
		relatedMap := make(map[int64]entity.Relation)
		relatedIDs := make([]int64, len(animeDB.Related))
		for i, r := range animeDB.Related {
			relatedIDs[i] = r.ID
			relatedMap[r.ID] = r.Relation
		}

		relates, code, err := s.anime.GetByIDs(ctx, relatedIDs)
		if err != nil {
			return nil, code, errors.Wrap(ctx, err)
		}

		anime.Related = make([]related, len(relates))
		for i, r := range relates {
			anime.Related[i] = related{
				ID:       r.ID,
				Title:    r.Title,
				Picture:  r.Picture,
				Relation: relatedMap[r.ID],
			}
		}
	}

	// Get studios.
	if len(animeDB.StudioIDs) > 0 {
		studios, code, err := s.studio.GetByIDs(ctx, animeDB.StudioIDs)
		if err != nil {
			return nil, code, errors.Wrap(ctx, err)
		}

		anime.Studios = make([]studio, len(studios))
		for i, s := range studios {
			anime.Studios[i] = studio{
				ID:   s.ID,
				Name: s.Name,
			}
		}
	}

	return &anime, http.StatusOK, nil
}

func (s *service) validateID(ctx context.Context, id int64) (int, error) {
	if id <= 0 {
		return http.StatusBadRequest, errors.Wrap(ctx, errors.ErrInvalidAnimeID)
	}

	if _, code, err := s.emptyID.Get(ctx, id); err != nil {
		if code == http.StatusNotFound {
			return http.StatusOK, nil
		}
		return code, errors.Wrap(ctx, err)
	}

	return http.StatusNotFound, errors.Wrap(ctx, errors.ErrAnimeNotFound)
}

// AnimeHistory is anime stats history.
type AnimeHistory struct {
	Year          int     `json:"year"`
	Month         int     `json:"month"`
	Week          int     `json:"week"`
	Mean          float64 `json:"mean"`
	Rank          int     `json:"rank"`
	Popularity    int     `json:"popularity"`
	Member        int     `json:"member"`
	Voter         int     `json:"voter"`
	UserWatching  int     `json:"user_watching"`
	UserCompleted int     `json:"user_completed"`
	UserOnHold    int     `json:"user_on_hold"`
	UserDropped   int     `json:"user_dropped"`
	UserPlanned   int     `json:"user_planned"`
}

// GetAnimeHistoriesRequest is get anime history request model.
type GetAnimeHistoriesRequest struct {
	StartDate string              `validate:"omitempty,datetime=2006-01-02" mod:"trim"`
	EndDate   string              `validate:"omitempty,datetime=2006-01-02" mod:"trim"`
	Group     entity.HistoryGroup `validate:"oneof=WEEKLY MONTHLY YEARLY" mod:"trim,ucase,default=MONTHLY"`
}

// GetAnimeHistoriesByID to get anime history by id.
func (s *service) GetAnimeHistoriesByID(ctx context.Context, id int64, data GetAnimeHistoriesRequest) ([]AnimeHistory, int, error) {
	if code, err := s.validateID(ctx, id); err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if data.StartDate == "" {
		switch data.Group {
		case entity.Yearly:
			data.StartDate = time.Now().AddDate(-5, 0, 0).Format("2006-01-02")
		case entity.Monthly:
			data.StartDate = time.Now().AddDate(0, -6, 0).Format("2006-01-02")
		case entity.Weekly:
			data.StartDate = time.Now().AddDate(0, -3, 0).Format("2006-01-02")
		}
	}

	histories, code, err := s.anime.GetHistories(ctx, entity.GetHistoriesRequest{
		AnimeID:   id,
		StartDate: utils.ParseToTimePtr("2006-01-02", data.StartDate),
		EndDate:   utils.ParseToTimePtr("2006-01-02", data.EndDate),
		Group:     data.Group,
	})
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	res := make([]AnimeHistory, len(histories))
	for i, h := range histories {
		res[i] = AnimeHistory{
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
