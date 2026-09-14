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
	ListSitesFunc                func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Site], error)
	ListDevicesFunc              func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Device], error)
	ListIPAddressesFunc          func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.IPAddress], error)
	ListPrefixesFunc             func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Prefix], error)
	ListVLANsFunc                func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VLAN], error)
	ListVirtualMachinesFunc      func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VirtualMachine], error)
	ListClustersFunc             func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Cluster], error)
	ListCircuitsFunc             func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Circuit], error)
	ListCircuitTerminationsFunc  func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.CircuitTermination], error)
	ListCablesFunc               func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Cable], error)
	ListRacksFunc                func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Rack], error)
	ListInterfacesFunc           func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Interface], error)
	ListVMInterfacesFunc         func(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VMInterface], error)
	GetObjectFunc                func(ctx context.Context, token string, objectType string, id int, params map[string]string) (domain.RawObject, error)
	CreateSiteFunc               func(ctx context.Context, token string, in domain.SiteWrite) (*domain.Site, error)
	UpdateSiteFunc               func(ctx context.Context, token string, id int, in domain.SiteWrite) (*domain.Site, error)
	DeleteSiteFunc               func(ctx context.Context, token string, id int) error
	CreateDeviceFunc             func(ctx context.Context, token string, in domain.DeviceWrite) (*domain.Device, error)
	UpdateDeviceFunc             func(ctx context.Context, token string, id int, in domain.DeviceWrite) (*domain.Device, error)
	DeleteDeviceFunc             func(ctx context.Context, token string, id int) error
	CreateIPAddressFunc          func(ctx context.Context, token string, in domain.IPAddressWrite) (*domain.IPAddress, error)
	UpdateIPAddressFunc          func(ctx context.Context, token string, id int, in domain.IPAddressWrite) (*domain.IPAddress, error)
	DeleteIPAddressFunc          func(ctx context.Context, token string, id int) error
	CreatePrefixFunc             func(ctx context.Context, token string, in domain.PrefixWrite) (*domain.Prefix, error)
	UpdatePrefixFunc             func(ctx context.Context, token string, id int, in domain.PrefixWrite) (*domain.Prefix, error)
	DeletePrefixFunc             func(ctx context.Context, token string, id int) error
	CreateVLANFunc               func(ctx context.Context, token string, in domain.VLANWrite) (*domain.VLAN, error)
	UpdateVLANFunc               func(ctx context.Context, token string, id int, in domain.VLANWrite) (*domain.VLAN, error)
	DeleteVLANFunc               func(ctx context.Context, token string, id int) error
	CreateVirtualMachineFunc     func(ctx context.Context, token string, in domain.VirtualMachineWrite) (*domain.VirtualMachine, error)
	UpdateVirtualMachineFunc     func(ctx context.Context, token string, id int, in domain.VirtualMachineWrite) (*domain.VirtualMachine, error)
	DeleteVirtualMachineFunc     func(ctx context.Context, token string, id int) error
	CreateClusterFunc            func(ctx context.Context, token string, in domain.ClusterWrite) (*domain.Cluster, error)
	UpdateClusterFunc            func(ctx context.Context, token string, id int, in domain.ClusterWrite) (*domain.Cluster, error)
	DeleteClusterFunc            func(ctx context.Context, token string, id int) error
	CreateCircuitFunc            func(ctx context.Context, token string, in domain.CircuitWrite) (*domain.Circuit, error)
	UpdateCircuitFunc            func(ctx context.Context, token string, id int, in domain.CircuitWrite) (*domain.Circuit, error)
	DeleteCircuitFunc            func(ctx context.Context, token string, id int) error
	CreateRackFunc               func(ctx context.Context, token string, in domain.RackWrite) (*domain.Rack, error)
	UpdateRackFunc               func(ctx context.Context, token string, id int, in domain.RackWrite) (*domain.Rack, error)
	DeleteRackFunc               func(ctx context.Context, token string, id int) error
	CreateInterfaceFunc          func(ctx context.Context, token string, in domain.InterfaceWrite) (*domain.Interface, error)
	UpdateInterfaceFunc          func(ctx context.Context, token string, id int, in domain.InterfaceWrite) (*domain.Interface, error)
	DeleteInterfaceFunc          func(ctx context.Context, token string, id int) error
	CreateCircuitTerminationFunc func(ctx context.Context, token string, in domain.CircuitTerminationWrite) (*domain.CircuitTermination, error)
	UpdateCircuitTerminationFunc func(ctx context.Context, token string, id int, in domain.CircuitTerminationWrite) (*domain.CircuitTermination, error)
	DeleteCircuitTerminationFunc func(ctx context.Context, token string, id int) error
	CreateCableFunc              func(ctx context.Context, token string, in domain.CableWrite) (*domain.Cable, error)
	UpdateCableFunc              func(ctx context.Context, token string, id int, in domain.CableWrite) (*domain.Cable, error)
	DeleteCableFunc              func(ctx context.Context, token string, id int) error
	CreateVMInterfaceFunc        func(ctx context.Context, token string, in domain.VMInterfaceWrite) (*domain.VMInterface, error)
	UpdateVMInterfaceFunc        func(ctx context.Context, token string, id int, in domain.VMInterfaceWrite) (*domain.VMInterface, error)
	DeleteVMInterfaceFunc        func(ctx context.Context, token string, id int) error
	CreateProviderFunc           func(ctx context.Context, token string, in domain.ProviderWrite) (*domain.Provider, error)
	UpdateProviderFunc           func(ctx context.Context, token string, id int, in domain.ProviderWrite) (*domain.Provider, error)
	DeleteProviderFunc           func(ctx context.Context, token string, id int) error
	CreateTenantFunc             func(ctx context.Context, token string, in domain.TenantWrite) (*domain.Tenant, error)
	UpdateTenantFunc             func(ctx context.Context, token string, id int, in domain.TenantWrite) (*domain.Tenant, error)
	DeleteTenantFunc             func(ctx context.Context, token string, id int) error
	CreateManufacturerFunc       func(ctx context.Context, token string, in domain.ManufacturerWrite) (*domain.Manufacturer, error)
	UpdateManufacturerFunc       func(ctx context.Context, token string, id int, in domain.ManufacturerWrite) (*domain.Manufacturer, error)
	DeleteManufacturerFunc       func(ctx context.Context, token string, id int) error
	CreateDeviceTypeFunc         func(ctx context.Context, token string, in domain.DeviceTypeWrite) (*domain.DeviceType, error)
	UpdateDeviceTypeFunc         func(ctx context.Context, token string, id int, in domain.DeviceTypeWrite) (*domain.DeviceType, error)
	DeleteDeviceTypeFunc         func(ctx context.Context, token string, id int) error
	CreateLocationFunc           func(ctx context.Context, token string, in domain.LocationWrite) (*domain.Location, error)
	UpdateLocationFunc           func(ctx context.Context, token string, id int, in domain.LocationWrite) (*domain.Location, error)
	DeleteLocationFunc           func(ctx context.Context, token string, id int) error
	CreateClusterTypeFunc        func(ctx context.Context, token string, in domain.ClusterTypeWrite) (*domain.ClusterType, error)
	UpdateClusterTypeFunc        func(ctx context.Context, token string, id int, in domain.ClusterTypeWrite) (*domain.ClusterType, error)
	DeleteClusterTypeFunc        func(ctx context.Context, token string, id int) error
	CreateClusterGroupFunc       func(ctx context.Context, token string, in domain.ClusterGroupWrite) (*domain.ClusterGroup, error)
	UpdateClusterGroupFunc       func(ctx context.Context, token string, id int, in domain.ClusterGroupWrite) (*domain.ClusterGroup, error)
	DeleteClusterGroupFunc       func(ctx context.Context, token string, id int) error
	CreateCircuitTypeFunc        func(ctx context.Context, token string, in domain.CircuitTypeWrite) (*domain.CircuitType, error)
	UpdateCircuitTypeFunc        func(ctx context.Context, token string, id int, in domain.CircuitTypeWrite) (*domain.CircuitType, error)
	DeleteCircuitTypeFunc        func(ctx context.Context, token string, id int) error
	CreateVrfFunc                func(ctx context.Context, token string, in domain.VrfWrite) (*domain.Vrf, error)
	UpdateVrfFunc                func(ctx context.Context, token string, id int, in domain.VrfWrite) (*domain.Vrf, error)
	DeleteVrfFunc                func(ctx context.Context, token string, id int) error
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

