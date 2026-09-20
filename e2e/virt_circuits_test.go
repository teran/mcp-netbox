//go:build e2e

// Package e2e contains the end-to-end NetBox test. This file registers CRUD
// coverage for the virtualization (cluster, cluster_type, cluster_group,
// vm_interface) and circuits (provider, circuit_type, circuit_termination)
// entities. Each helper follows the shared registry pattern: it is registered
// in init() and exercised by TestNetBoxE2E against a real NetBox.
package e2e

import "testing"

// init registers the virtualization + circuits entity CRUD helpers. They run in
// alphabetical order because TestNetBoxE2E sorts the global registry by name.
func init() {
	registerCRUD("circuit_termination", circuitTerminationCRUD)
	registerCRUD("circuit_type", circuitTypeCRUD)
	registerCRUD("cluster", clusterCRUD)
	registerCRUD("cluster_group", clusterGroupCRUD)
	registerCRUD("cluster_type", clusterTypeCRUD)
	registerCRUD("provider", providerCRUD)
	registerCRUD("vm_interface", vmInterfaceCRUD)
}

// providerCRUD is standalone: a provider requires only a name. NetBox 4.6
// requires an explicit slug (it is not auto-derived from the name).
func providerCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-provider-" + runID

	created := c.call("create_provider", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created provider name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("provider", id)); c.str(data, "name") != name {
		t.Fatalf("get provider name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_provider", map[string]interface{}{"id": id, "comments": "updated"})
	if data := c.data(c.getByID("provider", id)); c.str(data, "comments") != "updated" {
		t.Fatalf("get provider after update comments = %q", c.str(data, "comments"))
	}

	c.call("delete_provider", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "provider", "id": id})
}

// circuitTypeCRUD is standalone: a circuit type requires only a name (slug
// explicitly in NetBox 4.6).
func circuitTypeCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-ctype-" + runID

	created := c.call("create_circuit_type", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created circuit_type name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("circuit_type", id)); c.str(data, "name") != name {
		t.Fatalf("get circuit_type name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_circuit_type", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("circuit_type", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get circuit_type after update description = %q", c.str(data, "description"))
	}

	c.call("delete_circuit_type", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "circuit_type", "id": id})
}

// clusterTypeCRUD is standalone: a cluster type requires only a name (slug
// explicitly in NetBox 4.6).
func clusterTypeCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-ctype-" + runID

	created := c.call("create_cluster_type", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created cluster_type name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("cluster_type", id)); c.str(data, "name") != name {
		t.Fatalf("get cluster_type name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_cluster_type", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("cluster_type", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get cluster_type after update description = %q", c.str(data, "description"))
	}

	c.call("delete_cluster_type", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "cluster_type", "id": id})
}

// clusterGroupCRUD is standalone: a cluster group requires only a name (slug
// explicitly in NetBox 4.6).
func clusterGroupCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-cgroup-" + runID

	created := c.call("create_cluster_group", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created cluster_group name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("cluster_group", id)); c.str(data, "name") != name {
		t.Fatalf("get cluster_group name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_cluster_group", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("cluster_group", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get cluster_group after update description = %q", c.str(data, "description"))
	}

	c.call("delete_cluster_group", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "cluster_group", "id": id})
}

// clusterCRUD requires a cluster type, so create one as a prerequisite and clean
// it up afterwards.
func clusterCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()

	ctName := "e2e-clustertype-" + runID
	clusterType := c.call("create_cluster_type", map[string]interface{}{"name": ctName, "slug": ctName})
	clusterTypeID := c.mustID(clusterType)
	t.Cleanup(func() { c.call("delete_cluster_type", map[string]interface{}{"id": clusterTypeID}) })

	name := "e2e-cluster-" + runID
	created := c.call("create_cluster", map[string]interface{}{"name": name, "type": clusterTypeID})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created cluster name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("cluster", id)); c.str(data, "name") != name {
		t.Fatalf("get cluster name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_cluster", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("cluster", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get cluster after update description = %q", c.str(data, "description"))
	}

	c.call("delete_cluster", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "cluster", "id": id})
}

// vmInterfaceCRUD requires a virtual machine, which in turn requires a site (or
// cluster). Create a site, then a VM, as prerequisites and clean them up.
func vmInterfaceCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()

	// A virtual machine must be assigned to a site, cluster, or device.
	siteName := "e2e-vmisite-" + runID
	site := c.call("create_site", map[string]interface{}{"name": siteName, "slug": siteName})
	siteID := c.mustID(site)
	t.Cleanup(func() { c.call("delete_site", map[string]interface{}{"id": siteID}) })

	vmName := "e2e-vmi-vm-" + runID
	vm := c.call("create_virtual_machine", map[string]interface{}{"name": vmName, "site": siteID})
	vmID := c.mustID(vm)
	t.Cleanup(func() { c.call("delete_virtual_machine", map[string]interface{}{"id": vmID}) })

	name := "e2e-vmiface-" + runID
	created := c.call("create_vm_interface", map[string]interface{}{"name": name, "virtual_machine": vmID})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created vm_interface name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("vm_interface", id)); c.str(data, "name") != name {
		t.Fatalf("get vm_interface name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_vm_interface", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("vm_interface", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get vm_interface after update description = %q", c.str(data, "description"))
	}

	c.call("delete_vm_interface", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "vm_interface", "id": id})
}

// circuitTerminationCRUD requires a circuit (which itself needs a provider and a
// circuit type) and a site, plus the term_side ("A" or "Z"). Create all of them
// as prerequisites and clean them up.
func circuitTerminationCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()

	// circuit prerequisites: provider + circuit_type.
	provider := c.call("create_provider", map[string]interface{}{"name": "e2e-ct-prov-" + runID, "slug": "e2e-ct-prov-" + runID})
	providerID := c.mustID(provider)
	t.Cleanup(func() { c.call("delete_provider", map[string]interface{}{"id": providerID}) })

	circuitType := c.call("create_circuit_type", map[string]interface{}{"name": "e2e-ct-type-" + runID, "slug": "e2e-ct-type-" + runID})
	circuitTypeID := c.mustID(circuitType)
	t.Cleanup(func() { c.call("delete_circuit_type", map[string]interface{}{"id": circuitTypeID}) })

	circuit := c.call("create_circuit", map[string]interface{}{
		"cid":          "e2e-ct-cid-" + runID,
		"provider":     providerID,
		"circuit_type": circuitTypeID,
	})
	circuitID := c.mustID(circuit)
	t.Cleanup(func() { c.call("delete_circuit", map[string]interface{}{"id": circuitID}) })

	// A circuit termination must be associated with a site.
	siteName := "e2e-ctsite-" + runID
	site := c.call("create_site", map[string]interface{}{"name": siteName, "slug": siteName})
	siteID := c.mustID(site)
	t.Cleanup(func() { c.call("delete_site", map[string]interface{}{"id": siteID}) })

	created := c.call("create_circuit_termination", map[string]interface{}{
		"circuit":   circuitID,
		"site":      siteID,
		"term_side": "A",
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "term_side"); got != "A" {
		t.Fatalf("created circuit_termination term_side = %q, want %q", got, "A")
	}

	if data := c.data(c.getByID("circuit_termination", id)); c.str(data, "term_side") != "A" {
		t.Fatalf("get circuit_termination term_side = %q, want %q", c.str(data, "term_side"), "A")
	}

	c.call("update_circuit_termination", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("circuit_termination", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get circuit_termination after update description = %q", c.str(data, "description"))
	}

	c.call("delete_circuit_termination", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "circuit_termination", "id": id})
}
