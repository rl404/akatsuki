package pubsub

import (
	"context"
	"encoding/json"

	"github.com/rl404/akatsuki/internal/domain/publisher/entity"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/fairy/pubsub"
)

// Pubsub contains functions for pubsub.
type Pubsub struct {
	pubsub pubsub.PubSub
	topic  string
}

// New to create new pubsub.
func New(ps pubsub.PubSub, topic string) *Pubsub {
	return &Pubsub{
		pubsub: ps,
		topic:  topic,
	}
}

// PublishParseAnime to publish parse anime.
func (p *Pubsub) PublishParseAnime(ctx context.Context, data entity.ParseAnimeRequest) error {
	d, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(ctx, errors.ErrInternalServer, err)
	}

	if err := p.pubsub.Publish(p.topic, entity.Message{
		Type: entity.TypeParseAnime,
		Data: d,
	}); err != nil {
		return errors.Wrap(ctx, errors.ErrInternalServer, err)
	}

	return nil
}

// PublishParseUserAnime to publish parse user anime.
func (p *Pubsub) PublishParseUserAnime(ctx context.Context, data entity.ParseUserAnimeRequest) error {
	d, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(ctx, errors.ErrInternalServer, err)
	}

	if err := p.pubsub.Publish(p.topic, entity.Message{
		Type: entity.TypeParseUserAnime,
		Data: d,
	}); err != nil {
		return errors.Wrap(ctx, errors.ErrInternalServer, err)
	}

	return nil
}
