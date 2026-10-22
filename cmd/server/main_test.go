package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/teran/mcp-netbox/config"
	"github.com/teran/mcp-netbox/infrastructure/circuitbreaker"
)

func freePort() string {
	listener, err := net.Listen("tcp", ":0") //nolint:noctx,gosec
	if err != nil {
		return ":0"
	}
	defer func() { _ = listener.Close() }()

	return listener.Addr().String()
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

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
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
	internalAddr := freePort()

	cfg := config.Config{
		NetBoxURL:          "http://netbox.example.com",
		ListenAddr:         addr,
		InternalAddr:       internalAddr,
		RateLimitGlobal:    1000,
		RateLimitPerClient: 100,
		WriteTimeout:       300 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- Run(cfg)
	}()

	t.Cleanup(func() {
		// Send SIGTERM to our own process (simplified)
	})

	if !waitForServer("http://"+internalAddr+"/healthz", 2*time.Second) {
		t.Fatal("Server did not start")
	}

	resp, err := http.Get("http://" + internalAddr + "/healthz") //nolint:noctx
	if err != nil {
		t.Fatalf("Health check failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	// O01/O04: the liveness probe must NOT be served on the MCP address (:8080).
	if resp, err := http.Get("http://" + addr + "/healthz"); err == nil { //nolint:noctx
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("MCP /healthz status = %d, want %d", resp.StatusCode, http.StatusNotFound)
		}
	}
}

func TestRun_MCPHandlerPing(t *testing.T) {
	addr := freePort()
	internalAddr := freePort()

	cfg := config.Config{
		NetBoxURL:          "http://netbox.example.com",
		ListenAddr:         addr,
		InternalAddr:       internalAddr,
		RateLimitGlobal:    1000,
		RateLimitPerClient: 100,
		WriteTimeout:       300 * time.Second,
	}

	go func() {
		_ = Run(cfg)
	}()

	if !waitForServer("http://"+internalAddr+"/healthz", 2*time.Second) {
		t.Fatal("Server did not start")
	}

	reqBody := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "ping",
	}
	body, _ := json.Marshal(reqBody) //nolint:errchkjson

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"http://"+addr+"/mcp",
		bytes.NewReader(body),
	)
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
	internalAddr := freePort()

	cfg := config.Config{
		NetBoxURL:          "http://netbox.example.com",
		ListenAddr:         addr,
		InternalAddr:       internalAddr,
		RateLimitGlobal:    1000,
		RateLimitPerClient: 100,
		WriteTimeout:       300 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- Run(cfg)
	}()

	if !waitForServer("http://"+internalAddr+"/healthz", 2*time.Second) {
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
	rawConn := dialerRawConn(t)

	// Loopback is always rejected, even when private is not allowed.
	// Exercise Control directly (no real network dial): a weakened reject
	// check fails fast instead of hanging on a loopback connect.
	err := dialer.Control("tcp", "127.0.0.1:8080", rawConn)
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
	rawConn := dialerRawConn(t)

	// Exercise Control directly instead of DialContext: dialing the private
	// 10.0.0.1 would hang for dialTimeout (30s) whenever the reject check is
	// weakened, turning a gremlins mutant into TIMED OUT instead of a fast
	// KILLED.
	err := dialer.Control("tcp", "10.0.0.1:8080", rawConn)
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

	listener, err := net.Listen("tcp", "127.0.0.1:0") //nolint:noctx
	if err != nil {
		t.Fatalf("Listen error: %v", err)
	}
	defer func() { _ = listener.Close() }()

	conn, err := net.Dial("tcp", listener.Addr().String()) //nolint:noctx
	if err != nil {
		t.Fatalf("Dial error: %v", err)
	}
	defer func() { _ = conn.Close() }()

	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		t.Fatalf("conn type = %T, want *net.TCPConn", conn)
	}
	raw, err := tcpConn.SyscallConn()
	if err != nil {
		t.Fatalf("SyscallConn error: %v", err)
	}

	return raw
}

func TestNewPrivateIPCheckingDialer_RejectsLoopbackEvenWhenAllowed(t *testing.T) {
	t.Parallel()

	// Loopback/unspecified is rejected regardless of allowPrivate.
	dialer := newPrivateIPCheckingDialer(true)
	rawConn := dialerRawConn(t)

	err := dialer.Control("tcp", "127.0.0.1:8080", rawConn)
	if err == nil {
		t.Fatal("expected loopback dial to be rejected even with allowPrivate=true")
	}
}

func TestNewPrivateIPCheckingDialer_BadAddress(t *testing.T) {
	t.Parallel()

	dialer := newPrivateIPCheckingDialer(false)
	rawConn := dialerRawConn(t)
	// Missing port makes SplitHostPort fail inside Control.
	err := dialer.Control("tcp", "127.0.0.1", rawConn)
	if err == nil {
		t.Fatal("expected dial with malformed address to fail")
	}
}

