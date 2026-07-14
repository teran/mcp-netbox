package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/domain"
)

// — input/output types —

// SitesInput represents the input fields for the get_sites tool.
type SitesInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Region   string `json:"region,omitempty" jsonschema:"filter by region (slug or name)"`
	Status   string `json:"status,omitempty" jsonschema:"status: active, planned, staged, retired, decommissioning"`
	Tenant   string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug or name)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
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
	PageSize     int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
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
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
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
	Family   int    `json:"family,omitempty" jsonschema:"address family: 4 or 6"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
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
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
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
	PageSize     int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
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
	PageSize     int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
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
	PageSize    int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
}

// GetObjectInput represents the input fields for the get_object_by_id tool.
type GetObjectInput struct {
	ObjectType string `json:"object_type" jsonschema:"object type: site, device, prefix, ip_address, vlan, virtual_machine, cluster, circuit, provider, tenant, rack, manufacturer, device_type, location, cluster_type, cluster_group, circuit_type, vrf, vlan_group, role, contact, cable, interface, vm_interface, circuit_termination,required"`
	ID         int    `json:"id" jsonschema:"numeric ID of the object (positive integer),required"`
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
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
}

// VMInterfacesInput represents the input fields for the get_vm_interfaces tool.
type VMInterfacesInput struct {
	Q              string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	VirtualMachine string `json:"virtual_machine,omitempty" jsonschema:"filter by virtual machine (name)"`
	Name           string `json:"name,omitempty" jsonschema:"filter by name (case-insensitive partial match)"`
	Tag            string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	Page           int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize       int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
}

// CircuitTerminationsInput represents the input fields for the get_circuit_terminations tool.
type CircuitTerminationsInput struct {
	Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
	Circuit  string `json:"circuit,omitempty" jsonschema:"filter by circuit (ID or CID)"`
	Site     string `json:"site,omitempty" jsonschema:"filter by site (slug)"`
	Tag      string `json:"tag,omitempty" jsonschema:"filter by tag (slug)"`
	TermSide string `json:"term_side,omitempty" jsonschema:"filter by termination side: A or Z"`
	Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
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
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
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
	PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
}

// — helpers —

func paginationParams(page, pageSize int) map[string]string {
	params := make(map[string]string)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 25
	}
	params["offset"] = strconv.Itoa((page - 1) * pageSize)
	params["limit"] = strconv.Itoa(min(pageSize, 100))
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

// — handler factories —

type SitesOutput struct {
	Count    int           `json:"count"`
	Next     string        `json:"next"`
	Previous string        `json:"previous"`
	Results  []domain.Site `json:"results"`
}

func NewGetSitesHandler(svc *application.NetworkService) mcp.ToolHandlerFor[SitesInput, SitesOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in SitesInput) (*mcp.CallToolResult, SitesOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, SitesOutput{}, errServiceNotAvailable
		}

		params := paginationParams(in.Page, in.PageSize)
		addParam(params, "q", in.Q)
		addParam(params, "region", in.Region)
		addParam(params, "status", in.Status)
		addParam(params, "tenant", in.Tenant)
		addParam(params, "tag", in.Tag)

		resp, err := s.ListSites(ctx, params)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, SitesOutput{}, fmt.Errorf("list sites: %w", err)
		}

		return &mcp.CallToolResult{}, SitesOutput{Count: resp.Count, Next: resp.Next, Previous: resp.Previous, Results: resp.Results}, nil
	}
}

type DevicesOutput struct {
	Count    int             `json:"count"`
	Next     string          `json:"next"`
	Previous string          `json:"previous"`
	Results  []domain.Device `json:"results"`
}

func NewGetDevicesHandler(svc *application.NetworkService) mcp.ToolHandlerFor[DevicesInput, DevicesOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in DevicesInput) (*mcp.CallToolResult, DevicesOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, DevicesOutput{}, errServiceNotAvailable
		}

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

		resp, err := s.ListDevices(ctx, params)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, DevicesOutput{}, fmt.Errorf("list devices: %w", err)
		}

		return &mcp.CallToolResult{}, DevicesOutput{Count: resp.Count, Next: resp.Next, Previous: resp.Previous, Results: resp.Results}, nil
	}
}

type IPAddressesOutput struct {
	Count    int                `json:"count"`
	Next     string             `json:"next"`
	Previous string             `json:"previous"`
	Results  []domain.IPAddress `json:"results"`
}

func NewGetIPAddressesHandler(svc *application.NetworkService) mcp.ToolHandlerFor[IPAddressesInput, IPAddressesOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in IPAddressesInput) (*mcp.CallToolResult, IPAddressesOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, IPAddressesOutput{}, errServiceNotAvailable
		}

		params := paginationParams(in.Page, in.PageSize)
		addParam(params, "q", in.Q)
		addParam(params, "address", in.Address)
		addParam(params, "device", in.Device)
		addParam(params, "status", in.Status)
		addParam(params, "vrf", in.VRF)
		addParam(params, "role", in.Role)
		addParam(params, "tenant", in.Tenant)
		addParam(params, "tag", in.Tag)

		resp, err := s.ListIPAddresses(ctx, params)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, IPAddressesOutput{}, fmt.Errorf("list IP addresses: %w", err)
		}

		return &mcp.CallToolResult{}, IPAddressesOutput{Count: resp.Count, Next: resp.Next, Previous: resp.Previous, Results: resp.Results}, nil
	}
}

type PrefixesOutput struct {
	Count    int             `json:"count"`
	Next     string          `json:"next"`
	Previous string          `json:"previous"`
	Results  []domain.Prefix `json:"results"`
}

func NewGetPrefixesHandler(svc *application.NetworkService) mcp.ToolHandlerFor[PrefixesInput, PrefixesOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in PrefixesInput) (*mcp.CallToolResult, PrefixesOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, PrefixesOutput{}, errServiceNotAvailable
		}

		params := paginationParams(in.Page, in.PageSize)
		addParam(params, "q", in.Q)
		addParam(params, "prefix", in.Prefix)
		addParam(params, "site", in.Site)
		addParam(params, "vrf", in.VRF)
		addParam(params, "status", in.Status)
		addParam(params, "role", in.Role)
		addParam(params, "tenant", in.Tenant)
		addParam(params, "within", in.Within)
		addIntParam(params, "family", in.Family)
		addParam(params, "tag", in.Tag)

		resp, err := s.ListPrefixes(ctx, params)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, PrefixesOutput{}, fmt.Errorf("list prefixes: %w", err)
		}

		return &mcp.CallToolResult{}, PrefixesOutput{Count: resp.Count, Next: resp.Next, Previous: resp.Previous, Results: resp.Results}, nil
	}
}

type VLANsOutput struct {
	Count    int           `json:"count"`
	Next     string        `json:"next"`
	Previous string        `json:"previous"`
	Results  []domain.VLAN `json:"results"`
}

func NewGetVLANsHandler(svc *application.NetworkService) mcp.ToolHandlerFor[VLANsInput, VLANsOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in VLANsInput) (*mcp.CallToolResult, VLANsOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, VLANsOutput{}, errServiceNotAvailable
		}

		params := paginationParams(in.Page, in.PageSize)
		addParam(params, "q", in.Q)
		addParam(params, "site", in.Site)
		addParam(params, "group", in.Group)
		addParam(params, "status", in.Status)
		addParam(params, "tenant", in.Tenant)
		addIntParam(params, "vid", in.VID)
		addParam(params, "tag", in.Tag)

		resp, err := s.ListVLANs(ctx, params)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, VLANsOutput{}, fmt.Errorf("list VLANs: %w", err)
		}

		return &mcp.CallToolResult{}, VLANsOutput{Count: resp.Count, Next: resp.Next, Previous: resp.Previous, Results: resp.Results}, nil
	}
}

type VirtualMachinesOutput struct {
	Count    int                     `json:"count"`
	Next     string                  `json:"next"`
	Previous string                  `json:"previous"`
	Results  []domain.VirtualMachine `json:"results"`
}

func NewGetVirtualMachinesHandler(svc *application.NetworkService) mcp.ToolHandlerFor[VirtualMachinesInput, VirtualMachinesOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in VirtualMachinesInput) (*mcp.CallToolResult, VirtualMachinesOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, VirtualMachinesOutput{}, errServiceNotAvailable
		}

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

		resp, err := s.ListVirtualMachines(ctx, params)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, VirtualMachinesOutput{}, fmt.Errorf("list VMs: %w", err)
		}

		return &mcp.CallToolResult{}, VirtualMachinesOutput{Count: resp.Count, Next: resp.Next, Previous: resp.Previous, Results: resp.Results}, nil
	}
}

type ClustersOutput struct {
	Count    int              `json:"count"`
	Next     string           `json:"next"`
	Previous string           `json:"previous"`
	Results  []domain.Cluster `json:"results"`
}

func NewGetClustersHandler(svc *application.NetworkService) mcp.ToolHandlerFor[ClustersInput, ClustersOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in ClustersInput) (*mcp.CallToolResult, ClustersOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, ClustersOutput{}, errServiceNotAvailable
		}

		params := paginationParams(in.Page, in.PageSize)
		addParam(params, "q", in.Q)
		addParam(params, "type", in.ClusterType)
		addParam(params, "group", in.ClusterGroup)
		addParam(params, "site", in.Site)
		addParam(params, "tenant", in.Tenant)
		addParam(params, "name__ic", in.Name)
		addParam(params, "tag", in.Tag)

		resp, err := s.ListClusters(ctx, params)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ClustersOutput{}, fmt.Errorf("list clusters: %w", err)
		}

		return &mcp.CallToolResult{}, ClustersOutput{Count: resp.Count, Next: resp.Next, Previous: resp.Previous, Results: resp.Results}, nil
	}
}

type CircuitsOutput struct {
	Count    int              `json:"count"`
	Next     string           `json:"next"`
	Previous string           `json:"previous"`
	Results  []domain.Circuit `json:"results"`
}

func NewGetCircuitsHandler(svc *application.NetworkService) mcp.ToolHandlerFor[CircuitsInput, CircuitsOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in CircuitsInput) (*mcp.CallToolResult, CircuitsOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitsOutput{}, errServiceNotAvailable
		}

		params := paginationParams(in.Page, in.PageSize)
		addParam(params, "q", in.Q)
		addParam(params, "provider", in.Provider)
		addParam(params, "type", in.CircuitType)
		addParam(params, "site", in.Site)
		addParam(params, "status", in.Status)
		addParam(params, "tenant", in.Tenant)
		addParam(params, "tag", in.Tag)

		resp, err := s.ListCircuits(ctx, params)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CircuitsOutput{}, fmt.Errorf("list circuits: %w", err)
		}

		return &mcp.CallToolResult{}, CircuitsOutput{Count: resp.Count, Next: resp.Next, Previous: resp.Previous, Results: resp.Results}, nil
	}
}

type GetObjectOutput struct {
	Data domain.RawObject `json:"data"`
}

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

		data, err := s.GetObject(ctx, in.ObjectType, in.ID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, GetObjectOutput{}, fmt.Errorf("get object: %w", err)
		}

		return &mcp.CallToolResult{}, GetObjectOutput{Data: data}, nil
	}
}

type CircuitTerminationsOutput struct {
	Count    int                         `json:"count"`
	Next     string                      `json:"next"`
	Previous string                      `json:"previous"`
	Results  []domain.CircuitTermination `json:"results"`
}

func NewGetCircuitTerminationsHandler(svc *application.NetworkService) mcp.ToolHandlerFor[CircuitTerminationsInput, CircuitTerminationsOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in CircuitTerminationsInput) (*mcp.CallToolResult, CircuitTerminationsOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, CircuitTerminationsOutput{}, errServiceNotAvailable
		}

		params := paginationParams(in.Page, in.PageSize)
		addParam(params, "q", in.Q)
		addParam(params, "circuit", in.Circuit)
		addParam(params, "site", in.Site)
		addParam(params, "tag", in.Tag)
		addParam(params, "term_side", in.TermSide)

		resp, err := s.ListCircuitTerminations(ctx, params)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CircuitTerminationsOutput{}, fmt.Errorf("list circuit terminations: %w", err)
		}

		return &mcp.CallToolResult{}, CircuitTerminationsOutput{Count: resp.Count, Next: resp.Next, Previous: resp.Previous, Results: resp.Results}, nil
	}
}

type CablesOutput struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []domain.Cable `json:"results"`
}

func NewGetCablesHandler(svc *application.NetworkService) mcp.ToolHandlerFor[CablesInput, CablesOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in CablesInput) (*mcp.CallToolResult, CablesOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, CablesOutput{}, errServiceNotAvailable
		}

		params := paginationParams(in.Page, in.PageSize)
		addParam(params, "q", in.Q)
		addParam(params, "type", in.Type)
		addParam(params, "status", in.Status)
		addParam(params, "site", in.Site)
		addParam(params, "tag", in.Tag)
		addParam(params, "color", in.Color)
		addParam(params, "label__ic", in.Label)

		resp, err := s.ListCables(ctx, params)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CablesOutput{}, fmt.Errorf("list cables: %w", err)
		}

		return &mcp.CallToolResult{}, CablesOutput{Count: resp.Count, Next: resp.Next, Previous: resp.Previous, Results: resp.Results}, nil
	}
}

type RacksOutput struct {
	Count    int           `json:"count"`
	Next     string        `json:"next"`
	Previous string        `json:"previous"`
	Results  []domain.Rack `json:"results"`
}

func NewGetRacksHandler(svc *application.NetworkService) mcp.ToolHandlerFor[RacksInput, RacksOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in RacksInput) (*mcp.CallToolResult, RacksOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, RacksOutput{}, errServiceNotAvailable
		}

		params := paginationParams(in.Page, in.PageSize)
		addParam(params, "q", in.Q)
		addParam(params, "site", in.Site)
		addParam(params, "location", in.Location)
		addParam(params, "status", in.Status)
		addParam(params, "tenant", in.Tenant)
		addParam(params, "tag", in.Tag)

		resp, err := s.ListRacks(ctx, params)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, RacksOutput{}, fmt.Errorf("list racks: %w", err)
		}

		return &mcp.CallToolResult{}, RacksOutput{Count: resp.Count, Next: resp.Next, Previous: resp.Previous, Results: resp.Results}, nil
	}
}

type InterfacesOutput struct {
	Count    int                `json:"count"`
	Next     string             `json:"next"`
	Previous string             `json:"previous"`
	Results  []domain.Interface `json:"results"`
}

func NewGetInterfacesHandler(svc *application.NetworkService) mcp.ToolHandlerFor[InterfacesInput, InterfacesOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in InterfacesInput) (*mcp.CallToolResult, InterfacesOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, InterfacesOutput{}, errServiceNotAvailable
		}

		params := paginationParams(in.Page, in.PageSize)
		addParam(params, "q", in.Q)
		addParam(params, "device", in.Device)
		addParam(params, "type", in.Type)
		addParam(params, "name__ic", in.Name)
		addParam(params, "tag", in.Tag)
		if in.Enabled != nil {
			addParam(params, "enabled", fmt.Sprintf("%t", *in.Enabled))
		}

		resp, err := s.ListInterfaces(ctx, params)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, InterfacesOutput{}, fmt.Errorf("list interfaces: %w", err)
		}

		return &mcp.CallToolResult{}, InterfacesOutput{Count: resp.Count, Next: resp.Next, Previous: resp.Previous, Results: resp.Results}, nil
	}
}

type VMInterfacesOutput struct {
	Count    int                  `json:"count"`
	Next     string               `json:"next"`
	Previous string               `json:"previous"`
	Results  []domain.VMInterface `json:"results"`
}

func NewGetVMInterfacesHandler(svc *application.NetworkService) mcp.ToolHandlerFor[VMInterfacesInput, VMInterfacesOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in VMInterfacesInput) (*mcp.CallToolResult, VMInterfacesOutput, error) {
		s := resolveService(ctx, svc)
		if s == nil {
			return &mcp.CallToolResult{IsError: true}, VMInterfacesOutput{}, errServiceNotAvailable
		}

		params := paginationParams(in.Page, in.PageSize)
		addParam(params, "q", in.Q)
		addParam(params, "virtual_machine", in.VirtualMachine)
		addParam(params, "name__ic", in.Name)
		addParam(params, "tag", in.Tag)

		resp, err := s.ListVMInterfaces(ctx, params)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, VMInterfacesOutput{}, fmt.Errorf("list VM interfaces: %w", err)
		}

		return &mcp.CallToolResult{}, VMInterfacesOutput{Count: resp.Count, Next: resp.Next, Previous: resp.Previous, Results: resp.Results}, nil
	}
}
