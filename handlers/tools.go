package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/domain"
)

// — input types —

// SitesInput represents the input fields for the get_sites tool.
type SitesInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Region   string `json:"region,omitempty" jsonschema:"filter by region (slug or name)"`
	Status   string `json:"status,omitempty" jsonschema:"status: active, planned, staged, retired, decommissioning"`
	Tenant   string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug or name)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// RegionsInput represents the input fields for the get_regions tool.
type RegionsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Parent   string `json:"parent,omitempty" jsonschema:"filter by parent region ID"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// TenantGroupsInput represents the input fields for the get_tenant_groups tool.
type TenantGroupsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Parent   string `json:"parent,omitempty" jsonschema:"filter by parent tenant group ID"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// RolesInput represents the input fields for the get_roles tool.
type RolesInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// TenantsInput represents the input fields for the get_tenants tool.
type TenantsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Group    string `json:"group,omitempty" jsonschema:"filter by tenant group (ID)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// TagsInput represents the input fields for the get_tags tool.
type TagsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// SiteGroupsInput represents the input fields for the get_site_groups tool.
type SiteGroupsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Parent   string `json:"parent,omitempty" jsonschema:"filter by parent site group ID"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// RackGroupsInput represents the input fields for the get_rack_groups tool.
type RackGroupsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Site     string `json:"site,omitempty" jsonschema:"filter by site (slug)"`
	Parent   string `json:"parent,omitempty" jsonschema:"filter by parent rack group ID"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// RackRolesInput represents the input fields for the get_rack_roles tool.
type RackRolesInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Color    string `json:"color,omitempty" jsonschema:"filter by color (hex value)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// DeviceRolesInput represents the input fields for the get_device_roles tool.
type DeviceRolesInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Color    string `json:"color,omitempty" jsonschema:"filter by color (hex value)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// PlatformsInput represents the input fields for the get_platforms tool.
type PlatformsInput struct {
	Q            string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Manufacturer string `json:"manufacturer,omitempty" jsonschema:"filter by manufacturer (slug)"`
	Tag          string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page         int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize     int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// VirtualChassisInput represents the input fields for the get_virtual_chassis
// tool.
type VirtualChassisInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Name     string `json:"name,omitempty" jsonschema:"filter by name (case-insensitive partial match)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// PowerPanelsInput represents the input fields for the get_power_panels tool.
type PowerPanelsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Site     string `json:"site,omitempty" jsonschema:"filter by site (slug)"`
	Location string `json:"location,omitempty" jsonschema:"filter by location (slug)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// InventoryItemRolesInput represents the input fields for the
// get_inventory_item_roles tool.
type InventoryItemRolesInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Color    string `json:"color,omitempty" jsonschema:"filter by color (hex value)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// DevicesInput represents the input fields for the get_devices tool.
type DevicesInput struct {
	Q            string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Site         string `json:"site,omitempty" jsonschema:"filter by site (slug)"`
	Role         string `json:"role,omitempty" jsonschema:"filter by device role (slug)"`
	Manufacturer string `json:"manufacturer,omitempty" jsonschema:"filter by manufacturer (slug)"`
	DeviceType   string `json:"device_type,omitempty" jsonschema:"filter by device type slug (e.g. c-1250)"`
	Status       string `json:"status,omitempty" jsonschema:"status: active, offline, planned, staged, failed, inventory, decommissioning"` //nolint:lll // struct tag cannot be split across lines
	Name         string `json:"name,omitempty" jsonschema:"filter by name (case-insensitive partial match)"`
	Tenant       string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug)"`
	Rack         string `json:"rack,omitempty" jsonschema:"filter by rack (name)"`
	Cluster      string `json:"cluster,omitempty" jsonschema:"filter by cluster (name)"`
	Tag          string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page         int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize     int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// IPAddressesInput represents the input fields for the get_ip_addresses tool.
type IPAddressesInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Address  string `json:"address,omitempty" jsonschema:"filter by address (e.g. 192.168.1.0/24)"`
	Device   string `json:"device,omitempty" jsonschema:"filter by assigned device name"`
	Status   string `json:"status,omitempty" jsonschema:"status: active, reserved, deprecated, dhcp, slaac"`
	VRF      string `json:"vrf,omitempty" jsonschema:"filter by VRF (rd or name)"`
	Role     string `json:"role,omitempty" jsonschema:"filter by role (loopback, etc.)"`
	Tenant   string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// PrefixesInput represents the input fields for the get_prefixes tool.
type PrefixesInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Prefix   string `json:"prefix,omitempty" jsonschema:"filter by prefix (e.g. 10.0.0.0/8)"`
	Site     string `json:"site,omitempty" jsonschema:"filter by site (slug)"`
	VRF      string `json:"vrf,omitempty" jsonschema:"filter by VRF (rd or name)"`
	Status   string `json:"status,omitempty" jsonschema:"status: active, container, reserved, deprecated"`
	Role     string `json:"role,omitempty" jsonschema:"filter by role (slug)"`
	Tenant   string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Within   string `json:"within,omitempty" jsonschema:"find prefixes within a given prefix"`
	Family   *int   `json:"family,omitempty" jsonschema:"address family: 4 or 6"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// VLANsInput represents the input fields for the get_vlans tool.
type VLANsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Site     string `json:"site,omitempty" jsonschema:"filter by site (slug)"`
	Group    string `json:"group,omitempty" jsonschema:"filter by VLAN group (slug)"`
	Status   string `json:"status,omitempty" jsonschema:"status: active, reserved, deprecated"`
	Tenant   string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	VID      int    `json:"vid,omitempty" jsonschema:"filter by VLAN ID"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// VirtualMachinesInput represents the input fields for the get_virtual_machines tool.
type VirtualMachinesInput struct {
	Q            string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Cluster      string `json:"cluster,omitempty" jsonschema:"filter by cluster (name)"`
	ClusterGroup string `json:"cluster_group,omitempty" jsonschema:"filter by cluster group (slug)"`
	Role         string `json:"role,omitempty" jsonschema:"filter by VM role (slug)"`
	Status       string `json:"status,omitempty" jsonschema:"status: active, staged, offline, decommissioning"`
	Tenant       string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug)"`
	Name         string `json:"name,omitempty" jsonschema:"filter by name (case-insensitive partial match)"`
	Site         string `json:"site,omitempty" jsonschema:"filter by site (slug)"`
	Tag          string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page         int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize     int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// ClustersInput represents the input fields for the get_clusters tool.
type ClustersInput struct {
	Q            string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	ClusterType  string `json:"cluster_type,omitempty" jsonschema:"filter by cluster type (slug)"`
	ClusterGroup string `json:"cluster_group,omitempty" jsonschema:"filter by cluster group (slug)"`
	Site         string `json:"site,omitempty" jsonschema:"filter by site (slug)"`
	Tenant       string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug)"`
	Name         string `json:"name,omitempty" jsonschema:"filter by name (case-insensitive partial match)"`
	Tag          string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page         int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize     int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// CircuitsInput represents the input fields for the get_circuits tool.
type CircuitsInput struct {
	Q           string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Provider    string `json:"provider,omitempty" jsonschema:"filter by provider (slug)"`
	CircuitType string `json:"circuit_type,omitempty" jsonschema:"filter by circuit type (slug)"`
	Site        string `json:"site,omitempty" jsonschema:"filter by site (slug)"`
	Status      string `json:"status,omitempty" jsonschema:"status: active, planned, offline, decommissioning"`
	Tenant      string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug)"`
	Tag         string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page        int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize    int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// GetObjectInput represents the input fields for the get_object_by_id tool.
type GetObjectInput struct {
	ObjectType string            `json:"object_type" jsonschema:"object type: site, region, device, prefix, ip_address, vlan, virtual_machine, cluster, circuit, provider, tenant, rack, manufacturer, device_type, location, cluster_type, cluster_group, circuit_type, vrf, vlan_group, role, contact, cable, interface, vm_interface, circuit_termination, site_group, rack_group, rack_role, device_role, platform, virtual_chassis, power_panel, inventory_item_role, front_port, rear_port, console_port, console_server_port, power_port, power_outlet, power_feed, device_bay, rack_reservation, inventory_item, mac_address, rir, aggregate, ip_range, asn, asn_range, route_target, service, service_template, fhrp_group, fhrp_group_assignment, vlan_translation_policy, vlan_translation_rule, contact_group, contact_role, virtual_disk, virtual_machine_type, provider_account, provider_network, circuit_group, circuit_group_assignment, virtual_circuit, virtual_circuit_type, virtual_circuit_termination,required"` //nolint:lll // struct tag cannot be split across lines
	ID         int               `json:"id" jsonschema:"numeric ID of the object (positive integer),required"`
	Params     map[string]string `json:"params,omitempty" jsonschema:"additional query parameters to pass to NetBox (optional)"` //nolint:lll // struct tag cannot be split across lines
}

// InterfacesInput represents the input fields for the get_interfaces tool.
type InterfacesInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Device   string `json:"device,omitempty" jsonschema:"filter by device (name)"`
	Type     string `json:"type,omitempty" jsonschema:"filter by interface type (slug)"`
	Enabled  *bool  `json:"enabled,omitempty" jsonschema:"filter by enabled status"`
	Name     string `json:"name,omitempty" jsonschema:"filter by name (case-insensitive partial match)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// VMInterfacesInput represents the input fields for the get_vm_interfaces tool.
type VMInterfacesInput struct {
	Q              string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	VirtualMachine string `json:"virtual_machine,omitempty" jsonschema:"filter by virtual machine (name)"`
	Name           string `json:"name,omitempty" jsonschema:"filter by name (case-insensitive partial match)"`
	Tag            string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page           int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize       int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// CircuitTerminationsInput represents the input fields for the get_circuit_terminations tool.
type CircuitTerminationsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Circuit  string `json:"circuit,omitempty" jsonschema:"filter by circuit (ID or CID)"`
	Site     string `json:"site,omitempty" jsonschema:"filter by site (slug)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	TermSide string `json:"term_side,omitempty" jsonschema:"filter by termination side: A or Z"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// CablesInput represents the input fields for the get_cables tool.
type CablesInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Type     string `json:"type,omitempty" jsonschema:"filter by cable type (slug)"`
	Status   string `json:"status,omitempty" jsonschema:"filter by cable status: connected, planned, decommissioning"`
	Site     string `json:"site,omitempty" jsonschema:"filter by site (slug)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Color    string `json:"color,omitempty" jsonschema:"filter by color (slug)"`
	Label    string `json:"label,omitempty" jsonschema:"filter by label (case-insensitive partial match)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// RacksInput represents the input fields for the get_racks tool.
type RacksInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Site     string `json:"site,omitempty" jsonschema:"filter by site (slug)"`
	Location string `json:"location,omitempty" jsonschema:"filter by location (slug)"`
	Status   string `json:"status,omitempty" jsonschema:"status: active, planned, reserved, available"`
	Tenant   string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// SiteCreateInput represents the writable fields for creating a site.
type SiteCreateInput struct {
	Name            string         `json:"name" jsonschema:"site name (required)"`
	Slug            string         `json:"slug,omitempty" jsonschema:"URL-friendly unique identifier (defaults to a slugified name)"` //nolint:lll // struct tag cannot be split across lines
	Status          *string        `json:"status,omitempty" jsonschema:"status: planned, staged, active, decommissioning, retired"`   //nolint:lll // struct tag cannot be split across lines   //nolint:lll // struct tag cannot be split across lines
	Region          *int           `json:"region,omitempty" jsonschema:"region ID (from tenancy)"`
	Tenant          *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Facility        *string        `json:"facility,omitempty" jsonschema:"local facility ID or description"`
	TimeZone        *string        `json:"time_zone,omitempty" jsonschema:"IANA time zone, e.g. America/New_York"`
	Description     *string        `json:"description,omitempty" jsonschema:"short description"`
	PhysicalAddress *string        `json:"physical_address,omitempty" jsonschema:"physical address"`
	ShippingAddress *string        `json:"shipping_address,omitempty" jsonschema:"shipping address"`
	Comments        *string        `json:"comments,omitempty" jsonschema:"free-form comments"`
	Tags            []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields    map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// SiteUpdateInput represents the writable fields for updating a site. All fields
// are optional; only the explicitly provided ones are patched (PATCH merge).
type SiteUpdateInput struct {
	ID              int            `json:"id" jsonschema:"numeric ID of the site to update,required"`
	Name            *string        `json:"name,omitempty" jsonschema:"site name"`
	Slug            *string        `json:"slug,omitempty" jsonschema:"URL-friendly unique identifier"`
	Status          *string        `json:"status,omitempty" jsonschema:"status: planned, staged, active, decommissioning, retired"` //nolint:lll // struct tag cannot be split across lines
	Region          *int           `json:"region,omitempty" jsonschema:"region ID (from tenancy)"`
	Tenant          *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Facility        *string        `json:"facility,omitempty" jsonschema:"local facility ID or description"`
	TimeZone        *string        `json:"time_zone,omitempty" jsonschema:"IANA time zone, e.g. America/New_York"`
	Description     *string        `json:"description,omitempty" jsonschema:"short description"`
	PhysicalAddress *string        `json:"physical_address,omitempty" jsonschema:"physical address"`
	ShippingAddress *string        `json:"shipping_address,omitempty" jsonschema:"shipping address"`
	Comments        *string        `json:"comments,omitempty" jsonschema:"free-form comments"`
	Tags            []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields    map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// SiteDeleteInput represents the input fields for deleting a site.
type SiteDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the site to delete,required"`
}

// RegionCreateInput represents the writable fields for creating a region.
type RegionCreateInput struct {
	Name         string         `json:"name" jsonschema:"region name (required)"`
	Slug         string         `json:"slug,omitempty" jsonschema:"URL-friendly unique identifier (defaults to a slugified name)"` //nolint:lll // struct tag cannot be split across lines
	Parent       *int           `json:"parent,omitempty" jsonschema:"parent region ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// RegionUpdateInput represents the writable fields for updating a region. All fields
// are optional; only the explicitly provided ones are patched (PATCH merge).
type RegionUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the region to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"region name"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly unique identifier"`
	Parent       *int           `json:"parent,omitempty" jsonschema:"parent region ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// RegionDeleteInput represents the input fields for deleting a region.
type RegionDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the region to delete,required"`
}

// TenantGroupCreateInput represents the writable fields for creating a tenant
// group.
type TenantGroupCreateInput struct {
	Name         string         `json:"name" jsonschema:"tenant group name (required)"`
	Slug         string         `json:"slug,omitempty" jsonschema:"URL-friendly unique identifier (defaults to a slugified name)"` //nolint:lll // struct tag cannot be split across lines
	Parent       *int           `json:"parent,omitempty" jsonschema:"parent tenant group ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// TenantGroupUpdateInput represents the writable fields for updating a tenant
// group. All fields are optional; only the explicitly provided ones are
// patched (PATCH merge).
type TenantGroupUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the tenant group to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"tenant group name"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly unique identifier"`
	Parent       *int           `json:"parent,omitempty" jsonschema:"parent tenant group ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// TenantGroupDeleteInput represents the input fields for deleting a tenant
// group.
type TenantGroupDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the tenant group to delete,required"`
}

// SiteGroupCreateInput represents the writable fields for creating a site
// group.
type SiteGroupCreateInput struct {
	Name         string         `json:"name" jsonschema:"site group name (required)"`
	Slug         string         `json:"slug,omitempty" jsonschema:"URL-friendly unique identifier (defaults to a slugified name)"` //nolint:lll // struct tag cannot be split across lines
	Parent       *int           `json:"parent,omitempty" jsonschema:"parent site group ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// SiteGroupUpdateInput represents the writable fields for updating a site
// group. All fields are optional; only the explicitly provided ones are patched
// (PATCH merge).
type SiteGroupUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the site group to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"site group name"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly unique identifier"`
	Parent       *int           `json:"parent,omitempty" jsonschema:"parent site group ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// SiteGroupDeleteInput represents the input fields for deleting a site group.
type SiteGroupDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the site group to delete,required"`
}

// RackGroupCreateInput represents the writable fields for creating a rack
// group.
type RackGroupCreateInput struct {
	Name        string   `json:"name" jsonschema:"rack group name (required)"`
	Slug        string   `json:"slug,omitempty" jsonschema:"URL-friendly unique identifier (defaults to a slugified name)"` //nolint:lll // struct tag cannot be split across lines
	Site        *int     `json:"site,omitempty" jsonschema:"site ID"`
	Parent      *int     `json:"parent,omitempty" jsonschema:"parent rack group ID"`
	Description *string  `json:"description,omitempty" jsonschema:"short description"`
	Tags        []string `json:"tags,omitempty" jsonschema:"list of tag names"`
}

// RackGroupUpdateInput represents the writable fields for updating a rack
// group. All fields are optional; only the explicitly provided ones are patched
// (PATCH merge).
type RackGroupUpdateInput struct {
	ID          int      `json:"id" jsonschema:"numeric ID of the rack group to update,required"`
	Name        *string  `json:"name,omitempty" jsonschema:"rack group name"`
	Slug        *string  `json:"slug,omitempty" jsonschema:"URL-friendly unique identifier"`
	Site        *int     `json:"site,omitempty" jsonschema:"site ID"`
	Parent      *int     `json:"parent,omitempty" jsonschema:"parent rack group ID"`
	Description *string  `json:"description,omitempty" jsonschema:"short description"`
	Tags        []string `json:"tags,omitempty" jsonschema:"list of tag names"`
}

// RackGroupDeleteInput represents the input fields for deleting a rack group.
type RackGroupDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the rack group to delete,required"`
}

// RackRoleCreateInput represents the writable fields for creating a rack role.
type RackRoleCreateInput struct {
	Name         string         `json:"name" jsonschema:"rack role name (required)"`
	Slug         string         `json:"slug,omitempty" jsonschema:"URL-friendly unique identifier (defaults to a slugified name)"` //nolint:lll // struct tag cannot be split across lines
	Color        *string        `json:"color,omitempty" jsonschema:"hex color, e.g. 00ff00"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// RackRoleUpdateInput represents the writable fields for updating a rack role.
// All fields are optional; only the explicitly provided ones are patched (PATCH
// merge).
type RackRoleUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the rack role to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"rack role name"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly unique identifier"`
	Color        *string        `json:"color,omitempty" jsonschema:"hex color, e.g. 00ff00"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// RackRoleDeleteInput represents the input fields for deleting a rack role.
type RackRoleDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the rack role to delete,required"`
}

// DeviceRoleCreateInput represents the writable fields for creating a device
// role.
type DeviceRoleCreateInput struct {
	Name           string         `json:"name" jsonschema:"device role name (required)"`
	Slug           string         `json:"slug,omitempty" jsonschema:"URL-friendly unique identifier (defaults to a slugified name)"` //nolint:lll // struct tag cannot be split across lines
	Color          *string        `json:"color,omitempty" jsonschema:"hex color, e.g. 00ff00"`
	VMRole         *bool          `json:"vm_role,omitempty" jsonschema:"whether this role is a VM role"`
	ConfigTemplate *int           `json:"config_template,omitempty" jsonschema:"config template ID"`
	Description    *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags           []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields   map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
	Weight         *int           `json:"weight,omitempty" jsonschema:"sorting weight"`
}

// DeviceRoleUpdateInput represents the writable fields for updating a device
// role. All fields are optional; only the explicitly provided ones are patched
// (PATCH merge).
type DeviceRoleUpdateInput struct {
	ID             int            `json:"id" jsonschema:"numeric ID of the device role to update,required"`
	Name           *string        `json:"name,omitempty" jsonschema:"device role name"`
	Slug           *string        `json:"slug,omitempty" jsonschema:"URL-friendly unique identifier"`
	Color          *string        `json:"color,omitempty" jsonschema:"hex color, e.g. 00ff00"`
	VMRole         *bool          `json:"vm_role,omitempty" jsonschema:"whether this role is a VM role"`
	ConfigTemplate *int           `json:"config_template,omitempty" jsonschema:"config template ID"`
	Description    *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags           []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields   map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
	Weight         *int           `json:"weight,omitempty" jsonschema:"sorting weight"`
}

// DeviceRoleDeleteInput represents the input fields for deleting a device role.
type DeviceRoleDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the device role to delete,required"`
}

