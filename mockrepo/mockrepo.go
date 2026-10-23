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
	ListSitesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.Site], error)
	ListRegionsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.Region], error)
	ListDevicesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.Device], error)
	ListIPAddressesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.IPAddress], error)
	ListPrefixesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.Prefix], error)
	ListVLANsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.VLAN], error)
	ListVirtualMachinesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.VirtualMachine], error)
	ListClustersFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.Cluster], error)
	ListCircuitsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.Circuit], error)
	ListCircuitTerminationsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.CircuitTermination], error)
	ListCablesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.Cable], error)
	ListRacksFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.Rack], error)
	ListInterfacesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.Interface], error)
	ListVMInterfacesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.VMInterface], error)
	GetObjectFunc func(
		ctx context.Context, token string, objectType string, id int, params map[string]string,
	) (domain.RawObject, error)
	CreateSiteFunc func(
		ctx context.Context, token string, in domain.SiteWrite,
	) (*domain.Site, error)
	UpdateSiteFunc func(
		ctx context.Context, token string, id int, in domain.SiteWrite,
	) (*domain.Site, error)
	DeleteSiteFunc   func(ctx context.Context, token string, id int) error
	CreateRegionFunc func(
		ctx context.Context, token string, in domain.RegionWrite,
	) (*domain.Region, error)
	UpdateRegionFunc func(
		ctx context.Context, token string, id int, in domain.RegionWrite,
	) (*domain.Region, error)
	DeleteRegionFunc     func(ctx context.Context, token string, id int) error
	ListTenantGroupsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.TenantGroup], error)
	CreateTenantGroupFunc func(
		ctx context.Context, token string, in domain.TenantGroupWrite,
	) (*domain.TenantGroup, error)
	UpdateTenantGroupFunc func(
		ctx context.Context, token string, id int, in domain.TenantGroupWrite,
	) (*domain.TenantGroup, error)
	DeleteTenantGroupFunc func(ctx context.Context, token string, id int) error
	CreateDeviceFunc      func(
		ctx context.Context, token string, in domain.DeviceWrite,
	) (*domain.Device, error)
	UpdateDeviceFunc func(
		ctx context.Context, token string, id int, in domain.DeviceWrite,
	) (*domain.Device, error)
	DeleteDeviceFunc    func(ctx context.Context, token string, id int) error
	CreateIPAddressFunc func(
		ctx context.Context, token string, in domain.IPAddressWrite,
	) (*domain.IPAddress, error)
	UpdateIPAddressFunc func(
		ctx context.Context, token string, id int, in domain.IPAddressWrite,
	) (*domain.IPAddress, error)
	DeleteIPAddressFunc func(ctx context.Context, token string, id int) error
	CreatePrefixFunc    func(
		ctx context.Context, token string, in domain.PrefixWrite,
	) (*domain.Prefix, error)
	UpdatePrefixFunc func(
		ctx context.Context, token string, id int, in domain.PrefixWrite,
	) (*domain.Prefix, error)
	DeletePrefixFunc func(ctx context.Context, token string, id int) error
	CreateVLANFunc   func(
		ctx context.Context, token string, in domain.VLANWrite,
	) (*domain.VLAN, error)
	UpdateVLANFunc func(
		ctx context.Context, token string, id int, in domain.VLANWrite,
	) (*domain.VLAN, error)
	DeleteVLANFunc           func(ctx context.Context, token string, id int) error
	CreateVirtualMachineFunc func(
		ctx context.Context, token string, in domain.VirtualMachineWrite,
	) (*domain.VirtualMachine, error)
	UpdateVirtualMachineFunc func(
		ctx context.Context, token string, id int, in domain.VirtualMachineWrite,
	) (*domain.VirtualMachine, error)
	DeleteVirtualMachineFunc func(ctx context.Context, token string, id int) error
	CreateClusterFunc        func(
		ctx context.Context, token string, in domain.ClusterWrite,
	) (*domain.Cluster, error)
	UpdateClusterFunc func(
		ctx context.Context, token string, id int, in domain.ClusterWrite,
	) (*domain.Cluster, error)
	DeleteClusterFunc func(ctx context.Context, token string, id int) error
	CreateCircuitFunc func(
		ctx context.Context, token string, in domain.CircuitWrite,
	) (*domain.Circuit, error)
	UpdateCircuitFunc func(
		ctx context.Context, token string, id int, in domain.CircuitWrite,
	) (*domain.Circuit, error)
	DeleteCircuitFunc func(ctx context.Context, token string, id int) error
	CreateRackFunc    func(
		ctx context.Context, token string, in domain.RackWrite,
	) (*domain.Rack, error)
	UpdateRackFunc func(
		ctx context.Context, token string, id int, in domain.RackWrite,
	) (*domain.Rack, error)
	DeleteRackFunc      func(ctx context.Context, token string, id int) error
	CreateInterfaceFunc func(
		ctx context.Context, token string, in domain.InterfaceWrite,
	) (*domain.Interface, error)
	UpdateInterfaceFunc func(
		ctx context.Context, token string, id int, in domain.InterfaceWrite,
	) (*domain.Interface, error)
	DeleteInterfaceFunc          func(ctx context.Context, token string, id int) error
	CreateCircuitTerminationFunc func(
		ctx context.Context, token string, in domain.CircuitTerminationWrite,
	) (*domain.CircuitTermination, error)
	UpdateCircuitTerminationFunc func(
		ctx context.Context, token string, id int, in domain.CircuitTerminationWrite,
	) (*domain.CircuitTermination, error)
	DeleteCircuitTerminationFunc func(ctx context.Context, token string, id int) error
	CreateCableFunc              func(
		ctx context.Context, token string, in domain.CableWrite,
	) (*domain.Cable, error)
	UpdateCableFunc func(
		ctx context.Context, token string, id int, in domain.CableWrite,
	) (*domain.Cable, error)
	DeleteCableFunc       func(ctx context.Context, token string, id int) error
	CreateVMInterfaceFunc func(
		ctx context.Context, token string, in domain.VMInterfaceWrite,
	) (*domain.VMInterface, error)
	UpdateVMInterfaceFunc func(
		ctx context.Context, token string, id int, in domain.VMInterfaceWrite,
	) (*domain.VMInterface, error)
	DeleteVMInterfaceFunc func(ctx context.Context, token string, id int) error
	CreateProviderFunc    func(
		ctx context.Context, token string, in domain.ProviderWrite,
	) (*domain.Provider, error)
	UpdateProviderFunc func(
		ctx context.Context, token string, id int, in domain.ProviderWrite,
	) (*domain.Provider, error)
	DeleteProviderFunc func(ctx context.Context, token string, id int) error
	CreateTenantFunc   func(
		ctx context.Context, token string, in domain.TenantWrite,
	) (*domain.Tenant, error)
	UpdateTenantFunc func(
		ctx context.Context, token string, id int, in domain.TenantWrite,
	) (*domain.Tenant, error)
	DeleteTenantFunc       func(ctx context.Context, token string, id int) error
	CreateManufacturerFunc func(
		ctx context.Context, token string, in domain.ManufacturerWrite,
	) (*domain.Manufacturer, error)
	UpdateManufacturerFunc func(
		ctx context.Context, token string, id int, in domain.ManufacturerWrite,
	) (*domain.Manufacturer, error)
	DeleteManufacturerFunc func(ctx context.Context, token string, id int) error
	CreateDeviceTypeFunc   func(
		ctx context.Context, token string, in domain.DeviceTypeWrite,
	) (*domain.DeviceType, error)
	UpdateDeviceTypeFunc func(
		ctx context.Context, token string, id int, in domain.DeviceTypeWrite,
	) (*domain.DeviceType, error)
	DeleteDeviceTypeFunc func(ctx context.Context, token string, id int) error
	CreateLocationFunc   func(
		ctx context.Context, token string, in domain.LocationWrite,
	) (*domain.Location, error)
	UpdateLocationFunc func(
		ctx context.Context, token string, id int, in domain.LocationWrite,
	) (*domain.Location, error)
	DeleteLocationFunc    func(ctx context.Context, token string, id int) error
	CreateClusterTypeFunc func(
		ctx context.Context, token string, in domain.ClusterTypeWrite,
	) (*domain.ClusterType, error)
	UpdateClusterTypeFunc func(
		ctx context.Context, token string, id int, in domain.ClusterTypeWrite,
	) (*domain.ClusterType, error)
	DeleteClusterTypeFunc  func(ctx context.Context, token string, id int) error
	CreateClusterGroupFunc func(
		ctx context.Context, token string, in domain.ClusterGroupWrite,
	) (*domain.ClusterGroup, error)
	UpdateClusterGroupFunc func(
		ctx context.Context, token string, id int, in domain.ClusterGroupWrite,
	) (*domain.ClusterGroup, error)
	DeleteClusterGroupFunc func(ctx context.Context, token string, id int) error
	CreateCircuitTypeFunc  func(
		ctx context.Context, token string, in domain.CircuitTypeWrite,
	) (*domain.CircuitType, error)
	UpdateCircuitTypeFunc func(
		ctx context.Context, token string, id int, in domain.CircuitTypeWrite,
	) (*domain.CircuitType, error)
	DeleteCircuitTypeFunc func(ctx context.Context, token string, id int) error
	CreateVrfFunc         func(
		ctx context.Context, token string, in domain.VrfWrite,
	) (*domain.Vrf, error)
	UpdateVrfFunc func(
		ctx context.Context, token string, id int, in domain.VrfWrite,
	) (*domain.Vrf, error)
	DeleteVrfFunc       func(ctx context.Context, token string, id int) error
	CreateVlanGroupFunc func(
		ctx context.Context, token string, in domain.VlanGroupWrite,
	) (*domain.VlanGroup, error)
	UpdateVlanGroupFunc func(
		ctx context.Context, token string, id int, in domain.VlanGroupWrite,
	) (*domain.VlanGroup, error)
	DeleteVlanGroupFunc func(ctx context.Context, token string, id int) error
	CreateRoleFunc      func(
		ctx context.Context, token string, in domain.RoleWrite,
	) (*domain.Role, error)
	UpdateRoleFunc func(
		ctx context.Context, token string, id int, in domain.RoleWrite,
	) (*domain.Role, error)
	DeleteRoleFunc    func(ctx context.Context, token string, id int) error
	CreateContactFunc func(
		ctx context.Context, token string, in domain.ContactWrite,
	) (*domain.Contact, error)
	UpdateContactFunc func(
		ctx context.Context, token string, id int, in domain.ContactWrite,
	) (*domain.Contact, error)
	DeleteContactFunc func(ctx context.Context, token string, id int) error
}

