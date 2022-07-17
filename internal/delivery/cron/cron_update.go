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

	cnt1, _, err := c.service.UpdateOldReleasingAnime(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("updated %d old releasing anime", cnt1)

	cnt2, _, err := c.service.UpdateOldFinishedAnime(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("updated %d old finished anime", cnt2)

	cnt3, _, err := c.service.UpdateOldNotYetAnime(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("updated %d old not yet released anime", cnt3)

	return nil
}
