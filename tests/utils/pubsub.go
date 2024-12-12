package utils_test

import (
	"github.com/rl404/akatsuki/pkg/pubsub"
	_pubsub "github.com/rl404/fairy/pubsub"
)

var pubsubType = map[string]pubsub.PubsubType{
	"rabbitmq": pubsub.RabbitMQ,
	"redis":    pubsub.Redis,
	"google":   pubsub.Google,
}

func GetPubsub(cfg *config) (_pubsub.PubSub, error) {
	return pubsub.New(pubsubType[cfg.PubSub.Dialect], cfg.PubSub.Address, cfg.PubSub.Password)
}
