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

// DevicesInput represents the input fields for the get_devices tool.
type DevicesInput struct {
	Q            string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Site         string `json:"site,omitempty" jsonschema:"filter by site (slug)"`
	Role         string `json:"role,omitempty" jsonschema:"filter by device role (slug)"`
	Manufacturer string `json:"manufacturer,omitempty" jsonschema:"filter by manufacturer (slug)"`
	DeviceType   string `json:"device_type,omitempty" jsonschema:"filter by device type slug (e.g. c-1250)"`
	Status       string `json:"status,omitempty" jsonschema:"status: active, offline, planned, staged, failed, inventory, decommissioning"`
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
	ObjectType string            `json:"object_type" jsonschema:"object type: site, device, prefix, ip_address, vlan, virtual_machine, cluster, circuit, provider, tenant, rack, manufacturer, device_type, location, cluster_type, cluster_group, circuit_type, vrf, vlan_group, role, contact, cable, interface, vm_interface, circuit_termination,required"`
	ID         int               `json:"id" jsonschema:"numeric ID of the object (positive integer),required"`
	Params     map[string]string `json:"params,omitempty" jsonschema:"additional query parameters to pass to NetBox (optional)"`
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
	Slug            string         `json:"slug,omitempty" jsonschema:"URL-friendly unique identifier (defaults to a slugified name)"`
	Status          *string        `json:"status,omitempty" jsonschema:"status: planned, staged, active, decommissioning, retired"`
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
	Status          *string        `json:"status,omitempty" jsonschema:"status: planned, staged, active, decommissioning, retired"`
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
	Status       *string        `json:"status,omitempty" jsonschema:"status: offline, active, planned, staged, failed, inventory, decommissioning"`
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
	Status       *string        `json:"status,omitempty" jsonschema:"status: offline, active, planned, staged, failed, inventory, decommissioning"`
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
	Status             *string        `json:"status,omitempty" jsonschema:"status: active, reserved, deprecated, dhcp, slaac"`
	Role               *string        `json:"role,omitempty" jsonschema:"role: loopback, secondary, anycast, vip, vrrp, hsrp, glbp, carp"`
	VRF                *int           `json:"vrf,omitempty" jsonschema:"VRF ID"`
	Tenant             *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	DNSName            *string        `json:"dns_name,omitempty" jsonschema:"hostname or FQDN"`
	Description        *string        `json:"description,omitempty" jsonschema:"short description"`
	AssignedObjectType *string        `json:"assigned_object_type,omitempty" jsonschema:"assigned object type, e.g. dcim.interface (use with assigned_object_id)"`
	AssignedObjectID   *int           `json:"assigned_object_id,omitempty" jsonschema:"ID of the assigned object (use with assigned_object_type)"`
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
	Status             *string        `json:"status,omitempty" jsonschema:"status: active, reserved, deprecated, dhcp, slaac"`
	Role               *string        `json:"role,omitempty" jsonschema:"role: loopback, secondary, anycast, vip, vrrp, hsrp, glbp, carp"`
	VRF                *int           `json:"vrf,omitempty" jsonschema:"VRF ID"`
	Tenant             *int           `json:"tenant,omitempty" jsonschema:"tenant ID"`
	DNSName            *string        `json:"dns_name,omitempty" jsonschema:"hostname or FQDN"`
	Description        *string        `json:"description,omitempty" jsonschema:"short description"`
	AssignedObjectType *string        `json:"assigned_object_type,omitempty" jsonschema:"assigned object type, e.g. dcim.interface"`
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
	Status       *string        `json:"status,omitempty" jsonschema:"status: offline, active, planned, staged, failed, decommissioning"`
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
	Status       *string        `json:"status,omitempty" jsonschema:"status: offline, active, planned, staged, failed, decommissioning"`
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
	Status       *string        `json:"status,omitempty" jsonschema:"status: planned, provisioning, active, offline, decommissioning"`
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
	Status       *string        `json:"status,omitempty" jsonschema:"status: planned, provisioning, active, offline, decommissioning"`
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
	Status       *string        `json:"status,omitempty" jsonschema:"status: reserved, available, planned, active, decommissioning"`
	Role         *int           `json:"role,omitempty" jsonschema:"role ID (look up first)"`
	Serial       *string        `json:"serial,omitempty" jsonschema:"serial number"`
	AssetTag     *string        `json:"asset_tag,omitempty" jsonschema:"unique asset tag"`
	Type         *string        `json:"type,omitempty" jsonschema:"rack type: 2-post-frame, 4-post-frame, 4-post-cabinet, wall-frame, wall-cabinet"`
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
	Status       *string        `json:"status,omitempty" jsonschema:"status: reserved, available, planned, active, decommissioning"`
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
	Type         *string        `json:"type,omitempty" jsonschema:"interface type (required), e.g. 1000base-t, 10gbase-x-sfpp"`
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

// — output types —

// PaginatedOutput is a generic paginated response used by all list-oriented tools.
type PaginatedOutput[D any] struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []D    `json:"results"`
}

// GetObjectOutput represents the output for the get_object_by_id tool.
type GetObjectOutput struct {
	Data domain.RawObject `json:"data"`
}

