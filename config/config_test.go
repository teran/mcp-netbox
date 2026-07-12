package config

import (
	"testing"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func setenv(t *testing.T, key, value string) {
	t.Helper()
	t.Setenv(key, value)
}

func unsetenv(t *testing.T, key string) {
	t.Helper()
	t.Setenv(key, "") // no-op since t.Setenv cleans up; just keep for API compat
}

func TestLoad(t *testing.T) {
	t.Run("all env vars set correctly", func(t *testing.T) {
		setenv(t, "NETBOX_URL", "http://netbox:8000")
		defer unsetenv(t, "NETBOX_URL")
		setenv(t, "LISTEN_ADDR", ":9090")
		defer unsetenv(t, "LISTEN_ADDR")
		setenv(t, "PROMETHEUS_METRICS_ADDR", ":9091")
		defer unsetenv(t, "PROMETHEUS_METRICS_ADDR")
		setenv(t, "RATE_LIMIT_GLOBAL", "200")
		defer unsetenv(t, "RATE_LIMIT_GLOBAL")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}
		if cfg.NetBoxURL != "http://netbox:8000" {
			t.Errorf("NetBoxURL = %q, want %q", cfg.NetBoxURL, "http://netbox:8000")
		}
		if cfg.ListenAddr != ":9090" {
			t.Errorf("ListenAddr = %q, want %q", cfg.ListenAddr, ":9090")
		}
		if cfg.PrometheusMetricsAddr != ":9091" {
			t.Errorf("PrometheusMetricsAddr = %q, want %q", cfg.PrometheusMetricsAddr, ":9091")
		}
		if cfg.RateLimitGlobal != 200 {
			t.Errorf("RateLimitGlobal = %d, want %d", cfg.RateLimitGlobal, 200)
		}
	})

	t.Run("defaults when env vars are not set", func(t *testing.T) {
		setenv(t, "NETBOX_URL", "https://netbox.example.com")
		defer unsetenv(t, "NETBOX_URL")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}
		if cfg.ListenAddr != ":8080" {
			t.Errorf("ListenAddr = %q, want %q", cfg.ListenAddr, ":8080")
		}
		if cfg.RateLimitGlobal != 100 {
			t.Errorf("RateLimitGlobal = %d, want %d", cfg.RateLimitGlobal, 100)
		}
		if cfg.RateLimitPerClient != 10 {
			t.Errorf("RateLimitPerClient = %d, want %d", cfg.RateLimitPerClient, 10)
		}
		if cfg.WriteTimeout != 300*time.Second {
			t.Errorf("WriteTimeout = %v, want %v", cfg.WriteTimeout, 300*time.Second)
		}
	})
}

func TestLoad_Errors(t *testing.T) {
	t.Run("NETBOX_URL is required", func(t *testing.T) {
		_, err := Load()
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("NETBOX_URL invalid URL", func(t *testing.T) {
		setenv(t, "NETBOX_URL", "://invalid")
		defer unsetenv(t, "NETBOX_URL")
		_, err := Load()
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("NETBOX_URL must have http/https scheme", func(t *testing.T) {
		setenv(t, "NETBOX_URL", "ftp://netbox:8000")
		defer unsetenv(t, "NETBOX_URL")
		_, err := Load()
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("NETBOX_URL must have host", func(t *testing.T) {
		setenv(t, "NETBOX_URL", "http://")
		defer unsetenv(t, "NETBOX_URL")
		_, err := Load()
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("RATE_LIMIT_GLOBAL must be >= 1", func(t *testing.T) {
		setenv(t, "NETBOX_URL", "http://netbox:8000")
		defer unsetenv(t, "NETBOX_URL")
		setenv(t, "RATE_LIMIT_GLOBAL", "0")
		defer unsetenv(t, "RATE_LIMIT_GLOBAL")
		_, err := Load()
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("RATE_LIMIT_PER_CLIENT must be >= 1", func(t *testing.T) {
		setenv(t, "NETBOX_URL", "http://netbox:8000")
		defer unsetenv(t, "NETBOX_URL")
		setenv(t, "RATE_LIMIT_PER_CLIENT", "0")
		defer unsetenv(t, "RATE_LIMIT_PER_CLIENT")
		_, err := Load()
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestValidatePositiveInt(t *testing.T) {
	t.Run("non-int value", func(t *testing.T) {
		err := validatePositiveInt("not-an-int")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("validation rule err", func(t *testing.T) {
		err := validation.By(validatePositiveInt).Validate("string")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestValidateURLScheme(t *testing.T) {
	t.Run("parse error", func(t *testing.T) {
		err := validateURLScheme("://")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestValidateURLHost(t *testing.T) {
	t.Run("parse error", func(t *testing.T) {
		err := validateURLHost("://")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}
