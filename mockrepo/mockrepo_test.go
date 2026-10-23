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

// TestMockRepo_AllMethods exercises every MockRepo method through its function field.
func TestMockRepo_AllMethods(t *testing.T) { //nolint:funlen,gocognit,gocyclo,maintidx
	t.Parallel()

	t.Run("ListSites", func(t *testing.T) {
		got, err := (&MockRepo{
			ListSitesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.Site], error) {
				return &domain.PaginatedResponse[domain.Site]{Count: 1}, nil
			},
		}).ListSites(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListSites = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("ListRegions", func(t *testing.T) {
		got, err := (&MockRepo{
			ListRegionsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.Region], error) {
				return &domain.PaginatedResponse[domain.Region]{Count: 1}, nil
			},
		}).ListRegions(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListRegions = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("ListDevices", func(t *testing.T) {
		got, err := (&MockRepo{
			ListDevicesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.Device], error) {
				return &domain.PaginatedResponse[domain.Device]{Count: 2}, nil
			},
		}).ListDevices(context.Background(), "t", nil)
		if err != nil || got.Count != 2 {
			t.Errorf("ListDevices = (%v, %v)", got, err)
		}
	})
	t.Run("ListIPAddresses", func(t *testing.T) {
		got, err := (&MockRepo{
			ListIPAddressesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.IPAddress], error) {
				return &domain.PaginatedResponse[domain.IPAddress]{Count: 3}, nil
			},
		}).ListIPAddresses(context.Background(), "t", nil)
		if err != nil || got.Count != 3 {
			t.Errorf("ListIPAddresses = (%v, %v)", got, err)
		}
	})
	t.Run("ListPrefixes", func(t *testing.T) {
		got, err := (&MockRepo{
			ListPrefixesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.Prefix], error) {
				return &domain.PaginatedResponse[domain.Prefix]{Count: 4}, nil
			},
		}).ListPrefixes(context.Background(), "t", nil)
		if err != nil || got.Count != 4 {
			t.Errorf("ListPrefixes = (%v, %v)", got, err)
		}
	})
	t.Run("ListVLANs", func(t *testing.T) {
		got, err := (&MockRepo{
			ListVLANsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.VLAN], error) {
				return &domain.PaginatedResponse[domain.VLAN]{Count: 5}, nil
			},
		}).ListVLANs(context.Background(), "t", nil)
		if err != nil || got.Count != 5 {
			t.Errorf("ListVLANs = (%v, %v)", got, err)
		}
	})
	t.Run("ListVirtualMachines", func(t *testing.T) {
		got, err := (&MockRepo{
			ListVirtualMachinesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualMachine], error) {
				return &domain.PaginatedResponse[domain.VirtualMachine]{Count: 6}, nil
			},
		}).ListVirtualMachines(context.Background(), "t", nil)
		if err != nil || got.Count != 6 {
			t.Errorf("ListVirtualMachines = (%v, %v)", got, err)
		}
	})
	t.Run("ListClusters", func(t *testing.T) {
		got, err := (&MockRepo{
			ListClustersFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.Cluster], error) {
				return &domain.PaginatedResponse[domain.Cluster]{Count: 7}, nil
			},
		}).ListClusters(context.Background(), "t", nil)
		if err != nil || got.Count != 7 {
			t.Errorf("ListClusters = (%v, %v)", got, err)
		}
	})
	t.Run("ListCircuits", func(t *testing.T) {
		got, err := (&MockRepo{
			ListCircuitsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.Circuit], error) {
				return &domain.PaginatedResponse[domain.Circuit]{Count: 8}, nil
			},
		}).ListCircuits(context.Background(), "t", nil)
		if err != nil || got.Count != 8 {
			t.Errorf("ListCircuits = (%v, %v)", got, err)
		}
	})
	t.Run("ListCircuitTerminations", func(t *testing.T) {
		got, err := (&MockRepo{
			ListCircuitTerminationsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.CircuitTermination], error) {
				return &domain.PaginatedResponse[domain.CircuitTermination]{Count: 9}, nil
			},
		}).ListCircuitTerminations(context.Background(), "t", nil)
		if err != nil || got.Count != 9 {
			t.Errorf("ListCircuitTerminations = (%v, %v)", got, err)
		}
	})
	t.Run("ListCables", func(t *testing.T) {
		got, err := (&MockRepo{
			ListCablesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.Cable], error) {
				return &domain.PaginatedResponse[domain.Cable]{Count: 10}, nil
			},
		}).ListCables(context.Background(), "t", nil)
		if err != nil || got.Count != 10 {
			t.Errorf("ListCables = (%v, %v)", got, err)
		}
	})
	t.Run("ListRacks", func(t *testing.T) {
		got, err := (&MockRepo{
			ListRacksFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.Rack], error) {
				return &domain.PaginatedResponse[domain.Rack]{Count: 11}, nil
			},
		}).ListRacks(context.Background(), "t", nil)
		if err != nil || got.Count != 11 {
			t.Errorf("ListRacks = (%v, %v)", got, err)
		}
	})
	t.Run("ListInterfaces", func(t *testing.T) {
		got, err := (&MockRepo{
			ListInterfacesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.Interface], error) {
				return &domain.PaginatedResponse[domain.Interface]{Count: 12}, nil
			},
		}).ListInterfaces(context.Background(), "t", nil)
		if err != nil || got.Count != 12 {
			t.Errorf("ListInterfaces = (%v, %v)", got, err)
		}
	})
	t.Run("ListVMInterfaces", func(t *testing.T) {
		got, err := (&MockRepo{
			ListVMInterfacesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.VMInterface], error) {
				return &domain.PaginatedResponse[domain.VMInterface]{Count: 13}, nil
			},
		}).ListVMInterfaces(context.Background(), "t", nil)
		if err != nil || got.Count != 13 {
			t.Errorf("ListVMInterfaces = (%v, %v)", got, err)
		}
	})
	t.Run("GetObject", func(t *testing.T) {
		got, err := (&MockRepo{
			GetObjectFunc: func(
				context.Context, string, string, int, map[string]string,
			) (domain.RawObject, error) {
				return domain.RawObject(`{"id":1}`), nil
			},
		}).GetObject(context.Background(), "t", "site", 1, nil)
		if err != nil || string(got) != `{"id":1}` {
			t.Errorf("GetObject = (%v, %v)", got, err)
		}
	})
	t.Run("CreateSite", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateSiteFunc: func(
				_ context.Context, _ string, in domain.SiteWrite,
			) (*domain.Site, error) {
				return &domain.Site{ID: 1, Name: in.Name}, nil
			},
		}).CreateSite(context.Background(), "t", domain.SiteWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreateSite = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateSite", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateSiteFunc: func(
				_ context.Context, _ string, id int, in domain.SiteWrite,
			) (*domain.Site, error) {
				return &domain.Site{ID: id, Name: in.Name}, nil
			},
		}).UpdateSite(context.Background(), "t", 7, domain.SiteWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateSite = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteSite", func(t *testing.T) {
		err := (&MockRepo{
			DeleteSiteFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteSite(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteSite = %v", err)
		}
	})
	t.Run("CreateRegion", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateRegionFunc: func(
				_ context.Context, _ string, in domain.RegionWrite,
			) (*domain.Region, error) {
				return &domain.Region{ID: 1, Name: in.Name}, nil
			},
		}).CreateRegion(context.Background(), "t", domain.RegionWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreateRegion = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateRegion", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateRegionFunc: func(
				_ context.Context, _ string, id int, in domain.RegionWrite,
			) (*domain.Region, error) {
				return &domain.Region{ID: id, Name: in.Name}, nil
			},
		}).UpdateRegion(context.Background(), "t", 7, domain.RegionWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateRegion = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteRegion", func(t *testing.T) {
		err := (&MockRepo{
			DeleteRegionFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteRegion(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteRegion = %v", err)
		}
	})
	t.Run("CreateDevice", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateDeviceFunc: func(
				_ context.Context, _ string, in domain.DeviceWrite,
			) (*domain.Device, error) {
				return &domain.Device{ID: 1, Name: in.Name}, nil
			},
		}).CreateDevice(context.Background(), "t", domain.DeviceWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreateDevice = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateDevice", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateDeviceFunc: func(
				_ context.Context, _ string, id int, in domain.DeviceWrite,
			) (*domain.Device, error) {
				return &domain.Device{ID: id, Name: in.Name}, nil
			},
		}).UpdateDevice(context.Background(), "t", 7, domain.DeviceWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateDevice = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteDevice", func(t *testing.T) {
		err := (&MockRepo{
			DeleteDeviceFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteDevice(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteDevice = %v", err)
		}
	})
	t.Run("CreateIPAddress", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateIPAddressFunc: func(
				_ context.Context, _ string, in domain.IPAddressWrite,
			) (*domain.IPAddress, error) {
				return &domain.IPAddress{ID: 1, Address: in.Address}, nil
			},
		}).CreateIPAddress(context.Background(), "t", domain.IPAddressWrite{Address: "10.0.0.1/32"})
		if err != nil || got == nil || got.Address != "10.0.0.1/32" {
			t.Errorf("CreateIPAddress = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateIPAddress", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateIPAddressFunc: func(
				_ context.Context, _ string, id int, in domain.IPAddressWrite,
			) (*domain.IPAddress, error) {
				return &domain.IPAddress{ID: id, Address: in.Address}, nil
			},
		}).UpdateIPAddress(context.Background(), "t", 7, domain.IPAddressWrite{Address: "10.0.0.1/32"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateIPAddress = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteIPAddress", func(t *testing.T) {
		err := (&MockRepo{
			DeleteIPAddressFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteIPAddress(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteIPAddress = %v", err)
		}
	})
	t.Run("CreatePrefix", func(t *testing.T) {
		got, err := (&MockRepo{
			CreatePrefixFunc: func(
				_ context.Context, _ string, in domain.PrefixWrite,
			) (*domain.Prefix, error) {
				return &domain.Prefix{ID: 1, Prefix: in.Prefix}, nil
			},
		}).CreatePrefix(context.Background(), "t", domain.PrefixWrite{Prefix: "10.0.0.0/24"})
		if err != nil || got == nil || got.Prefix != "10.0.0.0/24" {
			t.Errorf("CreatePrefix = (%v, %v)", got, err)
		}
	})
	t.Run("UpdatePrefix", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdatePrefixFunc: func(
				_ context.Context, _ string, id int, in domain.PrefixWrite,
			) (*domain.Prefix, error) {
				return &domain.Prefix{ID: id, Prefix: in.Prefix}, nil
			},
		}).UpdatePrefix(context.Background(), "t", 7, domain.PrefixWrite{Prefix: "10.0.0.0/24"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdatePrefix = (%v, %v)", got, err)
		}
	})
	t.Run("DeletePrefix", func(t *testing.T) {
		err := (&MockRepo{
			DeletePrefixFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeletePrefix(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeletePrefix = %v", err)
		}
	})
	t.Run("CreateVLAN", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateVLANFunc: func(
				_ context.Context, _ string, in domain.VLANWrite,
			) (*domain.VLAN, error) {
				return &domain.VLAN{ID: 1, VID: in.VID, Name: in.Name}, nil
			},
		}).CreateVLAN(context.Background(), "t", domain.VLANWrite{VID: 100, Name: "mgmt"})
		if err != nil || got == nil || got.Name != "mgmt" {
			t.Errorf("CreateVLAN = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateVLAN", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateVLANFunc: func(
				_ context.Context, _ string, id int, in domain.VLANWrite,
			) (*domain.VLAN, error) {
				return &domain.VLAN{ID: id, VID: in.VID, Name: in.Name}, nil
			},
		}).UpdateVLAN(context.Background(), "t", 7, domain.VLANWrite{VID: 100, Name: "mgmt"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateVLAN = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteVLAN", func(t *testing.T) {
		err := (&MockRepo{
			DeleteVLANFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteVLAN(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteVLAN = %v", err)
		}
	})
	t.Run("CreateVirtualMachine", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateVirtualMachineFunc: func(
				_ context.Context, _ string, in domain.VirtualMachineWrite,
			) (*domain.VirtualMachine, error) {
				return &domain.VirtualMachine{ID: 1, Name: in.Name}, nil
			},
		}).CreateVirtualMachine(context.Background(), "t", domain.VirtualMachineWrite{Name: "web-01"})
		if err != nil || got == nil || got.Name != "web-01" {
			t.Errorf("CreateVirtualMachine = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateVirtualMachine", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateVirtualMachineFunc: func(
				_ context.Context, _ string, id int, in domain.VirtualMachineWrite,
			) (*domain.VirtualMachine, error) {
				return &domain.VirtualMachine{ID: id, Name: in.Name}, nil
			},
		}).UpdateVirtualMachine(context.Background(), "t", 7, domain.VirtualMachineWrite{Name: "web-01"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateVirtualMachine = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteVirtualMachine", func(t *testing.T) {
		err := (&MockRepo{
			DeleteVirtualMachineFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteVirtualMachine(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteVirtualMachine = %v", err)
		}
	})
	t.Run("CreateCluster", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateClusterFunc: func(
				_ context.Context, _ string, in domain.ClusterWrite,
			) (*domain.Cluster, error) {
				return &domain.Cluster{ID: 1, Name: in.Name}, nil
			},
		}).CreateCluster(context.Background(), "t", domain.ClusterWrite{Name: "prod"})
		if err != nil || got == nil || got.Name != "prod" {
			t.Errorf("CreateCluster = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateCluster", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateClusterFunc: func(
				_ context.Context, _ string, id int, in domain.ClusterWrite,
			) (*domain.Cluster, error) {
				return &domain.Cluster{ID: id, Name: in.Name}, nil
			},
		}).UpdateCluster(context.Background(), "t", 7, domain.ClusterWrite{Name: "prod"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateCluster = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteCluster", func(t *testing.T) {
		err := (&MockRepo{
			DeleteClusterFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteCluster(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteCluster = %v", err)
		}
	})
	t.Run("CreateRack", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateRackFunc: func(
				_ context.Context, _ string, in domain.RackWrite,
			) (*domain.Rack, error) {
				return &domain.Rack{ID: 1, Name: in.Name}, nil
			},
		}).CreateRack(context.Background(), "t", domain.RackWrite{Name: "R1"})
		if err != nil || got == nil || got.Name != "R1" {
			t.Errorf("CreateRack = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateRack", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateRackFunc: func(
				_ context.Context, _ string, id int, in domain.RackWrite,
			) (*domain.Rack, error) {
				return &domain.Rack{ID: id, Name: in.Name}, nil
			},
		}).UpdateRack(context.Background(), "t", 7, domain.RackWrite{Name: "R1"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateRack = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteRack", func(t *testing.T) {
		err := (&MockRepo{
			DeleteRackFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteRack(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteRack = %v", err)
		}
	})
	t.Run("CreateInterface", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateInterfaceFunc: func(
				_ context.Context, _ string, in domain.InterfaceWrite,
			) (*domain.Interface, error) {
				return &domain.Interface{ID: 1, Name: in.Name}, nil
			},
		}).CreateInterface(context.Background(), "t", domain.InterfaceWrite{Name: "eth0"})
		if err != nil || got == nil || got.Name != "eth0" {
			t.Errorf("CreateInterface = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateInterface", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateInterfaceFunc: func(
				_ context.Context, _ string, id int, in domain.InterfaceWrite,
			) (*domain.Interface, error) {
				return &domain.Interface{ID: id, Name: in.Name}, nil
			},
		}).UpdateInterface(context.Background(), "t", 7, domain.InterfaceWrite{Name: "eth0"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateInterface = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteInterface", func(t *testing.T) {
		err := (&MockRepo{
			DeleteInterfaceFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteInterface(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteInterface = %v", err)
		}
	})
	t.Run("CreateCircuitTermination", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateCircuitTerminationFunc: func(
				_ context.Context, _ string, in domain.CircuitTerminationWrite,
			) (*domain.CircuitTermination, error) {
				return &domain.CircuitTermination{ID: 1, TermSide: in.TermSide}, nil
			},
		}).CreateCircuitTermination(context.Background(), "t", domain.CircuitTerminationWrite{TermSide: "A"})
		if err != nil || got == nil || got.TermSide != "A" {
			t.Errorf("CreateCircuitTermination = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateCircuitTermination", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateCircuitTerminationFunc: func(
				_ context.Context, _ string, id int, in domain.CircuitTerminationWrite,
			) (*domain.CircuitTermination, error) {
				return &domain.CircuitTermination{ID: id, TermSide: in.TermSide}, nil
			},
		}).UpdateCircuitTermination(context.Background(), "t", 7, domain.CircuitTerminationWrite{TermSide: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateCircuitTermination = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteCircuitTermination", func(t *testing.T) {
		err := (&MockRepo{
			DeleteCircuitTerminationFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteCircuitTermination(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteCircuitTermination = %v", err)
		}
	})
	t.Run("CreateCable", func(t *testing.T) {
		label := "link-01"
		got, err := (&MockRepo{
			CreateCableFunc: func(
				_ context.Context, _ string, in domain.CableWrite,
			) (*domain.Cable, error) {
				return &domain.Cable{ID: 1, Label: *in.Label}, nil
			},
		}).CreateCable(context.Background(), "t", domain.CableWrite{Label: &label})
		if err != nil || got == nil || got.Label != "link-01" {
			t.Errorf("CreateCable = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateCable", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateCableFunc: func(
				_ context.Context, _ string, id int, _ domain.CableWrite,
			) (*domain.Cable, error) {
				return &domain.Cable{ID: id}, nil
			},
		}).UpdateCable(context.Background(), "t", 7, domain.CableWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateCable = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteCable", func(t *testing.T) {
		err := (&MockRepo{
			DeleteCableFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteCable(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteCable = %v", err)
		}
	})
	t.Run("CreateVMInterface", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateVMInterfaceFunc: func(
				_ context.Context, _ string, in domain.VMInterfaceWrite,
			) (*domain.VMInterface, error) {
				return &domain.VMInterface{ID: 1, Name: in.Name}, nil
			},
		}).CreateVMInterface(context.Background(), "t", domain.VMInterfaceWrite{Name: "eth0"})
		if err != nil || got == nil || got.Name != "eth0" {
			t.Errorf("CreateVMInterface = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateVMInterface", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateVMInterfaceFunc: func(
				_ context.Context, _ string, id int, _ domain.VMInterfaceWrite,
			) (*domain.VMInterface, error) {
				return &domain.VMInterface{ID: id}, nil
			},
		}).UpdateVMInterface(context.Background(), "t", 7, domain.VMInterfaceWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateVMInterface = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteVMInterface", func(t *testing.T) {
		err := (&MockRepo{
			DeleteVMInterfaceFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteVMInterface(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteVMInterface = %v", err)
		}
	})
	t.Run("CreateProvider", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateProviderFunc: func(
				_ context.Context, _ string, in domain.ProviderWrite,
			) (*domain.Provider, error) {
				return &domain.Provider{ID: 1, Name: in.Name}, nil
			},
		}).CreateProvider(context.Background(), "t", domain.ProviderWrite{Name: "ACME"})
		if err != nil || got == nil || got.Name != "ACME" {
			t.Errorf("CreateProvider = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateProvider", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateProviderFunc: func(
				_ context.Context, _ string, id int, _ domain.ProviderWrite,
			) (*domain.Provider, error) {
				return &domain.Provider{ID: id}, nil
			},
		}).UpdateProvider(context.Background(), "t", 7, domain.ProviderWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateProvider = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteProvider", func(t *testing.T) {
		err := (&MockRepo{
			DeleteProviderFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteProvider(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteProvider = %v", err)
		}
	})
	t.Run("CreateTenant", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateTenantFunc: func(
				_ context.Context, _ string, in domain.TenantWrite,
			) (*domain.Tenant, error) {
				return &domain.Tenant{ID: 1, Name: in.Name}, nil
			},
		}).CreateTenant(context.Background(), "t", domain.TenantWrite{Name: "ACME"})
		if err != nil || got == nil || got.Name != "ACME" {
			t.Errorf("CreateTenant = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateTenant", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateTenantFunc: func(
				_ context.Context, _ string, id int, _ domain.TenantWrite,
			) (*domain.Tenant, error) {
				return &domain.Tenant{ID: id}, nil
			},
		}).UpdateTenant(context.Background(), "t", 7, domain.TenantWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateTenant = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteTenant", func(t *testing.T) {
		err := (&MockRepo{
			DeleteTenantFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteTenant(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteTenant = %v", err)
		}
	})
	t.Run("CreateManufacturer", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateManufacturerFunc: func(
				_ context.Context, _ string, in domain.ManufacturerWrite,
			) (*domain.Manufacturer, error) {
				return &domain.Manufacturer{ID: 1, Name: in.Name}, nil
			},
		}).CreateManufacturer(context.Background(), "t", domain.ManufacturerWrite{Name: "Cisco"})
		if err != nil || got == nil || got.Name != "Cisco" {
			t.Errorf("CreateManufacturer = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateManufacturer", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateManufacturerFunc: func(
				_ context.Context, _ string, id int, _ domain.ManufacturerWrite,
			) (*domain.Manufacturer, error) {
				return &domain.Manufacturer{ID: id}, nil
			},
		}).UpdateManufacturer(context.Background(), "t", 7, domain.ManufacturerWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateManufacturer = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteManufacturer", func(t *testing.T) {
		err := (&MockRepo{
			DeleteManufacturerFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteManufacturer(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteManufacturer = %v", err)
		}
	})
	t.Run("CreateDeviceType", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateDeviceTypeFunc: func(
				_ context.Context, _ string, in domain.DeviceTypeWrite,
			) (*domain.DeviceType, error) {
				return &domain.DeviceType{ID: 1, Model: in.Model}, nil
			},
		}).CreateDeviceType(context.Background(), "t", domain.DeviceTypeWrite{Model: "C9300"})
		if err != nil || got == nil || got.Model != "C9300" {
			t.Errorf("CreateDeviceType = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateDeviceType", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateDeviceTypeFunc: func(
				_ context.Context, _ string, id int, _ domain.DeviceTypeWrite,
			) (*domain.DeviceType, error) {
				return &domain.DeviceType{ID: id}, nil
			},
		}).UpdateDeviceType(context.Background(), "t", 7, domain.DeviceTypeWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateDeviceType = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteDeviceType", func(t *testing.T) {
		err := (&MockRepo{
			DeleteDeviceTypeFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteDeviceType(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteDeviceType = %v", err)
		}
	})
	t.Run("CreateLocation", func(t *testing.T) {
		site := 5
		got, err := (&MockRepo{
			CreateLocationFunc: func(
				_ context.Context, _ string, in domain.LocationWrite,
			) (*domain.Location, error) {
				return &domain.Location{ID: 1, Name: in.Name}, nil
			},
		}).CreateLocation(context.Background(), "t", domain.LocationWrite{Name: "Row A", Site: &site})
		if err != nil || got == nil || got.Name != "Row A" {
			t.Errorf("CreateLocation = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateLocation", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateLocationFunc: func(
				_ context.Context, _ string, id int, _ domain.LocationWrite,
			) (*domain.Location, error) {
				return &domain.Location{ID: id}, nil
			},
		}).UpdateLocation(context.Background(), "t", 7, domain.LocationWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateLocation = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteLocation", func(t *testing.T) {
		err := (&MockRepo{
			DeleteLocationFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteLocation(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteLocation = %v", err)
		}
	})
	t.Run("CreateClusterType", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateClusterTypeFunc: func(
				_ context.Context, _ string, in domain.ClusterTypeWrite,
			) (*domain.ClusterType, error) {
				return &domain.ClusterType{ID: 1, Name: in.Name}, nil
			},
		}).CreateClusterType(context.Background(), "t", domain.ClusterTypeWrite{Name: "KVM"})
		if err != nil || got == nil || got.Name != "KVM" {
			t.Errorf("CreateClusterType = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateClusterType", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateClusterTypeFunc: func(
				_ context.Context, _ string, id int, _ domain.ClusterTypeWrite,
			) (*domain.ClusterType, error) {
				return &domain.ClusterType{ID: id}, nil
			},
		}).UpdateClusterType(context.Background(), "t", 7, domain.ClusterTypeWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateClusterType = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteClusterType", func(t *testing.T) {
		err := (&MockRepo{
			DeleteClusterTypeFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteClusterType(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteClusterType = %v", err)
		}
	})
	t.Run("CreateClusterGroup", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateClusterGroupFunc: func(
				_ context.Context, _ string, in domain.ClusterGroupWrite,
			) (*domain.ClusterGroup, error) {
				return &domain.ClusterGroup{ID: 1, Name: in.Name}, nil
			},
		}).CreateClusterGroup(context.Background(), "t", domain.ClusterGroupWrite{Name: "DC Clusters"})
		if err != nil || got == nil || got.Name != "DC Clusters" {
			t.Errorf("CreateClusterGroup = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateClusterGroup", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateClusterGroupFunc: func(
				_ context.Context, _ string, id int, _ domain.ClusterGroupWrite,
			) (*domain.ClusterGroup, error) {
				return &domain.ClusterGroup{ID: id}, nil
			},
		}).UpdateClusterGroup(context.Background(), "t", 7, domain.ClusterGroupWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateClusterGroup = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteClusterGroup", func(t *testing.T) {
		err := (&MockRepo{
			DeleteClusterGroupFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteClusterGroup(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteClusterGroup = %v", err)
		}
	})
	t.Run("CreateCircuitType", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateCircuitTypeFunc: func(
				_ context.Context, _ string, in domain.CircuitTypeWrite,
			) (*domain.CircuitType, error) {
				return &domain.CircuitType{ID: 1, Name: in.Name}, nil
			},
		}).CreateCircuitType(context.Background(), "t", domain.CircuitTypeWrite{Name: "Fiber"})
		if err != nil || got == nil || got.Name != "Fiber" {
			t.Errorf("CreateCircuitType = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateCircuitType", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateCircuitTypeFunc: func(
				_ context.Context, _ string, id int, _ domain.CircuitTypeWrite,
			) (*domain.CircuitType, error) {
				return &domain.CircuitType{ID: id}, nil
			},
		}).UpdateCircuitType(context.Background(), "t", 7, domain.CircuitTypeWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateCircuitType = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteCircuitType", func(t *testing.T) {
		err := (&MockRepo{
			DeleteCircuitTypeFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteCircuitType(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteCircuitType = %v", err)
		}
	})
	t.Run("CreateVrf", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateVrfFunc: func(
				_ context.Context, _ string, in domain.VrfWrite,
			) (*domain.Vrf, error) {
				return &domain.Vrf{ID: 1, Name: in.Name, Rd: in.Rd}, nil
			},
		}).CreateVrf(context.Background(), "t", domain.VrfWrite{Name: "prod", Rd: "65000:1"})
		if err != nil || got == nil || got.Name != "prod" || got.Rd != "65000:1" {
			t.Errorf("CreateVrf = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateVrf", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateVrfFunc: func(
				_ context.Context, _ string, id int, _ domain.VrfWrite,
			) (*domain.Vrf, error) {
				return &domain.Vrf{ID: id}, nil
			},
		}).UpdateVrf(context.Background(), "t", 7, domain.VrfWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateVrf = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteVrf", func(t *testing.T) {
		err := (&MockRepo{
			DeleteVrfFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteVrf(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteVrf = %v", err)
		}
	})
	t.Run("CreateVlanGroup", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateVlanGroupFunc: func(
				_ context.Context, _ string, in domain.VlanGroupWrite,
			) (*domain.VlanGroup, error) {
				return &domain.VlanGroup{ID: 1, Name: in.Name}, nil
			},
		}).CreateVlanGroup(context.Background(), "t", domain.VlanGroupWrite{Name: "DC VLANs"})
		if err != nil || got == nil || got.Name != "DC VLANs" {
			t.Errorf("CreateVlanGroup = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateVlanGroup", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateVlanGroupFunc: func(
				_ context.Context, _ string, id int, _ domain.VlanGroupWrite,
			) (*domain.VlanGroup, error) {
				return &domain.VlanGroup{ID: id}, nil
			},
		}).UpdateVlanGroup(context.Background(), "t", 7, domain.VlanGroupWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateVlanGroup = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteVlanGroup", func(t *testing.T) {
		err := (&MockRepo{
			DeleteVlanGroupFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteVlanGroup(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteVlanGroup = %v", err)
		}
	})
	t.Run("CreateRole", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateRoleFunc: func(
				_ context.Context, _ string, in domain.RoleWrite,
			) (*domain.Role, error) {
				return &domain.Role{ID: 1, Name: in.Name}, nil
			},
		}).CreateRole(context.Background(), "t", domain.RoleWrite{Name: "prod"})
		if err != nil || got == nil || got.Name != "prod" {
			t.Errorf("CreateRole = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateRole", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateRoleFunc: func(
				_ context.Context, _ string, id int, _ domain.RoleWrite,
			) (*domain.Role, error) {
				return &domain.Role{ID: id}, nil
			},
		}).UpdateRole(context.Background(), "t", 7, domain.RoleWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateRole = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteRole", func(t *testing.T) {
		err := (&MockRepo{
			DeleteRoleFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteRole(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteRole = %v", err)
		}
	})
	t.Run("CreateContact", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateContactFunc: func(
				_ context.Context, _ string, in domain.ContactWrite,
			) (*domain.Contact, error) {
				return &domain.Contact{ID: 1, Name: in.Name}, nil
			},
		}).CreateContact(context.Background(), "t", domain.ContactWrite{Name: "ops"})
		if err != nil || got == nil || got.Name != "ops" {
			t.Errorf("CreateContact = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateContact", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateContactFunc: func(
				_ context.Context, _ string, id int, _ domain.ContactWrite,
			) (*domain.Contact, error) {
				return &domain.Contact{ID: id}, nil
			},
		}).UpdateContact(context.Background(), "t", 7, domain.ContactWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateContact = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteContact", func(t *testing.T) {
		err := (&MockRepo{
			DeleteContactFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteContact(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteContact = %v", err)
		}
	})
	t.Run("CreateCircuit", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateCircuitFunc: func(
				_ context.Context, _ string, in domain.CircuitWrite,
			) (*domain.Circuit, error) {
				return &domain.Circuit{ID: 1, CID: in.CID}, nil
			},
		}).CreateCircuit(context.Background(), "t", domain.CircuitWrite{CID: "CIR-001"})
		if err != nil || got == nil || got.CID != "CIR-001" {
			t.Errorf("CreateCircuit = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateCircuit", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateCircuitFunc: func(
				_ context.Context, _ string, id int, in domain.CircuitWrite,
			) (*domain.Circuit, error) {
				return &domain.Circuit{ID: id, CID: in.CID}, nil
			},
		}).UpdateCircuit(context.Background(), "t", 7, domain.CircuitWrite{CID: "CIR-001"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateCircuit = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteCircuit", func(t *testing.T) {
		err := (&MockRepo{
			DeleteCircuitFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteCircuit(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteCircuit = %v", err)
		}
	})
	t.Run("ListSiteGroups", func(t *testing.T) {
		got, err := (&MockRepo{
			ListSiteGroupsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.SiteGroup], error) {
				return &domain.PaginatedResponse[domain.SiteGroup]{Count: 1}, nil
			},
		}).ListSiteGroups(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListSiteGroups = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateSiteGroup", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateSiteGroupFunc: func(
				_ context.Context, _ string, in domain.SiteGroupWrite,
			) (*domain.SiteGroup, error) {
				return &domain.SiteGroup{ID: 1, Name: in.Name}, nil
			},
		}).CreateSiteGroup(context.Background(), "t", domain.SiteGroupWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreateSiteGroup = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateSiteGroup", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateSiteGroupFunc: func(
				_ context.Context, _ string, id int, in domain.SiteGroupWrite,
			) (*domain.SiteGroup, error) {
				return &domain.SiteGroup{ID: id, Name: in.Name}, nil
			},
		}).UpdateSiteGroup(context.Background(), "t", 7, domain.SiteGroupWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateSiteGroup = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteSiteGroup", func(t *testing.T) {
		err := (&MockRepo{
			DeleteSiteGroupFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteSiteGroup(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteSiteGroup = %v", err)
		}
	})
	t.Run("ListRackGroups", func(t *testing.T) {
		got, err := (&MockRepo{
			ListRackGroupsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.RackGroup], error) {
				return &domain.PaginatedResponse[domain.RackGroup]{Count: 1}, nil
			},
		}).ListRackGroups(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListRackGroups = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateRackGroup", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateRackGroupFunc: func(
				_ context.Context, _ string, in domain.RackGroupWrite,
			) (*domain.RackGroup, error) {
				return &domain.RackGroup{ID: 1, Name: in.Name}, nil
			},
		}).CreateRackGroup(context.Background(), "t", domain.RackGroupWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreateRackGroup = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateRackGroup", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateRackGroupFunc: func(
				_ context.Context, _ string, id int, in domain.RackGroupWrite,
			) (*domain.RackGroup, error) {
				return &domain.RackGroup{ID: id, Name: in.Name}, nil
			},
		}).UpdateRackGroup(context.Background(), "t", 7, domain.RackGroupWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateRackGroup = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteRackGroup", func(t *testing.T) {
		err := (&MockRepo{
			DeleteRackGroupFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteRackGroup(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteRackGroup = %v", err)
		}
	})
	t.Run("ListRackRoles", func(t *testing.T) {
		got, err := (&MockRepo{
			ListRackRolesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.RackRole], error) {
				return &domain.PaginatedResponse[domain.RackRole]{Count: 1}, nil
			},
		}).ListRackRoles(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListRackRoles = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateRackRole", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateRackRoleFunc: func(
				_ context.Context, _ string, in domain.RackRoleWrite,
			) (*domain.RackRole, error) {
				return &domain.RackRole{ID: 1, Name: in.Name}, nil
			},
		}).CreateRackRole(context.Background(), "t", domain.RackRoleWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreateRackRole = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateRackRole", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateRackRoleFunc: func(
				_ context.Context, _ string, id int, in domain.RackRoleWrite,
			) (*domain.RackRole, error) {
				return &domain.RackRole{ID: id, Name: in.Name}, nil
			},
		}).UpdateRackRole(context.Background(), "t", 7, domain.RackRoleWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateRackRole = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteRackRole", func(t *testing.T) {
		err := (&MockRepo{
			DeleteRackRoleFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteRackRole(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteRackRole = %v", err)
		}
	})
	t.Run("ListDeviceRoles", func(t *testing.T) {
		got, err := (&MockRepo{
			ListDeviceRolesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.DeviceRole], error) {
				return &domain.PaginatedResponse[domain.DeviceRole]{Count: 1}, nil
			},
		}).ListDeviceRoles(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListDeviceRoles = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateDeviceRole", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateDeviceRoleFunc: func(
				_ context.Context, _ string, in domain.DeviceRoleWrite,
			) (*domain.DeviceRole, error) {
				return &domain.DeviceRole{ID: 1, Name: in.Name}, nil
			},
		}).CreateDeviceRole(context.Background(), "t", domain.DeviceRoleWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreateDeviceRole = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateDeviceRole", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateDeviceRoleFunc: func(
				_ context.Context, _ string, id int, in domain.DeviceRoleWrite,
			) (*domain.DeviceRole, error) {
				return &domain.DeviceRole{ID: id, Name: in.Name}, nil
			},
		}).UpdateDeviceRole(context.Background(), "t", 7, domain.DeviceRoleWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateDeviceRole = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteDeviceRole", func(t *testing.T) {
		err := (&MockRepo{
			DeleteDeviceRoleFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteDeviceRole(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteDeviceRole = %v", err)
		}
	})
	t.Run("ListPlatforms", func(t *testing.T) {
		got, err := (&MockRepo{
			ListPlatformsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.Platform], error) {
				return &domain.PaginatedResponse[domain.Platform]{Count: 1}, nil
			},
		}).ListPlatforms(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListPlatforms = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreatePlatform", func(t *testing.T) {
		got, err := (&MockRepo{
			CreatePlatformFunc: func(
				_ context.Context, _ string, in domain.PlatformWrite,
			) (*domain.Platform, error) {
				return &domain.Platform{ID: 1, Name: in.Name}, nil
			},
		}).CreatePlatform(context.Background(), "t", domain.PlatformWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreatePlatform = (%v, %v)", got, err)
		}
	})
	t.Run("UpdatePlatform", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdatePlatformFunc: func(
				_ context.Context, _ string, id int, in domain.PlatformWrite,
			) (*domain.Platform, error) {
				return &domain.Platform{ID: id, Name: in.Name}, nil
			},
		}).UpdatePlatform(context.Background(), "t", 7, domain.PlatformWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdatePlatform = (%v, %v)", got, err)
		}
	})
	t.Run("DeletePlatform", func(t *testing.T) {
		err := (&MockRepo{
			DeletePlatformFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeletePlatform(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeletePlatform = %v", err)
		}
	})
	t.Run("ListVirtualChassis", func(t *testing.T) {
		got, err := (&MockRepo{
			ListVirtualChassisFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualChassis], error) {
				return &domain.PaginatedResponse[domain.VirtualChassis]{Count: 1}, nil
			},
		}).ListVirtualChassis(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListVirtualChassis = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateVirtualChassis", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateVirtualChassisFunc: func(
				_ context.Context, _ string, in domain.VirtualChassisWrite,
			) (*domain.VirtualChassis, error) {
				return &domain.VirtualChassis{ID: 1, Name: in.Name}, nil
			},
		}).CreateVirtualChassis(context.Background(), "t", domain.VirtualChassisWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreateVirtualChassis = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateVirtualChassis", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateVirtualChassisFunc: func(
				_ context.Context, _ string, id int, in domain.VirtualChassisWrite,
			) (*domain.VirtualChassis, error) {
				return &domain.VirtualChassis{ID: id, Name: in.Name}, nil
			},
		}).UpdateVirtualChassis(context.Background(), "t", 7, domain.VirtualChassisWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateVirtualChassis = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteVirtualChassis", func(t *testing.T) {
		err := (&MockRepo{
			DeleteVirtualChassisFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteVirtualChassis(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteVirtualChassis = %v", err)
		}
	})
	t.Run("ListPowerPanels", func(t *testing.T) {
		got, err := (&MockRepo{
			ListPowerPanelsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.PowerPanel], error) {
				return &domain.PaginatedResponse[domain.PowerPanel]{Count: 1}, nil
			},
		}).ListPowerPanels(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListPowerPanels = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreatePowerPanel", func(t *testing.T) {
		got, err := (&MockRepo{
			CreatePowerPanelFunc: func(
				_ context.Context, _ string, in domain.PowerPanelWrite,
			) (*domain.PowerPanel, error) {
				return &domain.PowerPanel{ID: 1, Name: in.Name}, nil
			},
		}).CreatePowerPanel(context.Background(), "t", domain.PowerPanelWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreatePowerPanel = (%v, %v)", got, err)
		}
	})
	t.Run("UpdatePowerPanel", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdatePowerPanelFunc: func(
				_ context.Context, _ string, id int, in domain.PowerPanelWrite,
			) (*domain.PowerPanel, error) {
				return &domain.PowerPanel{ID: id, Name: in.Name}, nil
			},
		}).UpdatePowerPanel(context.Background(), "t", 7, domain.PowerPanelWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdatePowerPanel = (%v, %v)", got, err)
		}
	})
	t.Run("DeletePowerPanel", func(t *testing.T) {
		err := (&MockRepo{
			DeletePowerPanelFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeletePowerPanel(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeletePowerPanel = %v", err)
		}
	})
	t.Run("ListInventoryItemRoles", func(t *testing.T) {
		got, err := (&MockRepo{
			ListInventoryItemRolesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.InventoryItemRole], error) {
				return &domain.PaginatedResponse[domain.InventoryItemRole]{Count: 1}, nil
			},
		}).ListInventoryItemRoles(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListInventoryItemRoles = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateInventoryItemRole", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateInventoryItemRoleFunc: func(
				_ context.Context, _ string, in domain.InventoryItemRoleWrite,
			) (*domain.InventoryItemRole, error) {
				return &domain.InventoryItemRole{ID: 1, Name: in.Name}, nil
			},
		}).CreateInventoryItemRole(context.Background(), "t", domain.InventoryItemRoleWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreateInventoryItemRole = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateInventoryItemRole", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateInventoryItemRoleFunc: func(
				_ context.Context, _ string, id int, in domain.InventoryItemRoleWrite,
			) (*domain.InventoryItemRole, error) {
				return &domain.InventoryItemRole{ID: id, Name: in.Name}, nil
			},
		}).UpdateInventoryItemRole(context.Background(), "t", 7, domain.InventoryItemRoleWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateInventoryItemRole = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteInventoryItemRole", func(t *testing.T) {
		err := (&MockRepo{
			DeleteInventoryItemRoleFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteInventoryItemRole(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteInventoryItemRole = %v", err)
		}
	})
	t.Run("ListFrontPorts", func(t *testing.T) {
		got, err := (&MockRepo{
			ListFrontPortsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.FrontPort], error) {
				return &domain.PaginatedResponse[domain.FrontPort]{Count: 1}, nil
			},
		}).ListFrontPorts(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListFrontPorts = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateFrontPort", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateFrontPortFunc: func(
				_ context.Context, _ string, in domain.FrontPortWrite,
			) (*domain.FrontPort, error) {
				return &domain.FrontPort{ID: 1, Name: in.Name}, nil
			},
		}).CreateFrontPort(context.Background(), "t", domain.FrontPortWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreateFrontPort = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateFrontPort", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateFrontPortFunc: func(
				_ context.Context, _ string, id int, in domain.FrontPortWrite,
			) (*domain.FrontPort, error) {
				return &domain.FrontPort{ID: id, Name: in.Name}, nil
			},
		}).UpdateFrontPort(context.Background(), "t", 7, domain.FrontPortWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateFrontPort = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteFrontPort", func(t *testing.T) {
		err := (&MockRepo{
			DeleteFrontPortFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteFrontPort(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteFrontPort = %v", err)
		}
	})
	t.Run("ListRearPorts", func(t *testing.T) {
		got, err := (&MockRepo{
			ListRearPortsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.RearPort], error) {
				return &domain.PaginatedResponse[domain.RearPort]{Count: 1}, nil
			},
		}).ListRearPorts(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListRearPorts = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateRearPort", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateRearPortFunc: func(
				_ context.Context, _ string, in domain.RearPortWrite,
			) (*domain.RearPort, error) {
				return &domain.RearPort{ID: 1, Name: in.Name}, nil
			},
		}).CreateRearPort(context.Background(), "t", domain.RearPortWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreateRearPort = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateRearPort", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateRearPortFunc: func(
				_ context.Context, _ string, id int, in domain.RearPortWrite,
			) (*domain.RearPort, error) {
				return &domain.RearPort{ID: id, Name: in.Name}, nil
			},
		}).UpdateRearPort(context.Background(), "t", 7, domain.RearPortWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateRearPort = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteRearPort", func(t *testing.T) {
		err := (&MockRepo{
			DeleteRearPortFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteRearPort(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteRearPort = %v", err)
		}
	})
	t.Run("ListConsolePorts", func(t *testing.T) {
		got, err := (&MockRepo{
			ListConsolePortsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.ConsolePort], error) {
				return &domain.PaginatedResponse[domain.ConsolePort]{Count: 1}, nil
			},
		}).ListConsolePorts(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListConsolePorts = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateConsolePort", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateConsolePortFunc: func(
				_ context.Context, _ string, in domain.ConsolePortWrite,
			) (*domain.ConsolePort, error) {
				return &domain.ConsolePort{ID: 1, Name: in.Name}, nil
			},
		}).CreateConsolePort(context.Background(), "t", domain.ConsolePortWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreateConsolePort = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateConsolePort", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateConsolePortFunc: func(
				_ context.Context, _ string, id int, in domain.ConsolePortWrite,
			) (*domain.ConsolePort, error) {
				return &domain.ConsolePort{ID: id, Name: in.Name}, nil
			},
		}).UpdateConsolePort(context.Background(), "t", 7, domain.ConsolePortWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateConsolePort = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteConsolePort", func(t *testing.T) {
		err := (&MockRepo{
			DeleteConsolePortFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteConsolePort(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteConsolePort = %v", err)
		}
	})
	t.Run("ListConsoleServerPorts", func(t *testing.T) {
		got, err := (&MockRepo{
			ListConsoleServerPortsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.ConsoleServerPort], error) {
				return &domain.PaginatedResponse[domain.ConsoleServerPort]{Count: 1}, nil
			},
		}).ListConsoleServerPorts(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListConsoleServerPorts = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateConsoleServerPort", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateConsoleServerPortFunc: func(
				_ context.Context, _ string, in domain.ConsoleServerPortWrite,
			) (*domain.ConsoleServerPort, error) {
				return &domain.ConsoleServerPort{ID: 1, Name: in.Name}, nil
			},
		}).CreateConsoleServerPort(context.Background(), "t", domain.ConsoleServerPortWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreateConsoleServerPort = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateConsoleServerPort", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateConsoleServerPortFunc: func(
				_ context.Context, _ string, id int, in domain.ConsoleServerPortWrite,
			) (*domain.ConsoleServerPort, error) {
				return &domain.ConsoleServerPort{ID: id, Name: in.Name}, nil
			},
		}).UpdateConsoleServerPort(context.Background(), "t", 7, domain.ConsoleServerPortWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateConsoleServerPort = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteConsoleServerPort", func(t *testing.T) {
		err := (&MockRepo{
			DeleteConsoleServerPortFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteConsoleServerPort(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteConsoleServerPort = %v", err)
		}
	})
	t.Run("ListPowerPorts", func(t *testing.T) {
		got, err := (&MockRepo{
			ListPowerPortsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.PowerPort], error) {
				return &domain.PaginatedResponse[domain.PowerPort]{Count: 1}, nil
			},
		}).ListPowerPorts(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListPowerPorts = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreatePowerPort", func(t *testing.T) {
		got, err := (&MockRepo{
			CreatePowerPortFunc: func(
				_ context.Context, _ string, in domain.PowerPortWrite,
			) (*domain.PowerPort, error) {
				return &domain.PowerPort{ID: 1, Name: in.Name}, nil
			},
		}).CreatePowerPort(context.Background(), "t", domain.PowerPortWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreatePowerPort = (%v, %v)", got, err)
		}
	})
	t.Run("UpdatePowerPort", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdatePowerPortFunc: func(
				_ context.Context, _ string, id int, in domain.PowerPortWrite,
			) (*domain.PowerPort, error) {
				return &domain.PowerPort{ID: id, Name: in.Name}, nil
			},
		}).UpdatePowerPort(context.Background(), "t", 7, domain.PowerPortWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdatePowerPort = (%v, %v)", got, err)
		}
	})
	t.Run("DeletePowerPort", func(t *testing.T) {
		err := (&MockRepo{
			DeletePowerPortFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeletePowerPort(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeletePowerPort = %v", err)
		}
	})
	t.Run("ListPowerOutlets", func(t *testing.T) {
		got, err := (&MockRepo{
			ListPowerOutletsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.PowerOutlet], error) {
				return &domain.PaginatedResponse[domain.PowerOutlet]{Count: 1}, nil
			},
		}).ListPowerOutlets(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListPowerOutlets = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreatePowerOutlet", func(t *testing.T) {
		got, err := (&MockRepo{
			CreatePowerOutletFunc: func(
				_ context.Context, _ string, in domain.PowerOutletWrite,
			) (*domain.PowerOutlet, error) {
				return &domain.PowerOutlet{ID: 1, Name: in.Name}, nil
			},
		}).CreatePowerOutlet(context.Background(), "t", domain.PowerOutletWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreatePowerOutlet = (%v, %v)", got, err)
		}
	})
	t.Run("UpdatePowerOutlet", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdatePowerOutletFunc: func(
				_ context.Context, _ string, id int, in domain.PowerOutletWrite,
			) (*domain.PowerOutlet, error) {
				return &domain.PowerOutlet{ID: id, Name: in.Name}, nil
			},
		}).UpdatePowerOutlet(context.Background(), "t", 7, domain.PowerOutletWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdatePowerOutlet = (%v, %v)", got, err)
		}
	})
	t.Run("DeletePowerOutlet", func(t *testing.T) {
		err := (&MockRepo{
			DeletePowerOutletFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeletePowerOutlet(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeletePowerOutlet = %v", err)
		}
	})
	t.Run("ListPowerFeeds", func(t *testing.T) {
		got, err := (&MockRepo{
			ListPowerFeedsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.PowerFeed], error) {
				return &domain.PaginatedResponse[domain.PowerFeed]{Count: 1}, nil
			},
		}).ListPowerFeeds(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListPowerFeeds = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreatePowerFeed", func(t *testing.T) {
		got, err := (&MockRepo{
			CreatePowerFeedFunc: func(
				_ context.Context, _ string, in domain.PowerFeedWrite,
			) (*domain.PowerFeed, error) {
				return &domain.PowerFeed{ID: 1, Name: in.Name}, nil
			},
		}).CreatePowerFeed(context.Background(), "t", domain.PowerFeedWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreatePowerFeed = (%v, %v)", got, err)
		}
	})
	t.Run("UpdatePowerFeed", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdatePowerFeedFunc: func(
				_ context.Context, _ string, id int, in domain.PowerFeedWrite,
			) (*domain.PowerFeed, error) {
				return &domain.PowerFeed{ID: id, Name: in.Name}, nil
			},
		}).UpdatePowerFeed(context.Background(), "t", 7, domain.PowerFeedWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdatePowerFeed = (%v, %v)", got, err)
		}
	})
	t.Run("DeletePowerFeed", func(t *testing.T) {
		err := (&MockRepo{
			DeletePowerFeedFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeletePowerFeed(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeletePowerFeed = %v", err)
		}
	})
	t.Run("ListDeviceBays", func(t *testing.T) {
		got, err := (&MockRepo{
			ListDeviceBaysFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.DeviceBay], error) {
				return &domain.PaginatedResponse[domain.DeviceBay]{Count: 1}, nil
			},
		}).ListDeviceBays(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListDeviceBays = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateDeviceBay", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateDeviceBayFunc: func(
				_ context.Context, _ string, in domain.DeviceBayWrite,
			) (*domain.DeviceBay, error) {
				return &domain.DeviceBay{ID: 1, Name: in.Name}, nil
			},
		}).CreateDeviceBay(context.Background(), "t", domain.DeviceBayWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreateDeviceBay = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateDeviceBay", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateDeviceBayFunc: func(
				_ context.Context, _ string, id int, in domain.DeviceBayWrite,
			) (*domain.DeviceBay, error) {
				return &domain.DeviceBay{ID: id, Name: in.Name}, nil
			},
		}).UpdateDeviceBay(context.Background(), "t", 7, domain.DeviceBayWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateDeviceBay = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteDeviceBay", func(t *testing.T) {
		err := (&MockRepo{
			DeleteDeviceBayFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteDeviceBay(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteDeviceBay = %v", err)
		}
	})
	t.Run("ListRackReservations", func(t *testing.T) {
		got, err := (&MockRepo{
			ListRackReservationsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.RackReservation], error) {
				return &domain.PaginatedResponse[domain.RackReservation]{Count: 1}, nil
			},
		}).ListRackReservations(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListRackReservations = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateRackReservation", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateRackReservationFunc: func(
				_ context.Context, _ string, in domain.RackReservationWrite,
			) (*domain.RackReservation, error) {
				return &domain.RackReservation{ID: 1, Description: in.Description}, nil
			},
		}).CreateRackReservation(context.Background(), "t", domain.RackReservationWrite{Description: "reserved"})
		if err != nil || got == nil || got.Description != "reserved" {
			t.Errorf("CreateRackReservation = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateRackReservation", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateRackReservationFunc: func(
				_ context.Context, _ string, id int, in domain.RackReservationWrite,
			) (*domain.RackReservation, error) {
				return &domain.RackReservation{ID: id, Description: in.Description}, nil
			},
		}).UpdateRackReservation(context.Background(), "t", 7, domain.RackReservationWrite{Description: "reserved"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateRackReservation = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteRackReservation", func(t *testing.T) {
		err := (&MockRepo{
			DeleteRackReservationFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteRackReservation(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteRackReservation = %v", err)
		}
	})
	t.Run("ListInventoryItems", func(t *testing.T) {
		got, err := (&MockRepo{
			ListInventoryItemsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.InventoryItem], error) {
				return &domain.PaginatedResponse[domain.InventoryItem]{Count: 1}, nil
			},
		}).ListInventoryItems(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListInventoryItems = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateInventoryItem", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateInventoryItemFunc: func(
				_ context.Context, _ string, in domain.InventoryItemWrite,
			) (*domain.InventoryItem, error) {
				return &domain.InventoryItem{ID: 1, Name: in.Name}, nil
			},
		}).CreateInventoryItem(context.Background(), "t", domain.InventoryItemWrite{Name: "A"})
		if err != nil || got == nil || got.Name != "A" {
			t.Errorf("CreateInventoryItem = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateInventoryItem", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateInventoryItemFunc: func(
				_ context.Context, _ string, id int, in domain.InventoryItemWrite,
			) (*domain.InventoryItem, error) {
				return &domain.InventoryItem{ID: id, Name: in.Name}, nil
			},
		}).UpdateInventoryItem(context.Background(), "t", 7, domain.InventoryItemWrite{Name: "A"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateInventoryItem = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteInventoryItem", func(t *testing.T) {
		err := (&MockRepo{
			DeleteInventoryItemFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteInventoryItem(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteInventoryItem = %v", err)
		}
	})
	t.Run("ListMACAddresses", func(t *testing.T) {
		got, err := (&MockRepo{
			ListMACAddressesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.MACAddress], error) {
				return &domain.PaginatedResponse[domain.MACAddress]{Count: 1}, nil
			},
		}).ListMACAddresses(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListMACAddresses = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateMACAddress", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateMACAddressFunc: func(
				_ context.Context, _ string, in domain.MACAddressWrite,
			) (*domain.MACAddress, error) {
				return &domain.MACAddress{ID: 1, Address: in.Address}, nil
			},
		}).CreateMACAddress(context.Background(), "t", domain.MACAddressWrite{Address: "00:11:22:33:44:55"})
		if err != nil || got == nil || got.Address != "00:11:22:33:44:55" {
			t.Errorf("CreateMACAddress = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateMACAddress", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateMACAddressFunc: func(
				_ context.Context, _ string, id int, in domain.MACAddressWrite,
			) (*domain.MACAddress, error) {
				return &domain.MACAddress{ID: id, Address: in.Address}, nil
			},
		}).UpdateMACAddress(context.Background(), "t", 7, domain.MACAddressWrite{Address: "00:11:22:33:44:55"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateMACAddress = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteMACAddress", func(t *testing.T) {
		err := (&MockRepo{
			DeleteMACAddressFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteMACAddress(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteMACAddress = %v", err)
		}
	})
	t.Run("ListRIRs", func(t *testing.T) {
		got, err := (&MockRepo{
			ListRIRsFunc: func(context.Context, string, map[string]string) (*domain.PaginatedResponse[domain.RIR], error) {
				return &domain.PaginatedResponse[domain.RIR]{Count: 1}, nil
			},
		}).ListRIRs(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListRIRs = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateRIR", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateRIRFunc: func(_ context.Context, _ string, in domain.RIRWrite) (*domain.RIR, error) {
				return &domain.RIR{ID: 1, Name: in.Name}, nil
			},
		}).CreateRIR(context.Background(), "t", domain.RIRWrite{Name: "RIPE"})
		if err != nil || got == nil || got.Name != "RIPE" {
			t.Errorf("CreateRIR = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateRIR", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateRIRFunc: func(_ context.Context, _ string, id int, _ domain.RIRWrite) (*domain.RIR, error) {
				return &domain.RIR{ID: id}, nil
			},
		}).UpdateRIR(context.Background(), "t", 7, domain.RIRWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateRIR = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteRIR", func(t *testing.T) {
		err := (&MockRepo{
			DeleteRIRFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteRIR(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteRIR = %v", err)
		}
	})
	t.Run("ListAggregates", func(t *testing.T) {
		got, err := (&MockRepo{
			ListAggregatesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.Aggregate], error) {
				return &domain.PaginatedResponse[domain.Aggregate]{Count: 1}, nil
			},
		}).ListAggregates(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListAggregates = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateAggregate", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateAggregateFunc: func(_ context.Context, _ string, in domain.AggregateWrite) (*domain.Aggregate, error) {
				return &domain.Aggregate{ID: 1, Prefix: in.Prefix}, nil
			},
		}).CreateAggregate(context.Background(), "t", domain.AggregateWrite{Prefix: "10.0.0.0/8"})
		if err != nil || got == nil || got.Prefix != "10.0.0.0/8" {
			t.Errorf("CreateAggregate = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateAggregate", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateAggregateFunc: func(_ context.Context, _ string, id int, _ domain.AggregateWrite) (*domain.Aggregate, error) {
				return &domain.Aggregate{ID: id}, nil
			},
		}).UpdateAggregate(context.Background(), "t", 7, domain.AggregateWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateAggregate = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteAggregate", func(t *testing.T) {
		err := (&MockRepo{
			DeleteAggregateFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteAggregate(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteAggregate = %v", err)
		}
	})
	t.Run("ListIPRanges", func(t *testing.T) {
		got, err := (&MockRepo{
			ListIPRangesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.IPRange], error) {
				return &domain.PaginatedResponse[domain.IPRange]{Count: 1}, nil
			},
		}).ListIPRanges(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListIPRanges = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateIPRange", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateIPRangeFunc: func(_ context.Context, _ string, in domain.IPRangeWrite) (*domain.IPRange, error) {
				return &domain.IPRange{ID: 1, StartAddress: in.StartAddress}, nil
			},
		}).CreateIPRange(context.Background(), "t", domain.IPRangeWrite{StartAddress: "10.0.0.1"})
		if err != nil || got == nil || got.StartAddress != "10.0.0.1" {
			t.Errorf("CreateIPRange = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateIPRange", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateIPRangeFunc: func(_ context.Context, _ string, id int, _ domain.IPRangeWrite) (*domain.IPRange, error) {
				return &domain.IPRange{ID: id}, nil
			},
		}).UpdateIPRange(context.Background(), "t", 7, domain.IPRangeWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateIPRange = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteIPRange", func(t *testing.T) {
		err := (&MockRepo{
			DeleteIPRangeFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteIPRange(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteIPRange = %v", err)
		}
	})
	t.Run("ListASNs", func(t *testing.T) {
		got, err := (&MockRepo{
			ListASNsFunc: func(context.Context, string, map[string]string) (*domain.PaginatedResponse[domain.ASN], error) {
				return &domain.PaginatedResponse[domain.ASN]{Count: 1}, nil
			},
		}).ListASNs(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListASNs = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateASN", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateASNFunc: func(_ context.Context, _ string, in domain.ASNWrite) (*domain.ASN, error) {
				return &domain.ASN{ID: 1, ASN: in.ASN}, nil
			},
		}).CreateASN(context.Background(), "t", domain.ASNWrite{ASN: 64512})
		if err != nil || got == nil || got.ASN != 64512 {
			t.Errorf("CreateASN = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateASN", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateASNFunc: func(_ context.Context, _ string, id int, _ domain.ASNWrite) (*domain.ASN, error) {
				return &domain.ASN{ID: id}, nil
			},
		}).UpdateASN(context.Background(), "t", 7, domain.ASNWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateASN = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteASN", func(t *testing.T) {
		err := (&MockRepo{
			DeleteASNFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteASN(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteASN = %v", err)
		}
	})
	t.Run("ListASNRanges", func(t *testing.T) {
		got, err := (&MockRepo{
			ListASNRangesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.ASNRange], error) {
				return &domain.PaginatedResponse[domain.ASNRange]{Count: 1}, nil
			},
		}).ListASNRanges(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListASNRanges = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateASNRange", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateASNRangeFunc: func(_ context.Context, _ string, in domain.ASNRangeWrite) (*domain.ASNRange, error) {
				return &domain.ASNRange{ID: 1, Name: in.Name}, nil
			},
		}).CreateASNRange(context.Background(), "t", domain.ASNRangeWrite{Name: "PRIVATE"})
		if err != nil || got == nil || got.Name != "PRIVATE" {
			t.Errorf("CreateASNRange = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateASNRange", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateASNRangeFunc: func(_ context.Context, _ string, id int, _ domain.ASNRangeWrite) (*domain.ASNRange, error) {
				return &domain.ASNRange{ID: id}, nil
			},
		}).UpdateASNRange(context.Background(), "t", 7, domain.ASNRangeWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateASNRange = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteASNRange", func(t *testing.T) {
		err := (&MockRepo{
			DeleteASNRangeFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteASNRange(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteASNRange = %v", err)
		}
	})
	t.Run("ListRouteTargets", func(t *testing.T) {
		got, err := (&MockRepo{
			ListRouteTargetsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.RouteTarget], error) {
				return &domain.PaginatedResponse[domain.RouteTarget]{Count: 1}, nil
			},
		}).ListRouteTargets(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListRouteTargets = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateRouteTarget", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateRouteTargetFunc: func(_ context.Context, _ string, in domain.RouteTargetWrite) (*domain.RouteTarget, error) {
				return &domain.RouteTarget{ID: 1, Name: in.Name}, nil
			},
		}).CreateRouteTarget(context.Background(), "t", domain.RouteTargetWrite{Name: "65000:100"})
		if err != nil || got == nil || got.Name != "65000:100" {
			t.Errorf("CreateRouteTarget = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateRouteTarget", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateRouteTargetFunc: func(
				_ context.Context, _ string, id int, _ domain.RouteTargetWrite,
			) (*domain.RouteTarget, error) {
				return &domain.RouteTarget{ID: id}, nil
			},
		}).UpdateRouteTarget(context.Background(), "t", 7, domain.RouteTargetWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateRouteTarget = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteRouteTarget", func(t *testing.T) {
		err := (&MockRepo{
			DeleteRouteTargetFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteRouteTarget(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteRouteTarget = %v", err)
		}
	})
	t.Run("ListServices", func(t *testing.T) {
		got, err := (&MockRepo{
			ListServicesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.Service], error) {
				return &domain.PaginatedResponse[domain.Service]{Count: 1}, nil
			},
		}).ListServices(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListServices = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateService", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateServiceFunc: func(_ context.Context, _ string, in domain.ServiceWrite) (*domain.Service, error) {
				return &domain.Service{ID: 1, Name: in.Name}, nil
			},
		}).CreateService(context.Background(), "t", domain.ServiceWrite{Name: "web"})
		if err != nil || got == nil || got.Name != "web" {
			t.Errorf("CreateService = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateService", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateServiceFunc: func(_ context.Context, _ string, id int, _ domain.ServiceWrite) (*domain.Service, error) {
				return &domain.Service{ID: id}, nil
			},
		}).UpdateService(context.Background(), "t", 7, domain.ServiceWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateService = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteService", func(t *testing.T) {
		err := (&MockRepo{
			DeleteServiceFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteService(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteService = %v", err)
		}
	})
	t.Run("ListServiceTemplates", func(t *testing.T) {
		got, err := (&MockRepo{
			ListServiceTemplatesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.ServiceTemplate], error) {
				return &domain.PaginatedResponse[domain.ServiceTemplate]{Count: 1}, nil
			},
		}).ListServiceTemplates(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListServiceTemplates = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateServiceTemplate", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateServiceTemplateFunc: func(
				_ context.Context, _ string, in domain.ServiceTemplateWrite,
			) (*domain.ServiceTemplate, error) {
				return &domain.ServiceTemplate{ID: 1, Name: in.Name}, nil
			},
		}).CreateServiceTemplate(context.Background(), "t", domain.ServiceTemplateWrite{Name: "http"})
		if err != nil || got == nil || got.Name != "http" {
			t.Errorf("CreateServiceTemplate = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateServiceTemplate", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateServiceTemplateFunc: func(
				_ context.Context, _ string, id int, _ domain.ServiceTemplateWrite,
			) (*domain.ServiceTemplate, error) {
				return &domain.ServiceTemplate{ID: id}, nil
			},
		}).UpdateServiceTemplate(context.Background(), "t", 7, domain.ServiceTemplateWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateServiceTemplate = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteServiceTemplate", func(t *testing.T) {
		err := (&MockRepo{
			DeleteServiceTemplateFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteServiceTemplate(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteServiceTemplate = %v", err)
		}
	})
	t.Run("ListFHRPGroups", func(t *testing.T) {
		got, err := (&MockRepo{
			ListFHRPGroupsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.FHRPGroup], error) {
				return &domain.PaginatedResponse[domain.FHRPGroup]{Count: 1}, nil
			},
		}).ListFHRPGroups(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListFHRPGroups = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateFHRPGroup", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateFHRPGroupFunc: func(_ context.Context, _ string, in domain.FHRPGroupWrite) (*domain.FHRPGroup, error) {
				return &domain.FHRPGroup{ID: 1, Name: in.Name}, nil
			},
		}).CreateFHRPGroup(context.Background(), "t", domain.FHRPGroupWrite{Name: "gw"})
		if err != nil || got == nil || got.Name != "gw" {
			t.Errorf("CreateFHRPGroup = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateFHRPGroup", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateFHRPGroupFunc: func(_ context.Context, _ string, id int, _ domain.FHRPGroupWrite) (*domain.FHRPGroup, error) {
				return &domain.FHRPGroup{ID: id}, nil
			},
		}).UpdateFHRPGroup(context.Background(), "t", 7, domain.FHRPGroupWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateFHRPGroup = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteFHRPGroup", func(t *testing.T) {
		err := (&MockRepo{
			DeleteFHRPGroupFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteFHRPGroup(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteFHRPGroup = %v", err)
		}
	})
	t.Run("ListFHRPGroupAssignments", func(t *testing.T) {
		got, err := (&MockRepo{
			ListFHRPGroupAssignmentsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.FHRPGroupAssignment], error) {
				return &domain.PaginatedResponse[domain.FHRPGroupAssignment]{Count: 1}, nil
			},
		}).ListFHRPGroupAssignments(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListFHRPGroupAssignments = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateFHRPGroupAssignment", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateFHRPGroupAssignmentFunc: func(
				_ context.Context, _ string, _ domain.FHRPGroupAssignmentWrite,
			) (*domain.FHRPGroupAssignment, error) {
				return &domain.FHRPGroupAssignment{ID: 1}, nil
			},
		}).CreateFHRPGroupAssignment(context.Background(), "t", domain.FHRPGroupAssignmentWrite{Group: 5, InterfaceID: 6})
		if err != nil || got == nil || got.ID != 1 {
			t.Errorf("CreateFHRPGroupAssignment = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateFHRPGroupAssignment", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateFHRPGroupAssignmentFunc: func(
				_ context.Context, _ string, id int, _ domain.FHRPGroupAssignmentWrite,
			) (*domain.FHRPGroupAssignment, error) {
				return &domain.FHRPGroupAssignment{ID: id}, nil
			},
		}).UpdateFHRPGroupAssignment(context.Background(), "t", 7, domain.FHRPGroupAssignmentWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateFHRPGroupAssignment = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteFHRPGroupAssignment", func(t *testing.T) {
		err := (&MockRepo{
			DeleteFHRPGroupAssignmentFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteFHRPGroupAssignment(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteFHRPGroupAssignment = %v", err)
		}
	})
	t.Run("ListVLANTranslationPolicies", func(t *testing.T) {
		got, err := (&MockRepo{
			ListVLANTranslationPoliciesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.VLANTranslationPolicy], error) {
				return &domain.PaginatedResponse[domain.VLANTranslationPolicy]{Count: 1}, nil
			},
		}).ListVLANTranslationPolicies(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListVLANTranslationPolicies = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateVLANTranslationPolicy", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateVLANTranslationPolicyFunc: func(
				_ context.Context, _ string, in domain.VLANTranslationPolicyWrite,
			) (*domain.VLANTranslationPolicy, error) {
				return &domain.VLANTranslationPolicy{ID: 1, Name: in.Name}, nil
			},
		}).CreateVLANTranslationPolicy(context.Background(), "t", domain.VLANTranslationPolicyWrite{Name: "core"})
		if err != nil || got == nil || got.Name != "core" {
			t.Errorf("CreateVLANTranslationPolicy = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateVLANTranslationPolicy", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateVLANTranslationPolicyFunc: func(
				_ context.Context, _ string, id int, _ domain.VLANTranslationPolicyWrite,
			) (*domain.VLANTranslationPolicy, error) {
				return &domain.VLANTranslationPolicy{ID: id}, nil
			},
		}).UpdateVLANTranslationPolicy(context.Background(), "t", 7, domain.VLANTranslationPolicyWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateVLANTranslationPolicy = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteVLANTranslationPolicy", func(t *testing.T) {
		err := (&MockRepo{
			DeleteVLANTranslationPolicyFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteVLANTranslationPolicy(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteVLANTranslationPolicy = %v", err)
		}
	})
	t.Run("ListVLANTranslationRules", func(t *testing.T) {
		got, err := (&MockRepo{
			ListVLANTranslationRulesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.VLANTranslationRule], error) {
				return &domain.PaginatedResponse[domain.VLANTranslationRule]{Count: 1}, nil
			},
		}).ListVLANTranslationRules(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListVLANTranslationRules = (%v, %v), want count 1, nil", got, err)
		}
	})
	t.Run("CreateVLANTranslationRule", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateVLANTranslationRuleFunc: func(
				_ context.Context, _ string, _ domain.VLANTranslationRuleWrite,
			) (*domain.VLANTranslationRule, error) {
				return &domain.VLANTranslationRule{ID: 1}, nil
			},
		}).CreateVLANTranslationRule(context.Background(), "t", domain.VLANTranslationRuleWrite{Policy: 5})
		if err != nil || got == nil || got.ID != 1 {
			t.Errorf("CreateVLANTranslationRule = (%v, %v)", got, err)
		}
	})
	t.Run("UpdateVLANTranslationRule", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateVLANTranslationRuleFunc: func(
				_ context.Context, _ string, id int, _ domain.VLANTranslationRuleWrite,
			) (*domain.VLANTranslationRule, error) {
				return &domain.VLANTranslationRule{ID: id}, nil
			},
		}).UpdateVLANTranslationRule(context.Background(), "t", 7, domain.VLANTranslationRuleWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateVLANTranslationRule = (%v, %v)", got, err)
		}
	})
	t.Run("DeleteVLANTranslationRule", func(t *testing.T) {
		err := (&MockRepo{
			DeleteVLANTranslationRuleFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteVLANTranslationRule(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteVLANTranslationRule = %v", err)
		}
	})
}

func TestMockRepo_PropagatesError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("boom")
	_, err := (&MockRepo{
		ListSitesFunc: func(
			context.Context, string, map[string]string,
		) (*domain.PaginatedResponse[domain.Site], error) {
			return nil, wantErr
		},
	}).ListSites(context.Background(), "t", nil)
	if !errors.Is(err, wantErr) {
		t.Errorf("ListSites err = %v, want %v", err, wantErr)
	}
}

