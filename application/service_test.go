package application

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/teran/mcp-netbox/domain"
	"github.com/teran/mcp-netbox/internal/mockrepo"
)

func newTestService(repo *mockrepo.MockRepo) *NetworkService {
	return NewNetworkService(repo, "test-token")
}

func TestNetworkService_ListSites(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListSitesFunc: func(_ context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Site], error) {
				if token != "test-token" {
					t.Errorf("token = %q, want %q", token, "test-token")
				}
				return &domain.PaginatedResponse[domain.Site]{
					Count:   1,
					Results: []domain.Site{{ID: 1, Name: "Test"}},
				}, nil
			},
		})
		resp, err := svc.ListSites(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListSites() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("Count = %d, want %d", resp.Count, 1)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListSitesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Site], error) {
				return nil, errors.New("repo error")
			},
		})
		_, err := svc.ListSites(context.Background(), nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestNetworkService_ListDevices(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListDevicesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Device], error) {
				return &domain.PaginatedResponse[domain.Device]{Count: 0}, nil
			},
		})
		resp, err := svc.ListDevices(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListDevices() returned error: %v", err)
		}
		if resp.Count != 0 {
			t.Errorf("Count = %d, want %d", resp.Count, 0)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListDevicesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Device], error) {
				return nil, errors.New("repo error")
			},
		})
		_, err := svc.ListDevices(context.Background(), nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestNetworkService_ListIPAddresses(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListIPAddressesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.IPAddress], error) {
				return &domain.PaginatedResponse[domain.IPAddress]{Count: 0}, nil
			},
		})
		resp, err := svc.ListIPAddresses(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListIPAddresses() returned error: %v", err)
		}
		if resp.Count != 0 {
			t.Errorf("Count = %d, want %d", resp.Count, 0)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListIPAddressesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.IPAddress], error) {
				return nil, errors.New("repo error")
			},
		})
		_, err := svc.ListIPAddresses(context.Background(), nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestNetworkService_ListPrefixes(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListPrefixesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Prefix], error) {
				return &domain.PaginatedResponse[domain.Prefix]{Count: 0}, nil
			},
		})
		resp, err := svc.ListPrefixes(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListPrefixes() returned error: %v", err)
		}
		if resp.Count != 0 {
			t.Errorf("Count = %d, want %d", resp.Count, 0)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListPrefixesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Prefix], error) {
				return nil, errors.New("repo error")
			},
		})
		_, err := svc.ListPrefixes(context.Background(), nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestNetworkService_ListVLANs(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListVLANsFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.VLAN], error) {
				return &domain.PaginatedResponse[domain.VLAN]{Count: 0}, nil
			},
		})
		resp, err := svc.ListVLANs(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListVLANs() returned error: %v", err)
		}
		if resp.Count != 0 {
			t.Errorf("Count = %d, want %d", resp.Count, 0)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListVLANsFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.VLAN], error) {
				return nil, errors.New("repo error")
			},
		})
		_, err := svc.ListVLANs(context.Background(), nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestNetworkService_ListVirtualMachines(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListVirtualMachinesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.VirtualMachine], error) {
				return &domain.PaginatedResponse[domain.VirtualMachine]{Count: 0}, nil
			},
		})
		resp, err := svc.ListVirtualMachines(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListVirtualMachines() returned error: %v", err)
		}
		if resp.Count != 0 {
			t.Errorf("Count = %d, want %d", resp.Count, 0)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListVirtualMachinesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.VirtualMachine], error) {
				return nil, errors.New("repo error")
			},
		})
		_, err := svc.ListVirtualMachines(context.Background(), nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestNetworkService_ListClusters(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListClustersFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Cluster], error) {
				return &domain.PaginatedResponse[domain.Cluster]{Count: 0}, nil
			},
		})
		resp, err := svc.ListClusters(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListClusters() returned error: %v", err)
		}
		if resp.Count != 0 {
			t.Errorf("Count = %d, want %d", resp.Count, 0)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListClustersFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Cluster], error) {
				return nil, errors.New("repo error")
			},
		})
		_, err := svc.ListClusters(context.Background(), nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestNetworkService_ListCircuits(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListCircuitsFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Circuit], error) {
				return &domain.PaginatedResponse[domain.Circuit]{Count: 0}, nil
			},
		})
		resp, err := svc.ListCircuits(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListCircuits() returned error: %v", err)
		}
		if resp.Count != 0 {
			t.Errorf("Count = %d, want %d", resp.Count, 0)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListCircuitsFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Circuit], error) {
				return nil, errors.New("repo error")
			},
		})
		_, err := svc.ListCircuits(context.Background(), nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestNetworkService_ListRacks(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListRacksFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Rack], error) {
				return &domain.PaginatedResponse[domain.Rack]{Count: 0}, nil
			},
		})
		resp, err := svc.ListRacks(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListRacks() returned error: %v", err)
		}
		if resp.Count != 0 {
			t.Errorf("Count = %d, want %d", resp.Count, 0)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListRacksFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Rack], error) {
				return nil, errors.New("repo error")
			},
		})
		_, err := svc.ListRacks(context.Background(), nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestNetworkService_ListInterfaces(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListInterfacesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Interface], error) {
				return &domain.PaginatedResponse[domain.Interface]{Count: 0}, nil
			},
		})
		resp, err := svc.ListInterfaces(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListInterfaces() returned error: %v", err)
		}
		if resp.Count != 0 {
			t.Errorf("Count = %d, want %d", resp.Count, 0)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListInterfacesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Interface], error) {
				return nil, errors.New("repo error")
			},
		})
		_, err := svc.ListInterfaces(context.Background(), nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestNetworkService_ListVMInterfaces(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListVMInterfacesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.VMInterface], error) {
				return &domain.PaginatedResponse[domain.VMInterface]{Count: 0}, nil
			},
		})
		resp, err := svc.ListVMInterfaces(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListVMInterfaces() returned error: %v", err)
		}
		if resp.Count != 0 {
			t.Errorf("Count = %d, want %d", resp.Count, 0)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListVMInterfacesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.VMInterface], error) {
				return nil, errors.New("repo error")
			},
		})
		_, err := svc.ListVMInterfaces(context.Background(), nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestNetworkService_ListCircuitTerminations(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListCircuitTerminationsFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.CircuitTermination], error) {
				return &domain.PaginatedResponse[domain.CircuitTermination]{Count: 0}, nil
			},
		})
		resp, err := svc.ListCircuitTerminations(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListCircuitTerminations() returned error: %v", err)
		}
		if resp.Count != 0 {
			t.Errorf("Count = %d, want %d", resp.Count, 0)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListCircuitTerminationsFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.CircuitTermination], error) {
				return nil, errors.New("repo error")
			},
		})
		_, err := svc.ListCircuitTerminations(context.Background(), nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestNetworkService_ListCables(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListCablesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Cable], error) {
				return &domain.PaginatedResponse[domain.Cable]{Count: 0}, nil
			},
		})
		resp, err := svc.ListCables(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListCables() returned error: %v", err)
		}
		if resp.Count != 0 {
			t.Errorf("Count = %d, want %d", resp.Count, 0)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListCablesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Cable], error) {
				return nil, errors.New("repo error")
			},
		})
		_, err := svc.ListCables(context.Background(), nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestNetworkService_GetObject(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			GetObjectFunc: func(_ context.Context, _ string, objectType string, id int, _ map[string]string) (domain.RawObject, error) {
				return domain.RawObject(fmt.Sprintf(`{"id":%d,"name":"%s"}`, id, objectType)), nil
			},
		})
		resp, err := svc.GetObject(context.Background(), "site", 1, nil)
		if err != nil {
			t.Fatalf("GetObject() returned error: %v", err)
		}
		if len(resp) == 0 {
			t.Fatal("GetObject() returned empty response")
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			GetObjectFunc: func(_ context.Context, _ string, _ string, _ int, _ map[string]string) (domain.RawObject, error) {
				return nil, errors.New("repo error")
			},
		})
		_, err := svc.GetObject(context.Background(), "site", 1, nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestNetworkService_CreateSite(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateSiteFunc: func(_ context.Context, token string, in domain.SiteWrite) (*domain.Site, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.Site{ID: 1, Name: in.Name}, nil
			},
		})
		site, err := svc.CreateSite(context.Background(), domain.SiteWrite{Name: "A"})
		if err != nil {
			t.Fatalf("CreateSite() returned error: %v", err)
		}
		if site.ID != 1 || site.Name != "A" {
			t.Errorf("site = %+v, want id 1 name A", site)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"name":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreateSiteFunc: func(_ context.Context, _ string, _ domain.SiteWrite) (*domain.Site, error) {
				return nil, ve
			},
		})
		_, err := svc.CreateSite(context.Background(), domain.SiteWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
		if string(got.Body) != `{"name":["required"]}` {
			t.Errorf("Body = %s, want original body", got.Body)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateSiteFunc: func(_ context.Context, _ string, _ domain.SiteWrite) (*domain.Site, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateSite(context.Background(), domain.SiteWrite{Name: "A"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !containsService(err.Error(), "create site:") {
			t.Errorf("err = %q, want it to contain 'create site:'", err.Error())
		}
	})
}

func TestNetworkService_UpdateSite(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateSiteFunc: func(_ context.Context, _ string, id int, in domain.SiteWrite) (*domain.Site, error) {
				return &domain.Site{ID: id, Name: in.Name}, nil
			},
		})
		site, err := svc.UpdateSite(context.Background(), 7, domain.SiteWrite{Name: "Renamed"})
		if err != nil {
			t.Fatalf("UpdateSite() returned error: %v", err)
		}
		if site.ID != 7 || site.Name != "Renamed" {
			t.Errorf("site = %+v, want id 7 name Renamed", site)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"slug":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdateSiteFunc: func(_ context.Context, _ string, _ int, _ domain.SiteWrite) (*domain.Site, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdateSite(context.Background(), 7, domain.SiteWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateSiteFunc: func(_ context.Context, _ string, _ int, _ domain.SiteWrite) (*domain.Site, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateSite(context.Background(), 7, domain.SiteWrite{})
		if err == nil || !containsService(err.Error(), "update site:") {
			t.Errorf("err = %v, want it to contain 'update site:'", err)
		}
	})
}

func TestNetworkService_DeleteSite(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteSiteFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeleteSite(context.Background(), 3); err != nil {
			t.Fatalf("DeleteSite() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeleteSiteFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeleteSite(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteSiteFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeleteSite(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete site:") {
			t.Errorf("err = %v, want it to contain 'delete site:'", err)
		}
	})
}

func TestNetworkService_CreateDevice(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateDeviceFunc: func(_ context.Context, token string, in domain.DeviceWrite) (*domain.Device, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.Device{ID: 1, Name: in.Name}, nil
			},
		})
		device, err := svc.CreateDevice(context.Background(), domain.DeviceWrite{Name: "A"})
		if err != nil {
			t.Fatalf("CreateDevice() returned error: %v", err)
		}
		if device.ID != 1 || device.Name != "A" {
			t.Errorf("device = %+v, want id 1 name A", device)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"name":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreateDeviceFunc: func(_ context.Context, _ string, _ domain.DeviceWrite) (*domain.Device, error) {
				return nil, ve
			},
		})
		_, err := svc.CreateDevice(context.Background(), domain.DeviceWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
		if string(got.Body) != `{"name":["required"]}` {
			t.Errorf("Body = %s, want original body", got.Body)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateDeviceFunc: func(_ context.Context, _ string, _ domain.DeviceWrite) (*domain.Device, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateDevice(context.Background(), domain.DeviceWrite{Name: "A"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !containsService(err.Error(), "create device:") {
			t.Errorf("err = %q, want it to contain 'create device:'", err.Error())
		}
	})
}

func TestNetworkService_UpdateDevice(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateDeviceFunc: func(_ context.Context, _ string, id int, in domain.DeviceWrite) (*domain.Device, error) {
				return &domain.Device{ID: id, Name: in.Name}, nil
			},
		})
		device, err := svc.UpdateDevice(context.Background(), 7, domain.DeviceWrite{Name: "Renamed"})
		if err != nil {
			t.Fatalf("UpdateDevice() returned error: %v", err)
		}
		if device.ID != 7 || device.Name != "Renamed" {
			t.Errorf("device = %+v, want id 7 name Renamed", device)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"role":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdateDeviceFunc: func(_ context.Context, _ string, _ int, _ domain.DeviceWrite) (*domain.Device, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdateDevice(context.Background(), 7, domain.DeviceWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateDeviceFunc: func(_ context.Context, _ string, _ int, _ domain.DeviceWrite) (*domain.Device, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateDevice(context.Background(), 7, domain.DeviceWrite{})
		if err == nil || !containsService(err.Error(), "update device:") {
			t.Errorf("err = %v, want it to contain 'update device:'", err)
		}
	})
}

func TestNetworkService_DeleteDevice(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteDeviceFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeleteDevice(context.Background(), 3); err != nil {
			t.Fatalf("DeleteDevice() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeleteDeviceFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeleteDevice(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteDeviceFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeleteDevice(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete device:") {
			t.Errorf("err = %v, want it to contain 'delete device:'", err)
		}
	})
}

func TestNetworkService_CreateIPAddress(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateIPAddressFunc: func(_ context.Context, token string, in domain.IPAddressWrite) (*domain.IPAddress, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.IPAddress{ID: 1, Address: in.Address}, nil
			},
		})
		ip, err := svc.CreateIPAddress(context.Background(), domain.IPAddressWrite{Address: "A"})
		if err != nil {
			t.Fatalf("CreateIPAddress() returned error: %v", err)
		}
		if ip.ID != 1 || ip.Address != "A" {
			t.Errorf("ip = %+v, want id 1 address A", ip)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"address":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreateIPAddressFunc: func(_ context.Context, _ string, _ domain.IPAddressWrite) (*domain.IPAddress, error) {
				return nil, ve
			},
		})
		_, err := svc.CreateIPAddress(context.Background(), domain.IPAddressWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
		if string(got.Body) != `{"address":["required"]}` {
			t.Errorf("Body = %s, want original body", got.Body)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateIPAddressFunc: func(_ context.Context, _ string, _ domain.IPAddressWrite) (*domain.IPAddress, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateIPAddress(context.Background(), domain.IPAddressWrite{Address: "A"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !containsService(err.Error(), "create IP address:") {
			t.Errorf("err = %q, want it to contain 'create IP address:'", err.Error())
		}
	})
}

func TestNetworkService_UpdateIPAddress(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateIPAddressFunc: func(_ context.Context, _ string, id int, in domain.IPAddressWrite) (*domain.IPAddress, error) {
				return &domain.IPAddress{ID: id, Address: in.Address}, nil
			},
		})
		ip, err := svc.UpdateIPAddress(context.Background(), 7, domain.IPAddressWrite{Address: "Renamed"})
		if err != nil {
			t.Fatalf("UpdateIPAddress() returned error: %v", err)
		}
		if ip.ID != 7 || ip.Address != "Renamed" {
			t.Errorf("ip = %+v, want id 7 address Renamed", ip)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"dns_name":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdateIPAddressFunc: func(_ context.Context, _ string, _ int, _ domain.IPAddressWrite) (*domain.IPAddress, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdateIPAddress(context.Background(), 7, domain.IPAddressWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateIPAddressFunc: func(_ context.Context, _ string, _ int, _ domain.IPAddressWrite) (*domain.IPAddress, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateIPAddress(context.Background(), 7, domain.IPAddressWrite{})
		if err == nil || !containsService(err.Error(), "update IP address:") {
			t.Errorf("err = %v, want it to contain 'update IP address:'", err)
		}
	})
}

func TestNetworkService_DeleteIPAddress(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteIPAddressFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeleteIPAddress(context.Background(), 3); err != nil {
			t.Fatalf("DeleteIPAddress() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeleteIPAddressFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeleteIPAddress(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteIPAddressFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeleteIPAddress(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete IP address:") {
			t.Errorf("err = %v, want it to contain 'delete IP address:'", err)
		}
	})
}

func TestNetworkService_CreatePrefix(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreatePrefixFunc: func(_ context.Context, token string, in domain.PrefixWrite) (*domain.Prefix, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.Prefix{ID: 1, Prefix: in.Prefix}, nil
			},
		})
		prefix, err := svc.CreatePrefix(context.Background(), domain.PrefixWrite{Prefix: "10.0.0.0/24"})
		if err != nil {
			t.Fatalf("CreatePrefix() returned error: %v", err)
		}
		if prefix.ID != 1 || prefix.Prefix != "10.0.0.0/24" {
			t.Errorf("prefix = %+v, want id 1 prefix 10.0.0.0/24", prefix)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"prefix":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreatePrefixFunc: func(_ context.Context, _ string, _ domain.PrefixWrite) (*domain.Prefix, error) {
				return nil, ve
			},
		})
		_, err := svc.CreatePrefix(context.Background(), domain.PrefixWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreatePrefixFunc: func(_ context.Context, _ string, _ domain.PrefixWrite) (*domain.Prefix, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreatePrefix(context.Background(), domain.PrefixWrite{Prefix: "A"})
		if err == nil || !containsService(err.Error(), "create prefix:") {
			t.Errorf("err = %v, want it to contain 'create prefix:'", err)
		}
	})
}

