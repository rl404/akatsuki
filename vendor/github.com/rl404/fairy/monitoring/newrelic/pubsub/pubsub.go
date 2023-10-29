package pubsub

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/pubsub"
)

type client struct {
	dialect string
	pubsub  pubsub.PubSub
	nrApp   *newrelic.Application
}

// New to create new newrelic plugin for pubsub.
func New(d string, ps pubsub.PubSub, nrApp *newrelic.Application) pubsub.PubSub {
	return &client{
		dialect: d,
		pubsub:  ps,
		nrApp:   nrApp,
	}
}

// Use to use middelwares.
func (c *client) Use(middlewares ...func(pubsub.HandlerFunc) pubsub.HandlerFunc) {
	c.pubsub.Use(middlewares...)
}

// Publish to publish message.
func (c *client) Publish(ctx context.Context, topic string, data []byte) error {
	segment := newrelic.MessageProducerSegment{
		StartTime:       newrelic.FromContext(ctx).StartSegmentNow(),
		Library:         c.dialect,
		DestinationType: newrelic.MessageTopic,
		DestinationName: topic,
	}
	defer segment.End()

	return c.pubsub.Publish(ctx, topic, data)
}

// Subscribe to subscribe.
func (c *client) Subscribe(ctx context.Context, topic string, handlerFunc pubsub.HandlerFunc) error {
	return c.pubsub.Subscribe(ctx, topic, func(ctx context.Context, message []byte) error {
		tx := c.nrApp.StartTransaction("Consumer " + topic)
		defer tx.End()

		ctx = newrelic.NewContext(ctx, tx)

		if err := handlerFunc(ctx, message); err != nil {
			tx.NoticeError(err)
			return err
		}

		return nil
	})
}

// Close to close connection.
func (c *client) Close() error {
	return c.pubsub.Close()
}
