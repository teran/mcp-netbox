// Package config provides configuration loading from environment variables
// using kelseyhightower/envconfig and validation via ozzo-validation.
package config

import (
	"fmt"
	"net/url"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	NetBoxURL             string        `envconfig:"NETBOX_URL" required:"true"`
	ListenAddr            string        `envconfig:"LISTEN_ADDR" default:":8080"`
	PrometheusMetricsAddr string        `envconfig:"PROMETHEUS_METRICS_ADDR" default:":8081"`
	RateLimitGlobal       int           `envconfig:"RATE_LIMIT_GLOBAL" default:"100"`
	RateLimitPerClient    int           `envconfig:"RATE_LIMIT_PER_CLIENT" default:"10"`
	TrustedProxy          string        `envconfig:"TRUSTED_PROXY" default:""`
	WriteTimeout          time.Duration `envconfig:"WRITE_TIMEOUT" default:"60s"`
}

func (c Config) validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.NetBoxURL,
			validation.Required,
			validation.By(validateURLScheme),
			validation.By(validateURLHost),
		),
		validation.Field(&c.RateLimitGlobal, validation.By(validatePositiveInt)),
		validation.Field(&c.RateLimitPerClient, validation.By(validatePositiveInt)),
		validation.Field(&c.WriteTimeout, validation.Min(time.Second)),
	)
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