// ListSites delegates to ListSitesFunc.
func (m *MockRepo) ListSites(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.Site], error) {
	return m.ListSitesFunc(ctx, token, params)
}

// ListRegions delegates to ListRegionsFunc.
func (m *MockRepo) ListRegions(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.Region], error) {
	return m.ListRegionsFunc(ctx, token, params)
}

// ListDevices delegates to ListDevicesFunc.
func (m *MockRepo) ListDevices(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.Device], error) {
	return m.ListDevicesFunc(ctx, token, params)
}

// ListIPAddresses delegates to ListIPAddressesFunc.
func (m *MockRepo) ListIPAddresses(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.IPAddress], error) {
	return m.ListIPAddressesFunc(ctx, token, params)
}

// ListPrefixes delegates to ListPrefixesFunc.
func (m *MockRepo) ListPrefixes(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.Prefix], error) {
	return m.ListPrefixesFunc(ctx, token, params)
}

// ListVLANs delegates to ListVLANsFunc.
func (m *MockRepo) ListVLANs(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.VLAN], error) {
	return m.ListVLANsFunc(ctx, token, params)
}

// ListVirtualMachines delegates to ListVirtualMachinesFunc.
func (m *MockRepo) ListVirtualMachines(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.VirtualMachine], error) {
	return m.ListVirtualMachinesFunc(ctx, token, params)
}

// ListClusters delegates to ListClustersFunc.
func (m *MockRepo) ListClusters(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.Cluster], error) {
	return m.ListClustersFunc(ctx, token, params)
}

// ListCircuits delegates to ListCircuitsFunc.
func (m *MockRepo) ListCircuits(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.Circuit], error) {
	return m.ListCircuitsFunc(ctx, token, params)
}

// ListCircuitTerminations delegates to ListCircuitTerminationsFunc.
func (m *MockRepo) ListCircuitTerminations(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.CircuitTermination], error) {
	return m.ListCircuitTerminationsFunc(ctx, token, params)
}

// ListCables delegates to ListCablesFunc.
func (m *MockRepo) ListCables(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.Cable], error) {
	return m.ListCablesFunc(ctx, token, params)
}

// ListRacks delegates to ListRacksFunc.
func (m *MockRepo) ListRacks(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.Rack], error) {
	return m.ListRacksFunc(ctx, token, params)
}

// ListInterfaces delegates to ListInterfacesFunc.
func (m *MockRepo) ListInterfaces(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.Interface], error) {
	return m.ListInterfacesFunc(ctx, token, params)
}

// ListVMInterfaces delegates to ListVMInterfacesFunc.
func (m *MockRepo) ListVMInterfaces(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.VMInterface], error) {
	return m.ListVMInterfacesFunc(ctx, token, params)
}

// GetObject delegates to GetObjectFunc.
func (m *MockRepo) GetObject(
	ctx context.Context, token string, objectType string, id int, params map[string]string,
) (domain.RawObject, error) {
	return m.GetObjectFunc(ctx, token, objectType, id, params)
}

// CreateSite delegates to CreateSiteFunc.
func (m *MockRepo) CreateSite(
	ctx context.Context, token string, in domain.SiteWrite,
) (*domain.Site, error) {
	return m.CreateSiteFunc(ctx, token, in)
}

// UpdateSite delegates to UpdateSiteFunc.
func (m *MockRepo) UpdateSite(
	ctx context.Context, token string, id int, in domain.SiteWrite,
) (*domain.Site, error) {
	return m.UpdateSiteFunc(ctx, token, id, in)
}

// DeleteSite delegates to DeleteSiteFunc.
func (m *MockRepo) DeleteSite(ctx context.Context, token string, id int) error {
	return m.DeleteSiteFunc(ctx, token, id)
}

