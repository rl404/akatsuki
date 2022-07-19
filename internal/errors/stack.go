package errors

import (
	"context"

	"github.com/rl404/fairy/errors"
)

var stacker errors.ErrStacker

func init() {
	stacker = errors.New()
}

// Init to init context for error stack.
func Init(ctx context.Context) context.Context {
	return stacker.Init(ctx)
}

// Wrap to wrap error and put it in the stack.
func Wrap(ctx context.Context, err error, errs ...error) error {
	return stacker.Wrap(ctx, err, errs...)
}

// Get to get error stack.
func Get(ctx context.Context) []string {
	stacks := stacker.Get(ctx).([]string)

	for i, j := 0, len(stacks)-1; i < j; i, j = i+1, j-1 {
		stacks[i], stacks[j] = stacks[j], stacks[i]
	}

	return stacks
}
