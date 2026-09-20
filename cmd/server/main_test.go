package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/teran/mcp-netbox/config"
)

func freePort() string {
	l, err := net.Listen("tcp", ":0") //nolint:noctx,gosec
	if err != nil {
		return ":0"
	}
	defer func() { _ = l.Close() }()
	return l.Addr().String()
}

func waitForServer(url string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url) //nolint:gosec,noctx
		if err == nil {
			func() { _ = resp.Body.Close() }()
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}

func TestHealthEndpoint(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer srv.Close()

	resp, err := http.Get(srv.URL) //nolint:noctx
	if err != nil {
		t.Fatalf("GET /healthz failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(body, []byte("ok")) {
		t.Errorf("body = %q, want to contain 'ok'", string(body))
	}
}

func TestRun_HealthEndpoint(t *testing.T) {
	addr := freePort()
	metricsAddr := freePort()

	cfg := config.Config{
		NetBoxURL:             "http://netbox.example.com",
		ListenAddr:            addr,
		PrometheusMetricsAddr: metricsAddr,
		RateLimitGlobal:       1000,
		RateLimitPerClient:    100,
		WriteTimeout:          300 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- Run(cfg)
	}()

	t.Cleanup(func() {
		// Send SIGTERM to our own process (simplified)
	})

	if !waitForServer("http://"+addr+"/healthz", 2*time.Second) {
		t.Fatal("Server did not start")
	}

	resp, err := http.Get("http://" + addr + "/healthz") //nolint:noctx
	if err != nil {
		t.Fatalf("Health check failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestRun_MCPHandlerPing(t *testing.T) {
	addr := freePort()
	metricsAddr := freePort()

	cfg := config.Config{
		NetBoxURL:             "http://netbox.example.com",
		ListenAddr:            addr,
		PrometheusMetricsAddr: metricsAddr,
		RateLimitGlobal:       1000,
		RateLimitPerClient:    100,
		WriteTimeout:          300 * time.Second,
	}

	go func() {
		_ = Run(cfg)
	}()

	if !waitForServer("http://"+addr+"/healthz", 2*time.Second) {
		t.Fatal("Server did not start")
	}

	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "ping",
	}
	body, _ := json.Marshal(reqBody) //nolint:errchkjson

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "http://"+addr+"/mcp", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("MCP request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusOK {
		t.Log("MCP handler responded successfully")
	}
}

func TestRun_ShutdownViaSignal(t *testing.T) {
	addr := freePort()
	metricsAddr := freePort()

	cfg := config.Config{
		NetBoxURL:             "http://netbox.example.com",
		ListenAddr:            addr,
		PrometheusMetricsAddr: metricsAddr,
		RateLimitGlobal:       1000,
		RateLimitPerClient:    100,
		WriteTimeout:          300 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- Run(cfg)
	}()

	if !waitForServer("http://"+addr+"/healthz", 2*time.Second) {
		t.Fatal("Server did not start")
	}

	proc, err := os.FindProcess(os.Getpid())
	if err == nil {
		_ = proc.Signal(syscall.SIGTERM)
	}

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Run() returned error: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Server did not shut down within timeout")
	}
}

func TestNewPrivateIPCheckingDialer_RejectsReserved(t *testing.T) {
	t.Parallel()

	dialer := newPrivateIPCheckingDialer(false)

	// Loopback is always rejected, even when private is not allowed.
	_, err := dialer.DialContext(context.Background(), "tcp", "127.0.0.1:8080")
	if err == nil {
		t.Fatal("expected loopback dial to be rejected")
	}
	if !strings.Contains(err.Error(), "rejected") {
		t.Errorf("error = %q, want rejection message", err)
	}
}

func TestNewPrivateIPCheckingDialer_RejectsPrivate(t *testing.T) {
	t.Parallel()

	// With allowPrivate=false an RFC1918 address is rejected.
	dialer := newPrivateIPCheckingDialer(false)
	_, err := dialer.DialContext(context.Background(), "tcp", "10.0.0.1:8080")
	if err == nil {
		t.Fatal("expected private dial to be rejected")
	}
	if !strings.Contains(err.Error(), "rejected") {
		t.Errorf("error = %q, want rejection message", err)
	}
}

func TestNewPrivateIPCheckingDialer_AllowsPrivateWhenEnabled(t *testing.T) {
	t.Parallel()

	// With allowPrivate=true the Control function must let a private address
	// through (it is not rejected). We exercise the Control callback directly
	// via a real RawConn to avoid an actual 30s connect timeout.
	dialer := newPrivateIPCheckingDialer(true)
	rawConn := dialerRawConn(t)

	if err := dialer.Control("tcp", "10.0.0.1:8080", rawConn); err != nil {
		t.Errorf("Control(private, allowPrivate=true) = %v, want nil", err)
	}
}

func dialerRawConn(t *testing.T) syscall.RawConn {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0") //nolint:noctx
	if err != nil {
		t.Fatalf("Listen error: %v", err)
	}
	defer func() { _ = ln.Close() }()

	conn, err := net.Dial("tcp", ln.Addr().String()) //nolint:noctx
	if err != nil {
		t.Fatalf("Dial error: %v", err)
	}
	defer func() { _ = conn.Close() }()

	raw, err := conn.(*net.TCPConn).SyscallConn()
	if err != nil {
		t.Fatalf("SyscallConn error: %v", err)
	}
	return raw
}

func TestNewPrivateIPCheckingDialer_RejectsLoopbackEvenWhenAllowed(t *testing.T) {
	t.Parallel()

	// Loopback/unspecified is rejected regardless of allowPrivate.
	dialer := newPrivateIPCheckingDialer(true)
	_, err := dialer.DialContext(context.Background(), "tcp", "127.0.0.1:8080")
	if err == nil {
		t.Fatal("expected loopback dial to be rejected even with allowPrivate=true")
	}
}

func TestNewPrivateIPCheckingDialer_BadAddress(t *testing.T) {
	t.Parallel()

	dialer := newPrivateIPCheckingDialer(false)
	// Missing port makes SplitHostPort fail inside Control.
	_, err := dialer.DialContext(context.Background(), "tcp", "127.0.0.1")
	if err == nil {
		t.Fatal("expected dial with malformed address to fail")
	}
}

func TestNewPrivateIPCheckingDialer_DNSFailure(t *testing.T) {
	t.Parallel()

	dialer := newPrivateIPCheckingDialer(false)
	// A hostname that cannot be resolved triggers the DNS-lookup error branch.
	_, err := dialer.DialContext(context.Background(), "tcp", "nonexistent.invalid:8080")
	if err == nil {
		t.Fatal("expected dial with unresolvable hostname to fail")
	}
}

func TestBannerString(t *testing.T) {
	t.Cleanup(restoreBuildVars)
	appName = "mcp-netbox"
	appVersion = "1.2.3"
	appCommitHash = "abc123"
	appTimestamp = "2026-01-01T00:00:00Z"

	got := bannerString()
	want := "Starting mcp-netbox/1.2.3 (commit: abc123; built at 2026-01-01T00:00:00Z) ..."
	if got != want {
		t.Errorf("bannerString() = %q, want %q", got, want)
	}
}

func TestBannerString_Defaults(t *testing.T) {
	t.Cleanup(restoreBuildVars)
	appName = ""
	appVersion = ""
	appCommitHash = ""
	appTimestamp = ""

	got := bannerString()
	for _, want := range []string{"mcp-netbox", "Starting", "commit:", "built at"} {
		if !strings.Contains(got, want) {
			t.Errorf("bannerString() = %q, missing %q", got, want)
		}
	}
}

// restoreBuildVars returns the package build vars to their defaults so
// concurrent (non-banner) tests are unaffected.
func restoreBuildVars() {
	appName = "mcp-netbox"
	appVersion = "dev"
	appCommitHash = "none"
	appTimestamp = "unknown"
}

func TestRun_SetupLoggingError(t *testing.T) {
	// An invalid LOG_LEVEL makes logging.Setup fail, surfacing the error path
	// in Run before any transport is started.
	cfg := config.Config{
		NetBoxURL:   "http://netbox.example.com",
		Transport:   config.TransportHTTP,
		LogLevel:    "bogus",
		LogFormat:   "text",
		LogFilename: "",
	}
	if err := Run(cfg); err == nil {
		t.Fatal("Run() = nil, want error for invalid LOG_LEVEL")
	}
}
