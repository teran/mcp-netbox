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

	ListSiteGroupsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.SiteGroup], error)
	CreateSiteGroupFunc func(
		ctx context.Context, token string, in domain.SiteGroupWrite,
	) (*domain.SiteGroup, error)
	UpdateSiteGroupFunc func(
		ctx context.Context, token string, id int, in domain.SiteGroupWrite,
	) (*domain.SiteGroup, error)
	DeleteSiteGroupFunc func(ctx context.Context, token string, id int) error

	ListRackGroupsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.RackGroup], error)
	CreateRackGroupFunc func(
		ctx context.Context, token string, in domain.RackGroupWrite,
	) (*domain.RackGroup, error)
	UpdateRackGroupFunc func(
		ctx context.Context, token string, id int, in domain.RackGroupWrite,
	) (*domain.RackGroup, error)
	DeleteRackGroupFunc func(ctx context.Context, token string, id int) error

	ListRackRolesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.RackRole], error)
	CreateRackRoleFunc func(
		ctx context.Context, token string, in domain.RackRoleWrite,
	) (*domain.RackRole, error)
	UpdateRackRoleFunc func(
		ctx context.Context, token string, id int, in domain.RackRoleWrite,
	) (*domain.RackRole, error)
	DeleteRackRoleFunc func(ctx context.Context, token string, id int) error

	ListDeviceRolesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.DeviceRole], error)
	CreateDeviceRoleFunc func(
		ctx context.Context, token string, in domain.DeviceRoleWrite,
	) (*domain.DeviceRole, error)
	UpdateDeviceRoleFunc func(
		ctx context.Context, token string, id int, in domain.DeviceRoleWrite,
	) (*domain.DeviceRole, error)
	DeleteDeviceRoleFunc func(ctx context.Context, token string, id int) error

	ListPlatformsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.Platform], error)
	CreatePlatformFunc func(
		ctx context.Context, token string, in domain.PlatformWrite,
	) (*domain.Platform, error)
	UpdatePlatformFunc func(
		ctx context.Context, token string, id int, in domain.PlatformWrite,
	) (*domain.Platform, error)
	DeletePlatformFunc func(ctx context.Context, token string, id int) error

	ListVirtualChassisFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.VirtualChassis], error)
	CreateVirtualChassisFunc func(
		ctx context.Context, token string, in domain.VirtualChassisWrite,
	) (*domain.VirtualChassis, error)
	UpdateVirtualChassisFunc func(
		ctx context.Context, token string, id int, in domain.VirtualChassisWrite,
	) (*domain.VirtualChassis, error)
	DeleteVirtualChassisFunc func(ctx context.Context, token string, id int) error

	ListPowerPanelsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.PowerPanel], error)
	CreatePowerPanelFunc func(
		ctx context.Context, token string, in domain.PowerPanelWrite,
	) (*domain.PowerPanel, error)
	UpdatePowerPanelFunc func(
		ctx context.Context, token string, id int, in domain.PowerPanelWrite,
	) (*domain.PowerPanel, error)
	DeletePowerPanelFunc func(ctx context.Context, token string, id int) error

	ListInventoryItemRolesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.InventoryItemRole], error)
	CreateInventoryItemRoleFunc func(
		ctx context.Context, token string, in domain.InventoryItemRoleWrite,
	) (*domain.InventoryItemRole, error)
	UpdateInventoryItemRoleFunc func(
		ctx context.Context, token string, id int, in domain.InventoryItemRoleWrite,
	) (*domain.InventoryItemRole, error)
	DeleteInventoryItemRoleFunc func(ctx context.Context, token string, id int) error

	ListFrontPortsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.FrontPort], error)
	CreateFrontPortFunc func(
		ctx context.Context, token string, in domain.FrontPortWrite,
	) (*domain.FrontPort, error)
	UpdateFrontPortFunc func(
		ctx context.Context, token string, id int, in domain.FrontPortWrite,
	) (*domain.FrontPort, error)
	DeleteFrontPortFunc func(ctx context.Context, token string, id int) error

	ListRearPortsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.RearPort], error)
	CreateRearPortFunc func(
		ctx context.Context, token string, in domain.RearPortWrite,
	) (*domain.RearPort, error)
	UpdateRearPortFunc func(
		ctx context.Context, token string, id int, in domain.RearPortWrite,
	) (*domain.RearPort, error)
	DeleteRearPortFunc func(ctx context.Context, token string, id int) error

	ListConsolePortsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.ConsolePort], error)
	CreateConsolePortFunc func(
		ctx context.Context, token string, in domain.ConsolePortWrite,
	) (*domain.ConsolePort, error)
	UpdateConsolePortFunc func(
		ctx context.Context, token string, id int, in domain.ConsolePortWrite,
	) (*domain.ConsolePort, error)
	DeleteConsolePortFunc func(ctx context.Context, token string, id int) error

	ListConsoleServerPortsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.ConsoleServerPort], error)
	CreateConsoleServerPortFunc func(
		ctx context.Context, token string, in domain.ConsoleServerPortWrite,
	) (*domain.ConsoleServerPort, error)
	UpdateConsoleServerPortFunc func(
		ctx context.Context, token string, id int, in domain.ConsoleServerPortWrite,
	) (*domain.ConsoleServerPort, error)
	DeleteConsoleServerPortFunc func(ctx context.Context, token string, id int) error

	ListPowerPortsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.PowerPort], error)
	CreatePowerPortFunc func(
		ctx context.Context, token string, in domain.PowerPortWrite,
	) (*domain.PowerPort, error)
	UpdatePowerPortFunc func(
		ctx context.Context, token string, id int, in domain.PowerPortWrite,
	) (*domain.PowerPort, error)
	DeletePowerPortFunc func(ctx context.Context, token string, id int) error

	ListPowerOutletsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.PowerOutlet], error)
	CreatePowerOutletFunc func(
		ctx context.Context, token string, in domain.PowerOutletWrite,
	) (*domain.PowerOutlet, error)
	UpdatePowerOutletFunc func(
		ctx context.Context, token string, id int, in domain.PowerOutletWrite,
	) (*domain.PowerOutlet, error)
	DeletePowerOutletFunc func(ctx context.Context, token string, id int) error

	ListPowerFeedsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.PowerFeed], error)
	CreatePowerFeedFunc func(
		ctx context.Context, token string, in domain.PowerFeedWrite,
	) (*domain.PowerFeed, error)
	UpdatePowerFeedFunc func(
		ctx context.Context, token string, id int, in domain.PowerFeedWrite,
	) (*domain.PowerFeed, error)
	DeletePowerFeedFunc func(ctx context.Context, token string, id int) error

	ListDeviceBaysFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.DeviceBay], error)
	CreateDeviceBayFunc func(
		ctx context.Context, token string, in domain.DeviceBayWrite,
	) (*domain.DeviceBay, error)
	UpdateDeviceBayFunc func(
		ctx context.Context, token string, id int, in domain.DeviceBayWrite,
	) (*domain.DeviceBay, error)
	DeleteDeviceBayFunc func(ctx context.Context, token string, id int) error

	ListRackReservationsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.RackReservation], error)
	CreateRackReservationFunc func(
		ctx context.Context, token string, in domain.RackReservationWrite,
	) (*domain.RackReservation, error)
	UpdateRackReservationFunc func(
		ctx context.Context, token string, id int, in domain.RackReservationWrite,
	) (*domain.RackReservation, error)
	DeleteRackReservationFunc func(ctx context.Context, token string, id int) error

	ListInventoryItemsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.InventoryItem], error)
	CreateInventoryItemFunc func(
		ctx context.Context, token string, in domain.InventoryItemWrite,
	) (*domain.InventoryItem, error)
	UpdateInventoryItemFunc func(
		ctx context.Context, token string, id int, in domain.InventoryItemWrite,
	) (*domain.InventoryItem, error)
	DeleteInventoryItemFunc func(ctx context.Context, token string, id int) error

	ListMACAddressesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.MACAddress], error)
	CreateMACAddressFunc func(
		ctx context.Context, token string, in domain.MACAddressWrite,
	) (*domain.MACAddress, error)
	UpdateMACAddressFunc func(
		ctx context.Context, token string, id int, in domain.MACAddressWrite,
	) (*domain.MACAddress, error)
	DeleteMACAddressFunc func(ctx context.Context, token string, id int) error

	ListRIRsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.RIR], error)
	CreateRIRFunc func(
		ctx context.Context, token string, in domain.RIRWrite,
	) (*domain.RIR, error)
	UpdateRIRFunc func(
		ctx context.Context, token string, id int, in domain.RIRWrite,
	) (*domain.RIR, error)
	DeleteRIRFunc func(ctx context.Context, token string, id int) error

	ListAggregatesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.Aggregate], error)
	CreateAggregateFunc func(
		ctx context.Context, token string, in domain.AggregateWrite,
	) (*domain.Aggregate, error)
	UpdateAggregateFunc func(
		ctx context.Context, token string, id int, in domain.AggregateWrite,
	) (*domain.Aggregate, error)
	DeleteAggregateFunc func(ctx context.Context, token string, id int) error

	ListIPRangesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.IPRange], error)
	CreateIPRangeFunc func(
		ctx context.Context, token string, in domain.IPRangeWrite,
	) (*domain.IPRange, error)
	UpdateIPRangeFunc func(
		ctx context.Context, token string, id int, in domain.IPRangeWrite,
	) (*domain.IPRange, error)
	DeleteIPRangeFunc func(ctx context.Context, token string, id int) error

	ListASNsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.ASN], error)
	CreateASNFunc func(
		ctx context.Context, token string, in domain.ASNWrite,
	) (*domain.ASN, error)
	UpdateASNFunc func(
		ctx context.Context, token string, id int, in domain.ASNWrite,
	) (*domain.ASN, error)
	DeleteASNFunc func(ctx context.Context, token string, id int) error

	ListASNRangesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.ASNRange], error)
	CreateASNRangeFunc func(
		ctx context.Context, token string, in domain.ASNRangeWrite,
	) (*domain.ASNRange, error)
	UpdateASNRangeFunc func(
		ctx context.Context, token string, id int, in domain.ASNRangeWrite,
	) (*domain.ASNRange, error)
	DeleteASNRangeFunc func(ctx context.Context, token string, id int) error

	ListRouteTargetsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.RouteTarget], error)
	CreateRouteTargetFunc func(
		ctx context.Context, token string, in domain.RouteTargetWrite,
	) (*domain.RouteTarget, error)
	UpdateRouteTargetFunc func(
		ctx context.Context, token string, id int, in domain.RouteTargetWrite,
	) (*domain.RouteTarget, error)
	DeleteRouteTargetFunc func(ctx context.Context, token string, id int) error

	ListServicesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.Service], error)
	CreateServiceFunc func(
		ctx context.Context, token string, in domain.ServiceWrite,
	) (*domain.Service, error)
	UpdateServiceFunc func(
		ctx context.Context, token string, id int, in domain.ServiceWrite,
	) (*domain.Service, error)
	DeleteServiceFunc func(ctx context.Context, token string, id int) error

	ListServiceTemplatesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.ServiceTemplate], error)
	CreateServiceTemplateFunc func(
		ctx context.Context, token string, in domain.ServiceTemplateWrite,
	) (*domain.ServiceTemplate, error)
	UpdateServiceTemplateFunc func(
		ctx context.Context, token string, id int, in domain.ServiceTemplateWrite,
	) (*domain.ServiceTemplate, error)
	DeleteServiceTemplateFunc func(ctx context.Context, token string, id int) error

	ListFHRPGroupsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.FHRPGroup], error)
	CreateFHRPGroupFunc func(
		ctx context.Context, token string, in domain.FHRPGroupWrite,
	) (*domain.FHRPGroup, error)
	UpdateFHRPGroupFunc func(
		ctx context.Context, token string, id int, in domain.FHRPGroupWrite,
	) (*domain.FHRPGroup, error)
	DeleteFHRPGroupFunc func(ctx context.Context, token string, id int) error

	ListFHRPGroupAssignmentsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.FHRPGroupAssignment], error)
	CreateFHRPGroupAssignmentFunc func(
		ctx context.Context, token string, in domain.FHRPGroupAssignmentWrite,
	) (*domain.FHRPGroupAssignment, error)
	UpdateFHRPGroupAssignmentFunc func(
		ctx context.Context, token string, id int, in domain.FHRPGroupAssignmentWrite,
	) (*domain.FHRPGroupAssignment, error)
	DeleteFHRPGroupAssignmentFunc func(ctx context.Context, token string, id int) error

	ListVLANTranslationPoliciesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.VLANTranslationPolicy], error)
	CreateVLANTranslationPolicyFunc func(
		ctx context.Context, token string, in domain.VLANTranslationPolicyWrite,
	) (*domain.VLANTranslationPolicy, error)
	UpdateVLANTranslationPolicyFunc func(
		ctx context.Context, token string, id int, in domain.VLANTranslationPolicyWrite,
	) (*domain.VLANTranslationPolicy, error)
	DeleteVLANTranslationPolicyFunc func(ctx context.Context, token string, id int) error

	ListVLANTranslationRulesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.VLANTranslationRule], error)
	CreateVLANTranslationRuleFunc func(
		ctx context.Context, token string, in domain.VLANTranslationRuleWrite,
	) (*domain.VLANTranslationRule, error)
	UpdateVLANTranslationRuleFunc func(
		ctx context.Context, token string, id int, in domain.VLANTranslationRuleWrite,
	) (*domain.VLANTranslationRule, error)
	DeleteVLANTranslationRuleFunc func(ctx context.Context, token string, id int) error

	ListContactGroupsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.ContactGroup], error)
	CreateContactGroupFunc func(
		ctx context.Context, token string, in domain.ContactGroupWrite,
	) (*domain.ContactGroup, error)
	UpdateContactGroupFunc func(
		ctx context.Context, token string, id int, in domain.ContactGroupWrite,
	) (*domain.ContactGroup, error)
	DeleteContactGroupFunc func(ctx context.Context, token string, id int) error

	ListContactRolesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.ContactRole], error)
	CreateContactRoleFunc func(
		ctx context.Context, token string, in domain.ContactRoleWrite,
	) (*domain.ContactRole, error)
	UpdateContactRoleFunc func(
		ctx context.Context, token string, id int, in domain.ContactRoleWrite,
	) (*domain.ContactRole, error)
	DeleteContactRoleFunc func(ctx context.Context, token string, id int) error

	ListVirtualDisksFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.VirtualDisk], error)
	CreateVirtualDiskFunc func(
		ctx context.Context, token string, in domain.VirtualDiskWrite,
	) (*domain.VirtualDisk, error)
	UpdateVirtualDiskFunc func(
		ctx context.Context, token string, id int, in domain.VirtualDiskWrite,
	) (*domain.VirtualDisk, error)
	DeleteVirtualDiskFunc func(ctx context.Context, token string, id int) error

	ListVirtualMachineTypesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.VirtualMachineType], error)
	CreateVirtualMachineTypeFunc func(
		ctx context.Context, token string, in domain.VirtualMachineTypeWrite,
	) (*domain.VirtualMachineType, error)
	UpdateVirtualMachineTypeFunc func(
		ctx context.Context, token string, id int, in domain.VirtualMachineTypeWrite,
	) (*domain.VirtualMachineType, error)
	DeleteVirtualMachineTypeFunc func(ctx context.Context, token string, id int) error
	ListProviderAccountsFunc     func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.ProviderAccount], error)
	CreateProviderAccountFunc func(
		ctx context.Context, token string, in domain.ProviderAccountWrite,
	) (*domain.ProviderAccount, error)
	UpdateProviderAccountFunc func(
		ctx context.Context, token string, id int, in domain.ProviderAccountWrite,
	) (*domain.ProviderAccount, error)
	DeleteProviderAccountFunc func(ctx context.Context, token string, id int) error
	ListProviderNetworksFunc  func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.ProviderNetwork], error)
	CreateProviderNetworkFunc func(
		ctx context.Context, token string, in domain.ProviderNetworkWrite,
	) (*domain.ProviderNetwork, error)
	UpdateProviderNetworkFunc func(
		ctx context.Context, token string, id int, in domain.ProviderNetworkWrite,
	) (*domain.ProviderNetwork, error)
	DeleteProviderNetworkFunc func(ctx context.Context, token string, id int) error
	ListCircuitGroupsFunc     func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.CircuitGroup], error)
	CreateCircuitGroupFunc func(
		ctx context.Context, token string, in domain.CircuitGroupWrite,
	) (*domain.CircuitGroup, error)
	UpdateCircuitGroupFunc func(
		ctx context.Context, token string, id int, in domain.CircuitGroupWrite,
	) (*domain.CircuitGroup, error)
	DeleteCircuitGroupFunc          func(ctx context.Context, token string, id int) error
	ListCircuitGroupAssignmentsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.CircuitGroupAssignment], error)
	CreateCircuitGroupAssignmentFunc func(
		ctx context.Context, token string, in domain.CircuitGroupAssignmentWrite,
	) (*domain.CircuitGroupAssignment, error)
	UpdateCircuitGroupAssignmentFunc func(
		ctx context.Context, token string, id int, in domain.CircuitGroupAssignmentWrite,
	) (*domain.CircuitGroupAssignment, error)
	DeleteCircuitGroupAssignmentFunc func(ctx context.Context, token string, id int) error
	ListVirtualCircuitsFunc          func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.VirtualCircuit], error)
	CreateVirtualCircuitFunc func(
		ctx context.Context, token string, in domain.VirtualCircuitWrite,
	) (*domain.VirtualCircuit, error)
	UpdateVirtualCircuitFunc func(
		ctx context.Context, token string, id int, in domain.VirtualCircuitWrite,
	) (*domain.VirtualCircuit, error)
	DeleteVirtualCircuitFunc    func(ctx context.Context, token string, id int) error
	ListVirtualCircuitTypesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.VirtualCircuitType], error)
	CreateVirtualCircuitTypeFunc func(
		ctx context.Context, token string, in domain.VirtualCircuitTypeWrite,
	) (*domain.VirtualCircuitType, error)
	UpdateVirtualCircuitTypeFunc func(
		ctx context.Context, token string, id int, in domain.VirtualCircuitTypeWrite,
	) (*domain.VirtualCircuitType, error)
	DeleteVirtualCircuitTypeFunc       func(ctx context.Context, token string, id int) error
	ListVirtualCircuitTerminationsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.VirtualCircuitTermination], error)
	CreateVirtualCircuitTerminationFunc func(
		ctx context.Context, token string, in domain.VirtualCircuitTerminationWrite,
	) (*domain.VirtualCircuitTermination, error)
	UpdateVirtualCircuitTerminationFunc func(
		ctx context.Context, token string, id int, in domain.VirtualCircuitTerminationWrite,
	) (*domain.VirtualCircuitTermination, error)
	DeleteVirtualCircuitTerminationFunc func(ctx context.Context, token string, id int) error

	ListRolesFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.Role], error)
	ListTenantsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.Tenant], error)

	ListTagsFunc func(
		ctx context.Context, token string, params map[string]string,
	) (*domain.PaginatedResponse[domain.Tag], error)
	CreateTagFunc func(ctx context.Context, token string, in domain.TagWrite) (*domain.Tag, error)
	UpdateTagFunc func(ctx context.Context, token string, id int, in domain.TagWrite) (*domain.Tag, error)
	DeleteTagFunc func(ctx context.Context, token string, id int) error
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

