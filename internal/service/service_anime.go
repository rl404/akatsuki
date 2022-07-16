package service

import (
	"context"
	"net/http"
	"time"

	animeEntity "github.com/rl404/akatsuki/internal/domain/anime/entity"
	"github.com/rl404/akatsuki/internal/domain/publisher/entity"
	"github.com/rl404/akatsuki/internal/errors"
)

// Anime is anime model.
type Anime struct {
	ID                int64              `json:"id"`
	Title             string             `json:"title"`
	AlternativeTitles alternativeTitles  `json:"alternative_titles"`
	Picture           string             `json:"picture"`
	StartDate         date               `json:"start_date"`
	EndDate           date               `json:"end_date"`
	Synopsis          string             `json:"synopsis"`
	Background        string             `json:"background"`
	NSFW              bool               `json:"nsfw"`
	Type              animeEntity.Type   `json:"type"`
	Status            animeEntity.Status `json:"status"`
	Episode           episode            `json:"episode"`
	Season            *season            `json:"season"`
	Broadcast         *broadcast         `json:"broadcast"`
	Source            animeEntity.Source `json:"source"`
	Rating            animeEntity.Rating `json:"rating"`
	Mean              float64            `json:"mean"`
	Rank              int                `json:"rank"`
	Popularity        int                `json:"popularity"`
	Member            int                `json:"member"`
	Voter             int                `json:"voter"`
	Stats             stats              `json:"stats"`
	Genres            []genre            `json:"genres"`
	Pictures          []string           `json:"pictures"`
	Related           []related          `json:"related"`
	Studios           []studio           `json:"studio"`
	UpdatedAt         time.Time          `json:"updated_at"`
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
			if err := s.publisher.PublishParseAnime(ctx, entity.ParseAnimeRequest{ID: id}); err != nil {
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
		relatedMap := make(map[int64]animeEntity.Relation)
		relatedIDs := make([]int64, len(animeDB.Related))
		for i, r := range animeDB.Related {
			relatedIDs[i] = r.ID
			relatedMap[r.ID] = r.Relation
		}

		relateds, code, err := s.anime.GetByIDs(ctx, relatedIDs)
		if err != nil {
			return nil, code, errors.Wrap(ctx, err)
		}

		anime.Related = make([]related, len(relateds))
		for i, r := range relateds {
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

	isEmpty, code, err := s.emptyID.IsEmpty(ctx, id)
	if err != nil {
		return code, errors.Wrap(ctx, err)
	}

	if isEmpty {
		return http.StatusNotFound, errors.Wrap(ctx, errors.ErrAnimeNotFound)
	}

	return http.StatusOK, nil
}
