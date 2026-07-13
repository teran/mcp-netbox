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
		ID              int     `json:"id"`
		Name            string  `json:"name"`
		Slug            string  `json:"slug"`
		Status          *Label  `json:"status"`
		Region          *Nested `json:"region"`
		Tenant          *Nested `json:"tenant"`
		Facility        string  `json:"facility,omitempty"`
		TimeZone        string  `json:"time_zone,omitempty"`
		Description     string  `json:"description,omitempty"`
		PhysicalAddress string  `json:"physical_address,omitempty"`
		ShippingAddress string  `json:"shipping_address,omitempty"`
		Comments        string  `json:"comments,omitempty"`
		Tags            []Tag   `json:"tags,omitempty"`
		Created         string  `json:"created"`
		LastUpdated     string  `json:"last_updated"`
	}

	Device struct {
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		DeviceType  *Nested `json:"device_type"`
		Role        *Nested `json:"role"`
		Tenant      *Nested `json:"tenant"`
		Platform    *Nested `json:"platform"`
		Serial      string  `json:"serial,omitempty"`
		AssetTag    string  `json:"asset_tag,omitempty"`
		Site        *Nested `json:"site"`
		Rack        *Nested `json:"rack"`
		Position    float64 `json:"position,omitempty"`
		Face        *Label  `json:"face"`
		Status      *Label  `json:"status"`
		Cluster     *Nested `json:"cluster"`
		Comments    string  `json:"comments,omitempty"`
		Tags        []Tag   `json:"tags,omitempty"`
		Created     string  `json:"created"`
		LastUpdated string  `json:"last_updated"`
	}

	IPAddress struct {
		ID                 int             `json:"id"`
		Address            string          `json:"address"`
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
		Created            string          `json:"created"`
		LastUpdated        string          `json:"last_updated"`
	}

	Prefix struct {
		ID          int     `json:"id"`
		Prefix      string  `json:"prefix"`
		Site        *Nested `json:"site"`
		VRF         *Nested `json:"vrf"`
		Tenant      *Nested `json:"tenant"`
		VLAN        *Nested `json:"vlan"`
		Status      *Label  `json:"status"`
		Role        *Nested `json:"role"`
		IsPool      bool    `json:"is_pool"`
		Description string  `json:"description,omitempty"`
		Comments    string  `json:"comments,omitempty"`
		Tags        []Tag   `json:"tags,omitempty"`
		Created     string  `json:"created"`
		LastUpdated string  `json:"last_updated"`
		Children    int     `json:"children"`
		Depth       int     `json:"_depth"`
		Family      *Family `json:"family"`
	}

	VLAN struct {
		ID          int     `json:"id"`
		Site        *Nested `json:"site"`
		Group       *Nested `json:"group"`
		VID         int     `json:"vid"`
		Name        string  `json:"name"`
		Tenant      *Nested `json:"tenant"`
		Status      *Label  `json:"status"`
		Role        *Nested `json:"role"`
		Description string  `json:"description,omitempty"`
		Comments    string  `json:"comments,omitempty"`
		Tags        []Tag   `json:"tags,omitempty"`
		Created     string  `json:"created"`
		LastUpdated string  `json:"last_updated"`
	}

	VirtualMachine struct {
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		Cluster     *Nested `json:"cluster"`
		Role        *Nested `json:"role"`
		Tenant      *Nested `json:"tenant"`
		Platform    *Nested `json:"platform"`
		Status      *Label  `json:"status"`
		Site        *Nested `json:"site"`
		VCPUs       float64 `json:"vcpus,omitempty"`
		Memory      int     `json:"memory,omitempty"`
		Disk        int     `json:"disk,omitempty"`
		Comments    string  `json:"comments,omitempty"`
		Tags        []Tag   `json:"tags,omitempty"`
		Created     string  `json:"created"`
		LastUpdated string  `json:"last_updated"`
	}

	Cluster struct {
		ID           int     `json:"id"`
		Name         string  `json:"name"`
		ClusterType  *Nested `json:"type"`
		ClusterGroup *Nested `json:"group"`
		Site         *Nested `json:"site"`
		Tenant       *Nested `json:"tenant"`
		Description  string  `json:"description,omitempty"`
		Comments     string  `json:"comments,omitempty"`
		Tags         []Tag   `json:"tags,omitempty"`
		Created      string  `json:"created"`
		LastUpdated  string  `json:"last_updated"`
	}

	Circuit struct {
		ID          int     `json:"id"`
		CID         string  `json:"cid"`
		Provider    *Nested `json:"provider"`
		CircuitType *Nested `json:"circuit_type"`
		Tenant      *Nested `json:"tenant"`
		Status      *Label  `json:"status"`
		Description string  `json:"description,omitempty"`
		Comments    string  `json:"comments,omitempty"`
		Tags        []Tag   `json:"tags,omitempty"`
		Created     string  `json:"created"`
		LastUpdated string  `json:"last_updated"`
		InstallDate string  `json:"install_date,omitempty"`
		CommitRate  int     `json:"commit_rate,omitempty"`
	}

	Rack struct {
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		FacilityID  string  `json:"facility_id,omitempty"`
		Site        *Nested `json:"site"`
		Location    *Nested `json:"location"`
		Tenant      *Nested `json:"tenant"`
		Status      *Label  `json:"status"`
		Role        *Nested `json:"role"`
		Serial      string  `json:"serial,omitempty"`
		AssetTag    string  `json:"asset_tag,omitempty"`
		Type        *Label  `json:"type"`
		Width       int     `json:"width,omitempty"`
		UHeight     int     `json:"u_height,omitempty"`
		Comments    string  `json:"comments,omitempty"`
		Tags        []Tag   `json:"tags,omitempty"`
		Created     string  `json:"created"`
		LastUpdated string  `json:"last_updated"`
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