func TestNetworkService_UpdatePrefix(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdatePrefixFunc: func(_ context.Context, _ string, id int, in domain.PrefixWrite) (*domain.Prefix, error) {
				return &domain.Prefix{ID: id, Prefix: in.Prefix}, nil
			},
		})
		prefix, err := svc.UpdatePrefix(context.Background(), 7, domain.PrefixWrite{Prefix: "10.0.0.0/24"})
		if err != nil {
			t.Fatalf("UpdatePrefix() returned error: %v", err)
		}
		if prefix.ID != 7 {
			t.Errorf("prefix = %+v, want id 7", prefix)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"status":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdatePrefixFunc: func(_ context.Context, _ string, _ int, _ domain.PrefixWrite) (*domain.Prefix, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdatePrefix(context.Background(), 7, domain.PrefixWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdatePrefixFunc: func(_ context.Context, _ string, _ int, _ domain.PrefixWrite) (*domain.Prefix, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdatePrefix(context.Background(), 7, domain.PrefixWrite{})
		if err == nil || !containsService(err.Error(), "update prefix:") {
			t.Errorf("err = %v, want it to contain 'update prefix:'", err)
		}
	})
}

func TestNetworkService_DeletePrefix(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeletePrefixFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeletePrefix(context.Background(), 3); err != nil {
			t.Fatalf("DeletePrefix() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeletePrefixFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeletePrefix(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeletePrefixFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeletePrefix(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete prefix:") {
			t.Errorf("err = %v, want it to contain 'delete prefix:'", err)
		}
	})
}