// PlatformCreateInput represents the writable fields for creating a platform.
type PlatformCreateInput struct {
	Name           string         `json:"name" jsonschema:"platform name (required)"`
	Slug           string         `json:"slug,omitempty" jsonschema:"URL-friendly unique identifier (defaults to a slugified name)"` //nolint:lll // struct tag cannot be split across lines
	Manufacturer   *int           `json:"manufacturer,omitempty" jsonschema:"manufacturer ID (look up first)"`
	ConfigTemplate *int           `json:"config_template,omitempty" jsonschema:"config template ID"`
	Description    *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags           []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields   map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// PlatformUpdateInput represents the writable fields for updating a platform.
// All fields are optional; only the explicitly provided ones are patched (PATCH
// merge).
type PlatformUpdateInput struct {
	ID             int            `json:"id" jsonschema:"numeric ID of the platform to update,required"`
	Name           *string        `json:"name,omitempty" jsonschema:"platform name"`
	Slug           *string        `json:"slug,omitempty" jsonschema:"URL-friendly unique identifier"`
	Manufacturer   *int           `json:"manufacturer,omitempty" jsonschema:"manufacturer ID (look up first)"`
	ConfigTemplate *int           `json:"config_template,omitempty" jsonschema:"config template ID"`
	Description    *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags           []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields   map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// PlatformDeleteInput represents the input fields for deleting a platform.
type PlatformDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the platform to delete,required"`
}

// VirtualChassisCreateInput represents the writable fields for creating a
// virtual chassis. Name is optional on create for this entity.
type VirtualChassisCreateInput struct {
	Name         string         `json:"name,omitempty" jsonschema:"virtual chassis name (optional)"`
	Domain       *string        `json:"domain,omitempty" jsonschema:"domain name"`
	Master       *int           `json:"master,omitempty" jsonschema:"master device ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VirtualChassisUpdateInput represents the writable fields for updating a
// virtual chassis. All fields are optional; only the explicitly provided ones
// are patched (PATCH merge).
type VirtualChassisUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the virtual chassis to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"virtual chassis name"`
	Domain       *string        `json:"domain,omitempty" jsonschema:"domain name"`
	Master       *int           `json:"master,omitempty" jsonschema:"master device ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VirtualChassisDeleteInput represents the input fields for deleting a virtual
// chassis.
type VirtualChassisDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the virtual chassis to delete,required"`
}

// PowerPanelCreateInput represents the writable fields for creating a power
// panel.
type PowerPanelCreateInput struct {
	Site         *int           `json:"site,omitempty" jsonschema:"site ID"`
	Location     *int           `json:"location,omitempty" jsonschema:"location ID"`
	Name         string         `json:"name" jsonschema:"power panel name (required)"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// PowerPanelUpdateInput represents the writable fields for updating a power
// panel. All fields are optional; only the explicitly provided ones are patched
// (PATCH merge).
type PowerPanelUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the power panel to update,required"`
	Site         *int           `json:"site,omitempty" jsonschema:"site ID"`
	Location     *int           `json:"location,omitempty" jsonschema:"location ID"`
	Name         *string        `json:"name,omitempty" jsonschema:"power panel name"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// PowerPanelDeleteInput represents the input fields for deleting a power panel.
type PowerPanelDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the power panel to delete,required"`
}

// InventoryItemRoleCreateInput represents the writable fields for creating an
// inventory item role.
type InventoryItemRoleCreateInput struct {
	Name         string         `json:"name" jsonschema:"inventory item role name (required)"`
	Slug         string         `json:"slug,omitempty" jsonschema:"URL-friendly unique identifier (defaults to a slugified name)"` //nolint:lll // struct tag cannot be split across lines
	Color        *string        `json:"color,omitempty" jsonschema:"hex color, e.g. 00ff00"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// InventoryItemRoleUpdateInput represents the writable fields for updating an
// inventory item role. All fields are optional; only the explicitly provided
// ones are patched (PATCH merge).
type InventoryItemRoleUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the inventory item role to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"inventory item role name"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly unique identifier"`
	Color        *string        `json:"color,omitempty" jsonschema:"hex color, e.g. 00ff00"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// InventoryItemRoleDeleteInput represents the input fields for deleting an
// inventory item role.
type InventoryItemRoleDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the inventory item role to delete,required"`
}

// DeviceCreateInput represents the writable fields for creating a device.
type DeviceCreateInput struct {
	Name         string         `json:"name" jsonschema:"device name (required)"`
	DeviceType   *int           `json:"device_type,omitempty" jsonschema:"device type ID (look up first)"`
	Role         *int           `json:"role,omitempty" jsonschema:"device role ID (look up first)"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Platform     *int           `json:"platform,omitempty" jsonschema:"platform ID"`
	Serial       *string        `json:"serial,omitempty" jsonschema:"serial number"`
	AssetTag     *string        `json:"asset_tag,omitempty" jsonschema:"asset tag"`
	Site         *int           `json:"site,omitempty" jsonschema:"site ID"`
	Rack         *int           `json:"rack,omitempty" jsonschema:"rack ID"`
	Position     *float64       `json:"position,omitempty" jsonschema:"position within the rack (U height)"`
	Face         *string        `json:"face,omitempty" jsonschema:"rack face: front or rear"`
	Status       *string        `json:"status,omitempty" jsonschema:"status: offline, active, planned, staged, failed, inventory, decommissioning"` //nolint:lll // struct tag cannot be split across lines
	Cluster      *int           `json:"cluster,omitempty" jsonschema:"cluster ID"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"free-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// DeviceUpdateInput represents the writable fields for updating a device. All
// fields are optional; only the explicitly provided ones are patched (PATCH
// merge).
type DeviceUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the device to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"device name"`
	DeviceType   *int           `json:"device_type,omitempty" jsonschema:"device type ID"`
	Role         *int           `json:"role,omitempty" jsonschema:"device role ID"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Platform     *int           `json:"platform,omitempty" jsonschema:"platform ID"`
	Serial       *string        `json:"serial,omitempty" jsonschema:"serial number"`
	AssetTag     *string        `json:"asset_tag,omitempty" jsonschema:"asset tag"`
	Site         *int           `json:"site,omitempty" jsonschema:"site ID"`
	Rack         *int           `json:"rack,omitempty" jsonschema:"rack ID"`
	Position     *float64       `json:"position,omitempty" jsonschema:"position within the rack (U height)"`
	Face         *string        `json:"face,omitempty" jsonschema:"rack face: front or rear"`
	Status       *string        `json:"status,omitempty" jsonschema:"status: offline, active, planned, staged, failed, inventory, decommissioning"` //nolint:lll // struct tag cannot be split across lines
	Cluster      *int           `json:"cluster,omitempty" jsonschema:"cluster ID"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"free-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// DeviceDeleteInput represents the input fields for deleting a device.
type DeviceDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the device to delete,required"`
}

// IPAddressCreateInput represents the writable fields for creating an IP address.
type IPAddressCreateInput struct {
	Address            string         `json:"address" jsonschema:"IP address or prefix, e.g. 192.168.1.1/24 (required)"`
	Status             *string        `json:"status,omitempty" jsonschema:"status: active, reserved, deprecated, dhcp, slaac"`             //nolint:lll // struct tag cannot be split across lines
	Role               *string        `json:"role,omitempty" jsonschema:"role: loopback, secondary, anycast, vip, vrrp, hsrp, glbp, carp"` //nolint:lll // struct tag cannot be split across lines
	VRF                *int           `json:"vrf,omitempty" jsonschema:"VRF ID"`
	Tenant             *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	DNSName            *string        `json:"dns_name,omitempty" jsonschema:"hostname or FQDN"`
	Description        *string        `json:"description,omitempty" jsonschema:"short description"`
	AssignedObjectType *string        `json:"assigned_object_type,omitempty" jsonschema:"assigned object type, e.g. dcim.interface (use with assigned_object_id)"` //nolint:lll // struct tag cannot be split across lines
	AssignedObjectID   *int           `json:"assigned_object_id,omitempty" jsonschema:"ID of the assigned object (use with assigned_object_type)"`                 //nolint:lll // struct tag cannot be split across lines
	Comments           *string        `json:"comments,omitempty" jsonschema:"free-form comments"`
	Tags               []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields       map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// IPAddressUpdateInput represents the writable fields for updating an IP
// address. All fields are optional; only the explicitly provided ones are
// patched (PATCH merge).
type IPAddressUpdateInput struct {
	ID                 int            `json:"id" jsonschema:"numeric ID of the IP address to update,required"`
	Address            *string        `json:"address,omitempty" jsonschema:"IP address or prefix, e.g. 192.168.1.1/24"`
	Status             *string        `json:"status,omitempty" jsonschema:"status: active, reserved, deprecated, dhcp, slaac"`             //nolint:lll // struct tag cannot be split across lines
	Role               *string        `json:"role,omitempty" jsonschema:"role: loopback, secondary, anycast, vip, vrrp, hsrp, glbp, carp"` //nolint:lll // struct tag cannot be split across lines
	VRF                *int           `json:"vrf,omitempty" jsonschema:"VRF ID"`
	Tenant             *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	DNSName            *string        `json:"dns_name,omitempty" jsonschema:"hostname or FQDN"`
	Description        *string        `json:"description,omitempty" jsonschema:"short description"`
	AssignedObjectType *string        `json:"assigned_object_type,omitempty" jsonschema:"assigned object type, e.g. dcim.interface"` //nolint:lll // struct tag cannot be split across lines
	AssignedObjectID   *int           `json:"assigned_object_id,omitempty" jsonschema:"ID of the assigned object"`
	Comments           *string        `json:"comments,omitempty" jsonschema:"free-form comments"`
	Tags               []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields       map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// IPAddressDeleteInput represents the input fields for deleting an IP address.
type IPAddressDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the IP address to delete,required"`
}

// PrefixCreateInput represents the writable fields for creating a prefix.
type PrefixCreateInput struct {
	Prefix       string         `json:"prefix" jsonschema:"network prefix in CIDR notation, e.g. 10.0.0.0/24 (required)"`
	Site         *int           `json:"site,omitempty" jsonschema:"site ID"`
	VRF          *int           `json:"vrf,omitempty" jsonschema:"VRF ID"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	VLAN         *int           `json:"vlan,omitempty" jsonschema:"VLAN ID"`
	Status       *string        `json:"status,omitempty" jsonschema:"status: container, active, reserved, deprecated"`
	Role         *int           `json:"role,omitempty" jsonschema:"role ID (look up first)"`
	IsPool       *bool          `json:"is_pool,omitempty" jsonschema:"whether this prefix is a pool"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"free-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// PrefixUpdateInput represents the writable fields for updating a prefix. All
// fields are optional; only the explicitly provided ones are patched (PATCH
// merge).
type PrefixUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the prefix to update,required"`
	Prefix       *string        `json:"prefix,omitempty" jsonschema:"network prefix in CIDR notation, e.g. 10.0.0.0/24"`
	Site         *int           `json:"site,omitempty" jsonschema:"site ID"`
	VRF          *int           `json:"vrf,omitempty" jsonschema:"VRF ID"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	VLAN         *int           `json:"vlan,omitempty" jsonschema:"VLAN ID"`
	Status       *string        `json:"status,omitempty" jsonschema:"status: container, active, reserved, deprecated"`
	Role         *int           `json:"role,omitempty" jsonschema:"role ID"`
	IsPool       *bool          `json:"is_pool,omitempty" jsonschema:"whether this prefix is a pool"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"free-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// PrefixDeleteInput represents the input fields for deleting a prefix.
type PrefixDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the prefix to delete,required"`
}

// VLANCreateInput represents the writable fields for creating a VLAN.
type VLANCreateInput struct {
	VID          int            `json:"vid" jsonschema:"VLAN ID, 1-4094 (required)"`
	Name         string         `json:"name" jsonschema:"VLAN name (required)"`
	Site         *int           `json:"site,omitempty" jsonschema:"site ID"`
	Group        *int           `json:"group,omitempty" jsonschema:"VLAN group ID"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Status       *string        `json:"status,omitempty" jsonschema:"status: active, reserved, deprecated"`
	Role         *int           `json:"role,omitempty" jsonschema:"role ID (look up first)"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"free-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VLANUpdateInput represents the writable fields for updating a VLAN. All
// fields are optional; only the explicitly provided ones are patched (PATCH
// merge).
type VLANUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the VLAN to update,required"`
	VID          *int           `json:"vid,omitempty" jsonschema:"VLAN ID, 1-4094"`
	Name         *string        `json:"name,omitempty" jsonschema:"VLAN name"`
	Site         *int           `json:"site,omitempty" jsonschema:"site ID"`
	Group        *int           `json:"group,omitempty" jsonschema:"VLAN group ID"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Status       *string        `json:"status,omitempty" jsonschema:"status: active, reserved, deprecated"`
	Role         *int           `json:"role,omitempty" jsonschema:"role ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"free-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VLANDeleteInput represents the input fields for deleting a VLAN.
type VLANDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the VLAN to delete,required"`
}

// VirtualMachineCreateInput represents the writable fields for creating a
// virtual machine.
type VirtualMachineCreateInput struct {
	Name         string         `json:"name" jsonschema:"VM name (required)"`
	Cluster      *int           `json:"cluster,omitempty" jsonschema:"cluster ID"`
	Role         *int           `json:"role,omitempty" jsonschema:"role ID (look up first)"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Platform     *int           `json:"platform,omitempty" jsonschema:"platform ID"`
	Status       *string        `json:"status,omitempty" jsonschema:"status: offline, active, planned, staged, failed, decommissioning"` //nolint:lll // struct tag cannot be split across lines
	Site         *int           `json:"site,omitempty" jsonschema:"site ID"`
	VCPUs        *float64       `json:"vcpus,omitempty" jsonschema:"number of virtual CPUs"`
	Memory       *int           `json:"memory,omitempty" jsonschema:"memory in MB"`
	Disk         *int           `json:"disk,omitempty" jsonschema:"disk size in GB"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"free-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VirtualMachineUpdateInput represents the writable fields for updating a
// virtual machine. All fields are optional; only the explicitly provided ones
// are patched (PATCH merge).
type VirtualMachineUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the virtual machine to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"VM name"`
	Cluster      *int           `json:"cluster,omitempty" jsonschema:"cluster ID"`
	Role         *int           `json:"role,omitempty" jsonschema:"role ID"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Platform     *int           `json:"platform,omitempty" jsonschema:"platform ID"`
	Status       *string        `json:"status,omitempty" jsonschema:"status: offline, active, planned, staged, failed, decommissioning"` //nolint:lll // struct tag cannot be split across lines
	Site         *int           `json:"site,omitempty" jsonschema:"site ID"`
	VCPUs        *float64       `json:"vcpus,omitempty" jsonschema:"number of virtual CPUs"`
	Memory       *int           `json:"memory,omitempty" jsonschema:"memory in MB"`
	Disk         *int           `json:"disk,omitempty" jsonschema:"disk size in GB"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"free-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VirtualMachineDeleteInput represents the input fields for deleting a virtual
// machine.
type VirtualMachineDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the virtual machine to delete,required"`
}

// ClusterCreateInput represents the writable fields for creating a cluster.
type ClusterCreateInput struct {
	Name         string         `json:"name" jsonschema:"cluster name (required)"`
	ClusterType  *int           `json:"type,omitempty" jsonschema:"cluster type ID (required)"`
	ClusterGroup *int           `json:"group,omitempty" jsonschema:"cluster group ID"`
	Site         *int           `json:"site,omitempty" jsonschema:"site ID"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"free-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ClusterUpdateInput represents the writable fields for updating a cluster. All
// fields are optional; only the explicitly provided ones are patched (PATCH
// merge).
type ClusterUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the cluster to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"cluster name"`
	ClusterType  *int           `json:"type,omitempty" jsonschema:"cluster type ID"`
	ClusterGroup *int           `json:"group,omitempty" jsonschema:"cluster group ID"`
	Site         *int           `json:"site,omitempty" jsonschema:"site ID"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"free-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ClusterDeleteInput represents the input fields for deleting a cluster.
type ClusterDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the cluster to delete,required"`
}

// CircuitCreateInput represents the writable fields for creating a circuit.
type CircuitCreateInput struct {
	CID          string         `json:"cid" jsonschema:"circuit ID (required)"`
	Provider     *int           `json:"provider,omitempty" jsonschema:"provider ID (required)"`
	CircuitType  *int           `json:"circuit_type,omitempty" jsonschema:"circuit type ID (required)"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Status       *string        `json:"status,omitempty" jsonschema:"status: planned, provisioning, active, offline, decommissioning"` //nolint:lll // struct tag cannot be split across lines
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"free-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
	InstallDate  *string        `json:"install_date,omitempty" jsonschema:"install date (YYYY-MM-DD)"`
	CommitRate   *int           `json:"commit_rate,omitempty" jsonschema:"committed rate in kbps"`
}

// CircuitUpdateInput represents the writable fields for updating a circuit. All
// fields are optional; only the explicitly provided ones are patched (PATCH
// merge).
type CircuitUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the circuit to update,required"`
	CID          *string        `json:"cid,omitempty" jsonschema:"circuit ID"`
	Provider     *int           `json:"provider,omitempty" jsonschema:"provider ID"`
	CircuitType  *int           `json:"circuit_type,omitempty" jsonschema:"circuit type ID"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Status       *string        `json:"status,omitempty" jsonschema:"status: planned, provisioning, active, offline, decommissioning"` //nolint:lll // struct tag cannot be split across lines
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"free-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
	InstallDate  *string        `json:"install_date,omitempty" jsonschema:"install date (YYYY-MM-DD)"`
	CommitRate   *int           `json:"commit_rate,omitempty" jsonschema:"committed rate in kbps"`
}

// CircuitDeleteInput represents the input fields for deleting a circuit.
type CircuitDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the circuit to delete,required"`
}

// RackCreateInput represents the writable fields for creating a rack.
type RackCreateInput struct {
	Name         string         `json:"name" jsonschema:"rack name (required)"`
	FacilityID   *string        `json:"facility_id,omitempty" jsonschema:"facility-assigned ID"`
	Site         *int           `json:"site,omitempty" jsonschema:"site ID"`
	Location     *int           `json:"location,omitempty" jsonschema:"location ID"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Status       *string        `json:"status,omitempty" jsonschema:"status: reserved, available, planned, active, decommissioning"` //nolint:lll // struct tag cannot be split across lines
	Role         *int           `json:"role,omitempty" jsonschema:"role ID (look up first)"`
	Serial       *string        `json:"serial,omitempty" jsonschema:"serial number"`
	AssetTag     *string        `json:"asset_tag,omitempty" jsonschema:"unique asset tag"`
	Type         *string        `json:"type,omitempty" jsonschema:"rack type: 2-post-frame, 4-post-frame, 4-post-cabinet, wall-frame, wall-cabinet"` //nolint:lll // struct tag cannot be split across lines
	Width        *int           `json:"width,omitempty" jsonschema:"rail-to-rail width in inches: 10, 19, 21, 23"`
	UHeight      *int           `json:"u_height,omitempty" jsonschema:"height in rack units"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"free-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// RackUpdateInput represents the writable fields for updating a rack. All
// fields are optional; only the explicitly provided ones are patched (PATCH
// merge).
type RackUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the rack to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"rack name"`
	FacilityID   *string        `json:"facility_id,omitempty" jsonschema:"facility-assigned ID"`
	Site         *int           `json:"site,omitempty" jsonschema:"site ID"`
	Location     *int           `json:"location,omitempty" jsonschema:"location ID"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Status       *string        `json:"status,omitempty" jsonschema:"status: reserved, available, planned, active, decommissioning"` //nolint:lll // struct tag cannot be split across lines
	Role         *int           `json:"role,omitempty" jsonschema:"role ID"`
	Serial       *string        `json:"serial,omitempty" jsonschema:"serial number"`
	AssetTag     *string        `json:"asset_tag,omitempty" jsonschema:"unique asset tag"`
	Type         *string        `json:"type,omitempty" jsonschema:"rack type"`
	Width        *int           `json:"width,omitempty" jsonschema:"rail-to-rail width in inches"`
	UHeight      *int           `json:"u_height,omitempty" jsonschema:"height in rack units"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"free-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// RackDeleteInput represents the input fields for deleting a rack.
type RackDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the rack to delete,required"`
}

// InterfaceCreateInput represents the writable fields for creating a device
// interface.
type InterfaceCreateInput struct {
	Device       *int           `json:"device,omitempty" jsonschema:"parent device ID (required)"`
	Name         string         `json:"name" jsonschema:"interface name (required)"`
	Type         *string        `json:"type,omitempty" jsonschema:"interface type (required), e.g. 1000base-t, 10gbase-x-sfpp"` //nolint:lll // struct tag cannot be split across lines
	Enabled      *bool          `json:"enabled,omitempty" jsonschema:"whether the interface is enabled"`
	MTU          *int           `json:"mtu,omitempty" jsonschema:"MTU"`
	MACAddress   *string        `json:"mac_address,omitempty" jsonschema:"MAC address"`
	Speed        *int           `json:"speed,omitempty" jsonschema:"interface speed in bps"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// InterfaceUpdateInput represents the writable fields for updating a device
// interface. All fields are optional; only the explicitly provided ones are
// patched (PATCH merge).
type InterfaceUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the interface to update,required"`
	Device       *int           `json:"device,omitempty" jsonschema:"parent device ID"`
	Name         *string        `json:"name,omitempty" jsonschema:"interface name"`
	Type         *string        `json:"type,omitempty" jsonschema:"interface type"`
	Enabled      *bool          `json:"enabled,omitempty" jsonschema:"whether the interface is enabled"`
	MTU          *int           `json:"mtu,omitempty" jsonschema:"MTU"`
	MACAddress   *string        `json:"mac_address,omitempty" jsonschema:"MAC address"`
	Speed        *int           `json:"speed,omitempty" jsonschema:"interface speed in bps"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// InterfaceDeleteInput represents the input fields for deleting an interface.
type InterfaceDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the interface to delete,required"`
}

// CircuitTerminationCreateInput represents the writable fields for creating a
// circuit termination.
type CircuitTerminationCreateInput struct {
	Circuit       *int           `json:"circuit,omitempty" jsonschema:"circuit ID (required)"`
	TermSide      string         `json:"term_side" jsonschema:"termination side, A or Z (required)"`
	Site          *int           `json:"site,omitempty" jsonschema:"site ID (required)"`
	Speed         *int           `json:"speed,omitempty" jsonschema:"termination speed in kbps"`
	UpstreamSpeed *int           `json:"upstream_speed,omitempty" jsonschema:"upstream speed in kbps"`
	Description   *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags          []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields  map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// CircuitTerminationUpdateInput represents the writable fields for updating a
// circuit termination. All fields are optional; only the explicitly provided
// ones are patched (PATCH merge).
type CircuitTerminationUpdateInput struct {
	ID            int            `json:"id" jsonschema:"numeric ID of the circuit termination to update,required"`
	Circuit       *int           `json:"circuit,omitempty" jsonschema:"circuit ID"`
	TermSide      *string        `json:"term_side,omitempty" jsonschema:"termination side, A or Z"`
	Site          *int           `json:"site,omitempty" jsonschema:"site ID"`
	Speed         *int           `json:"speed,omitempty" jsonschema:"termination speed in kbps"`
	UpstreamSpeed *int           `json:"upstream_speed,omitempty" jsonschema:"upstream speed in kbps"`
	Description   *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags          []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields  map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// CircuitTerminationDeleteInput represents the input fields for deleting a
// circuit termination.
type CircuitTerminationDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the circuit termination to delete,required"`
}

// CableTerminationInput represents one end of a cable connection for the
// create/update tools. NetBox requires the object_type (e.g. "dcim.interface")
// and the numeric object_id.
type CableTerminationInput struct {
	ObjectType string `json:"object_type" jsonschema:"type of the terminated object, e.g. dcim.interface (required)"`
	ObjectID   int    `json:"object_id" jsonschema:"numeric ID of the terminated object (required)"`
}

// CableCreateInput represents the writable fields for creating a cable.
type CableCreateInput struct {
	TerminationA *CableTerminationInput `json:"termination_a,omitempty" jsonschema:"first end of the cable (required)"`
	TerminationB *CableTerminationInput `json:"termination_b,omitempty" jsonschema:"second end of the cable (required)"`
	Type         *string                `json:"type,omitempty" jsonschema:"cable type (choice string)"`
	Status       *string                `json:"status,omitempty" jsonschema:"cable status (choice string)"`
	Label        *string                `json:"label,omitempty" jsonschema:"cable label"`
	Color        *string                `json:"color,omitempty" jsonschema:"cable color (hex or name)"`
	Length       *float64               `json:"length,omitempty" jsonschema:"cable length"`
	LengthUnit   *string                `json:"length_unit,omitempty" jsonschema:"length unit (choice string)"`
	Description  *string                `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string               `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any         `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// CableUpdateInput represents the writable fields for updating a cable. All
// fields are optional; only the explicitly provided ones are patched (PATCH
// merge).
type CableUpdateInput struct {
	ID           int                    `json:"id" jsonschema:"numeric ID of the cable to update,required"`
	TerminationA *CableTerminationInput `json:"termination_a,omitempty" jsonschema:"first end of the cable"`
	TerminationB *CableTerminationInput `json:"termination_b,omitempty" jsonschema:"second end of the cable"`
	Type         *string                `json:"type,omitempty" jsonschema:"cable type (choice string)"`
	Status       *string                `json:"status,omitempty" jsonschema:"cable status (choice string)"`
	Label        *string                `json:"label,omitempty" jsonschema:"cable label"`
	Color        *string                `json:"color,omitempty" jsonschema:"cable color (hex or name)"`
	Length       *float64               `json:"length,omitempty" jsonschema:"cable length"`
	LengthUnit   *string                `json:"length_unit,omitempty" jsonschema:"length unit (choice string)"`
	Description  *string                `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string               `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any         `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// CableDeleteInput represents the input fields for deleting a cable.
type CableDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the cable to delete,required"`
}

// VMInterfaceCreateInput represents the writable fields for creating a
// virtual-machine interface.
type VMInterfaceCreateInput struct {
	VirtualMachine *int           `json:"virtual_machine,omitempty" jsonschema:"parent virtual machine ID (required)"`
	Name           string         `json:"name" jsonschema:"interface name (required)"`
	Enabled        *bool          `json:"enabled,omitempty" jsonschema:"whether the interface is enabled"`
	MTU            *int           `json:"mtu,omitempty" jsonschema:"MTU"`
	MACAddress     *string        `json:"mac_address,omitempty" jsonschema:"MAC address"`
	Description    *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags           []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields   map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VMInterfaceUpdateInput represents the writable fields for updating a
// virtual-machine interface. All fields are optional; only the explicitly
// provided ones are patched (PATCH merge).
type VMInterfaceUpdateInput struct {
	ID             int            `json:"id" jsonschema:"numeric ID of the VM interface to update,required"`
	VirtualMachine *int           `json:"virtual_machine,omitempty" jsonschema:"parent virtual machine ID"`
	Name           *string        `json:"name,omitempty" jsonschema:"interface name"`
	Enabled        *bool          `json:"enabled,omitempty" jsonschema:"whether the interface is enabled"`
	MTU            *int           `json:"mtu,omitempty" jsonschema:"MTU"`
	MACAddress     *string        `json:"mac_address,omitempty" jsonschema:"MAC address"`
	Description    *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags           []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields   map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VMInterfaceDeleteInput represents the input fields for deleting a
// virtual-machine interface.
type VMInterfaceDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the VM interface to delete,required"`
}

// ProviderCreateInput represents the writable fields for creating a provider.
type ProviderCreateInput struct {
	Name         string         `json:"name" jsonschema:"provider name (required)"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Asn          *int           `json:"asn,omitempty" jsonschema:"autonomous system number"`
	Account      *string        `json:"account,omitempty" jsonschema:"provider account number"`
	PortalURL    *string        `json:"portal_url,omitempty" jsonschema:"portal URL"`
	NocContact   *string        `json:"noc_contact,omitempty" jsonschema:"NOC contact"`
	AdminContact *string        `json:"admin_contact,omitempty" jsonschema:"administrative contact"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"long-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ProviderUpdateInput represents the writable fields for updating a provider.
// All fields are optional; only the explicitly provided ones are patched (PATCH
// merge).
type ProviderUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the provider to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"provider name"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Asn          *int           `json:"asn,omitempty" jsonschema:"autonomous system number"`
	Account      *string        `json:"account,omitempty" jsonschema:"provider account number"`
	PortalURL    *string        `json:"portal_url,omitempty" jsonschema:"portal URL"`
	NocContact   *string        `json:"noc_contact,omitempty" jsonschema:"NOC contact"`
	AdminContact *string        `json:"admin_contact,omitempty" jsonschema:"administrative contact"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"long-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ProviderDeleteInput represents the input fields for deleting a provider.
type ProviderDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the provider to delete,required"`
}

// TenantCreateInput represents the writable fields for creating a tenant.
type TenantCreateInput struct {
	Name         string         `json:"name" jsonschema:"tenant name (required)"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Group        *int           `json:"group,omitempty" jsonschema:"tenant group ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"long-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// TenantUpdateInput represents the writable fields for updating a tenant. All
// fields are optional; only the explicitly provided ones are patched (PATCH
// merge).
type TenantUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the tenant to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"tenant name"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Group        *int           `json:"group,omitempty" jsonschema:"tenant group ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"long-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// TenantDeleteInput represents the input fields for deleting a tenant.
type TenantDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the tenant to delete,required"`
}

// TagCreateInput represents the writable fields for creating a tag.
type TagCreateInput struct {
	Name        string  `json:"name" jsonschema:"tag name (required)"`
	Slug        *string `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Color       *string `json:"color,omitempty" jsonschema:"color (hex value, e.g. ff0000)"`
	Description *string `json:"description,omitempty" jsonschema:"short description"`
}

// TagUpdateInput represents the writable fields for updating a tag. All fields
// are optional; only the explicitly provided ones are patched (PATCH merge).
type TagUpdateInput struct {
	ID          int     `json:"id" jsonschema:"numeric ID of the tag to update,required"`
	Name        *string `json:"name,omitempty" jsonschema:"tag name"`
	Slug        *string `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Color       *string `json:"color,omitempty" jsonschema:"color (hex value, e.g. ff0000)"`
	Description *string `json:"description,omitempty" jsonschema:"short description"`
}

// TagDeleteInput represents the input fields for deleting a tag.
type TagDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the tag to delete,required"`
}

// ManufacturerCreateInput represents the writable fields for creating a
// manufacturer.
type ManufacturerCreateInput struct {
	Name         string         `json:"name" jsonschema:"manufacturer name (required)"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ManufacturerUpdateInput represents the writable fields for updating a
// manufacturer. All fields are optional; only the explicitly provided ones are
// patched (PATCH merge).
type ManufacturerUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the manufacturer to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"manufacturer name"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ManufacturerDeleteInput represents the input fields for deleting a
// manufacturer.
type ManufacturerDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the manufacturer to delete,required"`
}

// DeviceTypeCreateInput represents the writable fields for creating a device
// type. Manufacturer (numeric NetBox ID) and Model are required on create.
type DeviceTypeCreateInput struct {
	Manufacturer  *int           `json:"manufacturer,omitempty" jsonschema:"numeric ID of the manufacturer (required)"`
	Model         string         `json:"model" jsonschema:"device model name (required)"`
	Slug          *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	PartNumber    *string        `json:"part_number,omitempty" jsonschema:"part number"`
	UHeight       *float64       `json:"u_height,omitempty" jsonschema:"height in rack units (U)"`
	IsFullDepth   *bool          `json:"is_full_depth,omitempty" jsonschema:"whether the device is full depth"`
	SubdeviceRole *string        `json:"subdevice_role,omitempty" jsonschema:"parent/child subdevice role"`
	Comments      *string        `json:"comments,omitempty" jsonschema:"long-form comments"`
	Tags          []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields  map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// DeviceTypeUpdateInput represents the writable fields for updating a device
// type. All fields are optional; only the explicitly provided ones are patched
// (PATCH merge).
type DeviceTypeUpdateInput struct {
	ID            int            `json:"id" jsonschema:"numeric ID of the device type to update,required"`
	Manufacturer  *int           `json:"manufacturer,omitempty" jsonschema:"numeric ID of the manufacturer"`
	Model         *string        `json:"model,omitempty" jsonschema:"device model name"`
	Slug          *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	PartNumber    *string        `json:"part_number,omitempty" jsonschema:"part number"`
	UHeight       *float64       `json:"u_height,omitempty" jsonschema:"height in rack units (U)"`
	IsFullDepth   *bool          `json:"is_full_depth,omitempty" jsonschema:"whether the device is full depth"`
	SubdeviceRole *string        `json:"subdevice_role,omitempty" jsonschema:"parent/child subdevice role"`
	Comments      *string        `json:"comments,omitempty" jsonschema:"long-form comments"`
	Tags          []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields  map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// DeviceTypeDeleteInput represents the input fields for deleting a device type.
type DeviceTypeDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the device type to delete,required"`
}

// LocationCreateInput represents the writable fields for creating a location.
// Name and Site (numeric NetBox ID) are required on create.
type LocationCreateInput struct {
	Name         string         `json:"name" jsonschema:"location name (required)"`
	Site         *int           `json:"site,omitempty" jsonschema:"numeric ID of the parent site (required)"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Parent       *int           `json:"parent,omitempty" jsonschema:"numeric ID of a parent location"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Status       *string        `json:"status,omitempty" jsonschema:"location status (e.g. active, planned)"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// LocationUpdateInput represents the writable fields for updating a location.
// All fields are optional; only the explicitly provided ones are patched
// (PATCH merge).
type LocationUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the location to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"location name"`
	Site         *int           `json:"site,omitempty" jsonschema:"numeric ID of the parent site"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Parent       *int           `json:"parent,omitempty" jsonschema:"numeric ID of a parent location"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Status       *string        `json:"status,omitempty" jsonschema:"location status (e.g. active, planned)"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// LocationDeleteInput represents the input fields for deleting a location.
type LocationDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the location to delete,required"`
}

// ClusterTypeCreateInput represents the writable fields for creating a cluster
// type. Name is required on create.
type ClusterTypeCreateInput struct {
	Name         string         `json:"name" jsonschema:"cluster type name (required)"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ClusterTypeUpdateInput represents the writable fields for updating a cluster
// type. All fields are optional; only the explicitly provided ones are patched
// (PATCH merge).
type ClusterTypeUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the cluster type to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"cluster type name"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ClusterTypeDeleteInput represents the input fields for deleting a cluster type.
type ClusterTypeDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the cluster type to delete,required"`
}

// ClusterGroupCreateInput represents the writable fields for creating a cluster
// group. Name is required on create.
type ClusterGroupCreateInput struct {
	Name         string         `json:"name" jsonschema:"cluster group name (required)"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ClusterGroupUpdateInput represents the writable fields for updating a cluster
// group. All fields are optional; only the explicitly provided ones are patched
// (PATCH merge).
type ClusterGroupUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the cluster group to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"cluster group name"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ClusterGroupDeleteInput represents the input fields for deleting a cluster group.
type ClusterGroupDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the cluster group to delete,required"`
}

// CircuitTypeCreateInput represents the writable fields for creating a circuit
// type. Name is required on create.
type CircuitTypeCreateInput struct {
	Name         string         `json:"name" jsonschema:"circuit type name (required)"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// CircuitTypeUpdateInput represents the writable fields for updating a circuit
// type. All fields are optional; only the explicitly provided ones are patched
// (PATCH merge).
type CircuitTypeUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the circuit type to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"circuit type name"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// CircuitTypeDeleteInput represents the input fields for deleting a circuit type.
type CircuitTypeDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the circuit type to delete,required"`
}

// VrfCreateInput represents the writable fields for creating a VRF. Name and rd
// are required on create.
type VrfCreateInput struct {
	Name         string         `json:"name" jsonschema:"VRF name (required)"`
	Rd           string         `json:"rd" jsonschema:"route distinguisher (required)"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"numeric tenant ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VrfUpdateInput represents the writable fields for updating a VRF. All fields
// are optional; only the explicitly provided ones are patched (PATCH merge).
type VrfUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the VRF to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"VRF name"`
	Rd           *string        `json:"rd,omitempty" jsonschema:"route distinguisher"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"numeric tenant ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VrfDeleteInput represents the input fields for deleting a VRF.
type VrfDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the VRF to delete,required"`
}

// VlanGroupCreateInput represents the writable fields for creating a VLAN
// group. Name is required on create.
type VlanGroupCreateInput struct {
	Name         string         `json:"name" jsonschema:"VLAN group name (required)"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	ScopeType    *string        `json:"scope_type,omitempty" jsonschema:"scope content type (e.g. dcim.site)"`
	ScopeID      *int           `json:"scope_id,omitempty" jsonschema:"numeric ID of the scoped object"`
	MinVID       *int           `json:"min_vid,omitempty" jsonschema:"lowest permissible child VLAN ID"`
	MaxVID       *int           `json:"max_vid,omitempty" jsonschema:"highest permissible child VLAN ID"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VlanGroupUpdateInput represents the writable fields for updating a VLAN
// group. All fields are optional; only the explicitly provided ones are patched
// (PATCH merge).
type VlanGroupUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the VLAN group to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"VLAN group name"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	ScopeType    *string        `json:"scope_type,omitempty" jsonschema:"scope content type (e.g. dcim.site)"`
	ScopeID      *int           `json:"scope_id,omitempty" jsonschema:"numeric ID of the scoped object"`
	MinVID       *int           `json:"min_vid,omitempty" jsonschema:"lowest permissible child VLAN ID"`
	MaxVID       *int           `json:"max_vid,omitempty" jsonschema:"highest permissible child VLAN ID"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VlanGroupDeleteInput represents the input fields for deleting a VLAN group.
type VlanGroupDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the VLAN group to delete,required"`
}

// RoleCreateInput represents the writable fields for creating an IPAM role.
// name is required on create.
type RoleCreateInput struct {
	Name         string         `json:"name" jsonschema:"role name (required)"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Weight       *int           `json:"weight,omitempty" jsonschema:"sorting weight"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// RoleUpdateInput represents the writable fields for updating an IPAM role.
// All fields are optional; only the explicitly provided ones are patched
// (PATCH merge).
type RoleUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the role to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"role name"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Weight       *int           `json:"weight,omitempty" jsonschema:"sorting weight"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// RoleDeleteInput represents the input fields for deleting an IPAM role.
type RoleDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the role to delete,required"`
}

// ContactCreateInput represents the writable fields for creating a tenancy
// contact. name is required on create; group is optional.
type ContactCreateInput struct {
	Name         string         `json:"name" jsonschema:"contact name (required)"`
	Group        *int           `json:"group,omitempty" jsonschema:"numeric contact group ID"`
	Title        *string        `json:"title,omitempty" jsonschema:"job title or role"`
	Phone        *string        `json:"phone,omitempty" jsonschema:"phone number"`
	Email        *string        `json:"email,omitempty" jsonschema:"email address"`
	Address      *string        `json:"address,omitempty" jsonschema:"physical address"`
	Link         *string        `json:"link,omitempty" jsonschema:"link/URL"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"long-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ContactUpdateInput represents the writable fields for updating a tenancy
// contact. All fields are optional; only the explicitly provided ones are
// patched (PATCH merge).
type ContactUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the contact to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"contact name"`
	Group        *int           `json:"group,omitempty" jsonschema:"numeric contact group ID"`
	Title        *string        `json:"title,omitempty" jsonschema:"job title or role"`
	Phone        *string        `json:"phone,omitempty" jsonschema:"phone number"`
	Email        *string        `json:"email,omitempty" jsonschema:"email address"`
	Address      *string        `json:"address,omitempty" jsonschema:"physical address"`
	Link         *string        `json:"link,omitempty" jsonschema:"link/URL"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"long-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ContactDeleteInput represents the input fields for deleting a tenancy contact.
type ContactDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the contact to delete,required"`
}

// — output types —

// PaginatedOutput is a generic paginated response used by all list-oriented tools.
type PaginatedOutput[D any] struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []D    `json:"results"`
}

// GetObjectOutput represents the output for the get_object_by_id tool.
//
// Data is typed as `any` (rather than domain.RawObject) so the jsonschema
// inferred for the tool's output accepts any JSON value. domain.RawObject is
// json.RawMessage ([]byte), whose inferred schema only admits "null, array",
// which would make the server reject object payloads during output validation.
// The handler still stores a domain.RawObject here, whose MarshalJSON emits the
// raw object inline.
type GetObjectOutput struct {
	Data any `json:"data"`
}

// SiteOutput represents the output for the create/update site tools.
type SiteOutput struct {
	Data domain.Site `json:"data"`
}

// RegionOutput represents the output for the create/update region tools.
type RegionOutput struct {
	Data domain.Region `json:"data"`
}

// TenantGroupOutput represents the output for the create/update tenant group
// tools.
type TenantGroupOutput struct {
	Data domain.TenantGroup `json:"data"`
}

// SiteGroupOutput represents the output for the create/update site group tools.
type SiteGroupOutput struct {
	Data domain.SiteGroup `json:"data"`
}

// RackGroupOutput represents the output for the create/update rack group tools.
type RackGroupOutput struct {
	Data domain.RackGroup `json:"data"`
}

// RackRoleOutput represents the output for the create/update rack role tools.
type RackRoleOutput struct {
	Data domain.RackRole `json:"data"`
}

// DeviceRoleOutput represents the output for the create/update device role
// tools.
type DeviceRoleOutput struct {
	Data domain.DeviceRole `json:"data"`
}

// PlatformOutput represents the output for the create/update platform tools.
type PlatformOutput struct {
	Data domain.Platform `json:"data"`
}

// VirtualChassisOutput represents the output for the create/update virtual
// chassis tools.
type VirtualChassisOutput struct {
	Data domain.VirtualChassis `json:"data"`
}

// PowerPanelOutput represents the output for the create/update power panel
// tools.
type PowerPanelOutput struct {
	Data domain.PowerPanel `json:"data"`
}

// InventoryItemRoleOutput represents the output for the create/update inventory
// item role tools.
type InventoryItemRoleOutput struct {
	Data domain.InventoryItemRole `json:"data"`
}

// DeviceOutput represents the output for the create/update device tools.
type DeviceOutput struct {
	Data domain.Device `json:"data"`
}

// IPAddressOutput represents the output for the create/update IP address tools.
type IPAddressOutput struct {
	Data domain.IPAddress `json:"data"`
}

// PrefixOutput represents the output for the create/update prefix tools.
type PrefixOutput struct {
	Data domain.Prefix `json:"data"`
}

// VLANOutput represents the output for the create/update VLAN tools.
type VLANOutput struct {
	Data domain.VLAN `json:"data"`
}

// VirtualMachineOutput represents the output for the create/update virtual
// machine tools.
type VirtualMachineOutput struct {
	Data domain.VirtualMachine `json:"data"`
}

// ClusterOutput represents the output for the create/update cluster tools.
type ClusterOutput struct {
	Data domain.Cluster `json:"data"`
}

// CircuitOutput represents the output for the create/update circuit tools.
type CircuitOutput struct {
	Data domain.Circuit `json:"data"`
}

// RackOutput represents the output for the create/update rack tools.
type RackOutput struct {
	Data domain.Rack `json:"data"`
}

// InterfaceOutput represents the output for the create/update interface tools.
type InterfaceOutput struct {
	Data domain.Interface `json:"data"`
}

// CircuitTerminationOutput represents the output for the create/update circuit
// termination tools.
type CircuitTerminationOutput struct {
	Data domain.CircuitTermination `json:"data"`
}

// CableOutput represents the output for the create/update cable tools.
type CableOutput struct {
	Data domain.Cable `json:"data"`
}

// VMInterfaceOutput represents the output for the create/update VM interface
// tools.
type VMInterfaceOutput struct {
	Data domain.VMInterface `json:"data"`
}

// ProviderOutput represents the output for the create/update provider tools.
type ProviderOutput struct {
	Data domain.Provider `json:"data"`
}

// TenantOutput represents the output for the create/update tenant tools.
type TenantOutput struct {
	Data domain.Tenant `json:"data"`
}

// TagOutput represents the output for the create/update tag tools.
type TagOutput struct {
	Data domain.Tag `json:"data"`
}

// ManufacturerOutput represents the output for the create/update manufacturer
// tools.
type ManufacturerOutput struct {
	Data domain.Manufacturer `json:"data"`
}

// DeviceTypeOutput represents the output for the create/update device type
// tools.
type DeviceTypeOutput struct {
	Data domain.DeviceType `json:"data"`
}

// LocationOutput represents the output for the create/update location tools.
type LocationOutput struct {
	Data domain.Location `json:"data"`
}

// ClusterTypeOutput represents the output for the create/update cluster type
// tools.
type ClusterTypeOutput struct {
	Data domain.ClusterType `json:"data"`
}

// ClusterGroupOutput represents the output for the create/update cluster group
// tools.
type ClusterGroupOutput struct {
	Data domain.ClusterGroup `json:"data"`
}

// CircuitTypeOutput represents the output for the create/update circuit type
// tools.
type CircuitTypeOutput struct {
	Data domain.CircuitType `json:"data"`
}

// VrfOutput represents the output for the create/update VRF tools.
type VrfOutput struct {
	Data domain.Vrf `json:"data"`
}

// VlanGroupOutput represents the output for the create/update VLAN group tools.
type VlanGroupOutput struct {
	Data domain.VlanGroup `json:"data"`
}

// RoleOutput represents the output for the create/update role tools.
type RoleOutput struct {
	Data domain.Role `json:"data"`
}

// ContactOutput represents the output for the create/update contact tools.
type ContactOutput struct {
	Data domain.Contact `json:"data"`
}

// — write helpers —

// writeErrorResult converts an error into an MCP error result. A
// domain.ValidationError is surfaced as IsError with the NetBox body delivered
// via StructuredContent (not through the error string, so it never leaks into
// logs); all other errors are returned as a plain error to be surfaced normally.
func writeErrorResult(err error) (*mcp.CallToolResult, error) {
	if ve, ok := errors.AsType[*domain.ValidationError](err); ok {
		return &mcp.CallToolResult{
			IsError: true,
			StructuredContent: map[string]any{
				"validation_errors": json.RawMessage(ve.Body),
			},
		}, nil
	}

	return &mcp.CallToolResult{IsError: true}, err
}

// siteWriteFromCreate builds a domain.SiteWrite from a create input.
func siteWriteFromCreate(input SiteCreateInput) domain.SiteWrite {
	return domain.SiteWrite{
		Name:            input.Name,
		Slug:            strPtr(input.Slug),
		Status:          input.Status,
		Region:          input.Region,
		Tenant:          input.Tenant,
		Facility:        input.Facility,
		TimeZone:        input.TimeZone,
		Description:     input.Description,
		PhysicalAddress: input.PhysicalAddress,
		ShippingAddress: input.ShippingAddress,
		Comments:        input.Comments,
		Tags:            input.Tags,
		CustomFields:    input.CustomFields,
	}
}

// siteWriteFromUpdate builds a domain.SiteWrite from an update input, mapping
// only the explicitly-provided (non-nil) fields so the PATCH body is minimal.
func siteWriteFromUpdate(input SiteUpdateInput) domain.SiteWrite {
	write := domain.SiteWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Status = input.Status
	write.Region = input.Region
	write.Tenant = input.Tenant
	write.Facility = input.Facility
	write.TimeZone = input.TimeZone
	write.Description = input.Description
	write.PhysicalAddress = input.PhysicalAddress
	write.ShippingAddress = input.ShippingAddress
	write.Comments = input.Comments
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// regionWriteFromCreate builds a domain.RegionWrite from a create input.
func regionWriteFromCreate(input RegionCreateInput) domain.RegionWrite {
	return domain.RegionWrite{
		Name:         input.Name,
		Slug:         strPtr(input.Slug),
		Parent:       input.Parent,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// regionWriteFromUpdate builds a domain.RegionWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func regionWriteFromUpdate(input RegionUpdateInput) domain.RegionWrite {
	write := domain.RegionWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Parent = input.Parent
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// tenantGroupWriteFromCreate builds a domain.TenantGroupWrite from a create
// input.
func tenantGroupWriteFromCreate(input TenantGroupCreateInput) domain.TenantGroupWrite {
	return domain.TenantGroupWrite{
		Name:         input.Name,
		Slug:         strPtr(input.Slug),
		Parent:       input.Parent,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// tenantGroupWriteFromUpdate builds a domain.TenantGroupWrite from an update
// input, mapping only the explicitly-provided (non-nil) fields so the PATCH
// body is minimal.
func tenantGroupWriteFromUpdate(input TenantGroupUpdateInput) domain.TenantGroupWrite {
	write := domain.TenantGroupWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Parent = input.Parent
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}

// siteGroupWriteFromCreate builds a domain.SiteGroupWrite from a create input.
func siteGroupWriteFromCreate(input SiteGroupCreateInput) domain.SiteGroupWrite {
	return domain.SiteGroupWrite{
		Name:         input.Name,
		Slug:         strPtr(input.Slug),
		Parent:       input.Parent,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// siteGroupWriteFromUpdate builds a domain.SiteGroupWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func siteGroupWriteFromUpdate(input SiteGroupUpdateInput) domain.SiteGroupWrite {
	write := domain.SiteGroupWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Parent = input.Parent
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// rackGroupWriteFromCreate builds a domain.RackGroupWrite from a create input.
func rackGroupWriteFromCreate(input RackGroupCreateInput) domain.RackGroupWrite {
	return domain.RackGroupWrite{
		Name:        input.Name,
		Slug:        strPtr(input.Slug),
		Site:        input.Site,
		Parent:      input.Parent,
		Description: input.Description,
		Tags:        input.Tags,
	}
}

// rackGroupWriteFromUpdate builds a domain.RackGroupWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func rackGroupWriteFromUpdate(input RackGroupUpdateInput) domain.RackGroupWrite {
	write := domain.RackGroupWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Site = input.Site
	write.Parent = input.Parent
	write.Description = input.Description
	write.Tags = input.Tags

	return write
}

// rackRoleWriteFromCreate builds a domain.RackRoleWrite from a create input.
func rackRoleWriteFromCreate(input RackRoleCreateInput) domain.RackRoleWrite {
	return domain.RackRoleWrite{
		Name:         input.Name,
		Slug:         strPtr(input.Slug),
		Color:        input.Color,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// rackRoleWriteFromUpdate builds a domain.RackRoleWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func rackRoleWriteFromUpdate(input RackRoleUpdateInput) domain.RackRoleWrite {
	write := domain.RackRoleWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Color = input.Color
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// deviceRoleWriteFromCreate builds a domain.DeviceRoleWrite from a create
// input.
func deviceRoleWriteFromCreate(input DeviceRoleCreateInput) domain.DeviceRoleWrite {
	return domain.DeviceRoleWrite{
		Name:           input.Name,
		Slug:           strPtr(input.Slug),
		Color:          input.Color,
		VMRole:         input.VMRole,
		ConfigTemplate: input.ConfigTemplate,
		Description:    input.Description,
		Tags:           input.Tags,
		CustomFields:   input.CustomFields,
		Weight:         input.Weight,
	}
}

// deviceRoleWriteFromUpdate builds a domain.DeviceRoleWrite from an update
// input, mapping only the explicitly-provided (non-nil) fields so the PATCH
// body is minimal.
func deviceRoleWriteFromUpdate(input DeviceRoleUpdateInput) domain.DeviceRoleWrite {
	write := domain.DeviceRoleWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Color = input.Color
	write.VMRole = input.VMRole
	write.ConfigTemplate = input.ConfigTemplate
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields
	write.Weight = input.Weight

	return write
}

// platformWriteFromCreate builds a domain.PlatformWrite from a create input.
func platformWriteFromCreate(input PlatformCreateInput) domain.PlatformWrite {
	return domain.PlatformWrite{
		Name:           input.Name,
		Slug:           strPtr(input.Slug),
		Manufacturer:   input.Manufacturer,
		ConfigTemplate: input.ConfigTemplate,
		Description:    input.Description,
		Tags:           input.Tags,
		CustomFields:   input.CustomFields,
	}
}

// platformWriteFromUpdate builds a domain.PlatformWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func platformWriteFromUpdate(input PlatformUpdateInput) domain.PlatformWrite {
	write := domain.PlatformWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Manufacturer = input.Manufacturer
	write.ConfigTemplate = input.ConfigTemplate
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// virtualChassisWriteFromCreate builds a domain.VirtualChassisWrite from a
// create input. Name is optional for this entity.
func virtualChassisWriteFromCreate(input VirtualChassisCreateInput) domain.VirtualChassisWrite {
	return domain.VirtualChassisWrite{
		Name:         input.Name,
		Domain:       input.Domain,
		Master:       input.Master,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// virtualChassisWriteFromUpdate builds a domain.VirtualChassisWrite from an
// update input, mapping only the explicitly-provided (non-nil) fields so the
// PATCH body is minimal.
func virtualChassisWriteFromUpdate(input VirtualChassisUpdateInput) domain.VirtualChassisWrite {
	write := domain.VirtualChassisWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Domain = input.Domain
	write.Master = input.Master
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// powerPanelWriteFromCreate builds a domain.PowerPanelWrite from a create
// input.
func powerPanelWriteFromCreate(input PowerPanelCreateInput) domain.PowerPanelWrite {
	return domain.PowerPanelWrite{
		Site:         input.Site,
		Location:     input.Location,
		Name:         input.Name,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// powerPanelWriteFromUpdate builds a domain.PowerPanelWrite from an update
// input, mapping only the explicitly-provided (non-nil) fields so the PATCH
// body is minimal.
func powerPanelWriteFromUpdate(input PowerPanelUpdateInput) domain.PowerPanelWrite {
	write := domain.PowerPanelWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Site = input.Site
	write.Location = input.Location
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// inventoryItemRoleWriteFromCreate builds a domain.InventoryItemRoleWrite from
// a create input.
func inventoryItemRoleWriteFromCreate(input InventoryItemRoleCreateInput) domain.InventoryItemRoleWrite {
	return domain.InventoryItemRoleWrite{
		Name:         input.Name,
		Slug:         strPtr(input.Slug),
		Color:        input.Color,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// inventoryItemRoleWriteFromUpdate builds a domain.InventoryItemRoleWrite from
// an update input, mapping only the explicitly-provided (non-nil) fields so the
// PATCH body is minimal.
func inventoryItemRoleWriteFromUpdate(input InventoryItemRoleUpdateInput) domain.InventoryItemRoleWrite {
	write := domain.InventoryItemRoleWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Color = input.Color
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// deviceWriteFromCreate builds a domain.DeviceWrite from a create input.
func deviceWriteFromCreate(input DeviceCreateInput) domain.DeviceWrite {
	return domain.DeviceWrite{
		Name:         input.Name,
		DeviceType:   input.DeviceType,
		Role:         input.Role,
		Tenant:       input.Tenant,
		Platform:     input.Platform,
		Serial:       input.Serial,
		AssetTag:     input.AssetTag,
		Site:         input.Site,
		Rack:         input.Rack,
		Position:     input.Position,
		Face:         input.Face,
		Status:       input.Status,
		Cluster:      input.Cluster,
		Comments:     input.Comments,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// deviceWriteFromUpdate builds a domain.DeviceWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func deviceWriteFromUpdate(input DeviceUpdateInput) domain.DeviceWrite {
	write := domain.DeviceWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.DeviceType = input.DeviceType
	write.Role = input.Role
	write.Tenant = input.Tenant
	write.Platform = input.Platform
	write.Serial = input.Serial
	write.AssetTag = input.AssetTag
	write.Site = input.Site
	write.Rack = input.Rack
	write.Position = input.Position
	write.Face = input.Face
	write.Status = input.Status
	write.Cluster = input.Cluster
	write.Comments = input.Comments
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// ipAddressWriteFromCreate builds a domain.IPAddressWrite from a create input.
func ipAddressWriteFromCreate(input IPAddressCreateInput) domain.IPAddressWrite {
	return domain.IPAddressWrite{
		Address:            input.Address,
		Status:             input.Status,
		Role:               input.Role,
		VRF:                input.VRF,
		Tenant:             input.Tenant,
		DNSName:            input.DNSName,
		Description:        input.Description,
		AssignedObjectType: input.AssignedObjectType,
		AssignedObjectID:   input.AssignedObjectID,
		Comments:           input.Comments,
		Tags:               input.Tags,
		CustomFields:       input.CustomFields,
	}
}

// ipAddressWriteFromUpdate builds a domain.IPAddressWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func ipAddressWriteFromUpdate(input IPAddressUpdateInput) domain.IPAddressWrite {
	write := domain.IPAddressWrite{}
	if input.Address != nil {
		write.Address = *input.Address
	}
	write.Status = input.Status
	write.Role = input.Role
	write.VRF = input.VRF
	write.Tenant = input.Tenant
	write.DNSName = input.DNSName
	write.Description = input.Description
	write.AssignedObjectType = input.AssignedObjectType
	write.AssignedObjectID = input.AssignedObjectID
	write.Comments = input.Comments
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// prefixWriteFromCreate builds a domain.PrefixWrite from a create input.
func prefixWriteFromCreate(input PrefixCreateInput) domain.PrefixWrite {
	return domain.PrefixWrite{
		Prefix:       input.Prefix,
		Site:         input.Site,
		VRF:          input.VRF,
		Tenant:       input.Tenant,
		VLAN:         input.VLAN,
		Status:       input.Status,
		Role:         input.Role,
		IsPool:       input.IsPool,
		Description:  input.Description,
		Comments:     input.Comments,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// prefixWriteFromUpdate builds a domain.PrefixWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func prefixWriteFromUpdate(input PrefixUpdateInput) domain.PrefixWrite {
	write := domain.PrefixWrite{}
	if input.Prefix != nil {
		write.Prefix = *input.Prefix
	}
	write.Site = input.Site
	write.VRF = input.VRF
	write.Tenant = input.Tenant
	write.VLAN = input.VLAN
	write.Status = input.Status
	write.Role = input.Role
	write.IsPool = input.IsPool
	write.Description = input.Description
	write.Comments = input.Comments
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// vlanWriteFromCreate builds a domain.VLANWrite from a create input.
func vlanWriteFromCreate(input VLANCreateInput) domain.VLANWrite {
	return domain.VLANWrite{
		VID:          input.VID,
		Name:         input.Name,
		Site:         input.Site,
		Group:        input.Group,
		Tenant:       input.Tenant,
		Status:       input.Status,
		Role:         input.Role,
		Description:  input.Description,
		Comments:     input.Comments,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// vlanWriteFromUpdate builds a domain.VLANWrite from an update input, mapping
// only the explicitly-provided (non-nil) fields so the PATCH body is minimal.
func vlanWriteFromUpdate(input VLANUpdateInput) domain.VLANWrite {
	write := domain.VLANWrite{}
	if input.VID != nil {
		write.VID = *input.VID
	}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Site = input.Site
	write.Group = input.Group
	write.Tenant = input.Tenant
	write.Status = input.Status
	write.Role = input.Role
	write.Description = input.Description
	write.Comments = input.Comments
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// virtualMachineWriteFromCreate builds a domain.VirtualMachineWrite from a
// create input.
func virtualMachineWriteFromCreate(input VirtualMachineCreateInput) domain.VirtualMachineWrite {
	return domain.VirtualMachineWrite{
		Name:         input.Name,
		Cluster:      input.Cluster,
		Role:         input.Role,
		Tenant:       input.Tenant,
		Platform:     input.Platform,
		Status:       input.Status,
		Site:         input.Site,
		VCPUs:        input.VCPUs,
		Memory:       input.Memory,
		Disk:         input.Disk,
		Comments:     input.Comments,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// virtualMachineWriteFromUpdate builds a domain.VirtualMachineWrite from an
// update input, mapping only the explicitly-provided (non-nil) fields so the
// PATCH body is minimal.
func virtualMachineWriteFromUpdate(input VirtualMachineUpdateInput) domain.VirtualMachineWrite {
	write := domain.VirtualMachineWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Cluster = input.Cluster
	write.Role = input.Role
	write.Tenant = input.Tenant
	write.Platform = input.Platform
	write.Status = input.Status
	write.Site = input.Site
	write.VCPUs = input.VCPUs
	write.Memory = input.Memory
	write.Disk = input.Disk
	write.Comments = input.Comments
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// clusterWriteFromCreate builds a domain.ClusterWrite from a create input.
func clusterWriteFromCreate(input ClusterCreateInput) domain.ClusterWrite {
	return domain.ClusterWrite{
		Name:         input.Name,
		ClusterType:  input.ClusterType,
		ClusterGroup: input.ClusterGroup,
		Site:         input.Site,
		Tenant:       input.Tenant,
		Description:  input.Description,
		Comments:     input.Comments,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// clusterWriteFromUpdate builds a domain.ClusterWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func clusterWriteFromUpdate(input ClusterUpdateInput) domain.ClusterWrite {
	write := domain.ClusterWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.ClusterType = input.ClusterType
	write.ClusterGroup = input.ClusterGroup
	write.Site = input.Site
	write.Tenant = input.Tenant
	write.Description = input.Description
	write.Comments = input.Comments
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// circuitWriteFromCreate builds a domain.CircuitWrite from a create input.
func circuitWriteFromCreate(input CircuitCreateInput) domain.CircuitWrite {
	return domain.CircuitWrite{
		CID:          input.CID,
		Provider:     input.Provider,
		CircuitType:  input.CircuitType,
		Tenant:       input.Tenant,
		Status:       input.Status,
		Description:  input.Description,
		Comments:     input.Comments,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
		InstallDate:  input.InstallDate,
		CommitRate:   input.CommitRate,
	}
}

// circuitWriteFromUpdate builds a domain.CircuitWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func circuitWriteFromUpdate(input CircuitUpdateInput) domain.CircuitWrite {
	write := domain.CircuitWrite{}
	if input.CID != nil {
		write.CID = *input.CID
	}
	write.Provider = input.Provider
	write.CircuitType = input.CircuitType
	write.Tenant = input.Tenant
	write.Status = input.Status
	write.Description = input.Description
	write.Comments = input.Comments
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields
	write.InstallDate = input.InstallDate
	write.CommitRate = input.CommitRate

	return write
}

// rackWriteFromCreate builds a domain.RackWrite from a create input.
func rackWriteFromCreate(input RackCreateInput) domain.RackWrite {
	return domain.RackWrite{
		Name:         input.Name,
		FacilityID:   input.FacilityID,
		Site:         input.Site,
		Location:     input.Location,
		Tenant:       input.Tenant,
		Status:       input.Status,
		Role:         input.Role,
		Serial:       input.Serial,
		AssetTag:     input.AssetTag,
		Type:         input.Type,
		Width:        input.Width,
		UHeight:      input.UHeight,
		Comments:     input.Comments,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// rackWriteFromUpdate builds a domain.RackWrite from an update input, mapping
// only the explicitly-provided (non-nil) fields so the PATCH body is minimal.
func rackWriteFromUpdate(input RackUpdateInput) domain.RackWrite {
	write := domain.RackWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.FacilityID = input.FacilityID
	write.Site = input.Site
	write.Location = input.Location
	write.Tenant = input.Tenant
	write.Status = input.Status
	write.Role = input.Role
	write.Serial = input.Serial
	write.AssetTag = input.AssetTag
	write.Type = input.Type
	write.Width = input.Width
	write.UHeight = input.UHeight
	write.Comments = input.Comments
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// interfaceWriteFromCreate builds a domain.InterfaceWrite from a create input.
func interfaceWriteFromCreate(input InterfaceCreateInput) domain.InterfaceWrite {
	return domain.InterfaceWrite{
		Device:       input.Device,
		Name:         input.Name,
		Type:         input.Type,
		Enabled:      input.Enabled,
		MTU:          input.MTU,
		MACAddress:   input.MACAddress,
		Speed:        input.Speed,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// interfaceWriteFromUpdate builds a domain.InterfaceWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func interfaceWriteFromUpdate(input InterfaceUpdateInput) domain.InterfaceWrite {
	write := domain.InterfaceWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Device = input.Device
	write.Type = input.Type
	write.Enabled = input.Enabled
	write.MTU = input.MTU
	write.MACAddress = input.MACAddress
	write.Speed = input.Speed
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// circuitTerminationWriteFromCreate builds a domain.CircuitTerminationWrite
// from a create input.
func circuitTerminationWriteFromCreate(input CircuitTerminationCreateInput) domain.CircuitTerminationWrite {
	return domain.CircuitTerminationWrite{
		Circuit:       input.Circuit,
		TermSide:      input.TermSide,
		Site:          input.Site,
		Speed:         input.Speed,
		UpstreamSpeed: input.UpstreamSpeed,
		Description:   input.Description,
		Tags:          input.Tags,
		CustomFields:  input.CustomFields,
	}
}

// circuitTerminationWriteFromUpdate builds a domain.CircuitTerminationWrite
// from an update input, mapping only the explicitly-provided (non-nil) fields
// so the PATCH body is minimal.
func circuitTerminationWriteFromUpdate(input CircuitTerminationUpdateInput) domain.CircuitTerminationWrite {
	write := domain.CircuitTerminationWrite{}
	if input.TermSide != nil {
		write.TermSide = *input.TermSide
	}
	write.Circuit = input.Circuit
	write.Site = input.Site
	write.Speed = input.Speed
	write.UpstreamSpeed = input.UpstreamSpeed
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// cableTerminationFromInput converts a handler termination input into a domain
// cable termination DTO.
func cableTerminationFromInput(input *CableTerminationInput) *domain.CableTerminationWrite {
	if input == nil {
		return nil
	}

	return &domain.CableTerminationWrite{
		ObjectType: input.ObjectType,
		ObjectID:   input.ObjectID,
	}
}

// cableWriteFromCreate builds a domain.CableWrite from a create input.
func cableWriteFromCreate(input CableCreateInput) domain.CableWrite {
	return domain.CableWrite{
		TerminationA: cableTerminationFromInput(input.TerminationA),
		TerminationB: cableTerminationFromInput(input.TerminationB),
		Type:         input.Type,
		Status:       input.Status,
		Label:        input.Label,
		Color:        input.Color,
		Length:       input.Length,
		LengthUnit:   input.LengthUnit,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// cableWriteFromUpdate builds a domain.CableWrite from an update input, mapping
// only the explicitly-provided (non-nil) fields so the PATCH body is minimal.
func cableWriteFromUpdate(input CableUpdateInput) domain.CableWrite {
	return domain.CableWrite{
		TerminationA: cableTerminationFromInput(input.TerminationA),
		TerminationB: cableTerminationFromInput(input.TerminationB),
		Type:         input.Type,
		Status:       input.Status,
		Label:        input.Label,
		Color:        input.Color,
		Length:       input.Length,
		LengthUnit:   input.LengthUnit,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// vmInterfaceWriteFromCreate builds a domain.VMInterfaceWrite from a create
// input.
func vmInterfaceWriteFromCreate(input VMInterfaceCreateInput) domain.VMInterfaceWrite {
	return domain.VMInterfaceWrite{
		VirtualMachine: input.VirtualMachine,
		Name:           input.Name,
		Enabled:        input.Enabled,
		MTU:            input.MTU,
		MACAddress:     input.MACAddress,
		Description:    input.Description,
		Tags:           input.Tags,
		CustomFields:   input.CustomFields,
	}
}

// vmInterfaceWriteFromUpdate builds a domain.VMInterfaceWrite from an update
// input, mapping only the explicitly-provided (non-nil) fields so the PATCH
// body is minimal.
func vmInterfaceWriteFromUpdate(input VMInterfaceUpdateInput) domain.VMInterfaceWrite {
	write := domain.VMInterfaceWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.VirtualMachine = input.VirtualMachine
	write.Enabled = input.Enabled
	write.MTU = input.MTU
	write.MACAddress = input.MACAddress
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// providerWriteFromCreate builds a domain.ProviderWrite from a create input.
func providerWriteFromCreate(input ProviderCreateInput) domain.ProviderWrite {
	return domain.ProviderWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Asn:          input.Asn,
		Account:      input.Account,
		PortalURL:    input.PortalURL,
		NocContact:   input.NocContact,
		AdminContact: input.AdminContact,
		Comments:     input.Comments,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// providerWriteFromUpdate builds a domain.ProviderWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func providerWriteFromUpdate(input ProviderUpdateInput) domain.ProviderWrite {
	write := domain.ProviderWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Asn = input.Asn
	write.Account = input.Account
	write.PortalURL = input.PortalURL
	write.NocContact = input.NocContact
	write.AdminContact = input.AdminContact
	write.Comments = input.Comments
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// tenantWriteFromCreate builds a domain.TenantWrite from a create input.
func tenantWriteFromCreate(input TenantCreateInput) domain.TenantWrite {
	return domain.TenantWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Group:        input.Group,
		Description:  input.Description,
		Comments:     input.Comments,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// tenantWriteFromUpdate builds a domain.TenantWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func tenantWriteFromUpdate(input TenantUpdateInput) domain.TenantWrite {
	write := domain.TenantWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Group = input.Group
	write.Description = input.Description
	write.Comments = input.Comments
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// tagWriteFromCreate builds a domain.TagWrite from a create input.
func tagWriteFromCreate(input TagCreateInput) domain.TagWrite {
	return domain.TagWrite{
		Name:        input.Name,
		Slug:        input.Slug,
		Color:       input.Color,
		Description: input.Description,
	}
}

// tagWriteFromUpdate builds a domain.TagWrite from an update input, mapping
// only the explicitly-provided (non-nil) fields so the PATCH body is minimal.
func tagWriteFromUpdate(input TagUpdateInput) domain.TagWrite {
	write := domain.TagWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Color = input.Color
	write.Description = input.Description

	return write
}

// manufacturerWriteFromCreate builds a domain.ManufacturerWrite from a create
// input.
func manufacturerWriteFromCreate(input ManufacturerCreateInput) domain.ManufacturerWrite {
	return domain.ManufacturerWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// manufacturerWriteFromUpdate builds a domain.ManufacturerWrite from an update
// input, mapping only the explicitly-provided (non-nil) fields so the PATCH
// body is minimal.
func manufacturerWriteFromUpdate(input ManufacturerUpdateInput) domain.ManufacturerWrite {
	write := domain.ManufacturerWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// deviceTypeWriteFromCreate builds a domain.DeviceTypeWrite from a create input.
func deviceTypeWriteFromCreate(input DeviceTypeCreateInput) domain.DeviceTypeWrite {
	return domain.DeviceTypeWrite{
		Manufacturer:  input.Manufacturer,
		Model:         input.Model,
		Slug:          input.Slug,
		PartNumber:    input.PartNumber,
		UHeight:       input.UHeight,
		IsFullDepth:   input.IsFullDepth,
		SubdeviceRole: input.SubdeviceRole,
		Comments:      input.Comments,
		Tags:          input.Tags,
		CustomFields:  input.CustomFields,
	}
}

// deviceTypeWriteFromUpdate builds a domain.DeviceTypeWrite from an update
// input, mapping only the explicitly-provided (non-nil) fields so the PATCH
// body is minimal.
func deviceTypeWriteFromUpdate(input DeviceTypeUpdateInput) domain.DeviceTypeWrite {
	write := domain.DeviceTypeWrite{}
	if input.Model != nil {
		write.Model = *input.Model
	}
	write.Manufacturer = input.Manufacturer
	write.Slug = input.Slug
	write.PartNumber = input.PartNumber
	write.UHeight = input.UHeight
	write.IsFullDepth = input.IsFullDepth
	write.SubdeviceRole = input.SubdeviceRole
	write.Comments = input.Comments
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// locationWriteFromCreate builds a domain.LocationWrite from a create input.
func locationWriteFromCreate(input LocationCreateInput) domain.LocationWrite {
	return domain.LocationWrite{
		Name:         input.Name,
		Site:         input.Site,
		Slug:         input.Slug,
		Parent:       input.Parent,
		Description:  input.Description,
		Status:       input.Status,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// locationWriteFromUpdate builds a domain.LocationWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func locationWriteFromUpdate(input LocationUpdateInput) domain.LocationWrite {
	write := domain.LocationWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Site = input.Site
	write.Slug = input.Slug
	write.Parent = input.Parent
	write.Description = input.Description
	write.Status = input.Status
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// clusterTypeWriteFromCreate builds a domain.ClusterTypeWrite from a create input.
func clusterTypeWriteFromCreate(input ClusterTypeCreateInput) domain.ClusterTypeWrite {
	return domain.ClusterTypeWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// clusterTypeWriteFromUpdate builds a domain.ClusterTypeWrite from an update
// input, mapping only the explicitly-provided (non-nil) fields so the PATCH
// body is minimal.
func clusterTypeWriteFromUpdate(input ClusterTypeUpdateInput) domain.ClusterTypeWrite {
	write := domain.ClusterTypeWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// clusterGroupWriteFromCreate builds a domain.ClusterGroupWrite from a create input.
func clusterGroupWriteFromCreate(input ClusterGroupCreateInput) domain.ClusterGroupWrite {
	return domain.ClusterGroupWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// circuitTypeWriteFromCreate builds a domain.CircuitTypeWrite from a create input.
func circuitTypeWriteFromCreate(input CircuitTypeCreateInput) domain.CircuitTypeWrite {
	return domain.CircuitTypeWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// circuitTypeWriteFromUpdate builds a domain.CircuitTypeWrite from an update
// input, mapping only the explicitly-provided (non-nil) fields so the PATCH
// body is minimal.
func circuitTypeWriteFromUpdate(input CircuitTypeUpdateInput) domain.CircuitTypeWrite {
	write := domain.CircuitTypeWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// vrfWriteFromCreate builds a domain.VrfWrite from a create input.
func vrfWriteFromCreate(input VrfCreateInput) domain.VrfWrite {
	return domain.VrfWrite{
		Name:         input.Name,
		Rd:           input.Rd,
		Tenant:       input.Tenant,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// vrfWriteFromUpdate builds a domain.VrfWrite from an update input, mapping
// only the explicitly-provided (non-nil) fields so the PATCH body is minimal.
func vrfWriteFromUpdate(input VrfUpdateInput) domain.VrfWrite {
	write := domain.VrfWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	if input.Rd != nil {
		write.Rd = *input.Rd
	}
	write.Tenant = input.Tenant
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// vlanGroupWriteFromCreate builds a domain.VlanGroupWrite from a create input.
func vlanGroupWriteFromCreate(input VlanGroupCreateInput) domain.VlanGroupWrite {
	return domain.VlanGroupWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Description:  input.Description,
		ScopeType:    input.ScopeType,
		ScopeID:      input.ScopeID,
		MinVID:       input.MinVID,
		MaxVID:       input.MaxVID,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// vlanGroupWriteFromUpdate builds a domain.VlanGroupWrite from an update input,
// roleWriteFromCreate builds a domain.RoleWrite from a create input.
func roleWriteFromCreate(input RoleCreateInput) domain.RoleWrite {
	return domain.RoleWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Weight:       input.Weight,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// roleWriteFromUpdate builds a domain.RoleWrite from an update input, mapping
// only the explicitly-provided (non-nil) fields so the PATCH body is minimal.
func roleWriteFromUpdate(input RoleUpdateInput) domain.RoleWrite {
	write := domain.RoleWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Weight = input.Weight
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// contactWriteFromCreate builds a domain.ContactWrite from a create input.
func contactWriteFromCreate(input ContactCreateInput) domain.ContactWrite {
	return domain.ContactWrite{
		Name:         input.Name,
		Group:        input.Group,
		Title:        input.Title,
		Phone:        input.Phone,
		Email:        input.Email,
		Address:      input.Address,
		Link:         input.Link,
		Description:  input.Description,
		Comments:     input.Comments,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// contactWriteFromUpdate builds a domain.ContactWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func contactWriteFromUpdate(input ContactUpdateInput) domain.ContactWrite {
	write := domain.ContactWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Group = input.Group
	write.Title = input.Title
	write.Phone = input.Phone
	write.Email = input.Email
	write.Address = input.Address
	write.Link = input.Link
	write.Description = input.Description
	write.Comments = input.Comments
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// vlanGroupWriteFromUpdate builds a domain.VlanGroupWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func vlanGroupWriteFromUpdate(input VlanGroupUpdateInput) domain.VlanGroupWrite {
	write := domain.VlanGroupWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Description = input.Description
	write.ScopeType = input.ScopeType
	write.ScopeID = input.ScopeID
	write.MinVID = input.MinVID
	write.MaxVID = input.MaxVID
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// clusterGroupWriteFromUpdate builds a domain.ClusterGroupWrite from an update
// input, mapping only the explicitly-provided (non-nil) fields so the PATCH
// body is minimal.
func clusterGroupWriteFromUpdate(input ClusterGroupUpdateInput) domain.ClusterGroupWrite {
	write := domain.ClusterGroupWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// — handler factories —

// maxPageSize is the maximum allowed page size for paginated requests.
// Increased from 100 to 1000 for better AI assistant UX (fewer pagination rounds).
const maxPageSize = 1000

func paginationParams(page, pageSize int) map[string]string {
	params := make(map[string]string)
	if page <= 0 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 25
	}
	params["offset"] = strconv.Itoa((page - 1) * pageSize)
	params["limit"] = strconv.Itoa(min(pageSize, maxPageSize))

	return params
}

func addParam(m map[string]string, key, value string) {
	if value != "" {
		m[key] = value
	}
}

func addIntParam(m map[string]string, key string, value int) {
	if value != 0 {
		m[key] = strconv.Itoa(value)
	}
}

var errServiceNotAvailable = fmt.Errorf("service not available in request context")

// resolveService returns the service from the request context if available,
// falling back to the captured service parameter. This allows per-request
// service injection (with per-request tokens) while preserving backward
// compatibility with tests that pass a service directly.
func resolveService(ctx context.Context, svc *application.NetworkService) *application.NetworkService {
	if s := ServiceFromContext(ctx); s != nil {
		return s
	}

	return svc
}

// — generic list handler —

// listHandlerConfig configures a generic list-oriented MCP tool handler.
type listHandlerConfig[I, D any] struct {
	svc         *application.NetworkService
	listFunc    func(context.Context, map[string]string) (*domain.PaginatedResponse[D], error)
	buildParams func(I) map[string]string
	errorLabel  string
}

// newListHandler creates a generic MCP tool handler for list-type tools.
// It replaces the repetitive boilerplate found in all NewGetXxxHandler functions.
func newListHandler[I, D any](cfg listHandlerConfig[I, D]) mcp.ToolHandlerFor[I, PaginatedOutput[D]] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input I) (*mcp.CallToolResult, PaginatedOutput[D], error) {
		s := resolveService(ctx, cfg.svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, PaginatedOutput[D]{}, errServiceNotAvailable
		}

		params := cfg.buildParams(input)
		resp, err := cfg.listFunc(ctx, params)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, PaginatedOutput[D]{}, fmt.Errorf("%s: %w", cfg.errorLabel, err)
		}
		if resp == nil {
			return &mcp.CallToolResult{IsError: true}, PaginatedOutput[D]{}, fmt.Errorf(
				"%s: %w", cfg.errorLabel, errors.New("received nil response from service"),
			)
		}

		return &mcp.CallToolResult{}, PaginatedOutput[D]{
			Count: resp.Count, Next: resp.Next, Previous: resp.Previous, Results: resp.Results,
		}, nil
	}
}

// — handler factories —

// NewGetSitesHandler creates a handler for the get_sites tool.
func NewGetSitesHandler(svc *application.NetworkService) mcp.ToolHandlerFor[SitesInput, PaginatedOutput[domain.Site]] {
	return newListHandler(listHandlerConfig[SitesInput, domain.Site]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Site], error) {
			return resolveService(ctx, svc).ListSites(ctx, params)
		},
		buildParams: func(input SitesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "region", input.Region)
			addParam(params, "status", input.Status)
			addParam(params, "tenant", input.Tenant)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list sites",
	})
}

// NewGetRegionsHandler creates a handler for the get_regions tool.
func NewGetRegionsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[RegionsInput, PaginatedOutput[domain.Region]] {
	return newListHandler(listHandlerConfig[RegionsInput, domain.Region]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Region], error) {
			return resolveService(ctx, svc).ListRegions(ctx, params)
		},
		buildParams: func(input RegionsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "parent", input.Parent)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list regions",
	})
}

// NewGetTenantGroupsHandler creates a handler for the get_tenant_groups tool.
func NewGetTenantGroupsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[TenantGroupsInput, PaginatedOutput[domain.TenantGroup]] {
	return newListHandler(listHandlerConfig[TenantGroupsInput, domain.TenantGroup]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.TenantGroup], error) {
			return resolveService(ctx, svc).ListTenantGroups(ctx, params)
		},
		buildParams: func(input TenantGroupsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "parent", input.Parent)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list tenant groups",
	})
}

// NewGetRolesHandler creates a handler for the get_roles tool.
func NewGetRolesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[RolesInput, PaginatedOutput[domain.Role]] {
	return newListHandler(listHandlerConfig[RolesInput, domain.Role]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Role], error) {
			return resolveService(ctx, svc).ListRoles(ctx, params)
		},
		buildParams: func(input RolesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list roles",
	})
}

// NewGetTenantsHandler creates a handler for the get_tenants tool.
func NewGetTenantsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[TenantsInput, PaginatedOutput[domain.Tenant]] {
	return newListHandler(listHandlerConfig[TenantsInput, domain.Tenant]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Tenant], error) {
			return resolveService(ctx, svc).ListTenants(ctx, params)
		},
		buildParams: func(input TenantsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "group", input.Group)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list tenants",
	})
}

// NewGetTagsHandler creates a handler for the get_tags tool.
func NewGetTagsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[TagsInput, PaginatedOutput[domain.Tag]] {
	return newListHandler(listHandlerConfig[TagsInput, domain.Tag]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Tag], error) {
			return resolveService(ctx, svc).ListTags(ctx, params)
		},
		buildParams: func(input TagsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list tags",
	})
}

// NewGetSiteGroupsHandler creates a handler for the get_site_groups tool.
func NewGetSiteGroupsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[SiteGroupsInput, PaginatedOutput[domain.SiteGroup]] {
	return newListHandler(listHandlerConfig[SiteGroupsInput, domain.SiteGroup]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.SiteGroup], error) {
			return resolveService(ctx, svc).ListSiteGroups(ctx, params)
		},
		buildParams: func(input SiteGroupsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "parent", input.Parent)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list site groups",
	})
}

// NewGetRackGroupsHandler creates a handler for the get_rack_groups tool.
func NewGetRackGroupsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[RackGroupsInput, PaginatedOutput[domain.RackGroup]] {
	return newListHandler(listHandlerConfig[RackGroupsInput, domain.RackGroup]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.RackGroup], error) {
			return resolveService(ctx, svc).ListRackGroups(ctx, params)
		},
		buildParams: func(input RackGroupsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "site", input.Site)
			addParam(params, "parent", input.Parent)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list rack groups",
	})
}

// NewGetRackRolesHandler creates a handler for the get_rack_roles tool.
func NewGetRackRolesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[RackRolesInput, PaginatedOutput[domain.RackRole]] {
	return newListHandler(listHandlerConfig[RackRolesInput, domain.RackRole]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.RackRole], error) {
			return resolveService(ctx, svc).ListRackRoles(ctx, params)
		},
		buildParams: func(input RackRolesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "color", input.Color)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list rack roles",
	})
}

// NewGetDeviceRolesHandler creates a handler for the get_device_roles tool.
func NewGetDeviceRolesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[DeviceRolesInput, PaginatedOutput[domain.DeviceRole]] {
	return newListHandler(listHandlerConfig[DeviceRolesInput, domain.DeviceRole]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.DeviceRole], error) {
			return resolveService(ctx, svc).ListDeviceRoles(ctx, params)
		},
		buildParams: func(input DeviceRolesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "color", input.Color)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list device roles",
	})
}

// NewGetPlatformsHandler creates a handler for the get_platforms tool.
func NewGetPlatformsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[PlatformsInput, PaginatedOutput[domain.Platform]] {
	return newListHandler(listHandlerConfig[PlatformsInput, domain.Platform]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Platform], error) {
			return resolveService(ctx, svc).ListPlatforms(ctx, params)
		},
		buildParams: func(input PlatformsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "manufacturer", input.Manufacturer)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list platforms",
	})
}

// NewGetVirtualChassisHandler creates a handler for the get_virtual_chassis
// tool.
func NewGetVirtualChassisHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualChassisInput, PaginatedOutput[domain.VirtualChassis]] {
	return newListHandler(listHandlerConfig[VirtualChassisInput, domain.VirtualChassis]{
		svc: svc,
		listFunc: func(
			ctx context.Context, params map[string]string,
		) (*domain.PaginatedResponse[domain.VirtualChassis], error) {
			return resolveService(ctx, svc).ListVirtualChassis(ctx, params)
		},
		buildParams: func(input VirtualChassisInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "name", input.Name)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list virtual chassis",
	})
}

// NewGetPowerPanelsHandler creates a handler for the get_power_panels tool.
func NewGetPowerPanelsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[PowerPanelsInput, PaginatedOutput[domain.PowerPanel]] {
	return newListHandler(listHandlerConfig[PowerPanelsInput, domain.PowerPanel]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.PowerPanel], error) {
			return resolveService(ctx, svc).ListPowerPanels(ctx, params)
		},
		buildParams: func(input PowerPanelsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "site", input.Site)
			addParam(params, "location", input.Location)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list power panels",
	})
}

// NewGetInventoryItemRolesHandler creates a handler for the
// get_inventory_item_roles tool.
func NewGetInventoryItemRolesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[InventoryItemRolesInput, PaginatedOutput[domain.InventoryItemRole]] {
	return newListHandler(listHandlerConfig[InventoryItemRolesInput, domain.InventoryItemRole]{
		svc: svc,
		listFunc: func(
			ctx context.Context, params map[string]string,
		) (*domain.PaginatedResponse[domain.InventoryItemRole], error) {
			return resolveService(ctx, svc).ListInventoryItemRoles(ctx, params)
		},
		buildParams: func(input InventoryItemRolesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "color", input.Color)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list inventory item roles",
	})
}

// NewGetDevicesHandler creates a handler for the get_devices tool.
func NewGetDevicesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[DevicesInput, PaginatedOutput[domain.Device]] {
	return newListHandler(listHandlerConfig[DevicesInput, domain.Device]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Device], error) {
			return resolveService(ctx, svc).ListDevices(ctx, params)
		},
		buildParams: func(input DevicesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "site", input.Site)
			addParam(params, "role", input.Role)
			addParam(params, "manufacturer", input.Manufacturer)
			addParam(params, "device_type", input.DeviceType)
			addParam(params, "status", input.Status)
			addParam(params, "name__ic", input.Name)
			addParam(params, "tenant", input.Tenant)
			addParam(params, "rack", input.Rack)
			addParam(params, "cluster", input.Cluster)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list devices",
	})
}

// NewGetIPAddressesHandler creates a handler for the get_ip_addresses tool.
func NewGetIPAddressesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[IPAddressesInput, PaginatedOutput[domain.IPAddress]] {
	return newListHandler(listHandlerConfig[IPAddressesInput, domain.IPAddress]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.IPAddress], error) {
			return resolveService(ctx, svc).ListIPAddresses(ctx, params)
		},
		buildParams: func(input IPAddressesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "address", input.Address)
			addParam(params, "device", input.Device)
			addParam(params, "status", input.Status)
			addParam(params, "vrf", input.VRF)
			addParam(params, "role", input.Role)
			addParam(params, "tenant", input.Tenant)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list IP addresses",
	})
}

// NewGetPrefixesHandler creates a handler for the get_prefixes tool.
func NewGetPrefixesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[PrefixesInput, PaginatedOutput[domain.Prefix]] {
	return newListHandler(listHandlerConfig[PrefixesInput, domain.Prefix]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Prefix], error) {
			return resolveService(ctx, svc).ListPrefixes(ctx, params)
		},
		buildParams: func(input PrefixesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "prefix", input.Prefix)
			addParam(params, "site", input.Site)
			addParam(params, "vrf", input.VRF)
			addParam(params, "status", input.Status)
			addParam(params, "role", input.Role)
			addParam(params, "tenant", input.Tenant)
			addParam(params, "within", input.Within)
			if input.Family != nil {
				addParam(params, "family", strconv.Itoa(*input.Family))
			}
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list prefixes",
	})
}

// NewGetVLANsHandler creates a handler for the get_vlans tool.
func NewGetVLANsHandler(svc *application.NetworkService) mcp.ToolHandlerFor[VLANsInput, PaginatedOutput[domain.VLAN]] {
	return newListHandler(listHandlerConfig[VLANsInput, domain.VLAN]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.VLAN], error) {
			return resolveService(ctx, svc).ListVLANs(ctx, params)
		},
		buildParams: func(input VLANsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "site", input.Site)
			addParam(params, "group", input.Group)
			addParam(params, "status", input.Status)
			addParam(params, "tenant", input.Tenant)
			addIntParam(params, "vid", input.VID)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list VLANs",
	})
}

// NewGetVirtualMachinesHandler creates a handler for the get_virtual_machines tool.
func NewGetVirtualMachinesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualMachinesInput, PaginatedOutput[domain.VirtualMachine]] {
	return newListHandler(listHandlerConfig[VirtualMachinesInput, domain.VirtualMachine]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (
			*domain.PaginatedResponse[domain.VirtualMachine], error,
		) {
			return resolveService(ctx, svc).ListVirtualMachines(ctx, params)
		},
		buildParams: func(input VirtualMachinesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "cluster", input.Cluster)
			addParam(params, "cluster_group", input.ClusterGroup)
			addParam(params, "role", input.Role)
			addParam(params, "status", input.Status)
			addParam(params, "tenant", input.Tenant)
			addParam(params, "name__ic", input.Name)
			addParam(params, "site", input.Site)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list VMs",
	})
}

// NewGetClustersHandler creates a handler for the get_clusters tool.
func NewGetClustersHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ClustersInput, PaginatedOutput[domain.Cluster]] {
	return newListHandler(listHandlerConfig[ClustersInput, domain.Cluster]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Cluster], error) {
			return resolveService(ctx, svc).ListClusters(ctx, params)
		},
		buildParams: func(input ClustersInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "type", input.ClusterType)
			addParam(params, "group", input.ClusterGroup)
			addParam(params, "site", input.Site)
			addParam(params, "tenant", input.Tenant)
			addParam(params, "name__ic", input.Name)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list clusters",
	})
}

// NewGetCircuitsHandler creates a handler for the get_circuits tool.
func NewGetCircuitsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[CircuitsInput, PaginatedOutput[domain.Circuit]] {
	return newListHandler(listHandlerConfig[CircuitsInput, domain.Circuit]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Circuit], error) {
			return resolveService(ctx, svc).ListCircuits(ctx, params)
		},
		buildParams: func(input CircuitsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "provider", input.Provider)
			addParam(params, "type", input.CircuitType)
			addParam(params, "site", input.Site)
			addParam(params, "status", input.Status)
			addParam(params, "tenant", input.Tenant)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list circuits",
	})
}

// NewGetRacksHandler creates a handler for the get_racks tool.
func NewGetRacksHandler(svc *application.NetworkService) mcp.ToolHandlerFor[RacksInput, PaginatedOutput[domain.Rack]] {
	return newListHandler(listHandlerConfig[RacksInput, domain.Rack]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Rack], error) {
			return resolveService(ctx, svc).ListRacks(ctx, params)
		},
		buildParams: func(input RacksInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "site", input.Site)
			addParam(params, "location", input.Location)
			addParam(params, "status", input.Status)
			addParam(params, "tenant", input.Tenant)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list racks",
	})
}

