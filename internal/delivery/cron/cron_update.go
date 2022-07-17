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

	cnt1, _, err := c.service.UpdateAiringAnime(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("updated %d airing anime", cnt1)

	cnt2, _, err := c.service.UpdateOldData(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("updated %d old data", cnt2)

	return nil
}
