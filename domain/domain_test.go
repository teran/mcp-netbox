package domain

import (
	"context"
	"encoding/json"
	"testing"
)

func TestSite_Populate(t *testing.T) {
	site := Site{
		ID:     1,
		Name:   "Test Site",
		Slug:   "test-site",
		Status: &Label{Value: "active", Label: "Active"},
	}
	if site.ID != 1 {
		t.Errorf("ID = %d, want %d", site.ID, 1)
	}
	if site.Name != "Test Site" {
		t.Errorf("Name = %q, want %q", site.Name, "Test Site")
	}
	if site.Slug != "test-site" {
		t.Errorf("Slug = %q, want %q", site.Slug, "test-site")
	}
	if site.Status.Value != "active" {
		t.Errorf("Status.Value = %q, want %q", site.Status.Value, "active")
	}
}

func TestSite_ZeroValues(t *testing.T) {
	var site Site
	if site.ID != 0 {
		t.Errorf("ID = %d, want 0", site.ID)
	}
	if site.Name != "" {
		t.Errorf("Name = %q, want empty", site.Name)
	}
	if site.Status != nil {
		t.Errorf("Status = %v, want nil", site.Status)
	}
}

func TestDevice_Populate(t *testing.T) {
	device := Device{
		ID:         1,
		Name:       "router1",
		DeviceType: &Nested{ID: 2, Name: "MX480"},
		Site:       &Nested{ID: 3, Name: "Main DC"},
		Status:     &Label{Value: "active", Label: "Active"},
	}
	if device.ID != 1 {
		t.Errorf("ID = %d, want %d", device.ID, 1)
	}
	if device.Name != "router1" {
		t.Errorf("Name = %q, want %q", device.Name, "router1")
	}
	if device.DeviceType.Name != "MX480" {
		t.Errorf("DeviceType.Name = %q, want %q", device.DeviceType.Name, "MX480")
	}
}

func TestDevice_ZeroValues(t *testing.T) {
	var device Device
	if device.ID != 0 {
		t.Errorf("ID = %d, want 0", device.ID)
	}
	if device.Name != "" {
		t.Errorf("Name = %q, want empty", device.Name)
	}
}

func TestIPAddress_Populate(t *testing.T) {
	ipAddr := IPAddress{
		ID:      1,
		Address: "10.0.0.1/24",
		Status:  &Label{Value: "active", Label: "Active"},
		DNSName: "server.example.com",
	}
	if ipAddr.Address != "10.0.0.1/24" {
		t.Errorf("Address = %q, want %q", ipAddr.Address, "10.0.0.1/24")
	}
	if ipAddr.DNSName != "server.example.com" {
		t.Errorf("DNSName = %q, want %q", ipAddr.DNSName, "server.example.com")
	}
}

func TestIPAddress_ZeroValues(t *testing.T) {
	var ipAddr IPAddress
	if ipAddr.Address != "" {
		t.Errorf("Address = %q, want empty", ipAddr.Address)
	}
}

func TestPrefix_Populate(t *testing.T) {
	prefix := Prefix{
		ID:     1,
		Prefix: "10.0.0.0/8",
		Site:   &Nested{ID: 1, Name: "DC1"},
		Status: &Label{Value: "active", Label: "Active"},
		Family: &Family{Value: 4, Label: "IPv4"},
	}
	if prefix.Prefix != "10.0.0.0/8" {
		t.Errorf("Prefix = %q, want %q", prefix.Prefix, "10.0.0.0/8")
	}
	if prefix.Family.Value != 4 {
		t.Errorf("Family.Value = %d, want %d", prefix.Family.Value, 4)
	}
}

func TestPrefix_ZeroValues(t *testing.T) {
	var prefix Prefix
	if prefix.Prefix != "" {
		t.Errorf("Prefix = %q, want empty", prefix.Prefix)
	}
}

func TestVLAN_Populate(t *testing.T) {
	vlan := VLAN{
		ID:   1,
		VID:  100,
		Name: "Users",
		Site: &Nested{ID: 1, Name: "DC1"},
	}
	if vlan.VID != 100 {
		t.Errorf("VID = %d, want %d", vlan.VID, 100)
	}
	if vlan.Name != "Users" {
		t.Errorf("Name = %q, want %q", vlan.Name, "Users")
	}
}

