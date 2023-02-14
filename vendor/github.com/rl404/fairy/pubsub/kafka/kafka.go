// Package kafka is a wrapper of the original "github.com/segmentio/kafka-go" library.
//
// Only contains basic publish, subscribe, and close methods.
// Data will be encoded to JSON before publishing the message.
package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

// Client is kafka pubsub client.
type Client struct {
	url    string
	writer *kafka.Writer
}

// Channel is kafka subscription channel.
type Channel struct {
	reader *kafka.Reader
}

// New to create new kafka pubsub client.
func New(url string) (*Client, error) {
	return &Client{
		url: url,
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(url),
			AllowAutoTopicCreation: true,
		},
	}, nil
}

// Publish to publish message.
func (c *Client) Publish(ctx context.Context, topic string, data interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return c.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: j,
	})
}

// Subscribe to subscribe queue.
//
// Need to convert the return type to pubsub.Channel.
func (c *Client) Subscribe(ctx context.Context, topic string) (interface{}, error) {
	return &Channel{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{c.url},
			Topic:          topic,
			GroupID:        topic + "-consumer-group",
			GroupBalancers: []kafka.GroupBalancer{kafka.RoundRobinGroupBalancer{}, kafka.RangeGroupBalancer{}},
		}),
	}, nil
}

// Close to close pubsub connection.
func (c *Client) Close() error {
	return c.writer.Close()
}

// Read to read incoming message.
func (c *Channel) Read(ctx context.Context, model interface{}) (<-chan interface{}, <-chan error) {
	msgChan, errChan := make(chan interface{}), make(chan error)
	go func() {
		for {
			if err := c.reader.SetOffsetAt(ctx, time.Now()); err != nil {
				errChan <- err
				return
			}

			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				errChan <- err
			} else {
				if err := json.Unmarshal(msg.Value, &model); err != nil {
					errChan <- err
				} else {
					msgChan <- model
				}
			}
		}
	}()
	return (<-chan interface{})(msgChan), (<-chan error)(errChan)
}

// Close to close subscription.
func (c *Channel) Close() error {
	return c.reader.Close()
}
