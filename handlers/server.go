package handlers

import (
	"context"
	"net/http"

	"golang.org/x/time/rate"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/config"
	infra "github.com/teran/mcp-netbox/infrastructure/netbox"
)

const svcContextKey ContextKey = "netbox_service"

// ServiceFromContext extracts the NetworkService from the request context.
// Returns nil if the service is not present in the context.
func ServiceFromContext(r *http.Request) *application.NetworkService {
	if svc, ok := r.Context().Value(svcContextKey).(*application.NetworkService); ok {
		return svc
	}
	return nil
}

// NewMux builds the HTTP mux with all middleware and routes configured.
func NewMux(cfg config.Config, metrics *Metrics, sharedHTTPClient *http.Client, mcpHandler http.Handler) *http.ServeMux {
	injectClientMW := injectClientMiddleware(cfg.NetBoxURL, sharedHTTPClient)
	rateLimitMW, _ := RateLimitMiddleware(RateLimiterConfig{
		GlobalLimit:    rate.Limit(cfg.RateLimitGlobal),
		GlobalBurst:    cfg.RateLimitGlobal * 2,
		PerClientLimit: rate.Limit(cfg.RateLimitPerClient),
		PerClientBurst: cfg.RateLimitPerClient * 2,
	})

	handler := RecoveryMiddleware(
		MetricsMiddleware(metrics)(
			rateLimitMW(
				BodyLimitMiddleware(DefaultMaxRequestBodySize)(
					LoggingMiddleware(
						TokenMiddleware(
							injectClientMW(mcpHandler),
						),
					),
				),
			),
		),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	mux.Handle("/mcp", handler)

	return mux
}

// injectClientMiddleware creates a middleware that injects a per-request
// NetBox API client and NetworkService into the request context.
func injectClientMiddleware(netboxURL string, httpClient *http.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := r.Context().Value(TokenContextKey).(string)
			if !ok || token == "" {
				http.Error(w, `{"error":"authorization token missing from context"}`, http.StatusUnauthorized)
				return
			}

			netboxClient := infra.NewClient(netboxURL, httpClient)
			svc := application.NewNetworkService(netboxClient, token)

			ctx := context.WithValue(r.Context(), svcContextKey, svc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
