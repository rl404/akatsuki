package service

import (
	"context"
	"net/http"
	"strings"

	"github.com/rl404/akatsuki/internal/domain/mal/entity"
	userEntity "github.com/rl404/akatsuki/internal/domain/user_anime/entity"
	"github.com/rl404/fairy/errors/stack"
)

// UpdateUserAnime to update user anime.
func (s *service) UpdateUserAnime(ctx context.Context, username string) (int, error) {
	if err := s.publisher.PublishParseUserAnime(ctx, strings.ToLower(username), "", true); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err)
	}
	return http.StatusAccepted, nil
}

func (s *service) updateUserAnime(ctx context.Context, username, status string) (int, error) {
	username = strings.ToLower(username)

	var ids []int64
	limit, offset := 500, 0
	for {
		// Call mal api.
		anime, code, err := s.mal.GetUserAnime(ctx, entity.GetUserAnimeRequest{
			Username: username,
			Status:   status,
			Limit:    limit + 1,
			Offset:   offset,
		})
		if err != nil {
			if code == http.StatusNotFound || code == http.StatusForbidden {
				// Delete existing data.
				if code, err := s.userAnime.DeleteByUsername(ctx, username); err != nil {
					return code, stack.Wrap(ctx, err)
				}
				return http.StatusOK, nil
			}
			return code, stack.Wrap(ctx, err)
		}

		for _, a := range anime {
			ids = append(ids, int64(a.Anime.ID))

			// Update user anime data.
			if code, err := s.userAnime.Update(ctx, userEntity.UserAnimeFromMal(ctx, username, a)); err != nil {
				return code, stack.Wrap(ctx, err)
			}

			// Queue related anime.
			if err := s.publisher.PublishParseAnime(ctx, int64(a.Anime.ID), false); err != nil {
				return http.StatusInternalServerError, stack.Wrap(ctx, err)
			}
		}

		if len(anime) <= limit || len(anime) == 0 {
			break
		}

		offset += limit
	}

	// Delete anime not in list.
	if code, err := s.userAnime.DeleteNotInList(ctx, username, ids, userEntity.StrToStatus(status)); err != nil {
		return code, stack.Wrap(ctx, err)
	}

	return http.StatusOK, nil
}
