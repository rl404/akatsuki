package cron

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/akatsuki/internal/utils"
	"github.com/rl404/fairy/errors/stack"
)

// Update to update old data.
func (c *Cron) Update(limit int) error {
	ctx := stack.Init(context.Background())
	defer c.log(ctx)

	tx := c.nrApp.StartTransaction("Cron update")
	defer tx.End()

	ctx = newrelic.NewContext(ctx, tx)

	if err := c.queueOldReleasingAnime(ctx, limit); err != nil {
		return stack.Wrap(ctx, err)
	}

	if err := c.queueOldFinishedAnime(ctx, limit); err != nil {
		return stack.Wrap(ctx, err)
	}

	if err := c.queueOldNotYetAnime(ctx, limit); err != nil {
		return stack.Wrap(ctx, err)
	}

	if err := c.queueOldUsername(ctx, limit); err != nil {
		return stack.Wrap(ctx, err)
	}

	return nil
}

func (c *Cron) queueOldReleasingAnime(ctx context.Context, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldReleasingAnime").End()

	cnt, _, err := c.service.QueueOldReleasingAnime(ctx, limit)
	if err != nil {
		return stack.Wrap(ctx, err)
	}

	utils.Info("queued %d old releasing anime", cnt)
	c.nrApp.RecordCustomEvent("QueueOldReleasingAnime", map[string]interface{}{"count": cnt})

	return nil
}

func (c *Cron) queueOldFinishedAnime(ctx context.Context, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldFinishedAnime").End()

	cnt, _, err := c.service.QueueOldFinishedAnime(ctx, limit)
	if err != nil {
		return stack.Wrap(ctx, err)
	}

	utils.Info("queued %d old finished anime", cnt)
	c.nrApp.RecordCustomEvent("QueueOldFinishedAnime", map[string]interface{}{"count": cnt})

	return nil
}

func (c *Cron) queueOldNotYetAnime(ctx context.Context, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldNotYetAnime").End()

	cnt, _, err := c.service.QueueOldNotYetAnime(ctx, limit)
	if err != nil {
		return stack.Wrap(ctx, err)
	}

	utils.Info("queued %d old not yet released anime", cnt)
	c.nrApp.RecordCustomEvent("QueueOldNotYetAnime", map[string]interface{}{"count": cnt})

	return nil
}

func (c *Cron) queueOldUsername(ctx context.Context, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldUsername").End()

	cnt, _, err := c.service.QueueOldUserAnime(ctx, limit)
	if err != nil {
		return stack.Wrap(ctx, err)
	}

	utils.Info("queued %d old username", cnt)
	c.nrApp.RecordCustomEvent("QueueOldUsername", map[string]interface{}{"count": cnt})

	return nil
}