func (m *MockRepo) CreateVLAN(ctx context.Context, token string, in domain.VLANWrite) (*domain.VLAN, error) {
	return m.CreateVLANFunc(ctx, token, in)
}

func (m *MockRepo) UpdateVLAN(ctx context.Context, token string, id int, in domain.VLANWrite) (*domain.VLAN, error) {
	return m.UpdateVLANFunc(ctx, token, id, in)
}

func (m *MockRepo) DeleteVLAN(ctx context.Context, token string, id int) error {
	return m.DeleteVLANFunc(ctx, token, id)
}

func (m *MockRepo) CreateVirtualMachine(ctx context.Context, token string, in domain.VirtualMachineWrite) (*domain.VirtualMachine, error) {
	return m.CreateVirtualMachineFunc(ctx, token, in)
}

func (m *MockRepo) UpdateVirtualMachine(ctx context.Context, token string, id int, in domain.VirtualMachineWrite) (*domain.VirtualMachine, error) {
	return m.UpdateVirtualMachineFunc(ctx, token, id, in)
}

func (m *MockRepo) DeleteVirtualMachine(ctx context.Context, token string, id int) error {
	return m.DeleteVirtualMachineFunc(ctx, token, id)
}

func (m *MockRepo) CreateCluster(ctx context.Context, token string, in domain.ClusterWrite) (*domain.Cluster, error) {
	return m.CreateClusterFunc(ctx, token, in)
}

