package limit

import (
	"errors"
	"time"

	"github.com/rl404/fairy/limit/atomic"
	"github.com/rl404/fairy/limit/mutex"
)

// Limiter is rate limiting interface.
//
// See usage example in example folder.
type Limiter interface {
	// To add count to rate limiter.
	// Will block until the rate in below the
	// limit.
	Take()
}

// LimitType is type for rate limit.
type LimitType int8

// Available types for rate limit.
const (
	Mutex LimitType = iota
	Atomic
)

// ErrInvalidLimitType is error for invalid rate limit type.
var ErrInvalidLimitype = errors.New("invalid rate limit type")

// New to create new rate limiter.
func New(limitType LimitType, rate int, interval time.Duration) (Limiter, error) {
	switch limitType {
	case Mutex:
		return mutex.New(rate, interval), nil
	case Atomic:
		return atomic.New(rate, interval), nil
	default:
		return nil, ErrInvalidLimitype
	}
}
