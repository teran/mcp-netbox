//go:build e2e

package e2e

// CRUD helpers for the entities added in feature/expand-entities Batch 4:
// contact_group, contact_role, virtual_disk, virtual_machine_type. Each
// registers itself in init() so TestNetBoxE2E picks it up.
//
// virtual_disk requires a virtual_machine (created as a prerequisite);
// contact_group can nest under a parent contact group.

import (
	"fmt"
	"testing"
)

func init() {
	registerCRUD("contact_group", contactGroupCRUD)
	registerCRUD("contact_role", contactRoleCRUD)
	registerCRUD("virtual_disk", virtualDiskCRUD)
	registerCRUD("virtual_machine_type", virtualMachineTypeCRUD)
}

// createContactGroup creates a contact group (optionally under a parent),
// registers its cleanup and returns its id.
func createContactGroup(t *testing.T, c *e2eClient, name string, parent int) int {
	t.Helper()
	args := map[string]interface{}{"name": name, "slug": name}
	if parent != 0 {
		args["parent"] = parent
	}
	r := c.call("create_contact_group", args)
	id := c.mustID(r)
	t.Cleanup(func() { c.call("delete_contact_group", map[string]interface{}{"id": id}) })
	return id
}

// createVirtualMachine creates a VM hosted on a site, registers its cleanup and
// returns its id.
func createVirtualMachine(t *testing.T, c *e2eClient, name string) int {
	t.Helper()
	siteID := createSite(t, c, name+"-site")
	vm := c.call("create_virtual_machine", map[string]interface{}{"name": name, "site": siteID})
	id := c.mustID(vm)
	t.Cleanup(func() { c.call("delete_virtual_machine", map[string]interface{}{"id": id}) })
	return id
}

// contactGroupCRUD covers create_contact_group/update_contact_group/
// delete_contact_group. A contact group can nest under a parent.
func contactGroupCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-cg-" + runID
	parentID := createContactGroup(t, c, name+"-parent", 0)

	created := c.call("create_contact_group", map[string]interface{}{
		"name": name, "slug": name, "parent": parentID,
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created contact_group name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("contact_group", id)); c.str(data, "name") != name {
		t.Fatalf("get contact_group name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_contact_group", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("contact_group", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get contact_group after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_contact_group", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "contact_group", "id": id})
}

// contactRoleCRUD covers create_contact_role/update_contact_role/
// delete_contact_role. A contact role is standalone.
func contactRoleCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-cr-" + runID

	created := c.call("create_contact_role", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created contact_role name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("contact_role", id)); c.str(data, "name") != name {
		t.Fatalf("get contact_role name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_contact_role", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("contact_role", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get contact_role after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_contact_role", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "contact_role", "id": id})
}

// virtualDiskCRUD covers create_virtual_disk/update_virtual_disk/
// delete_virtual_disk. A virtual disk requires a virtual machine.
func virtualDiskCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-vd-" + runID
	vmID := createVirtualMachine(t, c, name+"-vm")

	created := c.call("create_virtual_disk", map[string]interface{}{
		"name": name, "virtual_machine": vmID, "size": 100,
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created virtual_disk name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("virtual_disk", id)); c.str(data, "name") != name {
		t.Fatalf("get virtual_disk name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_virtual_disk", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("virtual_disk", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get virtual_disk after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_virtual_disk", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "virtual_disk", "id": id})
}

// virtualMachineTypeCRUD covers create_virtual_machine_type/
// update_virtual_machine_type/delete_virtual_machine_type. A virtual machine
// type is standalone.
func virtualMachineTypeCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-vmt-" + runID
	// NetBox 4.6 requires an explicit slug for virtual machine types.
	slug := fmt.Sprintf("vmt-%s", runID)

	created := c.call("create_virtual_machine_type", map[string]interface{}{"name": name, "slug": slug})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created virtual_machine_type name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("virtual_machine_type", id)); c.str(data, "name") != name {
		t.Fatalf("get virtual_machine_type name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_virtual_machine_type", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("virtual_machine_type", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get virtual_machine_type after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_virtual_machine_type", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "virtual_machine_type", "id": id})
}
