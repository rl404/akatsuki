// Package stack is an error wrapper and put it to a stack
// using passed context.
package stack

import (
	"context"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type errctx struct{}
type errctxvalue struct {
	m map[string]struct{}
	l []string
}

// Stacker is error stack client.
type Stacker struct{}

// New to create new error stacker.
func New() *Stacker {
	return &Stacker{}
}

// Init to init context so it can be used
// as an error stack.
func (s *Stacker) Init(ctx context.Context) context.Context {
	if _, ok := ctx.Value(errctx{}).(*errctxvalue); ok {
		return ctx
	}
	return context.WithValue(ctx, errctx{}, &errctxvalue{})
}

// Wrap to wrap the errors and put them in the stack in the
// initiated context.
func (s *Stacker) Wrap(ctx context.Context, err error, errs ...error) error {
	if err != nil {
		s.wrap(ctx, err, errs...)
	}
	return err
}

func (s *Stacker) wrap(ctx context.Context, err error, errs ...error) {
	out, ok := ctx.Value(errctx{}).(*errctxvalue)
	if !ok {
		return
	}

	if out.l == nil {
		out.l = make([]string, 0)
	}

	if out.m == nil {
		out.m = make(map[string]struct{})
	}

	_, f, l, _ := runtime.Caller(3)
	caller := s.filename(f) + ":" + strconv.Itoa(l)

	if err.Error() == "" {
		out.l = append(out.l, caller)
		return
	}

	errs = s.prependErr(errs, err)

	for i := len(errs) - 1; i >= 0; i-- {
		eStr := errs[i].Error()
		if _, ok := out.m[eStr]; !ok {
			out.m[eStr] = struct{}{}
			out.l = append(out.l, caller+" "+eStr)
		} else {
			out.l = append(out.l, caller)
		}
	}
}

func (s *Stacker) filename(fpath string) string {
	if i := strings.LastIndexByte(fpath, filepath.Separator); i >= 0 {
		return fpath[i+1:]
	}
	return fpath
}

func (s *Stacker) prependErr(errs []error, err error) []error {
	return append([]error{err}, errs...)
}

// Get returns the errors stack from context.
// Will return stack starts from the first/deepest wrapped error.
func (s *Stacker) Get(ctx context.Context) interface{} {
	if value, ok := ctx.Value(errctx{}).(*errctxvalue); ok {
		return value.l
	}
	return nil
}
