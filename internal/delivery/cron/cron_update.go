package cron

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/utils"
)

// Update to update old data.
func (c *Cron) Update(nrApp *newrelic.Application, limit int) error {
	ctx := errors.Init(context.Background())
	defer c.log(ctx)

	tx := nrApp.StartTransaction("Cron update")
	defer tx.End()

	ctx = newrelic.NewContext(ctx, tx)

	if err := c.queueOldReleasingAnime(ctx, nrApp, limit); err != nil {
		return errors.Wrap(ctx, err)
	}

	if err := c.queueOldFinishedAnime(ctx, nrApp, limit); err != nil {
		return errors.Wrap(ctx, err)
	}

	if err := c.queueOldNotYetAnime(ctx, nrApp, limit); err != nil {
		return errors.Wrap(ctx, err)
	}

	if err := c.queueOldUsername(ctx, nrApp, limit); err != nil {
		return errors.Wrap(ctx, err)
	}

	return nil
}

func (c *Cron) queueOldReleasingAnime(ctx context.Context, nrApp *newrelic.Application, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldReleasingAnime").End()

	cnt, _, err := c.service.QueueOldReleasingAnime(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("queued %d old releasing anime", cnt)
	nrApp.RecordCustomEvent("QueueOldReleasingAnime", map[string]interface{}{"count": cnt})

	return nil
}

func (c *Cron) queueOldFinishedAnime(ctx context.Context, nrApp *newrelic.Application, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldFinishedAnime").End()

	cnt, _, err := c.service.QueueOldFinishedAnime(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("queued %d old finished anime", cnt)
	nrApp.RecordCustomEvent("QueueOldFinishedAnime", map[string]interface{}{"count": cnt})

	return nil
}

func (c *Cron) queueOldNotYetAnime(ctx context.Context, nrApp *newrelic.Application, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldNotYetAnime").End()

	cnt, _, err := c.service.QueueOldNotYetAnime(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("queued %d old not yet released anime", cnt)
	nrApp.RecordCustomEvent("QueueOldNotYetAnime", map[string]interface{}{"count": cnt})

	return nil
}

func (c *Cron) queueOldUsername(ctx context.Context, nrApp *newrelic.Application, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldUsername").End()

	cnt, _, err := c.service.QueueOldUserAnime(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("queued %d old username", cnt)
	nrApp.RecordCustomEvent("QueueOldUsername", map[string]interface{}{"count": cnt})

	return nil
}
