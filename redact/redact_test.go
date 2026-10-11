package redact

import (
	"encoding/json"
	"strings"
	"testing"
)

type flat struct {
	Name   string `json:"name"`
	Token  string `json:"token" secret:"true"`
	Secret string `json:"secret" secret:"true"`
	Public string `json:"public"`
}

type nested struct {
	Top    string `json:"top"`
	Inner  flat   `json:"inner"`
	InnerP *flat  `json:"inner_ptr"`
}

type withTags struct {
	A string `secret:"true"`
	B string `secret:"false"`
	C string
	D int    `secret:"true"`
	E string `secret:"true"`
}

func TestRedactFlatStruct(t *testing.T) {
	t.Parallel()

	input := flat{Name: "srv", Token: "tok-123", Secret: "s3cr3t", Public: "visible"}
	out, ok := Redact(input).(flat)
	if !ok {
		t.Fatalf("Redact() type = %T, want flat", Redact(input))
	}
	if out.Token != Mask {
		t.Errorf("Token = %q, want %q", out.Token, Mask)
	}
	if out.Secret != Mask {
		t.Errorf("Secret = %q, want %q", out.Secret, Mask)
	}
	if out.Name != "srv" {
		t.Errorf("Name = %q, want %q", out.Name, "srv")
	}
	if out.Public != "visible" {
		t.Errorf("Public = %q, want %q", out.Public, "visible")
	}
}

func TestRedactNestedStruct(t *testing.T) {
	t.Parallel()

	input := nested{
		Top:    "top",
		Inner:  flat{Token: "inner-tok", Public: "inner-pub"},
		InnerP: &flat{Token: "ptr-tok", Public: "ptr-pub"},
	}
	redacted, ok := Redact(input).(nested)
	if !ok {
		t.Fatalf("Redact type = %T, want nested", Redact(input))
	}
	out := redacted

	if out.Inner.Token != Mask {
		t.Errorf("nested Inner.Token = %q, want %q", out.Inner.Token, Mask)
	}
	if out.Inner.Public != "inner-pub" {
		t.Errorf("nested Inner.Public = %q, want %q", out.Inner.Public, "inner-pub")
	}
	if out.InnerP == nil {
		t.Fatal("InnerP = nil, want non-nil pointer")
	}
	if out.InnerP.Token != Mask {
		t.Errorf("nested InnerP.Token = %q, want %q", out.InnerP.Token, Mask)
	}
	if out.InnerP.Public != "ptr-pub" {
		t.Errorf("nested InnerP.Public = %q, want %q", out.InnerP.Public, "ptr-pub")
	}
}

func TestRedactPointer(t *testing.T) {
	t.Parallel()

	t.Run("non-nil pointer", func(t *testing.T) {
		pointer := &flat{Token: "ptr-secret", Public: "ok"}
		out, ok := Redact(pointer).(*flat)
		if !ok {
			t.Fatalf("Redact() type = %T, want *flat", Redact(pointer))
		}
		if out == pointer {
			t.Error("Redact returned the same pointer, want a copy")
		}
		if out.Token != Mask {
			t.Errorf("Token = %q, want %q", out.Token, Mask)
		}
		if pointer.Token != "ptr-secret" {
			t.Errorf("original mutated: Token = %q, want %q", pointer.Token, "ptr-secret")
		}
	})

	t.Run("nil pointer", func(t *testing.T) {
		var pointer *flat
		out, ok := Redact(pointer).(*flat)
		if !ok {
			t.Fatalf("Redact(nil ptr) type = %T, want *flat", Redact(pointer))
		}
		if out != nil {
			t.Errorf("Redact(nil ptr) = %v, want nil", out)
		}
	})
}

