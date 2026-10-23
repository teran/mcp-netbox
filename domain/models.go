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

	// Region is a NetBox region (a geographic grouping of sites).
	Region struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		Parent       *Nested        `json:"parent"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// TenantGroup is a NetBox tenancy tenant group (a grouping of tenants).
	TenantGroup struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		Parent       *Nested        `json:"parent"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// Device is a NetBox device (a physical piece of network infrastructure).
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

	// IPAddress is a NetBox IP address entry.
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

	// Prefix is a NetBox IP prefix (CIDR block).
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

	// VLAN is a NetBox VLAN.
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

	// VirtualMachine is a NetBox virtual machine.
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

	// VirtualDisk is a NetBox virtualization virtual disk attached to a virtual
	// machine.
	VirtualDisk struct {
		ID             int            `json:"id"`
		URL            string         `json:"url,omitempty"`
		Name           string         `json:"name"`
		Display        string         `json:"display,omitempty"`
		VirtualMachine *Nested        `json:"virtual_machine"`
		Size           int            `json:"size"`
		Description    string         `json:"description,omitempty"`
		Tags           []Tag          `json:"tags,omitempty"`
		CustomFields   map[string]any `json:"custom_fields,omitempty"`
		Created        string         `json:"created"`
		LastUpdated    string         `json:"last_updated"`
	}

	// VirtualMachineType is a NetBox virtualization virtual machine type (a
	// reusable hardware/software profile for virtual machines).
	VirtualMachineType struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// Cluster is a NetBox virtualization cluster.
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

	// Circuit is a NetBox circuit.
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

	// CircuitTermination is a NetBox circuit termination.
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

	// Cable is a NetBox cable connection.
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

	// Rack is a NetBox equipment rack.
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

	// Interface is a NetBox device interface.
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

	// VMInterface is a NetBox virtual-machine interface.
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

	// Provider is a circuits provider (a secondary entity: it has write tools
	// but no read tool).
	Provider struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug,omitempty"`
		Display      string         `json:"display,omitempty"`
		Asn          int            `json:"asn,omitempty"`
		Account      string         `json:"account,omitempty"`
		PortalURL    string         `json:"portal_url,omitempty"`
		NocContact   string         `json:"noc_contact,omitempty"`
		AdminContact string         `json:"admin_contact,omitempty"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// Tenant is a tenancy tenant (a secondary entity: it has write tools but
	// no read tool).
	Tenant struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug,omitempty"`
		Display      string         `json:"display,omitempty"`
		Group        *Nested        `json:"group"`
		Description  string         `json:"description,omitempty"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// Vrf is a NetBox VRF (a secondary entity: it has write tools but no read
	// tool).
	Vrf struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Rd           string         `json:"rd,omitempty"`
		Tenant       *Nested        `json:"tenant"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// VlanGroup is a NetBox VLAN group (a secondary entity: it has write tools
	// but no read tool).
	VlanGroup struct {
		ID           int             `json:"id"`
		URL          string          `json:"url,omitempty"`
		Name         string          `json:"name"`
		Slug         string          `json:"slug,omitempty"`
		Display      string          `json:"display,omitempty"`
		Description  string          `json:"description,omitempty"`
		Scope        *VlanGroupScope `json:"scope"`
		MinVID       *int            `json:"min_vid"`
		MaxVID       *int            `json:"max_vid"`
		Tags         []Tag           `json:"tags,omitempty"`
		CustomFields map[string]any  `json:"custom_fields,omitempty"`
		Created      string          `json:"created"`
		LastUpdated  string          `json:"last_updated"`
	}

	// VlanGroupScope is the polymorphic scope reference on a VLAN group.
	VlanGroupScope struct {
		ObjectType string `json:"object_type"`
		ObjectID   int    `json:"object_id"`
	}

	// Role is a NetBox IPAM role (a secondary entity: it has write tools but
	// no read tool).
	Role struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug,omitempty"`
		Weight       int            `json:"weight"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// Contact is a NetBox tenancy contact (a secondary entity: it has write
	// tools but no read tool).
	Contact struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Group        *Nested        `json:"group"`
		Title        string         `json:"title,omitempty"`
		Phone        string         `json:"phone,omitempty"`
		Email        string         `json:"email,omitempty"`
		Address      string         `json:"address,omitempty"`
		Link         string         `json:"link,omitempty"`
		Description  string         `json:"description,omitempty"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// ContactGroup is a NetBox tenancy contact group (a grouping of contacts).
	ContactGroup struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug"`
		Display      string         `json:"display,omitempty"`
		Parent       *Nested        `json:"parent"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// ContactRole is a NetBox tenancy contact role (e.g. an administrative or
	// operational role assigned to contacts).
	ContactRole struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// Manufacturer is a device manufacturer (a secondary entity: it has write
	// tools but no read tool).
	Manufacturer struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug,omitempty"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// DeviceType is a device type (a secondary entity: it has write tools but
	// no read tool). NetBox names the model field "model", not "name".
	DeviceType struct {
		ID            int            `json:"id"`
		URL           string         `json:"url,omitempty"`
		Manufacturer  *Nested        `json:"manufacturer"`
		Model         string         `json:"model"`
		Slug          string         `json:"slug,omitempty"`
		Display       string         `json:"display,omitempty"`
		PartNumber    string         `json:"part_number,omitempty"`
		UHeight       float64        `json:"u_height,omitempty"`
		IsFullDepth   bool           `json:"is_full_depth,omitempty"`
		SubdeviceRole *Label         `json:"subdevice_role,omitempty"`
		Comments      string         `json:"comments,omitempty"`
		Tags          []Tag          `json:"tags,omitempty"`
		CustomFields  map[string]any `json:"custom_fields,omitempty"`
		Created       string         `json:"created"`
		LastUpdated   string         `json:"last_updated"`
	}

	// Location is a NetBox site location (a secondary entity: it has write
	// tools but no read tool).
	Location struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug,omitempty"`
		Display      string         `json:"display,omitempty"`
		Site         *Nested        `json:"site"`
		Parent       *Nested        `json:"parent"`
		Status       *Label         `json:"status"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// ClusterType is a NetBox cluster type (a secondary entity: it has write
	// tools but no read tool).
	ClusterType struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug,omitempty"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// ClusterGroup is a NetBox cluster group (a secondary entity: it has write
	// tools but no read tool).
	ClusterGroup struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug,omitempty"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// CircuitType is a NetBox circuit type (a secondary entity: it has write
	// tools but no read tool).
	CircuitType struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug,omitempty"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		CircuitCount int            `json:"circuit_count"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// SiteGroup is a NetBox DCIM site group (a grouping of sites).
	SiteGroup struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug"`
		Display      string         `json:"display,omitempty"`
		Parent       *Nested        `json:"parent"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// RackGroup is a NetBox DCIM rack group (a grouping of racks).
	RackGroup struct {
		ID          int     `json:"id"`
		URL         string  `json:"url,omitempty"`
		Name        string  `json:"name"`
		Slug        string  `json:"slug"`
		Display     string  `json:"display,omitempty"`
		Site        *Nested `json:"site"`
		Parent      *Nested `json:"parent"`
		Description string  `json:"description,omitempty"`
		Tags        []Tag   `json:"tags,omitempty"`
		Created     string  `json:"created"`
		LastUpdated string  `json:"last_updated"`
	}

	// RackRole is a NetBox DCIM rack role (a functional role for racks).
	RackRole struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug"`
		Color        string         `json:"color,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// DeviceRole is a NetBox DCIM device role (a functional role for devices).
	DeviceRole struct {
		ID             int            `json:"id"`
		URL            string         `json:"url,omitempty"`
		Name           string         `json:"name"`
		Slug           string         `json:"slug"`
		Color          string         `json:"color,omitempty"`
		VMRole         bool           `json:"vm_role"`
		ConfigTemplate *Nested        `json:"config_template"`
		Description    string         `json:"description,omitempty"`
		Tags           []Tag          `json:"tags,omitempty"`
		CustomFields   map[string]any `json:"custom_fields,omitempty"`
		Created        string         `json:"created"`
		LastUpdated    string         `json:"last_updated"`
		Weight         int            `json:"weight"`
	}

	// Platform is a NetBox DCIM platform (a software/firmware platform).
	Platform struct {
		ID             int            `json:"id"`
		URL            string         `json:"url,omitempty"`
		Name           string         `json:"name"`
		Slug           string         `json:"slug"`
		Manufacturer   *Nested        `json:"manufacturer"`
		ConfigTemplate *Nested        `json:"config_template"`
		Description    string         `json:"description,omitempty"`
		Tags           []Tag          `json:"tags,omitempty"`
		CustomFields   map[string]any `json:"custom_fields,omitempty"`
		Created        string         `json:"created"`
		LastUpdated    string         `json:"last_updated"`
	}

	// VirtualChassis is a NetBox DCIM virtual chassis (a grouping of devices).
	VirtualChassis struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Domain       string         `json:"domain,omitempty"`
		Master       *Nested        `json:"master"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
		MemberCount  int            `json:"member_count"`
	}

	// PowerPanel is a NetBox DCIM power panel.
	PowerPanel struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Site         *Nested        `json:"site"`
		Location     *Nested        `json:"location"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// InventoryItemRole is a NetBox DCIM inventory item role.
	InventoryItemRole struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug"`
		Color        string         `json:"color,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// FrontPort is a NetBox DCIM front port.
	FrontPort struct {
		ID               int            `json:"id"`
		URL              string         `json:"url,omitempty"`
		Name             string         `json:"name"`
		Display          string         `json:"display,omitempty"`
		Device           *Nested        `json:"device"`
		Type             *Label         `json:"type"`
		RearPort         *Nested        `json:"rear_port"`
		RearPortPosition int            `json:"rear_port_position,omitempty"`
		Description      string         `json:"description,omitempty"`
		Tags             []Tag          `json:"tags,omitempty"`
		CustomFields     map[string]any `json:"custom_fields,omitempty"`
		Created          string         `json:"created"`
		LastUpdated      string         `json:"last_updated"`
	}

	// RearPort is a NetBox DCIM rear port.
	RearPort struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Device       *Nested        `json:"device"`
		Type         *Label         `json:"type"`
		Positions    int            `json:"positions"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// ConsolePort is a NetBox DCIM console port.
	ConsolePort struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Device       *Nested        `json:"device"`
		Type         *Label         `json:"type"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// ConsoleServerPort is a NetBox DCIM console server port.
	ConsoleServerPort struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Device       *Nested        `json:"device"`
		Type         *Label         `json:"type"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// PowerPort is a NetBox DCIM power port.
	PowerPort struct {
		ID            int            `json:"id"`
		URL           string         `json:"url,omitempty"`
		Name          string         `json:"name"`
		Display       string         `json:"display,omitempty"`
		Device        *Nested        `json:"device"`
		Type          *Label         `json:"type"`
		MaximumDraw   int            `json:"maximum_draw,omitempty"`
		AllocatedDraw int            `json:"allocated_draw,omitempty"`
		Description   string         `json:"description,omitempty"`
		Tags          []Tag          `json:"tags,omitempty"`
		CustomFields  map[string]any `json:"custom_fields,omitempty"`
		Created       string         `json:"created"`
		LastUpdated   string         `json:"last_updated"`
	}

	// PowerOutlet is a NetBox DCIM power outlet.
	PowerOutlet struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Device       *Nested        `json:"device"`
		Type         *Label         `json:"type"`
		PowerPort    *Nested        `json:"power_port"`
		FeedLeg      *Label         `json:"feed_leg"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// PowerFeed is a NetBox DCIM power feed.
	PowerFeed struct {
		ID             int            `json:"id"`
		URL            string         `json:"url,omitempty"`
		Name           string         `json:"name"`
		Display        string         `json:"display,omitempty"`
		PowerPanel     *Nested        `json:"power_panel"`
		Status         *Label         `json:"status"`
		Type           *Label         `json:"type"`
		Supply         *Label         `json:"supply"`
		Phase          *Label         `json:"phase"`
		Voltage        int            `json:"voltage"`
		Amperage       int            `json:"amperage"`
		MaxUtilization int            `json:"max_utilization,omitempty"`
		Description    string         `json:"description,omitempty"`
		Tags           []Tag          `json:"tags,omitempty"`
		CustomFields   map[string]any `json:"custom_fields,omitempty"`
		Created        string         `json:"created"`
		LastUpdated    string         `json:"last_updated"`
	}

	// DeviceBay is a NetBox DCIM device bay.
	DeviceBay struct {
		ID              int            `json:"id"`
		URL             string         `json:"url,omitempty"`
		Name            string         `json:"name"`
		Display         string         `json:"display,omitempty"`
		Device          *Nested        `json:"device"`
		InstalledDevice *Nested        `json:"installed_device"`
		Description     string         `json:"description,omitempty"`
		Tags            []Tag          `json:"tags,omitempty"`
		CustomFields    map[string]any `json:"custom_fields,omitempty"`
		Created         string         `json:"created"`
		LastUpdated     string         `json:"last_updated"`
	}

	// RackReservation is a NetBox DCIM rack reservation.
	RackReservation struct {
		ID          int     `json:"id"`
		URL         string  `json:"url,omitempty"`
		Rack        *Nested `json:"rack"`
		Units       []int   `json:"units,omitempty"`
		Description string  `json:"description"`
		User        *Nested `json:"user"`
		Created     string  `json:"created"`
		LastUpdated string  `json:"last_updated"`
	}

	// InventoryItem is a NetBox DCIM inventory item.
	InventoryItem struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Device       *Nested        `json:"device"`
		Role         *Nested        `json:"role"`
		Manufacturer *Nested        `json:"manufacturer"`
		Serial       string         `json:"serial,omitempty"`
		PartID       string         `json:"part_id,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// MACAddress is a NetBox DCIM MAC address.
	MACAddress struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Address      string         `json:"address"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// RIR is a NetBox IPAM RIR (regional internet registry).
	RIR struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug,omitempty"`
		Display      string         `json:"display,omitempty"`
		IsPrivate    bool           `json:"is_private"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// Aggregate is a NetBox IPAM aggregate (a top-level prefix block).
	Aggregate struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Prefix       string         `json:"prefix"`
		Display      string         `json:"display,omitempty"`
		RIR          *Nested        `json:"rir"`
		Tenant       *Nested        `json:"tenant"`
		DateAdded    string         `json:"date_added,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// IPRange is a NetBox IPAM IP range (a contiguous block of addresses).
	IPRange struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		StartAddress string         `json:"start_address"`
		EndAddress   string         `json:"end_address"`
		Display      string         `json:"display,omitempty"`
		VRF          *Nested        `json:"vrf"`
		Tenant       *Nested        `json:"tenant"`
		Status       *Label         `json:"status"`
		Role         *Nested        `json:"role"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// ASN is a NetBox IPAM autonomous system number.
	ASN struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		ASN          int            `json:"asn"`
		Display      string         `json:"display,omitempty"`
		RIR          *Nested        `json:"rir"`
		Tenant       *Nested        `json:"tenant"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// ASNRange is a NetBox IPAM ASN range (a contiguous block of AS numbers).
	ASNRange struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Start        int            `json:"start"`
		End          int            `json:"end"`
		RIR          *Nested        `json:"rir"`
		Tenant       *Nested        `json:"tenant"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// RouteTarget is a NetBox IPAM route target.
	RouteTarget struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Tenant       *Nested        `json:"tenant"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// Service is a NetBox IPAM service (a network service on a device/VM).
	Service struct {
		ID               int            `json:"id"`
		URL              string         `json:"url,omitempty"`
		Name             string         `json:"name"`
		Display          string         `json:"display,omitempty"`
		ParentObjectType string         `json:"parent_object_type"`
		ParentObjectID   int            `json:"parent_object_id"`
		Parent           *Nested        `json:"parent"`
		Protocol         *Label         `json:"protocol"`
		Ports            []int          `json:"ports,omitempty"`
		Description      string         `json:"description,omitempty"`
		Tags             []Tag          `json:"tags,omitempty"`
		CustomFields     map[string]any `json:"custom_fields,omitempty"`
		Created          string         `json:"created"`
		LastUpdated      string         `json:"last_updated"`
	}

	// ServiceTemplate is a NetBox IPAM service template (a reusable service).
	ServiceTemplate struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Protocol     *Label         `json:"protocol"`
		Ports        []int          `json:"ports,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// FHRPGroup is a NetBox IPAM FHRP (first-hop redundancy protocol) group.
	FHRPGroup struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Protocol     *Label         `json:"protocol"`
		GroupID      int            `json:"group_id"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// FHRPGroupAssignment is a NetBox IPAM FHRP group assignment (the mapping
	// of an FHRP group to an IP address).
	FHRPGroupAssignment struct {
		ID            int     `json:"id"`
		URL           string  `json:"url,omitempty"`
		Group         *Nested `json:"group"`
		InterfaceType string  `json:"interface_type"`
		InterfaceID   int     `json:"interface_id"`
		Interface     *Nested `json:"interface"`
		Priority      int     `json:"priority"`
		Created       string  `json:"created"`
		LastUpdated   string  `json:"last_updated"`
	}

	// VLANTranslationPolicy is a NetBox IPAM VLAN translation policy.
	VLANTranslationPolicy struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// VLANTranslationRule is a NetBox IPAM VLAN translation rule.
	VLANTranslationRule struct {
		ID          int    `json:"id"`
		URL         string `json:"url,omitempty"`
		Display     string `json:"display,omitempty"`
		Policy      int    `json:"policy"`
		LocalVID    int    `json:"local_vid"`
		RemoteVID   int    `json:"remote_vid"`
		Description string `json:"description,omitempty"`
		Created     string `json:"created"`
		LastUpdated string `json:"last_updated"`
	}

	// ProviderAccount is a NetBox circuits provider account (a billing/technical
	// account held with a provider).
	ProviderAccount struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Provider     *Nested        `json:"provider"`
		Account      string         `json:"account,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// ProviderNetwork is a NetBox circuits provider network (a logical network
	// offered by a provider).
	ProviderNetwork struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Provider     *Nested        `json:"provider"`
		Description  string         `json:"description,omitempty"`
		Comments     string         `json:"comments,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// CircuitGroup is a NetBox circuits circuit group (a logical grouping of
	// circuits).
	CircuitGroup struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// CircuitGroupAssignment is a NetBox circuits circuit group assignment (the
	// membership of a circuit in a circuit group).
	CircuitGroupAssignment struct {
		ID          int     `json:"id"`
		URL         string  `json:"url,omitempty"`
		Group       *Nested `json:"group"`
		Member      *Nested `json:"member"`
		Priority    *Label  `json:"priority,omitempty"`
		Created     string  `json:"created"`
		LastUpdated string  `json:"last_updated"`
	}

	// VirtualCircuit is a NetBox circuits virtual circuit (a logical circuit
	// carried over a provider network).
	VirtualCircuit struct {
		ID              int            `json:"id"`
		URL             string         `json:"url,omitempty"`
		CID             string         `json:"cid,omitempty"`
		Type            *Nested        `json:"type"`
		ProviderNetwork *Nested        `json:"provider_network"`
		Status          *Label         `json:"status,omitempty"`
		Description     string         `json:"description,omitempty"`
		Tags            []Tag          `json:"tags,omitempty"`
		CustomFields    map[string]any `json:"custom_fields,omitempty"`
		Created         string         `json:"created"`
		LastUpdated     string         `json:"last_updated"`
	}

	// VirtualCircuitType is a NetBox circuits virtual circuit type (e.g. VPLS,
	// VXLAN).
	VirtualCircuitType struct {
		ID           int            `json:"id"`
		URL          string         `json:"url,omitempty"`
		Name         string         `json:"name"`
		Slug         string         `json:"slug"`
		Display      string         `json:"display,omitempty"`
		Description  string         `json:"description,omitempty"`
		Tags         []Tag          `json:"tags,omitempty"`
		CustomFields map[string]any `json:"custom_fields,omitempty"`
		Created      string         `json:"created"`
		LastUpdated  string         `json:"last_updated"`
	}

	// VirtualCircuitTermination is a NetBox circuits virtual circuit termination
	// (the attachment of a virtual circuit to an interface).
	VirtualCircuitTermination struct {
		ID             int     `json:"id"`
		URL            string  `json:"url,omitempty"`
		VirtualCircuit *Nested `json:"virtual_circuit"`
		Role           *Label  `json:"role,omitempty"`
		Interface      *Nested `json:"interface"`
		Created        string  `json:"created"`
		LastUpdated    string  `json:"last_updated"`
	}
)

