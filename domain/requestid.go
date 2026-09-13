package domain

import "context"

// requestIDContextKey is the context key carrying the per-request correlation
// identifier. Using an unexported type prevents collisions with other packages.
type requestIDContextKey struct{}

// WithRequestID returns a copy of ctx carrying the given request_id. It is the
// domain-level helper used across the request lifecycle (L9/G11) so components
// can read the ID without threading it through every signature.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey{}, requestID)
}

// RequestIDFromContext returns the request_id stored in ctx, or the empty
// string when none is present.
func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDContextKey{}).(string); ok {
		return id
	}
	return ""
}