// TestNilRepo_AllMethods exercises every NilRepo method.
func TestNilRepo_AllMethods(t *testing.T) { //nolint:funlen,gocognit,gocyclo,maintidx
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
	if r, err := repo.ListRegions(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListRegions", r.Count, err)
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
	if r, err := repo.CreateRegion(ctx, "t", domain.RegionWrite{Name: "A"}); err != nil || r.Name != "A" {
		t.Errorf("CreateRegion = (%v, %v)", r, err)
	}
	if r, err := repo.UpdateRegion(ctx, "t", 7, domain.RegionWrite{Name: "A"}); err != nil || r.ID != 7 {
		t.Errorf("UpdateRegion = (%v, %v)", r, err)
	}
	if err := repo.DeleteRegion(ctx, "t", 3); err != nil {
		t.Errorf("DeleteRegion = %v", err)
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
	ipAddr, err := repo.CreateIPAddress(ctx, "t", domain.IPAddressWrite{Address: "10.0.0.1/32"})
	if err != nil || ipAddr.Address != "10.0.0.1/32" {
		t.Errorf("CreateIPAddress = (%v, %v)", ipAddr, err)
	}
	ipAddr, err = repo.UpdateIPAddress(ctx, "t", 7, domain.IPAddressWrite{Address: "10.0.0.1/32"})
	if err != nil || ipAddr.ID != 7 {
		t.Errorf("UpdateIPAddress = (%v, %v)", ipAddr, err)
	}
	if err := repo.DeleteIPAddress(ctx, "t", 3); err != nil {
		t.Errorf("DeleteIPAddress = %v", err)
	}
	prefix, err := repo.CreatePrefix(ctx, "t", domain.PrefixWrite{Prefix: "10.0.0.0/24"})
	if err != nil || prefix.Prefix != "10.0.0.0/24" {
		t.Errorf("CreatePrefix = (%v, %v)", prefix, err)
	}
	prefix, err = repo.UpdatePrefix(ctx, "t", 7, domain.PrefixWrite{Prefix: "10.0.0.0/24"})
	if err != nil || prefix.ID != 7 {
		t.Errorf("UpdatePrefix = (%v, %v)", prefix, err)
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
	virt, err := repo.CreateVirtualMachine(ctx, "t", domain.VirtualMachineWrite{Name: "web-01"})
	if err != nil || virt.Name != "web-01" {
		t.Errorf("CreateVirtualMachine = (%v, %v)", virt, err)
	}
	virt, err = repo.UpdateVirtualMachine(ctx, "t", 7, domain.VirtualMachineWrite{Name: "web-01"})
	if err != nil || virt.ID != 7 {
		t.Errorf("UpdateVirtualMachine = (%v, %v)", virt, err)
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
	ctTerm, err := repo.CreateCircuitTermination(ctx, "t", domain.CircuitTerminationWrite{TermSide: "A"})
	if err != nil || ctTerm.TermSide != "A" {
		t.Errorf("CreateCircuitTermination = (%v, %v)", ctTerm, err)
	}
	ctTerm, err = repo.UpdateCircuitTermination(ctx, "t", 7, domain.CircuitTerminationWrite{TermSide: "A"})
	if err != nil || ctTerm.ID != 7 {
		t.Errorf("UpdateCircuitTermination = (%v, %v)", ctTerm, err)
	}
	if err := repo.DeleteCircuitTermination(ctx, "t", 3); err != nil {
		t.Errorf("DeleteCircuitTermination = %v", err)
	}
	label := "link-01"
	if c, err := repo.CreateCable(ctx, "t", domain.CableWrite{Label: &label}); err != nil || c.Label != "link-01" {
		t.Errorf("CreateCable = (%v, %v)", c, err)
	}
	c, err := repo.UpdateCable(ctx, "t", 7, domain.CableWrite{Label: &label})
	if err != nil || c.ID != 7 || c.Label != "link-01" {
		t.Errorf("UpdateCable = (%v, %v)", c, err)
	}
	if err := repo.DeleteCable(ctx, "t", 3); err != nil {
		t.Errorf("DeleteCable = %v", err)
	}
	vmIface, err := repo.CreateVMInterface(ctx, "t", domain.VMInterfaceWrite{Name: "eth0"})
	if err != nil || vmIface.Name != "eth0" {
		t.Errorf("CreateVMInterface = (%v, %v)", vmIface, err)
	}
	vmIface, err = repo.UpdateVMInterface(ctx, "t", 7, domain.VMInterfaceWrite{Name: "eth0"})
	if err != nil || vmIface.ID != 7 {
		t.Errorf("UpdateVMInterface = (%v, %v)", vmIface, err)
	}
	if err := repo.DeleteVMInterface(ctx, "t", 3); err != nil {
		t.Errorf("DeleteVMInterface = %v", err)
	}
	if p, err := repo.CreateProvider(ctx, "t", domain.ProviderWrite{Name: "ACME"}); err != nil || p.Name != "ACME" {
		t.Errorf("CreateProvider = (%v, %v)", p, err)
	}
	if p, err := repo.UpdateProvider(ctx, "t", 7, domain.ProviderWrite{Name: "ACME"}); err != nil || p.ID != 7 {
		t.Errorf("UpdateProvider = (%v, %v)", p, err)
	}
	if err := repo.DeleteProvider(ctx, "t", 3); err != nil {
		t.Errorf("DeleteProvider = %v", err)
	}
	if tn, err := repo.CreateTenant(ctx, "t", domain.TenantWrite{Name: "ACME"}); err != nil || tn.Name != "ACME" {
		t.Errorf("CreateTenant = (%v, %v)", tn, err)
	}
	if tn, err := repo.UpdateTenant(ctx, "t", 7, domain.TenantWrite{Name: "ACME"}); err != nil || tn.ID != 7 {
		t.Errorf("UpdateTenant = (%v, %v)", tn, err)
	}
	if err := repo.DeleteTenant(ctx, "t", 3); err != nil {
		t.Errorf("DeleteTenant = %v", err)
	}
	manuf, err := repo.CreateManufacturer(ctx, "t", domain.ManufacturerWrite{Name: "Cisco"})
	if err != nil || manuf.Name != "Cisco" {
		t.Errorf("CreateManufacturer = (%v, %v)", manuf, err)
	}
	manuf, err = repo.UpdateManufacturer(ctx, "t", 7, domain.ManufacturerWrite{Name: "Cisco"})
	if err != nil || manuf.ID != 7 {
		t.Errorf("UpdateManufacturer = (%v, %v)", manuf, err)
	}
	if err := repo.DeleteManufacturer(ctx, "t", 3); err != nil {
		t.Errorf("DeleteManufacturer = %v", err)
	}
	devType, err := repo.CreateDeviceType(ctx, "t", domain.DeviceTypeWrite{Model: "C9300"})
	if err != nil || devType.Model != "C9300" {
		t.Errorf("CreateDeviceType = (%v, %v)", devType, err)
	}
	devType, err = repo.UpdateDeviceType(ctx, "t", 7, domain.DeviceTypeWrite{Model: "C9300"})
	if err != nil || devType.ID != 7 {
		t.Errorf("UpdateDeviceType = (%v, %v)", devType, err)
	}
	if err := repo.DeleteDeviceType(ctx, "t", 3); err != nil {
		t.Errorf("DeleteDeviceType = %v", err)
	}
	if loc, err := repo.CreateLocation(ctx, "t", domain.LocationWrite{Name: "Row A"}); err != nil || loc.Name != "Row A" {
		t.Errorf("CreateLocation = (%v, %v)", loc, err)
	}
	if loc, err := repo.UpdateLocation(ctx, "t", 7, domain.LocationWrite{Name: "Row A"}); err != nil || loc.ID != 7 {
		t.Errorf("UpdateLocation = (%v, %v)", loc, err)
	}
	if err := repo.DeleteLocation(ctx, "t", 3); err != nil {
		t.Errorf("DeleteLocation = %v", err)
	}
	if ct, err := repo.CreateClusterType(ctx, "t", domain.ClusterTypeWrite{Name: "KVM"}); err != nil || ct.Name != "KVM" {
		t.Errorf("CreateClusterType = (%v, %v)", ct, err)
	}
	if ct, err := repo.UpdateClusterType(ctx, "t", 7, domain.ClusterTypeWrite{Name: "KVM"}); err != nil || ct.ID != 7 {
		t.Errorf("UpdateClusterType = (%v, %v)", ct, err)
	}
	if err := repo.DeleteClusterType(ctx, "t", 3); err != nil {
		t.Errorf("DeleteClusterType = %v", err)
	}
	clusterGroup, err := repo.CreateClusterGroup(ctx, "t", domain.ClusterGroupWrite{Name: "DC Clusters"})
	if err != nil || clusterGroup.Name != "DC Clusters" {
		t.Errorf("CreateClusterGroup = (%v, %v)", clusterGroup, err)
	}
	clusterGroup, err = repo.UpdateClusterGroup(ctx, "t", 7, domain.ClusterGroupWrite{Name: "DC Clusters"})
	if err != nil || clusterGroup.ID != 7 {
		t.Errorf("UpdateClusterGroup = (%v, %v)", clusterGroup, err)
	}
	if err := repo.DeleteClusterGroup(ctx, "t", 3); err != nil {
		t.Errorf("DeleteClusterGroup = %v", err)
	}
	crt, err := repo.CreateCircuitType(ctx, "t", domain.CircuitTypeWrite{Name: "Fiber"})
	if err != nil || crt.Name != "Fiber" {
		t.Errorf("CreateCircuitType = (%v, %v)", crt, err)
	}
	crt, err = repo.UpdateCircuitType(ctx, "t", 7, domain.CircuitTypeWrite{Name: "Fiber"})
	if err != nil || crt.ID != 7 {
		t.Errorf("UpdateCircuitType = (%v, %v)", crt, err)
	}
	if err := repo.DeleteCircuitType(ctx, "t", 3); err != nil {
		t.Errorf("DeleteCircuitType = %v", err)
	}
	rv, err := repo.CreateVrf(ctx, "t", domain.VrfWrite{Name: "prod", Rd: "65000:1"})
	if err != nil || rv.Name != "prod" || rv.Rd != "65000:1" {
		t.Errorf("CreateVrf = (%v, %v)", rv, err)
	}
	if v, err := repo.UpdateVrf(ctx, "t", 7, domain.VrfWrite{Name: "prod"}); err != nil || v.ID != 7 {
		t.Errorf("UpdateVrf = (%v, %v)", v, err)
	}
	if err := repo.DeleteVrf(ctx, "t", 3); err != nil {
		t.Errorf("DeleteVrf = %v", err)
	}
	vg, err := repo.CreateVlanGroup(ctx, "t", domain.VlanGroupWrite{Name: "DC VLANs"})
	if err != nil || vg.Name != "DC VLANs" {
		t.Errorf("CreateVlanGroup = (%v, %v)", vg, err)
	}
	if vg, err := repo.UpdateVlanGroup(ctx, "t", 7, domain.VlanGroupWrite{Name: "DC VLANs"}); err != nil || vg.ID != 7 {
		t.Errorf("UpdateVlanGroup = (%v, %v)", vg, err)
	}
	if err := repo.DeleteVlanGroup(ctx, "t", 3); err != nil {
		t.Errorf("DeleteVlanGroup = %v", err)
	}
	if v, err := repo.CreateRole(ctx, "t", domain.RoleWrite{Name: "prod"}); err != nil || v.Name != "prod" {
		t.Errorf("CreateRole = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateRole(ctx, "t", 7, domain.RoleWrite{Name: "prod"}); err != nil || v.ID != 7 {
		t.Errorf("UpdateRole = (%v, %v)", v, err)
	}
	if err := repo.DeleteRole(ctx, "t", 3); err != nil {
		t.Errorf("DeleteRole = %v", err)
	}
	if v, err := repo.CreateContact(ctx, "t", domain.ContactWrite{Name: "ops"}); err != nil || v.Name != "ops" {
		t.Errorf("CreateContact = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateContact(ctx, "t", 7, domain.ContactWrite{Name: "ops"}); err != nil || v.ID != 7 {
		t.Errorf("UpdateContact = (%v, %v)", v, err)
	}
	if err := repo.DeleteContact(ctx, "t", 3); err != nil {
		t.Errorf("DeleteContact = %v", err)
	}
	if r, err := repo.ListSiteGroups(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListSiteGroups", r.Count, err)
	}
	if v, err := repo.CreateSiteGroup(ctx, "t", domain.SiteGroupWrite{Name: "A"}); err != nil || v.Name != "A" {
		t.Errorf("CreateSiteGroup = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateSiteGroup(ctx, "t", 7, domain.SiteGroupWrite{Name: "A"}); err != nil || v.ID != 7 {
		t.Errorf("UpdateSiteGroup = (%v, %v)", v, err)
	}
	if err := repo.DeleteSiteGroup(ctx, "t", 3); err != nil {
		t.Errorf("DeleteSiteGroup = %v", err)
	}
	if r, err := repo.ListRackGroups(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListRackGroups", r.Count, err)
	}
	if v, err := repo.CreateRackGroup(ctx, "t", domain.RackGroupWrite{Name: "A"}); err != nil || v.Name != "A" {
		t.Errorf("CreateRackGroup = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateRackGroup(ctx, "t", 7, domain.RackGroupWrite{Name: "A"}); err != nil || v.ID != 7 {
		t.Errorf("UpdateRackGroup = (%v, %v)", v, err)
	}
	if err := repo.DeleteRackGroup(ctx, "t", 3); err != nil {
		t.Errorf("DeleteRackGroup = %v", err)
	}
	if r, err := repo.ListRackRoles(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListRackRoles", r.Count, err)
	}
	if v, err := repo.CreateRackRole(ctx, "t", domain.RackRoleWrite{Name: "A"}); err != nil || v.Name != "A" {
		t.Errorf("CreateRackRole = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateRackRole(ctx, "t", 7, domain.RackRoleWrite{Name: "A"}); err != nil || v.ID != 7 {
		t.Errorf("UpdateRackRole = (%v, %v)", v, err)
	}
	if err := repo.DeleteRackRole(ctx, "t", 3); err != nil {
		t.Errorf("DeleteRackRole = %v", err)
	}
	if r, err := repo.ListDeviceRoles(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListDeviceRoles", r.Count, err)
	}
	if v, err := repo.CreateDeviceRole(ctx, "t", domain.DeviceRoleWrite{Name: "A"}); err != nil || v.Name != "A" {
		t.Errorf("CreateDeviceRole = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateDeviceRole(ctx, "t", 7, domain.DeviceRoleWrite{Name: "A"}); err != nil || v.ID != 7 {
		t.Errorf("UpdateDeviceRole = (%v, %v)", v, err)
	}
	if err := repo.DeleteDeviceRole(ctx, "t", 3); err != nil {
		t.Errorf("DeleteDeviceRole = %v", err)
	}
	if r, err := repo.ListPlatforms(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListPlatforms", r.Count, err)
	}
	if v, err := repo.CreatePlatform(ctx, "t", domain.PlatformWrite{Name: "A"}); err != nil || v.Name != "A" {
		t.Errorf("CreatePlatform = (%v, %v)", v, err)
	}
	if v, err := repo.UpdatePlatform(ctx, "t", 7, domain.PlatformWrite{Name: "A"}); err != nil || v.ID != 7 {
		t.Errorf("UpdatePlatform = (%v, %v)", v, err)
	}
	if err := repo.DeletePlatform(ctx, "t", 3); err != nil {
		t.Errorf("DeletePlatform = %v", err)
	}
	if r, err := repo.ListVirtualChassis(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListVirtualChassis", r.Count, err)
	}
	if v, err := repo.CreateVirtualChassis(ctx, "t", domain.VirtualChassisWrite{Name: "A"}); err != nil || v.Name != "A" {
		t.Errorf("CreateVirtualChassis = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateVirtualChassis(ctx, "t", 7, domain.VirtualChassisWrite{Name: "A"}); err != nil || v.ID != 7 {
		t.Errorf("UpdateVirtualChassis = (%v, %v)", v, err)
	}
	if err := repo.DeleteVirtualChassis(ctx, "t", 3); err != nil {
		t.Errorf("DeleteVirtualChassis = %v", err)
	}
	if r, err := repo.ListPowerPanels(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListPowerPanels", r.Count, err)
	}
	if v, err := repo.CreatePowerPanel(ctx, "t", domain.PowerPanelWrite{Name: "A"}); err != nil || v.Name != "A" {
		t.Errorf("CreatePowerPanel = (%v, %v)", v, err)
	}
	if v, err := repo.UpdatePowerPanel(ctx, "t", 7, domain.PowerPanelWrite{Name: "A"}); err != nil || v.ID != 7 {
		t.Errorf("UpdatePowerPanel = (%v, %v)", v, err)
	}
	if err := repo.DeletePowerPanel(ctx, "t", 3); err != nil {
		t.Errorf("DeletePowerPanel = %v", err)
	}
	if r, err := repo.ListInventoryItemRoles(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListInventoryItemRoles", r.Count, err)
	}
	inv, err := repo.CreateInventoryItemRole(ctx, "t", domain.InventoryItemRoleWrite{Name: "A"})
	if err != nil || inv.Name != "A" {
		t.Errorf("CreateInventoryItemRole = (%v, %v)", inv, err)
	}
	inv, err = repo.UpdateInventoryItemRole(ctx, "t", 7, domain.InventoryItemRoleWrite{Name: "A"})
	if err != nil || inv.ID != 7 {
		t.Errorf("UpdateInventoryItemRole = (%v, %v)", inv, err)
	}
	if err := repo.DeleteInventoryItemRole(ctx, "t", 3); err != nil {
		t.Errorf("DeleteInventoryItemRole = %v", err)
	}
	if r, err := repo.ListFrontPorts(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListFrontPorts", r.Count, err)
	}
	if fp, err := repo.CreateFrontPort(ctx, "t", domain.FrontPortWrite{Name: "A"}); err != nil || fp.Name != "A" {
		t.Errorf("CreateFrontPort = (%v, %v)", fp, err)
	}
	if fp, err := repo.UpdateFrontPort(ctx, "t", 7, domain.FrontPortWrite{Name: "A"}); err != nil || fp.ID != 7 {
		t.Errorf("UpdateFrontPort = (%v, %v)", fp, err)
	}
	if err := repo.DeleteFrontPort(ctx, "t", 3); err != nil {
		t.Errorf("DeleteFrontPort = %v", err)
	}
	if r, err := repo.ListRearPorts(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListRearPorts", r.Count, err)
	}
	if rp, err := repo.CreateRearPort(ctx, "t", domain.RearPortWrite{Name: "A"}); err != nil || rp.Name != "A" {
		t.Errorf("CreateRearPort = (%v, %v)", rp, err)
	}
	if rp, err := repo.UpdateRearPort(ctx, "t", 7, domain.RearPortWrite{Name: "A"}); err != nil || rp.ID != 7 {
		t.Errorf("UpdateRearPort = (%v, %v)", rp, err)
	}
	if err := repo.DeleteRearPort(ctx, "t", 3); err != nil {
		t.Errorf("DeleteRearPort = %v", err)
	}
	if r, err := repo.ListConsolePorts(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListConsolePorts", r.Count, err)
	}
	if cp, err := repo.CreateConsolePort(ctx, "t", domain.ConsolePortWrite{Name: "A"}); err != nil || cp.Name != "A" {
		t.Errorf("CreateConsolePort = (%v, %v)", cp, err)
	}
	if cp, err := repo.UpdateConsolePort(ctx, "t", 7, domain.ConsolePortWrite{Name: "A"}); err != nil || cp.ID != 7 {
		t.Errorf("UpdateConsolePort = (%v, %v)", cp, err)
	}
	if err := repo.DeleteConsolePort(ctx, "t", 3); err != nil {
		t.Errorf("DeleteConsolePort = %v", err)
	}
	if r, err := repo.ListConsoleServerPorts(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListConsoleServerPorts", r.Count, err)
	}
	if csp, err := repo.CreateConsoleServerPort(ctx, "t", domain.ConsoleServerPortWrite{Name: "A"}); err != nil ||
		csp.Name != "A" {
		t.Errorf("CreateConsoleServerPort = (%v, %v)", csp, err)
	}
	if csp, err := repo.UpdateConsoleServerPort(ctx, "t", 7, domain.ConsoleServerPortWrite{Name: "A"}); err != nil ||
		csp.ID != 7 {
		t.Errorf("UpdateConsoleServerPort = (%v, %v)", csp, err)
	}
	if err := repo.DeleteConsoleServerPort(ctx, "t", 3); err != nil {
		t.Errorf("DeleteConsoleServerPort = %v", err)
	}
	if r, err := repo.ListPowerPorts(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListPowerPorts", r.Count, err)
	}
	if pp, err := repo.CreatePowerPort(ctx, "t", domain.PowerPortWrite{Name: "A"}); err != nil || pp.Name != "A" {
		t.Errorf("CreatePowerPort = (%v, %v)", pp, err)
	}
	if pp, err := repo.UpdatePowerPort(ctx, "t", 7, domain.PowerPortWrite{Name: "A"}); err != nil || pp.ID != 7 {
		t.Errorf("UpdatePowerPort = (%v, %v)", pp, err)
	}
	if err := repo.DeletePowerPort(ctx, "t", 3); err != nil {
		t.Errorf("DeletePowerPort = %v", err)
	}
	if r, err := repo.ListPowerOutlets(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListPowerOutlets", r.Count, err)
	}
	if po, err := repo.CreatePowerOutlet(ctx, "t", domain.PowerOutletWrite{Name: "A"}); err != nil || po.Name != "A" {
		t.Errorf("CreatePowerOutlet = (%v, %v)", po, err)
	}
	if po, err := repo.UpdatePowerOutlet(ctx, "t", 7, domain.PowerOutletWrite{Name: "A"}); err != nil || po.ID != 7 {
		t.Errorf("UpdatePowerOutlet = (%v, %v)", po, err)
	}
	if err := repo.DeletePowerOutlet(ctx, "t", 3); err != nil {
		t.Errorf("DeletePowerOutlet = %v", err)
	}
	if r, err := repo.ListPowerFeeds(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListPowerFeeds", r.Count, err)
	}
	if pf, err := repo.CreatePowerFeed(ctx, "t", domain.PowerFeedWrite{Name: "A"}); err != nil || pf.Name != "A" {
		t.Errorf("CreatePowerFeed = (%v, %v)", pf, err)
	}
	if pf, err := repo.UpdatePowerFeed(ctx, "t", 7, domain.PowerFeedWrite{Name: "A"}); err != nil || pf.ID != 7 {
		t.Errorf("UpdatePowerFeed = (%v, %v)", pf, err)
	}
	if err := repo.DeletePowerFeed(ctx, "t", 3); err != nil {
		t.Errorf("DeletePowerFeed = %v", err)
	}
	if r, err := repo.ListDeviceBays(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListDeviceBays", r.Count, err)
	}
	if db, err := repo.CreateDeviceBay(ctx, "t", domain.DeviceBayWrite{Name: "A"}); err != nil || db.Name != "A" {
		t.Errorf("CreateDeviceBay = (%v, %v)", db, err)
	}
	if db, err := repo.UpdateDeviceBay(ctx, "t", 7, domain.DeviceBayWrite{Name: "A"}); err != nil || db.ID != 7 {
		t.Errorf("UpdateDeviceBay = (%v, %v)", db, err)
	}
	if err := repo.DeleteDeviceBay(ctx, "t", 3); err != nil {
		t.Errorf("DeleteDeviceBay = %v", err)
	}
	if r, err := repo.ListRackReservations(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListRackReservations", r.Count, err)
	}
	if rr, err := repo.CreateRackReservation(ctx, "t", domain.RackReservationWrite{Description: "reserved"}); err != nil ||
		rr.Description != "reserved" {
		t.Errorf("CreateRackReservation = (%v, %v)", rr, err)
	}
	rr, err := repo.UpdateRackReservation(ctx, "t", 7, domain.RackReservationWrite{Description: "reserved"})
	if err != nil || rr.ID != 7 {
		t.Errorf("UpdateRackReservation = (%v, %v)", rr, err)
	}
	if err := repo.DeleteRackReservation(ctx, "t", 3); err != nil {
		t.Errorf("DeleteRackReservation = %v", err)
	}
	if r, err := repo.ListInventoryItems(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListInventoryItems", r.Count, err)
	}
	if ii, err := repo.CreateInventoryItem(ctx, "t", domain.InventoryItemWrite{Name: "A"}); err != nil || ii.Name != "A" {
		t.Errorf("CreateInventoryItem = (%v, %v)", ii, err)
	}
	if ii, err := repo.UpdateInventoryItem(ctx, "t", 7, domain.InventoryItemWrite{Name: "A"}); err != nil || ii.ID != 7 {
		t.Errorf("UpdateInventoryItem = (%v, %v)", ii, err)
	}
	if err := repo.DeleteInventoryItem(ctx, "t", 3); err != nil {
		t.Errorf("DeleteInventoryItem = %v", err)
	}
	if r, err := repo.ListMACAddresses(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListMACAddresses", r.Count, err)
	}
	if ma, err := repo.CreateMACAddress(ctx, "t", domain.MACAddressWrite{Address: "00:11:22:33:44:55"}); err != nil ||
		ma.Address != "00:11:22:33:44:55" {
		t.Errorf("CreateMACAddress = (%v, %v)", ma, err)
	}
	if ma, err := repo.UpdateMACAddress(ctx, "t", 7, domain.MACAddressWrite{Address: "00:11:22:33:44:55"}); err != nil ||
		ma.ID != 7 {
		t.Errorf("UpdateMACAddress = (%v, %v)", ma, err)
	}
	if err := repo.DeleteMACAddress(ctx, "t", 3); err != nil {
		t.Errorf("DeleteMACAddress = %v", err)
	}
	if r, err := repo.ListRIRs(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListRIRs", r.Count, err)
	}
	if v, err := repo.CreateRIR(ctx, "t", domain.RIRWrite{Name: "RIPE"}); err != nil || v.Name != "RIPE" {
		t.Errorf("CreateRIR = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateRIR(ctx, "t", 7, domain.RIRWrite{Name: "RIPE"}); err != nil || v.ID != 7 {
		t.Errorf("UpdateRIR = (%v, %v)", v, err)
	}
	if err := repo.DeleteRIR(ctx, "t", 3); err != nil {
		t.Errorf("DeleteRIR = %v", err)
	}
	if r, err := repo.ListAggregates(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListAggregates", r.Count, err)
	}
	if v, err := repo.CreateAggregate(ctx, "t", domain.AggregateWrite{Prefix: "10.0.0.0/8"}); err != nil ||
		v.Prefix != "10.0.0.0/8" {
		t.Errorf("CreateAggregate = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateAggregate(ctx, "t", 7, domain.AggregateWrite{Prefix: "10.0.0.0/8"}); err != nil ||
		v.ID != 7 {
		t.Errorf("UpdateAggregate = (%v, %v)", v, err)
	}
	if err := repo.DeleteAggregate(ctx, "t", 3); err != nil {
		t.Errorf("DeleteAggregate = %v", err)
	}
	if r, err := repo.ListIPRanges(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListIPRanges", r.Count, err)
	}
	if v, err := repo.CreateIPRange(ctx, "t", domain.IPRangeWrite{StartAddress: "10.0.0.1"}); err != nil ||
		v.StartAddress != "10.0.0.1" {
		t.Errorf("CreateIPRange = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateIPRange(ctx, "t", 7, domain.IPRangeWrite{StartAddress: "10.0.0.1"}); err != nil ||
		v.ID != 7 {
		t.Errorf("UpdateIPRange = (%v, %v)", v, err)
	}
	if err := repo.DeleteIPRange(ctx, "t", 3); err != nil {
		t.Errorf("DeleteIPRange = %v", err)
	}
	if r, err := repo.ListASNs(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListASNs", r.Count, err)
	}
	if v, err := repo.CreateASN(ctx, "t", domain.ASNWrite{ASN: 64512}); err != nil || v.ASN != 64512 {
		t.Errorf("CreateASN = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateASN(ctx, "t", 7, domain.ASNWrite{ASN: 64512}); err != nil || v.ID != 7 {
		t.Errorf("UpdateASN = (%v, %v)", v, err)
	}
	if err := repo.DeleteASN(ctx, "t", 3); err != nil {
		t.Errorf("DeleteASN = %v", err)
	}
	if r, err := repo.ListASNRanges(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListASNRanges", r.Count, err)
	}
	if v, err := repo.CreateASNRange(ctx, "t", domain.ASNRangeWrite{Name: "PRIVATE"}); err != nil ||
		v.Name != "PRIVATE" {
		t.Errorf("CreateASNRange = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateASNRange(ctx, "t", 7, domain.ASNRangeWrite{Name: "PRIVATE"}); err != nil ||
		v.ID != 7 {
		t.Errorf("UpdateASNRange = (%v, %v)", v, err)
	}
	if err := repo.DeleteASNRange(ctx, "t", 3); err != nil {
		t.Errorf("DeleteASNRange = %v", err)
	}
	if r, err := repo.ListRouteTargets(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListRouteTargets", r.Count, err)
	}
	if v, err := repo.CreateRouteTarget(ctx, "t", domain.RouteTargetWrite{Name: "65000:100"}); err != nil ||
		v.Name != "65000:100" {
		t.Errorf("CreateRouteTarget = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateRouteTarget(ctx, "t", 7, domain.RouteTargetWrite{Name: "65000:100"}); err != nil ||
		v.ID != 7 {
		t.Errorf("UpdateRouteTarget = (%v, %v)", v, err)
	}
	if err := repo.DeleteRouteTarget(ctx, "t", 3); err != nil {
		t.Errorf("DeleteRouteTarget = %v", err)
	}
	if r, err := repo.ListServices(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListServices", r.Count, err)
	}
	if v, err := repo.CreateService(ctx, "t", domain.ServiceWrite{Name: "web"}); err != nil || v.Name != "web" {
		t.Errorf("CreateService = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateService(ctx, "t", 7, domain.ServiceWrite{Name: "web"}); err != nil || v.ID != 7 {
		t.Errorf("UpdateService = (%v, %v)", v, err)
	}
	if err := repo.DeleteService(ctx, "t", 3); err != nil {
		t.Errorf("DeleteService = %v", err)
	}
	if r, err := repo.ListServiceTemplates(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListServiceTemplates", r.Count, err)
	}
	if v, err := repo.CreateServiceTemplate(ctx, "t", domain.ServiceTemplateWrite{Name: "http"}); err != nil ||
		v.Name != "http" {
		t.Errorf("CreateServiceTemplate = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateServiceTemplate(ctx, "t", 7, domain.ServiceTemplateWrite{Name: "http"}); err != nil ||
		v.ID != 7 {
		t.Errorf("UpdateServiceTemplate = (%v, %v)", v, err)
	}
	if err := repo.DeleteServiceTemplate(ctx, "t", 3); err != nil {
		t.Errorf("DeleteServiceTemplate = %v", err)
	}
	if r, err := repo.ListFHRPGroups(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListFHRPGroups", r.Count, err)
	}
	if v, err := repo.CreateFHRPGroup(ctx, "t", domain.FHRPGroupWrite{Name: "gw"}); err != nil || v.Name != "gw" {
		t.Errorf("CreateFHRPGroup = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateFHRPGroup(ctx, "t", 7, domain.FHRPGroupWrite{Name: "gw"}); err != nil || v.ID != 7 {
		t.Errorf("UpdateFHRPGroup = (%v, %v)", v, err)
	}
	if err := repo.DeleteFHRPGroup(ctx, "t", 3); err != nil {
		t.Errorf("DeleteFHRPGroup = %v", err)
	}
	if r, err := repo.ListFHRPGroupAssignments(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListFHRPGroupAssignments", r.Count, err)
	}
	if v, err := repo.CreateFHRPGroupAssignment(
		ctx, "t", domain.FHRPGroupAssignmentWrite{Group: 5, InterfaceID: 6},
	); err != nil || v.ID != 1 {
		t.Errorf("CreateFHRPGroupAssignment = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateFHRPGroupAssignment(ctx, "t", 7, domain.FHRPGroupAssignmentWrite{Group: 5}); err != nil ||
		v.ID != 7 {
		t.Errorf("UpdateFHRPGroupAssignment = (%v, %v)", v, err)
	}
	if err := repo.DeleteFHRPGroupAssignment(ctx, "t", 3); err != nil {
		t.Errorf("DeleteFHRPGroupAssignment = %v", err)
	}
	if r, err := repo.ListVLANTranslationPolicies(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListVLANTranslationPolicies", r.Count, err)
	}
	if v, err := repo.CreateVLANTranslationPolicy(ctx, "t", domain.VLANTranslationPolicyWrite{Name: "core"}); err != nil ||
		v.Name != "core" {
		t.Errorf("CreateVLANTranslationPolicy = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateVLANTranslationPolicy(
		ctx, "t", 7, domain.VLANTranslationPolicyWrite{Name: "core"},
	); err != nil || v.ID != 7 {
		t.Errorf("UpdateVLANTranslationPolicy = (%v, %v)", v, err)
	}
	if err := repo.DeleteVLANTranslationPolicy(ctx, "t", 3); err != nil {
		t.Errorf("DeleteVLANTranslationPolicy = %v", err)
	}
	if r, err := repo.ListVLANTranslationRules(ctx, "t", nil); err != nil || r.Count != 0 {
		assertEmpty("ListVLANTranslationRules", r.Count, err)
	}
	if v, err := repo.CreateVLANTranslationRule(ctx, "t", domain.VLANTranslationRuleWrite{Policy: 5}); err != nil ||
		v.ID != 1 {
		t.Errorf("CreateVLANTranslationRule = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateVLANTranslationRule(ctx, "t", 7, domain.VLANTranslationRuleWrite{Policy: 5}); err != nil ||
		v.ID != 7 {
		t.Errorf("UpdateVLANTranslationRule = (%v, %v)", v, err)
	}
	if err := repo.DeleteVLANTranslationRule(ctx, "t", 3); err != nil {
		t.Errorf("DeleteVLANTranslationRule = %v", err)
	}
}
