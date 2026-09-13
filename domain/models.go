// Package domain provides the core domain models and repository interfaces
// for the NetBox MCP server.
//
// These types are free of external dependencies and represent the pure
// business objects of the application. JSON struct tags are present for
// serialization when returned as MCP tool output — encoding/json is a
// Go standard library type and does not compromise domain purity.
package domain

// Site represents a NetBox site (physical location or data center).
type (
	Site struct {
		ID              int            `json:"id"`
		URL             string         `json:"url,omitempty"`
		Name            string         `json:"name"`
		Slug            string         `json:"slug"`
		Display         string         `json:"display,omitempty"`
		Status          *Label         `json:"status"`
		Region          *Nested        `json:"region"`
		Tenant          *Nested        `json:"tenant"`
		Facility        string         `json:"facility,omitempty"`
		TimeZone        string         `json:"time_zone,omitempty"`
		Description     string         `json:"description,omitempty"`
		PhysicalAddress string         `json:"physical_address,omitempty"`
		ShippingAddress string         `json:"shipping_address,omitempty"`
		Comments        string         `json:"comments,omitempty"`
		Tags            []Tag          `json:"tags,omitempty"`
		CustomFields    map[string]any `json:"custom_fields,omitempty"`
		Created         string         `json:"created"`
		LastUpdated     string         `json:"last_updated"`
	}

	Device struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		DeviceType   *Nested        `json:"device_type"`
		Role         *Nested        `json:"role"`
		Tenant       *Nested        `json:"tenant"`
		Platform     *Nested        `json:"platform"`
		Serial       string         `json:"serial,omitempty"`
		AssetTag     string         `json:"asset_tag,omitempty"`
		Site         *Nested        `json:"site"`
		Rack         *Nested        `json:"rack"`
		Position     float64        `json:"position,omitempty"`
		Face         *Label         `json:"face"`
		Status       *Label         `json:"status"`
		Cluster      *Nested        `json:"cluster"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	IPAddress struct {
		ID                 int             `json:"id"`
		URL                string          `json:"url,omitempty"`
		Address            string          `json:"address"`
		Display            string          `json:"display,omitempty"`
		VRF                *Nested         `json:"vrf"`
		Tenant             *Nested         `json:"tenant"`
		Status             *Label          `json:"status"`
		Role               *Label          `json:"role"`
		AssignedObjectType string          `json:"assigned_object_type,omitempty"`
		AssignedObjectID   int             `json:"assigned_object_id,omitempty"`
		AssignedObject     *AssignedObject `json:"assigned_object,omitempty"`
		DNSName            string          `json:"dns_name,omitempty"`
		Description        string          `json:"description,omitempty"`
		Comments           string          `json:"comments,omitempty"`
		Tags               []Tag           `json:"tags,omitempty"`
		CustomFields       map[string]any  `json:"custom_fields,omitempty"`
		Created            string          `json:"created"`
		LastUpdated        string          `json:"last_updated"`
	}

	Prefix struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Prefix       string         `json:"prefix"`
		Display      string         `json:"display,omitempty"`
		Site         *Nested        `json:"site"`
		VRF          *Nested        `json:"vrf"`
		Tenant       *Nested        `json:"tenant"`
		VLAN         *Nested        `json:"vlan"`
		Status       *Label         `json:"status"`
		Role         *Nested        `json:"role"`
		IsPool       bool           `json:"is_pool"`
		Description  string         `json:"description,omitempty"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
		Children     int            `json:"children"`
		Depth        int            `json:"_depth"`
		Family       *Family        `json:"family"`
	}

	VLAN struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Site         *Nested        `json:"site"`
		Group        *Nested        `json:"group"`
		VID          int            `json:"vid"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Tenant       *Nested        `json:"tenant"`
		Status       *Label         `json:"status"`
		Role         *Nested        `json:"role"`
		Description  string         `json:"description,omitempty"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	VirtualMachine struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Cluster      *Nested        `json:"cluster"`
		Role         *Nested        `json:"role"`
		Tenant       *Nested        `json:"tenant"`
		Platform     *Nested        `json:"platform"`
		Status       *Label         `json:"status"`
		Site         *Nested        `json:"site"`
		VCPUs        float64        `json:"vcpus,omitempty"`
		Memory       int            `json:"memory,omitempty"`
		Disk         int            `json:"disk,omitempty"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	Cluster struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		ClusterType  *Nested        `json:"type"`
		ClusterGroup *Nested        `json:"group"`
		Site         *Nested        `json:"site"`
		Tenant       *Nested        `json:"tenant"`
		Description  string         `json:"description,omitempty"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	Circuit struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		CID          string         `json:"cid"`
		Display      string         `json:"display,omitempty"`
		Provider     *Nested        `json:"provider"`
		CircuitType  *Nested        `json:"circuit_type"`
		Tenant       *Nested        `json:"tenant"`
		Status       *Label         `json:"status"`
		Description  string         `json:"description,omitempty"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
		InstallDate  string         `json:"install_date,omitempty"`
		CommitRate   int            `json:"commit_rate,omitempty"`
	}

	CircuitTermination struct {
		ID            int            `json:"id"`
		URL           string         `json:"url,omitempty"`
		Display       string         `json:"display,omitempty"`
		Circuit       *Nested        `json:"circuit"`
		TermSide      string         `json:"term_side"`
		Site          *Nested        `json:"site"`
		Speed         int            `json:"speed,omitempty"`
		UpstreamSpeed int            `json:"upstream_speed,omitempty"`
		Description   string         `json:"description,omitempty"`
		Tags          []Tag          `json:"tags,omitempty"`
		CustomFields  map[string]any `json:"custom_fields,omitempty"`
		Created       string         `json:"created"`
		LastUpdated   string         `json:"last_updated"`
	}

	Cable struct {
		ID           int               `json:"id"`
		URL          string            `json:"url,omitempty"`
		Display      string            `json:"display,omitempty"`
		Type         *Label            `json:"type"`
		Status       *Label            `json:"status"`
		Label        string            `json:"label,omitempty"`
		Color        string            `json:"color,omitempty"`
		Length       float64           `json:"length,omitempty"`
		LengthUnit   *Label            `json:"length_unit,omitempty"`
		TerminationA *CableTermination `json:"termination_a,omitempty"`
		TerminationB *CableTermination `json:"termination_b,omitempty"`
		Description  string            `json:"description,omitempty"`
		Tags         []Tag             `json:"tags,omitempty"`
		CustomFields map[string]any    `json:"custom_fields,omitempty"`
		Created      string            `json:"created"`
		LastUpdated  string            `json:"last_updated"`
	}

	// CableTermination represents one end of a cable connection — the type
	// of the terminated object (e.g. "dcim.interface") and its raw data from
	// the NetBox API.
	CableTermination struct {
		ID   int    `json:"id"`
		URL  string `json:"url"`
		Name string `json:"name"`
		Type string `json:"_type"`
	}

	Rack struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		FacilityID   string         `json:"facility_id,omitempty"`
		Site         *Nested        `json:"site"`
		Location     *Nested        `json:"location"`
		Tenant       *Nested        `json:"tenant"`
		Status       *Label         `json:"status"`
		Role         *Nested        `json:"role"`
		Serial       string         `json:"serial,omitempty"`
		AssetTag     string         `json:"asset_tag,omitempty"`
		Type         *Label         `json:"type"`
		Width        int            `json:"width,omitempty"`
		UHeight      int            `json:"u_height,omitempty"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	Interface struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Device       *Nested        `json:"device"`
		Type         *Label         `json:"type"`
		Enabled      bool           `json:"enabled"`
		MTU          int            `json:"mtu,omitempty"`
		MACAddress   string         `json:"mac_address,omitempty"`
		Speed        int            `json:"speed,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	VMInterface struct {
		ID             int            `json:"id"`
		URL            string         `json:"url,omitempty"`
		Name           string         `json:"name"`
		Display        string         `json:"display,omitempty"`
		VirtualMachine *Nested        `json:"virtual_machine"`
		Enabled        bool           `json:"enabled"`
		MTU            int            `json:"mtu,omitempty"`
		MACAddress     string         `json:"mac_address,omitempty"`
		Description    string         `json:"description,omitempty"`
		Tags           []Tag          `json:"tags,omitempty"`
		CustomFields   map[string]any `json:"custom_fields,omitempty"`
		Created        string         `json:"created"`
		LastUpdated    string         `json:"last_updated"`
	}
)

type (
	Nested struct {
		ID   int    `json:"id"`
		URL  string `json:"url"`
		Name string `json:"name"`
		Slug string `json:"slug,omitempty"`
	}

	Label struct {
		Value string `json:"value"`
		Label string `json:"label"`
	}

	Tag struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
		URL  string `json:"url"`
	}

	AssignedObject struct {
		ID     int     `json:"id"`
		URL    string  `json:"url"`
		Device *Nested `json:"device"`
		Name   string  `json:"name"`
	}

	Family struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}
)

// SiteWrite is the domain DTO for creating or updating a NetBox site. It
// carries only the writable fields.
//
// Optional scalar fields are pointers so that a partial-update (PATCH) can
// distinguish "not provided" (nil) from "explicitly set to empty" (pointer to
// ""). Name is a plain string because it is required on create; the omitempty
// tag drops it from a PATCH body when not provided.
type SiteWrite struct {
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

// DeviceWrite is the domain DTO for creating or updating a NetBox device. It
// carries only the writable fields.
//
// Optional scalar fields are pointers so that a partial-update (PATCH) can
// distinguish "not provided" (nil) from "explicitly set to empty" (pointer to
// ""). Name is a plain string because it is required on create; the omitempty
// tag drops it from a PATCH body when not provided. References to related
// objects (device type, role, tenant, etc.) are their numeric NetBox IDs.
type DeviceWrite struct {
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

// IPAddressWrite is the domain DTO for creating or updating a NetBox IP
// address. It carries only the writable fields.
//
// Optional scalar fields are pointers so that a partial-update (PATCH) can
// distinguish "not provided" (nil) from "explicitly set to empty" (pointer to
// ""). Address is a plain string because it is required on create; the
// omitempty tag drops it from a PATCH body when not provided. References to
// related objects (VRF, tenant, assigned object) are their numeric NetBox IDs.
type IPAddressWrite struct {
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

// PrefixWrite is the domain DTO for creating or updating a NetBox prefix. It
// carries only the writable fields.
//
// Optional scalar fields are pointers so that a partial-update (PATCH) can
// distinguish "not provided" (nil) from "explicitly set to empty" (pointer to
// ""). Prefix is a plain string because it is required on create; the omitempty
// tag drops it from a PATCH body when not provided. References to related
// objects (site, VRF, tenant, VLAN, role) are their numeric NetBox IDs. IsPool
// is a pointer so an update can set it to false explicitly.
type PrefixWrite struct {
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

// VLANWrite is the domain DTO for creating or updating a NetBox VLAN. It
// carries only the writable fields.
//
// VID and Name are plain scalars because both are required on create; their
// omitempty tags drop them from a PATCH body when not provided. Other optional
// scalar fields are pointers so that a partial-update (PATCH) can distinguish
// "not provided" (nil) from "explicitly set to empty". References to related
// objects (site, group, tenant, role) are their numeric NetBox IDs.
type VLANWrite struct {
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

// VirtualMachineWrite is the domain DTO for creating or updating a NetBox
// virtual machine. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". References to related objects (cluster,
// role, tenant, platform, site) are their numeric NetBox IDs.
type VirtualMachineWrite struct {
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

// ClusterWrite is the domain DTO for creating or updating a NetBox cluster. It
// carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. ClusterType (type) is required
// on create. Optional scalar fields are pointers so that a partial-update
// (PATCH) can distinguish "not provided" (nil) from "explicitly set to empty".
// References to related objects (type, group, site, tenant) are their numeric
// NetBox IDs.
type ClusterWrite struct {
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

// CircuitWrite is the domain DTO for creating or updating a NetBox circuit. It
// carries only the writable fields.
//
// CID is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Provider and CircuitType are
// required on create. Optional scalar fields are pointers so that a
// partial-update (PATCH) can distinguish "not provided" (nil) from "explicitly
// set to empty". References to related objects (provider, circuit_type, tenant)
// are their numeric NetBox IDs.
type CircuitWrite struct {
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

// RackWrite is the domain DTO for creating or updating a NetBox rack. It
// carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty" (pointer to ""). References to related
// objects (site, location, tenant, role) are their numeric NetBox IDs. Status
// and type are choice-string values.
type RackWrite struct {
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

// InterfaceWrite is the domain DTO for creating or updating a NetBox device
// interface. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Device and Type are required on
// create. Optional scalar fields are pointers so that a partial-update (PATCH)
// can distinguish "not provided" (nil) from "explicitly set to empty". Device
// is the numeric NetBox ID of the parent device; type is a choice-string value.
type InterfaceWrite struct {
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

// CircuitTerminationWrite is the domain DTO for creating or updating a NetBox
// circuit termination. It carries only the writable fields.
//
// TermSide is a plain string because it is required on create; the omitempty
// tag drops it from a PATCH body when not provided. Circuit and Site are
// required on create. Optional scalar fields are pointers so that a
// partial-update (PATCH) can distinguish "not provided" (nil) from "explicitly
// set to empty". Circuit and Site are numeric NetBox IDs.
type CircuitTerminationWrite struct {
	Circuit       *int           `json:"circuit,omitempty"`
	TermSide      string         `json:"term_side,omitempty"`
	Site          *int           `json:"site,omitempty"`
	Speed         *int           `json:"speed,omitempty"`
	UpstreamSpeed *int           `json:"upstream_speed,omitempty"`
	Description   *string        `json:"description,omitempty"`
	Tags          []string       `json:"tags,omitempty"`
	CustomFields  map[string]any `json:"custom_fields,omitempty"`
}

// PaginatedResponse wraps a paginated API response with the total count,
// navigation URLs and the current page's typed results.
type PaginatedResponse[T any] struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []T    `json:"results"`
}
