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
		resp, err := svc.GetObject(context.Background(), "site", 1)
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
		_, err := svc.GetObject(context.Background(), "site", 1)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}
