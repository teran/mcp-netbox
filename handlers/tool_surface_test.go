package handlers

import (
	"testing"
)

// TestToolSurface_MatchesToolDefs verifies that ToolSurface() derives exactly
// the tools registered by RegisterTools: the same count as toolDefs() and in
// the same Name order (the read -> write -> delete grouping order of toolDefs).
func TestToolSurface_MatchesToolDefs(t *testing.T) {
	defs := toolDefs()
	surface := ToolSurface()

	if len(surface) != len(defs) {
		t.Fatalf("ToolSurface() returned %d entries, want %d (== toolDefs count)", len(surface), len(defs))
	}

	for i, def := range defs {
		if surface[i].Name != def.Name {
			t.Errorf("ToolSurface()[%d].Name = %q, want %q (order must match toolDefs)", i, surface[i].Name, def.Name)
		}
		if surface[i].Title != def.Title {
			t.Errorf("ToolSurface()[%d].Title = %q, want %q", i, surface[i].Title, def.Title)
		}
	}
}

// TestToolSurface_EntriesWellFormed verifies every surface entry carries a
// non-empty Name and Title and is marked AlwaysPresent (all mcp-netbox tools
// are static, L11).
func TestToolSurface_EntriesWellFormed(t *testing.T) {
	for _, entry := range ToolSurface() {
		if entry.Name == "" {
			t.Error("ToolSurface() returned an entry with empty Name")
		}
		if entry.Title == "" {
			t.Errorf("ToolSurface() entry %q has empty Title", entry.Name)
		}
		if !entry.AlwaysPresent {
			t.Errorf("ToolSurface() entry %q: AlwaysPresent = false, want true (static tool)", entry.Name)
		}
	}
}

// TestToolSurface_ContainsKnownTools spot-checks that the surface set includes
// representative read, create, update, delete and generic tools.
func TestToolSurface_ContainsKnownTools(t *testing.T) {
	names := make(map[string]bool, len(ToolSurface()))
	for _, entry := range ToolSurface() {
		names[entry.Name] = true
	}

	for _, want := range []string{"get_sites", "create_site", "update_site", "delete_site", "get_object_by_id"} {
		if !names[want] {
			t.Errorf("ToolSurface() is missing known tool %q", want)
		}
	}
}
