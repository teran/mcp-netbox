package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestNewInternalMux_Metrics(t *testing.T) {
	t.Parallel()

	reg := prometheus.NewRegistry()
	metricsHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux := NewInternalMux(metricsHandler, nil)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/metrics", http.NoBody)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("metrics status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestNewInternalMux_NilMetrics(t *testing.T) {
	t.Parallel()

	// With a nil metrics handler the /metrics route must not be registered.
	mux := NewInternalMux(nil, nil)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/metrics", http.NoBody)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("metrics status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestNewInternalMux_StartupProbe(t *testing.T) {
	t.Parallel()

	mux := NewInternalMux(nil, nil)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/startup", http.NoBody)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("startup status = %d, want %d", rr.Code, http.StatusOK)
	}
	if rr.Body.String() != `{"status":"ok"}` {
		t.Errorf("startup body = %q, want %q", rr.Body.String(), `{"status":"ok"}`)
	}
}

func TestNewInternalMux_Pprof(t *testing.T) {
	t.Parallel()

	mux := NewInternalMux(nil, nil)

	for _, path := range []string{"/debug/pprof/", "/debug/pprof/cmdline"} {
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, path, http.NoBody)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("%s status = %d, want %d", path, rr.Code, http.StatusOK)
		}
	}
}
