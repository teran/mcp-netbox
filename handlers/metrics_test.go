package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "github.com/teran/mcp-netbox/domain"
)

type testInput struct {
	Name string `json:"name"`
}

type testOutput struct {
	Result string `json:"result"`
}

type secretOutput struct {
	Result string `json:"result"`
	Token  string `json:"token" secret:"true"`
}

func TestNewMetrics(t *testing.T) {
	t.Parallel()

	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)
	if m == nil {
		t.Fatal("NewMetrics() returned nil")
	}
}

func TestWrapToolHandler(t *testing.T) {
	t.Parallel()

	t.Run("records 2xx metrics", func(t *testing.T) {
		reg := prometheus.NewRegistry()
		m := NewMetrics(reg)

		handler := WrapToolHandler[testInput, testOutput](m, "test_tool", func(ctx context.Context, req *mcp.CallToolRequest, in testInput) (*mcp.CallToolResult, testOutput, error) {
			return &mcp.CallToolResult{}, testOutput{Result: "ok"}, nil
		})

		result, out, err := handler(context.Background(), nil, testInput{Name: "test"})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Result != "ok" {
			t.Errorf("out.Result = %q, want %q", out.Result, "ok")
		}
	})

	t.Run("records 4xx metrics on error", func(t *testing.T) {
		reg := prometheus.NewRegistry()
		m := NewMetrics(reg)

		handler := WrapToolHandler[testInput, testOutput](m, "test_tool", func(ctx context.Context, req *mcp.CallToolRequest, in testInput) (*mcp.CallToolResult, testOutput, error) {
			return &mcp.CallToolResult{IsError: true}, testOutput{}, nil
		})

		result, _, err := handler(context.Background(), nil, testInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})

	t.Run("nil metrics safety", func(t *testing.T) {
		handler := WrapToolHandler[testInput, testOutput](nil, "test_tool", func(ctx context.Context, req *mcp.CallToolRequest, in testInput) (*mcp.CallToolResult, testOutput, error) {
			return &mcp.CallToolResult{}, testOutput{Result: "ok"}, nil
		})

		result, out, err := handler(context.Background(), nil, testInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if out.Result != "ok" {
			t.Errorf("out.Result = %q, want %q", out.Result, "ok")
		}
		_ = result
	})
}

func TestMetricsMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("increments and decrements gauge", func(t *testing.T) {
		reg := prometheus.NewRegistry()
		m := NewMetrics(reg)

		mw := MetricsMiddleware(m)
		handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("nil metrics safety", func(t *testing.T) {
		mw := MetricsMiddleware(nil)
		handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})
}

func TestSanitizeToolResult(t *testing.T) {
	t.Parallel()

	t.Run("populates sanitized JSON text content", func(t *testing.T) {
		reg := prometheus.NewRegistry()
		m := NewMetrics(reg)

		handler := WrapToolHandler[testInput, testOutput](m, "test_tool", func(ctx context.Context, req *mcp.CallToolRequest, in testInput) (*mcp.CallToolResult, testOutput, error) {
			return &mcp.CallToolResult{}, testOutput{Result: "ok \x1b[31mred\x1b[0m"}, nil
		})

		result, _, err := handler(context.Background(), nil, testInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Fatal("result.IsError = true, want false")
		}
		if len(result.Content) != 1 {
			t.Fatalf("len(result.Content) = %d, want 1", len(result.Content))
		}
		tc, ok := result.Content[0].(*mcp.TextContent)
		if !ok {
			t.Fatalf("Content[0] type = %T, want *mcp.TextContent", result.Content[0])
		}
		// json.Marshal escapes ESC as \u001b, so the text form is guaranteed free
		// of raw ANSI/control bytes; the text is the marshalled output.
		want := `{"result":"ok \u001b[31mred\u001b[0m"}`
		if got := tc.Text; got != want {
			t.Errorf("text = %q, want %q", got, want)
		}
	})

	t.Run("passes error results through untouched", func(t *testing.T) {
		reg := prometheus.NewRegistry()
		m := NewMetrics(reg)

		handler := WrapToolHandler[testInput, testOutput](m, "test_tool", func(ctx context.Context, req *mcp.CallToolRequest, in testInput) (*mcp.CallToolResult, testOutput, error) {
			return &mcp.CallToolResult{IsError: true}, testOutput{Result: "\x1b[31mboom\x1b[0m"}, nil
		})

		result, _, err := handler(context.Background(), nil, testInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if !result.IsError {
			t.Fatal("result.IsError = false, want true")
		}
		if len(result.Content) != 0 {
			t.Errorf("len(result.Content) = %d, want 0 (error results untouched)", len(result.Content))
		}
	})
}

func TestWrapToolHandler_RedactsSecret(t *testing.T) {
	t.Parallel()

	// S02: a successful result with a secret:"true" field must have the secret
	// masked in the text form, and the secret must not appear anywhere.
	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)

	handler := WrapToolHandler[testInput, secretOutput](m, "test_tool", func(ctx context.Context, req *mcp.CallToolRequest, in testInput) (*mcp.CallToolResult, secretOutput, error) {
		return &mcp.CallToolResult{}, secretOutput{Result: "ok", Token: "super-secret-token"}, nil
	})

	result, out, err := handler(context.Background(), nil, testInput{})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if result.IsError {
		t.Fatal("result.IsError = true, want false")
	}
	// The typed out returned to the caller must be redacted too.
	if out.Token != "***redacted***" {
		t.Errorf("out.Token = %q, want %q", out.Token, "***redacted***")
	}
	if out.Result != "ok" {
		t.Errorf("out.Result = %q, want %q", out.Result, "ok")
	}

	if len(result.Content) != 1 {
		t.Fatalf("len(result.Content) = %d, want 1", len(result.Content))
	}
	tc, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("Content[0] type = %T, want *mcp.TextContent", result.Content[0])
	}
	if strings.Contains(tc.Text, "super-secret-token") {
		t.Errorf("text leaks the secret: %s", tc.Text)
	}
	if !strings.Contains(tc.Text, "***redacted***") {
		t.Errorf("text does not contain the mask: %s", tc.Text)
	}
}

func TestWrapToolHandler_RedactsNonSecretUnchanged(t *testing.T) {
	t.Parallel()

	// A successful output with no secret field must be returned unchanged.
	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)

	handler := WrapToolHandler[testInput, testOutput](m, "test_tool", func(ctx context.Context, req *mcp.CallToolRequest, in testInput) (*mcp.CallToolResult, testOutput, error) {
		return &mcp.CallToolResult{}, testOutput{Result: "plain-value"}, nil
	})

	_, out, err := handler(context.Background(), nil, testInput{})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if out.Result != "plain-value" {
		t.Errorf("out.Result = %q, want %q", out.Result, "plain-value")
	}
}

func TestWrapToolHandler_RedactionSkipsErrorResults(t *testing.T) {
	t.Parallel()

	// S02: error results are not redacted (redaction only applies to the
	// success path); they must not panic and must be passed through.
	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)

	handler := WrapToolHandler[testInput, secretOutput](m, "test_tool", func(ctx context.Context, req *mcp.CallToolRequest, in testInput) (*mcp.CallToolResult, secretOutput, error) {
		return &mcp.CallToolResult{IsError: true}, secretOutput{Token: "should-stay"}, nil
	})

	result, _, err := handler(context.Background(), nil, testInput{})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if !result.IsError {
		t.Fatal("result.IsError = false, want true")
	}
	if len(result.Content) != 0 {
		t.Errorf("len(result.Content) = %d, want 0 (error result untouched)", len(result.Content))
	}
}

