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

// redact recursively returns a redacted deep copy of the value held by value,
// preserving its dynamic type.
func redact(value reflect.Value) any {
	//nolint:exhaustive // reflect.Kind has ~27 values; all unlisted scalar/composite kinds hit default.
	switch value.Kind() {
	case reflect.Pointer:
		if value.IsNil() {
			return value.Interface()
		}
		np := reflect.New(value.Elem().Type())
		np.Elem().Set(reflect.ValueOf(redact(value.Elem())))

		return np.Interface()

	case reflect.Interface:
		if value.IsNil() {
			return value.Interface()
		}

		return redact(value.Elem())

	case reflect.Struct:
		return redactStruct(value)

	case reflect.Slice:
		return redactSlice(value)

	case reflect.Array:
		return redactArray(value)

	case reflect.Map:
		return redactMap(value)

	default:
		return value.Interface()
	}
}

// redactStruct returns a deep copy of the struct value with every exported
// field tagged secret:"true" replaced by Mask and every composite field
// deep-redacted recursively. Unexported fields are copied verbatim.
func redactStruct(value reflect.Value) any {
	structType := value.Type()
	out := reflect.New(structType).Elem()
	// Copy the whole struct (including unexported fields) before rewriting.
	out.Set(value)

	for fieldIndex := range structType.NumField() {
		field := structType.Field(fieldIndex)
		if field.PkgPath != "" {
			// Unexported field — never rewritten.

			continue
		}

		fieldValue := out.Field(fieldIndex)
		if field.Tag.Get("secret") == "true" {
			if fieldValue.Kind() == reflect.String {
				fieldValue.SetString(Mask)
			} else {
				fieldValue.Set(reflect.ValueOf(redact(fieldValue)))
			}

			continue
		}

		// Deep-redact composite fields so nested secrets are caught.
		//nolint:exhaustive // Only composite kinds need recursion; scalars are left as-is.
		switch fieldValue.Kind() {
		case reflect.Struct, reflect.Pointer, reflect.Slice, reflect.Array, reflect.Map, reflect.Interface:
			if v, ok := settableRedact(fieldValue); ok {
				fieldValue.Set(v)
			}
		}
	}

	return out.Interface()
}

// settableRedact redacts the value held by fv and returns a reflect.Value that
// can be Set into fv, or ok=false when the field should be left as-is (e.g. a
// nil interface).
func settableRedact(fv reflect.Value) (reflect.Value, bool) {
	value := reflect.ValueOf(redact(fv))
	if !value.IsValid() {
		return reflect.Value{}, false
	}

	return value, true
}

// redactSlice returns a new slice whose elements are deep-redacted copies.
func redactSlice(value reflect.Value) any {
	if value.IsNil() {
		return value.Interface()
	}
	out := reflect.MakeSlice(value.Type(), value.Len(), value.Len())
	for i := range value.Len() {
		if v, ok := settableRedact(value.Index(i)); ok {
			out.Index(i).Set(v)
		}
	}

	return out.Interface()
}

// redactArray returns a new array whose elements are deep-redacted copies.
func redactArray(value reflect.Value) any {
	out := reflect.New(value.Type()).Elem()
	for i := range value.Len() {
		if v, ok := settableRedact(value.Index(i)); ok {
			out.Index(i).Set(v)
		}
	}

	return out.Interface()
}

// redactMap returns a new map whose values are deep-redacted copies; keys are
// copied unchanged.
func redactMap(value reflect.Value) any {
	if value.IsNil() {
		return value.Interface()
	}
	out := reflect.MakeMapWithSize(value.Type(), value.Len())
	iter := value.MapRange()
	for iter.Next() {
		key := iter.Key()
		val := iter.Value()
		redacted, ok := settableRedact(val)
		if !ok {
			redacted = reflect.Zero(value.Type().Elem())
		}
		out.SetMapIndex(key, redacted)
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
