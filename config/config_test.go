package config

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"net/netip"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestValidateURLScheme(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     any
		wantError bool
	}{
		{"valid http", "http://netbox:8000", false},
		{"valid https", "https://netbox:8000", false},
		{"invalid scheme", "ftp://netbox:8000", true},
		{"non-string input (int)", 123, true},
		{"non-string input (nil)", nil, true},
		{"invalid URL", "://invalid", true},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := validateURLScheme(testCase.input)
			if testCase.wantError && err == nil {
				t.Errorf("validateURLScheme(%v) = nil, want error", testCase.input)
			}
			if !testCase.wantError && err != nil {
				t.Errorf("validateURLScheme(%v) = %v, want nil", testCase.input, err)
			}
		})
	}
}

func TestValidateURLHost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     any
		wantError bool
	}{
		{"valid URL", "http://netbox:8000", false},
		{"empty host", "http://", true},
		{"non-string input (int)", 456, true},
		{"non-string input (nil)", nil, true},
		{"invalid URL", "://invalid", true},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := validateURLHost(testCase.input)
			if testCase.wantError && err == nil {
				t.Errorf("validateURLHost(%v) = nil, want error", testCase.input)
			}
			if !testCase.wantError && err != nil {
				t.Errorf("validateURLHost(%v) = %v, want nil", testCase.input, err)
			}
		})
	}
}

func TestValidatePositiveInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     any
		wantError bool
	}{
		{"positive int", 1, false},
		{"zero", 0, true},
		{"negative", -1, true},
		{"non-int input (string)", "abc", true},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := validatePositiveInt(testCase.input)
			if testCase.wantError && err == nil {
				t.Errorf("validatePositiveInt(%v) = nil, want error", testCase.input)
			}
			if !testCase.wantError && err != nil {
				t.Errorf("validatePositiveInt(%v) = %v, want nil", testCase.input, err)
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
		input     any
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := strictValidator(testCase.input)
			if testCase.wantError && err == nil {
				t.Errorf("validateURLNotPrivate(%v) = nil, want error", testCase.input)
			}
			if !testCase.wantError && err != nil {
				t.Errorf("validateURLNotPrivate(%v) = %v, want nil", testCase.input, err)
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			addr := mustParseAddr(t, testCase.addr)
			if got := IsPrivateAddr(addr); got != testCase.want {
				t.Errorf("IsPrivateAddr(%s) = %v, want %v", testCase.addr, got, testCase.want)
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var ip net.IP
			if testCase.ip != "" {
				ip = net.ParseIP(testCase.ip)
			}
			if got := IsPrivateIP(ip); got != testCase.want {
				t.Errorf("IsPrivateIP(%v) = %v, want %v", testCase.ip, got, testCase.want)
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
		input     any
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := validateLogLevel(testCase.input)
			if testCase.wantError && err == nil {
				t.Errorf("validateLogLevel(%v) = nil, want error", testCase.input)
			}
			if !testCase.wantError && err != nil {
				t.Errorf("validateLogLevel(%v) = %v, want nil", testCase.input, err)
			}
		})
	}
}

func TestValidateLogFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     any
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := validateLogFormat(testCase.input)
			if testCase.wantError && err == nil {
				t.Errorf("validateLogFormat(%v) = nil, want error", testCase.input)
			}
			if !testCase.wantError && err != nil {
				t.Errorf("validateLogFormat(%v) = %v, want nil", testCase.input, err)
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

func TestLoadCertPool_ValidPEM(t *testing.T) {
	t.Parallel()

	validPEM := testSelfSignedCertPEM(t)

	path := writeTempFile(t, validPEM)
	pool, err := LoadCertPool(path)
	if err != nil {
		t.Fatalf("LoadCertPool(valid) = error %v, want nil", err)
	}
	if pool == nil {
		t.Fatal("LoadCertPool(valid) = nil pool, want non-nil")
	}
	// The returned pool must contain exactly the single custom certificate.
	expected := x509.NewCertPool()
	expected.AppendCertsFromPEM(validPEM)
	if !expected.Equal(pool) {
		t.Error("LoadCertPool(valid) pool does not contain the custom certificate")
	}
}

func TestLoadCertPool_MissingFileFails(t *testing.T) {
	t.Parallel()

	_, err := LoadCertPool(filepath.Join(t.TempDir(), "does-not-exist.pem"))
	if err == nil {
		t.Error("LoadCertPool(missing) = nil, want error")
	}
}

func TestLoadCertPool_UnopenableDirectoryFails(t *testing.T) {
	t.Parallel()

	// The parent directory does not exist, so os.OpenRoot fails. This covers the
	// "open CA directory" error branch (distinct from the missing-file case,
	// where the directory exists but the file does not).
	_, err := LoadCertPool(filepath.Join(t.TempDir(), "no-such-dir", "ca.pem"))
	if err == nil {
		t.Error("LoadCertPool(unopenable dir) = nil, want error")
	}
}

func TestLoadCertPool_InvalidContentFails(t *testing.T) {
	t.Parallel()

	path := writeTempFile(t, []byte("this is not a PEM certificate"))
	if _, err := LoadCertPool(path); err == nil {
		t.Error("LoadCertPool(invalid) = nil, want error")
	}
}

func TestLoadCertPool_EmptyFileFails(t *testing.T) {
	t.Parallel()

	// An empty bundle has zero PEM certificates; AppendCertsFromPEM returns
	// false and LoadCertPool must fail rather than silently trusting nothing.
	path := writeTempFile(t, nil)
	if _, err := LoadCertPool(path); err == nil {
		t.Error("LoadCertPool(empty) = nil, want error")
	}
}

func TestLoadCertPool_DirectoryPathFailsOnRead(t *testing.T) {
	t.Parallel()

	// The CA "file" is actually a directory: root.Open succeeds but reading it
	// fails (EISDIR), exercising the io.ReadAll error branch.
	path := filepath.Join(t.TempDir(), "ca.pem")
	if err := os.Mkdir(path, 0o700); err != nil {
		t.Fatalf("Mkdir error: %v", err)
	}
	if _, err := LoadCertPool(path); err == nil {
		t.Error("LoadCertPool(directory) = nil, want error")
	}
}

func TestLoad_NetBoxCAFile(t *testing.T) {
	validPEM := testSelfSignedCertPEM(t)

	t.Run("unset uses empty path", func(t *testing.T) {
		t.Setenv("NETBOX_URL", "http://8.8.8.8")
		t.Setenv("NETBOX_CA_FILE", "")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error: %v", err)
		}
		if cfg.NetBoxCAFile != "" {
			t.Errorf("NetBoxCAFile = %q, want empty", cfg.NetBoxCAFile)
		}
	})

	t.Run("set to a valid PEM path", func(t *testing.T) {
		path := writeTempFile(t, validPEM)
		t.Setenv("NETBOX_URL", "http://8.8.8.8")
		t.Setenv("NETBOX_CA_FILE", path)

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error: %v", err)
		}
		if cfg.NetBoxCAFile != path {
			t.Errorf("NetBoxCAFile = %q, want %q", cfg.NetBoxCAFile, path)
		}
	})
}

// writeTempFile writes data to a temp file and returns its path.
func writeTempFile(t *testing.T, data []byte) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "ca.pem")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	return path
}

// testSelfSignedCertPEM generates a self-signed certificate and returns its
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
