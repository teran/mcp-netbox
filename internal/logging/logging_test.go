package logging

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/teran/mcp-netbox/config"
)

func baseConfig() config.Config {
	return config.Config{
		Transport: config.TransportHTTP,
		LogLevel:  "info",
		LogFormat: "text",
	}
}

func TestSetup_EmptyLevelDisablesLogging(t *testing.T) {
	t.Parallel()

	cfg := baseConfig()
	cfg.LogLevel = ""

	l, err := Setup(cfg)
	if err != nil {
		t.Fatalf("Setup = %v, want nil", err)
	}
	if l == nil {
		t.Fatal("Setup returned nil logger")
	}
	if l.Out != io.Discard {
		t.Errorf("logger output = %T, want io.Discard", l.Out)
	}
}

func TestSetup_LevelAndTextFormat(t *testing.T) {
	t.Parallel()

	l, err := Setup(baseConfig())
	if err != nil {
		t.Fatalf("Setup = %v, want nil", err)
	}
	if l.GetLevel() != logrus.InfoLevel {
		t.Errorf("level = %v, want info", l.GetLevel())
	}
	if _, ok := l.Formatter.(*logrus.TextFormatter); !ok {
		t.Errorf("formatter = %T, want TextFormatter", l.Formatter)
	}
}

func TestSetup_JSONFormat(t *testing.T) {
	t.Parallel()

	cfg := baseConfig()
	cfg.LogFormat = "json"

	l, err := Setup(cfg)
	if err != nil {
		t.Fatalf("Setup = %v, want nil", err)
	}
	if _, ok := l.Formatter.(*logrus.JSONFormatter); !ok {
		t.Errorf("formatter = %T, want JSONFormatter", l.Formatter)
	}
}

func TestSetup_InvalidLevel(t *testing.T) {
	t.Parallel()

	cfg := baseConfig()
	cfg.LogLevel = "bogus"

	if _, err := Setup(cfg); err == nil {
		t.Fatal("Setup = nil, want error for invalid level")
	}
}

func TestSetup_InvalidFormat(t *testing.T) {
	t.Parallel()

	cfg := baseConfig()
	cfg.LogFormat = "xml"

	if _, err := Setup(cfg); err == nil {
		t.Fatal("Setup = nil, want error for invalid format")
	}
}

func TestSetup_StdioWritesToFile(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "mcp-netbox.log")

	cfg := baseConfig()
	cfg.Transport = config.TransportStdio
	cfg.LogFilename = filename

	l, err := Setup(cfg)
	if err != nil {
		t.Fatalf("Setup = %v, want nil", err)
	}
	if l == nil {
		t.Fatal("Setup returned nil logger")
	}

	// Logging a line must land in the file.
	l.Info("hello from test")

	fi, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("stat log file: %v", err)
	}
	if fi.Size() == 0 {
		t.Error("log file is empty after logging")
	}
	if perm := fi.Mode().Perm(); perm != 0o600 {
		t.Errorf("log file perm = %o, want 600", perm)
	}
}

func TestSetup_StdioDefaultFilename(t *testing.T) {
	// LOG_FILENAME empty in stdio mode must fall back to the default path.
	cfg := baseConfig()
	cfg.Transport = config.TransportStdio
	cfg.LogFilename = ""

	l, err := Setup(cfg)
	if err != nil {
		t.Fatalf("Setup = %v, want nil", err)
	}
	if l == nil {
		t.Fatal("Setup returned nil logger")
	}
}

func TestSetup_HTTPOutputsToStdout(t *testing.T) {
	t.Parallel()

	l, err := Setup(baseConfig())
	if err != nil {
		t.Fatalf("Setup = %v, want nil", err)
	}
	if l.Out != os.Stdout {
		t.Errorf("logger output = %T, want os.Stdout", l.Out)
	}
}
