package handlers

import (
	"fmt"
	"reflect"
	"strings"
)

// validatePositiveID rejects non-positive identifiers for object-targeted tools.
func validatePositiveID(id int) error {
	if id <= 0 {
		return fmt.Errorf("id must be a positive integer")
	}
	return nil
}

// sanitizeControl strips CR and LF control characters, preventing header/query
// injection via crafted payload values.
func sanitizeControl(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '\r' || r == '\n' {
			return -1
		}
		return r
	}, s)
}

// sanitizeOutput strips ANSI escape sequences and C0 control characters from
// tool-result text (S09/N23). It is applied to the text form of a tool result
// so that terminal escape sequences or control characters embedded in NetBox
// data (e.g. in device names or descriptions) can never reach the model. It
// operates on bytes: only ASCII control bytes (< 0x20, 0x7f) and whole ANSI
// escape sequences are removed, so valid multi-byte UTF-8 and printable ASCII
// pass through untouched.
func sanitizeOutput(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); {
		c := s[i]
		if c == 0x1b { // ESC — ANSI escape sequence
			// CSI: ESC '[' ... final byte in 0x40-0x7e. When no final byte is
			// present before the end of the string, the whole remainder is part
			// of the incomplete escape and is dropped.
			if i+1 < len(s) && s[i+1] == '[' {
				j := i + 2
				for j < len(s) && (s[j] < 0x40 || s[j] > 0x7e) {
					j++
				}
				i = j + 1
				continue
			}
			// Two-character escape: ESC + a printable byte.
			i += 2
			continue
		}
		if c < 0x20 || c == 0x7f { // C0 control char or DEL
			i++
			continue
		}
		b.WriteByte(c)
		i++
	}
	return b.String()
}

// sanitizeMapKeysValues applies sanitizeControl to every key and value of a
// string map, returning a new map.
func sanitizeMapKeysValues(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[sanitizeControl(k)] = sanitizeControl(v)
	}
	return out
}

// sanitizeAllStrings recursively sanitizes every settable string field of a
// struct value (including through pointers). It is applied to write payloads
// before json.Marshal so CRLF never reaches the NetBox request body.
func sanitizeAllStrings(v reflect.Value) {
	if !v.IsValid() {
		return
	}
	switch v.Kind() {
	case reflect.Pointer:
		if !v.IsNil() {
			sanitizeAllStrings(v.Elem())
		}
	case reflect.Struct:
		for i := range v.NumField() {
			sanitizeAllStrings(v.Field(i))
		}
	case reflect.String:
		if v.CanSet() {
			v.SetString(sanitizeControl(v.String()))
		}
	}
}
