package redact

import (
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

type hookSecret struct {
	Token string `secret:"true"`
}

func TestLogrusHook_Levels(t *testing.T) {
	t.Parallel()

	h := NewLogrusHook()
	levels := h.Levels()
	if len(levels) != len(logrus.AllLevels) {
		t.Fatalf("Levels() len = %d, want %d", len(levels), len(logrus.AllLevels))
	}
}

func TestLogrusHook_FireRedactsNonStrings(t *testing.T) {
	t.Parallel()

	e := logrus.NewEntry(logrus.New())
	e.Data["obj"] = hookSecret{Token: "hook-secret-value"}
	e.Data["flat"] = "plain-string"

	if err := NewLogrusHook().Fire(e); err != nil {
		t.Fatalf("Fire: %v", err)
	}

	out, ok := e.Data["obj"].(hookSecret)
	if !ok {
		t.Fatalf("obj type = %T, want hookSecret", e.Data["obj"])
	}
	if out.Token != Mask {
		t.Errorf("obj.Token = %q, want %q", out.Token, Mask)
	}
	if strings.Contains(e.Data["obj"].(hookSecret).Token, "hook-secret-value") {
		t.Errorf("obj leaks the secret")
	}

	if got, _ := e.Data["flat"].(string); got != "plain-string" {
		t.Errorf("flat = %q, want %q", got, "plain-string")
	}
}

func TestLogrusHook_FireStringValueSkipped(t *testing.T) {
	t.Parallel()

	e := logrus.NewEntry(logrus.New())
	// A struct secret inside a field that is itself a string cannot occur, but
	// a string field must pass through the hook unchanged.
	e.Data["key"] = "some-string"

	if err := NewLogrusHook().Fire(e); err != nil {
		t.Fatalf("Fire: %v", err)
	}
	if got, _ := e.Data["key"].(string); got != "some-string" {
		t.Errorf("key = %q, want %q", got, "some-string")
	}
}
