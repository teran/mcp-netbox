package netbox

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

// TestRestyLoggerLevelRouting verifies that restyLogger forwards Errorf/Warnf/
// Debugf to logrus at the matching level, so each record surfaces under the
// corresponding LOG_LEVEL.
func TestRestyLoggerLevelRouting(t *testing.T) {
	cases := []struct {
		name      string
		call      func()
		wantLevel logrus.Level
	}{
		{"error", func() { restyLogger{}.Errorf("resty err %d", 1) }, logrus.ErrorLevel},
		{"warn", func() { restyLogger{}.Warnf("resty warn %d", 2) }, logrus.WarnLevel},
		{"debug", func() { restyLogger{}.Debugf("resty dbg %d", 3) }, logrus.DebugLevel},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Cleanup(func() {
				SetLogger(nil)
			})

			var buf bytes.Buffer
			logger := logrus.New()
			logger.SetOutput(&buf)
			logger.SetLevel(testCase.wantLevel)
			SetLogger(logger)

			testCase.call()

			out := buf.String()
			if out == "" {
				t.Fatalf("expected a log record at level %v, got empty output", testCase.wantLevel)
			}
			if !strings.Contains(out, testCase.wantLevel.String()) {
				t.Errorf("output %q missing level marker %q", out, testCase.wantLevel.String())
			}
		})
	}
}

// TestRestyLoggerSuppressedAtInfo verifies that resty's debug records are gated
// by the logrus level: at info level a Debugf record must be dropped while an
// Errorf record is still written.
func TestRestyLoggerSuppressedAtInfo(t *testing.T) {
	t.Cleanup(func() {
		SetLogger(nil)
	})

	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetLevel(logrus.InfoLevel)
	SetLogger(logger)

	restyLogger{}.Debugf("resty dbg %d", 1)
	if got := buf.String(); got != "" {
		t.Errorf("Debugf emitted at info level: %q", got)
	}

	restyLogger{}.Errorf("resty err %d", 2)
	if !strings.Contains(buf.String(), "resty err 2") {
		t.Errorf("Errorf not emitted at info level: %q", buf.String())
	}
}

// TestNewClientWiresRestyLogger verifies that the resty client created by
// NewClient carries the restyLogger adapter and has debug mode enabled, so
// request/response detail records are produced and routed to logrus.
func TestNewClientWiresRestyLogger(t *testing.T) {
	client := NewClient("http://example.com", http.DefaultClient)
	if !client.client.IsDebug() {
		t.Error("expected resty debug mode to be enabled")
	}
	if _, ok := client.client.Logger().(restyLogger); !ok {
		t.Errorf("resty logger = %T, want restyLogger", client.client.Logger())
	}
}

// TestClient_RestyDebugDoesNotLeakToken performs a real outbound call with
// debug logging enabled and asserts that the resty request/response dump is
// emitted at debug level while the NetBox token and the Authorization header
// value never appear (resty redacts the header before formatting the dump).
func TestClient_RestyDebugDoesNotLeakToken(t *testing.T) {
	t.Cleanup(func() {
		SetLogger(nil)
	})

	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetLevel(logrus.DebugLevel)
	SetLogger(logger)

	const token = "super-secret-token-abc123"

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"count":0,"results":[]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	if _, err := client.ListSites(context.Background(), token, nil); err != nil {
		t.Fatalf("ListSites() returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "~~~ REQUEST ~~~") {
		t.Errorf("expected resty request dump in debug output: %q", out)
	}
	if !strings.Contains(out, "REDACTED") {
		t.Errorf("expected Authorization header to be redacted: %q", out)
	}
	if strings.Contains(out, token) {
		t.Errorf("debug output leaked the token: %q", out)
	}
}
