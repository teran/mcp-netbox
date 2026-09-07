package handlers

import (
	"io"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

// defaultLogger is the package-level logger used by the HTTP middleware and
// handlers. It is set once at startup from the configured logger (see
// SetLogger). Until then, and whenever LOG_LEVEL is disabled, it discards all
// output.
var defaultLogger atomic.Pointer[logrus.Logger]

// SetLogger configures the package-level logger used by the HTTP middleware.
// A nil logger is treated as a discard logger.
func SetLogger(l *logrus.Logger) {
	if l == nil {
		l = discardLogger()
	}
	defaultLogger.Store(l)
}

// getLogger returns the configured logger, or a discard logger when none has
// been set (e.g. in tests).
func getLogger() *logrus.Logger {
	if l := defaultLogger.Load(); l != nil {
		return l
	}
	return discardLogger()
}

// discardLogger returns a logger that writes nowhere.
func discardLogger() *logrus.Logger {
	l := logrus.New()
	l.SetOutput(io.Discard)
	return l
}
