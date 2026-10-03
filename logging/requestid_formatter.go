package logging

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// requestIDKey is the logrus field key that carries a request id.
const requestIDKey = "request_id"

// requestIDFormatter wraps a *logrus.TextFormatter and hoists the request_id
// field to the very beginning of the rendered text line, so all log lines
// belonging to a single request can be collected with one
// `grep 'request_id=<id>'` over the log stream.
//
// Lines that carry no request_id (e.g. the SDK session lines) are rendered
// exactly as before by the wrapped TextFormatter.
type requestIDFormatter struct {
	inner *logrus.TextFormatter
}

// Format renders entry, prepending `request_id=<value> ` when the field is
// present. The field is removed from a copy of entry.Data before delegating to
// the wrapped TextFormatter so it is not emitted a second time at the end of
// the line.
func (f *requestIDFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	raw, ok := entry.Data[requestIDKey]
	if !ok {
		return f.inner.Format(entry)
	}

	copyData := make(logrus.Fields, len(entry.Data)-1)
	for k, v := range entry.Data {
		if k != requestIDKey {
			copyData[k] = v
		}
	}

	clone := *entry
	clone.Data = copyData

	body, err := f.inner.Format(&clone)
	if err != nil {
		return nil, err
	}

	prefix := requestIDKey + "=" + fmt.Sprintf("%v", raw) + " "
	out := make([]byte, 0, len(prefix)+len(body))
	out = append(out, prefix...)
	out = append(out, body...)
	return out, nil
}
