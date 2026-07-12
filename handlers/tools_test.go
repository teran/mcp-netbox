package handlers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/domain"
	"github.com/teran/mcp-netbox/handlers"
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

func TestGetSitesHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockRepo{
			listSitesFunc: func(_ context.Context, _ string, params map[string]string) (*domain.PaginatedResponse[domain.Site], error) {
				return &domain.PaginatedResponse[domain.Site]{
					Count:   1,
					Results: []domain.Site{{ID: 1, Name: "Test Site"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetSitesHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.SitesInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 {
			t.Fatalf("len(Results) = %d, want %d", len(out.Results), 1)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := application.NewNetworkService(&mockRepo{
			listSitesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Site], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetSitesHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.SitesInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestGetDevicesHandler(t *testing.T) {
	t.Parallel()

	svc := application.NewNetworkService(&mockRepo{
		listDevicesFunc: func(_ context.Context, _ string, params map[string]string) (*domain.PaginatedResponse[domain.Device], error) {
			return &domain.PaginatedResponse[domain.Device]{
				Count: 1,
				Results: []domain.Device{{
					ID:   1,
					Name: "router1",
					Site: &domain.Nested{ID: 1, Name: "DC1"},
				}},
			}, nil
		},
	}, "token")

	handler := handlers.NewGetDevicesHandler(svc)
	result, out, err := handler(context.Background(), nil, handlers.DevicesInput{Site: "dc1"})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if result.IsError {
		t.Error("result.IsError = true, want false")
	}
	if out.Count != 1 {
		t.Errorf("Count = %d, want %d", out.Count, 1)
	}
}

func TestGetDevicesHandler_Error(t *testing.T) {
	t.Parallel()
	svc := application.NewNetworkService(&mockRepo{
		listDevicesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Device], error) {
			return nil, errors.New("repo error")
		},
	}, "token")

	handler := handlers.NewGetDevicesHandler(svc)
	result, _, _ := handler(context.Background(), nil, handlers.DevicesInput{})
	if !result.IsError {
		t.Error("result.IsError = false, want true")
	}
}

func TestGetIPAddressesHandler(t *testing.T) {
	t.Parallel()

	svc := application.NewNetworkService(&mockRepo{
		listIPAddressesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.IPAddress], error) {
			return &domain.PaginatedResponse[domain.IPAddress]{Count: 0}, nil
		},
	}, "token")

	handler := handlers.NewGetIPAddressesHandler(svc)
	result, out, err := handler(context.Background(), nil, handlers.IPAddressesInput{})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if result.IsError {
		t.Error("result.IsError = true, want false")
	}
	if out.Count != 0 {
		t.Errorf("Count = %d, want %d", out.Count, 0)
	}
}

func TestGetPrefixesHandler(t *testing.T) {
	t.Parallel()

	svc := application.NewNetworkService(&mockRepo{
		listPrefixesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Prefix], error) {
			return &domain.PaginatedResponse[domain.Prefix]{Count: 0}, nil
		},
	}, "token")

	handler := handlers.NewGetPrefixesHandler(svc)
	result, out, err := handler(context.Background(), nil, handlers.PrefixesInput{})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if result.IsError {
		t.Error("result.IsError = true, want false")
	}
	if out.Count != 0 {
		t.Errorf("Count = %d, want %d", out.Count, 0)
	}
}

func TestGetVLANsHandler(t *testing.T) {
	t.Parallel()

	svc := application.NewNetworkService(&mockRepo{
		listVLANsFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.VLAN], error) {
			return &domain.PaginatedResponse[domain.VLAN]{Count: 0}, nil
		},
	}, "token")

	handler := handlers.NewGetVLANsHandler(svc)
	result, out, err := handler(context.Background(), nil, handlers.VLANsInput{})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if result.IsError {
		t.Error("result.IsError = true, want false")
	}
	if out.Count != 0 {
		t.Errorf("Count = %d, want %d", out.Count, 0)
	}
}

