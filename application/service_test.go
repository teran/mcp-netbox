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
