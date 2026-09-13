package mockrepo

import (
	"context"
	"errors"
	"testing"

	"github.com/teran/mcp-netbox/domain"
)

// compile-time check that MockRepo and NilRepo satisfy the repository interface.
var (
	_ domain.NetworkRepository = (*MockRepo)(nil)
	_ domain.NetworkRepository = (*NilRepo)(nil)
)

//nolint:gocognit,gocyclo,maintidx // exhaustive per-method table test over every repository method
func TestMockRepo_AllMethods(t *testing.T) {
	t.Parallel()

	t.Run("ListSites", func(t *testing.T) {
		got, err := (&MockRepo{ListSitesFunc: func(context.Context, string, map[string]string) (*domain.PaginatedResponse[domain.Site], error) {
			return &domain.PaginatedResponse[domain.Site]{Count: 1}, nil
		}}).ListSites(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListSites = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("ListDevices", func(t *testing.T) {
		got, err := (&MockRepo{ListDevicesFunc: func(context.Context, string, map[string]string) (*domain.PaginatedResponse[domain.Device], error) {
			return &domain.PaginatedResponse[domain.Device]{Count: 2}, nil
		}}).ListDevices(context.Background(), "t", nil)
		if err != nil || got.Count != 2 {
			t.Errorf("ListDevices = (%v, %v)", got, err)
		}
	})
	t.Run("ListIPAddresses", func(t *testing.T) {
		got, err := (&MockRepo{ListIPAddressesFunc: func(context.Context, string, map[string]string) (*domain.PaginatedResponse[domain.IPAddress], error) {
			return &domain.PaginatedResponse[domain.IPAddress]{Count: 3}, nil
		}}).ListIPAddresses(context.Background(), "t", nil)
		if err != nil || got.Count != 3 {
			t.Errorf("ListIPAddresses = (%v, %v)", got, err)
		}
	})
	t.Run("ListPrefixes", func(t *testing.T) {
		got, err := (&MockRepo{ListPrefixesFunc: func(context.Context, string, map[string]string) (*domain.PaginatedResponse[domain.Prefix], error) {
			return &domain.PaginatedResponse[domain.Prefix]{Count: 4}, nil
		}}).ListPrefixes(context.Background(), "t", nil)
		if err != nil || got.Count != 4 {
			t.Errorf("ListPrefixes = (%v, %v)", got, err)
		}
	})
	t.Run("ListVLANs", func(t *testing.T) {
		got, err := (&MockRepo{ListVLANsFunc: func(context.Context, string, map[string]string) (*domain.PaginatedResponse[domain.VLAN], error) {
			return &domain.PaginatedResponse[domain.VLAN]{Count: 5}, nil
		}}).ListVLANs(context.Background(), "t", nil)
		if err != nil || got.Count != 5 {
			t.Errorf("ListVLANs = (%v, %v)", got, err)
		}
	})
	t.Run("ListVirtualMachines", func(t *testing.T) {
		got, err := (&MockRepo{ListVirtualMachinesFunc: func(context.Context, string, map[string]string) (*domain.PaginatedResponse[domain.VirtualMachine], error) {
			return &domain.PaginatedResponse[domain.VirtualMachine]{Count: 6}, nil
		}}).ListVirtualMachines(context.Background(), "t", nil)
		if err != nil || got.Count != 6 {
			t.Errorf("ListVirtualMachines = (%v, %v)", got, err)
		}
	})
	t.Run("ListClusters", func(t *testing.T) {
		got, err := (&MockRepo{ListClustersFunc: func(context.Context, string, map[string]string) (*domain.PaginatedResponse[domain.Cluster], error) {
			return &domain.PaginatedResponse[domain.Cluster]{Count: 7}, nil
		}}).ListClusters(context.Background(), "t", nil)
		if err != nil || got.Count != 7 {
			t.Errorf("ListClusters = (%v, %v)", got, err)
		}
	})
	t.Run("ListCircuits", func(t *testing.T) {
		got, err := (&MockRepo{ListCircuitsFunc: func(context.Context, string, map[string]string) (*domain.PaginatedResponse[domain.Circuit], error) {
			return &domain.PaginatedResponse[domain.Circuit]{Count: 8}, nil
		}}).ListCircuits(context.Background(), "t", nil)
		if err != nil || got.Count != 8 {
			t.Errorf("ListCircuits = (%v, %v)", got, err)
		}
	})
	t.Run("ListCircuitTerminations", func(t *testing.T) {
		got, err := (&MockRepo{ListCircuitTerminationsFunc: func(context.Context, string, map[string]string) (*domain.PaginatedResponse[domain.CircuitTermination], error) {
			return &domain.PaginatedResponse[domain.CircuitTermination]{Count: 9}, nil
		}}).ListCircuitTerminations(context.Background(), "t", nil)
		if err != nil || got.Count != 9 {
			t.Errorf("ListCircuitTerminations = (%v, %v)", got, err)
		}
	})
	t.Run("ListCables", func(t *testing.T) {
		got, err := (&MockRepo{ListCablesFunc: func(context.Context, string, map[string]string) (*domain.PaginatedResponse[domain.Cable], error) {
			return &domain.PaginatedResponse[domain.Cable]{Count: 10}, nil
		}}).ListCables(context.Background(), "t", nil)
		if err != nil || got.Count != 10 {
			t.Errorf("ListCables = (%v, %v)", got, err)
		}
	})
	t.Run("ListRacks", func(t *testing.T) {
		got, err := (&MockRepo{ListRacksFunc: func(context.Context, string, map[string]string) (*domain.PaginatedResponse[domain.Rack], error) {
			return &domain.PaginatedResponse[domain.Rack]{Count: 11}, nil
		}}).ListRacks(context.Background(), "t", nil)
		if err != nil || got.Count != 11 {
			t.Errorf("ListRacks = (%v, %v)", got, err)
		}
	})
	t.Run("ListInterfaces", func(t *testing.T) {
		got, err := (&MockRepo{ListInterfacesFunc: func(context.Context, string, map[string]string) (*domain.PaginatedResponse[domain.Interface], error) {
			return &domain.PaginatedResponse[domain.Interface]{Count: 12}, nil
		}}).ListInterfaces(context.Background(), "t", nil)
		if err != nil || got.Count != 12 {
			t.Errorf("ListInterfaces = (%v, %v)", got, err)
		}
	})
	t.Run("ListVMInterfaces", func(t *testing.T) {
		got, err := (&MockRepo{ListVMInterfacesFunc: func(context.Context, string, map[string]string) (*domain.PaginatedResponse[domain.VMInterface], error) {
			return &domain.PaginatedResponse[domain.VMInterface]{Count: 13}, nil
		}}).ListVMInterfaces(context.Background(), "t", nil)
		if err != nil || got.Count != 13 {
			t.Errorf("ListVMInterfaces = (%v, %v)", got, err)
		}
	})
	t.Run("GetObject", func(t *testing.T) {
		got, err := (&MockRepo{GetObjectFunc: func(context.Context, string, string, int, map[string]string) (domain.RawObject, error) {
			return domain.RawObject(`{"id":1}`), nil
		}}).GetObject(context.Background(), "t", "site", 1, nil)
		if err != nil || string(got) != `{"id":1}` {
			t.Errorf("GetObject = (%v, %v)", got, err)
		}
	})
	t.Run("CreateSite", func(t *testing.T) {
		got, err := (&MockRepo{CreateSiteFunc: func(_ context.Context, _ string, in domain.SiteWrite) (*domain.Site, error) {
			return &domain.Site{ID: 1, Name: in.Name}, nil
		}}).CreateSite(context.Background(), "t", domain.SiteWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreateSite = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateSite", func(t *testing.T) {
		got, err := (&MockRepo{UpdateSiteFunc: func(_ context.Context, _ string, id int, in domain.SiteWrite) (*domain.Site, error) {
			return &domain.Site{ID: id, Name: in.Name}, nil
		}}).UpdateSite(context.Background(), "t", 7, domain.SiteWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateSite = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteSite", func(t *testing.T) {
		err := (&MockRepo{DeleteSiteFunc: func(_ context.Context, _ string, id int) error {
			if id != 3 {
				t.Errorf("id = %d, want 3", id)
			}
			return nil
		}}).DeleteSite(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteSite = %v", err)
		}
	})
	t.Run("CreateDevice", func(t *testing.T) {
		got, err := (&MockRepo{CreateDeviceFunc: func(_ context.Context, _ string, in domain.DeviceWrite) (*domain.Device, error) {
			return &domain.Device{ID: 1, Name: in.Name}, nil
		}}).CreateDevice(context.Background(), "t", domain.DeviceWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreateDevice = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateDevice", func(t *testing.T) {
		got, err := (&MockRepo{UpdateDeviceFunc: func(_ context.Context, _ string, id int, in domain.DeviceWrite) (*domain.Device, error) {
			return &domain.Device{ID: id, Name: in.Name}, nil
		}}).UpdateDevice(context.Background(), "t", 7, domain.DeviceWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateDevice = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteDevice", func(t *testing.T) {
		err := (&MockRepo{DeleteDeviceFunc: func(_ context.Context, _ string, id int) error {
			if id != 3 {
				t.Errorf("id = %d, want 3", id)
			}
			return nil
		}}).DeleteDevice(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteDevice = %v", err)
		}
	})
	t.Run("CreateIPAddress", func(t *testing.T) {
		got, err := (&MockRepo{CreateIPAddressFunc: func(_ context.Context, _ string, in domain.IPAddressWrite) (*domain.IPAddress, error) {
			return &domain.IPAddress{ID: 1, Address: in.Address}, nil
		}}).CreateIPAddress(context.Background(), "t", domain.IPAddressWrite{Address: "10.0.0.1/32"})
		if err != nil || got == nil || got.Address != "10.0.0.1/32" {
			t.Errorf("CreateIPAddress = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateIPAddress", func(t *testing.T) {
		got, err := (&MockRepo{UpdateIPAddressFunc: func(_ context.Context, _ string, id int, in domain.IPAddressWrite) (*domain.IPAddress, error) {
			return &domain.IPAddress{ID: id, Address: in.Address}, nil
		}}).UpdateIPAddress(context.Background(), "t", 7, domain.IPAddressWrite{Address: "10.0.0.1/32"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateIPAddress = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteIPAddress", func(t *testing.T) {
		err := (&MockRepo{DeleteIPAddressFunc: func(_ context.Context, _ string, id int) error {
			if id != 3 {
				t.Errorf("id = %d, want 3", id)
			}
			return nil
		}}).DeleteIPAddress(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteIPAddress = %v", err)
		}
	})
	t.Run("CreatePrefix", func(t *testing.T) {
		got, err := (&MockRepo{CreatePrefixFunc: func(_ context.Context, _ string, in domain.PrefixWrite) (*domain.Prefix, error) {
			return &domain.Prefix{ID: 1, Prefix: in.Prefix}, nil
		}}).CreatePrefix(context.Background(), "t", domain.PrefixWrite{Prefix: "10.0.0.0/24"})
		if err != nil || got == nil || got.Prefix != "10.0.0.0/24" {
			t.Errorf("CreatePrefix = (%v, %v)", got, err)
		}
	})
	t.Run("UpdatePrefix", func(t *testing.T) {
		got, err := (&MockRepo{UpdatePrefixFunc: func(_ context.Context, _ string, id int, in domain.PrefixWrite) (*domain.Prefix, error) {
			return &domain.Prefix{ID: id, Prefix: in.Prefix}, nil
		}}).UpdatePrefix(context.Background(), "t", 7, domain.PrefixWrite{Prefix: "10.0.0.0/24"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdatePrefix = (%v, %v)", got, err)
		}
	})
	t.Run("DeletePrefix", func(t *testing.T) {
		err := (&MockRepo{DeletePrefixFunc: func(_ context.Context, _ string, id int) error {
			if id != 3 {
				t.Errorf("id = %d, want 3", id)
			}
			return nil
		}}).DeletePrefix(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeletePrefix = %v", err)
		}
	})
	t.Run("CreateVLAN", func(t *testing.T) {
		got, err := (&MockRepo{CreateVLANFunc: func(_ context.Context, _ string, in domain.VLANWrite) (*domain.VLAN, error) {
			return &domain.VLAN{ID: 1, VID: in.VID, Name: in.Name}, nil
		}}).CreateVLAN(context.Background(), "t", domain.VLANWrite{VID: 100, Name: "mgmt"})
		if err != nil || got == nil || got.Name != "mgmt" {
			t.Errorf("CreateVLAN = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateVLAN", func(t *testing.T) {
		got, err := (&MockRepo{UpdateVLANFunc: func(_ context.Context, _ string, id int, in domain.VLANWrite) (*domain.VLAN, error) {
			return &domain.VLAN{ID: id, VID: in.VID, Name: in.Name}, nil
		}}).UpdateVLAN(context.Background(), "t", 7, domain.VLANWrite{VID: 100, Name: "mgmt"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateVLAN = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteVLAN", func(t *testing.T) {
		err := (&MockRepo{DeleteVLANFunc: func(_ context.Context, _ string, id int) error {
			if id != 3 {
				t.Errorf("id = %d, want 3", id)
			}
			return nil
		}}).DeleteVLAN(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteVLAN = %v", err)
		}
	})
	t.Run("CreateVirtualMachine", func(t *testing.T) {
		got, err := (&MockRepo{CreateVirtualMachineFunc: func(_ context.Context, _ string, in domain.VirtualMachineWrite) (*domain.VirtualMachine, error) {
			return &domain.VirtualMachine{ID: 1, Name: in.Name}, nil
		}}).CreateVirtualMachine(context.Background(), "t", domain.VirtualMachineWrite{Name: "web-01"})
		if err != nil || got == nil || got.Name != "web-01" {
			t.Errorf("CreateVirtualMachine = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateVirtualMachine", func(t *testing.T) {
		got, err := (&MockRepo{UpdateVirtualMachineFunc: func(_ context.Context, _ string, id int, in domain.VirtualMachineWrite) (*domain.VirtualMachine, error) {
			return &domain.VirtualMachine{ID: id, Name: in.Name}, nil
		}}).UpdateVirtualMachine(context.Background(), "t", 7, domain.VirtualMachineWrite{Name: "web-01"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateVirtualMachine = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteVirtualMachine", func(t *testing.T) {
		err := (&MockRepo{DeleteVirtualMachineFunc: func(_ context.Context, _ string, id int) error {
			if id != 3 {
				t.Errorf("id = %d, want 3", id)
			}
			return nil
		}}).DeleteVirtualMachine(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteVirtualMachine = %v", err)
		}
	})
	t.Run("CreateCluster", func(t *testing.T) {
		got, err := (&MockRepo{CreateClusterFunc: func(_ context.Context, _ string, in domain.ClusterWrite) (*domain.Cluster, error) {
			return &domain.Cluster{ID: 1, Name: in.Name}, nil
		}}).CreateCluster(context.Background(), "t", domain.ClusterWrite{Name: "prod"})
		if err != nil || got == nil || got.Name != "prod" {
			t.Errorf("CreateCluster = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateCluster", func(t *testing.T) {
		got, err := (&MockRepo{UpdateClusterFunc: func(_ context.Context, _ string, id int, in domain.ClusterWrite) (*domain.Cluster, error) {
			return &domain.Cluster{ID: id, Name: in.Name}, nil
		}}).UpdateCluster(context.Background(), "t", 7, domain.ClusterWrite{Name: "prod"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateCluster = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteCluster", func(t *testing.T) {
		err := (&MockRepo{DeleteClusterFunc: func(_ context.Context, _ string, id int) error {
			if id != 3 {
				t.Errorf("id = %d, want 3", id)
			}
			return nil
		}}).DeleteCluster(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteCluster = %v", err)
		}
	})
	t.Run("CreateRack", func(t *testing.T) {
		got, err := (&MockRepo{CreateRackFunc: func(_ context.Context, _ string, in domain.RackWrite) (*domain.Rack, error) {
			return &domain.Rack{ID: 1, Name: in.Name}, nil
		}}).CreateRack(context.Background(), "t", domain.RackWrite{Name: "R1"})
		if err != nil || got == nil || got.Name != "R1" {
			t.Errorf("CreateRack = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateRack", func(t *testing.T) {
		got, err := (&MockRepo{UpdateRackFunc: func(_ context.Context, _ string, id int, in domain.RackWrite) (*domain.Rack, error) {
			return &domain.Rack{ID: id, Name: in.Name}, nil
		}}).UpdateRack(context.Background(), "t", 7, domain.RackWrite{Name: "R1"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateRack = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteRack", func(t *testing.T) {
		err := (&MockRepo{DeleteRackFunc: func(_ context.Context, _ string, id int) error {
			if id != 3 {
				t.Errorf("id = %d, want 3", id)
			}
			return nil
		}}).DeleteRack(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteRack = %v", err)
		}
	})
	t.Run("CreateInterface", func(t *testing.T) {
		got, err := (&MockRepo{CreateInterfaceFunc: func(_ context.Context, _ string, in domain.InterfaceWrite) (*domain.Interface, error) {
			return &domain.Interface{ID: 1, Name: in.Name}, nil
		}}).CreateInterface(context.Background(), "t", domain.InterfaceWrite{Name: "eth0"})
		if err != nil || got == nil || got.Name != "eth0" {
			t.Errorf("CreateInterface = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateInterface", func(t *testing.T) {
		got, err := (&MockRepo{UpdateInterfaceFunc: func(_ context.Context, _ string, id int, in domain.InterfaceWrite) (*domain.Interface, error) {
			return &domain.Interface{ID: id, Name: in.Name}, nil
		}}).UpdateInterface(context.Background(), "t", 7, domain.InterfaceWrite{Name: "eth0"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateInterface = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteInterface", func(t *testing.T) {
		err := (&MockRepo{DeleteInterfaceFunc: func(_ context.Context, _ string, id int) error {
			if id != 3 {
				t.Errorf("id = %d, want 3", id)
			}
			return nil
		}}).DeleteInterface(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteInterface = %v", err)
		}
	})
	t.Run("CreateCircuit", func(t *testing.T) {
		got, err := (&MockRepo{CreateCircuitFunc: func(_ context.Context, _ string, in domain.CircuitWrite) (*domain.Circuit, error) {
			return &domain.Circuit{ID: 1, CID: in.CID}, nil
		}}).CreateCircuit(context.Background(), "t", domain.CircuitWrite{CID: "CIR-001"})
		if err != nil || got == nil || got.CID != "CIR-001" {
			t.Errorf("CreateCircuit = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateCircuit", func(t *testing.T) {
		got, err := (&MockRepo{UpdateCircuitFunc: func(_ context.Context, _ string, id int, in domain.CircuitWrite) (*domain.Circuit, error) {
			return &domain.Circuit{ID: id, CID: in.CID}, nil
		}}).UpdateCircuit(context.Background(), "t", 7, domain.CircuitWrite{CID: "CIR-001"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateCircuit = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteCircuit", func(t *testing.T) {
		err := (&MockRepo{DeleteCircuitFunc: func(_ context.Context, _ string, id int) error {
			if id != 3 {
				t.Errorf("id = %d, want 3", id)
			}
			return nil
		}}).DeleteCircuit(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteCircuit = %v", err)
		}
	})
}

func TestMockRepo_PropagatesError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("boom")
	_, err := (&MockRepo{ListSitesFunc: func(context.Context, string, map[string]string) (*domain.PaginatedResponse[domain.Site], error) {
		return nil, wantErr
	}}).ListSites(context.Background(), "t", nil)
	if err != wantErr {
		t.Errorf("ListSites err = %v, want %v", err, wantErr)
	}
}

//nolint:gocognit,gocyclo,maintidx // exhaustive per-method table test over every repository method
func TestNilRepo_AllMethods(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := &NilRepo{}

	assertEmpty := func(name string, count int, err error) {
		t.Helper()
		if err != nil {
			t.Errorf("%s returned error: %v", name, err)
		}
		if count != 0 {
			t.Errorf("%s count = %d, want 0", name, count)
		}
	}

	if r, err := repo.ListSites(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListSites", r.Count, err)
	}
	if r, err := repo.ListDevices(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListDevices", r.Count, err)
	}
	if r, err := repo.ListIPAddresses(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListIPAddresses", r.Count, err)
	}
	if r, err := repo.ListPrefixes(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListPrefixes", r.Count, err)
	}
	if r, err := repo.ListVLANs(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListVLANs", r.Count, err)
	}
	if r, err := repo.ListVirtualMachines(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListVirtualMachines", r.Count, err)
	}
	if r, err := repo.ListClusters(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListClusters", r.Count, err)
	}
	if r, err := repo.ListCircuits(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListCircuits", r.Count, err)
	}
	if r, err := repo.ListCircuitTerminations(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListCircuitTerminations", r.Count, err)
	}
	if r, err := repo.ListCables(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListCables", r.Count, err)
	}
	if r, err := repo.ListRacks(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListRacks", r.Count, err)
	}
	if r, err := repo.ListInterfaces(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListInterfaces", r.Count, err)
	}
	if r, err := repo.ListVMInterfaces(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListVMInterfaces", r.Count, err)
	}
	obj, err := repo.GetObject(ctx, "t", "site", 1, nil)
	if err != nil {
		t.Errorf("GetObject returned error: %v", err)
	}
	if string(obj) != `{"id":1}` {
		t.Errorf("GetObject = %s, want %s", obj, `{"id":1}`)
	}
	if s, err := repo.CreateSite(ctx, "t", domain.SiteWrite{Name: "A"}); err != nil || s.Name != "A" {
		t.Errorf("CreateSite = (%v, %v)", s, err)
	}
	if s, err := repo.UpdateSite(ctx, "t", 7, domain.SiteWrite{Name: "A"}); err != nil || s.ID != 7 {
		t.Errorf("UpdateSite = (%v, %v)", s, err)
	}
	if err := repo.DeleteSite(ctx, "t", 3); err != nil {
		t.Errorf("DeleteSite = %v", err)
	}
	if d, err := repo.CreateDevice(ctx, "t", domain.DeviceWrite{Name: "A"}); err != nil || d.Name != "A" {
		t.Errorf("CreateDevice = (%v, %v)", d, err)
	}
	if d, err := repo.UpdateDevice(ctx, "t", 7, domain.DeviceWrite{Name: "A"}); err != nil || d.ID != 7 {
		t.Errorf("UpdateDevice = (%v, %v)", d, err)
	}
	if err := repo.DeleteDevice(ctx, "t", 3); err != nil {
		t.Errorf("DeleteDevice = %v", err)
	}
	if ip, err := repo.CreateIPAddress(ctx, "t", domain.IPAddressWrite{Address: "10.0.0.1/32"}); err != nil || ip.Address != "10.0.0.1/32" {
		t.Errorf("CreateIPAddress = (%v, %v)", ip, err)
	}
	if ip, err := repo.UpdateIPAddress(ctx, "t", 7, domain.IPAddressWrite{Address: "10.0.0.1/32"}); err != nil || ip.ID != 7 {
		t.Errorf("UpdateIPAddress = (%v, %v)", ip, err)
	}
	if err := repo.DeleteIPAddress(ctx, "t", 3); err != nil {
		t.Errorf("DeleteIPAddress = %v", err)
	}
	if p, err := repo.CreatePrefix(ctx, "t", domain.PrefixWrite{Prefix: "10.0.0.0/24"}); err != nil || p.Prefix != "10.0.0.0/24" {
		t.Errorf("CreatePrefix = (%v, %v)", p, err)
	}
	if p, err := repo.UpdatePrefix(ctx, "t", 7, domain.PrefixWrite{Prefix: "10.0.0.0/24"}); err != nil || p.ID != 7 {
		t.Errorf("UpdatePrefix = (%v, %v)", p, err)
	}
	if err := repo.DeletePrefix(ctx, "t", 3); err != nil {
		t.Errorf("DeletePrefix = %v", err)
	}
	if v, err := repo.CreateVLAN(ctx, "t", domain.VLANWrite{VID: 100, Name: "mgmt"}); err != nil || v.Name != "mgmt" {
		t.Errorf("CreateVLAN = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateVLAN(ctx, "t", 7, domain.VLANWrite{VID: 100, Name: "mgmt"}); err != nil || v.ID != 7 {
		t.Errorf("UpdateVLAN = (%v, %v)", v, err)
	}
	if err := repo.DeleteVLAN(ctx, "t", 3); err != nil {
		t.Errorf("DeleteVLAN = %v", err)
	}
	if v, err := repo.CreateVirtualMachine(ctx, "t", domain.VirtualMachineWrite{Name: "web-01"}); err != nil || v.Name != "web-01" {
		t.Errorf("CreateVirtualMachine = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateVirtualMachine(ctx, "t", 7, domain.VirtualMachineWrite{Name: "web-01"}); err != nil || v.ID != 7 {
		t.Errorf("UpdateVirtualMachine = (%v, %v)", v, err)
	}
	if err := repo.DeleteVirtualMachine(ctx, "t", 3); err != nil {
		t.Errorf("DeleteVirtualMachine = %v", err)
	}
	if c, err := repo.CreateCluster(ctx, "t", domain.ClusterWrite{Name: "prod"}); err != nil || c.Name != "prod" {
		t.Errorf("CreateCluster = (%v, %v)", c, err)
	}
	if c, err := repo.UpdateCluster(ctx, "t", 7, domain.ClusterWrite{Name: "prod"}); err != nil || c.ID != 7 {
		t.Errorf("UpdateCluster = (%v, %v)", c, err)
	}
	if err := repo.DeleteCluster(ctx, "t", 3); err != nil {
		t.Errorf("DeleteCluster = %v", err)
	}
	if c, err := repo.CreateCircuit(ctx, "t", domain.CircuitWrite{CID: "CIR-001"}); err != nil || c.CID != "CIR-001" {
		t.Errorf("CreateCircuit = (%v, %v)", c, err)
	}
	if c, err := repo.UpdateCircuit(ctx, "t", 7, domain.CircuitWrite{CID: "CIR-001"}); err != nil || c.ID != 7 {
		t.Errorf("UpdateCircuit = (%v, %v)", c, err)
	}
	if err := repo.DeleteCircuit(ctx, "t", 3); err != nil {
		t.Errorf("DeleteCircuit = %v", err)
	}
	if r, err := repo.CreateRack(ctx, "t", domain.RackWrite{Name: "R1"}); err != nil || r.Name != "R1" {
		t.Errorf("CreateRack = (%v, %v)", r, err)
	}
	if r, err := repo.UpdateRack(ctx, "t", 7, domain.RackWrite{Name: "R1"}); err != nil || r.ID != 7 {
		t.Errorf("UpdateRack = (%v, %v)", r, err)
	}
	if err := repo.DeleteRack(ctx, "t", 3); err != nil {
		t.Errorf("DeleteRack = %v", err)
	}
	if i, err := repo.CreateInterface(ctx, "t", domain.InterfaceWrite{Name: "eth0"}); err != nil || i.Name != "eth0" {
		t.Errorf("CreateInterface = (%v, %v)", i, err)
	}
	if i, err := repo.UpdateInterface(ctx, "t", 7, domain.InterfaceWrite{Name: "eth0"}); err != nil || i.ID != 7 {
		t.Errorf("UpdateInterface = (%v, %v)", i, err)
	}
	if err := repo.DeleteInterface(ctx, "t", 3); err != nil {
		t.Errorf("DeleteInterface = %v", err)
	}
}