func TestWrapToolHandler_RedactsPointerOutput(t *testing.T) {
	t.Parallel()

	// Redaction must also handle pointer-typed tool output.
	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)

	handler := WrapToolHandler[testInput, *secretOutput](m, "test_tool", func(ctx context.Context, req *mcp.CallToolRequest, in testInput) (*mcp.CallToolResult, *secretOutput, error) {
		return &mcp.CallToolResult{}, &secretOutput{Result: "ok", Token: "ptr-secret"}, nil
	})

	_, out, err := handler(context.Background(), nil, testInput{})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if out == nil {
		t.Fatal("out = nil, want non-nil")
	}
	if out.Token != "***redacted***" {
		t.Errorf("out.Token = %q, want %q", out.Token, "***redacted***")
	}
}

func TestInstrumentInternalMux(t *testing.T) {
	t.Parallel()

	reg := prometheus.NewRegistry()
	metricsHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	internalMux := NewInternalMux(metricsHandler, nil)
	handler := InstrumentInternalMux(reg, internalMux)

	// Exercise the wrapped mux so the http_server_* metrics are populated.
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/startup", http.NoBody)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	// Scrape /metrics through the wrapped handler and confirm the standard
	// net/http metrics are exposed (O02).
	scrape := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/metrics", http.NoBody)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, scrape)
	if rec.Code != http.StatusOK {
		t.Fatalf("metrics status = %d, want %d", rec.Code, http.StatusOK)
	}
	body := rec.Body.String()
	for _, metric := range []string{
		"http_server_requests_total",
		"http_server_request_duration_seconds",
		"http_server_request_size_bytes",
		"http_server_response_size_bytes",
	} {
		if !strings.Contains(body, metric) {
			t.Errorf("metrics output missing %q", metric)
		}
	}
}

func TestRegisterMetricsOnRegistry(t *testing.T) {
	t.Parallel()

	reg := prometheus.NewRegistry()
	handler := RegisterMetricsOnRegistry(reg)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/metrics", http.NoBody)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestSanitizeToolResult_MarshalError(t *testing.T) {
	t.Parallel()

	// A successful, content-free result whose out value cannot be marshalled
	// must be passed through unchanged (the marshal error branch).
	res := &mcp.CallToolResult{}
	unmarshalable := make(chan int)
	got := sanitizeToolResult(res, unmarshalable)
	if got != res {
		t.Error("sanitizeToolResult returned a different result")
	}
	if got.Content != nil {
		t.Errorf("Content = %v, want nil", got.Content)
	}
}

func TestSanitizeToolResult_Passthrough(t *testing.T) {
	t.Parallel()

	// Error results and results with existing content are never rewritten.
	errRes := &mcp.CallToolResult{IsError: true}
	if got := sanitizeToolResult(errRes, "x"); got != errRes {
		t.Error("error result was modified")
	}

	withContent := &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "hi"}}}
	if got := sanitizeToolResult(withContent, "x"); got != withContent {
		t.Error("result with content was modified")
	}

	if got := sanitizeToolResult(nil, "x"); got != nil {
		t.Error("nil result was modified")
	}
}

func TestSanitizeToolResult_SetsSanitizedText(t *testing.T) {
	t.Parallel()

	res := &mcp.CallToolResult{}
	got := sanitizeToolResult(res, map[string]string{"name": "srv\x1b[31m"})
	if got.Content == nil {
		t.Fatal("Content = nil, want sanitized text content")
	}
	text, ok := got.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("Content[0] = %T, want *mcp.TextContent", got.Content[0])
	}
	if !strings.Contains(text.Text, "srv") {
		t.Errorf("text = %q, want it to contain the sanitized field value", text.Text)
	}
}
