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

//nolint:funlen // Single flat assertion over one struct state; splitting is not sensible.
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
	sample := payload{
		Name:     "name\r\nx",
		Desc:     &desc,
		Tags:     []string{"t\r\n1", "ok"},
		Count:    5,
		Enabled:  true,
		Custom:   map[string]any{"k\r\n": "v\r\nx", "nested": []string{"a\r\nb"}},
		Inner:    nested{Note: &note},
		InnerPtr: &nested{Note: &note},
	}

	sanitizeAllStrings(reflect.ValueOf(&sample))

	if sample.Name != "namex" {
		t.Errorf("Name = %q, want 'namex'", sample.Name)
	}
	if sample.Desc == nil || *sample.Desc != "line1line2" {
		t.Errorf("Desc = %v, want 'line1line2'", sample.Desc)
	}
	// slices are sanitized recursively by the struct walker
	if sample.Tags[0] != "t1" {
		t.Errorf("Tags[0] = %q, want 't1'", sample.Tags[0])
	}
	if sample.Tags[1] != "ok" {
		t.Errorf("Tags[1] = %q, want 'ok'", sample.Tags[1])
	}
	// map keys AND values are sanitized (including nested string values)
	if _, ok := sample.Custom["k"]; !ok {
		t.Errorf("Custom = %v, want sanitized key 'k'", sample.Custom)
	}
	if sample.Custom["k"] != "vx" {
		t.Errorf("Custom['k'] = %q, want 'vx'", sample.Custom["k"])
	}
	if nested, ok := sample.Custom["nested"].([]string); !ok || nested[0] != "ab" {
		t.Errorf("Custom['nested'] = %v, want ['ab']", sample.Custom["nested"])
	}
	if sample.Count != 5 || !sample.Enabled {
		t.Error("non-string fields must be left intact")
	}
	if sample.Inner.Note == nil || *sample.Inner.Note != "note" {
		t.Errorf("Inner.Note = %v, want 'note'", sample.Inner.Note)
	}
	if sample.InnerPtr.Note == nil || *sample.InnerPtr.Note != "note" {
		t.Errorf("InnerPtr.Note = %v, want 'note'", sample.InnerPtr.Note)
	}
}

//nolint:funlen // Table of subtests covering the nested-collection branches; splitting is not sensible.
func TestSanitizeAllStrings_NestedCollections(t *testing.T) {
	t.Parallel()

	t.Run("slice of strings sanitized in place", func(t *testing.T) {
		sl := []string{"a\r\nb", "ok", "\n"}
		sanitizeAllStrings(reflect.ValueOf(&sl))

		want := []string{"ab", "ok", ""}
		if !reflect.DeepEqual(sl, want) {
			t.Errorf("slice = %#v, want %#v", sl, want)
		}
	})

	t.Run("nil slice is a no-op", func(t *testing.T) {
		var sl []string
		sanitizeAllStrings(reflect.ValueOf(&sl))
		if sl != nil {
			t.Errorf("nil slice became %#v, want nil", sl)
		}
	})

	t.Run("array of strings sanitized", func(t *testing.T) {
		arr := [2]string{"x\r\ny", "clean"}
		sanitizeAllStrings(reflect.ValueOf(&arr))

		want := [2]string{"xy", "clean"}
		if arr != want {
			t.Errorf("array = %#v, want %#v", arr, want)
		}
	})

	t.Run("map[string]string keys and values sanitized", func(t *testing.T) {
		mapping := map[string]string{"k\r\n1": "v\r\n1", "clean": "ok"}
		sanitizeAllStrings(reflect.ValueOf(&mapping))

		if got, ok := mapping["k1"]; !ok || got != "v1" {
			t.Errorf("map = %#v, want sanitized 'k1' -> 'v1'", mapping)
		}
		if got := mapping["clean"]; got != "ok" {
			t.Errorf("map['clean'] = %q, want 'ok'", got)
		}
		if _, ok := mapping["k\r\n1"]; ok {
			t.Errorf("map = %#v, stale unsanitized key 'k\\r\\n1' still present", mapping)
		}
	})

	t.Run("map[string]any nested values sanitized", func(t *testing.T) {
		type inner struct {
			Note *string
		}
		note := "n\r\no"
		mapping := map[string]any{
			"s":      "str\r\nval",
			"slice":  []string{"a\r\nb", "c\n"},
			"struct": inner{Note: &note},
			"int":    42,
		}
		sanitizeAllStrings(reflect.ValueOf(&mapping))

		if mapping["s"] != "strval" {
			t.Errorf("map['s'] = %q, want 'strval'", mapping["s"])
		}
		if sl, ok := mapping["slice"].([]string); !ok || !reflect.DeepEqual(sl, []string{"ab", "c"}) {
			t.Errorf("map['slice'] = %#v, want ['ab' 'c']", mapping["slice"])
		}
		if st, ok := mapping["struct"].(inner); !ok || st.Note == nil || *st.Note != "no" {
			t.Errorf("map['struct'] = %#v, want Note 'no'", mapping["struct"])
		}
		if mapping["int"] != 42 {
			t.Errorf("map['int'] = %v, want 42", mapping["int"])
		}
	})

	t.Run("nil map is a no-op", func(t *testing.T) {
		var mapping map[string]string
		sanitizeAllStrings(reflect.ValueOf(&mapping))
		if mapping != nil {
			t.Errorf("nil map became %#v, want nil", mapping)
		}
	})

	t.Run("interface field holding a string sanitized", func(t *testing.T) {
		var v any = "iface\r\nval"
		sanitizeAllStrings(reflect.ValueOf(&v))

		if got, ok := v.(string); !ok || got != "ifaceval" {
			t.Errorf("interface value = %#v, want 'ifaceval'", v)
		}
	})

	t.Run("map with non-string key type: values stripped, keys untouched", func(t *testing.T) {
		m := map[int]string{1: "a\r\nb", 2: "clean"}
		sanitizeAllStrings(reflect.ValueOf(&m))

		want := map[int]string{1: "ab", 2: "clean"}
		if !reflect.DeepEqual(m, want) {
			t.Errorf("map = %#v, want %#v", m, want)
		}
	})
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
		{"CSI final byte at lower bound 0x40 '@'", "\x1b[@X", "X"},
		{"CSI final byte at upper bound 0x7e '~'", "\x1b[~X", "X"},
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
