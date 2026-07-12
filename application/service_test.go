package application

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/teran/mcp-netbox/domain"
)

type mockRepo struct {
	listSitesFunc           func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Site], error)
	listDevicesFunc         func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Device], error)
	listIPAddressesFunc     func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.IPAddress], error)
	listPrefixesFunc        func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Prefix], error)
	listVLANsFunc           func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VLAN], error)
	listVirtualMachinesFunc func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VirtualMachine], error)
	listClustersFunc        func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Cluster], error)
	listCircuitsFunc        func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Circuit], error)
	listRacksFunc           func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Rack], error)
	getObjectFunc           func(ctx context.Context, token string, objectType string, id int, params map[string]string) (domain.RawObject, error)
}

func (m *mockRepo) ListSites(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Site], error) {
	return m.listSitesFunc(ctx, token, params)
}
func (m *mockRepo) ListDevices(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Device], error) {
	return m.listDevicesFunc(ctx, token, params)
}
func (m *mockRepo) ListIPAddresses(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.IPAddress], error) {
	return m.listIPAddressesFunc(ctx, token, params)
}
func (m *mockRepo) ListPrefixes(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Prefix], error) {
	return m.listPrefixesFunc(ctx, token, params)
}
func (m *mockRepo) ListVLANs(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VLAN], error) {
	return m.listVLANsFunc(ctx, token, params)
}
func (m *mockRepo) ListVirtualMachines(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VirtualMachine], error) {
	return m.listVirtualMachinesFunc(ctx, token, params)
}
func (m *mockRepo) ListClusters(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Cluster], error) {
	return m.listClustersFunc(ctx, token, params)
}
func (m *mockRepo) ListCircuits(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Circuit], error) {
	return m.listCircuitsFunc(ctx, token, params)
}
func (m *mockRepo) ListRacks(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Rack], error) {
	return m.listRacksFunc(ctx, token, params)
}
func (m *mockRepo) GetObject(ctx context.Context, token string, objectType string, id int, params map[string]string) (domain.RawObject, error) {
	return m.getObjectFunc(ctx, token, objectType, id, params)
}

func newTestService(repo *mockRepo) *NetworkService {
	return NewNetworkService(repo, "test-token")
}

func TestNetworkService_ListSites(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockRepo{
			listSitesFunc: func(_ context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Site], error) {
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
		svc := newTestService(&mockRepo{
			listSitesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Site], error) {
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
		svc := newTestService(&mockRepo{
			listDevicesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Device], error) {
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
		svc := newTestService(&mockRepo{
			listDevicesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Device], error) {
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
		svc := newTestService(&mockRepo{
			listIPAddressesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.IPAddress], error) {
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
		svc := newTestService(&mockRepo{
			listIPAddressesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.IPAddress], error) {
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
		svc := newTestService(&mockRepo{
			listPrefixesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Prefix], error) {
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
		svc := newTestService(&mockRepo{
			listPrefixesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Prefix], error) {
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
		svc := newTestService(&mockRepo{
			listVLANsFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.VLAN], error) {
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
		svc := newTestService(&mockRepo{
			listVLANsFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.VLAN], error) {
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
		svc := newTestService(&mockRepo{
			listVirtualMachinesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.VirtualMachine], error) {
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
		svc := newTestService(&mockRepo{
			listVirtualMachinesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.VirtualMachine], error) {
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
		svc := newTestService(&mockRepo{
			listClustersFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Cluster], error) {
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
		svc := newTestService(&mockRepo{
			listClustersFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Cluster], error) {
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
		svc := newTestService(&mockRepo{
			listCircuitsFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Circuit], error) {
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
		svc := newTestService(&mockRepo{
			listCircuitsFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Circuit], error) {
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
		svc := newTestService(&mockRepo{
			listRacksFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Rack], error) {
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
		svc := newTestService(&mockRepo{
			listRacksFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Rack], error) {
				return nil, errors.New("repo error")
			},
		})
		_, err := svc.ListRacks(context.Background(), nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestNetworkService_GetObject(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockRepo{
			getObjectFunc: func(_ context.Context, _ string, objectType string, id int, _ map[string]string) (domain.RawObject, error) {
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
		svc := newTestService(&mockRepo{
			getObjectFunc: func(_ context.Context, _ string, _ string, _ int, _ map[string]string) (domain.RawObject, error) {
				return nil, errors.New("repo error")
			},
		})
		_, err := svc.GetObject(context.Background(), "site", 1)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}