// ListSiteGroups delegates to ListSiteGroupsFunc.
func (m *MockRepo) ListSiteGroups(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.SiteGroup], error) {
	return m.ListSiteGroupsFunc(ctx, token, params)
}

// CreateSiteGroup delegates to CreateSiteGroupFunc.
func (m *MockRepo) CreateSiteGroup(
	ctx context.Context, token string, in domain.SiteGroupWrite,
) (*domain.SiteGroup, error) {
	return m.CreateSiteGroupFunc(ctx, token, in)
}

// UpdateSiteGroup delegates to UpdateSiteGroupFunc.
func (m *MockRepo) UpdateSiteGroup(
	ctx context.Context, token string, id int, in domain.SiteGroupWrite,
) (*domain.SiteGroup, error) {
	return m.UpdateSiteGroupFunc(ctx, token, id, in)
}

// DeleteSiteGroup delegates to DeleteSiteGroupFunc.
func (m *MockRepo) DeleteSiteGroup(ctx context.Context, token string, id int) error {
	return m.DeleteSiteGroupFunc(ctx, token, id)
}

// ListRackGroups delegates to ListRackGroupsFunc.
func (m *MockRepo) ListRackGroups(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.RackGroup], error) {
	return m.ListRackGroupsFunc(ctx, token, params)
}

// CreateRackGroup delegates to CreateRackGroupFunc.
func (m *MockRepo) CreateRackGroup(
	ctx context.Context, token string, in domain.RackGroupWrite,
) (*domain.RackGroup, error) {
	return m.CreateRackGroupFunc(ctx, token, in)
}

// UpdateRackGroup delegates to UpdateRackGroupFunc.
func (m *MockRepo) UpdateRackGroup(
	ctx context.Context, token string, id int, in domain.RackGroupWrite,
) (*domain.RackGroup, error) {
	return m.UpdateRackGroupFunc(ctx, token, id, in)
}

// DeleteRackGroup delegates to DeleteRackGroupFunc.
func (m *MockRepo) DeleteRackGroup(ctx context.Context, token string, id int) error {
	return m.DeleteRackGroupFunc(ctx, token, id)
}

// ListRackRoles delegates to ListRackRolesFunc.
func (m *MockRepo) ListRackRoles(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.RackRole], error) {
	return m.ListRackRolesFunc(ctx, token, params)
}

// CreateRackRole delegates to CreateRackRoleFunc.
func (m *MockRepo) CreateRackRole(
	ctx context.Context, token string, in domain.RackRoleWrite,
) (*domain.RackRole, error) {
	return m.CreateRackRoleFunc(ctx, token, in)
}

// UpdateRackRole delegates to UpdateRackRoleFunc.
func (m *MockRepo) UpdateRackRole(
	ctx context.Context, token string, id int, in domain.RackRoleWrite,
) (*domain.RackRole, error) {
	return m.UpdateRackRoleFunc(ctx, token, id, in)
}

// DeleteRackRole delegates to DeleteRackRoleFunc.
func (m *MockRepo) DeleteRackRole(ctx context.Context, token string, id int) error {
	return m.DeleteRackRoleFunc(ctx, token, id)
}

// ListDeviceRoles delegates to ListDeviceRolesFunc.
func (m *MockRepo) ListDeviceRoles(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.DeviceRole], error) {
	return m.ListDeviceRolesFunc(ctx, token, params)
}

// CreateDeviceRole delegates to CreateDeviceRoleFunc.
func (m *MockRepo) CreateDeviceRole(
	ctx context.Context, token string, in domain.DeviceRoleWrite,
) (*domain.DeviceRole, error) {
	return m.CreateDeviceRoleFunc(ctx, token, in)
}

// UpdateDeviceRole delegates to UpdateDeviceRoleFunc.
func (m *MockRepo) UpdateDeviceRole(
	ctx context.Context, token string, id int, in domain.DeviceRoleWrite,
) (*domain.DeviceRole, error) {
	return m.UpdateDeviceRoleFunc(ctx, token, id, in)
}

// DeleteDeviceRole delegates to DeleteDeviceRoleFunc.
func (m *MockRepo) DeleteDeviceRole(ctx context.Context, token string, id int) error {
	return m.DeleteDeviceRoleFunc(ctx, token, id)
}

// ListPlatforms delegates to ListPlatformsFunc.
func (m *MockRepo) ListPlatforms(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.Platform], error) {
	return m.ListPlatformsFunc(ctx, token, params)
}

// CreatePlatform delegates to CreatePlatformFunc.
func (m *MockRepo) CreatePlatform(
	ctx context.Context, token string, in domain.PlatformWrite,
) (*domain.Platform, error) {
	return m.CreatePlatformFunc(ctx, token, in)
}

// UpdatePlatform delegates to UpdatePlatformFunc.
func (m *MockRepo) UpdatePlatform(
	ctx context.Context, token string, id int, in domain.PlatformWrite,
) (*domain.Platform, error) {
	return m.UpdatePlatformFunc(ctx, token, id, in)
}

// DeletePlatform delegates to DeletePlatformFunc.
func (m *MockRepo) DeletePlatform(ctx context.Context, token string, id int) error {
	return m.DeletePlatformFunc(ctx, token, id)
}

// ListVirtualChassis delegates to ListVirtualChassisFunc.
func (m *MockRepo) ListVirtualChassis(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.VirtualChassis], error) {
	return m.ListVirtualChassisFunc(ctx, token, params)
}

// CreateVirtualChassis delegates to CreateVirtualChassisFunc.
func (m *MockRepo) CreateVirtualChassis(
	ctx context.Context, token string, in domain.VirtualChassisWrite,
) (*domain.VirtualChassis, error) {
	return m.CreateVirtualChassisFunc(ctx, token, in)
}

// UpdateVirtualChassis delegates to UpdateVirtualChassisFunc.
func (m *MockRepo) UpdateVirtualChassis(
	ctx context.Context, token string, id int, in domain.VirtualChassisWrite,
) (*domain.VirtualChassis, error) {
	return m.UpdateVirtualChassisFunc(ctx, token, id, in)
}

// DeleteVirtualChassis delegates to DeleteVirtualChassisFunc.
func (m *MockRepo) DeleteVirtualChassis(ctx context.Context, token string, id int) error {
	return m.DeleteVirtualChassisFunc(ctx, token, id)
}

// ListPowerPanels delegates to ListPowerPanelsFunc.
func (m *MockRepo) ListPowerPanels(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.PowerPanel], error) {
	return m.ListPowerPanelsFunc(ctx, token, params)
}

// CreatePowerPanel delegates to CreatePowerPanelFunc.
func (m *MockRepo) CreatePowerPanel(
	ctx context.Context, token string, in domain.PowerPanelWrite,
) (*domain.PowerPanel, error) {
	return m.CreatePowerPanelFunc(ctx, token, in)
}

// UpdatePowerPanel delegates to UpdatePowerPanelFunc.
func (m *MockRepo) UpdatePowerPanel(
	ctx context.Context, token string, id int, in domain.PowerPanelWrite,
) (*domain.PowerPanel, error) {
	return m.UpdatePowerPanelFunc(ctx, token, id, in)
}

// DeletePowerPanel delegates to DeletePowerPanelFunc.
func (m *MockRepo) DeletePowerPanel(ctx context.Context, token string, id int) error {
	return m.DeletePowerPanelFunc(ctx, token, id)
}

// ListInventoryItemRoles delegates to ListInventoryItemRolesFunc.
func (m *MockRepo) ListInventoryItemRoles(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.InventoryItemRole], error) {
	return m.ListInventoryItemRolesFunc(ctx, token, params)
}

// CreateInventoryItemRole delegates to CreateInventoryItemRoleFunc.
func (m *MockRepo) CreateInventoryItemRole(
	ctx context.Context, token string, in domain.InventoryItemRoleWrite,
) (*domain.InventoryItemRole, error) {
	return m.CreateInventoryItemRoleFunc(ctx, token, in)
}

// UpdateInventoryItemRole delegates to UpdateInventoryItemRoleFunc.
func (m *MockRepo) UpdateInventoryItemRole(
	ctx context.Context, token string, id int, in domain.InventoryItemRoleWrite,
) (*domain.InventoryItemRole, error) {
	return m.UpdateInventoryItemRoleFunc(ctx, token, id, in)
}

// DeleteInventoryItemRole delegates to DeleteInventoryItemRoleFunc.
func (m *MockRepo) DeleteInventoryItemRole(ctx context.Context, token string, id int) error {
	return m.DeleteInventoryItemRoleFunc(ctx, token, id)
}

// ListFrontPorts delegates to ListFrontPortsFunc.
func (m *MockRepo) ListFrontPorts(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.FrontPort], error) {
	return m.ListFrontPortsFunc(ctx, token, params)
}

