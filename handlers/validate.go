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
func sanitizeControl(value string) string {
	return strings.Map(func(r rune) rune {
		if r == '\r' || r == '\n' {
			return -1
		}

		return r
	}, value)
}

const (
	// escByte is the ASCII ESC byte that begins an ANSI escape sequence.
	escByte = 0x1b
	// csiOffset is the length of the "ESC [" CSI prefix (escByte + '[').
	csiOffset = 2
)

// sanitizeOutput strips ANSI escape sequences and C0 control characters from
// tool-result text (S09/N23). It is applied to the text form of a tool result
// so that terminal escape sequences or control characters embedded in NetBox
// data (e.g. in device names or descriptions) can never reach the model. It
// operates on bytes: only ASCII control bytes (< 0x20, 0x7f) and whole ANSI
// escape sequences are removed, so valid multi-byte UTF-8 and printable ASCII
// pass through untouched.
func sanitizeOutput(value string) string {
	var builder strings.Builder
	builder.Grow(len(value))
	for idx := 0; idx < len(value); {
		currentByte := value[idx]
		if currentByte == escByte { // ESC — ANSI escape sequence
			// CSI: ESC '[' ... final byte in 0x40-0x7e. When no final byte is
			// present before the end of the string, the whole remainder is part
			// of the incomplete escape and is dropped.
			if idx+1 < len(value) && value[idx+1] == '[' {
				j := idx + csiOffset
				for j < len(value) && (value[j] < 0x40 || value[j] > 0x7e) {
					j++
				}

				idx = j + 1

				continue
			}
			// Two-character escape: ESC + a printable byte.
			idx += 2

			continue
		}
		if currentByte < 0x20 || currentByte == 0x7f { // C0 control char or DEL
			idx++

			continue
		}
		builder.WriteByte(currentByte)
		idx++
	}

	return builder.String()
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
func sanitizeAllStrings(value reflect.Value) {
	if !value.IsValid() {
		return
	}
	//nolint:exhaustive // reflect.Kind has ~27 values; only pointer/struct/string are handled here.
	switch value.Kind() {
	case reflect.Pointer:
		if !value.IsNil() {
			sanitizeAllStrings(value.Elem())
		}
	case reflect.Struct:
		for _, field := range value.Fields() {
			sanitizeAllStrings(field)
		}
	case reflect.String:
		if value.CanSet() {
			value.SetString(sanitizeControl(value.String()))
		}
	}
}
