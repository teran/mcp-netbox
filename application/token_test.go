package application

import (
	"encoding/json"
	"testing"
)

func TestNewToken(t *testing.T) {
	t.Parallel()

	t.Run("creates token for non-empty string", func(t *testing.T) {
		tok := NewToken("my-secret-token")
		if tok == nil {
			t.Fatal("NewToken() = nil, want non-nil")
		}
	})

	t.Run("returns nil for empty string", func(t *testing.T) {
		tok := NewToken("")
		if tok != nil {
			t.Fatal("NewToken() = non-nil, want nil")
		}
	})
}

func TestTokenValue(t *testing.T) {
	t.Parallel()

	t.Run("returns original value", func(t *testing.T) {
		tok := NewToken("supersecret123")
		if tok.Value() != "supersecret123" {
			t.Errorf("Value() = %q, want %q", tok.Value(), "supersecret123")
		}
	})

	t.Run("returns empty string for nil token", func(t *testing.T) {
		var tok *Token
		if tok.Value() != "" {
			t.Errorf("Value() = %q, want empty", tok.Value())
		}
	})
}

func TestTokenString(t *testing.T) {
	t.Parallel()

	t.Run("redacts long token", func(t *testing.T) {
		tok := NewToken("abcdefgh12345678")
		s := tok.String()
		if s == "abcdefgh12345678" {
			t.Error("String() returned raw token, expected redacted")
		}
		if len(s) < 10 || len(s) > 30 {
			t.Errorf("String() length = %d, want reasonable redacted length", len(s))
		}
	})

	t.Run("short token is redacted", func(t *testing.T) {
		tok := NewToken("short")
		s := tok.String()
		if s == "short" {
			t.Error("String() returned raw token for short token")
		}
	})

	t.Run("nil token returns empty", func(t *testing.T) {
		var tok *Token
		if tok.String() != "" {
			t.Errorf("String() = %q, want empty", tok.String())
		}
	})
}

func TestTokenGoString(t *testing.T) {
	t.Parallel()

	tok := NewToken("secret123")
	// Test that GoString does not leak the raw token.
	// Using fmt.Sprint with %#v would be the direct test, but that triggers
	// vet checks, so we test the GoString method directly instead.
	if got := tok.GoString(); got != "application.Token(***redacted***)" {
		t.Errorf("GoString() = %q, want %q", got, "application.Token(***redacted***)")
	}
}

func TestTokenMarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("marshal redacts token", func(t *testing.T) {
		tok := NewToken("my-secret-key")
		data, err := json.Marshal(tok)
		if err != nil {
			t.Fatalf("MarshalJSON() returned error: %v", err)
		}
		if string(data) != `"***redacted***"` {
			t.Errorf("MarshalJSON() = %s, want %s", string(data), `"***redacted***"`)
		}
	})

	t.Run("nil token marshals as null", func(t *testing.T) {
		var tok *Token
		data, err := json.Marshal(tok)
		if err != nil {
			t.Fatalf("MarshalJSON() for nil returned error: %v", err)
		}
		if string(data) != "null" {
			t.Errorf("MarshalJSON() for nil = %s, want %s", string(data), "null")
		}
	})
}
