package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	toolRequestsTotal *prometheus.CounterVec
	toolDuration      *prometheus.HistogramVec
	activeRequests    prometheus.Gauge
}

func NewMetrics(reg *prometheus.Registry) *Metrics {
	m := &Metrics{
		toolRequestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "mcp_tool_requests_total",
			Help: "Per-tool request count.",
		}, []string{"tool", "status_class"}),
		toolDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "mcp_tool_duration_seconds",
			Help:    "Per-tool request duration.",
			Buckets: prometheus.DefBuckets,
		}, []string{"tool"}),
		activeRequests: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "mcp_active_requests",
			Help: "Current in-flight MCP requests.",
		}),
	}
	reg.MustRegister(m.toolRequestsTotal, m.toolDuration, m.activeRequests)
	return m
}

// WrapToolHandler wraps an MCP tool handler with metrics recording.
func WrapToolHandler[I, O any](metrics *Metrics, toolName string, handler mcp.ToolHandlerFor[I, O]) mcp.ToolHandlerFor[I, O] {
	if metrics == nil {
		return handler
	}
	return func(ctx context.Context, req *mcp.CallToolRequest, in I) (*mcp.CallToolResult, O, error) {
		start := time.Now()
		result, out, err := handler(ctx, req, in)
		duration := time.Since(start)

		statusClass := "2xx"
		if err != nil || (result != nil && result.IsError) {
			statusClass = "error"
		}

		metrics.toolRequestsTotal.WithLabelValues(toolName, statusClass).Inc()
		metrics.toolDuration.WithLabelValues(toolName).Observe(duration.Seconds())

		return result, out, err
	}
}

// MetricsMiddleware tracks active request count via gauge.
func MetricsMiddleware(metrics *Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if metrics != nil {
				metrics.activeRequests.Inc()
				defer metrics.activeRequests.Dec()
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RegisterMetricsOnRegistry registers Go runtime metrics and returns an HTTP handler.
func RegisterMetricsOnRegistry(reg *prometheus.Registry) http.Handler {
	reg.MustRegister(collectors.NewGoCollector())
	return promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
}