func TestVLAN_ZeroValues(t *testing.T) {
	var vlan VLAN
	if vlan.VID != 0 {
		t.Errorf("VID = %d, want 0", vlan.VID)
	}
}

func TestVirtualMachine_Populate(t *testing.T) {
	virtualMachine := VirtualMachine{
		ID:      1,
		Name:    "web-01",
		Cluster: &Nested{ID: 1, Name: "Prod Cluster"},
		VCPUs:   4,
		Memory:  16384,
		Disk:    100,
	}
	if virtualMachine.Name != "web-01" {
		t.Errorf("Name = %q, want %q", virtualMachine.Name, "web-01")
	}
	if virtualMachine.VCPUs != 4 {
		t.Errorf("VCPUs = %f, want %f", virtualMachine.VCPUs, 4.0)
	}
}

func TestVirtualMachine_ZeroValues(t *testing.T) {
	var virtualMachine VirtualMachine
	if virtualMachine.VCPUs != 0 {
		t.Errorf("VCPUs = %f, want 0", virtualMachine.VCPUs)
	}
}

func TestCluster_Populate(t *testing.T) {
	cluster := Cluster{
		ID:          1,
		Name:        "Production",
		ClusterType: &Nested{Name: "VMware"},
		Site:        &Nested{Name: "DC1"},
	}
	if cluster.Name != "Production" {
		t.Errorf("Name = %q, want %q", cluster.Name, "Production")
	}
	if cluster.ClusterType.Name != "VMware" {
		t.Errorf("ClusterType.Name = %q, want %q", cluster.ClusterType.Name, "VMware")
	}
}

func TestCluster_ZeroValues(t *testing.T) {
	var cluster Cluster
	if cluster.ID != 0 {
		t.Errorf("ID = %d, want 0", cluster.ID)
	}
}

func TestCircuit_Populate(t *testing.T) {
	circuit := Circuit{
		ID:          1,
		CID:         "CIR-001",
		Provider:    &Nested{Name: "ATT"},
		CircuitType: &Nested{Name: "Dark Fiber"},
		Status:      &Label{Value: "active", Label: "Active"},
	}
	if circuit.CID != "CIR-001" {
		t.Errorf("CID = %q, want %q", circuit.CID, "CIR-001")
	}
	if circuit.Provider.Name != "ATT" {
		t.Errorf("Provider.Name = %q, want %q", circuit.Provider.Name, "ATT")
	}
}

func TestCircuit_ZeroValues(t *testing.T) {
	var circuit Circuit
	if circuit.ID != 0 {
		t.Errorf("ID = %d, want 0", circuit.ID)
	}
}

func TestRack_Populate(t *testing.T) {
	rack := Rack{
		ID:      1,
		Name:    "Rack-01",
		Site:    &Nested{Name: "DC1"},
		UHeight: 42,
		Width:   19,
	}
	if rack.Name != "Rack-01" {
		t.Errorf("Name = %q, want %q", rack.Name, "Rack-01")
	}
	if rack.UHeight != 42 {
		t.Errorf("UHeight = %d, want %d", rack.UHeight, 42)
	}
}

func TestRack_ZeroValues(t *testing.T) {
	var rack Rack
	if rack.UHeight != 0 {
		t.Errorf("UHeight = %d, want 0", rack.UHeight)
	}
}

func TestInterface_Populate(t *testing.T) {
	iface := Interface{
		ID:         1,
		Name:       "ge-0/0/0",
		Device:     &Nested{ID: 1, Name: "router1"},
		Type:       &Label{Value: "1000base-t", Label: "1000BASE-T"},
		Enabled:    true,
		MTU:        1500,
		MACAddress: "aa:bb:cc:dd:ee:ff",
		Speed:      1000000,
	}
	if iface.Name != "ge-0/0/0" {
		t.Errorf("Name = %q, want %q", iface.Name, "ge-0/0/0")
	}
	if iface.Device.Name != "router1" {
		t.Errorf("Device.Name = %q, want %q", iface.Device.Name, "router1")
	}
	if iface.Type.Value != "1000base-t" {
		t.Errorf("Type.Value = %q, want %q", iface.Type.Value, "1000base-t")
	}
	if !iface.Enabled {
		t.Error("Enabled = false, want true")
	}
}

