package pubsub

import (
	"errors"

	"github.com/rl404/fairy/pubsub"
	"github.com/rl404/fairy/pubsub/google"
	"github.com/rl404/fairy/pubsub/rabbitmq"
	"github.com/rl404/fairy/pubsub/redis"
)

// PubsubType is type for pubsub.
type PubsubType int8

// Available types for pubsub.
const (
	Redis PubsubType = iota + 1
	RabbitMQ
	Google
)

// ErrInvalidPubsubType is error for invalid pubsub type.
var ErrInvalidPubsubType = errors.New("invalid pubsub type")

// New to create new pubsub client depends on the type.
func New(pubsubType PubsubType, address string, password string) (pubsub.PubSub, error) {
	switch pubsubType {
	case Redis:
		return redis.New(address, password)
	case RabbitMQ:
		return rabbitmq.New(address)
	case Google:
		return google.New(address, password)
	default:
		return nil, ErrInvalidPubsubType
	}
}