func TestRedactSlice(t *testing.T) {
	t.Parallel()

	t.Run("slice of structs", func(t *testing.T) {
		input := []flat{{Token: "a", Public: "1"}, {Token: "encoded", Public: "2"}}
		out, ok := Redact(input).([]flat)
		if !ok {
			t.Fatalf("Redact() type = %T, want []flat", Redact(input))
		}
		if len(out) != 2 {
			t.Fatalf("len(out) = %d, want 2", len(out))
		}
		for i := range out {
			if out[i].Token != Mask {
				t.Errorf("out[%d].Token = %q, want %q", i, out[i].Token, Mask)
			}
		}
		if input[0].Token != "a" {
			t.Errorf("original mutated: input[0].Token = %q", input[0].Token)
		}
	})

	t.Run("nil slice", func(t *testing.T) {
		var input []flat
		out, ok := Redact(input).([]flat)
		if !ok {
			t.Fatalf("Redact(nil slice) type = %T, want []flat", Redact(input))
		}
		if out != nil {
			t.Errorf("Redact(nil slice) = %v, want nil", out)
		}
	})
}

func TestRedactArray(t *testing.T) {
	t.Parallel()

	input := [2]flat{{Token: "a", Public: "1"}, {Token: "encoded", Public: "2"}}
	out, ok := Redact(input).([2]flat)
	if !ok {
		t.Fatalf("Redact() type = %T, want [2]flat", Redact(input))
	}
	for i := range out {
		if out[i].Token != Mask {
			t.Errorf("out[%d].Token = %q, want %q", i, out[i].Token, Mask)
		}
	}
	if input[0].Token != "a" {
		t.Errorf("original mutated: input[0].Token = %q", input[0].Token)
	}
}

func TestRedactMap(t *testing.T) {
	t.Parallel()

	t.Run("values redacted, keys untouched", func(t *testing.T) {
		input := map[string]flat{
			"key-a":       {Token: "tok-a", Public: "pub-a"},
			"key-encoded": {Token: "tok-encoded", Public: "pub-encoded"},
		}
		out, ok := Redact(input).(map[string]flat)
		if !ok {
			t.Fatalf("Redact() type = %T, want map[string]flat", Redact(input))
		}
		if len(out) != 2 {
			t.Fatalf("len(out) = %d, want 2", len(out))
		}
		for k, v := range out {
			if k != "key-a" && k != "key-encoded" {
				t.Errorf("unexpected key %q", k)
			}
			if v.Token != Mask {
				t.Errorf("out[%q].Token = %q, want %q", k, v.Token, Mask)
			}
		}
		if input["key-a"].Token != "tok-a" {
			t.Errorf("original mutated: input[key-a].Token = %q", input["key-a"].Token)
		}
	})

	t.Run("nil map", func(t *testing.T) {
		var input map[string]flat
		out, ok := Redact(input).(map[string]flat)
		if !ok {
			t.Fatalf("Redact(nil map) type = %T, want map[string]flat", Redact(input))
		}
		if out != nil {
			t.Errorf("Redact(nil map) = %v, want nil", out)
		}
	})
}

func TestRedactInterface(t *testing.T) {
	t.Parallel()

	var input any = flat{Token: "iface-secret", Public: "iface-pub"}
	out, ok := Redact(input).(flat)
	if !ok {
		t.Fatalf("Redact() type = %T, want flat", Redact(input))
	}
	if out.Token != Mask {
		t.Errorf("interface Token = %q, want %q", out.Token, Mask)
	}
	if out.Public != "iface-pub" {
		t.Errorf("interface Public = %q, want %q", out.Public, "iface-pub")
	}
}

func TestRedactUnexportedFieldsSkipped(t *testing.T) {
	t.Parallel()

	type str struct {
		Exported string `secret:"true"`
		hidden   string
	}
	input := str{Exported: "exported-secret", hidden: "hidden-value"}
	out, ok := Redact(input).(str)
	if !ok {
		t.Fatalf("Redact() type = %T, want str", Redact(input))
	}
	if out.Exported != Mask {
		t.Errorf("Exported = %q, want %q", out.Exported, Mask)
	}
	// Unexported field is preserved (copied verbatim).
	if out.hidden != "hidden-value" {
		t.Errorf("hidden = %q, want %q", out.hidden, "hidden-value")
	}
}