func TestNetworkService_CreateVLAN(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVLANFunc: func(_ context.Context, token string, in domain.VLANWrite) (*domain.VLAN, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.VLAN{ID: 1, VID: in.VID, Name: in.Name}, nil
			},
		})
		vlan, err := svc.CreateVLAN(context.Background(), domain.VLANWrite{VID: 100, Name: "mgmt"})
		if err != nil {
			t.Fatalf("CreateVLAN() returned error: %v", err)
		}
		if vlan.ID != 1 || vlan.VID != 100 || vlan.Name != "mgmt" {
			t.Errorf("vlan = %+v, want id 1 vid 100 name mgmt", vlan)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"name":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreateVLANFunc: func(_ context.Context, _ string, _ domain.VLANWrite) (*domain.VLAN, error) {
				return nil, ve
			},
		})
		_, err := svc.CreateVLAN(context.Background(), domain.VLANWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVLANFunc: func(_ context.Context, _ string, _ domain.VLANWrite) (*domain.VLAN, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateVLAN(context.Background(), domain.VLANWrite{VID: 100, Name: "mgmt"})
		if err == nil || !containsService(err.Error(), "create VLAN:") {
			t.Errorf("err = %v, want it to contain 'create VLAN:'", err)
		}
	})
}

func TestNetworkService_UpdateVLAN(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVLANFunc: func(_ context.Context, _ string, id int, in domain.VLANWrite) (*domain.VLAN, error) {
				return &domain.VLAN{ID: id, VID: in.VID, Name: in.Name}, nil
			},
		})
		vlan, err := svc.UpdateVLAN(context.Background(), 7, domain.VLANWrite{Name: "Renamed"})
		if err != nil {
			t.Fatalf("UpdateVLAN() returned error: %v", err)
		}
		if vlan.ID != 7 {
			t.Errorf("vlan = %+v, want id 7", vlan)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"status":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVLANFunc: func(_ context.Context, _ string, _ int, _ domain.VLANWrite) (*domain.VLAN, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdateVLAN(context.Background(), 7, domain.VLANWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVLANFunc: func(_ context.Context, _ string, _ int, _ domain.VLANWrite) (*domain.VLAN, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateVLAN(context.Background(), 7, domain.VLANWrite{})
		if err == nil || !containsService(err.Error(), "update VLAN:") {
			t.Errorf("err = %v, want it to contain 'update VLAN:'", err)
		}
	})
}

func TestNetworkService_DeleteVLAN(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVLANFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeleteVLAN(context.Background(), 3); err != nil {
			t.Fatalf("DeleteVLAN() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVLANFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeleteVLAN(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVLANFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeleteVLAN(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete VLAN:") {
			t.Errorf("err = %v, want it to contain 'delete VLAN:'", err)
		}
	})
}

func TestNetworkService_CreateVirtualMachine(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVirtualMachineFunc: func(_ context.Context, token string, in domain.VirtualMachineWrite) (*domain.VirtualMachine, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.VirtualMachine{ID: 1, Name: in.Name}, nil
			},
		})
		vm, err := svc.CreateVirtualMachine(context.Background(), domain.VirtualMachineWrite{Name: "web-01"})
		if err != nil {
			t.Fatalf("CreateVirtualMachine() returned error: %v", err)
		}
		if vm.ID != 1 || vm.Name != "web-01" {
			t.Errorf("vm = %+v, want id 1 name web-01", vm)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"name":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreateVirtualMachineFunc: func(_ context.Context, _ string, _ domain.VirtualMachineWrite) (*domain.VirtualMachine, error) {
				return nil, ve
			},
		})
		_, err := svc.CreateVirtualMachine(context.Background(), domain.VirtualMachineWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVirtualMachineFunc: func(_ context.Context, _ string, _ domain.VirtualMachineWrite) (*domain.VirtualMachine, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateVirtualMachine(context.Background(), domain.VirtualMachineWrite{Name: "web-01"})
		if err == nil || !containsService(err.Error(), "create virtual machine:") {
			t.Errorf("err = %v, want it to contain 'create virtual machine:'", err)
		}
	})
}

func TestNetworkService_UpdateVirtualMachine(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVirtualMachineFunc: func(_ context.Context, _ string, id int, in domain.VirtualMachineWrite) (*domain.VirtualMachine, error) {
				return &domain.VirtualMachine{ID: id, Name: in.Name}, nil
			},
		})
		vm, err := svc.UpdateVirtualMachine(context.Background(), 7, domain.VirtualMachineWrite{Name: "Renamed"})
		if err != nil {
			t.Fatalf("UpdateVirtualMachine() returned error: %v", err)
		}
		if vm.ID != 7 {
			t.Errorf("vm = %+v, want id 7", vm)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"status":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVirtualMachineFunc: func(_ context.Context, _ string, _ int, _ domain.VirtualMachineWrite) (*domain.VirtualMachine, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdateVirtualMachine(context.Background(), 7, domain.VirtualMachineWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVirtualMachineFunc: func(_ context.Context, _ string, _ int, _ domain.VirtualMachineWrite) (*domain.VirtualMachine, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateVirtualMachine(context.Background(), 7, domain.VirtualMachineWrite{})
		if err == nil || !containsService(err.Error(), "update virtual machine:") {
			t.Errorf("err = %v, want it to contain 'update virtual machine:'", err)
		}
	})
}

func TestNetworkService_DeleteVirtualMachine(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVirtualMachineFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeleteVirtualMachine(context.Background(), 3); err != nil {
			t.Fatalf("DeleteVirtualMachine() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVirtualMachineFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeleteVirtualMachine(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVirtualMachineFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeleteVirtualMachine(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete virtual machine:") {
			t.Errorf("err = %v, want it to contain 'delete virtual machine:'", err)
		}
	})
}

func TestNetworkService_CreateCluster(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateClusterFunc: func(_ context.Context, token string, in domain.ClusterWrite) (*domain.Cluster, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.Cluster{ID: 1, Name: in.Name}, nil
			},
		})
		cluster, err := svc.CreateCluster(context.Background(), domain.ClusterWrite{Name: "prod"})
		if err != nil {
			t.Fatalf("CreateCluster() returned error: %v", err)
		}
		if cluster.ID != 1 || cluster.Name != "prod" {
			t.Errorf("cluster = %+v, want id 1 name prod", cluster)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"name":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreateClusterFunc: func(_ context.Context, _ string, _ domain.ClusterWrite) (*domain.Cluster, error) {
				return nil, ve
			},
		})
		_, err := svc.CreateCluster(context.Background(), domain.ClusterWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateClusterFunc: func(_ context.Context, _ string, _ domain.ClusterWrite) (*domain.Cluster, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateCluster(context.Background(), domain.ClusterWrite{Name: "prod"})
		if err == nil || !containsService(err.Error(), "create cluster:") {
			t.Errorf("err = %v, want it to contain 'create cluster:'", err)
		}
	})
}

