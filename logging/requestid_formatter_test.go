package logging

import (
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestRequestIDFormatter_WithRequestID(t *testing.T) {
	t.Parallel()

	f := &requestIDFormatter{inner: &logrus.TextFormatter{FullTimestamp: true}}
	entry := &logrus.Entry{
		Data:    logrus.Fields{"request_id": "abc-123", "method": "GET"},
		Level:   logrus.InfoLevel,
		Message: "mcp_request",
	}

	out, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Format = %v, want nil", err)
	}

	line := string(out)
	if !strings.HasPrefix(line, "request_id=abc-123 ") {
		t.Errorf("line should start with request_id=abc-123, got: %q", line)
	}
	if strings.Contains(line, "msg=mcp_request") == false {
		t.Errorf("line missing message field: %q", line)
	}
	// The field must not be duplicated at the end of the line.
	if strings.HasSuffix(strings.TrimSpace(line), "request_id=abc-123") {
		t.Errorf("request_id duplicated at end of line: %q", line)
	}
	if !strings.Contains(line, "method=GET") {
		t.Errorf("line missing method field: %q", line)
	}
}

func TestRequestIDFormatter_WithoutRequestID(t *testing.T) {
	t.Parallel()

	f := &requestIDFormatter{inner: &logrus.TextFormatter{FullTimestamp: true}}
	entry := &logrus.Entry{
		Data:    logrus.Fields{"method": "GET"},
		Level:   logrus.InfoLevel,
		Message: "session connected",
	}

	out, err := f.Format(entry)
	if err != nil {
		t.Fatalf("Format = %v, want nil", err)
	}

	line := string(out)
	if strings.HasPrefix(line, "request_id=") {
		t.Errorf("line should not start with request_id prefix: %q", line)
	}
	if strings.HasPrefix(line, "time=") == false {
		t.Errorf("line should start with the standard time field, got: %q", line)
	}
	if !strings.Contains(line, "session connected") {
		t.Errorf("line missing message: %q", line)
	}
}

func TestRequestIDFormatter_LeavesEntryDataUntouched(t *testing.T) {
	t.Parallel()

	f := &requestIDFormatter{inner: &logrus.TextFormatter{FullTimestamp: true}}
	entry := &logrus.Entry{
		Data:    logrus.Fields{"request_id": "xyz"},
		Level:   logrus.InfoLevel,
		Message: "mcp_request",
	}

	if _, err := f.Format(entry); err != nil {
		t.Fatalf("Format = %v, want nil", err)
	}

	if _, ok := entry.Data["request_id"]; !ok {
		t.Error("Format mutated the caller's entry.Data: request_id removed")
	}
}