type (
	// Nested is a compact reference to a related NetBox object (id + url + name + slug).
	Nested struct {
		ID   int    `json:"id"`
		URL  string `json:"url"`
		Name string `json:"name"`
		Slug string `json:"slug,omitempty"`
	}

	// Label is a NetBox choice field carrying both the value and its display label.
	Label struct {
		Value string `json:"value"`
		Label string `json:"label"`
	}

	// Tag is a NetBox tag attached to an object.
	Tag struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
		URL  string `json:"url"`
	}

	// CustomField is a NetBox extras custom field definition.
	CustomField struct {
		ID           int      `json:"id"`
		URL          string   `json:"url"`
		Name         string   `json:"name"`
		Label        string   `json:"label"`
		Type         string   `json:"type"`
		ContentTypes []string `json:"object_types"`
		GroupName    string   `json:"group_name"`
		Description  string   `json:"description"`
		Required     bool     `json:"required"`
		Default      string   `json:"default"`
		FilterLogic  string   `json:"filter_logic"`
		Weight       int      `json:"weight"`
		UIVisibility string   `json:"ui_visibility"`
		Choices      []string `json:"choices"`
	}

	// AssignedObject is the polymorphic reference to the object an IP address is assigned to.
	AssignedObject struct {
		ID     int     `json:"id"`
		URL    string  `json:"url"`
		Device *Nested `json:"device"`
		Name   string  `json:"name"`
	}

	// Family is the IP address family (IPv4/IPv6) of a prefix.
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