// CreateFrontPort delegates to CreateFrontPortFunc.
func (m *MockRepo) CreateFrontPort(
	ctx context.Context, token string, in domain.FrontPortWrite,
) (*domain.FrontPort, error) {
	return m.CreateFrontPortFunc(ctx, token, in)
}

// UpdateFrontPort delegates to UpdateFrontPortFunc.
func (m *MockRepo) UpdateFrontPort(
	ctx context.Context, token string, id int, in domain.FrontPortWrite,
) (*domain.FrontPort, error) {
	return m.UpdateFrontPortFunc(ctx, token, id, in)
}

// DeleteFrontPort delegates to DeleteFrontPortFunc.
func (m *MockRepo) DeleteFrontPort(ctx context.Context, token string, id int) error {
	return m.DeleteFrontPortFunc(ctx, token, id)
}

// ListRearPorts delegates to ListRearPortsFunc.
func (m *MockRepo) ListRearPorts(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.RearPort], error) {
	return m.ListRearPortsFunc(ctx, token, params)
}

// CreateRearPort delegates to CreateRearPortFunc.
func (m *MockRepo) CreateRearPort(
	ctx context.Context, token string, in domain.RearPortWrite,
) (*domain.RearPort, error) {
	return m.CreateRearPortFunc(ctx, token, in)
}

// UpdateRearPort delegates to UpdateRearPortFunc.
func (m *MockRepo) UpdateRearPort(
	ctx context.Context, token string, id int, in domain.RearPortWrite,
) (*domain.RearPort, error) {
	return m.UpdateRearPortFunc(ctx, token, id, in)
}

// DeleteRearPort delegates to DeleteRearPortFunc.
func (m *MockRepo) DeleteRearPort(ctx context.Context, token string, id int) error {
	return m.DeleteRearPortFunc(ctx, token, id)
}

// ListConsolePorts delegates to ListConsolePortsFunc.
func (m *MockRepo) ListConsolePorts(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.ConsolePort], error) {
	return m.ListConsolePortsFunc(ctx, token, params)
}

// CreateConsolePort delegates to CreateConsolePortFunc.
func (m *MockRepo) CreateConsolePort(
	ctx context.Context, token string, in domain.ConsolePortWrite,
) (*domain.ConsolePort, error) {
	return m.CreateConsolePortFunc(ctx, token, in)
}

// UpdateConsolePort delegates to UpdateConsolePortFunc.
func (m *MockRepo) UpdateConsolePort(
	ctx context.Context, token string, id int, in domain.ConsolePortWrite,
) (*domain.ConsolePort, error) {
	return m.UpdateConsolePortFunc(ctx, token, id, in)
}

// DeleteConsolePort delegates to DeleteConsolePortFunc.
func (m *MockRepo) DeleteConsolePort(ctx context.Context, token string, id int) error {
	return m.DeleteConsolePortFunc(ctx, token, id)
}

// ListConsoleServerPorts delegates to ListConsoleServerPortsFunc.
func (m *MockRepo) ListConsoleServerPorts(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.ConsoleServerPort], error) {
	return m.ListConsoleServerPortsFunc(ctx, token, params)
}

// CreateConsoleServerPort delegates to CreateConsoleServerPortFunc.
func (m *MockRepo) CreateConsoleServerPort(
	ctx context.Context, token string, in domain.ConsoleServerPortWrite,
) (*domain.ConsoleServerPort, error) {
	return m.CreateConsoleServerPortFunc(ctx, token, in)
}

// UpdateConsoleServerPort delegates to UpdateConsoleServerPortFunc.
func (m *MockRepo) UpdateConsoleServerPort(
	ctx context.Context, token string, id int, in domain.ConsoleServerPortWrite,
) (*domain.ConsoleServerPort, error) {
	return m.UpdateConsoleServerPortFunc(ctx, token, id, in)
}

// DeleteConsoleServerPort delegates to DeleteConsoleServerPortFunc.
func (m *MockRepo) DeleteConsoleServerPort(ctx context.Context, token string, id int) error {
	return m.DeleteConsoleServerPortFunc(ctx, token, id)
}

// ListPowerPorts delegates to ListPowerPortsFunc.
func (m *MockRepo) ListPowerPorts(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.PowerPort], error) {
	return m.ListPowerPortsFunc(ctx, token, params)
}

// CreatePowerPort delegates to CreatePowerPortFunc.
func (m *MockRepo) CreatePowerPort(
	ctx context.Context, token string, in domain.PowerPortWrite,
) (*domain.PowerPort, error) {
	return m.CreatePowerPortFunc(ctx, token, in)
}

// UpdatePowerPort delegates to UpdatePowerPortFunc.
func (m *MockRepo) UpdatePowerPort(
	ctx context.Context, token string, id int, in domain.PowerPortWrite,
) (*domain.PowerPort, error) {
	return m.UpdatePowerPortFunc(ctx, token, id, in)
}

// DeletePowerPort delegates to DeletePowerPortFunc.
func (m *MockRepo) DeletePowerPort(ctx context.Context, token string, id int) error {
	return m.DeletePowerPortFunc(ctx, token, id)
}

// ListPowerOutlets delegates to ListPowerOutletsFunc.
func (m *MockRepo) ListPowerOutlets(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.PowerOutlet], error) {
	return m.ListPowerOutletsFunc(ctx, token, params)
}

// CreatePowerOutlet delegates to CreatePowerOutletFunc.
func (m *MockRepo) CreatePowerOutlet(
	ctx context.Context, token string, in domain.PowerOutletWrite,
) (*domain.PowerOutlet, error) {
	return m.CreatePowerOutletFunc(ctx, token, in)
}

// UpdatePowerOutlet delegates to UpdatePowerOutletFunc.
func (m *MockRepo) UpdatePowerOutlet(
	ctx context.Context, token string, id int, in domain.PowerOutletWrite,
) (*domain.PowerOutlet, error) {
	return m.UpdatePowerOutletFunc(ctx, token, id, in)
}

// DeletePowerOutlet delegates to DeletePowerOutletFunc.
func (m *MockRepo) DeletePowerOutlet(ctx context.Context, token string, id int) error {
	return m.DeletePowerOutletFunc(ctx, token, id)
}

// ListPowerFeeds delegates to ListPowerFeedsFunc.
func (m *MockRepo) ListPowerFeeds(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.PowerFeed], error) {
	return m.ListPowerFeedsFunc(ctx, token, params)
}

// CreatePowerFeed delegates to CreatePowerFeedFunc.
func (m *MockRepo) CreatePowerFeed(
	ctx context.Context, token string, in domain.PowerFeedWrite,
) (*domain.PowerFeed, error) {
	return m.CreatePowerFeedFunc(ctx, token, in)
}

// UpdatePowerFeed delegates to UpdatePowerFeedFunc.
func (m *MockRepo) UpdatePowerFeed(
	ctx context.Context, token string, id int, in domain.PowerFeedWrite,
) (*domain.PowerFeed, error) {
	return m.UpdatePowerFeedFunc(ctx, token, id, in)
}

// DeletePowerFeed delegates to DeletePowerFeedFunc.
func (m *MockRepo) DeletePowerFeed(ctx context.Context, token string, id int) error {
	return m.DeletePowerFeedFunc(ctx, token, id)
}

// ListDeviceBays delegates to ListDeviceBaysFunc.
func (m *MockRepo) ListDeviceBays(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.DeviceBay], error) {
	return m.ListDeviceBaysFunc(ctx, token, params)
}

// CreateDeviceBay delegates to CreateDeviceBayFunc.
func (m *MockRepo) CreateDeviceBay(
	ctx context.Context, token string, in domain.DeviceBayWrite,
) (*domain.DeviceBay, error) {
	return m.CreateDeviceBayFunc(ctx, token, in)
}

// UpdateDeviceBay delegates to UpdateDeviceBayFunc.
func (m *MockRepo) UpdateDeviceBay(
	ctx context.Context, token string, id int, in domain.DeviceBayWrite,
) (*domain.DeviceBay, error) {
	return m.UpdateDeviceBayFunc(ctx, token, id, in)
}

// DeleteDeviceBay delegates to DeleteDeviceBayFunc.
func (m *MockRepo) DeleteDeviceBay(ctx context.Context, token string, id int) error {
	return m.DeleteDeviceBayFunc(ctx, token, id)
}

// ListRackReservations delegates to ListRackReservationsFunc.
func (m *MockRepo) ListRackReservations(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.RackReservation], error) {
	return m.ListRackReservationsFunc(ctx, token, params)
}

// CreateRackReservation delegates to CreateRackReservationFunc.
func (m *MockRepo) CreateRackReservation(
	ctx context.Context, token string, in domain.RackReservationWrite,
) (*domain.RackReservation, error) {
	return m.CreateRackReservationFunc(ctx, token, in)
}

// UpdateRackReservation delegates to UpdateRackReservationFunc.
func (m *MockRepo) UpdateRackReservation(
	ctx context.Context, token string, id int, in domain.RackReservationWrite,
) (*domain.RackReservation, error) {
	return m.UpdateRackReservationFunc(ctx, token, id, in)
}

// DeleteRackReservation delegates to DeleteRackReservationFunc.
func (m *MockRepo) DeleteRackReservation(ctx context.Context, token string, id int) error {
	return m.DeleteRackReservationFunc(ctx, token, id)
}

// ListInventoryItems delegates to ListInventoryItemsFunc.
func (m *MockRepo) ListInventoryItems(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.InventoryItem], error) {
	return m.ListInventoryItemsFunc(ctx, token, params)
}

// CreateInventoryItem delegates to CreateInventoryItemFunc.
func (m *MockRepo) CreateInventoryItem(
	ctx context.Context, token string, in domain.InventoryItemWrite,
) (*domain.InventoryItem, error) {
	return m.CreateInventoryItemFunc(ctx, token, in)
}

// UpdateInventoryItem delegates to UpdateInventoryItemFunc.
func (m *MockRepo) UpdateInventoryItem(
	ctx context.Context, token string, id int, in domain.InventoryItemWrite,
) (*domain.InventoryItem, error) {
	return m.UpdateInventoryItemFunc(ctx, token, id, in)
}

// DeleteInventoryItem delegates to DeleteInventoryItemFunc.
func (m *MockRepo) DeleteInventoryItem(ctx context.Context, token string, id int) error {
	return m.DeleteInventoryItemFunc(ctx, token, id)
}

// ListMACAddresses delegates to ListMACAddressesFunc.
func (m *MockRepo) ListMACAddresses(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.MACAddress], error) {
	return m.ListMACAddressesFunc(ctx, token, params)
}

// CreateMACAddress delegates to CreateMACAddressFunc.
func (m *MockRepo) CreateMACAddress(
	ctx context.Context, token string, in domain.MACAddressWrite,
) (*domain.MACAddress, error) {
	return m.CreateMACAddressFunc(ctx, token, in)
}

// UpdateMACAddress delegates to UpdateMACAddressFunc.
func (m *MockRepo) UpdateMACAddress(
	ctx context.Context, token string, id int, in domain.MACAddressWrite,
) (*domain.MACAddress, error) {
	return m.UpdateMACAddressFunc(ctx, token, id, in)
}

// DeleteMACAddress delegates to DeleteMACAddressFunc.
func (m *MockRepo) DeleteMACAddress(ctx context.Context, token string, id int) error {
	return m.DeleteMACAddressFunc(ctx, token, id)
}

// ListRIRs delegates to ListRIRsFunc.
func (m *MockRepo) ListRIRs(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.RIR], error) {
	return m.ListRIRsFunc(ctx, token, params)
}

// CreateRIR delegates to CreateRIRFunc.
func (m *MockRepo) CreateRIR(
	ctx context.Context, token string, in domain.RIRWrite,
) (*domain.RIR, error) {
	return m.CreateRIRFunc(ctx, token, in)
}

// UpdateRIR delegates to UpdateRIRFunc.
func (m *MockRepo) UpdateRIR(
	ctx context.Context, token string, id int, in domain.RIRWrite,
) (*domain.RIR, error) {
	return m.UpdateRIRFunc(ctx, token, id, in)
}

// DeleteRIR delegates to DeleteRIRFunc.
func (m *MockRepo) DeleteRIR(ctx context.Context, token string, id int) error {
	return m.DeleteRIRFunc(ctx, token, id)
}

// ListAggregates delegates to ListAggregatesFunc.
func (m *MockRepo) ListAggregates(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.Aggregate], error) {
	return m.ListAggregatesFunc(ctx, token, params)
}

// CreateAggregate delegates to CreateAggregateFunc.
func (m *MockRepo) CreateAggregate(
	ctx context.Context, token string, in domain.AggregateWrite,
) (*domain.Aggregate, error) {
	return m.CreateAggregateFunc(ctx, token, in)
}

// UpdateAggregate delegates to UpdateAggregateFunc.
func (m *MockRepo) UpdateAggregate(
	ctx context.Context, token string, id int, in domain.AggregateWrite,
) (*domain.Aggregate, error) {
	return m.UpdateAggregateFunc(ctx, token, id, in)
}

// DeleteAggregate delegates to DeleteAggregateFunc.
func (m *MockRepo) DeleteAggregate(ctx context.Context, token string, id int) error {
	return m.DeleteAggregateFunc(ctx, token, id)
}

// ListIPRanges delegates to ListIPRangesFunc.
func (m *MockRepo) ListIPRanges(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.IPRange], error) {
	return m.ListIPRangesFunc(ctx, token, params)
}

// CreateIPRange delegates to CreateIPRangeFunc.
func (m *MockRepo) CreateIPRange(
	ctx context.Context, token string, in domain.IPRangeWrite,
) (*domain.IPRange, error) {
	return m.CreateIPRangeFunc(ctx, token, in)
}

// UpdateIPRange delegates to UpdateIPRangeFunc.
func (m *MockRepo) UpdateIPRange(
	ctx context.Context, token string, id int, in domain.IPRangeWrite,
) (*domain.IPRange, error) {
	return m.UpdateIPRangeFunc(ctx, token, id, in)
}

