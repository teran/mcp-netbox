//go:build e2e

package e2e

// CRUD helpers for the circuits entities added in feature/expand-entities
// Batch 5: provider_account, provider_network, circuit_group,
// circuit_group_assignment, virtual_circuit, virtual_circuit_type,
// virtual_circuit_termination. Each registers itself in init() so
// TestNetBoxE2E picks it up.
//
// Prerequisites are seeded via the MCP tools themselves:
//   - provider_account / provider_network need a provider
//   - circuit_group_assignment needs a circuit_group and a circuit
//   - virtual_circuit needs a virtual_circuit_type and a provider_network
//   - virtual_circuit_termination needs a virtual_circuit and an interface
//     (dcim.interface)

import "testing"

func init() {
	registerCRUD("provider_account", providerAccountCRUD)
	registerCRUD("provider_network", providerNetworkCRUD)
	registerCRUD("circuit_group", circuitGroupCRUD)
	registerCRUD("circuit_group_assignment", circuitGroupAssignmentCRUD)
	registerCRUD("virtual_circuit", virtualCircuitCRUD)
	registerCRUD("virtual_circuit_type", virtualCircuitTypeCRUD)
	registerCRUD("virtual_circuit_termination", virtualCircuitTerminationCRUD)
}

// createProvider creates a provider with an explicit slug (required by NetBox
// 4.6), registers its cleanup and returns its id.
func createProvider(t *testing.T, c *e2eClient, name string) int {
	t.Helper()
	p := c.call("create_provider", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(p)
	t.Cleanup(func() { c.call("delete_provider", map[string]interface{}{"id": id}) })
	return id
}

// createCircuitGroup creates a circuit group (slug required in NetBox 4.6) and
// returns its id (with cleanup).
func createCircuitGroup(t *testing.T, c *e2eClient, name string) int {
	t.Helper()
	g := c.call("create_circuit_group", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(g)
	t.Cleanup(func() { c.call("delete_circuit_group", map[string]interface{}{"id": id}) })
	return id
}

// createVirtualCircuitType creates a virtual circuit type (slug required in
// NetBox 4.6) and returns its id (with cleanup).
func createVirtualCircuitType(t *testing.T, c *e2eClient, name string) int {
	t.Helper()
	vct := c.call("create_virtual_circuit_type", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(vct)
	t.Cleanup(func() { c.call("delete_virtual_circuit_type", map[string]interface{}{"id": id}) })
	return id
}

// createProviderNetwork creates a provider network under the given provider and
// returns its id (with cleanup).
func createProviderNetwork(t *testing.T, c *e2eClient, name string, providerID int) int {
	t.Helper()
	pn := c.call("create_provider_network", map[string]interface{}{"name": name, "provider": providerID})
	id := c.mustID(pn)
	t.Cleanup(func() { c.call("delete_provider_network", map[string]interface{}{"id": id}) })
	return id
}

// createVirtualCircuit creates a virtual circuit with a type, provider network
// and status, and returns its id (with cleanup).
func createVirtualCircuit(t *testing.T, c *e2eClient, name string, typeID, providerNetworkID int) int {
	t.Helper()
	vc := c.call("create_virtual_circuit", map[string]interface{}{
		"cid": name, "type": typeID, "provider_network": providerNetworkID, "status": "planned",
	})
	id := c.mustID(vc)
	t.Cleanup(func() { c.call("delete_virtual_circuit", map[string]interface{}{"id": id}) })
	return id
}

// createVirtualInterface seeds a virtual (is_virtual=true) DCIM interface on the
// given device via the NetBox REST API and returns its id. Virtual circuit
// terminations may only attach to virtual interfaces, and the create_interface
// tool does not expose is_virtual, so it is seeded out-of-band.
func createVirtualInterface(t *testing.T, c *e2eClient, deviceID int, name string) int {
	t.Helper()
	return restPOST(c, "/api/dcim/interfaces/", map[string]interface{}{
		"name": name, "device": deviceID, "type": "virtual", "is_virtual": true,
	})
}

// providerAccountCRUD covers create_provider_account/update_provider_account/
// delete_provider_account. It requires a provider, created as a prerequisite.
func providerAccountCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-pacct-" + runID
	providerID := createProvider(t, c, name+"-prov")

	created := c.call("create_provider_account", map[string]interface{}{
		"name": name, "provider": providerID, "account": "acc-1",
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created provider_account name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("provider_account", id)); c.str(data, "name") != name {
		t.Fatalf("get provider_account name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_provider_account", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("provider_account", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get provider_account after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_provider_account", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "provider_account", "id": id})
}

// providerNetworkCRUD covers create_provider_network/update_provider_network/
// delete_provider_network. It requires a provider, created as a prerequisite.
func providerNetworkCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-pnet-" + runID
	providerID := createProvider(t, c, name+"-prov")

	created := c.call("create_provider_network", map[string]interface{}{"name": name, "provider": providerID})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created provider_network name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("provider_network", id)); c.str(data, "name") != name {
		t.Fatalf("get provider_network name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_provider_network", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("provider_network", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get provider_network after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_provider_network", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "provider_network", "id": id})
}

// circuitGroupCRUD covers create_circuit_group/update_circuit_group/
// delete_circuit_group. A circuit group is standalone.
func circuitGroupCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-cgrp-" + runID

	created := c.call("create_circuit_group", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created circuit_group name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("circuit_group", id)); c.str(data, "name") != name {
		t.Fatalf("get circuit_group name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_circuit_group", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("circuit_group", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get circuit_group after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_circuit_group", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "circuit_group", "id": id})
}

// circuitGroupAssignmentCRUD covers create_circuit_group_assignment/
// update_circuit_group_assignment/delete_circuit_group_assignment. It requires
// a circuit group and a circuit (provider + circuit_type), both created as
// prerequisites.
func circuitGroupAssignmentCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-casg-" + runID
	groupID := createCircuitGroup(t, c, name+"-g")

	providerID := createProvider(t, c, name+"-prov")
	circuitType := c.call("create_circuit_type", map[string]interface{}{
		"name": name + "-ct", "slug": name + "-ct",
	})
	circuitTypeID := c.mustID(circuitType)
	t.Cleanup(func() { c.call("delete_circuit_type", map[string]interface{}{"id": circuitTypeID}) })

	circuit := c.call("create_circuit", map[string]interface{}{
		"cid": name + "-cid", "provider": providerID, "circuit_type": circuitTypeID,
	})
	circuitID := c.mustID(circuit)
	t.Cleanup(func() { c.call("delete_circuit", map[string]interface{}{"id": circuitID}) })

	created := c.call("create_circuit_group_assignment", map[string]interface{}{
		"group": groupID, "member_type": "circuits.circuit", "member_id": circuitID, "priority": "primary",
	})
	id := c.mustID(created)

	// An assignment has no name; verify it can be fetched and updated.
	if data := c.data(c.getByID("circuit_group_assignment", id)); len(data) == 0 {
		t.Fatalf("get circuit_group_assignment returned empty data")
	}

	c.call("update_circuit_group_assignment", map[string]interface{}{"id": id, "priority": "secondary"})
	if data := c.data(c.getByID("circuit_group_assignment", id)); len(data) == 0 {
		t.Fatalf("get circuit_group_assignment after update returned empty data")
	}

	c.call("delete_circuit_group_assignment", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "circuit_group_assignment", "id": id})
}

// virtualCircuitTypeCRUD covers create_virtual_circuit_type/
// update_virtual_circuit_type/delete_virtual_circuit_type. A virtual circuit
// type is standalone.
func virtualCircuitTypeCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-vct-" + runID

	created := c.call("create_virtual_circuit_type", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created virtual_circuit_type name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("virtual_circuit_type", id)); c.str(data, "name") != name {
		t.Fatalf("get virtual_circuit_type name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_virtual_circuit_type", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("virtual_circuit_type", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get virtual_circuit_type after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_virtual_circuit_type", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "virtual_circuit_type", "id": id})
}

// virtualCircuitCRUD covers create_virtual_circuit/update_virtual_circuit/
// delete_virtual_circuit. It requires a virtual circuit type and a provider
// network (which needs a provider), all created as prerequisites.
func virtualCircuitCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-vc-" + runID
	typeID := createVirtualCircuitType(t, c, name+"-type")
	providerID := createProvider(t, c, name+"-prov")
	providerNetworkID := createProviderNetwork(t, c, name+"-pn", providerID)

	created := c.call("create_virtual_circuit", map[string]interface{}{
		"cid": name, "type": typeID, "provider_network": providerNetworkID, "status": "planned",
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "cid"); got != name {
		t.Fatalf("created virtual_circuit cid = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("virtual_circuit", id)); c.str(data, "cid") != name {
		t.Fatalf("get virtual_circuit cid = %q, want %q", c.str(data, "cid"), name)
	}

	c.call("update_virtual_circuit", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("virtual_circuit", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get virtual_circuit after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_virtual_circuit", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "virtual_circuit", "id": id})
}

// virtualCircuitTerminationCRUD covers create_virtual_circuit_termination/
// update_virtual_circuit_termination/delete_virtual_circuit_termination. It
// requires a virtual circuit and a device interface (interface_type
// "dcim.interface"), both created as prerequisites.
func virtualCircuitTerminationCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-vcterm-" + runID
	typeID := createVirtualCircuitType(t, c, name+"-type")
	providerID := createProvider(t, c, name+"-prov")
	providerNetworkID := createProviderNetwork(t, c, name+"-pn", providerID)
	vcID := createVirtualCircuit(t, c, name+"-vc", typeID, providerNetworkID)

	deviceID := createDevice(t, c, name+"-dev")
	ifaceID := createVirtualInterface(t, c, deviceID, "veth0")

	created := c.call("create_virtual_circuit_termination", map[string]interface{}{
		"virtual_circuit": vcID, "role": "peer", "interface": ifaceID,
	})
	id := c.mustID(created)

	// A termination has no name; verify it can be fetched and updated.
	if data := c.data(c.getByID("virtual_circuit_termination", id)); len(data) == 0 {
		t.Fatalf("get virtual_circuit_termination returned empty data")
	}

	c.call("update_virtual_circuit_termination", map[string]interface{}{"id": id, "role": "hub"})
	if data := c.data(c.getByID("virtual_circuit_termination", id)); len(data) == 0 {
		t.Fatalf("get virtual_circuit_termination after update returned empty data")
	}

	c.call("delete_virtual_circuit_termination", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "virtual_circuit_termination", "id": id})
}
