package handlers

import (
	"net/http"
	"net/http/pprof"

	"github.com/teran/mcp-netbox/infrastructure/circuitbreaker"
)

// healthzHandler is the liveness probe. It returns {"status":"ok"} with HTTP
// 200 to indicate the process is alive.
func healthzHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
}

// readyzHandler is the readiness probe. It reports degraded (HTTP 503) when
// the NetBox circuit breaker is open.
func readyzHandler(cb *circuitbreaker.Breaker) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		status := http.StatusOK
		body := `{"status":"ok"}`
		if cb != nil && cb.State() == circuitbreaker.StateOpen {
			status = http.StatusServiceUnavailable
			body = `{"status":"degraded","circuit_breaker":"open"}`
		}

		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	})
}

// startupHandler is the startup probe. It returns {"status":"ok"} with HTTP
// 200 once the server is up.
func startupHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
}

// NewInternalMux builds the internal observability mux served on INTERNAL_ADDR
// (:8081 by default). It is fully separate from the MCP mux (LISTEN_ADDR,
// :8080) and hosts the Prometheus /metrics endpoint, the
// healthz/readyz/startup probes, and the net/http/pprof endpoints
// (/debug/pprof/*) (O01/O04). metricsHandler may be nil (e.g. in tests) in
// which case /metrics is not registered.
func NewInternalMux(metricsHandler http.Handler, cb *circuitbreaker.Breaker) *http.ServeMux {
	mux := http.NewServeMux()

	if metricsHandler != nil {
		mux.Handle("GET /metrics", metricsHandler)
	}

	mux.Handle("GET /healthz", healthzHandler())
	mux.Handle("GET /readyz", readyzHandler(cb))
	mux.Handle("GET /startup", startupHandler())

	// pprof — registered explicitly on this internal mux only, so it is never
	// exposed on the MCP mux (:8080).
	mux.HandleFunc("GET /debug/pprof/", pprof.Index)
	mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)

	return mux
}
