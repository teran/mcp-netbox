package domain

import (
	"context"
	"encoding/json"
	"testing"
)

func TestSite_Populate(t *testing.T) {
	s := Site{
		ID:     1,
		Name:   "Test Site",
		Slug:   "test-site",
		Status: &Label{Value: "active", Label: "Active"},
	}
	if s.ID != 1 {
		t.Errorf("ID = %d, want %d", s.ID, 1)
	}
	if s.Name != "Test Site" {
		t.Errorf("Name = %q, want %q", s.Name, "Test Site")
	}
	if s.Slug != "test-site" {
		t.Errorf("Slug = %q, want %q", s.Slug, "test-site")
	}
	if s.Status.Value != "active" {
		t.Errorf("Status.Value = %q, want %q", s.Status.Value, "active")
	}
}

func TestSite_ZeroValues(t *testing.T) {
	var s Site
	if s.ID != 0 {
		t.Errorf("ID = %d, want 0", s.ID)
	}
	if s.Name != "" {
		t.Errorf("Name = %q, want empty", s.Name)
	}
	if s.Status != nil {
		t.Errorf("Status = %v, want nil", s.Status)
	}
}

func TestDevice_Populate(t *testing.T) {
	d := Device{
		ID:         1,
		Name:       "router1",
		DeviceType: &Nested{ID: 2, Name: "MX480"},
		Site:       &Nested{ID: 3, Name: "Main DC"},
		Status:     &Label{Value: "active", Label: "Active"},
	}
	if d.ID != 1 {
		t.Errorf("ID = %d, want %d", d.ID, 1)
	}
	if d.Name != "router1" {
		t.Errorf("Name = %q, want %q", d.Name, "router1")
	}
	if d.DeviceType.Name != "MX480" {
		t.Errorf("DeviceType.Name = %q, want %q", d.DeviceType.Name, "MX480")
	}
}

func TestDevice_ZeroValues(t *testing.T) {
	var d Device
	if d.ID != 0 {
		t.Errorf("ID = %d, want 0", d.ID)
	}
	if d.Name != "" {
		t.Errorf("Name = %q, want empty", d.Name)
	}
}

func TestIPAddress_Populate(t *testing.T) {
	ip := IPAddress{
		ID:      1,
		Address: "10.0.0.1/24",
		Status:  &Label{Value: "active", Label: "Active"},
		DNSName: "server.example.com",
	}
	if ip.Address != "10.0.0.1/24" {
		t.Errorf("Address = %q, want %q", ip.Address, "10.0.0.1/24")
	}
	if ip.DNSName != "server.example.com" {
		t.Errorf("DNSName = %q, want %q", ip.DNSName, "server.example.com")
	}
}

func TestIPAddress_ZeroValues(t *testing.T) {
	var ip IPAddress
	if ip.Address != "" {
		t.Errorf("Address = %q, want empty", ip.Address)
	}
}

func TestPrefix_Populate(t *testing.T) {
	p := Prefix{
		ID:     1,
		Prefix: "10.0.0.0/8",
		Site:   &Nested{ID: 1, Name: "DC1"},
		Status: &Label{Value: "active", Label: "Active"},
		Family: &Family{Value: 4, Label: "IPv4"},
	}
	if p.Prefix != "10.0.0.0/8" {
		t.Errorf("Prefix = %q, want %q", p.Prefix, "10.0.0.0/8")
	}
	if p.Family.Value != 4 {
		t.Errorf("Family.Value = %d, want %d", p.Family.Value, 4)
	}
}

func TestPrefix_ZeroValues(t *testing.T) {
	var p Prefix
	if p.Prefix != "" {
		t.Errorf("Prefix = %q, want empty", p.Prefix)
	}
}

func TestVLAN_Populate(t *testing.T) {
	v := VLAN{
		ID:   1,
		VID:  100,
		Name: "Users",
		Site: &Nested{ID: 1, Name: "DC1"},
	}
	if v.VID != 100 {
		t.Errorf("VID = %d, want %d", v.VID, 100)
	}
	if v.Name != "Users" {
		t.Errorf("Name = %q, want %q", v.Name, "Users")
	}
}

func TestVLAN_ZeroValues(t *testing.T) {
	var v VLAN
	if v.VID != 0 {
		t.Errorf("VID = %d, want 0", v.VID)
	}
}

func TestVirtualMachine_Populate(t *testing.T) {
	vm := VirtualMachine{
		ID:      1,
		Name:    "web-01",
		Cluster: &Nested{ID: 1, Name: "Prod Cluster"},
		VCPUs:   4,
		Memory:  16384,
		Disk:    100,
	}
	if vm.Name != "web-01" {
		t.Errorf("Name = %q, want %q", vm.Name, "web-01")
	}
	if vm.VCPUs != 4 {
		t.Errorf("VCPUs = %f, want %f", vm.VCPUs, 4.0)
	}
}

func TestVirtualMachine_ZeroValues(t *testing.T) {
	var vm VirtualMachine
	if vm.VCPUs != 0 {
		t.Errorf("VCPUs = %f, want 0", vm.VCPUs)
	}
}

func TestCluster_Populate(t *testing.T) {
	c := Cluster{
		ID:          1,
		Name:        "Production",
		ClusterType: &Nested{Name: "VMware"},
		Site:        &Nested{Name: "DC1"},
	}
	if c.Name != "Production" {
		t.Errorf("Name = %q, want %q", c.Name, "Production")
	}
	if c.ClusterType.Name != "VMware" {
		t.Errorf("ClusterType.Name = %q, want %q", c.ClusterType.Name, "VMware")
	}
}

