package pubsub

import (
	"context"
	"errors"

	"github.com/rl404/fairy/pubsub/google"
	"github.com/rl404/fairy/pubsub/nsq"
	"github.com/rl404/fairy/pubsub/rabbitmq"
	"github.com/rl404/fairy/pubsub/redis"
)

// PubSub is pubsub interface.
//
// For subscribe function, you have to convert
// the return type to Channel.
//
// See usage example in example folder.
type PubSub interface {
	// Publish message to specific topic/channel.
	// Data will be encoded first before publishing.
	Publish(ctx context.Context, topic string, data interface{}) error
	// Subscribe to specific topic/channel.
	Subscribe(ctx context.Context, topic string) (interface{}, error)
	// Close pubsub client connection.
	Close() error
}

// Channel is channel interface.
//
// See usage example in example folder.
type Channel interface {
	// Read and process incoming message. Param `data` should
	// be a pointer just like when using json.Unmarshal.
	Read(ctx context.Context, data interface{}) (<-chan interface{}, <-chan error)
	// Close subscription.
	Close() error
}

// PubsubType is type for pubsub.
type PubsubType int8

// Available types for pubsub.
const (
	Redis PubsubType = iota + 1
	RabbitMQ
	NSQ
	Google
)

// ErrInvalidPubsubType is error for invalid pubsub type.
var ErrInvalidPubsubType = errors.New("invalid pubsub type")

// New to create new pubsub client depends on the type.
func New(pubsubType PubsubType, address string, password string) (PubSub, error) {
	switch pubsubType {
	case Redis:
		return redis.New(address, password)
	case RabbitMQ:
		return rabbitmq.New(address)
	case NSQ:
		return nsq.New(address)
	case Google:
		return google.New(address, password)
	default:
		return nil, ErrInvalidPubsubType
	}
}
