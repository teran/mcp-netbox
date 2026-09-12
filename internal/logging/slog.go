package logging

import (
	"context"
	"log/slog"

	"github.com/sirupsen/logrus"
)

// slogHandler adapts a logrus.Logger to the slog.Handler interface, so that
// the MCP go-sdk's internal slog logger (L7/G10) can be wired into the
// server's logrus logger and its events become visible in the server logs.
type slogHandler struct {
	logger *logrus.Logger
	attrs  []slog.Attr
}

// newSlogHandler builds a slog.Handler that forwards records to logger.
func newSlogHandler(logger *logrus.Logger) slog.Handler {
	return &slogHandler{logger: logger}
}

// Enabled reports whether the logrus logger would emit the given level.
func (h *slogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return toLogrusLevel(level) <= h.logger.GetLevel()
}

// Handle forwards a single slog record to logrus.
func (h *slogHandler) Handle(_ context.Context, r slog.Record) error {
	entry := h.logger.WithFields(attrsToFields(h.attrs))
	r.Attrs(func(a slog.Attr) bool {
		entry = entry.WithField(attrKey(a.Key), a.Value.Any())
		return true
	})
	entry.Log(toLogrusLevel(r.Level), r.Message)
	return nil
}

// WithAttrs returns a handler whose Handle will be called with the given attrs
// prepended to the record.
func (h *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	combined := make([]slog.Attr, 0, len(h.attrs)+len(attrs))
	combined = append(combined, h.attrs...)
	combined = append(combined, attrs...)
	return &slogHandler{logger: h.logger, attrs: combined}
}

// WithGroup returns a handler. Attribute groups are flattened into dotted keys
// by slog before reaching Handle, so no special handling is required here.
func (h *slogHandler) WithGroup(string) slog.Handler {
	return h
}

// NewSlogLogger returns an *slog.Logger backed by the given logrus logger.
func NewSlogLogger(l *logrus.Logger) *slog.Logger {
	return slog.New(newSlogHandler(l))
}

func toLogrusLevel(level slog.Level) logrus.Level {
	switch {
	case level >= slog.LevelError:
		return logrus.ErrorLevel
	case level >= slog.LevelWarn:
		return logrus.WarnLevel
	case level >= slog.LevelInfo:
		return logrus.InfoLevel
	default:
		return logrus.DebugLevel
	}
}

func attrsToFields(attrs []slog.Attr) logrus.Fields {
	fields := make(logrus.Fields, len(attrs))
	for _, a := range attrs {
		fields[attrKey(a.Key)] = a.Value.Any()
	}
	return fields
}

func attrKey(key string) string {
	if key == "" {
		return "message"
	}
	return key
}
