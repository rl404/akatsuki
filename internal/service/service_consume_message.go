package service

import (
	"context"
	"encoding/json"

	"github.com/rl404/akatsuki/internal/domain/publisher/entity"
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

		isOld, _, err := s.anime.IsOld(ctx, req.ID)
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

	if _, err := s.updateData(ctx, req.ID); err != nil {
		return errors.Wrap(ctx, err)
	}

	return nil
}
