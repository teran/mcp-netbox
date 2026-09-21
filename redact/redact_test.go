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

	in := flat{Name: "srv", Token: "tok-123", Secret: "s3cr3t", Public: "visible"}
	out, ok := Redact(in).(flat)
	if !ok {
		t.Fatalf("Redact() type = %T, want flat", Redact(in))
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

	in := nested{
		Top:    "top",
		Inner:  flat{Token: "inner-tok", Public: "inner-pub"},
		InnerP: &flat{Token: "ptr-tok", Public: "ptr-pub"},
	}
	out := Redact(in).(nested)

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
		p := &flat{Token: "ptr-secret", Public: "ok"}
		out, ok := Redact(p).(*flat)
		if !ok {
			t.Fatalf("Redact() type = %T, want *flat", Redact(p))
		}
		if out == p {
			t.Error("Redact returned the same pointer, want a copy")
		}
		if out.Token != Mask {
			t.Errorf("Token = %q, want %q", out.Token, Mask)
		}
		if p.Token != "ptr-secret" {
			t.Errorf("original mutated: Token = %q, want %q", p.Token, "ptr-secret")
		}
	})

	t.Run("nil pointer", func(t *testing.T) {
		var p *flat
		out, ok := Redact(p).(*flat)
		if !ok {
			t.Fatalf("Redact(nil ptr) type = %T, want *flat", Redact(p))
		}
		if out != nil {
			t.Errorf("Redact(nil ptr) = %v, want nil", out)
		}
	})
}

func TestRedactSlice(t *testing.T) {
	t.Parallel()

	t.Run("slice of structs", func(t *testing.T) {
		in := []flat{{Token: "a", Public: "1"}, {Token: "b", Public: "2"}}
		out, ok := Redact(in).([]flat)
		if !ok {
			t.Fatalf("Redact() type = %T, want []flat", Redact(in))
		}
		if len(out) != 2 {
			t.Fatalf("len(out) = %d, want 2", len(out))
		}
		for i := range out {
			if out[i].Token != Mask {
				t.Errorf("out[%d].Token = %q, want %q", i, out[i].Token, Mask)
			}
		}
		if in[0].Token != "a" {
			t.Errorf("original mutated: in[0].Token = %q", in[0].Token)
		}
	})

	t.Run("nil slice", func(t *testing.T) {
		var in []flat
		out, ok := Redact(in).([]flat)
		if !ok {
			t.Fatalf("Redact(nil slice) type = %T, want []flat", Redact(in))
		}
		if out != nil {
			t.Errorf("Redact(nil slice) = %v, want nil", out)
		}
	})
}

func TestRedactArray(t *testing.T) {
	t.Parallel()

	in := [2]flat{{Token: "a", Public: "1"}, {Token: "b", Public: "2"}}
	out, ok := Redact(in).([2]flat)
	if !ok {
		t.Fatalf("Redact() type = %T, want [2]flat", Redact(in))
	}
	for i := range out {
		if out[i].Token != Mask {
			t.Errorf("out[%d].Token = %q, want %q", i, out[i].Token, Mask)
		}
	}
	if in[0].Token != "a" {
		t.Errorf("original mutated: in[0].Token = %q", in[0].Token)
	}
}

func TestRedactMap(t *testing.T) {
	t.Parallel()

	t.Run("values redacted, keys untouched", func(t *testing.T) {
		in := map[string]flat{
			"key-a": {Token: "tok-a", Public: "pub-a"},
			"key-b": {Token: "tok-b", Public: "pub-b"},
		}
		out, ok := Redact(in).(map[string]flat)
		if !ok {
			t.Fatalf("Redact() type = %T, want map[string]flat", Redact(in))
		}
		if len(out) != 2 {
			t.Fatalf("len(out) = %d, want 2", len(out))
		}
		for k, v := range out {
			if k != "key-a" && k != "key-b" {
				t.Errorf("unexpected key %q", k)
			}
			if v.Token != Mask {
				t.Errorf("out[%q].Token = %q, want %q", k, v.Token, Mask)
			}
		}
		if in["key-a"].Token != "tok-a" {
			t.Errorf("original mutated: in[key-a].Token = %q", in["key-a"].Token)
		}
	})

	t.Run("nil map", func(t *testing.T) {
		var in map[string]flat
		out, ok := Redact(in).(map[string]flat)
		if !ok {
			t.Fatalf("Redact(nil map) type = %T, want map[string]flat", Redact(in))
		}
		if out != nil {
			t.Errorf("Redact(nil map) = %v, want nil", out)
		}
	})
}

func TestRedactInterface(t *testing.T) {
	t.Parallel()

	var in any = flat{Token: "iface-secret", Public: "iface-pub"}
	out, ok := Redact(in).(flat)
	if !ok {
		t.Fatalf("Redact() type = %T, want flat", Redact(in))
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

	type s struct {
		Exported string `secret:"true"`
		hidden   string
	}
	in := s{Exported: "exported-secret", hidden: "hidden-value"}
	out, ok := Redact(in).(s)
	if !ok {
		t.Fatalf("Redact() type = %T, want s", Redact(in))
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
	in := nested{
		Top:    "top",
		Inner:  flat{Token: "inner-tok", Public: "inner-pub"},
		InnerP: inner,
	}

	_ = Redact(in)

	// The original value, its nested structs, and its pointer target must be
	// completely untouched.
	if in.Top != "top" {
		t.Errorf("Top mutated: %q", in.Top)
	}
	if in.Inner.Token != "inner-tok" {
		t.Errorf("Inner.Token mutated: %q", in.Inner.Token)
	}
	if in.Inner.Public != "inner-pub" {
		t.Errorf("Inner.Public mutated: %q", in.Inner.Public)
	}
	if in.InnerP == nil {
		t.Fatal("InnerP became nil")
	}
	if in.InnerP.Token != "deep-tok" {
		t.Errorf("InnerP.Token mutated: %q", in.InnerP.Token)
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

	in := withTags{A: "secret-a", B: "not-secret", C: "plain", D: 7, E: "secret-e"}
	out, ok := Redact(in).(withTags)
	if !ok {
		t.Fatalf("Redact() type = %T, want withTags", Redact(in))
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

	in := flat{Name: "srv", Token: "tok-json", Secret: "s3cr3t-json", Public: "pub-json"}
	b, err := MarshalJSON(in)
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}
	raw := string(b)
	if strings.Contains(raw, "tok-json") {
		t.Errorf("MarshalJSON leaks Token: %s", raw)
	}
	if strings.Contains(raw, "s3cr3t-json") {
		t.Errorf("MarshalJSON leaks Secret: %s", raw)
	}
	if !strings.Contains(raw, Mask) {
		t.Errorf("MarshalJSON missing mask, got: %s", raw)
	}

	// Must be valid JSON and decodable back.
	var decoded map[string]string
	if err := json.Unmarshal(b, &decoded); err != nil {
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

	in := flat{Name: "srv", Token: "tok-str", Secret: "s3cr3t-str", Public: "pub"}
	s := String(in)
	if strings.Contains(s, "tok-str") {
		t.Errorf("String leaks Token: %s", s)
	}
	if strings.Contains(s, "s3cr3t-str") {
		t.Errorf("String leaks Secret: %s", s)
	}
	if !strings.Contains(s, Mask) {
		t.Errorf("String missing mask, got: %s", s)
	}
}