// SiteOutput represents the output for the create/update site tools.
type SiteOutput struct {
	Data domain.Site `json:"data"`
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

// — write helpers —

// writeErrorResult converts an error into an MCP error result. A
// domain.ValidationError is surfaced as IsError with the NetBox body delivered
// via StructuredContent (not through the error string, so it never leaks into
// logs); all other errors are returned as a plain error to be surfaced normally.
func writeErrorResult(err error) (*mcp.CallToolResult, error) {
	var ve *domain.ValidationError
	if errors.As(err, &ve) {
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
func siteWriteFromCreate(in SiteCreateInput) domain.SiteWrite {
	return domain.SiteWrite{
		Name:            in.Name,
		Slug:            strPtr(in.Slug),
		Status:          in.Status,
		Region:          in.Region,
		Tenant:          in.Tenant,
		Facility:        in.Facility,
		TimeZone:        in.TimeZone,
		Description:     in.Description,
		PhysicalAddress: in.PhysicalAddress,
		ShippingAddress: in.ShippingAddress,
		Comments:        in.Comments,
		Tags:            in.Tags,
		CustomFields:    in.CustomFields,
	}
}

// siteWriteFromUpdate builds a domain.SiteWrite from an update input, mapping
// only the explicitly-provided (non-nil) fields so the PATCH body is minimal.
func siteWriteFromUpdate(in SiteUpdateInput) domain.SiteWrite {
	write := domain.SiteWrite{}
	if in.Name != nil {
		write.Name = *in.Name
	}
	write.Slug = in.Slug
	write.Status = in.Status
	write.Region = in.Region
	write.Tenant = in.Tenant
	write.Facility = in.Facility
	write.TimeZone = in.TimeZone
	write.Description = in.Description
	write.PhysicalAddress = in.PhysicalAddress
	write.ShippingAddress = in.ShippingAddress
	write.Comments = in.Comments
	write.Tags = in.Tags
	write.CustomFields = in.CustomFields
	return write
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// deviceWriteFromCreate builds a domain.DeviceWrite from a create input.
func deviceWriteFromCreate(in DeviceCreateInput) domain.DeviceWrite {
	return domain.DeviceWrite{
		Name:         in.Name,
		DeviceType:   in.DeviceType,
		Role:         in.Role,
		Tenant:       in.Tenant,
		Platform:     in.Platform,
		Serial:       in.Serial,
		AssetTag:     in.AssetTag,
		Site:         in.Site,
		Rack:         in.Rack,
		Position:     in.Position,
		Face:         in.Face,
		Status:       in.Status,
		Cluster:      in.Cluster,
		Comments:     in.Comments,
		Tags:         in.Tags,
		CustomFields: in.CustomFields,
	}
}

// deviceWriteFromUpdate builds a domain.DeviceWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func deviceWriteFromUpdate(in DeviceUpdateInput) domain.DeviceWrite {
	write := domain.DeviceWrite{}
	if in.Name != nil {
		write.Name = *in.Name
	}
	write.DeviceType = in.DeviceType
	write.Role = in.Role
	write.Tenant = in.Tenant
	write.Platform = in.Platform
	write.Serial = in.Serial
	write.AssetTag = in.AssetTag
	write.Site = in.Site
	write.Rack = in.Rack
	write.Position = in.Position
	write.Face = in.Face
	write.Status = in.Status
	write.Cluster = in.Cluster
	write.Comments = in.Comments
	write.Tags = in.Tags
	write.CustomFields = in.CustomFields
	return write
}

// ipAddressWriteFromCreate builds a domain.IPAddressWrite from a create input.
func ipAddressWriteFromCreate(in IPAddressCreateInput) domain.IPAddressWrite {
	return domain.IPAddressWrite{
		Address:            in.Address,
		Status:             in.Status,
		Role:               in.Role,
		VRF:                in.VRF,
		Tenant:             in.Tenant,
		DNSName:            in.DNSName,
		Description:        in.Description,
		AssignedObjectType: in.AssignedObjectType,
		AssignedObjectID:   in.AssignedObjectID,
		Comments:           in.Comments,
		Tags:               in.Tags,
		CustomFields:       in.CustomFields,
	}
}

// ipAddressWriteFromUpdate builds a domain.IPAddressWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func ipAddressWriteFromUpdate(in IPAddressUpdateInput) domain.IPAddressWrite {
	write := domain.IPAddressWrite{}
	if in.Address != nil {
		write.Address = *in.Address
	}
	write.Status = in.Status
	write.Role = in.Role
	write.VRF = in.VRF
	write.Tenant = in.Tenant
	write.DNSName = in.DNSName
	write.Description = in.Description
	write.AssignedObjectType = in.AssignedObjectType
	write.AssignedObjectID = in.AssignedObjectID
	write.Comments = in.Comments
	write.Tags = in.Tags
	write.CustomFields = in.CustomFields
	return write
}

// prefixWriteFromCreate builds a domain.PrefixWrite from a create input.
func prefixWriteFromCreate(in PrefixCreateInput) domain.PrefixWrite {
	return domain.PrefixWrite{
		Prefix:       in.Prefix,
		Site:         in.Site,
		VRF:          in.VRF,
		Tenant:       in.Tenant,
		VLAN:         in.VLAN,
		Status:       in.Status,
		Role:         in.Role,
		IsPool:       in.IsPool,
		Description:  in.Description,
		Comments:     in.Comments,
		Tags:         in.Tags,
		CustomFields: in.CustomFields,
	}
}

// prefixWriteFromUpdate builds a domain.PrefixWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func prefixWriteFromUpdate(in PrefixUpdateInput) domain.PrefixWrite {
	write := domain.PrefixWrite{}
	if in.Prefix != nil {
		write.Prefix = *in.Prefix
	}
	write.Site = in.Site
	write.VRF = in.VRF
	write.Tenant = in.Tenant
	write.VLAN = in.VLAN
	write.Status = in.Status
	write.Role = in.Role
	write.IsPool = in.IsPool
	write.Description = in.Description
	write.Comments = in.Comments
	write.Tags = in.Tags
	write.CustomFields = in.CustomFields
	return write
}

// vlanWriteFromCreate builds a domain.VLANWrite from a create input.
func vlanWriteFromCreate(in VLANCreateInput) domain.VLANWrite {
	return domain.VLANWrite{
		VID:          in.VID,
		Name:         in.Name,
		Site:         in.Site,
		Group:        in.Group,
		Tenant:       in.Tenant,
		Status:       in.Status,
		Role:         in.Role,
		Description:  in.Description,
		Comments:     in.Comments,
		Tags:         in.Tags,
		CustomFields: in.CustomFields,
	}
}

// vlanWriteFromUpdate builds a domain.VLANWrite from an update input, mapping
// only the explicitly-provided (non-nil) fields so the PATCH body is minimal.
func vlanWriteFromUpdate(in VLANUpdateInput) domain.VLANWrite {
	write := domain.VLANWrite{}
	if in.VID != nil {
		write.VID = *in.VID
	}
	if in.Name != nil {
		write.Name = *in.Name
	}
	write.Site = in.Site
	write.Group = in.Group
	write.Tenant = in.Tenant
	write.Status = in.Status
	write.Role = in.Role
	write.Description = in.Description
	write.Comments = in.Comments
	write.Tags = in.Tags
	write.CustomFields = in.CustomFields
	return write
}

// virtualMachineWriteFromCreate builds a domain.VirtualMachineWrite from a
// create input.
func virtualMachineWriteFromCreate(in VirtualMachineCreateInput) domain.VirtualMachineWrite {
	return domain.VirtualMachineWrite{
		Name:         in.Name,
		Cluster:      in.Cluster,
		Role:         in.Role,
		Tenant:       in.Tenant,
		Platform:     in.Platform,
		Status:       in.Status,
		Site:         in.Site,
		VCPUs:        in.VCPUs,
		Memory:       in.Memory,
		Disk:         in.Disk,
		Comments:     in.Comments,
		Tags:         in.Tags,
		CustomFields: in.CustomFields,
	}
}

// virtualMachineWriteFromUpdate builds a domain.VirtualMachineWrite from an
// update input, mapping only the explicitly-provided (non-nil) fields so the
// PATCH body is minimal.
func virtualMachineWriteFromUpdate(in VirtualMachineUpdateInput) domain.VirtualMachineWrite {
	write := domain.VirtualMachineWrite{}
	if in.Name != nil {
		write.Name = *in.Name
	}
	write.Cluster = in.Cluster
	write.Role = in.Role
	write.Tenant = in.Tenant
	write.Platform = in.Platform
	write.Status = in.Status
	write.Site = in.Site
	write.VCPUs = in.VCPUs
	write.Memory = in.Memory
	write.Disk = in.Disk
	write.Comments = in.Comments
	write.Tags = in.Tags
	write.CustomFields = in.CustomFields
	return write
}

// clusterWriteFromCreate builds a domain.ClusterWrite from a create input.
func clusterWriteFromCreate(in ClusterCreateInput) domain.ClusterWrite {
	return domain.ClusterWrite{
		Name:         in.Name,
		ClusterType:  in.ClusterType,
		ClusterGroup: in.ClusterGroup,
		Site:         in.Site,
		Tenant:       in.Tenant,
		Description:  in.Description,
		Comments:     in.Comments,
		Tags:         in.Tags,
		CustomFields: in.CustomFields,
	}
}

// clusterWriteFromUpdate builds a domain.ClusterWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func clusterWriteFromUpdate(in ClusterUpdateInput) domain.ClusterWrite {
	write := domain.ClusterWrite{}
	if in.Name != nil {
		write.Name = *in.Name
	}
	write.ClusterType = in.ClusterType
	write.ClusterGroup = in.ClusterGroup
	write.Site = in.Site
	write.Tenant = in.Tenant
	write.Description = in.Description
	write.Comments = in.Comments
	write.Tags = in.Tags
	write.CustomFields = in.CustomFields
	return write
}

