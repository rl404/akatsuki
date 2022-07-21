package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rl404/akatsuki/internal/domain/publisher/entity"
	"github.com/rl404/akatsuki/internal/errors"
)

// ConsumeMessage to consume message from queue.
// Each message type will be handled differently.
func (s *service) ConsumeMessage(ctx context.Context, data entity.Message) error {
	switch data.Type {
	case entity.TypeParseAnime:
		return errors.Wrap(ctx, s.consumeParseAnime(ctx, data.Data))
	case entity.TypeParseUserAnime:
		return errors.Wrap(ctx, s.consumeParseUserAnime(ctx, data.Data))
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
		if code, err := s.validateID(ctx, req.ID); err != nil {
			if code == http.StatusNotFound {
				return nil
			}
			return errors.Wrap(ctx, err)
		}

		isOld, _, err := s.anime.IsOld(ctx, req.ID)
		if err != nil {
			return errors.Wrap(ctx, err)
		}

		if !isOld {
			return nil
		}
	} else {
		// Delete existing empty id.
		if _, err := s.emptyID.Delete(ctx, req.ID); err != nil {
			return errors.Wrap(ctx, err)
		}
	}

	if _, err := s.updateAnime(ctx, req.ID); err != nil {
		return errors.Wrap(ctx, err)
	}

	return nil
}

func (s *service) consumeParseUserAnime(ctx context.Context, data []byte) error {
	var req entity.ParseUserAnimeRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return errors.Wrap(ctx, errors.ErrInvalidRequestFormat)
	}

	if !req.Forced {
		isOld, _, err := s.userAnime.IsOld(ctx, req.Username)
		if err != nil {
			return errors.Wrap(ctx, err)
		}
		if !isOld {
			return nil
		}
	}

	if _, err := s.updateUserAnime(ctx, req.Username); err != nil {
		return errors.Wrap(ctx, err)
	}

	return nil
}