// NewGetInterfacesHandler creates a handler for the get_interfaces tool.
func NewGetInterfacesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[InterfacesInput, PaginatedOutput[domain.Interface]] {
	return newListHandler(listHandlerConfig[InterfacesInput, domain.Interface]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Interface], error) {
			return resolveService(ctx, svc).ListInterfaces(ctx, params)
		},
		buildParams: func(input InterfacesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "device", input.Device)
			addParam(params, "type", input.Type)
			addParam(params, "name__ic", input.Name)
			addParam(params, "tag", input.Tag)
			if input.Enabled != nil {
				addParam(params, "enabled", fmt.Sprintf("%t", *input.Enabled))
			}

			return params
		},
		errorLabel: "list interfaces",
	})
}

// NewGetVMInterfacesHandler creates a handler for the get_vm_interfaces tool.
func NewGetVMInterfacesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VMInterfacesInput, PaginatedOutput[domain.VMInterface]] {
	return newListHandler(listHandlerConfig[VMInterfacesInput, domain.VMInterface]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.VMInterface], error) {
			return resolveService(ctx, svc).ListVMInterfaces(ctx, params)
		},
		buildParams: func(input VMInterfacesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "virtual_machine", input.VirtualMachine)
			addParam(params, "name__ic", input.Name)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list VM interfaces",
	})
}

// NewGetCircuitTerminationsHandler creates a handler for the get_circuit_terminations tool.
func NewGetCircuitTerminationsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[CircuitTerminationsInput, PaginatedOutput[domain.CircuitTermination]] {
	return newListHandler(listHandlerConfig[CircuitTerminationsInput, domain.CircuitTermination]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (
			*domain.PaginatedResponse[domain.CircuitTermination], error,
		) {
			return resolveService(ctx, svc).ListCircuitTerminations(ctx, params)
		},
		buildParams: func(input CircuitTerminationsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "circuit", input.Circuit)
			addParam(params, "site", input.Site)
			addParam(params, "tag", input.Tag)
			addParam(params, "term_side", input.TermSide)

			return params
		},
		errorLabel: "list circuit terminations",
	})
}

// NewGetCablesHandler creates a handler for the get_cables tool.
func NewGetCablesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[CablesInput, PaginatedOutput[domain.Cable]] {
	return newListHandler(listHandlerConfig[CablesInput, domain.Cable]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Cable], error) {
			return resolveService(ctx, svc).ListCables(ctx, params)
		},
		buildParams: func(input CablesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "type", input.Type)
			addParam(params, "status", input.Status)
			addParam(params, "site", input.Site)
			addParam(params, "tag", input.Tag)
			addParam(params, "color", input.Color)
			addParam(params, "label__ic", input.Label)

			return params
		},
		errorLabel: "list cables",
	})
}

// NewGetObjectByIDHandler creates a handler for the get_object_by_id tool.
func NewGetObjectByIDHandler(svc *application.NetworkService) mcp.ToolHandlerFor[GetObjectInput, GetObjectOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input GetObjectInput,
	) (*mcp.CallToolResult, GetObjectOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, GetObjectOutput{}, errServiceNotAvailable
		}

		if input.ObjectType == "" {
			return &mcp.CallToolResult{IsError: true}, GetObjectOutput{}, fmt.Errorf("object_type is required")
		}
		if input.ID <= 0 {
			return &mcp.CallToolResult{IsError: true}, GetObjectOutput{}, fmt.Errorf("id must be a positive integer")
		}

		// Sanitize params to prevent CRLF injection
		sanitizedParams := make(map[string]string, len(input.Params))
		for key, val := range input.Params {
			sanitizedKey := strings.Map(func(r rune) rune {
				if r == '\r' || r == '\n' {
					return -1
				}

				return r
			}, key)
			sanitizedVal := strings.Map(func(r rune) rune {
				if r == '\r' || r == '\n' {
					return -1
				}

				return r
			}, val)
			sanitizedParams[sanitizedKey] = sanitizedVal
		}

		data, err := service.GetObject(ctx, input.ObjectType, input.ID, sanitizedParams)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, GetObjectOutput{}, fmt.Errorf("get object: %w", err)
		}

		return &mcp.CallToolResult{}, GetObjectOutput{Data: data}, nil
	}
}

