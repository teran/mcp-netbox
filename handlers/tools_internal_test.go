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
