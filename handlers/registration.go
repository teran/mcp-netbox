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

// writeTool returns tool annotations for a non-destructive write tool
// (create/update). These are not idempotent by default (a create produces a new
// object each call) and not read-only.
func writeTool(title string) *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{
		Title:           title,
		ReadOnlyHint:    false,
		IdempotentHint:  false,
		DestructiveHint: boolPtr(false),
		OpenWorldHint:   boolPtr(false),
	}
}

// destructiveTool returns tool annotations for a destructive tool (delete).
// Deletion is idempotent (deleting a non-existent object is a no-op from the
// caller's perspective) and flagged as destructive so clients can guard it.
func destructiveTool(title string) *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{
		Title:           title,
		ReadOnlyHint:    false,
		IdempotentHint:  true,
		DestructiveHint: boolPtr(true),
		OpenWorldHint:   boolPtr(false),
	}
}

// toolDefIndex returns a map of tool name -> definition for lookup during
// registration, so each tool's Description carries its per-tool Instructions
// (M4) through the SDK's only per-tool guidance channel.
func toolDefIndex() map[string]toolDef {
	idx := make(map[string]toolDef, len(toolDefs()))
	for _, d := range toolDefs() {
		idx[d.Name] = d
	}
	return idx
}

// registerTool is a helper that builds and registers a single read-only tool,
// combining its description with the per-tool instructions.
func registerTool[I, O any](s *mcp.Server, def toolDef, metrics *Metrics, name string, handler mcp.ToolHandlerFor[I, O]) {
	description := def.Description
	if def.Instructions != "" {
		description = description + "\n\nInstructions: " + def.Instructions
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        name,
		Description: description,
		Annotations: readOnlyTool(def.Title),
	}, handler)
}

// registerWriteTool registers a non-destructive write tool with writeTool
// annotations.
func registerWriteTool[I, O any](s *mcp.Server, def toolDef, metrics *Metrics, name string, handler mcp.ToolHandlerFor[I, O]) {
	description := def.Description
	if def.Instructions != "" {
		description = description + "\n\nInstructions: " + def.Instructions
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        name,
		Description: description,
		Annotations: writeTool(def.Title),
	}, handler)
}

// registerDeleteTool registers a destructive tool with destructiveTool
// annotations.
func registerDeleteTool[I, O any](s *mcp.Server, def toolDef, metrics *Metrics, name string, handler mcp.ToolHandlerFor[I, O]) {
	description := def.Description
	if def.Instructions != "" {
		description = description + "\n\nInstructions: " + def.Instructions
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        name,
		Description: description,
		Annotations: destructiveTool(def.Title),
	}, handler)
}

