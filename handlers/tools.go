package handlers

import (
	"context"
	"errors"
	"fmt"
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

// — helpers —

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
