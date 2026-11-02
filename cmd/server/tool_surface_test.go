package main

import (
	"fmt"
	"io"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/teran/mcp-netbox/handlers"
)

// capturingHook is a minimal logrus.Hook that records the raw *logrus.Entry
// values so tests can assert on the structured fields (tool/title/
// always_present) rather than parsing rendered text output.
type capturingHook struct {
	mu      sync.Mutex
	entries []*logrus.Entry
}

func (h *capturingHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *capturingHook) Fire(e *logrus.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = append(h.entries, e)

	return nil
}

func (h *capturingHook) snapshot() []*logrus.Entry {
	h.mu.Lock()
	defer h.mu.Unlock()

	out := make([]*logrus.Entry, len(h.entries))
	copy(out, h.entries)

	return out
}

// TestLogToolSurface_EmitsRegisteredToolLine verifies logToolSurface emits an
// info-level "registered tool" line for a known tool (get_sites) whose
// structured fields carry the tool name, a non-empty title and always_present
// == true.
func TestLogToolSurface_EmitsRegisteredToolLine(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	logger.SetLevel(logrus.InfoLevel)
	hook := &capturingHook{}
	logger.AddHook(hook)

	logToolSurface(logger)

	var found bool
	for _, entry := range hook.snapshot() {
		if entry.Level != logrus.InfoLevel || entry.Message != "registered tool" {
			continue
		}
		if entry.Data["tool"] != "get_sites" {
			continue
		}

		found = true

		if title, ok := entry.Data["title"].(string); !ok || title == "" {
			t.Errorf("registered tool line for get_sites: title field = %v, want non-empty string", entry.Data["title"])
		}
		if got := fmt.Sprint(entry.Data["always_present"]); got != "true" {
			t.Errorf("registered tool line for get_sites: always_present = %v, want true", entry.Data["always_present"])
		}
	}

	if !found {
		t.Error("no info-level 'registered tool' log line found for get_sites")
	}
}

// TestLogToolSurface_EmitsOneLinePerTool verifies logToolSurface emits exactly
// one "registered tool" line per entry in handlers.ToolSurface().
func TestLogToolSurface_EmitsOneLinePerTool(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	logger.SetLevel(logrus.InfoLevel)
	hook := &capturingHook{}
	logger.AddHook(hook)

	logToolSurface(logger)

	count := 0
	for _, entry := range hook.snapshot() {
		if entry.Level == logrus.InfoLevel && entry.Message == "registered tool" {
			count++
		}
	}

	if want := len(handlers.ToolSurface()); count != want {
		t.Errorf("number of 'registered tool' lines = %d, want %d (== ToolSurface count)", count, want)
	}
}
