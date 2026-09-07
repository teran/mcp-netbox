package netbox

import (
	"io"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestSetLoggerAndGetLogger(t *testing.T) {
	t.Cleanup(func() {
		SetLogger(nil)
	})

	l := logrus.New()
	l.SetOutput(io.Discard)
	l.SetLevel(logrus.DebugLevel)

	SetLogger(l)

	if got := getLogger(); got != l {
		t.Errorf("getLogger() = %p, want %p", got, l)
	}
}

func TestGetLoggerDefaultDiscards(t *testing.T) {
	t.Cleanup(func() {
		SetLogger(nil)
	})

	SetLogger(nil)

	l := getLogger()
	if l == nil {
		t.Fatal("getLogger() returned nil")
	}
	if l.Out != io.Discard {
		t.Errorf("default logger output = %T, want io.Discard", l.Out)
	}
}