// DeleteIPRange delegates to DeleteIPRangeFunc.
func (m *MockRepo) DeleteIPRange(ctx context.Context, token string, id int) error {
	return m.DeleteIPRangeFunc(ctx, token, id)
}

// ListASNs delegates to ListASNsFunc.
func (m *MockRepo) ListASNs(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.ASN], error) {
	return m.ListASNsFunc(ctx, token, params)
}

// CreateASN delegates to CreateASNFunc.
func (m *MockRepo) CreateASN(
	ctx context.Context, token string, in domain.ASNWrite,
) (*domain.ASN, error) {
	return m.CreateASNFunc(ctx, token, in)
}

// UpdateASN delegates to UpdateASNFunc.
func (m *MockRepo) UpdateASN(
	ctx context.Context, token string, id int, in domain.ASNWrite,
) (*domain.ASN, error) {
	return m.UpdateASNFunc(ctx, token, id, in)
}

// DeleteASN delegates to DeleteASNFunc.
func (m *MockRepo) DeleteASN(ctx context.Context, token string, id int) error {
	return m.DeleteASNFunc(ctx, token, id)
}

// ListASNRanges delegates to ListASNRangesFunc.
func (m *MockRepo) ListASNRanges(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.ASNRange], error) {
	return m.ListASNRangesFunc(ctx, token, params)
}

// CreateASNRange delegates to CreateASNRangeFunc.
func (m *MockRepo) CreateASNRange(
	ctx context.Context, token string, in domain.ASNRangeWrite,
) (*domain.ASNRange, error) {
	return m.CreateASNRangeFunc(ctx, token, in)
}

// UpdateASNRange delegates to UpdateASNRangeFunc.
func (m *MockRepo) UpdateASNRange(
	ctx context.Context, token string, id int, in domain.ASNRangeWrite,
) (*domain.ASNRange, error) {
	return m.UpdateASNRangeFunc(ctx, token, id, in)
}

// DeleteASNRange delegates to DeleteASNRangeFunc.
func (m *MockRepo) DeleteASNRange(ctx context.Context, token string, id int) error {
	return m.DeleteASNRangeFunc(ctx, token, id)
}

// ListRouteTargets delegates to ListRouteTargetsFunc.
func (m *MockRepo) ListRouteTargets(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.RouteTarget], error) {
	return m.ListRouteTargetsFunc(ctx, token, params)
}

// CreateRouteTarget delegates to CreateRouteTargetFunc.
func (m *MockRepo) CreateRouteTarget(
	ctx context.Context, token string, in domain.RouteTargetWrite,
) (*domain.RouteTarget, error) {
	return m.CreateRouteTargetFunc(ctx, token, in)
}

// UpdateRouteTarget delegates to UpdateRouteTargetFunc.
func (m *MockRepo) UpdateRouteTarget(
	ctx context.Context, token string, id int, in domain.RouteTargetWrite,
) (*domain.RouteTarget, error) {
	return m.UpdateRouteTargetFunc(ctx, token, id, in)
}

// DeleteRouteTarget delegates to DeleteRouteTargetFunc.
func (m *MockRepo) DeleteRouteTarget(ctx context.Context, token string, id int) error {
	return m.DeleteRouteTargetFunc(ctx, token, id)
}

// ListServices delegates to ListServicesFunc.
func (m *MockRepo) ListServices(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.Service], error) {
	return m.ListServicesFunc(ctx, token, params)
}

// CreateService delegates to CreateServiceFunc.
func (m *MockRepo) CreateService(
	ctx context.Context, token string, in domain.ServiceWrite,
) (*domain.Service, error) {
	return m.CreateServiceFunc(ctx, token, in)
}

// UpdateService delegates to UpdateServiceFunc.
func (m *MockRepo) UpdateService(
	ctx context.Context, token string, id int, in domain.ServiceWrite,
) (*domain.Service, error) {
	return m.UpdateServiceFunc(ctx, token, id, in)
}

// DeleteService delegates to DeleteServiceFunc.
func (m *MockRepo) DeleteService(ctx context.Context, token string, id int) error {
	return m.DeleteServiceFunc(ctx, token, id)
}

// ListServiceTemplates delegates to ListServiceTemplatesFunc.
func (m *MockRepo) ListServiceTemplates(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.ServiceTemplate], error) {
	return m.ListServiceTemplatesFunc(ctx, token, params)
}

// CreateServiceTemplate delegates to CreateServiceTemplateFunc.
func (m *MockRepo) CreateServiceTemplate(
	ctx context.Context, token string, in domain.ServiceTemplateWrite,
) (*domain.ServiceTemplate, error) {
	return m.CreateServiceTemplateFunc(ctx, token, in)
}

// UpdateServiceTemplate delegates to UpdateServiceTemplateFunc.
func (m *MockRepo) UpdateServiceTemplate(
	ctx context.Context, token string, id int, in domain.ServiceTemplateWrite,
) (*domain.ServiceTemplate, error) {
	return m.UpdateServiceTemplateFunc(ctx, token, id, in)
}

// DeleteServiceTemplate delegates to DeleteServiceTemplateFunc.
func (m *MockRepo) DeleteServiceTemplate(ctx context.Context, token string, id int) error {
	return m.DeleteServiceTemplateFunc(ctx, token, id)
}

// ListFHRPGroups delegates to ListFHRPGroupsFunc.
func (m *MockRepo) ListFHRPGroups(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.FHRPGroup], error) {
	return m.ListFHRPGroupsFunc(ctx, token, params)
}

// CreateFHRPGroup delegates to CreateFHRPGroupFunc.
func (m *MockRepo) CreateFHRPGroup(
	ctx context.Context, token string, in domain.FHRPGroupWrite,
) (*domain.FHRPGroup, error) {
	return m.CreateFHRPGroupFunc(ctx, token, in)
}

// UpdateFHRPGroup delegates to UpdateFHRPGroupFunc.
func (m *MockRepo) UpdateFHRPGroup(
	ctx context.Context, token string, id int, in domain.FHRPGroupWrite,
) (*domain.FHRPGroup, error) {
	return m.UpdateFHRPGroupFunc(ctx, token, id, in)
}

// DeleteFHRPGroup delegates to DeleteFHRPGroupFunc.
func (m *MockRepo) DeleteFHRPGroup(ctx context.Context, token string, id int) error {
	return m.DeleteFHRPGroupFunc(ctx, token, id)
}

// ListFHRPGroupAssignments delegates to ListFHRPGroupAssignmentsFunc.
func (m *MockRepo) ListFHRPGroupAssignments(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.FHRPGroupAssignment], error) {
	return m.ListFHRPGroupAssignmentsFunc(ctx, token, params)
}

// CreateFHRPGroupAssignment delegates to CreateFHRPGroupAssignmentFunc.
func (m *MockRepo) CreateFHRPGroupAssignment(
	ctx context.Context, token string, in domain.FHRPGroupAssignmentWrite,
) (*domain.FHRPGroupAssignment, error) {
	return m.CreateFHRPGroupAssignmentFunc(ctx, token, in)
}

// UpdateFHRPGroupAssignment delegates to UpdateFHRPGroupAssignmentFunc.
func (m *MockRepo) UpdateFHRPGroupAssignment(
	ctx context.Context, token string, id int, in domain.FHRPGroupAssignmentWrite,
) (*domain.FHRPGroupAssignment, error) {
	return m.UpdateFHRPGroupAssignmentFunc(ctx, token, id, in)
}

// DeleteFHRPGroupAssignment delegates to DeleteFHRPGroupAssignmentFunc.
func (m *MockRepo) DeleteFHRPGroupAssignment(ctx context.Context, token string, id int) error {
	return m.DeleteFHRPGroupAssignmentFunc(ctx, token, id)
}

// ListVLANTranslationPolicies delegates to ListVLANTranslationPoliciesFunc.
func (m *MockRepo) ListVLANTranslationPolicies(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.VLANTranslationPolicy], error) {
	return m.ListVLANTranslationPoliciesFunc(ctx, token, params)
}

// CreateVLANTranslationPolicy delegates to CreateVLANTranslationPolicyFunc.
func (m *MockRepo) CreateVLANTranslationPolicy(
	ctx context.Context, token string, in domain.VLANTranslationPolicyWrite,
) (*domain.VLANTranslationPolicy, error) {
	return m.CreateVLANTranslationPolicyFunc(ctx, token, in)
}

// UpdateVLANTranslationPolicy delegates to UpdateVLANTranslationPolicyFunc.
func (m *MockRepo) UpdateVLANTranslationPolicy(
	ctx context.Context, token string, id int, in domain.VLANTranslationPolicyWrite,
) (*domain.VLANTranslationPolicy, error) {
	return m.UpdateVLANTranslationPolicyFunc(ctx, token, id, in)
}

// DeleteVLANTranslationPolicy delegates to DeleteVLANTranslationPolicyFunc.
func (m *MockRepo) DeleteVLANTranslationPolicy(ctx context.Context, token string, id int) error {
	return m.DeleteVLANTranslationPolicyFunc(ctx, token, id)
}

// ListVLANTranslationRules delegates to ListVLANTranslationRulesFunc.
func (m *MockRepo) ListVLANTranslationRules(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.VLANTranslationRule], error) {
	return m.ListVLANTranslationRulesFunc(ctx, token, params)
}

// CreateVLANTranslationRule delegates to CreateVLANTranslationRuleFunc.
func (m *MockRepo) CreateVLANTranslationRule(
	ctx context.Context, token string, in domain.VLANTranslationRuleWrite,
) (*domain.VLANTranslationRule, error) {
	return m.CreateVLANTranslationRuleFunc(ctx, token, in)
}

// UpdateVLANTranslationRule delegates to UpdateVLANTranslationRuleFunc.
func (m *MockRepo) UpdateVLANTranslationRule(
	ctx context.Context, token string, id int, in domain.VLANTranslationRuleWrite,
) (*domain.VLANTranslationRule, error) {
	return m.UpdateVLANTranslationRuleFunc(ctx, token, id, in)
}

// DeleteVLANTranslationRule delegates to DeleteVLANTranslationRuleFunc.
func (m *MockRepo) DeleteVLANTranslationRule(ctx context.Context, token string, id int) error {
	return m.DeleteVLANTranslationRuleFunc(ctx, token, id)
}

// ListContactGroups delegates to ListContactGroupsFunc.
func (m *MockRepo) ListContactGroups(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.ContactGroup], error) {
	return m.ListContactGroupsFunc(ctx, token, params)
}

// CreateContactGroup delegates to CreateContactGroupFunc.
func (m *MockRepo) CreateContactGroup(
	ctx context.Context, token string, in domain.ContactGroupWrite,
) (*domain.ContactGroup, error) {
	return m.CreateContactGroupFunc(ctx, token, in)
}

// UpdateContactGroup delegates to UpdateContactGroupFunc.
func (m *MockRepo) UpdateContactGroup(
	ctx context.Context, token string, id int, in domain.ContactGroupWrite,
) (*domain.ContactGroup, error) {
	return m.UpdateContactGroupFunc(ctx, token, id, in)
}

// DeleteContactGroup delegates to DeleteContactGroupFunc.
func (m *MockRepo) DeleteContactGroup(ctx context.Context, token string, id int) error {
	return m.DeleteContactGroupFunc(ctx, token, id)
}

// ListContactRoles delegates to ListContactRolesFunc.
func (m *MockRepo) ListContactRoles(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.ContactRole], error) {
	return m.ListContactRolesFunc(ctx, token, params)
}

// CreateContactRole delegates to CreateContactRoleFunc.
func (m *MockRepo) CreateContactRole(
	ctx context.Context, token string, in domain.ContactRoleWrite,
) (*domain.ContactRole, error) {
	return m.CreateContactRoleFunc(ctx, token, in)
}

// UpdateContactRole delegates to UpdateContactRoleFunc.
func (m *MockRepo) UpdateContactRole(
	ctx context.Context, token string, id int, in domain.ContactRoleWrite,
) (*domain.ContactRole, error) {
	return m.UpdateContactRoleFunc(ctx, token, id, in)
}

// DeleteContactRole delegates to DeleteContactRoleFunc.
func (m *MockRepo) DeleteContactRole(ctx context.Context, token string, id int) error {
	return m.DeleteContactRoleFunc(ctx, token, id)
}

// ListVirtualDisks delegates to ListVirtualDisksFunc.
func (m *MockRepo) ListVirtualDisks(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.VirtualDisk], error) {
	return m.ListVirtualDisksFunc(ctx, token, params)
}

// CreateVirtualDisk delegates to CreateVirtualDiskFunc.
func (m *MockRepo) CreateVirtualDisk(
	ctx context.Context, token string, in domain.VirtualDiskWrite,
) (*domain.VirtualDisk, error) {
	return m.CreateVirtualDiskFunc(ctx, token, in)
}

// UpdateVirtualDisk delegates to UpdateVirtualDiskFunc.
func (m *MockRepo) UpdateVirtualDisk(
	ctx context.Context, token string, id int, in domain.VirtualDiskWrite,
) (*domain.VirtualDisk, error) {
	return m.UpdateVirtualDiskFunc(ctx, token, id, in)
}

// DeleteVirtualDisk delegates to DeleteVirtualDiskFunc.
func (m *MockRepo) DeleteVirtualDisk(ctx context.Context, token string, id int) error {
	return m.DeleteVirtualDiskFunc(ctx, token, id)
}

// ListVirtualMachineTypes delegates to ListVirtualMachineTypesFunc.
func (m *MockRepo) ListVirtualMachineTypes(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.VirtualMachineType], error) {
	return m.ListVirtualMachineTypesFunc(ctx, token, params)
}

// CreateVirtualMachineType delegates to CreateVirtualMachineTypeFunc.
func (m *MockRepo) CreateVirtualMachineType(
	ctx context.Context, token string, in domain.VirtualMachineTypeWrite,
) (*domain.VirtualMachineType, error) {
	return m.CreateVirtualMachineTypeFunc(ctx, token, in)
}

