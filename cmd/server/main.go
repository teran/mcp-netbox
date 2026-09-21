package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/config"
	"github.com/teran/mcp-netbox/handlers"
	"github.com/teran/mcp-netbox/infrastructure/circuitbreaker"
	infra "github.com/teran/mcp-netbox/infrastructure/netbox"
	"github.com/teran/mcp-netbox/logging"
)

// Build-time variables injected by goreleaser (via ldflags, B2).
var (
	appName       = "mcp-netbox"
	appVersion    = "dev"
	appCommitHash = "none"
	appTimestamp  = "unknown"
)

// bannerString returns the B5 startup banner built from the ldflags-embedded
// build metadata. It is emitted as the very first log line when logging is
// enabled (L6).
func bannerString() string {
	name := appName
	if name == "" {
		name = "mcp-netbox"
	}
	return fmt.Sprintf("Starting %s/%s (commit: %s; built at %s) ...", name, appVersion, appCommitHash, appTimestamp)
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}

	if err := Run(*cfg); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
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
	logger, err := logging.Setup(cfg)
	if err != nil {
		return fmt.Errorf("setup logging: %w", err)
	}

	// Share the configured logger with the HTTP middleware and the NetBox
	// client so that all logging honours LOG_LEVEL/LOG_FORMAT/LOG_FILENAME.
	handlers.SetLogger(logger)
	infra.SetLogger(logger)

	switch cfg.Transport {
	case config.TransportStdio:
		return runStdio(cfg, logger)
	case "", config.TransportHTTP:
		return runHTTP(cfg, logger)
	default:
		return fmt.Errorf("unsupported transport %q", cfg.Transport)
	}
}

// runStdio runs the MCP server over stdin/stdout. Because there is no HTTP
// Authorization header in this mode, the NetBox token must be provided via
// the NETBOX_TOKEN environment variable and is used by a single shared
// NetworkService for every tool call.
func runStdio(cfg config.Config, logger *logrus.Logger) error {
	return runStdioWithTransport(cfg, logger, &mcp.StdioTransport{})
}

// runStdioWithTransport is runStdio with an injectable transport, so that
// tests can drive the stdio server in-process over an in-memory pipe.
func runStdioWithTransport(cfg config.Config, logger *logrus.Logger, transport mcp.Transport) error {
	srv := newMCPServer(logger)

	httpClient, _ := newNetBoxHTTPClient(cfg)
	netboxClient := infra.NewClient(cfg.NetBoxURL, httpClient)
	svc := application.NewNetworkService(netboxClient, cfg.NetBoxToken)

	// In stdio mode there is no HTTP middleware chain, so no per-request
	// service injection. Register the shared service directly.
	handlers.RegisterTools(srv, nil, svc)

	// B5/L6: when logging is enabled, the startup banner is the very first log line.
	logger.Info(bannerString())

	conn, err := srv.Connect(context.Background(), transport, nil)
	if err != nil {
		return fmt.Errorf("connect stdio transport: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.WithError(err).Error("stdio transport close error")
		}
	}()

	return conn.Wait()
}

// runHTTP runs the MCP server over Streamable HTTP plus a separate internal
// observability server (metrics, pprof, probes). Tokens are read per request
// from the Authorization header.
func runHTTP(cfg config.Config, logger *logrus.Logger) error {
	srv := newMCPServer(logger)

	sharedHTTPClient, breaker := newNetBoxHTTPClient(cfg)

	promRegistry := prometheus.NewRegistry()
	metrics := handlers.NewMetrics(promRegistry)
	upstreamMetrics := infra.NewUpstreamMetrics(promRegistry)

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

	mux, stopRateLimit := handlers.NewMux(cfg, metrics, sharedHTTPClient, breaker, mcpHandler, upstreamMetrics)
	defer stopRateLimit()

	// B5/L6: when logging is enabled, the startup banner is the very first log line.
	logger.Info(bannerString())

	mainServer := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           mux,
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       120 * time.Second,
	}

	metricsHandler := handlers.RegisterMetricsOnRegistry(promRegistry)
	internalMux := handlers.NewInternalMux(metricsHandler, breaker)
	internalHandler := handlers.InstrumentInternalMux(promRegistry, internalMux)

	internalServer := &http.Server{
		Addr:              cfg.InternalAddr,
		Handler:           internalHandler,
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
		logger.WithField("addr", handlers.SanitizeLog(cfg.ListenAddr)).Info("mcp server listening")
		if err := mainServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()
	go func() {
		logger.WithField("addr", handlers.SanitizeLog(cfg.InternalAddr)).Info("internal observability server listening")
		if err := internalServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case sig := <-quit:
		logger.WithField("signal", sig.String()).Info("shutting down")
	case err := <-errCh:
		logger.WithError(err).Error("server error")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := mainServer.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Error("main server shutdown error")
	}
	if err := internalServer.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Error("internal server shutdown error")
	}

	logger.Info("server stopped gracefully")
	return nil
}

// serverInstructions are the general usage rules advertised to MCP clients.
const serverInstructions = `This server provides CRUD access to the NetBox inventory.

- Read tools (get_*) are safe and never modify NetBox.
- Write tools create and update objects (create_* / update_*); they are not
  idempotent. Update tools perform a partial PATCH merge: only the fields you
  provide are changed.
- Delete tools (delete_*) permanently and irreversibly remove an object and are
  marked with a destructive hint. Never call a delete tool unless the user has
  explicitly asked to delete that specific object.
- Validation is enforced by NetBox: a rejected create/update returns the DRF
  validation body to the model so it can correct the payload.
- All list tools are paginated: use page (1-based) and page_size (max 1000)
  to navigate results. An empty results array means no objects matched your
  filters — it is not an error.
- Filters are additive: provide only the fields relevant to your query.
- Authentication: in HTTP/SSE mode the NetBox token is passed in the
  Authorization: Bearer header. In STDIO mode it is provided via the
  NETBOX_TOKEN environment variable at server startup. Authorization is
  enforced entirely by NetBox; this server does not validate tokens.`

// newMCPServer builds the shared MCP server with the configured capabilities,
// instructions, and an SDK logger wired into the server's logrus logger so
// SDK-level events are visible in the logs (L7/G10).
func newMCPServer(logger *logrus.Logger) *mcp.Server {
	return mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-netbox",
		Version: appVersion,
	}, &mcp.ServerOptions{
		Capabilities: &mcp.ServerCapabilities{
			Tools: &mcp.ToolCapabilities{ListChanged: false},
		},
		Instructions: serverInstructions,
		Logger:       logging.NewSlogLogger(logger),
	})
}

// newNetBoxHTTPClient builds the shared HTTP client used to talk to NetBox,
// with DNS-rebinding protection and a circuit breaker transport. It returns
// the client and the circuit breaker (used by the internal readyz endpoint).
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
