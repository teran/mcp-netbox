//go:build e2e

package e2e

// CRUD helpers for the IPAM entities added in feature/expand-entities Batch 3:
// rir, aggregate, ip_range, asn, asn_range, route_target, service,
// service_template, fhrp_group, fhrp_group_assignment, vlan_translation_policy,
// vlan_translation_rule. Each registers itself in init() so TestNetBoxE2E picks
// it up.
//
// Several entities require prerequisites that are seeded via the MCP tools
// themselves: aggregate needs an rir, asn_range needs an rir, service needs a
// device (or virtual machine), fhrp_group_assignment needs an fhrp_group and an
// ip_address, and vlan_translation_rule needs a policy. Those prerequisites are
// created (and cleaned up) inline.

import (
	"fmt"
	"testing"
)

func init() {
	registerCRUD("rir", rirCRUD)
	registerCRUD("aggregate", aggregateCRUD)
	registerCRUD("ip_range", ipRangeCRUD)
	registerCRUD("asn", asnCRUD)
	registerCRUD("asn_range", asnRangeCRUD)
	registerCRUD("route_target", routeTargetCRUD)
	registerCRUD("service", serviceCRUD)
	registerCRUD("service_template", serviceTemplateCRUD)
	registerCRUD("fhrp_group", fhrpGroupCRUD)
	registerCRUD("fhrp_group_assignment", fhrpGroupAssignmentCRUD)
	registerCRUD("vlan_translation_policy", vlanTranslationPolicyCRUD)
	registerCRUD("vlan_translation_rule", vlanTranslationRuleCRUD)
}

