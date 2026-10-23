package netbox

// WireSiteWrite is the wire request model for site create/update (POST/PATCH
// body). Its json tags intentionally mirror NetBox's writable fields.
//
// Wire types mirror NetBox JSON API responses. The domain.PaginatedResponse
// type is used for both deserialization and output — wire types are embedded
// via type parameter instantiation (e.g. domain.PaginatedResponse[WireSite]).
type WireSiteWrite struct {
	Name            string         `json:"name,omitempty"`
	Slug            *string        `json:"slug,omitempty"`
	Status          *string        `json:"status,omitempty"`
	Region          *int           `json:"region,omitempty"`
	Tenant          *int           `json:"tenant,omitempty"`
	Facility        *string        `json:"facility,omitempty"`
	TimeZone        *string        `json:"time_zone,omitempty"`
	Description     *string        `json:"description,omitempty"`
	PhysicalAddress *string        `json:"physical_address,omitempty"`
	ShippingAddress *string        `json:"shipping_address,omitempty"`
	Comments        *string        `json:"comments,omitempty"`
	Tags            []string       `json:"tags,omitempty"`
	CustomFields    map[string]any `json:"custom_fields,omitempty"`
}

// WireRegionWrite is the wire request model for region create/update (POST/PATCH
// body). Its json tags intentionally mirror NetBox's writable fields.
type WireRegionWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Parent       *int           `json:"parent,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireTenantGroupWrite is the wire request model for tenant group create/update
// (POST/PATCH body). Its json tags intentionally mirror NetBox's writable
// fields.
type WireTenantGroupWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Parent       *int           `json:"parent,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireDeviceWrite is the wire request model for device create/update (POST/PATCH
// body). Its json tags intentionally mirror NetBox's writable fields. Related
// objects are expressed as numeric IDs.
type WireDeviceWrite struct {
	Name         string         `json:"name,omitempty"`
	DeviceType   *int           `json:"device_type,omitempty"`
	Role         *int           `json:"role,omitempty"`
	Tenant       *int           `json:"tenant,omitempty"`
	Platform     *int           `json:"platform,omitempty"`
	Serial       *string        `json:"serial,omitempty"`
	AssetTag     *string        `json:"asset_tag,omitempty"`
	Site         *int           `json:"site,omitempty"`
	Rack         *int           `json:"rack,omitempty"`
	Position     *float64       `json:"position,omitempty"`
	Face         *string        `json:"face,omitempty"`
	Status       *string        `json:"status,omitempty"`
	Cluster      *int           `json:"cluster,omitempty"`
	Comments     *string        `json:"comments,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireIPAddressWrite is the wire request model for IP address create/update
// (POST/PATCH body). Its json tags intentionally mirror NetBox's writable
// fields. References to related objects are numeric IDs.
type WireIPAddressWrite struct {
	Address            string         `json:"address,omitempty"`
	Status             *string        `json:"status,omitempty"`
	Role               *string        `json:"role,omitempty"`
	VRF                *int           `json:"vrf,omitempty"`
	Tenant             *int           `json:"tenant,omitempty"`
	DNSName            *string        `json:"dns_name,omitempty"`
	Description        *string        `json:"description,omitempty"`
	AssignedObjectType *string        `json:"assigned_object_type,omitempty"`
	AssignedObjectID   *int           `json:"assigned_object_id,omitempty"`
	Comments           *string        `json:"comments,omitempty"`
	Tags               []string       `json:"tags,omitempty"`
	CustomFields       map[string]any `json:"custom_fields,omitempty"`
}

// WirePrefixWrite is the wire request model for prefix create/update (POST/PATCH
// body). Its json tags intentionally mirror NetBox's writable fields.
// References to related objects are numeric IDs.
type WirePrefixWrite struct {
	Prefix       string         `json:"prefix,omitempty"`
	Site         *int           `json:"site,omitempty"`
	VRF          *int           `json:"vrf,omitempty"`
	Tenant       *int           `json:"tenant,omitempty"`
	VLAN         *int           `json:"vlan,omitempty"`
	Status       *string        `json:"status,omitempty"`
	Role         *int           `json:"role,omitempty"`
	IsPool       *bool          `json:"is_pool,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Comments     *string        `json:"comments,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireVLANWrite is the wire request model for VLAN create/update (POST/PATCH
// body). Its json tags intentionally mirror NetBox's writable fields.
// References to related objects are numeric IDs.
type WireVLANWrite struct {
	VID          int            `json:"vid,omitempty"`
	Name         string         `json:"name,omitempty"`
	Site         *int           `json:"site,omitempty"`
	Group        *int           `json:"group,omitempty"`
	Tenant       *int           `json:"tenant,omitempty"`
	Status       *string        `json:"status,omitempty"`
	Role         *int           `json:"role,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Comments     *string        `json:"comments,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireVirtualMachineWrite is the wire request model for virtual machine
// create/update (POST/PATCH body). Its json tags intentionally mirror NetBox's
// writable fields. References to related objects are numeric IDs.
type WireVirtualMachineWrite struct {
	Name         string         `json:"name,omitempty"`
	Cluster      *int           `json:"cluster,omitempty"`
	Role         *int           `json:"role,omitempty"`
	Tenant       *int           `json:"tenant,omitempty"`
	Platform     *int           `json:"platform,omitempty"`
	Status       *string        `json:"status,omitempty"`
	Site         *int           `json:"site,omitempty"`
	VCPUs        *float64       `json:"vcpus,omitempty"`
	Memory       *int           `json:"memory,omitempty"`
	Disk         *int           `json:"disk,omitempty"`
	Comments     *string        `json:"comments,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireClusterWrite is the wire request model for cluster create/update
// (POST/PATCH body). Its json tags intentionally mirror NetBox's writable
// fields. References to related objects (type, group, site, tenant) are numeric
// IDs.
type WireClusterWrite struct {
	Name         string         `json:"name,omitempty"`
	ClusterType  *int           `json:"type,omitempty"`
	ClusterGroup *int           `json:"group,omitempty"`
	Site         *int           `json:"site,omitempty"`
	Tenant       *int           `json:"tenant,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Comments     *string        `json:"comments,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireCircuitWrite is the wire request model for circuit create/update
// (POST/PATCH body). Its json tags intentionally mirror NetBox's writable
// fields. References to related objects (provider, circuit_type, tenant) are
// numeric IDs.
type WireCircuitWrite struct {
	CID          string         `json:"cid,omitempty"`
	Provider     *int           `json:"provider,omitempty"`
	CircuitType  *int           `json:"type,omitempty"`
	Tenant       *int           `json:"tenant,omitempty"`
	Status       *string        `json:"status,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Comments     *string        `json:"comments,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
	InstallDate  *string        `json:"install_date,omitempty"`
	CommitRate   *int           `json:"commit_rate,omitempty"`
}

// WireRackWrite is the wire request model for rack create/update (POST/PATCH
// body). Its json tags intentionally mirror NetBox's writable fields.
// References to related objects (site, location, tenant, role) are numeric IDs.
type WireRackWrite struct {
	Name         string         `json:"name,omitempty"`
	FacilityID   *string        `json:"facility_id,omitempty"`
	Site         *int           `json:"site,omitempty"`
	Location     *int           `json:"location,omitempty"`
	Tenant       *int           `json:"tenant,omitempty"`
	Status       *string        `json:"status,omitempty"`
	Role         *int           `json:"role,omitempty"`
	Serial       *string        `json:"serial,omitempty"`
	AssetTag     *string        `json:"asset_tag,omitempty"`
	Type         *string        `json:"type,omitempty"`
	Width        *int           `json:"width,omitempty"`
	UHeight      *int           `json:"u_height,omitempty"`
	Comments     *string        `json:"comments,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireInterfaceWrite is the wire request model for interface create/update
// (POST/PATCH body). Its json tags intentionally mirror NetBox's writable
// fields. Device is the numeric ID of the parent device; type is a choice
// string.
type WireInterfaceWrite struct {
	Device       *int           `json:"device,omitempty"`
	Name         string         `json:"name,omitempty"`
	Type         *string        `json:"type,omitempty"`
	Enabled      *bool          `json:"enabled,omitempty"`
	MTU          *int           `json:"mtu,omitempty"`
	MACAddress   *string        `json:"mac_address,omitempty"`
	Speed        *int           `json:"speed,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireCircuitTerminationWrite is the wire request model for circuit
// termination create/update (POST/PATCH body). Its json tags intentionally
// mirror NetBox's writable fields. Circuit is a numeric ID; term_side is A or
// Z. NetBox 4.6 represents the terminating object generically: termination_type
// is the content type (e.g. "dcim.site") and termination_id is the object's
// numeric ID.
type WireCircuitTerminationWrite struct {
	Circuit         *int           `json:"circuit,omitempty"`
	TermSide        string         `json:"term_side,omitempty"`
	TerminationType string         `json:"termination_type,omitempty"`
	TerminationID   *int           `json:"termination_id,omitempty"`
	Speed           *int           `json:"speed,omitempty"`
	UpstreamSpeed   *int           `json:"upstream_speed,omitempty"`
	Description     *string        `json:"description,omitempty"`
	Tags            []string       `json:"tags,omitempty"`
	CustomFields    map[string]any `json:"custom_fields,omitempty"`
}

// WireCableWriteTermination is the wire request model for one end of a cable
// create/update. NetBox accepts a nested object with object_type and object_id.
type WireCableWriteTermination struct {
	ObjectType string `json:"object_type"`
	ObjectID   int    `json:"object_id"`
}

// WireCableWrite is the wire request model for cable create/update (POST/PATCH
// body). Its json tags intentionally mirror NetBox's writable fields.
type WireCableWrite struct {
	TerminationA *WireCableWriteTermination `json:"termination_a,omitempty"`
	TerminationB *WireCableWriteTermination `json:"termination_b,omitempty"`
	Type         *string                    `json:"type,omitempty"`
	Status       *string                    `json:"status,omitempty"`
	Label        *string                    `json:"label,omitempty"`
	Color        *string                    `json:"color,omitempty"`
	Length       *float64                   `json:"length,omitempty"`
	LengthUnit   *string                    `json:"length_unit,omitempty"`
	Description  *string                    `json:"description,omitempty"`
	Tags         []string                   `json:"tags,omitempty"`
	CustomFields map[string]any             `json:"custom_fields,omitempty"`
}

// WireVMInterfaceWrite is the wire request model for VM interface
// create/update (POST/PATCH body). Its json tags intentionally mirror NetBox's
// writable fields. VirtualMachine is a numeric ID.
type WireVMInterfaceWrite struct {
	VirtualMachine *int           `json:"virtual_machine,omitempty"`
	Name           string         `json:"name,omitempty"`
	Enabled        *bool          `json:"enabled,omitempty"`
	MTU            *int           `json:"mtu,omitempty"`
	MACAddress     *string        `json:"mac_address,omitempty"`
	Description    *string        `json:"description,omitempty"`
	Tags           []string       `json:"tags,omitempty"`
	CustomFields   map[string]any `json:"custom_fields,omitempty"`
}

// WireProviderWrite is the wire request model for provider create/update
// (POST/PATCH body). Its json tags intentionally mirror NetBox's writable
// fields. Asn is a numeric AS number.
type WireProviderWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Asn          *int           `json:"asn,omitempty"`
	Account      *string        `json:"account,omitempty"`
	PortalURL    *string        `json:"portal_url,omitempty"`
	NocContact   *string        `json:"noc_contact,omitempty"`
	AdminContact *string        `json:"admin_contact,omitempty"`
	Comments     *string        `json:"comments,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireTenantWrite is the wire request model for tenant create/update
// (POST/PATCH body). Its json tags intentionally mirror NetBox's writable
// fields.
type WireTenantWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Group        *int           `json:"group,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Comments     *string        `json:"comments,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireManufacturerWrite is the wire request model for manufacturer create/update
// (POST/PATCH body). Its json tags intentionally mirror NetBox's writable
// fields.
type WireManufacturerWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireDeviceTypeWrite is the wire request model for device type create/update
// (POST/PATCH body). Its json tags intentionally mirror NetBox's writable
// fields. The model field is named "model" (not "name").
type WireDeviceTypeWrite struct {
	Manufacturer  *int           `json:"manufacturer,omitempty"`
	Model         string         `json:"model,omitempty"`
	Slug          *string        `json:"slug,omitempty"`
	PartNumber    *string        `json:"part_number,omitempty"`
	UHeight       *float64       `json:"u_height,omitempty"`
	IsFullDepth   *bool          `json:"is_full_depth,omitempty"`
	SubdeviceRole *string        `json:"subdevice_role,omitempty"`
	Comments      *string        `json:"comments,omitempty"`
	Tags          []string       `json:"tags,omitempty"`
	CustomFields  map[string]any `json:"custom_fields,omitempty"`
}

// WireLocationWrite is the wire request model for location create/update
// (POST/PATCH body). Its json tags intentionally mirror NetBox's writable
// fields. Site (numeric NetBox ID) is required on create.
type WireLocationWrite struct {
	Name         string         `json:"name,omitempty"`
	Site         *int           `json:"site,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Parent       *int           `json:"parent,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Status       *string        `json:"status,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireClusterTypeWrite is the wire request model for cluster type create/update
// (POST/PATCH body). Its json tags intentionally mirror NetBox's writable
// fields.
type WireClusterTypeWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireClusterGroupWrite is the wire request model for cluster group create/update
// (POST/PATCH body). Its json tags intentionally mirror NetBox's writable
// fields.
type WireClusterGroupWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireCircuitTypeWrite is the wire request model for circuit type create/update
// (POST/PATCH body). Its json tags intentionally mirror NetBox's writable
// fields.
type WireCircuitTypeWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireVrfWrite is the wire request model for VRF create/update (POST/PATCH
// body). Its json tags intentionally mirror NetBox's writable fields.
type WireVrfWrite struct {
	Name         string         `json:"name,omitempty"`
	Rd           string         `json:"rd,omitempty"`
	Tenant       *int           `json:"tenant,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireVlanGroupWrite is the wire request model for VLAN group create/update
// (POST/PATCH body). Its json tags intentionally mirror NetBox's writable
// fields.
type WireVlanGroupWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Description  *string        `json:"description,omitempty"`
	ScopeType    *string        `json:"scope_type,omitempty"`
	ScopeID      *int           `json:"scope_id,omitempty"`
	MinVID       *int           `json:"min_vid,omitempty"`
	MaxVID       *int           `json:"max_vid,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireRoleWrite is the wire request model for IPAM role create/update (POST/
// PATCH body). Its json tags intentionally mirror NetBox's writable fields.
type WireRoleWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Weight       *int           `json:"weight,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireContactWrite is the wire request model for tenancy contact create/update
// (POST/PATCH body). Its json tags intentionally mirror NetBox's writable
// fields.
type WireContactWrite struct {
	Name         string         `json:"name,omitempty"`
	Group        *int           `json:"group,omitempty"`
	Title        *string        `json:"title,omitempty"`
	Phone        *string        `json:"phone,omitempty"`
	Email        *string        `json:"email,omitempty"`
	Address      *string        `json:"address,omitempty"`
	Link         *string        `json:"link,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Comments     *string        `json:"comments,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireContactGroupWrite is the wire request model for tenancy contact group
// create/update (POST/PATCH body). Its json tags intentionally mirror NetBox's
// writable fields. Parent is the numeric NetBox ID of the parent contact group.
type WireContactGroupWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Parent       *int           `json:"parent,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireContactRoleWrite is the wire request model for tenancy contact role
// create/update (POST/PATCH body). Its json tags intentionally mirror NetBox's
// writable fields.
type WireContactRoleWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireVirtualDiskWrite is the wire request model for virtualization virtual
// disk create/update (POST/PATCH body). Its json tags intentionally mirror
// NetBox's writable fields. VirtualMachine is the numeric NetBox ID of the
// parent virtual machine.
type WireVirtualDiskWrite struct {
	Name           string         `json:"name,omitempty"`
	VirtualMachine *int           `json:"virtual_machine,omitempty"`
	Size           *int           `json:"size,omitempty"`
	Description    *string        `json:"description,omitempty"`
	Tags           []string       `json:"tags,omitempty"`
	CustomFields   map[string]any `json:"custom_fields,omitempty"`
}

// WireVirtualMachineTypeWrite is the wire request model for virtualization
// virtual machine type create/update (POST/PATCH body). Its json tags
// intentionally mirror NetBox's writable fields.
type WireVirtualMachineTypeWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireSiteGroupWrite is the wire request model for DCIM site group
// create/update (POST/PATCH body). Its json tags intentionally mirror NetBox's
// writable fields.
type WireSiteGroupWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Parent       *int           `json:"parent,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireRackGroupWrite is the wire request model for DCIM rack group
// create/update (POST/PATCH body). Its json tags intentionally mirror NetBox's
// writable fields.
type WireRackGroupWrite struct {
	Name        string   `json:"name,omitempty"`
	Slug        *string  `json:"slug,omitempty"`
	Site        *int     `json:"site,omitempty"`
	Parent      *int     `json:"parent,omitempty"`
	Description *string  `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// WireRackRoleWrite is the wire request model for DCIM rack role create/update
// (POST/PATCH body). Its json tags intentionally mirror NetBox's writable
// fields.
type WireRackRoleWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Color        *string        `json:"color,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireDeviceRoleWrite is the wire request model for DCIM device role
// create/update (POST/PATCH body). Its json tags intentionally mirror NetBox's
// writable fields.
type WireDeviceRoleWrite struct {
	Name           string         `json:"name,omitempty"`
	Slug           *string        `json:"slug,omitempty"`
	Color          *string        `json:"color,omitempty"`
	VMRole         *bool          `json:"vm_role,omitempty"`
	ConfigTemplate *int           `json:"config_template,omitempty"`
	Description    *string        `json:"description,omitempty"`
	Tags           []string       `json:"tags,omitempty"`
	CustomFields   map[string]any `json:"custom_fields,omitempty"`
	Weight         *int           `json:"weight,omitempty"`
}

// WirePlatformWrite is the wire request model for DCIM platform create/update
// (POST/PATCH body). Its json tags intentionally mirror NetBox's writable
// fields.
type WirePlatformWrite struct {
	Name           string         `json:"name,omitempty"`
	Slug           *string        `json:"slug,omitempty"`
	Manufacturer   *int           `json:"manufacturer,omitempty"`
	ConfigTemplate *int           `json:"config_template,omitempty"`
	Description    *string        `json:"description,omitempty"`
	Tags           []string       `json:"tags,omitempty"`
	CustomFields   map[string]any `json:"custom_fields,omitempty"`
}

// WireVirtualChassisWrite is the wire request model for DCIM virtual chassis
// create/update (POST/PATCH body). Its json tags intentionally mirror NetBox's
// writable fields.
type WireVirtualChassisWrite struct {
	Name         string         `json:"name,omitempty"`
	Domain       *string        `json:"domain,omitempty"`
	Master       *int           `json:"master,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WirePowerPanelWrite is the wire request model for DCIM power panel
// create/update (POST/PATCH body). Its json tags intentionally mirror NetBox's
// writable fields.
type WirePowerPanelWrite struct {
	Site         *int           `json:"site,omitempty"`
	Location     *int           `json:"location,omitempty"`
	Name         string         `json:"name,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireInventoryItemRoleWrite is the wire request model for DCIM inventory item
// role create/update (POST/PATCH body). Its json tags intentionally mirror
// NetBox's writable fields.
type WireInventoryItemRoleWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Color        *string        `json:"color,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireFrontPortWrite is the wire request model for DCIM front port
// create/update (POST/PATCH body). Its json tags intentionally mirror NetBox's
// writable fields.
type WireFrontPortWrite struct {
	Name             string         `json:"name,omitempty"`
	Device           *int           `json:"device,omitempty"`
	Type             *string        `json:"type,omitempty"`
	RearPort         *int           `json:"rear_port,omitempty"`
	RearPortPosition *int           `json:"rear_port_position,omitempty"`
	Description      *string        `json:"description,omitempty"`
	Tags             []string       `json:"tags,omitempty"`
	CustomFields     map[string]any `json:"custom_fields,omitempty"`
}

// WireRearPortWrite is the wire request model for DCIM rear port create/update
// (POST/PATCH body). Its json tags intentionally mirror NetBox's writable
// fields.
type WireRearPortWrite struct {
	Name         string         `json:"name,omitempty"`
	Device       *int           `json:"device,omitempty"`
	Type         *string        `json:"type,omitempty"`
	Positions    *int           `json:"positions,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireConsolePortWrite is the wire request model for DCIM console port
// create/update (POST/PATCH body). Its json tags intentionally mirror NetBox's
// writable fields.
type WireConsolePortWrite struct {
	Name         string         `json:"name,omitempty"`
	Device       *int           `json:"device,omitempty"`
	Type         *string        `json:"type,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireConsoleServerPortWrite is the wire request model for DCIM console server
// port create/update (POST/PATCH body). Its json tags intentionally mirror
// NetBox's writable fields.
type WireConsoleServerPortWrite struct {
	Name         string         `json:"name,omitempty"`
	Device       *int           `json:"device,omitempty"`
	Type         *string        `json:"type,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WirePowerPortWrite is the wire request model for DCIM power port
// create/update (POST/PATCH body). Its json tags intentionally mirror NetBox's
// writable fields.
type WirePowerPortWrite struct {
	Name          string         `json:"name,omitempty"`
	Device        *int           `json:"device,omitempty"`
	Type          *string        `json:"type,omitempty"`
	MaximumDraw   *int           `json:"maximum_draw,omitempty"`
	AllocatedDraw *int           `json:"allocated_draw,omitempty"`
	Description   *string        `json:"description,omitempty"`
	Tags          []string       `json:"tags,omitempty"`
	CustomFields  map[string]any `json:"custom_fields,omitempty"`
}

// WirePowerOutletWrite is the wire request model for DCIM power outlet
// create/update (POST/PATCH body). Its json tags intentionally mirror NetBox's
// writable fields.
type WirePowerOutletWrite struct {
	Name         string         `json:"name,omitempty"`
	Device       *int           `json:"device,omitempty"`
	Type         *string        `json:"type,omitempty"`
	PowerPort    *int           `json:"power_port,omitempty"`
	FeedLeg      *string        `json:"feed_leg,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WirePowerFeedWrite is the wire request model for DCIM power feed
// create/update (POST/PATCH body). Its json tags intentionally mirror NetBox's
// writable fields.
type WirePowerFeedWrite struct {
	Name           string         `json:"name,omitempty"`
	PowerPanel     *int           `json:"power_panel,omitempty"`
	Status         *string        `json:"status,omitempty"`
	Type           *string        `json:"type,omitempty"`
	Supply         *string        `json:"supply,omitempty"`
	Phase          *string        `json:"phase,omitempty"`
	Voltage        *int           `json:"voltage,omitempty"`
	Amperage       *int           `json:"amperage,omitempty"`
	MaxUtilization *int           `json:"max_utilization,omitempty"`
	Description    *string        `json:"description,omitempty"`
	Tags           []string       `json:"tags,omitempty"`
	CustomFields   map[string]any `json:"custom_fields,omitempty"`
}

// WireDeviceBayWrite is the wire request model for DCIM device bay
// create/update (POST/PATCH body). Its json tags intentionally mirror NetBox's
// writable fields.
type WireDeviceBayWrite struct {
	Name         string         `json:"name,omitempty"`
	Device       *int           `json:"device,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireRackReservationWrite is the wire request model for DCIM rack reservation
// create/update (POST/PATCH body). Its json tags intentionally mirror NetBox's
// writable fields.
type WireRackReservationWrite struct {
	Description  string         `json:"description,omitempty"`
	Rack         *int           `json:"rack,omitempty"`
	Units        []int          `json:"units,omitempty"`
	User         *int           `json:"user,omitempty"`
	Tenant       *int           `json:"tenant,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireInventoryItemWrite is the wire request model for DCIM inventory item
// create/update (POST/PATCH body). Its json tags intentionally mirror NetBox's
// writable fields.
type WireInventoryItemWrite struct {
	Name         string         `json:"name,omitempty"`
	Device       *int           `json:"device,omitempty"`
	Role         *int           `json:"role,omitempty"`
	Manufacturer *int           `json:"manufacturer,omitempty"`
	Serial       *string        `json:"serial,omitempty"`
	PartID       *string        `json:"part_id,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireMACAddressWrite is the wire request model for DCIM MAC address
// create/update (POST/PATCH body). Its json tags intentionally mirror NetBox's
// writable fields.
type WireMACAddressWrite struct {
	Address            string         `json:"mac_address,omitempty"`
	AssignedObjectType *string        `json:"assigned_object_type,omitempty"`
	AssignedObjectID   *int           `json:"assigned_object_id,omitempty"`
	Description        *string        `json:"description,omitempty"`
	Tags               []string       `json:"tags,omitempty"`
	CustomFields       map[string]any `json:"custom_fields,omitempty"`
}

// WireRIRWrite is the wire request model for IPAM RIR create/update. Its json
// tags intentionally mirror NetBox's writable fields.
type WireRIRWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	IsPrivate    *bool          `json:"is_private,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireAggregateWrite is the wire request model for IPAM aggregate
// create/update. Its json tags intentionally mirror NetBox's writable fields.
type WireAggregateWrite struct {
	Prefix       string         `json:"prefix,omitempty"`
	RIR          *int           `json:"rir,omitempty"`
	Tenant       *int           `json:"tenant,omitempty"`
	DateAdded    *string        `json:"date_added,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireIPRangeWrite is the wire request model for IPAM IP range create/update.
// Its json tags intentionally mirror NetBox's writable fields.
type WireIPRangeWrite struct {
	StartAddress string         `json:"start_address,omitempty"`
	EndAddress   string         `json:"end_address,omitempty"`
	VRF          *int           `json:"vrf,omitempty"`
	Tenant       *int           `json:"tenant,omitempty"`
	Status       *string        `json:"status,omitempty"`
	Role         *int           `json:"role,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireASNWrite is the wire request model for IPAM ASN create/update. Its json
// tags intentionally mirror NetBox's writable fields.
type WireASNWrite struct {
	ASN          int            `json:"asn,omitempty"`
	RIR          *int           `json:"rir,omitempty"`
	Tenant       *int           `json:"tenant,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireASNRangeWrite is the wire request model for IPAM ASN range create/update.
// Its json tags intentionally mirror NetBox's writable fields.
type WireASNRangeWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Start        *int           `json:"start,omitempty"`
	End          *int           `json:"end,omitempty"`
	RIR          *int           `json:"rir,omitempty"`
	Tenant       *int           `json:"tenant,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireRouteTargetWrite is the wire request model for IPAM route target
// create/update. Its json tags intentionally mirror NetBox's writable fields.
type WireRouteTargetWrite struct {
	Name         string         `json:"name,omitempty"`
	Tenant       *int           `json:"tenant,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireServiceWrite is the wire request model for IPAM service create/update.
// Its json tags intentionally mirror NetBox's writable fields.
type WireServiceWrite struct {
	Name             string         `json:"name,omitempty"`
	ParentObjectType *string        `json:"parent_object_type,omitempty"`
	ParentObjectID   *int           `json:"parent_object_id,omitempty"`
	Protocol         *string        `json:"protocol,omitempty"`
	Ports            []int          `json:"ports,omitempty"`
	Description      *string        `json:"description,omitempty"`
	Tags             []string       `json:"tags,omitempty"`
	CustomFields     map[string]any `json:"custom_fields,omitempty"`
}

// WireServiceTemplateWrite is the wire request model for IPAM service template
// create/update. Its json tags intentionally mirror NetBox's writable fields.
type WireServiceTemplateWrite struct {
	Name         string         `json:"name,omitempty"`
	Protocol     *string        `json:"protocol,omitempty"`
	Ports        []int          `json:"ports,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireFHRPGroupWrite is the wire request model for IPAM FHRP group
// create/update. Its json tags intentionally mirror NetBox's writable fields.
type WireFHRPGroupWrite struct {
	Name         string         `json:"name,omitempty"`
	Protocol     *string        `json:"protocol,omitempty"`
	GroupID      *int           `json:"group_id,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireFHRPGroupAssignmentWrite is the wire request model for IPAM FHRP group
// assignment create/update. Its json tags intentionally mirror NetBox's
// writable fields.
type WireFHRPGroupAssignmentWrite struct {
	Group         int    `json:"group,omitempty"`
	InterfaceType string `json:"interface_type,omitempty"`
	InterfaceID   int    `json:"interface_id,omitempty"`
	Priority      *int   `json:"priority,omitempty"`
}

// WireVLANTranslationPolicyWrite is the wire request model for IPAM VLAN
// translation policy create/update. Its json tags intentionally mirror NetBox's
// writable fields.
type WireVLANTranslationPolicyWrite struct {
	Name         string         `json:"name,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// WireVLANTranslationRuleWrite is the wire request model for IPAM VLAN
// translation rule create/update. Its json tags intentionally mirror NetBox's
// writable fields.
type WireVLANTranslationRuleWrite struct {
	Policy      int     `json:"policy,omitempty"`
	LocalVID    *int    `json:"local_vid,omitempty"`
	RemoteVID   *int    `json:"remote_vid,omitempty"`
	Description *string `json:"description,omitempty"`
}

// Wire types mirror the NetBox JSON API response objects.
type (
	// WireSite mirrors the corresponding NetBox JSON API response object.
	WireSite struct {
		ID              int            `json:"id"`
		URL             string         `json:"url"`
		Name            string         `json:"name"`
		Slug            string         `json:"slug"`
		Display         string         `json:"display,omitempty"`
		Status          *WireLabel     `json:"status"`
		Region          *WireNested    `json:"region"`
		Tenant          *WireNested    `json:"tenant"`
		Facility        string         `json:"facility,omitempty"`
		TimeZone        string         `json:"time_zone,omitempty"`
		Description     string         `json:"description,omitempty"`
		PhysicalAddress string         `json:"physical_address,omitempty"`
		ShippingAddress string         `json:"shipping_address,omitempty"`
		Comments        string         `json:"comments,omitempty"`
		Tags            []WireTag      `json:"tags,omitempty"`
		CustomFields    map[string]any `json:"custom_fields,omitempty"`
		Created         string         `json:"created"`
		LastUpdated     string         `json:"last_updated"`
	}

	// WireRegion mirrors the corresponding NetBox JSON API response object.
	WireRegion struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		Parent       *WireNested    `json:"parent"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireTenantGroup mirrors the corresponding NetBox JSON API response object.
	WireTenantGroup struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		Parent       *WireNested    `json:"parent"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireDevice mirrors the corresponding NetBox JSON API response object.
	WireDevice struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		DeviceType   *WireNested    `json:"device_type"`
		Role         *WireNested    `json:"role"`
		Tenant       *WireNested    `json:"tenant"`
		Platform     *WireNested    `json:"platform"`
		Serial       string         `json:"serial,omitempty"`
		AssetTag     string         `json:"asset_tag,omitempty"`
		Site         *WireNested    `json:"site"`
		Rack         *WireNested    `json:"rack"`
		Position     float64        `json:"position,omitempty"`
		Face         *WireLabel     `json:"face"`
		Status       *WireLabel     `json:"status"`
		Cluster      *WireNested    `json:"cluster"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireIPAddress mirrors the corresponding NetBox JSON API response object.
	WireIPAddress struct {
		ID                 int              `json:"id"`
		URL                string           `json:"url"`
		Address            string           `json:"address"`
		Display            string           `json:"display,omitempty"`
		VRF                *WireNested      `json:"vrf"`
		Tenant             *WireNested      `json:"tenant"`
		Status             *WireLabel       `json:"status"`
		Role               *WireLabel       `json:"role"`
		AssignedObjectType string           `json:"assigned_object_type,omitempty"`
		AssignedObjectID   int              `json:"assigned_object_id,omitempty"`
		AssignedObject     *WireAssignedObj `json:"assigned_object,omitempty"`
		DNSName            string           `json:"dns_name,omitempty"`
		Description        string           `json:"description,omitempty"`
		Comments           string           `json:"comments,omitempty"`
		Tags               []WireTag        `json:"tags,omitempty"`
		CustomFields       map[string]any   `json:"custom_fields,omitempty"`
		Created            string           `json:"created"`
		LastUpdated        string           `json:"last_updated"`
	}

	// WirePrefix mirrors the corresponding NetBox JSON API response object.
	WirePrefix struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Prefix       string         `json:"prefix"`
		Display      string         `json:"display,omitempty"`
		Site         *WireNested    `json:"site"`
		VRF          *WireNested    `json:"vrf"`
		Tenant       *WireNested    `json:"tenant"`
		VLAN         *WireNested    `json:"vlan"`
		Status       *WireLabel     `json:"status"`
		Role         *WireNested    `json:"role"`
		IsPool       bool           `json:"is_pool"`
		Description  string         `json:"description,omitempty"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
		Children     int            `json:"children"`
		Depth        int            `json:"_depth"`
		Family       *WireFamily    `json:"family"`
	}

	// WireVLAN mirrors the corresponding NetBox JSON API response object.
	WireVLAN struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Site         *WireNested    `json:"site"`
		Group        *WireNested    `json:"group"`
		VID          int            `json:"vid"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Tenant       *WireNested    `json:"tenant"`
		Status       *WireLabel     `json:"status"`
		Role         *WireNested    `json:"role"`
		Description  string         `json:"description,omitempty"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireVirtualMachine mirrors the corresponding NetBox JSON API response object.
	WireVirtualMachine struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Cluster      *WireNested    `json:"cluster"`
		Role         *WireNested    `json:"role"`
		Tenant       *WireNested    `json:"tenant"`
		Platform     *WireNested    `json:"platform"`
		Status       *WireLabel     `json:"status"`
		Site         *WireNested    `json:"site"`
		VCPUs        float64        `json:"vcpus,omitempty"`
		Memory       int            `json:"memory,omitempty"`
		Disk         int            `json:"disk,omitempty"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireCluster mirrors the corresponding NetBox JSON API response object.
	WireCluster struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		ClusterType  *WireNested    `json:"type"`
		ClusterGroup *WireNested    `json:"group"`
		Site         *WireNested    `json:"site"`
		Tenant       *WireNested    `json:"tenant"`
		Description  string         `json:"description,omitempty"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireCircuit mirrors the corresponding NetBox JSON API response object.
	WireCircuit struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		CID          string         `json:"cid"`
		Display      string         `json:"display,omitempty"`
		Provider     *WireNested    `json:"provider"`
		CircuitType  *WireNested    `json:"type"`
		Tenant       *WireNested    `json:"tenant"`
		Status       *WireLabel     `json:"status"`
		Description  string         `json:"description,omitempty"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
		InstallDate  string         `json:"install_date,omitempty"`
		CommitRate   int            `json:"commit_rate,omitempty"`
	}

	// WireRack mirrors the corresponding NetBox JSON API response object.
	WireRack struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		FacilityID   string         `json:"facility_id,omitempty"`
		Site         *WireNested    `json:"site"`
		Location     *WireNested    `json:"location"`
		Tenant       *WireNested    `json:"tenant"`
		Status       *WireLabel     `json:"status"`
		Role         *WireNested    `json:"role"`
		Serial       string         `json:"serial,omitempty"`
		AssetTag     string         `json:"asset_tag,omitempty"`
		Type         *WireLabel     `json:"type"`
		Width        int            `json:"width,omitempty"`
		UHeight      int            `json:"u_height,omitempty"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireCircuitTermination mirrors the corresponding NetBox JSON API response object.
	WireCircuitTermination struct {
		ID            int            `json:"id"`
		URL           string         `json:"url"`
		Display       string         `json:"display,omitempty"`
		Circuit       *WireNested    `json:"circuit"`
		TermSide      string         `json:"term_side"`
		Site          *WireNested    `json:"site"`
		Speed         int            `json:"speed,omitempty"`
		UpstreamSpeed int            `json:"upstream_speed,omitempty"`
		Description   string         `json:"description,omitempty"`
		Tags          []WireTag      `json:"tags,omitempty"`
		CustomFields  map[string]any `json:"custom_fields,omitempty"`
		Created       string         `json:"created"`
		LastUpdated   string         `json:"last_updated"`
	}

	// WireCable mirrors the corresponding NetBox JSON API response object.
	WireCable struct {
		ID           int                   `json:"id"`
		URL          string                `json:"url"`
		Display      string                `json:"display,omitempty"`
		Type         *WireLabel            `json:"type"`
		Status       *WireLabel            `json:"status"`
		Label        string                `json:"label,omitempty"`
		Color        string                `json:"color,omitempty"`
		Length       float64               `json:"length,omitempty"`
		LengthUnit   *WireLabel            `json:"length_unit,omitempty"`
		TerminationA *WireCableTermination `json:"termination_a,omitempty"`
		TerminationB *WireCableTermination `json:"termination_b,omitempty"`
		Description  string                `json:"description,omitempty"`
		Tags         []WireTag             `json:"tags,omitempty"`
		CustomFields map[string]any        `json:"custom_fields,omitempty"`
		Created      string                `json:"created"`
		LastUpdated  string                `json:"last_updated"`
	}

	// WireCableTermination mirrors the corresponding NetBox JSON API response object.
	WireCableTermination struct {
		ID   int    `json:"id"`
		URL  string `json:"url"`
		Name string `json:"name"`
		Type string `json:"_type"`
	}

	// WireInterface mirrors the corresponding NetBox JSON API response object.
	WireInterface struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Device       *WireNested    `json:"device"`
		Type         *WireLabel     `json:"type"`
		Enabled      bool           `json:"enabled"`
		MTU          int            `json:"mtu,omitempty"`
		MACAddress   string         `json:"mac_address,omitempty"`
		Speed        int            `json:"speed,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireVMInterface mirrors the corresponding NetBox JSON API response object.
	WireVMInterface struct {
		ID             int            `json:"id"`
		URL            string         `json:"url"`
		Name           string         `json:"name"`
		Display        string         `json:"display,omitempty"`
		VirtualMachine *WireNested    `json:"virtual_machine"`
		Enabled        bool           `json:"enabled"`
		MTU            int            `json:"mtu,omitempty"`
		MACAddress     string         `json:"mac_address,omitempty"`
		Description    string         `json:"description,omitempty"`
		Tags           []WireTag      `json:"tags,omitempty"`
		CustomFields   map[string]any `json:"custom_fields,omitempty"`
		Created        string         `json:"created"`
		LastUpdated    string         `json:"last_updated"`
	}

	// WireProvider mirrors the corresponding NetBox JSON API response object.
	WireProvider struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug,omitempty"`
		Display      string         `json:"display,omitempty"`
		Asn          int            `json:"asn,omitempty"`
		Account      string         `json:"account,omitempty"`
		PortalURL    string         `json:"portal_url,omitempty"`
		NocContact   string         `json:"noc_contact,omitempty"`
		AdminContact string         `json:"admin_contact,omitempty"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireTenant mirrors the corresponding NetBox JSON API response object.
	WireTenant struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug,omitempty"`
		Display      string         `json:"display,omitempty"`
		Group        *WireNested    `json:"group"`
		Description  string         `json:"description,omitempty"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireManufacturer mirrors the corresponding NetBox JSON API response object.
	WireManufacturer struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug,omitempty"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireDeviceType mirrors the corresponding NetBox JSON API response object.
	WireDeviceType struct {
		ID            int            `json:"id"`
		URL           string         `json:"url"`
		Manufacturer  *WireNested    `json:"manufacturer"`
		Model         string         `json:"model"`
		Slug          string         `json:"slug,omitempty"`
		Display       string         `json:"display,omitempty"`
		PartNumber    string         `json:"part_number,omitempty"`
		UHeight       float64        `json:"u_height,omitempty"`
		IsFullDepth   bool           `json:"is_full_depth,omitempty"`
		SubdeviceRole *WireLabel     `json:"subdevice_role,omitempty"`
		Comments      string         `json:"comments,omitempty"`
		Tags          []WireTag      `json:"tags,omitempty"`
		CustomFields  map[string]any `json:"custom_fields,omitempty"`
		Created       string         `json:"created"`
		LastUpdated   string         `json:"last_updated"`
	}

	// WireLocation mirrors the corresponding NetBox JSON API response object.
	WireLocation struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug,omitempty"`
		Display      string         `json:"display,omitempty"`
		Site         *WireNested    `json:"site"`
		Parent       *WireNested    `json:"parent"`
		Status       *WireLabel     `json:"status"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireClusterType mirrors the corresponding NetBox JSON API response object.
	WireClusterType struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug,omitempty"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireClusterGroup mirrors the corresponding NetBox JSON API response object.
	WireClusterGroup struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug,omitempty"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireCircuitType mirrors the corresponding NetBox JSON API response object.
	WireCircuitType struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug,omitempty"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		CircuitCount int            `json:"circuit_count"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireVrf mirrors the corresponding NetBox JSON API response object.
	WireVrf struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Rd           string         `json:"rd,omitempty"`
		Tenant       *WireNested    `json:"tenant"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireVlanGroup mirrors the corresponding NetBox JSON API response object.
	WireVlanGroup struct {
		ID           int                 `json:"id"`
		URL          string              `json:"url"`
		Name         string              `json:"name"`
		Slug         string              `json:"slug,omitempty"`
		Display      string              `json:"display,omitempty"`
		Description  string              `json:"description,omitempty"`
		Scope        *WireVlanGroupScope `json:"scope"`
		MinVID       *int                `json:"min_vid"`
		MaxVID       *int                `json:"max_vid"`
		Tags         []WireTag           `json:"tags,omitempty"`
		CustomFields map[string]any      `json:"custom_fields,omitempty"`
		Created      string              `json:"created"`
		LastUpdated  string              `json:"last_updated"`
	}

	// WireVlanGroupScope mirrors the corresponding NetBox JSON API response object.
	WireVlanGroupScope struct {
		ObjectType string `json:"object_type"`
		ObjectID   int    `json:"object_id"`
	}

	// WireRole mirrors the corresponding NetBox JSON API response object.
	WireRole struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug,omitempty"`
		Weight       int            `json:"weight"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireContact mirrors the corresponding NetBox JSON API response object.
	WireContact struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Group        *WireNested    `json:"group"`
		Title        string         `json:"title,omitempty"`
		Phone        string         `json:"phone,omitempty"`
		Email        string         `json:"email,omitempty"`
		Address      string         `json:"address,omitempty"`
		Link         string         `json:"link,omitempty"`
		Description  string         `json:"description,omitempty"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireContactGroup mirrors the corresponding NetBox JSON API response object.
	WireContactGroup struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug"`
		Display      string         `json:"display,omitempty"`
		Parent       *WireNested    `json:"parent"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireContactRole mirrors the corresponding NetBox JSON API response object.
	WireContactRole struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireVirtualDisk mirrors the corresponding NetBox JSON API response object.
	WireVirtualDisk struct {
		ID             int            `json:"id"`
		URL            string         `json:"url"`
		Name           string         `json:"name"`
		Display        string         `json:"display,omitempty"`
		VirtualMachine *WireNested    `json:"virtual_machine"`
		Size           int            `json:"size"`
		Description    string         `json:"description,omitempty"`
		Tags           []WireTag      `json:"tags,omitempty"`
		CustomFields   map[string]any `json:"custom_fields,omitempty"`
		Created        string         `json:"created"`
		LastUpdated    string         `json:"last_updated"`
	}

	// WireVirtualMachineType mirrors the corresponding NetBox JSON API response
	// object.
	WireVirtualMachineType struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireSiteGroup mirrors the corresponding NetBox JSON API response object.
	WireSiteGroup struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug"`
		Display      string         `json:"display,omitempty"`
		Parent       *WireNested    `json:"parent"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireRackGroup mirrors the corresponding NetBox JSON API response object.
	WireRackGroup struct {
		ID          int         `json:"id"`
		URL         string      `json:"url"`
		Name        string      `json:"name"`
		Slug        string      `json:"slug"`
		Display     string      `json:"display,omitempty"`
		Site        *WireNested `json:"site"`
		Parent      *WireNested `json:"parent"`
		Description string      `json:"description,omitempty"`
		Tags        []WireTag   `json:"tags,omitempty"`
		Created     string      `json:"created"`
		LastUpdated string      `json:"last_updated"`
	}

	// WireRackRole mirrors the corresponding NetBox JSON API response object.
	WireRackRole struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug"`
		Color        string         `json:"color,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireDeviceRole mirrors the corresponding NetBox JSON API response object.
	WireDeviceRole struct {
		ID             int            `json:"id"`
		URL            string         `json:"url"`
		Name           string         `json:"name"`
		Slug           string         `json:"slug"`
		Color          string         `json:"color,omitempty"`
		VMRole         bool           `json:"vm_role"`
		ConfigTemplate *WireNested    `json:"config_template"`
		Description    string         `json:"description,omitempty"`
		Tags           []WireTag      `json:"tags,omitempty"`
		CustomFields   map[string]any `json:"custom_fields,omitempty"`
		Created        string         `json:"created"`
		LastUpdated    string         `json:"last_updated"`
		Weight         int            `json:"weight"`
	}

	// WirePlatform mirrors the corresponding NetBox JSON API response object.
	WirePlatform struct {
		ID             int            `json:"id"`
		URL            string         `json:"url"`
		Name           string         `json:"name"`
		Slug           string         `json:"slug"`
		Manufacturer   *WireNested    `json:"manufacturer"`
		ConfigTemplate *WireNested    `json:"config_template"`
		Description    string         `json:"description,omitempty"`
		Tags           []WireTag      `json:"tags,omitempty"`
		CustomFields   map[string]any `json:"custom_fields,omitempty"`
		Created        string         `json:"created"`
		LastUpdated    string         `json:"last_updated"`
	}

	// WireVirtualChassis mirrors the corresponding NetBox JSON API response
	// object.
	WireVirtualChassis struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Domain       string         `json:"domain,omitempty"`
		Master       *WireNested    `json:"master"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
		MemberCount  int            `json:"member_count"`
	}

	// WirePowerPanel mirrors the corresponding NetBox JSON API response object.
	WirePowerPanel struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Site         *WireNested    `json:"site"`
		Location     *WireNested    `json:"location"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireInventoryItemRole mirrors the corresponding NetBox JSON API response
	// object.
	WireInventoryItemRole struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug"`
		Color        string         `json:"color,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireFrontPort mirrors the corresponding NetBox JSON API response object.
	WireFrontPort struct {
		ID               int            `json:"id"`
		URL              string         `json:"url"`
		Name             string         `json:"name"`
		Display          string         `json:"display,omitempty"`
		Device           *WireNested    `json:"device"`
		Type             *WireLabel     `json:"type"`
		RearPort         *WireNested    `json:"rear_port"`
		RearPortPosition int            `json:"rear_port_position,omitempty"`
		Description      string         `json:"description,omitempty"`
		Tags             []WireTag      `json:"tags,omitempty"`
		CustomFields     map[string]any `json:"custom_fields,omitempty"`
		Created          string         `json:"created"`
		LastUpdated      string         `json:"last_updated"`
	}

	// WireRearPort mirrors the corresponding NetBox JSON API response object.
	WireRearPort struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Device       *WireNested    `json:"device"`
		Type         *WireLabel     `json:"type"`
		Positions    int            `json:"positions"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireConsolePort mirrors the corresponding NetBox JSON API response object.
	WireConsolePort struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Device       *WireNested    `json:"device"`
		Type         *WireLabel     `json:"type"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireConsoleServerPort mirrors the corresponding NetBox JSON API response
	// object.
	WireConsoleServerPort struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Device       *WireNested    `json:"device"`
		Type         *WireLabel     `json:"type"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WirePowerPort mirrors the corresponding NetBox JSON API response object.
	WirePowerPort struct {
		ID            int            `json:"id"`
		URL           string         `json:"url"`
		Name          string         `json:"name"`
		Display       string         `json:"display,omitempty"`
		Device        *WireNested    `json:"device"`
		Type          *WireLabel     `json:"type"`
		MaximumDraw   int            `json:"maximum_draw,omitempty"`
		AllocatedDraw int            `json:"allocated_draw,omitempty"`
		Description   string         `json:"description,omitempty"`
		Tags          []WireTag      `json:"tags,omitempty"`
		CustomFields  map[string]any `json:"custom_fields,omitempty"`
		Created       string         `json:"created"`
		LastUpdated   string         `json:"last_updated"`
	}

	// WirePowerOutlet mirrors the corresponding NetBox JSON API response object.
	WirePowerOutlet struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Device       *WireNested    `json:"device"`
		Type         *WireLabel     `json:"type"`
		PowerPort    *WireNested    `json:"power_port"`
		FeedLeg      *WireLabel     `json:"feed_leg"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WirePowerFeed mirrors the corresponding NetBox JSON API response object.
	WirePowerFeed struct {
		ID             int            `json:"id"`
		URL            string         `json:"url"`
		Name           string         `json:"name"`
		Display        string         `json:"display,omitempty"`
		PowerPanel     *WireNested    `json:"power_panel"`
		Status         *WireLabel     `json:"status"`
		Type           *WireLabel     `json:"type"`
		Supply         *WireLabel     `json:"supply"`
		Phase          *WireLabel     `json:"phase"`
		Voltage        int            `json:"voltage"`
		Amperage       int            `json:"amperage"`
		MaxUtilization int            `json:"max_utilization,omitempty"`
		Description    string         `json:"description,omitempty"`
		Tags           []WireTag      `json:"tags,omitempty"`
		CustomFields   map[string]any `json:"custom_fields,omitempty"`
		Created        string         `json:"created"`
		LastUpdated    string         `json:"last_updated"`
	}

	// WireDeviceBay mirrors the corresponding NetBox JSON API response object.
	WireDeviceBay struct {
		ID              int            `json:"id"`
		URL             string         `json:"url"`
		Name            string         `json:"name"`
		Display         string         `json:"display,omitempty"`
		Device          *WireNested    `json:"device"`
		InstalledDevice *WireNested    `json:"installed_device"`
		Description     string         `json:"description,omitempty"`
		Tags            []WireTag      `json:"tags,omitempty"`
		CustomFields    map[string]any `json:"custom_fields,omitempty"`
		Created         string         `json:"created"`
		LastUpdated     string         `json:"last_updated"`
	}

	// WireRackReservation mirrors the corresponding NetBox JSON API response
	// object.
	WireRackReservation struct {
		ID          int         `json:"id"`
		URL         string      `json:"url"`
		Rack        *WireNested `json:"rack"`
		Units       []int       `json:"units,omitempty"`
		Description string      `json:"description"`
		User        *WireNested `json:"user"`
		Created     string      `json:"created"`
		LastUpdated string      `json:"last_updated"`
	}

	// WireInventoryItem mirrors the corresponding NetBox JSON API response
	// object.
	WireInventoryItem struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Device       *WireNested    `json:"device"`
		Role         *WireNested    `json:"role"`
		Manufacturer *WireNested    `json:"manufacturer"`
		Serial       string         `json:"serial,omitempty"`
		PartID       string         `json:"part_id,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireMACAddress mirrors the corresponding NetBox JSON API response object.
	WireMACAddress struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Address      string         `json:"mac_address"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireRIR mirrors the corresponding NetBox JSON API response object.
	WireRIR struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug"`
		Display      string         `json:"display,omitempty"`
		IsPrivate    bool           `json:"is_private"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireAggregate mirrors the corresponding NetBox JSON API response object.
	WireAggregate struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Prefix       string         `json:"prefix"`
		Display      string         `json:"display,omitempty"`
		RIR          *WireNested    `json:"rir"`
		Tenant       *WireNested    `json:"tenant"`
		DateAdded    string         `json:"date_added"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireIPRange mirrors the corresponding NetBox JSON API response object.
	WireIPRange struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		StartAddress string         `json:"start_address"`
		EndAddress   string         `json:"end_address"`
		Display      string         `json:"display,omitempty"`
		VRF          *WireNested    `json:"vrf"`
		Tenant       *WireNested    `json:"tenant"`
		Status       *WireLabel     `json:"status"`
		Role         *WireNested    `json:"role"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireASN mirrors the corresponding NetBox JSON API response object.
	WireASN struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		ASN          int            `json:"asn"`
		Display      string         `json:"display,omitempty"`
		RIR          *WireNested    `json:"rir"`
		Tenant       *WireNested    `json:"tenant"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireASNRange mirrors the corresponding NetBox JSON API response object.
	WireASNRange struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Start        int            `json:"start"`
		End          int            `json:"end"`
		RIR          *WireNested    `json:"rir"`
		Tenant       *WireNested    `json:"tenant"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireRouteTarget mirrors the corresponding NetBox JSON API response object.
	WireRouteTarget struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Tenant       *WireNested    `json:"tenant"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireService mirrors the corresponding NetBox JSON API response object.
	WireService struct {
		ID               int            `json:"id"`
		URL              string         `json:"url"`
		Name             string         `json:"name"`
		Display          string         `json:"display,omitempty"`
		ParentObjectType string         `json:"parent_object_type"`
		ParentObjectID   int            `json:"parent_object_id"`
		Parent           *WireNested    `json:"parent"`
		Protocol         *WireLabel     `json:"protocol"`
		Ports            []int          `json:"ports,omitempty"`
		Description      string         `json:"description,omitempty"`
		Tags             []WireTag      `json:"tags,omitempty"`
		CustomFields     map[string]any `json:"custom_fields,omitempty"`
		Created          string         `json:"created"`
		LastUpdated      string         `json:"last_updated"`
	}

	// WireServiceTemplate mirrors the corresponding NetBox JSON API response
	// object.
	WireServiceTemplate struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Protocol     *WireLabel     `json:"protocol"`
		Ports        []int          `json:"ports,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireFHRPGroup mirrors the corresponding NetBox JSON API response object.
	WireFHRPGroup struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Protocol     string         `json:"protocol"`
		GroupID      int            `json:"group_id"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireFHRPGroupAssignment mirrors the corresponding NetBox JSON API
	// response object.
	WireFHRPGroupAssignment struct {
		ID            int         `json:"id"`
		URL           string      `json:"url"`
		Group         *WireNested `json:"group"`
		InterfaceType string      `json:"interface_type"`
		InterfaceID   int         `json:"interface_id"`
		Interface     *WireNested `json:"interface"`
		Priority      int         `json:"priority"`
		Created       string      `json:"created"`
		LastUpdated   string      `json:"last_updated"`
	}

	// WireVLANTranslationPolicy mirrors the corresponding NetBox JSON API
	// response object.
	WireVLANTranslationPolicy struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// WireVLANTranslationRule mirrors the corresponding NetBox JSON API response
	// object.
	WireVLANTranslationRule struct {
		ID          int    `json:"id"`
		URL         string `json:"url"`
		Display     string `json:"display,omitempty"`
		Policy      int    `json:"policy"`
		LocalVID    int    `json:"local_vid"`
		RemoteVID   int    `json:"remote_vid"`
		Description string `json:"description,omitempty"`
		Created     string `json:"created"`
		LastUpdated string `json:"last_updated"`
	}
)

// Wire types for the small nested, label and tag sub-objects embedded in the
// response objects above.
type (
	// WireNested mirrors the corresponding NetBox JSON API response object.
	WireNested struct {
		ID   int    `json:"id"`
		URL  string `json:"url"`
		Name string `json:"name"`
		Slug string `json:"slug,omitempty"`
	}

	// WireLabel mirrors the corresponding NetBox JSON API response object.
	WireLabel struct {
		Value string `json:"value"`
		Label string `json:"label"`
	}

	// WireTag mirrors the corresponding NetBox JSON API response object.
	WireTag struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
		URL  string `json:"url"`
	}

	// WireAssignedObj mirrors the corresponding NetBox JSON API response object.
	WireAssignedObj struct {
		ID     int         `json:"id"`
		URL    string      `json:"url"`
		Device *WireNested `json:"device"`
		Name   string      `json:"name"`
	}

	// WireFamily mirrors the corresponding NetBox JSON API response object.
	WireFamily struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}
)
