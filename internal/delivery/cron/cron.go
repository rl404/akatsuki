package cron

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/akatsuki/internal/service"
	"github.com/rl404/akatsuki/internal/utils"
	"github.com/rl404/akatsuki/pkg/log"
	"github.com/rl404/fairy/errors/stack"
)

// Cron contains functions for cron.
type Cron struct {
	service service.Service
	nrApp   *newrelic.Application
}

// New to create new cron.
func New(service service.Service, nrApp *newrelic.Application) *Cron {
	return &Cron{
		service: service,
		nrApp:   nrApp,
	}
}

func (c *Cron) log(ctx context.Context) {
	if rvr := recover(); rvr != nil {
		stack.Wrap(ctx, fmt.Errorf("%s", debug.Stack()), fmt.Errorf("%v", rvr), fmt.Errorf("panic"))
	}

	errStack := stack.Get(ctx)
	if len(errStack) > 0 {
		utils.Log(map[string]interface{}{
			"level": log.ErrorLevel,
			"error": errStack,
		})
	}
}
