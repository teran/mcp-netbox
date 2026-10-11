package handlers

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBoolPtr(t *testing.T) {
	t.Parallel()

	if b := new(true); b == nil || !*b {
		t.Errorf("boolPtr(true) = %v, want pointer to true", b)
	}
	if b := new(false); b == nil || *b {
		t.Errorf("boolPtr(false) = %v, want pointer to false", b)
	}
}

func TestReadOnlyToolAnnotations(t *testing.T) {
	t.Parallel()

	annotations := readOnlyTool("List Sites")
	if annotations == nil {
		t.Fatal("readOnlyTool returned nil")
	}
	if annotations.Title != "List Sites" {
		t.Errorf("Title = %q, want %q", annotations.Title, "List Sites")
	}
	if !annotations.ReadOnlyHint {
		t.Error("ReadOnlyHint = false, want true")
	}
	if !annotations.IdempotentHint {
		t.Error("IdempotentHint = false, want true")
	}
	if annotations.DestructiveHint == nil || *annotations.DestructiveHint {
		t.Errorf("DestructiveHint = %v, want false", annotations.DestructiveHint)
	}
	if annotations.OpenWorldHint == nil || *annotations.OpenWorldHint {
		t.Errorf("OpenWorldHint = %v, want false (closed world)", annotations.OpenWorldHint)
	}
}

func TestWriteToolAnnotations(t *testing.T) {
	t.Parallel()

	annotations := writeTool("Create Site")
	if annotations == nil {
		t.Fatal("writeTool returned nil")
	}
	if annotations.Title != "Create Site" {
		t.Errorf("Title = %q, want %q", annotations.Title, "Create Site")
	}
	if annotations.ReadOnlyHint {
		t.Error("ReadOnlyHint = true, want false")
	}
	if annotations.IdempotentHint {
		t.Error("IdempotentHint = true, want false (create is not idempotent)")
	}
	if annotations.DestructiveHint == nil || *annotations.DestructiveHint {
		t.Errorf("DestructiveHint = %v, want false", annotations.DestructiveHint)
	}
	if annotations.OpenWorldHint == nil || *annotations.OpenWorldHint {
		t.Errorf("OpenWorldHint = %v, want false (closed world)", annotations.OpenWorldHint)
	}
}

func TestDestructiveToolAnnotations(t *testing.T) {
	t.Parallel()

	annotations := destructiveTool("Delete Site")
	if annotations == nil {
		t.Fatal("destructiveTool returned nil")
	}
	if annotations.Title != "Delete Site" {
		t.Errorf("Title = %q, want %q", annotations.Title, "Delete Site")
	}
	if annotations.ReadOnlyHint {
		t.Error("ReadOnlyHint = true, want false")
	}
	if !annotations.IdempotentHint {
		t.Error("IdempotentHint = false, want true (delete is idempotent)")
	}
	if annotations.DestructiveHint == nil || !*annotations.DestructiveHint {
		t.Errorf("DestructiveHint = %v, want true", annotations.DestructiveHint)
	}
	if annotations.OpenWorldHint == nil || *annotations.OpenWorldHint {
		t.Errorf("OpenWorldHint = %v, want false (closed world)", annotations.OpenWorldHint)
	}
}

// TestDeleteToolAnnotationJSON verifies the destructive hint serialises as
// destructiveHint:true in the JSON annotations so clients can gate destructive
// tools.
func TestDeleteToolAnnotationJSON(t *testing.T) {
	t.Parallel()

	annotations := destructiveTool("Delete Site")
	data, err := json.Marshal(annotations)
	if err != nil {
		t.Fatalf("Marshal annotations: %v", err)
	}
	if !strings.Contains(string(data), `"destructiveHint":true`) {
		t.Errorf("annotations JSON = %s, want destructiveHint:true", data)
	}
	if strings.Contains(string(data), `"readOnlyHint":true`) {
		t.Errorf("annotations JSON = %s, want readOnlyHint not true", data)
	}
}
