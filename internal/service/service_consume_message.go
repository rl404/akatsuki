package service

import (
	"context"
	"encoding/json"
	"net/http"

	animeEntity "github.com/rl404/akatsuki/internal/domain/anime/entity"
	genreEntity "github.com/rl404/akatsuki/internal/domain/genre/entity"
	"github.com/rl404/akatsuki/internal/domain/publisher/entity"
	studioEntity "github.com/rl404/akatsuki/internal/domain/studio/entity"
	"github.com/rl404/akatsuki/internal/errors"
)

// ConsumeMessage to consume message from queue.
// Each message type will be handled differently.
func (s *service) ConsumeMessage(ctx context.Context, data entity.Message) error {
	switch data.Type {
	case entity.TypeParseAnime:
		return errors.Wrap(ctx, s.consumeParseAnime(ctx, data.Data))
	default:
		return errors.Wrap(ctx, errors.ErrInvalidMessageType)
	}
}

func (s *service) consumeParseAnime(ctx context.Context, data []byte) error {
	var req entity.ParseAnimeRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return errors.Wrap(ctx, errors.ErrInvalidRequestFormat)
	}

	if !req.Forced {
		if _, err := s.validateID(ctx, req.ID); err != nil {
			return errors.Wrap(ctx, err)
		}

		// Check if data is still new.
		isOld, _, err := s.anime.IsDataOld(ctx, req.ID)
		if err != nil {
			return errors.Wrap(ctx, err)
		}

		if !isOld {
			return errors.Wrap(ctx, errors.ErrDataStillNew)
		}
	} else {
		// Delete existing empty id.
		if _, err := s.emptyID.Delete(ctx, req.ID); err != nil {
			return errors.Wrap(ctx, err)
		}
	}

	// Call mal api.
	anime, code, err := s.mal.GetAnimeByID(ctx, int(req.ID))
	if err != nil {
		// If the id is empty.
		if code == http.StatusNotFound {
			if _, err := s.emptyID.Create(ctx, req.ID); err != nil {
				return errors.Wrap(ctx, err)
			}
		}
		return errors.Wrap(ctx, err)
	}

	// Update genre data.
	genres := make([]genreEntity.Genre, len(anime.Genres))
	for i, g := range anime.Genres {
		genres[i] = genreEntity.Genre{
			ID:   int64(g.ID),
			Name: g.Name,
		}
	}

	if _, err := s.genre.BatchUpdate(ctx, genres); err != nil {
		return errors.Wrap(ctx, err)
	}

	// Update studio data.
	studios := make([]studioEntity.Studio, len(anime.Studios))
	for i, s := range anime.Studios {
		studios[i] = studioEntity.Studio{
			ID:   int64(s.ID),
			Name: s.Name,
		}
	}

	if _, err := s.studio.BatchUpdate(ctx, studios); err != nil {
		return errors.Wrap(ctx, err)
	}

	// Update anime data.
	animeE, err := animeEntity.AnimeFromMal(anime)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	if _, err := s.anime.Update(ctx, *animeE); err != nil {
		return errors.Wrap(ctx, err)
	}

	// Queue related anime.
	for _, r := range anime.RelatedAnime {
		if err := s.publisher.PublishParseAnime(ctx, entity.ParseAnimeRequest{ID: int64(r.Node.ID)}); err != nil {
			return errors.Wrap(ctx, err)
		}
	}

	return nil
}