// circuitWriteFromCreate builds a domain.CircuitWrite from a create input.
func circuitWriteFromCreate(in CircuitCreateInput) domain.CircuitWrite {
	return domain.CircuitWrite{
		CID:          in.CID,
		Provider:     in.Provider,
		CircuitType:  in.CircuitType,
		Tenant:       in.Tenant,
		Status:       in.Status,
		Description:  in.Description,
		Comments:     in.Comments,
		Tags:         in.Tags,
		CustomFields: in.CustomFields,
		InstallDate:  in.InstallDate,
		CommitRate:   in.CommitRate,
	}
}

// circuitWriteFromUpdate builds a domain.CircuitWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func circuitWriteFromUpdate(in CircuitUpdateInput) domain.CircuitWrite {
	write := domain.CircuitWrite{}
	if in.CID != nil {
		write.CID = *in.CID
	}
	write.Provider = in.Provider
	write.CircuitType = in.CircuitType
	write.Tenant = in.Tenant
	write.Status = in.Status
	write.Description = in.Description
	write.Comments = in.Comments
	write.Tags = in.Tags
	write.CustomFields = in.CustomFields
	write.InstallDate = in.InstallDate
	write.CommitRate = in.CommitRate
	return write
}

// rackWriteFromCreate builds a domain.RackWrite from a create input.
func rackWriteFromCreate(in RackCreateInput) domain.RackWrite {
	return domain.RackWrite{
		Name:         in.Name,
		FacilityID:   in.FacilityID,
		Site:         in.Site,
		Location:     in.Location,
		Tenant:       in.Tenant,
		Status:       in.Status,
		Role:         in.Role,
		Serial:       in.Serial,
		AssetTag:     in.AssetTag,
		Type:         in.Type,
		Width:        in.Width,
		UHeight:      in.UHeight,
		Comments:     in.Comments,
		Tags:         in.Tags,
		CustomFields: in.CustomFields,
	}
}

// rackWriteFromUpdate builds a domain.RackWrite from an update input, mapping
// only the explicitly-provided (non-nil) fields so the PATCH body is minimal.
func rackWriteFromUpdate(in RackUpdateInput) domain.RackWrite {
	write := domain.RackWrite{}
	if in.Name != nil {
		write.Name = *in.Name
	}
	write.FacilityID = in.FacilityID
	write.Site = in.Site
	write.Location = in.Location
	write.Tenant = in.Tenant
	write.Status = in.Status
	write.Role = in.Role
	write.Serial = in.Serial
	write.AssetTag = in.AssetTag
	write.Type = in.Type
	write.Width = in.Width
	write.UHeight = in.UHeight
	write.Comments = in.Comments
	write.Tags = in.Tags
	write.CustomFields = in.CustomFields
	return write
}

// interfaceWriteFromCreate builds a domain.InterfaceWrite from a create input.
func interfaceWriteFromCreate(in InterfaceCreateInput) domain.InterfaceWrite {
	return domain.InterfaceWrite{
		Device:       in.Device,
		Name:         in.Name,
		Type:         in.Type,
		Enabled:      in.Enabled,
		MTU:          in.MTU,
		MACAddress:   in.MACAddress,
		Speed:        in.Speed,
		Description:  in.Description,
		Tags:         in.Tags,
		CustomFields: in.CustomFields,
	}
}

// interfaceWriteFromUpdate builds a domain.InterfaceWrite from an update input,
// mapping only the explicitly-provided (non-nil) fields so the PATCH body is
// minimal.
func interfaceWriteFromUpdate(in InterfaceUpdateInput) domain.InterfaceWrite {
	write := domain.InterfaceWrite{}
	if in.Name != nil {
		write.Name = *in.Name
	}
	write.Device = in.Device
	write.Type = in.Type
	write.Enabled = in.Enabled
	write.MTU = in.MTU
	write.MACAddress = in.MACAddress
	write.Speed = in.Speed
	write.Description = in.Description
	write.Tags = in.Tags
	write.CustomFields = in.CustomFields
	return write
}

// circuitTerminationWriteFromCreate builds a domain.CircuitTerminationWrite
// from a create input.
func circuitTerminationWriteFromCreate(in CircuitTerminationCreateInput) domain.CircuitTerminationWrite {
	return domain.CircuitTerminationWrite{
		Circuit:       in.Circuit,
		TermSide:      in.TermSide,
		Site:          in.Site,
		Speed:         in.Speed,
		UpstreamSpeed: in.UpstreamSpeed,
		Description:   in.Description,
		Tags:          in.Tags,
		CustomFields:  in.CustomFields,
	}
}

// circuitTerminationWriteFromUpdate builds a domain.CircuitTerminationWrite
// from an update input, mapping only the explicitly-provided (non-nil) fields
// so the PATCH body is minimal.
func circuitTerminationWriteFromUpdate(in CircuitTerminationUpdateInput) domain.CircuitTerminationWrite {
	write := domain.CircuitTerminationWrite{}
	if in.TermSide != nil {
		write.TermSide = *in.TermSide
	}
	write.Circuit = in.Circuit
	write.Site = in.Site
	write.Speed = in.Speed
	write.UpstreamSpeed = in.UpstreamSpeed
	write.Description = in.Description
	write.Tags = in.Tags
	write.CustomFields = in.CustomFields
	return write
}

// cableTerminationFromInput converts a handler termination input into a domain
// cable termination DTO.
func cableTerminationFromInput(in *CableTerminationInput) *domain.CableTerminationWrite {
	if in == nil {
		return nil
	}
	return &domain.CableTerminationWrite{
		ObjectType: in.ObjectType,
		ObjectID:   in.ObjectID,
	}
}

