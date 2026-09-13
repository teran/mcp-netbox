package handlers

import (
	"context"
	"net/http"

	"golang.org/x/time/rate"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/config"
	"github.com/teran/mcp-netbox/infrastructure/circuitbreaker"
	infra "github.com/teran/mcp-netbox/infrastructure/netbox"
)

const svcContextKey contextKey = "netbox_service"

// ServiceFromContext extracts the NetworkService from the request context.
// Returns nil if the service is not present in the context.
func ServiceFromContext(ctx context.Context) *application.NetworkService {
	if svc, ok := ctx.Value(svcContextKey).(*application.NetworkService); ok {
		return svc
	}
	return nil
}

// NewMux builds the HTTP mux with all middleware and routes configured.
// Returns the mux and a stop function for background goroutines (e.g. rate limiter eviction).
func NewMux(cfg config.Config, metrics *Metrics, sharedHTTPClient *http.Client, cb *circuitbreaker.Breaker, mcpHandler http.Handler) (*http.ServeMux, func()) {
	injectClientMW := injectClientMiddleware(cfg.NetBoxURL, sharedHTTPClient)
	rateLimitMW, stopRateLimit := RateLimitMiddleware(RateLimiterConfig{
		GlobalLimit:    rate.Limit(cfg.RateLimitGlobal),
		GlobalBurst:    cfg.RateLimitGlobal * 2,
		PerClientLimit: rate.Limit(cfg.RateLimitPerClient),
		PerClientBurst: cfg.RateLimitPerClient * 2,
	}, cfg.TrustedProxy)

	handler := RequestIDMiddleware(
		RecoveryMiddleware(
			SecurityHeadersMiddleware(
				HostValidationMiddleware(
					rateLimitMW(
						MetricsMiddleware(metrics)(
							BodyLimitMiddleware(DefaultMaxRequestBodySize)(
								LoggingMiddleware(
									TokenMiddleware(
										injectClientMW(mcpHandler),
									),
								),
							),
						),
					),
				),
			),
		),
	)

	mux := http.NewServeMux()
	mux.Handle("GET /healthz", RecoveryMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})))
	mux.Handle("GET /readyz", RecoveryMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		status := http.StatusOK
		body := `{"status":"ok"}`
		if cb != nil && cb.State() == circuitbreaker.StateOpen {
			status = http.StatusServiceUnavailable
			body = `{"status":"degraded","circuit_breaker":"open"}`
		}

		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	})))
	mux.Handle("/mcp", handler)

	return mux, stopRateLimit
}

// injectClientMiddleware creates a middleware that injects a per-request
// NetBox API client and NetworkService into the request context.
func injectClientMiddleware(netboxURL string, httpClient *http.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t, ok := r.Context().Value(tokenContextKey).(*application.Token)
			if !ok || t == nil || t.Value() == "" {
				http.Error(w, `{"error":"authorization token missing from context"}`, http.StatusUnauthorized)
				return
			}

			netboxClient := infra.NewClient(netboxURL, httpClient)
			svc := application.NewNetworkService(netboxClient, t.Value())

			ctx := context.WithValue(r.Context(), svcContextKey, svc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
