// Package rabbitmq is a wrapper of the original "github.com/streadway/amqp" library.
//
// todo: add reconnect feature.
package rabbitmq

import (
	"context"

	"github.com/rl404/fairy/pubsub"
	"github.com/streadway/amqp"
)

// Client is rabbitmq pubsub client.
type Client struct {
	client      *amqp.Connection
	middlewares []func(pubsub.HandlerFunc) pubsub.HandlerFunc
}

// New to create new rabbitmq pubsub client.
func New(url string) (*Client, error) {
	c, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Client{client: c}, nil
}

// Use to add pubsub middlewares.
func (c *Client) Use(middlewares ...func(pubsub.HandlerFunc) pubsub.HandlerFunc) {
	c.middlewares = append(c.middlewares, middlewares...)
}

func (c *Client) applyMiddlewares(handlerFunc pubsub.HandlerFunc) pubsub.HandlerFunc {
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		handlerFunc = c.middlewares[i](handlerFunc)
	}
	return handlerFunc
}

// Publish to publish message.
func (c *Client) Publish(ctx context.Context, queue string, data []byte) error {
	ch, err := c.client.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	if err := ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        data,
	}); err != nil {
		return err
	}

	return nil
}

// Subscribe to subscribe queue.
func (c *Client) Subscribe(ctx context.Context, queue string, handlerFunc pubsub.HandlerFunc) error {
	ch, err := c.client.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func(msgs <-chan amqp.Delivery, h pubsub.HandlerFunc) {
		h = c.applyMiddlewares(h)

		for msg := range msgs {
			h(ctx, msg.Body)
		}
	}(msgs, handlerFunc)

	return nil
}

// Close to close pubsub connection.
func (c *Client) Close() error {
	return c.client.Close()
}
