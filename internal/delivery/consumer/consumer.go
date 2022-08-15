package consumer

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/akatsuki/internal/domain/publisher/entity"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/service"
	"github.com/rl404/akatsuki/internal/utils"
	"github.com/rl404/fairy/log"
	"github.com/rl404/fairy/pubsub"
)

// Consumer contains functions for consumer.
type Consumer struct {
	service service.Service
	channel pubsub.Channel
}

// New to create new consumer.
func New(service service.Service, ps pubsub.PubSub, topic string) (*Consumer, error) {
	s, err := ps.Subscribe(topic)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		service: service,
		channel: s.(pubsub.Channel),
	}, nil
}

// Subscribe to start subscribing to topic.
func (c *Consumer) Subscribe(nrApp *newrelic.Application) error {
	var msg entity.Message
	msgs, errChan := c.channel.Read(&msg)

	go func() {
		for {
			func() {
				select {
				case <-msgs:
					var err error
					ctx, start := errors.Init(context.Background()), time.Now()
					defer func() {
						if rvr := recover(); rvr != nil {
							err = fmt.Errorf("%v", rvr)
							errors.Wrap(ctx, err, fmt.Errorf("%s", debug.Stack()))
						}

						c.log(ctx, msg, start, err)
					}()

					tx := nrApp.StartTransaction("Consumer " + string(msg.Type))
					defer tx.End()

					ctx = newrelic.NewContext(ctx, tx)

					err = errors.Wrap(ctx, c.service.ConsumeMessage(ctx, msg))
				case err := <-errChan:
					utils.Error(err.Error())
				}
			}()
		}
	}()

	return nil
}

func (c *Consumer) log(ctx context.Context, msg entity.Message, start time.Time, err error) {
	m := map[string]interface{}{
		"level":    log.InfoLevel,
		"type":     msg.Type,
		"data":     string(msg.Data),
		"duration": time.Since(start).String(),
	}

	if err != nil {
		m["level"] = log.ErrorLevel
		errStack := errors.Get(ctx)
		if len(errStack) > 0 {
			m["error"] = strings.Join(errStack, ",")
		}
	}

	utils.Log(m)
}

// Close to stop consumer connection.
func (c *Consumer) Close() error {
	return c.channel.Close()
}
