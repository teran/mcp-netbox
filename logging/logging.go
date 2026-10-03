// Package logging configures a logrus logger from server configuration.
//
// The behaviour is driven by environment variables:
//
//   - LOG_LEVEL: logrus level (trace, debug, info, warn, error, fatal, panic).
//     Enablement is launch-mode dependent (L2): in HTTP mode an empty LOG_LEVEL
//     defaults to `info` (always enabled); in STDIO mode an empty LOG_LEVEL
//     disables logging entirely (output goes to io.Discard).
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
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/teran/mcp-netbox/config"
	"github.com/teran/mcp-netbox/redact"
)

// DefaultLogFilename is used for the STDIO log file when LOG_FILENAME is unset.
const DefaultLogFilename = "/tmp/mcp-netbox.log"

// Setup builds a *logrus.Logger according to cfg.
//
// The returned logger owns an open file handle in STDIO mode that lives for
// the process lifetime. Callers should not call Close on it.
func Setup(cfg config.Config) (*logrus.Logger, error) {
	l := logrus.New()

	// L2: enablement is launch-mode dependent.
	//   - HTTP: logging is ALWAYS enabled; an empty LOG_LEVEL defaults to info.
	//   - STDIO: logging is enabled only when LOG_LEVEL is set; unset => disabled.
	levelStr := strings.TrimSpace(cfg.LogLevel)
	if cfg.Transport != config.TransportStdio && levelStr == "" {
		levelStr = logrus.InfoLevel.String()
	}
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
		l.SetFormatter(&requestIDFormatter{inner: &logrus.TextFormatter{FullTimestamp: true}})
	case "json":
		l.SetFormatter(&logrus.JSONFormatter{})
	default:
		return nil, fmt.Errorf("invalid LOG_FORMAT %q (want text or json)", cfg.LogFormat)
	}

	// S02: redact any struct-typed (secret:"true") field before it is written,
	// so annotated secrets never leak into log lines.
	l.AddHook(redact.NewLogrusHook())

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
		return openLogFile(filename)
	default: // "" and http
		return os.Stdout, nil
	}
}

// openLogFile opens (creating if needed) the log file for STDIO transport and
// returns a writer that must not be closed by the caller.
//
// The path is canonicalized with filepath.Clean/Abs so a relative or sloppy
// LOG_FILENAME cannot reference a different file than the operator intended.
// Access is then scoped to the file's parent directory via os.Root (G304),
// which never follows symlinks, and the opened handle is verified to be a
// regular file so log output can never be redirected to a sensitive target.
func openLogFile(filename string) (*os.File, error) {
	path, err := filepath.Abs(filepath.Clean(filename))
	if err != nil {
		return nil, fmt.Errorf("resolve log path %q: %w", filename, err)
	}

	dir := filepath.Dir(path)
	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, fmt.Errorf("open log directory %q: %w", dir, err)
	}
	defer func() { _ = root.Close() }()

	name := filepath.Base(path)
	f, err := root.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open log file %q: %w", path, err)
	}
	if err := root.Chmod(name, 0o600); err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("chmod log file %q: %w", path, err)
	}

	fi, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("stat log file %q: %w", path, err)
	}
	if !fi.Mode().IsRegular() {
		_ = f.Close()
		return nil, fmt.Errorf("log path %q is not a regular file", path)
	}

	return f, nil
}
