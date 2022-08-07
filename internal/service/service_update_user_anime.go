package service

import (
	"context"
	"net/http"

	"github.com/rl404/akatsuki/internal/domain/mal/entity"
	publisherEntity "github.com/rl404/akatsuki/internal/domain/publisher/entity"
	userEntity "github.com/rl404/akatsuki/internal/domain/user_anime/entity"
	"github.com/rl404/akatsuki/internal/errors"
)

// UpdateUserAnime to update user anime.
func (s *service) UpdateUserAnime(ctx context.Context, username string) (int, error) {
	if err := s.publisher.PublishParseUserAnime(ctx, publisherEntity.ParseUserAnimeRequest{
		Username: username,
		Forced:   true,
	}); err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, err)
	}
	return http.StatusAccepted, nil
}

func (s *service) updateUserAnime(ctx context.Context, username string) (int, error) {
	var ids []int64
	limit, offset := 500, 0
	for {
		// Call mal api.
		anime, code, err := s.mal.GetUserAnime(ctx, entity.GetUserAnimeRequest{
			Username: username,
			Limit:    limit + 1,
			Offset:   offset,
		})
		if err != nil {
			return code, errors.Wrap(ctx, err)
		}

		for _, a := range anime {
			ids = append(ids, int64(a.Anime.ID))

			// Update user anime data.
			animeE, err := userEntity.UserAnimeFromMal(ctx, username, a)
			if err != nil {
				return http.StatusInternalServerError, errors.Wrap(ctx, err)
			}

			if code, err := s.userAnime.Update(ctx, *animeE); err != nil {
				return code, errors.Wrap(ctx, err)
			}

			// Queue related anime.
			if err := s.publisher.PublishParseAnime(ctx, publisherEntity.ParseAnimeRequest{ID: int64(a.Anime.ID)}); err != nil {
				return http.StatusInternalServerError, errors.Wrap(ctx, err)
			}
		}

		if len(anime) <= limit || len(anime) == 0 {
			break
		}

		offset += limit
	}

	// Delete anime not in list.
	if code, err := s.userAnime.DeleteNotInList(ctx, username, ids); err != nil {
		return code, errors.Wrap(ctx, err)
	}

	return http.StatusOK, nil
}