// cableWriteFromCreate builds a domain.CableWrite from a create input.
func cableWriteFromCreate(in CableCreateInput) domain.CableWrite {
	return domain.CableWrite{
		TerminationA: cableTerminationFromInput(in.TerminationA),
		TerminationB: cableTerminationFromInput(in.TerminationB),
		Type:         in.Type,
		Status:       in.Status,
		Label:        in.Label,
		Color:        in.Color,
		Length:       in.Length,
		LengthUnit:   in.LengthUnit,
		Description:  in.Description,
		Tags:         in.Tags,
		CustomFields: in.CustomFields,
	}
}

// cableWriteFromUpdate builds a domain.CableWrite from an update input, mapping
// only the explicitly-provided (non-nil) fields so the PATCH body is minimal.
func cableWriteFromUpdate(in CableUpdateInput) domain.CableWrite {
	return domain.CableWrite{
		TerminationA: cableTerminationFromInput(in.TerminationA),
		TerminationB: cableTerminationFromInput(in.TerminationB),
		Type:         in.Type,
		Status:       in.Status,
		Label:        in.Label,
		Color:        in.Color,
		Length:       in.Length,
		LengthUnit:   in.LengthUnit,
		Description:  in.Description,
		Tags:         in.Tags,
		CustomFields: in.CustomFields,
	}
}

// — handler factories —

// maxPageSize is the maximum allowed page size for paginated requests.
// Increased from 100 to 1000 for better AI assistant UX (fewer pagination rounds).
const maxPageSize = 1000

func paginationParams(page, pageSize int) map[string]string {
	params := make(map[string]string)
	if page < 1 {
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
	return func(ctx context.Context, req *mcp.CallToolRequest, in I) (*mcp.CallToolResult, PaginatedOutput[D], error) {
		s := resolveService(ctx, cfg.svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, PaginatedOutput[D]{}, errServiceNotAvailable
		}

		params := cfg.buildParams(in)
		resp, err := cfg.listFunc(ctx, params)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, PaginatedOutput[D]{}, fmt.Errorf("%s: %w", cfg.errorLabel, err)
		}
		if resp == nil {
			return &mcp.CallToolResult{IsError: true}, PaginatedOutput[D]{}, fmt.Errorf("%s: %w", cfg.errorLabel, errors.New("received nil response from service"))
		}

		return &mcp.CallToolResult{}, PaginatedOutput[D]{Count: resp.Count, Next: resp.Next, Previous: resp.Previous, Results: resp.Results}, nil
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
		buildParams: func(in SitesInput) map[string]string {
			params := paginationParams(in.Page, in.PageSize)
			addParam(params, "q", in.Q)
			addParam(params, "region", in.Region)
			addParam(params, "status", in.Status)
			addParam(params, "tenant", in.Tenant)
			addParam(params, "tag", in.Tag)
			return params
		},
		errorLabel: "list sites",
	})
}

// NewGetDevicesHandler creates a handler for the get_devices tool.
func NewGetDevicesHandler(svc *application.NetworkService) mcp.ToolHandlerFor[DevicesInput, PaginatedOutput[domain.Device]] {
	return newListHandler(listHandlerConfig[DevicesInput, domain.Device]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Device], error) {
			return resolveService(ctx, svc).ListDevices(ctx, params)
		},
		buildParams: func(in DevicesInput) map[string]string {
			params := paginationParams(in.Page, in.PageSize)
			addParam(params, "q", in.Q)
			addParam(params, "site", in.Site)
			addParam(params, "role", in.Role)
			addParam(params, "manufacturer", in.Manufacturer)
			addParam(params, "device_type", in.DeviceType)
			addParam(params, "status", in.Status)
			addParam(params, "name__ic", in.Name)
			addParam(params, "tenant", in.Tenant)
			addParam(params, "rack", in.Rack)
			addParam(params, "cluster", in.Cluster)
			addParam(params, "tag", in.Tag)
			return params
		},
		errorLabel: "list devices",
	})
}

// NewGetIPAddressesHandler creates a handler for the get_ip_addresses tool.
func NewGetIPAddressesHandler(svc *application.NetworkService) mcp.ToolHandlerFor[IPAddressesInput, PaginatedOutput[domain.IPAddress]] {
	return newListHandler(listHandlerConfig[IPAddressesInput, domain.IPAddress]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.IPAddress], error) {
			return resolveService(ctx, svc).ListIPAddresses(ctx, params)
		},
		buildParams: func(in IPAddressesInput) map[string]string {
			params := paginationParams(in.Page, in.PageSize)
			addParam(params, "q", in.Q)
			addParam(params, "address", in.Address)
			addParam(params, "device", in.Device)
			addParam(params, "status", in.Status)
			addParam(params, "vrf", in.VRF)
			addParam(params, "role", in.Role)
			addParam(params, "tenant", in.Tenant)
			addParam(params, "tag", in.Tag)
			return params
		},
		errorLabel: "list IP addresses",
	})
}

// NewGetPrefixesHandler creates a handler for the get_prefixes tool.
func NewGetPrefixesHandler(svc *application.NetworkService) mcp.ToolHandlerFor[PrefixesInput, PaginatedOutput[domain.Prefix]] {
	return newListHandler(listHandlerConfig[PrefixesInput, domain.Prefix]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Prefix], error) {
			return resolveService(ctx, svc).ListPrefixes(ctx, params)
		},
		buildParams: func(in PrefixesInput) map[string]string {
			params := paginationParams(in.Page, in.PageSize)
			addParam(params, "q", in.Q)
			addParam(params, "prefix", in.Prefix)
			addParam(params, "site", in.Site)
			addParam(params, "vrf", in.VRF)
			addParam(params, "status", in.Status)
			addParam(params, "role", in.Role)
			addParam(params, "tenant", in.Tenant)
			addParam(params, "within", in.Within)
			if in.Family != nil {
				addParam(params, "family", strconv.Itoa(*in.Family))
			}
			addParam(params, "tag", in.Tag)
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
		buildParams: func(in VLANsInput) map[string]string {
			params := paginationParams(in.Page, in.PageSize)
			addParam(params, "q", in.Q)
			addParam(params, "site", in.Site)
			addParam(params, "group", in.Group)
			addParam(params, "status", in.Status)
			addParam(params, "tenant", in.Tenant)
			addIntParam(params, "vid", in.VID)
			addParam(params, "tag", in.Tag)
			return params
		},
		errorLabel: "list VLANs",
	})
}

// NewGetVirtualMachinesHandler creates a handler for the get_virtual_machines tool.
func NewGetVirtualMachinesHandler(svc *application.NetworkService) mcp.ToolHandlerFor[VirtualMachinesInput, PaginatedOutput[domain.VirtualMachine]] {
	return newListHandler(listHandlerConfig[VirtualMachinesInput, domain.VirtualMachine]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.VirtualMachine], error) {
			return resolveService(ctx, svc).ListVirtualMachines(ctx, params)
		},
		buildParams: func(in VirtualMachinesInput) map[string]string {
			params := paginationParams(in.Page, in.PageSize)
			addParam(params, "q", in.Q)
			addParam(params, "cluster", in.Cluster)
			addParam(params, "cluster_group", in.ClusterGroup)
			addParam(params, "role", in.Role)
			addParam(params, "status", in.Status)
			addParam(params, "tenant", in.Tenant)
			addParam(params, "name__ic", in.Name)
			addParam(params, "site", in.Site)
			addParam(params, "tag", in.Tag)
			return params
		},
		errorLabel: "list VMs",
	})
}