// RegionWrite is the domain DTO for creating or updating a NetBox region. It
// carries only the writable fields.
//
// Optional scalar fields are pointers so that a partial-update (PATCH) can
// distinguish "not provided" (nil) from "explicitly set to empty" (pointer to
// ""). Name is a plain string because it is required on create; the omitempty
// tag drops it from a PATCH body when not provided. Parent references the
// numeric NetBox ID of the parent region.
type RegionWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Parent       *int           `json:"parent,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// TenantGroupWrite is the domain DTO for creating or updating a NetBox tenancy
// tenant group. It carries only the writable fields.
//
// Optional scalar fields are pointers so that a partial-update (PATCH) can
// distinguish "not provided" (nil) from "explicitly set to empty" (pointer to
// ""). Name is a plain string because it is required on create; the omitempty
// tag drops it from a PATCH body when not provided. Parent references the
// numeric NetBox ID of the parent tenant group.
type TenantGroupWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Parent       *int           `json:"parent,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
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

// CableTerminationWrite describes one end of a cable for create/update. NetBox
// represents each termination as a nested object carrying the object_type
// (e.g. "dcim.interface") and the numeric object_id. Both fields are required
// on create.
type CableTerminationWrite struct {
	ObjectType string `json:"object_type"`
	ObjectID   int    `json:"object_id"`
}

