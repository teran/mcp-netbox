//go:build e2e

package e2e

// CRUD helper for the extras CustomField entity added in feature/expand-entities
// Batch 5c. It registers itself in init() so TestNetBoxE2E picks it up, and
// drives create_custom_field/update_custom_field/delete_custom_field (plus
// get_custom_fields enumeration and get_object_by_id) over the full MCP protocol
// against a real NetBox.

import "testing"

func init() {
	registerCRUD("custom_field", customFieldCRUD)
}

// customFieldCRUD covers create_custom_field/update_custom_field/delete_custom_field.
// NetBox requires name, label, a valid type (e.g. "text") and at least one
// content type on create.
func customFieldCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e_cf_" + runID
	label := "E2E Custom Field"

	created := c.call("create_custom_field", map[string]interface{}{
		"name":         name,
		"label":        label,
		"type":         "text",
		"object_types": []interface{}{"ipam.prefix", "ipam.vlan"},
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created custom field name = %q, want %q", got, name)
	}
	if got := c.str(c.data(created), "label"); got != label {
		t.Fatalf("created custom field label = %q, want %q", got, label)
	}

	if data := c.data(c.getByID("custom_field", id)); c.str(data, "name") != name {
		t.Fatalf("get custom field name = %q, want %q", c.str(data, "name"), name)
	}

	// Also enumerate custom fields and confirm ours is present.
	found := false
	for _, item := range c.list(c.call("get_custom_fields", nil)) {
		obj, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if v, _ := obj["name"].(string); v == name {
			found = true

			break
		}
	}
	if !found {
		t.Fatalf("get_custom_fields did not return %q", name)
	}

	c.call("update_custom_field", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("custom_field", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get custom field after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_custom_field", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "custom_field", "id": id})
}
