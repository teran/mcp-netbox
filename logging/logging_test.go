package logging

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/config"
)

// secretLog is a struct with a secret:"true" field used to verify the
// redaction hook masks annotated secrets in log output.
type secretLog struct {
	Name  string `json:"name"`
	Token string `json:"token" secret:"true"`
}

func baseConfig() config.Config {
	return config.Config{
		Transport: config.TransportHTTP,
		LogLevel:  "info",
		LogFormat: "text",
	}
}

func TestSetup_HTTPEmptyLevelDefaultsToInfo(t *testing.T) {
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
	if l.Out == io.Discard {
		t.Errorf("HTTP logger output = io.Discard, want enabled")
	}
	if l.GetLevel() != logrus.InfoLevel {
		t.Errorf("level = %v, want info", l.GetLevel())
	}
}

func TestSetup_StdioEmptyLevelDisablesLogging(t *testing.T) {
	t.Parallel()

	cfg := baseConfig()
	cfg.Transport = config.TransportStdio
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
	f, ok := l.Formatter.(*requestIDFormatter)
	if !ok {
		t.Fatalf("formatter = %T, want *requestIDFormatter", l.Formatter)
	}
	if f.inner == nil {
		t.Error("wrapped TextFormatter is nil")
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

func TestSetup_StdioLogDirMissing(t *testing.T) {
	// A stdio transport with a LOG_FILENAME whose parent directory does not
	// exist must fail Setup (os.OpenRoot fails), surfacing the error path
	// through outputForTransport.
	cfg := baseConfig()
	cfg.Transport = config.TransportStdio
	cfg.LogFilename = filepath.Join(t.TempDir(), "does-not-exist", "mcp.log")

	if _, err := Setup(cfg); err == nil {
		t.Fatal("Setup = nil, want error for missing log directory")
	}
}

// captureLogger builds a logger via Setup(cfg) (so the redaction hook is
// installed) and redirects its output to a buffer so assertions can inspect
// the rendered text.
func captureLogger(t *testing.T, format string) (*logrus.Logger, *bytes.Buffer) {
	t.Helper()
	cfg := baseConfig()
	cfg.LogFormat = format
	l, err := Setup(cfg)
	if err != nil {
		t.Fatalf("Setup = %v, want nil", err)
	}
	var buf bytes.Buffer
	l.SetOutput(&buf)
	return l, &buf
}

func TestHook_RedactsStructSecret(t *testing.T) {
	t.Parallel()

	t.Run("text format", func(t *testing.T) {
		l, buf := captureLogger(t, "text")
		l.WithField("obj", secretLog{Name: "srv", Token: "super-secret"}).Info("msg")

		out := buf.String()
		if strings.Contains(out, "super-secret") {
			t.Errorf("text log leaks the secret: %s", out)
		}
		if !strings.Contains(out, "***redacted***") {
			t.Errorf("text log missing mask: %s", out)
		}
		if !strings.Contains(out, "srv") {
			t.Errorf("text log missing non-secret field: %s", out)
		}
	})

	t.Run("json format", func(t *testing.T) {
		l, buf := captureLogger(t, "json")
		l.WithField("obj", secretLog{Name: "srv", Token: "super-secret"}).Info("msg")

		out := buf.String()
		if strings.Contains(out, "super-secret") {
			t.Errorf("json log leaks the secret: %s", out)
		}
		if !strings.Contains(out, "***redacted***") {
			t.Errorf("json log missing mask: %s", out)
		}
	})
}

func TestHook_RedactsApplicationToken(t *testing.T) {
	t.Parallel()

	tok := application.NewToken("tok-abc-123")
	if tok == nil {
		t.Fatal("NewToken returned nil")
	}

	t.Run("text format", func(t *testing.T) {
		l, buf := captureLogger(t, "text")
		l.WithField("token", tok).Info("msg")

		out := buf.String()
		if strings.Contains(out, "tok-abc-123") {
			t.Errorf("text log leaks application.Token value: %s", out)
		}
		if !strings.Contains(out, "***redacted***") {
			t.Errorf("text log missing mask for token: %s", out)
		}
	})

	t.Run("json format", func(t *testing.T) {
		l, buf := captureLogger(t, "json")
		l.WithField("token", tok).Info("msg")

		out := buf.String()
		if strings.Contains(out, "tok-abc-123") {
			t.Errorf("json log leaks application.Token value: %s", out)
		}
		if !strings.Contains(out, "***redacted***") {
			t.Errorf("json log missing mask for token: %s", out)
		}
	})
}

func TestHook_LeavesFlatStringsUnchanged(t *testing.T) {
	t.Parallel()

	l, buf := captureLogger(t, "json")
	l.WithField("plain", "just-a-string").Info("msg")

	out := buf.String()
	if !strings.Contains(out, "just-a-string") {
		t.Errorf("json log modified flat string: %s", out)
	}
}