// NewCreateSiteHandler creates a handler for the create_site tool.
func NewCreateSiteHandler(svc *application.NetworkService) mcp.ToolHandlerFor[SiteCreateInput, SiteOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input SiteCreateInput,
	) (*mcp.CallToolResult, SiteOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, SiteOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, SiteOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		site, err := service.CreateSite(ctx, siteWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, SiteOutput{}, e
		}

		return &mcp.CallToolResult{}, SiteOutput{Data: *site}, nil
	}
}

// NewUpdateSiteHandler creates a handler for the update_site tool. It performs a
// partial-merge PATCH using only the explicitly provided fields.
func NewUpdateSiteHandler(svc *application.NetworkService) mcp.ToolHandlerFor[SiteUpdateInput, SiteOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input SiteUpdateInput,
	) (*mcp.CallToolResult, SiteOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, SiteOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, SiteOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		site, err := service.UpdateSite(ctx, input.ID, siteWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, SiteOutput{}, e
		}

		return &mcp.CallToolResult{}, SiteOutput{Data: *site}, nil
	}
}

// NewDeleteSiteHandler creates a handler for the delete_site tool.
func NewDeleteSiteHandler(svc *application.NetworkService) mcp.ToolHandlerFor[SiteDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input SiteDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteSite(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateRegionHandler creates a handler for the create_region tool.
func NewCreateRegionHandler(svc *application.NetworkService) mcp.ToolHandlerFor[RegionCreateInput, RegionOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RegionCreateInput,
	) (*mcp.CallToolResult, RegionOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, RegionOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, RegionOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		region, err := service.CreateRegion(ctx, regionWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, RegionOutput{}, e
		}

		return &mcp.CallToolResult{}, RegionOutput{Data: *region}, nil
	}
}

// NewUpdateRegionHandler creates a handler for the update_region tool. It
// performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateRegionHandler(svc *application.NetworkService) mcp.ToolHandlerFor[RegionUpdateInput, RegionOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RegionUpdateInput,
	) (*mcp.CallToolResult, RegionOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, RegionOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, RegionOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		region, err := service.UpdateRegion(ctx, input.ID, regionWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, RegionOutput{}, e
		}

		return &mcp.CallToolResult{}, RegionOutput{Data: *region}, nil
	}
}

// NewDeleteRegionHandler creates a handler for the delete_region tool.
func NewDeleteRegionHandler(svc *application.NetworkService) mcp.ToolHandlerFor[RegionDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RegionDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteRegion(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateTenantGroupHandler creates a handler for the create_tenant_group
// tool.
func NewCreateTenantGroupHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[TenantGroupCreateInput, TenantGroupOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input TenantGroupCreateInput,
	) (*mcp.CallToolResult, TenantGroupOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, TenantGroupOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, TenantGroupOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.CreateTenantGroup(ctx, tenantGroupWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, TenantGroupOutput{}, e
		}

		return &mcp.CallToolResult{}, TenantGroupOutput{Data: *result}, nil
	}
}

// NewUpdateTenantGroupHandler creates a handler for the update_tenant_group
// tool. It performs a partial-merge PATCH using only the explicitly provided
// fields.
func NewUpdateTenantGroupHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[TenantGroupUpdateInput, TenantGroupOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input TenantGroupUpdateInput,
	) (*mcp.CallToolResult, TenantGroupOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, TenantGroupOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, TenantGroupOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.UpdateTenantGroup(ctx, input.ID, tenantGroupWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, TenantGroupOutput{}, e
		}

		return &mcp.CallToolResult{}, TenantGroupOutput{Data: *result}, nil
	}
}

// NewDeleteTenantGroupHandler creates a handler for the delete_tenant_group
// tool.
func NewDeleteTenantGroupHandler(svc *application.NetworkService) mcp.ToolHandlerFor[TenantGroupDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input TenantGroupDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteTenantGroup(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateDeviceHandler creates a handler for the create_device tool.
func NewCreateDeviceHandler(svc *application.NetworkService) mcp.ToolHandlerFor[DeviceCreateInput, DeviceOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input DeviceCreateInput,
	) (*mcp.CallToolResult, DeviceOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, DeviceOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, DeviceOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		device, err := service.CreateDevice(ctx, deviceWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, DeviceOutput{}, e
		}

		return &mcp.CallToolResult{}, DeviceOutput{Data: *device}, nil
	}
}

// NewUpdateDeviceHandler creates a handler for the update_device tool. It
// performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateDeviceHandler(svc *application.NetworkService) mcp.ToolHandlerFor[DeviceUpdateInput, DeviceOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input DeviceUpdateInput,
	) (*mcp.CallToolResult, DeviceOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, DeviceOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, DeviceOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		device, err := service.UpdateDevice(ctx, input.ID, deviceWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, DeviceOutput{}, e
		}

		return &mcp.CallToolResult{}, DeviceOutput{Data: *device}, nil
	}
}

// NewDeleteDeviceHandler creates a handler for the delete_device tool.
func NewDeleteDeviceHandler(svc *application.NetworkService) mcp.ToolHandlerFor[DeviceDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input DeviceDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteDevice(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateIPAddressHandler creates a handler for the create_ip_address tool.
func NewCreateIPAddressHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[IPAddressCreateInput, IPAddressOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input IPAddressCreateInput,
	) (*mcp.CallToolResult, IPAddressOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, IPAddressOutput{}, errServiceNotAvailable
		}
		if input.Address == "" {
			return &mcp.CallToolResult{IsError: true}, IPAddressOutput{}, fmt.Errorf("address is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.CreateIPAddress(ctx, ipAddressWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, IPAddressOutput{}, e
		}

		return &mcp.CallToolResult{}, IPAddressOutput{Data: *result}, nil
	}
}

// NewUpdateIPAddressHandler creates a handler for the update_ip_address tool. It
// performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateIPAddressHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[IPAddressUpdateInput, IPAddressOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input IPAddressUpdateInput,
	) (*mcp.CallToolResult, IPAddressOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, IPAddressOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, IPAddressOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		ipAddr, err := service.UpdateIPAddress(ctx, input.ID, ipAddressWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, IPAddressOutput{}, e
		}

		return &mcp.CallToolResult{}, IPAddressOutput{Data: *ipAddr}, nil
	}
}

// NewDeleteIPAddressHandler creates a handler for the delete_ip_address tool.
func NewDeleteIPAddressHandler(svc *application.NetworkService) mcp.ToolHandlerFor[IPAddressDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input IPAddressDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteIPAddress(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreatePrefixHandler creates a handler for the create_prefix tool.
func NewCreatePrefixHandler(svc *application.NetworkService) mcp.ToolHandlerFor[PrefixCreateInput, PrefixOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input PrefixCreateInput,
	) (*mcp.CallToolResult, PrefixOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, PrefixOutput{}, errServiceNotAvailable
		}
		if input.Prefix == "" {
			return &mcp.CallToolResult{IsError: true}, PrefixOutput{}, fmt.Errorf("prefix is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		prefix, err := service.CreatePrefix(ctx, prefixWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, PrefixOutput{}, e
		}

		return &mcp.CallToolResult{}, PrefixOutput{Data: *prefix}, nil
	}
}

// NewUpdatePrefixHandler creates a handler for the update_prefix tool. It
// performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdatePrefixHandler(svc *application.NetworkService) mcp.ToolHandlerFor[PrefixUpdateInput, PrefixOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input PrefixUpdateInput,
	) (*mcp.CallToolResult, PrefixOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, PrefixOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, PrefixOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		prefix, err := service.UpdatePrefix(ctx, input.ID, prefixWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, PrefixOutput{}, e
		}

		return &mcp.CallToolResult{}, PrefixOutput{Data: *prefix}, nil
	}
}

// NewDeletePrefixHandler creates a handler for the delete_prefix tool.
func NewDeletePrefixHandler(svc *application.NetworkService) mcp.ToolHandlerFor[PrefixDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input PrefixDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeletePrefix(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateVLANHandler creates a handler for the create_vlan tool.
func NewCreateVLANHandler(svc *application.NetworkService) mcp.ToolHandlerFor[VLANCreateInput, VLANOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VLANCreateInput,
	) (*mcp.CallToolResult, VLANOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VLANOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, VLANOutput{}, fmt.Errorf("name is required")
		}
		if input.VID <= 0 {
			return &mcp.CallToolResult{IsError: true}, VLANOutput{}, fmt.Errorf("vid must be a positive integer")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		vlan, err := service.CreateVLAN(ctx, vlanWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VLANOutput{}, e
		}

		return &mcp.CallToolResult{}, VLANOutput{Data: *vlan}, nil
	}
}

// NewUpdateVLANHandler creates a handler for the update_vlan tool. It performs a
// partial-merge PATCH using only the explicitly provided fields.
func NewUpdateVLANHandler(svc *application.NetworkService) mcp.ToolHandlerFor[VLANUpdateInput, VLANOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VLANUpdateInput,
	) (*mcp.CallToolResult, VLANOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VLANOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, VLANOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		vlan, err := service.UpdateVLAN(ctx, input.ID, vlanWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VLANOutput{}, e
		}

		return &mcp.CallToolResult{}, VLANOutput{Data: *vlan}, nil
	}
}

// NewDeleteVLANHandler creates a handler for the delete_vlan tool.
func NewDeleteVLANHandler(svc *application.NetworkService) mcp.ToolHandlerFor[VLANDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VLANDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteVLAN(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateVirtualMachineHandler creates a handler for the
// create_virtual_machine tool.
func NewCreateVirtualMachineHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualMachineCreateInput, VirtualMachineOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VirtualMachineCreateInput,
	) (*mcp.CallToolResult, VirtualMachineOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VirtualMachineOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, VirtualMachineOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		vmRes, err := service.CreateVirtualMachine(ctx, virtualMachineWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VirtualMachineOutput{}, e
		}

		return &mcp.CallToolResult{}, VirtualMachineOutput{Data: *vmRes}, nil
	}
}

// NewUpdateVirtualMachineHandler creates a handler for the
// update_virtual_machine tool. It performs a partial-merge PATCH using only the
// explicitly provided fields.
func NewUpdateVirtualMachineHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualMachineUpdateInput, VirtualMachineOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VirtualMachineUpdateInput,
	) (*mcp.CallToolResult, VirtualMachineOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VirtualMachineOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, VirtualMachineOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		vmRes, err := service.UpdateVirtualMachine(ctx, input.ID, virtualMachineWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VirtualMachineOutput{}, e
		}

		return &mcp.CallToolResult{}, VirtualMachineOutput{Data: *vmRes}, nil
	}
}

// NewDeleteVirtualMachineHandler creates a handler for the
// delete_virtual_machine tool.
func NewDeleteVirtualMachineHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualMachineDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VirtualMachineDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteVirtualMachine(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateClusterHandler creates a handler for the create_cluster tool.
func NewCreateClusterHandler(svc *application.NetworkService) mcp.ToolHandlerFor[ClusterCreateInput, ClusterOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ClusterCreateInput,
	) (*mcp.CallToolResult, ClusterOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ClusterOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, ClusterOutput{}, fmt.Errorf("name is required")
		}
		if input.ClusterType == nil {
			return &mcp.CallToolResult{IsError: true}, ClusterOutput{}, fmt.Errorf("type (cluster type) is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		cluster, err := service.CreateCluster(ctx, clusterWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ClusterOutput{}, e
		}

		return &mcp.CallToolResult{}, ClusterOutput{Data: *cluster}, nil
	}
}

// NewUpdateClusterHandler creates a handler for the update_cluster tool. It
// performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateClusterHandler(svc *application.NetworkService) mcp.ToolHandlerFor[ClusterUpdateInput, ClusterOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ClusterUpdateInput,
	) (*mcp.CallToolResult, ClusterOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ClusterOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, ClusterOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		cluster, err := service.UpdateCluster(ctx, input.ID, clusterWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ClusterOutput{}, e
		}

		return &mcp.CallToolResult{}, ClusterOutput{Data: *cluster}, nil
	}
}

// NewDeleteClusterHandler creates a handler for the delete_cluster tool.
func NewDeleteClusterHandler(svc *application.NetworkService) mcp.ToolHandlerFor[ClusterDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ClusterDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteCluster(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateCircuitHandler creates a handler for the create_circuit tool.
func NewCreateCircuitHandler(svc *application.NetworkService) mcp.ToolHandlerFor[CircuitCreateInput, CircuitOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input CircuitCreateInput,
	) (*mcp.CallToolResult, CircuitOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitOutput{}, errServiceNotAvailable
		}
		if input.CID == "" {
			return &mcp.CallToolResult{IsError: true}, CircuitOutput{}, fmt.Errorf("cid is required")
		}
		if input.Provider == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitOutput{}, fmt.Errorf("provider is required")
		}
		if input.CircuitType == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitOutput{}, fmt.Errorf("circuit_type is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		circuit, err := service.CreateCircuit(ctx, circuitWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, CircuitOutput{}, e
		}

		return &mcp.CallToolResult{}, CircuitOutput{Data: *circuit}, nil
	}
}

// NewUpdateCircuitHandler creates a handler for the update_circuit tool. It
// performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateCircuitHandler(svc *application.NetworkService) mcp.ToolHandlerFor[CircuitUpdateInput, CircuitOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input CircuitUpdateInput,
	) (*mcp.CallToolResult, CircuitOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, CircuitOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		circuit, err := service.UpdateCircuit(ctx, input.ID, circuitWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, CircuitOutput{}, e
		}

		return &mcp.CallToolResult{}, CircuitOutput{Data: *circuit}, nil
	}
}

// NewDeleteCircuitHandler creates a handler for the delete_circuit tool.
func NewDeleteCircuitHandler(svc *application.NetworkService) mcp.ToolHandlerFor[CircuitDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input CircuitDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteCircuit(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateRackHandler creates a handler for the create_rack tool.
func NewCreateRackHandler(svc *application.NetworkService) mcp.ToolHandlerFor[RackCreateInput, RackOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RackCreateInput,
	) (*mcp.CallToolResult, RackOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, RackOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, RackOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		rack, err := service.CreateRack(ctx, rackWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, RackOutput{}, e
		}

		return &mcp.CallToolResult{}, RackOutput{Data: *rack}, nil
	}
}

// NewUpdateRackHandler creates a handler for the update_rack tool. It performs
// a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateRackHandler(svc *application.NetworkService) mcp.ToolHandlerFor[RackUpdateInput, RackOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RackUpdateInput,
	) (*mcp.CallToolResult, RackOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, RackOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, RackOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		rack, err := service.UpdateRack(ctx, input.ID, rackWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, RackOutput{}, e
		}

		return &mcp.CallToolResult{}, RackOutput{Data: *rack}, nil
	}
}

// NewDeleteRackHandler creates a handler for the delete_rack tool.
func NewDeleteRackHandler(svc *application.NetworkService) mcp.ToolHandlerFor[RackDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RackDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteRack(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateInterfaceHandler creates a handler for the create_interface tool.
func NewCreateInterfaceHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[InterfaceCreateInput, InterfaceOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input InterfaceCreateInput,
	) (*mcp.CallToolResult, InterfaceOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, InterfaceOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, InterfaceOutput{}, fmt.Errorf("name is required")
		}
		if input.Device == nil {
			return &mcp.CallToolResult{IsError: true}, InterfaceOutput{}, fmt.Errorf("device is required")
		}
		if input.Type == nil {
			return &mcp.CallToolResult{IsError: true}, InterfaceOutput{}, fmt.Errorf("type is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		iface, err := service.CreateInterface(ctx, interfaceWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, InterfaceOutput{}, e
		}

		return &mcp.CallToolResult{}, InterfaceOutput{Data: *iface}, nil
	}
}

// NewUpdateInterfaceHandler creates a handler for the update_interface tool. It
// performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateInterfaceHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[InterfaceUpdateInput, InterfaceOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input InterfaceUpdateInput,
	) (*mcp.CallToolResult, InterfaceOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, InterfaceOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, InterfaceOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		iface, err := service.UpdateInterface(ctx, input.ID, interfaceWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, InterfaceOutput{}, e
		}

		return &mcp.CallToolResult{}, InterfaceOutput{Data: *iface}, nil
	}
}

// NewDeleteInterfaceHandler creates a handler for the delete_interface tool.
func NewDeleteInterfaceHandler(svc *application.NetworkService) mcp.ToolHandlerFor[InterfaceDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input InterfaceDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteInterface(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateCircuitTerminationHandler creates a handler for the
// create_circuit_termination tool.
func NewCreateCircuitTerminationHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[CircuitTerminationCreateInput, CircuitTerminationOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input CircuitTerminationCreateInput,
	) (*mcp.CallToolResult, CircuitTerminationOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitTerminationOutput{}, errServiceNotAvailable
		}
		if input.TermSide == "" {
			return &mcp.CallToolResult{IsError: true}, CircuitTerminationOutput{}, fmt.Errorf("term_side is required")
		}
		if input.Circuit == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitTerminationOutput{}, fmt.Errorf("circuit is required")
		}
		if input.Site == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitTerminationOutput{}, fmt.Errorf("site is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.CreateCircuitTermination(ctx, circuitTerminationWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, CircuitTerminationOutput{}, e
		}

		return &mcp.CallToolResult{}, CircuitTerminationOutput{Data: *result}, nil
	}
}

// NewUpdateCircuitTerminationHandler creates a handler for the
// update_circuit_termination tool. It performs a partial-merge PATCH using only
// the explicitly provided fields.
func NewUpdateCircuitTerminationHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[CircuitTerminationUpdateInput, CircuitTerminationOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input CircuitTerminationUpdateInput,
	) (*mcp.CallToolResult, CircuitTerminationOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitTerminationOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, CircuitTerminationOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		ctRes, err := service.UpdateCircuitTermination(ctx, input.ID, circuitTerminationWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, CircuitTerminationOutput{}, e
		}

		return &mcp.CallToolResult{}, CircuitTerminationOutput{Data: *ctRes}, nil
	}
}

// NewDeleteCircuitTerminationHandler creates a handler for the
// delete_circuit_termination tool.
func NewDeleteCircuitTerminationHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[CircuitTerminationDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input CircuitTerminationDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteCircuitTermination(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateCableHandler creates a handler for the create_cable tool.
func NewCreateCableHandler(svc *application.NetworkService) mcp.ToolHandlerFor[CableCreateInput, CableOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input CableCreateInput,
	) (*mcp.CallToolResult, CableOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, CableOutput{}, errServiceNotAvailable
		}
		if input.TerminationA == nil {
			return &mcp.CallToolResult{IsError: true}, CableOutput{}, fmt.Errorf("termination_a is required")
		}
		if input.TerminationB == nil {
			return &mcp.CallToolResult{IsError: true}, CableOutput{}, fmt.Errorf("termination_b is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		cable, err := service.CreateCable(ctx, cableWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, CableOutput{}, e
		}

		return &mcp.CallToolResult{}, CableOutput{Data: *cable}, nil
	}
}

// NewUpdateCableHandler creates a handler for the update_cable tool. It performs
// a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateCableHandler(svc *application.NetworkService) mcp.ToolHandlerFor[CableUpdateInput, CableOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input CableUpdateInput,
	) (*mcp.CallToolResult, CableOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, CableOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, CableOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.UpdateCable(ctx, input.ID, cableWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, CableOutput{}, e
		}

		return &mcp.CallToolResult{}, CableOutput{Data: *result}, nil
	}
}

// NewDeleteCableHandler creates a handler for the delete_cable tool.
func NewDeleteCableHandler(svc *application.NetworkService) mcp.ToolHandlerFor[CableDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input CableDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteCable(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateVMInterfaceHandler creates a handler for the create_vm_interface
// tool.
func NewCreateVMInterfaceHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VMInterfaceCreateInput, VMInterfaceOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VMInterfaceCreateInput,
	) (*mcp.CallToolResult, VMInterfaceOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VMInterfaceOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, VMInterfaceOutput{}, fmt.Errorf("name is required")
		}
		if input.VirtualMachine == nil {
			return &mcp.CallToolResult{IsError: true}, VMInterfaceOutput{}, fmt.Errorf("virtual_machine is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.CreateVMInterface(ctx, vmInterfaceWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VMInterfaceOutput{}, e
		}

		return &mcp.CallToolResult{}, VMInterfaceOutput{Data: *result}, nil
	}
}

// NewUpdateVMInterfaceHandler creates a handler for the update_vm_interface
// tool. It performs a partial-merge PATCH using only the explicitly provided
// fields.
func NewUpdateVMInterfaceHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VMInterfaceUpdateInput, VMInterfaceOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VMInterfaceUpdateInput,
	) (*mcp.CallToolResult, VMInterfaceOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VMInterfaceOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, VMInterfaceOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		vmIf, err := service.UpdateVMInterface(ctx, input.ID, vmInterfaceWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VMInterfaceOutput{}, e
		}

		return &mcp.CallToolResult{}, VMInterfaceOutput{Data: *vmIf}, nil
	}
}

// NewDeleteVMInterfaceHandler creates a handler for the delete_vm_interface
// tool.
func NewDeleteVMInterfaceHandler(svc *application.NetworkService) mcp.ToolHandlerFor[VMInterfaceDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VMInterfaceDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteVMInterface(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateProviderHandler creates a handler for the create_provider tool.
func NewCreateProviderHandler(svc *application.NetworkService) mcp.ToolHandlerFor[ProviderCreateInput, ProviderOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ProviderCreateInput,
	) (*mcp.CallToolResult, ProviderOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ProviderOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, ProviderOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		provider, err := service.CreateProvider(ctx, providerWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ProviderOutput{}, e
		}

		return &mcp.CallToolResult{}, ProviderOutput{Data: *provider}, nil
	}
}

// NewUpdateProviderHandler creates a handler for the update_provider tool. It
// performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateProviderHandler(svc *application.NetworkService) mcp.ToolHandlerFor[ProviderUpdateInput, ProviderOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ProviderUpdateInput,
	) (*mcp.CallToolResult, ProviderOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ProviderOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, ProviderOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.UpdateProvider(ctx, input.ID, providerWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ProviderOutput{}, e
		}

		return &mcp.CallToolResult{}, ProviderOutput{Data: *result}, nil
	}
}

// NewDeleteProviderHandler creates a handler for the delete_provider tool.
func NewDeleteProviderHandler(svc *application.NetworkService) mcp.ToolHandlerFor[ProviderDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ProviderDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteProvider(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateTenantHandler creates a handler for the create_tenant tool.
func NewCreateTenantHandler(svc *application.NetworkService) mcp.ToolHandlerFor[TenantCreateInput, TenantOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input TenantCreateInput,
	) (*mcp.CallToolResult, TenantOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, TenantOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, TenantOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.CreateTenant(ctx, tenantWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, TenantOutput{}, e
		}

		return &mcp.CallToolResult{}, TenantOutput{Data: *result}, nil
	}
}

// NewUpdateTenantHandler creates a handler for the update_tenant tool. It
// performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateTenantHandler(svc *application.NetworkService) mcp.ToolHandlerFor[TenantUpdateInput, TenantOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input TenantUpdateInput,
	) (*mcp.CallToolResult, TenantOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, TenantOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, TenantOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.UpdateTenant(ctx, input.ID, tenantWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, TenantOutput{}, e
		}

		return &mcp.CallToolResult{}, TenantOutput{Data: *result}, nil
	}
}

// NewDeleteTenantHandler creates a handler for the delete_tenant tool.
func NewDeleteTenantHandler(svc *application.NetworkService) mcp.ToolHandlerFor[TenantDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input TenantDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteTenant(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateTagHandler creates a handler for the create_tag tool.
func NewCreateTagHandler(svc *application.NetworkService) mcp.ToolHandlerFor[TagCreateInput, TagOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input TagCreateInput,
	) (*mcp.CallToolResult, TagOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, TagOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, TagOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.CreateTag(ctx, tagWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, TagOutput{}, e
		}

		return &mcp.CallToolResult{}, TagOutput{Data: *result}, nil
	}
}

// NewUpdateTagHandler creates a handler for the update_tag tool. It performs a
// partial-merge PATCH using only the explicitly provided fields.
func NewUpdateTagHandler(svc *application.NetworkService) mcp.ToolHandlerFor[TagUpdateInput, TagOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input TagUpdateInput,
	) (*mcp.CallToolResult, TagOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, TagOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, TagOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.UpdateTag(ctx, input.ID, tagWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, TagOutput{}, e
		}

		return &mcp.CallToolResult{}, TagOutput{Data: *result}, nil
	}
}

// NewDeleteTagHandler creates a handler for the delete_tag tool.
func NewDeleteTagHandler(svc *application.NetworkService) mcp.ToolHandlerFor[TagDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input TagDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteTag(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateManufacturerHandler creates a handler for the create_manufacturer
// tool.
func NewCreateManufacturerHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ManufacturerCreateInput, ManufacturerOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ManufacturerCreateInput,
	) (*mcp.CallToolResult, ManufacturerOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ManufacturerOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, ManufacturerOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		manufacturer, err := service.CreateManufacturer(ctx, manufacturerWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ManufacturerOutput{}, e
		}

		return &mcp.CallToolResult{}, ManufacturerOutput{Data: *manufacturer}, nil
	}
}

// NewUpdateManufacturerHandler creates a handler for the update_manufacturer
// tool. It performs a partial-merge PATCH using only the explicitly provided
// fields.
func NewUpdateManufacturerHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ManufacturerUpdateInput, ManufacturerOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ManufacturerUpdateInput,
	) (*mcp.CallToolResult, ManufacturerOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ManufacturerOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, ManufacturerOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.UpdateManufacturer(ctx, input.ID, manufacturerWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ManufacturerOutput{}, e
		}

		return &mcp.CallToolResult{}, ManufacturerOutput{Data: *result}, nil
	}
}

// NewDeleteManufacturerHandler creates a handler for the delete_manufacturer
// tool.
func NewDeleteManufacturerHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ManufacturerDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ManufacturerDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteManufacturer(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateDeviceTypeHandler creates a handler for the create_device_type tool.
func NewCreateDeviceTypeHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[DeviceTypeCreateInput, DeviceTypeOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input DeviceTypeCreateInput,
	) (*mcp.CallToolResult, DeviceTypeOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, DeviceTypeOutput{}, errServiceNotAvailable
		}
		if input.Manufacturer == nil {
			return &mcp.CallToolResult{IsError: true}, DeviceTypeOutput{}, fmt.Errorf("manufacturer is required")
		}
		if input.Model == "" {
			return &mcp.CallToolResult{IsError: true}, DeviceTypeOutput{}, fmt.Errorf("model is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		deviceType, err := service.CreateDeviceType(ctx, deviceTypeWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, DeviceTypeOutput{}, e
		}

		return &mcp.CallToolResult{}, DeviceTypeOutput{Data: *deviceType}, nil
	}
}

// NewUpdateDeviceTypeHandler creates a handler for the update_device_type tool.
// It performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateDeviceTypeHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[DeviceTypeUpdateInput, DeviceTypeOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input DeviceTypeUpdateInput,
	) (*mcp.CallToolResult, DeviceTypeOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, DeviceTypeOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, DeviceTypeOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		deviceType, err := service.UpdateDeviceType(ctx, input.ID, deviceTypeWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, DeviceTypeOutput{}, e
		}

		return &mcp.CallToolResult{}, DeviceTypeOutput{Data: *deviceType}, nil
	}
}

// NewDeleteDeviceTypeHandler creates a handler for the delete_device_type tool.
func NewDeleteDeviceTypeHandler(svc *application.NetworkService) mcp.ToolHandlerFor[DeviceTypeDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input DeviceTypeDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteDeviceType(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateLocationHandler creates a handler for the create_location tool.
func NewCreateLocationHandler(svc *application.NetworkService) mcp.ToolHandlerFor[LocationCreateInput, LocationOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input LocationCreateInput,
	) (*mcp.CallToolResult, LocationOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, LocationOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, LocationOutput{}, fmt.Errorf("name is required")
		}
		if input.Site == nil {
			return &mcp.CallToolResult{IsError: true}, LocationOutput{}, fmt.Errorf("site is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.CreateLocation(ctx, locationWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, LocationOutput{}, e
		}

		return &mcp.CallToolResult{}, LocationOutput{Data: *result}, nil
	}
}

// NewUpdateLocationHandler creates a handler for the update_location tool.
// It performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateLocationHandler(svc *application.NetworkService) mcp.ToolHandlerFor[LocationUpdateInput, LocationOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input LocationUpdateInput,
	) (*mcp.CallToolResult, LocationOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, LocationOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, LocationOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.UpdateLocation(ctx, input.ID, locationWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, LocationOutput{}, e
		}

		return &mcp.CallToolResult{}, LocationOutput{Data: *result}, nil
	}
}

// NewDeleteLocationHandler creates a handler for the delete_location tool.
func NewDeleteLocationHandler(svc *application.NetworkService) mcp.ToolHandlerFor[LocationDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input LocationDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteLocation(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateClusterTypeHandler creates a handler for the create_cluster_type tool.
func NewCreateClusterTypeHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ClusterTypeCreateInput, ClusterTypeOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ClusterTypeCreateInput,
	) (*mcp.CallToolResult, ClusterTypeOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ClusterTypeOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, ClusterTypeOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.CreateClusterType(ctx, clusterTypeWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ClusterTypeOutput{}, e
		}

		return &mcp.CallToolResult{}, ClusterTypeOutput{Data: *result}, nil
	}
}

// NewUpdateClusterTypeHandler creates a handler for the update_cluster_type tool.
// It performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateClusterTypeHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ClusterTypeUpdateInput, ClusterTypeOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ClusterTypeUpdateInput,
	) (*mcp.CallToolResult, ClusterTypeOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ClusterTypeOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, ClusterTypeOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.UpdateClusterType(ctx, input.ID, clusterTypeWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ClusterTypeOutput{}, e
		}

		return &mcp.CallToolResult{}, ClusterTypeOutput{Data: *result}, nil
	}
}

// NewDeleteClusterTypeHandler creates a handler for the delete_cluster_type tool.
func NewDeleteClusterTypeHandler(svc *application.NetworkService) mcp.ToolHandlerFor[ClusterTypeDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ClusterTypeDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteClusterType(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateClusterGroupHandler creates a handler for the create_cluster_group tool.
func NewCreateClusterGroupHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ClusterGroupCreateInput, ClusterGroupOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ClusterGroupCreateInput,
	) (*mcp.CallToolResult, ClusterGroupOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ClusterGroupOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, ClusterGroupOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.CreateClusterGroup(ctx, clusterGroupWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ClusterGroupOutput{}, e
		}

		return &mcp.CallToolResult{}, ClusterGroupOutput{Data: *result}, nil
	}
}

// NewUpdateClusterGroupHandler creates a handler for the update_cluster_group tool.
// It performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateClusterGroupHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ClusterGroupUpdateInput, ClusterGroupOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ClusterGroupUpdateInput,
	) (*mcp.CallToolResult, ClusterGroupOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ClusterGroupOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, ClusterGroupOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.UpdateClusterGroup(ctx, input.ID, clusterGroupWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ClusterGroupOutput{}, e
		}

		return &mcp.CallToolResult{}, ClusterGroupOutput{Data: *result}, nil
	}
}

// NewDeleteClusterGroupHandler creates a handler for the delete_cluster_group tool.
func NewDeleteClusterGroupHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ClusterGroupDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ClusterGroupDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteClusterGroup(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateCircuitTypeHandler creates a handler for the create_circuit_type tool.
func NewCreateCircuitTypeHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[CircuitTypeCreateInput, CircuitTypeOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input CircuitTypeCreateInput,
	) (*mcp.CallToolResult, CircuitTypeOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitTypeOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, CircuitTypeOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.CreateCircuitType(ctx, circuitTypeWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, CircuitTypeOutput{}, e
		}

		return &mcp.CallToolResult{}, CircuitTypeOutput{Data: *result}, nil
	}
}

// NewUpdateCircuitTypeHandler creates a handler for the update_circuit_type tool.
// It performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateCircuitTypeHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[CircuitTypeUpdateInput, CircuitTypeOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input CircuitTypeUpdateInput,
	) (*mcp.CallToolResult, CircuitTypeOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitTypeOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, CircuitTypeOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		circuitType, err := service.UpdateCircuitType(ctx, input.ID, circuitTypeWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, CircuitTypeOutput{}, e
		}

		return &mcp.CallToolResult{}, CircuitTypeOutput{Data: *circuitType}, nil
	}
}

// NewDeleteCircuitTypeHandler creates a handler for the delete_circuit_type tool.
func NewDeleteCircuitTypeHandler(svc *application.NetworkService) mcp.ToolHandlerFor[CircuitTypeDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input CircuitTypeDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteCircuitType(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateVrfHandler creates a handler for the create_vrf tool.
func NewCreateVrfHandler(svc *application.NetworkService) mcp.ToolHandlerFor[VrfCreateInput, VrfOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VrfCreateInput,
	) (*mcp.CallToolResult, VrfOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VrfOutput{}, errServiceNotAvailable
		}
		if input.Name == "" || input.Rd == "" {
			return &mcp.CallToolResult{IsError: true}, VrfOutput{}, fmt.Errorf("name and rd are required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.CreateVrf(ctx, vrfWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VrfOutput{}, e
		}

		return &mcp.CallToolResult{}, VrfOutput{Data: *result}, nil
	}
}

// NewUpdateVrfHandler creates a handler for the update_vrf tool. It performs a
// partial-merge PATCH using only the explicitly provided fields.
func NewUpdateVrfHandler(svc *application.NetworkService) mcp.ToolHandlerFor[VrfUpdateInput, VrfOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VrfUpdateInput,
	) (*mcp.CallToolResult, VrfOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VrfOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, VrfOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.UpdateVrf(ctx, input.ID, vrfWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VrfOutput{}, e
		}

		return &mcp.CallToolResult{}, VrfOutput{Data: *result}, nil
	}
}

// NewDeleteVrfHandler creates a handler for the delete_vrf tool.
func NewDeleteVrfHandler(svc *application.NetworkService) mcp.ToolHandlerFor[VrfDeleteInput, struct{}] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input VrfDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteVrf(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateVlanGroupHandler creates a handler for the create_vlan_group tool.
func NewCreateVlanGroupHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VlanGroupCreateInput, VlanGroupOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VlanGroupCreateInput,
	) (*mcp.CallToolResult, VlanGroupOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VlanGroupOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, VlanGroupOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		group, err := service.CreateVlanGroup(ctx, vlanGroupWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VlanGroupOutput{}, e
		}

		return &mcp.CallToolResult{}, VlanGroupOutput{Data: *group}, nil
	}
}

// NewUpdateVlanGroupHandler creates a handler for the update_vlan_group tool.
// It performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateVlanGroupHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VlanGroupUpdateInput, VlanGroupOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VlanGroupUpdateInput,
	) (*mcp.CallToolResult, VlanGroupOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VlanGroupOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, VlanGroupOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.UpdateVlanGroup(ctx, input.ID, vlanGroupWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VlanGroupOutput{}, e
		}

		return &mcp.CallToolResult{}, VlanGroupOutput{Data: *result}, nil
	}
}

// NewDeleteVlanGroupHandler creates a handler for the delete_vlan_group tool.
func NewDeleteVlanGroupHandler(svc *application.NetworkService) mcp.ToolHandlerFor[VlanGroupDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VlanGroupDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteVlanGroup(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateRoleHandler creates a handler for the create_role tool.
func NewCreateRoleHandler(svc *application.NetworkService) mcp.ToolHandlerFor[RoleCreateInput, RoleOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RoleCreateInput,
	) (*mcp.CallToolResult, RoleOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, RoleOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, RoleOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.CreateRole(ctx, roleWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, RoleOutput{}, e
		}

		return &mcp.CallToolResult{}, RoleOutput{Data: *result}, nil
	}
}

// NewUpdateRoleHandler creates a handler for the update_role tool. It performs
// a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateRoleHandler(svc *application.NetworkService) mcp.ToolHandlerFor[RoleUpdateInput, RoleOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RoleUpdateInput,
	) (*mcp.CallToolResult, RoleOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, RoleOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, RoleOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.UpdateRole(ctx, input.ID, roleWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, RoleOutput{}, e
		}

		return &mcp.CallToolResult{}, RoleOutput{Data: *result}, nil
	}
}

// NewDeleteRoleHandler creates a handler for the delete_role tool.
func NewDeleteRoleHandler(svc *application.NetworkService) mcp.ToolHandlerFor[RoleDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RoleDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteRole(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateContactHandler creates a handler for the create_contact tool.
func NewCreateContactHandler(svc *application.NetworkService) mcp.ToolHandlerFor[ContactCreateInput, ContactOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ContactCreateInput,
	) (*mcp.CallToolResult, ContactOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ContactOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, ContactOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.CreateContact(ctx, contactWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ContactOutput{}, e
		}

		return &mcp.CallToolResult{}, ContactOutput{Data: *result}, nil
	}
}

// NewUpdateContactHandler creates a handler for the update_contact tool. It
// performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateContactHandler(svc *application.NetworkService) mcp.ToolHandlerFor[ContactUpdateInput, ContactOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ContactUpdateInput,
	) (*mcp.CallToolResult, ContactOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ContactOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, ContactOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		result, err := service.UpdateContact(ctx, input.ID, contactWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ContactOutput{}, e
		}

		return &mcp.CallToolResult{}, ContactOutput{Data: *result}, nil
	}
}

// NewDeleteContactHandler creates a handler for the delete_contact tool.
func NewDeleteContactHandler(svc *application.NetworkService) mcp.ToolHandlerFor[ContactDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ContactDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteContact(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateSiteGroupHandler creates a handler for the create_site_group tool.
func NewCreateSiteGroupHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[SiteGroupCreateInput, SiteGroupOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input SiteGroupCreateInput,
	) (*mcp.CallToolResult, SiteGroupOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, SiteGroupOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, SiteGroupOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		siteGroup, err := service.CreateSiteGroup(ctx, siteGroupWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, SiteGroupOutput{}, e
		}

		return &mcp.CallToolResult{}, SiteGroupOutput{Data: *siteGroup}, nil
	}
}

// NewUpdateSiteGroupHandler creates a handler for the update_site_group tool.
// It performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateSiteGroupHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[SiteGroupUpdateInput, SiteGroupOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input SiteGroupUpdateInput,
	) (*mcp.CallToolResult, SiteGroupOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, SiteGroupOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, SiteGroupOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		siteGroup, err := service.UpdateSiteGroup(ctx, input.ID, siteGroupWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, SiteGroupOutput{}, e
		}

		return &mcp.CallToolResult{}, SiteGroupOutput{Data: *siteGroup}, nil
	}
}