func TestNewPrivateIPCheckingDialer_DNSFailure(t *testing.T) {
	t.Parallel()

	dialer := newPrivateIPCheckingDialer(false)
	rawConn := dialerRawConn(t)
	// A hostname that cannot be resolved triggers the DNS-lookup error branch.
	err := dialer.Control("tcp", "nonexistent.invalid:8080", rawConn)
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

// httpClientTransport unwraps the circuit-breaker RoundTripper to expose the
// underlying *http.Transport for TLS configuration assertions.
func httpClientTransport(t *testing.T, client *http.Client) *http.Transport {
	t.Helper()

	cbRT, ok := client.Transport.(*circuitbreaker.RoundTripper)
	if !ok {
		t.Fatalf("client.Transport type = %T, want *circuitbreaker.RoundTripper", client.Transport)
	}
	tr, ok := cbRT.RoundTripper.(*http.Transport)
	if !ok {
		t.Fatalf("wrapped transport type = %T, want *http.Transport", cbRT.RoundTripper)
	}

	return tr
}

func TestNewNetBoxHTTPClient_CustomCA(t *testing.T) {
	t.Parallel()

	certPEM := testSelfSignedCertPEM(t)
	path := writeTempCAFile(t, certPEM)

	cfg := config.Config{NetBoxCAFile: path}
	client, _, err := newNetBoxHTTPClient(cfg)
	if err != nil {
		t.Fatalf("newNetBoxHTTPClient(custom CA) error: %v", err)
	}

	transport := httpClientTransport(t, client)
	if transport.TLSClientConfig == nil {
		t.Fatal("TLSClientConfig = nil, want set with custom RootCAs")
	}
	pool := transport.TLSClientConfig.RootCAs
	if pool == nil {
		t.Fatal("RootCAs = nil, want custom pool")
	}
	expected := x509.NewCertPool()
	expected.AppendCertsFromPEM(certPEM)
	if !expected.Equal(pool) {
		t.Error("RootCAs does not contain the custom CA certificate")
	}
}

func TestNewNetBoxHTTPClient_EmptyCAUsesSystemRoots(t *testing.T) {
	t.Parallel()

	client, _, err := newNetBoxHTTPClient(config.Config{})
	if err != nil {
		t.Fatalf("newNetBoxHTTPClient(no CA) error: %v", err)
	}

	transport := httpClientTransport(t, client)
	if transport.TLSClientConfig == nil {
		// No custom config and no TLSClientConfig means system roots are used.
		t.Log("TLSClientConfig is nil (system roots)")

		return
	}
	if transport.TLSClientConfig.RootCAs != nil {
		t.Error("RootCAs = non-nil, want nil (system roots) when no CA file is set")
	}
}

func TestNewNetBoxHTTPClient_MissingCAFileFails(t *testing.T) {
	t.Parallel()

	cfg := config.Config{NetBoxCAFile: filepath.Join(t.TempDir(), "missing.pem")}
	if _, _, err := newNetBoxHTTPClient(cfg); err == nil {
		t.Fatal("newNetBoxHTTPClient(missing CA) = nil error, want error")
	}
}

// writeTempCAFile writes certPEM to a temp file and returns its path.
func writeTempCAFile(t *testing.T, certPEM []byte) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "ca.pem")
	if err := os.WriteFile(path, certPEM, 0o600); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	return path
}

// testSelfSignedCertPEM generates a self-signed CA certificate and returns its
// PEM-encoded certificate block.
func testSelfSignedCertPEM(t *testing.T) []byte {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey error: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "test-ca"},
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("CreateCertificate error: %v", err)
	}

	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
}

func TestResolveTransport_ModeFlagWins(t *testing.T) {
	t.Setenv("TRANSPORT", config.TransportHTTP)
	if got := resolveTransport(config.TransportStdio); got != config.TransportStdio {
		t.Errorf("resolveTransport(stdio with TRANSPORT=http) = %q, want stdio (flag wins)", got)
	}
	if got := resolveTransport(config.TransportHTTP); got != config.TransportHTTP {
		t.Errorf("resolveTransport(http) = %q, want http", got)
	}
}

func TestResolveTransport_EnvOverride(t *testing.T) {
	t.Setenv("TRANSPORT", config.TransportHTTP)
	if got := resolveTransport(""); got != config.TransportHTTP {
		t.Errorf("resolveTransport(\"\") with TRANSPORT=http = %q, want http", got)
	}
}

func TestResolveTransport_DefaultStdio(t *testing.T) {
	t.Setenv("TRANSPORT", "")
	if got := resolveTransport(""); got != config.TransportStdio {
		t.Errorf("resolveTransport(\"\") with TRANSPORT unset = %q, want stdio (M6 default)", got)
	}
}