func (m *MockRepo) UpdateCluster(ctx context.Context, token string, id int, in domain.ClusterWrite) (*domain.Cluster, error) {
	return m.UpdateClusterFunc(ctx, token, id, in)
}

func (m *MockRepo) DeleteCluster(ctx context.Context, token string, id int) error {
	return m.DeleteClusterFunc(ctx, token, id)
}

func (m *MockRepo) CreateCircuit(ctx context.Context, token string, in domain.CircuitWrite) (*domain.Circuit, error) {
	return m.CreateCircuitFunc(ctx, token, in)
}

func (m *MockRepo) UpdateCircuit(ctx context.Context, token string, id int, in domain.CircuitWrite) (*domain.Circuit, error) {
	return m.UpdateCircuitFunc(ctx, token, id, in)
}

func (m *MockRepo) DeleteCircuit(ctx context.Context, token string, id int) error {
	return m.DeleteCircuitFunc(ctx, token, id)
}

func (m *MockRepo) CreateRack(ctx context.Context, token string, in domain.RackWrite) (*domain.Rack, error) {
	return m.CreateRackFunc(ctx, token, in)
}

func (m *MockRepo) UpdateRack(ctx context.Context, token string, id int, in domain.RackWrite) (*domain.Rack, error) {
	return m.UpdateRackFunc(ctx, token, id, in)
}

func (m *MockRepo) DeleteRack(ctx context.Context, token string, id int) error {
	return m.DeleteRackFunc(ctx, token, id)
}

func (m *MockRepo) CreateInterface(ctx context.Context, token string, in domain.InterfaceWrite) (*domain.Interface, error) {
	return m.CreateInterfaceFunc(ctx, token, in)
}

func (m *MockRepo) UpdateInterface(ctx context.Context, token string, id int, in domain.InterfaceWrite) (*domain.Interface, error) {
	return m.UpdateInterfaceFunc(ctx, token, id, in)
}

func (m *MockRepo) DeleteInterface(ctx context.Context, token string, id int) error {
	return m.DeleteInterfaceFunc(ctx, token, id)
}

func (m *MockRepo) CreateCircuitTermination(ctx context.Context, token string, in domain.CircuitTerminationWrite) (*domain.CircuitTermination, error) {
	return m.CreateCircuitTerminationFunc(ctx, token, in)
}

func (m *MockRepo) UpdateCircuitTermination(ctx context.Context, token string, id int, in domain.CircuitTerminationWrite) (*domain.CircuitTermination, error) {
	return m.UpdateCircuitTerminationFunc(ctx, token, id, in)
}

func (m *MockRepo) DeleteCircuitTermination(ctx context.Context, token string, id int) error {
	return m.DeleteCircuitTerminationFunc(ctx, token, id)
}

