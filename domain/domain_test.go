package domain

import (
	"context"
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