// NewGetClustersHandler creates a handler for the get_clusters tool.
func NewGetClustersHandler(svc *application.NetworkService) mcp.ToolHandlerFor[ClustersInput, PaginatedOutput[domain.Cluster]] {
	return newListHandler(listHandlerConfig[ClustersInput, domain.Cluster]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Cluster], error) {
			return resolveService(ctx, svc).ListClusters(ctx, params)
		},
		buildParams: func(in ClustersInput) map[string]string {
			params := paginationParams(in.Page, in.PageSize)
			addParam(params, "q", in.Q)
			addParam(params, "type", in.ClusterType)
			addParam(params, "group", in.ClusterGroup)
			addParam(params, "site", in.Site)
			addParam(params, "tenant", in.Tenant)
			addParam(params, "name__ic", in.Name)
			addParam(params, "tag", in.Tag)
			return params
		},
		errorLabel: "list clusters",
	})
}

// NewGetCircuitsHandler creates a handler for the get_circuits tool.
func NewGetCircuitsHandler(svc *application.NetworkService) mcp.ToolHandlerFor[CircuitsInput, PaginatedOutput[domain.Circuit]] {
	return newListHandler(listHandlerConfig[CircuitsInput, domain.Circuit]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Circuit], error) {
			return resolveService(ctx, svc).ListCircuits(ctx, params)
		},
		buildParams: func(in CircuitsInput) map[string]string {
			params := paginationParams(in.Page, in.PageSize)
			addParam(params, "q", in.Q)
			addParam(params, "provider", in.Provider)
			addParam(params, "type", in.CircuitType)
			addParam(params, "site", in.Site)
			addParam(params, "status", in.Status)
			addParam(params, "tenant", in.Tenant)
			addParam(params, "tag", in.Tag)
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
		buildParams: func(in RacksInput) map[string]string {
			params := paginationParams(in.Page, in.PageSize)
			addParam(params, "q", in.Q)
			addParam(params, "site", in.Site)
			addParam(params, "location", in.Location)
			addParam(params, "status", in.Status)
			addParam(params, "tenant", in.Tenant)
			addParam(params, "tag", in.Tag)
			return params
		},
		errorLabel: "list racks",
	})
}

// NewGetInterfacesHandler creates a handler for the get_interfaces tool.
func NewGetInterfacesHandler(svc *application.NetworkService) mcp.ToolHandlerFor[InterfacesInput, PaginatedOutput[domain.Interface]] {
	return newListHandler(listHandlerConfig[InterfacesInput, domain.Interface]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Interface], error) {
			return resolveService(ctx, svc).ListInterfaces(ctx, params)
		},
		buildParams: func(in InterfacesInput) map[string]string {
			params := paginationParams(in.Page, in.PageSize)
			addParam(params, "q", in.Q)
			addParam(params, "device", in.Device)
			addParam(params, "type", in.Type)
			addParam(params, "name__ic", in.Name)
			addParam(params, "tag", in.Tag)
			if in.Enabled != nil {
				addParam(params, "enabled", fmt.Sprintf("%t", *in.Enabled))
			}
			return params
		},
		errorLabel: "list interfaces",
	})
}

// NewGetVMInterfacesHandler creates a handler for the get_vm_interfaces tool.
func NewGetVMInterfacesHandler(svc *application.NetworkService) mcp.ToolHandlerFor[VMInterfacesInput, PaginatedOutput[domain.VMInterface]] {
	return newListHandler(listHandlerConfig[VMInterfacesInput, domain.VMInterface]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.VMInterface], error) {
			return resolveService(ctx, svc).ListVMInterfaces(ctx, params)
		},
		buildParams: func(in VMInterfacesInput) map[string]string {
			params := paginationParams(in.Page, in.PageSize)
			addParam(params, "q", in.Q)
			addParam(params, "virtual_machine", in.VirtualMachine)
			addParam(params, "name__ic", in.Name)
			addParam(params, "tag", in.Tag)
			return params
		},
		errorLabel: "list VM interfaces",
	})
}

// NewGetCircuitTerminationsHandler creates a handler for the get_circuit_terminations tool.
func NewGetCircuitTerminationsHandler(svc *application.NetworkService) mcp.ToolHandlerFor[CircuitTerminationsInput, PaginatedOutput[domain.CircuitTermination]] {
	return newListHandler(listHandlerConfig[CircuitTerminationsInput, domain.CircuitTermination]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.CircuitTermination], error) {
			return resolveService(ctx, svc).ListCircuitTerminations(ctx, params)
		},
		buildParams: func(in CircuitTerminationsInput) map[string]string {
			params := paginationParams(in.Page, in.PageSize)
			addParam(params, "q", in.Q)
			addParam(params, "circuit", in.Circuit)
			addParam(params, "site", in.Site)
			addParam(params, "tag", in.Tag)
			addParam(params, "term_side", in.TermSide)
			return params
		},
		errorLabel: "list circuit terminations",
	})
}

// NewGetCablesHandler creates a handler for the get_cables tool.
func NewGetCablesHandler(svc *application.NetworkService) mcp.ToolHandlerFor[CablesInput, PaginatedOutput[domain.Cable]] {
	return newListHandler(listHandlerConfig[CablesInput, domain.Cable]{
		svc: svc,
		listFunc: func(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Cable], error) {
			return resolveService(ctx, svc).ListCables(ctx, params)
		},
		buildParams: func(in CablesInput) map[string]string {
			params := paginationParams(in.Page, in.PageSize)
			addParam(params, "q", in.Q)
			addParam(params, "type", in.Type)
			addParam(params, "status", in.Status)
			addParam(params, "site", in.Site)
			addParam(params, "tag", in.Tag)
			addParam(params, "color", in.Color)
			addParam(params, "label__ic", in.Label)
			return params
		},
		errorLabel: "list cables",
	})
}

// NewGetObjectByIDHandler creates a handler for the get_object_by_id tool.
func NewGetObjectByIDHandler(svc *application.NetworkService) mcp.ToolHandlerFor[GetObjectInput, GetObjectOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in GetObjectInput) (*mcp.CallToolResult, GetObjectOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, GetObjectOutput{}, errServiceNotAvailable
		}

		if in.ObjectType == "" {
			return &mcp.CallToolResult{IsError: true}, GetObjectOutput{}, fmt.Errorf("object_type is required")
		}
		if in.ID <= 0 {
			return &mcp.CallToolResult{IsError: true}, GetObjectOutput{}, fmt.Errorf("id must be a positive integer")
		}

		// Sanitize params to prevent CRLF injection
		sanitizedParams := make(map[string]string, len(in.Params))
		for k, v := range in.Params {
			sanitizedKey := strings.Map(func(r rune) rune {
				if r == '\r' || r == '\n' {
					return -1
				}
				return r
			}, k)
			sanitizedVal := strings.Map(func(r rune) rune {
				if r == '\r' || r == '\n' {
					return -1
				}
				return r
			}, v)
			sanitizedParams[sanitizedKey] = sanitizedVal
		}

		data, err := s.GetObject(ctx, in.ObjectType, in.ID, sanitizedParams)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, GetObjectOutput{}, fmt.Errorf("get object: %w", err)
		}

		return &mcp.CallToolResult{}, GetObjectOutput{Data: data}, nil
	}
}

