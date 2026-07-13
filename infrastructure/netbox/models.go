package netbox

// Wire types mirror NetBox JSON API responses. The domain.PaginatedResponse
// type is used for both deserialization and output — wire types are embedded
// via type parameter instantiation (e.g. domain.PaginatedResponse[WireSite]).
type (
	WireSite struct {
		ID              int         `json:"id"`
		URL             string      `json:"url"`
		Name            string      `json:"name"`
		Slug            string      `json:"slug"`
		Status          *WireLabel  `json:"status"`
		Region          *WireNested `json:"region"`
		Tenant          *WireNested `json:"tenant"`
		Facility        string      `json:"facility,omitempty"`
		TimeZone        string      `json:"time_zone,omitempty"`
		Description     string      `json:"description,omitempty"`
		PhysicalAddress string      `json:"physical_address,omitempty"`
		ShippingAddress string      `json:"shipping_address,omitempty"`
		Comments        string      `json:"comments,omitempty"`
		Tags            []WireTag   `json:"tags,omitempty"`
		Created         string      `json:"created"`
		LastUpdated     string      `json:"last_updated"`
	}

	WireDevice struct {
		ID          int         `json:"id"`
		URL         string      `json:"url"`
		Name        string      `json:"name"`
		DeviceType  *WireNested `json:"device_type"`
		Role        *WireNested `json:"role"`
		Tenant      *WireNested `json:"tenant"`
		Platform    *WireNested `json:"platform"`
		Serial      string      `json:"serial,omitempty"`
		AssetTag    string      `json:"asset_tag,omitempty"`
		Site        *WireNested `json:"site"`
		Rack        *WireNested `json:"rack"`
		Position    float64     `json:"position,omitempty"`
		Face        *WireLabel  `json:"face"`
		Status      *WireLabel  `json:"status"`
		Cluster     *WireNested `json:"cluster"`
		Comments    string      `json:"comments,omitempty"`
		Tags        []WireTag   `json:"tags,omitempty"`
		Created     string      `json:"created"`
		LastUpdated string      `json:"last_updated"`
	}

	WireIPAddress struct {
		ID                 int              `json:"id"`
		URL                string           `json:"url"`
		Address            string           `json:"address"`
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
		Created            string           `json:"created"`
		LastUpdated        string           `json:"last_updated"`
	}

	WirePrefix struct {
		ID          int         `json:"id"`
		URL         string      `json:"url"`
		Prefix      string      `json:"prefix"`
		Site        *WireNested `json:"site"`
		VRF         *WireNested `json:"vrf"`
		Tenant      *WireNested `json:"tenant"`
		VLAN        *WireNested `json:"vlan"`
		Status      *WireLabel  `json:"status"`
		Role        *WireNested `json:"role"`
		IsPool      bool        `json:"is_pool"`
		Description string      `json:"description,omitempty"`
		Comments    string      `json:"comments,omitempty"`
		Tags        []WireTag   `json:"tags,omitempty"`
		Created     string      `json:"created"`
		LastUpdated string      `json:"last_updated"`
		Children    int         `json:"children"`
		Depth       int         `json:"_depth"`
		Family      *WireFamily `json:"family"`
	}

	WireVLAN struct {
		ID          int         `json:"id"`
		URL         string      `json:"url"`
		Site        *WireNested `json:"site"`
		Group       *WireNested `json:"group"`
		VID         int         `json:"vid"`
		Name        string      `json:"name"`
		Tenant      *WireNested `json:"tenant"`
		Status      *WireLabel  `json:"status"`
		Role        *WireNested `json:"role"`
		Description string      `json:"description,omitempty"`
		Comments    string      `json:"comments,omitempty"`
		Tags        []WireTag   `json:"tags,omitempty"`
		Created     string      `json:"created"`
		LastUpdated string      `json:"last_updated"`
	}

	WireVirtualMachine struct {
		ID          int         `json:"id"`
		URL         string      `json:"url"`
		Name        string      `json:"name"`
		Cluster     *WireNested `json:"cluster"`
		Role        *WireNested `json:"role"`
		Tenant      *WireNested `json:"tenant"`
		Platform    *WireNested `json:"platform"`
		Status      *WireLabel  `json:"status"`
		Site        *WireNested `json:"site"`
		VCPUs       float64     `json:"vcpus,omitempty"`
		Memory      int         `json:"memory,omitempty"`
		Disk        int         `json:"disk,omitempty"`
		Comments    string      `json:"comments,omitempty"`
		Tags        []WireTag   `json:"tags,omitempty"`
		Created     string      `json:"created"`
		LastUpdated string      `json:"last_updated"`
	}

	WireCluster struct {
		ID           int         `json:"id"`
		URL          string      `json:"url"`
		Name         string      `json:"name"`
		ClusterType  *WireNested `json:"type"`
		ClusterGroup *WireNested `json:"group"`
		Site         *WireNested `json:"site"`
		Tenant       *WireNested `json:"tenant"`
		Description  string      `json:"description,omitempty"`
		Comments     string      `json:"comments,omitempty"`
		Tags         []WireTag   `json:"tags,omitempty"`
		Created      string      `json:"created"`
		LastUpdated  string      `json:"last_updated"`
	}

	WireCircuit struct {
		ID          int         `json:"id"`
		URL         string      `json:"url"`
		CID         string      `json:"cid"`
		Provider    *WireNested `json:"provider"`
		CircuitType *WireNested `json:"circuit_type"`
		Tenant      *WireNested `json:"tenant"`
		Status      *WireLabel  `json:"status"`
		Description string      `json:"description,omitempty"`
		Comments    string      `json:"comments,omitempty"`
		Tags        []WireTag   `json:"tags,omitempty"`
		Created     string      `json:"created"`
		LastUpdated string      `json:"last_updated"`
		InstallDate string      `json:"install_date,omitempty"`
		CommitRate  int         `json:"commit_rate,omitempty"`
	}

	WireRack struct {
		ID          int         `json:"id"`
		URL         string      `json:"url"`
		Name        string      `json:"name"`
		FacilityID  string      `json:"facility_id,omitempty"`
		Site        *WireNested `json:"site"`
		Location    *WireNested `json:"location"`
		Tenant      *WireNested `json:"tenant"`
		Status      *WireLabel  `json:"status"`
		Role        *WireNested `json:"role"`
		Serial      string      `json:"serial,omitempty"`
		AssetTag    string      `json:"asset_tag,omitempty"`
		Type        *WireLabel  `json:"type"`
		Width       int         `json:"width,omitempty"`
		UHeight     int         `json:"u_height,omitempty"`
		Comments    string      `json:"comments,omitempty"`
		Tags        []WireTag   `json:"tags,omitempty"`
		Created     string      `json:"created"`
		LastUpdated string      `json:"last_updated"`
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
