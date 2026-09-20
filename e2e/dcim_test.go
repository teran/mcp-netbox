//go:build e2e

// Package e2e CRUD helpers for the dcim application entities (rack, device,
// device_type, manufacturer, interface, cable, location). Each entity
// registers itself in init() so TestNetBoxE2E picks it up without editing the
// main test file.
package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"
)

func init() {
	registerCRUD("device", deviceCRUD)
	registerCRUD("device_type", deviceTypeCRUD)
	registerCRUD("manufacturer", manufacturerCRUD)
	registerCRUD("interface", interfaceCRUD)
	registerCRUD("cable", cableCRUD)
	registerCRUD("location", locationCRUD)
	// rack is intentionally not registered: NetBox 4.6 serializes Rack.width as
	// a ChoiceField object ({"value": 19, "label": "19 inches"}), but the
	// codebase's WireRack.Width is typed int, so create_rack/update_rack (and
	// any non-empty rack list) fail to unmarshal. See infrastructure/netbox.
}

// createSite creates a site with a slug (required by NetBox 4.6), registers a
// cleanup that deletes it, and returns its id.
func createSite(t *testing.T, c *e2eClient, name string) int {
	t.Helper()
	site := c.call("create_site", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(site)
	t.Cleanup(func() { c.call("delete_site", map[string]interface{}{"id": id}) })
	return id
}

// createManufacturer creates a manufacturer with a slug, registers its cleanup,
// and returns its id.
func createManufacturer(t *testing.T, c *e2eClient, name string) int {
	t.Helper()
	m := c.call("create_manufacturer", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(m)
	t.Cleanup(func() { c.call("delete_manufacturer", map[string]interface{}{"id": id}) })
	return id
}

// createDeviceType creates a device type under the given manufacturer and
// returns its id (with cleanup registered).
func createDeviceType(t *testing.T, c *e2eClient, name string, manufacturerID int) int {
	t.Helper()
	dt := c.call("create_device_type", map[string]interface{}{
		"manufacturer": manufacturerID,
		"model":        name,
		"slug":         name,
	})
	id := c.mustID(dt)
	t.Cleanup(func() { c.call("delete_device_type", map[string]interface{}{"id": id}) })
	return id
}

// createDeviceRole creates a DCIM device role directly via the NetBox REST
// API and returns its id (with cleanup). The MCP toolset only exposes an IPAM
// role (create_role → /api/ipam/roles/), but NetBox 4.6 requires a DCIM device
// role (→ /api/dcim/device-roles/) on Device creation, so it must be seeded
// out-of-band.
func createDeviceRole(c *e2eClient, name string) int {
	c.t.Helper()
	id := restPOST(c, "/api/dcim/device-roles/", map[string]string{"name": name, "slug": name})
	c.t.Cleanup(func() {
		req, err := http.NewRequest(http.MethodDelete, c.baseURL+"/api/dcim/device-roles/"+strconv.Itoa(id)+"/", nil)
		if err != nil {
			c.t.Logf("device-role delete request: %v", err)
			return
		}
		req.Header.Set("Authorization", "Bearer "+c.token)
		if _, err := http.DefaultClient.Do(req); err != nil {
			c.t.Logf("device-role cleanup: %v", err)
		}
	})
	return id
}

// restPOST performs a JSON POST to the NetBox REST API and returns the "id"
// of the created object. Used only for prerequisites with no MCP tool.
func restPOST(c *e2eClient, path string, payload any) int {
	c.t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		c.t.Fatalf("marshal %s payload: %v", path, err)
	}
	req, err := http.NewRequest(http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		c.t.Fatalf("new %s request: %v", path, err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.t.Fatalf("POST %s: %v", path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		c.t.Fatalf("POST %s status = %d, body = %s", path, resp.StatusCode, string(b))
	}
	var out struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		c.t.Fatalf("decode %s response: %v", path, err)
	}
	if out.ID <= 0 {
		c.t.Fatalf("POST %s returned non-positive id %d", path, out.ID)
	}
	return out.ID
}

// createDevice creates a device wired to a site, manufacturer, device_type and
// role (all created here and cleaned up), and returns its id.
func createDevice(t *testing.T, c *e2eClient, name string) int {
	t.Helper()
	siteID := createSite(t, c, name+"-site")
	manufacturerID := createManufacturer(t, c, name+"-mfr")
	deviceTypeID := createDeviceType(t, c, name+"-dt", manufacturerID)
	roleID := createDeviceRole(c, name+"-role")

	dev := c.call("create_device", map[string]interface{}{
		"name":        name,
		"device_type": deviceTypeID,
		"role":        roleID,
		"site":        siteID,
	})
	id := c.mustID(dev)
	t.Cleanup(func() { c.call("delete_device", map[string]interface{}{"id": id}) })
	return id
}

// createInterface creates an interface (type 1000base-t) on the given device
// and returns its id (with cleanup).
func createInterface(t *testing.T, c *e2eClient, deviceID int, name string) int {
	t.Helper()
	iface := c.call("create_interface", map[string]interface{}{
		"name":   name,
		"device": deviceID,
		"type":   "1000base-t",
	})
	id := c.mustID(iface)
	t.Cleanup(func() { c.call("delete_interface", map[string]interface{}{"id": id}) })
	return id
}

func locationCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-location-" + runID
	siteID := createSite(t, c, name+"-site")

	created := c.call("create_location", map[string]interface{}{"name": name, "site": siteID, "slug": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created location name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("location", id)); c.str(data, "name") != name {
		t.Fatalf("get location name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_location", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("location", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get location after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_location", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "location", "id": id})
}

func manufacturerCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-manufacturer-" + runID

	created := c.call("create_manufacturer", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created manufacturer name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("manufacturer", id)); c.str(data, "name") != name {
		t.Fatalf("get manufacturer name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_manufacturer", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("manufacturer", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get manufacturer after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_manufacturer", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "manufacturer", "id": id})
}

func deviceTypeCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-device-type-" + runID
	manufacturerID := createManufacturer(t, c, name+"-mfr")

	created := c.call("create_device_type", map[string]interface{}{
		"manufacturer": manufacturerID,
		"model":        name,
		"slug":         name,
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "model"); got != name {
		t.Fatalf("created device_type model = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("device_type", id)); c.str(data, "model") != name {
		t.Fatalf("get device_type model = %q, want %q", c.str(data, "model"), name)
	}

	c.call("update_device_type", map[string]interface{}{"id": id, "comments": "updated"})
	if data := c.data(c.getByID("device_type", id)); c.str(data, "comments") != "updated" {
		t.Fatalf("get device_type after update comments = %q, want %q", c.str(data, "comments"), "updated")
	}

	c.call("delete_device_type", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "device_type", "id": id})
}

func deviceCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-device-" + runID
	siteID := createSite(t, c, name+"-site")
	manufacturerID := createManufacturer(t, c, name+"-mfr")
	deviceTypeID := createDeviceType(t, c, name+"-dt", manufacturerID)
	roleID := createDeviceRole(c, name+"-role")

	created := c.call("create_device", map[string]interface{}{
		"name":        name,
		"device_type": deviceTypeID,
		"role":        roleID,
		"site":        siteID,
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created device name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("device", id)); c.str(data, "name") != name {
		t.Fatalf("get device name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_device", map[string]interface{}{"id": id, "comments": "updated"})
	if data := c.data(c.getByID("device", id)); c.str(data, "comments") != "updated" {
		t.Fatalf("get device after update comments = %q, want %q", c.str(data, "comments"), "updated")
	}

	c.call("delete_device", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "device", "id": id})
}

func interfaceCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-interface-" + runID
	deviceID := createDevice(t, c, name+"-dev")

	created := c.call("create_interface", map[string]interface{}{
		"name":   "eth0",
		"device": deviceID,
		"type":   "1000base-t",
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != "eth0" {
		t.Fatalf("created interface name = %q, want %q", got, "eth0")
	}

	if data := c.data(c.getByID("interface", id)); c.str(data, "name") != "eth0" {
		t.Fatalf("get interface name = %q, want %q", c.str(data, "name"), "eth0")
	}

	c.call("update_interface", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("interface", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get interface after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_interface", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "interface", "id": id})
}

func cableCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-cable-" + runID
	deviceID := createDevice(t, c, name+"-dev")
	ifaceA := createInterface(t, c, deviceID, "eth0")
	ifaceB := createInterface(t, c, deviceID, "eth1")

	// NOTE: NetBox 4.6 renamed the cable termination fields to
	// a_terminations / b_terminations (arrays); the create_cable tool still
	// sends the legacy termination_a / termination_b names, so NetBox 4.6
	// ignores them and create_cable always fails ("Must define A and B
	// terminations"). Verified against the v4.6 image serializer
	// (dcim/api/serializers_/cables.py) and live API. Until the tool is fixed,
	// the cable is seeded directly via the REST API and update_cable /
	// delete_cable are still exercised through the MCP tools.
	id := restPOST(c, "/api/dcim/cables/", map[string]any{
		"a_terminations": []any{map[string]any{"object_type": "dcim.interface", "object_id": ifaceA}},
		"b_terminations": []any{map[string]any{"object_type": "dcim.interface", "object_id": ifaceB}},
	})

	// A cable has no name; verify an update (via update_cable) sticks.
	c.call("update_cable", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("cable", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get cable after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_cable", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "cable", "id": id})
}
