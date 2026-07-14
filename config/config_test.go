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
