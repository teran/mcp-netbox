package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/domain"
)

type stubRepo struct{}

func (s *stubRepo) ListSites(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Site], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListDevices(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Device], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListIPAddresses(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.IPAddress], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListPrefixes(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Prefix], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListVLANs(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VLAN], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListVirtualMachines(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VirtualMachine], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListClusters(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Cluster], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListCircuits(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Circuit], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListRacks(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Rack], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListCircuitTerminations(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.CircuitTermination], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListCables(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Cable], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListInterfaces(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Interface], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListVMInterfaces(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VMInterface], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) GetObject(ctx context.Context, token string, objectType string, id int, params map[string]string) (domain.RawObject, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) CreateSite(ctx context.Context, token string, in domain.SiteWrite) (*domain.Site, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateSite(ctx context.Context, token string, id int, in domain.SiteWrite) (*domain.Site, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteSite(ctx context.Context, token string, id int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateDevice(ctx context.Context, token string, in domain.DeviceWrite) (*domain.Device, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateDevice(ctx context.Context, token string, id int, in domain.DeviceWrite) (*domain.Device, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteDevice(ctx context.Context, token string, id int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateIPAddress(ctx context.Context, token string, in domain.IPAddressWrite) (*domain.IPAddress, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateIPAddress(ctx context.Context, token string, id int, in domain.IPAddressWrite) (*domain.IPAddress, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteIPAddress(ctx context.Context, token string, id int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreatePrefix(ctx context.Context, token string, in domain.PrefixWrite) (*domain.Prefix, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdatePrefix(ctx context.Context, token string, id int, in domain.PrefixWrite) (*domain.Prefix, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeletePrefix(ctx context.Context, token string, id int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateVLAN(ctx context.Context, token string, in domain.VLANWrite) (*domain.VLAN, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateVLAN(ctx context.Context, token string, id int, in domain.VLANWrite) (*domain.VLAN, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteVLAN(ctx context.Context, token string, id int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateVirtualMachine(ctx context.Context, token string, in domain.VirtualMachineWrite) (*domain.VirtualMachine, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateVirtualMachine(ctx context.Context, token string, id int, in domain.VirtualMachineWrite) (*domain.VirtualMachine, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteVirtualMachine(ctx context.Context, token string, id int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateCluster(ctx context.Context, token string, in domain.ClusterWrite) (*domain.Cluster, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateCluster(ctx context.Context, token string, id int, in domain.ClusterWrite) (*domain.Cluster, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteCluster(ctx context.Context, token string, id int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateCircuit(ctx context.Context, token string, in domain.CircuitWrite) (*domain.Circuit, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateCircuit(ctx context.Context, token string, id int, in domain.CircuitWrite) (*domain.Circuit, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteCircuit(ctx context.Context, token string, id int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateRack(ctx context.Context, token string, in domain.RackWrite) (*domain.Rack, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateRack(ctx context.Context, token string, id int, in domain.RackWrite) (*domain.Rack, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteRack(ctx context.Context, token string, id int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateInterface(ctx context.Context, token string, in domain.InterfaceWrite) (*domain.Interface, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateInterface(ctx context.Context, token string, id int, in domain.InterfaceWrite) (*domain.Interface, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteInterface(ctx context.Context, token string, id int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateCircuitTermination(ctx context.Context, token string, in domain.CircuitTerminationWrite) (*domain.CircuitTermination, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateCircuitTermination(ctx context.Context, token string, id int, in domain.CircuitTerminationWrite) (*domain.CircuitTermination, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteCircuitTermination(ctx context.Context, token string, id int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateCable(ctx context.Context, token string, in domain.CableWrite) (*domain.Cable, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateCable(ctx context.Context, token string, id int, in domain.CableWrite) (*domain.Cable, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteCable(ctx context.Context, token string, id int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateVMInterface(ctx context.Context, token string, in domain.VMInterfaceWrite) (*domain.VMInterface, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateVMInterface(ctx context.Context, token string, id int, in domain.VMInterfaceWrite) (*domain.VMInterface, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteVMInterface(ctx context.Context, token string, id int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateProvider(ctx context.Context, token string, in domain.ProviderWrite) (*domain.Provider, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateProvider(ctx context.Context, token string, id int, in domain.ProviderWrite) (*domain.Provider, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteProvider(ctx context.Context, token string, id int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateTenant(ctx context.Context, token string, in domain.TenantWrite) (*domain.Tenant, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateTenant(ctx context.Context, token string, id int, in domain.TenantWrite) (*domain.Tenant, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteTenant(ctx context.Context, token string, id int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateManufacturer(ctx context.Context, token string, in domain.ManufacturerWrite) (*domain.Manufacturer, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateManufacturer(ctx context.Context, token string, id int, in domain.ManufacturerWrite) (*domain.Manufacturer, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteManufacturer(ctx context.Context, token string, id int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateDeviceType(ctx context.Context, token string, in domain.DeviceTypeWrite) (*domain.DeviceType, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateDeviceType(ctx context.Context, token string, id int, in domain.DeviceTypeWrite) (*domain.DeviceType, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteDeviceType(ctx context.Context, token string, id int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateLocation(ctx context.Context, token string, in domain.LocationWrite) (*domain.Location, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateLocation(ctx context.Context, token string, id int, in domain.LocationWrite) (*domain.Location, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteLocation(ctx context.Context, token string, id int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateClusterType(ctx context.Context, token string, in domain.ClusterTypeWrite) (*domain.ClusterType, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateClusterType(ctx context.Context, token string, id int, in domain.ClusterTypeWrite) (*domain.ClusterType, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteClusterType(ctx context.Context, token string, id int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateClusterGroup(ctx context.Context, token string, in domain.ClusterGroupWrite) (*domain.ClusterGroup, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateClusterGroup(ctx context.Context, token string, id int, in domain.ClusterGroupWrite) (*domain.ClusterGroup, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteClusterGroup(ctx context.Context, token string, id int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateCircuitType(ctx context.Context, token string, in domain.CircuitTypeWrite) (*domain.CircuitType, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateCircuitType(ctx context.Context, token string, id int, in domain.CircuitTypeWrite) (*domain.CircuitType, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteCircuitType(ctx context.Context, token string, id int) error {
	return errors.New("not implemented")
}

func TestPaginationParams(t *testing.T) {
	t.Parallel()

	t.Run("default values when page=0 and pageSize=0", func(t *testing.T) {
		params := paginationParams(0, 0)
		if params["offset"] != "0" {
			t.Errorf("offset = %q, want %q", params["offset"], "0")
		}
		if params["limit"] != "25" {
			t.Errorf("limit = %q, want %q", params["limit"], "25")
		}
	})

	t.Run("page=1 and pageSize=25", func(t *testing.T) {
		params := paginationParams(1, 25)
		if params["offset"] != "0" {
			t.Errorf("offset = %q, want %q", params["offset"], "0")
		}
		if params["limit"] != "25" {
			t.Errorf("limit = %q, want %q", params["limit"], "25")
		}
	})

	t.Run("page=3 and pageSize=10", func(t *testing.T) {
		params := paginationParams(3, 10)
		if params["offset"] != "20" {
			t.Errorf("offset = %q, want %q", params["offset"], "20")
		}
		if params["limit"] != "10" {
			t.Errorf("limit = %q, want %q", params["limit"], "10")
		}
	})

	t.Run("pageSize capped at 1000", func(t *testing.T) {
		params := paginationParams(1, 2000)
		if params["limit"] != "1000" {
			t.Errorf("limit = %q, want %q", params["limit"], "1000")
		}
	})

	t.Run("negative page defaults to 1", func(t *testing.T) {
		params := paginationParams(-5, 25)
		if params["offset"] != "0" {
			t.Errorf("offset = %q, want %q", params["offset"], "0")
		}
	})

	t.Run("negative pageSize defaults to 25", func(t *testing.T) {
		params := paginationParams(1, -5)
		if params["limit"] != "25" {
			t.Errorf("limit = %q, want %q", params["limit"], "25")
		}
	})
}

func TestAddIntParam(t *testing.T) {
	t.Parallel()

	t.Run("non-zero value added", func(t *testing.T) {
		m := make(map[string]string)
		addIntParam(m, "family", 4)
		if m["family"] != "4" {
			t.Errorf("family = %q, want %q", m["family"], "4")
		}
	})

	t.Run("zero value skipped", func(t *testing.T) {
		m := make(map[string]string)
		addIntParam(m, "family", 0)
		if _, ok := m["family"]; ok {
			t.Error("family param should not be set for zero value")
		}
	})
}

func TestResolveService(t *testing.T) {
	t.Parallel()

	t.Run("returns context service when present", func(t *testing.T) {
		svc := application.NewNetworkService(&stubRepo{}, "ctx-token")
		ctx := context.WithValue(context.Background(), svcContextKey, svc)

		got := resolveService(ctx, nil)
		if got == nil {
			t.Fatal("resolveService() = nil, want non-nil")
		}
	})

	t.Run("returns fallback svc when context has none", func(t *testing.T) {
		fallback := application.NewNetworkService(&stubRepo{}, "fallback-token")
		got := resolveService(context.Background(), fallback)
		if got == nil {
			t.Fatal("resolveService() = nil, want non-nil")
		}
	})

	t.Run("returns nil when neither context nor fallback has service", func(t *testing.T) {
		got := resolveService(context.Background(), nil)
		if got != nil {
			t.Fatal("resolveService() = non-nil, want nil")
		}
	})
}

func TestGetObjectByID_CRLFSanitization(t *testing.T) {
	t.Parallel()

	t.Run("strips CRLF from param keys", func(t *testing.T) {
		svc := application.NewNetworkService(&stubRepo{}, "token")
		handler := NewGetObjectByIDHandler(svc)

		// Should not panic or produce CRLF in the URL.
		_, _, _ = handler(context.Background(), nil, GetObjectInput{
			ObjectType: "site",
			ID:         1,
			Params: map[string]string{
				"key\r\ninjected": "value",
			},
		})
	})

	t.Run("strips CRLF from param values", func(t *testing.T) {
		svc := application.NewNetworkService(&stubRepo{}, "token")
		handler := NewGetObjectByIDHandler(svc)

		// Should not panic or produce CRLF in the URL.
		_, _, _ = handler(context.Background(), nil, GetObjectInput{
			ObjectType: "site",
			ID:         1,
			Params: map[string]string{
				"key": "value\r\ninjected",
			},
		})
	})

	t.Run("strips nested CRLF from multiple params", func(t *testing.T) {
		svc := application.NewNetworkService(&stubRepo{}, "token")
		handler := NewGetObjectByIDHandler(svc)

		// Should not panic or produce CRLF in the URL.
		_, _, _ = handler(context.Background(), nil, GetObjectInput{
			ObjectType: "site",
			ID:         1,
			Params: map[string]string{
				"key\r\n1":  "val\r\n1",
				"key\r\n2":  "val\r\n2",
				"clean_key": "clean_val",
			},
		})
	})

	t.Run("empty params map is handled gracefully", func(t *testing.T) {
		svc := application.NewNetworkService(&stubRepo{}, "token")
		handler := NewGetObjectByIDHandler(svc)

		// Nil params should not panic.
		_, _, _ = handler(context.Background(), nil, GetObjectInput{
			ObjectType: "site",
			ID:         1,
			Params:     nil,
		})
	})

	t.Run("object_type validation before params processing", func(t *testing.T) {
		svc := application.NewNetworkService(&stubRepo{}, "token")
		handler := NewGetObjectByIDHandler(svc)

		_, _, err := handler(context.Background(), nil, GetObjectInput{
			ObjectType: "",
			ID:         1,
		})
		if err == nil {
			t.Error("expected error for empty object_type")
		}
	})

	t.Run("id validation before params processing", func(t *testing.T) {
		svc := application.NewNetworkService(&stubRepo{}, "token")
		handler := NewGetObjectByIDHandler(svc)

		_, _, err := handler(context.Background(), nil, GetObjectInput{
			ObjectType: "site",
			ID:         0,
		})
		if err == nil {
			t.Error("expected error for zero id")
		}
	})
}

func TestAddParam(t *testing.T) {
	t.Parallel()

	t.Run("non-empty value added", func(t *testing.T) {
		m := make(map[string]string)
		addParam(m, "site", "dc1")
		if m["site"] != "dc1" {
			t.Errorf("site = %q, want %q", m["site"], "dc1")
		}
	})

	t.Run("empty value skipped", func(t *testing.T) {
		m := make(map[string]string)
		addParam(m, "site", "")
		if _, ok := m["site"]; ok {
			t.Error("site param should not be set for empty value")
		}
	})
}
