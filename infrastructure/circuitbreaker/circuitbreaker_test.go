package circuitbreaker

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBreaker_InitialState(t *testing.T) {
	b := &Breaker{config: DefaultConfig(), state: StateClosed}
	if b.State() != StateClosed {
		t.Errorf("initial state = %d, want %d", b.State(), StateClosed)
	}
}

func TestBreaker_AllowWhenClosed(t *testing.T) {
	b := &Breaker{config: DefaultConfig(), state: StateClosed}
	if err := b.Allow(); err != nil {
		t.Errorf("Allow() = %v, want nil", err)
	}
}

func TestBreaker_OpensAfterFailures(t *testing.T) {
	b := &Breaker{
		config: Config{FailureThreshold: 3, Timeout: 1 * time.Hour},
		state:  StateClosed,
	}

	// First two failures
	b.RecordResult(errors.New("connection refused"))
	b.RecordResult(errors.New("connection refused"))
	if b.State() != StateClosed {
		t.Errorf("state after 2 failures = %d, want %d", b.State(), StateClosed)
	}

	// Third failure should open
	b.RecordResult(errors.New("connection refused"))
	if b.State() != StateOpen {
		t.Errorf("state after 3 failures = %d, want %d", b.State(), StateOpen)
	}

	// Allow should fail
	if err := b.Allow(); err != ErrCircuitOpen {
		t.Errorf("Allow() = %v, want %v", err, ErrCircuitOpen)
	}
}

func TestBreaker_ClosesAfterSuccess(t *testing.T) {
	b := &Breaker{
		config: Config{FailureThreshold: 2, Timeout: 1 * time.Hour},
		state:  StateClosed,
	}

	b.RecordResult(errors.New("timeout"))
	b.RecordResult(errors.New("timeout"))
	if b.State() != StateOpen {
		t.Errorf("state after 2 failures = %d, want %d", b.State(), StateOpen)
	}

	// Force half-open (simulating timeout)
	b.lastFailureTime = time.Now().Add(-2 * time.Hour)
	if err := b.Allow(); err != nil {
		t.Errorf("Allow() after timeout = %v, want nil", err)
	}
	if b.State() != StateHalfOpen {
		t.Errorf("state after half-open = %d, want %d", b.State(), StateHalfOpen)
	}

	// Success should close
	b.RecordResult(nil)
	if b.State() != StateClosed {
		t.Errorf("state after success = %d, want %d", b.State(), StateClosed)
	}
}

func TestBreaker_SuccessResetsCounter(t *testing.T) {
	b := &Breaker{
		config: Config{FailureThreshold: 3, Timeout: 1 * time.Hour},
		state:  StateClosed,
	}

	b.RecordResult(errors.New("error"))
	b.RecordResult(nil) // success resets
	b.RecordResult(errors.New("error"))
	b.RecordResult(errors.New("error"))

	if b.State() != StateClosed {
		t.Errorf("state = %d, want %d (count should have been reset)", b.State(), StateClosed)
	}
}

func TestRoundTripper_PassesThroughOnSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	rt := NewRoundTripper(http.DefaultTransport, DefaultConfig())
	client := &http.Client{Transport: rt}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if rt.Breaker().State() != StateClosed {
		t.Errorf("breaker state = %d, want %d", rt.Breaker().State(), StateClosed)
	}
}

func TestRoundTripper_HTTPErrorsDoNotTripBreaker(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	rt := NewRoundTripper(http.DefaultTransport, Config{
		FailureThreshold: 3,
		Timeout:          1 * time.Hour,
	})
	client := &http.Client{Transport: rt}

	// HTTP 500 errors should NOT trip the circuit breaker
	for range 5 {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
		if err != nil {
			t.Fatalf("NewRequest error: %v", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("GET failed: %v", err)
		}
		_ = resp.Body.Close()
	}

	if rt.Breaker().State() != StateClosed {
		t.Errorf("breaker state = %d, want %d (HTTP errors should not trip)", rt.Breaker().State(), StateClosed)
	}
}

func TestIsConnectionError(t *testing.T) {
	t.Parallel()

	if isConnectionError(nil) {
		t.Error("isConnectionError(nil) = true, want false")
	}
	if isConnectionError(context.Canceled) {
		t.Error("isConnectionError(context.Canceled) = true, want false")
	}
	if isConnectionError(context.DeadlineExceeded) {
		t.Error("isConnectionError(context.DeadlineExceeded) = true, want false")
	}
	if !isConnectionError(errors.New("connection refused")) {
		t.Error("isConnectionError(conn refused) = false, want true")
	}
}

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	// Kills the ARITHMETIC_BASE mutant on `30 * time.Second` (which would
	// produce a 30ns timeout).
	cfg := DefaultConfig()
	if cfg.FailureThreshold != 5 {
		t.Errorf("FailureThreshold = %d, want 5", cfg.FailureThreshold)
	}
	if cfg.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want %v", cfg.Timeout, 30*time.Second)
	}
}

func TestBreaker_NotHalfOpenAtExactTimeout(t *testing.T) {
	t.Parallel()

	// The circuit transitions to half-open only strictly after Timeout has
	// elapsed; at exactly Timeout it must remain open. Kills the
	// CONDITIONALS_BOUNDARY on `now().Sub(b.lastFailureTime) > b.config.Timeout`.
	base := time.Now()
	b := &Breaker{
		config:          Config{FailureThreshold: 1, Timeout: 10 * time.Second},
		state:           StateOpen,
		lastFailureTime: base.Add(-10 * time.Second),
		now:             func() time.Time { return base },
	}
	if err := b.Allow(); err != ErrCircuitOpen {
		t.Errorf("Allow() = %v, want %v (exactly at timeout boundary)", err, ErrCircuitOpen)
	}
	if b.State() != StateOpen {
		t.Errorf("state = %d, want %d", b.State(), StateOpen)
	}
}

func TestRoundTripper_FailsFastWhenOpen(t *testing.T) {
	b := &Breaker{
		config:          Config{FailureThreshold: 1, Timeout: 1 * time.Hour},
		state:           StateOpen,
		lastFailureTime: time.Now(),
	}
	rt := NewRoundTripper(http.DefaultTransport, DefaultConfig())
	rt.breaker = b

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}

	// With an open circuit the inner transport must never be reached.
	resp, err := rt.RoundTrip(req)
	if resp != nil {
		_ = resp.Body.Close()
	}
	if err != ErrCircuitOpen {
		t.Errorf("RoundTrip() err = %v, want %v", err, ErrCircuitOpen)
	}
	if resp != nil {
		t.Errorf("RoundTrip() resp = %v, want nil", resp)
	}
}
