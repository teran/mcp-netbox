package config

import (
	"net"
	"net/netip"
	"testing"
	"time"
)

func TestValidateURLScheme(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{"valid http", "http://netbox:8000", false},
		{"valid https", "https://netbox:8000", false},
		{"invalid scheme", "ftp://netbox:8000", true},
		{"non-string input (int)", 123, true},
		{"non-string input (nil)", nil, true},
		{"invalid URL", "://invalid", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateURLScheme(tc.input)
			if tc.wantError && err == nil {
				t.Errorf("validateURLScheme(%v) = nil, want error", tc.input)
			}
			if !tc.wantError && err != nil {
				t.Errorf("validateURLScheme(%v) = %v, want nil", tc.input, err)
			}
		})
	}
}

func TestValidateURLHost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{"valid URL", "http://netbox:8000", false},
		{"empty host", "http://", true},
		{"non-string input (int)", 456, true},
		{"non-string input (nil)", nil, true},
		{"invalid URL", "://invalid", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateURLHost(tc.input)
			if tc.wantError && err == nil {
				t.Errorf("validateURLHost(%v) = nil, want error", tc.input)
			}
			if !tc.wantError && err != nil {
				t.Errorf("validateURLHost(%v) = %v, want nil", tc.input, err)
			}
		})
	}
}

func TestValidatePositiveInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{"positive int", 1, false},
		{"zero", 0, true},
		{"negative", -1, true},
		{"non-int input (string)", "abc", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePositiveInt(tc.input)
			if tc.wantError && err == nil {
				t.Errorf("validatePositiveInt(%v) = nil, want error", tc.input)
			}
			if !tc.wantError && err != nil {
				t.Errorf("validatePositiveInt(%v) = %v, want nil", tc.input, err)
			}
		})
	}
}

func TestValidateURLNotPrivate(t *testing.T) {
	t.Parallel()

	// Use the strict (non-allow) variant for all tests.
	strictValidator := validateURLNotPrivate(false)

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{"public IP", "http://8.8.8.8:8000", false},
		{"public IP (cloudflare)", "http://1.1.1.1", false},
		{"hostname", "http://google.com", false},
		{"hostname with path", "https://google.com/api", false},
		{"loopback IPv4", "http://127.0.0.1:8000", true},
		{"loopback IPv6", "http://[::1]:8000", true},
		{"private IPv4 (10.x)", "http://10.0.0.1:8000", true},
		{"private IPv4 (172.16.x)", "http://172.16.0.1:8000", true},
		{"private IPv4 (192.168.x)", "http://192.168.1.1:8000", true},
		{"link-local IPv4", "http://169.254.1.1:8000", true},
		{"link-local IPv6", "http://[fe80::1]:8000", true},
		{"unspecified IPv4", "http://0.0.0.0:8000", true},
		{"unspecified IPv6", "http://[::]:8000", true},
		{"non-string input", 123, true},
		{"non-string input (nil)", nil, true},
		{"invalid URL", "://invalid", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := strictValidator(tc.input)
			if tc.wantError && err == nil {
				t.Errorf("validateURLNotPrivate(%v) = nil, want error", tc.input)
			}
			if !tc.wantError && err != nil {
				t.Errorf("validateURLNotPrivate(%v) = %v, want nil", tc.input, err)
			}
		})
	}
}

func TestValidateURLNotPrivate_AllowPrivate(t *testing.T) {
	t.Parallel()

	allowValidator := validateURLNotPrivate(true)

	t.Run("bypasses validation when allowPrivate=true", func(t *testing.T) {
		if err := allowValidator("http://127.0.0.1:8000"); err != nil {
			t.Errorf("allowValidator(loopback) = %v, want nil", err)
		}
		if err := allowValidator("http://192.168.1.1:8000"); err != nil {
			t.Errorf("allowValidator(private) = %v, want nil", err)
		}
	})
}

func TestValidateConfig(t *testing.T) {
	t.Parallel()

	t.Run("valid config", func(t *testing.T) {
		cfg := Config{
			NetBoxURL:          "http://google.com",
			RateLimitGlobal:    100,
			RateLimitPerClient: 10,
			WriteTimeout:       300 * time.Second,
		}
		if err := cfg.validate(); err != nil {
			t.Errorf("validate() = %v, want nil", err)
		}
	})

	t.Run("missing URL fails validation", func(t *testing.T) {
		cfg := Config{
			NetBoxURL: "",
		}
		if err := cfg.validate(); err == nil {
			t.Error("validate() = nil, want error for empty URL")
		}
	})

	t.Run("zero rate limit fails validation", func(t *testing.T) {
		cfg := Config{
			NetBoxURL:          "http://netbox:8000",
			RateLimitGlobal:    0,
			RateLimitPerClient: 1,
		}
		if err := cfg.validate(); err == nil {
			t.Error("validate() = nil, want error for zero rate limit")
		}
	})
}

