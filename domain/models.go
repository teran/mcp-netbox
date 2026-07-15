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
		ID            int               `json:"id"`
		URL           string            `json:"url,omitempty"`
		Display       string            `json:"display,omitempty"`
		Type          *Label            `json:"type"`
		Status        *Label            `json:"status"`
		Label         string            `json:"label,omitempty"`
		Color         string            `json:"color,omitempty"`
		Length        float64           `json:"length,omitempty"`
		LengthUnit    *Label            `json:"length_unit,omitempty"`
		TerminationA  *CableTermination `json:"termination_a,omitempty"`
		TerminationB  *CableTermination `json:"termination_b,omitempty"`
		Description   string            `json:"description,omitempty"`
		Tags          []Tag             `json:"tags,omitempty"`
		CustomFields  map[string]any    `json:"custom_fields,omitempty"`
		Created       string            `json:"created"`
		LastUpdated   string            `json:"last_updated"`
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

// PaginatedResponse wraps a paginated API response with the total count,
// navigation URLs and the current page's typed results.
type PaginatedResponse[T any] struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []T    `json:"results"`
}
