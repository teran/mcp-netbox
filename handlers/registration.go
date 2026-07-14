// Package handlers provides the HTTP transport layer for the MCP NetBox server.
//
// It implements the MCP Streamable HTTP protocol, including:
// - Middleware chain: recovery, security headers, host validation, rate limiting,
//   metrics, body limit, logging, token extraction, and service injection.
// - Tool handler factories that translate MCP tool calls into NetBox API queries.
// - Prometheus metrics collection for monitoring.
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

// RegisterTools registers all NetBox MCP tools on the given server.
func RegisterTools(s *mcp.Server, metrics *Metrics, svc *application.NetworkService) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_sites",
		Description: "List sites in NetBox with optional filters.",
	}, WrapToolHandler[SitesInput, PaginatedOutput[domain.Site]](metrics, "get_sites", func(ctx context.Context, req *mcp.CallToolRequest, in SitesInput) (*mcp.CallToolResult, PaginatedOutput[domain.Site], error) {
		return NewGetSitesHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_devices",
		Description: "List devices in NetBox with optional filters.",
	}, WrapToolHandler[DevicesInput, PaginatedOutput[domain.Device]](metrics, "get_devices", func(ctx context.Context, req *mcp.CallToolRequest, in DevicesInput) (*mcp.CallToolResult, PaginatedOutput[domain.Device], error) {
		return NewGetDevicesHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_ip_addresses",
		Description: "Search IP addresses in NetBox.",
	}, WrapToolHandler[IPAddressesInput, PaginatedOutput[domain.IPAddress]](metrics, "get_ip_addresses", func(ctx context.Context, req *mcp.CallToolRequest, in IPAddressesInput) (*mcp.CallToolResult, PaginatedOutput[domain.IPAddress], error) {
		return NewGetIPAddressesHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_prefixes",
		Description: "Search IP prefixes in NetBox.",
	}, WrapToolHandler[PrefixesInput, PaginatedOutput[domain.Prefix]](metrics, "get_prefixes", func(ctx context.Context, req *mcp.CallToolRequest, in PrefixesInput) (*mcp.CallToolResult, PaginatedOutput[domain.Prefix], error) {
		return NewGetPrefixesHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_vlans",
		Description: "List VLANs in NetBox.",
	}, WrapToolHandler[VLANsInput, PaginatedOutput[domain.VLAN]](metrics, "get_vlans", func(ctx context.Context, req *mcp.CallToolRequest, in VLANsInput) (*mcp.CallToolResult, PaginatedOutput[domain.VLAN], error) {
		return NewGetVLANsHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_virtual_machines",
		Description: "List virtual machines in NetBox.",
	}, WrapToolHandler[VirtualMachinesInput, PaginatedOutput[domain.VirtualMachine]](metrics, "get_virtual_machines", func(ctx context.Context, req *mcp.CallToolRequest, in VirtualMachinesInput) (*mcp.CallToolResult, PaginatedOutput[domain.VirtualMachine], error) {
		return NewGetVirtualMachinesHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_clusters",
		Description: "List clusters in NetBox.",
	}, WrapToolHandler[ClustersInput, PaginatedOutput[domain.Cluster]](metrics, "get_clusters", func(ctx context.Context, req *mcp.CallToolRequest, in ClustersInput) (*mcp.CallToolResult, PaginatedOutput[domain.Cluster], error) {
		return NewGetClustersHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_circuits",
		Description: "List circuits in NetBox.",
	}, WrapToolHandler[CircuitsInput, PaginatedOutput[domain.Circuit]](metrics, "get_circuits", func(ctx context.Context, req *mcp.CallToolRequest, in CircuitsInput) (*mcp.CallToolResult, PaginatedOutput[domain.Circuit], error) {
		return NewGetCircuitsHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_object_by_id",
		Description: "Retrieve any NetBox object by its type and numeric ID.",
	}, WrapToolHandler[GetObjectInput, GetObjectOutput](metrics, "get_object_by_id", func(ctx context.Context, req *mcp.CallToolRequest, in GetObjectInput) (*mcp.CallToolResult, GetObjectOutput, error) {
		return NewGetObjectByIDHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_racks",
		Description: "List racks in NetBox.",
	}, WrapToolHandler[RacksInput, PaginatedOutput[domain.Rack]](metrics, "get_racks", func(ctx context.Context, req *mcp.CallToolRequest, in RacksInput) (*mcp.CallToolResult, PaginatedOutput[domain.Rack], error) {
		return NewGetRacksHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_interfaces",
		Description: "List device interfaces in NetBox with optional filters.",
	}, WrapToolHandler[InterfacesInput, PaginatedOutput[domain.Interface]](metrics, "get_interfaces", func(ctx context.Context, req *mcp.CallToolRequest, in InterfacesInput) (*mcp.CallToolResult, PaginatedOutput[domain.Interface], error) {
		return NewGetInterfacesHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_circuit_terminations",
		Description: "List circuit terminations in NetBox with optional filters.",
	}, WrapToolHandler[CircuitTerminationsInput, PaginatedOutput[domain.CircuitTermination]](metrics, "get_circuit_terminations", func(ctx context.Context, req *mcp.CallToolRequest, in CircuitTerminationsInput) (*mcp.CallToolResult, PaginatedOutput[domain.CircuitTermination], error) {
		return NewGetCircuitTerminationsHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_cables",
		Description: "List cables in NetBox with optional filters.",
	}, WrapToolHandler[CablesInput, PaginatedOutput[domain.Cable]](metrics, "get_cables", func(ctx context.Context, req *mcp.CallToolRequest, in CablesInput) (*mcp.CallToolResult, PaginatedOutput[domain.Cable], error) {
		return NewGetCablesHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_vm_interfaces",
		Description: "List VM interfaces in NetBox with optional filters.",
	}, WrapToolHandler[VMInterfacesInput, PaginatedOutput[domain.VMInterface]](metrics, "get_vm_interfaces", func(ctx context.Context, req *mcp.CallToolRequest, in VMInterfacesInput) (*mcp.CallToolResult, PaginatedOutput[domain.VMInterface], error) {
		return NewGetVMInterfacesHandler(svc)(ctx, req, in)
	}))
}
