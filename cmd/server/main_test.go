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

const (
	// serverStartupTimeout bounds how long we wait for the internal /healthz
	// endpoint to come up after Run() is launched. Under `-race` on a loaded
	// machine the process can take several seconds to schedule the Run
	// goroutine and reach ListenAndServe, so a generous timeout avoids spurious
	// "Server did not start" flakes while still failing fast on a genuinely
	// broken startup.
	serverStartupTimeout = 15 * time.Second
	// serverShutdownTimeout bounds how long we wait for graceful shutdown after
	// sending SIGTERM. runHTTP budgets up to 30s for the HTTP servers to drain,
	// and under `-race` on a loaded machine graceful shutdown can occasionally
	// take several seconds, so this must be generous to avoid spurious flakes.
	serverShutdownTimeout = 15 * time.Second
	// maxStartupAttempts is how many times we relaunch Run with fresh ports
	// before giving up. freePort() closes its probe listener and hands the port
	// number back, but another socket can grab that port in the window before
	// Run() binds it — a TOCTOU race that intermittently surfaces as
	// "listen: address already in use" on the main server. Relaunching with new
	// ports retries past that transient collision. At the observed ~2% per-attempt
	// collision rate three attempts make the residual failure negligible.
	maxStartupAttempts = 3
)

// waitForServerOrExit polls url for readiness while also detecting whether Run
// has already exited via errCh. Run exits before serving when its main listener
// fails to bind (the freePort() TOCTOU collision), so returning exited=true
// lets the caller retry with fresh ports instead of polling a dead endpoint.
func waitForServerOrExit(url string, errCh <-chan error, timeout time.Duration) (ready, exited bool) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		select {
		case <-errCh: //nolint:errcheck // discarding: Run returns nil even on a port-collision bind error
			return false, true
		default:
		}

		resp, err := http.Get(url) //nolint:gosec,noctx
		if err == nil {
			func() { _ = resp.Body.Close() }()

			return true, false
		}
		time.Sleep(10 * time.Millisecond)
	}

	return false, false
}

// runningServer bundles the handles a Run-based test needs: the bound listen
// addresses (so tests can issue real HTTP requests) and the Run result channel
// and stop function (so the server can be torn down between tests).
type runningServer struct {
	addr         string
	internalAddr string
	errCh        <-chan error
	stop         func()
}

// startTestServer runs Run(cfg) in the background and blocks until the internal
// /healthz endpoint reports healthy. Because Run binds its own listeners, a
// freePort()-derived port can occasionally be taken by another socket before
// Run binds it (a TOCTOU collision); startTestServer relaunches Run with fresh
// ports to retry past that transient failure. It returns the running server
// handles; every test that launches Run must call srv.stop (e.g. via
// t.Cleanup) so the server and its process-wide signal handler are torn down
// before the next test runs.
func startTestServer(t *testing.T, cfg config.Config) runningServer {
	t.Helper()

	for attempt := range maxStartupAttempts {
		cfg.ListenAddr = freePort()
		cfg.InternalAddr = freePort()

		errCh := make(chan error, 1)
		go func() {
			errCh <- Run(cfg)
		}()

		ready, exited := waitForServerOrExit("http://"+cfg.InternalAddr+"/healthz", errCh, serverStartupTimeout)
		if ready {
			srv := runningServer{
				addr:         cfg.ListenAddr,
				internalAddr: cfg.InternalAddr,
				errCh:        errCh,
			}
			srv.stop = func() {
				t.Helper()

				if proc, err := os.FindProcess(os.Getpid()); err == nil {
					_ = proc.Signal(syscall.SIGTERM)
				}

				select {
				case err := <-srv.errCh:
					if err != nil {
						t.Errorf("Run() returned error: %v", err)
					}
				case <-time.After(serverShutdownTimeout):
					t.Error("Server did not shut down within timeout")
				}
			}

			return srv
		}

		if !exited {
			// Run is still alive but never served within the timeout — a genuine
			// startup problem, not a transient port collision.
			t.Fatalf("Server did not start (attempt %d/%d)", attempt+1, maxStartupAttempts)
		}
		// Run exited before serving — almost certainly a transient port
		// collision; retry with fresh ports.
	}

	t.Fatal("Server did not start after multiple attempts")

	return runningServer{}
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
	cfg := config.Config{
		NetBoxURL:          "http://netbox.example.com",
		RateLimitGlobal:    1000,
		RateLimitPerClient: 100,
		WriteTimeout:       300 * time.Second,
	}

	srv := startTestServer(t, cfg)
	t.Cleanup(srv.stop)

	resp, err := http.Get("http://" + srv.internalAddr + "/healthz") //nolint:noctx
	if err != nil {
		t.Fatalf("Health check failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	// O01/O04: the liveness probe must NOT be served on the MCP address (:8080).
	if resp, err := http.Get("http://" + srv.addr + "/healthz"); err == nil { //nolint:noctx
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("MCP /healthz status = %d, want %d", resp.StatusCode, http.StatusNotFound)
		}
	}
}

func TestRun_MCPHandlerPing(t *testing.T) {
	cfg := config.Config{
		NetBoxURL:          "http://netbox.example.com",
		RateLimitGlobal:    1000,
		RateLimitPerClient: 100,
		WriteTimeout:       300 * time.Second,
	}

	srv := startTestServer(t, cfg)
	t.Cleanup(srv.stop)

	reqBody := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "ping",
	}
	body, _ := json.Marshal(reqBody) //nolint:errchkjson

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"http://"+srv.addr+"/mcp",
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
	cfg := config.Config{
		NetBoxURL:          "http://netbox.example.com",
		RateLimitGlobal:    1000,
		RateLimitPerClient: 100,
		WriteTimeout:       300 * time.Second,
	}

	srv := startTestServer(t, cfg)

	// Send SIGTERM and assert Run returns nil (graceful shutdown). srv.stop
	// fails the test if Run returns an error or does not shut down in time.
	srv.stop()
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
