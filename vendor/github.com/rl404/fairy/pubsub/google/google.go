// Package google is a wrapper of the original "cloud.google.com/go/pubsub" library.
//
// Only contains basic publish, subscribe, and close methods.
// Data will be encoded to JSON before publishing the message.
package google

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"cloud.google.com/go/pubsub"
)

// Client is google pubsub client.
type Client struct {
	sync.Mutex
	client            *pubsub.Client
	topicExist        map[string]bool
	subscriptionExist map[string]string
}

// Channel is google pubsub channel.
type Channel struct {
	subscription *pubsub.Subscription
}

// New to create new google pubsub client.
//
// Required google service account credential.
// https://cloud.google.com/pubsub/docs/publish-receive-messages-client-library.
//
// If haven't set env "GOOGLE_APPLICATION_CREDENTIALS", you can provide the
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
		client:            client,
		topicExist:        make(map[string]bool),
		subscriptionExist: make(map[string]string),
	}, nil
}

// Publish to publish message.
func (c *Client) Publish(ctx context.Context, topic string, data interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	t, err := c.getTopic(topic)
	if err != nil {
		return err
	}

	res := t.Publish(ctx, &pubsub.Message{
		Data: j,
	})

	if _, err := res.Get(ctx); err != nil {
		return err
	}

	return nil
}

// Subscribe to subscribe to a topic.
//
// Need to convert the return type to pubsub.Channel.
func (c *Client) Subscribe(ctx context.Context, topic string) (interface{}, error) {
	subscription, err := c.getSubscription(topic)
	if err != nil {
		return nil, err
	}

	// Limit to 1 so you can have multiple consumer
	// for the same topic.
	subscription.ReceiveSettings.NumGoroutines = 1
	subscription.ReceiveSettings.MaxOutstandingMessages = 1

	return &Channel{
		subscription: subscription,
	}, nil
}

// Close to close pubsub connection.
func (c *Client) Close() error {
	return c.client.Close()
}

// Read to read incoming message.
func (c *Channel) Read(ctx context.Context, model interface{}) (<-chan interface{}, <-chan error) {
	msgChan, errChan := make(chan interface{}), make(chan error)
	go func() {
		if err := c.subscription.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
			if err := json.Unmarshal(msg.Data, &model); err != nil {
				errChan <- err
			} else {
				msgChan <- model
			}
			msg.Ack()
		}); err != nil {
			errChan <- err
		}
	}()
	return (<-chan interface{})(msgChan), (<-chan error)(errChan)
}

// Close to close subscription.
func (c *Channel) Close() error {
	return nil
}
