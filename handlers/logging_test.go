package handlers

import (
	"io"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestSetLoggerAndGetLogger(t *testing.T) {
	t.Cleanup(func() {
		// Reset so other tests keep using the discard logger.
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

	SetLogger(nil) // ensure reset

	l := getLogger()
	if l == nil {
		t.Fatal("getLogger() returned nil")
	}
	if l.Out != io.Discard {
		t.Errorf("default logger output = %T, want io.Discard", l.Out)
	}
}

func TestSetLoggerNilResetsToDiscard(t *testing.T) {
	t.Cleanup(func() {
		SetLogger(nil)
	})

	SetLogger(logrus.New())
	SetLogger(nil)

	if l := getLogger(); l.Out != io.Discard {
		t.Errorf("after SetLogger(nil) output = %T, want io.Discard", l.Out)
	}
}