// CreateRegion delegates to CreateRegionFunc.
func (m *MockRepo) CreateRegion(
	ctx context.Context, token string, in domain.RegionWrite,
) (*domain.Region, error) {
	return m.CreateRegionFunc(ctx, token, in)
}

// UpdateRegion delegates to UpdateRegionFunc.
func (m *MockRepo) UpdateRegion(
	ctx context.Context, token string, id int, in domain.RegionWrite,
) (*domain.Region, error) {
	return m.UpdateRegionFunc(ctx, token, id, in)
}

// DeleteRegion delegates to DeleteRegionFunc.
func (m *MockRepo) DeleteRegion(ctx context.Context, token string, id int) error {
	return m.DeleteRegionFunc(ctx, token, id)
}

// ListTenantGroups delegates to ListTenantGroupsFunc.
func (m *MockRepo) ListTenantGroups(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.TenantGroup], error) {
	return m.ListTenantGroupsFunc(ctx, token, params)
}

// CreateTenantGroup delegates to CreateTenantGroupFunc.
func (m *MockRepo) CreateTenantGroup(
	ctx context.Context, token string, in domain.TenantGroupWrite,
) (*domain.TenantGroup, error) {
	return m.CreateTenantGroupFunc(ctx, token, in)
}

// UpdateTenantGroup delegates to UpdateTenantGroupFunc.
func (m *MockRepo) UpdateTenantGroup(
	ctx context.Context, token string, id int, in domain.TenantGroupWrite,
) (*domain.TenantGroup, error) {
	return m.UpdateTenantGroupFunc(ctx, token, id, in)
}

// DeleteTenantGroup delegates to DeleteTenantGroupFunc.
func (m *MockRepo) DeleteTenantGroup(ctx context.Context, token string, id int) error {
	return m.DeleteTenantGroupFunc(ctx, token, id)
}

// CreateDevice delegates to CreateDeviceFunc.
func (m *MockRepo) CreateDevice(
	ctx context.Context, token string, in domain.DeviceWrite,
) (*domain.Device, error) {
	return m.CreateDeviceFunc(ctx, token, in)
}

// UpdateDevice delegates to UpdateDeviceFunc.
func (m *MockRepo) UpdateDevice(
	ctx context.Context, token string, id int, in domain.DeviceWrite,
) (*domain.Device, error) {
	return m.UpdateDeviceFunc(ctx, token, id, in)
}

// DeleteDevice delegates to DeleteDeviceFunc.
func (m *MockRepo) DeleteDevice(ctx context.Context, token string, id int) error {
	return m.DeleteDeviceFunc(ctx, token, id)
}

// CreateIPAddress delegates to CreateIPAddressFunc.
func (m *MockRepo) CreateIPAddress(
	ctx context.Context, token string, in domain.IPAddressWrite,
) (*domain.IPAddress, error) {
	return m.CreateIPAddressFunc(ctx, token, in)
}

// UpdateIPAddress delegates to UpdateIPAddressFunc.
func (m *MockRepo) UpdateIPAddress(
	ctx context.Context, token string, id int, in domain.IPAddressWrite,
) (*domain.IPAddress, error) {
	return m.UpdateIPAddressFunc(ctx, token, id, in)
}

// DeleteIPAddress delegates to DeleteIPAddressFunc.
func (m *MockRepo) DeleteIPAddress(ctx context.Context, token string, id int) error {
	return m.DeleteIPAddressFunc(ctx, token, id)
}

// CreatePrefix delegates to CreatePrefixFunc.
func (m *MockRepo) CreatePrefix(
	ctx context.Context, token string, in domain.PrefixWrite,
) (*domain.Prefix, error) {
	return m.CreatePrefixFunc(ctx, token, in)
}

// UpdatePrefix delegates to UpdatePrefixFunc.
func (m *MockRepo) UpdatePrefix(
	ctx context.Context, token string, id int, in domain.PrefixWrite,
) (*domain.Prefix, error) {
	return m.UpdatePrefixFunc(ctx, token, id, in)
}

// DeletePrefix delegates to DeletePrefixFunc.
func (m *MockRepo) DeletePrefix(ctx context.Context, token string, id int) error {
	return m.DeletePrefixFunc(ctx, token, id)
}

// CreateVLAN delegates to CreateVLANFunc.
func (m *MockRepo) CreateVLAN(
	ctx context.Context, token string, in domain.VLANWrite,
) (*domain.VLAN, error) {
	return m.CreateVLANFunc(ctx, token, in)
}

// UpdateVLAN delegates to UpdateVLANFunc.
func (m *MockRepo) UpdateVLAN(
	ctx context.Context, token string, id int, in domain.VLANWrite,
) (*domain.VLAN, error) {
	return m.UpdateVLANFunc(ctx, token, id, in)
}

// DeleteVLAN delegates to DeleteVLANFunc.
func (m *MockRepo) DeleteVLAN(ctx context.Context, token string, id int) error {
	return m.DeleteVLANFunc(ctx, token, id)
}

// CreateVirtualMachine delegates to CreateVirtualMachineFunc.
func (m *MockRepo) CreateVirtualMachine(
	ctx context.Context, token string, in domain.VirtualMachineWrite,
) (*domain.VirtualMachine, error) {
	return m.CreateVirtualMachineFunc(ctx, token, in)
}

// UpdateVirtualMachine delegates to UpdateVirtualMachineFunc.
func (m *MockRepo) UpdateVirtualMachine(
	ctx context.Context, token string, id int, in domain.VirtualMachineWrite,
) (*domain.VirtualMachine, error) {
	return m.UpdateVirtualMachineFunc(ctx, token, id, in)
}

// DeleteVirtualMachine delegates to DeleteVirtualMachineFunc.
func (m *MockRepo) DeleteVirtualMachine(ctx context.Context, token string, id int) error {
	return m.DeleteVirtualMachineFunc(ctx, token, id)
}

// CreateCluster delegates to CreateClusterFunc.
func (m *MockRepo) CreateCluster(
	ctx context.Context, token string, in domain.ClusterWrite,
) (*domain.Cluster, error) {
	return m.CreateClusterFunc(ctx, token, in)
}

// UpdateCluster delegates to UpdateClusterFunc.
func (m *MockRepo) UpdateCluster(
	ctx context.Context, token string, id int, in domain.ClusterWrite,
) (*domain.Cluster, error) {
	return m.UpdateClusterFunc(ctx, token, id, in)
}

// DeleteCluster delegates to DeleteClusterFunc.
func (m *MockRepo) DeleteCluster(ctx context.Context, token string, id int) error {
	return m.DeleteClusterFunc(ctx, token, id)
}

// CreateCircuit delegates to CreateCircuitFunc.
func (m *MockRepo) CreateCircuit(
	ctx context.Context, token string, in domain.CircuitWrite,
) (*domain.Circuit, error) {
	return m.CreateCircuitFunc(ctx, token, in)
}

// UpdateCircuit delegates to UpdateCircuitFunc.
func (m *MockRepo) UpdateCircuit(
	ctx context.Context, token string, id int, in domain.CircuitWrite,
) (*domain.Circuit, error) {
	return m.UpdateCircuitFunc(ctx, token, id, in)
}