func TestInterface_ZeroValues(t *testing.T) {
	var iface Interface
	if iface.ID != 0 {
		t.Errorf("ID = %d, want 0", iface.ID)
	}
	if iface.Enabled {
		t.Error("Enabled = true, want false")
	}
}

func TestVMInterface_Populate(t *testing.T) {
	vmi := VMInterface{
		ID:             1,
		Name:           "eth0",
		VirtualMachine: &Nested{ID: 1, Name: "web-01"},
		Enabled:        true,
		MTU:            1500,
		MACAddress:     "11:22:33:44:55:66",
	}
	if vmi.Name != "eth0" {
		t.Errorf("Name = %q, want %q", vmi.Name, "eth0")
	}
	if vmi.VirtualMachine.Name != "web-01" {
		t.Errorf("VirtualMachine.Name = %q, want %q", vmi.VirtualMachine.Name, "web-01")
	}
	if !vmi.Enabled {
		t.Error("Enabled = false, want true")
	}
}

func TestVMInterface_ZeroValues(t *testing.T) {
	var vmi VMInterface
	if vmi.MTU != 0 {
		t.Errorf("MTU = %d, want 0", vmi.MTU)
	}
}

func TestCircuitTermination_Populate(t *testing.T) {
	circuitTermination := CircuitTermination{
		ID:            1,
		Circuit:       &Nested{ID: 1, Name: "CIR-001"},
		TermSide:      "A",
		Site:          &Nested{ID: 1, Name: "DC1"},
		Speed:         10000000,
		UpstreamSpeed: 10000000,
	}
	if circuitTermination.TermSide != "A" {
		t.Errorf("TermSide = %q, want %q", circuitTermination.TermSide, "A")
	}
	if circuitTermination.Circuit.Name != "CIR-001" {
		t.Errorf("Circuit.Name = %q, want %q", circuitTermination.Circuit.Name, "CIR-001")
	}
	if circuitTermination.Speed != 10000000 {
		t.Errorf("Speed = %d, want %d", circuitTermination.Speed, 10000000)
	}
}

func TestCircuitTermination_ZeroValues(t *testing.T) {
	var circuitTermination CircuitTermination
	if circuitTermination.TermSide != "" {
		t.Errorf("TermSide = %q, want empty", circuitTermination.TermSide)
	}
}

func TestCable_Populate(t *testing.T) {
	cable := Cable{
		ID:         1,
		Type:       &Label{Value: "cat6a", Label: "CAT6a"},
		Status:     &Label{Value: "connected", Label: "Connected"},
		Label:      "link-01",
		Color:      "blue",
		Length:     10.5,
		LengthUnit: &Label{Value: "m", Label: "Meters"},
	}
	if cable.Label != "link-01" {
		t.Errorf("Label = %q, want %q", cable.Label, "link-01")
	}
	if cable.Type.Value != "cat6a" {
		t.Errorf("Type.Value = %q, want %q", cable.Type.Value, "cat6a")
	}
	if cable.Status.Value != "connected" {
		t.Errorf("Status.Value = %q, want %q", cable.Status.Value, "connected")
	}
	if cable.Length != 10.5 {
		t.Errorf("Length = %f, want %f", cable.Length, 10.5)
	}
}

func TestCable_ZeroValues(t *testing.T) {
	var cable Cable
	if cable.Label != "" {
		t.Errorf("Label = %q, want empty", cable.Label)
	}
	if cable.Type != nil {
		t.Errorf("Type = %+v, want nil", cable.Type)
	}
}

