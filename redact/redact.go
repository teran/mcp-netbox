// Package redact provides annotation-based deep redaction of arbitrary Go
// values. Any struct field tagged with secret:"true" is replaced with Mask in
// a deep copy of the value, so secrets never leak into logs, metrics text, or
// serialized output.
//
// The package is intentionally stdlib-only (reflection based) apart from the
// optional logrus hook (hook.go) that wires redaction into logging.
package redact

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Mask is the replacement value applied to every struct field tagged
// secret:"true".
const Mask = "***redacted***"

// Redact returns a deep copy of v in which every struct field carrying the
// struct tag secret:"true" is replaced with Mask. The input is never mutated.
//
// The copy preserves the dynamic type of v, so callers can type-assert the
// result back to the original concrete type. Supported shapes:
//
//   - nested structs
//   - pointers (nil-safe)
//   - slices and arrays
//   - maps (values are redacted; keys are not)
//   - interfaces
//
// Unexported struct fields are skipped (copied verbatim, never rewritten).
// Redact(nil) and scalar values (numbers, bools, strings without a secret
// tag) are returned unchanged.
func Redact(v any) any {
	if v == nil {
		return nil
	}
	return redact(reflect.ValueOf(v))
}

// redact recursively returns a redacted deep copy of the value held by rv,
// preserving its dynamic type.
func redact(rv reflect.Value) any {
	switch rv.Kind() {
	case reflect.Pointer:
		if rv.IsNil() {
			return rv.Interface()
		}
		np := reflect.New(rv.Elem().Type())
		np.Elem().Set(reflect.ValueOf(redact(rv.Elem())))
		return np.Interface()

	case reflect.Interface:
		if rv.IsNil() {
			return rv.Interface()
		}
		return redact(rv.Elem())

	case reflect.Struct:
		return redactStruct(rv)

	case reflect.Slice:
		return redactSlice(rv)

	case reflect.Array:
		return redactArray(rv)

	case reflect.Map:
		return redactMap(rv)

	default:
		return rv.Interface()
	}
}

// redactStruct returns a deep copy of the struct value rv with every exported
// field tagged secret:"true" replaced by Mask and every composite field
// deep-redacted recursively. Unexported fields are copied verbatim.
func redactStruct(rv reflect.Value) any {
	t := rv.Type()
	out := reflect.New(t).Elem()
	// Copy the whole struct (including unexported fields) before rewriting.
	out.Set(rv)

	for i := range t.NumField() {
		field := t.Field(i)
		if field.PkgPath != "" {
			// Unexported field — never rewritten.
			continue
		}

		fv := out.Field(i)
		if field.Tag.Get("secret") == "true" {
			if fv.Kind() == reflect.String {
				fv.SetString(Mask)
			} else {
				fv.Set(reflect.ValueOf(redact(fv)))
			}
			continue
		}

		// Deep-redact composite fields so nested secrets are caught.
		switch fv.Kind() {
		case reflect.Struct, reflect.Pointer, reflect.Slice, reflect.Array, reflect.Map, reflect.Interface:
			if v, ok := settableRedact(fv); ok {
				fv.Set(v)
			}
		}
	}

	return out.Interface()
}

// settableRedact redacts the value held by fv and returns a reflect.Value that
// can be Set into fv, or ok=false when the field should be left as-is (e.g. a
// nil interface).
func settableRedact(fv reflect.Value) (reflect.Value, bool) {
	rv := reflect.ValueOf(redact(fv))
	if !rv.IsValid() {
		return reflect.Value{}, false
	}
	return rv, true
}

// redactSlice returns a new slice whose elements are deep-redacted copies.
func redactSlice(rv reflect.Value) any {
	if rv.IsNil() {
		return rv.Interface()
	}
	out := reflect.MakeSlice(rv.Type(), rv.Len(), rv.Len())
	for i := range rv.Len() {
		if v, ok := settableRedact(rv.Index(i)); ok {
			out.Index(i).Set(v)
		}
	}
	return out.Interface()
}

// redactArray returns a new array whose elements are deep-redacted copies.
func redactArray(rv reflect.Value) any {
	out := reflect.New(rv.Type()).Elem()
	for i := range rv.Len() {
		if v, ok := settableRedact(rv.Index(i)); ok {
			out.Index(i).Set(v)
		}
	}
	return out.Interface()
}

// redactMap returns a new map whose values are deep-redacted copies; keys are
// copied unchanged.
func redactMap(rv reflect.Value) any {
	if rv.IsNil() {
		return rv.Interface()
	}
	out := reflect.MakeMapWithSize(rv.Type(), rv.Len())
	iter := rv.MapRange()
	for iter.Next() {
		k := iter.Key()
		val := iter.Value()
		redacted, ok := settableRedact(val)
		if !ok {
			redacted = reflect.Zero(rv.Type().Elem())
		}
		out.SetMapIndex(k, redacted)
	}
	return out.Interface()
}

// MarshalJSON marshals Redact(v) to JSON, so any secret:"true" field is masked
// before serialization.
func MarshalJSON(v any) ([]byte, error) {
	return json.Marshal(Redact(v))
}

// String returns a string representation of Redact(v). Secrets are replaced by
// Mask, so the result never contains a tagged secret.
func String(v any) string {
	return fmt.Sprintf("%v", Redact(v))
}