func TestGetVirtualMachinesHandler(t *testing.T) {
	t.Parallel()

	svc := application.NewNetworkService(&mockRepo{
		listVirtualMachinesFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.VirtualMachine], error) {
			return &domain.PaginatedResponse[domain.VirtualMachine]{Count: 0}, nil
		},
	}, "token")

	handler := handlers.NewGetVirtualMachinesHandler(svc)
	result, out, err := handler(context.Background(), nil, handlers.VirtualMachinesInput{})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if result.IsError {
		t.Error("result.IsError = true, want false")
	}
	if out.Count != 0 {
		t.Errorf("Count = %d, want %d", out.Count, 0)
	}
}

func TestGetClustersHandler(t *testing.T) {
	t.Parallel()

	svc := application.NewNetworkService(&mockRepo{
		listClustersFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Cluster], error) {
			return &domain.PaginatedResponse[domain.Cluster]{Count: 0}, nil
		},
	}, "token")

	handler := handlers.NewGetClustersHandler(svc)
	result, out, err := handler(context.Background(), nil, handlers.ClustersInput{})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if result.IsError {
		t.Error("result.IsError = true, want false")
	}
	if out.Count != 0 {
		t.Errorf("Count = %d, want %d", out.Count, 0)
	}
}

func TestGetCircuitsHandler(t *testing.T) {
	t.Parallel()

	svc := application.NewNetworkService(&mockRepo{
		listCircuitsFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Circuit], error) {
			return &domain.PaginatedResponse[domain.Circuit]{Count: 0}, nil
		},
	}, "token")

	handler := handlers.NewGetCircuitsHandler(svc)
	result, out, err := handler(context.Background(), nil, handlers.CircuitsInput{})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if result.IsError {
		t.Error("result.IsError = true, want false")
	}
	if out.Count != 0 {
		t.Errorf("Count = %d, want %d", out.Count, 0)
	}
}

func TestGetRacksHandler(t *testing.T) {
	t.Parallel()

	svc := application.NewNetworkService(&mockRepo{
		listRacksFunc: func(_ context.Context, _ string, _ map[string]string) (*domain.PaginatedResponse[domain.Rack], error) {
			return &domain.PaginatedResponse[domain.Rack]{Count: 0}, nil
		},
	}, "token")

	handler := handlers.NewGetRacksHandler(svc)
	result, out, err := handler(context.Background(), nil, handlers.RacksInput{})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if result.IsError {
		t.Error("result.IsError = true, want false")
	}
	if out.Count != 0 {
		t.Errorf("Count = %d, want %d", out.Count, 0)
	}
}

func TestGetObjectByIDHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockRepo{
			getObjectFunc: func(_ context.Context, _ string, objectType string, id int, _ map[string]string) (domain.RawObject, error) {
				return domain.RawObject(`{"id":1,"name":"Test"}`), nil
			},
		}, "token")

		handler := handlers.NewGetObjectByIDHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.GetObjectInput{ObjectType: "site", ID: 1})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data == nil {
			t.Fatal("Data is nil")
		}
	})

	t.Run("missing object_type", func(t *testing.T) {
		svc := application.NewNetworkService(&mockRepo{}, "token")

		handler := handlers.NewGetObjectByIDHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.GetObjectInput{ObjectType: "", ID: 1})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockRepo{
			getObjectFunc: func(_ context.Context, _ string, _ string, _ int, _ map[string]string) (domain.RawObject, error) {
				return nil, errors.New("not found")
			},
		}, "token")
		handler := handlers.NewGetObjectByIDHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.GetObjectInput{ObjectType: "site", ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})

	t.Run("error from service", func(t *testing.T) {
		svc := application.NewNetworkService(&mockRepo{
			getObjectFunc: func(_ context.Context, _ string, _ string, _ int, _ map[string]string) (domain.RawObject, error) {
				return nil, errors.New("not found")
			},
		}, "token")

		handler := handlers.NewGetObjectByIDHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.GetObjectInput{ObjectType: "site", ID: 999})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestPaginationParams(t *testing.T) {
	t.Parallel()

	// Test from internal_test.go would go here
	// For now, the handlers package test aligns with tool usage
}
