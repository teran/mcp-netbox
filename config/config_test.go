package config

import (
	"testing"
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
		{"hostname", "http://netbox.example.com", false},
		{"hostname with path", "https://netbox.internal/api", false},
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
			NetBoxURL:          "http://netbox:8000",
			RateLimitGlobal:    100,
			RateLimitPerClient: 10,
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
