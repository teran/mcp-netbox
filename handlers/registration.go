// Package handlers provides the HTTP transport layer for the MCP NetBox server.
//
// It implements the MCP Streamable HTTP protocol, including:
//   - Middleware chain: recovery, security headers, host validation, rate limiting,
//     metrics, body limit, logging, token extraction, and service injection.
//   - Tool handler factories that translate MCP tool calls into NetBox API queries.
//   - Prometheus metrics collection for monitoring.
//
// The middleware chain is assembled in NewMux() (server.go) and applies to
// all requests on the /mcp path. The healthz and readyz endpoints bypass
// authentication middleware for health-check purposes.
package handlers

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/domain"
)

// boolPtr returns a pointer to b. ToolAnnotations uses *bool for
// DestructiveHint and OpenWorldHint because they default to true.
func boolPtr(b bool) *bool {
	return &b
}

// readOnlyTool returns tool annotations for one of the read-only NetBox query
// tools. Every tool in this server is read-only and operates on a closed
// domain (the NetBox inventory), so these hints are constant across tools.
func readOnlyTool(title string) *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{
		Title:           title,
		ReadOnlyHint:    true,
		IdempotentHint:  true,
		DestructiveHint: boolPtr(false),
		OpenWorldHint:   boolPtr(false),
	}
}

