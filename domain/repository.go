package domain

import (
	"context"
	"encoding/json"
)

// RawObject is a JSON-preserving type for generic object retrieval via get_object_by_id.
// It preserves the exact JSON structure from NetBox, avoiding map[string]interface{} lossiness.
//
// It implements json.Marshaler and json.Unmarshaler by delegating to the
// underlying json.RawMessage, ensuring correct JSON-in-JSON serialization
// (rather than base64-encoded output) when used as a field in a struct.
type RawObject json.RawMessage

// MarshalJSON implements json.Marshaler so RawObject serialises as inline JSON
// instead of a base64-encoded string.
func (r RawObject) MarshalJSON() ([]byte, error) {
	return json.RawMessage(r).MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler.
func (r *RawObject) UnmarshalJSON(b []byte) error {
	return (*json.RawMessage)(r).UnmarshalJSON(b)
}

// NetworkRepository defines the port (interface) for NetBox data access.
// Implementations provide CRUD operations for sites, devices, IP addresses,
// prefixes, VLANs, virtual machines, clusters, circuits, and racks.
//
//nolint:interfacebloat // Single port mirrors NetBox entity surface; splitting would ripple.
type NetworkRepository interface {
	ListSites(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[Site], error)
	ListRegions(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[Region], error)
	ListTenantGroups(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[TenantGroup], error)
	ListDevices(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[Device], error)
	ListIPAddresses(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[IPAddress], error)
	ListPrefixes(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[Prefix], error)
	ListVLANs(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[VLAN], error)
	ListVirtualMachines(
		ctx context.Context, token string, params map[string]string,
	) (*PaginatedResponse[VirtualMachine], error)
	ListClusters(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[Cluster], error)
	ListCircuits(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[Circuit], error)
	ListRacks(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[Rack], error)
	ListInterfaces(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[Interface], error)
	ListVMInterfaces(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[VMInterface], error)
	ListCircuitTerminations(
		ctx context.Context, token string, params map[string]string,
	) (*PaginatedResponse[CircuitTermination], error)
	ListCables(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[Cable], error)
	GetObject(ctx context.Context, token string, objectType string, id int, params map[string]string) (RawObject, error)

	CreateSite(ctx context.Context, token string, in SiteWrite) (*Site, error)
	UpdateSite(ctx context.Context, token string, id int, in SiteWrite) (*Site, error)
	DeleteSite(ctx context.Context, token string, id int) error

	CreateRegion(ctx context.Context, token string, in RegionWrite) (*Region, error)
	UpdateRegion(ctx context.Context, token string, id int, in RegionWrite) (*Region, error)
	DeleteRegion(ctx context.Context, token string, id int) error

	CreateTenantGroup(ctx context.Context, token string, in TenantGroupWrite) (*TenantGroup, error)
	UpdateTenantGroup(ctx context.Context, token string, id int, in TenantGroupWrite) (*TenantGroup, error)
	DeleteTenantGroup(ctx context.Context, token string, id int) error

	CreateDevice(ctx context.Context, token string, in DeviceWrite) (*Device, error)
	UpdateDevice(ctx context.Context, token string, id int, in DeviceWrite) (*Device, error)
	DeleteDevice(ctx context.Context, token string, id int) error

	CreateIPAddress(ctx context.Context, token string, in IPAddressWrite) (*IPAddress, error)
	UpdateIPAddress(ctx context.Context, token string, id int, in IPAddressWrite) (*IPAddress, error)
	DeleteIPAddress(ctx context.Context, token string, id int) error

	CreatePrefix(ctx context.Context, token string, in PrefixWrite) (*Prefix, error)
	UpdatePrefix(ctx context.Context, token string, id int, in PrefixWrite) (*Prefix, error)
	DeletePrefix(ctx context.Context, token string, id int) error

	CreateVLAN(ctx context.Context, token string, in VLANWrite) (*VLAN, error)
	UpdateVLAN(ctx context.Context, token string, id int, in VLANWrite) (*VLAN, error)
	DeleteVLAN(ctx context.Context, token string, id int) error

	CreateVirtualMachine(ctx context.Context, token string, in VirtualMachineWrite) (*VirtualMachine, error)
	UpdateVirtualMachine(ctx context.Context, token string, id int, in VirtualMachineWrite) (*VirtualMachine, error)
	DeleteVirtualMachine(ctx context.Context, token string, id int) error

	CreateCluster(ctx context.Context, token string, in ClusterWrite) (*Cluster, error)
	UpdateCluster(ctx context.Context, token string, id int, in ClusterWrite) (*Cluster, error)
	DeleteCluster(ctx context.Context, token string, id int) error

	CreateCircuit(ctx context.Context, token string, in CircuitWrite) (*Circuit, error)
	UpdateCircuit(ctx context.Context, token string, id int, in CircuitWrite) (*Circuit, error)
	DeleteCircuit(ctx context.Context, token string, id int) error

	CreateRack(ctx context.Context, token string, in RackWrite) (*Rack, error)
	UpdateRack(ctx context.Context, token string, id int, in RackWrite) (*Rack, error)
	DeleteRack(ctx context.Context, token string, id int) error

	CreateInterface(ctx context.Context, token string, in InterfaceWrite) (*Interface, error)
	UpdateInterface(ctx context.Context, token string, id int, in InterfaceWrite) (*Interface, error)
	DeleteInterface(ctx context.Context, token string, id int) error

	CreateCircuitTermination(ctx context.Context, token string, in CircuitTerminationWrite) (*CircuitTermination, error)
	UpdateCircuitTermination(
		ctx context.Context, token string, id int, in CircuitTerminationWrite,
	) (*CircuitTermination, error)
	DeleteCircuitTermination(ctx context.Context, token string, id int) error

	CreateCable(ctx context.Context, token string, in CableWrite) (*Cable, error)
	UpdateCable(ctx context.Context, token string, id int, in CableWrite) (*Cable, error)
	DeleteCable(ctx context.Context, token string, id int) error

	CreateVMInterface(ctx context.Context, token string, in VMInterfaceWrite) (*VMInterface, error)
	UpdateVMInterface(ctx context.Context, token string, id int, in VMInterfaceWrite) (*VMInterface, error)
	DeleteVMInterface(ctx context.Context, token string, id int) error

	CreateProvider(ctx context.Context, token string, in ProviderWrite) (*Provider, error)
	UpdateProvider(ctx context.Context, token string, id int, in ProviderWrite) (*Provider, error)
	DeleteProvider(ctx context.Context, token string, id int) error

	CreateTenant(ctx context.Context, token string, in TenantWrite) (*Tenant, error)
	UpdateTenant(ctx context.Context, token string, id int, in TenantWrite) (*Tenant, error)
	DeleteTenant(ctx context.Context, token string, id int) error

	CreateManufacturer(ctx context.Context, token string, in ManufacturerWrite) (*Manufacturer, error)
	UpdateManufacturer(ctx context.Context, token string, id int, in ManufacturerWrite) (*Manufacturer, error)
	DeleteManufacturer(ctx context.Context, token string, id int) error

	CreateDeviceType(ctx context.Context, token string, in DeviceTypeWrite) (*DeviceType, error)
	UpdateDeviceType(ctx context.Context, token string, id int, in DeviceTypeWrite) (*DeviceType, error)
	DeleteDeviceType(ctx context.Context, token string, id int) error

	CreateLocation(ctx context.Context, token string, in LocationWrite) (*Location, error)
	UpdateLocation(ctx context.Context, token string, id int, in LocationWrite) (*Location, error)
	DeleteLocation(ctx context.Context, token string, id int) error

	CreateClusterType(ctx context.Context, token string, in ClusterTypeWrite) (*ClusterType, error)
	UpdateClusterType(ctx context.Context, token string, id int, in ClusterTypeWrite) (*ClusterType, error)
	DeleteClusterType(ctx context.Context, token string, id int) error

	CreateClusterGroup(ctx context.Context, token string, in ClusterGroupWrite) (*ClusterGroup, error)
	UpdateClusterGroup(ctx context.Context, token string, id int, in ClusterGroupWrite) (*ClusterGroup, error)
	DeleteClusterGroup(ctx context.Context, token string, id int) error

	CreateCircuitType(ctx context.Context, token string, in CircuitTypeWrite) (*CircuitType, error)
	UpdateCircuitType(ctx context.Context, token string, id int, in CircuitTypeWrite) (*CircuitType, error)
	DeleteCircuitType(ctx context.Context, token string, id int) error

	CreateVrf(ctx context.Context, token string, in VrfWrite) (*Vrf, error)
	UpdateVrf(ctx context.Context, token string, id int, in VrfWrite) (*Vrf, error)
	DeleteVrf(ctx context.Context, token string, id int) error

	CreateVlanGroup(ctx context.Context, token string, in VlanGroupWrite) (*VlanGroup, error)
	UpdateVlanGroup(ctx context.Context, token string, id int, in VlanGroupWrite) (*VlanGroup, error)
	DeleteVlanGroup(ctx context.Context, token string, id int) error

	CreateRole(ctx context.Context, token string, in RoleWrite) (*Role, error)
	UpdateRole(ctx context.Context, token string, id int, in RoleWrite) (*Role, error)
	DeleteRole(ctx context.Context, token string, id int) error

	CreateContact(ctx context.Context, token string, in ContactWrite) (*Contact, error)
	UpdateContact(ctx context.Context, token string, id int, in ContactWrite) (*Contact, error)
	DeleteContact(ctx context.Context, token string, id int) error

	ListSiteGroups(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[SiteGroup], error)
	CreateSiteGroup(ctx context.Context, token string, in SiteGroupWrite) (*SiteGroup, error)
	UpdateSiteGroup(ctx context.Context, token string, id int, in SiteGroupWrite) (*SiteGroup, error)
	DeleteSiteGroup(ctx context.Context, token string, id int) error

	ListRackGroups(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[RackGroup], error)
	CreateRackGroup(ctx context.Context, token string, in RackGroupWrite) (*RackGroup, error)
	UpdateRackGroup(ctx context.Context, token string, id int, in RackGroupWrite) (*RackGroup, error)
	DeleteRackGroup(ctx context.Context, token string, id int) error

	ListRackRoles(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[RackRole], error)
	CreateRackRole(ctx context.Context, token string, in RackRoleWrite) (*RackRole, error)
	UpdateRackRole(ctx context.Context, token string, id int, in RackRoleWrite) (*RackRole, error)
	DeleteRackRole(ctx context.Context, token string, id int) error

	ListDeviceRoles(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[DeviceRole], error)
	CreateDeviceRole(ctx context.Context, token string, in DeviceRoleWrite) (*DeviceRole, error)
	UpdateDeviceRole(ctx context.Context, token string, id int, in DeviceRoleWrite) (*DeviceRole, error)
	DeleteDeviceRole(ctx context.Context, token string, id int) error

	ListPlatforms(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[Platform], error)
	CreatePlatform(ctx context.Context, token string, in PlatformWrite) (*Platform, error)
	UpdatePlatform(ctx context.Context, token string, id int, in PlatformWrite) (*Platform, error)
	DeletePlatform(ctx context.Context, token string, id int) error

	ListVirtualChassis(
		ctx context.Context, token string, params map[string]string,
	) (*PaginatedResponse[VirtualChassis], error)
	CreateVirtualChassis(ctx context.Context, token string, in VirtualChassisWrite) (*VirtualChassis, error)
	UpdateVirtualChassis(ctx context.Context, token string, id int, in VirtualChassisWrite) (*VirtualChassis, error)
	DeleteVirtualChassis(ctx context.Context, token string, id int) error

	ListPowerPanels(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[PowerPanel], error)
	CreatePowerPanel(ctx context.Context, token string, in PowerPanelWrite) (*PowerPanel, error)
	UpdatePowerPanel(ctx context.Context, token string, id int, in PowerPanelWrite) (*PowerPanel, error)
	DeletePowerPanel(ctx context.Context, token string, id int) error

	ListInventoryItemRoles(
		ctx context.Context, token string, params map[string]string,
	) (*PaginatedResponse[InventoryItemRole], error)
	CreateInventoryItemRole(
		ctx context.Context, token string, in InventoryItemRoleWrite,
	) (*InventoryItemRole, error)
	UpdateInventoryItemRole(
		ctx context.Context, token string, id int, in InventoryItemRoleWrite,
	) (*InventoryItemRole, error)
	DeleteInventoryItemRole(ctx context.Context, token string, id int) error

	ListFrontPorts(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[FrontPort], error)
	CreateFrontPort(ctx context.Context, token string, in FrontPortWrite) (*FrontPort, error)
	UpdateFrontPort(ctx context.Context, token string, id int, in FrontPortWrite) (*FrontPort, error)
	DeleteFrontPort(ctx context.Context, token string, id int) error

	ListRearPorts(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[RearPort], error)
	CreateRearPort(ctx context.Context, token string, in RearPortWrite) (*RearPort, error)
	UpdateRearPort(ctx context.Context, token string, id int, in RearPortWrite) (*RearPort, error)
	DeleteRearPort(ctx context.Context, token string, id int) error

	ListConsolePorts(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[ConsolePort], error)
	CreateConsolePort(ctx context.Context, token string, in ConsolePortWrite) (*ConsolePort, error)
	UpdateConsolePort(ctx context.Context, token string, id int, in ConsolePortWrite) (*ConsolePort, error)
	DeleteConsolePort(ctx context.Context, token string, id int) error

	ListConsoleServerPorts(
		ctx context.Context, token string, params map[string]string,
	) (*PaginatedResponse[ConsoleServerPort], error)
	CreateConsoleServerPort(ctx context.Context, token string, in ConsoleServerPortWrite) (*ConsoleServerPort, error)
	UpdateConsoleServerPort(
		ctx context.Context, token string, id int, in ConsoleServerPortWrite,
	) (*ConsoleServerPort, error)
	DeleteConsoleServerPort(ctx context.Context, token string, id int) error

	ListPowerPorts(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[PowerPort], error)
	CreatePowerPort(ctx context.Context, token string, in PowerPortWrite) (*PowerPort, error)
	UpdatePowerPort(ctx context.Context, token string, id int, in PowerPortWrite) (*PowerPort, error)
	DeletePowerPort(ctx context.Context, token string, id int) error

	ListPowerOutlets(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[PowerOutlet], error)
	CreatePowerOutlet(ctx context.Context, token string, in PowerOutletWrite) (*PowerOutlet, error)
	UpdatePowerOutlet(ctx context.Context, token string, id int, in PowerOutletWrite) (*PowerOutlet, error)
	DeletePowerOutlet(ctx context.Context, token string, id int) error

	ListPowerFeeds(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[PowerFeed], error)
	CreatePowerFeed(ctx context.Context, token string, in PowerFeedWrite) (*PowerFeed, error)
	UpdatePowerFeed(ctx context.Context, token string, id int, in PowerFeedWrite) (*PowerFeed, error)
	DeletePowerFeed(ctx context.Context, token string, id int) error

	ListDeviceBays(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[DeviceBay], error)
	CreateDeviceBay(ctx context.Context, token string, in DeviceBayWrite) (*DeviceBay, error)
	UpdateDeviceBay(ctx context.Context, token string, id int, in DeviceBayWrite) (*DeviceBay, error)
	DeleteDeviceBay(ctx context.Context, token string, id int) error

	ListRackReservations(
		ctx context.Context, token string, params map[string]string,
	) (*PaginatedResponse[RackReservation], error)
	CreateRackReservation(ctx context.Context, token string, in RackReservationWrite) (*RackReservation, error)
	UpdateRackReservation(
		ctx context.Context, token string, id int, in RackReservationWrite,
	) (*RackReservation, error)
	DeleteRackReservation(ctx context.Context, token string, id int) error

	ListInventoryItems(
		ctx context.Context, token string, params map[string]string,
	) (*PaginatedResponse[InventoryItem], error)
	CreateInventoryItem(ctx context.Context, token string, in InventoryItemWrite) (*InventoryItem, error)
	UpdateInventoryItem(ctx context.Context, token string, id int, in InventoryItemWrite) (*InventoryItem, error)
	DeleteInventoryItem(ctx context.Context, token string, id int) error

	ListMACAddresses(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[MACAddress], error)
	CreateMACAddress(ctx context.Context, token string, in MACAddressWrite) (*MACAddress, error)
	UpdateMACAddress(ctx context.Context, token string, id int, in MACAddressWrite) (*MACAddress, error)
	DeleteMACAddress(ctx context.Context, token string, id int) error
}
