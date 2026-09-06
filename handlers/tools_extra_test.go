package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/config"
	"github.com/teran/mcp-netbox/domain"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/teran/mcp-netbox/infrastructure/circuitbreaker"
)

// newTestMetrics returns a Metrics bound to a fresh registry.
func newTestMetrics() *Metrics {
	return NewMetrics(prometheus.NewRegistry())
}

func TestReadyz_Healthy(t *testing.T) {
	t.Parallel()

	// cb == nil -> readyz reports healthy.
	mux, stop := NewMux(testConfig(), newTestMetrics(), http.DefaultClient, nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer stop()

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("readyz status = %d, want %d", rr.Code, http.StatusOK)
	}
	if !bytes.Contains(rr.Body.Bytes(), []byte("ok")) {
		t.Errorf("readyz body = %q, want to contain 'ok'", rr.Body.String())
	}
}

func TestReadyz_DegradedWhenCircuitOpen(t *testing.T) {
	t.Parallel()

	cb := openCircuitBreaker(t)

	mux, stop := NewMux(testConfig(), newTestMetrics(), http.DefaultClient, cb, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer stop()

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("readyz status = %d, want %d", rr.Code, http.StatusServiceUnavailable)
	}
	if !bytes.Contains(rr.Body.Bytes(), []byte("degraded")) {
		t.Errorf("readyz body = %q, want to contain 'degraded'", rr.Body.String())
	}
}

// openCircuitBreaker returns a circuit breaker in the open state.
func openCircuitBreaker(t *testing.T) *circuitbreaker.Breaker {
	t.Helper()

	failing := failingRoundTripper{}
	rt := circuitbreaker.NewRoundTripper(failing, circuitbreaker.Config{
		FailureThreshold: 1,
		Timeout:          1 * time.Hour,
	})
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	if _, err := rt.RoundTrip(req); err == nil {
		t.Fatal("expected failing RoundTrip to return error")
	}
	if rt.Breaker().State() != circuitbreaker.StateOpen {
		t.Fatal("breaker did not open after failure")
	}
	return rt.Breaker()
}

// failingRoundTripper always fails to satisfy http.RoundTripper.
type failingRoundTripper struct{}

func (failingRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("connection refused")
}

func TestHealthz_Endpoint(t *testing.T) {
	t.Parallel()

	mux, stop := NewMux(testConfig(), newTestMetrics(), http.DefaultClient, nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer stop()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("healthz status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestInjectClientMiddleware_MissingToken(t *testing.T) {
	t.Parallel()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	mw := injectClientMiddleware("http://netbox.example.com", http.DefaultClient)
	rr := httptest.NewRecorder()
	mw(next).ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/mcp", nil))

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestInjectClientMiddleware_WithToken(t *testing.T) {
	t.Parallel()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ServiceFromContext(r.Context()) == nil {
			t.Error("ServiceFromContext() = nil, want non-nil after injection")
		}
		w.WriteHeader(http.StatusOK)
	})

	mw := injectClientMiddleware("http://netbox.example.com", http.DefaultClient)
	ctx := context.WithValue(context.Background(), tokenContextKey, application.NewToken("test-token"))
	rr := httptest.NewRecorder()
	mw(next).ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/mcp", nil).WithContext(ctx))

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestCheckBatchSize_ValidBatch(t *testing.T) {
	t.Parallel()

	// A valid, small batch array must pass checkBatchSize without error.
	body := []byte(`[{"jsonrpc":"2.0","id":1,"method":"ping"},{"jsonrpc":"2.0","id":2,"method":"ping"}]`)
	if err := checkBatchSize(body); err != nil {
		t.Errorf("checkBatchSize(valid batch) = %v, want nil", err)
	}
}

func TestExtractFirstIP_Empty(t *testing.T) {
	t.Parallel()

	if got := extractFirstIP(""); got != "" {
		t.Errorf("extractFirstIP(\"\") = %q, want empty", got)
	}
}

// testConfig returns a minimal Config suitable for NewMux.
func testConfig() config.Config {
	return config.Config{
		NetBoxURL:          "http://netbox.example.com",
		RateLimitGlobal:    100000,
		RateLimitPerClient: 100000,
	}
}

func TestNewListHandler_NoService(t *testing.T) {
	t.Parallel()

	h := NewGetSitesHandler(nil)
	_, _, err := h(context.Background(), nil, SitesInput{})
	if err != errServiceNotAvailable {
		t.Errorf("err = %v, want %v", err, errServiceNotAvailable)
	}
}

func TestNewListHandler_NilResponse(t *testing.T) {
	t.Parallel()

	h := newListHandler[struct{}, domain.Site](listHandlerConfig[struct{}, domain.Site]{
		svc: application.NewNetworkService(&stubRepo{}, "t"),
		listFunc: func(context.Context, map[string]string) (*domain.PaginatedResponse[domain.Site], error) {
			return nil, nil
		},
		buildParams: func(struct{}) map[string]string { return nil },
		errorLabel:  "list sites",
	})

	_, _, err := h(context.Background(), nil, struct{}{})
	if err == nil {
		t.Fatal("err = nil, want error for nil response")
	}
}

func TestNewGetObjectByIDHandler_NoService(t *testing.T) {
	t.Parallel()

	h := NewGetObjectByIDHandler(nil)
	_, _, err := h(context.Background(), nil, GetObjectInput{ObjectType: "site", ID: 1})
	if err != errServiceNotAvailable {
		t.Errorf("err = %v, want %v", err, errServiceNotAvailable)
	}
}

func TestBuildParams_Branches(t *testing.T) {
	t.Parallel()

	svc := application.NewNetworkService(&stubRepo{}, "t")

	t.Run("prefixes family set", func(t *testing.T) {
		four := 4
		// Must not panic and must route the family param into buildParams.
		NewGetPrefixesHandler(svc)(context.Background(), nil, PrefixesInput{Family: &four})
	})

	t.Run("interfaces enabled set", func(t *testing.T) {
		enabled := true
		NewGetInterfacesHandler(svc)(context.Background(), nil, InterfacesInput{Enabled: &enabled})
	})
}