// RegisterTools registers all NetBox MCP tools on the given server.
func RegisterTools(s *mcp.Server, metrics *Metrics, svc *application.NetworkService) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_sites",
		Description: "List sites in NetBox with optional filters. Use q for free-text search, or filter by region, status, tenant, or tag. Results are paginated: set page and page_size (max 1000). An empty results array means no sites match the filters, not an error.",
		Annotations: readOnlyTool("List Sites"),
	}, WrapToolHandler[SitesInput, PaginatedOutput[domain.Site]](metrics, "get_sites", func(ctx context.Context, req *mcp.CallToolRequest, in SitesInput) (*mcp.CallToolResult, PaginatedOutput[domain.Site], error) {
		return NewGetSitesHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_devices",
		Description: "List devices in NetBox with optional filters. Filter by site, role, manufacturer, device_type, status, name, tenant, rack, cluster, or tag, or use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no devices match.",
		Annotations: readOnlyTool("List Devices"),
	}, WrapToolHandler[DevicesInput, PaginatedOutput[domain.Device]](metrics, "get_devices", func(ctx context.Context, req *mcp.CallToolRequest, in DevicesInput) (*mcp.CallToolResult, PaginatedOutput[domain.Device], error) {
		return NewGetDevicesHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_ip_addresses",
		Description: "Search IP addresses in NetBox. Filter by address (e.g. 192.168.1.0/24), assigned device, status, VRF, role, tenant, or tag, or use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no addresses match.",
		Annotations: readOnlyTool("List IP Addresses"),
	}, WrapToolHandler[IPAddressesInput, PaginatedOutput[domain.IPAddress]](metrics, "get_ip_addresses", func(ctx context.Context, req *mcp.CallToolRequest, in IPAddressesInput) (*mcp.CallToolResult, PaginatedOutput[domain.IPAddress], error) {
		return NewGetIPAddressesHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_prefixes",
		Description: "Search IP prefixes in NetBox. Filter by prefix (e.g. 10.0.0.0/8), site, VRF, status, role, tenant, family, or tag, or use within to find prefixes nested inside a given prefix. Results are paginated: set page and page_size (max 1000). An empty results array means no prefixes match.",
		Annotations: readOnlyTool("List Prefixes"),
	}, WrapToolHandler[PrefixesInput, PaginatedOutput[domain.Prefix]](metrics, "get_prefixes", func(ctx context.Context, req *mcp.CallToolRequest, in PrefixesInput) (*mcp.CallToolResult, PaginatedOutput[domain.Prefix], error) {
		return NewGetPrefixesHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_vlans",
		Description: "List VLANs in NetBox. Filter by site, group, status, tenant, or tag, or by a specific VLAN ID (vid). Use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no VLANs match.",
		Annotations: readOnlyTool("List VLANs"),
	}, WrapToolHandler[VLANsInput, PaginatedOutput[domain.VLAN]](metrics, "get_vlans", func(ctx context.Context, req *mcp.CallToolRequest, in VLANsInput) (*mcp.CallToolResult, PaginatedOutput[domain.VLAN], error) {
		return NewGetVLANsHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_virtual_machines",
		Description: "List virtual machines in NetBox. Filter by cluster, cluster_group, role, status, tenant, name, or site, or use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no VMs match.",
		Annotations: readOnlyTool("List Virtual Machines"),
	}, WrapToolHandler[VirtualMachinesInput, PaginatedOutput[domain.VirtualMachine]](metrics, "get_virtual_machines", func(ctx context.Context, req *mcp.CallToolRequest, in VirtualMachinesInput) (*mcp.CallToolResult, PaginatedOutput[domain.VirtualMachine], error) {
		return NewGetVirtualMachinesHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_clusters",
		Description: "List clusters in NetBox. Filter by cluster_type, cluster_group, site, tenant, or name, or use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no clusters match.",
		Annotations: readOnlyTool("List Clusters"),
	}, WrapToolHandler[ClustersInput, PaginatedOutput[domain.Cluster]](metrics, "get_clusters", func(ctx context.Context, req *mcp.CallToolRequest, in ClustersInput) (*mcp.CallToolResult, PaginatedOutput[domain.Cluster], error) {
		return NewGetClustersHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_circuits",
		Description: "List circuits in NetBox. Filter by provider, circuit_type, site, status, or tenant, or use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no circuits match.",
		Annotations: readOnlyTool("List Circuits"),
	}, WrapToolHandler[CircuitsInput, PaginatedOutput[domain.Circuit]](metrics, "get_circuits", func(ctx context.Context, req *mcp.CallToolRequest, in CircuitsInput) (*mcp.CallToolResult, PaginatedOutput[domain.Circuit], error) {
		return NewGetCircuitsHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_object_by_id",
		Description: "Retrieve any single NetBox object by its type and numeric ID. object_type is required (e.g. site, device, prefix, ip_address, vlan, virtual_machine, cluster, circuit, provider, tenant, rack, cable, interface). Returns a 404-style error if the object does not exist or the token lacks permission.",
		Annotations: readOnlyTool("Get Object by ID"),
	}, WrapToolHandler[GetObjectInput, GetObjectOutput](metrics, "get_object_by_id", func(ctx context.Context, req *mcp.CallToolRequest, in GetObjectInput) (*mcp.CallToolResult, GetObjectOutput, error) {
		return NewGetObjectByIDHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_racks",
		Description: "List racks in NetBox. Filter by site, location, status, or tenant, or use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no racks match.",
		Annotations: readOnlyTool("List Racks"),
	}, WrapToolHandler[RacksInput, PaginatedOutput[domain.Rack]](metrics, "get_racks", func(ctx context.Context, req *mcp.CallToolRequest, in RacksInput) (*mcp.CallToolResult, PaginatedOutput[domain.Rack], error) {
		return NewGetRacksHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_interfaces",
		Description: "List device interfaces in NetBox. Filter by device, type, enabled, or name, or use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no interfaces match.",
		Annotations: readOnlyTool("List Interfaces"),
	}, WrapToolHandler[InterfacesInput, PaginatedOutput[domain.Interface]](metrics, "get_interfaces", func(ctx context.Context, req *mcp.CallToolRequest, in InterfacesInput) (*mcp.CallToolResult, PaginatedOutput[domain.Interface], error) {
		return NewGetInterfacesHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_circuit_terminations",
		Description: "List circuit terminations in NetBox. Filter by circuit, site, term_side (A or Z), or tag, or use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no terminations match.",
		Annotations: readOnlyTool("List Circuit Terminations"),
	}, WrapToolHandler[CircuitTerminationsInput, PaginatedOutput[domain.CircuitTermination]](metrics, "get_circuit_terminations", func(ctx context.Context, req *mcp.CallToolRequest, in CircuitTerminationsInput) (*mcp.CallToolResult, PaginatedOutput[domain.CircuitTermination], error) {
		return NewGetCircuitTerminationsHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_cables",
		Description: "List cables in NetBox. Filter by type, status, site, color, or label, or use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no cables match.",
		Annotations: readOnlyTool("List Cables"),
	}, WrapToolHandler[CablesInput, PaginatedOutput[domain.Cable]](metrics, "get_cables", func(ctx context.Context, req *mcp.CallToolRequest, in CablesInput) (*mcp.CallToolResult, PaginatedOutput[domain.Cable], error) {
		return NewGetCablesHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_vm_interfaces",
		Description: "List VM interfaces in NetBox. Filter by virtual_machine, name, or tag, or use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no VM interfaces match.",
		Annotations: readOnlyTool("List VM Interfaces"),
	}, WrapToolHandler[VMInterfacesInput, PaginatedOutput[domain.VMInterface]](metrics, "get_vm_interfaces", func(ctx context.Context, req *mcp.CallToolRequest, in VMInterfacesInput) (*mcp.CallToolResult, PaginatedOutput[domain.VMInterface], error) {
		return NewGetVMInterfacesHandler(svc)(ctx, req, in)
	}))
}
