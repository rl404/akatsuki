package errors

import (
	"context"

	"github.com/rl404/fairy/errors/stack"
)

// ErrStacker is error stack interface.
//
// See usage example in example folder.
type ErrStacker interface {
	// Init the context so it can be used for stack.
	Init(ctx context.Context) context.Context
	// Wrap the error and put it in the error stack.
	Wrap(ctx context.Context, err error, errs ...error) error
	// Get the error stack.
	Get(ctx context.Context) interface{}
}

// New to create new error stack.
func New() ErrStacker {
	return stack.New()
}
