package redact

import (
	"github.com/sirupsen/logrus"
)

// NewLogrusHook returns a logrus.Hook that redacts every non-string value in
// the entry's Data map before the entry is formatted and written. This ensures
// struct-typed fields (for example the application.Token wrapper or any value
// with secret:"true" fields) cannot leak a secret into log output.
//
// Scalar strings are passed through unchanged: a bare string has no struct tag
// for Redact to act on, and rewriting it would corrupt ordinary log fields.
func NewLogrusHook() logrus.Hook {
	return &secretHook{}
}

type secretHook struct{}

// Levels returns every log level so the hook always runs.
func (h *secretHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire redacts every non-string value in e.Data in place.
func (h *secretHook) Fire(e *logrus.Entry) error {
	for k, v := range e.Data {
		if _, ok := v.(string); ok {
			continue
		}
		e.Data[k] = Redact(v)
	}
	return nil
}
