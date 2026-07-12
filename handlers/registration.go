package handlers

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/teran/mcp-netbox/application"
)

// RegisterTools registers all NetBox MCP tools on the given server.
func RegisterTools(s *mcp.Server, metrics *Metrics, svc *application.NetworkService) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_sites",
		Description: "List sites in NetBox with optional filters.",
	}, WrapToolHandler[SitesInput, SitesOutput](metrics, "get_sites", func(ctx context.Context, req *mcp.CallToolRequest, in SitesInput) (*mcp.CallToolResult, SitesOutput, error) {
		return NewGetSitesHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_devices",
		Description: "List devices in NetBox with optional filters.",
	}, WrapToolHandler[DevicesInput, DevicesOutput](metrics, "get_devices", func(ctx context.Context, req *mcp.CallToolRequest, in DevicesInput) (*mcp.CallToolResult, DevicesOutput, error) {
		return NewGetDevicesHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_ip_addresses",
		Description: "Search IP addresses in NetBox.",
	}, WrapToolHandler[IPAddressesInput, IPAddressesOutput](metrics, "get_ip_addresses", func(ctx context.Context, req *mcp.CallToolRequest, in IPAddressesInput) (*mcp.CallToolResult, IPAddressesOutput, error) {
		return NewGetIPAddressesHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_prefixes",
		Description: "Search IP prefixes in NetBox.",
	}, WrapToolHandler[PrefixesInput, PrefixesOutput](metrics, "get_prefixes", func(ctx context.Context, req *mcp.CallToolRequest, in PrefixesInput) (*mcp.CallToolResult, PrefixesOutput, error) {
		return NewGetPrefixesHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_vlans",
		Description: "List VLANs in NetBox.",
	}, WrapToolHandler[VLANsInput, VLANsOutput](metrics, "get_vlans", func(ctx context.Context, req *mcp.CallToolRequest, in VLANsInput) (*mcp.CallToolResult, VLANsOutput, error) {
		return NewGetVLANsHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_virtual_machines",
		Description: "List virtual machines in NetBox.",
	}, WrapToolHandler[VirtualMachinesInput, VirtualMachinesOutput](metrics, "get_virtual_machines", func(ctx context.Context, req *mcp.CallToolRequest, in VirtualMachinesInput) (*mcp.CallToolResult, VirtualMachinesOutput, error) {
		return NewGetVirtualMachinesHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_clusters",
		Description: "List clusters in NetBox.",
	}, WrapToolHandler[ClustersInput, ClustersOutput](metrics, "get_clusters", func(ctx context.Context, req *mcp.CallToolRequest, in ClustersInput) (*mcp.CallToolResult, ClustersOutput, error) {
		return NewGetClustersHandler(svc)(ctx, req, in)
	}))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_circuits",
		Description: "List circuits in NetBox.",
	}, WrapToolHandler[CircuitsInput, CircuitsOutput](metrics, "get_circuits", func(ctx context.Context, req *mcp.CallToolRequest, in CircuitsInput) (*mcp.CallToolResult, CircuitsOutput, error) {
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
	}, WrapToolHandler[RacksInput, RacksOutput](metrics, "get_racks", func(ctx context.Context, req *mcp.CallToolRequest, in RacksInput) (*mcp.CallToolResult, RacksOutput, error) {
		return NewGetRacksHandler(svc)(ctx, req, in)
	}))
}