func TestCableWrite_Populate(t *testing.T) {
	label := "link-01"
	length := 10.5
	cableWrite := CableWrite{
		TerminationA: &CableTerminationWrite{ObjectType: "dcim.interface", ObjectID: 101},
		TerminationB: &CableTerminationWrite{ObjectType: "dcim.interface", ObjectID: 102},
		Label:        &label,
		Length:       &length,
	}
	if cableWrite.TerminationA.ObjectType != "dcim.interface" || cableWrite.TerminationA.ObjectID != 101 {
		t.Errorf("TerminationA = %+v, want dcim.interface 101", cableWrite.TerminationA)
	}
	if cableWrite.TerminationB.ObjectID != 102 {
		t.Errorf("TerminationB = %+v, want object_id 102", cableWrite.TerminationB)
	}
	if cableWrite.Label == nil || *cableWrite.Label != "link-01" ||
		cableWrite.Length == nil || *cableWrite.Length != 10.5 {
		t.Errorf("cableWrite = %+v, want label link-01 length 10.5", cableWrite)
	}
}

func TestCableWrite_ZeroValues(t *testing.T) {
	var cableWrite CableWrite
	if cableWrite.TerminationA != nil || cableWrite.TerminationB != nil {
		t.Errorf("cableWrite = %+v, want nil terminations", cableWrite)
	}
}

func TestVMInterfaceWrite_Populate(t *testing.T) {
	vmID := 5
	mtu := 1500
	viw := VMInterfaceWrite{VirtualMachine: &vmID, Name: "eth0", MTU: &mtu}
	if viw.Name != "eth0" || viw.VirtualMachine == nil || *viw.VirtualMachine != 5 || viw.MTU == nil || *viw.MTU != 1500 {
		t.Errorf("viw = %+v, want name eth0 virtual_machine 5 mtu 1500", viw)
	}
}

func TestVMInterfaceWrite_ZeroValues(t *testing.T) {
	var viw VMInterfaceWrite
	if viw.Name != "" || viw.VirtualMachine != nil {
		t.Errorf("viw = %+v, want empty name and nil virtual_machine", viw)
	}
}

func TestProvider_Populate(t *testing.T) {
	provider := Provider{
		ID:           1,
		Name:         "ACME",
		Slug:         "acme",
		Asn:          64512,
		Account:      "ACC-1",
		PortalURL:    "https://portal.example.com",
		NocContact:   "noc@example.com",
		AdminContact: "admin@example.com",
		Comments:     "comments",
	}
	if provider.Name != "ACME" || provider.Asn != 64512 || provider.PortalURL != "https://portal.example.com" {
		t.Errorf("provider = %+v, want name ACME asn 64512 portal url", provider)
	}
}

func TestProvider_ZeroValues(t *testing.T) {
	var provider Provider
	if provider.Name != "" || provider.Asn != 0 {
		t.Errorf("provider = %+v, want empty name and zero asn", provider)
	}
}

func TestProviderWrite_Populate(t *testing.T) {
	asn := 64512
	portal := "https://portal.example.com"
	providerWrite := ProviderWrite{Name: "ACME", Asn: &asn, PortalURL: &portal}
	if providerWrite.Name != "ACME" || providerWrite.Asn == nil ||
		*providerWrite.Asn != 64512 || providerWrite.PortalURL == nil || *providerWrite.PortalURL != portal {
		t.Errorf("providerWrite = %+v, want name ACME asn 64512 portal url", providerWrite)
	}
}

func TestProviderWrite_ZeroValues(t *testing.T) {
	var providerWrite ProviderWrite
	if providerWrite.Name != "" || providerWrite.Asn != nil {
		t.Errorf("providerWrite = %+v, want empty name and nil asn", providerWrite)
	}
}

func TestTenant_Populate(t *testing.T) {
	tenant := Tenant{
		ID:          1,
		Name:        "ACME",
		Slug:        "acme",
		Description: "desc",
		Comments:    "comments",
	}
	if tenant.Name != "ACME" || tenant.Slug != "acme" || tenant.Description != "desc" {
		t.Errorf("tenant = %+v, want name ACME slug acme description desc", tenant)
	}
}

func TestTenant_ZeroValues(t *testing.T) {
	var tenant Tenant
	if tenant.Name != "" || tenant.Slug != "" {
		t.Errorf("tenant = %+v, want empty name and slug", tenant)
	}
}

