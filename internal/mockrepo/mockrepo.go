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
	ListSitesFunc               func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Site], error)
	ListDevicesFunc             func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Device], error)
	ListIPAddressesFunc         func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.IPAddress], error)
	ListPrefixesFunc            func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Prefix], error)
	ListVLANsFunc               func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VLAN], error)
	ListVirtualMachinesFunc     func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VirtualMachine], error)
	ListClustersFunc            func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Cluster], error)
	ListCircuitsFunc            func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Circuit], error)
	ListCircuitTerminationsFunc func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.CircuitTermination], error)
	ListCablesFunc              func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Cable], error)
	ListRacksFunc               func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Rack], error)
	ListInterfacesFunc          func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Interface], error)
	ListVMInterfacesFunc        func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VMInterface], error)
	GetObjectFunc               func(ctx context.Context, token string, objectType string, id int, params map[string]string) (domain.RawObject, error)
	CreateSiteFunc              func(ctx context.Context, token string, in domain.SiteWrite) (*domain.Site, error)
	UpdateSiteFunc              func(ctx context.Context, token string, id int, in domain.SiteWrite) (*domain.Site, error)
	DeleteSiteFunc              func(ctx context.Context, token string, id int) error
	CreateDeviceFunc            func(ctx context.Context, token string, in domain.DeviceWrite) (*domain.Device, error)
	UpdateDeviceFunc            func(ctx context.Context, token string, id int, in domain.DeviceWrite) (*domain.Device, error)
	DeleteDeviceFunc            func(ctx context.Context, token string, id int) error
	CreateIPAddressFunc         func(ctx context.Context, token string, in domain.IPAddressWrite) (*domain.IPAddress, error)
	UpdateIPAddressFunc         func(ctx context.Context, token string, id int, in domain.IPAddressWrite) (*domain.IPAddress, error)
	DeleteIPAddressFunc         func(ctx context.Context, token string, id int) error
	CreatePrefixFunc            func(ctx context.Context, token string, in domain.PrefixWrite) (*domain.Prefix, error)
	UpdatePrefixFunc            func(ctx context.Context, token string, id int, in domain.PrefixWrite) (*domain.Prefix, error)
	DeletePrefixFunc            func(ctx context.Context, token string, id int) error
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

func (m *MockRepo) ListCircuitTerminations(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.CircuitTermination], error) {
	return m.ListCircuitTerminationsFunc(ctx, token, params)
}

func (m *MockRepo) ListCables(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Cable], error) {
	return m.ListCablesFunc(ctx, token, params)
}

func (m *MockRepo) ListRacks(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Rack], error) {
	return m.ListRacksFunc(ctx, token, params)
}

func (m *MockRepo) ListInterfaces(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Interface], error) {
	return m.ListInterfacesFunc(ctx, token, params)
}

func (m *MockRepo) ListVMInterfaces(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VMInterface], error) {
	return m.ListVMInterfacesFunc(ctx, token, params)
}

func (m *MockRepo) GetObject(ctx context.Context, token string, objectType string, id int, params map[string]string) (domain.RawObject, error) {
	return m.GetObjectFunc(ctx, token, objectType, id, params)
}

func (m *MockRepo) CreateSite(ctx context.Context, token string, in domain.SiteWrite) (*domain.Site, error) {
	return m.CreateSiteFunc(ctx, token, in)
}

func (m *MockRepo) UpdateSite(ctx context.Context, token string, id int, in domain.SiteWrite) (*domain.Site, error) {
	return m.UpdateSiteFunc(ctx, token, id, in)
}

func (m *MockRepo) DeleteSite(ctx context.Context, token string, id int) error {
	return m.DeleteSiteFunc(ctx, token, id)
}

func (m *MockRepo) CreateDevice(ctx context.Context, token string, in domain.DeviceWrite) (*domain.Device, error) {
	return m.CreateDeviceFunc(ctx, token, in)
}

func (m *MockRepo) UpdateDevice(ctx context.Context, token string, id int, in domain.DeviceWrite) (*domain.Device, error) {
	return m.UpdateDeviceFunc(ctx, token, id, in)
}

func (m *MockRepo) DeleteDevice(ctx context.Context, token string, id int) error {
	return m.DeleteDeviceFunc(ctx, token, id)
}

func (m *MockRepo) CreateIPAddress(ctx context.Context, token string, in domain.IPAddressWrite) (*domain.IPAddress, error) {
	return m.CreateIPAddressFunc(ctx, token, in)
}

