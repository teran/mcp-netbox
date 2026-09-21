package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

		mux, stop := handlers.NewMux(cfg, metrics, http.DefaultClient, nil, mcpHandler, nil)
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

		mux, stop := handlers.NewMux(cfg, metrics, http.DefaultClient, nil, mcpHandler, nil)
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

		mux, stop := handlers.NewMux(cfg, metrics, http.DefaultClient, nil, unreachable, nil)
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

		mux, stop := handlers.NewMux(cfg, metrics, http.DefaultClient, nil, unreachable, nil)
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

	t.Run("/healthz returns 200 OK on the internal mux", func(t *testing.T) {
		// The liveness probe lives on the internal observability mux
		// (NewInternalMux), not on the MCP mux.
		mux := handlers.NewInternalMux(nil, nil)

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

	t.Run("MCP mux does not expose probes", func(t *testing.T) {
		// O01/O04: probes must not be reachable on the MCP mux (:8080).
		cfg := config.Config{
			NetBoxURL:          "http://localhost:1",
			RateLimitGlobal:    100,
			RateLimitPerClient: 10,
		}
		metrics := handlers.NewMetrics(prometheus.NewRegistry())
		mcpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		mux, stop := handlers.NewMux(cfg, metrics, http.DefaultClient, nil, mcpHandler, nil)
		defer stop()

		for _, path := range []string{"/healthz", "/readyz", "/startup", "/metrics", "/debug/pprof/"} {
			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, path, http.NoBody)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != http.StatusNotFound {
				t.Errorf("MCP mux %s = %d, want 404", path, rec.Code)
			}
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

		mux, stop := handlers.NewMux(cfg, metrics, http.DefaultClient, nil, mcpHandler, nil)
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

		mux, stop := handlers.NewMux(cfg, metrics, http.DefaultClient, nil, mcpHandler, nil)

		// Signal that the goroutine should have been stopped.
		done := make(chan struct{})
		go func() {
			// We can't directly observe the evictExpired goroutine, but we can
			// verify that stop() does not block indefinitely and the mux still
			// serves requests afterwards.
			stop()
			close(done)
		}()

		select {
		case <-done:
			// stop() returned — goroutine was cleaned up.
		case <-time.After(5 * time.Second):
			t.Fatal("stop() did not return within 5 seconds — goroutine may leak")
		}

		// The mux should still work after stopping the rate limiter. With no
		// Authorization header the /mcp route returns 401, which proves the
		// mux is still serving after stop().
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/mcp", http.NoBody)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})
}