// NewDeleteSiteGroupHandler creates a handler for the delete_site_group tool.
func NewDeleteSiteGroupHandler(svc *application.NetworkService) mcp.ToolHandlerFor[SiteGroupDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input SiteGroupDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteSiteGroup(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateRackGroupHandler creates a handler for the create_rack_group tool.
func NewCreateRackGroupHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[RackGroupCreateInput, RackGroupOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RackGroupCreateInput,
	) (*mcp.CallToolResult, RackGroupOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, RackGroupOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, RackGroupOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		rackGroup, err := service.CreateRackGroup(ctx, rackGroupWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, RackGroupOutput{}, e
		}

		return &mcp.CallToolResult{}, RackGroupOutput{Data: *rackGroup}, nil
	}
}

// NewUpdateRackGroupHandler creates a handler for the update_rack_group tool.
// It performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateRackGroupHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[RackGroupUpdateInput, RackGroupOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RackGroupUpdateInput,
	) (*mcp.CallToolResult, RackGroupOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, RackGroupOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, RackGroupOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		rackGroup, err := service.UpdateRackGroup(ctx, input.ID, rackGroupWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, RackGroupOutput{}, e
		}

		return &mcp.CallToolResult{}, RackGroupOutput{Data: *rackGroup}, nil
	}
}

// NewDeleteRackGroupHandler creates a handler for the delete_rack_group tool.
func NewDeleteRackGroupHandler(svc *application.NetworkService) mcp.ToolHandlerFor[RackGroupDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RackGroupDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteRackGroup(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateRackRoleHandler creates a handler for the create_rack_role tool.
func NewCreateRackRoleHandler(svc *application.NetworkService) mcp.ToolHandlerFor[RackRoleCreateInput, RackRoleOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RackRoleCreateInput,
	) (*mcp.CallToolResult, RackRoleOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, RackRoleOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, RackRoleOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		rackRole, err := service.CreateRackRole(ctx, rackRoleWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, RackRoleOutput{}, e
		}

		return &mcp.CallToolResult{}, RackRoleOutput{Data: *rackRole}, nil
	}
}

// NewUpdateRackRoleHandler creates a handler for the update_rack_role tool. It
// performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateRackRoleHandler(svc *application.NetworkService) mcp.ToolHandlerFor[RackRoleUpdateInput, RackRoleOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RackRoleUpdateInput,
	) (*mcp.CallToolResult, RackRoleOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, RackRoleOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, RackRoleOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		rackRole, err := service.UpdateRackRole(ctx, input.ID, rackRoleWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, RackRoleOutput{}, e
		}

		return &mcp.CallToolResult{}, RackRoleOutput{Data: *rackRole}, nil
	}
}

// NewDeleteRackRoleHandler creates a handler for the delete_rack_role tool.
func NewDeleteRackRoleHandler(svc *application.NetworkService) mcp.ToolHandlerFor[RackRoleDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RackRoleDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteRackRole(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateDeviceRoleHandler creates a handler for the create_device_role tool.
func NewCreateDeviceRoleHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[DeviceRoleCreateInput, DeviceRoleOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input DeviceRoleCreateInput,
	) (*mcp.CallToolResult, DeviceRoleOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, DeviceRoleOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, DeviceRoleOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		deviceRole, err := service.CreateDeviceRole(ctx, deviceRoleWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, DeviceRoleOutput{}, e
		}

		return &mcp.CallToolResult{}, DeviceRoleOutput{Data: *deviceRole}, nil
	}
}

// NewUpdateDeviceRoleHandler creates a handler for the update_device_role tool.
// It performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateDeviceRoleHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[DeviceRoleUpdateInput, DeviceRoleOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input DeviceRoleUpdateInput,
	) (*mcp.CallToolResult, DeviceRoleOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, DeviceRoleOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, DeviceRoleOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		deviceRole, err := service.UpdateDeviceRole(ctx, input.ID, deviceRoleWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, DeviceRoleOutput{}, e
		}

		return &mcp.CallToolResult{}, DeviceRoleOutput{Data: *deviceRole}, nil
	}
}

// NewDeleteDeviceRoleHandler creates a handler for the delete_device_role tool.
func NewDeleteDeviceRoleHandler(svc *application.NetworkService) mcp.ToolHandlerFor[DeviceRoleDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input DeviceRoleDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteDeviceRole(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreatePlatformHandler creates a handler for the create_platform tool.
func NewCreatePlatformHandler(svc *application.NetworkService) mcp.ToolHandlerFor[PlatformCreateInput, PlatformOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input PlatformCreateInput,
	) (*mcp.CallToolResult, PlatformOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, PlatformOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, PlatformOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		platform, err := service.CreatePlatform(ctx, platformWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, PlatformOutput{}, e
		}

		return &mcp.CallToolResult{}, PlatformOutput{Data: *platform}, nil
	}
}

// NewUpdatePlatformHandler creates a handler for the update_platform tool. It
// performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdatePlatformHandler(svc *application.NetworkService) mcp.ToolHandlerFor[PlatformUpdateInput, PlatformOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input PlatformUpdateInput,
	) (*mcp.CallToolResult, PlatformOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, PlatformOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, PlatformOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		platform, err := service.UpdatePlatform(ctx, input.ID, platformWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, PlatformOutput{}, e
		}

		return &mcp.CallToolResult{}, PlatformOutput{Data: *platform}, nil
	}
}

// NewDeletePlatformHandler creates a handler for the delete_platform tool.
func NewDeletePlatformHandler(svc *application.NetworkService) mcp.ToolHandlerFor[PlatformDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input PlatformDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeletePlatform(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateVirtualChassisHandler creates a handler for the create_virtual_chassis
// tool.
func NewCreateVirtualChassisHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualChassisCreateInput, VirtualChassisOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VirtualChassisCreateInput,
	) (*mcp.CallToolResult, VirtualChassisOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VirtualChassisOutput{}, errServiceNotAvailable
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		virtualChassis, err := service.CreateVirtualChassis(ctx, virtualChassisWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VirtualChassisOutput{}, e
		}

		return &mcp.CallToolResult{}, VirtualChassisOutput{Data: *virtualChassis}, nil
	}
}

// NewUpdateVirtualChassisHandler creates a handler for the
// update_virtual_chassis tool. It performs a partial-merge PATCH using only the
// explicitly provided fields.
func NewUpdateVirtualChassisHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualChassisUpdateInput, VirtualChassisOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VirtualChassisUpdateInput,
	) (*mcp.CallToolResult, VirtualChassisOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VirtualChassisOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, VirtualChassisOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		virtualChassis, err := service.UpdateVirtualChassis(ctx, input.ID, virtualChassisWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VirtualChassisOutput{}, e
		}

		return &mcp.CallToolResult{}, VirtualChassisOutput{Data: *virtualChassis}, nil
	}
}

// NewDeleteVirtualChassisHandler creates a handler for the
// delete_virtual_chassis tool.
func NewDeleteVirtualChassisHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualChassisDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VirtualChassisDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteVirtualChassis(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreatePowerPanelHandler creates a handler for the create_power_panel tool.
func NewCreatePowerPanelHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[PowerPanelCreateInput, PowerPanelOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input PowerPanelCreateInput,
	) (*mcp.CallToolResult, PowerPanelOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, PowerPanelOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, PowerPanelOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		powerPanel, err := service.CreatePowerPanel(ctx, powerPanelWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, PowerPanelOutput{}, e
		}

		return &mcp.CallToolResult{}, PowerPanelOutput{Data: *powerPanel}, nil
	}
}

// NewUpdatePowerPanelHandler creates a handler for the update_power_panel tool.
// It performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdatePowerPanelHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[PowerPanelUpdateInput, PowerPanelOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input PowerPanelUpdateInput,
	) (*mcp.CallToolResult, PowerPanelOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, PowerPanelOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, PowerPanelOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		powerPanel, err := service.UpdatePowerPanel(ctx, input.ID, powerPanelWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, PowerPanelOutput{}, e
		}

		return &mcp.CallToolResult{}, PowerPanelOutput{Data: *powerPanel}, nil
	}
}

// NewDeletePowerPanelHandler creates a handler for the delete_power_panel tool.
func NewDeletePowerPanelHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[PowerPanelDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input PowerPanelDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeletePowerPanel(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateInventoryItemRoleHandler creates a handler for the
// create_inventory_item_role tool.
func NewCreateInventoryItemRoleHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[InventoryItemRoleCreateInput, InventoryItemRoleOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input InventoryItemRoleCreateInput,
	) (*mcp.CallToolResult, InventoryItemRoleOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, InventoryItemRoleOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, InventoryItemRoleOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		invItemRole, err := service.CreateInventoryItemRole(ctx, inventoryItemRoleWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, InventoryItemRoleOutput{}, e
		}

		return &mcp.CallToolResult{}, InventoryItemRoleOutput{Data: *invItemRole}, nil
	}
}

// NewUpdateInventoryItemRoleHandler creates a handler for the
// update_inventory_item_role tool. It performs a partial-merge PATCH using only
// the explicitly provided fields.
func NewUpdateInventoryItemRoleHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[InventoryItemRoleUpdateInput, InventoryItemRoleOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input InventoryItemRoleUpdateInput,
	) (*mcp.CallToolResult, InventoryItemRoleOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, InventoryItemRoleOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, InventoryItemRoleOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		invItemRole, err := service.UpdateInventoryItemRole(ctx, input.ID, inventoryItemRoleWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, InventoryItemRoleOutput{}, e
		}

		return &mcp.CallToolResult{}, InventoryItemRoleOutput{Data: *invItemRole}, nil
	}
}

