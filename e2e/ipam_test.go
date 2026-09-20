//go:build e2e

package e2e

import (
	"fmt"
	"testing"
)

// This file extends the e2e CRUD registry with the IPAM and misc entities that
// have write tools (create_/update_/delete_) but no dedicated read tool:
// vlan, vlan_group, vrf, role, and contact. Each helper follows the same
// registry pattern as siteCRUD/prefixCRUD in netbox_e2e_test.go: create ->
// fetch via getByID -> update -> fetch -> delete -> expect not-found.

func init() {
	registerCRUD("vlan", vlanCRUD)
	registerCRUD("vlan_group", vlanGroupCRUD)
	registerCRUD("vrf", vrfCRUD)
	registerCRUD("role", roleCRUD)
	registerCRUD("contact", contactCRUD)
}

// vlanCRUD covers create_vlan/update_vlan/delete_vlan. A VLAN requires a site
// or a vlan_group in NetBox 4.6, so a site is created as a prerequisite (and
// cleaned up via t.Cleanup) and referenced by the VLAN.
func vlanCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()

	vid := 1000 + hashSuffix(runID)%2000

	// Prerequisite: a site to host the VLAN.
	siteName := "e2e-vlansite-" + runID
	site := c.call("create_site", map[string]interface{}{"name": siteName, "slug": siteName})
	siteID := c.mustID(site)
	t.Cleanup(func() { c.call("delete_site", map[string]interface{}{"id": siteID}) })

	name := "e2e-vlan-" + runID
	created := c.call("create_vlan", map[string]interface{}{"vid": vid, "name": name, "site": siteID})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created vlan name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("vlan", id)); c.str(data, "name") != name {
		t.Fatalf("get vlan name = %q, want %q", c.str(data, "name"), name)
	} else if got, ok := num(data["vid"]); !ok || got != vid {
		t.Fatalf("get vlan vid = %v, want %d", data["vid"], vid)
	}

	c.call("update_vlan", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("vlan", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get vlan after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_vlan", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "vlan", "id": id})
}

// vlanGroupCRUD covers create_vlan_group/update_vlan_group/delete_vlan_group.
// In NetBox 4.6 a VLAN group must be scoped to a site, so a site is created as
// a prerequisite and passed via scope_type/scope_id.
func vlanGroupCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()

	// Prerequisite: a site to scope the VLAN group.
	siteName := "e2e-vgsite-" + runID
	site := c.call("create_site", map[string]interface{}{"name": siteName, "slug": siteName})
	siteID := c.mustID(site)
	t.Cleanup(func() { c.call("delete_site", map[string]interface{}{"id": siteID}) })

	name := "e2e-vlangroup-" + runID
	created := c.call("create_vlan_group", map[string]interface{}{
		"name":       name,
		"slug":       name,
		"scope_type": "dcim.site",
		"scope_id":   siteID,
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created vlan_group name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("vlan_group", id)); c.str(data, "name") != name {
		t.Fatalf("get vlan_group name = %q, want %q", c.str(data, "name"), name)
	} else if got := c.str(data, "slug"); got != name {
		t.Fatalf("get vlan_group slug = %q, want %q", got, name)
	}

	c.call("update_vlan_group", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("vlan_group", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get vlan_group after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_vlan_group", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "vlan_group", "id": id})
}

// vrfCRUD covers create_vrf/update_vrf/delete_vrf. A VRF is standalone and
// requires name and rd on create.
func vrfCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()

	name := "e2e-vrf-" + runID
	rd := fmt.Sprintf("65000:%d", 100+hashSuffix(runID)%89900)
	created := c.call("create_vrf", map[string]interface{}{"name": name, "rd": rd})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created vrf name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("vrf", id)); c.str(data, "name") != name {
		t.Fatalf("get vrf name = %q, want %q", c.str(data, "name"), name)
	} else if got := c.str(data, "rd"); got != rd {
		t.Fatalf("get vrf rd = %q, want %q", got, rd)
	}

	c.call("update_vrf", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("vrf", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get vrf after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_vrf", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "vrf", "id": id})
}

// roleCRUD covers create_role/update_role/delete_role. An IPAM role is
// standalone and requires name; NetBox 4.6 also requires an explicit slug.
func roleCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()

	name := "e2e-role-" + runID
	created := c.call("create_role", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created role name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("role", id)); c.str(data, "name") != name {
		t.Fatalf("get role name = %q, want %q", c.str(data, "name"), name)
	} else if got := c.str(data, "slug"); got != name {
		t.Fatalf("get role slug = %q, want %q", got, name)
	}

	c.call("update_role", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("role", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get role after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_role", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "role", "id": id})
}

// contactCRUD covers create_contact/update_contact/delete_contact. A tenancy
// contact is standalone and requires only name on create.
func contactCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()

	name := "e2e-contact-" + runID
	email := "e2e-" + runID + "@example.test"
	created := c.call("create_contact", map[string]interface{}{"name": name, "email": email})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created contact name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("contact", id)); c.str(data, "name") != name {
		t.Fatalf("get contact name = %q, want %q", c.str(data, "name"), name)
	} else if got := c.str(data, "email"); got != email {
		t.Fatalf("get contact email = %q, want %q", got, email)
	}

	c.call("update_contact", map[string]interface{}{"id": id, "email": "updated@example.test"})
	if data := c.data(c.getByID("contact", id)); c.str(data, "email") != "updated@example.test" {
		t.Fatalf("get contact after update email = %q, want %q", c.str(data, "email"), "updated@example.test")
	}

	c.call("delete_contact", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "contact", "id": id})
}
