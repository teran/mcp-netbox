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

	"github.com/teran/mcp-netbox/config"
	"github.com/teran/mcp-netbox/handlers"
	"github.com/teran/mcp-netbox/infrastructure/circuitbreaker"
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

func Run(cfg config.Config) error {
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

	sharedHTTPClient := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: cbTransport,
	}

	srv := mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-netbox",
		Version: version,
	}, &mcp.ServerOptions{
		Capabilities: &mcp.ServerCapabilities{
			Tools: &mcp.ToolCapabilities{ListChanged: false},
		},
	})

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

	// Register tools on the MCP server
	handlers.RegisterTools(srv, metrics, nil)

	mux, stopRateLimit := handlers.NewMux(cfg, metrics, sharedHTTPClient, cbTransport.Breaker(), mcpHandler)
	defer stopRateLimit()

	u, _ := url.Parse(cfg.NetBoxURL)
	slog.Info("starting server", "netbox_url", handlers.SanitizeLog(u.Redacted()))
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