// UpdateVirtualMachineType delegates to UpdateVirtualMachineTypeFunc.
func (m *MockRepo) UpdateVirtualMachineType(
	ctx context.Context, token string, id int, in domain.VirtualMachineTypeWrite,
) (*domain.VirtualMachineType, error) {
	return m.UpdateVirtualMachineTypeFunc(ctx, token, id, in)
}

// DeleteVirtualMachineType delegates to DeleteVirtualMachineTypeFunc.
func (m *MockRepo) DeleteVirtualMachineType(ctx context.Context, token string, id int) error {
	return m.DeleteVirtualMachineTypeFunc(ctx, token, id)
}

// ListProviderAccounts delegates to ListProviderAccountsFunc.
func (m *MockRepo) ListProviderAccounts(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.ProviderAccount], error) {
	return m.ListProviderAccountsFunc(ctx, token, params)
}

// CreateProviderAccount delegates to CreateProviderAccountFunc.
func (m *MockRepo) CreateProviderAccount(
	ctx context.Context, token string, in domain.ProviderAccountWrite,
) (*domain.ProviderAccount, error) {
	return m.CreateProviderAccountFunc(ctx, token, in)
}

// UpdateProviderAccount delegates to UpdateProviderAccountFunc.
func (m *MockRepo) UpdateProviderAccount(
	ctx context.Context, token string, id int, in domain.ProviderAccountWrite,
) (*domain.ProviderAccount, error) {
	return m.UpdateProviderAccountFunc(ctx, token, id, in)
}

// DeleteProviderAccount delegates to DeleteProviderAccountFunc.
func (m *MockRepo) DeleteProviderAccount(ctx context.Context, token string, id int) error {
	return m.DeleteProviderAccountFunc(ctx, token, id)
}

// ListProviderNetworks delegates to ListProviderNetworksFunc.
func (m *MockRepo) ListProviderNetworks(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.ProviderNetwork], error) {
	return m.ListProviderNetworksFunc(ctx, token, params)
}

// CreateProviderNetwork delegates to CreateProviderNetworkFunc.
func (m *MockRepo) CreateProviderNetwork(
	ctx context.Context, token string, in domain.ProviderNetworkWrite,
) (*domain.ProviderNetwork, error) {
	return m.CreateProviderNetworkFunc(ctx, token, in)
}

// UpdateProviderNetwork delegates to UpdateProviderNetworkFunc.
func (m *MockRepo) UpdateProviderNetwork(
	ctx context.Context, token string, id int, in domain.ProviderNetworkWrite,
) (*domain.ProviderNetwork, error) {
	return m.UpdateProviderNetworkFunc(ctx, token, id, in)
}

// DeleteProviderNetwork delegates to DeleteProviderNetworkFunc.
func (m *MockRepo) DeleteProviderNetwork(ctx context.Context, token string, id int) error {
	return m.DeleteProviderNetworkFunc(ctx, token, id)
}

// ListCircuitGroups delegates to ListCircuitGroupsFunc.
func (m *MockRepo) ListCircuitGroups(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.CircuitGroup], error) {
	return m.ListCircuitGroupsFunc(ctx, token, params)
}

// CreateCircuitGroup delegates to CreateCircuitGroupFunc.
func (m *MockRepo) CreateCircuitGroup(
	ctx context.Context, token string, in domain.CircuitGroupWrite,
) (*domain.CircuitGroup, error) {
	return m.CreateCircuitGroupFunc(ctx, token, in)
}

// UpdateCircuitGroup delegates to UpdateCircuitGroupFunc.
func (m *MockRepo) UpdateCircuitGroup(
	ctx context.Context, token string, id int, in domain.CircuitGroupWrite,
) (*domain.CircuitGroup, error) {
	return m.UpdateCircuitGroupFunc(ctx, token, id, in)
}

// DeleteCircuitGroup delegates to DeleteCircuitGroupFunc.
func (m *MockRepo) DeleteCircuitGroup(ctx context.Context, token string, id int) error {
	return m.DeleteCircuitGroupFunc(ctx, token, id)
}

// ListCircuitGroupAssignments delegates to ListCircuitGroupAssignmentsFunc.
func (m *MockRepo) ListCircuitGroupAssignments(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.CircuitGroupAssignment], error) {
	return m.ListCircuitGroupAssignmentsFunc(ctx, token, params)
}

// CreateCircuitGroupAssignment delegates to CreateCircuitGroupAssignmentFunc.
func (m *MockRepo) CreateCircuitGroupAssignment(
	ctx context.Context, token string, in domain.CircuitGroupAssignmentWrite,
) (*domain.CircuitGroupAssignment, error) {
	return m.CreateCircuitGroupAssignmentFunc(ctx, token, in)
}

// UpdateCircuitGroupAssignment delegates to UpdateCircuitGroupAssignmentFunc.
func (m *MockRepo) UpdateCircuitGroupAssignment(
	ctx context.Context, token string, id int, in domain.CircuitGroupAssignmentWrite,
) (*domain.CircuitGroupAssignment, error) {
	return m.UpdateCircuitGroupAssignmentFunc(ctx, token, id, in)
}

// DeleteCircuitGroupAssignment delegates to DeleteCircuitGroupAssignmentFunc.
func (m *MockRepo) DeleteCircuitGroupAssignment(ctx context.Context, token string, id int) error {
	return m.DeleteCircuitGroupAssignmentFunc(ctx, token, id)
}

// ListVirtualCircuits delegates to ListVirtualCircuitsFunc.
func (m *MockRepo) ListVirtualCircuits(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.VirtualCircuit], error) {
	return m.ListVirtualCircuitsFunc(ctx, token, params)
}

// CreateVirtualCircuit delegates to CreateVirtualCircuitFunc.
func (m *MockRepo) CreateVirtualCircuit(
	ctx context.Context, token string, in domain.VirtualCircuitWrite,
) (*domain.VirtualCircuit, error) {
	return m.CreateVirtualCircuitFunc(ctx, token, in)
}

// UpdateVirtualCircuit delegates to UpdateVirtualCircuitFunc.
func (m *MockRepo) UpdateVirtualCircuit(
	ctx context.Context, token string, id int, in domain.VirtualCircuitWrite,
) (*domain.VirtualCircuit, error) {
	return m.UpdateVirtualCircuitFunc(ctx, token, id, in)
}

// DeleteVirtualCircuit delegates to DeleteVirtualCircuitFunc.
func (m *MockRepo) DeleteVirtualCircuit(ctx context.Context, token string, id int) error {
	return m.DeleteVirtualCircuitFunc(ctx, token, id)
}

// ListVirtualCircuitTypes delegates to ListVirtualCircuitTypesFunc.
func (m *MockRepo) ListVirtualCircuitTypes(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.VirtualCircuitType], error) {
	return m.ListVirtualCircuitTypesFunc(ctx, token, params)
}

// CreateVirtualCircuitType delegates to CreateVirtualCircuitTypeFunc.
func (m *MockRepo) CreateVirtualCircuitType(
	ctx context.Context, token string, in domain.VirtualCircuitTypeWrite,
) (*domain.VirtualCircuitType, error) {
	return m.CreateVirtualCircuitTypeFunc(ctx, token, in)
}

// UpdateVirtualCircuitType delegates to UpdateVirtualCircuitTypeFunc.
func (m *MockRepo) UpdateVirtualCircuitType(
	ctx context.Context, token string, id int, in domain.VirtualCircuitTypeWrite,
) (*domain.VirtualCircuitType, error) {
	return m.UpdateVirtualCircuitTypeFunc(ctx, token, id, in)
}

// DeleteVirtualCircuitType delegates to DeleteVirtualCircuitTypeFunc.
func (m *MockRepo) DeleteVirtualCircuitType(ctx context.Context, token string, id int) error {
	return m.DeleteVirtualCircuitTypeFunc(ctx, token, id)
}

// ListVirtualCircuitTerminations delegates to ListVirtualCircuitTerminationsFunc.
func (m *MockRepo) ListVirtualCircuitTerminations(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.VirtualCircuitTermination], error) {
	return m.ListVirtualCircuitTerminationsFunc(ctx, token, params)
}

// CreateVirtualCircuitTermination delegates to CreateVirtualCircuitTerminationFunc.
func (m *MockRepo) CreateVirtualCircuitTermination(
	ctx context.Context, token string, in domain.VirtualCircuitTerminationWrite,
) (*domain.VirtualCircuitTermination, error) {
	return m.CreateVirtualCircuitTerminationFunc(ctx, token, in)
}

// UpdateVirtualCircuitTermination delegates to UpdateVirtualCircuitTerminationFunc.
func (m *MockRepo) UpdateVirtualCircuitTermination(
	ctx context.Context, token string, id int, in domain.VirtualCircuitTerminationWrite,
) (*domain.VirtualCircuitTermination, error) {
	return m.UpdateVirtualCircuitTerminationFunc(ctx, token, id, in)
}

// DeleteVirtualCircuitTermination delegates to DeleteVirtualCircuitTerminationFunc.
func (m *MockRepo) DeleteVirtualCircuitTermination(ctx context.Context, token string, id int) error {
	return m.DeleteVirtualCircuitTerminationFunc(ctx, token, id)
}

// ListRoles delegates to ListRolesFunc.
func (m *MockRepo) ListRoles(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.Role], error) {
	return m.ListRolesFunc(ctx, token, params)
}

// ListTenants delegates to ListTenantsFunc.
func (m *MockRepo) ListTenants(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.Tenant], error) {
	return m.ListTenantsFunc(ctx, token, params)
}

// ListTags delegates to ListTagsFunc.
func (m *MockRepo) ListTags(
	ctx context.Context, token string, params map[string]string,
) (*domain.PaginatedResponse[domain.Tag], error) {
	return m.ListTagsFunc(ctx, token, params)
}

// CreateTag delegates to CreateTagFunc.
func (m *MockRepo) CreateTag(
	ctx context.Context, token string, in domain.TagWrite,
) (*domain.Tag, error) {
	return m.CreateTagFunc(ctx, token, in)
}

// UpdateTag delegates to UpdateTagFunc.
func (m *MockRepo) UpdateTag(
	ctx context.Context, token string, id int, in domain.TagWrite,
) (*domain.Tag, error) {
	return m.UpdateTagFunc(ctx, token, id, in)
}

// DeleteTag delegates to DeleteTagFunc.
func (m *MockRepo) DeleteTag(ctx context.Context, token string, id int) error {
	return m.DeleteTagFunc(ctx, token, id)
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

// ListSiteGroups returns an empty site group page.
func (r *NilRepo) ListSiteGroups(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.SiteGroup], error) {
	return &domain.PaginatedResponse[domain.SiteGroup]{Count: 0, Results: []domain.SiteGroup{}}, nil
}

// CreateSiteGroup returns a site group with id 1.
func (r *NilRepo) CreateSiteGroup(
	_ context.Context, _ string, in domain.SiteGroupWrite,
) (*domain.SiteGroup, error) {
	return &domain.SiteGroup{ID: 1, Name: in.Name}, nil
}

// UpdateSiteGroup returns the site group with the given id.
func (r *NilRepo) UpdateSiteGroup(
	_ context.Context, _ string, id int, in domain.SiteGroupWrite,
) (*domain.SiteGroup, error) {
	return &domain.SiteGroup{ID: id, Name: in.Name}, nil
}

// DeleteSiteGroup always succeeds.
func (r *NilRepo) DeleteSiteGroup(_ context.Context, _ string, _ int) error {
	return nil
}

// ListRackGroups returns an empty rack group page.
func (r *NilRepo) ListRackGroups(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.RackGroup], error) {
	return &domain.PaginatedResponse[domain.RackGroup]{Count: 0, Results: []domain.RackGroup{}}, nil
}

// CreateRackGroup returns a rack group with id 1.
func (r *NilRepo) CreateRackGroup(
	_ context.Context, _ string, in domain.RackGroupWrite,
) (*domain.RackGroup, error) {
	return &domain.RackGroup{ID: 1, Name: in.Name}, nil
}

// UpdateRackGroup returns the rack group with the given id.
func (r *NilRepo) UpdateRackGroup(
	_ context.Context, _ string, id int, in domain.RackGroupWrite,
) (*domain.RackGroup, error) {
	return &domain.RackGroup{ID: id, Name: in.Name}, nil
}

// DeleteRackGroup always succeeds.
func (r *NilRepo) DeleteRackGroup(_ context.Context, _ string, _ int) error {
	return nil
}

// ListRackRoles returns an empty rack role page.
func (r *NilRepo) ListRackRoles(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.RackRole], error) {
	return &domain.PaginatedResponse[domain.RackRole]{Count: 0, Results: []domain.RackRole{}}, nil
}

// CreateRackRole returns a rack role with id 1.
func (r *NilRepo) CreateRackRole(
	_ context.Context, _ string, in domain.RackRoleWrite,
) (*domain.RackRole, error) {
	return &domain.RackRole{ID: 1, Name: in.Name}, nil
}

// UpdateRackRole returns the rack role with the given id.
func (r *NilRepo) UpdateRackRole(
	_ context.Context, _ string, id int, in domain.RackRoleWrite,
) (*domain.RackRole, error) {
	return &domain.RackRole{ID: id, Name: in.Name}, nil
}

// DeleteRackRole always succeeds.
func (r *NilRepo) DeleteRackRole(_ context.Context, _ string, _ int) error {
	return nil
}

// ListDeviceRoles returns an empty device role page.
func (r *NilRepo) ListDeviceRoles(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.DeviceRole], error) {
	return &domain.PaginatedResponse[domain.DeviceRole]{Count: 0, Results: []domain.DeviceRole{}}, nil
}

// CreateDeviceRole returns a device role with id 1.
func (r *NilRepo) CreateDeviceRole(
	_ context.Context, _ string, in domain.DeviceRoleWrite,
) (*domain.DeviceRole, error) {
	return &domain.DeviceRole{ID: 1, Name: in.Name}, nil
}

// UpdateDeviceRole returns the device role with the given id.
func (r *NilRepo) UpdateDeviceRole(
	_ context.Context, _ string, id int, in domain.DeviceRoleWrite,
) (*domain.DeviceRole, error) {
	return &domain.DeviceRole{ID: id, Name: in.Name}, nil
}

