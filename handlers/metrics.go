package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/teran/mcp-netbox/redact"
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

// WrapToolHandler wraps an MCP tool handler with metrics recording and output
// sanitization. The metrics recording is skipped when metrics is nil (used by
// tests), but the text-form sanitization is always applied so tool output is
// consistently scrubbed of ANSI/control characters regardless of caller.
func WrapToolHandler[I, O any](metrics *Metrics, toolName string, handler mcp.ToolHandlerFor[I, O]) mcp.ToolHandlerFor[I, O] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in I) (*mcp.CallToolResult, O, error) {
		start := time.Now()
		result, out, err := handler(ctx, req, in)
		duration := time.Since(start)

		if metrics != nil {
			statusClass := "2xx"
			if err != nil || (result != nil && result.IsError) {
				statusClass = "error"
			}

			metrics.toolRequestsTotal.WithLabelValues(toolName, statusClass).Inc()
			metrics.toolDuration.WithLabelValues(toolName).Observe(duration.Seconds())
		}

		// S02: redact a deep copy of the successful typed output before it is
		// turned into the tool's text form, so any secret:"true" field is masked
		// in the model-visible output. The original out is not mutated. If the
		// redaction cannot be asserted back to the concrete type O, the
		// original value is kept.
		if err == nil && result != nil && !result.IsError {
			if redacted, ok := redact.Redact(out).(O); ok {
				out = redacted
			}
		}

		return sanitizeToolResult(result, out), out, err
	}
}

// sanitizeToolResult ensures the text form of a successful tool result is free
// of ANSI escape sequences and control characters (S09/N23). The SDK derives
// the text fallback from the same typed output it marshals into
// StructuredContent; we pre-populate the sanitized text so that a crafted
// string field in NetBox data can never inject terminal escapes into the model
// output. The typed StructuredContent is left untouched (it remains the raw
// JSON), so no valid data is corrupted. Error results and results that already
// carry explicit content are passed through unchanged.
func sanitizeToolResult(res *mcp.CallToolResult, out any) *mcp.CallToolResult {
	if res == nil || res.IsError || res.Content != nil {
		return res
	}
	jsonBytes, err := json.Marshal(out)
	if err != nil {
		return res
	}
	res.Content = []mcp.Content{&mcp.TextContent{Text: sanitizeOutput(string(jsonBytes))}}
	return res
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