func TestRedactDoesNotMutateOriginal(t *testing.T) {
	t.Parallel()

	inner := &flat{Token: "deep-tok", Public: "deep-pub"}
	input := nested{
		Top:    "top",
		Inner:  flat{Token: "inner-tok", Public: "inner-pub"},
		InnerP: inner,
	}

	_ = Redact(input)

	// The original value, its nested structs, and its pointer target must be
	// completely untouched.
	if input.Top != "top" {
		t.Errorf("Top mutated: %q", input.Top)
	}
	if input.Inner.Token != "inner-tok" {
		t.Errorf("Inner.Token mutated: %q", input.Inner.Token)
	}
	if input.Inner.Public != "inner-pub" {
		t.Errorf("Inner.Public mutated: %q", input.Inner.Public)
	}
	if input.InnerP == nil {
		t.Fatal("InnerP became nil")
	}
	if input.InnerP.Token != "deep-tok" {
		t.Errorf("InnerP.Token mutated: %q", input.InnerP.Token)
	}
	if inner.Token != "deep-tok" {
		t.Errorf("shared pointer target mutated: %q", inner.Token)
	}
	if inner.Public != "deep-pub" {
		t.Errorf("shared pointer Public mutated: %q", inner.Public)
	}
}

func TestRedactNilAndScalars(t *testing.T) {
	t.Parallel()

	if Redact(nil) != nil {
		t.Errorf("Redact(nil) = %v, want nil", Redact(nil))
	}

	if got := Redact(42); got != 42 {
		t.Errorf("Redact(42) = %v, want 42", got)
	}
	if got := Redact(3.14); got != 3.14 {
		t.Errorf("Redact(3.14) = %v, want 3.14", got)
	}
	if got := Redact(true); got != true {
		t.Errorf("Redact(true) = %v, want true", got)
	}
	if got := Redact("plain"); got != "plain" {
		t.Errorf(`Redact("plain") = %v, want "plain"`, got)
	}
}

func TestRedactSecretFalseDoesNotMask(t *testing.T) {
	t.Parallel()

	input := withTags{A: "secret-a", B: "not-secret", C: "plain", D: 7, E: "secret-e"}
	out, ok := Redact(input).(withTags)
	if !ok {
		t.Fatalf("Redact() type = %T, want withTags", Redact(input))
	}
	if out.A != Mask {
		t.Errorf("A = %q, want %q (secret:true)", out.A, Mask)
	}
	if out.E != Mask {
		t.Errorf("E = %q, want %q (secret:true)", out.E, Mask)
	}
	// secret:"false" and no-tag fields are NOT masked.
	if out.B != "not-secret" {
		t.Errorf("B = %q, want %q (secret:false)", out.B, "not-secret")
	}
	if out.C != "plain" {
		t.Errorf("C = %q, want %q (no tag)", out.C, "plain")
	}
	// Non-string secret field: recursively redacted, but the value is an int so
	// it is left as its (non-string) copy — it must not be turned into Mask.
	if out.D != 7 {
		t.Errorf("D = %d, want %d", out.D, 7)
	}
}

func TestMarshalJSON(t *testing.T) {
	t.Parallel()

	input := flat{Name: "srv", Token: "tok-json", Secret: "s3cr3t-json", Public: "pub-json"}
	encoded, err := MarshalJSON(input)
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}
	raw := string(encoded)
	if strings.Contains(raw, "tok-json") {
		t.Errorf("MarshalJSON leaks Token: %str", raw)
	}
	if strings.Contains(raw, "s3cr3t-json") {
		t.Errorf("MarshalJSON leaks Secret: %str", raw)
	}
	if !strings.Contains(raw, Mask) {
		t.Errorf("MarshalJSON missing mask, got: %str", raw)
	}

	// Must be valid JSON and decodable back.
	var decoded map[string]string
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("MarshalJSON output not valid JSON: %v", err)
	}
	if decoded["token"] != Mask {
		t.Errorf(`decoded["token"] = %q, want %q`, decoded["token"], Mask)
	}
	if decoded["public"] != "pub-json" {
		t.Errorf(`decoded["public"] = %q, want %q`, decoded["public"], "pub-json")
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	input := flat{Name: "srv", Token: "tok-str", Secret: "s3cr3t-str", Public: "pub"}
	str := String(input)
	if strings.Contains(str, "tok-str") {
		t.Errorf("String leaks Token: %str", str)
	}
	if strings.Contains(str, "s3cr3t-str") {
		t.Errorf("String leaks Secret: %str", str)
	}
	if !strings.Contains(str, Mask) {
		t.Errorf("String missing mask, got: %str", str)
	}
}