func TestIsPrivateAddr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		addr string
		want bool
	}{
		{"loopback", "127.0.0.1", true},
		{"loopback v6", "::1", true},
		{"private 10", "10.0.0.1", true},
		{"private 172", "172.16.0.1", true},
		{"private 192.168", "192.168.1.1", true},
		{"link-local unicast", "169.254.1.1", true},
		{"link-local multicast", "ff02::1", true},
		{"unspecified", "0.0.0.0", true},
		{"unspecified v6", "::", true},
		{"public", "8.8.8.8", false},
		{"public v6", "2606:4700:4700::1111", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			addr := mustParseAddr(t, tc.addr)
			if got := IsPrivateAddr(addr); got != tc.want {
				t.Errorf("IsPrivateAddr(%s) = %v, want %v", tc.addr, got, tc.want)
			}
		})
	}
}

func TestIsPrivateIP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ip   string // empty means nil IP
		want bool
	}{
		{"nil", "", false},
		{"loopback", "127.0.0.1", true},
		{"private", "10.1.2.3", true},
		{"public", "8.8.8.8", false},
		{"invalid", "not-an-ip", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var ip net.IP
			if tc.ip != "" {
				ip = net.ParseIP(tc.ip)
			}
			if got := IsPrivateIP(ip); got != tc.want {
				t.Errorf("IsPrivateIP(%v) = %v, want %v", tc.ip, got, tc.want)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Setenv("NETBOX_URL", "http://8.8.8.8")
		t.Setenv("LISTEN_ADDR", ":9090")
		t.Setenv("RATE_LIMIT_GLOBAL", "50")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error: %v", err)
		}
		if cfg.NetBoxURL != "http://8.8.8.8" {
			t.Errorf("NetBoxURL = %q, want %q", cfg.NetBoxURL, "http://8.8.8.8")
		}
		if cfg.ListenAddr != ":9090" {
			t.Errorf("ListenAddr = %q, want %q", cfg.ListenAddr, ":9090")
		}
		if cfg.RateLimitGlobal != 50 {
			t.Errorf("RateLimitGlobal = %d, want 50", cfg.RateLimitGlobal)
		}
	})

	t.Run("validation failure for private URL", func(t *testing.T) {
		t.Setenv("NETBOX_URL", "http://127.0.0.1")
		if _, err := Load(); err == nil {
			t.Error("Load() = nil, want error for private URL")
		}
	})

	t.Run("missing required URL", func(t *testing.T) {
		t.Setenv("NETBOX_URL", "")
		if _, err := Load(); err == nil {
			t.Error("Load() = nil, want error for empty URL")
		}
	})
}

// mustParseAddr parses a string into a netip.Addr, failing the test on error.
func mustParseAddr(t *testing.T, s string) netip.Addr {
	t.Helper()
	a, err := netip.ParseAddr(s)
	if err != nil {
		t.Fatalf("ParseAddr(%q) error: %v", s, err)
	}
	return a
}

func TestLoad_EnvconfigError(t *testing.T) {
	// RATE_LIMIT_GLOBAL is an int field; a non-numeric value triggers an
	// envconfig.Process error before validation runs.
	t.Setenv("NETBOX_URL", "http://8.8.8.8")
	t.Setenv("RATE_LIMIT_GLOBAL", "not-a-number")

	if _, err := Load(); err == nil {
		t.Error("Load() = nil, want error for invalid env value")
	}
}

func TestValidateLogLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{"empty disables logging", "", false},
		{"whitespace trimmed", "   ", false},
		{"valid level", "debug", false},
		{"valid uppercase level", "INFO", false},
		{"invalid level", "bogus", true},
		{"non-string input", 123, true},
		{"non-string input (nil)", nil, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateLogLevel(tc.input)
			if tc.wantError && err == nil {
				t.Errorf("validateLogLevel(%v) = nil, want error", tc.input)
			}
			if !tc.wantError && err != nil {
				t.Errorf("validateLogLevel(%v) = %v, want nil", tc.input, err)
			}
		})
	}
}

func TestValidateLogFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{"empty defaults to text", "", false},
		{"text", "text", false},
		{"json", "json", false},
		{"case-insensitive", "JSON", false},
		{"invalid format", "xml", true},
		{"non-string input", 123, true},
		{"non-string input (nil)", nil, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateLogFormat(tc.input)
			if tc.wantError && err == nil {
				t.Errorf("validateLogFormat(%v) = nil, want error", tc.input)
			}
			if !tc.wantError && err != nil {
				t.Errorf("validateLogFormat(%v) = %v, want nil", tc.input, err)
			}
		})
	}
}

func TestValidateTransport_NonString(t *testing.T) {
	t.Parallel()

	v := validateTransport("token")
	if err := v(123); err == nil {
		t.Error("validateTransport(non-string) = nil, want error")
	}
	if err := v(nil); err == nil {
		t.Error("validateTransport(nil) = nil, want error")
	}
}