// NewCreateSiteHandler creates a handler for the create_site tool.
func NewCreateSiteHandler(svc *application.NetworkService) mcp.ToolHandlerFor[SiteCreateInput, SiteOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in SiteCreateInput) (*mcp.CallToolResult, SiteOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, SiteOutput{}, errServiceNotAvailable
		}
		if in.Name == "" {
			return &mcp.CallToolResult{IsError: true}, SiteOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		site, err := s.CreateSite(ctx, siteWriteFromCreate(in))
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
	return func(ctx context.Context, req *mcp.CallToolRequest, in SiteUpdateInput) (*mcp.CallToolResult, SiteOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, SiteOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, SiteOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		site, err := s.UpdateSite(ctx, in.ID, siteWriteFromUpdate(in))
		if err != nil {
			res, e := writeErrorResult(err)
			return res, SiteOutput{}, e
		}

		return &mcp.CallToolResult{}, SiteOutput{Data: *site}, nil
	}
}

// NewDeleteSiteHandler creates a handler for the delete_site tool.
func NewDeleteSiteHandler(svc *application.NetworkService) mcp.ToolHandlerFor[SiteDeleteInput, struct{}] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in SiteDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := s.DeleteSite(ctx, in.ID); err != nil {
			res, e := writeErrorResult(err)
			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateDeviceHandler creates a handler for the create_device tool.
func NewCreateDeviceHandler(svc *application.NetworkService) mcp.ToolHandlerFor[DeviceCreateInput, DeviceOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in DeviceCreateInput) (*mcp.CallToolResult, DeviceOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, DeviceOutput{}, errServiceNotAvailable
		}
		if in.Name == "" {
			return &mcp.CallToolResult{IsError: true}, DeviceOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		device, err := s.CreateDevice(ctx, deviceWriteFromCreate(in))
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
	return func(ctx context.Context, req *mcp.CallToolRequest, in DeviceUpdateInput) (*mcp.CallToolResult, DeviceOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, DeviceOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, DeviceOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		device, err := s.UpdateDevice(ctx, in.ID, deviceWriteFromUpdate(in))
		if err != nil {
			res, e := writeErrorResult(err)
			return res, DeviceOutput{}, e
		}

		return &mcp.CallToolResult{}, DeviceOutput{Data: *device}, nil
	}
}

// NewDeleteDeviceHandler creates a handler for the delete_device tool.
func NewDeleteDeviceHandler(svc *application.NetworkService) mcp.ToolHandlerFor[DeviceDeleteInput, struct{}] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in DeviceDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := s.DeleteDevice(ctx, in.ID); err != nil {
			res, e := writeErrorResult(err)
			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateIPAddressHandler creates a handler for the create_ip_address tool.
func NewCreateIPAddressHandler(svc *application.NetworkService) mcp.ToolHandlerFor[IPAddressCreateInput, IPAddressOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in IPAddressCreateInput) (*mcp.CallToolResult, IPAddressOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, IPAddressOutput{}, errServiceNotAvailable
		}
		if in.Address == "" {
			return &mcp.CallToolResult{IsError: true}, IPAddressOutput{}, fmt.Errorf("address is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		ip, err := s.CreateIPAddress(ctx, ipAddressWriteFromCreate(in))
		if err != nil {
			res, e := writeErrorResult(err)
			return res, IPAddressOutput{}, e
		}

		return &mcp.CallToolResult{}, IPAddressOutput{Data: *ip}, nil
	}
}

// NewUpdateIPAddressHandler creates a handler for the update_ip_address tool. It
// performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateIPAddressHandler(svc *application.NetworkService) mcp.ToolHandlerFor[IPAddressUpdateInput, IPAddressOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in IPAddressUpdateInput) (*mcp.CallToolResult, IPAddressOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, IPAddressOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, IPAddressOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		ip, err := s.UpdateIPAddress(ctx, in.ID, ipAddressWriteFromUpdate(in))
		if err != nil {
			res, e := writeErrorResult(err)
			return res, IPAddressOutput{}, e
		}

		return &mcp.CallToolResult{}, IPAddressOutput{Data: *ip}, nil
	}
}

