// Package logging configures a logrus logger from server configuration.
//
// The behaviour is driven by environment variables:
//
//   - LOG_LEVEL: logrus level (trace, debug, info, warn, error, fatal, panic).
//     When empty, logging is disabled entirely (output goes to io.Discard).
//   - LOG_FORMAT: "text" (default) or "json".
//   - LOG_FILENAME: path of the log file used in STDIO mode. Ignored in HTTP
//     mode. Defaults to /tmp/mcp-netbox.log. The file is created with mode
//     0600.
//
// The log channel is selected by transport (L1):
//
//   - HTTP transport logs to stdout (12-factor style).
//   - STDIO transport logs to a file, because stdout carries the MCP protocol
//     itself and must not be polluted with log lines.
package logging

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/teran/mcp-netbox/config"
)

// DefaultLogFilename is used for the STDIO log file when LOG_FILENAME is unset.
const DefaultLogFilename = "/tmp/mcp-netbox.log"

// Setup builds a *logrus.Logger according to cfg.
//
// The returned logger owns an open file handle in STDIO mode that lives for
// the process lifetime. Callers should not call Close on it.
func Setup(cfg config.Config) (*logrus.Logger, error) {
	l := logrus.New()

	// L2: an empty LOG_LEVEL disables logging entirely.
	levelStr := strings.TrimSpace(cfg.LogLevel)
	if levelStr == "" {
		l.SetOutput(io.Discard)
		return l, nil
	}

	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		return nil, fmt.Errorf("invalid LOG_LEVEL %q: %w", cfg.LogLevel, err)
	}
	l.SetLevel(level)

	// L4: log format, text by default.
	switch strings.ToLower(strings.TrimSpace(cfg.LogFormat)) {
	case "", "text":
		l.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	case "json":
		l.SetFormatter(&logrus.JSONFormatter{})
	default:
		return nil, fmt.Errorf("invalid LOG_FORMAT %q (want text or json)", cfg.LogFormat)
	}

	// L1: output channel by transport.
	out, err := outputForTransport(cfg)
	if err != nil {
		return nil, err
	}
	l.SetOutput(out)

	return l, nil
}

// outputForTransport selects the log output: stdout for HTTP transport, a
// chmod-0600 file for STDIO transport (L3: LOG_FILENAME overrides the default).
func outputForTransport(cfg config.Config) (io.Writer, error) {
	switch cfg.Transport {
	case config.TransportStdio:
		filename := strings.TrimSpace(cfg.LogFilename)
		if filename == "" {
			filename = DefaultLogFilename
		}
		//nolint:gosec // LOG_FILENAME is operator-controlled config, not untrusted input.
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
		if err != nil {
			return nil, fmt.Errorf("open log file %q: %w", filename, err)
		}
		if err := os.Chmod(filename, 0o600); err != nil {
			return nil, fmt.Errorf("chmod log file %q: %w", filename, err)
		}
		return f, nil
	default: // "" and http
		return os.Stdout, nil
	}
}