// DeleteCircuit delegates to DeleteCircuitFunc.
func (m *MockRepo) DeleteCircuit(ctx context.Context, token string, id int) error {
	return m.DeleteCircuitFunc(ctx, token, id)
}

// CreateRack delegates to CreateRackFunc.
func (m *MockRepo) CreateRack(
	ctx context.Context, token string, in domain.RackWrite,
) (*domain.Rack, error) {
	return m.CreateRackFunc(ctx, token, in)
}

// UpdateRack delegates to UpdateRackFunc.
func (m *MockRepo) UpdateRack(
	ctx context.Context, token string, id int, in domain.RackWrite,
) (*domain.Rack, error) {
	return m.UpdateRackFunc(ctx, token, id, in)
}

// DeleteRack delegates to DeleteRackFunc.
func (m *MockRepo) DeleteRack(ctx context.Context, token string, id int) error {
	return m.DeleteRackFunc(ctx, token, id)
}

// CreateInterface delegates to CreateInterfaceFunc.
func (m *MockRepo) CreateInterface(
	ctx context.Context, token string, in domain.InterfaceWrite,
) (*domain.Interface, error) {
	return m.CreateInterfaceFunc(ctx, token, in)
}

// UpdateInterface delegates to UpdateInterfaceFunc.
func (m *MockRepo) UpdateInterface(
	ctx context.Context, token string, id int, in domain.InterfaceWrite,
) (*domain.Interface, error) {
	return m.UpdateInterfaceFunc(ctx, token, id, in)
}

// DeleteInterface delegates to DeleteInterfaceFunc.
func (m *MockRepo) DeleteInterface(ctx context.Context, token string, id int) error {
	return m.DeleteInterfaceFunc(ctx, token, id)
}

// CreateCircuitTermination delegates to CreateCircuitTerminationFunc.
func (m *MockRepo) CreateCircuitTermination(
	ctx context.Context, token string, in domain.CircuitTerminationWrite,
) (*domain.CircuitTermination, error) {
	return m.CreateCircuitTerminationFunc(ctx, token, in)
}

// UpdateCircuitTermination delegates to UpdateCircuitTerminationFunc.
func (m *MockRepo) UpdateCircuitTermination(
	ctx context.Context, token string, id int, in domain.CircuitTerminationWrite,
) (*domain.CircuitTermination, error) {
	return m.UpdateCircuitTerminationFunc(ctx, token, id, in)
}

// DeleteCircuitTermination delegates to DeleteCircuitTerminationFunc.
func (m *MockRepo) DeleteCircuitTermination(ctx context.Context, token string, id int) error {
	return m.DeleteCircuitTerminationFunc(ctx, token, id)
}

// CreateCable delegates to CreateCableFunc.
func (m *MockRepo) CreateCable(
	ctx context.Context, token string, in domain.CableWrite,
) (*domain.Cable, error) {
	return m.CreateCableFunc(ctx, token, in)
}

// UpdateCable delegates to UpdateCableFunc.
func (m *MockRepo) UpdateCable(
	ctx context.Context, token string, id int, in domain.CableWrite,
) (*domain.Cable, error) {
	return m.UpdateCableFunc(ctx, token, id, in)
}

// DeleteCable delegates to DeleteCableFunc.
func (m *MockRepo) DeleteCable(ctx context.Context, token string, id int) error {
	return m.DeleteCableFunc(ctx, token, id)
}

// CreateVMInterface delegates to CreateVMInterfaceFunc.
func (m *MockRepo) CreateVMInterface(
	ctx context.Context, token string, in domain.VMInterfaceWrite,
) (*domain.VMInterface, error) {
	return m.CreateVMInterfaceFunc(ctx, token, in)
}

// UpdateVMInterface delegates to UpdateVMInterfaceFunc.
func (m *MockRepo) UpdateVMInterface(
	ctx context.Context, token string, id int, in domain.VMInterfaceWrite,
) (*domain.VMInterface, error) {
	return m.UpdateVMInterfaceFunc(ctx, token, id, in)
}

// DeleteVMInterface delegates to DeleteVMInterfaceFunc.
func (m *MockRepo) DeleteVMInterface(ctx context.Context, token string, id int) error {
	return m.DeleteVMInterfaceFunc(ctx, token, id)
}

// CreateProvider delegates to CreateProviderFunc.
func (m *MockRepo) CreateProvider(
	ctx context.Context, token string, in domain.ProviderWrite,
) (*domain.Provider, error) {
	return m.CreateProviderFunc(ctx, token, in)
}

// UpdateProvider delegates to UpdateProviderFunc.
func (m *MockRepo) UpdateProvider(
	ctx context.Context, token string, id int, in domain.ProviderWrite,
) (*domain.Provider, error) {
	return m.UpdateProviderFunc(ctx, token, id, in)
}

// DeleteProvider delegates to DeleteProviderFunc.
func (m *MockRepo) DeleteProvider(ctx context.Context, token string, id int) error {
	return m.DeleteProviderFunc(ctx, token, id)
}

// CreateTenant delegates to CreateTenantFunc.
func (m *MockRepo) CreateTenant(
	ctx context.Context, token string, in domain.TenantWrite,
) (*domain.Tenant, error) {
	return m.CreateTenantFunc(ctx, token, in)
}

// UpdateTenant delegates to UpdateTenantFunc.
func (m *MockRepo) UpdateTenant(
	ctx context.Context, token string, id int, in domain.TenantWrite,
) (*domain.Tenant, error) {
	return m.UpdateTenantFunc(ctx, token, id, in)
}

// DeleteTenant delegates to DeleteTenantFunc.
func (m *MockRepo) DeleteTenant(ctx context.Context, token string, id int) error {
	return m.DeleteTenantFunc(ctx, token, id)
}

// CreateManufacturer delegates to CreateManufacturerFunc.
func (m *MockRepo) CreateManufacturer(
	ctx context.Context, token string, in domain.ManufacturerWrite,
) (*domain.Manufacturer, error) {
	return m.CreateManufacturerFunc(ctx, token, in)
}

// UpdateManufacturer delegates to UpdateManufacturerFunc.
func (m *MockRepo) UpdateManufacturer(
	ctx context.Context, token string, id int, in domain.ManufacturerWrite,
) (*domain.Manufacturer, error) {
	return m.UpdateManufacturerFunc(ctx, token, id, in)
}

// DeleteManufacturer delegates to DeleteManufacturerFunc.
func (m *MockRepo) DeleteManufacturer(ctx context.Context, token string, id int) error {
	return m.DeleteManufacturerFunc(ctx, token, id)
}

// CreateDeviceType delegates to CreateDeviceTypeFunc.
func (m *MockRepo) CreateDeviceType(
	ctx context.Context, token string, in domain.DeviceTypeWrite,
) (*domain.DeviceType, error) {
	return m.CreateDeviceTypeFunc(ctx, token, in)
}

// UpdateDeviceType delegates to UpdateDeviceTypeFunc.
func (m *MockRepo) UpdateDeviceType(
	ctx context.Context, token string, id int, in domain.DeviceTypeWrite,
) (*domain.DeviceType, error) {
	return m.UpdateDeviceTypeFunc(ctx, token, id, in)
}

