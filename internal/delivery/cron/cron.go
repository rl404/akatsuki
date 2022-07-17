package cron

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/service"
	"github.com/rl404/akatsuki/internal/utils"
)

// Cron contains functions for cron.
type Cron struct {
	service service.Service
}

// New to create new cron.
func New(service service.Service) *Cron {
	return &Cron{
		service: service,
	}
}

func (c *Cron) log(ctx context.Context) {
	if rvr := recover(); rvr != nil {
		errors.Wrap(ctx, fmt.Errorf("%v", rvr), fmt.Errorf("%s", debug.Stack()))
	}

	errStack := errors.Get(ctx)
	if len(errStack) > 0 {
		utils.Error(strings.Join(errStack, ","))
	}
}
