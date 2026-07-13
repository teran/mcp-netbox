// Package application provides the business logic / use case layer.
package application

import (
	"encoding/json"
)

// Token is a sensitive value that redacts itself when formatted or serialized.
// Use this instead of a raw string for tokens to prevent accidental leakage
// through logging, debugging, or serialization.
type Token struct {
	value string
}

// NewToken creates a new Token. Returns nil if token is empty.
func NewToken(token string) *Token {
	if token == "" {
		return nil
	}
	return &Token{value: token}
}

// String returns a redacted representation of the token.
func (t *Token) String() string {
	if t == nil {
		return ""
	}
	if len(t.value) <= 8 {
		return "***redacted***"
	}
	return t.value[:4] + "***redacted***" + t.value[len(t.value)-4:]
}

// MarshalJSON implements json.Marshaler to prevent token leakage via JSON.
func (t *Token) MarshalJSON() ([]byte, error) {
	return json.Marshal("***redacted***")
}

// Value returns the actual token string. Use this for the actual API call.
func (t *Token) Value() string {
	if t == nil {
		return ""
	}
	return t.value
}

// GoString implements fmt.GoStringer to prevent token leakage via %#v formatting.
func (t *Token) GoString() string {
	return "application.Token(***redacted***)"
}
