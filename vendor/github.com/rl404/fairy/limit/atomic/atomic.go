// Package atomic is a rate limiter based on "github.com/uber-go/ratelimit" library.
package atomic

import (
	"sync/atomic"
	"time"
	"unsafe"
)

type atomicLimit struct {
	state unsafe.Pointer
	// Padding is unused but it is crucial to maintain performance
	// of this rate limiter in case of collocation with other frequently accessed memory.
	padding [56]byte // cache line size - state pointer size = 64 - 8; created to avoid false sharing.

	perRequest time.Duration
	maxSlack   time.Duration
	clock      clocker
}

type clocker interface {
	Now() time.Time
	Sleep(time.Duration)
}

type state struct {
	last     time.Time
	sleepFor time.Duration
}

// New to create new atomic based limiter.
func New(rate int, interval time.Duration) *atomicLimit {
	perRequest := interval / time.Duration(rate)
	l := &atomicLimit{
		perRequest: perRequest,
		maxSlack:   -1 * time.Duration(10) * perRequest,
		clock:      newClock(),
	}
	initialState := state{
		last:     time.Time{},
		sleepFor: 0,
	}
	atomic.StorePointer(&l.state, unsafe.Pointer(&initialState))
	return l
}

// Take blocks to ensure that the time spent between multiple
// Take calls is on average time.Second/rate.
func (t *atomicLimit) Take() {
	var (
		newState state
		taken    bool
		interval time.Duration
	)
	for !taken {
		now := t.clock.Now()

		previousStatePointer := atomic.LoadPointer(&t.state)
		oldState := (*state)(previousStatePointer)

		newState = state{
			last:     now,
			sleepFor: oldState.sleepFor,
		}

		// If this is our first request, then we allow it.
		if oldState.last.IsZero() {
			taken = atomic.CompareAndSwapPointer(&t.state, previousStatePointer, unsafe.Pointer(&newState))
			continue
		}

		// sleepFor calculates how much time we should sleep based on
		// the perRequest budget and how long the last request took.
		// Since the request may take longer than the budget, this number
		// can get negative, and is summed across requests.
		newState.sleepFor += t.perRequest - now.Sub(oldState.last)
		// We shouldn't allow sleepFor to get too negative, since it would mean that
		// a service that slowed down a lot for a short period of time would get
		// a much higher RPS following that.
		if newState.sleepFor < t.maxSlack {
			newState.sleepFor = t.maxSlack
		}
		if newState.sleepFor > 0 {
			newState.last = newState.last.Add(newState.sleepFor)
			interval, newState.sleepFor = newState.sleepFor, 0
		}
		taken = atomic.CompareAndSwapPointer(&t.state, previousStatePointer, unsafe.Pointer(&newState))
	}
	t.clock.Sleep(interval)
}
