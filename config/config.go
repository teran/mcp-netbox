// Package config provides configuration loading from environment variables
// using kelseyhightower/envconfig and validation via ozzo-validation.
package config

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// IsPrivateAddr reports whether addr is a loopback, private, link-local, or
// unspecified address. It uses netip.Addr which is the modern Go IP type.
func IsPrivateAddr(addr netip.Addr) bool {
	return addr.IsLoopback() || addr.IsPrivate() || addr.IsLinkLocalUnicast() || addr.IsLinkLocalMulticast() || addr.IsUnspecified()
}

// IsPrivateIP reports whether ip is a loopback, private, link-local, or
// unspecified address.
func IsPrivateIP(ip net.IP) bool {
	if ip == nil {
		return false
	}
	a, err := netip.ParseAddr(ip.String())
	if err != nil {
		return false
	}
	return IsPrivateAddr(a)
}

// Transport modes supported by the server. HTTP serves the MCP protocol over
// Streamable HTTP (the token is taken from the per-request Authorization
// header), while STDIO serves it over stdin/stdout (the token is taken from
// the NETBOX_TOKEN environment variable).
const (
	TransportHTTP  = "http"
	TransportStdio = "stdio"
)

type Config struct {
	NetBoxURL   string `envconfig:"NETBOX_URL" required:"true"`
	NetBoxToken string `envconfig:"NETBOX_TOKEN" default:"" secret:"true"`
	Transport   string `envconfig:"TRANSPORT" default:"http"`
	ListenAddr  string `envconfig:"LISTEN_ADDR" default:":8080"`
	// InternalAddr is the internal observability endpoint (default :8081). It
	// is separate from ListenAddr and serves Prometheus metrics, pprof, and the
	// healthz/readyz/startup probes (O01/O04).
	InternalAddr       string        `envconfig:"INTERNAL_ADDR" default:":8081"`
	RateLimitGlobal    int           `envconfig:"RATE_LIMIT_GLOBAL" default:"100"`
	RateLimitPerClient int           `envconfig:"RATE_LIMIT_PER_CLIENT" default:"10"`
	TrustedProxy       string        `envconfig:"TRUSTED_PROXY" default:""`
	WriteTimeout       time.Duration `envconfig:"WRITE_TIMEOUT" default:"300s"`
	AllowPrivateNetBox bool          `envconfig:"ALLOW_PRIVATE_NETBOX" default:"false"`
	LogLevel           string        `envconfig:"LOG_LEVEL" default:""`
	LogFormat          string        `envconfig:"LOG_FORMAT" default:"text"`
	LogFilename        string        `envconfig:"LOG_FILENAME" default:"/tmp/mcp-netbox.log"`
}

func (c Config) validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.NetBoxURL,
			validation.Required,
			validation.By(validateURLScheme),
			validation.By(validateURLHost),
			validation.By(validateURLNotPrivate(c.AllowPrivateNetBox)),
		),
		validation.Field(&c.Transport, validation.By(validateTransport(c.NetBoxToken))),
		validation.Field(&c.RateLimitGlobal, validation.By(validatePositiveInt)),
		validation.Field(&c.RateLimitPerClient, validation.By(validatePositiveInt)),
		validation.Field(&c.WriteTimeout, validation.Min(time.Second)),
		validation.Field(&c.LogLevel, validation.By(validateLogLevel)),
		validation.Field(&c.LogFormat, validation.By(validateLogFormat)),
	)
}

// validateLogLevel checks that LOG_LEVEL, when set, is a valid logrus level.
// An empty value means logging is disabled.
func validateLogLevel(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("must be a string")
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if _, err := logrus.ParseLevel(s); err != nil {
		return fmt.Errorf("invalid log level %q (want one of trace, debug, info, warn, error, fatal, panic)", s)
	}
	return nil
}

// validateLogFormat checks that LOG_FORMAT is one of the supported formats.
func validateLogFormat(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("must be a string")
	}
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "text", "json":
		return nil
	default:
		return fmt.Errorf("invalid log format %q (want text or json)", s)
	}
}

// validateTransport checks that the transport is one of the supported modes
// and that a NETBOX_TOKEN is provided for the STDIO transport (which has no
// HTTP Authorization header to carry it per request).
func validateTransport(netboxToken string) validation.RuleFunc {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("must be a string")
		}
		// The empty value means HTTP, which is the default and is what callers
		// constructing config.Config directly (without envconfig defaults) get.
		switch s {
		case "", TransportHTTP:
			return nil
		case TransportStdio:
			if strings.TrimSpace(netboxToken) == "" {
				return fmt.Errorf("%s transport requires NETBOX_TOKEN to be set", TransportStdio)
			}
			return nil
		default:
			return fmt.Errorf("must be %q or %q (got %q)", TransportHTTP, TransportStdio, s)
		}
	}
}

func validatePositiveInt(value interface{}) error {
	n, ok := value.(int)
	if !ok {
		return fmt.Errorf("must be an integer")
	}
	if n < 1 {
		return fmt.Errorf("must be at least 1 (got %d)", n)
	}
	return nil
}

func validateURLScheme(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("must be a string")
	}
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("must use http or https scheme (got %q)", u.Scheme)
	}
	return nil
}

func validateURLHost(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("must be a string")
	}
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	if u.Host == "" {
		return fmt.Errorf("must include a host (e.g. http://netbox:8000)")
	}
	return nil
}

func validateURLNotPrivate(allowPrivate bool) validation.RuleFunc {
	return func(value interface{}) error {
		if allowPrivate {
			return nil
		}
		return validateURLNotPrivateCheck(value)
	}
}

func validateURLNotPrivateCheck(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("must be a string")
	}
	u, err := url.Parse(s)
	if err != nil {
		return err
	}

	host := u.Hostname()
	// Strip IPv6 brackets if present
	host = strings.Trim(host, "[]")

	// Try to parse as IP address
	ip := net.ParseIP(host)
	if ip == nil {
		// Hostname — resolve DNS to check for private IPs.
		// This prevents SSRF bypass where a hostname resolves to a private IP.
		ips, err := net.DefaultResolver.LookupNetIP(context.Background(), "ip", host)
		if err != nil {
			// DNS resolution failed — reject the hostname to prevent SSRF bypass.
			// The operator must either fix DNS, use a literal IP, or set ALLOW_PRIVATE_NETBOX=true.
			return fmt.Errorf("hostname %q DNS lookup failed (%w); set ALLOW_PRIVATE_NETBOX=true to bypass", host, err)
		}
		for _, resolvedIP := range ips {
			if IsPrivateAddr(resolvedIP) {
				return fmt.Errorf("hostname %q resolves to private/reserved IP %q (set ALLOW_PRIVATE_NETBOX=true to bypass)", host, resolvedIP)
			}
		}
		return nil
	}

	if ip.IsLoopback() {
		return fmt.Errorf("must not be a loopback address (got %q)", host)
	}
	if ip.IsPrivate() {
		return fmt.Errorf("must not be a private IP address (got %q)", host)
	}
	if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return fmt.Errorf("must not be a link-local address (got %q)", host)
	}
	if ip.IsUnspecified() {
		return fmt.Errorf("must not be an unspecified address (got %q)", host)
	}

	return nil
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}
	return &cfg, nil
}
