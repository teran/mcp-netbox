package domain

import (
	"context"
	"testing"
)

func TestRequestID_RoundTrip(t *testing.T) {
	t.Parallel()

	ctx := WithRequestID(context.Background(), "req-123")
	if got := RequestIDFromContext(ctx); got != "req-123" {
		t.Errorf("RequestIDFromContext = %q, want req-123", got)
	}
}

func TestRequestIDFromContext_DefaultEmpty(t *testing.T) {
	t.Parallel()

	if got := RequestIDFromContext(context.Background()); got != "" {
		t.Errorf("RequestIDFromContext(empty) = %q, want empty string", got)
	}
}

func TestRequestID_WithValueOverrides(t *testing.T) {
	t.Parallel()

	ctx := WithRequestID(context.Background(), "first")
	ctx = WithRequestID(ctx, "second")
	if got := RequestIDFromContext(ctx); got != "second" {
		t.Errorf("RequestIDFromContext = %q, want second", got)
	}
}
