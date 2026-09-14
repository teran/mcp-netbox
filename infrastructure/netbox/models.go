package netbox

// Wire types mirror NetBox JSON API responses. The domain.PaginatedResponse
// type is used for both deserialization and output — wire types are embedded
// via type parameter instantiation (e.g. domain.PaginatedResponse[WireSite]).
//
// WireSiteWrite is the wire request model for site create/update (POST/PATCH
// body). Its json tags intentionally mirror NetBox's writable fields.
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
	CircuitType  *int           `json:"circuit_type,omitempty"`
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
// mirror NetBox's writable fields. Circuit and Site are numeric IDs; term_side
// is A or Z.
type WireCircuitTerminationWrite struct {
	Circuit       *int           `json:"circuit,omitempty"`
	TermSide      string         `json:"term_side,omitempty"`
	Site          *int           `json:"site,omitempty"`
	Speed         *int           `json:"speed,omitempty"`
	UpstreamSpeed *int           `json:"upstream_speed,omitempty"`
	Description   *string        `json:"description,omitempty"`
	Tags          []string       `json:"tags,omitempty"`
	CustomFields  map[string]any `json:"custom_fields,omitempty"`
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

type (
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

	WireCircuit struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		CID          string         `json:"cid"`
		Display      string         `json:"display,omitempty"`
		Provider     *WireNested    `json:"provider"`
		CircuitType  *WireNested    `json:"circuit_type"`
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

	WireCableTermination struct {
		ID   int    `json:"id"`
		URL  string `json:"url"`
		Name string `json:"name"`
		Type string `json:"_type"`
	}

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

	WireTenant struct {
		ID           int            `json:"id"`
		URL          string         `json:"url"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug,omitempty"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []WireTag      `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

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
)

type (
	WireNested struct {
		ID   int    `json:"id"`
		URL  string `json:"url"`
		Name string `json:"name"`
		Slug string `json:"slug,omitempty"`
	}

	WireLabel struct {
		Value string `json:"value"`
		Label string `json:"label"`
	}

	WireTag struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
		URL  string `json:"url"`
	}

	WireAssignedObj struct {
		ID     int         `json:"id"`
		URL    string      `json:"url"`
		Device *WireNested `json:"device"`
		Name   string      `json:"name"`
	}

	WireFamily struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}
)