func TestNetworkService_UpdateCluster(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateClusterFunc: func(_ context.Context, _ string, id int, in domain.ClusterWrite) (*domain.Cluster, error) {
				return &domain.Cluster{ID: id, Name: in.Name}, nil
			},
		})
		cluster, err := svc.UpdateCluster(context.Background(), 7, domain.ClusterWrite{Name: "Renamed"})
		if err != nil {
			t.Fatalf("UpdateCluster() returned error: %v", err)
		}
		if cluster.ID != 7 {
			t.Errorf("cluster = %+v, want id 7", cluster)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"type":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdateClusterFunc: func(_ context.Context, _ string, _ int, _ domain.ClusterWrite) (*domain.Cluster, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdateCluster(context.Background(), 7, domain.ClusterWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateClusterFunc: func(_ context.Context, _ string, _ int, _ domain.ClusterWrite) (*domain.Cluster, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateCluster(context.Background(), 7, domain.ClusterWrite{})
		if err == nil || !containsService(err.Error(), "update cluster:") {
			t.Errorf("err = %v, want it to contain 'update cluster:'", err)
		}
	})
}

func TestNetworkService_DeleteCluster(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteClusterFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeleteCluster(context.Background(), 3); err != nil {
			t.Fatalf("DeleteCluster() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeleteClusterFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeleteCluster(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteClusterFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeleteCluster(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete cluster:") {
			t.Errorf("err = %v, want it to contain 'delete cluster:'", err)
		}
	})
}

func TestNetworkService_CreateCircuit(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateCircuitFunc: func(_ context.Context, token string, in domain.CircuitWrite) (*domain.Circuit, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.Circuit{ID: 1, CID: in.CID}, nil
			},
		})
		circuit, err := svc.CreateCircuit(context.Background(), domain.CircuitWrite{CID: "CIR-001"})
		if err != nil {
			t.Fatalf("CreateCircuit() returned error: %v", err)
		}
		if circuit.ID != 1 || circuit.CID != "CIR-001" {
			t.Errorf("circuit = %+v, want id 1 cid CIR-001", circuit)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"cid":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreateCircuitFunc: func(_ context.Context, _ string, _ domain.CircuitWrite) (*domain.Circuit, error) {
				return nil, ve
			},
		})
		_, err := svc.CreateCircuit(context.Background(), domain.CircuitWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateCircuitFunc: func(_ context.Context, _ string, _ domain.CircuitWrite) (*domain.Circuit, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateCircuit(context.Background(), domain.CircuitWrite{CID: "CIR-001"})
		if err == nil || !containsService(err.Error(), "create circuit:") {
			t.Errorf("err = %v, want it to contain 'create circuit:'", err)
		}
	})
}

func TestNetworkService_UpdateCircuit(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateCircuitFunc: func(_ context.Context, _ string, id int, in domain.CircuitWrite) (*domain.Circuit, error) {
				return &domain.Circuit{ID: id, CID: in.CID}, nil
			},
		})
		circuit, err := svc.UpdateCircuit(context.Background(), 7, domain.CircuitWrite{CID: "Renamed"})
		if err != nil {
			t.Fatalf("UpdateCircuit() returned error: %v", err)
		}
		if circuit.ID != 7 {
			t.Errorf("circuit = %+v, want id 7", circuit)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"provider":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdateCircuitFunc: func(_ context.Context, _ string, _ int, _ domain.CircuitWrite) (*domain.Circuit, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdateCircuit(context.Background(), 7, domain.CircuitWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateCircuitFunc: func(_ context.Context, _ string, _ int, _ domain.CircuitWrite) (*domain.Circuit, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateCircuit(context.Background(), 7, domain.CircuitWrite{})
		if err == nil || !containsService(err.Error(), "update circuit:") {
			t.Errorf("err = %v, want it to contain 'update circuit:'", err)
		}
	})
}

func TestNetworkService_DeleteCircuit(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteCircuitFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeleteCircuit(context.Background(), 3); err != nil {
			t.Fatalf("DeleteCircuit() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeleteCircuitFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeleteCircuit(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteCircuitFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeleteCircuit(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete circuit:") {
			t.Errorf("err = %v, want it to contain 'delete circuit:'", err)
		}
	})
}

func TestNetworkService_CreateRack(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateRackFunc: func(_ context.Context, token string, in domain.RackWrite) (*domain.Rack, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.Rack{ID: 1, Name: in.Name}, nil
			},
		})
		rack, err := svc.CreateRack(context.Background(), domain.RackWrite{Name: "R1"})
		if err != nil {
			t.Fatalf("CreateRack() returned error: %v", err)
		}
		if rack.ID != 1 || rack.Name != "R1" {
			t.Errorf("rack = %+v, want id 1 name R1", rack)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"name":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreateRackFunc: func(_ context.Context, _ string, _ domain.RackWrite) (*domain.Rack, error) {
				return nil, ve
			},
		})
		_, err := svc.CreateRack(context.Background(), domain.RackWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateRackFunc: func(_ context.Context, _ string, _ domain.RackWrite) (*domain.Rack, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateRack(context.Background(), domain.RackWrite{Name: "R1"})
		if err == nil || !containsService(err.Error(), "create rack:") {
			t.Errorf("err = %v, want it to contain 'create rack:'", err)
		}
	})
}

func TestNetworkService_UpdateRack(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateRackFunc: func(_ context.Context, _ string, id int, in domain.RackWrite) (*domain.Rack, error) {
				return &domain.Rack{ID: id, Name: in.Name}, nil
			},
		})
		rack, err := svc.UpdateRack(context.Background(), 7, domain.RackWrite{Name: "Renamed"})
		if err != nil {
			t.Fatalf("UpdateRack() returned error: %v", err)
		}
		if rack.ID != 7 {
			t.Errorf("rack = %+v, want id 7", rack)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"status":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdateRackFunc: func(_ context.Context, _ string, _ int, _ domain.RackWrite) (*domain.Rack, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdateRack(context.Background(), 7, domain.RackWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateRackFunc: func(_ context.Context, _ string, _ int, _ domain.RackWrite) (*domain.Rack, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateRack(context.Background(), 7, domain.RackWrite{})
		if err == nil || !containsService(err.Error(), "update rack:") {
			t.Errorf("err = %v, want it to contain 'update rack:'", err)
		}
	})
}

func TestNetworkService_DeleteRack(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteRackFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeleteRack(context.Background(), 3); err != nil {
			t.Fatalf("DeleteRack() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeleteRackFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeleteRack(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteRackFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeleteRack(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete rack:") {
			t.Errorf("err = %v, want it to contain 'delete rack:'", err)
		}
	})
}

func TestNetworkService_CreateInterface(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateInterfaceFunc: func(_ context.Context, token string, in domain.InterfaceWrite) (*domain.Interface, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.Interface{ID: 1, Name: in.Name}, nil
			},
		})
		iface, err := svc.CreateInterface(context.Background(), domain.InterfaceWrite{Name: "eth0"})
		if err != nil {
			t.Fatalf("CreateInterface() returned error: %v", err)
		}
		if iface.ID != 1 || iface.Name != "eth0" {
			t.Errorf("iface = %+v, want id 1 name eth0", iface)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"name":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreateInterfaceFunc: func(_ context.Context, _ string, _ domain.InterfaceWrite) (*domain.Interface, error) {
				return nil, ve
			},
		})
		_, err := svc.CreateInterface(context.Background(), domain.InterfaceWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateInterfaceFunc: func(_ context.Context, _ string, _ domain.InterfaceWrite) (*domain.Interface, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateInterface(context.Background(), domain.InterfaceWrite{Name: "eth0"})
		if err == nil || !containsService(err.Error(), "create interface:") {
			t.Errorf("err = %v, want it to contain 'create interface:'", err)
		}
	})
}

func TestNetworkService_UpdateInterface(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateInterfaceFunc: func(_ context.Context, _ string, id int, in domain.InterfaceWrite) (*domain.Interface, error) {
				return &domain.Interface{ID: id, Name: in.Name}, nil
			},
		})
		iface, err := svc.UpdateInterface(context.Background(), 7, domain.InterfaceWrite{Name: "Renamed"})
		if err != nil {
			t.Fatalf("UpdateInterface() returned error: %v", err)
		}
		if iface.ID != 7 {
			t.Errorf("iface = %+v, want id 7", iface)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"type":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdateInterfaceFunc: func(_ context.Context, _ string, _ int, _ domain.InterfaceWrite) (*domain.Interface, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdateInterface(context.Background(), 7, domain.InterfaceWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateInterfaceFunc: func(_ context.Context, _ string, _ int, _ domain.InterfaceWrite) (*domain.Interface, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateInterface(context.Background(), 7, domain.InterfaceWrite{})
		if err == nil || !containsService(err.Error(), "update interface:") {
			t.Errorf("err = %v, want it to contain 'update interface:'", err)
		}
	})
}

func TestNetworkService_DeleteInterface(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteInterfaceFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeleteInterface(context.Background(), 3); err != nil {
			t.Fatalf("DeleteInterface() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeleteInterfaceFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeleteInterface(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteInterfaceFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeleteInterface(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete interface:") {
			t.Errorf("err = %v, want it to contain 'delete interface:'", err)
		}
	})
}