// DeleteDeviceRole always succeeds.
func (r *NilRepo) DeleteDeviceRole(_ context.Context, _ string, _ int) error {
	return nil
}

// ListPlatforms returns an empty platform page.
func (r *NilRepo) ListPlatforms(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Platform], error) {
	return &domain.PaginatedResponse[domain.Platform]{Count: 0, Results: []domain.Platform{}}, nil
}

// CreatePlatform returns a platform with id 1.
func (r *NilRepo) CreatePlatform(
	_ context.Context, _ string, in domain.PlatformWrite,
) (*domain.Platform, error) {
	return &domain.Platform{ID: 1, Name: in.Name}, nil
}

// UpdatePlatform returns the platform with the given id.
func (r *NilRepo) UpdatePlatform(
	_ context.Context, _ string, id int, in domain.PlatformWrite,
) (*domain.Platform, error) {
	return &domain.Platform{ID: id, Name: in.Name}, nil
}

// DeletePlatform always succeeds.
func (r *NilRepo) DeletePlatform(_ context.Context, _ string, _ int) error {
	return nil
}

// ListVirtualChassis returns an empty virtual chassis page.
func (r *NilRepo) ListVirtualChassis(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.VirtualChassis], error) {
	return &domain.PaginatedResponse[domain.VirtualChassis]{Count: 0, Results: []domain.VirtualChassis{}}, nil
}

// CreateVirtualChassis returns a virtual chassis with id 1.
func (r *NilRepo) CreateVirtualChassis(
	_ context.Context, _ string, in domain.VirtualChassisWrite,
) (*domain.VirtualChassis, error) {
	return &domain.VirtualChassis{ID: 1, Name: in.Name}, nil
}

// UpdateVirtualChassis returns the virtual chassis with the given id.
func (r *NilRepo) UpdateVirtualChassis(
	_ context.Context, _ string, id int, in domain.VirtualChassisWrite,
) (*domain.VirtualChassis, error) {
	return &domain.VirtualChassis{ID: id, Name: in.Name}, nil
}

// DeleteVirtualChassis always succeeds.
func (r *NilRepo) DeleteVirtualChassis(_ context.Context, _ string, _ int) error {
	return nil
}

// ListPowerPanels returns an empty power panel page.
func (r *NilRepo) ListPowerPanels(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.PowerPanel], error) {
	return &domain.PaginatedResponse[domain.PowerPanel]{Count: 0, Results: []domain.PowerPanel{}}, nil
}

// CreatePowerPanel returns a power panel with id 1.
func (r *NilRepo) CreatePowerPanel(
	_ context.Context, _ string, in domain.PowerPanelWrite,
) (*domain.PowerPanel, error) {
	return &domain.PowerPanel{ID: 1, Name: in.Name}, nil
}

// UpdatePowerPanel returns the power panel with the given id.
func (r *NilRepo) UpdatePowerPanel(
	_ context.Context, _ string, id int, in domain.PowerPanelWrite,
) (*domain.PowerPanel, error) {
	return &domain.PowerPanel{ID: id, Name: in.Name}, nil
}

// DeletePowerPanel always succeeds.
func (r *NilRepo) DeletePowerPanel(_ context.Context, _ string, _ int) error {
	return nil
}

// ListInventoryItemRoles returns an empty inventory item role page.
func (r *NilRepo) ListInventoryItemRoles(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.InventoryItemRole], error) {
	return &domain.PaginatedResponse[domain.InventoryItemRole]{Count: 0, Results: []domain.InventoryItemRole{}}, nil
}

// CreateInventoryItemRole returns an inventory item role with id 1.
func (r *NilRepo) CreateInventoryItemRole(
	_ context.Context, _ string, in domain.InventoryItemRoleWrite,
) (*domain.InventoryItemRole, error) {
	return &domain.InventoryItemRole{ID: 1, Name: in.Name}, nil
}

// UpdateInventoryItemRole returns the inventory item role with the given id.
func (r *NilRepo) UpdateInventoryItemRole(
	_ context.Context, _ string, id int, in domain.InventoryItemRoleWrite,
) (*domain.InventoryItemRole, error) {
	return &domain.InventoryItemRole{ID: id, Name: in.Name}, nil
}

// DeleteInventoryItemRole always succeeds.
func (r *NilRepo) DeleteInventoryItemRole(_ context.Context, _ string, _ int) error {
	return nil
}

// ListFrontPorts returns an empty front port page.
func (r *NilRepo) ListFrontPorts(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.FrontPort], error) {
	return &domain.PaginatedResponse[domain.FrontPort]{Count: 0, Results: []domain.FrontPort{}}, nil
}

// CreateFrontPort returns a front port with id 1.
func (r *NilRepo) CreateFrontPort(
	_ context.Context, _ string, in domain.FrontPortWrite,
) (*domain.FrontPort, error) {
	return &domain.FrontPort{ID: 1, Name: in.Name}, nil
}

// UpdateFrontPort returns the front port with the given id.
func (r *NilRepo) UpdateFrontPort(
	_ context.Context, _ string, id int, in domain.FrontPortWrite,
) (*domain.FrontPort, error) {
	return &domain.FrontPort{ID: id, Name: in.Name}, nil
}

// DeleteFrontPort always succeeds.
func (r *NilRepo) DeleteFrontPort(_ context.Context, _ string, _ int) error {
	return nil
}

// ListRearPorts returns an empty rear port page.
func (r *NilRepo) ListRearPorts(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.RearPort], error) {
	return &domain.PaginatedResponse[domain.RearPort]{Count: 0, Results: []domain.RearPort{}}, nil
}

// CreateRearPort returns a rear port with id 1.
func (r *NilRepo) CreateRearPort(
	_ context.Context, _ string, in domain.RearPortWrite,
) (*domain.RearPort, error) {
	return &domain.RearPort{ID: 1, Name: in.Name}, nil
}

// UpdateRearPort returns the rear port with the given id.
func (r *NilRepo) UpdateRearPort(
	_ context.Context, _ string, id int, in domain.RearPortWrite,
) (*domain.RearPort, error) {
	return &domain.RearPort{ID: id, Name: in.Name}, nil
}

// DeleteRearPort always succeeds.
func (r *NilRepo) DeleteRearPort(_ context.Context, _ string, _ int) error {
	return nil
}

// ListConsolePorts returns an empty console port page.
func (r *NilRepo) ListConsolePorts(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.ConsolePort], error) {
	return &domain.PaginatedResponse[domain.ConsolePort]{Count: 0, Results: []domain.ConsolePort{}}, nil
}

// CreateConsolePort returns a console port with id 1.
func (r *NilRepo) CreateConsolePort(
	_ context.Context, _ string, in domain.ConsolePortWrite,
) (*domain.ConsolePort, error) {
	return &domain.ConsolePort{ID: 1, Name: in.Name}, nil
}

// UpdateConsolePort returns the console port with the given id.
func (r *NilRepo) UpdateConsolePort(
	_ context.Context, _ string, id int, in domain.ConsolePortWrite,
) (*domain.ConsolePort, error) {
	return &domain.ConsolePort{ID: id, Name: in.Name}, nil
}

// DeleteConsolePort always succeeds.
func (r *NilRepo) DeleteConsolePort(_ context.Context, _ string, _ int) error {
	return nil
}

// ListConsoleServerPorts returns an empty console server port page.
func (r *NilRepo) ListConsoleServerPorts(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.ConsoleServerPort], error) {
	return &domain.PaginatedResponse[domain.ConsoleServerPort]{Count: 0, Results: []domain.ConsoleServerPort{}}, nil
}

// CreateConsoleServerPort returns a console server port with id 1.
func (r *NilRepo) CreateConsoleServerPort(
	_ context.Context, _ string, in domain.ConsoleServerPortWrite,
) (*domain.ConsoleServerPort, error) {
	return &domain.ConsoleServerPort{ID: 1, Name: in.Name}, nil
}

// UpdateConsoleServerPort returns the console server port with the given id.
func (r *NilRepo) UpdateConsoleServerPort(
	_ context.Context, _ string, id int, in domain.ConsoleServerPortWrite,
) (*domain.ConsoleServerPort, error) {
	return &domain.ConsoleServerPort{ID: id, Name: in.Name}, nil
}

// DeleteConsoleServerPort always succeeds.
func (r *NilRepo) DeleteConsoleServerPort(_ context.Context, _ string, _ int) error {
	return nil
}

// ListPowerPorts returns an empty power port page.
func (r *NilRepo) ListPowerPorts(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.PowerPort], error) {
	return &domain.PaginatedResponse[domain.PowerPort]{Count: 0, Results: []domain.PowerPort{}}, nil
}

// CreatePowerPort returns a power port with id 1.
func (r *NilRepo) CreatePowerPort(
	_ context.Context, _ string, in domain.PowerPortWrite,
) (*domain.PowerPort, error) {
	return &domain.PowerPort{ID: 1, Name: in.Name}, nil
}

// UpdatePowerPort returns the power port with the given id.
func (r *NilRepo) UpdatePowerPort(
	_ context.Context, _ string, id int, in domain.PowerPortWrite,
) (*domain.PowerPort, error) {
	return &domain.PowerPort{ID: id, Name: in.Name}, nil
}

// DeletePowerPort always succeeds.
func (r *NilRepo) DeletePowerPort(_ context.Context, _ string, _ int) error {
	return nil
}

// ListPowerOutlets returns an empty power outlet page.
func (r *NilRepo) ListPowerOutlets(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.PowerOutlet], error) {
	return &domain.PaginatedResponse[domain.PowerOutlet]{Count: 0, Results: []domain.PowerOutlet{}}, nil
}

// CreatePowerOutlet returns a power outlet with id 1.
func (r *NilRepo) CreatePowerOutlet(
	_ context.Context, _ string, in domain.PowerOutletWrite,
) (*domain.PowerOutlet, error) {
	return &domain.PowerOutlet{ID: 1, Name: in.Name}, nil
}

// UpdatePowerOutlet returns the power outlet with the given id.
func (r *NilRepo) UpdatePowerOutlet(
	_ context.Context, _ string, id int, in domain.PowerOutletWrite,
) (*domain.PowerOutlet, error) {
	return &domain.PowerOutlet{ID: id, Name: in.Name}, nil
}

// DeletePowerOutlet always succeeds.
func (r *NilRepo) DeletePowerOutlet(_ context.Context, _ string, _ int) error {
	return nil
}

// ListPowerFeeds returns an empty power feed page.
func (r *NilRepo) ListPowerFeeds(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.PowerFeed], error) {
	return &domain.PaginatedResponse[domain.PowerFeed]{Count: 0, Results: []domain.PowerFeed{}}, nil
}

// CreatePowerFeed returns a power feed with id 1.
func (r *NilRepo) CreatePowerFeed(
	_ context.Context, _ string, in domain.PowerFeedWrite,
) (*domain.PowerFeed, error) {
	return &domain.PowerFeed{ID: 1, Name: in.Name}, nil
}

// UpdatePowerFeed returns the power feed with the given id.
func (r *NilRepo) UpdatePowerFeed(
	_ context.Context, _ string, id int, in domain.PowerFeedWrite,
) (*domain.PowerFeed, error) {
	return &domain.PowerFeed{ID: id, Name: in.Name}, nil
}

// DeletePowerFeed always succeeds.
func (r *NilRepo) DeletePowerFeed(_ context.Context, _ string, _ int) error {
	return nil
}

// ListDeviceBays returns an empty device bay page.
func (r *NilRepo) ListDeviceBays(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.DeviceBay], error) {
	return &domain.PaginatedResponse[domain.DeviceBay]{Count: 0, Results: []domain.DeviceBay{}}, nil
}

// CreateDeviceBay returns a device bay with id 1.
func (r *NilRepo) CreateDeviceBay(
	_ context.Context, _ string, in domain.DeviceBayWrite,
) (*domain.DeviceBay, error) {
	return &domain.DeviceBay{ID: 1, Name: in.Name}, nil
}

// UpdateDeviceBay returns the device bay with the given id.
func (r *NilRepo) UpdateDeviceBay(
	_ context.Context, _ string, id int, in domain.DeviceBayWrite,
) (*domain.DeviceBay, error) {
	return &domain.DeviceBay{ID: id, Name: in.Name}, nil
}

// DeleteDeviceBay always succeeds.
func (r *NilRepo) DeleteDeviceBay(_ context.Context, _ string, _ int) error {
	return nil
}

// ListRackReservations returns an empty rack reservation page.
func (r *NilRepo) ListRackReservations(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.RackReservation], error) {
	return &domain.PaginatedResponse[domain.RackReservation]{Count: 0, Results: []domain.RackReservation{}}, nil
}

// CreateRackReservation returns a rack reservation with id 1.
func (r *NilRepo) CreateRackReservation(
	_ context.Context, _ string, in domain.RackReservationWrite,
) (*domain.RackReservation, error) {
	return &domain.RackReservation{ID: 1, Description: in.Description}, nil
}

// UpdateRackReservation returns the rack reservation with the given id.
func (r *NilRepo) UpdateRackReservation(
	_ context.Context, _ string, id int, in domain.RackReservationWrite,
) (*domain.RackReservation, error) {
	return &domain.RackReservation{ID: id, Description: in.Description}, nil
}

// DeleteRackReservation always succeeds.
func (r *NilRepo) DeleteRackReservation(_ context.Context, _ string, _ int) error {
	return nil
}

// ListInventoryItems returns an empty inventory item page.
func (r *NilRepo) ListInventoryItems(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.InventoryItem], error) {
	return &domain.PaginatedResponse[domain.InventoryItem]{Count: 0, Results: []domain.InventoryItem{}}, nil
}

// CreateInventoryItem returns an inventory item with id 1.
func (r *NilRepo) CreateInventoryItem(
	_ context.Context, _ string, in domain.InventoryItemWrite,
) (*domain.InventoryItem, error) {
	return &domain.InventoryItem{ID: 1, Name: in.Name}, nil
}