func (m *MockRepo) CreateCable(ctx context.Context, token string, in domain.CableWrite) (*domain.Cable, error) {
	return m.CreateCableFunc(ctx, token, in)
}

func (m *MockRepo) UpdateCable(ctx context.Context, token string, id int, in domain.CableWrite) (*domain.Cable, error) {
	return m.UpdateCableFunc(ctx, token, id, in)
}

func (m *MockRepo) DeleteCable(ctx context.Context, token string, id int) error {
	return m.DeleteCableFunc(ctx, token, id)
}

func (m *MockRepo) CreateVMInterface(ctx context.Context, token string, in domain.VMInterfaceWrite) (*domain.VMInterface, error) {
	return m.CreateVMInterfaceFunc(ctx, token, in)
}

func (m *MockRepo) UpdateVMInterface(ctx context.Context, token string, id int, in domain.VMInterfaceWrite) (*domain.VMInterface, error) {
	return m.UpdateVMInterfaceFunc(ctx, token, id, in)
}

func (m *MockRepo) DeleteVMInterface(ctx context.Context, token string, id int) error {
	return m.DeleteVMInterfaceFunc(ctx, token, id)
}

func (m *MockRepo) CreateProvider(ctx context.Context, token string, in domain.ProviderWrite) (*domain.Provider, error) {
	return m.CreateProviderFunc(ctx, token, in)
}

func (m *MockRepo) UpdateProvider(ctx context.Context, token string, id int, in domain.ProviderWrite) (*domain.Provider, error) {
	return m.UpdateProviderFunc(ctx, token, id, in)
}

func (m *MockRepo) DeleteProvider(ctx context.Context, token string, id int) error {
	return m.DeleteProviderFunc(ctx, token, id)
}

func (m *MockRepo) CreateTenant(ctx context.Context, token string, in domain.TenantWrite) (*domain.Tenant, error) {
	return m.CreateTenantFunc(ctx, token, in)
}

func (m *MockRepo) UpdateTenant(ctx context.Context, token string, id int, in domain.TenantWrite) (*domain.Tenant, error) {
	return m.UpdateTenantFunc(ctx, token, id, in)
}

func (m *MockRepo) DeleteTenant(ctx context.Context, token string, id int) error {
	return m.DeleteTenantFunc(ctx, token, id)
}

func (m *MockRepo) CreateManufacturer(ctx context.Context, token string, in domain.ManufacturerWrite) (*domain.Manufacturer, error) {
	return m.CreateManufacturerFunc(ctx, token, in)
}

func (m *MockRepo) UpdateManufacturer(ctx context.Context, token string, id int, in domain.ManufacturerWrite) (*domain.Manufacturer, error) {
	return m.UpdateManufacturerFunc(ctx, token, id, in)
}

func (m *MockRepo) DeleteManufacturer(ctx context.Context, token string, id int) error {
	return m.DeleteManufacturerFunc(ctx, token, id)
}

func (m *MockRepo) CreateDeviceType(ctx context.Context, token string, in domain.DeviceTypeWrite) (*domain.DeviceType, error) {
	return m.CreateDeviceTypeFunc(ctx, token, in)
}

func (m *MockRepo) UpdateDeviceType(ctx context.Context, token string, id int, in domain.DeviceTypeWrite) (*domain.DeviceType, error) {
	return m.UpdateDeviceTypeFunc(ctx, token, id, in)
}

func (m *MockRepo) DeleteDeviceType(ctx context.Context, token string, id int) error {
	return m.DeleteDeviceTypeFunc(ctx, token, id)
}

func (m *MockRepo) CreateLocation(ctx context.Context, token string, in domain.LocationWrite) (*domain.Location, error) {
	return m.CreateLocationFunc(ctx, token, in)
}

func (m *MockRepo) UpdateLocation(ctx context.Context, token string, id int, in domain.LocationWrite) (*domain.Location, error) {
	return m.UpdateLocationFunc(ctx, token, id, in)
}

func (m *MockRepo) DeleteLocation(ctx context.Context, token string, id int) error {
	return m.DeleteLocationFunc(ctx, token, id)
}