func TestTenantWrite_Populate(t *testing.T) {
	slug := "acme"
	desc := "desc"
	tenantWrite := TenantWrite{Name: "ACME", Slug: &slug, Description: &desc}
	if tenantWrite.Name != "ACME" || tenantWrite.Slug == nil || *tenantWrite.Slug != "acme" {
		t.Errorf("tenantWrite = %+v, want name ACME slug acme", tenantWrite)
	}
}

func TestTenantWrite_ZeroValues(t *testing.T) {
	var tenantWrite TenantWrite
	if tenantWrite.Name != "" || tenantWrite.Slug != nil {
		t.Errorf("tenantWrite = %+v, want empty name and nil slug", tenantWrite)
	}
}

func TestManufacturer_Populate(t *testing.T) {
	manufacturer := Manufacturer{
		ID:          1,
		Name:        "Cisco",
		Slug:        "cisco",
		Description: "desc",
	}
	if manufacturer.Name != "Cisco" || manufacturer.Slug != "cisco" || manufacturer.Description != "desc" {
		t.Errorf("manufacturer = %+v, want name Cisco slug cisco description desc", manufacturer)
	}
}

func TestManufacturer_ZeroValues(t *testing.T) {
	var manufacturer Manufacturer
	if manufacturer.Name != "" || manufacturer.Slug != "" {
		t.Errorf("manufacturer = %+v, want empty name and slug", manufacturer)
	}
}

func TestManufacturerWrite_Populate(t *testing.T) {
	slug := "cisco"
	manufacturerWrite := ManufacturerWrite{Name: "Cisco", Slug: &slug}
	if manufacturerWrite.Name != "Cisco" || manufacturerWrite.Slug == nil || *manufacturerWrite.Slug != "cisco" {
		t.Errorf("manufacturerWrite = %+v, want name Cisco slug cisco", manufacturerWrite)
	}
}

func TestManufacturerWrite_ZeroValues(t *testing.T) {
	var manufacturerWrite ManufacturerWrite
	if manufacturerWrite.Name != "" || manufacturerWrite.Slug != nil {
		t.Errorf("manufacturerWrite = %+v, want empty name and nil slug", manufacturerWrite)
	}
}

func TestDeviceType_Populate(t *testing.T) {
	deviceType := DeviceType{
		ID:         1,
		Model:      "C9300-24T",
		PartNumber: "C9300-24T",
		UHeight:    1,
	}
	if deviceType.Model != "C9300-24T" || deviceType.PartNumber != "C9300-24T" || deviceType.UHeight != 1 {
		t.Errorf("deviceType = %+v, want model C9300-24T", deviceType)
	}
}

func TestDeviceType_ZeroValues(t *testing.T) {
	var deviceType DeviceType
	if deviceType.Model != "" || deviceType.Manufacturer != nil {
		t.Errorf("deviceType = %+v, want empty model and nil manufacturer", deviceType)
	}
}

func TestDeviceTypeWrite_Populate(t *testing.T) {
	manufacturerID := 3
	uHeight := 1.0
	dtw := DeviceTypeWrite{Manufacturer: &manufacturerID, Model: "C9300-24T", UHeight: &uHeight}
	if dtw.Model != "C9300-24T" || dtw.Manufacturer == nil || *dtw.Manufacturer != 3 {
		t.Errorf("dtw = %+v, want model C9300-24T manufacturer 3", dtw)
	}
}

func TestDeviceTypeWrite_ZeroValues(t *testing.T) {
	var dtw DeviceTypeWrite
	if dtw.Model != "" || dtw.Manufacturer != nil {
		t.Errorf("dtw = %+v, want empty model and nil manufacturer", dtw)
	}
}

func TestNested_Populate(t *testing.T) {
	nested := Nested{
		ID:   1,
		Name: "Test",
		Slug: "test",
		URL:  "http://netbox/api/test/1/",
	}
	if nested.ID != 1 {
		t.Errorf("ID = %d, want %d", nested.ID, 1)
	}
	if nested.Name != "Test" {
		t.Errorf("Name = %q, want %q", nested.Name, "Test")
	}
}

func TestNested_ZeroValues(t *testing.T) {
	var nested Nested
	if nested.ID != 0 {
		t.Errorf("ID = %d, want 0", nested.ID)
	}
}

