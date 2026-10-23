//go:build e2e

// Package e2e CRUD helpers for the DCIM port/bay entities added in
// feature/expand-entities Batch 2: front_port, rear_port, console_port,
// console_server_port, power_port, power_outlet, power_feed, device_bay,
// rack_reservation, inventory_item, mac_address. Each registers itself in
// init() so TestNetBoxE2E picks it up.
package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func init() {
	// These entities are covered by e2e CRUD helpers below. front_port and
	// rear_port use a device (+ a backing rear_port for front_port), and
	// mac_address needs a device interface.
	registerCRUD("front_port", frontPortCRUD)
	registerCRUD("rear_port", rearPortCRUD)
	registerCRUD("mac_address", macAddressCRUD)
	registerCRUD("console_port", consolePortCRUD)
	registerCRUD("console_server_port", consoleServerPortCRUD)
	registerCRUD("power_port", powerPortCRUD)
	registerCRUD("power_outlet", powerOutletCRUD)
	registerCRUD("power_feed", powerFeedCRUD)
	registerCRUD("device_bay", deviceBayCRUD)
	registerCRUD("rack_reservation", rackReservationCRUD)
	registerCRUD("inventory_item", inventoryItemCRUD)
}

// createParentDevice creates a device whose device_type has subdevice_role
// "Parent", which NetBox requires before a device bay can be added. Returns
// the device id (with cleanup registered).
func createParentDevice(t *testing.T, c *e2eClient, name string) int {
	t.Helper()
	siteID := createSite(t, c, name+"-site")
	manufacturerID := createManufacturer(t, c, name+"-mfr")
	dt := c.call("create_device_type", map[string]interface{}{
		"manufacturer":   manufacturerID,
		"model":          name + "-dt",
		"slug":           name + "-dt",
		"subdevice_role": "parent",
	})
	dtID := c.mustID(dt)
	t.Cleanup(func() { c.call("delete_device_type", map[string]interface{}{"id": dtID}) })
	roleID := createDeviceRole(c, name+"-role")

	dev := c.call("create_device", map[string]interface{}{
		"name":        name,
		"device_type": dtID,
		"role":        roleID,
		"site":        siteID,
	})
	id := c.mustID(dev)
	t.Cleanup(func() { c.call("delete_device", map[string]interface{}{"id": id}) })
	return id
}

func frontPortCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-fp-" + runID
	deviceID := createDevice(t, c, name+"-dev")

	// A front port must be backed by a rear port on the same device.
	rear := c.call("create_rear_port", map[string]interface{}{
		"device":    deviceID,
		"name":      "rp1",
		"type":      "lc",
		"positions": 1,
	})
	rearID := c.mustID(rear)
	t.Cleanup(func() { c.call("delete_rear_port", map[string]interface{}{"id": rearID}) })

	created := c.call("create_front_port", map[string]interface{}{
		"device":             deviceID,
		"name":               "fp0",
		"type":               "lc",
		"rear_port":          rearID,
		"rear_port_position": 1,
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != "fp0" {
		t.Fatalf("created front port name = %q, want %q", got, "fp0")
	}

	if data := c.data(c.getByID("front_port", id)); c.str(data, "name") != "fp0" {
		t.Fatalf("get front port name = %q, want %q", c.str(data, "name"), "fp0")
	}

	c.call("update_front_port", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("front_port", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get front port after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_front_port", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "front_port", "id": id})
}

func rearPortCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-rp-" + runID
	deviceID := createDevice(t, c, name+"-dev")

	created := c.call("create_rear_port", map[string]interface{}{
		"device":    deviceID,
		"name":      "rp0",
		"type":      "lc",
		"positions": 1,
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != "rp0" {
		t.Fatalf("created rear port name = %q, want %q", got, "rp0")
	}

	if data := c.data(c.getByID("rear_port", id)); c.str(data, "name") != "rp0" {
		t.Fatalf("get rear port name = %q, want %q", c.str(data, "name"), "rp0")
	}

	c.call("update_rear_port", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("rear_port", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get rear port after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_rear_port", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "rear_port", "id": id})
}

func consolePortCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-cp-" + runID
	deviceID := createDevice(t, c, name+"-dev")

	created := c.call("create_console_port", map[string]interface{}{
		"device": deviceID,
		"name":   "con0",
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != "con0" {
		t.Fatalf("created console port name = %q, want %q", got, "con0")
	}

	if data := c.data(c.getByID("console_port", id)); c.str(data, "name") != "con0" {
		t.Fatalf("get console port name = %q, want %q", c.str(data, "name"), "con0")
	}

	c.call("update_console_port", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("console_port", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get console port after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_console_port", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "console_port", "id": id})
}

func consoleServerPortCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-csp-" + runID
	deviceID := createDevice(t, c, name+"-dev")

	created := c.call("create_console_server_port", map[string]interface{}{
		"device": deviceID,
		"name":   "csp0",
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != "csp0" {
		t.Fatalf("created console server port name = %q, want %q", got, "csp0")
	}

	if data := c.data(c.getByID("console_server_port", id)); c.str(data, "name") != "csp0" {
		t.Fatalf("get console server port name = %q, want %q", c.str(data, "name"), "csp0")
	}

	c.call("update_console_server_port", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("console_server_port", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get console server port after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_console_server_port", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "console_server_port", "id": id})
}

func powerPortCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-pp-" + runID
	deviceID := createDevice(t, c, name+"-dev")

	created := c.call("create_power_port", map[string]interface{}{
		"device": deviceID,
		"name":   "psu0",
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != "psu0" {
		t.Fatalf("created power port name = %q, want %q", got, "psu0")
	}

	if data := c.data(c.getByID("power_port", id)); c.str(data, "name") != "psu0" {
		t.Fatalf("get power port name = %q, want %q", c.str(data, "name"), "psu0")
	}

	c.call("update_power_port", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("power_port", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get power port after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_power_port", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "power_port", "id": id})
}

func powerOutletCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-po-" + runID
	deviceID := createDevice(t, c, name+"-dev")

	created := c.call("create_power_outlet", map[string]interface{}{
		"device": deviceID,
		"name":   "out0",
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != "out0" {
		t.Fatalf("created power outlet name = %q, want %q", got, "out0")
	}

	if data := c.data(c.getByID("power_outlet", id)); c.str(data, "name") != "out0" {
		t.Fatalf("get power outlet name = %q, want %q", c.str(data, "name"), "out0")
	}

	c.call("update_power_outlet", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("power_outlet", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get power outlet after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_power_outlet", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "power_outlet", "id": id})
}

func powerFeedCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-pf-" + runID
	siteID := createSite(t, c, name+"-site")

	panel := c.call("create_power_panel", map[string]interface{}{"name": name + "-panel", "site": siteID})
	panelID := c.mustID(panel)
	t.Cleanup(func() { c.call("delete_power_panel", map[string]interface{}{"id": panelID}) })

	created := c.call("create_power_feed", map[string]interface{}{
		"name":        "feed0",
		"power_panel": panelID,
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != "feed0" {
		t.Fatalf("created power feed name = %q, want %q", got, "feed0")
	}

	if data := c.data(c.getByID("power_feed", id)); c.str(data, "name") != "feed0" {
		t.Fatalf("get power feed name = %q, want %q", c.str(data, "name"), "feed0")
	}

	c.call("update_power_feed", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("power_feed", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get power feed after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_power_feed", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "power_feed", "id": id})
}

func deviceBayCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-db-" + runID
	deviceID := createParentDevice(t, c, name+"-dev")

	created := c.call("create_device_bay", map[string]interface{}{
		"device": deviceID,
		"name":   "bay0",
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != "bay0" {
		t.Fatalf("created device bay name = %q, want %q", got, "bay0")
	}

	if data := c.data(c.getByID("device_bay", id)); c.str(data, "name") != "bay0" {
		t.Fatalf("get device bay name = %q, want %q", c.str(data, "name"), "bay0")
	}

	c.call("update_device_bay", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("device_bay", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get device bay after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_device_bay", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "device_bay", "id": id})
}

func rackReservationCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-rr-" + runID
	siteID := createSite(t, c, name+"-site")

	// create_rack fails against NetBox 4.6 because WireRack.Width is typed int
	// while NetBox serializes width as a ChoiceField object, so seed the rack
	// directly via the REST API (see the note in e2e/dcim_test.go).
	rackID := restPOST(c, "/api/dcim/racks/", map[string]any{"name": name + "-rack", "site": siteID})
	t.Cleanup(func() {
		req, err := http.NewRequest(http.MethodDelete, c.baseURL+"/api/dcim/racks/"+strconv.Itoa(rackID)+"/", nil)
		if err != nil {
			c.t.Logf("rack delete request: %v", err)
			return
		}
		req.Header.Set("Authorization", "Bearer "+c.token)
		if _, err := http.DefaultClient.Do(req); err != nil {
			c.t.Logf("rack cleanup: %v", err)
		}
	})

	// Rack reservations require a user; reuse the token's (super)user.
	userID := firstUserID(c)

	created := c.call("create_rack_reservation", map[string]interface{}{
		"description": "reserved",
		"rack":        rackID,
		"units":       []int{1, 2},
		"user":        userID,
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "description"); got != "reserved" {
		t.Fatalf("created rack reservation description = %q, want %q", got, "reserved")
	}

	if data := c.data(c.getByID("rack_reservation", id)); c.str(data, "description") != "reserved" {
		t.Fatalf("get rack reservation description = %q, want %q", c.str(data, "description"), "reserved")
	}

	c.call("update_rack_reservation", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("rack_reservation", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get rack reservation after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_rack_reservation", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "rack_reservation", "id": id})
}

func inventoryItemCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-ii-" + runID
	deviceID := createDevice(t, c, name+"-dev")

	created := c.call("create_inventory_item", map[string]interface{}{
		"device": deviceID,
		"name":   "PSU-1",
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != "PSU-1" {
		t.Fatalf("created inventory item name = %q, want %q", got, "PSU-1")
	}

	if data := c.data(c.getByID("inventory_item", id)); c.str(data, "name") != "PSU-1" {
		t.Fatalf("get inventory item name = %q, want %q", c.str(data, "name"), "PSU-1")
	}

	c.call("update_inventory_item", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("inventory_item", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get inventory item after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_inventory_item", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "inventory_item", "id": id})
}

func macAddressCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-mac-" + runID
	deviceID := createDevice(t, c, name+"-dev")
	ifaceID := createInterface(t, c, deviceID, "eth0")

	mac := fmt.Sprintf(
		"02:00:00:%02x:%02x:%02x",
		hashSuffix(runID)%256, (hashSuffix(runID)>>8)%256, (hashSuffix(runID)>>16)%256,
	)

	created := c.call("create_mac_address", map[string]interface{}{
		"address":              mac,
		"assigned_object_type": "dcim.interface",
		"assigned_object_id":   ifaceID,
	})
	id := c.mustID(created)
	// NetBox normalizes MAC addresses to uppercase hex, so compare case-insensitively.
	if got := c.str(c.data(created), "address"); !strings.EqualFold(got, mac) {
		t.Fatalf("created mac address = %q, want %q", got, mac)
	}

	// get_object_by_id returns the raw NetBox JSON, where NetBox 4.6 uses the
	// JSON key "mac_address" (not "address"). The create response above is the
	// typed domain model, which marshals Address under "address".
	if data := c.data(c.getByID("mac_address", id)); !strings.EqualFold(c.str(data, "mac_address"), mac) {
		t.Fatalf("get mac address = %q, want %q", c.str(data, "mac_address"), mac)
	}

	c.call("update_mac_address", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("mac_address", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get mac address after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_mac_address", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "mac_address", "id": id})
}

// firstUserID returns the id of the first NetBox user (used as the required
// `user` FK on rack reservations).
func firstUserID(c *e2eClient) int {
	c.t.Helper()
	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/api/users/users/?limit=1", nil)
	if err != nil {
		c.t.Fatalf("users request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.t.Fatalf("GET users: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.t.Fatalf("GET users status = %d", resp.StatusCode)
	}
	var out struct {
		Results []struct {
			ID int `json:"id"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		c.t.Fatalf("decode users response: %v", err)
	}
	if len(out.Results) == 0 {
		c.t.Fatalf("no users found")
	}
	return out.Results[0].ID
}
