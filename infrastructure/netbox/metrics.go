package netbox

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// UpstreamMetrics records Prometheus metrics about HTTP calls made to the
// NetBox upstream. It is nil-safe: all methods no-op when the receiver is nil,
// so the client can be used without metrics instrumentation (O03).
type UpstreamMetrics struct {
	requestsTotal *prometheus.CounterVec
	duration      prometheus.Histogram
	responseSize  prometheus.Histogram
}

// upstreamSizeBucketStart is the first bucket of the upstream size histogram.
const upstreamSizeBucketStart = 1024

// upstreamSizeBucketBase is the multiplicative factor between upstream size buckets.
const upstreamSizeBucketBase = 10

// upstreamSizeBucketCount is the number of buckets in the upstream size histogram.
const upstreamSizeBucketCount = 8

// NewUpstreamMetrics creates and registers the NetBox upstream metrics on the
// supplied registry, returning a nil-safe handle.
func NewUpstreamMetrics(reg *prometheus.Registry) *UpstreamMetrics {
	metrics := &UpstreamMetrics{
		requestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "netbox_upstream_requests_total",
			Help: "Total number of HTTP requests made to the NetBox upstream, partitioned by response status code.",
		}, []string{"status"}),
		duration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "netbox_upstream_request_duration_seconds",
			Help:    "Time spent waiting for NetBox upstream responses.",
			Buckets: prometheus.DefBuckets,
		}),
		responseSize: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "netbox_upstream_response_size_bytes",
			Help:    "Size of NetBox upstream response bodies.",
			Buckets: prometheus.ExponentialBuckets(upstreamSizeBucketStart, upstreamSizeBucketBase, upstreamSizeBucketCount),
		}),
	}
	reg.MustRegister(metrics.requestsTotal, metrics.duration, metrics.responseSize)

	return metrics
}

// observe records one upstream call: the response status counter, the request
// latency and the response body size. It is nil-safe.
func (m *UpstreamMetrics) observe(status int, duration time.Duration, size int) {
	if m == nil {
		return
	}
	m.requestsTotal.WithLabelValues(strconv.Itoa(status)).Inc()
	m.duration.Observe(duration.Seconds())
	m.responseSize.Observe(float64(size))
}
