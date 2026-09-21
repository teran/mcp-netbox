package netbox

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

// histogramSampleCount gathers a histogram collector and returns its observed
// sample count, so tests can assert how many observations were recorded.
func histogramSampleCount(c prometheus.Collector) uint64 {
	reg := prometheus.NewRegistry()
	reg.MustRegister(c)
	mfs, err := reg.Gather()
	if err != nil {
		return 0
	}
	for _, mf := range mfs {
		for _, m := range mf.GetMetric() {
			if h := m.GetHistogram(); h != nil {
				return h.GetSampleCount()
			}
		}
	}
	return 0
}

func TestUpstreamMetrics_Observe(t *testing.T) {
	t.Parallel()

	reg := prometheus.NewRegistry()
	um := NewUpstreamMetrics(reg)

	// A nil handle must be safe to call (no panic, no-op).
	var nilMetrics *UpstreamMetrics
	nilMetrics.observe(200, time.Millisecond, 512)

	um.observe(200, 100*time.Millisecond, 512)
	um.observe(404, 200*time.Millisecond, 64)
	um.observe(200, 300*time.Millisecond, 256)

	if got := testutil.ToFloat64(um.requestsTotal.WithLabelValues("200")); got != 2 {
		t.Errorf("requests_total{status=200} = %v, want 2", got)
	}
	if got := testutil.ToFloat64(um.requestsTotal.WithLabelValues("404")); got != 1 {
		t.Errorf("requests_total{status=404} = %v, want 1", got)
	}
	if got := histogramSampleCount(um.duration); got != 3 {
		t.Errorf("duration histogram sample count = %d, want 3", got)
	}
	if got := histogramSampleCount(um.responseSize); got != 3 {
		t.Errorf("response_size histogram sample count = %d, want 3", got)
	}
}

func TestClient_UpstreamMetrics(t *testing.T) {
	t.Parallel()

	reg := prometheus.NewRegistry()
	um := NewUpstreamMetrics(reg)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/dcim/sites/" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"count":0,"next":null,"previous":null,"results":[]}`))
			return
		}
		// Any other path returns 404 to exercise the error branch.
		http.NotFound(w, r)
	}))
	defer srv.Close()

	client := NewClientWithMetrics(srv.URL, http.DefaultClient, um)

	// Success path -> 200.
	if _, err := client.ListSites(context.Background(), "token", nil); err != nil {
		t.Fatalf("ListSites() returned error: %v", err)
	}
	// Error path -> 404.
	if _, err := client.GetObject(context.Background(), "token", "site", 999, nil); err == nil {
		t.Fatal("GetObject() expected error for missing object, got nil")
	}

	if got := testutil.ToFloat64(um.requestsTotal.WithLabelValues("200")); got != 1 {
		t.Errorf("requests_total{status=200} = %v, want 1", got)
	}
	if got := testutil.ToFloat64(um.requestsTotal.WithLabelValues("404")); got != 1 {
		t.Errorf("requests_total{status=404} = %v, want 1", got)
	}
	if got := histogramSampleCount(um.duration); got != 2 {
		t.Errorf("duration histogram sample count = %d, want 2", got)
	}
	if got := histogramSampleCount(um.responseSize); got != 2 {
		t.Errorf("response_size histogram sample count = %d, want 2", got)
	}
}

func TestNewUpstreamMetrics_Registers(t *testing.T) {
	t.Parallel()

	reg := prometheus.NewRegistry()
	um := NewUpstreamMetrics(reg)

	// Observe once so the counter vector has a value and the metric families
	// are emitted on gather.
	um.observe(200, time.Millisecond, 1)

	mfs, err := reg.Gather()
	if err != nil {
		t.Fatalf("Gather() returned error: %v", err)
	}
	gathered := make(map[string]bool)
	for _, mf := range mfs {
		gathered[mf.GetName()] = true
	}
	for _, name := range []string{
		"netbox_upstream_requests_total",
		"netbox_upstream_request_duration_seconds",
		"netbox_upstream_response_size_bytes",
	} {
		if !gathered[name] {
			t.Errorf("metric %q not present on registry", name)
		}
	}
}
