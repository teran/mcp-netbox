//go:build e2e

// Package e2e CRUD helpers for the DCIM entities added in the
// feature/expand-entities branch: site_group, rack_group, rack_role,
// device_role, platform, virtual_chassis, power_panel, inventory_item_role.
// Each registers itself in init() so TestNetBoxE2E picks it up.
package e2e

import (
	"testing"
)

func init() {
	registerCRUD("site_group", siteGroupCRUD)
	registerCRUD("rack_group", rackGroupCRUD)
	registerCRUD("rack_role", rackRoleCRUD)
	registerCRUD("device_role", deviceRoleCRUD)
	registerCRUD("platform", platformCRUD)
	registerCRUD("virtual_chassis", virtualChassisCRUD)
	registerCRUD("power_panel", powerPanelCRUD)
	registerCRUD("inventory_item_role", inventoryItemRoleCRUD)
}

func siteGroupCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-site-group-" + runID

	created := c.call("create_site_group", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created site group name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("site_group", id)); c.str(data, "name") != name {
		t.Fatalf("get site group name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_site_group", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("site_group", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get site group after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_site_group", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "site_group", "id": id})
}

func rackGroupCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-rack-group-" + runID

	created := c.call("create_rack_group", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created rack group name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("rack_group", id)); c.str(data, "name") != name {
		t.Fatalf("get rack group name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_rack_group", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("rack_group", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get rack group after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_rack_group", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "rack_group", "id": id})
}

func rackRoleCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-rack-role-" + runID

	created := c.call("create_rack_role", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created rack role name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("rack_role", id)); c.str(data, "name") != name {
		t.Fatalf("get rack role name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_rack_role", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("rack_role", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get rack role after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_rack_role", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "rack_role", "id": id})
}

func deviceRoleCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-device-role-" + runID

	created := c.call("create_device_role", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created device role name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("device_role", id)); c.str(data, "name") != name {
		t.Fatalf("get device role name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_device_role", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("device_role", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get device role after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_device_role", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "device_role", "id": id})
}

func platformCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-platform-" + runID

	created := c.call("create_platform", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created platform name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("platform", id)); c.str(data, "name") != name {
		t.Fatalf("get platform name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_platform", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("platform", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get platform after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_platform", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "platform", "id": id})
}

func virtualChassisCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-vc-" + runID

	// Virtual chassis has no slug field; name is optional but provided here.
	created := c.call("create_virtual_chassis", map[string]interface{}{"name": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created virtual chassis name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("virtual_chassis", id)); c.str(data, "name") != name {
		t.Fatalf("get virtual chassis name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_virtual_chassis", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("virtual_chassis", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get virtual chassis after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_virtual_chassis", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "virtual_chassis", "id": id})
}

func powerPanelCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-power-panel-" + runID

	// Power panel has no slug field; it needs a site on create.
	siteID := createSite(t, c, name+"-site")

	created := c.call("create_power_panel", map[string]interface{}{"name": name, "site": siteID})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created power panel name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("power_panel", id)); c.str(data, "name") != name {
		t.Fatalf("get power panel name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_power_panel", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("power_panel", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get power panel after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_power_panel", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "power_panel", "id": id})
}

func inventoryItemRoleCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-inv-role-" + runID

	created := c.call("create_inventory_item_role", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created inventory item role name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("inventory_item_role", id)); c.str(data, "name") != name {
		t.Fatalf("get inventory item role name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_inventory_item_role", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("inventory_item_role", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get inventory item role after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_inventory_item_role", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "inventory_item_role", "id": id})
}
