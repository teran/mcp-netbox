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
//
//nolint:maintidx // one registration block per tool; the size is inherent to the tool inventory
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

	registerWriteTool(s, toolDefIndex()["create_prefix"], metrics, "create_prefix", WrapToolHandler[PrefixCreateInput, PrefixOutput](metrics, "create_prefix", func(ctx context.Context, req *mcp.CallToolRequest, in PrefixCreateInput) (*mcp.CallToolResult, PrefixOutput, error) {
		return NewCreatePrefixHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["update_prefix"], metrics, "update_prefix", WrapToolHandler[PrefixUpdateInput, PrefixOutput](metrics, "update_prefix", func(ctx context.Context, req *mcp.CallToolRequest, in PrefixUpdateInput) (*mcp.CallToolResult, PrefixOutput, error) {
		return NewUpdatePrefixHandler(svc)(ctx, req, in)
	}))

	registerDeleteTool(s, toolDefIndex()["delete_prefix"], metrics, "delete_prefix", WrapToolHandler[PrefixDeleteInput, struct{}](metrics, "delete_prefix", func(ctx context.Context, req *mcp.CallToolRequest, in PrefixDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		return NewDeletePrefixHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["create_vlan"], metrics, "create_vlan", WrapToolHandler[VLANCreateInput, VLANOutput](metrics, "create_vlan", func(ctx context.Context, req *mcp.CallToolRequest, in VLANCreateInput) (*mcp.CallToolResult, VLANOutput, error) {
		return NewCreateVLANHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["update_vlan"], metrics, "update_vlan", WrapToolHandler[VLANUpdateInput, VLANOutput](metrics, "update_vlan", func(ctx context.Context, req *mcp.CallToolRequest, in VLANUpdateInput) (*mcp.CallToolResult, VLANOutput, error) {
		return NewUpdateVLANHandler(svc)(ctx, req, in)
	}))

	registerDeleteTool(s, toolDefIndex()["delete_vlan"], metrics, "delete_vlan", WrapToolHandler[VLANDeleteInput, struct{}](metrics, "delete_vlan", func(ctx context.Context, req *mcp.CallToolRequest, in VLANDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		return NewDeleteVLANHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["create_virtual_machine"], metrics, "create_virtual_machine", WrapToolHandler[VirtualMachineCreateInput, VirtualMachineOutput](metrics, "create_virtual_machine", func(ctx context.Context, req *mcp.CallToolRequest, in VirtualMachineCreateInput) (*mcp.CallToolResult, VirtualMachineOutput, error) {
		return NewCreateVirtualMachineHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["update_virtual_machine"], metrics, "update_virtual_machine", WrapToolHandler[VirtualMachineUpdateInput, VirtualMachineOutput](metrics, "update_virtual_machine", func(ctx context.Context, req *mcp.CallToolRequest, in VirtualMachineUpdateInput) (*mcp.CallToolResult, VirtualMachineOutput, error) {
		return NewUpdateVirtualMachineHandler(svc)(ctx, req, in)
	}))

	registerDeleteTool(s, toolDefIndex()["delete_virtual_machine"], metrics, "delete_virtual_machine", WrapToolHandler[VirtualMachineDeleteInput, struct{}](metrics, "delete_virtual_machine", func(ctx context.Context, req *mcp.CallToolRequest, in VirtualMachineDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		return NewDeleteVirtualMachineHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["create_cluster"], metrics, "create_cluster", WrapToolHandler[ClusterCreateInput, ClusterOutput](metrics, "create_cluster", func(ctx context.Context, req *mcp.CallToolRequest, in ClusterCreateInput) (*mcp.CallToolResult, ClusterOutput, error) {
		return NewCreateClusterHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["update_cluster"], metrics, "update_cluster", WrapToolHandler[ClusterUpdateInput, ClusterOutput](metrics, "update_cluster", func(ctx context.Context, req *mcp.CallToolRequest, in ClusterUpdateInput) (*mcp.CallToolResult, ClusterOutput, error) {
		return NewUpdateClusterHandler(svc)(ctx, req, in)
	}))

	registerDeleteTool(s, toolDefIndex()["delete_cluster"], metrics, "delete_cluster", WrapToolHandler[ClusterDeleteInput, struct{}](metrics, "delete_cluster", func(ctx context.Context, req *mcp.CallToolRequest, in ClusterDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		return NewDeleteClusterHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["create_circuit"], metrics, "create_circuit", WrapToolHandler[CircuitCreateInput, CircuitOutput](metrics, "create_circuit", func(ctx context.Context, req *mcp.CallToolRequest, in CircuitCreateInput) (*mcp.CallToolResult, CircuitOutput, error) {
		return NewCreateCircuitHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["update_circuit"], metrics, "update_circuit", WrapToolHandler[CircuitUpdateInput, CircuitOutput](metrics, "update_circuit", func(ctx context.Context, req *mcp.CallToolRequest, in CircuitUpdateInput) (*mcp.CallToolResult, CircuitOutput, error) {
		return NewUpdateCircuitHandler(svc)(ctx, req, in)
	}))

	registerDeleteTool(s, toolDefIndex()["delete_circuit"], metrics, "delete_circuit", WrapToolHandler[CircuitDeleteInput, struct{}](metrics, "delete_circuit", func(ctx context.Context, req *mcp.CallToolRequest, in CircuitDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		return NewDeleteCircuitHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["create_rack"], metrics, "create_rack", WrapToolHandler[RackCreateInput, RackOutput](metrics, "create_rack", func(ctx context.Context, req *mcp.CallToolRequest, in RackCreateInput) (*mcp.CallToolResult, RackOutput, error) {
		return NewCreateRackHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["update_rack"], metrics, "update_rack", WrapToolHandler[RackUpdateInput, RackOutput](metrics, "update_rack", func(ctx context.Context, req *mcp.CallToolRequest, in RackUpdateInput) (*mcp.CallToolResult, RackOutput, error) {
		return NewUpdateRackHandler(svc)(ctx, req, in)
	}))

	registerDeleteTool(s, toolDefIndex()["delete_rack"], metrics, "delete_rack", WrapToolHandler[RackDeleteInput, struct{}](metrics, "delete_rack", func(ctx context.Context, req *mcp.CallToolRequest, in RackDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		return NewDeleteRackHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["create_interface"], metrics, "create_interface", WrapToolHandler[InterfaceCreateInput, InterfaceOutput](metrics, "create_interface", func(ctx context.Context, req *mcp.CallToolRequest, in InterfaceCreateInput) (*mcp.CallToolResult, InterfaceOutput, error) {
		return NewCreateInterfaceHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["update_interface"], metrics, "update_interface", WrapToolHandler[InterfaceUpdateInput, InterfaceOutput](metrics, "update_interface", func(ctx context.Context, req *mcp.CallToolRequest, in InterfaceUpdateInput) (*mcp.CallToolResult, InterfaceOutput, error) {
		return NewUpdateInterfaceHandler(svc)(ctx, req, in)
	}))

	registerDeleteTool(s, toolDefIndex()["delete_interface"], metrics, "delete_interface", WrapToolHandler[InterfaceDeleteInput, struct{}](metrics, "delete_interface", func(ctx context.Context, req *mcp.CallToolRequest, in InterfaceDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		return NewDeleteInterfaceHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["create_circuit_termination"], metrics, "create_circuit_termination", WrapToolHandler[CircuitTerminationCreateInput, CircuitTerminationOutput](metrics, "create_circuit_termination", func(ctx context.Context, req *mcp.CallToolRequest, in CircuitTerminationCreateInput) (*mcp.CallToolResult, CircuitTerminationOutput, error) {
		return NewCreateCircuitTerminationHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["update_circuit_termination"], metrics, "update_circuit_termination", WrapToolHandler[CircuitTerminationUpdateInput, CircuitTerminationOutput](metrics, "update_circuit_termination", func(ctx context.Context, req *mcp.CallToolRequest, in CircuitTerminationUpdateInput) (*mcp.CallToolResult, CircuitTerminationOutput, error) {
		return NewUpdateCircuitTerminationHandler(svc)(ctx, req, in)
	}))

	registerDeleteTool(s, toolDefIndex()["delete_circuit_termination"], metrics, "delete_circuit_termination", WrapToolHandler[CircuitTerminationDeleteInput, struct{}](metrics, "delete_circuit_termination", func(ctx context.Context, req *mcp.CallToolRequest, in CircuitTerminationDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		return NewDeleteCircuitTerminationHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["create_cable"], metrics, "create_cable", WrapToolHandler[CableCreateInput, CableOutput](metrics, "create_cable", func(ctx context.Context, req *mcp.CallToolRequest, in CableCreateInput) (*mcp.CallToolResult, CableOutput, error) {
		return NewCreateCableHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["update_cable"], metrics, "update_cable", WrapToolHandler[CableUpdateInput, CableOutput](metrics, "update_cable", func(ctx context.Context, req *mcp.CallToolRequest, in CableUpdateInput) (*mcp.CallToolResult, CableOutput, error) {
		return NewUpdateCableHandler(svc)(ctx, req, in)
	}))

	registerDeleteTool(s, toolDefIndex()["delete_cable"], metrics, "delete_cable", WrapToolHandler[CableDeleteInput, struct{}](metrics, "delete_cable", func(ctx context.Context, req *mcp.CallToolRequest, in CableDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		return NewDeleteCableHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["create_vm_interface"], metrics, "create_vm_interface", WrapToolHandler[VMInterfaceCreateInput, VMInterfaceOutput](metrics, "create_vm_interface", func(ctx context.Context, req *mcp.CallToolRequest, in VMInterfaceCreateInput) (*mcp.CallToolResult, VMInterfaceOutput, error) {
		return NewCreateVMInterfaceHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["update_vm_interface"], metrics, "update_vm_interface", WrapToolHandler[VMInterfaceUpdateInput, VMInterfaceOutput](metrics, "update_vm_interface", func(ctx context.Context, req *mcp.CallToolRequest, in VMInterfaceUpdateInput) (*mcp.CallToolResult, VMInterfaceOutput, error) {
		return NewUpdateVMInterfaceHandler(svc)(ctx, req, in)
	}))

	registerDeleteTool(s, toolDefIndex()["delete_vm_interface"], metrics, "delete_vm_interface", WrapToolHandler[VMInterfaceDeleteInput, struct{}](metrics, "delete_vm_interface", func(ctx context.Context, req *mcp.CallToolRequest, in VMInterfaceDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		return NewDeleteVMInterfaceHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["create_provider"], metrics, "create_provider", WrapToolHandler[ProviderCreateInput, ProviderOutput](metrics, "create_provider", func(ctx context.Context, req *mcp.CallToolRequest, in ProviderCreateInput) (*mcp.CallToolResult, ProviderOutput, error) {
		return NewCreateProviderHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["update_provider"], metrics, "update_provider", WrapToolHandler[ProviderUpdateInput, ProviderOutput](metrics, "update_provider", func(ctx context.Context, req *mcp.CallToolRequest, in ProviderUpdateInput) (*mcp.CallToolResult, ProviderOutput, error) {
		return NewUpdateProviderHandler(svc)(ctx, req, in)
	}))

	registerDeleteTool(s, toolDefIndex()["delete_provider"], metrics, "delete_provider", WrapToolHandler[ProviderDeleteInput, struct{}](metrics, "delete_provider", func(ctx context.Context, req *mcp.CallToolRequest, in ProviderDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		return NewDeleteProviderHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["create_tenant"], metrics, "create_tenant", WrapToolHandler[TenantCreateInput, TenantOutput](metrics, "create_tenant", func(ctx context.Context, req *mcp.CallToolRequest, in TenantCreateInput) (*mcp.CallToolResult, TenantOutput, error) {
		return NewCreateTenantHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["update_tenant"], metrics, "update_tenant", WrapToolHandler[TenantUpdateInput, TenantOutput](metrics, "update_tenant", func(ctx context.Context, req *mcp.CallToolRequest, in TenantUpdateInput) (*mcp.CallToolResult, TenantOutput, error) {
		return NewUpdateTenantHandler(svc)(ctx, req, in)
	}))

	registerDeleteTool(s, toolDefIndex()["delete_tenant"], metrics, "delete_tenant", WrapToolHandler[TenantDeleteInput, struct{}](metrics, "delete_tenant", func(ctx context.Context, req *mcp.CallToolRequest, in TenantDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		return NewDeleteTenantHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["create_manufacturer"], metrics, "create_manufacturer", WrapToolHandler[ManufacturerCreateInput, ManufacturerOutput](metrics, "create_manufacturer", func(ctx context.Context, req *mcp.CallToolRequest, in ManufacturerCreateInput) (*mcp.CallToolResult, ManufacturerOutput, error) {
		return NewCreateManufacturerHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["update_manufacturer"], metrics, "update_manufacturer", WrapToolHandler[ManufacturerUpdateInput, ManufacturerOutput](metrics, "update_manufacturer", func(ctx context.Context, req *mcp.CallToolRequest, in ManufacturerUpdateInput) (*mcp.CallToolResult, ManufacturerOutput, error) {
		return NewUpdateManufacturerHandler(svc)(ctx, req, in)
	}))

	registerDeleteTool(s, toolDefIndex()["delete_manufacturer"], metrics, "delete_manufacturer", WrapToolHandler[ManufacturerDeleteInput, struct{}](metrics, "delete_manufacturer", func(ctx context.Context, req *mcp.CallToolRequest, in ManufacturerDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		return NewDeleteManufacturerHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["create_device_type"], metrics, "create_device_type", WrapToolHandler[DeviceTypeCreateInput, DeviceTypeOutput](metrics, "create_device_type", func(ctx context.Context, req *mcp.CallToolRequest, in DeviceTypeCreateInput) (*mcp.CallToolResult, DeviceTypeOutput, error) {
		return NewCreateDeviceTypeHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["update_device_type"], metrics, "update_device_type", WrapToolHandler[DeviceTypeUpdateInput, DeviceTypeOutput](metrics, "update_device_type", func(ctx context.Context, req *mcp.CallToolRequest, in DeviceTypeUpdateInput) (*mcp.CallToolResult, DeviceTypeOutput, error) {
		return NewUpdateDeviceTypeHandler(svc)(ctx, req, in)
	}))

	registerDeleteTool(s, toolDefIndex()["delete_device_type"], metrics, "delete_device_type", WrapToolHandler[DeviceTypeDeleteInput, struct{}](metrics, "delete_device_type", func(ctx context.Context, req *mcp.CallToolRequest, in DeviceTypeDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		return NewDeleteDeviceTypeHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["create_location"], metrics, "create_location", WrapToolHandler[LocationCreateInput, LocationOutput](metrics, "create_location", func(ctx context.Context, req *mcp.CallToolRequest, in LocationCreateInput) (*mcp.CallToolResult, LocationOutput, error) {
		return NewCreateLocationHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["update_location"], metrics, "update_location", WrapToolHandler[LocationUpdateInput, LocationOutput](metrics, "update_location", func(ctx context.Context, req *mcp.CallToolRequest, in LocationUpdateInput) (*mcp.CallToolResult, LocationOutput, error) {
		return NewUpdateLocationHandler(svc)(ctx, req, in)
	}))

	registerDeleteTool(s, toolDefIndex()["delete_location"], metrics, "delete_location", WrapToolHandler[LocationDeleteInput, struct{}](metrics, "delete_location", func(ctx context.Context, req *mcp.CallToolRequest, in LocationDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		return NewDeleteLocationHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["create_cluster_type"], metrics, "create_cluster_type", WrapToolHandler[ClusterTypeCreateInput, ClusterTypeOutput](metrics, "create_cluster_type", func(ctx context.Context, req *mcp.CallToolRequest, in ClusterTypeCreateInput) (*mcp.CallToolResult, ClusterTypeOutput, error) {
		return NewCreateClusterTypeHandler(svc)(ctx, req, in)
	}))

	registerWriteTool(s, toolDefIndex()["update_cluster_type"], metrics, "update_cluster_type", WrapToolHandler[ClusterTypeUpdateInput, ClusterTypeOutput](metrics, "update_cluster_type", func(ctx context.Context, req *mcp.CallToolRequest, in ClusterTypeUpdateInput) (*mcp.CallToolResult, ClusterTypeOutput, error) {
		return NewUpdateClusterTypeHandler(svc)(ctx, req, in)
	}))

	registerDeleteTool(s, toolDefIndex()["delete_cluster_type"], metrics, "delete_cluster_type", WrapToolHandler[ClusterTypeDeleteInput, struct{}](metrics, "delete_cluster_type", func(ctx context.Context, req *mcp.CallToolRequest, in ClusterTypeDeleteInput) (*mcp.CallToolResult, struct{}, error) {
		return NewDeleteClusterTypeHandler(svc)(ctx, req, in)
	}))
}