func TestCluster_ZeroValues(t *testing.T) {
	var c Cluster
	if c.ID != 0 {
		t.Errorf("ID = %d, want 0", c.ID)
	}
}

func TestCircuit_Populate(t *testing.T) {
	c := Circuit{
		ID:          1,
		CID:         "CIR-001",
		Provider:    &Nested{Name: "ATT"},
		CircuitType: &Nested{Name: "Dark Fiber"},
		Status:      &Label{Value: "active", Label: "Active"},
	}
	if c.CID != "CIR-001" {
		t.Errorf("CID = %q, want %q", c.CID, "CIR-001")
	}
	if c.Provider.Name != "ATT" {
		t.Errorf("Provider.Name = %q, want %q", c.Provider.Name, "ATT")
	}
}

func TestCircuit_ZeroValues(t *testing.T) {
	var c Circuit
	if c.ID != 0 {
		t.Errorf("ID = %d, want 0", c.ID)
	}
}

func TestRack_Populate(t *testing.T) {
	r := Rack{
		ID:      1,
		Name:    "Rack-01",
		Site:    &Nested{Name: "DC1"},
		UHeight: 42,
		Width:   19,
	}
	if r.Name != "Rack-01" {
		t.Errorf("Name = %q, want %q", r.Name, "Rack-01")
	}
	if r.UHeight != 42 {
		t.Errorf("UHeight = %d, want %d", r.UHeight, 42)
	}
}

func TestRack_ZeroValues(t *testing.T) {
	var r Rack
	if r.UHeight != 0 {
		t.Errorf("UHeight = %d, want 0", r.UHeight)
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
	ct := CircuitTermination{
		ID:            1,
		Circuit:       &Nested{ID: 1, Name: "CIR-001"},
		TermSide:      "A",
		Site:          &Nested{ID: 1, Name: "DC1"},
		Speed:         10000000,
		UpstreamSpeed: 10000000,
	}
	if ct.TermSide != "A" {
		t.Errorf("TermSide = %q, want %q", ct.TermSide, "A")
	}
	if ct.Circuit.Name != "CIR-001" {
		t.Errorf("Circuit.Name = %q, want %q", ct.Circuit.Name, "CIR-001")
	}
	if ct.Speed != 10000000 {
		t.Errorf("Speed = %d, want %d", ct.Speed, 10000000)
	}
}

func TestCircuitTermination_ZeroValues(t *testing.T) {
	var ct CircuitTermination
	if ct.TermSide != "" {
		t.Errorf("TermSide = %q, want empty", ct.TermSide)
	}
}

func TestCable_Populate(t *testing.T) {
	c := Cable{
		ID:         1,
		Type:       &Label{Value: "cat6a", Label: "CAT6a"},
		Status:     &Label{Value: "connected", Label: "Connected"},
		Label:      "link-01",
		Color:      "blue",
		Length:     10.5,
		LengthUnit: &Label{Value: "m", Label: "Meters"},
	}
	if c.Label != "link-01" {
		t.Errorf("Label = %q, want %q", c.Label, "link-01")
	}
	if c.Type.Value != "cat6a" {
		t.Errorf("Type.Value = %q, want %q", c.Type.Value, "cat6a")
	}
	if c.Status.Value != "connected" {
		t.Errorf("Status.Value = %q, want %q", c.Status.Value, "connected")
	}
	if c.Length != 10.5 {
		t.Errorf("Length = %f, want %f", c.Length, 10.5)
	}
}

func TestCable_ZeroValues(t *testing.T) {
	var c Cable
	if c.Label != "" {
		t.Errorf("Label = %q, want empty", c.Label)
	}
	if c.Type != nil {
		t.Errorf("Type = %+v, want nil", c.Type)
	}
}

func TestNested_Populate(t *testing.T) {
	n := Nested{
		ID:   1,
		Name: "Test",
		Slug: "test",
		URL:  "http://netbox/api/test/1/",
	}
	if n.ID != 1 {
		t.Errorf("ID = %d, want %d", n.ID, 1)
	}
	if n.Name != "Test" {
		t.Errorf("Name = %q, want %q", n.Name, "Test")
	}
}

func TestNested_ZeroValues(t *testing.T) {
	var n Nested
	if n.ID != 0 {
		t.Errorf("ID = %d, want 0", n.ID)
	}
}

func TestLabel_Populate(t *testing.T) {
	l := Label{Value: "active", Label: "Active"}
	if l.Value != "active" {
		t.Errorf("Value = %q, want %q", l.Value, "active")
	}
	if l.Label != "Active" {
		t.Errorf("Label = %q, want %q", l.Label, "Active")
	}
}

func TestLabel_ZeroValues(t *testing.T) {
	var l Label
	if l.Value != "" {
		t.Errorf("Value = %q, want empty", l.Value)
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

func (m *mockRepo) ListSites(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[Site], error) {
	return m.listSitesFunc(ctx, token, params)
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.in)
			if err != nil {
				t.Fatalf("MarshalJSON() error: %v", err)
			}
			if string(b) != tc.want {
				t.Errorf("MarshalJSON() = %s, want %s", b, tc.want)
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
	b, err := json.Marshal(w)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	// The raw object must be inlined as JSON, not base64.
	if string(b) != `{"name":"n","data":{"x":1}}` {
		t.Errorf("Marshal() = %s, want inlined JSON", b)
	}
}
