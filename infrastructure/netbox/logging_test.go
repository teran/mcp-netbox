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

	logger := logrus.New()
	logger.SetOutput(io.Discard)
	logger.SetLevel(logrus.DebugLevel)

	SetLogger(logger)

	if got := getLogger(); got != logger {
		t.Errorf("getLogger() = %p, want %p", got, logger)
	}
}

func TestGetLoggerDefaultDiscards(t *testing.T) {
	t.Cleanup(func() {
		SetLogger(nil)
	})

	SetLogger(nil)

	logger := getLogger()
	if logger == nil {
		t.Fatal("getLogger() returned nil")
	}
	if logger.Out != io.Discard {
		t.Errorf("default logger output = %T, want io.Discard", logger.Out)
	}
}

func TestClient_LogsRequestIDAtDebug(t *testing.T) {
	t.Cleanup(func() {
		SetLogger(nil)
	})

	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetLevel(logrus.DebugLevel)
	SetLogger(logger)

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"count":0,"results":[]}`))
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