func TestNetworkService_CreateCircuitTermination(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateCircuitTerminationFunc: func(_ context.Context, token string, in domain.CircuitTerminationWrite) (*domain.CircuitTermination, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.CircuitTermination{ID: 1, TermSide: in.TermSide}, nil
			},
		})
		ct, err := svc.CreateCircuitTermination(context.Background(), domain.CircuitTerminationWrite{TermSide: "A"})
		if err != nil {
			t.Fatalf("CreateCircuitTermination() returned error: %v", err)
		}
		if ct.ID != 1 || ct.TermSide != "A" {
			t.Errorf("ct = %+v, want id 1 term_side A", ct)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"term_side":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreateCircuitTerminationFunc: func(_ context.Context, _ string, _ domain.CircuitTerminationWrite) (*domain.CircuitTermination, error) {
				return nil, ve
			},
		})
		_, err := svc.CreateCircuitTermination(context.Background(), domain.CircuitTerminationWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateCircuitTerminationFunc: func(_ context.Context, _ string, _ domain.CircuitTerminationWrite) (*domain.CircuitTermination, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateCircuitTermination(context.Background(), domain.CircuitTerminationWrite{TermSide: "A"})
		if err == nil || !containsService(err.Error(), "create circuit termination:") {
			t.Errorf("err = %v, want it to contain 'create circuit termination:'", err)
		}
	})
}

func TestNetworkService_UpdateCircuitTermination(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateCircuitTerminationFunc: func(_ context.Context, _ string, id int, in domain.CircuitTerminationWrite) (*domain.CircuitTermination, error) {
				return &domain.CircuitTermination{ID: id, TermSide: in.TermSide}, nil
			},
		})
		ct, err := svc.UpdateCircuitTermination(context.Background(), 7, domain.CircuitTerminationWrite{TermSide: "Z"})
		if err != nil {
			t.Fatalf("UpdateCircuitTermination() returned error: %v", err)
		}
		if ct.ID != 7 {
			t.Errorf("ct = %+v, want id 7", ct)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"site":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdateCircuitTerminationFunc: func(_ context.Context, _ string, _ int, _ domain.CircuitTerminationWrite) (*domain.CircuitTermination, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdateCircuitTermination(context.Background(), 7, domain.CircuitTerminationWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateCircuitTerminationFunc: func(_ context.Context, _ string, _ int, _ domain.CircuitTerminationWrite) (*domain.CircuitTermination, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateCircuitTermination(context.Background(), 7, domain.CircuitTerminationWrite{})
		if err == nil || !containsService(err.Error(), "update circuit termination:") {
			t.Errorf("err = %v, want it to contain 'update circuit termination:'", err)
		}
	})
}

func TestNetworkService_DeleteCircuitTermination(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteCircuitTerminationFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeleteCircuitTermination(context.Background(), 3); err != nil {
			t.Fatalf("DeleteCircuitTermination() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeleteCircuitTerminationFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeleteCircuitTermination(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteCircuitTerminationFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeleteCircuitTermination(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete circuit termination:") {
			t.Errorf("err = %v, want it to contain 'delete circuit termination:'", err)
		}
	})
}

func TestNetworkService_CreateCable(t *testing.T) {
	t.Parallel()

	label := "link-01"

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateCableFunc: func(_ context.Context, token string, in domain.CableWrite) (*domain.Cable, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.Cable{ID: 1, Label: *in.Label}, nil
			},
		})
		c, err := svc.CreateCable(context.Background(), domain.CableWrite{Label: &label})
		if err != nil {
			t.Fatalf("CreateCable() returned error: %v", err)
		}
		if c.ID != 1 || c.Label != "link-01" {
			t.Errorf("c = %+v, want id 1 label link-01", c)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"termination_a":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreateCableFunc: func(_ context.Context, _ string, _ domain.CableWrite) (*domain.Cable, error) {
				return nil, ve
			},
		})
		_, err := svc.CreateCable(context.Background(), domain.CableWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateCableFunc: func(_ context.Context, _ string, _ domain.CableWrite) (*domain.Cable, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateCable(context.Background(), domain.CableWrite{Label: &label})
		if err == nil || !containsService(err.Error(), "create cable:") {
			t.Errorf("err = %v, want it to contain 'create cable:'", err)
		}
	})
}

func TestNetworkService_UpdateCable(t *testing.T) {
	t.Parallel()

	label := "link-02"

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateCableFunc: func(_ context.Context, _ string, id int, _ domain.CableWrite) (*domain.Cable, error) {
				return &domain.Cable{ID: id, Label: label}, nil
			},
		})
		c, err := svc.UpdateCable(context.Background(), 7, domain.CableWrite{Label: &label})
		if err != nil {
			t.Fatalf("UpdateCable() returned error: %v", err)
		}
		if c.ID != 7 {
			t.Errorf("c = %+v, want id 7", c)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"status":["invalid"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdateCableFunc: func(_ context.Context, _ string, _ int, _ domain.CableWrite) (*domain.Cable, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdateCable(context.Background(), 7, domain.CableWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateCableFunc: func(_ context.Context, _ string, _ int, _ domain.CableWrite) (*domain.Cable, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateCable(context.Background(), 7, domain.CableWrite{})
		if err == nil || !containsService(err.Error(), "update cable:") {
			t.Errorf("err = %v, want it to contain 'update cable:'", err)
		}
	})
}

func TestNetworkService_DeleteCable(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteCableFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeleteCable(context.Background(), 3); err != nil {
			t.Fatalf("DeleteCable() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeleteCableFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeleteCable(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteCableFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeleteCable(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete cable:") {
			t.Errorf("err = %v, want it to contain 'delete cable:'", err)
		}
	})
}

func TestNetworkService_CreateVMInterface(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVMInterfaceFunc: func(_ context.Context, token string, in domain.VMInterfaceWrite) (*domain.VMInterface, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.VMInterface{ID: 1, Name: in.Name}, nil
			},
		})
		vi, err := svc.CreateVMInterface(context.Background(), domain.VMInterfaceWrite{Name: "eth0"})
		if err != nil {
			t.Fatalf("CreateVMInterface() returned error: %v", err)
		}
		if vi.ID != 1 || vi.Name != "eth0" {
			t.Errorf("vi = %+v, want id 1 name eth0", vi)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"name":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreateVMInterfaceFunc: func(_ context.Context, _ string, _ domain.VMInterfaceWrite) (*domain.VMInterface, error) {
				return nil, ve
			},
		})
		_, err := svc.CreateVMInterface(context.Background(), domain.VMInterfaceWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVMInterfaceFunc: func(_ context.Context, _ string, _ domain.VMInterfaceWrite) (*domain.VMInterface, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateVMInterface(context.Background(), domain.VMInterfaceWrite{Name: "eth0"})
		if err == nil || !containsService(err.Error(), "create vm interface:") {
			t.Errorf("err = %v, want it to contain 'create vm interface:'", err)
		}
	})
}

func TestNetworkService_UpdateVMInterface(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVMInterfaceFunc: func(_ context.Context, _ string, id int, _ domain.VMInterfaceWrite) (*domain.VMInterface, error) {
				return &domain.VMInterface{ID: id, Name: "eth1"}, nil
			},
		})
		vi, err := svc.UpdateVMInterface(context.Background(), 7, domain.VMInterfaceWrite{Name: "eth1"})
		if err != nil {
			t.Fatalf("UpdateVMInterface() returned error: %v", err)
		}
		if vi.ID != 7 {
			t.Errorf("vi = %+v, want id 7", vi)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"mac_address":["invalid"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVMInterfaceFunc: func(_ context.Context, _ string, _ int, _ domain.VMInterfaceWrite) (*domain.VMInterface, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdateVMInterface(context.Background(), 7, domain.VMInterfaceWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVMInterfaceFunc: func(_ context.Context, _ string, _ int, _ domain.VMInterfaceWrite) (*domain.VMInterface, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateVMInterface(context.Background(), 7, domain.VMInterfaceWrite{})
		if err == nil || !containsService(err.Error(), "update vm interface:") {
			t.Errorf("err = %v, want it to contain 'update vm interface:'", err)
		}
	})
}

func TestNetworkService_DeleteVMInterface(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVMInterfaceFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeleteVMInterface(context.Background(), 3); err != nil {
			t.Fatalf("DeleteVMInterface() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVMInterfaceFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeleteVMInterface(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVMInterfaceFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeleteVMInterface(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete vm interface:") {
			t.Errorf("err = %v, want it to contain 'delete vm interface:'", err)
		}
	})
}

func TestNetworkService_CreateProvider(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateProviderFunc: func(_ context.Context, token string, in domain.ProviderWrite) (*domain.Provider, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.Provider{ID: 1, Name: in.Name}, nil
			},
		})
		p, err := svc.CreateProvider(context.Background(), domain.ProviderWrite{Name: "ACME"})
		if err != nil {
			t.Fatalf("CreateProvider() returned error: %v", err)
		}
		if p.ID != 1 || p.Name != "ACME" {
			t.Errorf("p = %+v, want id 1 name ACME", p)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"name":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreateProviderFunc: func(_ context.Context, _ string, _ domain.ProviderWrite) (*domain.Provider, error) {
				return nil, ve
			},
		})
		_, err := svc.CreateProvider(context.Background(), domain.ProviderWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateProviderFunc: func(_ context.Context, _ string, _ domain.ProviderWrite) (*domain.Provider, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateProvider(context.Background(), domain.ProviderWrite{Name: "ACME"})
		if err == nil || !containsService(err.Error(), "create provider:") {
			t.Errorf("err = %v, want it to contain 'create provider:'", err)
		}
	})
}

