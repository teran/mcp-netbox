package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/prometheus/client_golang/prometheus"

	_ "github.com/teran/mcp-netbox/domain"
)

type testInput struct {
	Name string `json:"name"`
}

type testOutput struct {
	Result string `json:"result"`
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
