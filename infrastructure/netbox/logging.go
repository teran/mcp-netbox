package netbox

import (
	"io"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

// defaultLogger is the package-level logger used by the NetBox client. It is
// set once at startup via SetLogger; until then it discards output.
var defaultLogger atomic.Pointer[logrus.Logger]

// SetLogger configures the package-level logger used by the NetBox client.
func SetLogger(l *logrus.Logger) {
	if l == nil {
		l = discardLogger()
	}
	defaultLogger.Store(l)
}

// getLogger returns the configured logger, or a discard logger when none has
// been set.
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