// NewDeleteInventoryItemRoleHandler creates a handler for the
// delete_inventory_item_role tool.
func NewDeleteInventoryItemRoleHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[InventoryItemRoleDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input InventoryItemRoleDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteInventoryItemRole(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// FrontPortsInput represents the input fields for the get_front_ports tool.
type FrontPortsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Device   string `json:"device,omitempty" jsonschema:"filter by device (name)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// RearPortsInput represents the input fields for the get_rear_ports tool.
type RearPortsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Device   string `json:"device,omitempty" jsonschema:"filter by device (name)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// ConsolePortsInput represents the input fields for the get_console_ports tool.
type ConsolePortsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Device   string `json:"device,omitempty" jsonschema:"filter by device (name)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// ConsoleServerPortsInput represents the input fields for the
// get_console_server_ports tool.
type ConsoleServerPortsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Device   string `json:"device,omitempty" jsonschema:"filter by device (name)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// PowerPortsInput represents the input fields for the get_power_ports tool.
type PowerPortsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Device   string `json:"device,omitempty" jsonschema:"filter by device (name)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// PowerOutletsInput represents the input fields for the get_power_outlets tool.
type PowerOutletsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Device   string `json:"device,omitempty" jsonschema:"filter by device (name)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// PowerFeedsInput represents the input fields for the get_power_feeds tool.
type PowerFeedsInput struct {
	Q          string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	PowerPanel string `json:"power_panel,omitempty" jsonschema:"filter by power panel (name)"`
	Tag        string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page       int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize   int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// DeviceBaysInput represents the input fields for the get_device_bays tool.
type DeviceBaysInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Device   string `json:"device,omitempty" jsonschema:"filter by device (name)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// RackReservationsInput represents the input fields for the
// get_rack_reservations tool.
type RackReservationsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Rack     string `json:"rack,omitempty" jsonschema:"filter by rack (name)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// InventoryItemsInput represents the input fields for the get_inventory_items
// tool.
type InventoryItemsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Device   string `json:"device,omitempty" jsonschema:"filter by device (name)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// MACAddressesInput represents the input fields for the get_mac_addresses tool.
type MACAddressesInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// RIRsInput represents the input fields for the get_rirs tool.
type RIRsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// ContactGroupsInput represents the input fields for the get_contact_groups
// tool.
type ContactGroupsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Parent   string `json:"parent,omitempty" jsonschema:"filter by parent contact group ID"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// ContactRolesInput represents the input fields for the get_contact_roles tool.
type ContactRolesInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// VirtualDisksInput represents the input fields for the get_virtual_disks tool.
type VirtualDisksInput struct {
	Q              string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	VirtualMachine string `json:"virtual_machine,omitempty" jsonschema:"filter by virtual machine ID"`
	Tag            string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page           int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize       int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// VirtualMachineTypesInput represents the input fields for the
// get_virtual_machine_types tool.
type VirtualMachineTypesInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// AggregatesInput represents the input fields for the get_aggregates tool.
type AggregatesInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	RIR      string `json:"rir,omitempty" jsonschema:"filter by RIR (slug)"`
	Tenant   string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// IPRangesInput represents the input fields for the get_ip_ranges tool.
type IPRangesInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	VRF      string `json:"vrf,omitempty" jsonschema:"filter by VRF"`
	Tenant   string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug)"`
	Status   string `json:"status,omitempty" jsonschema:"filter by status"`
	Role     string `json:"role,omitempty" jsonschema:"filter by role (slug)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// ASNsInput represents the input fields for the get_asns tool.
type ASNsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	RIR      string `json:"rir,omitempty" jsonschema:"filter by RIR (slug)"`
	Tenant   string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// ASNRangesInput represents the input fields for the get_asn_ranges tool.
type ASNRangesInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	RIR      string `json:"rir,omitempty" jsonschema:"filter by RIR (slug)"`
	Tenant   string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// RouteTargetsInput represents the input fields for the get_route_targets tool.
type RouteTargetsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Tenant   string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// ServicesInput represents the input fields for the get_services tool.
type ServicesInput struct {
	Q              string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Device         string `json:"device,omitempty" jsonschema:"filter by device (name)"`
	VirtualMachine string `json:"virtual_machine,omitempty" jsonschema:"filter by virtual machine (name)"`
	Tag            string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page           int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize       int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// ServiceTemplatesInput represents the input fields for the
// get_service_templates tool.
type ServiceTemplatesInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Protocol string `json:"protocol,omitempty" jsonschema:"filter by protocol"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// FHRPGroupsInput represents the input fields for the get_fhrp_groups tool.
type FHRPGroupsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Protocol string `json:"protocol,omitempty" jsonschema:"filter by protocol"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// FHRPGroupAssignmentsInput represents the input fields for the
// get_fhrp_group_assignments tool.
type FHRPGroupAssignmentsInput struct {
	Q         string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Group     string `json:"group,omitempty" jsonschema:"filter by FHRP group ID"`
	IPAddress string `json:"ip_address,omitempty" jsonschema:"filter by IP address ID"`
	Page      int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize  int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// VLANTranslationPoliciesInput represents the input fields for the
// get_vlan_translation_policies tool.
type VLANTranslationPoliciesInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// VLANTranslationRulesInput represents the input fields for the
// get_vlan_translation_rules tool.
type VLANTranslationRulesInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Policy   string `json:"policy,omitempty" jsonschema:"filter by policy ID"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// FrontPortCreateInput represents the writable fields for creating a front port.
type FrontPortCreateInput struct {
	Device           *int           `json:"device,omitempty" jsonschema:"parent device ID (required)"`
	Name             string         `json:"name" jsonschema:"front port name (required)"`
	Type             *string        `json:"type,omitempty" jsonschema:"front port type (choice string)"`
	RearPort         *int           `json:"rear_port,omitempty" jsonschema:"rear port ID"`
	RearPortPosition *int           `json:"rear_port_position,omitempty" jsonschema:"rear port position"`
	Description      *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags             []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields     map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// FrontPortUpdateInput represents the writable fields for updating a front
// port. All fields are optional; only the explicitly provided ones are patched
// (PATCH merge).
type FrontPortUpdateInput struct {
	ID               int            `json:"id" jsonschema:"numeric ID of the front port to update,required"`
	Device           *int           `json:"device,omitempty" jsonschema:"parent device ID"`
	Name             *string        `json:"name,omitempty" jsonschema:"front port name"`
	Type             *string        `json:"type,omitempty" jsonschema:"front port type"`
	RearPort         *int           `json:"rear_port,omitempty" jsonschema:"rear port ID"`
	RearPortPosition *int           `json:"rear_port_position,omitempty" jsonschema:"rear port position"`
	Description      *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags             []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields     map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// FrontPortDeleteInput represents the input fields for deleting a front port.
type FrontPortDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the front port to delete,required"`
}

// FrontPortOutput represents the output for the create/update front port tools.
type FrontPortOutput struct {
	Data domain.FrontPort `json:"data"`
}

// RearPortCreateInput represents the writable fields for creating a rear port.
type RearPortCreateInput struct {
	Device       *int           `json:"device,omitempty" jsonschema:"parent device ID (required)"`
	Name         string         `json:"name" jsonschema:"rear port name (required)"`
	Type         *string        `json:"type,omitempty" jsonschema:"rear port type (choice string)"`
	Positions    *int           `json:"positions,omitempty" jsonschema:"number of positions"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// RearPortUpdateInput represents the writable fields for updating a rear port.
// All fields are optional; only the explicitly provided ones are patched (PATCH
// merge).
type RearPortUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the rear port to update,required"`
	Device       *int           `json:"device,omitempty" jsonschema:"parent device ID"`
	Name         *string        `json:"name,omitempty" jsonschema:"rear port name"`
	Type         *string        `json:"type,omitempty" jsonschema:"rear port type"`
	Positions    *int           `json:"positions,omitempty" jsonschema:"number of positions"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// RearPortDeleteInput represents the input fields for deleting a rear port.
type RearPortDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the rear port to delete,required"`
}

// RearPortOutput represents the output for the create/update rear port tools.
type RearPortOutput struct {
	Data domain.RearPort `json:"data"`
}

// ConsolePortCreateInput represents the writable fields for creating a console
// port.
type ConsolePortCreateInput struct {
	Device       *int           `json:"device,omitempty" jsonschema:"parent device ID (required)"`
	Name         string         `json:"name" jsonschema:"console port name (required)"`
	Type         *string        `json:"type,omitempty" jsonschema:"console port type (choice string)"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ConsolePortUpdateInput represents the writable fields for updating a console
// port. All fields are optional; only the explicitly provided ones are patched
// (PATCH merge).
type ConsolePortUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the console port to update,required"`
	Device       *int           `json:"device,omitempty" jsonschema:"parent device ID"`
	Name         *string        `json:"name,omitempty" jsonschema:"console port name"`
	Type         *string        `json:"type,omitempty" jsonschema:"console port type"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ConsolePortDeleteInput represents the input fields for deleting a console
// port.
type ConsolePortDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the console port to delete,required"`
}

// ConsolePortOutput represents the output for the create/update console port
// tools.
type ConsolePortOutput struct {
	Data domain.ConsolePort `json:"data"`
}

// ConsoleServerPortCreateInput represents the writable fields for creating a
// console server port.
type ConsoleServerPortCreateInput struct {
	Device       *int           `json:"device,omitempty" jsonschema:"parent device ID (required)"`
	Name         string         `json:"name" jsonschema:"console server port name (required)"`
	Type         *string        `json:"type,omitempty" jsonschema:"console server port type (choice string)"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ConsoleServerPortUpdateInput represents the writable fields for updating a
// console server port. All fields are optional; only the explicitly provided
// ones are patched (PATCH merge).
type ConsoleServerPortUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the console server port to update,required"`
	Device       *int           `json:"device,omitempty" jsonschema:"parent device ID"`
	Name         *string        `json:"name,omitempty" jsonschema:"console server port name"`
	Type         *string        `json:"type,omitempty" jsonschema:"console server port type"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ConsoleServerPortDeleteInput represents the input fields for deleting a
// console server port.
type ConsoleServerPortDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the console server port to delete,required"`
}

// ConsoleServerPortOutput represents the output for the create/update console
// server port tools.
type ConsoleServerPortOutput struct {
	Data domain.ConsoleServerPort `json:"data"`
}

// PowerPortCreateInput represents the writable fields for creating a power
// port.
type PowerPortCreateInput struct {
	Device        *int           `json:"device,omitempty" jsonschema:"parent device ID (required)"`
	Name          string         `json:"name" jsonschema:"power port name (required)"`
	Type          *string        `json:"type,omitempty" jsonschema:"power port type (choice string)"`
	MaximumDraw   *int           `json:"maximum_draw,omitempty" jsonschema:"maximum power draw (VA)"`
	AllocatedDraw *int           `json:"allocated_draw,omitempty" jsonschema:"allocated power draw (VA)"`
	Description   *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags          []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields  map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// PowerPortUpdateInput represents the writable fields for updating a power
// port. All fields are optional; only the explicitly provided ones are patched
// (PATCH merge).
type PowerPortUpdateInput struct {
	ID            int            `json:"id" jsonschema:"numeric ID of the power port to update,required"`
	Device        *int           `json:"device,omitempty" jsonschema:"parent device ID"`
	Name          *string        `json:"name,omitempty" jsonschema:"power port name"`
	Type          *string        `json:"type,omitempty" jsonschema:"power port type"`
	MaximumDraw   *int           `json:"maximum_draw,omitempty" jsonschema:"maximum power draw (VA)"`
	AllocatedDraw *int           `json:"allocated_draw,omitempty" jsonschema:"allocated power draw (VA)"`
	Description   *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags          []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields  map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// PowerPortDeleteInput represents the input fields for deleting a power port.
type PowerPortDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the power port to delete,required"`
}

// PowerPortOutput represents the output for the create/update power port tools.
type PowerPortOutput struct {
	Data domain.PowerPort `json:"data"`
}

// PowerOutletCreateInput represents the writable fields for creating a power
// outlet.
type PowerOutletCreateInput struct {
	Device       *int           `json:"device,omitempty" jsonschema:"parent device ID (required)"`
	Name         string         `json:"name" jsonschema:"power outlet name (required)"`
	Type         *string        `json:"type,omitempty" jsonschema:"power outlet type (choice string)"`
	PowerPort    *int           `json:"power_port,omitempty" jsonschema:"linked power port ID"`
	FeedLeg      *string        `json:"feed_leg,omitempty" jsonschema:"feed leg (choice string)"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// PowerOutletUpdateInput represents the writable fields for updating a power
// outlet. All fields are optional; only the explicitly provided ones are
// patched (PATCH merge).
type PowerOutletUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the power outlet to update,required"`
	Device       *int           `json:"device,omitempty" jsonschema:"parent device ID"`
	Name         *string        `json:"name,omitempty" jsonschema:"power outlet name"`
	Type         *string        `json:"type,omitempty" jsonschema:"power outlet type"`
	PowerPort    *int           `json:"power_port,omitempty" jsonschema:"linked power port ID"`
	FeedLeg      *string        `json:"feed_leg,omitempty" jsonschema:"feed leg (choice string)"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// PowerOutletDeleteInput represents the input fields for deleting a power
// outlet.
type PowerOutletDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the power outlet to delete,required"`
}

// PowerOutletOutput represents the output for the create/update power outlet
// tools.
type PowerOutletOutput struct {
	Data domain.PowerOutlet `json:"data"`
}

// PowerFeedCreateInput represents the writable fields for creating a power
// feed.
type PowerFeedCreateInput struct {
	Name           string         `json:"name" jsonschema:"power feed name (required)"`
	PowerPanel     *int           `json:"power_panel,omitempty" jsonschema:"power panel ID"`
	Status         *string        `json:"status,omitempty" jsonschema:"status (choice string)"`
	Type           *string        `json:"type,omitempty" jsonschema:"power feed type (choice string)"`
	Supply         *string        `json:"supply,omitempty" jsonschema:"supply (choice string)"`
	Phase          *string        `json:"phase,omitempty" jsonschema:"phase (choice string)"`
	Voltage        *int           `json:"voltage,omitempty" jsonschema:"voltage"`
	Amperage       *int           `json:"amperage,omitempty" jsonschema:"amperage"`
	MaxUtilization *int           `json:"max_utilization,omitempty" jsonschema:"maximum utilization (percent)"`
	Description    *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags           []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields   map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// PowerFeedUpdateInput represents the writable fields for updating a power
// feed. All fields are optional; only the explicitly provided ones are patched
// (PATCH merge).
type PowerFeedUpdateInput struct {
	ID             int            `json:"id" jsonschema:"numeric ID of the power feed to update,required"`
	Name           *string        `json:"name,omitempty" jsonschema:"power feed name"`
	PowerPanel     *int           `json:"power_panel,omitempty" jsonschema:"power panel ID"`
	Status         *string        `json:"status,omitempty" jsonschema:"status (choice string)"`
	Type           *string        `json:"type,omitempty" jsonschema:"power feed type (choice string)"`
	Supply         *string        `json:"supply,omitempty" jsonschema:"supply (choice string)"`
	Phase          *string        `json:"phase,omitempty" jsonschema:"phase (choice string)"`
	Voltage        *int           `json:"voltage,omitempty" jsonschema:"voltage"`
	Amperage       *int           `json:"amperage,omitempty" jsonschema:"amperage"`
	MaxUtilization *int           `json:"max_utilization,omitempty" jsonschema:"maximum utilization (percent)"`
	Description    *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags           []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields   map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// PowerFeedDeleteInput represents the input fields for deleting a power feed.
type PowerFeedDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the power feed to delete,required"`
}

// PowerFeedOutput represents the output for the create/update power feed tools.
type PowerFeedOutput struct {
	Data domain.PowerFeed `json:"data"`
}

// DeviceBayCreateInput represents the writable fields for creating a device
// bay.
type DeviceBayCreateInput struct {
	Device       *int           `json:"device,omitempty" jsonschema:"parent device ID (required)"`
	Name         string         `json:"name" jsonschema:"device bay name (required)"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// DeviceBayUpdateInput represents the writable fields for updating a device
// bay. All fields are optional; only the explicitly provided ones are patched
// (PATCH merge).
type DeviceBayUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the device bay to update,required"`
	Device       *int           `json:"device,omitempty" jsonschema:"parent device ID"`
	Name         *string        `json:"name,omitempty" jsonschema:"device bay name"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// DeviceBayDeleteInput represents the input fields for deleting a device bay.
type DeviceBayDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the device bay to delete,required"`
}

// DeviceBayOutput represents the output for the create/update device bay tools.
type DeviceBayOutput struct {
	Data domain.DeviceBay `json:"data"`
}

// RackReservationCreateInput represents the writable fields for creating a rack
// reservation.
type RackReservationCreateInput struct {
	Description  string         `json:"description" jsonschema:"description (required)"`
	Rack         *int           `json:"rack,omitempty" jsonschema:"rack ID"`
	Units        []int          `json:"units,omitempty" jsonschema:"reserved rack units"`
	User         *int           `json:"user,omitempty" jsonschema:"user ID"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// RackReservationUpdateInput represents the writable fields for updating a rack
// reservation. All fields are optional; only the explicitly provided ones are
// patched (PATCH merge).
type RackReservationUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the rack reservation to update,required"`
	Description  *string        `json:"description,omitempty" jsonschema:"description"`
	Rack         *int           `json:"rack,omitempty" jsonschema:"rack ID"`
	Units        []int          `json:"units,omitempty" jsonschema:"reserved rack units"`
	User         *int           `json:"user,omitempty" jsonschema:"user ID"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// RackReservationDeleteInput represents the input fields for deleting a rack
// reservation.
type RackReservationDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the rack reservation to delete,required"`
}

// RackReservationOutput represents the output for the create/update rack
// reservation tools.
type RackReservationOutput struct {
	Data domain.RackReservation `json:"data"`
}

// InventoryItemCreateInput represents the writable fields for creating an
// inventory item.
type InventoryItemCreateInput struct {
	Device       *int           `json:"device,omitempty" jsonschema:"parent device ID (required)"`
	Name         string         `json:"name" jsonschema:"inventory item name (required)"`
	Role         *int           `json:"role,omitempty" jsonschema:"inventory item role ID"`
	Manufacturer *int           `json:"manufacturer,omitempty" jsonschema:"manufacturer ID"`
	Serial       *string        `json:"serial,omitempty" jsonschema:"serial number"`
	PartID       *string        `json:"part_id,omitempty" jsonschema:"part ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// InventoryItemUpdateInput represents the writable fields for updating an
// inventory item. All fields are optional; only the explicitly provided ones
// are patched (PATCH merge).
type InventoryItemUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the inventory item to update,required"`
	Device       *int           `json:"device,omitempty" jsonschema:"parent device ID"`
	Name         *string        `json:"name,omitempty" jsonschema:"inventory item name"`
	Role         *int           `json:"role,omitempty" jsonschema:"inventory item role ID"`
	Manufacturer *int           `json:"manufacturer,omitempty" jsonschema:"manufacturer ID"`
	Serial       *string        `json:"serial,omitempty" jsonschema:"serial number"`
	PartID       *string        `json:"part_id,omitempty" jsonschema:"part ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// InventoryItemDeleteInput represents the input fields for deleting an
// inventory item.
type InventoryItemDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the inventory item to delete,required"`
}

// InventoryItemOutput represents the output for the create/update inventory
// item tools.
type InventoryItemOutput struct {
	Data domain.InventoryItem `json:"data"`
}

// MACAddressCreateInput represents the writable fields for creating a MAC
// address.
type MACAddressCreateInput struct {
	Address            string         `json:"address" jsonschema:"MAC address (required)"`
	AssignedObjectType *string        `json:"assigned_object_type,omitempty" jsonschema:"assigned object content type, e.g. dcim.interface"` //nolint:lll // struct tag cannot be split across lines
	AssignedObjectID   *int           `json:"assigned_object_id,omitempty" jsonschema:"assigned object ID"`
	Description        *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags               []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields       map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// MACAddressUpdateInput represents the writable fields for updating a MAC
// address. All fields are optional; only the explicitly provided ones are
// patched (PATCH merge).
type MACAddressUpdateInput struct {
	ID                 int            `json:"id" jsonschema:"numeric ID of the MAC address to update,required"`
	Address            *string        `json:"address,omitempty" jsonschema:"MAC address"`
	AssignedObjectType *string        `json:"assigned_object_type,omitempty" jsonschema:"assigned object content type, e.g. dcim.interface"` //nolint:lll // struct tag cannot be split across lines
	AssignedObjectID   *int           `json:"assigned_object_id,omitempty" jsonschema:"assigned object ID"`
	Description        *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags               []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields       map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// MACAddressDeleteInput represents the input fields for deleting a MAC address.
type MACAddressDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the MAC address to delete,required"`
}

// MACAddressOutput represents the output for the create/update MAC address
// tools.
type MACAddressOutput struct {
	Data domain.MACAddress `json:"data"`
}

// RIRCreateInput represents the writable fields for creating an RIR.
type RIRCreateInput struct {
	Name         string         `json:"name" jsonschema:"RIR name (required)"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	IsPrivate    *bool          `json:"is_private,omitempty" jsonschema:"whether this RIR represents private use"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// RIRUpdateInput represents the writable fields for updating an RIR. All fields
// are optional; only the explicitly provided ones are patched (PATCH merge).
type RIRUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the RIR to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"RIR name"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	IsPrivate    *bool          `json:"is_private,omitempty" jsonschema:"whether this RIR represents private use"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// RIRDeleteInput represents the input fields for deleting an RIR.
type RIRDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the RIR to delete,required"`
}

// RIROutput represents the output for the create/update RIR tools.
type RIROutput struct {
	Data domain.RIR `json:"data"`
}

// ContactGroupCreateInput represents the writable fields for creating a contact
// group.
type ContactGroupCreateInput struct {
	Name         string         `json:"name" jsonschema:"contact group name (required)"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Parent       *int           `json:"parent,omitempty" jsonschema:"ID of the parent contact group"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ContactGroupUpdateInput represents the writable fields for updating a contact
// group. All fields are optional; only the explicitly provided ones are patched
// (PATCH merge).
type ContactGroupUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the contact group to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"contact group name"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Parent       *int           `json:"parent,omitempty" jsonschema:"ID of the parent contact group"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ContactGroupDeleteInput represents the input fields for deleting a contact
// group.
type ContactGroupDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the contact group to delete,required"`
}

// ContactGroupOutput represents the output for the create/update contact group
// tools.
type ContactGroupOutput struct {
	Data domain.ContactGroup `json:"data"`
}

// ContactRoleCreateInput represents the writable fields for creating a contact
// role.
type ContactRoleCreateInput struct {
	Name         string         `json:"name" jsonschema:"contact role name (required)"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ContactRoleUpdateInput represents the writable fields for updating a contact
// role. All fields are optional; only the explicitly provided ones are patched
// (PATCH merge).
type ContactRoleUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the contact role to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"contact role name"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ContactRoleDeleteInput represents the input fields for deleting a contact
// role.
type ContactRoleDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the contact role to delete,required"`
}

// ContactRoleOutput represents the output for the create/update contact role
// tools.
type ContactRoleOutput struct {
	Data domain.ContactRole `json:"data"`
}

// VirtualDiskCreateInput represents the writable fields for creating a virtual
// disk.
type VirtualDiskCreateInput struct {
	Name           string         `json:"name" jsonschema:"virtual disk name (required)"`
	VirtualMachine *int           `json:"virtual_machine,omitempty" jsonschema:"ID of the parent virtual machine"`
	Size           *int           `json:"size,omitempty" jsonschema:"disk size in MB"`
	Description    *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags           []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields   map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VirtualDiskUpdateInput represents the writable fields for updating a virtual
// disk. All fields are optional; only the explicitly provided ones are patched
// (PATCH merge).
type VirtualDiskUpdateInput struct {
	ID             int            `json:"id" jsonschema:"numeric ID of the virtual disk to update,required"`
	Name           *string        `json:"name,omitempty" jsonschema:"virtual disk name"`
	VirtualMachine *int           `json:"virtual_machine,omitempty" jsonschema:"ID of the parent virtual machine"`
	Size           *int           `json:"size,omitempty" jsonschema:"disk size in MB"`
	Description    *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags           []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields   map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VirtualDiskDeleteInput represents the input fields for deleting a virtual
// disk.
type VirtualDiskDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the virtual disk to delete,required"`
}

// VirtualDiskOutput represents the output for the create/update virtual disk
// tools.
type VirtualDiskOutput struct {
	Data domain.VirtualDisk `json:"data"`
}

// VirtualMachineTypeCreateInput represents the writable fields for creating a
// virtual machine type.
type VirtualMachineTypeCreateInput struct {
	Name         string         `json:"name" jsonschema:"virtual machine type name (required)"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VirtualMachineTypeUpdateInput represents the writable fields for updating a
// virtual machine type. All fields are optional; only the explicitly provided
// ones are patched (PATCH merge).
type VirtualMachineTypeUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the virtual machine type to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"virtual machine type name"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VirtualMachineTypeDeleteInput represents the input fields for deleting a
// virtual machine type.
type VirtualMachineTypeDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the virtual machine type to delete,required"`
}

// VirtualMachineTypeOutput represents the output for the create/update virtual
// machine type tools.
type VirtualMachineTypeOutput struct {
	Data domain.VirtualMachineType `json:"data"`
}

// ProviderAccountsInput represents the input fields for the
// get_provider_accounts tool.
type ProviderAccountsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Provider string `json:"provider,omitempty" jsonschema:"filter by provider ID"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// ProviderAccountCreateInput represents the writable fields for creating a
// provider account.
type ProviderAccountCreateInput struct {
	Name         string         `json:"name" jsonschema:"provider account name (required)"`
	Provider     *int           `json:"provider,omitempty" jsonschema:"provider ID"`
	Account      *string        `json:"account,omitempty" jsonschema:"account number/identifier"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ProviderAccountUpdateInput represents the writable fields for updating a
// provider account. All fields are optional; only the explicitly provided ones
// are patched (PATCH merge).
type ProviderAccountUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the provider account to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"provider account name"`
	Provider     *int           `json:"provider,omitempty" jsonschema:"provider ID"`
	Account      *string        `json:"account,omitempty" jsonschema:"account number/identifier"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ProviderAccountDeleteInput represents the input fields for deleting a
// provider account.
type ProviderAccountDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the provider account to delete,required"`
}

// ProviderAccountOutput represents the output for the create/update provider
// account tools.
type ProviderAccountOutput struct {
	Data domain.ProviderAccount `json:"data"`
}

// ProviderNetworksInput represents the input fields for the
// get_provider_networks tool.
type ProviderNetworksInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Provider string `json:"provider,omitempty" jsonschema:"filter by provider ID"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// ProviderNetworkCreateInput represents the writable fields for creating a
// provider network.
type ProviderNetworkCreateInput struct {
	Name         string         `json:"name" jsonschema:"provider network name (required)"`
	Provider     *int           `json:"provider,omitempty" jsonschema:"provider ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"free-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ProviderNetworkUpdateInput represents the writable fields for updating a
// provider network. All fields are optional; only the explicitly provided ones
// are patched (PATCH merge).
type ProviderNetworkUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the provider network to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"provider network name"`
	Provider     *int           `json:"provider,omitempty" jsonschema:"provider ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Comments     *string        `json:"comments,omitempty" jsonschema:"free-form comments"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ProviderNetworkDeleteInput represents the input fields for deleting a
// provider network.
type ProviderNetworkDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the provider network to delete,required"`
}

// ProviderNetworkOutput represents the output for the create/update provider
// network tools.
type ProviderNetworkOutput struct {
	Data domain.ProviderNetwork `json:"data"`
}

// CircuitGroupsInput represents the input fields for the get_circuit_groups
// tool.
type CircuitGroupsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// CircuitGroupCreateInput represents the writable fields for creating a circuit
// group.
type CircuitGroupCreateInput struct {
	Name         string         `json:"name" jsonschema:"circuit group name (required)"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// CircuitGroupUpdateInput represents the writable fields for updating a circuit
// group. All fields are optional; only the explicitly provided ones are patched
// (PATCH merge).
type CircuitGroupUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the circuit group to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"circuit group name"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// CircuitGroupDeleteInput represents the input fields for deleting a circuit
// group.
type CircuitGroupDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the circuit group to delete,required"`
}

// CircuitGroupOutput represents the output for the create/update circuit group
// tools.
type CircuitGroupOutput struct {
	Data domain.CircuitGroup `json:"data"`
}

// CircuitGroupAssignmentsInput represents the input fields for the
// get_circuit_group_assignments tool.
type CircuitGroupAssignmentsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Group    string `json:"group,omitempty" jsonschema:"filter by circuit group ID"`
	Circuit  string `json:"circuit,omitempty" jsonschema:"filter by circuit ID"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// CircuitGroupAssignmentCreateInput represents the writable fields for creating
// a circuit group assignment.
type CircuitGroupAssignmentCreateInput struct {
	Group      int     `json:"group" jsonschema:"circuit group ID (required)"`
	MemberType string  `json:"member_type" jsonschema:"member content type (e.g. circuits.circuit),required"`
	MemberID   int     `json:"member_id" jsonschema:"member object ID (required)"`
	Priority   *string `json:"priority,omitempty" jsonschema:"priority (choice: primary/secondary/tertiary/inactive)"`
}

// CircuitGroupAssignmentUpdateInput represents the writable fields for updating
// a circuit group assignment. All fields are optional; only the explicitly
// provided ones are patched (PATCH merge).
type CircuitGroupAssignmentUpdateInput struct {
	ID         int     `json:"id" jsonschema:"numeric ID of the circuit group assignment to update,required"`
	Group      *int    `json:"group,omitempty" jsonschema:"circuit group ID"`
	MemberType *string `json:"member_type,omitempty" jsonschema:"member content type (e.g. circuits.circuit)"`
	MemberID   *int    `json:"member_id,omitempty" jsonschema:"member object ID"`
	Priority   *string `json:"priority,omitempty" jsonschema:"priority (choice: primary/secondary/tertiary/inactive)"`
}

// CircuitGroupAssignmentDeleteInput represents the input fields for deleting a
// circuit group assignment.
type CircuitGroupAssignmentDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the circuit group assignment to delete,required"`
}

// CircuitGroupAssignmentOutput represents the output for the create/update
// circuit group assignment tools.
type CircuitGroupAssignmentOutput struct {
	Data domain.CircuitGroupAssignment `json:"data"`
}

// VirtualCircuitsInput represents the input fields for the get_virtual_circuits
// tool.
type VirtualCircuitsInput struct {
	Q               string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Type            string `json:"type,omitempty" jsonschema:"filter by virtual circuit type ID"`
	ProviderNetwork string `json:"provider_network,omitempty" jsonschema:"filter by provider network ID"`
	Status          string `json:"status,omitempty" jsonschema:"filter by status (choice string)"`
	Tag             string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page            int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize        int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// VirtualCircuitCreateInput represents the writable fields for creating a
// virtual circuit.
type VirtualCircuitCreateInput struct {
	CID             string         `json:"cid" jsonschema:"unique circuit ID (required)"`
	Type            *int           `json:"type,omitempty" jsonschema:"virtual circuit type ID"`
	ProviderNetwork *int           `json:"provider_network,omitempty" jsonschema:"provider network ID"`
	Status          *string        `json:"status,omitempty" jsonschema:"status (choice string)"`
	Description     *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags            []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields    map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VirtualCircuitUpdateInput represents the writable fields for updating a
// virtual circuit. All fields are optional; only the explicitly provided ones
// are patched (PATCH merge).
type VirtualCircuitUpdateInput struct {
	ID              int            `json:"id" jsonschema:"numeric ID of the virtual circuit to update,required"`
	CID             *string        `json:"cid,omitempty" jsonschema:"unique circuit ID"`
	Type            *int           `json:"type,omitempty" jsonschema:"virtual circuit type ID"`
	ProviderNetwork *int           `json:"provider_network,omitempty" jsonschema:"provider network ID"`
	Status          *string        `json:"status,omitempty" jsonschema:"status (choice string)"`
	Description     *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags            []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields    map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VirtualCircuitDeleteInput represents the input fields for deleting a virtual
// circuit.
type VirtualCircuitDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the virtual circuit to delete,required"`
}

// VirtualCircuitOutput represents the output for the create/update virtual
// circuit tools.
type VirtualCircuitOutput struct {
	Data domain.VirtualCircuit `json:"data"`
}

// VirtualCircuitTypesInput represents the input fields for the
// get_virtual_circuit_types tool.
type VirtualCircuitTypesInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// VirtualCircuitTypeCreateInput represents the writable fields for creating a
// virtual circuit type.
type VirtualCircuitTypeCreateInput struct {
	Name         string         `json:"name" jsonschema:"virtual circuit type name (required)"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VirtualCircuitTypeUpdateInput represents the writable fields for updating a
// virtual circuit type. All fields are optional; only the explicitly provided
// ones are patched (PATCH merge).
type VirtualCircuitTypeUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the virtual circuit type to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"virtual circuit type name"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VirtualCircuitTypeDeleteInput represents the input fields for deleting a
// virtual circuit type.
type VirtualCircuitTypeDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the virtual circuit type to delete,required"`
}

// VirtualCircuitTypeOutput represents the output for the create/update virtual
// circuit type tools.
type VirtualCircuitTypeOutput struct {
	Data domain.VirtualCircuitType `json:"data"`
}

// VirtualCircuitTerminationsInput represents the input fields for the
// get_virtual_circuit_terminations tool.
type VirtualCircuitTerminationsInput struct {
	Q              string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	VirtualCircuit string `json:"virtual_circuit,omitempty" jsonschema:"filter by virtual circuit ID"`
	InterfaceType  string `json:"interface_type,omitempty" jsonschema:"filter by interface type (e.g. dcim.interface)"`
	InterfaceID    string `json:"interface_id,omitempty" jsonschema:"filter by interface ID"`
	Page           int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize       int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, maximum: 1000)"`
}

// VirtualCircuitTerminationCreateInput represents the writable fields for
// creating a virtual circuit termination.
type VirtualCircuitTerminationCreateInput struct {
	VirtualCircuit int     `json:"virtual_circuit" jsonschema:"virtual circuit ID (required)"`
	Role           *string `json:"role,omitempty" jsonschema:"role (choice: peer/hub/spoke)"`
	Interface      int     `json:"interface" jsonschema:"dcim interface ID (required)"`
}

// VirtualCircuitTerminationUpdateInput represents the writable fields for
// updating a virtual circuit termination. All fields are optional; only the
// explicitly provided ones are patched (PATCH merge).
type VirtualCircuitTerminationUpdateInput struct {
	ID             int     `json:"id" jsonschema:"numeric ID of the virtual circuit termination to update,required"`
	VirtualCircuit *int    `json:"virtual_circuit,omitempty" jsonschema:"virtual circuit ID"`
	Role           *string `json:"role,omitempty" jsonschema:"role (choice: peer/hub/spoke)"`
	Interface      *int    `json:"interface,omitempty" jsonschema:"dcim interface ID"`
}

// VirtualCircuitTerminationDeleteInput represents the input fields for deleting
// a virtual circuit termination.
type VirtualCircuitTerminationDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the virtual circuit termination to delete,required"`
}

// VirtualCircuitTerminationOutput represents the output for the create/update
// virtual circuit termination tools.
type VirtualCircuitTerminationOutput struct {
	Data domain.VirtualCircuitTermination `json:"data"`
}

// AggregateCreateInput represents the writable fields for creating an aggregate.
type AggregateCreateInput struct {
	Prefix       string         `json:"prefix" jsonschema:"aggregate prefix (required)"`
	RIR          *int           `json:"rir,omitempty" jsonschema:"RIR ID"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	DateAdded    *string        `json:"date_added,omitempty" jsonschema:"date the aggregate was added"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// AggregateUpdateInput represents the writable fields for updating an aggregate.
// All fields are optional; only the explicitly provided ones are patched (PATCH
// merge).
type AggregateUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the aggregate to update,required"`
	Prefix       *string        `json:"prefix,omitempty" jsonschema:"aggregate prefix"`
	RIR          *int           `json:"rir,omitempty" jsonschema:"RIR ID"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	DateAdded    *string        `json:"date_added,omitempty" jsonschema:"date the aggregate was added"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// AggregateDeleteInput represents the input fields for deleting an aggregate.
type AggregateDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the aggregate to delete,required"`
}

// AggregateOutput represents the output for the create/update aggregate tools.
type AggregateOutput struct {
	Data domain.Aggregate `json:"data"`
}

// IPRangeCreateInput represents the writable fields for creating an IP range.
type IPRangeCreateInput struct {
	StartAddress string         `json:"start_address" jsonschema:"start address (required)"`
	EndAddress   string         `json:"end_address" jsonschema:"end address (required)"`
	VRF          *int           `json:"vrf,omitempty" jsonschema:"VRF ID"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Status       *string        `json:"status,omitempty" jsonschema:"status (choice string)"`
	Role         *int           `json:"role,omitempty" jsonschema:"role ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// IPRangeUpdateInput represents the writable fields for updating an IP range.
// All fields are optional; only the explicitly provided ones are patched (PATCH
// merge).
type IPRangeUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the IP range to update,required"`
	StartAddress *string        `json:"start_address,omitempty" jsonschema:"start address"`
	EndAddress   *string        `json:"end_address,omitempty" jsonschema:"end address"`
	VRF          *int           `json:"vrf,omitempty" jsonschema:"VRF ID"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Status       *string        `json:"status,omitempty" jsonschema:"status (choice string)"`
	Role         *int           `json:"role,omitempty" jsonschema:"role ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// IPRangeDeleteInput represents the input fields for deleting an IP range.
type IPRangeDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the IP range to delete,required"`
}

// IPRangeOutput represents the output for the create/update IP range tools.
type IPRangeOutput struct {
	Data domain.IPRange `json:"data"`
}

// ASNCreateInput represents the writable fields for creating an ASN.
type ASNCreateInput struct {
	ASN          int            `json:"asn" jsonschema:"autonomous system number (required)"`
	RIR          *int           `json:"rir,omitempty" jsonschema:"RIR ID"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ASNUpdateInput represents the writable fields for updating an ASN. All fields
// are optional; only the explicitly provided ones are patched (PATCH merge).
type ASNUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the ASN to update,required"`
	ASN          *int           `json:"asn,omitempty" jsonschema:"autonomous system number"`
	RIR          *int           `json:"rir,omitempty" jsonschema:"RIR ID"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ASNDeleteInput represents the input fields for deleting an ASN.
type ASNDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the ASN to delete,required"`
}

// ASNOutput represents the output for the create/update ASN tools.
type ASNOutput struct {
	Data domain.ASN `json:"data"`
}

// ASNRangeCreateInput represents the writable fields for creating an ASN range.
type ASNRangeCreateInput struct {
	Name         string         `json:"name" jsonschema:"ASN range name (required)"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Start        *int           `json:"start,omitempty" jsonschema:"first ASN in the range"`
	End          *int           `json:"end,omitempty" jsonschema:"last ASN in the range"`
	RIR          *int           `json:"rir,omitempty" jsonschema:"RIR ID"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ASNRangeUpdateInput represents the writable fields for updating an ASN range.
// All fields are optional; only the explicitly provided ones are patched (PATCH
// merge).
type ASNRangeUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the ASN range to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"ASN range name"`
	Slug         *string        `json:"slug,omitempty" jsonschema:"URL-friendly slug"`
	Start        *int           `json:"start,omitempty" jsonschema:"first ASN in the range"`
	End          *int           `json:"end,omitempty" jsonschema:"last ASN in the range"`
	RIR          *int           `json:"rir,omitempty" jsonschema:"RIR ID"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ASNRangeDeleteInput represents the input fields for deleting an ASN range.
type ASNRangeDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the ASN range to delete,required"`
}

// ASNRangeOutput represents the output for the create/update ASN range tools.
type ASNRangeOutput struct {
	Data domain.ASNRange `json:"data"`
}

// RouteTargetCreateInput represents the writable fields for creating a route
// target.
type RouteTargetCreateInput struct {
	Name         string         `json:"name" jsonschema:"route target name (required)"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// RouteTargetUpdateInput represents the writable fields for updating a route
// target. All fields are optional; only the explicitly provided ones are
// patched (PATCH merge).
type RouteTargetUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the route target to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"route target name"`
	Tenant       *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// RouteTargetDeleteInput represents the input fields for deleting a route
// target.
type RouteTargetDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the route target to delete,required"`
}

// RouteTargetOutput represents the output for the create/update route target
// tools.
type RouteTargetOutput struct {
	Data domain.RouteTarget `json:"data"`
}

// ServiceCreateInput represents the writable fields for creating a service.
type ServiceCreateInput struct {
	Name             string         `json:"name" jsonschema:"service name (required)"`
	ParentObjectType *string        `json:"parent_object_type,omitempty" jsonschema:"parent object type (dcim.device)"`
	ParentObjectID   *int           `json:"parent_object_id,omitempty" jsonschema:"assigned parent object ID"`
	Protocol         *string        `json:"protocol,omitempty" jsonschema:"protocol (choice string, e.g. tcp)"`
	Ports            []int          `json:"ports,omitempty" jsonschema:"list of port numbers"`
	Description      *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags             []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields     map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ServiceUpdateInput represents the writable fields for updating a service. All
// fields are optional; only the explicitly provided ones are patched (PATCH
// merge).
type ServiceUpdateInput struct {
	ID               int            `json:"id" jsonschema:"numeric ID of the service to update,required"`
	Name             *string        `json:"name,omitempty" jsonschema:"service name"`
	ParentObjectType *string        `json:"parent_object_type,omitempty" jsonschema:"assigned parent object type"`
	ParentObjectID   *int           `json:"parent_object_id,omitempty" jsonschema:"assigned parent object ID"`
	Protocol         *string        `json:"protocol,omitempty" jsonschema:"protocol (choice string, e.g. tcp)"`
	Ports            []int          `json:"ports,omitempty" jsonschema:"list of port numbers"`
	Description      *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags             []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields     map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ServiceDeleteInput represents the input fields for deleting a service.
type ServiceDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the service to delete,required"`
}

// ServiceOutput represents the output for the create/update service tools.
type ServiceOutput struct {
	Data domain.Service `json:"data"`
}

// ServiceTemplateCreateInput represents the writable fields for creating a
// service template.
type ServiceTemplateCreateInput struct {
	Name         string         `json:"name" jsonschema:"service template name (required)"`
	Protocol     *string        `json:"protocol,omitempty" jsonschema:"protocol (choice string, e.g. tcp)"`
	Ports        []int          `json:"ports,omitempty" jsonschema:"list of port numbers"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ServiceTemplateUpdateInput represents the writable fields for updating a
// service template. All fields are optional; only the explicitly provided ones
// are patched (PATCH merge).
type ServiceTemplateUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the service template to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"service template name"`
	Protocol     *string        `json:"protocol,omitempty" jsonschema:"protocol (choice string, e.g. tcp)"`
	Ports        []int          `json:"ports,omitempty" jsonschema:"list of port numbers"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// ServiceTemplateDeleteInput represents the input fields for deleting a service
// template.
type ServiceTemplateDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the service template to delete,required"`
}

// ServiceTemplateOutput represents the output for the create/update service
// template tools.
type ServiceTemplateOutput struct {
	Data domain.ServiceTemplate `json:"data"`
}

// FHRPGroupCreateInput represents the writable fields for creating an FHRP
// group.
type FHRPGroupCreateInput struct {
	Name         string         `json:"name" jsonschema:"FHRP group name (required)"`
	Protocol     *string        `json:"protocol,omitempty" jsonschema:"protocol (choice string, e.g. vrrp)"`
	GroupID      *int           `json:"group_id,omitempty" jsonschema:"FHRP group number"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// FHRPGroupUpdateInput represents the writable fields for updating an FHRP
// group. All fields are optional; only the explicitly provided ones are patched
// (PATCH merge).
type FHRPGroupUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the FHRP group to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"FHRP group name"`
	Protocol     *string        `json:"protocol,omitempty" jsonschema:"protocol (choice string, e.g. vrrp)"`
	GroupID      *int           `json:"group_id,omitempty" jsonschema:"FHRP group number"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// FHRPGroupDeleteInput represents the input fields for deleting an FHRP group.
type FHRPGroupDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the FHRP group to delete,required"`
}

// FHRPGroupOutput represents the output for the create/update FHRP group tools.
type FHRPGroupOutput struct {
	Data domain.FHRPGroup `json:"data"`
}

// FHRPGroupAssignmentCreateInput represents the writable fields for creating an
// FHRP group assignment.
type FHRPGroupAssignmentCreateInput struct {
	Group         int    `json:"group" jsonschema:"FHRP group ID (required)"`
	InterfaceType string `json:"interface_type" jsonschema:"assigned interface type (e.g. dcim.interface) (required)"`
	InterfaceID   int    `json:"interface_id" jsonschema:"assigned interface ID (required)"`
	Priority      *int   `json:"priority,omitempty" jsonschema:"priority value"`
}

// FHRPGroupAssignmentUpdateInput represents the writable fields for updating an
// FHRP group assignment. All fields are optional; only the explicitly provided
// ones are patched (PATCH merge).
type FHRPGroupAssignmentUpdateInput struct {
	ID            int     `json:"id" jsonschema:"numeric ID of the FHRP group assignment to update,required"`
	Group         *int    `json:"group,omitempty" jsonschema:"FHRP group ID"`
	InterfaceType *string `json:"interface_type,omitempty" jsonschema:"assigned interface type"`
	InterfaceID   *int    `json:"interface_id,omitempty" jsonschema:"assigned interface ID"`
	Priority      *int    `json:"priority,omitempty" jsonschema:"priority value"`
}

// FHRPGroupAssignmentDeleteInput represents the input fields for deleting an
// FHRP group assignment.
type FHRPGroupAssignmentDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the FHRP group assignment to delete,required"`
}

// FHRPGroupAssignmentOutput represents the output for the create/update FHRP
// group assignment tools.
type FHRPGroupAssignmentOutput struct {
	Data domain.FHRPGroupAssignment `json:"data"`
}

// VLANTranslationPolicyCreateInput represents the writable fields for creating
// a VLAN translation policy.
type VLANTranslationPolicyCreateInput struct {
	Name         string         `json:"name" jsonschema:"policy name (required)"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VLANTranslationPolicyUpdateInput represents the writable fields for updating
// a VLAN translation policy. All fields are optional; only the explicitly
// provided ones are patched (PATCH merge).
type VLANTranslationPolicyUpdateInput struct {
	ID           int            `json:"id" jsonschema:"numeric ID of the policy to update,required"`
	Name         *string        `json:"name,omitempty" jsonschema:"policy name"`
	Description  *string        `json:"description,omitempty" jsonschema:"short description"`
	Tags         []string       `json:"tags,omitempty" jsonschema:"list of tag names"`
	CustomFields map[string]any `json:"custom_fields,omitempty" jsonschema:"custom field values keyed by name"`
}

// VLANTranslationPolicyDeleteInput represents the input fields for deleting a
// VLAN translation policy.
type VLANTranslationPolicyDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the policy to delete,required"`
}

// VLANTranslationPolicyOutput represents the output for the create/update VLAN
// translation policy tools.
type VLANTranslationPolicyOutput struct {
	Data domain.VLANTranslationPolicy `json:"data"`
}

// VLANTranslationRuleCreateInput represents the writable fields for creating a
// VLAN translation rule.
type VLANTranslationRuleCreateInput struct {
	Policy      int     `json:"policy" jsonschema:"policy ID (required)"`
	LocalVID    *int    `json:"local_vid,omitempty" jsonschema:"local VLAN VID"`
	RemoteVID   *int    `json:"remote_vid,omitempty" jsonschema:"remote VLAN VID"`
	Description *string `json:"description,omitempty" jsonschema:"short description"`
}

// VLANTranslationRuleUpdateInput represents the writable fields for updating a
// VLAN translation rule. All fields are optional; only the explicitly provided
// ones are patched (PATCH merge).
type VLANTranslationRuleUpdateInput struct {
	ID          int     `json:"id" jsonschema:"numeric ID of the rule to update,required"`
	Policy      *int    `json:"policy,omitempty" jsonschema:"policy ID"`
	LocalVID    *int    `json:"local_vid,omitempty" jsonschema:"local VLAN VID"`
	RemoteVID   *int    `json:"remote_vid,omitempty" jsonschema:"remote VLAN VID"`
	Description *string `json:"description,omitempty" jsonschema:"short description"`
}

// VLANTranslationRuleDeleteInput represents the input fields for deleting a VLAN
// translation rule.
type VLANTranslationRuleDeleteInput struct {
	ID int `json:"id" jsonschema:"numeric ID of the rule to delete,required"`
}

// VLANTranslationRuleOutput represents the output for the create/update VLAN
// translation rule tools.
type VLANTranslationRuleOutput struct {
	Data domain.VLANTranslationRule `json:"data"`
}

// frontPortWriteFromCreate builds a domain.FrontPortWrite from a create input.
func frontPortWriteFromCreate(input FrontPortCreateInput) domain.FrontPortWrite {
	return domain.FrontPortWrite{
		Name:             input.Name,
		Device:           input.Device,
		Type:             input.Type,
		RearPort:         input.RearPort,
		RearPortPosition: input.RearPortPosition,
		Description:      input.Description,
		Tags:             input.Tags,
		CustomFields:     input.CustomFields,
	}
}

// frontPortWriteFromUpdate builds a domain.FrontPortWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func frontPortWriteFromUpdate(input FrontPortUpdateInput) domain.FrontPortWrite {
	write := domain.FrontPortWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Device = input.Device
	write.Type = input.Type
	write.RearPort = input.RearPort
	write.RearPortPosition = input.RearPortPosition
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// rearPortWriteFromCreate builds a domain.RearPortWrite from a create input.
func rearPortWriteFromCreate(input RearPortCreateInput) domain.RearPortWrite {
	return domain.RearPortWrite{
		Name:         input.Name,
		Device:       input.Device,
		Type:         input.Type,
		Positions:    input.Positions,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// rearPortWriteFromUpdate builds a domain.RearPortWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func rearPortWriteFromUpdate(input RearPortUpdateInput) domain.RearPortWrite {
	write := domain.RearPortWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Device = input.Device
	write.Type = input.Type
	write.Positions = input.Positions
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// consolePortWriteFromCreate builds a domain.ConsolePortWrite from a create
// input.
func consolePortWriteFromCreate(input ConsolePortCreateInput) domain.ConsolePortWrite {
	return domain.ConsolePortWrite{
		Name:         input.Name,
		Device:       input.Device,
		Type:         input.Type,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// consolePortWriteFromUpdate builds a domain.ConsolePortWrite from an update
// input, mapping only the explicitly-provided (non-nil) fields so the PATCH
// body is minimal.
func consolePortWriteFromUpdate(input ConsolePortUpdateInput) domain.ConsolePortWrite {
	write := domain.ConsolePortWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Device = input.Device
	write.Type = input.Type
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// consoleServerPortWriteFromCreate builds a domain.ConsoleServerPortWrite from
// a create input.
func consoleServerPortWriteFromCreate(input ConsoleServerPortCreateInput) domain.ConsoleServerPortWrite {
	return domain.ConsoleServerPortWrite{
		Name:         input.Name,
		Device:       input.Device,
		Type:         input.Type,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// consoleServerPortWriteFromUpdate builds a domain.ConsoleServerPortWrite from
// an update input, mapping only the explicitly-provided (non-nil) fields so the
// PATCH body is minimal.
func consoleServerPortWriteFromUpdate(input ConsoleServerPortUpdateInput) domain.ConsoleServerPortWrite {
	write := domain.ConsoleServerPortWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Device = input.Device
	write.Type = input.Type
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// powerPortWriteFromCreate builds a domain.PowerPortWrite from a create input.
func powerPortWriteFromCreate(input PowerPortCreateInput) domain.PowerPortWrite {
	return domain.PowerPortWrite{
		Name:          input.Name,
		Device:        input.Device,
		Type:          input.Type,
		MaximumDraw:   input.MaximumDraw,
		AllocatedDraw: input.AllocatedDraw,
		Description:   input.Description,
		Tags:          input.Tags,
		CustomFields:  input.CustomFields,
	}
}

// powerPortWriteFromUpdate builds a domain.PowerPortWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func powerPortWriteFromUpdate(input PowerPortUpdateInput) domain.PowerPortWrite {
	write := domain.PowerPortWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Device = input.Device
	write.Type = input.Type
	write.MaximumDraw = input.MaximumDraw
	write.AllocatedDraw = input.AllocatedDraw
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// powerOutletWriteFromCreate builds a domain.PowerOutletWrite from a create
// input.
func powerOutletWriteFromCreate(input PowerOutletCreateInput) domain.PowerOutletWrite {
	return domain.PowerOutletWrite{
		Name:         input.Name,
		Device:       input.Device,
		Type:         input.Type,
		PowerPort:    input.PowerPort,
		FeedLeg:      input.FeedLeg,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// powerOutletWriteFromUpdate builds a domain.PowerOutletWrite from an update
// input, mapping only the explicitly-provided (non-nil) fields so the PATCH
// body is minimal.
func powerOutletWriteFromUpdate(input PowerOutletUpdateInput) domain.PowerOutletWrite {
	write := domain.PowerOutletWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Device = input.Device
	write.Type = input.Type
	write.PowerPort = input.PowerPort
	write.FeedLeg = input.FeedLeg
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// powerFeedWriteFromCreate builds a domain.PowerFeedWrite from a create input.
func powerFeedWriteFromCreate(input PowerFeedCreateInput) domain.PowerFeedWrite {
	return domain.PowerFeedWrite{
		Name:           input.Name,
		PowerPanel:     input.PowerPanel,
		Status:         input.Status,
		Type:           input.Type,
		Supply:         input.Supply,
		Phase:          input.Phase,
		Voltage:        input.Voltage,
		Amperage:       input.Amperage,
		MaxUtilization: input.MaxUtilization,
		Description:    input.Description,
		Tags:           input.Tags,
		CustomFields:   input.CustomFields,
	}
}

// powerFeedWriteFromUpdate builds a domain.PowerFeedWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func powerFeedWriteFromUpdate(input PowerFeedUpdateInput) domain.PowerFeedWrite {
	write := domain.PowerFeedWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.PowerPanel = input.PowerPanel
	write.Status = input.Status
	write.Type = input.Type
	write.Supply = input.Supply
	write.Phase = input.Phase
	write.Voltage = input.Voltage
	write.Amperage = input.Amperage
	write.MaxUtilization = input.MaxUtilization
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// deviceBayWriteFromCreate builds a domain.DeviceBayWrite from a create input.
func deviceBayWriteFromCreate(input DeviceBayCreateInput) domain.DeviceBayWrite {
	return domain.DeviceBayWrite{
		Name:         input.Name,
		Device:       input.Device,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// deviceBayWriteFromUpdate builds a domain.DeviceBayWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func deviceBayWriteFromUpdate(input DeviceBayUpdateInput) domain.DeviceBayWrite {
	write := domain.DeviceBayWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Device = input.Device
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// rackReservationWriteFromCreate builds a domain.RackReservationWrite from a
// create input.
func rackReservationWriteFromCreate(input RackReservationCreateInput) domain.RackReservationWrite {
	return domain.RackReservationWrite{
		Description:  input.Description,
		Rack:         input.Rack,
		Units:        input.Units,
		User:         input.User,
		Tenant:       input.Tenant,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// rackReservationWriteFromUpdate builds a domain.RackReservationWrite from an
// update input, mapping only the explicitly-provided (non-nil) fields so the
// PATCH body is minimal.
func rackReservationWriteFromUpdate(input RackReservationUpdateInput) domain.RackReservationWrite {
	write := domain.RackReservationWrite{}
	if input.Description != nil {
		write.Description = *input.Description
	}
	write.Rack = input.Rack
	write.Units = input.Units
	write.User = input.User
	write.Tenant = input.Tenant
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// inventoryItemWriteFromCreate builds a domain.InventoryItemWrite from a create
// input.
func inventoryItemWriteFromCreate(input InventoryItemCreateInput) domain.InventoryItemWrite {
	return domain.InventoryItemWrite{
		Name:         input.Name,
		Device:       input.Device,
		Role:         input.Role,
		Manufacturer: input.Manufacturer,
		Serial:       input.Serial,
		PartID:       input.PartID,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// inventoryItemWriteFromUpdate builds a domain.InventoryItemWrite from an
// update input, mapping only the explicitly-provided (non-nil) fields so the
// PATCH body is minimal.
func inventoryItemWriteFromUpdate(input InventoryItemUpdateInput) domain.InventoryItemWrite {
	write := domain.InventoryItemWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Device = input.Device
	write.Role = input.Role
	write.Manufacturer = input.Manufacturer
	write.Serial = input.Serial
	write.PartID = input.PartID
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// macAddressWriteFromCreate builds a domain.MACAddressWrite from a create
// input.
func macAddressWriteFromCreate(input MACAddressCreateInput) domain.MACAddressWrite {
	return domain.MACAddressWrite{
		Address:            input.Address,
		AssignedObjectType: input.AssignedObjectType,
		AssignedObjectID:   input.AssignedObjectID,
		Description:        input.Description,
		Tags:               input.Tags,
		CustomFields:       input.CustomFields,
	}
}

// macAddressWriteFromUpdate builds a domain.MACAddressWrite from an update
// input, mapping only the explicitly-provided (non-nil) fields so the PATCH
// body is minimal.
func macAddressWriteFromUpdate(input MACAddressUpdateInput) domain.MACAddressWrite {
	write := domain.MACAddressWrite{}
	if input.Address != nil {
		write.Address = *input.Address
	}
	write.AssignedObjectType = input.AssignedObjectType
	write.AssignedObjectID = input.AssignedObjectID
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// rirWriteFromCreate builds a domain.RIRWrite from a create input.
func rirWriteFromCreate(input RIRCreateInput) domain.RIRWrite {
	return domain.RIRWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		IsPrivate:    input.IsPrivate,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// rirWriteFromUpdate builds a domain.RIRWrite from an update input, mapping only
// the explicitly-provided (non-nil) fields so the PATCH body is minimal.
func rirWriteFromUpdate(input RIRUpdateInput) domain.RIRWrite {
	write := domain.RIRWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.IsPrivate = input.IsPrivate
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// contactGroupWriteFromCreate builds a domain.ContactGroupWrite from a create
// input.
func contactGroupWriteFromCreate(input ContactGroupCreateInput) domain.ContactGroupWrite {
	return domain.ContactGroupWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Parent:       input.Parent,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// contactGroupWriteFromUpdate builds a domain.ContactGroupWrite from an update
// input, mapping only the explicitly-provided (non-nil) fields so the PATCH
// body is minimal.
func contactGroupWriteFromUpdate(input ContactGroupUpdateInput) domain.ContactGroupWrite {
	write := domain.ContactGroupWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Parent = input.Parent
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// contactRoleWriteFromCreate builds a domain.ContactRoleWrite from a create
// input.
func contactRoleWriteFromCreate(input ContactRoleCreateInput) domain.ContactRoleWrite {
	return domain.ContactRoleWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// contactRoleWriteFromUpdate builds a domain.ContactRoleWrite from an update
// input, mapping only the explicitly-provided (non-nil) fields so the PATCH
// body is minimal.
func contactRoleWriteFromUpdate(input ContactRoleUpdateInput) domain.ContactRoleWrite {
	write := domain.ContactRoleWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// virtualDiskWriteFromCreate builds a domain.VirtualDiskWrite from a create
// input.
func virtualDiskWriteFromCreate(input VirtualDiskCreateInput) domain.VirtualDiskWrite {
	return domain.VirtualDiskWrite{
		Name:           input.Name,
		VirtualMachine: input.VirtualMachine,
		Size:           input.Size,
		Description:    input.Description,
		Tags:           input.Tags,
		CustomFields:   input.CustomFields,
	}
}

// virtualDiskWriteFromUpdate builds a domain.VirtualDiskWrite from an update
// input, mapping only the explicitly-provided (non-nil) fields so the PATCH
// body is minimal.
func virtualDiskWriteFromUpdate(input VirtualDiskUpdateInput) domain.VirtualDiskWrite {
	write := domain.VirtualDiskWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.VirtualMachine = input.VirtualMachine
	write.Size = input.Size
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// virtualMachineTypeWriteFromCreate builds a domain.VirtualMachineTypeWrite from
// a create input.
func virtualMachineTypeWriteFromCreate(input VirtualMachineTypeCreateInput) domain.VirtualMachineTypeWrite {
	return domain.VirtualMachineTypeWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// virtualMachineTypeWriteFromUpdate builds a domain.VirtualMachineTypeWrite from
// an update input, mapping only the explicitly-provided (non-nil) fields so the
// PATCH body is minimal.
func virtualMachineTypeWriteFromUpdate(input VirtualMachineTypeUpdateInput) domain.VirtualMachineTypeWrite {
	write := domain.VirtualMachineTypeWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// aggregateWriteFromCreate builds a domain.AggregateWrite from a create input.
func aggregateWriteFromCreate(input AggregateCreateInput) domain.AggregateWrite {
	return domain.AggregateWrite{
		Prefix:       input.Prefix,
		RIR:          input.RIR,
		Tenant:       input.Tenant,
		DateAdded:    input.DateAdded,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// aggregateWriteFromUpdate builds a domain.AggregateWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func aggregateWriteFromUpdate(input AggregateUpdateInput) domain.AggregateWrite {
	write := domain.AggregateWrite{}
	if input.Prefix != nil {
		write.Prefix = *input.Prefix
	}
	write.RIR = input.RIR
	write.Tenant = input.Tenant
	write.DateAdded = input.DateAdded
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// ipRangeWriteFromCreate builds a domain.IPRangeWrite from a create input.
func ipRangeWriteFromCreate(input IPRangeCreateInput) domain.IPRangeWrite {
	return domain.IPRangeWrite{
		StartAddress: input.StartAddress,
		EndAddress:   input.EndAddress,
		VRF:          input.VRF,
		Tenant:       input.Tenant,
		Status:       input.Status,
		Role:         input.Role,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// ipRangeWriteFromUpdate builds a domain.IPRangeWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func ipRangeWriteFromUpdate(input IPRangeUpdateInput) domain.IPRangeWrite {
	write := domain.IPRangeWrite{}
	if input.StartAddress != nil {
		write.StartAddress = *input.StartAddress
	}
	if input.EndAddress != nil {
		write.EndAddress = *input.EndAddress
	}
	write.VRF = input.VRF
	write.Tenant = input.Tenant
	write.Status = input.Status
	write.Role = input.Role
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// asnWriteFromCreate builds a domain.ASNWrite from a create input.
func asnWriteFromCreate(input ASNCreateInput) domain.ASNWrite {
	return domain.ASNWrite{
		ASN:          input.ASN,
		RIR:          input.RIR,
		Tenant:       input.Tenant,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// asnWriteFromUpdate builds a domain.ASNWrite from an update input, mapping only
// the explicitly-provided (non-nil) fields so the PATCH body is minimal.
func asnWriteFromUpdate(input ASNUpdateInput) domain.ASNWrite {
	write := domain.ASNWrite{}
	if input.ASN != nil {
		write.ASN = *input.ASN
	}
	write.RIR = input.RIR
	write.Tenant = input.Tenant
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// asnRangeWriteFromCreate builds a domain.ASNRangeWrite from a create input.
func asnRangeWriteFromCreate(input ASNRangeCreateInput) domain.ASNRangeWrite {
	return domain.ASNRangeWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Start:        input.Start,
		End:          input.End,
		RIR:          input.RIR,
		Tenant:       input.Tenant,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// asnRangeWriteFromUpdate builds a domain.ASNRangeWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func asnRangeWriteFromUpdate(input ASNRangeUpdateInput) domain.ASNRangeWrite {
	write := domain.ASNRangeWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Start = input.Start
	write.End = input.End
	write.RIR = input.RIR
	write.Tenant = input.Tenant
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// routeTargetWriteFromCreate builds a domain.RouteTargetWrite from a create
// input.
func routeTargetWriteFromCreate(input RouteTargetCreateInput) domain.RouteTargetWrite {
	return domain.RouteTargetWrite{
		Name:         input.Name,
		Tenant:       input.Tenant,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// routeTargetWriteFromUpdate builds a domain.RouteTargetWrite from an update
// input, mapping only the explicitly-provided (non-nil) fields so the PATCH body
// is minimal.
func routeTargetWriteFromUpdate(input RouteTargetUpdateInput) domain.RouteTargetWrite {
	write := domain.RouteTargetWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Tenant = input.Tenant
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// serviceWriteFromCreate builds a domain.ServiceWrite from a create input.
func serviceWriteFromCreate(input ServiceCreateInput) domain.ServiceWrite {
	return domain.ServiceWrite{
		Name:             input.Name,
		ParentObjectType: input.ParentObjectType,
		ParentObjectID:   input.ParentObjectID,
		Protocol:         input.Protocol,
		Ports:            input.Ports,
		Description:      input.Description,
		Tags:             input.Tags,
		CustomFields:     input.CustomFields,
	}
}

// serviceWriteFromUpdate builds a domain.ServiceWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func serviceWriteFromUpdate(input ServiceUpdateInput) domain.ServiceWrite {
	write := domain.ServiceWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	if input.ParentObjectType != nil {
		write.ParentObjectType = input.ParentObjectType
	}
	if input.ParentObjectID != nil {
		write.ParentObjectID = input.ParentObjectID
	}
	write.Protocol = input.Protocol
	write.Ports = input.Ports
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// serviceTemplateWriteFromCreate builds a domain.ServiceTemplateWrite from a
// create input.
func serviceTemplateWriteFromCreate(input ServiceTemplateCreateInput) domain.ServiceTemplateWrite {
	return domain.ServiceTemplateWrite{
		Name:         input.Name,
		Protocol:     input.Protocol,
		Ports:        input.Ports,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// serviceTemplateWriteFromUpdate builds a domain.ServiceTemplateWrite from an
// update input, mapping only the explicitly-provided (non-nil) fields so the
// PATCH body is minimal.
func serviceTemplateWriteFromUpdate(input ServiceTemplateUpdateInput) domain.ServiceTemplateWrite {
	write := domain.ServiceTemplateWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Protocol = input.Protocol
	write.Ports = input.Ports
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// fhrpGroupWriteFromCreate builds a domain.FHRPGroupWrite from a create input.
func fhrpGroupWriteFromCreate(input FHRPGroupCreateInput) domain.FHRPGroupWrite {
	return domain.FHRPGroupWrite{
		Name:         input.Name,
		Protocol:     input.Protocol,
		GroupID:      input.GroupID,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// fhrpGroupWriteFromUpdate builds a domain.FHRPGroupWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func fhrpGroupWriteFromUpdate(input FHRPGroupUpdateInput) domain.FHRPGroupWrite {
	write := domain.FHRPGroupWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Protocol = input.Protocol
	write.GroupID = input.GroupID
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// fhrpGroupAssignmentWriteFromCreate builds a domain.FHRPGroupAssignmentWrite
// from a create input.
func fhrpGroupAssignmentWriteFromCreate(input FHRPGroupAssignmentCreateInput) domain.FHRPGroupAssignmentWrite {
	return domain.FHRPGroupAssignmentWrite{
		Group:         input.Group,
		InterfaceType: input.InterfaceType,
		InterfaceID:   input.InterfaceID,
		Priority:      input.Priority,
	}
}

// fhrpGroupAssignmentWriteFromUpdate builds a domain.FHRPGroupAssignmentWrite
// from an update input, mapping only the explicitly-provided (non-nil) fields so
// the PATCH body is minimal.
func fhrpGroupAssignmentWriteFromUpdate(input FHRPGroupAssignmentUpdateInput) domain.FHRPGroupAssignmentWrite {
	write := domain.FHRPGroupAssignmentWrite{}
	if input.Group != nil {
		write.Group = *input.Group
	}
	if input.InterfaceType != nil {
		write.InterfaceType = *input.InterfaceType
	}
	if input.InterfaceID != nil {
		write.InterfaceID = *input.InterfaceID
	}
	write.Priority = input.Priority

	return write
}

// providerAccountWriteFromCreate builds a domain.ProviderAccountWrite from a
// create input.
func providerAccountWriteFromCreate(input ProviderAccountCreateInput) domain.ProviderAccountWrite {
	return domain.ProviderAccountWrite{
		Name:         input.Name,
		Provider:     input.Provider,
		Account:      input.Account,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// providerAccountWriteFromUpdate builds a domain.ProviderAccountWrite from an
// update input, mapping only the explicitly-provided (non-nil) fields so the
// PATCH body is minimal.
func providerAccountWriteFromUpdate(input ProviderAccountUpdateInput) domain.ProviderAccountWrite {
	write := domain.ProviderAccountWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Provider = input.Provider
	write.Account = input.Account
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// providerNetworkWriteFromCreate builds a domain.ProviderNetworkWrite from a
// create input.
func providerNetworkWriteFromCreate(input ProviderNetworkCreateInput) domain.ProviderNetworkWrite {
	return domain.ProviderNetworkWrite{
		Name:         input.Name,
		Provider:     input.Provider,
		Description:  input.Description,
		Comments:     input.Comments,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// providerNetworkWriteFromUpdate builds a domain.ProviderNetworkWrite from an
// update input, mapping only the explicitly-provided (non-nil) fields so the
// PATCH body is minimal.
func providerNetworkWriteFromUpdate(input ProviderNetworkUpdateInput) domain.ProviderNetworkWrite {
	write := domain.ProviderNetworkWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Provider = input.Provider
	write.Description = input.Description
	write.Comments = input.Comments
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// circuitGroupWriteFromCreate builds a domain.CircuitGroupWrite from a create
// input.
func circuitGroupWriteFromCreate(input CircuitGroupCreateInput) domain.CircuitGroupWrite {
	return domain.CircuitGroupWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// circuitGroupWriteFromUpdate builds a domain.CircuitGroupWrite from an update
// input, mapping only the explicitly-provided (non-nil) fields so the PATCH
// body is minimal.
func circuitGroupWriteFromUpdate(input CircuitGroupUpdateInput) domain.CircuitGroupWrite {
	write := domain.CircuitGroupWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// circuitGroupAssignmentWriteFromCreate builds a domain.CircuitGroupAssignmentWrite
// from a create input.
func circuitGroupAssignmentWriteFromCreate(input CircuitGroupAssignmentCreateInput) domain.CircuitGroupAssignmentWrite {
	return domain.CircuitGroupAssignmentWrite{
		Group:      input.Group,
		MemberType: input.MemberType,
		MemberID:   input.MemberID,
		Priority:   input.Priority,
	}
}

// circuitGroupAssignmentWriteFromUpdate builds a domain.CircuitGroupAssignmentWrite
// from an update input, mapping only the explicitly-provided (non-nil) fields
// so the PATCH body is minimal.
func circuitGroupAssignmentWriteFromUpdate(input CircuitGroupAssignmentUpdateInput) domain.CircuitGroupAssignmentWrite {
	write := domain.CircuitGroupAssignmentWrite{}
	if input.Group != nil {
		write.Group = *input.Group
	}
	if input.MemberType != nil {
		write.MemberType = *input.MemberType
	}
	if input.MemberID != nil {
		write.MemberID = *input.MemberID
	}
	write.Priority = input.Priority

	return write
}

// virtualCircuitWriteFromCreate builds a domain.VirtualCircuitWrite from a
// create input.
func virtualCircuitWriteFromCreate(input VirtualCircuitCreateInput) domain.VirtualCircuitWrite {
	return domain.VirtualCircuitWrite{
		CID:             input.CID,
		Type:            input.Type,
		ProviderNetwork: input.ProviderNetwork,
		Status:          input.Status,
		Description:     input.Description,
		Tags:            input.Tags,
		CustomFields:    input.CustomFields,
	}
}

// virtualCircuitWriteFromUpdate builds a domain.VirtualCircuitWrite from an
// update input, mapping only the explicitly-provided (non-nil) fields so the
// PATCH body is minimal.
func virtualCircuitWriteFromUpdate(input VirtualCircuitUpdateInput) domain.VirtualCircuitWrite {
	write := domain.VirtualCircuitWrite{}
	if input.CID != nil {
		write.CID = *input.CID
	}
	write.Type = input.Type
	write.ProviderNetwork = input.ProviderNetwork
	write.Status = input.Status
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// virtualCircuitTypeWriteFromCreate builds a domain.VirtualCircuitTypeWrite
// from a create input.
func virtualCircuitTypeWriteFromCreate(input VirtualCircuitTypeCreateInput) domain.VirtualCircuitTypeWrite {
	return domain.VirtualCircuitTypeWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// virtualCircuitTypeWriteFromUpdate builds a domain.VirtualCircuitTypeWrite
// from an update input, mapping only the explicitly-provided (non-nil) fields
// so the PATCH body is minimal.
func virtualCircuitTypeWriteFromUpdate(input VirtualCircuitTypeUpdateInput) domain.VirtualCircuitTypeWrite {
	write := domain.VirtualCircuitTypeWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Slug = input.Slug
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// virtualCircuitTerminationWriteFromCreate builds a
// domain.VirtualCircuitTerminationWrite from a create input.
func virtualCircuitTerminationWriteFromCreate(
	input VirtualCircuitTerminationCreateInput,
) domain.VirtualCircuitTerminationWrite {
	return domain.VirtualCircuitTerminationWrite{
		VirtualCircuit: input.VirtualCircuit,
		Role:           input.Role,
		Interface:      input.Interface,
	}
}

// virtualCircuitTerminationWriteFromUpdate builds a
// domain.VirtualCircuitTerminationWrite from an update input, mapping only the
// explicitly-provided (non-nil) fields so the PATCH body is minimal.
func virtualCircuitTerminationWriteFromUpdate(
	input VirtualCircuitTerminationUpdateInput,
) domain.VirtualCircuitTerminationWrite {
	write := domain.VirtualCircuitTerminationWrite{}
	if input.VirtualCircuit != nil {
		write.VirtualCircuit = *input.VirtualCircuit
	}
	write.Role = input.Role
	if input.Interface != nil {
		write.Interface = *input.Interface
	}

	return write
}

// vlanTranslationPolicyWriteFromCreate builds a domain.VLANTranslationPolicyWrite
// from a create input.
func vlanTranslationPolicyWriteFromCreate(input VLANTranslationPolicyCreateInput) domain.VLANTranslationPolicyWrite {
	return domain.VLANTranslationPolicyWrite{
		Name:         input.Name,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// vlanTranslationPolicyWriteFromUpdate builds a domain.VLANTranslationPolicyWrite
// from an update input, mapping only the explicitly-provided (non-nil) fields so
// the PATCH body is minimal.
func vlanTranslationPolicyWriteFromUpdate(input VLANTranslationPolicyUpdateInput) domain.VLANTranslationPolicyWrite {
	write := domain.VLANTranslationPolicyWrite{}
	if input.Name != nil {
		write.Name = *input.Name
	}
	write.Description = input.Description
	write.Tags = input.Tags
	write.CustomFields = input.CustomFields

	return write
}

// vlanTranslationRuleWriteFromCreate builds a domain.VLANTranslationRuleWrite
// from a create input.
func vlanTranslationRuleWriteFromCreate(input VLANTranslationRuleCreateInput) domain.VLANTranslationRuleWrite {
	return domain.VLANTranslationRuleWrite{
		Policy:      input.Policy,
		LocalVID:    input.LocalVID,
		RemoteVID:   input.RemoteVID,
		Description: input.Description,
	}
}

// vlanTranslationRuleWriteFromUpdate builds a domain.VLANTranslationRuleWrite
// from an update input, mapping only the explicitly-provided (non-nil) fields so
// the PATCH body is minimal.
func vlanTranslationRuleWriteFromUpdate(input VLANTranslationRuleUpdateInput) domain.VLANTranslationRuleWrite {
	write := domain.VLANTranslationRuleWrite{}
	if input.Policy != nil {
		write.Policy = *input.Policy
	}
	write.LocalVID = input.LocalVID
	write.RemoteVID = input.RemoteVID
	write.Description = input.Description

	return write
}

// NewGetFrontPortsHandler creates a handler for the get_front_ports tool.
func NewGetFrontPortsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[FrontPortsInput, PaginatedOutput[domain.FrontPort]] {
	return newListHandler(listHandlerConfig[FrontPortsInput, domain.FrontPort]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.FrontPort], error) {
			return resolveService(ctx, svc).ListFrontPorts(ctx, params)
		},
		buildParams: func(input FrontPortsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "device", input.Device)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list front ports",
	})
}

// NewCreateFrontPortHandler creates a handler for the create_front_port tool.
func NewCreateFrontPortHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[FrontPortCreateInput, FrontPortOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input FrontPortCreateInput,
	) (*mcp.CallToolResult, FrontPortOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, FrontPortOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, FrontPortOutput{}, fmt.Errorf("name is required")
		}
		if input.Device == nil {
			return &mcp.CallToolResult{IsError: true}, FrontPortOutput{}, fmt.Errorf("device is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		frontPort, err := service.CreateFrontPort(ctx, frontPortWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, FrontPortOutput{}, e
		}

		return &mcp.CallToolResult{}, FrontPortOutput{Data: *frontPort}, nil
	}
}

// NewUpdateFrontPortHandler creates a handler for the update_front_port tool.
// It performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateFrontPortHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[FrontPortUpdateInput, FrontPortOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input FrontPortUpdateInput,
	) (*mcp.CallToolResult, FrontPortOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, FrontPortOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, FrontPortOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		frontPort, err := service.UpdateFrontPort(ctx, input.ID, frontPortWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, FrontPortOutput{}, e
		}

		return &mcp.CallToolResult{}, FrontPortOutput{Data: *frontPort}, nil
	}
}

// NewDeleteFrontPortHandler creates a handler for the delete_front_port tool.
func NewDeleteFrontPortHandler(svc *application.NetworkService) mcp.ToolHandlerFor[FrontPortDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input FrontPortDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteFrontPort(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetRearPortsHandler creates a handler for the get_rear_ports tool.
func NewGetRearPortsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[RearPortsInput, PaginatedOutput[domain.RearPort]] {
	return newListHandler(listHandlerConfig[RearPortsInput, domain.RearPort]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.RearPort], error) {
			return resolveService(ctx, svc).ListRearPorts(ctx, params)
		},
		buildParams: func(input RearPortsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "device", input.Device)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list rear ports",
	})
}

// NewCreateRearPortHandler creates a handler for the create_rear_port tool.
func NewCreateRearPortHandler(svc *application.NetworkService) mcp.ToolHandlerFor[RearPortCreateInput, RearPortOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RearPortCreateInput,
	) (*mcp.CallToolResult, RearPortOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, RearPortOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, RearPortOutput{}, fmt.Errorf("name is required")
		}
		if input.Device == nil {
			return &mcp.CallToolResult{IsError: true}, RearPortOutput{}, fmt.Errorf("device is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		rearPort, err := service.CreateRearPort(ctx, rearPortWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, RearPortOutput{}, e
		}

		return &mcp.CallToolResult{}, RearPortOutput{Data: *rearPort}, nil
	}
}

// NewUpdateRearPortHandler creates a handler for the update_rear_port tool. It
// performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateRearPortHandler(svc *application.NetworkService) mcp.ToolHandlerFor[RearPortUpdateInput, RearPortOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RearPortUpdateInput,
	) (*mcp.CallToolResult, RearPortOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, RearPortOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, RearPortOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		rearPort, err := service.UpdateRearPort(ctx, input.ID, rearPortWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, RearPortOutput{}, e
		}

		return &mcp.CallToolResult{}, RearPortOutput{Data: *rearPort}, nil
	}
}

// NewDeleteRearPortHandler creates a handler for the delete_rear_port tool.
func NewDeleteRearPortHandler(svc *application.NetworkService) mcp.ToolHandlerFor[RearPortDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RearPortDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteRearPort(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetConsolePortsHandler creates a handler for the get_console_ports tool.
func NewGetConsolePortsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ConsolePortsInput, PaginatedOutput[domain.ConsolePort]] {
	return newListHandler(listHandlerConfig[ConsolePortsInput, domain.ConsolePort]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.ConsolePort], error) {
			return resolveService(ctx, svc).ListConsolePorts(ctx, params)
		},
		buildParams: func(input ConsolePortsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "device", input.Device)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list console ports",
	})
}

// NewCreateConsolePortHandler creates a handler for the create_console_port
// tool.
func NewCreateConsolePortHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ConsolePortCreateInput, ConsolePortOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ConsolePortCreateInput,
	) (*mcp.CallToolResult, ConsolePortOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ConsolePortOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, ConsolePortOutput{}, fmt.Errorf("name is required")
		}
		if input.Device == nil {
			return &mcp.CallToolResult{IsError: true}, ConsolePortOutput{}, fmt.Errorf("device is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		consolePort, err := service.CreateConsolePort(ctx, consolePortWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ConsolePortOutput{}, e
		}

		return &mcp.CallToolResult{}, ConsolePortOutput{Data: *consolePort}, nil
	}
}

// NewUpdateConsolePortHandler creates a handler for the update_console_port
// tool. It performs a partial-merge PATCH using only the explicitly provided
// fields.
func NewUpdateConsolePortHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ConsolePortUpdateInput, ConsolePortOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ConsolePortUpdateInput,
	) (*mcp.CallToolResult, ConsolePortOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ConsolePortOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, ConsolePortOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		consolePort, err := service.UpdateConsolePort(ctx, input.ID, consolePortWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ConsolePortOutput{}, e
		}

		return &mcp.CallToolResult{}, ConsolePortOutput{Data: *consolePort}, nil
	}
}

// NewDeleteConsolePortHandler creates a handler for the delete_console_port
// tool.
func NewDeleteConsolePortHandler(svc *application.NetworkService) mcp.ToolHandlerFor[ConsolePortDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ConsolePortDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteConsolePort(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetConsoleServerPortsHandler creates a handler for the
// get_console_server_ports tool.
func NewGetConsoleServerPortsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ConsoleServerPortsInput, PaginatedOutput[domain.ConsoleServerPort]] {
	return newListHandler(listHandlerConfig[ConsoleServerPortsInput, domain.ConsoleServerPort]{
		svc: svc,
		listFunc: func(
			ctx context.Context, params map[string]string,
		) (*domain.PaginatedResponse[domain.ConsoleServerPort], error) {
			return resolveService(ctx, svc).ListConsoleServerPorts(ctx, params)
		},
		buildParams: func(input ConsoleServerPortsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "device", input.Device)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list console server ports",
	})
}

// NewCreateConsoleServerPortHandler creates a handler for the
// create_console_server_port tool.
func NewCreateConsoleServerPortHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ConsoleServerPortCreateInput, ConsoleServerPortOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ConsoleServerPortCreateInput,
	) (*mcp.CallToolResult, ConsoleServerPortOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ConsoleServerPortOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, ConsoleServerPortOutput{}, fmt.Errorf("name is required")
		}
		if input.Device == nil {
			return &mcp.CallToolResult{IsError: true}, ConsoleServerPortOutput{}, fmt.Errorf("device is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		csp, err := service.CreateConsoleServerPort(ctx, consoleServerPortWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ConsoleServerPortOutput{}, e
		}

		return &mcp.CallToolResult{}, ConsoleServerPortOutput{Data: *csp}, nil
	}
}

// NewUpdateConsoleServerPortHandler creates a handler for the
// update_console_server_port tool. It performs a partial-merge PATCH using only
// the explicitly provided fields.
func NewUpdateConsoleServerPortHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ConsoleServerPortUpdateInput, ConsoleServerPortOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ConsoleServerPortUpdateInput,
	) (*mcp.CallToolResult, ConsoleServerPortOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ConsoleServerPortOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, ConsoleServerPortOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		csp, err := service.UpdateConsoleServerPort(ctx, input.ID, consoleServerPortWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ConsoleServerPortOutput{}, e
		}

		return &mcp.CallToolResult{}, ConsoleServerPortOutput{Data: *csp}, nil
	}
}

// NewDeleteConsoleServerPortHandler creates a handler for the
// delete_console_server_port tool.
func NewDeleteConsoleServerPortHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ConsoleServerPortDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ConsoleServerPortDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteConsoleServerPort(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetPowerPortsHandler creates a handler for the get_power_ports tool.
func NewGetPowerPortsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[PowerPortsInput, PaginatedOutput[domain.PowerPort]] {
	return newListHandler(listHandlerConfig[PowerPortsInput, domain.PowerPort]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.PowerPort], error) {
			return resolveService(ctx, svc).ListPowerPorts(ctx, params)
		},
		buildParams: func(input PowerPortsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "device", input.Device)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list power ports",
	})
}

// NewCreatePowerPortHandler creates a handler for the create_power_port tool.
func NewCreatePowerPortHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[PowerPortCreateInput, PowerPortOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input PowerPortCreateInput,
	) (*mcp.CallToolResult, PowerPortOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, PowerPortOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, PowerPortOutput{}, fmt.Errorf("name is required")
		}
		if input.Device == nil {
			return &mcp.CallToolResult{IsError: true}, PowerPortOutput{}, fmt.Errorf("device is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		powerPort, err := service.CreatePowerPort(ctx, powerPortWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, PowerPortOutput{}, e
		}

		return &mcp.CallToolResult{}, PowerPortOutput{Data: *powerPort}, nil
	}
}

// NewUpdatePowerPortHandler creates a handler for the update_power_port tool.
// It performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdatePowerPortHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[PowerPortUpdateInput, PowerPortOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input PowerPortUpdateInput,
	) (*mcp.CallToolResult, PowerPortOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, PowerPortOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, PowerPortOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		powerPort, err := service.UpdatePowerPort(ctx, input.ID, powerPortWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, PowerPortOutput{}, e
		}

		return &mcp.CallToolResult{}, PowerPortOutput{Data: *powerPort}, nil
	}
}

// NewDeletePowerPortHandler creates a handler for the delete_power_port tool.
func NewDeletePowerPortHandler(svc *application.NetworkService) mcp.ToolHandlerFor[PowerPortDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input PowerPortDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeletePowerPort(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetPowerOutletsHandler creates a handler for the get_power_outlets tool.
func NewGetPowerOutletsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[PowerOutletsInput, PaginatedOutput[domain.PowerOutlet]] {
	return newListHandler(listHandlerConfig[PowerOutletsInput, domain.PowerOutlet]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.PowerOutlet], error) {
			return resolveService(ctx, svc).ListPowerOutlets(ctx, params)
		},
		buildParams: func(input PowerOutletsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "device", input.Device)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list power outlets",
	})
}

// NewCreatePowerOutletHandler creates a handler for the create_power_outlet
// tool.
func NewCreatePowerOutletHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[PowerOutletCreateInput, PowerOutletOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input PowerOutletCreateInput,
	) (*mcp.CallToolResult, PowerOutletOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, PowerOutletOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, PowerOutletOutput{}, fmt.Errorf("name is required")
		}
		if input.Device == nil {
			return &mcp.CallToolResult{IsError: true}, PowerOutletOutput{}, fmt.Errorf("device is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		powerOutlet, err := service.CreatePowerOutlet(ctx, powerOutletWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, PowerOutletOutput{}, e
		}

		return &mcp.CallToolResult{}, PowerOutletOutput{Data: *powerOutlet}, nil
	}
}

// NewUpdatePowerOutletHandler creates a handler for the update_power_outlet
// tool. It performs a partial-merge PATCH using only the explicitly provided
// fields.
func NewUpdatePowerOutletHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[PowerOutletUpdateInput, PowerOutletOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input PowerOutletUpdateInput,
	) (*mcp.CallToolResult, PowerOutletOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, PowerOutletOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, PowerOutletOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		powerOutlet, err := service.UpdatePowerOutlet(ctx, input.ID, powerOutletWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, PowerOutletOutput{}, e
		}

		return &mcp.CallToolResult{}, PowerOutletOutput{Data: *powerOutlet}, nil
	}
}

// NewDeletePowerOutletHandler creates a handler for the delete_power_outlet
// tool.
func NewDeletePowerOutletHandler(svc *application.NetworkService) mcp.ToolHandlerFor[PowerOutletDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input PowerOutletDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeletePowerOutlet(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetPowerFeedsHandler creates a handler for the get_power_feeds tool.
func NewGetPowerFeedsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[PowerFeedsInput, PaginatedOutput[domain.PowerFeed]] {
	return newListHandler(listHandlerConfig[PowerFeedsInput, domain.PowerFeed]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.PowerFeed], error) {
			return resolveService(ctx, svc).ListPowerFeeds(ctx, params)
		},
		buildParams: func(input PowerFeedsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "power_panel", input.PowerPanel)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list power feeds",
	})
}

// NewCreatePowerFeedHandler creates a handler for the create_power_feed tool.
func NewCreatePowerFeedHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[PowerFeedCreateInput, PowerFeedOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input PowerFeedCreateInput,
	) (*mcp.CallToolResult, PowerFeedOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, PowerFeedOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, PowerFeedOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		powerFeed, err := service.CreatePowerFeed(ctx, powerFeedWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, PowerFeedOutput{}, e
		}

		return &mcp.CallToolResult{}, PowerFeedOutput{Data: *powerFeed}, nil
	}
}

// NewUpdatePowerFeedHandler creates a handler for the update_power_feed tool.
// It performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdatePowerFeedHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[PowerFeedUpdateInput, PowerFeedOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input PowerFeedUpdateInput,
	) (*mcp.CallToolResult, PowerFeedOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, PowerFeedOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, PowerFeedOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		powerFeed, err := service.UpdatePowerFeed(ctx, input.ID, powerFeedWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, PowerFeedOutput{}, e
		}

		return &mcp.CallToolResult{}, PowerFeedOutput{Data: *powerFeed}, nil
	}
}

// NewDeletePowerFeedHandler creates a handler for the delete_power_feed tool.
func NewDeletePowerFeedHandler(svc *application.NetworkService) mcp.ToolHandlerFor[PowerFeedDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input PowerFeedDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeletePowerFeed(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetDeviceBaysHandler creates a handler for the get_device_bays tool.
func NewGetDeviceBaysHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[DeviceBaysInput, PaginatedOutput[domain.DeviceBay]] {
	return newListHandler(listHandlerConfig[DeviceBaysInput, domain.DeviceBay]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.DeviceBay], error) {
			return resolveService(ctx, svc).ListDeviceBays(ctx, params)
		},
		buildParams: func(input DeviceBaysInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "device", input.Device)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list device bays",
	})
}

// NewCreateDeviceBayHandler creates a handler for the create_device_bay tool.
func NewCreateDeviceBayHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[DeviceBayCreateInput, DeviceBayOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input DeviceBayCreateInput,
	) (*mcp.CallToolResult, DeviceBayOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, DeviceBayOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, DeviceBayOutput{}, fmt.Errorf("name is required")
		}
		if input.Device == nil {
			return &mcp.CallToolResult{IsError: true}, DeviceBayOutput{}, fmt.Errorf("device is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		deviceBay, err := service.CreateDeviceBay(ctx, deviceBayWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, DeviceBayOutput{}, e
		}

		return &mcp.CallToolResult{}, DeviceBayOutput{Data: *deviceBay}, nil
	}
}

// NewUpdateDeviceBayHandler creates a handler for the update_device_bay tool.
// It performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateDeviceBayHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[DeviceBayUpdateInput, DeviceBayOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input DeviceBayUpdateInput,
	) (*mcp.CallToolResult, DeviceBayOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, DeviceBayOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, DeviceBayOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		deviceBay, err := service.UpdateDeviceBay(ctx, input.ID, deviceBayWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, DeviceBayOutput{}, e
		}

		return &mcp.CallToolResult{}, DeviceBayOutput{Data: *deviceBay}, nil
	}
}

// NewDeleteDeviceBayHandler creates a handler for the delete_device_bay tool.
func NewDeleteDeviceBayHandler(svc *application.NetworkService) mcp.ToolHandlerFor[DeviceBayDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input DeviceBayDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteDeviceBay(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetRackReservationsHandler creates a handler for the
// get_rack_reservations tool.
func NewGetRackReservationsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[RackReservationsInput, PaginatedOutput[domain.RackReservation]] {
	return newListHandler(listHandlerConfig[RackReservationsInput, domain.RackReservation]{
		svc: svc,
		listFunc: func(
			ctx context.Context, params map[string]string,
		) (*domain.PaginatedResponse[domain.RackReservation], error) {
			return resolveService(ctx, svc).ListRackReservations(ctx, params)
		},
		buildParams: func(input RackReservationsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "rack", input.Rack)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list rack reservations",
	})
}

// NewCreateRackReservationHandler creates a handler for the
// create_rack_reservation tool.
func NewCreateRackReservationHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[RackReservationCreateInput, RackReservationOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RackReservationCreateInput,
	) (*mcp.CallToolResult, RackReservationOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, RackReservationOutput{}, errServiceNotAvailable
		}
		if input.Description == "" {
			return &mcp.CallToolResult{IsError: true}, RackReservationOutput{}, fmt.Errorf("description is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		rackReservation, err := service.CreateRackReservation(ctx, rackReservationWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, RackReservationOutput{}, e
		}

		return &mcp.CallToolResult{}, RackReservationOutput{Data: *rackReservation}, nil
	}
}

// NewUpdateRackReservationHandler creates a handler for the
// update_rack_reservation tool. It performs a partial-merge PATCH using only
// the explicitly provided fields.
func NewUpdateRackReservationHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[RackReservationUpdateInput, RackReservationOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RackReservationUpdateInput,
	) (*mcp.CallToolResult, RackReservationOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, RackReservationOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, RackReservationOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		rackReservation, err := service.UpdateRackReservation(ctx, input.ID, rackReservationWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, RackReservationOutput{}, e
		}

		return &mcp.CallToolResult{}, RackReservationOutput{Data: *rackReservation}, nil
	}
}

// NewDeleteRackReservationHandler creates a handler for the
// delete_rack_reservation tool.
func NewDeleteRackReservationHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[RackReservationDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RackReservationDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteRackReservation(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetInventoryItemsHandler creates a handler for the get_inventory_items
// tool.
func NewGetInventoryItemsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[InventoryItemsInput, PaginatedOutput[domain.InventoryItem]] {
	return newListHandler(listHandlerConfig[InventoryItemsInput, domain.InventoryItem]{
		svc: svc,
		listFunc: func(
			ctx context.Context, params map[string]string,
		) (*domain.PaginatedResponse[domain.InventoryItem], error) {
			return resolveService(ctx, svc).ListInventoryItems(ctx, params)
		},
		buildParams: func(input InventoryItemsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "device", input.Device)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list inventory items",
	})
}

// NewCreateInventoryItemHandler creates a handler for the create_inventory_item
// tool.
func NewCreateInventoryItemHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[InventoryItemCreateInput, InventoryItemOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input InventoryItemCreateInput,
	) (*mcp.CallToolResult, InventoryItemOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, InventoryItemOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, InventoryItemOutput{}, fmt.Errorf("name is required")
		}
		if input.Device == nil {
			return &mcp.CallToolResult{IsError: true}, InventoryItemOutput{}, fmt.Errorf("device is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		inventoryItem, err := service.CreateInventoryItem(ctx, inventoryItemWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, InventoryItemOutput{}, e
		}

		return &mcp.CallToolResult{}, InventoryItemOutput{Data: *inventoryItem}, nil
	}
}

// NewUpdateInventoryItemHandler creates a handler for the update_inventory_item
// tool. It performs a partial-merge PATCH using only the explicitly provided
// fields.
func NewUpdateInventoryItemHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[InventoryItemUpdateInput, InventoryItemOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input InventoryItemUpdateInput,
	) (*mcp.CallToolResult, InventoryItemOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, InventoryItemOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, InventoryItemOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		inventoryItem, err := service.UpdateInventoryItem(ctx, input.ID, inventoryItemWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, InventoryItemOutput{}, e
		}

		return &mcp.CallToolResult{}, InventoryItemOutput{Data: *inventoryItem}, nil
	}
}

// NewDeleteInventoryItemHandler creates a handler for the delete_inventory_item
// tool.
func NewDeleteInventoryItemHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[InventoryItemDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input InventoryItemDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteInventoryItem(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetMACAddressesHandler creates a handler for the get_mac_addresses tool.
func NewGetMACAddressesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[MACAddressesInput, PaginatedOutput[domain.MACAddress]] {
	return newListHandler(listHandlerConfig[MACAddressesInput, domain.MACAddress]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.MACAddress], error) {
			return resolveService(ctx, svc).ListMACAddresses(ctx, params)
		},
		buildParams: func(input MACAddressesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list MAC addresses",
	})
}

// NewCreateMACAddressHandler creates a handler for the create_mac_address tool.
func NewCreateMACAddressHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[MACAddressCreateInput, MACAddressOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input MACAddressCreateInput,
	) (*mcp.CallToolResult, MACAddressOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, MACAddressOutput{}, errServiceNotAvailable
		}
		if input.Address == "" {
			return &mcp.CallToolResult{IsError: true}, MACAddressOutput{}, fmt.Errorf("address is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		macAddr, err := service.CreateMACAddress(ctx, macAddressWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, MACAddressOutput{}, e
		}

		return &mcp.CallToolResult{}, MACAddressOutput{Data: *macAddr}, nil
	}
}

// NewUpdateMACAddressHandler creates a handler for the update_mac_address tool.
// It performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateMACAddressHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[MACAddressUpdateInput, MACAddressOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input MACAddressUpdateInput,
	) (*mcp.CallToolResult, MACAddressOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, MACAddressOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, MACAddressOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		macAddr, err := service.UpdateMACAddress(ctx, input.ID, macAddressWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, MACAddressOutput{}, e
		}

		return &mcp.CallToolResult{}, MACAddressOutput{Data: *macAddr}, nil
	}
}

// NewDeleteMACAddressHandler creates a handler for the delete_mac_address tool.
func NewDeleteMACAddressHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[MACAddressDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input MACAddressDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteMACAddress(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetRIRsHandler creates a handler for the get_rirs tool.
func NewGetRIRsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[RIRsInput, PaginatedOutput[domain.RIR]] {
	return newListHandler(listHandlerConfig[RIRsInput, domain.RIR]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.RIR], error) {
			return resolveService(ctx, svc).ListRIRs(ctx, params)
		},
		buildParams: func(input RIRsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list RIRs",
	})
}

// NewCreateRIRHandler creates a handler for the create_rir tool.
func NewCreateRIRHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[RIRCreateInput, RIROutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RIRCreateInput,
	) (*mcp.CallToolResult, RIROutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, RIROutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, RIROutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateRIR(ctx, rirWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, RIROutput{}, e
		}

		return &mcp.CallToolResult{}, RIROutput{Data: *obj}, nil
	}
}

// NewUpdateRIRHandler creates a handler for the update_rir tool.
func NewUpdateRIRHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[RIRUpdateInput, RIROutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RIRUpdateInput,
	) (*mcp.CallToolResult, RIROutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, RIROutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, RIROutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateRIR(ctx, input.ID, rirWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, RIROutput{}, e
		}

		return &mcp.CallToolResult{}, RIROutput{Data: *obj}, nil
	}
}

// NewDeleteRIRHandler creates a handler for the delete_rir tool.
func NewDeleteRIRHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[RIRDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RIRDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteRIR(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetContactGroupsHandler creates a handler for the get_contact_groups tool.
func NewGetContactGroupsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ContactGroupsInput, PaginatedOutput[domain.ContactGroup]] {
	return newListHandler(listHandlerConfig[ContactGroupsInput, domain.ContactGroup]{
		svc: svc,
		listFunc: func(
			ctx context.Context, params map[string]string,
		) (*domain.PaginatedResponse[domain.ContactGroup], error) {
			return resolveService(ctx, svc).ListContactGroups(ctx, params)
		},
		buildParams: func(input ContactGroupsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "parent", input.Parent)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list contact groups",
	})
}

// NewCreateContactGroupHandler creates a handler for the create_contact_group
// tool.
func NewCreateContactGroupHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ContactGroupCreateInput, ContactGroupOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ContactGroupCreateInput,
	) (*mcp.CallToolResult, ContactGroupOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ContactGroupOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, ContactGroupOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateContactGroup(ctx, contactGroupWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ContactGroupOutput{}, e
		}

		return &mcp.CallToolResult{}, ContactGroupOutput{Data: *obj}, nil
	}
}

// NewUpdateContactGroupHandler creates a handler for the update_contact_group
// tool.
func NewUpdateContactGroupHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ContactGroupUpdateInput, ContactGroupOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ContactGroupUpdateInput,
	) (*mcp.CallToolResult, ContactGroupOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ContactGroupOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, ContactGroupOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateContactGroup(ctx, input.ID, contactGroupWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ContactGroupOutput{}, e
		}

		return &mcp.CallToolResult{}, ContactGroupOutput{Data: *obj}, nil
	}
}

// NewDeleteContactGroupHandler creates a handler for the delete_contact_group
// tool.
func NewDeleteContactGroupHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ContactGroupDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ContactGroupDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteContactGroup(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetContactRolesHandler creates a handler for the get_contact_roles tool.
func NewGetContactRolesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ContactRolesInput, PaginatedOutput[domain.ContactRole]] {
	return newListHandler(listHandlerConfig[ContactRolesInput, domain.ContactRole]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.ContactRole], error) {
			return resolveService(ctx, svc).ListContactRoles(ctx, params)
		},
		buildParams: func(input ContactRolesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list contact roles",
	})
}

// NewCreateContactRoleHandler creates a handler for the create_contact_role
// tool.
func NewCreateContactRoleHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ContactRoleCreateInput, ContactRoleOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ContactRoleCreateInput,
	) (*mcp.CallToolResult, ContactRoleOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ContactRoleOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, ContactRoleOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateContactRole(ctx, contactRoleWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ContactRoleOutput{}, e
		}

		return &mcp.CallToolResult{}, ContactRoleOutput{Data: *obj}, nil
	}
}

// NewUpdateContactRoleHandler creates a handler for the update_contact_role
// tool.
func NewUpdateContactRoleHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ContactRoleUpdateInput, ContactRoleOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ContactRoleUpdateInput,
	) (*mcp.CallToolResult, ContactRoleOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ContactRoleOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, ContactRoleOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateContactRole(ctx, input.ID, contactRoleWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ContactRoleOutput{}, e
		}

		return &mcp.CallToolResult{}, ContactRoleOutput{Data: *obj}, nil
	}
}

// NewDeleteContactRoleHandler creates a handler for the delete_contact_role
// tool.
func NewDeleteContactRoleHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ContactRoleDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ContactRoleDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteContactRole(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetVirtualDisksHandler creates a handler for the get_virtual_disks tool.
func NewGetVirtualDisksHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualDisksInput, PaginatedOutput[domain.VirtualDisk]] {
	return newListHandler(listHandlerConfig[VirtualDisksInput, domain.VirtualDisk]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.VirtualDisk], error) {
			return resolveService(ctx, svc).ListVirtualDisks(ctx, params)
		},
		buildParams: func(input VirtualDisksInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "virtual_machine", input.VirtualMachine)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list virtual disks",
	})
}

// NewCreateVirtualDiskHandler creates a handler for the create_virtual_disk
// tool.
func NewCreateVirtualDiskHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualDiskCreateInput, VirtualDiskOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VirtualDiskCreateInput,
	) (*mcp.CallToolResult, VirtualDiskOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VirtualDiskOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, VirtualDiskOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateVirtualDisk(ctx, virtualDiskWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VirtualDiskOutput{}, e
		}

		return &mcp.CallToolResult{}, VirtualDiskOutput{Data: *obj}, nil
	}
}

// NewUpdateVirtualDiskHandler creates a handler for the update_virtual_disk
// tool.
func NewUpdateVirtualDiskHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualDiskUpdateInput, VirtualDiskOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VirtualDiskUpdateInput,
	) (*mcp.CallToolResult, VirtualDiskOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VirtualDiskOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, VirtualDiskOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateVirtualDisk(ctx, input.ID, virtualDiskWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VirtualDiskOutput{}, e
		}

		return &mcp.CallToolResult{}, VirtualDiskOutput{Data: *obj}, nil
	}
}

// NewDeleteVirtualDiskHandler creates a handler for the delete_virtual_disk
// tool.
func NewDeleteVirtualDiskHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualDiskDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VirtualDiskDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteVirtualDisk(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetVirtualMachineTypesHandler creates a handler for the
// get_virtual_machine_types tool.
func NewGetVirtualMachineTypesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualMachineTypesInput, PaginatedOutput[domain.VirtualMachineType]] {
	return newListHandler(listHandlerConfig[VirtualMachineTypesInput, domain.VirtualMachineType]{
		svc: svc,
		listFunc: func(
			ctx context.Context, params map[string]string,
		) (*domain.PaginatedResponse[domain.VirtualMachineType], error) {
			return resolveService(ctx, svc).ListVirtualMachineTypes(ctx, params)
		},
		buildParams: func(input VirtualMachineTypesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list virtual machine types",
	})
}

// NewCreateVirtualMachineTypeHandler creates a handler for the
// create_virtual_machine_type tool.
func NewCreateVirtualMachineTypeHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualMachineTypeCreateInput, VirtualMachineTypeOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VirtualMachineTypeCreateInput,
	) (*mcp.CallToolResult, VirtualMachineTypeOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VirtualMachineTypeOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, VirtualMachineTypeOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateVirtualMachineType(ctx, virtualMachineTypeWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VirtualMachineTypeOutput{}, e
		}

		return &mcp.CallToolResult{}, VirtualMachineTypeOutput{Data: *obj}, nil
	}
}

// NewUpdateVirtualMachineTypeHandler creates a handler for the
// update_virtual_machine_type tool.
func NewUpdateVirtualMachineTypeHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualMachineTypeUpdateInput, VirtualMachineTypeOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VirtualMachineTypeUpdateInput,
	) (*mcp.CallToolResult, VirtualMachineTypeOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VirtualMachineTypeOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, VirtualMachineTypeOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateVirtualMachineType(ctx, input.ID, virtualMachineTypeWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VirtualMachineTypeOutput{}, e
		}

		return &mcp.CallToolResult{}, VirtualMachineTypeOutput{Data: *obj}, nil
	}
}

// NewDeleteVirtualMachineTypeHandler creates a handler for the
// delete_virtual_machine_type tool.
func NewDeleteVirtualMachineTypeHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualMachineTypeDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VirtualMachineTypeDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteVirtualMachineType(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetAggregatesHandler creates a handler for the get_aggregates tool.
func NewGetAggregatesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[AggregatesInput, PaginatedOutput[domain.Aggregate]] {
	return newListHandler(listHandlerConfig[AggregatesInput, domain.Aggregate]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Aggregate], error) {
			return resolveService(ctx, svc).ListAggregates(ctx, params)
		},
		buildParams: func(input AggregatesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "rir", input.RIR)
			addParam(params, "tenant", input.Tenant)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list aggregates",
	})
}

// NewCreateAggregateHandler creates a handler for the create_aggregate tool.
func NewCreateAggregateHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[AggregateCreateInput, AggregateOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input AggregateCreateInput,
	) (*mcp.CallToolResult, AggregateOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, AggregateOutput{}, errServiceNotAvailable
		}
		if input.Prefix == "" {
			return &mcp.CallToolResult{IsError: true}, AggregateOutput{}, fmt.Errorf("prefix is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateAggregate(ctx, aggregateWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, AggregateOutput{}, e
		}

		return &mcp.CallToolResult{}, AggregateOutput{Data: *obj}, nil
	}
}

// NewUpdateAggregateHandler creates a handler for the update_aggregate tool.
func NewUpdateAggregateHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[AggregateUpdateInput, AggregateOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input AggregateUpdateInput,
	) (*mcp.CallToolResult, AggregateOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, AggregateOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, AggregateOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateAggregate(ctx, input.ID, aggregateWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, AggregateOutput{}, e
		}

		return &mcp.CallToolResult{}, AggregateOutput{Data: *obj}, nil
	}
}

