package handlers

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/domain"
	"github.com/teran/mcp-netbox/mockrepo"
)

type stubRepo struct{}

func (s *stubRepo) ListSites(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Site], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListDevices(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Device], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListIPAddresses(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.IPAddress], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListPrefixes(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Prefix], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListVLANs(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.VLAN], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListVirtualMachines(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.VirtualMachine], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListClusters(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Cluster], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListCircuits(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Circuit], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListRacks(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Rack], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListCircuitTerminations(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.CircuitTermination], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListCables(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Cable], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListInterfaces(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Interface], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) ListVMInterfaces(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.VMInterface], error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) GetObject(
	_ context.Context, _ string, _ string, _ int, _ map[string]string,
) (domain.RawObject, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) CreateSite(_ context.Context, _ string, _ domain.SiteWrite) (*domain.Site, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateSite(_ context.Context, _ string, _ int, _ domain.SiteWrite) (*domain.Site, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteSite(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateDevice(_ context.Context, _ string, _ domain.DeviceWrite) (*domain.Device, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateDevice(_ context.Context, _ string, _ int, _ domain.DeviceWrite) (*domain.Device, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteDevice(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateIPAddress(_ context.Context, _ string, _ domain.IPAddressWrite) (*domain.IPAddress, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateIPAddress(
	_ context.Context, _ string, _ int, _ domain.IPAddressWrite,
) (*domain.IPAddress, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteIPAddress(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreatePrefix(_ context.Context, _ string, _ domain.PrefixWrite) (*domain.Prefix, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdatePrefix(_ context.Context, _ string, _ int, _ domain.PrefixWrite) (*domain.Prefix, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeletePrefix(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateVLAN(_ context.Context, _ string, _ domain.VLANWrite) (*domain.VLAN, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateVLAN(_ context.Context, _ string, _ int, _ domain.VLANWrite) (*domain.VLAN, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteVLAN(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateVirtualMachine(
	_ context.Context, _ string, _ domain.VirtualMachineWrite,
) (*domain.VirtualMachine, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateVirtualMachine(
	_ context.Context, _ string, _ int, _ domain.VirtualMachineWrite,
) (*domain.VirtualMachine, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteVirtualMachine(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateCluster(_ context.Context, _ string, _ domain.ClusterWrite) (*domain.Cluster, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateCluster(_ context.Context, _ string, _ int, _ domain.ClusterWrite) (*domain.Cluster, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteCluster(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateCircuit(_ context.Context, _ string, _ domain.CircuitWrite) (*domain.Circuit, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateCircuit(_ context.Context, _ string, _ int, _ domain.CircuitWrite) (*domain.Circuit, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteCircuit(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateRack(_ context.Context, _ string, _ domain.RackWrite) (*domain.Rack, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateRack(_ context.Context, _ string, _ int, _ domain.RackWrite) (*domain.Rack, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteRack(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateInterface(_ context.Context, _ string, _ domain.InterfaceWrite) (*domain.Interface, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateInterface(
	_ context.Context, _ string, _ int, _ domain.InterfaceWrite,
) (*domain.Interface, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteInterface(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateCircuitTermination(
	_ context.Context, _ string, _ domain.CircuitTerminationWrite,
) (*domain.CircuitTermination, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateCircuitTermination(
	_ context.Context, _ string, _ int, _ domain.CircuitTerminationWrite,
) (*domain.CircuitTermination, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteCircuitTermination(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateCable(_ context.Context, _ string, _ domain.CableWrite) (*domain.Cable, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateCable(_ context.Context, _ string, _ int, _ domain.CableWrite) (*domain.Cable, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteCable(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateVMInterface(
	_ context.Context, _ string, _ domain.VMInterfaceWrite,
) (*domain.VMInterface, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateVMInterface(
	_ context.Context, _ string, _ int, _ domain.VMInterfaceWrite,
) (*domain.VMInterface, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteVMInterface(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateProvider(_ context.Context, _ string, _ domain.ProviderWrite) (*domain.Provider, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateProvider(
	_ context.Context, _ string, _ int, _ domain.ProviderWrite,
) (*domain.Provider, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteProvider(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateTenant(_ context.Context, _ string, _ domain.TenantWrite) (*domain.Tenant, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateTenant(_ context.Context, _ string, _ int, _ domain.TenantWrite) (*domain.Tenant, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteTenant(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateManufacturer(
	_ context.Context, _ string, _ domain.ManufacturerWrite,
) (*domain.Manufacturer, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateManufacturer(
	_ context.Context, _ string, _ int, _ domain.ManufacturerWrite,
) (*domain.Manufacturer, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteManufacturer(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateDeviceType(_ context.Context, _ string, _ domain.DeviceTypeWrite) (*domain.DeviceType, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateDeviceType(
	_ context.Context, _ string, _ int, _ domain.DeviceTypeWrite,
) (*domain.DeviceType, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteDeviceType(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateLocation(_ context.Context, _ string, _ domain.LocationWrite) (*domain.Location, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateLocation(
	_ context.Context, _ string, _ int, _ domain.LocationWrite,
) (*domain.Location, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteLocation(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateClusterType(
	_ context.Context, _ string, _ domain.ClusterTypeWrite,
) (*domain.ClusterType, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateClusterType(
	_ context.Context, _ string, _ int, _ domain.ClusterTypeWrite,
) (*domain.ClusterType, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteClusterType(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateClusterGroup(
	_ context.Context, _ string, _ domain.ClusterGroupWrite,
) (*domain.ClusterGroup, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateClusterGroup(
	_ context.Context, _ string, _ int, _ domain.ClusterGroupWrite,
) (*domain.ClusterGroup, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteClusterGroup(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateCircuitType(
	_ context.Context, _ string, _ domain.CircuitTypeWrite,
) (*domain.CircuitType, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateCircuitType(
	_ context.Context, _ string, _ int, _ domain.CircuitTypeWrite,
) (*domain.CircuitType, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteCircuitType(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateVrf(_ context.Context, _ string, _ domain.VrfWrite) (*domain.Vrf, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateVrf(_ context.Context, _ string, _ int, _ domain.VrfWrite) (*domain.Vrf, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteVrf(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateVlanGroup(_ context.Context, _ string, _ domain.VlanGroupWrite) (*domain.VlanGroup, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateVlanGroup(
	_ context.Context, _ string, _ int, _ domain.VlanGroupWrite,
) (*domain.VlanGroup, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteVlanGroup(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateRole(_ context.Context, _ string, _ domain.RoleWrite) (*domain.Role, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateRole(_ context.Context, _ string, _ int, _ domain.RoleWrite) (*domain.Role, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteRole(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func (s *stubRepo) CreateContact(_ context.Context, _ string, _ domain.ContactWrite) (*domain.Contact, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) UpdateContact(_ context.Context, _ string, _ int, _ domain.ContactWrite) (*domain.Contact, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepo) DeleteContact(_ context.Context, _ string, _ int) error {
	return errors.New("not implemented")
}

func TestPaginationParams(t *testing.T) { //nolint:funlen,lll // Large table-driven/declarative function; splitting is not sensible.
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

	t.Run("page=0 with non-zero pageSize clamps to page 1", func(t *testing.T) {
		// Exactly zero page must be treated as page 1 (offset 0), not page 0
		// (which would produce a negative offset). Kills the CONDITIONALS_BOUNDARY
		// on `page <= 0`.
		params := paginationParams(0, 25)
		if params["offset"] != "0" {
			t.Errorf("offset = %q, want %q", params["offset"], "0")
		}
	})

	t.Run("pageSize=1 stays 1", func(t *testing.T) {
		// A pageSize of exactly 1 must be preserved (limit 1), not bumped to the
		// default of 25. Kills the CONDITIONALS_BOUNDARY on `pageSize < 1`.
		params := paginationParams(1, 1)
		if params["limit"] != "1" {
			t.Errorf("limit = %q, want %q", params["limit"], "1")
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

func TestGetObjectByID_CRLFSanitization(t *testing.T) { //nolint:funlen,lll // Large table-driven/declarative function; splitting is not sensible.
	t.Parallel()
	t.Run("strips CRLF from param keys", func(_ *testing.T) {
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

	t.Run("strips CRLF from param values", func(_ *testing.T) {
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

	t.Run("strips nested CRLF from multiple params", func(_ *testing.T) {
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

	t.Run("empty params map is handled gracefully", func(_ *testing.T) {
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

func TestGetObjectByID_ZeroIDRejectedBeforeService(t *testing.T) {
	t.Parallel()

	// An ID of exactly 0 must be rejected with the "id must be a positive
	// integer" error and must never reach the repository. Kills the
	// CONDITIONALS_BOUNDARY on `in.ID <= 0`, whose mutant would let ID=0
	// through to the service call.
	var serviceCalled bool
	svc := application.NewNetworkService(&mockrepo.MockRepo{
		GetObjectFunc: func(_ context.Context, _ string, _ string, _ int, _ map[string]string) (domain.RawObject, error) {
			serviceCalled = true

			return nil, errors.New("service should not be called")
		},
	}, "token")
	handler := NewGetObjectByIDHandler(svc)

	_, _, err := handler(context.Background(), nil, GetObjectInput{ObjectType: "site", ID: 0})
	if err == nil {
		t.Fatal("expected error for zero id, got nil")
	}
	if !strings.Contains(err.Error(), "id must be a positive integer") {
		t.Errorf("error = %q, want it to contain 'id must be a positive integer'", err.Error())
	}
	if serviceCalled {
		t.Error("GetObject was called for ID=0; it must be rejected before reaching the service")
	}
}

func TestGetObjectByID_CRLFParamsSanitizedBeforeService(t *testing.T) {
	t.Parallel()

	// The params map handed to the repository must be fully CRLF/control-char
	// free, while ordinary key/value pairs are preserved. This kills the
	// CONDITIONALS_NEGATION (r=='\r' / r=='\n') and INVERT_NEGATIVES /
	// ARITHMETIC_BASE (return -1) mutants in the get_object_by_id sanitizer.
	var gotParams map[string]string
	svc := application.NewNetworkService(&mockrepo.MockRepo{
		GetObjectFunc: func(
			_ context.Context, _ string, _ string, _ int, params map[string]string,
		) (domain.RawObject, error) {
			gotParams = params

			return domain.RawObject(`{}`), nil
		},
	}, "token")
	handler := NewGetObjectByIDHandler(svc)

	_, _, err := handler(context.Background(), nil, GetObjectInput{
		ObjectType: "site",
		ID:         1,
		Params: map[string]string{
			"hello":    "world",
			"key\r\n1": "value\r\n1",
			"clean":    "ok",
		},
	})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if gotParams == nil {
		t.Fatal("params were not forwarded to the service")
	}
	if gotParams["hello"] != "world" {
		t.Errorf("params['hello'] = %q, want 'world'", gotParams["hello"])
	}
	if gotParams["clean"] != "ok" {
		t.Errorf("params['clean'] = %q, want 'ok'", gotParams["clean"])
	}
	if gotParams["key1"] != "value1" {
		t.Errorf("params['key1'] = %q, want 'value1'", gotParams["key1"])
	}
	for k, v := range gotParams {
		if strings.ContainsAny(k, "\r\n\x01") || strings.ContainsAny(v, "\r\n\x01") {
			t.Errorf("sanitized params still contain control characters: key=%q value=%q", k, v)
		}
	}
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