// UpdateInventoryItem returns the inventory item with the given id.
func (r *NilRepo) UpdateInventoryItem(
	_ context.Context, _ string, id int, in domain.InventoryItemWrite,
) (*domain.InventoryItem, error) {
	return &domain.InventoryItem{ID: id, Name: in.Name}, nil
}

// DeleteInventoryItem always succeeds.
func (r *NilRepo) DeleteInventoryItem(_ context.Context, _ string, _ int) error {
	return nil
}

// ListMACAddresses returns an empty MAC address page.
func (r *NilRepo) ListMACAddresses(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.MACAddress], error) {
	return &domain.PaginatedResponse[domain.MACAddress]{Count: 0, Results: []domain.MACAddress{}}, nil
}

// CreateMACAddress returns a MAC address with id 1.
func (r *NilRepo) CreateMACAddress(
	_ context.Context, _ string, in domain.MACAddressWrite,
) (*domain.MACAddress, error) {
	return &domain.MACAddress{ID: 1, Address: in.Address}, nil
}

// UpdateMACAddress returns the MAC address with the given id.
func (r *NilRepo) UpdateMACAddress(
	_ context.Context, _ string, id int, in domain.MACAddressWrite,
) (*domain.MACAddress, error) {
	return &domain.MACAddress{ID: id, Address: in.Address}, nil
}

// DeleteMACAddress always succeeds.
func (r *NilRepo) DeleteMACAddress(_ context.Context, _ string, _ int) error {
	return nil
}

// ListRIRs returns an empty RIR page.
func (r *NilRepo) ListRIRs(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.RIR], error) {
	return &domain.PaginatedResponse[domain.RIR]{Count: 0, Results: []domain.RIR{}}, nil
}

// CreateRIR returns an RIR with id 1.
func (r *NilRepo) CreateRIR(
	_ context.Context, _ string, in domain.RIRWrite,
) (*domain.RIR, error) {
	return &domain.RIR{ID: 1, Name: in.Name}, nil
}

// UpdateRIR returns the RIR with the given id.
func (r *NilRepo) UpdateRIR(
	_ context.Context, _ string, id int, in domain.RIRWrite,
) (*domain.RIR, error) {
	return &domain.RIR{ID: id, Name: in.Name}, nil
}

// DeleteRIR always succeeds.
func (r *NilRepo) DeleteRIR(_ context.Context, _ string, _ int) error {
	return nil
}

// ListAggregates returns an empty aggregate page.
func (r *NilRepo) ListAggregates(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Aggregate], error) {
	return &domain.PaginatedResponse[domain.Aggregate]{Count: 0, Results: []domain.Aggregate{}}, nil
}

// CreateAggregate returns an aggregate with id 1.
func (r *NilRepo) CreateAggregate(
	_ context.Context, _ string, in domain.AggregateWrite,
) (*domain.Aggregate, error) {
	return &domain.Aggregate{ID: 1, Prefix: in.Prefix}, nil
}

// UpdateAggregate returns the aggregate with the given id.
func (r *NilRepo) UpdateAggregate(
	_ context.Context, _ string, id int, in domain.AggregateWrite,
) (*domain.Aggregate, error) {
	return &domain.Aggregate{ID: id, Prefix: in.Prefix}, nil
}

// DeleteAggregate always succeeds.
func (r *NilRepo) DeleteAggregate(_ context.Context, _ string, _ int) error {
	return nil
}

// ListIPRanges returns an empty IP range page.
func (r *NilRepo) ListIPRanges(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.IPRange], error) {
	return &domain.PaginatedResponse[domain.IPRange]{Count: 0, Results: []domain.IPRange{}}, nil
}

// CreateIPRange returns an IP range with id 1.
func (r *NilRepo) CreateIPRange(
	_ context.Context, _ string, in domain.IPRangeWrite,
) (*domain.IPRange, error) {
	return &domain.IPRange{ID: 1, StartAddress: in.StartAddress, EndAddress: in.EndAddress}, nil
}

// UpdateIPRange returns the IP range with the given id.
func (r *NilRepo) UpdateIPRange(
	_ context.Context, _ string, id int, in domain.IPRangeWrite,
) (*domain.IPRange, error) {
	return &domain.IPRange{ID: id, StartAddress: in.StartAddress, EndAddress: in.EndAddress}, nil
}

// DeleteIPRange always succeeds.
func (r *NilRepo) DeleteIPRange(_ context.Context, _ string, _ int) error {
	return nil
}

// ListASNs returns an empty ASN page.
func (r *NilRepo) ListASNs(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.ASN], error) {
	return &domain.PaginatedResponse[domain.ASN]{Count: 0, Results: []domain.ASN{}}, nil
}

// CreateASN returns an ASN with id 1.
func (r *NilRepo) CreateASN(
	_ context.Context, _ string, in domain.ASNWrite,
) (*domain.ASN, error) {
	return &domain.ASN{ID: 1, ASN: in.ASN}, nil
}

// UpdateASN returns the ASN with the given id.
func (r *NilRepo) UpdateASN(
	_ context.Context, _ string, id int, in domain.ASNWrite,
) (*domain.ASN, error) {
	return &domain.ASN{ID: id, ASN: in.ASN}, nil
}

// DeleteASN always succeeds.
func (r *NilRepo) DeleteASN(_ context.Context, _ string, _ int) error {
	return nil
}

// ListASNRanges returns an empty ASN range page.
func (r *NilRepo) ListASNRanges(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.ASNRange], error) {
	return &domain.PaginatedResponse[domain.ASNRange]{Count: 0, Results: []domain.ASNRange{}}, nil
}

// CreateASNRange returns an ASN range with id 1.
func (r *NilRepo) CreateASNRange(
	_ context.Context, _ string, in domain.ASNRangeWrite,
) (*domain.ASNRange, error) {
	return &domain.ASNRange{ID: 1, Name: in.Name}, nil
}

// UpdateASNRange returns the ASN range with the given id.
func (r *NilRepo) UpdateASNRange(
	_ context.Context, _ string, id int, in domain.ASNRangeWrite,
) (*domain.ASNRange, error) {
	return &domain.ASNRange{ID: id, Name: in.Name}, nil
}

// DeleteASNRange always succeeds.
func (r *NilRepo) DeleteASNRange(_ context.Context, _ string, _ int) error {
	return nil
}

// ListRouteTargets returns an empty route target page.
func (r *NilRepo) ListRouteTargets(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.RouteTarget], error) {
	return &domain.PaginatedResponse[domain.RouteTarget]{Count: 0, Results: []domain.RouteTarget{}}, nil
}

// CreateRouteTarget returns a route target with id 1.
func (r *NilRepo) CreateRouteTarget(
	_ context.Context, _ string, in domain.RouteTargetWrite,
) (*domain.RouteTarget, error) {
	return &domain.RouteTarget{ID: 1, Name: in.Name}, nil
}

// UpdateRouteTarget returns the route target with the given id.
func (r *NilRepo) UpdateRouteTarget(
	_ context.Context, _ string, id int, in domain.RouteTargetWrite,
) (*domain.RouteTarget, error) {
	return &domain.RouteTarget{ID: id, Name: in.Name}, nil
}

// DeleteRouteTarget always succeeds.
func (r *NilRepo) DeleteRouteTarget(_ context.Context, _ string, _ int) error {
	return nil
}

// ListServices returns an empty service page.
func (r *NilRepo) ListServices(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Service], error) {
	return &domain.PaginatedResponse[domain.Service]{Count: 0, Results: []domain.Service{}}, nil
}

// CreateService returns a service with id 1.
func (r *NilRepo) CreateService(
	_ context.Context, _ string, in domain.ServiceWrite,
) (*domain.Service, error) {
	return &domain.Service{ID: 1, Name: in.Name, Ports: in.Ports}, nil
}

// UpdateService returns the service with the given id.
func (r *NilRepo) UpdateService(
	_ context.Context, _ string, id int, in domain.ServiceWrite,
) (*domain.Service, error) {
	return &domain.Service{ID: id, Name: in.Name, Ports: in.Ports}, nil
}

// DeleteService always succeeds.
func (r *NilRepo) DeleteService(_ context.Context, _ string, _ int) error {
	return nil
}

// ListServiceTemplates returns an empty service template page.
func (r *NilRepo) ListServiceTemplates(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.ServiceTemplate], error) {
	return &domain.PaginatedResponse[domain.ServiceTemplate]{Count: 0, Results: []domain.ServiceTemplate{}}, nil
}

// CreateServiceTemplate returns a service template with id 1.
func (r *NilRepo) CreateServiceTemplate(
	_ context.Context, _ string, in domain.ServiceTemplateWrite,
) (*domain.ServiceTemplate, error) {
	return &domain.ServiceTemplate{ID: 1, Name: in.Name, Ports: in.Ports}, nil
}

// UpdateServiceTemplate returns the service template with the given id.
func (r *NilRepo) UpdateServiceTemplate(
	_ context.Context, _ string, id int, in domain.ServiceTemplateWrite,
) (*domain.ServiceTemplate, error) {
	return &domain.ServiceTemplate{ID: id, Name: in.Name, Ports: in.Ports}, nil
}

// DeleteServiceTemplate always succeeds.
func (r *NilRepo) DeleteServiceTemplate(_ context.Context, _ string, _ int) error {
	return nil
}

// ListFHRPGroups returns an empty FHRP group page.
func (r *NilRepo) ListFHRPGroups(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.FHRPGroup], error) {
	return &domain.PaginatedResponse[domain.FHRPGroup]{Count: 0, Results: []domain.FHRPGroup{}}, nil
}

// CreateFHRPGroup returns an FHRP group with id 1.
func (r *NilRepo) CreateFHRPGroup(
	_ context.Context, _ string, in domain.FHRPGroupWrite,
) (*domain.FHRPGroup, error) {
	return &domain.FHRPGroup{ID: 1, Name: in.Name}, nil
}

// UpdateFHRPGroup returns the FHRP group with the given id.
func (r *NilRepo) UpdateFHRPGroup(
	_ context.Context, _ string, id int, in domain.FHRPGroupWrite,
) (*domain.FHRPGroup, error) {
	return &domain.FHRPGroup{ID: id, Name: in.Name}, nil
}

// DeleteFHRPGroup always succeeds.
func (r *NilRepo) DeleteFHRPGroup(_ context.Context, _ string, _ int) error {
	return nil
}

// ListFHRPGroupAssignments returns an empty FHRP group assignment page.
func (r *NilRepo) ListFHRPGroupAssignments(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.FHRPGroupAssignment], error) {
	return &domain.PaginatedResponse[domain.FHRPGroupAssignment]{Count: 0, Results: []domain.FHRPGroupAssignment{}}, nil
}

// CreateFHRPGroupAssignment returns an FHRP group assignment with id 1.
func (r *NilRepo) CreateFHRPGroupAssignment(
	_ context.Context, _ string, in domain.FHRPGroupAssignmentWrite,
) (*domain.FHRPGroupAssignment, error) {
	return &domain.FHRPGroupAssignment{ID: 1, Group: &domain.Nested{ID: in.Group}}, nil
}

// UpdateFHRPGroupAssignment returns the FHRP group assignment with the given id.
func (r *NilRepo) UpdateFHRPGroupAssignment(
	_ context.Context, _ string, id int, in domain.FHRPGroupAssignmentWrite,
) (*domain.FHRPGroupAssignment, error) {
	return &domain.FHRPGroupAssignment{ID: id, Group: &domain.Nested{ID: in.Group}}, nil
}

// DeleteFHRPGroupAssignment always succeeds.
func (r *NilRepo) DeleteFHRPGroupAssignment(_ context.Context, _ string, _ int) error {
	return nil
}

// ListVLANTranslationPolicies returns an empty VLAN translation policy page.
func (r *NilRepo) ListVLANTranslationPolicies(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.VLANTranslationPolicy], error) {
	return &domain.PaginatedResponse[domain.VLANTranslationPolicy]{
		Count: 0, Results: []domain.VLANTranslationPolicy{},
	}, nil
}

// CreateVLANTranslationPolicy returns a VLAN translation policy with id 1.
func (r *NilRepo) CreateVLANTranslationPolicy(
	_ context.Context, _ string, in domain.VLANTranslationPolicyWrite,
) (*domain.VLANTranslationPolicy, error) {
	return &domain.VLANTranslationPolicy{ID: 1, Name: in.Name}, nil
}

// UpdateVLANTranslationPolicy returns the VLAN translation policy with the
// given id.
func (r *NilRepo) UpdateVLANTranslationPolicy(
	_ context.Context, _ string, id int, in domain.VLANTranslationPolicyWrite,
) (*domain.VLANTranslationPolicy, error) {
	return &domain.VLANTranslationPolicy{ID: id, Name: in.Name}, nil
}

// DeleteVLANTranslationPolicy always succeeds.
func (r *NilRepo) DeleteVLANTranslationPolicy(_ context.Context, _ string, _ int) error {
	return nil
}

// ListVLANTranslationRules returns an empty VLAN translation rule page.
func (r *NilRepo) ListVLANTranslationRules(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.VLANTranslationRule], error) {
	return &domain.PaginatedResponse[domain.VLANTranslationRule]{Count: 0, Results: []domain.VLANTranslationRule{}}, nil
}

// defaultLocalVID and defaultRemoteVID are stable sample VIDs returned by the
// NilRepo's VLAN translation rule stubs.
const (
	defaultLocalVID  = 100
	defaultRemoteVID = 200
)

// CreateVLANTranslationRule returns a VLAN translation rule with id 1.
func (r *NilRepo) CreateVLANTranslationRule(
	_ context.Context, _ string, in domain.VLANTranslationRuleWrite,
) (*domain.VLANTranslationRule, error) {
	return &domain.VLANTranslationRule{
		ID: 1, Policy: in.Policy, LocalVID: defaultLocalVID, RemoteVID: defaultRemoteVID,
	}, nil
}

// UpdateVLANTranslationRule returns the VLAN translation rule with the given id.
func (r *NilRepo) UpdateVLANTranslationRule(
	_ context.Context, _ string, id int, in domain.VLANTranslationRuleWrite,
) (*domain.VLANTranslationRule, error) {
	return &domain.VLANTranslationRule{
		ID: id, Policy: in.Policy, LocalVID: defaultLocalVID, RemoteVID: defaultRemoteVID,
	}, nil
}

