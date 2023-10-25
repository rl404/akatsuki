package service

import (
	"context"

	"github.com/rl404/akatsuki/internal/domain/publisher/entity"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/nagato"
)

// ConsumeMessage to consume message from queue.
// Each message type will be handled differently.
func (s *service) ConsumeMessage(ctx context.Context, data entity.Message) error {
	switch data.Type {
	case entity.TypeParseAnime:
		return stack.Wrap(ctx, s.consumeParseAnime(ctx, data))
	case entity.TypeParseUserAnime:
		return stack.Wrap(ctx, s.consumeParseUserAnime(ctx, data))
	default:
		return stack.Wrap(ctx, errors.ErrInvalidMessageType)
	}
}

func (s *service) consumeParseAnime(ctx context.Context, data entity.Message) error {
	if !data.Forced {
		isOld, _, err := s.anime.IsOld(ctx, data.ID)
		if err != nil {
			return stack.Wrap(ctx, err)
		}

		if !isOld {
			return nil
		}
	}

	// Delete existing empty id.
	if _, err := s.emptyID.Delete(ctx, data.ID); err != nil {
		return stack.Wrap(ctx, err)
	}

	if _, err := s.updateAnime(ctx, data.ID); err != nil {
		return stack.Wrap(ctx, err)
	}

	return nil
}

func (s *service) consumeParseUserAnime(ctx context.Context, data entity.Message) error {
	if !data.Forced {
		isOld, _, err := s.userAnime.IsOld(ctx, data.Username)
		if err != nil {
			return stack.Wrap(ctx, err)
		}

		if !isOld {
			return nil
		}
	}

	if data.Status != "" {
		if _, err := s.updateUserAnime(ctx, data.Username, data.Status); err != nil {
			return stack.Wrap(ctx, err)
		}
		return nil
	}

	statuses := []nagato.UserAnimeStatusType{
		nagato.UserAnimeStatusWatching,
		nagato.UserAnimeStatusCompleted,
		nagato.UserAnimeStatusOnHold,
		nagato.UserAnimeStatusDropped,
		nagato.UserAnimeStatusPlanToWatch,
	}

	for _, status := range statuses {
		if err := s.publisher.PublishParseUserAnime(ctx, data.Username, string(status), true); err != nil {
			return stack.Wrap(ctx, err)
		}
	}

	return nil
}
