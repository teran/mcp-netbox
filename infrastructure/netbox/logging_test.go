package netbox

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/teran/mcp-netbox/domain"
)

func TestSetLoggerAndGetLogger(t *testing.T) {
	t.Cleanup(func() {
		SetLogger(nil)
	})

	l := logrus.New()
	l.SetOutput(io.Discard)
	l.SetLevel(logrus.DebugLevel)

	SetLogger(l)

	if got := getLogger(); got != l {
		t.Errorf("getLogger() = %p, want %p", got, l)
	}
}

func TestGetLoggerDefaultDiscards(t *testing.T) {
	t.Cleanup(func() {
		SetLogger(nil)
	})

	SetLogger(nil)

	l := getLogger()
	if l == nil {
		t.Fatal("getLogger() returned nil")
	}
	if l.Out != io.Discard {
		t.Errorf("default logger output = %T, want io.Discard", l.Out)
	}
}

func TestClient_LogsRequestIDAtDebug(t *testing.T) {
	t.Cleanup(func() {
		SetLogger(nil)
	})

	var buf bytes.Buffer
	l := logrus.New()
	l.SetOutput(&buf)
	l.SetLevel(logrus.DebugLevel)
	SetLogger(l)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"count":0,"results":[]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	ctx := domain.WithRequestID(context.Background(), "corr-debug-99")
	if _, err := client.ListSites(ctx, "token", nil); err != nil {
		t.Fatalf("ListSites() returned error: %v", err)
	}

	// Kills the CONDITIONALS_NEGATION on `requestID != ""` in the outbound
	// debug log: when a request_id is present in the context it must be
	// attached to the emitted log line.
	out := buf.String()
	if !strings.Contains(out, "corr-debug-99") {
		t.Errorf("debug log missing request_id: %q", out)
	}
}
