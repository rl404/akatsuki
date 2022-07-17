package cron

import (
	"context"

	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/utils"
)

// Fill to fill missing anime.
func (c *Cron) Fill(limit int) error {
	ctx := errors.Init(context.Background())
	defer c.log(ctx)

	cnt, _, err := c.service.QueueMissingAnime(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("queued %d anime", cnt)

	return nil
}
