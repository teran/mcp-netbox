package handlers

// ToolSurfaceEntry describes one tool in the registered surface for the
// startup log (L10/L05GO). mcp-netbox registers only static tools, so
// AlwaysPresent is always true (L11).
type ToolSurfaceEntry struct {
	Name          string
	Title         string
	AlwaysPresent bool
}

// ToolSurface returns the exact set of tools registered by RegisterTools —
// the set returned by tools/list. It is derived from toolDefs() and drives
// the L10 startup log of the registered tool surface.
func ToolSurface() []ToolSurfaceEntry {
	defs := toolDefs()
	out := make([]ToolSurfaceEntry, 0, len(defs))
	for _, d := range defs {
		out = append(out, ToolSurfaceEntry{Name: d.Name, Title: d.Title, AlwaysPresent: true})
	}

	return out
}