func TestNetworkService_UpdateProvider(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateProviderFunc: func(_ context.Context, _ string, id int, _ domain.ProviderWrite) (*domain.Provider, error) {
				return &domain.Provider{ID: id, Name: "ACME2"}, nil
			},
		})
		p, err := svc.UpdateProvider(context.Background(), 7, domain.ProviderWrite{Name: "ACME2"})
		if err != nil {
			t.Fatalf("UpdateProvider() returned error: %v", err)
		}
		if p.ID != 7 {
			t.Errorf("p = %+v, want id 7", p)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"slug":["invalid"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdateProviderFunc: func(_ context.Context, _ string, _ int, _ domain.ProviderWrite) (*domain.Provider, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdateProvider(context.Background(), 7, domain.ProviderWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateProviderFunc: func(_ context.Context, _ string, _ int, _ domain.ProviderWrite) (*domain.Provider, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateProvider(context.Background(), 7, domain.ProviderWrite{})
		if err == nil || !containsService(err.Error(), "update provider:") {
			t.Errorf("err = %v, want it to contain 'update provider:'", err)
		}
	})
}

func TestNetworkService_DeleteProvider(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteProviderFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeleteProvider(context.Background(), 3); err != nil {
			t.Fatalf("DeleteProvider() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeleteProviderFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeleteProvider(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteProviderFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeleteProvider(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete provider:") {
			t.Errorf("err = %v, want it to contain 'delete provider:'", err)
		}
	})
}

func TestNetworkService_CreateTenant(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateTenantFunc: func(_ context.Context, token string, in domain.TenantWrite) (*domain.Tenant, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.Tenant{ID: 1, Name: in.Name}, nil
			},
		})
		p, err := svc.CreateTenant(context.Background(), domain.TenantWrite{Name: "ACME"})
		if err != nil {
			t.Fatalf("CreateTenant() returned error: %v", err)
		}
		if p.ID != 1 || p.Name != "ACME" {
			t.Errorf("p = %+v, want id 1 name ACME", p)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"name":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreateTenantFunc: func(_ context.Context, _ string, _ domain.TenantWrite) (*domain.Tenant, error) {
				return nil, ve
			},
		})
		_, err := svc.CreateTenant(context.Background(), domain.TenantWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateTenantFunc: func(_ context.Context, _ string, _ domain.TenantWrite) (*domain.Tenant, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateTenant(context.Background(), domain.TenantWrite{Name: "ACME"})
		if err == nil || !containsService(err.Error(), "create tenant:") {
			t.Errorf("err = %v, want it to contain 'create tenant:'", err)
		}
	})
}

func TestNetworkService_UpdateTenant(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateTenantFunc: func(_ context.Context, _ string, id int, _ domain.TenantWrite) (*domain.Tenant, error) {
				return &domain.Tenant{ID: id, Name: "ACME2"}, nil
			},
		})
		p, err := svc.UpdateTenant(context.Background(), 7, domain.TenantWrite{Name: "ACME2"})
		if err != nil {
			t.Fatalf("UpdateTenant() returned error: %v", err)
		}
		if p.ID != 7 {
			t.Errorf("p = %+v, want id 7", p)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"slug":["invalid"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdateTenantFunc: func(_ context.Context, _ string, _ int, _ domain.TenantWrite) (*domain.Tenant, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdateTenant(context.Background(), 7, domain.TenantWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateTenantFunc: func(_ context.Context, _ string, _ int, _ domain.TenantWrite) (*domain.Tenant, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateTenant(context.Background(), 7, domain.TenantWrite{})
		if err == nil || !containsService(err.Error(), "update tenant:") {
			t.Errorf("err = %v, want it to contain 'update tenant:'", err)
		}
	})
}

func TestNetworkService_DeleteTenant(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteTenantFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeleteTenant(context.Background(), 3); err != nil {
			t.Fatalf("DeleteTenant() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeleteTenantFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeleteTenant(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteTenantFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeleteTenant(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete tenant:") {
			t.Errorf("err = %v, want it to contain 'delete tenant:'", err)
		}
	})
}

func TestNetworkService_CreateManufacturer(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateManufacturerFunc: func(_ context.Context, token string, in domain.ManufacturerWrite) (*domain.Manufacturer, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.Manufacturer{ID: 1, Name: in.Name}, nil
			},
		})
		p, err := svc.CreateManufacturer(context.Background(), domain.ManufacturerWrite{Name: "Cisco"})
		if err != nil {
			t.Fatalf("CreateManufacturer() returned error: %v", err)
		}
		if p.ID != 1 || p.Name != "Cisco" {
			t.Errorf("p = %+v, want id 1 name Cisco", p)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"name":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreateManufacturerFunc: func(_ context.Context, _ string, _ domain.ManufacturerWrite) (*domain.Manufacturer, error) {
				return nil, ve
			},
		})
		_, err := svc.CreateManufacturer(context.Background(), domain.ManufacturerWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateManufacturerFunc: func(_ context.Context, _ string, _ domain.ManufacturerWrite) (*domain.Manufacturer, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateManufacturer(context.Background(), domain.ManufacturerWrite{Name: "Cisco"})
		if err == nil || !containsService(err.Error(), "create manufacturer:") {
			t.Errorf("err = %v, want it to contain 'create manufacturer:'", err)
		}
	})
}

func TestNetworkService_UpdateManufacturer(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateManufacturerFunc: func(_ context.Context, _ string, id int, _ domain.ManufacturerWrite) (*domain.Manufacturer, error) {
				return &domain.Manufacturer{ID: id, Name: "Cisco2"}, nil
			},
		})
		p, err := svc.UpdateManufacturer(context.Background(), 7, domain.ManufacturerWrite{Name: "Cisco2"})
		if err != nil {
			t.Fatalf("UpdateManufacturer() returned error: %v", err)
		}
		if p.ID != 7 {
			t.Errorf("p = %+v, want id 7", p)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"slug":["invalid"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdateManufacturerFunc: func(_ context.Context, _ string, _ int, _ domain.ManufacturerWrite) (*domain.Manufacturer, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdateManufacturer(context.Background(), 7, domain.ManufacturerWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateManufacturerFunc: func(_ context.Context, _ string, _ int, _ domain.ManufacturerWrite) (*domain.Manufacturer, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateManufacturer(context.Background(), 7, domain.ManufacturerWrite{})
		if err == nil || !containsService(err.Error(), "update manufacturer:") {
			t.Errorf("err = %v, want it to contain 'update manufacturer:'", err)
		}
	})
}

func TestNetworkService_DeleteManufacturer(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteManufacturerFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeleteManufacturer(context.Background(), 3); err != nil {
			t.Fatalf("DeleteManufacturer() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeleteManufacturerFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeleteManufacturer(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteManufacturerFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeleteManufacturer(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete manufacturer:") {
			t.Errorf("err = %v, want it to contain 'delete manufacturer:'", err)
		}
	})
}

func TestNetworkService_CreateDeviceType(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateDeviceTypeFunc: func(_ context.Context, token string, in domain.DeviceTypeWrite) (*domain.DeviceType, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.DeviceType{ID: 1, Model: in.Model}, nil
			},
		})
		p, err := svc.CreateDeviceType(context.Background(), domain.DeviceTypeWrite{Model: "C9300"})
		if err != nil {
			t.Fatalf("CreateDeviceType() returned error: %v", err)
		}
		if p.ID != 1 || p.Model != "C9300" {
			t.Errorf("p = %+v, want id 1 model C9300", p)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"model":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreateDeviceTypeFunc: func(_ context.Context, _ string, _ domain.DeviceTypeWrite) (*domain.DeviceType, error) {
				return nil, ve
			},
		})
		_, err := svc.CreateDeviceType(context.Background(), domain.DeviceTypeWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateDeviceTypeFunc: func(_ context.Context, _ string, _ domain.DeviceTypeWrite) (*domain.DeviceType, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateDeviceType(context.Background(), domain.DeviceTypeWrite{Model: "C9300"})
		if err == nil || !containsService(err.Error(), "create device type:") {
			t.Errorf("err = %v, want it to contain 'create device type:'", err)
		}
	})
}

