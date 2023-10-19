package cron

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/akatsuki/internal/utils"
	"github.com/rl404/fairy/errors/stack"
)

// Fill to fill missing anime.
func (c *Cron) Fill(limit int) error {
	ctx := stack.Init(context.Background())
	defer c.log(ctx)

	tx := c.nrApp.StartTransaction("Cron fill")
	defer tx.End()

	ctx = newrelic.NewContext(ctx, tx)

	if err := c.queueMissingAnime(ctx, limit); err != nil {
		tx.NoticeError(err)
		return stack.Wrap(ctx, err)
	}

	return nil
}

func (c *Cron) queueMissingAnime(ctx context.Context, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueMissingAnime").End()

	cnt, _, err := c.service.QueueMissingAnime(ctx, limit)
	if err != nil {
		return stack.Wrap(ctx, err)
	}

	utils.Info("queued %d anime", cnt)
	c.nrApp.RecordCustomEvent("QueueMissingAnime", map[string]interface{}{"count": cnt})

	return nil
}