// RegisterTools registers all NetBox MCP tools on the given server.
func RegisterTools(s *mcp.Server, metrics *Metrics, svc *application.NetworkService) {
	registerTool(s, toolDefIndex()["get_sites"], metrics, "get_sites", WrapToolHandler[SitesInput, PaginatedOutput[domain.Site]](metrics, "get_sites", func(ctx context.Context, req *mcp.CallToolRequest, in SitesInput) (*mcp.CallToolResult, PaginatedOutput[domain.Site], error) {
		return NewGetSitesHandler(svc)(ctx, req, in)
	}))

	registerTool(s, toolDefIndex()["get_devices"], metrics, "get_devices", WrapToolHandler[DevicesInput, PaginatedOutput[domain.Device]](metrics, "get_devices", func(ctx context.Context, req *mcp.CallToolRequest, in DevicesInput) (*mcp.CallToolResult, PaginatedOutput[domain.Device], error) {
		return NewGetDevicesHandler(svc)(ctx, req, in)
	}))

	registerTool(s, toolDefIndex()["get_ip_addresses"], metrics, "get_ip_addresses", WrapToolHandler[IPAddressesInput, PaginatedOutput[domain.IPAddress]](metrics, "get_ip_addresses", func(ctx context.Context, req *mcp.CallToolRequest, in IPAddressesInput) (*mcp.CallToolResult, PaginatedOutput[domain.IPAddress], error) {
		return NewGetIPAddressesHandler(svc)(ctx, req, in)
	}))

	registerTool(s, toolDefIndex()["get_prefixes"], metrics, "get_prefixes", WrapToolHandler[PrefixesInput, PaginatedOutput[domain.Prefix]](metrics, "get_prefixes", func(ctx context.Context, req *mcp.CallToolRequest, in PrefixesInput) (*mcp.CallToolResult, PaginatedOutput[domain.Prefix], error) {
		return NewGetPrefixesHandler(svc)(ctx, req, in)
	}))

	registerTool(s, toolDefIndex()["get_vlans"], metrics, "get_vlans", WrapToolHandler[VLANsInput, PaginatedOutput[domain.VLAN]](metrics, "get_vlans", func(ctx context.Context, req *mcp.CallToolRequest, in VLANsInput) (*mcp.CallToolResult, PaginatedOutput[domain.VLAN], error) {
		return NewGetVLANsHandler(svc)(ctx, req, in)
	}))

	registerTool(s, toolDefIndex()["get_virtual_machines"], metrics, "get_virtual_machines", WrapToolHandler[VirtualMachinesInput, PaginatedOutput[domain.VirtualMachine]](metrics, "get_virtual_machines", func(ctx context.Context, req *mcp.CallToolRequest, in VirtualMachinesInput) (*mcp.CallToolResult, PaginatedOutput[domain.VirtualMachine], error) {
		return NewGetVirtualMachinesHandler(svc)(ctx, req, in)
	}))

	registerTool(s, toolDefIndex()["get_clusters"], metrics, "get_clusters", WrapToolHandler[ClustersInput, PaginatedOutput[domain.Cluster]](metrics, "get_clusters", func(ctx context.Context, req *mcp.CallToolRequest, in ClustersInput) (*mcp.CallToolResult, PaginatedOutput[domain.Cluster], error) {
		return NewGetClustersHandler(svc)(ctx, req, in)
	}))

	registerTool(s, toolDefIndex()["get_circuits"], metrics, "get_circuits", WrapToolHandler[CircuitsInput, PaginatedOutput[domain.Circuit]](metrics, "get_circuits", func(ctx context.Context, req *mcp.CallToolRequest, in CircuitsInput) (*mcp.CallToolResult, PaginatedOutput[domain.Circuit], error) {
		return NewGetCircuitsHandler(svc)(ctx, req, in)
	}))

	registerTool(s, toolDefIndex()["get_object_by_id"], metrics, "get_object_by_id", WrapToolHandler[GetObjectInput, GetObjectOutput](metrics, "get_object_by_id", func(ctx context.Context, req *mcp.CallToolRequest, in GetObjectInput) (*mcp.CallToolResult, GetObjectOutput, error) {
		return NewGetObjectByIDHandler(svc)(ctx, req, in)
	}))

	registerTool(s, toolDefIndex()["get_racks"], metrics, "get_racks", WrapToolHandler[RacksInput, PaginatedOutput[domain.Rack]](metrics, "get_racks", func(ctx context.Context, req *mcp.CallToolRequest, in RacksInput) (*mcp.CallToolResult, PaginatedOutput[domain.Rack], error) {
		return NewGetRacksHandler(svc)(ctx, req, in)
	}))

	registerTool(s, toolDefIndex()["get_interfaces"], metrics, "get_interfaces", WrapToolHandler[InterfacesInput, PaginatedOutput[domain.Interface]](metrics, "get_interfaces", func(ctx context.Context, req *mcp.CallToolRequest, in InterfacesInput) (*mcp.CallToolResult, PaginatedOutput[domain.Interface], error) {
		return NewGetInterfacesHandler(svc)(ctx, req, in)
	}))

	registerTool(s, toolDefIndex()["get_circuit_terminations"], metrics, "get_circuit_terminations", WrapToolHandler[CircuitTerminationsInput, PaginatedOutput[domain.CircuitTermination]](metrics, "get_circuit_terminations", func(ctx context.Context, req *mcp.CallToolRequest, in CircuitTerminationsInput) (*mcp.CallToolResult, PaginatedOutput[domain.CircuitTermination], error) {
		return NewGetCircuitTerminationsHandler(svc)(ctx, req, in)
	}))

	registerTool(s, toolDefIndex()["get_cables"], metrics, "get_cables", WrapToolHandler[CablesInput, PaginatedOutput[domain.Cable]](metrics, "get_cables", func(ctx context.Context, req *mcp.CallToolRequest, in CablesInput) (*mcp.CallToolResult, PaginatedOutput[domain.Cable], error) {
		return NewGetCablesHandler(svc)(ctx, req, in)
	}))

	registerTool(s, toolDefIndex()["get_vm_interfaces"], metrics, "get_vm_interfaces", WrapToolHandler[VMInterfacesInput, PaginatedOutput[domain.VMInterface]](metrics, "get_vm_interfaces", func(ctx context.Context, req *mcp.CallToolRequest, in VMInterfacesInput) (*mcp.CallToolResult, PaginatedOutput[domain.VMInterface], error) {
		return NewGetVMInterfacesHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["create_site"], metrics, "create_site", WrapToolHandler[SiteCreateInput, SiteOutput](metrics, "create_site", func(ctx context.Context, req *mcp.CallToolRequest, in SiteCreateInput) (*mcp.CallToolResult, SiteOutput, error) {
		return NewCreateSiteHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["update_site"], metrics, "update_site", WrapToolHandler[SiteUpdateInput, SiteOutput](metrics, "update_site", func(ctx context.Context, req *mcp.CallToolRequest, in SiteUpdateInput) (*mcp.CallToolResult, SiteOutput, error) {
		return NewUpdateSiteHandler(svc)(ctx, req, in)
	}))

	registerDeleteTool(s, toolDefIndex()["delete_site"], metrics, "delete_site", WrapToolHandler[SiteDeleteInput, struct{}](metrics, "delete_site", func(ctx context.Context, req *mcp.CallToolRequest, in SiteDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		return NewDeleteSiteHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["create_device"], metrics, "create_device", WrapToolHandler[DeviceCreateInput, DeviceOutput](metrics, "create_device", func(ctx context.Context, req *mcp.CallToolRequest, in DeviceCreateInput) (*mcp.CallToolResult, DeviceOutput, error) {
		return NewCreateDeviceHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["update_device"], metrics, "update_device", WrapToolHandler[DeviceUpdateInput, DeviceOutput](metrics, "update_device", func(ctx context.Context, req *mcp.CallToolRequest, in DeviceUpdateInput) (*mcp.CallToolResult, DeviceOutput, error) {
		return NewUpdateDeviceHandler(svc)(ctx, req, in)
	}))

	registerDeleteTool(s, toolDefIndex()["delete_device"], metrics, "delete_device", WrapToolHandler[DeviceDeleteInput, struct{}](metrics, "delete_device", func(ctx context.Context, req *mcp.CallToolRequest, in DeviceDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		return NewDeleteDeviceHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["create_ip_address"], metrics, "create_ip_address", WrapToolHandler[IPAddressCreateInput, IPAddressOutput](metrics, "create_ip_address", func(ctx context.Context, req *mcp.CallToolRequest, in IPAddressCreateInput) (*mcp.CallToolResult, IPAddressOutput, error) {
		return NewCreateIPAddressHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["update_ip_address"], metrics, "update_ip_address", WrapToolHandler[IPAddressUpdateInput, IPAddressOutput](metrics, "update_ip_address", func(ctx context.Context, req *mcp.CallToolRequest, in IPAddressUpdateInput) (*mcp.CallToolResult, IPAddressOutput, error) {
		return NewUpdateIPAddressHandler(svc)(ctx, req, in)
	}))

	registerDeleteTool(s, toolDefIndex()["delete_ip_address"], metrics, "delete_ip_address", WrapToolHandler[IPAddressDeleteInput, struct{}](metrics, "delete_ip_address", func(ctx context.Context, req *mcp.CallToolRequest, in IPAddressDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		return NewDeleteIPAddressHandler(svc)(ctx, req, in)
	}))
}
