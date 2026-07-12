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

type (
	SitesInput struct {
		Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
		Region   string `json:"region,omitempty" jsonschema:"filter by region (slug or name)"`
		Status   string `json:"status,omitempty" jsonschema:"filter by status (active, planned, etc.)"`
		Tenant   string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug or name)"`
		Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
		PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
	}

	DevicesInput struct {
		Q            string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
		Site         string `json:"site,omitempty" jsonschema:"filter by site (slug)"`
		Role         string `json:"role,omitempty" jsonschema:"filter by device role (slug)"`
		Manufacturer string `json:"manufacturer,omitempty" jsonschema:"filter by manufacturer (slug)"`
		DeviceType   string `json:"device_type,omitempty" jsonschema:"filter by device type slug (e.g. c-1250)"`
		Status       string `json:"status,omitempty" jsonschema:"filter by status (active, planned, etc.)"`
		Name         string `json:"name,omitempty" jsonschema:"filter by name (partial match)"`
		Tenant       string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug)"`
		Rack         string `json:"rack,omitempty" jsonschema:"filter by rack (name)"`
		Cluster      string `json:"cluster,omitempty" jsonschema:"filter by cluster (name)"`
		Page         int    `json:"page,omitempty" jsonschema:"page number,default=1"`
		PageSize     int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
	}

	IPAddressesInput struct {
		Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
		Address  string `json:"address,omitempty" jsonschema:"filter by address (e.g. 192.168.1.0/24)"`
		Device   string `json:"device,omitempty" jsonschema:"filter by assigned device name"`
		Status   string `json:"status,omitempty" jsonschema:"filter by status (active, reserved, etc.)"`
		VRF      string `json:"vrf,omitempty" jsonschema:"filter by VRF (rd or name)"`
		Role     string `json:"role,omitempty" jsonschema:"filter by role (loopback, etc.)"`
		Tenant   string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug)"`
		Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
		PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
	}

	PrefixesInput struct {
		Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
		Prefix   string `json:"prefix,omitempty" jsonschema:"filter by prefix (e.g. 10.0.0.0/8)"`
		Site     string `json:"site,omitempty" jsonschema:"filter by site (slug)"`
		VRF      string `json:"vrf,omitempty" jsonschema:"filter by VRF (rd or name)"`
		Status   string `json:"status,omitempty" jsonschema:"filter by status (active, container, etc.)"`
		Role     string `json:"role,omitempty" jsonschema:"filter by role (slug)"`
		Tenant   string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug)"`
		Within   string `json:"within,omitempty" jsonschema:"find prefixes within a given prefix"`
		Family   int    `json:"family,omitempty" jsonschema:"address family: 4 or 6"`
		Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
		PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
	}

	VLANsInput struct {
		Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
		Site     string `json:"site,omitempty" jsonschema:"filter by site (slug)"`
		Group    string `json:"group,omitempty" jsonschema:"filter by VLAN group (slug)"`
		Status   string `json:"status,omitempty" jsonschema:"filter by status (active, reserved, etc.)"`
		Tenant   string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug)"`
		VID      int    `json:"vid,omitempty" jsonschema:"filter by VLAN ID"`
		Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
		PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
	}

	VirtualMachinesInput struct {
		Q            string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
		Cluster      string `json:"cluster,omitempty" jsonschema:"filter by cluster (name)"`
		ClusterGroup string `json:"cluster_group,omitempty" jsonschema:"filter by cluster group (slug)"`
		Role         string `json:"role,omitempty" jsonschema:"filter by VM role (slug)"`
		Status       string `json:"status,omitempty" jsonschema:"filter by status (active, staged, etc.)"`
		Tenant       string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug)"`
		Name         string `json:"name,omitempty" jsonschema:"filter by name (partial match)"`
		Site         string `json:"site,omitempty" jsonschema:"filter by site (slug)"`
		Page         int    `json:"page,omitempty" jsonschema:"page number,default=1"`
		PageSize     int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
	}

	ClustersInput struct {
		Q            string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
		ClusterType  string `json:"cluster_type,omitempty" jsonschema:"filter by cluster type (slug)"`
		ClusterGroup string `json:"cluster_group,omitempty" jsonschema:"filter by cluster group (slug)"`
		Site         string `json:"site,omitempty" jsonschema:"filter by site (slug)"`
		Tenant       string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug)"`
		Name         string `json:"name,omitempty" jsonschema:"filter by name (partial match)"`
		Page         int    `json:"page,omitempty" jsonschema:"page number,default=1"`
		PageSize     int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
	}

	CircuitsInput struct {
		Q           string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
		Provider    string `json:"provider,omitempty" jsonschema:"filter by provider (slug)"`
		CircuitType string `json:"circuit_type,omitempty" jsonschema:"filter by circuit type (slug)"`
		Site        string `json:"site,omitempty" jsonschema:"filter by site (slug)"`
		Status      string `json:"status,omitempty" jsonschema:"filter by status (active, planned, etc.)"`
		Tenant      string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug)"`
		Page        int    `json:"page,omitempty" jsonschema:"page number,default=1"`
		PageSize    int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
	}

	GetObjectInput struct {
		ObjectType string `json:"object_type" jsonschema:"object type: site, device, prefix, ip_address, vlan, virtual_machine, cluster, circuit, provider, tenant, rack, manufacturer, device_type, location, cluster_type, cluster_group, circuit_type, vrf, vlan_group, role, contact, cable,required"`
		ID         int    `json:"id" jsonschema:"numeric ID of the object,required"`
	}

	RacksInput struct {
		Q        string `json:"q,omitempty" jsonschema:"free-text search across all fields"`
		Site     string `json:"site,omitempty" jsonschema:"filter by site (slug)"`
		Location string `json:"location,omitempty" jsonschema:"filter by location (slug)"`
		Status   string `json:"status,omitempty" jsonschema:"filter by status (active, planned, etc.)"`
		Tenant   string `json:"tenant,omitempty" jsonschema:"filter by tenant (slug)"`
		Page     int    `json:"page,omitempty" jsonschema:"page number,default=1"`
		PageSize int    `json:"page_size,omitempty" jsonschema:"results per page (default: 25, max: 100)"`
	}
)

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

