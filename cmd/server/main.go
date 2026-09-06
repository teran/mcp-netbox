package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/config"
	"github.com/teran/mcp-netbox/handlers"
	"github.com/teran/mcp-netbox/infrastructure/circuitbreaker"
	infra "github.com/teran/mcp-netbox/infrastructure/netbox"
)

// Build-time variables injected by goreleaser (via ldflags).
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := Run(*cfg); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// newPrivateIPCheckingDialer returns a net.Dialer whose DialContext rejects
// connections to loopback, private, and link-local IPs. This prevents DNS
// rebinding attacks where a hostname resolves to a private IP at dial time.
//
// When allowPrivate is true, only loopback, link-local, and unspecified
// addresses are rejected (private RFC 1918 ranges are allowed).
func newPrivateIPCheckingDialer(allowPrivate bool) *net.Dialer {
	rejectAllPrivate := !allowPrivate
	return &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		Control: func(network, address string, c syscall.RawConn) error {
			host, _, err := net.SplitHostPort(address)
			if err != nil {
				return err
			}
			ips, err := net.DefaultResolver.LookupNetIP(context.Background(), "ip", host)
			if err != nil {
				return fmt.Errorf("DNS lookup for %q failed: %w", host, err)
			}
			for _, ip := range ips {
				if config.IsPrivateAddr(ip) {
					// Always reject loopback and unspecified, even when private is allowed.
					if ip.IsLoopback() || ip.IsUnspecified() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
						return fmt.Errorf("connection to reserved IP %q rejected (DNS rebinding protection)", ip)
					}
					if rejectAllPrivate {
						return fmt.Errorf("connection to private IP %q rejected (DNS rebinding protection; set ALLOW_PRIVATE_NETBOX=true to bypass)", ip)
					}
				}
			}
			return nil
		},
	}
}

// Run starts the server using the configured transport.
//
//   - TransportHTTP: serves MCP over Streamable HTTP. The NetBox token is
//     taken from the per-request Authorization: Bearer header and injected
//     per request by the HTTP middleware chain.
//   - TransportStdio: serves MCP over stdin/stdout (newline-delimited JSON).
//     The NetBox token is taken from the NETBOX_TOKEN environment variable at
//     startup and passed to a single shared NetworkService.
//
// The empty transport string is treated as HTTP for backward compatibility
// with callers that construct config.Config directly.
func Run(cfg config.Config) error {
	switch cfg.Transport {
	case config.TransportStdio:
		return runStdio(cfg)
	case "", config.TransportHTTP:
		return runHTTP(cfg)
	default:
		return fmt.Errorf("unsupported transport %q", cfg.Transport)
	}
}

// runStdio runs the MCP server over stdin/stdout. Because there is no HTTP
// Authorization header in this mode, the NetBox token must be provided via
// the NETBOX_TOKEN environment variable and is used by a single shared
// NetworkService for every tool call.
func runStdio(cfg config.Config) error {
	return runStdioWithTransport(cfg, &mcp.StdioTransport{})
}

// runStdioWithTransport is runStdio with an injectable transport, so that
// tests can drive the stdio server in-process over an in-memory pipe.
func runStdioWithTransport(cfg config.Config, transport mcp.Transport) error {
	srv := newMCPServer()

	httpClient, _ := newNetBoxHTTPClient(cfg)
	netboxClient := infra.NewClient(cfg.NetBoxURL, httpClient)
	svc := application.NewNetworkService(netboxClient, cfg.NetBoxToken)

	// In stdio mode there is no HTTP middleware chain, so no per-request
	// service injection. Register the shared service directly.
	handlers.RegisterTools(srv, nil, svc)

	slog.Info("starting stdio MCP server", "netbox_url", handlers.SanitizeLog(redactedURL(cfg.NetBoxURL)))

	conn, err := srv.Connect(context.Background(), transport, nil)
	if err != nil {
		return fmt.Errorf("connect stdio transport: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			slog.Error("stdio transport close error", "error", err)
		}
	}()

	return conn.Wait()
}

// runHTTP runs the MCP server over Streamable HTTP plus a Prometheus metrics
// server. Tokens are read per request from the Authorization header.
func runHTTP(cfg config.Config) error {
	srv := newMCPServer()

	sharedHTTPClient, breaker := newNetBoxHTTPClient(cfg)

	promRegistry := prometheus.NewRegistry()
	metrics := handlers.NewMetrics(promRegistry)

	mcpHandler := mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server {
			return srv
		},
		&mcp.StreamableHTTPOptions{
			Stateless: true,
		},
	)

	// Register tools on the MCP server. The service is nil in HTTP mode; it is
	// injected per request by the middleware chain using the request token.
	handlers.RegisterTools(srv, metrics, nil)

	mux, stopRateLimit := handlers.NewMux(cfg, metrics, sharedHTTPClient, breaker, mcpHandler)
	defer stopRateLimit()

	slog.Info("starting server", "netbox_url", handlers.SanitizeLog(redactedURL(cfg.NetBoxURL)))
	slog.Info("build info", "version", version, "commit", commit, "date", date)

	mainServer := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           mux,
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       120 * time.Second,
	}

	metricsHandler := handlers.RegisterMetricsOnRegistry(promRegistry)
	metricsMux := http.NewServeMux()
	metricsMux.Handle("GET /metrics", metricsHandler)

	metricsServer := &http.Server{
		Addr:              cfg.PrometheusMetricsAddr,
		Handler:           metricsMux,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(quit)

	// Start servers and wait for signal or error
	errCh := make(chan error, 2)
	go func() {
		slog.Info("mcp server listening", "addr", handlers.SanitizeLog(cfg.ListenAddr))
		if err := mainServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()
	go func() {
		slog.Info("metrics server listening", "addr", handlers.SanitizeLog(cfg.PrometheusMetricsAddr))
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case sig := <-quit:
		slog.Info("shutting down", "signal", sig)
	case err := <-errCh:
		slog.Error("server error", "error", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := mainServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("main server shutdown error", "error", err)
	}
	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("metrics server shutdown error", "error", err)
	}

	slog.Info("server stopped gracefully")
	return nil
}

// newMCPServer builds the shared MCP server with the configured capabilities
// and instructions.
func newMCPServer() *mcp.Server {
	return mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-netbox",
		Version: version,
	}, &mcp.ServerOptions{
		Capabilities: &mcp.ServerCapabilities{
			Tools: &mcp.ToolCapabilities{ListChanged: false},
		},
	})
}

// newNetBoxHTTPClient builds the shared HTTP client used to talk to NetBox,
// with DNS-rebinding protection and a circuit breaker transport. It returns
// the client and the circuit breaker (used by the HTTP readyz endpoint).
func newNetBoxHTTPClient(cfg config.Config) (*http.Client, *circuitbreaker.Breaker) {
	dialer := newPrivateIPCheckingDialer(cfg.AllowPrivateNetBox)

	cbTransport := circuitbreaker.NewRoundTripper(
		&http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  false,
			DisableKeepAlives:   false,
			DialContext:         dialer.DialContext,
		},
		circuitbreaker.DefaultConfig(),
	)

	return &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: cbTransport,
	}, cbTransport.Breaker()
}

// redactedURL returns the NetBox URL with any credentials redacted for safe
// logging.
func redactedURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	return u.Redacted()
}
