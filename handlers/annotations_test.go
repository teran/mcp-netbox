package handlers

import (
	"encoding/json"
	"strings"
	"testing"
)

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

func TestWriteToolAnnotations(t *testing.T) {
	t.Parallel()

	a := writeTool("Create Site")
	if a == nil {
		t.Fatal("writeTool returned nil")
	}
	if a.Title != "Create Site" {
		t.Errorf("Title = %q, want %q", a.Title, "Create Site")
	}
	if a.ReadOnlyHint {
		t.Error("ReadOnlyHint = true, want false")
	}
	if a.IdempotentHint {
		t.Error("IdempotentHint = true, want false (create is not idempotent)")
	}
	if a.DestructiveHint == nil || *a.DestructiveHint {
		t.Errorf("DestructiveHint = %v, want false", a.DestructiveHint)
	}
	if a.OpenWorldHint == nil || *a.OpenWorldHint {
		t.Errorf("OpenWorldHint = %v, want false (closed world)", a.OpenWorldHint)
	}
}

func TestDestructiveToolAnnotations(t *testing.T) {
	t.Parallel()

	a := destructiveTool("Delete Site")
	if a == nil {
		t.Fatal("destructiveTool returned nil")
	}
	if a.Title != "Delete Site" {
		t.Errorf("Title = %q, want %q", a.Title, "Delete Site")
	}
	if a.ReadOnlyHint {
		t.Error("ReadOnlyHint = true, want false")
	}
	if !a.IdempotentHint {
		t.Error("IdempotentHint = false, want true (delete is idempotent)")
	}
	if a.DestructiveHint == nil || !*a.DestructiveHint {
		t.Errorf("DestructiveHint = %v, want true", a.DestructiveHint)
	}
	if a.OpenWorldHint == nil || *a.OpenWorldHint {
		t.Errorf("OpenWorldHint = %v, want false (closed world)", a.OpenWorldHint)
	}
}

// TestDeleteToolAnnotationJSON verifies the destructive hint serialises as
// destructiveHint:true in the JSON annotations so clients can gate destructive
// tools.
func TestDeleteToolAnnotationJSON(t *testing.T) {
	t.Parallel()

	a := destructiveTool("Delete Site")
	b, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("Marshal annotations: %v", err)
	}
	if !strings.Contains(string(b), `"destructiveHint":true`) {
		t.Errorf("annotations JSON = %s, want destructiveHint:true", b)
	}
	if strings.Contains(string(b), `"readOnlyHint":true`) {
		t.Errorf("annotations JSON = %s, want readOnlyHint not true", b)
	}
}