func (m *MockRepo) UpdateIPAddress(ctx context.Context, token string, id int, in domain.IPAddressWrite) (*domain.IPAddress, error) {
	return m.UpdateIPAddressFunc(ctx, token, id, in)
}

func (m *MockRepo) DeleteIPAddress(ctx context.Context, token string, id int) error {
	return m.DeleteIPAddressFunc(ctx, token, id)
}

func (m *MockRepo) CreatePrefix(ctx context.Context, token string, in domain.PrefixWrite) (*domain.Prefix, error) {
	return m.CreatePrefixFunc(ctx, token, in)
}

func (m *MockRepo) UpdatePrefix(ctx context.Context, token string, id int, in domain.PrefixWrite) (*domain.Prefix, error) {
	return m.UpdatePrefixFunc(ctx, token, id, in)
}

func (m *MockRepo) DeletePrefix(ctx context.Context, token string, id int) error {
	return m.DeletePrefixFunc(ctx, token, id)
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

func (r *NilRepo) ListCircuitTerminations(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.CircuitTermination], error) {
	return &domain.PaginatedResponse[domain.CircuitTermination]{Count: 0, Results: []domain.CircuitTermination{}}, nil
}

func (r *NilRepo) ListCables(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Cable], error) {
	return &domain.PaginatedResponse[domain.Cable]{Count: 0, Results: []domain.Cable{}}, nil
}

func (r *NilRepo) ListRacks(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Rack], error) {
	return &domain.PaginatedResponse[domain.Rack]{Count: 0, Results: []domain.Rack{}}, nil
}

func (r *NilRepo) ListInterfaces(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Interface], error) {
	return &domain.PaginatedResponse[domain.Interface]{Count: 0, Results: []domain.Interface{}}, nil
}

func (r *NilRepo) ListVMInterfaces(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VMInterface], error) {
	return &domain.PaginatedResponse[domain.VMInterface]{Count: 0, Results: []domain.VMInterface{}}, nil
}

func (r *NilRepo) GetObject(ctx context.Context, token string, objectType string, id int, params map[string]string) (domain.RawObject, error) {
	return domain.RawObject(fmt.Sprintf(`{"id":%d}`, id)), nil
}

func (r *NilRepo) CreateSite(ctx context.Context, token string, in domain.SiteWrite) (*domain.Site, error) {
	return &domain.Site{ID: 1, Name: in.Name}, nil
}

func (r *NilRepo) UpdateSite(ctx context.Context, token string, id int, in domain.SiteWrite) (*domain.Site, error) {
	return &domain.Site{ID: id, Name: in.Name}, nil
}

func (r *NilRepo) DeleteSite(ctx context.Context, token string, id int) error {
	return nil
}

func (r *NilRepo) CreateDevice(ctx context.Context, token string, in domain.DeviceWrite) (*domain.Device, error) {
	return &domain.Device{ID: 1, Name: in.Name}, nil
}

func (r *NilRepo) UpdateDevice(ctx context.Context, token string, id int, in domain.DeviceWrite) (*domain.Device, error) {
	return &domain.Device{ID: id, Name: in.Name}, nil
}

func (r *NilRepo) DeleteDevice(ctx context.Context, token string, id int) error {
	return nil
}

func (r *NilRepo) CreateIPAddress(ctx context.Context, token string, in domain.IPAddressWrite) (*domain.IPAddress, error) {
	return &domain.IPAddress{ID: 1, Address: in.Address}, nil
}

func (r *NilRepo) UpdateIPAddress(ctx context.Context, token string, id int, in domain.IPAddressWrite) (*domain.IPAddress, error) {
	return &domain.IPAddress{ID: id, Address: in.Address}, nil
}

func (r *NilRepo) DeleteIPAddress(ctx context.Context, token string, id int) error {
	return nil
}

func (r *NilRepo) CreatePrefix(ctx context.Context, token string, in domain.PrefixWrite) (*domain.Prefix, error) {
	return &domain.Prefix{ID: 1, Prefix: in.Prefix}, nil
}

func (r *NilRepo) UpdatePrefix(ctx context.Context, token string, id int, in domain.PrefixWrite) (*domain.Prefix, error) {
	return &domain.Prefix{ID: id, Prefix: in.Prefix}, nil
}

func (r *NilRepo) DeletePrefix(ctx context.Context, token string, id int) error {
	return nil
}
