// Package mockrepo provides test helpers for domain.NetworkRepository.
package mockrepo

import (
	"context"
	"fmt"

	"github.com/teran/mcp-netbox/domain"
)

// MockRepo is a configurable mock implementation of domain.NetworkRepository.
// Each method can be set independently via function fields.
type MockRepo struct {
	ListSitesFunc           func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Site], error)
	ListDevicesFunc         func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Device], error)
	ListIPAddressesFunc     func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.IPAddress], error)
	ListPrefixesFunc        func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Prefix], error)
	ListVLANsFunc           func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VLAN], error)
	ListVirtualMachinesFunc func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VirtualMachine], error)
	ListClustersFunc        func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Cluster], error)
	ListCircuitsFunc        func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Circuit], error)
	ListRacksFunc           func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Rack], error)
	GetObjectFunc           func(ctx context.Context, token string, objectType string, id int, params map[string]string) (domain.RawObject, error)
}

func (m *MockRepo) ListSites(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Site], error) {
	return m.ListSitesFunc(ctx, token, params)
}

func (m *MockRepo) ListDevices(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Device], error) {
	return m.ListDevicesFunc(ctx, token, params)
}

func (m *MockRepo) ListIPAddresses(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.IPAddress], error) {
	return m.ListIPAddressesFunc(ctx, token, params)
}

func (m *MockRepo) ListPrefixes(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Prefix], error) {
	return m.ListPrefixesFunc(ctx, token, params)
}

func (m *MockRepo) ListVLANs(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VLAN], error) {
	return m.ListVLANsFunc(ctx, token, params)
}

func (m *MockRepo) ListVirtualMachines(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VirtualMachine], error) {
	return m.ListVirtualMachinesFunc(ctx, token, params)
}

func (m *MockRepo) ListClusters(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Cluster], error) {
	return m.ListClustersFunc(ctx, token, params)
}

func (m *MockRepo) ListCircuits(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Circuit], error) {
	return m.ListCircuitsFunc(ctx, token, params)
}

func (m *MockRepo) ListRacks(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Rack], error) {
	return m.ListRacksFunc(ctx, token, params)
}

func (m *MockRepo) GetObject(ctx context.Context, token string, objectType string, id int, params map[string]string) (domain.RawObject, error) {
	return m.GetObjectFunc(ctx, token, objectType, id, params)
}

// NilRepo is a mock implementation of domain.NetworkRepository
// that returns empty results for every method.
type NilRepo struct{}

func (r *NilRepo) ListSites(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Site], error) {
	return &domain.PaginatedResponse[domain.Site]{Count: 0, Results: []domain.Site{}}, nil
}

func (r *NilRepo) ListDevices(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Device], error) {
	return &domain.PaginatedResponse[domain.Device]{Count: 0, Results: []domain.Device{}}, nil
}

func (r *NilRepo) ListIPAddresses(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.IPAddress], error) {
	return &domain.PaginatedResponse[domain.IPAddress]{Count: 0, Results: []domain.IPAddress{}}, nil
}

func (r *NilRepo) ListPrefixes(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Prefix], error) {
	return &domain.PaginatedResponse[domain.Prefix]{Count: 0, Results: []domain.Prefix{}}, nil
}

func (r *NilRepo) ListVLANs(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VLAN], error) {
	return &domain.PaginatedResponse[domain.VLAN]{Count: 0, Results: []domain.VLAN{}}, nil
}

func (r *NilRepo) ListVirtualMachines(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VirtualMachine], error) {
	return &domain.PaginatedResponse[domain.VirtualMachine]{Count: 0, Results: []domain.VirtualMachine{}}, nil
}

func (r *NilRepo) ListClusters(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Cluster], error) {
	return &domain.PaginatedResponse[domain.Cluster]{Count: 0, Results: []domain.Cluster{}}, nil
}

func (r *NilRepo) ListCircuits(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Circuit], error) {
	return &domain.PaginatedResponse[domain.Circuit]{Count: 0, Results: []domain.Circuit{}}, nil
}

func (r *NilRepo) ListRacks(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Rack], error) {
	return &domain.PaginatedResponse[domain.Rack]{Count: 0, Results: []domain.Rack{}}, nil
}

func (r *NilRepo) GetObject(ctx context.Context, token string, objectType string, id int, params map[string]string) (domain.RawObject, error) {
	return domain.RawObject(fmt.Sprintf(`{"id":%d}`, id)), nil
}