func TestNetworkService_UpdateDeviceType(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateDeviceTypeFunc: func(_ context.Context, _ string, id int, _ domain.DeviceTypeWrite) (*domain.DeviceType, error) {
				return &domain.DeviceType{ID: id, Model: "C9300-2"}, nil
			},
		})
		p, err := svc.UpdateDeviceType(context.Background(), 7, domain.DeviceTypeWrite{Model: "C9300-2"})
		if err != nil {
			t.Fatalf("UpdateDeviceType() returned error: %v", err)
		}
		if p.ID != 7 {
			t.Errorf("p = %+v, want id 7", p)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"slug":["invalid"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdateDeviceTypeFunc: func(_ context.Context, _ string, _ int, _ domain.DeviceTypeWrite) (*domain.DeviceType, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdateDeviceType(context.Background(), 7, domain.DeviceTypeWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateDeviceTypeFunc: func(_ context.Context, _ string, _ int, _ domain.DeviceTypeWrite) (*domain.DeviceType, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateDeviceType(context.Background(), 7, domain.DeviceTypeWrite{})
		if err == nil || !containsService(err.Error(), "update device type:") {
			t.Errorf("err = %v, want it to contain 'update device type:'", err)
		}
	})
}

func TestNetworkService_DeleteDeviceType(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteDeviceTypeFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeleteDeviceType(context.Background(), 3); err != nil {
			t.Fatalf("DeleteDeviceType() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeleteDeviceTypeFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeleteDeviceType(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteDeviceTypeFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeleteDeviceType(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete device type:") {
			t.Errorf("err = %v, want it to contain 'delete device type:'", err)
		}
	})
}

// containsService is a tiny substring helper scoped to this test package.
func containsService(s, sub string) bool {
	return len(s) >= len(sub) && indexOfService(s, sub) >= 0
}

func indexOfService(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func TestNetworkService_CreateLocation(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateLocationFunc: func(_ context.Context, token string, in domain.LocationWrite) (*domain.Location, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.Location{ID: 1, Name: in.Name}, nil
			},
		})
		site := 5
		p, err := svc.CreateLocation(context.Background(), domain.LocationWrite{Name: "Row A", Site: &site})
		if err != nil {
			t.Fatalf("CreateLocation() returned error: %v", err)
		}
		if p.ID != 1 || p.Name != "Row A" {
			t.Errorf("p = %+v, want id 1 name 'Row A'", p)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"site":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreateLocationFunc: func(_ context.Context, _ string, _ domain.LocationWrite) (*domain.Location, error) {
				return nil, ve
			},
		})
		_, err := svc.CreateLocation(context.Background(), domain.LocationWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateLocationFunc: func(_ context.Context, _ string, _ domain.LocationWrite) (*domain.Location, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateLocation(context.Background(), domain.LocationWrite{Name: "Row A"})
		if err == nil || !containsService(err.Error(), "create location:") {
			t.Errorf("err = %v, want it to contain 'create location:'", err)
		}
	})
}

func TestNetworkService_UpdateLocation(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateLocationFunc: func(_ context.Context, _ string, id int, _ domain.LocationWrite) (*domain.Location, error) {
				return &domain.Location{ID: id, Name: "Row B"}, nil
			},
		})
		p, err := svc.UpdateLocation(context.Background(), 7, domain.LocationWrite{Name: "Row B"})
		if err != nil {
			t.Fatalf("UpdateLocation() returned error: %v", err)
		}
		if p.ID != 7 {
			t.Errorf("p = %+v, want id 7", p)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"slug":["invalid"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdateLocationFunc: func(_ context.Context, _ string, _ int, _ domain.LocationWrite) (*domain.Location, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdateLocation(context.Background(), 7, domain.LocationWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateLocationFunc: func(_ context.Context, _ string, _ int, _ domain.LocationWrite) (*domain.Location, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateLocation(context.Background(), 7, domain.LocationWrite{})
		if err == nil || !containsService(err.Error(), "update location:") {
			t.Errorf("err = %v, want it to contain 'update location:'", err)
		}
	})
}

func TestNetworkService_DeleteLocation(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteLocationFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeleteLocation(context.Background(), 3); err != nil {
			t.Fatalf("DeleteLocation() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeleteLocationFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeleteLocation(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteLocationFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeleteLocation(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete location:") {
			t.Errorf("err = %v, want it to contain 'delete location:'", err)
		}
	})
}

func TestNetworkService_CreateClusterType(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateClusterTypeFunc: func(_ context.Context, token string, in domain.ClusterTypeWrite) (*domain.ClusterType, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.ClusterType{ID: 1, Name: in.Name}, nil
			},
		})
		p, err := svc.CreateClusterType(context.Background(), domain.ClusterTypeWrite{Name: "KVM"})
		if err != nil {
			t.Fatalf("CreateClusterType() returned error: %v", err)
		}
		if p.ID != 1 || p.Name != "KVM" {
			t.Errorf("p = %+v, want id 1 name 'KVM'", p)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"name":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreateClusterTypeFunc: func(_ context.Context, _ string, _ domain.ClusterTypeWrite) (*domain.ClusterType, error) {
				return nil, ve
			},
		})
		_, err := svc.CreateClusterType(context.Background(), domain.ClusterTypeWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateClusterTypeFunc: func(_ context.Context, _ string, _ domain.ClusterTypeWrite) (*domain.ClusterType, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateClusterType(context.Background(), domain.ClusterTypeWrite{Name: "KVM"})
		if err == nil || !containsService(err.Error(), "create cluster type:") {
			t.Errorf("err = %v, want it to contain 'create cluster type:'", err)
		}
	})
}

func TestNetworkService_UpdateClusterType(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateClusterTypeFunc: func(_ context.Context, _ string, id int, _ domain.ClusterTypeWrite) (*domain.ClusterType, error) {
				return &domain.ClusterType{ID: id, Name: "KVM"}, nil
			},
		})
		p, err := svc.UpdateClusterType(context.Background(), 7, domain.ClusterTypeWrite{Name: "KVM"})
		if err != nil {
			t.Fatalf("UpdateClusterType() returned error: %v", err)
		}
		if p.ID != 7 {
			t.Errorf("p = %+v, want id 7", p)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"slug":["invalid"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdateClusterTypeFunc: func(_ context.Context, _ string, _ int, _ domain.ClusterTypeWrite) (*domain.ClusterType, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdateClusterType(context.Background(), 7, domain.ClusterTypeWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateClusterTypeFunc: func(_ context.Context, _ string, _ int, _ domain.ClusterTypeWrite) (*domain.ClusterType, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateClusterType(context.Background(), 7, domain.ClusterTypeWrite{})
		if err == nil || !containsService(err.Error(), "update cluster type:") {
			t.Errorf("err = %v, want it to contain 'update cluster type:'", err)
		}
	})
}

func TestNetworkService_DeleteClusterType(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteClusterTypeFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeleteClusterType(context.Background(), 3); err != nil {
			t.Fatalf("DeleteClusterType() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeleteClusterTypeFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeleteClusterType(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteClusterTypeFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeleteClusterType(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete cluster type:") {
			t.Errorf("err = %v, want it to contain 'delete cluster type:'", err)
		}
	})
}

func TestNetworkService_CreateClusterGroup(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateClusterGroupFunc: func(_ context.Context, token string, in domain.ClusterGroupWrite) (*domain.ClusterGroup, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.ClusterGroup{ID: 1, Name: in.Name}, nil
			},
		})
		p, err := svc.CreateClusterGroup(context.Background(), domain.ClusterGroupWrite{Name: "DC Clusters"})
		if err != nil {
			t.Fatalf("CreateClusterGroup() returned error: %v", err)
		}
		if p.ID != 1 || p.Name != "DC Clusters" {
			t.Errorf("p = %+v, want id 1 name 'DC Clusters'", p)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"name":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreateClusterGroupFunc: func(_ context.Context, _ string, _ domain.ClusterGroupWrite) (*domain.ClusterGroup, error) {
				return nil, ve
			},
		})
		_, err := svc.CreateClusterGroup(context.Background(), domain.ClusterGroupWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateClusterGroupFunc: func(_ context.Context, _ string, _ domain.ClusterGroupWrite) (*domain.ClusterGroup, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateClusterGroup(context.Background(), domain.ClusterGroupWrite{Name: "DC Clusters"})
		if err == nil || !containsService(err.Error(), "create cluster group:") {
			t.Errorf("err = %v, want it to contain 'create cluster group:'", err)
		}
	})
}

