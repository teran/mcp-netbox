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

// burstForLimit computes the token-bucket burst for a given per-second limit.
// The burst is double the limit so that short request spikes are absorbed
// without tripping the rate limiter while the steady-state rate stays capped.
func burstForLimit(limit int) int {
	return limit * burstMultiplier
}

// burstMultiplier scales a per-second rate limit into a token-bucket burst.
const burstMultiplier = 2

// NewMux builds the MCP HTTP mux (served on LISTEN_ADDR, :8080) with all
// middleware and routes configured. Only the MCP endpoint (/mcp) is served
// here; observability (metrics, pprof, probes) lives on the separate internal
// mux (see NewInternalMux) on INTERNAL_ADDR (:8081) (O01/O04).
// Returns the mux and a stop function for background goroutines (e.g. rate limiter eviction).
func NewMux(
	cfg config.Config,
	metrics *Metrics,
	sharedHTTPClient *http.Client,
	_ *circuitbreaker.Breaker,
	mcpHandler http.Handler,
	upstreamMetrics *infra.UpstreamMetrics,
) (*http.ServeMux, func()) {
	injectClientMW := injectClientMiddleware(cfg.NetBoxURL, sharedHTTPClient, upstreamMetrics)
	rateLimitMW, stopRateLimit := RateLimitMiddleware(RateLimiterConfig{
		GlobalLimit:    rate.Limit(cfg.RateLimitGlobal),
		GlobalBurst:    burstForLimit(cfg.RateLimitGlobal),
		PerClientLimit: rate.Limit(cfg.RateLimitPerClient),
		PerClientBurst: burstForLimit(cfg.RateLimitPerClient),
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
	mux.Handle("/mcp", handler)

	return mux, stopRateLimit
}

// injectClientMiddleware creates a middleware that injects a per-request
// NetBox API client and NetworkService into the request context.
func injectClientMiddleware(
	netboxURL string,
	httpClient *http.Client,
	upstreamMetrics *infra.UpstreamMetrics,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			token, ok := request.Context().Value(tokenContextKey).(*application.Token)
			if !ok || token == nil || token.Value() == "" {
				http.Error(writer, `{"error":"authorization token missing from context"}`, http.StatusUnauthorized)

				return
			}

			netboxClient := infra.NewClientWithMetrics(netboxURL, httpClient, upstreamMetrics)
			svc := application.NewNetworkService(netboxClient, token.Value())

			ctx := context.WithValue(request.Context(), svcContextKey, svc)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}
