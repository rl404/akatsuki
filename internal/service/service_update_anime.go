package service

import (
	"context"
	"net/http"

	animeEntity "github.com/rl404/akatsuki/internal/domain/anime/entity"
	genreEntity "github.com/rl404/akatsuki/internal/domain/genre/entity"
	studioEntity "github.com/rl404/akatsuki/internal/domain/studio/entity"
	"github.com/rl404/fairy/errors/stack"
)

// UpdateAnimeByID to update anime by id.
func (s *service) UpdateAnimeByID(ctx context.Context, id int64) (int, error) {
	if err := s.publisher.PublishParseAnime(ctx, id, true); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err)
	}
	return http.StatusAccepted, nil
}

func (s *service) updateAnime(ctx context.Context, id int64) (int, error) {
	// Call mal api.
	anime, code, err := s.mal.GetAnimeByID(ctx, int(id))
	if err != nil {
		if code == http.StatusNotFound {
			// Insert empty id.
			if code, err := s.emptyID.Create(ctx, id); err != nil {
				return code, stack.Wrap(ctx, err)
			}

			// Delete existing data.
			if code, err := s.anime.DeleteByID(ctx, id); err != nil {
				return code, stack.Wrap(ctx, err)
			}

			if code, err := s.userAnime.DeleteByAnimeID(ctx, id); err != nil {
				return code, stack.Wrap(ctx, err)
			}
		}
		return code, stack.Wrap(ctx, err)
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
			return code, stack.Wrap(ctx, err)
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
			return code, stack.Wrap(ctx, err)
		}
	}

	// Update anime data.
	if code, err := s.anime.Update(ctx, animeEntity.AnimeFromMal(ctx, anime)); err != nil {
		return code, stack.Wrap(ctx, err)
	}

	// Queue related anime.
	for _, r := range anime.RelatedAnime {
		if err := s.publisher.PublishParseAnime(ctx, int64(r.Anime.ID), false); err != nil {
			return http.StatusInternalServerError, stack.Wrap(ctx, err)
		}
	}

	return http.StatusOK, nil
}
