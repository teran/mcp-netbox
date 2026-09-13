package domain

import "fmt"

// ValidationError represents a NetBox 4xx validation failure (most commonly
// HTTP 400 Bad Request). It carries the raw response body so the handler can
// surface the NetBox field-level errors to the model through structured
// content.
//
// The Body is intentionally NOT included in Error(): the error string is what
// flows into logs and generic error messages, and dumping an unbounded body
// there could leak sensitive detail or bloat log lines. The body is instead
// delivered to the caller via structured content (see handlers).
type ValidationError struct {
	StatusCode int
	Body       []byte
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("netbox validation error: status %d", e.StatusCode)
}