// CableWrite is the domain DTO for creating or updating a NetBox cable. It
// carries only the writable fields.
//
// TerminationA and TerminationB are required on create; each is a nested
// object with object_type and object_id. They are pointers so a partial-update
// (PATCH) can omit them (moving a termination is done via other NetBox
// operations). Optional scalar fields are pointers so that a PATCH can
// distinguish "not provided" (nil) from "explicitly set to empty". Type, Status
// and LengthUnit are choice-string values.
type CableWrite struct {
	TerminationA *CableTerminationWrite `json:"termination_a,omitempty"`
	TerminationB *CableTerminationWrite `json:"termination_b,omitempty"`
	Type         *string                `json:"type,omitempty"`
	Status       *string                `json:"status,omitempty"`
	Label        *string                `json:"label,omitempty"`
	Color        *string                `json:"color,omitempty"`
	Length       *float64               `json:"length,omitempty"`
	LengthUnit   *string                `json:"length_unit,omitempty"`
	Description  *string                `json:"description,omitempty"`
	Tags         []string               `json:"tags,omitempty"`
	CustomFields map[string]any         `json:"custom_fields,omitempty"`
}

// VMInterfaceWrite is the domain DTO for creating or updating a NetBox
// virtual-machine interface. It carries only the writable fields.
//
// VirtualMachine is the numeric NetBox ID of the parent VM and Name is required
// on create. Optional scalar fields are pointers so that a partial-update
// (PATCH) can distinguish "not provided" (nil) from "explicitly set to empty".
type VMInterfaceWrite struct {
	VirtualMachine *int           `json:"virtual_machine,omitempty"`
	Name           string         `json:"name,omitempty"`
	Enabled        *bool          `json:"enabled,omitempty"`
	MTU            *int           `json:"mtu,omitempty"`
	MACAddress     *string        `json:"mac_address,omitempty"`
	Description    *string        `json:"description,omitempty"`
	Tags           []string       `json:"tags,omitempty"`
	CustomFields   map[string]any `json:"custom_fields,omitempty"`
}

