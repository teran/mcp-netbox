package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/domain"
)

func TestSanitizeLog(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"hello\nworld", "helloworld"},
		{"hello\tworld", "hello\tworld"},
		{"hello\x00world", "helloworld"},
		{"hello\x1fworld", "helloworld"},
		{"hello\x7fworld", "helloworld"},
	}
	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			got := SanitizeLog(tc.input)
			if got != tc.expected {
				t.Errorf("SanitizeLog(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	t.Parallel()

	handler := RecoveryMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}
	if body := rec.Body.String(); body != `{"error":"internal server error"}` {
		t.Errorf("body = %q, want %q", body, `{"error":"internal server error"}`)
	}
}

func TestBodyLimitMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("within limit", func(t *testing.T) {
		handler := BodyLimitMiddleware(1024)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", strings.NewReader("small body"))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("exceeds limit", func(t *testing.T) {
		handler := BodyLimitMiddleware(10)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := io.ReadAll(r.Body)
			if err != nil {
				if MaxBytesError(err) {
					http.Error(w, "request too large", http.StatusRequestEntityTooLarge)
					return
				}
			}
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", strings.NewReader("this body is way too long for the limit"))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusRequestEntityTooLarge {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusRequestEntityTooLarge)
		}
	})
}

func TestMaxBytesError(t *testing.T) {
	t.Parallel()

	if MaxBytesError(nil) {
		t.Error("MaxBytesError(nil) = true, want false")
	}

	err := &http.MaxBytesError{}
	if !MaxBytesError(err) {
		t.Error("MaxBytesError(http.MaxBytesError) = false, want true")
	}
}

func TestMCPRequestMethod(t *testing.T) {
	t.Parallel()

	tests := []struct {
		body     string
		expected string
	}{
		{`{"method":"ping"}`, "ping"},
		{`{"method":"tools/call/search"}`, "tools/call/search"},
		{``, "empty_body"},
		{`invalid`, "parse_error"},
		{`{"method":""}`, "no_method"},
		{`{}`, "no_method"},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			got := mcpRequestMethod([]byte(tc.body))
			if got != tc.expected {
				t.Errorf("mcpRequestMethod(%q) = %q, want %q", tc.body, got, tc.expected)
			}
		})
	}
}

func TestLoggingMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("logs request", func(t *testing.T) {
		handler := LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"result":"ok"}`))
		}))

		body := `{"method":"ping"}`
		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/mcp", strings.NewReader(body))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("rejects large batch", func(t *testing.T) {
		handler := LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		items := make([]string, MaxBatchSize+1)
		for i := range items {
			items[i] = `{}`
		}
		body := `[` + strings.Join(items, ",") + `]`
		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/mcp", strings.NewReader(body))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("tracks response size and status", func(t *testing.T) {
		handler := LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"not found"}`))
		}))

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", strings.NewReader(`{"method":"test"}`))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})
}

func TestCheckBatchSize(t *testing.T) {
	t.Parallel()

	if err := checkBatchSize(nil); err != nil {
		t.Errorf("checkBatchSize(nil) = %v, want nil", err)
	}

	if err := checkBatchSize([]byte("{}")); err != nil {
		t.Errorf("checkBatchSize({}) = %v, want nil", err)
	}

	items := make([]string, MaxBatchSize+1)
	for i := range items {
		items[i] = `{}`
	}
	body := []byte(`[` + strings.Join(items, ",") + `]`)
	if err := checkBatchSize(body); err == nil {
		t.Error("checkBatchSize(large batch) = nil, want error")
	}

	// Test invalid JSON body
	if err := checkBatchSize([]byte("[")); err == nil {
		t.Error("checkBatchSize(invalid JSON) = nil, want error")
	}
}

func TestLoggingMiddleware_BodyTooLarge(t *testing.T) {
	t.Parallel()

	handler := BodyLimitMiddleware(10)(
		LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})),
	)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", strings.NewReader("this body is way too long"))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusRequestEntityTooLarge)
	}
}

func TestLoggingResponseWriter(t *testing.T) {
	t.Parallel()

	t.Run("captures status code and body size", func(t *testing.T) {
		rec := httptest.NewRecorder()
		lrw := &loggingResponseWriter{ResponseWriter: rec, statusCode: http.StatusOK}

		lrw.WriteHeader(http.StatusTeapot)
		n, err := lrw.Write([]byte("hello"))
		if err != nil {
			t.Fatalf("Write() returned error: %v", err)
		}

		if lrw.statusCode != http.StatusTeapot {
			t.Errorf("statusCode = %d, want %d", lrw.statusCode, http.StatusTeapot)
		}
		if lrw.bodySize.Load() != 5 {
			t.Errorf("bodySize = %d, want %d", lrw.bodySize.Load(), 5)
		}
		if n != 5 {
			t.Errorf("Write returned %d, want %d", n, 5)
		}
	})
}

func TestTokenMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("valid bearer token", func(t *testing.T) {
		handler := TokenMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Context().Value(tokenContextKey).(*application.Token)
			if token == nil {
				t.Fatal("token is nil")
			}
			if token.Value() != "test-token" {
				t.Errorf("token = %q, want %q", token.Value(), "test-token")
			}
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
		req.Header.Set("Authorization", "Bearer test-token")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("valid token scheme", func(t *testing.T) {
		handler := TokenMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Context().Value(tokenContextKey).(*application.Token)
			if token == nil {
				t.Fatal("token is nil")
			}
			if token.Value() != "tok-value" {
				t.Errorf("token = %q, want %q", token.Value(), "tok-value")
			}
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
		req.Header.Set("Authorization", "Token tok-value")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("missing header", func(t *testing.T) {
		handler := TokenMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called")
		}))

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("empty token", func(t *testing.T) {
		handler := TokenMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called")
		}))

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
		req.Header.Set("Authorization", "Bearer ")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("unsupported scheme", func(t *testing.T) {
		handler := TokenMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called")
		}))

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
		req.Header.Set("Authorization", "Basic dGVzdDp0ZXN0")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("no space after scheme", func(t *testing.T) {
		handler := TokenMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called")
		}))

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
		req.Header.Set("Authorization", "BearerToken")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("token too long", func(t *testing.T) {
		handler := TokenMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called")
		}))

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
		req.Header.Set("Authorization", "Bearer "+strings.Repeat("a", MaxTokenLength+1))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})
}

func TestHostValidationMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("empty host returns 400", func(t *testing.T) {
		handler := HostValidationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called when host is empty")
		}))

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
		req.Host = "" // Ensure empty host
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("non-empty host passes through", func(t *testing.T) {
		handler := HostValidationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Host != "example.com" {
				t.Errorf("Host = %q, want %q", r.Host, "example.com")
			}
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
		req.Host = "example.com"
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	t.Parallel()

	handler := SecurityHeadersMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("X-Content-Type-Options = %q, want %q", rec.Header().Get("X-Content-Type-Options"), "nosniff")
	}
	if rec.Header().Get("X-Frame-Options") != "DENY" {
		t.Errorf("X-Frame-Options = %q, want %q", rec.Header().Get("X-Frame-Options"), "DENY")
	}
	if rec.Header().Get("Referrer-Policy") != "strict-origin-when-cross-origin" {
		t.Errorf("Referrer-Policy = %q, want %q", rec.Header().Get("Referrer-Policy"), "strict-origin-when-cross-origin")
	}
}

func TestLoggingMiddleware_ToolCallMethod(t *testing.T) {
	t.Parallel()

	handler := LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/mcp", strings.NewReader(`{"method":"tools/call/get_sites"}`))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestRequestIDMiddleware_GeneratesID(t *testing.T) {
	t.Parallel()

	var got string
	handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = domain.RequestIDFromContext(r.Context())
	}))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got == "" {
		t.Error("RequestIDFromContext = empty, want a generated id")
	}
}

func TestRequestIDMiddleware_ReusesInboundHeader(t *testing.T) {
	t.Parallel()

	var got string
	handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = domain.RequestIDFromContext(r.Context())
	}))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
	req.Header.Set("X-Request-ID", "inbound-42")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got != "inbound-42" {
		t.Errorf("RequestIDFromContext = %q, want inbound-42", got)
	}
}

func TestWithSession_AddsRequestIDField(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logrus.New()
	l.SetOutput(&buf)
	l.SetLevel(logrus.InfoLevel)

	ctx := domain.WithRequestID(context.Background(), "req-xyz")
	WithSession(ctx, l).Info("hello")

	out := buf.String()
	if !strings.Contains(out, "req-xyz") {
		t.Errorf("log output = %q, want it to contain the request_id", out)
	}
}

func TestWithSession_NoRequestID(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logrus.New()
	l.SetOutput(&buf)
	l.SetLevel(logrus.InfoLevel)

	WithSession(context.Background(), l).Info("hello")

	if strings.Contains(buf.String(), "request_id") {
		t.Errorf("log output = %q, want no request_id field", buf.String())
	}
}

func TestRequestSource_Chain(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
	req.Header.Set("X-Real-IP", "1.1.1.1")
	req.Header.Set("X-Forwarded-For", "2.2.2.2, 3.3.3.3")
	req.RemoteAddr = "4.4.4.4:1234"

	got := requestSource(req)
	if got != "1.1.1.1, 2.2.2.2, 3.3.3.3, 4.4.4.4:1234" {
		t.Errorf("requestSource = %q", got)
	}
}

func TestRequestSource_OnlyPeer(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
	req.RemoteAddr = "10.0.0.1:99"

	if got := requestSource(req); got != "10.0.0.1:99" {
		t.Errorf("requestSource = %q, want peer only", got)
	}
}

func TestRequestArgs_ExtractsParams(t *testing.T) {
	t.Parallel()

	body := []byte(`{"method":"tools/call/get_sites","params":{"q":"dc"}}`)
	got := requestArgs(body)
	if string(got.(json.RawMessage)) != `{"q":"dc"}` {
		t.Errorf("requestArgs = %s, want params object", got)
	}
}

func TestRequestArgs_Empty(t *testing.T) {
	t.Parallel()

	if got := requestArgs(nil); got != nil {
		t.Errorf("requestArgs(nil) = %v, want nil", got)
	}
	if got := requestArgs([]byte(`{"method":"ping"}`)); got != nil {
		t.Errorf("requestArgs(no params) = %v, want nil", got)
	}
}

func TestLoggingMiddleware_InfoLevelWithFields(t *testing.T) {
	t.Cleanup(func() {
		SetLogger(nil)
	})

	var buf bytes.Buffer
	l := logrus.New()
	l.SetOutput(&buf)
	// L8: the per-request access log must be visible at the default info level
	// (no debug gate), so it is tested at InfoLevel.
	l.SetLevel(logrus.InfoLevel)
	SetLogger(l)

	handler := LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result":"ok"}`))
	}))

	body := `{"method":"tools/call/get_sites","params":{"q":"dc"}}`
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/mcp", strings.NewReader(body))
	req.RemoteAddr = "5.5.5.5:80"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	out := buf.String()
	for _, want := range []string{"get_sites", "5.5.5.5:80", "success", "mcp_request"} {
		if !strings.Contains(out, want) {
			t.Errorf("log output missing %q; got: %s", want, out)
		}
	}
	if strings.Contains(out, "Authorization") {
		t.Errorf("log output must not contain Authorization header; got: %s", out)
	}
}
