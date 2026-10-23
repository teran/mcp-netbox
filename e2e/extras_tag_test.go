//go:build e2e

package e2e

// CRUD helper for the extras Tag entity added in feature/expand-entities
// Batch 5b. It registers itself in init() so TestNetBoxE2E picks it up, and
// drives create_tag/update_tag/delete_tag (plus get_tags enumeration via
// readTools and get_object_by_id) over the full MCP protocol against a real
// NetBox.

import "testing"

func init() {
	registerCRUD("tag", tagCRUD)
}

// tagCRUD covers create_tag/update_tag/delete_tag. A tag is standalone in
// NetBox's extras app (only name is required), so it can be created bare.
func tagCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-tag-" + runID

	created := c.call("create_tag", map[string]interface{}{"name": name, "slug": name, "color": "ff0000"})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created tag name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("tag", id)); c.str(data, "name") != name {
		t.Fatalf("get tag name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_tag", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("tag", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get tag after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_tag", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "tag", "id": id})
}