// NewDeleteAggregateHandler creates a handler for the delete_aggregate tool.
func NewDeleteAggregateHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[AggregateDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input AggregateDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteAggregate(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetIPRangesHandler creates a handler for the get_ip_ranges tool.
func NewGetIPRangesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[IPRangesInput, PaginatedOutput[domain.IPRange]] {
	return newListHandler(listHandlerConfig[IPRangesInput, domain.IPRange]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.IPRange], error) {
			return resolveService(ctx, svc).ListIPRanges(ctx, params)
		},
		buildParams: func(input IPRangesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "vrf", input.VRF)
			addParam(params, "tenant", input.Tenant)
			addParam(params, "status", input.Status)
			addParam(params, "role", input.Role)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list IP ranges",
	})
}

// NewCreateIPRangeHandler creates a handler for the create_ip_range tool.
func NewCreateIPRangeHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[IPRangeCreateInput, IPRangeOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input IPRangeCreateInput,
	) (*mcp.CallToolResult, IPRangeOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, IPRangeOutput{}, errServiceNotAvailable
		}
		if input.StartAddress == "" {
			return &mcp.CallToolResult{IsError: true}, IPRangeOutput{}, fmt.Errorf("start_address is required")
		}
		if input.EndAddress == "" {
			return &mcp.CallToolResult{IsError: true}, IPRangeOutput{}, fmt.Errorf("end_address is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateIPRange(ctx, ipRangeWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, IPRangeOutput{}, e
		}

		return &mcp.CallToolResult{}, IPRangeOutput{Data: *obj}, nil
	}
}

