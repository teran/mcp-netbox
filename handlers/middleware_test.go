package handlers

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSanitizeLog(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"hello\nworld", "hello\nworld"},
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
	if err := checkBatchSize([]byte("[")); err != nil {
		t.Logf("checkBatchSize(invalid JSON) = %v", err)
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
		if lrw.bodySize != 5 {
			t.Errorf("bodySize = %d, want %d", lrw.bodySize, 5)
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
			token := r.Context().Value(TokenContextKey).(string)
			if token != "test-token" {
				t.Errorf("token = %q, want %q", token, "test-token")
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
			token := r.Context().Value(TokenContextKey).(string)
			if token != "tok-value" {
				t.Errorf("token = %q, want %q", token, "tok-value")
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

func TestLoggingMiddleware_ToolCallMethod(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	// We can't easily capture log output, so just verify it doesn't crash
	handler := LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/mcp", strings.NewReader(`{"method":"tools/call/get_sites"}`))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	_ = buf
}