// DeleteDeviceType delegates to DeleteDeviceTypeFunc.
func (m *MockRepo) DeleteDeviceType(ctx context.Context, token string, id int) error {
	return m.DeleteDeviceTypeFunc(ctx, token, id)
}

// CreateLocation delegates to CreateLocationFunc.
func (m *MockRepo) CreateLocation(
	ctx context.Context, token string, in domain.LocationWrite,
) (*domain.Location, error) {
	return m.CreateLocationFunc(ctx, token, in)
}

// UpdateLocation delegates to UpdateLocationFunc.
func (m *MockRepo) UpdateLocation(
	ctx context.Context, token string, id int, in domain.LocationWrite,
) (*domain.Location, error) {
	return m.UpdateLocationFunc(ctx, token, id, in)
}

// DeleteLocation delegates to DeleteLocationFunc.
func (m *MockRepo) DeleteLocation(ctx context.Context, token string, id int) error {
	return m.DeleteLocationFunc(ctx, token, id)
}

// CreateClusterType delegates to CreateClusterTypeFunc.
func (m *MockRepo) CreateClusterType(
	ctx context.Context, token string, in domain.ClusterTypeWrite,
) (*domain.ClusterType, error) {
	return m.CreateClusterTypeFunc(ctx, token, in)
}

// UpdateClusterType delegates to UpdateClusterTypeFunc.
func (m *MockRepo) UpdateClusterType(
	ctx context.Context, token string, id int, in domain.ClusterTypeWrite,
) (*domain.ClusterType, error) {
	return m.UpdateClusterTypeFunc(ctx, token, id, in)
}

// DeleteClusterType delegates to DeleteClusterTypeFunc.
func (m *MockRepo) DeleteClusterType(ctx context.Context, token string, id int) error {
	return m.DeleteClusterTypeFunc(ctx, token, id)
}

// CreateClusterGroup delegates to CreateClusterGroupFunc.
func (m *MockRepo) CreateClusterGroup(
	ctx context.Context, token string, in domain.ClusterGroupWrite,
) (*domain.ClusterGroup, error) {
	return m.CreateClusterGroupFunc(ctx, token, in)
}

// UpdateClusterGroup delegates to UpdateClusterGroupFunc.
func (m *MockRepo) UpdateClusterGroup(
	ctx context.Context, token string, id int, in domain.ClusterGroupWrite,
) (*domain.ClusterGroup, error) {
	return m.UpdateClusterGroupFunc(ctx, token, id, in)
}

// DeleteClusterGroup delegates to DeleteClusterGroupFunc.
func (m *MockRepo) DeleteClusterGroup(ctx context.Context, token string, id int) error {
	return m.DeleteClusterGroupFunc(ctx, token, id)
}

// CreateCircuitType delegates to CreateCircuitTypeFunc.
func (m *MockRepo) CreateCircuitType(
	ctx context.Context, token string, in domain.CircuitTypeWrite,
) (*domain.CircuitType, error) {
	return m.CreateCircuitTypeFunc(ctx, token, in)
}

// UpdateCircuitType delegates to UpdateCircuitTypeFunc.
func (m *MockRepo) UpdateCircuitType(
	ctx context.Context, token string, id int, in domain.CircuitTypeWrite,
) (*domain.CircuitType, error) {
	return m.UpdateCircuitTypeFunc(ctx, token, id, in)
}

// DeleteCircuitType delegates to DeleteCircuitTypeFunc.
func (m *MockRepo) DeleteCircuitType(ctx context.Context, token string, id int) error {
	return m.DeleteCircuitTypeFunc(ctx, token, id)
}

// CreateVrf delegates to CreateVrfFunc.
func (m *MockRepo) CreateVrf(
	ctx context.Context, token string, in domain.VrfWrite,
) (*domain.Vrf, error) {
	return m.CreateVrfFunc(ctx, token, in)
}

// UpdateVrf delegates to UpdateVrfFunc.
func (m *MockRepo) UpdateVrf(
	ctx context.Context, token string, id int, in domain.VrfWrite,
) (*domain.Vrf, error) {
	return m.UpdateVrfFunc(ctx, token, id, in)
}

// DeleteVrf delegates to DeleteVrfFunc.
func (m *MockRepo) DeleteVrf(ctx context.Context, token string, id int) error {
	return m.DeleteVrfFunc(ctx, token, id)
}

// CreateVlanGroup delegates to CreateVlanGroupFunc.
func (m *MockRepo) CreateVlanGroup(
	ctx context.Context, token string, in domain.VlanGroupWrite,
) (*domain.VlanGroup, error) {
	return m.CreateVlanGroupFunc(ctx, token, in)
}

// UpdateVlanGroup delegates to UpdateVlanGroupFunc.
func (m *MockRepo) UpdateVlanGroup(
	ctx context.Context, token string, id int, in domain.VlanGroupWrite,
) (*domain.VlanGroup, error) {
	return m.UpdateVlanGroupFunc(ctx, token, id, in)
}

// DeleteVlanGroup delegates to DeleteVlanGroupFunc.
func (m *MockRepo) DeleteVlanGroup(ctx context.Context, token string, id int) error {
	return m.DeleteVlanGroupFunc(ctx, token, id)
}

// CreateRole delegates to CreateRoleFunc.
func (m *MockRepo) CreateRole(
	ctx context.Context, token string, in domain.RoleWrite,
) (*domain.Role, error) {
	return m.CreateRoleFunc(ctx, token, in)
}

// UpdateRole delegates to UpdateRoleFunc.
func (m *MockRepo) UpdateRole(
	ctx context.Context, token string, id int, in domain.RoleWrite,
) (*domain.Role, error) {
	return m.UpdateRoleFunc(ctx, token, id, in)
}

// DeleteRole delegates to DeleteRoleFunc.
func (m *MockRepo) DeleteRole(ctx context.Context, token string, id int) error {
	return m.DeleteRoleFunc(ctx, token, id)
}

// CreateContact delegates to CreateContactFunc.
func (m *MockRepo) CreateContact(
	ctx context.Context, token string, in domain.ContactWrite,
) (*domain.Contact, error) {
	return m.CreateContactFunc(ctx, token, in)
}

// UpdateContact delegates to UpdateContactFunc.
func (m *MockRepo) UpdateContact(
	ctx context.Context, token string, id int, in domain.ContactWrite,
) (*domain.Contact, error) {
	return m.UpdateContactFunc(ctx, token, id, in)
}

// DeleteContact delegates to DeleteContactFunc.
func (m *MockRepo) DeleteContact(ctx context.Context, token string, id int) error {
	return m.DeleteContactFunc(ctx, token, id)
}

// NilRepo is a mock implementation of domain.NetworkRepository
// that returns empty results for every method.
type NilRepo struct{}

// ListSites returns an empty site page.
func (r *NilRepo) ListSites(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Site], error) {
	return &domain.PaginatedResponse[domain.Site]{Count: 0, Results: []domain.Site{}}, nil
}

// ListRegions returns an empty region page.
func (r *NilRepo) ListRegions(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Region], error) {
	return &domain.PaginatedResponse[domain.Region]{Count: 0, Results: []domain.Region{}}, nil
}