func TestLabel_Populate(t *testing.T) {
	label := Label{Value: "active", Label: "Active"}
	if label.Value != "active" {
		t.Errorf("Value = %q, want %q", label.Value, "active")
	}
	if label.Label != "Active" {
		t.Errorf("Label = %q, want %q", label.Label, "Active")
	}
}

func TestLabel_ZeroValues(t *testing.T) {
	var label Label
	if label.Value != "" {
		t.Errorf("Value = %q, want empty", label.Value)
	}
}

func TestPaginatedResponse_Empty(t *testing.T) {
	resp := PaginatedResponse[Site]{}
	if resp.Count != 0 {
		t.Errorf("Count = %d, want 0", resp.Count)
	}
	if len(resp.Results) != 0 {
		t.Errorf("len(Results) = %d, want 0", len(resp.Results))
	}
}

func TestPaginatedResponse_Populated(t *testing.T) {
	resp := PaginatedResponse[Device]{
		Count: 1,
		Results: []Device{
			{ID: 1, Name: "device1"},
		},
	}
	if resp.Count != 1 {
		t.Errorf("Count = %d, want %d", resp.Count, 1)
	}
	if len(resp.Results) != 1 {
		t.Errorf("len(Results) = %d, want %d", len(resp.Results), 1)
	}
	if resp.Results[0].Name != "device1" {
		t.Errorf("Results[0].Name = %q, want %q", resp.Results[0].Name, "device1")
	}
}

type mockRepo struct {
	listSitesFunc func(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[Site], error)
}

func (repo *mockRepo) ListSites(
	ctx context.Context, token string, params map[string]string,
) (*PaginatedResponse[Site], error) {
	return repo.listSitesFunc(ctx, token, params)
}

func TestNetworkRepositoryInterface(t *testing.T) {
	t.Parallel()

	t.Run("mock repository returns expected response", func(t *testing.T) {
		repo := &mockRepo{
			listSitesFunc: func(_ context.Context, _ string, params map[string]string) (*PaginatedResponse[Site], error) {
				return &PaginatedResponse[Site]{
					Count:   1,
					Results: []Site{{ID: 1, Name: params["name"]}},
				}, nil
			},
		}

		resp, err := repo.ListSites(context.Background(), "token", map[string]string{"name": "Test"})
		if err != nil {
			t.Fatalf("ListSites() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("Count = %d, want %d", resp.Count, 1)
		}
		if resp.Results[0].Name != "Test" {
			t.Errorf("Results[0].Name = %q, want %q", resp.Results[0].Name, "Test")
		}
	})
}

func TestRawObject_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   RawObject
		want string
	}{
		{"object", RawObject(`{"id":1,"name":"x"}`), `{"id":1,"name":"x"}`},
		{"array", RawObject(`[1,2,3]`), `[1,2,3]`},
		{"null", RawObject(`null`), `null`},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := json.Marshal(testCase.in)
			if err != nil {
				t.Fatalf("MarshalJSON() error: %v", err)
			}
			if string(got) != testCase.want {
				t.Errorf("MarshalJSON() = %s, want %s", got, testCase.want)
			}
		})
	}
}

func TestRawObject_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		var out RawObject
		if err := json.Unmarshal([]byte(`{"a":1}`), &out); err != nil {
			t.Fatalf("UnmarshalJSON() error: %v", err)
		}
		if string(out) != `{"a":1}` {
			t.Errorf("UnmarshalJSON() = %s, want %s", out, `{"a":1}`)
		}
	})

	t.Run("invalid json returns error", func(t *testing.T) {
		var out RawObject
		if err := json.Unmarshal([]byte(`{invalid`), &out); err == nil {
			t.Error("UnmarshalJSON() = nil, want error")
		}
	})
}

func TestRawObject_EmbeddedInStruct(t *testing.T) {
	t.Parallel()

	type wrapper struct {
		Name string    `json:"name"`
		Data RawObject `json:"data"`
	}

	w := wrapper{Name: "n", Data: RawObject(`{"x":1}`)}
	data, err := json.Marshal(w)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	// The raw object must be inlined as JSON, not base64.
	if string(data) != `{"name":"n","data":{"x":1}}` {
		t.Errorf("Marshal() = %s, want inlined JSON", data)
	}
}
