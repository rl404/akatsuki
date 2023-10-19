package pubsub

import (
	"context"
	"encoding/json"

	"github.com/rl404/akatsuki/internal/domain/publisher/entity"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/fairy/errors/stack"
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
func (p *Pubsub) PublishParseAnime(ctx context.Context, id int64, forced bool) error {
	d, err := json.Marshal(entity.Message{
		Type:   entity.TypeParseAnime,
		ID:     id,
		Forced: forced,
	})
	if err != nil {
		return stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	if err := p.pubsub.Publish(ctx, p.topic, d); err != nil {
		return stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	return nil
}

// PublishParseUserAnime to publish parse user anime.
func (p *Pubsub) PublishParseUserAnime(ctx context.Context, username, status string, forced bool) error {
	d, err := json.Marshal(entity.Message{
		Type:     entity.TypeParseUserAnime,
		Username: username,
		Status:   status,
		Forced:   forced,
	})
	if err != nil {
		return stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	if err := p.pubsub.Publish(ctx, p.topic, d); err != nil {
		return stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	return nil
}