// ListDevices returns an empty device page.
func (r *NilRepo) ListDevices(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Device], error) {
	return &domain.PaginatedResponse[domain.Device]{Count: 0, Results: []domain.Device{}}, nil
}

// ListIPAddresses returns an empty IP address page.
func (r *NilRepo) ListIPAddresses(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.IPAddress], error) {
	return &domain.PaginatedResponse[domain.IPAddress]{Count: 0, Results: []domain.IPAddress{}}, nil
}

// ListPrefixes returns an empty prefix page.
func (r *NilRepo) ListPrefixes(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Prefix], error) {
	return &domain.PaginatedResponse[domain.Prefix]{Count: 0, Results: []domain.Prefix{}}, nil
}

// ListVLANs returns an empty VLAN page.
func (r *NilRepo) ListVLANs(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.VLAN], error) {
	return &domain.PaginatedResponse[domain.VLAN]{Count: 0, Results: []domain.VLAN{}}, nil
}

// ListVirtualMachines returns an empty virtual machine page.
func (r *NilRepo) ListVirtualMachines(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.VirtualMachine], error) {
	return &domain.PaginatedResponse[domain.VirtualMachine]{Count: 0, Results: []domain.VirtualMachine{}}, nil
}

// ListClusters returns an empty cluster page.
func (r *NilRepo) ListClusters(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Cluster], error) {
	return &domain.PaginatedResponse[domain.Cluster]{Count: 0, Results: []domain.Cluster{}}, nil
}

// ListCircuits returns an empty circuit page.
func (r *NilRepo) ListCircuits(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Circuit], error) {
	return &domain.PaginatedResponse[domain.Circuit]{Count: 0, Results: []domain.Circuit{}}, nil
}

// ListCircuitTerminations returns an empty circuit termination page.
func (r *NilRepo) ListCircuitTerminations(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.CircuitTermination], error) {
	return &domain.PaginatedResponse[domain.CircuitTermination]{Count: 0, Results: []domain.CircuitTermination{}}, nil
}

// ListCables returns an empty cable page.
func (r *NilRepo) ListCables(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Cable], error) {
	return &domain.PaginatedResponse[domain.Cable]{Count: 0, Results: []domain.Cable{}}, nil
}

// ListRacks returns an empty rack page.
func (r *NilRepo) ListRacks(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Rack], error) {
	return &domain.PaginatedResponse[domain.Rack]{Count: 0, Results: []domain.Rack{}}, nil
}

// ListInterfaces returns an empty interface page.
func (r *NilRepo) ListInterfaces(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Interface], error) {
	return &domain.PaginatedResponse[domain.Interface]{Count: 0, Results: []domain.Interface{}}, nil
}

// ListVMInterfaces returns an empty VM interface page.
func (r *NilRepo) ListVMInterfaces(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.VMInterface], error) {
	return &domain.PaginatedResponse[domain.VMInterface]{Count: 0, Results: []domain.VMInterface{}}, nil
}

// GetObject returns a raw object with the given id.
func (r *NilRepo) GetObject(
	_ context.Context, _ string, _ string, id int, _ map[string]string,
) (domain.RawObject, error) {
	return domain.RawObject(fmt.Sprintf(`{"id":%d}`, id)), nil
}

// CreateSite returns a site with id 1.
func (r *NilRepo) CreateSite(
	_ context.Context, _ string, in domain.SiteWrite,
) (*domain.Site, error) {
	return &domain.Site{ID: 1, Name: in.Name}, nil
}

// UpdateSite returns the site with the given id.
func (r *NilRepo) UpdateSite(
	_ context.Context, _ string, id int, in domain.SiteWrite,
) (*domain.Site, error) {
	return &domain.Site{ID: id, Name: in.Name}, nil
}

// DeleteSite always succeeds.
func (r *NilRepo) DeleteSite(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateRegion returns a region with id 1.
func (r *NilRepo) CreateRegion(
	_ context.Context, _ string, in domain.RegionWrite,
) (*domain.Region, error) {
	return &domain.Region{ID: 1, Name: in.Name}, nil
}

// UpdateRegion returns the region with the given id.
func (r *NilRepo) UpdateRegion(
	_ context.Context, _ string, id int, in domain.RegionWrite,
) (*domain.Region, error) {
	return &domain.Region{ID: id, Name: in.Name}, nil
}

// DeleteRegion always succeeds.
func (r *NilRepo) DeleteRegion(_ context.Context, _ string, _ int) error {
	return nil
}

// ListTenantGroups returns an empty tenant group page.
func (r *NilRepo) ListTenantGroups(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.TenantGroup], error) {
	return &domain.PaginatedResponse[domain.TenantGroup]{Count: 0, Results: []domain.TenantGroup{}}, nil
}

// CreateTenantGroup returns a tenant group with id 1.
func (r *NilRepo) CreateTenantGroup(
	_ context.Context, _ string, in domain.TenantGroupWrite,
) (*domain.TenantGroup, error) {
	return &domain.TenantGroup{ID: 1, Name: in.Name}, nil
}

// UpdateTenantGroup returns the tenant group with the given id.
func (r *NilRepo) UpdateTenantGroup(
	_ context.Context, _ string, id int, in domain.TenantGroupWrite,
) (*domain.TenantGroup, error) {
	return &domain.TenantGroup{ID: id, Name: in.Name}, nil
}

