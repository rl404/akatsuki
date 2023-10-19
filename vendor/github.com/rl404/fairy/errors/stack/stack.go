// Package stack is an error wrapper and put it to a stack
// using passed context.
package stack

import (
	"context"
	"fmt"
	"runtime"
)

type errStackCtx struct{}

type errStackCtxValue struct {
	stacks []ErrStack
}

// ErrStack is error stack.
type ErrStack struct {
	File     string `json:"file"`
	Function string `json:"function"`
	Message  string `json:"message"`
}

// Init to init context so it can be used as error stack.
func Init(ctx context.Context) context.Context {
	if _, ok := ctx.Value(errStackCtx{}).(*errStackCtxValue); ok {
		return ctx
	}
	return context.WithValue(ctx, errStackCtx{}, &errStackCtxValue{})
}

// Wrap to put the errors in the stack in the initiated context.
//
// Will return the last of maskedErrs.
func Wrap(ctx context.Context, originalErr error, maskedErrs ...error) error {
	if originalErr == nil && len(maskedErrs) == 0 {
		return nil
	}

	errs := append([]error{originalErr}, maskedErrs...)
	lastErr := errs[len(errs)-1]

	ctxValue, ok := ctx.Value(errStackCtx{}).(*errStackCtxValue)
	if !ok {
		// If ctx is not initiated, just return the last error.
		return lastErr
	}

	pc, file, line, _ := runtime.Caller(1)

	funcName := runtime.FuncForPC(pc).Name()
	fileLine := fmt.Sprintf("%s:%d", file, line)

	// Add to stack.
	for _, err := range errs {
		msg := ""
		if err != nil {
			msg = err.Error()
		}

		ctxValue.stacks = append(ctxValue.stacks, ErrStack{
			File:     fileLine,
			Function: funcName,
			Message:  msg,
		})
	}

	return lastErr
}

// Get to get error stacks.
func Get(ctx context.Context) []ErrStack {
	ctxValue, ok := ctx.Value(errStackCtx{}).(*errStackCtxValue)
	if !ok {
		return nil
	}
	return ctxValue.stacks
}
