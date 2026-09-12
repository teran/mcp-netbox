package logging

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

// TestNewSlogLogger forwards SDK log events to the underlying logrus logger.
func TestNewSlogLogger_ForwardsToLogrus(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logrus.New()
	l.SetOutput(&buf)
	l.SetLevel(logrus.InfoLevel)

	sl := NewSlogLogger(l)
	sl.Info("sdk event", "tool", "get_sites", "session", "sess-1")

	out := buf.String()
	if !strings.Contains(out, "sdk event") {
		t.Errorf("logrus output = %q, want it to contain the message", out)
	}
	if !strings.Contains(out, "get_sites") {
		t.Errorf("logrus output = %q, want it to contain the tool attribute", out)
	}
}

// TestNewSlogLogger_LevelGating verifies that slog records below the logrus
// level are dropped.
func TestNewSlogLogger_LevelGating(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logrus.New()
	l.SetOutput(&buf)
	l.SetLevel(logrus.InfoLevel)

	sl := NewSlogLogger(l)
	sl.Debug("should be dropped")

	if buf.Len() != 0 {
		t.Errorf("logrus output = %q, want empty (debug dropped at info level)", buf.String())
	}
}

// TestNewSlogLogger_ErrorLevel verifies errors map to the logrus error level
// and are emitted.
func TestNewSlogLogger_ErrorLevel(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logrus.New()
	l.SetOutput(&buf)
	l.SetLevel(logrus.InfoLevel)

	sl := NewSlogLogger(l)
	sl.Error("boom")

	out := buf.String()
	if !strings.Contains(out, "boom") {
		t.Errorf("logrus output = %q, want it to contain the message", out)
	}
}

// TestSlogHandler_Enabled verifies Enabled honours the configured logrus level.
func TestSlogHandler_Enabled(t *testing.T) {
	t.Parallel()

	l := logrus.New()
	l.SetLevel(logrus.WarnLevel)

	h := newSlogHandler(l)
	if h.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("Enabled(Info) = true, want false at warn level")
	}
	if !h.Enabled(context.Background(), slog.LevelWarn) {
		t.Error("Enabled(Warn) = false, want true at warn level")
	}
}

// TestSlogHandler_WithAttrs verifies attrs attached via WithAttrs are emitted.
func TestSlogHandler_WithAttrs(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logrus.New()
	l.SetOutput(&buf)
	l.SetLevel(logrus.InfoLevel)

	sl := NewSlogLogger(l).With("env", "prod")
	sl.Info("with attrs")

	out := buf.String()
	if !strings.Contains(out, "prod") {
		t.Errorf("logrus output = %q, want it to contain the attr value", out)
	}
}

// TestSlogHandler_WithGroup verifies WithGroup is accepted and produces output.
func TestSlogHandler_WithGroup(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logrus.New()
	l.SetOutput(&buf)
	l.SetLevel(logrus.InfoLevel)

	sl := NewSlogLogger(l).WithGroup("mcp")
	sl.Info("grouped")

	if !strings.Contains(buf.String(), "grouped") {
		t.Errorf("logrus output = %q, want it to contain the message", buf.String())
	}
}

// TestToLogrusLevel verifies the level mapping boundaries.
func TestToLogrusLevel(t *testing.T) {
	t.Parallel()

	cases := []struct {
		slog slog.Level
		want logrus.Level
	}{
		{slog.LevelError, logrus.ErrorLevel},
		{slog.LevelError - 1, logrus.WarnLevel},
		{slog.LevelWarn, logrus.WarnLevel},
		{slog.LevelWarn - 1, logrus.InfoLevel},
		{slog.LevelInfo, logrus.InfoLevel},
		{slog.LevelDebug, logrus.DebugLevel},
		{slog.LevelDebug - 2, logrus.DebugLevel},
	}
	for _, c := range cases {
		if got := toLogrusLevel(c.slog); got != c.want {
			t.Errorf("toLogrusLevel(%v) = %v, want %v", c.slog, got, c.want)
		}
	}
}

// TestAttrKey verifies the empty-key fallback.
func TestAttrKey(t *testing.T) {
	t.Parallel()

	if got := attrKey(""); got != "message" {
		t.Errorf("attrKey(\"\") = %q, want message", got)
	}
	if got := attrKey("tool"); got != "tool" {
		t.Errorf("attrKey(tool) = %q, want tool", got)
	}
}