// — handler factories —

type SitesOutput struct {
	Count    int           `json:"count"`
	Next     string        `json:"next"`
	Previous string        `json:"previous"`
	Results  []domain.Site `json:"results"`
}

func NewGetSitesHandler(svc *application.NetworkService) mcp.ToolHandlerFor[SitesInput, SitesOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, in SitesInput) (*mcp.CallToolResult, SitesOutput, error) {
		params := paginationParams(in.Page, in.PageSize)
		addParam(params, "q", in.Q)
		addParam(params, "region", in.Region)
		addParam(params, "status", in.Status)
		addParam(params, "tenant", in.Tenant)

		resp, err := svc.ListSites(ctx, params)
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

		resp, err := svc.ListDevices(ctx, params)
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
		params := paginationParams(in.Page, in.PageSize)
		addParam(params, "q", in.Q)
		addParam(params, "address", in.Address)
		addParam(params, "device", in.Device)
		addParam(params, "status", in.Status)
		addParam(params, "vrf", in.VRF)
		addParam(params, "role", in.Role)
		addParam(params, "tenant", in.Tenant)

		resp, err := svc.ListIPAddresses(ctx, params)
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

		resp, err := svc.ListPrefixes(ctx, params)
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
		params := paginationParams(in.Page, in.PageSize)
		addParam(params, "q", in.Q)
		addParam(params, "site", in.Site)
		addParam(params, "group", in.Group)
		addParam(params, "status", in.Status)
		addParam(params, "tenant", in.Tenant)
		addIntParam(params, "vid", in.VID)

		resp, err := svc.ListVLANs(ctx, params)
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
		params := paginationParams(in.Page, in.PageSize)
		addParam(params, "q", in.Q)
		addParam(params, "cluster", in.Cluster)
		addParam(params, "cluster_group", in.ClusterGroup)
		addParam(params, "role", in.Role)
		addParam(params, "status", in.Status)
		addParam(params, "tenant", in.Tenant)
		addParam(params, "name__ic", in.Name)
		addParam(params, "site", in.Site)

		resp, err := svc.ListVirtualMachines(ctx, params)
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
		params := paginationParams(in.Page, in.PageSize)
		addParam(params, "q", in.Q)
		addParam(params, "type", in.ClusterType)
		addParam(params, "group", in.ClusterGroup)
		addParam(params, "site", in.Site)
		addParam(params, "tenant", in.Tenant)
		addParam(params, "name__ic", in.Name)

		resp, err := svc.ListClusters(ctx, params)
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
		params := paginationParams(in.Page, in.PageSize)
		addParam(params, "q", in.Q)
		addParam(params, "provider", in.Provider)
		addParam(params, "type", in.CircuitType)
		addParam(params, "site", in.Site)
		addParam(params, "status", in.Status)
		addParam(params, "tenant", in.Tenant)

		resp, err := svc.ListCircuits(ctx, params)
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
		if in.ObjectType == "" {
			return &mcp.CallToolResult{IsError: true}, GetObjectOutput{}, fmt.Errorf("object_type is required")
		}
		if in.ID <= 0 {
			return &mcp.CallToolResult{IsError: true}, GetObjectOutput{}, fmt.Errorf("id must be a positive integer")
		}

		data, err := svc.GetObject(ctx, in.ObjectType, in.ID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, GetObjectOutput{}, fmt.Errorf("get object: %w", err)
		}

		return &mcp.CallToolResult{}, GetObjectOutput{Data: data}, nil
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
		params := paginationParams(in.Page, in.PageSize)
		addParam(params, "q", in.Q)
		addParam(params, "site", in.Site)
		addParam(params, "location", in.Location)
		addParam(params, "status", in.Status)
		addParam(params, "tenant", in.Tenant)

		resp, err := svc.ListRacks(ctx, params)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, RacksOutput{}, fmt.Errorf("list racks: %w", err)
		}

		return &mcp.CallToolResult{}, RacksOutput{Count: resp.Count, Next: resp.Next, Previous: resp.Previous, Results: resp.Results}, nil
	}
}