// DeleteTenantGroup always succeeds.
func (r *NilRepo) DeleteTenantGroup(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateDevice returns a device with id 1.
func (r *NilRepo) CreateDevice(
	_ context.Context, _ string, in domain.DeviceWrite,
) (*domain.Device, error) {
	return &domain.Device{ID: 1, Name: in.Name}, nil
}

// UpdateDevice returns the device with the given id.
func (r *NilRepo) UpdateDevice(
	_ context.Context, _ string, id int, in domain.DeviceWrite,
) (*domain.Device, error) {
	return &domain.Device{ID: id, Name: in.Name}, nil
}

// DeleteDevice always succeeds.
func (r *NilRepo) DeleteDevice(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateIPAddress returns an IP address with id 1.
func (r *NilRepo) CreateIPAddress(
	_ context.Context, _ string, in domain.IPAddressWrite,
) (*domain.IPAddress, error) {
	return &domain.IPAddress{ID: 1, Address: in.Address}, nil
}

// UpdateIPAddress returns the IP address with the given id.
func (r *NilRepo) UpdateIPAddress(
	_ context.Context, _ string, id int, in domain.IPAddressWrite,
) (*domain.IPAddress, error) {
	return &domain.IPAddress{ID: id, Address: in.Address}, nil
}

// DeleteIPAddress always succeeds.
func (r *NilRepo) DeleteIPAddress(_ context.Context, _ string, _ int) error {
	return nil
}

// CreatePrefix returns a prefix with id 1.
func (r *NilRepo) CreatePrefix(
	_ context.Context, _ string, in domain.PrefixWrite,
) (*domain.Prefix, error) {
	return &domain.Prefix{ID: 1, Prefix: in.Prefix}, nil
}

// UpdatePrefix returns the prefix with the given id.
func (r *NilRepo) UpdatePrefix(
	_ context.Context, _ string, id int, in domain.PrefixWrite,
) (*domain.Prefix, error) {
	return &domain.Prefix{ID: id, Prefix: in.Prefix}, nil
}

// DeletePrefix always succeeds.
func (r *NilRepo) DeletePrefix(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateVLAN returns a VLAN with id 1.
func (r *NilRepo) CreateVLAN(
	_ context.Context, _ string, in domain.VLANWrite,
) (*domain.VLAN, error) {
	return &domain.VLAN{ID: 1, VID: in.VID, Name: in.Name}, nil
}

// UpdateVLAN returns the VLAN with the given id.
func (r *NilRepo) UpdateVLAN(
	_ context.Context, _ string, id int, in domain.VLANWrite,
) (*domain.VLAN, error) {
	return &domain.VLAN{ID: id, VID: in.VID, Name: in.Name}, nil
}

// DeleteVLAN always succeeds.
func (r *NilRepo) DeleteVLAN(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateVirtualMachine returns a virtual machine with id 1.
func (r *NilRepo) CreateVirtualMachine(
	_ context.Context, _ string, in domain.VirtualMachineWrite,
) (*domain.VirtualMachine, error) {
	return &domain.VirtualMachine{ID: 1, Name: in.Name}, nil
}

// UpdateVirtualMachine returns the virtual machine with the given id.
func (r *NilRepo) UpdateVirtualMachine(
	_ context.Context, _ string, id int, in domain.VirtualMachineWrite,
) (*domain.VirtualMachine, error) {
	return &domain.VirtualMachine{ID: id, Name: in.Name}, nil
}

// DeleteVirtualMachine always succeeds.
func (r *NilRepo) DeleteVirtualMachine(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateCluster returns a cluster with id 1.
func (r *NilRepo) CreateCluster(
	_ context.Context, _ string, in domain.ClusterWrite,
) (*domain.Cluster, error) {
	return &domain.Cluster{ID: 1, Name: in.Name}, nil
}

// UpdateCluster returns the cluster with the given id.
func (r *NilRepo) UpdateCluster(
	_ context.Context, _ string, id int, in domain.ClusterWrite,
) (*domain.Cluster, error) {
	return &domain.Cluster{ID: id, Name: in.Name}, nil
}

// DeleteCluster always succeeds.
func (r *NilRepo) DeleteCluster(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateCircuit returns a circuit with id 1.
func (r *NilRepo) CreateCircuit(
	_ context.Context, _ string, in domain.CircuitWrite,
) (*domain.Circuit, error) {
	return &domain.Circuit{ID: 1, CID: in.CID}, nil
}

// UpdateCircuit returns the circuit with the given id.
func (r *NilRepo) UpdateCircuit(
	_ context.Context, _ string, id int, in domain.CircuitWrite,
) (*domain.Circuit, error) {
	return &domain.Circuit{ID: id, CID: in.CID}, nil
}

// DeleteCircuit always succeeds.
func (r *NilRepo) DeleteCircuit(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateRack returns a rack with id 1.
func (r *NilRepo) CreateRack(
	_ context.Context, _ string, in domain.RackWrite,
) (*domain.Rack, error) {
	return &domain.Rack{ID: 1, Name: in.Name}, nil
}

// UpdateRack returns the rack with the given id.
func (r *NilRepo) UpdateRack(
	_ context.Context, _ string, id int, in domain.RackWrite,
) (*domain.Rack, error) {
	return &domain.Rack{ID: id, Name: in.Name}, nil
}

// DeleteRack always succeeds.
func (r *NilRepo) DeleteRack(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateInterface returns an interface with id 1.
func (r *NilRepo) CreateInterface(
	_ context.Context, _ string, in domain.InterfaceWrite,
) (*domain.Interface, error) {
	return &domain.Interface{ID: 1, Name: in.Name}, nil
}

// UpdateInterface returns the interface with the given id.
func (r *NilRepo) UpdateInterface(
	_ context.Context, _ string, id int, in domain.InterfaceWrite,
) (*domain.Interface, error) {
	return &domain.Interface{ID: id, Name: in.Name}, nil
}

// DeleteInterface always succeeds.
func (r *NilRepo) DeleteInterface(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateCircuitTermination returns a circuit termination with id 1.
func (r *NilRepo) CreateCircuitTermination(
	_ context.Context, _ string, in domain.CircuitTerminationWrite,
) (*domain.CircuitTermination, error) {
	return &domain.CircuitTermination{ID: 1, TermSide: in.TermSide}, nil
}

// UpdateCircuitTermination returns the circuit termination with the given id.
func (r *NilRepo) UpdateCircuitTermination(
	_ context.Context, _ string, id int, in domain.CircuitTerminationWrite,
) (*domain.CircuitTermination, error) {
	return &domain.CircuitTermination{ID: id, TermSide: in.TermSide}, nil
}

// DeleteCircuitTermination always succeeds.
func (r *NilRepo) DeleteCircuitTermination(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateCable returns a cable with id 1.
func (r *NilRepo) CreateCable(
	_ context.Context, _ string, in domain.CableWrite,
) (*domain.Cable, error) {
	c := &domain.Cable{ID: 1}
	if in.Label != nil {
		c.Label = *in.Label
	}

	return c, nil
}

// UpdateCable returns the cable with the given id.
func (r *NilRepo) UpdateCable(
	_ context.Context, _ string, id int, in domain.CableWrite,
) (*domain.Cable, error) {
	c := &domain.Cable{ID: id}
	if in.Label != nil {
		c.Label = *in.Label
	}

	return c, nil
}

// DeleteCable always succeeds.
func (r *NilRepo) DeleteCable(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateVMInterface returns a VM interface with id 1.
func (r *NilRepo) CreateVMInterface(
	_ context.Context, _ string, in domain.VMInterfaceWrite,
) (*domain.VMInterface, error) {
	return &domain.VMInterface{ID: 1, Name: in.Name}, nil
}

// UpdateVMInterface returns the VM interface with the given id.
func (r *NilRepo) UpdateVMInterface(
	_ context.Context, _ string, id int, in domain.VMInterfaceWrite,
) (*domain.VMInterface, error) {
	return &domain.VMInterface{ID: id, Name: in.Name}, nil
}

// DeleteVMInterface always succeeds.
func (r *NilRepo) DeleteVMInterface(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateProvider returns a provider with id 1.
func (r *NilRepo) CreateProvider(
	_ context.Context, _ string, in domain.ProviderWrite,
) (*domain.Provider, error) {
	return &domain.Provider{ID: 1, Name: in.Name}, nil
}

// UpdateProvider returns the provider with the given id.
func (r *NilRepo) UpdateProvider(
	_ context.Context, _ string, id int, in domain.ProviderWrite,
) (*domain.Provider, error) {
	return &domain.Provider{ID: id, Name: in.Name}, nil
}

// DeleteProvider always succeeds.
func (r *NilRepo) DeleteProvider(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateTenant returns a tenant with id 1.
func (r *NilRepo) CreateTenant(
	_ context.Context, _ string, in domain.TenantWrite,
) (*domain.Tenant, error) {
	return &domain.Tenant{ID: 1, Name: in.Name}, nil
}

// UpdateTenant returns the tenant with the given id.
func (r *NilRepo) UpdateTenant(
	_ context.Context, _ string, id int, in domain.TenantWrite,
) (*domain.Tenant, error) {
	return &domain.Tenant{ID: id, Name: in.Name}, nil
}

// DeleteTenant always succeeds.
func (r *NilRepo) DeleteTenant(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateManufacturer returns a manufacturer with id 1.
func (r *NilRepo) CreateManufacturer(
	_ context.Context, _ string, in domain.ManufacturerWrite,
) (*domain.Manufacturer, error) {
	return &domain.Manufacturer{ID: 1, Name: in.Name}, nil
}

// UpdateManufacturer returns the manufacturer with the given id.
func (r *NilRepo) UpdateManufacturer(
	_ context.Context, _ string, id int, in domain.ManufacturerWrite,
) (*domain.Manufacturer, error) {
	return &domain.Manufacturer{ID: id, Name: in.Name}, nil
}

// DeleteManufacturer always succeeds.
func (r *NilRepo) DeleteManufacturer(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateDeviceType returns a device type with id 1.
func (r *NilRepo) CreateDeviceType(
	_ context.Context, _ string, in domain.DeviceTypeWrite,
) (*domain.DeviceType, error) {
	return &domain.DeviceType{ID: 1, Model: in.Model}, nil
}

// UpdateDeviceType returns the device type with the given id.
func (r *NilRepo) UpdateDeviceType(
	_ context.Context, _ string, id int, in domain.DeviceTypeWrite,
) (*domain.DeviceType, error) {
	return &domain.DeviceType{ID: id, Model: in.Model}, nil
}

// DeleteDeviceType always succeeds.
func (r *NilRepo) DeleteDeviceType(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateLocation returns a location with id 1.
func (r *NilRepo) CreateLocation(
	_ context.Context, _ string, in domain.LocationWrite,
) (*domain.Location, error) {
	return &domain.Location{ID: 1, Name: in.Name}, nil
}

// UpdateLocation returns the location with the given id.
func (r *NilRepo) UpdateLocation(
	_ context.Context, _ string, id int, in domain.LocationWrite,
) (*domain.Location, error) {
	return &domain.Location{ID: id, Name: in.Name}, nil
}

// DeleteLocation always succeeds.
func (r *NilRepo) DeleteLocation(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateClusterType returns a cluster type with id 1.
func (r *NilRepo) CreateClusterType(
	_ context.Context, _ string, in domain.ClusterTypeWrite,
) (*domain.ClusterType, error) {
	return &domain.ClusterType{ID: 1, Name: in.Name}, nil
}

// UpdateClusterType returns the cluster type with the given id.
func (r *NilRepo) UpdateClusterType(
	_ context.Context, _ string, id int, in domain.ClusterTypeWrite,
) (*domain.ClusterType, error) {
	return &domain.ClusterType{ID: id, Name: in.Name}, nil
}

// DeleteClusterType always succeeds.
func (r *NilRepo) DeleteClusterType(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateClusterGroup returns a cluster group with id 1.
func (r *NilRepo) CreateClusterGroup(
	_ context.Context, _ string, in domain.ClusterGroupWrite,
) (*domain.ClusterGroup, error) {
	return &domain.ClusterGroup{ID: 1, Name: in.Name}, nil
}

// UpdateClusterGroup returns the cluster group with the given id.
func (r *NilRepo) UpdateClusterGroup(
	_ context.Context, _ string, id int, in domain.ClusterGroupWrite,
) (*domain.ClusterGroup, error) {
	return &domain.ClusterGroup{ID: id, Name: in.Name}, nil
}

// DeleteClusterGroup always succeeds.
func (r *NilRepo) DeleteClusterGroup(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateCircuitType returns a circuit type with id 1.
func (r *NilRepo) CreateCircuitType(
	_ context.Context, _ string, in domain.CircuitTypeWrite,
) (*domain.CircuitType, error) {
	return &domain.CircuitType{ID: 1, Name: in.Name}, nil
}

// UpdateCircuitType returns the circuit type with the given id.
func (r *NilRepo) UpdateCircuitType(
	_ context.Context, _ string, id int, in domain.CircuitTypeWrite,
) (*domain.CircuitType, error) {
	return &domain.CircuitType{ID: id, Name: in.Name}, nil
}

// DeleteCircuitType always succeeds.
func (r *NilRepo) DeleteCircuitType(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateVrf returns a VRF with id 1.
func (r *NilRepo) CreateVrf(
	_ context.Context, _ string, in domain.VrfWrite,
) (*domain.Vrf, error) {
	return &domain.Vrf{ID: 1, Name: in.Name, Rd: in.Rd}, nil
}

// UpdateVrf returns the VRF with the given id.
func (r *NilRepo) UpdateVrf(
	_ context.Context, _ string, id int, in domain.VrfWrite,
) (*domain.Vrf, error) {
	return &domain.Vrf{ID: id, Name: in.Name, Rd: in.Rd}, nil
}

// DeleteVrf always succeeds.
func (r *NilRepo) DeleteVrf(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateVlanGroup returns a VLAN group with id 1.
func (r *NilRepo) CreateVlanGroup(
	_ context.Context, _ string, in domain.VlanGroupWrite,
) (*domain.VlanGroup, error) {
	return &domain.VlanGroup{ID: 1, Name: in.Name}, nil
}

// UpdateVlanGroup returns the VLAN group with the given id.
func (r *NilRepo) UpdateVlanGroup(
	_ context.Context, _ string, id int, in domain.VlanGroupWrite,
) (*domain.VlanGroup, error) {
	return &domain.VlanGroup{ID: id, Name: in.Name}, nil
}

// DeleteVlanGroup always succeeds.
func (r *NilRepo) DeleteVlanGroup(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateRole returns a role with id 1.
func (r *NilRepo) CreateRole(
	_ context.Context, _ string, in domain.RoleWrite,
) (*domain.Role, error) {
	return &domain.Role{ID: 1, Name: in.Name}, nil
}

// UpdateRole returns the role with the given id.
func (r *NilRepo) UpdateRole(
	_ context.Context, _ string, id int, in domain.RoleWrite,
) (*domain.Role, error) {
	return &domain.Role{ID: id, Name: in.Name}, nil
}

// DeleteRole always succeeds.
func (r *NilRepo) DeleteRole(_ context.Context, _ string, _ int) error {
	return nil
}

// CreateContact returns a contact with id 1.
func (r *NilRepo) CreateContact(
	_ context.Context, _ string, in domain.ContactWrite,
) (*domain.Contact, error) {
	return &domain.Contact{ID: 1, Name: in.Name}, nil
}

// UpdateContact returns the contact with the given id.
func (r *NilRepo) UpdateContact(
	_ context.Context, _ string, id int, in domain.ContactWrite,
) (*domain.Contact, error) {
	return &domain.Contact{ID: id, Name: in.Name}, nil
}

// DeleteContact always succeeds.
func (r *NilRepo) DeleteContact(_ context.Context, _ string, _ int) error {
	return nil
}
