package handlers

import (
	"reflect"
	"testing"
)

func TestValidatePositiveID(t *testing.T) {
	t.Parallel()

	if err := validatePositiveID(1); err != nil {
		t.Errorf("validatePositiveID(1) = %v, want nil", err)
	}
	if err := validatePositiveID(0); err == nil {
		t.Error("validatePositiveID(0) = nil, want error")
	}
	if err := validatePositiveID(-5); err == nil {
		t.Error("validatePositiveID(-5) = nil, want error")
	}
}

func TestSanitizeControl(t *testing.T) {
	t.Parallel()

	cases := []struct {
		in   string
		want string
	}{
		{"plain", "plain"},
		{"a\r\nb", "ab"},
		{"\r\n", ""},
		{"a\nb\rc", "abc"},
		{"no controls here", "no controls here"},
	}
	for _, tc := range cases {
		if got := sanitizeControl(tc.in); got != tc.want {
			t.Errorf("sanitizeControl(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestSanitizeOutput(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		in   string
		want string
	}{
		{"plain text preserved", "plain", "plain"},
		{"ANSI CSI sequence stripped", "\x1b[31mred\x1b[0m", "red"},
		{"ANSI SGR with params stripped", "\x1b[1;31mBOLD\x1b[m", "BOLD"},
		{"two-char escape stripped", "\x1b(AUTF", "AUTF"},
		{"C0 control chars removed", "a\x00b\x07c", "abc"},
		{"DEL removed", "a\x7fb", "ab"},
		{"tab removed", "a\tb", "ab"},
		{"newline removed", "a\nb", "ab"},
		{"multibyte UTF-8 preserved", "привет", "привет"},
		{"unicode with controls", "caf\x1b[35mé", "café"},
		{"empty string", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := sanitizeOutput(tc.in); got != tc.want {
				t.Errorf("sanitizeOutput(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestSanitizeMapKeysValues(t *testing.T) {
	t.Parallel()

	t.Run("nil map", func(t *testing.T) {
		if got := sanitizeMapKeysValues(nil); got != nil {
			t.Errorf("sanitizeMapKeysValues(nil) = %v, want nil", got)
		}
	})

	t.Run("strips CRLF from keys and values", func(t *testing.T) {
		in := map[string]string{
			"key\r\n1": "val\r\n1",
			"clean":    "ok",
		}
		got := sanitizeMapKeysValues(in)
		if _, ok := got["key1"]; !ok {
			t.Errorf("result = %v, want sanitized key 'key1'", got)
		}
		if got["key1"] != "val1" {
			t.Errorf("result['key1'] = %q, want val1", got["key1"])
		}
		if got["clean"] != "ok" {
			t.Errorf("result['clean'] = %q, want ok", got["clean"])
		}
	})

	t.Run("does not mutate input", func(t *testing.T) {
		in := map[string]string{"a\r\n": "b\n"}
		_ = sanitizeMapKeysValues(in)
		if _, ok := in["a\r\n"]; !ok {
			t.Error("input map was mutated")
		}
	})
}

func TestSanitizeAllStrings(t *testing.T) {
	t.Parallel()

	type nested struct {
		Note *string
	}
	type payload struct {
		Name     string
		Desc     *string
		Tags     []string
		Count    int
		Enabled  bool
		Custom   map[string]any
		Inner    nested
		InnerPtr *nested
	}

	note := "note\r\n"
	desc := "line1\nline2\r"
	p := payload{
		Name:     "name\r\nx",
		Desc:     &desc,
		Tags:     []string{"t\r\n1", "ok"},
		Count:    5,
		Enabled:  true,
		Custom:   map[string]any{"k\r\n": "v"},
		Inner:    nested{Note: &note},
		InnerPtr: &nested{Note: &note},
	}

	sanitizeAllStrings(reflect.ValueOf(&p))

	if p.Name != "namex" {
		t.Errorf("Name = %q, want 'namex'", p.Name)
	}
	if p.Desc == nil || *p.Desc != "line1line2" {
		t.Errorf("Desc = %v, want 'line1line2'", p.Desc)
	}
	// slices/maps are not sanitized by the struct walker
	if p.Tags[0] != "t\r\n1" {
		t.Errorf("Tags[0] = %q, want untouched 't\\r\\n1'", p.Tags[0])
	}
	if p.Count != 5 || !p.Enabled {
		t.Error("non-string fields must be left intact")
	}
	if p.Inner.Note == nil || *p.Inner.Note != "note" {
		t.Errorf("Inner.Note = %v, want 'note'", p.Inner.Note)
	}
	if p.InnerPtr.Note == nil || *p.InnerPtr.Note != "note" {
		t.Errorf("InnerPtr.Note = %v, want 'note'", p.InnerPtr.Note)
	}
}

func TestSanitizeOutput_EdgeCases(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		in   string
		want string
	}{
		{"incomplete CSI at end of string", "\x1b[31", ""},
		{"lone ESC at end of string", "\x1b", ""},
		{"ESC followed by incomplete two-char", "\x1bA", ""},
		{"multiple escapes", "\x1b[1m\x1b[2mok", "ok"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := sanitizeOutput(tc.in); got != tc.want {
				t.Errorf("sanitizeOutput(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestSanitizeAllStrings_InvalidValue(t *testing.T) {
	t.Parallel()

	// A zero (invalid) reflect.Value must be a no-op rather than panicking.
	var zero reflect.Value
	sanitizeAllStrings(zero)
}
