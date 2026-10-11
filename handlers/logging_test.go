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

	logger := logrus.New()
	logger.SetOutput(io.Discard)
	logger.SetLevel(logrus.DebugLevel)

	SetLogger(logger)

	if got := getLogger(); got != logger {
		t.Errorf("getLogger() = %p, want %p", got, logger)
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