// ProviderWrite is the domain DTO for creating or updating a NetBox circuits
// provider. It carries only the writable fields.
//
// Name is required on create. Optional scalar fields are pointers so that a
// partial-update (PATCH) can distinguish "not provided" (nil) from "explicitly
// set to empty". Asn is a numeric AS number.
type ProviderWrite struct {
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

// TenantWrite is the domain DTO for creating or updating a NetBox tenancy
// tenant. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty" (pointer to ""). Group is a pointer to
// the numeric tenant-group ID and is optional (a tenant may exist without a
// group).
type TenantWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Group        *int           `json:"group,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Comments     *string        `json:"comments,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// ManufacturerWrite is the domain DTO for creating or updating a NetBox device
// manufacturer. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty" (pointer to "").
type ManufacturerWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// DeviceTypeWrite is the domain DTO for creating or updating a NetBox device
// type. It carries only the writable fields.
//
// Manufacturer (numeric NetBox ID) and Model are required on create. Model is a
// plain string because it is required on create; the omitempty tag drops it
// from a PATCH body when not provided. Optional scalar fields are pointers so
// that a partial-update (PATCH) can distinguish "not provided" (nil) from
// "explicitly set to empty" (pointer to ""). IsFullDepth is a pointer so an
// update can set it to false explicitly. SubdeviceRole is a choice-string
// value.
type DeviceTypeWrite struct {
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

// LocationWrite is the domain DTO for creating or updating a NetBox site
// location. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Site (numeric NetBox ID) is
// required on create. Optional scalar fields are pointers so that a
// partial-update (PATCH) can distinguish "not provided" (nil) from "explicitly
// set to empty" (pointer to ""). Parent is the numeric NetBox ID of a parent
// location. Status is a choice-string value.
type LocationWrite struct {
	Name         string         `json:"name,omitempty"`
	Site         *int           `json:"site,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Parent       *int           `json:"parent,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Status       *string        `json:"status,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// ClusterTypeWrite is the domain DTO for creating or updating a NetBox cluster
// type. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty" (pointer to "").
type ClusterTypeWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// ClusterGroupWrite is the domain DTO for creating or updating a NetBox cluster
// group. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty" (pointer to "").
type ClusterGroupWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// CircuitTypeWrite is the domain DTO for creating or updating a NetBox circuit
// type. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty" (pointer to "").
type CircuitTypeWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// VrfWrite is the domain DTO for creating or updating a NetBox VRF. It carries
// only the writable fields.
//
// Name and Rd are plain strings because they are required on create; the
// omitempty tags drop them from a PATCH body when not provided. Optional scalar
// fields are pointers so that a partial-update (PATCH) can distinguish "not
// provided" (nil) from "explicitly set to empty" (pointer to ""). Tenant is a
// pointer to the numeric NetBox tenant ID.
type VrfWrite struct {
	Name         string         `json:"name,omitempty"`
	Rd           string         `json:"rd,omitempty"`
	Tenant       *int           `json:"tenant,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// VlanGroupWrite is the domain DTO for creating or updating a NetBox VLAN
// group. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty" (pointer to ""). The polymorphic scope
// is written via scope_type (content type) and scope_id (object ID).
type VlanGroupWrite struct {
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

// RoleWrite is the domain DTO for creating or updating a NetBox IPAM role. It
// carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields (slug,
// weight, description) are pointers so that a partial-update (PATCH) can
// distinguish "not provided" (nil) from "explicitly set to empty".
type RoleWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Weight       *int           `json:"weight,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// TagWrite is the domain DTO for creating or updating a NetBox extras tag. It
// carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields (slug,
// color, description) are pointers so that a partial-update (PATCH) can
// distinguish "not provided" (nil) from "explicitly set to empty".
type TagWrite struct {
	Name        string  `json:"name,omitempty"`
	Slug        *string `json:"slug,omitempty"`
	Color       *string `json:"color,omitempty"`
	Description *string `json:"description,omitempty"`
}

// CustomFieldWrite is the domain DTO for creating or updating a NetBox extras
// custom field. It carries only the writable fields.
//
// Name and Label are plain strings because they are required on create; the
// omitempty tags drop them from a PATCH body when not provided. Optional scalar
// fields are pointers so that a partial-update (PATCH) can distinguish "not
// provided" (nil) from "explicitly set to empty".
type CustomFieldWrite struct {
	Name         string   `json:"name,omitempty"`
	Label        string   `json:"label,omitempty"`
	Type         *string  `json:"type,omitempty"`
	ContentTypes []string `json:"object_types,omitempty"`
	GroupName    *string  `json:"group_name,omitempty"`
	Description  *string  `json:"description,omitempty"`
	Required     *bool    `json:"required,omitempty"`
	Default      *string  `json:"default,omitempty"`
	FilterLogic  *string  `json:"filter_logic,omitempty"`
	Weight       *int     `json:"weight,omitempty"`
	UIVisibility *string  `json:"ui_visibility,omitempty"`
	Choices      []string `json:"choices,omitempty"`
}

// ContactWrite is the domain DTO for creating or updating a NetBox tenancy
// contact. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Group is a pointer to the
// numeric contact-group ID and is optional (a contact may exist without a
// group). Other optional scalar fields are pointers so that a partial-update
// (PATCH) can distinguish "not provided" (nil) from "explicitly set to empty".
type ContactWrite struct {
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

// ContactGroupWrite is the domain DTO for creating or updating a NetBox tenancy
// contact group. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". Parent references the numeric NetBox ID
// of the parent contact group.
type ContactGroupWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Parent       *int           `json:"parent,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// ContactRoleWrite is the domain DTO for creating or updating a NetBox tenancy
// contact role. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty".
type ContactRoleWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// VirtualDiskWrite is the domain DTO for creating or updating a NetBox
// virtualization virtual disk. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". VirtualMachine references the numeric
// NetBox ID of the parent virtual machine.
type VirtualDiskWrite struct {
	Name           string         `json:"name,omitempty"`
	VirtualMachine *int           `json:"virtual_machine,omitempty"`
	Size           *int           `json:"size,omitempty"`
	Description    *string        `json:"description,omitempty"`
	Tags           []string       `json:"tags,omitempty"`
	CustomFields   map[string]any `json:"custom_fields,omitempty"`
}

// VirtualMachineTypeWrite is the domain DTO for creating or updating a NetBox
// virtualization virtual machine type. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty".
type VirtualMachineTypeWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// SiteGroupWrite is the domain DTO for creating or updating a NetBox DCIM site
// group. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". Parent references the numeric NetBox ID
// of the parent site group.
type SiteGroupWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Parent       *int           `json:"parent,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// RackGroupWrite is the domain DTO for creating or updating a NetBox DCIM rack
// group. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". Site and Parent reference the numeric
// NetBox IDs of the parent site and rack group.
type RackGroupWrite struct {
	Name        string   `json:"name,omitempty"`
	Slug        *string  `json:"slug,omitempty"`
	Site        *int     `json:"site,omitempty"`
	Parent      *int     `json:"parent,omitempty"`
	Description *string  `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// RackRoleWrite is the domain DTO for creating or updating a NetBox DCIM rack
// role. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty".
type RackRoleWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Color        *string        `json:"color,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// DeviceRoleWrite is the domain DTO for creating or updating a NetBox DCIM
// device role. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". VMRole and Weight are pointers so an
// update can set them to zero/false explicitly. ConfigTemplate is the numeric
// NetBox ID of the config template.
type DeviceRoleWrite struct {
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

// PlatformWrite is the domain DTO for creating or updating a NetBox DCIM
// platform. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". Manufacturer and ConfigTemplate are the
// numeric NetBox IDs of the related objects.
type PlatformWrite struct {
	Name           string         `json:"name,omitempty"`
	Slug           *string        `json:"slug,omitempty"`
	Manufacturer   *int           `json:"manufacturer,omitempty"`
	ConfigTemplate *int           `json:"config_template,omitempty"`
	Description    *string        `json:"description,omitempty"`
	Tags           []string       `json:"tags,omitempty"`
	CustomFields   map[string]any `json:"custom_fields,omitempty"`
}

// VirtualChassisWrite is the domain DTO for creating or updating a NetBox DCIM
// virtual chassis. It carries only the writable fields.
//
// Name is a plain string; it is optional on create (unlike most entities), so
// the create handler does not require it. Optional scalar fields are pointers
// so that a partial-update (PATCH) can distinguish "not provided" (nil) from
// "explicitly set to empty". Master is the numeric NetBox ID of the master
// device.
type VirtualChassisWrite struct {
	Name         string         `json:"name,omitempty"`
	Domain       *string        `json:"domain,omitempty"`
	Master       *int           `json:"master,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// PowerPanelWrite is the domain DTO for creating or updating a NetBox DCIM
// power panel. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". Site and Location are the numeric
// NetBox IDs of the related objects.
type PowerPanelWrite struct {
	Site         *int           `json:"site,omitempty"`
	Location     *int           `json:"location,omitempty"`
	Name         string         `json:"name,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// InventoryItemRoleWrite is the domain DTO for creating or updating a NetBox
// DCIM inventory item role. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty".
type InventoryItemRoleWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Color        *string        `json:"color,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// FrontPortWrite is the domain DTO for creating or updating a NetBox DCIM
// front port. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". Device and RearPort are the numeric
// NetBox IDs of the parent device and rear port; Type is a choice-string value.
type FrontPortWrite struct {
	Name             string         `json:"name,omitempty"`
	Device           *int           `json:"device,omitempty"`
	Type             *string        `json:"type,omitempty"`
	RearPort         *int           `json:"rear_port,omitempty"`
	RearPortPosition *int           `json:"rear_port_position,omitempty"`
	Description      *string        `json:"description,omitempty"`
	Tags             []string       `json:"tags,omitempty"`
	CustomFields     map[string]any `json:"custom_fields,omitempty"`
}

// RearPortWrite is the domain DTO for creating or updating a NetBox DCIM rear
// port. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". Device is the numeric NetBox ID of the
// parent device; Type is a choice-string value.
type RearPortWrite struct {
	Name         string         `json:"name,omitempty"`
	Device       *int           `json:"device,omitempty"`
	Type         *string        `json:"type,omitempty"`
	Positions    *int           `json:"positions,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// ConsolePortWrite is the domain DTO for creating or updating a NetBox DCIM
// console port. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". Device is the numeric NetBox ID of the
// parent device; Type is a choice-string value.
type ConsolePortWrite struct {
	Name         string         `json:"name,omitempty"`
	Device       *int           `json:"device,omitempty"`
	Type         *string        `json:"type,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// ConsoleServerPortWrite is the domain DTO for creating or updating a NetBox
// DCIM console server port. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". Device is the numeric NetBox ID of the
// parent device; Type is a choice-string value.
type ConsoleServerPortWrite struct {
	Name         string         `json:"name,omitempty"`
	Device       *int           `json:"device,omitempty"`
	Type         *string        `json:"type,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// PowerPortWrite is the domain DTO for creating or updating a NetBox DCIM
// power port. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". Device is the numeric NetBox ID of the
// parent device; Type is a choice-string value.
type PowerPortWrite struct {
	Name          string         `json:"name,omitempty"`
	Device        *int           `json:"device,omitempty"`
	Type          *string        `json:"type,omitempty"`
	MaximumDraw   *int           `json:"maximum_draw,omitempty"`
	AllocatedDraw *int           `json:"allocated_draw,omitempty"`
	Description   *string        `json:"description,omitempty"`
	Tags          []string       `json:"tags,omitempty"`
	CustomFields  map[string]any `json:"custom_fields,omitempty"`
}

// PowerOutletWrite is the domain DTO for creating or updating a NetBox DCIM
// power outlet. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". Device and PowerPort are the numeric
// NetBox IDs of the parent device and power port; Type and FeedLeg are
// choice-string values.
type PowerOutletWrite struct {
	Name         string         `json:"name,omitempty"`
	Device       *int           `json:"device,omitempty"`
	Type         *string        `json:"type,omitempty"`
	PowerPort    *int           `json:"power_port,omitempty"`
	FeedLeg      *string        `json:"feed_leg,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// PowerFeedWrite is the domain DTO for creating or updating a NetBox DCIM
// power feed. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". PowerPanel is the numeric NetBox ID of
// the parent power panel; Status, Type, Supply and Phase are choice-string
// values.
type PowerFeedWrite struct {
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

// DeviceBayWrite is the domain DTO for creating or updating a NetBox DCIM
// device bay. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". Device is the numeric NetBox ID of the
// parent device.
type DeviceBayWrite struct {
	Name         string         `json:"name,omitempty"`
	Device       *int           `json:"device,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// RackReservationWrite is the domain DTO for creating or updating a NetBox DCIM
// rack reservation. It carries only the writable fields.
//
// Description is a plain string because it is required on create; the omitempty
// tag drops it from a PATCH body when not provided. Rack is the numeric NetBox
// ID of the rack; Units lists the reserved rack units.
type RackReservationWrite struct {
	Description  string         `json:"description,omitempty"`
	Rack         *int           `json:"rack,omitempty"`
	Units        []int          `json:"units,omitempty"`
	User         *int           `json:"user,omitempty"`
	Tenant       *int           `json:"tenant,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// InventoryItemWrite is the domain DTO for creating or updating a NetBox DCIM
// inventory item. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". Device, Role and Manufacturer are the
// numeric NetBox IDs of the related objects.
type InventoryItemWrite struct {
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

// MACAddressWrite is the domain DTO for creating or updating a NetBox DCIM MAC
// address. It carries only the writable fields.
//
// Address is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. AssignedObjectType is the
// content type of the assigned object (e.g. "dcim.interface") and
// AssignedObjectID is its numeric NetBox ID.
type MACAddressWrite struct {
	Address            string         `json:"address,omitempty"`
	AssignedObjectType *string        `json:"assigned_object_type,omitempty"`
	AssignedObjectID   *int           `json:"assigned_object_id,omitempty"`
	Description        *string        `json:"description,omitempty"`
	Tags               []string       `json:"tags,omitempty"`
	CustomFields       map[string]any `json:"custom_fields,omitempty"`
}

// RIRWrite is the domain DTO for creating or updating a NetBox IPAM RIR. It
// carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". IsPrivate is a pointer so an update can
// set it to false explicitly.
type RIRWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	IsPrivate    *bool          `json:"is_private,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// AggregateWrite is the domain DTO for creating or updating a NetBox IPAM
// aggregate. It carries only the writable fields.
//
// Prefix is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". RIR and Tenant are numeric NetBox IDs.
type AggregateWrite struct {
	Prefix       string         `json:"prefix,omitempty"`
	RIR          *int           `json:"rir,omitempty"`
	Tenant       *int           `json:"tenant,omitempty"`
	DateAdded    *string        `json:"date_added,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// IPRangeWrite is the domain DTO for creating or updating a NetBox IPAM IP
// range. It carries only the writable fields.
//
// StartAddress and EndAddress are plain strings because both are required on
// create; their omitempty tags drop them from a PATCH body when not provided.
// Optional scalar fields are pointers so that a partial-update (PATCH) can
// distinguish "not provided" (nil) from "explicitly set to empty". VRF, Tenant
// and Role are numeric NetBox IDs; Status is a choice-string value.
type IPRangeWrite struct {
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

// ASNWrite is the domain DTO for creating or updating a NetBox IPAM ASN. It
// carries only the writable fields.
//
// ASN is a plain number because it is required on create. Optional scalar
// fields are pointers so that a partial-update (PATCH) can distinguish "not
// provided" (nil) from "explicitly set to empty". RIR and Tenant are numeric
// NetBox IDs.
type ASNWrite struct {
	ASN          int            `json:"asn,omitempty"`
	RIR          *int           `json:"rir,omitempty"`
	Tenant       *int           `json:"tenant,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// ASNRangeWrite is the domain DTO for creating or updating a NetBox IPAM ASN
// range. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". RIR and Tenant are numeric NetBox IDs.
type ASNRangeWrite struct {
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

// RouteTargetWrite is the domain DTO for creating or updating a NetBox IPAM
// route target. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". Tenant is a numeric NetBox ID.
type RouteTargetWrite struct {
	Name         string         `json:"name,omitempty"`
	Tenant       *int           `json:"tenant,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// ServiceWrite is the domain DTO for creating or updating a NetBox IPAM
// service. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". ParentObjectType and ParentObjectID
// identify the assigned parent object (e.g. a device or virtual machine);
// Protocol is a choice-string value; Ports is a list of port numbers.
type ServiceWrite struct {
	Name             string         `json:"name,omitempty"`
	ParentObjectType *string        `json:"parent_object_type,omitempty"`
	ParentObjectID   *int           `json:"parent_object_id,omitempty"`
	Protocol         *string        `json:"protocol,omitempty"`
	Ports            []int          `json:"ports,omitempty"`
	Description      *string        `json:"description,omitempty"`
	Tags             []string       `json:"tags,omitempty"`
	CustomFields     map[string]any `json:"custom_fields,omitempty"`
}

// ServiceTemplateWrite is the domain DTO for creating or updating a NetBox IPAM
// service template. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". Protocol is a choice-string value;
// Ports is a list of port numbers.
type ServiceTemplateWrite struct {
	Name         string         `json:"name,omitempty"`
	Protocol     *string        `json:"protocol,omitempty"`
	Ports        []int          `json:"ports,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// FHRPGroupWrite is the domain DTO for creating or updating a NetBox IPAM FHRP
// group. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty". Protocol is a choice-string value.
type FHRPGroupWrite struct {
	Name         string         `json:"name,omitempty"`
	Protocol     *string        `json:"protocol,omitempty"`
	GroupID      *int           `json:"group_id,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// FHRPGroupAssignmentWrite is the domain DTO for creating or updating a NetBox
// IPAM FHRP group assignment. It carries only the writable fields.
//
// Group and IPAddress are required on create (the numeric NetBox IDs of the FHRP
// group and interface). Priority is optional.
type FHRPGroupAssignmentWrite struct {
	Group         int    `json:"group,omitempty"`
	InterfaceType string `json:"interface_type,omitempty"`
	InterfaceID   int    `json:"interface_id,omitempty"`
	Priority      *int   `json:"priority,omitempty"`
}

// VLANTranslationPolicyWrite is the domain DTO for creating or updating a NetBox
// IPAM VLAN translation policy. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided.
type VLANTranslationPolicyWrite struct {
	Name         string         `json:"name,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// VLANTranslationRuleWrite is the domain DTO for creating or updating a NetBox
// IPAM VLAN translation rule. It carries only the writable fields.
//
// Policy is required on create (the numeric NetBox ID of the parent policy).
// LocalVID and RemoteVID are the translated VLAN identifiers (VLAN VIDs);
// Description is a free-text note.
type VLANTranslationRuleWrite struct {
	Policy      int     `json:"policy,omitempty"`
	LocalVID    *int    `json:"local_vid,omitempty"`
	RemoteVID   *int    `json:"remote_vid,omitempty"`
	Description *string `json:"description,omitempty"`
}

// ProviderAccountWrite is the domain DTO for creating or updating a NetBox
// circuits provider account. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Provider is the numeric NetBox
// ID of the provider. Optional scalar fields are pointers so that a
// partial-update (PATCH) can distinguish "not provided" (nil) from "explicitly
// set to empty".
type ProviderAccountWrite struct {
	Name         string         `json:"name,omitempty"`
	Provider     *int           `json:"provider,omitempty"`
	Account      *string        `json:"account,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// ProviderNetworkWrite is the domain DTO for creating or updating a NetBox
// circuits provider network. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Provider is the numeric NetBox
// ID of the provider. Optional scalar fields are pointers so that a
// partial-update (PATCH) can distinguish "not provided" (nil) from "explicitly
// set to empty".
type ProviderNetworkWrite struct {
	Name         string         `json:"name,omitempty"`
	Provider     *int           `json:"provider,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Comments     *string        `json:"comments,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// CircuitGroupWrite is the domain DTO for creating or updating a NetBox
// circuits circuit group. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty".
type CircuitGroupWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// CircuitGroupAssignmentWrite is the domain DTO for creating or updating a
// NetBox circuits circuit group assignment. It carries only the writable
// fields.
//
// Group is required on create (the numeric NetBox ID of the circuit group).
// MemberType is the NetBox content type string (e.g. "circuits.circuit" or
// "circuits.virtualcircuit") and MemberID is the numeric NetBox ID of the
// member. Priority is an optional choice string (primary/secondary/tertiary/
// inactive).
type CircuitGroupAssignmentWrite struct {
	Group      int     `json:"group,omitempty"`
	MemberType string  `json:"member_type,omitempty"`
	MemberID   int     `json:"member_id,omitempty"`
	Priority   *string `json:"priority,omitempty"`
}

// VirtualCircuitWrite is the domain DTO for creating or updating a NetBox
// circuits virtual circuit. It carries only the writable fields.
//
// CID is a plain string because it is required on create (the unique circuit
// ID); the omitempty tag drops it from a PATCH body when not provided. Type and
// ProviderNetwork are the numeric NetBox IDs of the related objects; Status is
// a choice-string value. Optional scalar fields are pointers so that a
// partial-update (PATCH) can distinguish "not provided" (nil) from "explicitly
// set to empty".
type VirtualCircuitWrite struct {
	CID             string         `json:"cid,omitempty"`
	Type            *int           `json:"type,omitempty"`
	ProviderNetwork *int           `json:"provider_network,omitempty"`
	Status          *string        `json:"status,omitempty"`
	Description     *string        `json:"description,omitempty"`
	Tags            []string       `json:"tags,omitempty"`
	CustomFields    map[string]any `json:"custom_fields,omitempty"`
}

// VirtualCircuitTypeWrite is the domain DTO for creating or updating a NetBox
// circuits virtual circuit type. It carries only the writable fields.
//
// Name is a plain string because it is required on create; the omitempty tag
// drops it from a PATCH body when not provided. Optional scalar fields are
// pointers so that a partial-update (PATCH) can distinguish "not provided"
// (nil) from "explicitly set to empty".
type VirtualCircuitTypeWrite struct {
	Name         string         `json:"name,omitempty"`
	Slug         *string        `json:"slug,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// VirtualCircuitTerminationWrite is the domain DTO for creating or updating a
// NetBox circuits virtual circuit termination. It carries only the writable
// fields.
//
// VirtualCircuit is required on create (the numeric NetBox ID of the parent
// virtual circuit). Role is a choice-string value; Interface is the numeric
// NetBox ID of the dcim.Interface the virtual circuit terminates on.
type VirtualCircuitTerminationWrite struct {
	VirtualCircuit int     `json:"virtual_circuit,omitempty"`
	Role           *string `json:"role,omitempty"`
	Interface      int     `json:"interface,omitempty"`
}

// PaginatedResponse wraps a paginated API response with the total count,
// navigation URLs and the current page's typed results.
type PaginatedResponse[T any] struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []T    `json:"results"`
}
