package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/config"
	"github.com/teran/mcp-netbox/handlers"
)

// ---------------------------------------------------------------------------
// TestServiceFromContext
// ---------------------------------------------------------------------------

func TestServiceFromContext(t *testing.T) {
	t.Parallel()

	t.Run("returns nil when context has no service", func(t *testing.T) {
		svc := handlers.ServiceFromContext(context.Background())
		if svc != nil {
			t.Errorf("ServiceFromContext(background) = %v, want nil", svc)
		}
	})

	t.Run("returns service when context has one", func(t *testing.T) {
		// We can't use svcContextKey directly (unexported), so we route a real
		// request through NewMux and inspect the injected service in the handler.

		cfg := config.Config{
			NetBoxURL:          "http://localhost:1",
			RateLimitGlobal:    100,
			RateLimitPerClient: 10,
		}
		metrics := handlers.NewMetrics(prometheus.NewRegistry())

		var capturedSvc *application.NetworkService
		mcpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedSvc = handlers.ServiceFromContext(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		mux, stop := handlers.NewMux(cfg, metrics, http.DefaultClient, mcpHandler)
		defer stop()

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/mcp", http.NoBody)
		req.Header.Set("Authorization", "Bearer test-token")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if capturedSvc == nil {
			t.Error("ServiceFromContext returned nil, want non-nil service")
		}
	})
}

// ---------------------------------------------------------------------------
// TestInjectClientMiddleware (exercised via NewMux)
// ---------------------------------------------------------------------------

func TestInjectClientMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("injects service when valid token is provided", func(t *testing.T) {
		cfg := config.Config{
			NetBoxURL:          "http://localhost:1",
			RateLimitGlobal:    100,
			RateLimitPerClient: 10,
		}
		metrics := handlers.NewMetrics(prometheus.NewRegistry())

		var svc *application.NetworkService
		mcpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			svc = handlers.ServiceFromContext(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		mux, stop := handlers.NewMux(cfg, metrics, http.DefaultClient, mcpHandler)
		defer stop()

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/mcp", http.NoBody)
		req.Header.Set("Authorization", "Bearer some-valid-token")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if svc == nil {
			t.Error("service was not injected into context")
		}
	})

	t.Run("returns 401 when no authorization header is present", func(t *testing.T) {
		cfg := config.Config{
			NetBoxURL:          "http://localhost:1",
			RateLimitGlobal:    100,
			RateLimitPerClient: 10,
		}
		metrics := handlers.NewMetrics(prometheus.NewRegistry())

		unreachable := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("downstream handler should not be reached when token is missing")
		})

		mux, stop := handlers.NewMux(cfg, metrics, http.DefaultClient, unreachable)
		defer stop()

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/mcp", http.NoBody)
		// No Authorization header set
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("returns 401 on empty token", func(t *testing.T) {
		cfg := config.Config{
			NetBoxURL:          "http://localhost:1",
			RateLimitGlobal:    100,
			RateLimitPerClient: 10,
		}
		metrics := handlers.NewMetrics(prometheus.NewRegistry())

		unreachable := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("downstream handler should not be reached")
		})

		mux, stop := handlers.NewMux(cfg, metrics, http.DefaultClient, unreachable)
		defer stop()

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/mcp", http.NoBody)
		req.Header.Set("Authorization", "Bearer ")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})
}

// ---------------------------------------------------------------------------
// TestNewMux
// ---------------------------------------------------------------------------

func TestNewMux(t *testing.T) {
	t.Parallel()

	t.Run("/healthz returns 200 OK", func(t *testing.T) {
		cfg := config.Config{
			NetBoxURL:          "http://localhost:1",
			RateLimitGlobal:    100,
			RateLimitPerClient: 10,
		}
		metrics := handlers.NewMetrics(prometheus.NewRegistry())
		mcpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("MCP handler should not be called for /healthz")
		})

		mux, stop := handlers.NewMux(cfg, metrics, http.DefaultClient, mcpHandler)
		defer stop()

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/healthz", http.NoBody)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if rec.Body.String() != `{"status":"ok"}` {
			t.Errorf("body = %q, want %q", rec.Body.String(), `{"status":"ok"}`)
		}
		if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want %q", ct, "application/json")
		}
	})

	t.Run("unknown route returns 404", func(t *testing.T) {
		cfg := config.Config{
			NetBoxURL:          "http://localhost:1",
			RateLimitGlobal:    100,
			RateLimitPerClient: 10,
		}
		metrics := handlers.NewMetrics(prometheus.NewRegistry())
		mcpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("MCP handler should not be called for unknown routes")
		})

		mux, stop := handlers.NewMux(cfg, metrics, http.DefaultClient, mcpHandler)
		defer stop()

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/nonexistent", http.NoBody)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})

	t.Run("stop function stops rate limiter goroutine", func(t *testing.T) {
		cfg := config.Config{
			NetBoxURL:          "http://localhost:1",
			RateLimitGlobal:    100,
			RateLimitPerClient: 10,
		}
		metrics := handlers.NewMetrics(prometheus.NewRegistry())
		mcpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		mux, stop := handlers.NewMux(cfg, metrics, http.DefaultClient, mcpHandler)
		// Stopping should not panic and should clean up the rate limiter goroutine.
		stop()

		// The mux should still work after stopping the rate limiter.
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/healthz", http.NoBody)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})
}
