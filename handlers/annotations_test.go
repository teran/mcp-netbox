package handlers

import "testing"

func TestBoolPtr(t *testing.T) {
	t.Parallel()

	if b := boolPtr(true); b == nil || !*b {
		t.Errorf("boolPtr(true) = %v, want pointer to true", b)
	}
	if b := boolPtr(false); b == nil || *b {
		t.Errorf("boolPtr(false) = %v, want pointer to false", b)
	}
}

func TestReadOnlyToolAnnotations(t *testing.T) {
	t.Parallel()

	a := readOnlyTool("List Sites")
	if a == nil {
		t.Fatal("readOnlyTool returned nil")
	}
	if a.Title != "List Sites" {
		t.Errorf("Title = %q, want %q", a.Title, "List Sites")
	}
	if !a.ReadOnlyHint {
		t.Error("ReadOnlyHint = false, want true")
	}
	if !a.IdempotentHint {
		t.Error("IdempotentHint = false, want true")
	}
	if a.DestructiveHint == nil || *a.DestructiveHint {
		t.Errorf("DestructiveHint = %v, want false", a.DestructiveHint)
	}
	if a.OpenWorldHint == nil || *a.OpenWorldHint {
		t.Errorf("OpenWorldHint = %v, want false (closed world)", a.OpenWorldHint)
	}
}