func (m *MockRepo) CreateClusterType(ctx context.Context, token string, in domain.ClusterTypeWrite) (*domain.ClusterType, error) {
	return m.CreateClusterTypeFunc(ctx, token, in)
}

func (m *MockRepo) UpdateClusterType(ctx context.Context, token string, id int, in domain.ClusterTypeWrite) (*domain.ClusterType, error) {
	return m.UpdateClusterTypeFunc(ctx, token, id, in)
}

func (m *MockRepo) DeleteClusterType(ctx context.Context, token string, id int) error {
	return m.DeleteClusterTypeFunc(ctx, token, id)
}

func (m *MockRepo) CreateClusterGroup(ctx context.Context, token string, in domain.ClusterGroupWrite) (*domain.ClusterGroup, error) {
	return m.CreateClusterGroupFunc(ctx, token, in)
}

func (m *MockRepo) UpdateClusterGroup(ctx context.Context, token string, id int, in domain.ClusterGroupWrite) (*domain.ClusterGroup, error) {
	return m.UpdateClusterGroupFunc(ctx, token, id, in)
}

func (m *MockRepo) DeleteClusterGroup(ctx context.Context, token string, id int) error {
	return m.DeleteClusterGroupFunc(ctx, token, id)
}

func (m *MockRepo) CreateCircuitType(ctx context.Context, token string, in domain.CircuitTypeWrite) (*domain.CircuitType, error) {
	return m.CreateCircuitTypeFunc(ctx, token, in)
}

func (m *MockRepo) UpdateCircuitType(ctx context.Context, token string, id int, in domain.CircuitTypeWrite) (*domain.CircuitType, error) {
	return m.UpdateCircuitTypeFunc(ctx, token, id, in)
}

func (m *MockRepo) DeleteCircuitType(ctx context.Context, token string, id int) error {
	return m.DeleteCircuitTypeFunc(ctx, token, id)
}

func (m *MockRepo) CreateVrf(ctx context.Context, token string, in domain.VrfWrite) (*domain.Vrf, error) {
	return m.CreateVrfFunc(ctx, token, in)
}

func (m *MockRepo) UpdateVrf(ctx context.Context, token string, id int, in domain.VrfWrite) (*domain.Vrf, error) {
	return m.UpdateVrfFunc(ctx, token, id, in)
}

