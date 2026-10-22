package netbox

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestUserAgent_Default verifies that, before SetUserAgent is called, the
// outbound User-Agent header is a non-empty identifying value (N44): the app
// name with the repository reference, and never the raw go-resty default.
func TestUserAgent_Default(t *testing.T) {
	t.Cleanup(restoreUserAgent)

	const want = "mcp-netbox (https://github.com/teran/mcp-netbox)"

	if got := userAgent(); got != want {
		t.Errorf("userAgent() = %q, want %q", got, want)
	}

	got := captureOutboundUserAgent(t)
	if got != want {
		t.Errorf("outbound User-Agent = %q, want %q", got, want)
	}
}

// TestUserAgent_WithVersion verifies the full M10 format when a version is set.
func TestUserAgent_WithVersion(t *testing.T) {
	t.Cleanup(restoreUserAgent)

	SetUserAgent("1.2.3")

	const want = "mcp-netbox/1.2.3 (https://github.com/teran/mcp-netbox)"

	if got := userAgent(); got != want {
		t.Errorf("userAgent() = %q, want %q", got, want)
	}

	if got := captureOutboundUserAgent(t); got != want {
		t.Errorf("outbound User-Agent = %q, want %q", got, want)
	}
}

// TestUserAgent_EmptyVersion verifies the fallback: an empty version still
// yields a non-empty identifying User-Agent (N44), never an empty or malformed
// header.
func TestUserAgent_EmptyVersion(t *testing.T) {
	t.Cleanup(restoreUserAgent)

	SetUserAgent("")

	const want = "mcp-netbox (https://github.com/teran/mcp-netbox)"

	if got := userAgent(); got != want {
		t.Errorf("userAgent() = %q, want %q", got, want)
	}

	if got := captureOutboundUserAgent(t); got != want {
		t.Errorf("outbound User-Agent = %q, want %q", got, want)
	}
}

// TestUserAgent_DoesNotLeakToken verifies the User-Agent header never contains
// the NetBox token; the Authorization header is the only place the token goes.
func TestUserAgent_DoesNotLeakToken(t *testing.T) {
	t.Cleanup(restoreUserAgent)

	const token = "super-secret-token-abc123"

	got := captureOutboundUserAgentWithToken(t, token)
	if got == "" {
		t.Fatal("outbound User-Agent is empty")
	}
	if strings.Contains(got, token) {
		t.Errorf("User-Agent leaked the token: %q", got)
	}
}

// captureOutboundUserAgent performs a real outbound call with the default token
// and returns the User-Agent header the server received.
func captureOutboundUserAgent(t *testing.T) string {
	t.Helper()

	return captureOutboundUserAgentWithToken(t, "token")
}

// captureOutboundUserAgentWithToken performs a real outbound call with the given
// token and returns the User-Agent header the server received.
func captureOutboundUserAgentWithToken(t *testing.T, token string) string {
	t.Helper()

	var got string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("User-Agent")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"count":0,"results":[]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	if _, err := client.ListSites(context.Background(), token, nil); err != nil {
		t.Fatalf("ListSites() returned error: %v", err)
	}

	return got
}

// restoreUserAgent resets the package User-Agent to its unset state so other
// tests are unaffected.
func restoreUserAgent() {
	SetUserAgent("")
}