// DeleteVLANTranslationRule always succeeds.
func (r *NilRepo) DeleteVLANTranslationRule(_ context.Context, _ string, _ int) error {
	return nil
}

// ListContactGroups returns an empty contact group page.
func (r *NilRepo) ListContactGroups(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.ContactGroup], error) {
	return &domain.PaginatedResponse[domain.ContactGroup]{Count: 0, Results: []domain.ContactGroup{}}, nil
}

// CreateContactGroup returns a contact group with id 1.
func (r *NilRepo) CreateContactGroup(
	_ context.Context, _ string, in domain.ContactGroupWrite,
) (*domain.ContactGroup, error) {
	return &domain.ContactGroup{ID: 1, Name: in.Name}, nil
}

// UpdateContactGroup returns the contact group with the given id.
func (r *NilRepo) UpdateContactGroup(
	_ context.Context, _ string, id int, in domain.ContactGroupWrite,
) (*domain.ContactGroup, error) {
	return &domain.ContactGroup{ID: id, Name: in.Name}, nil
}

// DeleteContactGroup always succeeds.
func (r *NilRepo) DeleteContactGroup(_ context.Context, _ string, _ int) error {
	return nil
}

// ListContactRoles returns an empty contact role page.
func (r *NilRepo) ListContactRoles(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.ContactRole], error) {
	return &domain.PaginatedResponse[domain.ContactRole]{Count: 0, Results: []domain.ContactRole{}}, nil
}

// CreateContactRole returns a contact role with id 1.
func (r *NilRepo) CreateContactRole(
	_ context.Context, _ string, in domain.ContactRoleWrite,
) (*domain.ContactRole, error) {
	return &domain.ContactRole{ID: 1, Name: in.Name}, nil
}

// UpdateContactRole returns the contact role with the given id.
func (r *NilRepo) UpdateContactRole(
	_ context.Context, _ string, id int, in domain.ContactRoleWrite,
) (*domain.ContactRole, error) {
	return &domain.ContactRole{ID: id, Name: in.Name}, nil
}

// DeleteContactRole always succeeds.
func (r *NilRepo) DeleteContactRole(_ context.Context, _ string, _ int) error {
	return nil
}

// ListVirtualDisks returns an empty virtual disk page.
func (r *NilRepo) ListVirtualDisks(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.VirtualDisk], error) {
	return &domain.PaginatedResponse[domain.VirtualDisk]{Count: 0, Results: []domain.VirtualDisk{}}, nil
}

// CreateVirtualDisk returns a virtual disk with id 1.
func (r *NilRepo) CreateVirtualDisk(
	_ context.Context, _ string, in domain.VirtualDiskWrite,
) (*domain.VirtualDisk, error) {
	return &domain.VirtualDisk{ID: 1, Name: in.Name}, nil
}

// UpdateVirtualDisk returns the virtual disk with the given id.
func (r *NilRepo) UpdateVirtualDisk(
	_ context.Context, _ string, id int, in domain.VirtualDiskWrite,
) (*domain.VirtualDisk, error) {
	return &domain.VirtualDisk{ID: id, Name: in.Name}, nil
}

// DeleteVirtualDisk always succeeds.
func (r *NilRepo) DeleteVirtualDisk(_ context.Context, _ string, _ int) error {
	return nil
}

// ListVirtualMachineTypes returns an empty virtual machine type page.
func (r *NilRepo) ListVirtualMachineTypes(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.VirtualMachineType], error) {
	return &domain.PaginatedResponse[domain.VirtualMachineType]{Count: 0, Results: []domain.VirtualMachineType{}}, nil
}

// CreateVirtualMachineType returns a virtual machine type with id 1.
func (r *NilRepo) CreateVirtualMachineType(
	_ context.Context, _ string, in domain.VirtualMachineTypeWrite,
) (*domain.VirtualMachineType, error) {
	return &domain.VirtualMachineType{ID: 1, Name: in.Name}, nil
}

// UpdateVirtualMachineType returns the virtual machine type with the given id.
func (r *NilRepo) UpdateVirtualMachineType(
	_ context.Context, _ string, id int, in domain.VirtualMachineTypeWrite,
) (*domain.VirtualMachineType, error) {
	return &domain.VirtualMachineType{ID: id, Name: in.Name}, nil
}

// DeleteVirtualMachineType always succeeds.
func (r *NilRepo) DeleteVirtualMachineType(_ context.Context, _ string, _ int) error {
	return nil
}

// ListProviderAccounts returns an empty provider account page.
func (r *NilRepo) ListProviderAccounts(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.ProviderAccount], error) {
	return &domain.PaginatedResponse[domain.ProviderAccount]{Count: 0, Results: []domain.ProviderAccount{}}, nil
}

// CreateProviderAccount returns a provider account with id 1.
func (r *NilRepo) CreateProviderAccount(
	_ context.Context, _ string, in domain.ProviderAccountWrite,
) (*domain.ProviderAccount, error) {
	return &domain.ProviderAccount{ID: 1, Name: in.Name}, nil
}

// UpdateProviderAccount returns the provider account with the given id.
func (r *NilRepo) UpdateProviderAccount(
	_ context.Context, _ string, id int, in domain.ProviderAccountWrite,
) (*domain.ProviderAccount, error) {
	return &domain.ProviderAccount{ID: id, Name: in.Name}, nil
}

// DeleteProviderAccount always succeeds.
func (r *NilRepo) DeleteProviderAccount(_ context.Context, _ string, _ int) error {
	return nil
}

// ListProviderNetworks returns an empty provider network page.
func (r *NilRepo) ListProviderNetworks(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.ProviderNetwork], error) {
	return &domain.PaginatedResponse[domain.ProviderNetwork]{Count: 0, Results: []domain.ProviderNetwork{}}, nil
}

// CreateProviderNetwork returns a provider network with id 1.
func (r *NilRepo) CreateProviderNetwork(
	_ context.Context, _ string, in domain.ProviderNetworkWrite,
) (*domain.ProviderNetwork, error) {
	return &domain.ProviderNetwork{ID: 1, Name: in.Name}, nil
}

// UpdateProviderNetwork returns the provider network with the given id.
func (r *NilRepo) UpdateProviderNetwork(
	_ context.Context, _ string, id int, in domain.ProviderNetworkWrite,
) (*domain.ProviderNetwork, error) {
	return &domain.ProviderNetwork{ID: id, Name: in.Name}, nil
}

// DeleteProviderNetwork always succeeds.
func (r *NilRepo) DeleteProviderNetwork(_ context.Context, _ string, _ int) error {
	return nil
}

// ListCircuitGroups returns an empty circuit group page.
func (r *NilRepo) ListCircuitGroups(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.CircuitGroup], error) {
	return &domain.PaginatedResponse[domain.CircuitGroup]{Count: 0, Results: []domain.CircuitGroup{}}, nil
}

// CreateCircuitGroup returns a circuit group with id 1.
func (r *NilRepo) CreateCircuitGroup(
	_ context.Context, _ string, in domain.CircuitGroupWrite,
) (*domain.CircuitGroup, error) {
	return &domain.CircuitGroup{ID: 1, Name: in.Name}, nil
}

// UpdateCircuitGroup returns the circuit group with the given id.
func (r *NilRepo) UpdateCircuitGroup(
	_ context.Context, _ string, id int, in domain.CircuitGroupWrite,
) (*domain.CircuitGroup, error) {
	return &domain.CircuitGroup{ID: id, Name: in.Name}, nil
}

// DeleteCircuitGroup always succeeds.
func (r *NilRepo) DeleteCircuitGroup(_ context.Context, _ string, _ int) error {
	return nil
}

// ListCircuitGroupAssignments returns an empty circuit group assignment page.
func (r *NilRepo) ListCircuitGroupAssignments(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.CircuitGroupAssignment], error) {
	return &domain.PaginatedResponse[domain.CircuitGroupAssignment]{
		Count: 0, Results: []domain.CircuitGroupAssignment{},
	}, nil
}

// CreateCircuitGroupAssignment returns a circuit group assignment with id 1.
func (r *NilRepo) CreateCircuitGroupAssignment(
	_ context.Context, _ string, in domain.CircuitGroupAssignmentWrite,
) (*domain.CircuitGroupAssignment, error) {
	return &domain.CircuitGroupAssignment{ID: 1, Group: &domain.Nested{ID: in.Group}}, nil
}

// UpdateCircuitGroupAssignment returns the circuit group assignment with the
// given id.
func (r *NilRepo) UpdateCircuitGroupAssignment(
	_ context.Context, _ string, id int, _ domain.CircuitGroupAssignmentWrite,
) (*domain.CircuitGroupAssignment, error) {
	return &domain.CircuitGroupAssignment{ID: id}, nil
}

// DeleteCircuitGroupAssignment always succeeds.
func (r *NilRepo) DeleteCircuitGroupAssignment(_ context.Context, _ string, _ int) error {
	return nil
}

// ListVirtualCircuits returns an empty virtual circuit page.
func (r *NilRepo) ListVirtualCircuits(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.VirtualCircuit], error) {
	return &domain.PaginatedResponse[domain.VirtualCircuit]{Count: 0, Results: []domain.VirtualCircuit{}}, nil
}

// CreateVirtualCircuit returns a virtual circuit with id 1.
func (r *NilRepo) CreateVirtualCircuit(
	_ context.Context, _ string, in domain.VirtualCircuitWrite,
) (*domain.VirtualCircuit, error) {
	return &domain.VirtualCircuit{ID: 1, CID: in.CID}, nil
}

// UpdateVirtualCircuit returns the virtual circuit with the given id.
func (r *NilRepo) UpdateVirtualCircuit(
	_ context.Context, _ string, id int, in domain.VirtualCircuitWrite,
) (*domain.VirtualCircuit, error) {
	return &domain.VirtualCircuit{ID: id, CID: in.CID}, nil
}

// DeleteVirtualCircuit always succeeds.
func (r *NilRepo) DeleteVirtualCircuit(_ context.Context, _ string, _ int) error {
	return nil
}

// ListVirtualCircuitTypes returns an empty virtual circuit type page.
func (r *NilRepo) ListVirtualCircuitTypes(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.VirtualCircuitType], error) {
	return &domain.PaginatedResponse[domain.VirtualCircuitType]{Count: 0, Results: []domain.VirtualCircuitType{}}, nil
}

// CreateVirtualCircuitType returns a virtual circuit type with id 1.
func (r *NilRepo) CreateVirtualCircuitType(
	_ context.Context, _ string, in domain.VirtualCircuitTypeWrite,
) (*domain.VirtualCircuitType, error) {
	return &domain.VirtualCircuitType{ID: 1, Name: in.Name}, nil
}

// UpdateVirtualCircuitType returns the virtual circuit type with the given id.
func (r *NilRepo) UpdateVirtualCircuitType(
	_ context.Context, _ string, id int, in domain.VirtualCircuitTypeWrite,
) (*domain.VirtualCircuitType, error) {
	return &domain.VirtualCircuitType{ID: id, Name: in.Name}, nil
}

// DeleteVirtualCircuitType always succeeds.
func (r *NilRepo) DeleteVirtualCircuitType(_ context.Context, _ string, _ int) error {
	return nil
}

// ListVirtualCircuitTerminations returns an empty virtual circuit termination
// page.
func (r *NilRepo) ListVirtualCircuitTerminations(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.VirtualCircuitTermination], error) {
	return &domain.PaginatedResponse[domain.VirtualCircuitTermination]{
		Count: 0, Results: []domain.VirtualCircuitTermination{},
	}, nil
}

// CreateVirtualCircuitTermination returns a virtual circuit termination with id
// 1.
func (r *NilRepo) CreateVirtualCircuitTermination(
	_ context.Context, _ string, in domain.VirtualCircuitTerminationWrite,
) (*domain.VirtualCircuitTermination, error) {
	return &domain.VirtualCircuitTermination{ID: 1, VirtualCircuit: &domain.Nested{ID: in.VirtualCircuit}}, nil
}

// UpdateVirtualCircuitTermination returns the virtual circuit termination with
// the given id.
func (r *NilRepo) UpdateVirtualCircuitTermination(
	_ context.Context, _ string, id int, _ domain.VirtualCircuitTerminationWrite,
) (*domain.VirtualCircuitTermination, error) {
	return &domain.VirtualCircuitTermination{ID: id}, nil
}

// DeleteVirtualCircuitTermination always succeeds.
func (r *NilRepo) DeleteVirtualCircuitTermination(_ context.Context, _ string, _ int) error {
	return nil
}

// ListRoles returns an empty role page.
func (r *NilRepo) ListRoles(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Role], error) {
	return &domain.PaginatedResponse[domain.Role]{Count: 0, Results: []domain.Role{}}, nil
}

// ListTenants returns an empty tenant page.
func (r *NilRepo) ListTenants(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Tenant], error) {
	return &domain.PaginatedResponse[domain.Tenant]{Count: 0, Results: []domain.Tenant{}}, nil
}

// ListTags returns an empty tag page.
func (r *NilRepo) ListTags(
	_ context.Context, _ string, _ map[string]string,
) (*domain.PaginatedResponse[domain.Tag], error) {
	return &domain.PaginatedResponse[domain.Tag]{Count: 0, Results: []domain.Tag{}}, nil
}

// CreateTag returns a tag with id 1.
func (r *NilRepo) CreateTag(_ context.Context, _ string, in domain.TagWrite) (*domain.Tag, error) {
	return &domain.Tag{ID: 1, Name: in.Name}, nil
}

// UpdateTag returns the tag with the given id.
func (r *NilRepo) UpdateTag(_ context.Context, _ string, id int, in domain.TagWrite) (*domain.Tag, error) {
	return &domain.Tag{ID: id, Name: in.Name}, nil
}

// DeleteTag always succeeds.
func (r *NilRepo) DeleteTag(_ context.Context, _ string, _ int) error {
	return nil
}