// NewDeleteIPAddressHandler creates a handler for the delete_ip_address tool.
func NewDeleteIPAddressHandler(svc *application.NetworkService) mcp.ToolHandlerFor[IPAddressDeleteInput, struct{}] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in IPAddressDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := s.DeleteIPAddress(ctx, in.ID); err != nil {
			res, e := writeErrorResult(err)
			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreatePrefixHandler creates a handler for the create_prefix tool.
func NewCreatePrefixHandler(svc *application.NetworkService) mcp.ToolHandlerFor[PrefixCreateInput, PrefixOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in PrefixCreateInput) (*mcp.CallToolResult, PrefixOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, PrefixOutput{}, errServiceNotAvailable
		}
		if in.Prefix == "" {
			return &mcp.CallToolResult{IsError: true}, PrefixOutput{}, fmt.Errorf("prefix is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		prefix, err := s.CreatePrefix(ctx, prefixWriteFromCreate(in))
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
	return func(ctx context.Context, req *mcp.CallToolRequest, in PrefixUpdateInput) (*mcp.CallToolResult, PrefixOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, PrefixOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, PrefixOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		prefix, err := s.UpdatePrefix(ctx, in.ID, prefixWriteFromUpdate(in))
		if err != nil {
			res, e := writeErrorResult(err)
			return res, PrefixOutput{}, e
		}

		return &mcp.CallToolResult{}, PrefixOutput{Data: *prefix}, nil
	}
}

// NewDeletePrefixHandler creates a handler for the delete_prefix tool.
func NewDeletePrefixHandler(svc *application.NetworkService) mcp.ToolHandlerFor[PrefixDeleteInput, struct{}] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in PrefixDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := s.DeletePrefix(ctx, in.ID); err != nil {
			res, e := writeErrorResult(err)
			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateVLANHandler creates a handler for the create_vlan tool.
func NewCreateVLANHandler(svc *application.NetworkService) mcp.ToolHandlerFor[VLANCreateInput, VLANOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in VLANCreateInput) (*mcp.CallToolResult, VLANOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, VLANOutput{}, errServiceNotAvailable
		}
		if in.Name == "" {
			return &mcp.CallToolResult{IsError: true}, VLANOutput{}, fmt.Errorf("name is required")
		}
		if in.VID <= 0 {
			return &mcp.CallToolResult{IsError: true}, VLANOutput{}, fmt.Errorf("vid must be a positive integer")
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		vlan, err := s.CreateVLAN(ctx, vlanWriteFromCreate(in))
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
	return func(ctx context.Context, req *mcp.CallToolRequest, in VLANUpdateInput) (*mcp.CallToolResult, VLANOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, VLANOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, VLANOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		vlan, err := s.UpdateVLAN(ctx, in.ID, vlanWriteFromUpdate(in))
		if err != nil {
			res, e := writeErrorResult(err)
			return res, VLANOutput{}, e
		}

		return &mcp.CallToolResult{}, VLANOutput{Data: *vlan}, nil
	}
}

// NewDeleteVLANHandler creates a handler for the delete_vlan tool.
func NewDeleteVLANHandler(svc *application.NetworkService) mcp.ToolHandlerFor[VLANDeleteInput, struct{}] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in VLANDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := s.DeleteVLAN(ctx, in.ID); err != nil {
			res, e := writeErrorResult(err)
			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateVirtualMachineHandler creates a handler for the
// create_virtual_machine tool.
func NewCreateVirtualMachineHandler(svc *application.NetworkService) mcp.ToolHandlerFor[VirtualMachineCreateInput, VirtualMachineOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in VirtualMachineCreateInput) (*mcp.CallToolResult, VirtualMachineOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, VirtualMachineOutput{}, errServiceNotAvailable
		}
		if in.Name == "" {
			return &mcp.CallToolResult{IsError: true}, VirtualMachineOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		vm, err := s.CreateVirtualMachine(ctx, virtualMachineWriteFromCreate(in))
		if err != nil {
			res, e := writeErrorResult(err)
			return res, VirtualMachineOutput{}, e
		}

		return &mcp.CallToolResult{}, VirtualMachineOutput{Data: *vm}, nil
	}
}

// NewUpdateVirtualMachineHandler creates a handler for the
// update_virtual_machine tool. It performs a partial-merge PATCH using only the
// explicitly provided fields.
func NewUpdateVirtualMachineHandler(svc *application.NetworkService) mcp.ToolHandlerFor[VirtualMachineUpdateInput, VirtualMachineOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in VirtualMachineUpdateInput) (*mcp.CallToolResult, VirtualMachineOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, VirtualMachineOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, VirtualMachineOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		vm, err := s.UpdateVirtualMachine(ctx, in.ID, virtualMachineWriteFromUpdate(in))
		if err != nil {
			res, e := writeErrorResult(err)
			return res, VirtualMachineOutput{}, e
		}

		return &mcp.CallToolResult{}, VirtualMachineOutput{Data: *vm}, nil
	}
}

// NewDeleteVirtualMachineHandler creates a handler for the
// delete_virtual_machine tool.
func NewDeleteVirtualMachineHandler(svc *application.NetworkService) mcp.ToolHandlerFor[VirtualMachineDeleteInput, struct{}] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in VirtualMachineDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := s.DeleteVirtualMachine(ctx, in.ID); err != nil {
			res, e := writeErrorResult(err)
			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateClusterHandler creates a handler for the create_cluster tool.
func NewCreateClusterHandler(svc *application.NetworkService) mcp.ToolHandlerFor[ClusterCreateInput, ClusterOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in ClusterCreateInput) (*mcp.CallToolResult, ClusterOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, ClusterOutput{}, errServiceNotAvailable
		}
		if in.Name == "" {
			return &mcp.CallToolResult{IsError: true}, ClusterOutput{}, fmt.Errorf("name is required")
		}
		if in.ClusterType == nil {
			return &mcp.CallToolResult{IsError: true}, ClusterOutput{}, fmt.Errorf("type (cluster type) is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		cluster, err := s.CreateCluster(ctx, clusterWriteFromCreate(in))
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
	return func(ctx context.Context, req *mcp.CallToolRequest, in ClusterUpdateInput) (*mcp.CallToolResult, ClusterOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, ClusterOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, ClusterOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		cluster, err := s.UpdateCluster(ctx, in.ID, clusterWriteFromUpdate(in))
		if err != nil {
			res, e := writeErrorResult(err)
			return res, ClusterOutput{}, e
		}

		return &mcp.CallToolResult{}, ClusterOutput{Data: *cluster}, nil
	}
}

// NewDeleteClusterHandler creates a handler for the delete_cluster tool.
func NewDeleteClusterHandler(svc *application.NetworkService) mcp.ToolHandlerFor[ClusterDeleteInput, struct{}] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in ClusterDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := s.DeleteCluster(ctx, in.ID); err != nil {
			res, e := writeErrorResult(err)
			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateCircuitHandler creates a handler for the create_circuit tool.
func NewCreateCircuitHandler(svc *application.NetworkService) mcp.ToolHandlerFor[CircuitCreateInput, CircuitOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in CircuitCreateInput) (*mcp.CallToolResult, CircuitOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitOutput{}, errServiceNotAvailable
		}
		if in.CID == "" {
			return &mcp.CallToolResult{IsError: true}, CircuitOutput{}, fmt.Errorf("cid is required")
		}
		if in.Provider == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitOutput{}, fmt.Errorf("provider is required")
		}
		if in.CircuitType == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitOutput{}, fmt.Errorf("circuit_type is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		circuit, err := s.CreateCircuit(ctx, circuitWriteFromCreate(in))
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
	return func(ctx context.Context, req *mcp.CallToolRequest, in CircuitUpdateInput) (*mcp.CallToolResult, CircuitOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, CircuitOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		circuit, err := s.UpdateCircuit(ctx, in.ID, circuitWriteFromUpdate(in))
		if err != nil {
			res, e := writeErrorResult(err)
			return res, CircuitOutput{}, e
		}

		return &mcp.CallToolResult{}, CircuitOutput{Data: *circuit}, nil
	}
}

// NewDeleteCircuitHandler creates a handler for the delete_circuit tool.
func NewDeleteCircuitHandler(svc *application.NetworkService) mcp.ToolHandlerFor[CircuitDeleteInput, struct{}] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in CircuitDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := s.DeleteCircuit(ctx, in.ID); err != nil {
			res, e := writeErrorResult(err)
			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateRackHandler creates a handler for the create_rack tool.
func NewCreateRackHandler(svc *application.NetworkService) mcp.ToolHandlerFor[RackCreateInput, RackOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in RackCreateInput) (*mcp.CallToolResult, RackOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, RackOutput{}, errServiceNotAvailable
		}
		if in.Name == "" {
			return &mcp.CallToolResult{IsError: true}, RackOutput{}, fmt.Errorf("name is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		rack, err := s.CreateRack(ctx, rackWriteFromCreate(in))
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
	return func(ctx context.Context, req *mcp.CallToolRequest, in RackUpdateInput) (*mcp.CallToolResult, RackOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, RackOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, RackOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		rack, err := s.UpdateRack(ctx, in.ID, rackWriteFromUpdate(in))
		if err != nil {
			res, e := writeErrorResult(err)
			return res, RackOutput{}, e
		}

		return &mcp.CallToolResult{}, RackOutput{Data: *rack}, nil
	}
}

// NewDeleteRackHandler creates a handler for the delete_rack tool.
func NewDeleteRackHandler(svc *application.NetworkService) mcp.ToolHandlerFor[RackDeleteInput, struct{}] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in RackDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := s.DeleteRack(ctx, in.ID); err != nil {
			res, e := writeErrorResult(err)
			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateInterfaceHandler creates a handler for the create_interface tool.
func NewCreateInterfaceHandler(svc *application.NetworkService) mcp.ToolHandlerFor[InterfaceCreateInput, InterfaceOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in InterfaceCreateInput) (*mcp.CallToolResult, InterfaceOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, InterfaceOutput{}, errServiceNotAvailable
		}
		if in.Name == "" {
			return &mcp.CallToolResult{IsError: true}, InterfaceOutput{}, fmt.Errorf("name is required")
		}
		if in.Device == nil {
			return &mcp.CallToolResult{IsError: true}, InterfaceOutput{}, fmt.Errorf("device is required")
		}
		if in.Type == nil {
			return &mcp.CallToolResult{IsError: true}, InterfaceOutput{}, fmt.Errorf("type is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		iface, err := s.CreateInterface(ctx, interfaceWriteFromCreate(in))
		if err != nil {
			res, e := writeErrorResult(err)
			return res, InterfaceOutput{}, e
		}

		return &mcp.CallToolResult{}, InterfaceOutput{Data: *iface}, nil
	}
}

// NewUpdateInterfaceHandler creates a handler for the update_interface tool. It
// performs a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateInterfaceHandler(svc *application.NetworkService) mcp.ToolHandlerFor[InterfaceUpdateInput, InterfaceOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in InterfaceUpdateInput) (*mcp.CallToolResult, InterfaceOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, InterfaceOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, InterfaceOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		iface, err := s.UpdateInterface(ctx, in.ID, interfaceWriteFromUpdate(in))
		if err != nil {
			res, e := writeErrorResult(err)
			return res, InterfaceOutput{}, e
		}

		return &mcp.CallToolResult{}, InterfaceOutput{Data: *iface}, nil
	}
}

// NewDeleteInterfaceHandler creates a handler for the delete_interface tool.
func NewDeleteInterfaceHandler(svc *application.NetworkService) mcp.ToolHandlerFor[InterfaceDeleteInput, struct{}] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in InterfaceDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := s.DeleteInterface(ctx, in.ID); err != nil {
			res, e := writeErrorResult(err)
			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateCircuitTerminationHandler creates a handler for the
// create_circuit_termination tool.
func NewCreateCircuitTerminationHandler(svc *application.NetworkService) mcp.ToolHandlerFor[CircuitTerminationCreateInput, CircuitTerminationOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in CircuitTerminationCreateInput) (*mcp.CallToolResult, CircuitTerminationOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitTerminationOutput{}, errServiceNotAvailable
		}
		if in.TermSide == "" {
			return &mcp.CallToolResult{IsError: true}, CircuitTerminationOutput{}, fmt.Errorf("term_side is required")
		}
		if in.Circuit == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitTerminationOutput{}, fmt.Errorf("circuit is required")
		}
		if in.Site == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitTerminationOutput{}, fmt.Errorf("site is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		ct, err := s.CreateCircuitTermination(ctx, circuitTerminationWriteFromCreate(in))
		if err != nil {
			res, e := writeErrorResult(err)
			return res, CircuitTerminationOutput{}, e
		}

		return &mcp.CallToolResult{}, CircuitTerminationOutput{Data: *ct}, nil
	}
}

// NewUpdateCircuitTerminationHandler creates a handler for the
// update_circuit_termination tool. It performs a partial-merge PATCH using only
// the explicitly provided fields.
func NewUpdateCircuitTerminationHandler(svc *application.NetworkService) mcp.ToolHandlerFor[CircuitTerminationUpdateInput, CircuitTerminationOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in CircuitTerminationUpdateInput) (*mcp.CallToolResult, CircuitTerminationOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitTerminationOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, CircuitTerminationOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		ct, err := s.UpdateCircuitTermination(ctx, in.ID, circuitTerminationWriteFromUpdate(in))
		if err != nil {
			res, e := writeErrorResult(err)
			return res, CircuitTerminationOutput{}, e
		}

		return &mcp.CallToolResult{}, CircuitTerminationOutput{Data: *ct}, nil
	}
}

// NewDeleteCircuitTerminationHandler creates a handler for the
// delete_circuit_termination tool.
func NewDeleteCircuitTerminationHandler(svc *application.NetworkService) mcp.ToolHandlerFor[CircuitTerminationDeleteInput, struct{}] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in CircuitTerminationDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := s.DeleteCircuitTermination(ctx, in.ID); err != nil {
			res, e := writeErrorResult(err)
			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}

// NewCreateCableHandler creates a handler for the create_cable tool.
func NewCreateCableHandler(svc *application.NetworkService) mcp.ToolHandlerFor[CableCreateInput, CableOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in CableCreateInput) (*mcp.CallToolResult, CableOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, CableOutput{}, errServiceNotAvailable
		}
		if in.TerminationA == nil {
			return &mcp.CallToolResult{IsError: true}, CableOutput{}, fmt.Errorf("termination_a is required")
		}
		if in.TerminationB == nil {
			return &mcp.CallToolResult{IsError: true}, CableOutput{}, fmt.Errorf("termination_b is required")
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		c, err := s.CreateCable(ctx, cableWriteFromCreate(in))
		if err != nil {
			res, e := writeErrorResult(err)
			return res, CableOutput{}, e
		}

		return &mcp.CallToolResult{}, CableOutput{Data: *c}, nil
	}
}

// NewUpdateCableHandler creates a handler for the update_cable tool. It performs
// a partial-merge PATCH using only the explicitly provided fields.
func NewUpdateCableHandler(svc *application.NetworkService) mcp.ToolHandlerFor[CableUpdateInput, CableOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in CableUpdateInput) (*mcp.CallToolResult, CableOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, CableOutput{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, CableOutput{}, err
		}

		sanitizeAllStrings(reflect.ValueOf(&in))
		c, err := s.UpdateCable(ctx, in.ID, cableWriteFromUpdate(in))
		if err != nil {
			res, e := writeErrorResult(err)
			return res, CableOutput{}, e
		}

		return &mcp.CallToolResult{}, CableOutput{Data: *c}, nil
	}
}

// NewDeleteCableHandler creates a handler for the delete_cable tool.
func NewDeleteCableHandler(svc *application.NetworkService) mcp.ToolHandlerFor[CableDeleteInput, struct{}] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in CableDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, errServiceNotAvailable
		}
		if err := validatePositiveID(in.ID); err != nil {
			return &mcp.CallToolResult{IsError: true}, struct{}{}, err
		}

		if err := s.DeleteCable(ctx, in.ID); err != nil {
			res, e := writeErrorResult(err)
			return res, struct{}{}, e
		}

		return &mcp.CallToolResult{}, struct{}{}, nil
	}
}
