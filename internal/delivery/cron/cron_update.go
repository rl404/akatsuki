package cron

import (
	"context"

	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/utils"
)

// Update to update old data.
func (c *Cron) Update(limit int) error {
	ctx := errors.Init(context.Background())
	defer c.log(ctx)

	cnt1, _, err := c.service.QueueOldReleasingAnime(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("queued %d old releasing anime", cnt1)

	cnt2, _, err := c.service.QueueOldFinishedAnime(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("queued %d old finished anime", cnt2)

	cnt3, _, err := c.service.QueueOldNotYetAnime(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("queued %d old not yet released anime", cnt3)

	return nil
}