func (m *MockRepo) DeleteVrf(ctx context.Context, token string, id int) error {
	return m.DeleteVrfFunc(ctx, token, id)
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

func (r *NilRepo) CreateVLAN(ctx context.Context, token string, in domain.VLANWrite) (*domain.VLAN, error) {
	return &domain.VLAN{ID: 1, VID: in.VID, Name: in.Name}, nil
}

func (r *NilRepo) UpdateVLAN(ctx context.Context, token string, id int, in domain.VLANWrite) (*domain.VLAN, error) {
	return &domain.VLAN{ID: id, VID: in.VID, Name: in.Name}, nil
}

func (r *NilRepo) DeleteVLAN(ctx context.Context, token string, id int) error {
	return nil
}

func (r *NilRepo) CreateVirtualMachine(ctx context.Context, token string, in domain.VirtualMachineWrite) (*domain.VirtualMachine, error) {
	return &domain.VirtualMachine{ID: 1, Name: in.Name}, nil
}

func (r *NilRepo) UpdateVirtualMachine(ctx context.Context, token string, id int, in domain.VirtualMachineWrite) (*domain.VirtualMachine, error) {
	return &domain.VirtualMachine{ID: id, Name: in.Name}, nil
}

func (r *NilRepo) DeleteVirtualMachine(ctx context.Context, token string, id int) error {
	return nil
}

func (r *NilRepo) CreateCluster(ctx context.Context, token string, in domain.ClusterWrite) (*domain.Cluster, error) {
	return &domain.Cluster{ID: 1, Name: in.Name}, nil
}

func (r *NilRepo) UpdateCluster(ctx context.Context, token string, id int, in domain.ClusterWrite) (*domain.Cluster, error) {
	return &domain.Cluster{ID: id, Name: in.Name}, nil
}

func (r *NilRepo) DeleteCluster(ctx context.Context, token string, id int) error {
	return nil
}

func (r *NilRepo) CreateCircuit(ctx context.Context, token string, in domain.CircuitWrite) (*domain.Circuit, error) {
	return &domain.Circuit{ID: 1, CID: in.CID}, nil
}

func (r *NilRepo) UpdateCircuit(ctx context.Context, token string, id int, in domain.CircuitWrite) (*domain.Circuit, error) {
	return &domain.Circuit{ID: id, CID: in.CID}, nil
}

func (r *NilRepo) DeleteCircuit(ctx context.Context, token string, id int) error {
	return nil
}

func (r *NilRepo) CreateRack(ctx context.Context, token string, in domain.RackWrite) (*domain.Rack, error) {
	return &domain.Rack{ID: 1, Name: in.Name}, nil
}

func (r *NilRepo) UpdateRack(ctx context.Context, token string, id int, in domain.RackWrite) (*domain.Rack, error) {
	return &domain.Rack{ID: id, Name: in.Name}, nil
}

func (r *NilRepo) DeleteRack(ctx context.Context, token string, id int) error {
	return nil
}

func (r *NilRepo) CreateInterface(ctx context.Context, token string, in domain.InterfaceWrite) (*domain.Interface, error) {
	return &domain.Interface{ID: 1, Name: in.Name}, nil
}

func (r *NilRepo) UpdateInterface(ctx context.Context, token string, id int, in domain.InterfaceWrite) (*domain.Interface, error) {
	return &domain.Interface{ID: id, Name: in.Name}, nil
}

func (r *NilRepo) DeleteInterface(ctx context.Context, token string, id int) error {
	return nil
}

func (r *NilRepo) CreateCircuitTermination(ctx context.Context, token string, in domain.CircuitTerminationWrite) (*domain.CircuitTermination, error) {
	return &domain.CircuitTermination{ID: 1, TermSide: in.TermSide}, nil
}

func (r *NilRepo) UpdateCircuitTermination(ctx context.Context, token string, id int, in domain.CircuitTerminationWrite) (*domain.CircuitTermination, error) {
	return &domain.CircuitTermination{ID: id, TermSide: in.TermSide}, nil
}

func (r *NilRepo) DeleteCircuitTermination(ctx context.Context, token string, id int) error {
	return nil
}

func (r *NilRepo) CreateCable(ctx context.Context, token string, in domain.CableWrite) (*domain.Cable, error) {
	c := &domain.Cable{ID: 1}
	if in.Label != nil {
		c.Label = *in.Label
	}
	return c, nil
}

func (r *NilRepo) UpdateCable(ctx context.Context, token string, id int, in domain.CableWrite) (*domain.Cable, error) {
	c := &domain.Cable{ID: id}
	if in.Label != nil {
		c.Label = *in.Label
	}
	return c, nil
}

func (r *NilRepo) DeleteCable(ctx context.Context, token string, id int) error {
	return nil
}

func (r *NilRepo) CreateVMInterface(ctx context.Context, token string, in domain.VMInterfaceWrite) (*domain.VMInterface, error) {
	return &domain.VMInterface{ID: 1, Name: in.Name}, nil
}

func (r *NilRepo) UpdateVMInterface(ctx context.Context, token string, id int, in domain.VMInterfaceWrite) (*domain.VMInterface, error) {
	return &domain.VMInterface{ID: id, Name: in.Name}, nil
}

func (r *NilRepo) DeleteVMInterface(ctx context.Context, token string, id int) error {
	return nil
}

func (r *NilRepo) CreateProvider(ctx context.Context, token string, in domain.ProviderWrite) (*domain.Provider, error) {
	return &domain.Provider{ID: 1, Name: in.Name}, nil
}

func (r *NilRepo) UpdateProvider(ctx context.Context, token string, id int, in domain.ProviderWrite) (*domain.Provider, error) {
	return &domain.Provider{ID: id, Name: in.Name}, nil
}

func (r *NilRepo) DeleteProvider(ctx context.Context, token string, id int) error {
	return nil
}

func (r *NilRepo) CreateTenant(ctx context.Context, token string, in domain.TenantWrite) (*domain.Tenant, error) {
	return &domain.Tenant{ID: 1, Name: in.Name}, nil
}

func (r *NilRepo) UpdateTenant(ctx context.Context, token string, id int, in domain.TenantWrite) (*domain.Tenant, error) {
	return &domain.Tenant{ID: id, Name: in.Name}, nil
}

func (r *NilRepo) DeleteTenant(ctx context.Context, token string, id int) error {
	return nil
}

func (r *NilRepo) CreateManufacturer(ctx context.Context, token string, in domain.ManufacturerWrite) (*domain.Manufacturer, error) {
	return &domain.Manufacturer{ID: 1, Name: in.Name}, nil
}

func (r *NilRepo) UpdateManufacturer(ctx context.Context, token string, id int, in domain.ManufacturerWrite) (*domain.Manufacturer, error) {
	return &domain.Manufacturer{ID: id, Name: in.Name}, nil
}

func (r *NilRepo) DeleteManufacturer(ctx context.Context, token string, id int) error {
	return nil
}

func (r *NilRepo) CreateDeviceType(ctx context.Context, token string, in domain.DeviceTypeWrite) (*domain.DeviceType, error) {
	return &domain.DeviceType{ID: 1, Model: in.Model}, nil
}

func (r *NilRepo) UpdateDeviceType(ctx context.Context, token string, id int, in domain.DeviceTypeWrite) (*domain.DeviceType, error) {
	return &domain.DeviceType{ID: id, Model: in.Model}, nil
}

func (r *NilRepo) DeleteDeviceType(ctx context.Context, token string, id int) error {
	return nil
}

func (r *NilRepo) CreateLocation(ctx context.Context, token string, in domain.LocationWrite) (*domain.Location, error) {
	return &domain.Location{ID: 1, Name: in.Name}, nil
}

func (r *NilRepo) UpdateLocation(ctx context.Context, token string, id int, in domain.LocationWrite) (*domain.Location, error) {
	return &domain.Location{ID: id, Name: in.Name}, nil
}

func (r *NilRepo) DeleteLocation(ctx context.Context, token string, id int) error {
	return nil
}

func (r *NilRepo) CreateClusterType(ctx context.Context, token string, in domain.ClusterTypeWrite) (*domain.ClusterType, error) {
	return &domain.ClusterType{ID: 1, Name: in.Name}, nil
}

func (r *NilRepo) UpdateClusterType(ctx context.Context, token string, id int, in domain.ClusterTypeWrite) (*domain.ClusterType, error) {
	return &domain.ClusterType{ID: id, Name: in.Name}, nil
}

func (r *NilRepo) DeleteClusterType(ctx context.Context, token string, id int) error {
	return nil
}

func (r *NilRepo) CreateClusterGroup(ctx context.Context, token string, in domain.ClusterGroupWrite) (*domain.ClusterGroup, error) {
	return &domain.ClusterGroup{ID: 1, Name: in.Name}, nil
}

func (r *NilRepo) UpdateClusterGroup(ctx context.Context, token string, id int, in domain.ClusterGroupWrite) (*domain.ClusterGroup, error) {
	return &domain.ClusterGroup{ID: id, Name: in.Name}, nil
}

func (r *NilRepo) DeleteClusterGroup(ctx context.Context, token string, id int) error {
	return nil
}

func (r *NilRepo) CreateCircuitType(ctx context.Context, token string, in domain.CircuitTypeWrite) (*domain.CircuitType, error) {
	return &domain.CircuitType{ID: 1, Name: in.Name}, nil
}

func (r *NilRepo) UpdateCircuitType(ctx context.Context, token string, id int, in domain.CircuitTypeWrite) (*domain.CircuitType, error) {
	return &domain.CircuitType{ID: id, Name: in.Name}, nil
}

func (r *NilRepo) DeleteCircuitType(ctx context.Context, token string, id int) error {
	return nil
}

func (r *NilRepo) CreateVrf(ctx context.Context, token string, in domain.VrfWrite) (*domain.Vrf, error) {
	return &domain.Vrf{ID: 1, Name: in.Name, Rd: in.Rd}, nil
}

func (r *NilRepo) UpdateVrf(ctx context.Context, token string, id int, in domain.VrfWrite) (*domain.Vrf, error) {
	return &domain.Vrf{ID: id, Name: in.Name, Rd: in.Rd}, nil
}

func (r *NilRepo) DeleteVrf(ctx context.Context, token string, id int) error {
	return nil
}