func TestNetworkService_UpdateClusterGroup(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateClusterGroupFunc: func(_ context.Context, _ string, id int, _ domain.ClusterGroupWrite) (*domain.ClusterGroup, error) {
				return &domain.ClusterGroup{ID: id, Name: "DC Clusters"}, nil
			},
		})
		p, err := svc.UpdateClusterGroup(context.Background(), 7, domain.ClusterGroupWrite{Name: "DC Clusters"})
		if err != nil {
			t.Fatalf("UpdateClusterGroup() returned error: %v", err)
		}
		if p.ID != 7 {
			t.Errorf("p = %+v, want id 7", p)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"slug":["invalid"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdateClusterGroupFunc: func(_ context.Context, _ string, _ int, _ domain.ClusterGroupWrite) (*domain.ClusterGroup, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdateClusterGroup(context.Background(), 7, domain.ClusterGroupWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateClusterGroupFunc: func(_ context.Context, _ string, _ int, _ domain.ClusterGroupWrite) (*domain.ClusterGroup, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateClusterGroup(context.Background(), 7, domain.ClusterGroupWrite{})
		if err == nil || !containsService(err.Error(), "update cluster group:") {
			t.Errorf("err = %v, want it to contain 'update cluster group:'", err)
		}
	})
}

func TestNetworkService_DeleteClusterGroup(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteClusterGroupFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeleteClusterGroup(context.Background(), 3); err != nil {
			t.Fatalf("DeleteClusterGroup() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeleteClusterGroupFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeleteClusterGroup(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteClusterGroupFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeleteClusterGroup(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete cluster group:") {
			t.Errorf("err = %v, want it to contain 'delete cluster group:'", err)
		}
	})
}

func TestNetworkService_CreateCircuitType(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateCircuitTypeFunc: func(_ context.Context, token string, in domain.CircuitTypeWrite) (*domain.CircuitType, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.CircuitType{ID: 1, Name: in.Name}, nil
			},
		})
		p, err := svc.CreateCircuitType(context.Background(), domain.CircuitTypeWrite{Name: "Fiber"})
		if err != nil {
			t.Fatalf("CreateCircuitType() returned error: %v", err)
		}
		if p.ID != 1 || p.Name != "Fiber" {
			t.Errorf("p = %+v, want id 1 name 'Fiber'", p)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"name":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreateCircuitTypeFunc: func(_ context.Context, _ string, _ domain.CircuitTypeWrite) (*domain.CircuitType, error) {
				return nil, ve
			},
		})
		_, err := svc.CreateCircuitType(context.Background(), domain.CircuitTypeWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateCircuitTypeFunc: func(_ context.Context, _ string, _ domain.CircuitTypeWrite) (*domain.CircuitType, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateCircuitType(context.Background(), domain.CircuitTypeWrite{Name: "Fiber"})
		if err == nil || !containsService(err.Error(), "create circuit type:") {
			t.Errorf("err = %v, want it to contain 'create circuit type:'", err)
		}
	})
}

func TestNetworkService_UpdateCircuitType(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateCircuitTypeFunc: func(_ context.Context, _ string, id int, _ domain.CircuitTypeWrite) (*domain.CircuitType, error) {
				return &domain.CircuitType{ID: id, Name: "Fiber"}, nil
			},
		})
		p, err := svc.UpdateCircuitType(context.Background(), 7, domain.CircuitTypeWrite{Name: "Fiber"})
		if err != nil {
			t.Fatalf("UpdateCircuitType() returned error: %v", err)
		}
		if p.ID != 7 {
			t.Errorf("p = %+v, want id 7", p)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"slug":["invalid"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdateCircuitTypeFunc: func(_ context.Context, _ string, _ int, _ domain.CircuitTypeWrite) (*domain.CircuitType, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdateCircuitType(context.Background(), 7, domain.CircuitTypeWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateCircuitTypeFunc: func(_ context.Context, _ string, _ int, _ domain.CircuitTypeWrite) (*domain.CircuitType, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateCircuitType(context.Background(), 7, domain.CircuitTypeWrite{})
		if err == nil || !containsService(err.Error(), "update circuit type:") {
			t.Errorf("err = %v, want it to contain 'update circuit type:'", err)
		}
	})
}

func TestNetworkService_DeleteCircuitType(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteCircuitTypeFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeleteCircuitType(context.Background(), 3); err != nil {
			t.Fatalf("DeleteCircuitType() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeleteCircuitTypeFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeleteCircuitType(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteCircuitTypeFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeleteCircuitType(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete circuit type:") {
			t.Errorf("err = %v, want it to contain 'delete circuit type:'", err)
		}
	})
}

func TestNetworkService_CreateVrf(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVrfFunc: func(_ context.Context, token string, in domain.VrfWrite) (*domain.Vrf, error) {
				if token != "test-token" {
					t.Errorf("token = %q, want test-token", token)
				}
				return &domain.Vrf{ID: 1, Name: in.Name, Rd: in.Rd}, nil
			},
		})
		p, err := svc.CreateVrf(context.Background(), domain.VrfWrite{Name: "prod", Rd: "65000:1"})
		if err != nil {
			t.Fatalf("CreateVrf() returned error: %v", err)
		}
		if p.ID != 1 || p.Name != "prod" || p.Rd != "65000:1" {
			t.Errorf("p = %+v, want id 1 name 'prod' rd '65000:1'", p)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"name":["required"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			CreateVrfFunc: func(_ context.Context, _ string, _ domain.VrfWrite) (*domain.Vrf, error) {
				return nil, ve
			},
		})
		_, err := svc.CreateVrf(context.Background(), domain.VrfWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var got *domain.ValidationError
		if !errors.As(err, &got) {
			t.Fatalf("err = %v, want errors.As to find *domain.ValidationError", err)
		}
		if got != ve {
			t.Error("ValidationError was not passed through unchanged")
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVrfFunc: func(_ context.Context, _ string, _ domain.VrfWrite) (*domain.Vrf, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateVrf(context.Background(), domain.VrfWrite{Name: "prod", Rd: "65000:1"})
		if err == nil || !containsService(err.Error(), "create vrf:") {
			t.Errorf("err = %v, want it to contain 'create vrf:'", err)
		}
	})
}

func TestNetworkService_UpdateVrf(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVrfFunc: func(_ context.Context, _ string, id int, _ domain.VrfWrite) (*domain.Vrf, error) {
				return &domain.Vrf{ID: id, Name: "prod", Rd: "65000:1"}, nil
			},
		})
		p, err := svc.UpdateVrf(context.Background(), 7, domain.VrfWrite{Name: "prod"})
		if err != nil {
			t.Fatalf("UpdateVrf() returned error: %v", err)
		}
		if p.ID != 7 {
			t.Errorf("p = %+v, want id 7", p)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"rd":["invalid"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVrfFunc: func(_ context.Context, _ string, _ int, _ domain.VrfWrite) (*domain.Vrf, error) {
				return nil, ve
			},
		})
		_, err := svc.UpdateVrf(context.Background(), 7, domain.VrfWrite{})
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVrfFunc: func(_ context.Context, _ string, _ int, _ domain.VrfWrite) (*domain.Vrf, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateVrf(context.Background(), 7, domain.VrfWrite{})
		if err == nil || !containsService(err.Error(), "update vrf:") {
			t.Errorf("err = %v, want it to contain 'update vrf:'", err)
		}
	})
}

func TestNetworkService_DeleteVrf(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVrfFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		})
		if err := svc.DeleteVrf(context.Background(), 3); err != nil {
			t.Fatalf("DeleteVrf() returned error: %v", err)
		}
	})

	t.Run("propagates validation error via errors.As", func(t *testing.T) {
		ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"detail":["dependent"]}`)}
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVrfFunc: func(_ context.Context, _ string, _ int) error {
				return ve
			},
		})
		err := svc.DeleteVrf(context.Background(), 3)
		var got *domain.ValidationError
		if !errors.As(err, &got) || got != ve {
			t.Errorf("err = %v, want the ValidationError passed through", err)
		}
	})

	t.Run("wraps generic error with label", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVrfFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		err := svc.DeleteVrf(context.Background(), 3)
		if err == nil || !containsService(err.Error(), "delete vrf:") {
			t.Errorf("err = %v, want it to contain 'delete vrf:'", err)
		}
	})
}
