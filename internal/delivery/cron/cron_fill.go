package cron

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/utils"
)

// Fill to fill missing anime.
func (c *Cron) Fill(nrApp *newrelic.Application, limit int) error {
	ctx := errors.Init(context.Background())
	defer c.log(ctx)

	tx := nrApp.StartTransaction("Cron fill")
	defer tx.End()

	ctx = newrelic.NewContext(ctx, tx)

	if err := c.queueMissingAnime(ctx, nrApp, limit); err != nil {
		tx.NoticeError(err)
		return errors.Wrap(ctx, err)
	}

	return nil
}

func (c *Cron) queueMissingAnime(ctx context.Context, nrApp *newrelic.Application, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueMissingAnime").End()

	cnt, _, err := c.service.QueueMissingAnime(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("queued %d anime", cnt)
	nrApp.RecordCustomEvent("QueueMissingAnime", map[string]interface{}{"count": cnt})

	return nil
}
