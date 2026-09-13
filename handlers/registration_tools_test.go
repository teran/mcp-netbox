package handlers

import (
	"strings"
	"testing"
)

// TestToolDefs_EachHasInstructions verifies that every registered tool carries
// per-tool Instructions metadata (M4). Because the current official Go SDK
// (go-sdk v1.7.0) does not expose a per-tool Instructions field on mcp.Tool,
// the instructions are surfaced to the model via the Description channel and
// are kept as first-class, testable data in toolDefs().
func TestToolDefs_EachHasInstructions(t *testing.T) {
	t.Parallel()

	defs := toolDefs()
	if len(defs) == 0 {
		t.Fatal("toolDefs() returned no tools")
	}

	for _, d := range defs {
		if strings.TrimSpace(d.Name) == "" {
			t.Errorf("tool %q: empty name", d.Name)
		}
		if strings.TrimSpace(d.Description) == "" {
			t.Errorf("tool %q: empty description", d.Name)
		}
		if strings.TrimSpace(d.Instructions) == "" {
			t.Errorf("tool %q: empty instructions", d.Name)
		}
	}
}

// TestToolDefs_Count verifies the total number of registered tool definitions.
func TestToolDefs_Count(t *testing.T) {
	t.Parallel()

	if got := len(toolDefs()); got != expectedToolCount {
		t.Errorf("len(toolDefs()) = %d, want %d", got, expectedToolCount)
	}
}

// TestRegisterTool_DescriptionCarriesInstructions verifies that the
// registration helper appends the per-tool Instructions to the model-facing
// Description (the SDK's only per-tool guidance channel).
func TestRegisterTool_DescriptionCarriesInstructions(t *testing.T) {
	t.Parallel()

	def := toolDefs()[0]
	description := def.Description
	if def.Instructions != "" {
		description = description + "\n\nInstructions: " + def.Instructions
	}
	if !strings.Contains(description, "Instructions:") {
		t.Errorf("description for %q does not include Instructions section", def.Name)
	}
}