// createRIR creates an RIR via create_rir, registers its cleanup and returns
// its id. NetBox 4.6 requires an explicit slug.
func createRIR(t *testing.T, c *e2eClient, name string) int {
	t.Helper()
	r := c.call("create_rir", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(r)
	t.Cleanup(func() { c.call("delete_rir", map[string]interface{}{"id": id}) })
	return id
}

// rirCRUD covers create_rir/update_rir/delete_rir. An RIR is standalone.
func rirCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-rir-" + runID
	created := c.call("create_rir", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created rir name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("rir", id)); c.str(data, "name") != name {
		t.Fatalf("get rir name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_rir", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("rir", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get rir after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_rir", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "rir", "id": id})
}

// aggregateCRUD covers create_aggregate/update_aggregate/delete_aggregate. An
// aggregate requires an RIR in NetBox 4.6, so one is created as a prerequisite.
func aggregateCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-agg-" + runID
	rirID := createRIR(t, c, name+"-rir")

	prefix := fmt.Sprintf("10.%d.0.0/16", 100+hashSuffix(runID)%150)
	created := c.call("create_aggregate", map[string]interface{}{"prefix": prefix, "rir": rirID})
	id := c.mustID(created)
	if got := c.str(c.data(created), "prefix"); got != prefix {
		t.Fatalf("created aggregate prefix = %q, want %q", got, prefix)
	}

	if data := c.data(c.getByID("aggregate", id)); c.str(data, "prefix") != prefix {
		t.Fatalf("get aggregate prefix = %q, want %q", c.str(data, "prefix"), prefix)
	}

	c.call("update_aggregate", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("aggregate", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get aggregate after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_aggregate", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "aggregate", "id": id})
}

// ipRangeCRUD covers create_ip_range/update_ip_range/delete_ip_range. An IP
// range is standalone (vrf/tenant are optional), so it can be created bare.
func ipRangeCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	// Keep the second octet within valid IPv4 bounds (<=255) and disjoint from
	// the 192.168.x space used by ipAddressCRUD.
	base := 10 + hashSuffix(runID)%200
	start := fmt.Sprintf("10.%d.0.1", base)
	end := fmt.Sprintf("10.%d.0.10", base)
	// NetBox normalizes host addresses to include the /32 mask.
	startNormalized := start + "/32"
	endNormalized := end + "/32"

	created := c.call("create_ip_range", map[string]interface{}{"start_address": start, "end_address": end})
	id := c.mustID(created)
	if got := c.str(c.data(created), "start_address"); got != startNormalized {
		t.Fatalf("created ip_range start_address = %q, want %q", got, startNormalized)
	}

	if data := c.data(c.getByID("ip_range", id)); c.str(data, "start_address") != startNormalized {
		t.Fatalf("get ip_range start_address = %q, want %q", c.str(data, "start_address"), startNormalized)
	} else if c.str(data, "end_address") != endNormalized {
		t.Fatalf("get ip_range end_address = %q, want %q", c.str(data, "end_address"), endNormalized)
	}

	c.call("update_ip_range", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("ip_range", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get ip_range after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_ip_range", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "ip_range", "id": id})
}

// asnCRUD covers create_asn/update_asn/delete_asn. NetBox 4.6 requires an RIR
// on ASN creation, so one is created as a prerequisite.
func asnCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-asn-" + runID
	rirID := createRIR(t, c, name+"-rir")
	asn := 64512 + hashSuffix(runID)%1000

	created := c.call("create_asn", map[string]interface{}{"asn": asn, "rir": rirID})
	id := c.mustID(created)
	if got, ok := num(c.data(created)["asn"]); !ok || got != asn {
		t.Fatalf("created asn = %v, want %d", c.data(created)["asn"], asn)
	}

	if data := c.data(c.getByID("asn", id)); len(data) == 0 {
		t.Fatalf("get asn returned empty data: %v", data)
	} else if got, ok := num(data["asn"]); !ok || got != asn {
		t.Fatalf("get asn = %v, want %d", data["asn"], asn)
	}

	c.call("update_asn", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("asn", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get asn after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_asn", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "asn", "id": id})
}

// asnRangeCRUD covers create_asn_range/update_asn_range/delete_asn_range. An
// ASN range requires an RIR, so one is created as a prerequisite.
func asnRangeCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-asnrange-" + runID
	rirID := createRIR(t, c, name+"-rir")

	start := 64512 + hashSuffix(runID)%500
	end := start + 20
	created := c.call("create_asn_range", map[string]interface{}{
		"name": name, "slug": name, "start": start, "end": end, "rir": rirID,
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created asn_range name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("asn_range", id)); c.str(data, "name") != name {
		t.Fatalf("get asn_range name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_asn_range", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("asn_range", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get asn_range after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_asn_range", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "asn_range", "id": id})
}

// routeTargetCRUD covers create_route_target/update_route_target/
// delete_route_target. A route target is standalone.
func routeTargetCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := fmt.Sprintf("65000:%d", 100+hashSuffix(runID)%89900)
	created := c.call("create_route_target", map[string]interface{}{"name": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created route_target name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("route_target", id)); c.str(data, "name") != name {
		t.Fatalf("get route_target name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_route_target", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("route_target", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get route_target after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_route_target", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "route_target", "id": id})
}

// serviceCRUD covers create_service/update_service/delete_service. A service
// must be bound to a device or virtual machine, so a device is created as a
// prerequisite.
func serviceCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-svc-" + runID
	deviceID := createDevice(t, c, name+"-dev")

	created := c.call("create_service", map[string]interface{}{
		"name": name, "parent_object_type": "dcim.device", "parent_object_id": deviceID,
		"protocol": "tcp", "ports": []int{80, 443},
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created service name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("service", id)); c.str(data, "name") != name {
		t.Fatalf("get service name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_service", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("service", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get service after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_service", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "service", "id": id})
}

// serviceTemplateCRUD covers create_service_template/update_service_template/
// delete_service_template. A service template is standalone.
func serviceTemplateCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-svct-" + runID
	created := c.call("create_service_template", map[string]interface{}{
		"name": name, "protocol": "tcp", "ports": []int{443},
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created service_template name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("service_template", id)); c.str(data, "name") != name {
		t.Fatalf("get service_template name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_service_template", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("service_template", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get service_template after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_service_template", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "service_template", "id": id})
}

// fhrpGroupCRUD covers create_fhrp_group/update_fhrp_group/delete_fhrp_group.
// An FHRP group is standalone.
func fhrpGroupCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-fhrpg-" + runID
	created := c.call("create_fhrp_group", map[string]interface{}{
		"name": name, "protocol": "hsrp", "group_id": 10,
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created fhrp_group name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("fhrp_group", id)); c.str(data, "name") != name {
		t.Fatalf("get fhrp_group name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_fhrp_group", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("fhrp_group", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get fhrp_group after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_fhrp_group", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "fhrp_group", "id": id})
}

// fhrpGroupAssignmentCRUD covers create_fhrp_group_assignment/
// update_fhrp_group_assignment/delete_fhrp_group_assignment. An assignment
// requires an FHRP group and an interface (via interface_type/interface_id),
// both created as prerequisites.
func fhrpGroupAssignmentCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-fhrpa-" + runID
	group := c.call("create_fhrp_group", map[string]interface{}{
		"name": name + "-g", "protocol": "hsrp", "group_id": 20,
	})
	groupID := c.mustID(group)
	t.Cleanup(func() { c.call("delete_fhrp_group", map[string]interface{}{"id": groupID}) })

	deviceID := createDevice(t, c, name+"-dev")
	ifaceID := createInterface(t, c, deviceID, "eth0")

	created := c.call("create_fhrp_group_assignment", map[string]interface{}{
		"group": groupID, "interface_type": "dcim.interface", "interface_id": ifaceID, "priority": 100,
	})
	id := c.mustID(created)

	// An assignment has no name; verify it can be fetched and updated.
	if data := c.data(c.getByID("fhrp_group_assignment", id)); len(data) == 0 {
		t.Fatalf("get fhrp_group_assignment returned empty data")
	}

	c.call("update_fhrp_group_assignment", map[string]interface{}{"id": id, "priority": 200})
	if got, ok := num(c.data(c.getByID("fhrp_group_assignment", id))["priority"]); !ok || got != 200 {
		t.Fatalf("get fhrp_group_assignment priority after update = %v, want 200", got)
	}

	c.call("delete_fhrp_group_assignment", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "fhrp_group_assignment", "id": id})
}

// vlanTranslationPolicyCRUD covers create_vlan_translation_policy/
// update_vlan_translation_policy/delete_vlan_translation_policy. A policy is
// standalone.
func vlanTranslationPolicyCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-vtpol-" + runID
	created := c.call("create_vlan_translation_policy", map[string]interface{}{"name": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created vlan_translation_policy name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("vlan_translation_policy", id)); c.str(data, "name") != name {
		t.Fatalf("get vlan_translation_policy name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_vlan_translation_policy", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("vlan_translation_policy", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get vlan_translation_policy after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_vlan_translation_policy", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "vlan_translation_policy", "id": id})
}

// vlanTranslationRuleCRUD covers create_vlan_translation_rule/
// update_vlan_translation_rule/delete_vlan_translation_rule. A rule requires a
// policy, which is created as a prerequisite.
func vlanTranslationRuleCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-vtrule-" + runID
	pol := c.call("create_vlan_translation_policy", map[string]interface{}{"name": name + "-p"})
	polID := c.mustID(pol)
	t.Cleanup(func() { c.call("delete_vlan_translation_policy", map[string]interface{}{"id": polID}) })

	created := c.call("create_vlan_translation_rule", map[string]interface{}{
		"policy": polID, "local_vid": 100, "remote_vid": 200,
	})
	id := c.mustID(created)

	// A rule has no name; verify it can be fetched and updated (policy is
	// unchanged on update, which still exercises the PATCH path).
	if data := c.data(c.getByID("vlan_translation_rule", id)); len(data) == 0 {
		t.Fatalf("get vlan_translation_rule returned empty data")
	}

	c.call("update_vlan_translation_rule", map[string]interface{}{"id": id, "local_vid": 110})
	if data := c.data(c.getByID("vlan_translation_rule", id)); len(data) == 0 {
		t.Fatalf("get vlan_translation_rule after update returned empty data")
	}

	c.call("delete_vlan_translation_rule", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "vlan_translation_rule", "id": id})
}