// NewUpdateIPRangeHandler creates a handler for the update_ip_range tool.
func NewUpdateIPRangeHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[IPRangeUpdateInput, IPRangeOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input IPRangeUpdateInput,
	) (*mcp.CallToolResult, IPRangeOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, IPRangeOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, IPRangeOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateIPRange(ctx, input.ID, ipRangeWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, IPRangeOutput{}, e
		}

		return &mcp.CallToolResult{}, IPRangeOutput{Data: *obj}, nil
	}
}

// NewDeleteIPRangeHandler creates a handler for the delete_ip_range tool.
func NewDeleteIPRangeHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[IPRangeDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input IPRangeDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteIPRange(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetASNsHandler creates a handler for the get_asns tool.
func NewGetASNsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ASNsInput, PaginatedOutput[domain.ASN]] {
	return newListHandler(listHandlerConfig[ASNsInput, domain.ASN]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.ASN], error) {
			return resolveService(ctx, svc).ListASNs(ctx, params)
		},
		buildParams: func(input ASNsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "rir", input.RIR)
			addParam(params, "tenant", input.Tenant)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list ASNs",
	})
}

// NewCreateASNHandler creates a handler for the create_asn tool.
func NewCreateASNHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ASNCreateInput, ASNOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ASNCreateInput,
	) (*mcp.CallToolResult, ASNOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ASNOutput{}, errServiceNotAvailable
		}
		if input.ASN == 0 {
			return &mcp.CallToolResult{IsError: true}, ASNOutput{}, fmt.Errorf("asn is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateASN(ctx, asnWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ASNOutput{}, e
		}

		return &mcp.CallToolResult{}, ASNOutput{Data: *obj}, nil
	}
}

// NewUpdateASNHandler creates a handler for the update_asn tool.
func NewUpdateASNHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ASNUpdateInput, ASNOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ASNUpdateInput,
	) (*mcp.CallToolResult, ASNOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ASNOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, ASNOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateASN(ctx, input.ID, asnWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ASNOutput{}, e
		}

		return &mcp.CallToolResult{}, ASNOutput{Data: *obj}, nil
	}
}

// NewDeleteASNHandler creates a handler for the delete_asn tool.
func NewDeleteASNHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ASNDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ASNDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteASN(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetASNRangesHandler creates a handler for the get_asn_ranges tool.
func NewGetASNRangesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ASNRangesInput, PaginatedOutput[domain.ASNRange]] {
	return newListHandler(listHandlerConfig[ASNRangesInput, domain.ASNRange]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.ASNRange], error) {
			return resolveService(ctx, svc).ListASNRanges(ctx, params)
		},
		buildParams: func(input ASNRangesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "rir", input.RIR)
			addParam(params, "tenant", input.Tenant)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list ASN ranges",
	})
}

// NewCreateASNRangeHandler creates a handler for the create_asn_range tool.
func NewCreateASNRangeHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ASNRangeCreateInput, ASNRangeOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ASNRangeCreateInput,
	) (*mcp.CallToolResult, ASNRangeOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ASNRangeOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, ASNRangeOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateASNRange(ctx, asnRangeWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ASNRangeOutput{}, e
		}

		return &mcp.CallToolResult{}, ASNRangeOutput{Data: *obj}, nil
	}
}

// NewUpdateASNRangeHandler creates a handler for the update_asn_range tool.
func NewUpdateASNRangeHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ASNRangeUpdateInput, ASNRangeOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ASNRangeUpdateInput,
	) (*mcp.CallToolResult, ASNRangeOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ASNRangeOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, ASNRangeOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateASNRange(ctx, input.ID, asnRangeWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ASNRangeOutput{}, e
		}

		return &mcp.CallToolResult{}, ASNRangeOutput{Data: *obj}, nil
	}
}

// NewDeleteASNRangeHandler creates a handler for the delete_asn_range tool.
func NewDeleteASNRangeHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ASNRangeDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ASNRangeDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteASNRange(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetRouteTargetsHandler creates a handler for the get_route_targets tool.
func NewGetRouteTargetsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[RouteTargetsInput, PaginatedOutput[domain.RouteTarget]] {
	return newListHandler(listHandlerConfig[RouteTargetsInput, domain.RouteTarget]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.RouteTarget], error) {
			return resolveService(ctx, svc).ListRouteTargets(ctx, params)
		},
		buildParams: func(input RouteTargetsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "tenant", input.Tenant)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list route targets",
	})
}

// NewCreateRouteTargetHandler creates a handler for the create_route_target tool.
func NewCreateRouteTargetHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[RouteTargetCreateInput, RouteTargetOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RouteTargetCreateInput,
	) (*mcp.CallToolResult, RouteTargetOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, RouteTargetOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, RouteTargetOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateRouteTarget(ctx, routeTargetWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, RouteTargetOutput{}, e
		}

		return &mcp.CallToolResult{}, RouteTargetOutput{Data: *obj}, nil
	}
}

// NewUpdateRouteTargetHandler creates a handler for the update_route_target tool.
func NewUpdateRouteTargetHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[RouteTargetUpdateInput, RouteTargetOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RouteTargetUpdateInput,
	) (*mcp.CallToolResult, RouteTargetOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, RouteTargetOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, RouteTargetOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateRouteTarget(ctx, input.ID, routeTargetWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, RouteTargetOutput{}, e
		}

		return &mcp.CallToolResult{}, RouteTargetOutput{Data: *obj}, nil
	}
}

// NewDeleteRouteTargetHandler creates a handler for the delete_route_target tool.
func NewDeleteRouteTargetHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[RouteTargetDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input RouteTargetDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteRouteTarget(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetServicesHandler creates a handler for the get_services tool.
func NewGetServicesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ServicesInput, PaginatedOutput[domain.Service]] {
	return newListHandler(listHandlerConfig[ServicesInput, domain.Service]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Service], error) {
			return resolveService(ctx, svc).ListServices(ctx, params)
		},
		buildParams: func(input ServicesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "device", input.Device)
			addParam(params, "virtual_machine", input.VirtualMachine)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list services",
	})
}

// NewCreateServiceHandler creates a handler for the create_service tool.
func NewCreateServiceHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ServiceCreateInput, ServiceOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ServiceCreateInput,
	) (*mcp.CallToolResult, ServiceOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ServiceOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, ServiceOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateService(ctx, serviceWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ServiceOutput{}, e
		}

		return &mcp.CallToolResult{}, ServiceOutput{Data: *obj}, nil
	}
}

// NewUpdateServiceHandler creates a handler for the update_service tool.
func NewUpdateServiceHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ServiceUpdateInput, ServiceOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ServiceUpdateInput,
	) (*mcp.CallToolResult, ServiceOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ServiceOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, ServiceOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateService(ctx, input.ID, serviceWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ServiceOutput{}, e
		}

		return &mcp.CallToolResult{}, ServiceOutput{Data: *obj}, nil
	}
}

// NewDeleteServiceHandler creates a handler for the delete_service tool.
func NewDeleteServiceHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ServiceDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ServiceDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteService(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetServiceTemplatesHandler creates a handler for the get_service_templates tool.
func NewGetServiceTemplatesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ServiceTemplatesInput, PaginatedOutput[domain.ServiceTemplate]] {
	return newListHandler(listHandlerConfig[ServiceTemplatesInput, domain.ServiceTemplate]{
		svc: svc,
		listFunc: func(
			ctx context.Context, params map[string]string,
		) (*domain.PaginatedResponse[domain.ServiceTemplate], error) {
			return resolveService(ctx, svc).ListServiceTemplates(ctx, params)
		},
		buildParams: func(input ServiceTemplatesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "protocol", input.Protocol)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list service templates",
	})
}

// NewCreateServiceTemplateHandler creates a handler for the
// create_service_template tool.
func NewCreateServiceTemplateHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ServiceTemplateCreateInput, ServiceTemplateOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ServiceTemplateCreateInput,
	) (*mcp.CallToolResult, ServiceTemplateOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ServiceTemplateOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, ServiceTemplateOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateServiceTemplate(ctx, serviceTemplateWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ServiceTemplateOutput{}, e
		}

		return &mcp.CallToolResult{}, ServiceTemplateOutput{Data: *obj}, nil
	}
}

// NewUpdateServiceTemplateHandler creates a handler for the
// update_service_template tool.
func NewUpdateServiceTemplateHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ServiceTemplateUpdateInput, ServiceTemplateOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ServiceTemplateUpdateInput,
	) (*mcp.CallToolResult, ServiceTemplateOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ServiceTemplateOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, ServiceTemplateOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateServiceTemplate(ctx, input.ID, serviceTemplateWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ServiceTemplateOutput{}, e
		}

		return &mcp.CallToolResult{}, ServiceTemplateOutput{Data: *obj}, nil
	}
}

// NewDeleteServiceTemplateHandler creates a handler for the
// delete_service_template tool.
func NewDeleteServiceTemplateHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ServiceTemplateDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ServiceTemplateDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteServiceTemplate(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetFHRPGroupsHandler creates a handler for the get_fhrp_groups tool.
func NewGetFHRPGroupsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[FHRPGroupsInput, PaginatedOutput[domain.FHRPGroup]] {
	return newListHandler(listHandlerConfig[FHRPGroupsInput, domain.FHRPGroup]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.FHRPGroup], error) {
			return resolveService(ctx, svc).ListFHRPGroups(ctx, params)
		},
		buildParams: func(input FHRPGroupsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "protocol", input.Protocol)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list FHRP groups",
	})
}

// NewCreateFHRPGroupHandler creates a handler for the create_fhrp_group tool.
func NewCreateFHRPGroupHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[FHRPGroupCreateInput, FHRPGroupOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input FHRPGroupCreateInput,
	) (*mcp.CallToolResult, FHRPGroupOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, FHRPGroupOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, FHRPGroupOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateFHRPGroup(ctx, fhrpGroupWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, FHRPGroupOutput{}, e
		}

		return &mcp.CallToolResult{}, FHRPGroupOutput{Data: *obj}, nil
	}
}

// NewUpdateFHRPGroupHandler creates a handler for the update_fhrp_group tool.
func NewUpdateFHRPGroupHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[FHRPGroupUpdateInput, FHRPGroupOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input FHRPGroupUpdateInput,
	) (*mcp.CallToolResult, FHRPGroupOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, FHRPGroupOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, FHRPGroupOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateFHRPGroup(ctx, input.ID, fhrpGroupWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, FHRPGroupOutput{}, e
		}

		return &mcp.CallToolResult{}, FHRPGroupOutput{Data: *obj}, nil
	}
}

// NewDeleteFHRPGroupHandler creates a handler for the delete_fhrp_group tool.
func NewDeleteFHRPGroupHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[FHRPGroupDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input FHRPGroupDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteFHRPGroup(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetFHRPGroupAssignmentsHandler creates a handler for the
// get_fhrp_group_assignments tool.
func NewGetFHRPGroupAssignmentsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[FHRPGroupAssignmentsInput, PaginatedOutput[domain.FHRPGroupAssignment]] {
	return newListHandler(listHandlerConfig[FHRPGroupAssignmentsInput, domain.FHRPGroupAssignment]{
		svc: svc,
		listFunc: func(
			ctx context.Context, params map[string]string,
		) (*domain.PaginatedResponse[domain.FHRPGroupAssignment], error) {
			return resolveService(ctx, svc).ListFHRPGroupAssignments(ctx, params)
		},
		buildParams: func(input FHRPGroupAssignmentsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "group", input.Group)
			addParam(params, "ip_address", input.IPAddress)

			return params
		},
		errorLabel: "list FHRP group assignments",
	})
}

// NewCreateFHRPGroupAssignmentHandler creates a handler for the
// create_fhrp_group_assignment tool.
func NewCreateFHRPGroupAssignmentHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[FHRPGroupAssignmentCreateInput, FHRPGroupAssignmentOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input FHRPGroupAssignmentCreateInput,
	) (*mcp.CallToolResult, FHRPGroupAssignmentOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, FHRPGroupAssignmentOutput{}, errServiceNotAvailable
		}
		if input.Group == 0 {
			return &mcp.CallToolResult{IsError: true}, FHRPGroupAssignmentOutput{}, fmt.Errorf("group is required")
		}
		if input.InterfaceID == 0 {
			return &mcp.CallToolResult{IsError: true}, FHRPGroupAssignmentOutput{}, fmt.Errorf("interface_id is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateFHRPGroupAssignment(ctx, fhrpGroupAssignmentWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, FHRPGroupAssignmentOutput{}, e
		}

		return &mcp.CallToolResult{}, FHRPGroupAssignmentOutput{Data: *obj}, nil
	}
}

// NewUpdateFHRPGroupAssignmentHandler creates a handler for the
// update_fhrp_group_assignment tool.
func NewUpdateFHRPGroupAssignmentHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[FHRPGroupAssignmentUpdateInput, FHRPGroupAssignmentOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input FHRPGroupAssignmentUpdateInput,
	) (*mcp.CallToolResult, FHRPGroupAssignmentOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, FHRPGroupAssignmentOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, FHRPGroupAssignmentOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateFHRPGroupAssignment(ctx, input.ID, fhrpGroupAssignmentWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, FHRPGroupAssignmentOutput{}, e
		}

		return &mcp.CallToolResult{}, FHRPGroupAssignmentOutput{Data: *obj}, nil
	}
}

// NewDeleteFHRPGroupAssignmentHandler creates a handler for the
// delete_fhrp_group_assignment tool.
func NewDeleteFHRPGroupAssignmentHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[FHRPGroupAssignmentDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input FHRPGroupAssignmentDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteFHRPGroupAssignment(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetVLANTranslationPoliciesHandler creates a handler for the
// get_vlan_translation_policies tool.
func NewGetVLANTranslationPoliciesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VLANTranslationPoliciesInput, PaginatedOutput[domain.VLANTranslationPolicy]] {
	return newListHandler(listHandlerConfig[VLANTranslationPoliciesInput, domain.VLANTranslationPolicy]{
		svc: svc,
		listFunc: func(
			ctx context.Context, params map[string]string,
		) (*domain.PaginatedResponse[domain.VLANTranslationPolicy], error) {
			return resolveService(ctx, svc).ListVLANTranslationPolicies(ctx, params)
		},
		buildParams: func(input VLANTranslationPoliciesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list VLAN translation policies",
	})
}

// NewCreateVLANTranslationPolicyHandler creates a handler for the
// create_vlan_translation_policy tool.
func NewCreateVLANTranslationPolicyHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VLANTranslationPolicyCreateInput, VLANTranslationPolicyOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VLANTranslationPolicyCreateInput,
	) (*mcp.CallToolResult, VLANTranslationPolicyOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VLANTranslationPolicyOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, VLANTranslationPolicyOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateVLANTranslationPolicy(ctx, vlanTranslationPolicyWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VLANTranslationPolicyOutput{}, e
		}

		return &mcp.CallToolResult{}, VLANTranslationPolicyOutput{Data: *obj}, nil
	}
}

// NewUpdateVLANTranslationPolicyHandler creates a handler for the
// update_vlan_translation_policy tool.
func NewUpdateVLANTranslationPolicyHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VLANTranslationPolicyUpdateInput, VLANTranslationPolicyOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VLANTranslationPolicyUpdateInput,
	) (*mcp.CallToolResult, VLANTranslationPolicyOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VLANTranslationPolicyOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, VLANTranslationPolicyOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateVLANTranslationPolicy(ctx, input.ID, vlanTranslationPolicyWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VLANTranslationPolicyOutput{}, e
		}

		return &mcp.CallToolResult{}, VLANTranslationPolicyOutput{Data: *obj}, nil
	}
}

// NewDeleteVLANTranslationPolicyHandler creates a handler for the
// delete_vlan_translation_policy tool.
func NewDeleteVLANTranslationPolicyHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VLANTranslationPolicyDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VLANTranslationPolicyDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteVLANTranslationPolicy(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetVLANTranslationRulesHandler creates a handler for the
// get_vlan_translation_rules tool.
func NewGetVLANTranslationRulesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VLANTranslationRulesInput, PaginatedOutput[domain.VLANTranslationRule]] {
	return newListHandler(listHandlerConfig[VLANTranslationRulesInput, domain.VLANTranslationRule]{
		svc: svc,
		listFunc: func(
			ctx context.Context, params map[string]string,
		) (*domain.PaginatedResponse[domain.VLANTranslationRule], error) {
			return resolveService(ctx, svc).ListVLANTranslationRules(ctx, params)
		},
		buildParams: func(input VLANTranslationRulesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "policy", input.Policy)

			return params
		},
		errorLabel: "list VLAN translation rules",
	})
}

// NewCreateVLANTranslationRuleHandler creates a handler for the
// create_vlan_translation_rule tool.
func NewCreateVLANTranslationRuleHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VLANTranslationRuleCreateInput, VLANTranslationRuleOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VLANTranslationRuleCreateInput,
	) (*mcp.CallToolResult, VLANTranslationRuleOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VLANTranslationRuleOutput{}, errServiceNotAvailable
		}
		if input.Policy == 0 {
			return &mcp.CallToolResult{IsError: true}, VLANTranslationRuleOutput{}, fmt.Errorf("policy is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateVLANTranslationRule(ctx, vlanTranslationRuleWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VLANTranslationRuleOutput{}, e
		}

		return &mcp.CallToolResult{}, VLANTranslationRuleOutput{Data: *obj}, nil
	}
}

// NewUpdateVLANTranslationRuleHandler creates a handler for the
// update_vlan_translation_rule tool.
func NewUpdateVLANTranslationRuleHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VLANTranslationRuleUpdateInput, VLANTranslationRuleOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VLANTranslationRuleUpdateInput,
	) (*mcp.CallToolResult, VLANTranslationRuleOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VLANTranslationRuleOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, VLANTranslationRuleOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateVLANTranslationRule(ctx, input.ID, vlanTranslationRuleWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VLANTranslationRuleOutput{}, e
		}

		return &mcp.CallToolResult{}, VLANTranslationRuleOutput{Data: *obj}, nil
	}
}

// NewDeleteVLANTranslationRuleHandler creates a handler for the
// delete_vlan_translation_rule tool.
func NewDeleteVLANTranslationRuleHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VLANTranslationRuleDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VLANTranslationRuleDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteVLANTranslationRule(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetProviderAccountsHandler creates a handler for the get_provider_accounts
// tool.
func NewGetProviderAccountsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ProviderAccountsInput, PaginatedOutput[domain.ProviderAccount]] {
	return newListHandler(listHandlerConfig[ProviderAccountsInput, domain.ProviderAccount]{
		svc: svc,
		listFunc: func(
			ctx context.Context, params map[string]string,
		) (*domain.PaginatedResponse[domain.ProviderAccount], error) {
			return resolveService(ctx, svc).ListProviderAccounts(ctx, params)
		},
		buildParams: func(input ProviderAccountsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "provider", input.Provider)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list provider accounts",
	})
}

// NewCreateProviderAccountHandler creates a handler for the
// create_provider_account tool.
func NewCreateProviderAccountHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ProviderAccountCreateInput, ProviderAccountOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ProviderAccountCreateInput,
	) (*mcp.CallToolResult, ProviderAccountOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ProviderAccountOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, ProviderAccountOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateProviderAccount(ctx, providerAccountWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ProviderAccountOutput{}, e
		}

		return &mcp.CallToolResult{}, ProviderAccountOutput{Data: *obj}, nil
	}
}

// NewUpdateProviderAccountHandler creates a handler for the
// update_provider_account tool.
func NewUpdateProviderAccountHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ProviderAccountUpdateInput, ProviderAccountOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ProviderAccountUpdateInput,
	) (*mcp.CallToolResult, ProviderAccountOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ProviderAccountOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, ProviderAccountOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateProviderAccount(ctx, input.ID, providerAccountWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ProviderAccountOutput{}, e
		}

		return &mcp.CallToolResult{}, ProviderAccountOutput{Data: *obj}, nil
	}
}

// NewDeleteProviderAccountHandler creates a handler for the
// delete_provider_account tool.
func NewDeleteProviderAccountHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ProviderAccountDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ProviderAccountDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteProviderAccount(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetProviderNetworksHandler creates a handler for the
// get_provider_networks tool.
func NewGetProviderNetworksHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ProviderNetworksInput, PaginatedOutput[domain.ProviderNetwork]] {
	return newListHandler(listHandlerConfig[ProviderNetworksInput, domain.ProviderNetwork]{
		svc: svc,
		listFunc: func(
			ctx context.Context, params map[string]string,
		) (*domain.PaginatedResponse[domain.ProviderNetwork], error) {
			return resolveService(ctx, svc).ListProviderNetworks(ctx, params)
		},
		buildParams: func(input ProviderNetworksInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "provider", input.Provider)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list provider networks",
	})
}

// NewCreateProviderNetworkHandler creates a handler for the
// create_provider_network tool.
func NewCreateProviderNetworkHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ProviderNetworkCreateInput, ProviderNetworkOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ProviderNetworkCreateInput,
	) (*mcp.CallToolResult, ProviderNetworkOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ProviderNetworkOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, ProviderNetworkOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateProviderNetwork(ctx, providerNetworkWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ProviderNetworkOutput{}, e
		}

		return &mcp.CallToolResult{}, ProviderNetworkOutput{Data: *obj}, nil
	}
}

// NewUpdateProviderNetworkHandler creates a handler for the
// update_provider_network tool.
func NewUpdateProviderNetworkHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ProviderNetworkUpdateInput, ProviderNetworkOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ProviderNetworkUpdateInput,
	) (*mcp.CallToolResult, ProviderNetworkOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, ProviderNetworkOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, ProviderNetworkOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateProviderNetwork(ctx, input.ID, providerNetworkWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, ProviderNetworkOutput{}, e
		}

		return &mcp.CallToolResult{}, ProviderNetworkOutput{Data: *obj}, nil
	}
}

// NewDeleteProviderNetworkHandler creates a handler for the
// delete_provider_network tool.
func NewDeleteProviderNetworkHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[ProviderNetworkDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input ProviderNetworkDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteProviderNetwork(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetCircuitGroupsHandler creates a handler for the get_circuit_groups tool.
func NewGetCircuitGroupsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[CircuitGroupsInput, PaginatedOutput[domain.CircuitGroup]] {
	return newListHandler(listHandlerConfig[CircuitGroupsInput, domain.CircuitGroup]{
		svc: svc,
		listFunc: func(
			ctx context.Context, params map[string]string,
		) (*domain.PaginatedResponse[domain.CircuitGroup], error) {
			return resolveService(ctx, svc).ListCircuitGroups(ctx, params)
		},
		buildParams: func(input CircuitGroupsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list circuit groups",
	})
}

// NewCreateCircuitGroupHandler creates a handler for the create_circuit_group
// tool.
func NewCreateCircuitGroupHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[CircuitGroupCreateInput, CircuitGroupOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input CircuitGroupCreateInput,
	) (*mcp.CallToolResult, CircuitGroupOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitGroupOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, CircuitGroupOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateCircuitGroup(ctx, circuitGroupWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, CircuitGroupOutput{}, e
		}

		return &mcp.CallToolResult{}, CircuitGroupOutput{Data: *obj}, nil
	}
}

// NewUpdateCircuitGroupHandler creates a handler for the update_circuit_group
// tool.
func NewUpdateCircuitGroupHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[CircuitGroupUpdateInput, CircuitGroupOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input CircuitGroupUpdateInput,
	) (*mcp.CallToolResult, CircuitGroupOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitGroupOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, CircuitGroupOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateCircuitGroup(ctx, input.ID, circuitGroupWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, CircuitGroupOutput{}, e
		}

		return &mcp.CallToolResult{}, CircuitGroupOutput{Data: *obj}, nil
	}
}

// NewDeleteCircuitGroupHandler creates a handler for the delete_circuit_group
// tool.
func NewDeleteCircuitGroupHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[CircuitGroupDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input CircuitGroupDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteCircuitGroup(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetCircuitGroupAssignmentsHandler creates a handler for the
// get_circuit_group_assignments tool.
func NewGetCircuitGroupAssignmentsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[CircuitGroupAssignmentsInput, PaginatedOutput[domain.CircuitGroupAssignment]] {
	return newListHandler(listHandlerConfig[CircuitGroupAssignmentsInput, domain.CircuitGroupAssignment]{
		svc: svc,
		listFunc: func(
			ctx context.Context, params map[string]string,
		) (*domain.PaginatedResponse[domain.CircuitGroupAssignment], error) {
			return resolveService(ctx, svc).ListCircuitGroupAssignments(ctx, params)
		},
		buildParams: func(input CircuitGroupAssignmentsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "group", input.Group)
			addParam(params, "circuit", input.Circuit)

			return params
		},
		errorLabel: "list circuit group assignments",
	})
}

// NewCreateCircuitGroupAssignmentHandler creates a handler for the
// create_circuit_group_assignment tool.
func NewCreateCircuitGroupAssignmentHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[CircuitGroupAssignmentCreateInput, CircuitGroupAssignmentOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input CircuitGroupAssignmentCreateInput,
	) (*mcp.CallToolResult, CircuitGroupAssignmentOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitGroupAssignmentOutput{}, errServiceNotAvailable
		}
		if input.Group == 0 {
			return &mcp.CallToolResult{IsError: true}, CircuitGroupAssignmentOutput{}, fmt.Errorf("group is required")
		}
		if input.MemberType == "" || input.MemberID == 0 {
			return &mcp.CallToolResult{IsError: true}, CircuitGroupAssignmentOutput{},
				fmt.Errorf("member_type and member_id are required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateCircuitGroupAssignment(ctx, circuitGroupAssignmentWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, CircuitGroupAssignmentOutput{}, e
		}

		return &mcp.CallToolResult{}, CircuitGroupAssignmentOutput{Data: *obj}, nil
	}
}

// NewUpdateCircuitGroupAssignmentHandler creates a handler for the
// update_circuit_group_assignment tool.
func NewUpdateCircuitGroupAssignmentHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[CircuitGroupAssignmentUpdateInput, CircuitGroupAssignmentOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input CircuitGroupAssignmentUpdateInput,
	) (*mcp.CallToolResult, CircuitGroupAssignmentOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitGroupAssignmentOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, CircuitGroupAssignmentOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateCircuitGroupAssignment(ctx, input.ID, circuitGroupAssignmentWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, CircuitGroupAssignmentOutput{}, e
		}

		return &mcp.CallToolResult{}, CircuitGroupAssignmentOutput{Data: *obj}, nil
	}
}

// NewDeleteCircuitGroupAssignmentHandler creates a handler for the
// delete_circuit_group_assignment tool.
func NewDeleteCircuitGroupAssignmentHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[CircuitGroupAssignmentDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input CircuitGroupAssignmentDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteCircuitGroupAssignment(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetVirtualCircuitsHandler creates a handler for the get_virtual_circuits
// tool.
func NewGetVirtualCircuitsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualCircuitsInput, PaginatedOutput[domain.VirtualCircuit]] {
	return newListHandler(listHandlerConfig[VirtualCircuitsInput, domain.VirtualCircuit]{
		svc: svc,
		listFunc: func(
			ctx context.Context, params map[string]string,
		) (*domain.PaginatedResponse[domain.VirtualCircuit], error) {
			return resolveService(ctx, svc).ListVirtualCircuits(ctx, params)
		},
		buildParams: func(input VirtualCircuitsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "type", input.Type)
			addParam(params, "provider_network", input.ProviderNetwork)
			addParam(params, "status", input.Status)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list virtual circuits",
	})
}

// NewCreateVirtualCircuitHandler creates a handler for the
// create_virtual_circuit tool.
func NewCreateVirtualCircuitHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualCircuitCreateInput, VirtualCircuitOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VirtualCircuitCreateInput,
	) (*mcp.CallToolResult, VirtualCircuitOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VirtualCircuitOutput{}, errServiceNotAvailable
		}
		if input.CID == "" {
			return &mcp.CallToolResult{IsError: true}, VirtualCircuitOutput{}, fmt.Errorf("cid is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateVirtualCircuit(ctx, virtualCircuitWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VirtualCircuitOutput{}, e
		}

		return &mcp.CallToolResult{}, VirtualCircuitOutput{Data: *obj}, nil
	}
}

// NewUpdateVirtualCircuitHandler creates a handler for the
// update_virtual_circuit tool.
func NewUpdateVirtualCircuitHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualCircuitUpdateInput, VirtualCircuitOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VirtualCircuitUpdateInput,
	) (*mcp.CallToolResult, VirtualCircuitOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VirtualCircuitOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, VirtualCircuitOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateVirtualCircuit(ctx, input.ID, virtualCircuitWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VirtualCircuitOutput{}, e
		}

		return &mcp.CallToolResult{}, VirtualCircuitOutput{Data: *obj}, nil
	}
}

// NewDeleteVirtualCircuitHandler creates a handler for the
// delete_virtual_circuit tool.
func NewDeleteVirtualCircuitHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualCircuitDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VirtualCircuitDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteVirtualCircuit(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetVirtualCircuitTypesHandler creates a handler for the
// get_virtual_circuit_types tool.
func NewGetVirtualCircuitTypesHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualCircuitTypesInput, PaginatedOutput[domain.VirtualCircuitType]] {
	return newListHandler(listHandlerConfig[VirtualCircuitTypesInput, domain.VirtualCircuitType]{
		svc: svc,
		listFunc: func(
			ctx context.Context, params map[string]string,
		) (*domain.PaginatedResponse[domain.VirtualCircuitType], error) {
			return resolveService(ctx, svc).ListVirtualCircuitTypes(ctx, params)
		},
		buildParams: func(input VirtualCircuitTypesInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "tag", input.Tag)

			return params
		},
		errorLabel: "list virtual circuit types",
	})
}

// NewCreateVirtualCircuitTypeHandler creates a handler for the
// create_virtual_circuit_type tool.
func NewCreateVirtualCircuitTypeHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualCircuitTypeCreateInput, VirtualCircuitTypeOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VirtualCircuitTypeCreateInput,
	) (*mcp.CallToolResult, VirtualCircuitTypeOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VirtualCircuitTypeOutput{}, errServiceNotAvailable
		}
		if input.Name == "" {
			return &mcp.CallToolResult{IsError: true}, VirtualCircuitTypeOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateVirtualCircuitType(ctx, virtualCircuitTypeWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VirtualCircuitTypeOutput{}, e
		}

		return &mcp.CallToolResult{}, VirtualCircuitTypeOutput{Data: *obj}, nil
	}
}

// NewUpdateVirtualCircuitTypeHandler creates a handler for the
// update_virtual_circuit_type tool.
func NewUpdateVirtualCircuitTypeHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualCircuitTypeUpdateInput, VirtualCircuitTypeOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VirtualCircuitTypeUpdateInput,
	) (*mcp.CallToolResult, VirtualCircuitTypeOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VirtualCircuitTypeOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, VirtualCircuitTypeOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateVirtualCircuitType(ctx, input.ID, virtualCircuitTypeWriteFromUpdate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VirtualCircuitTypeOutput{}, e
		}

		return &mcp.CallToolResult{}, VirtualCircuitTypeOutput{Data: *obj}, nil
	}
}

// NewDeleteVirtualCircuitTypeHandler creates a handler for the
// delete_virtual_circuit_type tool.
func NewDeleteVirtualCircuitTypeHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualCircuitTypeDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VirtualCircuitTypeDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteVirtualCircuitType(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewGetVirtualCircuitTerminationsHandler creates a handler for the
// get_virtual_circuit_terminations tool.
func NewGetVirtualCircuitTerminationsHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualCircuitTerminationsInput, PaginatedOutput[domain.VirtualCircuitTermination]] {
	return newListHandler(listHandlerConfig[VirtualCircuitTerminationsInput, domain.VirtualCircuitTermination]{
		svc: svc,
		listFunc: func(
			ctx context.Context, params map[string]string,
		) (*domain.PaginatedResponse[domain.VirtualCircuitTermination], error) {
			return resolveService(ctx, svc).ListVirtualCircuitTerminations(ctx, params)
		},
		buildParams: func(input VirtualCircuitTerminationsInput) map[string]string {
			params := paginationParams(input.Page, input.PageSize)
			addParam(params, "q", input.Q)
			addParam(params, "virtual_circuit", input.VirtualCircuit)
			addParam(params, "interface_type", input.InterfaceType)
			addParam(params, "interface_id", input.InterfaceID)

			return params
		},
		errorLabel: "list virtual circuit terminations",
	})
}

// NewCreateVirtualCircuitTerminationHandler creates a handler for the
// create_virtual_circuit_termination tool.
func NewCreateVirtualCircuitTerminationHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualCircuitTerminationCreateInput, VirtualCircuitTerminationOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VirtualCircuitTerminationCreateInput,
	) (*mcp.CallToolResult, VirtualCircuitTerminationOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VirtualCircuitTerminationOutput{}, errServiceNotAvailable
		}
		if input.VirtualCircuit == 0 {
			return &mcp.CallToolResult{IsError: true}, VirtualCircuitTerminationOutput{},
				fmt.Errorf("virtual_circuit is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.CreateVirtualCircuitTermination(ctx, virtualCircuitTerminationWriteFromCreate(input))
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VirtualCircuitTerminationOutput{}, e
		}

		return &mcp.CallToolResult{}, VirtualCircuitTerminationOutput{Data: *obj}, nil
	}
}

// NewUpdateVirtualCircuitTerminationHandler creates a handler for the
// update_virtual_circuit_termination tool.
func NewUpdateVirtualCircuitTerminationHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualCircuitTerminationUpdateInput, VirtualCircuitTerminationOutput] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VirtualCircuitTerminationUpdateInput,
	) (*mcp.CallToolResult, VirtualCircuitTerminationOutput, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, VirtualCircuitTerminationOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, VirtualCircuitTerminationOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&input))
		obj, err := service.UpdateVirtualCircuitTermination(
			ctx, input.ID, virtualCircuitTerminationWriteFromUpdate(input),
		)
		if err != nil {
			res, e := writeErrorResult(err)

			return res, VirtualCircuitTerminationOutput{}, e
		}

		return &mcp.CallToolResult{}, VirtualCircuitTerminationOutput{Data: *obj}, nil
	}
}

// NewDeleteVirtualCircuitTerminationHandler creates a handler for the
// delete_virtual_circuit_termination tool.
func NewDeleteVirtualCircuitTerminationHandler(
	svc *application.NetworkService,
) mcp.ToolHandlerFor[VirtualCircuitTerminationDeleteInput, struct{}] {
	return func(
		ctx context.Context, _ *mcp.CallToolRequest, input VirtualCircuitTerminationDeleteInput,
	) (*mcp.CallToolResult, struct{}, error) {
		service := resolveService(ctx, svc)
		if service == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(input.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := service.DeleteVirtualCircuitTermination(ctx, input.ID); err != nil {
			res, e := writeErrorResult(err)

			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}
