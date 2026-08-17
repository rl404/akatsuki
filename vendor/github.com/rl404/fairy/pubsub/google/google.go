// Package google is a wrapper of the original "cloud.google.com/go/pubsub" library.
package google

import (
	"context"
	"os"
	"sync"

	"cloud.google.com/go/pubsub/v2"
	_pubsub "github.com/rl404/fairy/pubsub"
)

// Client is google pubsub client.
type Client struct {
	sync.Mutex
	projectID         string
	client            *pubsub.Client
	middlewares       []func(_pubsub.HandlerFunc) _pubsub.HandlerFunc
	topicExist        map[string]bool
	subscriptionExist map[string]string
}

// New to create new google pubsub client.
//
// Required google service account credential.
// https://cloud.google.com/pubsub/docs/publish-receive-messages-client-library.
//
// If you haven't set env "GOOGLE_APPLICATION_CREDENTIALS", you can provide the
// credential json file path in the param.
func New(projectID, serviceAccountCredentialPath string) (*Client, error) {
	if serviceAccountCredentialPath != "" {
		if err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", serviceAccountCredentialPath); err != nil {
			return nil, err
		}
	}

	client, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		return nil, err
	}

	return &Client{
		projectID:         projectID,
		client:            client,
		topicExist:        make(map[string]bool),
		subscriptionExist: make(map[string]string),
	}, nil
}

// Use to add pubsub middlewares.
func (c *Client) Use(middlewares ...func(_pubsub.HandlerFunc) _pubsub.HandlerFunc) {
	c.middlewares = append(c.middlewares, middlewares...)
}

func (c *Client) applyMiddlewares(handlerFunc _pubsub.HandlerFunc) _pubsub.HandlerFunc {
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		handlerFunc = c.middlewares[i](handlerFunc)
	}
	return handlerFunc
}

// Publish to publish message.
func (c *Client) Publish(ctx context.Context, topic string, data []byte) error {
	publisher, err := c.getPublisher(topic)
	if err != nil {
		return err
	}

	if _, err := publisher.Publish(ctx, &pubsub.Message{
		Data: data,
	}).Get(ctx); err != nil {
		return err
	}

	return nil
}

// Subscribe to subscribe topic.
func (c *Client) Subscribe(ctx context.Context, topic string, handlerFunc _pubsub.HandlerFunc) error {
	subscriber, err := c.getSubscriber(topic)
	if err != nil {
		return err
	}

	// Limit to 1 so you can have multiple consumer
	// for the same topic.
	subscriber.ReceiveSettings.NumGoroutines = 1
	subscriber.ReceiveSettings.MaxOutstandingMessages = 1

	go func(s *pubsub.Subscriber, h _pubsub.HandlerFunc) {
		h = c.applyMiddlewares(h)
		if err := s.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			msg.Ack()
			h(ctx, msg.Data)
		}); err != nil {
			panic(err)
		}
	}(subscriber, handlerFunc)

	return nil
}

// Close to close client.
func (c *Client) Close() error {
	return c.client.Close()
}
