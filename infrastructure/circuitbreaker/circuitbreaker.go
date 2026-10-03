// Package circuitbreaker provides a simple circuit breaker pattern for HTTP clients.
//
// The circuit breaker wraps an http.RoundTripper and tracks consecutive failures.
// After a configurable threshold of consecutive failures, the circuit opens and
// subsequent requests fail fast without calling the underlying transport.
// After a configurable timeout, the circuit half-opens and allows one probe
// request through. If it succeeds, the circuit closes; otherwise it reopens.
package circuitbreaker

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"
)

// State represents the circuit breaker state.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // failing fast
	StateHalfOpen              // probe request allowed
)

// ErrCircuitOpen is returned when the circuit breaker is open.
var ErrCircuitOpen = errors.New("circuit breaker is open")

// Config configures the circuit breaker behaviour.
type Config struct {
	// FailureThreshold is the number of consecutive failures before the circuit opens.
	FailureThreshold int

	// Timeout is how long the circuit stays open before transitioning to half-open.
	Timeout time.Duration
}

// DefaultConfig returns a sensible default configuration.
func DefaultConfig() Config {
	return Config{
		FailureThreshold: 5,
		Timeout:          30 * time.Second,
	}
}

// Breaker implements the circuit breaker pattern for HTTP round trips.
type Breaker struct {
	config Config

	mu              sync.Mutex
	state           State
	failureCount    int
	lastFailureTime time.Time
	// now is an injectable clock used by Allow; when nil, time.Now is used.
	// Tests set it to a fixed value so the timeout boundary is deterministic.
	now func() time.Time
}

// State returns the current circuit breaker state.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

// Allow checks whether a request is allowed through the circuit breaker.
// Returns nil if allowed, ErrCircuitOpen if the circuit is open.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateOpen:
		now := b.now
		if now == nil {
			now = time.Now
		}
		if now().Sub(b.lastFailureTime) > b.config.Timeout {
			b.state = StateHalfOpen
			return nil
		}
		return ErrCircuitOpen
	default:
		return nil
	}
}

// RecordResult records the outcome of a request. Call this after making the
// actual HTTP request to update the circuit breaker state.
func (b *Breaker) RecordResult(err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if err == nil {
		// Success resets the breaker
		b.failureCount = 0
		if b.state == StateHalfOpen {
			b.state = StateClosed
		}
		return
	}

	// Only count connection-level failures, not HTTP application errors.
	// HTTP errors (4xx, 5xx) come from the NetBox server and are not
	// circuit-breaking events.
	if isConnectionError(err) {
		b.failureCount++
		b.lastFailureTime = time.Now()
		if b.failureCount >= b.config.FailureThreshold {
			b.state = StateOpen
		}
	}
}

// isConnectionError returns true if the error is likely a connection-level issue.
func isConnectionError(err error) bool {
	if err == nil {
		return false
	}

	// Connection refused, DNS failures, timeouts, etc. typically manifest
	// as http-level errors from the http.Client. Application-level errors
	// (4xx/5xx responses) return nil error from RoundTrip — they come
	// back as a non-nil Response with an error status code.
	// So any non-nil error from RoundTrip is a transport failure,
	// except for context cancellation which is not a connection issue.
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	return true
}

// RoundTripper wraps an http.RoundTripper with circuit breaker protection.
// Use this as the Transport in an http.Client.
type RoundTripper struct {
	http.RoundTripper

	breaker *Breaker
}

// NewRoundTripper wraps an http.RoundTripper with a circuit breaker.
func NewRoundTripper(inner http.RoundTripper, cfg Config) *RoundTripper {
	return &RoundTripper{
		RoundTripper: inner,
		breaker: &Breaker{
			config: cfg,
			state:  StateClosed,
		},
	}
}

// Breaker returns the underlying circuit breaker instance.
func (rt *RoundTripper) Breaker() *Breaker {
	return rt.breaker
}

// RoundTrip implements http.RoundTripper.
func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := rt.breaker.Allow(); err != nil {
		return nil, err
	}

	resp, err := rt.RoundTripper.RoundTrip(req)
	rt.breaker.RecordResult(err)

	//nolint:wrapcheck
	return resp, err
}
