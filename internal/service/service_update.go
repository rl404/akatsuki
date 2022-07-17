package service

import (
	"context"
	"net/http"

	animeEntity "github.com/rl404/akatsuki/internal/domain/anime/entity"
	genreEntity "github.com/rl404/akatsuki/internal/domain/genre/entity"
	publisherEntity "github.com/rl404/akatsuki/internal/domain/publisher/entity"
	studioEntity "github.com/rl404/akatsuki/internal/domain/studio/entity"
	"github.com/rl404/akatsuki/internal/errors"
)

// UpdateOldReleasingAnime to update old releasing anime data.
func (s *service) UpdateOldReleasingAnime(ctx context.Context, limit int) (int, int, error) {
	var cnt int

	releasingAnime, code, err := s.anime.GetOldReleasing(ctx, limit)
	if err != nil {
		return cnt, code, errors.Wrap(ctx, err)
	}

	for _, anime := range releasingAnime {
		if code, err := s.updateData(ctx, anime.ID); err != nil {
			return cnt, code, errors.Wrap(ctx, err)
		}
		cnt++
	}

	return cnt, http.StatusOK, nil
}

// UpdateOldFinishedAnime to update old finished anime data.
func (s *service) UpdateOldFinishedAnime(ctx context.Context, limit int) (int, int, error) {
	var cnt int

	finishedAnime, code, err := s.anime.GetOldFinished(ctx, limit)
	if err != nil {
		return cnt, code, errors.Wrap(ctx, err)
	}

	for _, anime := range finishedAnime {
		if code, err := s.updateData(ctx, anime.ID); err != nil {
			return cnt, code, errors.Wrap(ctx, err)
		}
		cnt++
	}

	return cnt, http.StatusOK, nil
}

// UpdateOldNotYetAnime to update old not yet released anime data.
func (s *service) UpdateOldNotYetAnime(ctx context.Context, limit int) (int, int, error) {
	var cnt int

	finishedAnime, code, err := s.anime.GetOldFinished(ctx, limit)
	if err != nil {
		return cnt, code, errors.Wrap(ctx, err)
	}

	for _, anime := range finishedAnime {
		if code, err := s.updateData(ctx, anime.ID); err != nil {
			return cnt, code, errors.Wrap(ctx, err)
		}
		cnt++
	}

	return cnt, http.StatusOK, nil
}

func (s *service) updateData(ctx context.Context, id int64) (int, error) {
	// Call mal api.
	anime, code, err := s.mal.GetAnimeByID(ctx, int(id))
	if err != nil {
		// If the id is empty.
		if code == http.StatusNotFound {
			if code, err := s.emptyID.Create(ctx, id); err != nil {
				return code, errors.Wrap(ctx, err)
			}
		}
		return code, errors.Wrap(ctx, err)
	}

	// Update genre data.
	if len(anime.Genres) > 0 {
		genres := make([]genreEntity.Genre, len(anime.Genres))
		for i, g := range anime.Genres {
			genres[i] = genreEntity.Genre{
				ID:   int64(g.ID),
				Name: g.Name,
			}
		}

		if code, err := s.genre.BatchUpdate(ctx, genres); err != nil {
			return code, errors.Wrap(ctx, err)
		}
	}

	// Update studio data.
	if len(anime.Studios) > 0 {
		studios := make([]studioEntity.Studio, len(anime.Studios))
		for i, s := range anime.Studios {
			studios[i] = studioEntity.Studio{
				ID:   int64(s.ID),
				Name: s.Name,
			}
		}

		if code, err := s.studio.BatchUpdate(ctx, studios); err != nil {
			return code, errors.Wrap(ctx, err)
		}
	}

	// Update anime data.
	animeE, err := animeEntity.AnimeFromMal(ctx, anime)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, err)
	}

	if code, err := s.anime.Update(ctx, *animeE); err != nil {
		return code, errors.Wrap(ctx, err)
	}

	// Queue related anime.
	for _, r := range anime.RelatedAnime {
		if err := s.publisher.PublishParseAnime(ctx, publisherEntity.ParseAnimeRequest{ID: int64(r.Node.ID)}); err != nil {
			return http.StatusInternalServerError, errors.Wrap(ctx, err)
		}
	}

	return http.StatusOK, nil
}
