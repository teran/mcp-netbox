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
//
//go:fix inline
//nolint:modernize // newexpr suggests new(b); new(false)/new(true) are not always valid (new(true) is a compile error).
func boolPtr(b bool) *bool {
	return &b
}

// readOnlyTool returns tool annotations for one of the read-only NetBox query
// tools. Every read tool operates on a closed domain (the NetBox inventory),
// so these hints are constant across the read tools.
func readOnlyTool(title string) *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{
		Title:           title,
		ReadOnlyHint:    true,
		IdempotentHint:  true,
		DestructiveHint: new(bool),
		OpenWorldHint:   new(bool),
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
		DestructiveHint: new(bool),
		OpenWorldHint:   new(bool),
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
		DestructiveHint: boolPtr(true), //nolint:modernize,lll // new(true) is a compile error; boolPtr(true) cannot be simplified.
		OpenWorldHint:   new(bool),
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

// toolDescription combines a tool definition's base Description with its
// per-tool Instructions (M4), which the SDK surfaces through the Description
// channel. When a tool carries no instructions the base description is
// returned unchanged.
func toolDescription(def toolDef) string {
	description := def.Description
	if def.Instructions != "" {
		description = description + "\n\nInstructions: " + def.Instructions
	}

	return description
}

// registerTool is a helper that builds and registers a single read-only tool,
// combining its description with the per-tool instructions.
func registerTool[I, O any](s *mcp.Server, def toolDef, _ *Metrics, name string, handler mcp.ToolHandlerFor[I, O]) {
	description := toolDescription(def)
	mcp.AddTool(s, &mcp.Tool{
		Name:        name,
		Description: description,
		Annotations: readOnlyTool(def.Title),
	}, handler)
}

// registerWriteTool registers a non-destructive write tool with writeTool
// annotations.
func registerWriteTool[I, O any](
	s *mcp.Server, def toolDef, _ *Metrics, name string, handler mcp.ToolHandlerFor[I, O],
) {
	description := toolDescription(def)
	mcp.AddTool(s, &mcp.Tool{
		Name:        name,
		Description: description,
		Annotations: writeTool(def.Title),
	}, handler)
}

// registerDeleteTool registers a destructive tool with destructiveTool
// annotations.
func registerDeleteTool[I, O any](
	s *mcp.Server, def toolDef, _ *Metrics, name string, handler mcp.ToolHandlerFor[I, O],
) {
	description := toolDescription(def)
	mcp.AddTool(s, &mcp.Tool{
		Name:        name,
		Description: description,
		Annotations: destructiveTool(def.Title),
	}, handler)
}

// RegisterTools registers all NetBox MCP tools on the given server.
//
//nolint:funlen,maintidx // one registration block per tool; the size is inherent to the tool inventory
func RegisterTools(server *mcp.Server, metrics *Metrics, svc *application.NetworkService) {
	registerTool(
		server,
		toolDefIndex()["get_sites"],
		metrics,
		"get_sites",
		WrapToolHandler[SitesInput, PaginatedOutput[domain.Site]](
			metrics,
			"get_sites",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in SitesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.Site], error) {
				return NewGetSitesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_regions"],
		metrics,
		"get_regions",
		WrapToolHandler[RegionsInput, PaginatedOutput[domain.Region]](
			metrics,
			"get_regions",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RegionsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.Region], error) {
				return NewGetRegionsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_tenant_groups"],
		metrics,
		"get_tenant_groups",
		WrapToolHandler[TenantGroupsInput, PaginatedOutput[domain.TenantGroup]](
			metrics,
			"get_tenant_groups",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in TenantGroupsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.TenantGroup], error) {
				return NewGetTenantGroupsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_devices"],
		metrics,
		"get_devices",
		WrapToolHandler[DevicesInput, PaginatedOutput[domain.Device]](
			metrics,
			"get_devices",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in DevicesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.Device], error) {
				return NewGetDevicesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_ip_addresses"],
		metrics,
		"get_ip_addresses",
		WrapToolHandler[IPAddressesInput, PaginatedOutput[domain.IPAddress]](
			metrics,
			"get_ip_addresses",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in IPAddressesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.IPAddress], error) {
				return NewGetIPAddressesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_prefixes"],
		metrics,
		"get_prefixes",
		WrapToolHandler[PrefixesInput, PaginatedOutput[domain.Prefix]](
			metrics,
			"get_prefixes",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PrefixesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.Prefix], error) {
				return NewGetPrefixesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_vlans"],
		metrics,
		"get_vlans",
		WrapToolHandler[VLANsInput, PaginatedOutput[domain.VLAN]](
			metrics,
			"get_vlans",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VLANsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.VLAN], error) {
				return NewGetVLANsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_virtual_machines"],
		metrics,
		"get_virtual_machines",
		WrapToolHandler[VirtualMachinesInput, PaginatedOutput[domain.VirtualMachine]](
			metrics,
			"get_virtual_machines",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualMachinesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.VirtualMachine], error) {
				return NewGetVirtualMachinesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_clusters"],
		metrics,
		"get_clusters",
		WrapToolHandler[ClustersInput, PaginatedOutput[domain.Cluster]](
			metrics,
			"get_clusters",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ClustersInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.Cluster], error) {
				return NewGetClustersHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_circuits"],
		metrics,
		"get_circuits",
		WrapToolHandler[CircuitsInput, PaginatedOutput[domain.Circuit]](
			metrics,
			"get_circuits",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CircuitsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.Circuit], error) {
				return NewGetCircuitsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_object_by_id"],
		metrics,
		"get_object_by_id",
		WrapToolHandler[GetObjectInput, GetObjectOutput](
			metrics,
			"get_object_by_id",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in GetObjectInput,
			) (*mcp.CallToolResult, GetObjectOutput, error) {
				return NewGetObjectByIDHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_racks"],
		metrics,
		"get_racks",
		WrapToolHandler[RacksInput, PaginatedOutput[domain.Rack]](
			metrics,
			"get_racks",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RacksInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.Rack], error) {
				return NewGetRacksHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_interfaces"],
		metrics,
		"get_interfaces",
		WrapToolHandler[InterfacesInput, PaginatedOutput[domain.Interface]](
			metrics,
			"get_interfaces",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in InterfacesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.Interface], error) {
				return NewGetInterfacesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_circuit_terminations"],
		metrics,
		"get_circuit_terminations",
		WrapToolHandler[CircuitTerminationsInput, PaginatedOutput[domain.CircuitTermination]](
			metrics,
			"get_circuit_terminations",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CircuitTerminationsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.CircuitTermination], error) {
				return NewGetCircuitTerminationsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_cables"],
		metrics,
		"get_cables",
		WrapToolHandler[CablesInput, PaginatedOutput[domain.Cable]](
			metrics,
			"get_cables",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CablesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.Cable], error) {
				return NewGetCablesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_vm_interfaces"],
		metrics,
		"get_vm_interfaces",
		WrapToolHandler[VMInterfacesInput, PaginatedOutput[domain.VMInterface]](
			metrics,
			"get_vm_interfaces",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VMInterfacesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.VMInterface], error) {
				return NewGetVMInterfacesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_site_groups"],
		metrics,
		"get_site_groups",
		WrapToolHandler[SiteGroupsInput, PaginatedOutput[domain.SiteGroup]](
			metrics,
			"get_site_groups",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in SiteGroupsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.SiteGroup], error) {
				return NewGetSiteGroupsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_rack_groups"],
		metrics,
		"get_rack_groups",
		WrapToolHandler[RackGroupsInput, PaginatedOutput[domain.RackGroup]](
			metrics,
			"get_rack_groups",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RackGroupsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.RackGroup], error) {
				return NewGetRackGroupsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_rack_roles"],
		metrics,
		"get_rack_roles",
		WrapToolHandler[RackRolesInput, PaginatedOutput[domain.RackRole]](
			metrics,
			"get_rack_roles",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RackRolesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.RackRole], error) {
				return NewGetRackRolesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_device_roles"],
		metrics,
		"get_device_roles",
		WrapToolHandler[DeviceRolesInput, PaginatedOutput[domain.DeviceRole]](
			metrics,
			"get_device_roles",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in DeviceRolesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.DeviceRole], error) {
				return NewGetDeviceRolesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_platforms"],
		metrics,
		"get_platforms",
		WrapToolHandler[PlatformsInput, PaginatedOutput[domain.Platform]](
			metrics,
			"get_platforms",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PlatformsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.Platform], error) {
				return NewGetPlatformsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_virtual_chassis"],
		metrics,
		"get_virtual_chassis",
		WrapToolHandler[VirtualChassisInput, PaginatedOutput[domain.VirtualChassis]](
			metrics,
			"get_virtual_chassis",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualChassisInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.VirtualChassis], error) {
				return NewGetVirtualChassisHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_power_panels"],
		metrics,
		"get_power_panels",
		WrapToolHandler[PowerPanelsInput, PaginatedOutput[domain.PowerPanel]](
			metrics,
			"get_power_panels",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PowerPanelsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.PowerPanel], error) {
				return NewGetPowerPanelsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_inventory_item_roles"],
		metrics,
		"get_inventory_item_roles",
		WrapToolHandler[InventoryItemRolesInput, PaginatedOutput[domain.InventoryItemRole]](
			metrics,
			"get_inventory_item_roles",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in InventoryItemRolesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.InventoryItemRole], error) {
				return NewGetInventoryItemRolesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_site"],
		metrics,
		"create_site",
		WrapToolHandler[SiteCreateInput, SiteOutput](
			metrics,
			"create_site",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in SiteCreateInput,
			) (*mcp.CallToolResult, SiteOutput, error) {
				return NewCreateSiteHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_site"],
		metrics,
		"update_site",
		WrapToolHandler[SiteUpdateInput, SiteOutput](
			metrics,
			"update_site",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in SiteUpdateInput,
			) (*mcp.CallToolResult, SiteOutput, error) {
				return NewUpdateSiteHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_site"],
		metrics,
		"delete_site",
		WrapToolHandler[SiteDeleteInput, struct{}](
			metrics,
			"delete_site",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in SiteDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteSiteHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_region"],
		metrics,
		"create_region",
		WrapToolHandler[RegionCreateInput, RegionOutput](
			metrics,
			"create_region",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RegionCreateInput,
			) (*mcp.CallToolResult, RegionOutput, error) {
				return NewCreateRegionHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_region"],
		metrics,
		"update_region",
		WrapToolHandler[RegionUpdateInput, RegionOutput](
			metrics,
			"update_region",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RegionUpdateInput,
			) (*mcp.CallToolResult, RegionOutput, error) {
				return NewUpdateRegionHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_region"],
		metrics,
		"delete_region",
		WrapToolHandler[RegionDeleteInput, struct{}](
			metrics,
			"delete_region",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RegionDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteRegionHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_tenant_group"],
		metrics,
		"create_tenant_group",
		WrapToolHandler[TenantGroupCreateInput, TenantGroupOutput](
			metrics,
			"create_tenant_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in TenantGroupCreateInput,
			) (*mcp.CallToolResult, TenantGroupOutput, error) {
				return NewCreateTenantGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_tenant_group"],
		metrics,
		"update_tenant_group",
		WrapToolHandler[TenantGroupUpdateInput, TenantGroupOutput](
			metrics,
			"update_tenant_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in TenantGroupUpdateInput,
			) (*mcp.CallToolResult, TenantGroupOutput, error) {
				return NewUpdateTenantGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_tenant_group"],
		metrics,
		"delete_tenant_group",
		WrapToolHandler[TenantGroupDeleteInput, struct{}](
			metrics,
			"delete_tenant_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in TenantGroupDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteTenantGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_device"],
		metrics,
		"create_device",
		WrapToolHandler[DeviceCreateInput, DeviceOutput](
			metrics,
			"create_device",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in DeviceCreateInput,
			) (*mcp.CallToolResult, DeviceOutput, error) {
				return NewCreateDeviceHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_device"],
		metrics,
		"update_device",
		WrapToolHandler[DeviceUpdateInput, DeviceOutput](
			metrics,
			"update_device",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in DeviceUpdateInput,
			) (*mcp.CallToolResult, DeviceOutput, error) {
				return NewUpdateDeviceHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_device"],
		metrics,
		"delete_device",
		WrapToolHandler[DeviceDeleteInput, struct{}](
			metrics,
			"delete_device",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in DeviceDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteDeviceHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_ip_address"],
		metrics,
		"create_ip_address",
		WrapToolHandler[IPAddressCreateInput, IPAddressOutput](
			metrics,
			"create_ip_address",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in IPAddressCreateInput,
			) (*mcp.CallToolResult, IPAddressOutput, error) {
				return NewCreateIPAddressHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_ip_address"],
		metrics,
		"update_ip_address",
		WrapToolHandler[IPAddressUpdateInput, IPAddressOutput](
			metrics,
			"update_ip_address",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in IPAddressUpdateInput,
			) (*mcp.CallToolResult, IPAddressOutput, error) {
				return NewUpdateIPAddressHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_ip_address"],
		metrics,
		"delete_ip_address",
		WrapToolHandler[IPAddressDeleteInput, struct{}](
			metrics,
			"delete_ip_address",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in IPAddressDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteIPAddressHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_prefix"],
		metrics,
		"create_prefix",
		WrapToolHandler[PrefixCreateInput, PrefixOutput](
			metrics,
			"create_prefix",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PrefixCreateInput,
			) (*mcp.CallToolResult, PrefixOutput, error) {
				return NewCreatePrefixHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_prefix"],
		metrics,
		"update_prefix",
		WrapToolHandler[PrefixUpdateInput, PrefixOutput](
			metrics,
			"update_prefix",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PrefixUpdateInput,
			) (*mcp.CallToolResult, PrefixOutput, error) {
				return NewUpdatePrefixHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_prefix"],
		metrics,
		"delete_prefix",
		WrapToolHandler[PrefixDeleteInput, struct{}](
			metrics,
			"delete_prefix",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PrefixDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeletePrefixHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_vlan"],
		metrics,
		"create_vlan",
		WrapToolHandler[VLANCreateInput, VLANOutput](
			metrics,
			"create_vlan",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VLANCreateInput,
			) (*mcp.CallToolResult, VLANOutput, error) {
				return NewCreateVLANHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_vlan"],
		metrics,
		"update_vlan",
		WrapToolHandler[VLANUpdateInput, VLANOutput](
			metrics,
			"update_vlan",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VLANUpdateInput,
			) (*mcp.CallToolResult, VLANOutput, error) {
				return NewUpdateVLANHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_vlan"],
		metrics,
		"delete_vlan",
		WrapToolHandler[VLANDeleteInput, struct{}](
			metrics,
			"delete_vlan",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VLANDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteVLANHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_virtual_machine"],
		metrics,
		"create_virtual_machine",
		WrapToolHandler[VirtualMachineCreateInput, VirtualMachineOutput](
			metrics,
			"create_virtual_machine",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualMachineCreateInput,
			) (*mcp.CallToolResult, VirtualMachineOutput, error) {
				return NewCreateVirtualMachineHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_virtual_machine"],
		metrics,
		"update_virtual_machine",
		WrapToolHandler[VirtualMachineUpdateInput, VirtualMachineOutput](
			metrics,
			"update_virtual_machine",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualMachineUpdateInput,
			) (*mcp.CallToolResult, VirtualMachineOutput, error) {
				return NewUpdateVirtualMachineHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_virtual_machine"],
		metrics,
		"delete_virtual_machine",
		WrapToolHandler[VirtualMachineDeleteInput, struct{}](
			metrics,
			"delete_virtual_machine",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualMachineDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteVirtualMachineHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_cluster"],
		metrics,
		"create_cluster",
		WrapToolHandler[ClusterCreateInput, ClusterOutput](
			metrics,
			"create_cluster",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ClusterCreateInput,
			) (*mcp.CallToolResult, ClusterOutput, error) {
				return NewCreateClusterHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_cluster"],
		metrics,
		"update_cluster",
		WrapToolHandler[ClusterUpdateInput, ClusterOutput](
			metrics,
			"update_cluster",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ClusterUpdateInput,
			) (*mcp.CallToolResult, ClusterOutput, error) {
				return NewUpdateClusterHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_cluster"],
		metrics,
		"delete_cluster",
		WrapToolHandler[ClusterDeleteInput, struct{}](
			metrics,
			"delete_cluster",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ClusterDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteClusterHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_circuit"],
		metrics,
		"create_circuit",
		WrapToolHandler[CircuitCreateInput, CircuitOutput](
			metrics,
			"create_circuit",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CircuitCreateInput,
			) (*mcp.CallToolResult, CircuitOutput, error) {
				return NewCreateCircuitHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_circuit"],
		metrics,
		"update_circuit",
		WrapToolHandler[CircuitUpdateInput, CircuitOutput](
			metrics,
			"update_circuit",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CircuitUpdateInput,
			) (*mcp.CallToolResult, CircuitOutput, error) {
				return NewUpdateCircuitHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_circuit"],
		metrics,
		"delete_circuit",
		WrapToolHandler[CircuitDeleteInput, struct{}](
			metrics,
			"delete_circuit",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CircuitDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteCircuitHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_rack"],
		metrics,
		"create_rack",
		WrapToolHandler[RackCreateInput, RackOutput](
			metrics,
			"create_rack",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RackCreateInput,
			) (*mcp.CallToolResult, RackOutput, error) {
				return NewCreateRackHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_rack"],
		metrics,
		"update_rack",
		WrapToolHandler[RackUpdateInput, RackOutput](
			metrics,
			"update_rack",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RackUpdateInput,
			) (*mcp.CallToolResult, RackOutput, error) {
				return NewUpdateRackHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_rack"],
		metrics,
		"delete_rack",
		WrapToolHandler[RackDeleteInput, struct{}](
			metrics,
			"delete_rack",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RackDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteRackHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_interface"],
		metrics,
		"create_interface",
		WrapToolHandler[InterfaceCreateInput, InterfaceOutput](
			metrics,
			"create_interface",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in InterfaceCreateInput,
			) (*mcp.CallToolResult, InterfaceOutput, error) {
				return NewCreateInterfaceHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_interface"],
		metrics,
		"update_interface",
		WrapToolHandler[InterfaceUpdateInput, InterfaceOutput](
			metrics,
			"update_interface",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in InterfaceUpdateInput,
			) (*mcp.CallToolResult, InterfaceOutput, error) {
				return NewUpdateInterfaceHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_interface"],
		metrics,
		"delete_interface",
		WrapToolHandler[InterfaceDeleteInput, struct{}](
			metrics,
			"delete_interface",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in InterfaceDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteInterfaceHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_circuit_termination"],
		metrics,
		"create_circuit_termination",
		WrapToolHandler[CircuitTerminationCreateInput, CircuitTerminationOutput](
			metrics,
			"create_circuit_termination",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CircuitTerminationCreateInput,
			) (*mcp.CallToolResult, CircuitTerminationOutput, error) {
				return NewCreateCircuitTerminationHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_circuit_termination"],
		metrics,
		"update_circuit_termination",
		WrapToolHandler[CircuitTerminationUpdateInput, CircuitTerminationOutput](
			metrics,
			"update_circuit_termination",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CircuitTerminationUpdateInput,
			) (*mcp.CallToolResult, CircuitTerminationOutput, error) {
				return NewUpdateCircuitTerminationHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_circuit_termination"],
		metrics,
		"delete_circuit_termination",
		WrapToolHandler[CircuitTerminationDeleteInput, struct{}](
			metrics,
			"delete_circuit_termination",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CircuitTerminationDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteCircuitTerminationHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_cable"],
		metrics,
		"create_cable",
		WrapToolHandler[CableCreateInput, CableOutput](
			metrics,
			"create_cable",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CableCreateInput,
			) (*mcp.CallToolResult, CableOutput, error) {
				return NewCreateCableHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_cable"],
		metrics,
		"update_cable",
		WrapToolHandler[CableUpdateInput, CableOutput](
			metrics,
			"update_cable",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CableUpdateInput,
			) (*mcp.CallToolResult, CableOutput, error) {
				return NewUpdateCableHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_cable"],
		metrics,
		"delete_cable",
		WrapToolHandler[CableDeleteInput, struct{}](
			metrics,
			"delete_cable",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CableDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteCableHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_vm_interface"],
		metrics,
		"create_vm_interface",
		WrapToolHandler[VMInterfaceCreateInput, VMInterfaceOutput](
			metrics,
			"create_vm_interface",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VMInterfaceCreateInput,
			) (*mcp.CallToolResult, VMInterfaceOutput, error) {
				return NewCreateVMInterfaceHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_vm_interface"],
		metrics,
		"update_vm_interface",
		WrapToolHandler[VMInterfaceUpdateInput, VMInterfaceOutput](
			metrics,
			"update_vm_interface",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VMInterfaceUpdateInput,
			) (*mcp.CallToolResult, VMInterfaceOutput, error) {
				return NewUpdateVMInterfaceHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_vm_interface"],
		metrics,
		"delete_vm_interface",
		WrapToolHandler[VMInterfaceDeleteInput, struct{}](
			metrics,
			"delete_vm_interface",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VMInterfaceDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteVMInterfaceHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_provider"],
		metrics,
		"create_provider",
		WrapToolHandler[ProviderCreateInput, ProviderOutput](
			metrics,
			"create_provider",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ProviderCreateInput,
			) (*mcp.CallToolResult, ProviderOutput, error) {
				return NewCreateProviderHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_provider"],
		metrics,
		"update_provider",
		WrapToolHandler[ProviderUpdateInput, ProviderOutput](
			metrics,
			"update_provider",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ProviderUpdateInput,
			) (*mcp.CallToolResult, ProviderOutput, error) {
				return NewUpdateProviderHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_provider"],
		metrics,
		"delete_provider",
		WrapToolHandler[ProviderDeleteInput, struct{}](
			metrics,
			"delete_provider",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ProviderDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteProviderHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_tenant"],
		metrics,
		"create_tenant",
		WrapToolHandler[TenantCreateInput, TenantOutput](
			metrics,
			"create_tenant",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in TenantCreateInput,
			) (*mcp.CallToolResult, TenantOutput, error) {
				return NewCreateTenantHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_tenant"],
		metrics,
		"update_tenant",
		WrapToolHandler[TenantUpdateInput, TenantOutput](
			metrics,
			"update_tenant",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in TenantUpdateInput,
			) (*mcp.CallToolResult, TenantOutput, error) {
				return NewUpdateTenantHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_tenant"],
		metrics,
		"delete_tenant",
		WrapToolHandler[TenantDeleteInput, struct{}](
			metrics,
			"delete_tenant",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in TenantDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteTenantHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_manufacturer"],
		metrics,
		"create_manufacturer",
		WrapToolHandler[ManufacturerCreateInput, ManufacturerOutput](
			metrics,
			"create_manufacturer",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ManufacturerCreateInput,
			) (*mcp.CallToolResult, ManufacturerOutput, error) {
				return NewCreateManufacturerHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_manufacturer"],
		metrics,
		"update_manufacturer",
		WrapToolHandler[ManufacturerUpdateInput, ManufacturerOutput](
			metrics,
			"update_manufacturer",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ManufacturerUpdateInput,
			) (*mcp.CallToolResult, ManufacturerOutput, error) {
				return NewUpdateManufacturerHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_manufacturer"],
		metrics,
		"delete_manufacturer",
		WrapToolHandler[ManufacturerDeleteInput, struct{}](
			metrics,
			"delete_manufacturer",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ManufacturerDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteManufacturerHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_device_type"],
		metrics,
		"create_device_type",
		WrapToolHandler[DeviceTypeCreateInput, DeviceTypeOutput](
			metrics,
			"create_device_type",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in DeviceTypeCreateInput,
			) (*mcp.CallToolResult, DeviceTypeOutput, error) {
				return NewCreateDeviceTypeHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_device_type"],
		metrics,
		"update_device_type",
		WrapToolHandler[DeviceTypeUpdateInput, DeviceTypeOutput](
			metrics,
			"update_device_type",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in DeviceTypeUpdateInput,
			) (*mcp.CallToolResult, DeviceTypeOutput, error) {
				return NewUpdateDeviceTypeHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_device_type"],
		metrics,
		"delete_device_type",
		WrapToolHandler[DeviceTypeDeleteInput, struct{}](
			metrics,
			"delete_device_type",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in DeviceTypeDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteDeviceTypeHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_location"],
		metrics,
		"create_location",
		WrapToolHandler[LocationCreateInput, LocationOutput](
			metrics,
			"create_location",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in LocationCreateInput,
			) (*mcp.CallToolResult, LocationOutput, error) {
				return NewCreateLocationHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_location"],
		metrics,
		"update_location",
		WrapToolHandler[LocationUpdateInput, LocationOutput](
			metrics,
			"update_location",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in LocationUpdateInput,
			) (*mcp.CallToolResult, LocationOutput, error) {
				return NewUpdateLocationHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_location"],
		metrics,
		"delete_location",
		WrapToolHandler[LocationDeleteInput, struct{}](
			metrics,
			"delete_location",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in LocationDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteLocationHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_cluster_type"],
		metrics,
		"create_cluster_type",
		WrapToolHandler[ClusterTypeCreateInput, ClusterTypeOutput](
			metrics,
			"create_cluster_type",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ClusterTypeCreateInput,
			) (*mcp.CallToolResult, ClusterTypeOutput, error) {
				return NewCreateClusterTypeHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_cluster_type"],
		metrics,
		"update_cluster_type",
		WrapToolHandler[ClusterTypeUpdateInput, ClusterTypeOutput](
			metrics,
			"update_cluster_type",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ClusterTypeUpdateInput,
			) (*mcp.CallToolResult, ClusterTypeOutput, error) {
				return NewUpdateClusterTypeHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_cluster_type"],
		metrics,
		"delete_cluster_type",
		WrapToolHandler[ClusterTypeDeleteInput, struct{}](
			metrics,
			"delete_cluster_type",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ClusterTypeDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteClusterTypeHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_cluster_group"],
		metrics,
		"create_cluster_group",
		WrapToolHandler[ClusterGroupCreateInput, ClusterGroupOutput](
			metrics,
			"create_cluster_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ClusterGroupCreateInput,
			) (*mcp.CallToolResult, ClusterGroupOutput, error) {
				return NewCreateClusterGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_cluster_group"],
		metrics,
		"update_cluster_group",
		WrapToolHandler[ClusterGroupUpdateInput, ClusterGroupOutput](
			metrics,
			"update_cluster_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ClusterGroupUpdateInput,
			) (*mcp.CallToolResult, ClusterGroupOutput, error) {
				return NewUpdateClusterGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_cluster_group"],
		metrics,
		"delete_cluster_group",
		WrapToolHandler[ClusterGroupDeleteInput, struct{}](
			metrics,
			"delete_cluster_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ClusterGroupDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteClusterGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_circuit_type"],
		metrics,
		"create_circuit_type",
		WrapToolHandler[CircuitTypeCreateInput, CircuitTypeOutput](
			metrics,
			"create_circuit_type",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CircuitTypeCreateInput,
			) (*mcp.CallToolResult, CircuitTypeOutput, error) {
				return NewCreateCircuitTypeHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_circuit_type"],
		metrics,
		"update_circuit_type",
		WrapToolHandler[CircuitTypeUpdateInput, CircuitTypeOutput](
			metrics,
			"update_circuit_type",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CircuitTypeUpdateInput,
			) (*mcp.CallToolResult, CircuitTypeOutput, error) {
				return NewUpdateCircuitTypeHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_circuit_type"],
		metrics,
		"delete_circuit_type",
		WrapToolHandler[CircuitTypeDeleteInput, struct{}](
			metrics,
			"delete_circuit_type",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CircuitTypeDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteCircuitTypeHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_vrf"],
		metrics,
		"create_vrf",
		WrapToolHandler[VrfCreateInput, VrfOutput](
			metrics,
			"create_vrf",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VrfCreateInput,
			) (*mcp.CallToolResult, VrfOutput, error) {
				return NewCreateVrfHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_vrf"],
		metrics,
		"update_vrf",
		WrapToolHandler[VrfUpdateInput, VrfOutput](
			metrics,
			"update_vrf",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VrfUpdateInput,
			) (*mcp.CallToolResult, VrfOutput, error) {
				return NewUpdateVrfHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_vrf"],
		metrics,
		"delete_vrf",
		WrapToolHandler[VrfDeleteInput, struct{}](
			metrics,
			"delete_vrf",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VrfDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteVrfHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_vlan_group"],
		metrics,
		"create_vlan_group",
		WrapToolHandler[VlanGroupCreateInput, VlanGroupOutput](
			metrics,
			"create_vlan_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VlanGroupCreateInput,
			) (*mcp.CallToolResult, VlanGroupOutput, error) {
				return NewCreateVlanGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_vlan_group"],
		metrics,
		"update_vlan_group",
		WrapToolHandler[VlanGroupUpdateInput, VlanGroupOutput](
			metrics,
			"update_vlan_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VlanGroupUpdateInput,
			) (*mcp.CallToolResult, VlanGroupOutput, error) {
				return NewUpdateVlanGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_vlan_group"],
		metrics,
		"delete_vlan_group",
		WrapToolHandler[VlanGroupDeleteInput, struct{}](
			metrics,
			"delete_vlan_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VlanGroupDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteVlanGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_role"],
		metrics,
		"create_role",
		WrapToolHandler[RoleCreateInput, RoleOutput](
			metrics,
			"create_role",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RoleCreateInput,
			) (*mcp.CallToolResult, RoleOutput, error) {
				return NewCreateRoleHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_role"],
		metrics,
		"update_role",
		WrapToolHandler[RoleUpdateInput, RoleOutput](
			metrics,
			"update_role",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RoleUpdateInput,
			) (*mcp.CallToolResult, RoleOutput, error) {
				return NewUpdateRoleHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_role"],
		metrics,
		"delete_role",
		WrapToolHandler[RoleDeleteInput, struct{}](
			metrics,
			"delete_role",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RoleDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteRoleHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_contact"],
		metrics,
		"create_contact",
		WrapToolHandler[ContactCreateInput, ContactOutput](
			metrics,
			"create_contact",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ContactCreateInput,
			) (*mcp.CallToolResult, ContactOutput, error) {
				return NewCreateContactHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_contact"],
		metrics,
		"update_contact",
		WrapToolHandler[ContactUpdateInput, ContactOutput](
			metrics,
			"update_contact",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ContactUpdateInput,
			) (*mcp.CallToolResult, ContactOutput, error) {
				return NewUpdateContactHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_contact"],
		metrics,
		"delete_contact",
		WrapToolHandler[ContactDeleteInput, struct{}](
			metrics,
			"delete_contact",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ContactDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteContactHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_site_group"],
		metrics,
		"create_site_group",
		WrapToolHandler[SiteGroupCreateInput, SiteGroupOutput](
			metrics,
			"create_site_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in SiteGroupCreateInput,
			) (*mcp.CallToolResult, SiteGroupOutput, error) {
				return NewCreateSiteGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_site_group"],
		metrics,
		"update_site_group",
		WrapToolHandler[SiteGroupUpdateInput, SiteGroupOutput](
			metrics,
			"update_site_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in SiteGroupUpdateInput,
			) (*mcp.CallToolResult, SiteGroupOutput, error) {
				return NewUpdateSiteGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_site_group"],
		metrics,
		"delete_site_group",
		WrapToolHandler[SiteGroupDeleteInput, struct{}](
			metrics,
			"delete_site_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in SiteGroupDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteSiteGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_rack_group"],
		metrics,
		"create_rack_group",
		WrapToolHandler[RackGroupCreateInput, RackGroupOutput](
			metrics,
			"create_rack_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RackGroupCreateInput,
			) (*mcp.CallToolResult, RackGroupOutput, error) {
				return NewCreateRackGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_rack_group"],
		metrics,
		"update_rack_group",
		WrapToolHandler[RackGroupUpdateInput, RackGroupOutput](
			metrics,
			"update_rack_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RackGroupUpdateInput,
			) (*mcp.CallToolResult, RackGroupOutput, error) {
				return NewUpdateRackGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_rack_group"],
		metrics,
		"delete_rack_group",
		WrapToolHandler[RackGroupDeleteInput, struct{}](
			metrics,
			"delete_rack_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RackGroupDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteRackGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_rack_role"],
		metrics,
		"create_rack_role",
		WrapToolHandler[RackRoleCreateInput, RackRoleOutput](
			metrics,
			"create_rack_role",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RackRoleCreateInput,
			) (*mcp.CallToolResult, RackRoleOutput, error) {
				return NewCreateRackRoleHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_rack_role"],
		metrics,
		"update_rack_role",
		WrapToolHandler[RackRoleUpdateInput, RackRoleOutput](
			metrics,
			"update_rack_role",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RackRoleUpdateInput,
			) (*mcp.CallToolResult, RackRoleOutput, error) {
				return NewUpdateRackRoleHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_rack_role"],
		metrics,
		"delete_rack_role",
		WrapToolHandler[RackRoleDeleteInput, struct{}](
			metrics,
			"delete_rack_role",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RackRoleDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteRackRoleHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_device_role"],
		metrics,
		"create_device_role",
		WrapToolHandler[DeviceRoleCreateInput, DeviceRoleOutput](
			metrics,
			"create_device_role",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in DeviceRoleCreateInput,
			) (*mcp.CallToolResult, DeviceRoleOutput, error) {
				return NewCreateDeviceRoleHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_device_role"],
		metrics,
		"update_device_role",
		WrapToolHandler[DeviceRoleUpdateInput, DeviceRoleOutput](
			metrics,
			"update_device_role",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in DeviceRoleUpdateInput,
			) (*mcp.CallToolResult, DeviceRoleOutput, error) {
				return NewUpdateDeviceRoleHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_device_role"],
		metrics,
		"delete_device_role",
		WrapToolHandler[DeviceRoleDeleteInput, struct{}](
			metrics,
			"delete_device_role",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in DeviceRoleDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteDeviceRoleHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_platform"],
		metrics,
		"create_platform",
		WrapToolHandler[PlatformCreateInput, PlatformOutput](
			metrics,
			"create_platform",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PlatformCreateInput,
			) (*mcp.CallToolResult, PlatformOutput, error) {
				return NewCreatePlatformHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_platform"],
		metrics,
		"update_platform",
		WrapToolHandler[PlatformUpdateInput, PlatformOutput](
			metrics,
			"update_platform",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PlatformUpdateInput,
			) (*mcp.CallToolResult, PlatformOutput, error) {
				return NewUpdatePlatformHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_platform"],
		metrics,
		"delete_platform",
		WrapToolHandler[PlatformDeleteInput, struct{}](
			metrics,
			"delete_platform",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PlatformDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeletePlatformHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_virtual_chassis"],
		metrics,
		"create_virtual_chassis",
		WrapToolHandler[VirtualChassisCreateInput, VirtualChassisOutput](
			metrics,
			"create_virtual_chassis",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualChassisCreateInput,
			) (*mcp.CallToolResult, VirtualChassisOutput, error) {
				return NewCreateVirtualChassisHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_virtual_chassis"],
		metrics,
		"update_virtual_chassis",
		WrapToolHandler[VirtualChassisUpdateInput, VirtualChassisOutput](
			metrics,
			"update_virtual_chassis",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualChassisUpdateInput,
			) (*mcp.CallToolResult, VirtualChassisOutput, error) {
				return NewUpdateVirtualChassisHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_virtual_chassis"],
		metrics,
		"delete_virtual_chassis",
		WrapToolHandler[VirtualChassisDeleteInput, struct{}](
			metrics,
			"delete_virtual_chassis",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualChassisDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteVirtualChassisHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_power_panel"],
		metrics,
		"create_power_panel",
		WrapToolHandler[PowerPanelCreateInput, PowerPanelOutput](
			metrics,
			"create_power_panel",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PowerPanelCreateInput,
			) (*mcp.CallToolResult, PowerPanelOutput, error) {
				return NewCreatePowerPanelHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_power_panel"],
		metrics,
		"update_power_panel",
		WrapToolHandler[PowerPanelUpdateInput, PowerPanelOutput](
			metrics,
			"update_power_panel",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PowerPanelUpdateInput,
			) (*mcp.CallToolResult, PowerPanelOutput, error) {
				return NewUpdatePowerPanelHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_power_panel"],
		metrics,
		"delete_power_panel",
		WrapToolHandler[PowerPanelDeleteInput, struct{}](
			metrics,
			"delete_power_panel",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PowerPanelDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeletePowerPanelHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_inventory_item_role"],
		metrics,
		"create_inventory_item_role",
		WrapToolHandler[InventoryItemRoleCreateInput, InventoryItemRoleOutput](
			metrics,
			"create_inventory_item_role",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in InventoryItemRoleCreateInput,
			) (*mcp.CallToolResult, InventoryItemRoleOutput, error) {
				return NewCreateInventoryItemRoleHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_inventory_item_role"],
		metrics,
		"update_inventory_item_role",
		WrapToolHandler[InventoryItemRoleUpdateInput, InventoryItemRoleOutput](
			metrics,
			"update_inventory_item_role",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in InventoryItemRoleUpdateInput,
			) (*mcp.CallToolResult, InventoryItemRoleOutput, error) {
				return NewUpdateInventoryItemRoleHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_inventory_item_role"],
		metrics,
		"delete_inventory_item_role",
		WrapToolHandler[InventoryItemRoleDeleteInput, struct{}](
			metrics,
			"delete_inventory_item_role",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in InventoryItemRoleDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteInventoryItemRoleHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_front_ports"],
		metrics,
		"get_front_ports",
		WrapToolHandler[FrontPortsInput, PaginatedOutput[domain.FrontPort]](
			metrics,
			"get_front_ports",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in FrontPortsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.FrontPort], error) {
				return NewGetFrontPortsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_front_port"],
		metrics,
		"create_front_port",
		WrapToolHandler[FrontPortCreateInput, FrontPortOutput](
			metrics,
			"create_front_port",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in FrontPortCreateInput,
			) (*mcp.CallToolResult, FrontPortOutput, error) {
				return NewCreateFrontPortHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_front_port"],
		metrics,
		"update_front_port",
		WrapToolHandler[FrontPortUpdateInput, FrontPortOutput](
			metrics,
			"update_front_port",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in FrontPortUpdateInput,
			) (*mcp.CallToolResult, FrontPortOutput, error) {
				return NewUpdateFrontPortHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_front_port"],
		metrics,
		"delete_front_port",
		WrapToolHandler[FrontPortDeleteInput, struct{}](
			metrics,
			"delete_front_port",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in FrontPortDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteFrontPortHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_rear_ports"],
		metrics,
		"get_rear_ports",
		WrapToolHandler[RearPortsInput, PaginatedOutput[domain.RearPort]](
			metrics,
			"get_rear_ports",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RearPortsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.RearPort], error) {
				return NewGetRearPortsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_rear_port"],
		metrics,
		"create_rear_port",
		WrapToolHandler[RearPortCreateInput, RearPortOutput](
			metrics,
			"create_rear_port",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RearPortCreateInput,
			) (*mcp.CallToolResult, RearPortOutput, error) {
				return NewCreateRearPortHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_rear_port"],
		metrics,
		"update_rear_port",
		WrapToolHandler[RearPortUpdateInput, RearPortOutput](
			metrics,
			"update_rear_port",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RearPortUpdateInput,
			) (*mcp.CallToolResult, RearPortOutput, error) {
				return NewUpdateRearPortHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_rear_port"],
		metrics,
		"delete_rear_port",
		WrapToolHandler[RearPortDeleteInput, struct{}](
			metrics,
			"delete_rear_port",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RearPortDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteRearPortHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_console_ports"],
		metrics,
		"get_console_ports",
		WrapToolHandler[ConsolePortsInput, PaginatedOutput[domain.ConsolePort]](
			metrics,
			"get_console_ports",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ConsolePortsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.ConsolePort], error) {
				return NewGetConsolePortsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_console_port"],
		metrics,
		"create_console_port",
		WrapToolHandler[ConsolePortCreateInput, ConsolePortOutput](
			metrics,
			"create_console_port",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ConsolePortCreateInput,
			) (*mcp.CallToolResult, ConsolePortOutput, error) {
				return NewCreateConsolePortHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_console_port"],
		metrics,
		"update_console_port",
		WrapToolHandler[ConsolePortUpdateInput, ConsolePortOutput](
			metrics,
			"update_console_port",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ConsolePortUpdateInput,
			) (*mcp.CallToolResult, ConsolePortOutput, error) {
				return NewUpdateConsolePortHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_console_port"],
		metrics,
		"delete_console_port",
		WrapToolHandler[ConsolePortDeleteInput, struct{}](
			metrics,
			"delete_console_port",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ConsolePortDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteConsolePortHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_console_server_ports"],
		metrics,
		"get_console_server_ports",
		WrapToolHandler[ConsoleServerPortsInput, PaginatedOutput[domain.ConsoleServerPort]](
			metrics,
			"get_console_server_ports",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ConsoleServerPortsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.ConsoleServerPort], error) {
				return NewGetConsoleServerPortsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_console_server_port"],
		metrics,
		"create_console_server_port",
		WrapToolHandler[ConsoleServerPortCreateInput, ConsoleServerPortOutput](
			metrics,
			"create_console_server_port",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ConsoleServerPortCreateInput,
			) (*mcp.CallToolResult, ConsoleServerPortOutput, error) {
				return NewCreateConsoleServerPortHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_console_server_port"],
		metrics,
		"update_console_server_port",
		WrapToolHandler[ConsoleServerPortUpdateInput, ConsoleServerPortOutput](
			metrics,
			"update_console_server_port",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ConsoleServerPortUpdateInput,
			) (*mcp.CallToolResult, ConsoleServerPortOutput, error) {
				return NewUpdateConsoleServerPortHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_console_server_port"],
		metrics,
		"delete_console_server_port",
		WrapToolHandler[ConsoleServerPortDeleteInput, struct{}](
			metrics,
			"delete_console_server_port",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ConsoleServerPortDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteConsoleServerPortHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_power_ports"],
		metrics,
		"get_power_ports",
		WrapToolHandler[PowerPortsInput, PaginatedOutput[domain.PowerPort]](
			metrics,
			"get_power_ports",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PowerPortsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.PowerPort], error) {
				return NewGetPowerPortsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_power_port"],
		metrics,
		"create_power_port",
		WrapToolHandler[PowerPortCreateInput, PowerPortOutput](
			metrics,
			"create_power_port",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PowerPortCreateInput,
			) (*mcp.CallToolResult, PowerPortOutput, error) {
				return NewCreatePowerPortHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_power_port"],
		metrics,
		"update_power_port",
		WrapToolHandler[PowerPortUpdateInput, PowerPortOutput](
			metrics,
			"update_power_port",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PowerPortUpdateInput,
			) (*mcp.CallToolResult, PowerPortOutput, error) {
				return NewUpdatePowerPortHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_power_port"],
		metrics,
		"delete_power_port",
		WrapToolHandler[PowerPortDeleteInput, struct{}](
			metrics,
			"delete_power_port",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PowerPortDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeletePowerPortHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_power_outlets"],
		metrics,
		"get_power_outlets",
		WrapToolHandler[PowerOutletsInput, PaginatedOutput[domain.PowerOutlet]](
			metrics,
			"get_power_outlets",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PowerOutletsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.PowerOutlet], error) {
				return NewGetPowerOutletsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_power_outlet"],
		metrics,
		"create_power_outlet",
		WrapToolHandler[PowerOutletCreateInput, PowerOutletOutput](
			metrics,
			"create_power_outlet",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PowerOutletCreateInput,
			) (*mcp.CallToolResult, PowerOutletOutput, error) {
				return NewCreatePowerOutletHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_power_outlet"],
		metrics,
		"update_power_outlet",
		WrapToolHandler[PowerOutletUpdateInput, PowerOutletOutput](
			metrics,
			"update_power_outlet",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PowerOutletUpdateInput,
			) (*mcp.CallToolResult, PowerOutletOutput, error) {
				return NewUpdatePowerOutletHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_power_outlet"],
		metrics,
		"delete_power_outlet",
		WrapToolHandler[PowerOutletDeleteInput, struct{}](
			metrics,
			"delete_power_outlet",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PowerOutletDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeletePowerOutletHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_power_feeds"],
		metrics,
		"get_power_feeds",
		WrapToolHandler[PowerFeedsInput, PaginatedOutput[domain.PowerFeed]](
			metrics,
			"get_power_feeds",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PowerFeedsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.PowerFeed], error) {
				return NewGetPowerFeedsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_power_feed"],
		metrics,
		"create_power_feed",
		WrapToolHandler[PowerFeedCreateInput, PowerFeedOutput](
			metrics,
			"create_power_feed",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PowerFeedCreateInput,
			) (*mcp.CallToolResult, PowerFeedOutput, error) {
				return NewCreatePowerFeedHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_power_feed"],
		metrics,
		"update_power_feed",
		WrapToolHandler[PowerFeedUpdateInput, PowerFeedOutput](
			metrics,
			"update_power_feed",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PowerFeedUpdateInput,
			) (*mcp.CallToolResult, PowerFeedOutput, error) {
				return NewUpdatePowerFeedHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_power_feed"],
		metrics,
		"delete_power_feed",
		WrapToolHandler[PowerFeedDeleteInput, struct{}](
			metrics,
			"delete_power_feed",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in PowerFeedDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeletePowerFeedHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_device_bays"],
		metrics,
		"get_device_bays",
		WrapToolHandler[DeviceBaysInput, PaginatedOutput[domain.DeviceBay]](
			metrics,
			"get_device_bays",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in DeviceBaysInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.DeviceBay], error) {
				return NewGetDeviceBaysHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_device_bay"],
		metrics,
		"create_device_bay",
		WrapToolHandler[DeviceBayCreateInput, DeviceBayOutput](
			metrics,
			"create_device_bay",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in DeviceBayCreateInput,
			) (*mcp.CallToolResult, DeviceBayOutput, error) {
				return NewCreateDeviceBayHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_device_bay"],
		metrics,
		"update_device_bay",
		WrapToolHandler[DeviceBayUpdateInput, DeviceBayOutput](
			metrics,
			"update_device_bay",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in DeviceBayUpdateInput,
			) (*mcp.CallToolResult, DeviceBayOutput, error) {
				return NewUpdateDeviceBayHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_device_bay"],
		metrics,
		"delete_device_bay",
		WrapToolHandler[DeviceBayDeleteInput, struct{}](
			metrics,
			"delete_device_bay",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in DeviceBayDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteDeviceBayHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_rack_reservations"],
		metrics,
		"get_rack_reservations",
		WrapToolHandler[RackReservationsInput, PaginatedOutput[domain.RackReservation]](
			metrics,
			"get_rack_reservations",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RackReservationsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.RackReservation], error) {
				return NewGetRackReservationsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_rack_reservation"],
		metrics,
		"create_rack_reservation",
		WrapToolHandler[RackReservationCreateInput, RackReservationOutput](
			metrics,
			"create_rack_reservation",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RackReservationCreateInput,
			) (*mcp.CallToolResult, RackReservationOutput, error) {
				return NewCreateRackReservationHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_rack_reservation"],
		metrics,
		"update_rack_reservation",
		WrapToolHandler[RackReservationUpdateInput, RackReservationOutput](
			metrics,
			"update_rack_reservation",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RackReservationUpdateInput,
			) (*mcp.CallToolResult, RackReservationOutput, error) {
				return NewUpdateRackReservationHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_rack_reservation"],
		metrics,
		"delete_rack_reservation",
		WrapToolHandler[RackReservationDeleteInput, struct{}](
			metrics,
			"delete_rack_reservation",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RackReservationDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteRackReservationHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_inventory_items"],
		metrics,
		"get_inventory_items",
		WrapToolHandler[InventoryItemsInput, PaginatedOutput[domain.InventoryItem]](
			metrics,
			"get_inventory_items",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in InventoryItemsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.InventoryItem], error) {
				return NewGetInventoryItemsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_inventory_item"],
		metrics,
		"create_inventory_item",
		WrapToolHandler[InventoryItemCreateInput, InventoryItemOutput](
			metrics,
			"create_inventory_item",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in InventoryItemCreateInput,
			) (*mcp.CallToolResult, InventoryItemOutput, error) {
				return NewCreateInventoryItemHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_inventory_item"],
		metrics,
		"update_inventory_item",
		WrapToolHandler[InventoryItemUpdateInput, InventoryItemOutput](
			metrics,
			"update_inventory_item",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in InventoryItemUpdateInput,
			) (*mcp.CallToolResult, InventoryItemOutput, error) {
				return NewUpdateInventoryItemHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_inventory_item"],
		metrics,
		"delete_inventory_item",
		WrapToolHandler[InventoryItemDeleteInput, struct{}](
			metrics,
			"delete_inventory_item",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in InventoryItemDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteInventoryItemHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_mac_addresses"],
		metrics,
		"get_mac_addresses",
		WrapToolHandler[MACAddressesInput, PaginatedOutput[domain.MACAddress]](
			metrics,
			"get_mac_addresses",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in MACAddressesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.MACAddress], error) {
				return NewGetMACAddressesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_mac_address"],
		metrics,
		"create_mac_address",
		WrapToolHandler[MACAddressCreateInput, MACAddressOutput](
			metrics,
			"create_mac_address",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in MACAddressCreateInput,
			) (*mcp.CallToolResult, MACAddressOutput, error) {
				return NewCreateMACAddressHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_mac_address"],
		metrics,
		"update_mac_address",
		WrapToolHandler[MACAddressUpdateInput, MACAddressOutput](
			metrics,
			"update_mac_address",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in MACAddressUpdateInput,
			) (*mcp.CallToolResult, MACAddressOutput, error) {
				return NewUpdateMACAddressHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_mac_address"],
		metrics,
		"delete_mac_address",
		WrapToolHandler[MACAddressDeleteInput, struct{}](
			metrics,
			"delete_mac_address",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in MACAddressDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteMACAddressHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_rirs"],
		metrics,
		"get_rirs",
		WrapToolHandler[RIRsInput, PaginatedOutput[domain.RIR]](
			metrics,
			"get_rirs",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RIRsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.RIR], error) {
				return NewGetRIRsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_rir"],
		metrics,
		"create_rir",
		WrapToolHandler[RIRCreateInput, RIROutput](
			metrics,
			"create_rir",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RIRCreateInput,
			) (*mcp.CallToolResult, RIROutput, error) {
				return NewCreateRIRHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_rir"],
		metrics,
		"update_rir",
		WrapToolHandler[RIRUpdateInput, RIROutput](
			metrics,
			"update_rir",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RIRUpdateInput,
			) (*mcp.CallToolResult, RIROutput, error) {
				return NewUpdateRIRHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_rir"],
		metrics,
		"delete_rir",
		WrapToolHandler[RIRDeleteInput, struct{}](
			metrics,
			"delete_rir",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RIRDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteRIRHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_contact_groups"],
		metrics,
		"get_contact_groups",
		WrapToolHandler[ContactGroupsInput, PaginatedOutput[domain.ContactGroup]](
			metrics,
			"get_contact_groups",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ContactGroupsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.ContactGroup], error) {
				return NewGetContactGroupsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_contact_group"],
		metrics,
		"create_contact_group",
		WrapToolHandler[ContactGroupCreateInput, ContactGroupOutput](
			metrics,
			"create_contact_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ContactGroupCreateInput,
			) (*mcp.CallToolResult, ContactGroupOutput, error) {
				return NewCreateContactGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_contact_group"],
		metrics,
		"update_contact_group",
		WrapToolHandler[ContactGroupUpdateInput, ContactGroupOutput](
			metrics,
			"update_contact_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ContactGroupUpdateInput,
			) (*mcp.CallToolResult, ContactGroupOutput, error) {
				return NewUpdateContactGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_contact_group"],
		metrics,
		"delete_contact_group",
		WrapToolHandler[ContactGroupDeleteInput, struct{}](
			metrics,
			"delete_contact_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ContactGroupDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteContactGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_contact_roles"],
		metrics,
		"get_contact_roles",
		WrapToolHandler[ContactRolesInput, PaginatedOutput[domain.ContactRole]](
			metrics,
			"get_contact_roles",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ContactRolesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.ContactRole], error) {
				return NewGetContactRolesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_contact_role"],
		metrics,
		"create_contact_role",
		WrapToolHandler[ContactRoleCreateInput, ContactRoleOutput](
			metrics,
			"create_contact_role",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ContactRoleCreateInput,
			) (*mcp.CallToolResult, ContactRoleOutput, error) {
				return NewCreateContactRoleHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_contact_role"],
		metrics,
		"update_contact_role",
		WrapToolHandler[ContactRoleUpdateInput, ContactRoleOutput](
			metrics,
			"update_contact_role",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ContactRoleUpdateInput,
			) (*mcp.CallToolResult, ContactRoleOutput, error) {
				return NewUpdateContactRoleHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_contact_role"],
		metrics,
		"delete_contact_role",
		WrapToolHandler[ContactRoleDeleteInput, struct{}](
			metrics,
			"delete_contact_role",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ContactRoleDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteContactRoleHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_virtual_disks"],
		metrics,
		"get_virtual_disks",
		WrapToolHandler[VirtualDisksInput, PaginatedOutput[domain.VirtualDisk]](
			metrics,
			"get_virtual_disks",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualDisksInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.VirtualDisk], error) {
				return NewGetVirtualDisksHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_virtual_disk"],
		metrics,
		"create_virtual_disk",
		WrapToolHandler[VirtualDiskCreateInput, VirtualDiskOutput](
			metrics,
			"create_virtual_disk",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualDiskCreateInput,
			) (*mcp.CallToolResult, VirtualDiskOutput, error) {
				return NewCreateVirtualDiskHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_virtual_disk"],
		metrics,
		"update_virtual_disk",
		WrapToolHandler[VirtualDiskUpdateInput, VirtualDiskOutput](
			metrics,
			"update_virtual_disk",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualDiskUpdateInput,
			) (*mcp.CallToolResult, VirtualDiskOutput, error) {
				return NewUpdateVirtualDiskHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_virtual_disk"],
		metrics,
		"delete_virtual_disk",
		WrapToolHandler[VirtualDiskDeleteInput, struct{}](
			metrics,
			"delete_virtual_disk",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualDiskDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteVirtualDiskHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_virtual_machine_types"],
		metrics,
		"get_virtual_machine_types",
		WrapToolHandler[VirtualMachineTypesInput, PaginatedOutput[domain.VirtualMachineType]](
			metrics,
			"get_virtual_machine_types",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualMachineTypesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.VirtualMachineType], error) {
				return NewGetVirtualMachineTypesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_virtual_machine_type"],
		metrics,
		"create_virtual_machine_type",
		WrapToolHandler[VirtualMachineTypeCreateInput, VirtualMachineTypeOutput](
			metrics,
			"create_virtual_machine_type",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualMachineTypeCreateInput,
			) (*mcp.CallToolResult, VirtualMachineTypeOutput, error) {
				return NewCreateVirtualMachineTypeHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_virtual_machine_type"],
		metrics,
		"update_virtual_machine_type",
		WrapToolHandler[VirtualMachineTypeUpdateInput, VirtualMachineTypeOutput](
			metrics,
			"update_virtual_machine_type",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualMachineTypeUpdateInput,
			) (*mcp.CallToolResult, VirtualMachineTypeOutput, error) {
				return NewUpdateVirtualMachineTypeHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_virtual_machine_type"],
		metrics,
		"delete_virtual_machine_type",
		WrapToolHandler[VirtualMachineTypeDeleteInput, struct{}](
			metrics,
			"delete_virtual_machine_type",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualMachineTypeDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteVirtualMachineTypeHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_aggregates"],
		metrics,
		"get_aggregates",
		WrapToolHandler[AggregatesInput, PaginatedOutput[domain.Aggregate]](
			metrics,
			"get_aggregates",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in AggregatesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.Aggregate], error) {
				return NewGetAggregatesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_aggregate"],
		metrics,
		"create_aggregate",
		WrapToolHandler[AggregateCreateInput, AggregateOutput](
			metrics,
			"create_aggregate",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in AggregateCreateInput,
			) (*mcp.CallToolResult, AggregateOutput, error) {
				return NewCreateAggregateHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_aggregate"],
		metrics,
		"update_aggregate",
		WrapToolHandler[AggregateUpdateInput, AggregateOutput](
			metrics,
			"update_aggregate",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in AggregateUpdateInput,
			) (*mcp.CallToolResult, AggregateOutput, error) {
				return NewUpdateAggregateHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_aggregate"],
		metrics,
		"delete_aggregate",
		WrapToolHandler[AggregateDeleteInput, struct{}](
			metrics,
			"delete_aggregate",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in AggregateDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteAggregateHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_ip_ranges"],
		metrics,
		"get_ip_ranges",
		WrapToolHandler[IPRangesInput, PaginatedOutput[domain.IPRange]](
			metrics,
			"get_ip_ranges",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in IPRangesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.IPRange], error) {
				return NewGetIPRangesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_ip_range"],
		metrics,
		"create_ip_range",
		WrapToolHandler[IPRangeCreateInput, IPRangeOutput](
			metrics,
			"create_ip_range",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in IPRangeCreateInput,
			) (*mcp.CallToolResult, IPRangeOutput, error) {
				return NewCreateIPRangeHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_ip_range"],
		metrics,
		"update_ip_range",
		WrapToolHandler[IPRangeUpdateInput, IPRangeOutput](
			metrics,
			"update_ip_range",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in IPRangeUpdateInput,
			) (*mcp.CallToolResult, IPRangeOutput, error) {
				return NewUpdateIPRangeHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_ip_range"],
		metrics,
		"delete_ip_range",
		WrapToolHandler[IPRangeDeleteInput, struct{}](
			metrics,
			"delete_ip_range",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in IPRangeDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteIPRangeHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_asns"],
		metrics,
		"get_asns",
		WrapToolHandler[ASNsInput, PaginatedOutput[domain.ASN]](
			metrics,
			"get_asns",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ASNsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.ASN], error) {
				return NewGetASNsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_asn"],
		metrics,
		"create_asn",
		WrapToolHandler[ASNCreateInput, ASNOutput](
			metrics,
			"create_asn",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ASNCreateInput,
			) (*mcp.CallToolResult, ASNOutput, error) {
				return NewCreateASNHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_asn"],
		metrics,
		"update_asn",
		WrapToolHandler[ASNUpdateInput, ASNOutput](
			metrics,
			"update_asn",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ASNUpdateInput,
			) (*mcp.CallToolResult, ASNOutput, error) {
				return NewUpdateASNHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_asn"],
		metrics,
		"delete_asn",
		WrapToolHandler[ASNDeleteInput, struct{}](
			metrics,
			"delete_asn",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ASNDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteASNHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_asn_ranges"],
		metrics,
		"get_asn_ranges",
		WrapToolHandler[ASNRangesInput, PaginatedOutput[domain.ASNRange]](
			metrics,
			"get_asn_ranges",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ASNRangesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.ASNRange], error) {
				return NewGetASNRangesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_asn_range"],
		metrics,
		"create_asn_range",
		WrapToolHandler[ASNRangeCreateInput, ASNRangeOutput](
			metrics,
			"create_asn_range",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ASNRangeCreateInput,
			) (*mcp.CallToolResult, ASNRangeOutput, error) {
				return NewCreateASNRangeHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_asn_range"],
		metrics,
		"update_asn_range",
		WrapToolHandler[ASNRangeUpdateInput, ASNRangeOutput](
			metrics,
			"update_asn_range",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ASNRangeUpdateInput,
			) (*mcp.CallToolResult, ASNRangeOutput, error) {
				return NewUpdateASNRangeHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_asn_range"],
		metrics,
		"delete_asn_range",
		WrapToolHandler[ASNRangeDeleteInput, struct{}](
			metrics,
			"delete_asn_range",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ASNRangeDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteASNRangeHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_route_targets"],
		metrics,
		"get_route_targets",
		WrapToolHandler[RouteTargetsInput, PaginatedOutput[domain.RouteTarget]](
			metrics,
			"get_route_targets",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RouteTargetsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.RouteTarget], error) {
				return NewGetRouteTargetsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_route_target"],
		metrics,
		"create_route_target",
		WrapToolHandler[RouteTargetCreateInput, RouteTargetOutput](
			metrics,
			"create_route_target",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RouteTargetCreateInput,
			) (*mcp.CallToolResult, RouteTargetOutput, error) {
				return NewCreateRouteTargetHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_route_target"],
		metrics,
		"update_route_target",
		WrapToolHandler[RouteTargetUpdateInput, RouteTargetOutput](
			metrics,
			"update_route_target",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RouteTargetUpdateInput,
			) (*mcp.CallToolResult, RouteTargetOutput, error) {
				return NewUpdateRouteTargetHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_route_target"],
		metrics,
		"delete_route_target",
		WrapToolHandler[RouteTargetDeleteInput, struct{}](
			metrics,
			"delete_route_target",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RouteTargetDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteRouteTargetHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_services"],
		metrics,
		"get_services",
		WrapToolHandler[ServicesInput, PaginatedOutput[domain.Service]](
			metrics,
			"get_services",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ServicesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.Service], error) {
				return NewGetServicesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_service"],
		metrics,
		"create_service",
		WrapToolHandler[ServiceCreateInput, ServiceOutput](
			metrics,
			"create_service",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ServiceCreateInput,
			) (*mcp.CallToolResult, ServiceOutput, error) {
				return NewCreateServiceHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_service"],
		metrics,
		"update_service",
		WrapToolHandler[ServiceUpdateInput, ServiceOutput](
			metrics,
			"update_service",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ServiceUpdateInput,
			) (*mcp.CallToolResult, ServiceOutput, error) {
				return NewUpdateServiceHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_service"],
		metrics,
		"delete_service",
		WrapToolHandler[ServiceDeleteInput, struct{}](
			metrics,
			"delete_service",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ServiceDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteServiceHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_service_templates"],
		metrics,
		"get_service_templates",
		WrapToolHandler[ServiceTemplatesInput, PaginatedOutput[domain.ServiceTemplate]](
			metrics,
			"get_service_templates",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ServiceTemplatesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.ServiceTemplate], error) {
				return NewGetServiceTemplatesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_service_template"],
		metrics,
		"create_service_template",
		WrapToolHandler[ServiceTemplateCreateInput, ServiceTemplateOutput](
			metrics,
			"create_service_template",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ServiceTemplateCreateInput,
			) (*mcp.CallToolResult, ServiceTemplateOutput, error) {
				return NewCreateServiceTemplateHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_service_template"],
		metrics,
		"update_service_template",
		WrapToolHandler[ServiceTemplateUpdateInput, ServiceTemplateOutput](
			metrics,
			"update_service_template",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ServiceTemplateUpdateInput,
			) (*mcp.CallToolResult, ServiceTemplateOutput, error) {
				return NewUpdateServiceTemplateHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_service_template"],
		metrics,
		"delete_service_template",
		WrapToolHandler[ServiceTemplateDeleteInput, struct{}](
			metrics,
			"delete_service_template",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ServiceTemplateDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteServiceTemplateHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_fhrp_groups"],
		metrics,
		"get_fhrp_groups",
		WrapToolHandler[FHRPGroupsInput, PaginatedOutput[domain.FHRPGroup]](
			metrics,
			"get_fhrp_groups",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in FHRPGroupsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.FHRPGroup], error) {
				return NewGetFHRPGroupsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_fhrp_group"],
		metrics,
		"create_fhrp_group",
		WrapToolHandler[FHRPGroupCreateInput, FHRPGroupOutput](
			metrics,
			"create_fhrp_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in FHRPGroupCreateInput,
			) (*mcp.CallToolResult, FHRPGroupOutput, error) {
				return NewCreateFHRPGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_fhrp_group"],
		metrics,
		"update_fhrp_group",
		WrapToolHandler[FHRPGroupUpdateInput, FHRPGroupOutput](
			metrics,
			"update_fhrp_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in FHRPGroupUpdateInput,
			) (*mcp.CallToolResult, FHRPGroupOutput, error) {
				return NewUpdateFHRPGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_fhrp_group"],
		metrics,
		"delete_fhrp_group",
		WrapToolHandler[FHRPGroupDeleteInput, struct{}](
			metrics,
			"delete_fhrp_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in FHRPGroupDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteFHRPGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_fhrp_group_assignments"],
		metrics,
		"get_fhrp_group_assignments",
		WrapToolHandler[FHRPGroupAssignmentsInput, PaginatedOutput[domain.FHRPGroupAssignment]](
			metrics,
			"get_fhrp_group_assignments",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in FHRPGroupAssignmentsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.FHRPGroupAssignment], error) {
				return NewGetFHRPGroupAssignmentsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_fhrp_group_assignment"],
		metrics,
		"create_fhrp_group_assignment",
		WrapToolHandler[FHRPGroupAssignmentCreateInput, FHRPGroupAssignmentOutput](
			metrics,
			"create_fhrp_group_assignment",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in FHRPGroupAssignmentCreateInput,
			) (*mcp.CallToolResult, FHRPGroupAssignmentOutput, error) {
				return NewCreateFHRPGroupAssignmentHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_fhrp_group_assignment"],
		metrics,
		"update_fhrp_group_assignment",
		WrapToolHandler[FHRPGroupAssignmentUpdateInput, FHRPGroupAssignmentOutput](
			metrics,
			"update_fhrp_group_assignment",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in FHRPGroupAssignmentUpdateInput,
			) (*mcp.CallToolResult, FHRPGroupAssignmentOutput, error) {
				return NewUpdateFHRPGroupAssignmentHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_fhrp_group_assignment"],
		metrics,
		"delete_fhrp_group_assignment",
		WrapToolHandler[FHRPGroupAssignmentDeleteInput, struct{}](
			metrics,
			"delete_fhrp_group_assignment",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in FHRPGroupAssignmentDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteFHRPGroupAssignmentHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_vlan_translation_policies"],
		metrics,
		"get_vlan_translation_policies",
		WrapToolHandler[VLANTranslationPoliciesInput, PaginatedOutput[domain.VLANTranslationPolicy]](
			metrics,
			"get_vlan_translation_policies",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VLANTranslationPoliciesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.VLANTranslationPolicy], error) {
				return NewGetVLANTranslationPoliciesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_vlan_translation_policy"],
		metrics,
		"create_vlan_translation_policy",
		WrapToolHandler[VLANTranslationPolicyCreateInput, VLANTranslationPolicyOutput](
			metrics,
			"create_vlan_translation_policy",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VLANTranslationPolicyCreateInput,
			) (*mcp.CallToolResult, VLANTranslationPolicyOutput, error) {
				return NewCreateVLANTranslationPolicyHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_vlan_translation_policy"],
		metrics,
		"update_vlan_translation_policy",
		WrapToolHandler[VLANTranslationPolicyUpdateInput, VLANTranslationPolicyOutput](
			metrics,
			"update_vlan_translation_policy",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VLANTranslationPolicyUpdateInput,
			) (*mcp.CallToolResult, VLANTranslationPolicyOutput, error) {
				return NewUpdateVLANTranslationPolicyHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_vlan_translation_policy"],
		metrics,
		"delete_vlan_translation_policy",
		WrapToolHandler[VLANTranslationPolicyDeleteInput, struct{}](
			metrics,
			"delete_vlan_translation_policy",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VLANTranslationPolicyDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteVLANTranslationPolicyHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_vlan_translation_rules"],
		metrics,
		"get_vlan_translation_rules",
		WrapToolHandler[VLANTranslationRulesInput, PaginatedOutput[domain.VLANTranslationRule]](
			metrics,
			"get_vlan_translation_rules",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VLANTranslationRulesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.VLANTranslationRule], error) {
				return NewGetVLANTranslationRulesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_vlan_translation_rule"],
		metrics,
		"create_vlan_translation_rule",
		WrapToolHandler[VLANTranslationRuleCreateInput, VLANTranslationRuleOutput](
			metrics,
			"create_vlan_translation_rule",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VLANTranslationRuleCreateInput,
			) (*mcp.CallToolResult, VLANTranslationRuleOutput, error) {
				return NewCreateVLANTranslationRuleHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_vlan_translation_rule"],
		metrics,
		"update_vlan_translation_rule",
		WrapToolHandler[VLANTranslationRuleUpdateInput, VLANTranslationRuleOutput](
			metrics,
			"update_vlan_translation_rule",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VLANTranslationRuleUpdateInput,
			) (*mcp.CallToolResult, VLANTranslationRuleOutput, error) {
				return NewUpdateVLANTranslationRuleHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_vlan_translation_rule"],
		metrics,
		"delete_vlan_translation_rule",
		WrapToolHandler[VLANTranslationRuleDeleteInput, struct{}](
			metrics,
			"delete_vlan_translation_rule",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VLANTranslationRuleDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteVLANTranslationRuleHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_provider_accounts"],
		metrics,
		"get_provider_accounts",
		WrapToolHandler[ProviderAccountsInput, PaginatedOutput[domain.ProviderAccount]](
			metrics,
			"get_provider_accounts",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ProviderAccountsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.ProviderAccount], error) {
				return NewGetProviderAccountsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_provider_account"],
		metrics,
		"create_provider_account",
		WrapToolHandler[ProviderAccountCreateInput, ProviderAccountOutput](
			metrics,
			"create_provider_account",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ProviderAccountCreateInput,
			) (*mcp.CallToolResult, ProviderAccountOutput, error) {
				return NewCreateProviderAccountHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_provider_account"],
		metrics,
		"update_provider_account",
		WrapToolHandler[ProviderAccountUpdateInput, ProviderAccountOutput](
			metrics,
			"update_provider_account",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ProviderAccountUpdateInput,
			) (*mcp.CallToolResult, ProviderAccountOutput, error) {
				return NewUpdateProviderAccountHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_provider_account"],
		metrics,
		"delete_provider_account",
		WrapToolHandler[ProviderAccountDeleteInput, struct{}](
			metrics,
			"delete_provider_account",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ProviderAccountDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteProviderAccountHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_provider_networks"],
		metrics,
		"get_provider_networks",
		WrapToolHandler[ProviderNetworksInput, PaginatedOutput[domain.ProviderNetwork]](
			metrics,
			"get_provider_networks",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ProviderNetworksInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.ProviderNetwork], error) {
				return NewGetProviderNetworksHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_provider_network"],
		metrics,
		"create_provider_network",
		WrapToolHandler[ProviderNetworkCreateInput, ProviderNetworkOutput](
			metrics,
			"create_provider_network",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ProviderNetworkCreateInput,
			) (*mcp.CallToolResult, ProviderNetworkOutput, error) {
				return NewCreateProviderNetworkHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_provider_network"],
		metrics,
		"update_provider_network",
		WrapToolHandler[ProviderNetworkUpdateInput, ProviderNetworkOutput](
			metrics,
			"update_provider_network",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ProviderNetworkUpdateInput,
			) (*mcp.CallToolResult, ProviderNetworkOutput, error) {
				return NewUpdateProviderNetworkHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_provider_network"],
		metrics,
		"delete_provider_network",
		WrapToolHandler[ProviderNetworkDeleteInput, struct{}](
			metrics,
			"delete_provider_network",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in ProviderNetworkDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteProviderNetworkHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_circuit_groups"],
		metrics,
		"get_circuit_groups",
		WrapToolHandler[CircuitGroupsInput, PaginatedOutput[domain.CircuitGroup]](
			metrics,
			"get_circuit_groups",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CircuitGroupsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.CircuitGroup], error) {
				return NewGetCircuitGroupsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_circuit_group"],
		metrics,
		"create_circuit_group",
		WrapToolHandler[CircuitGroupCreateInput, CircuitGroupOutput](
			metrics,
			"create_circuit_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CircuitGroupCreateInput,
			) (*mcp.CallToolResult, CircuitGroupOutput, error) {
				return NewCreateCircuitGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_circuit_group"],
		metrics,
		"update_circuit_group",
		WrapToolHandler[CircuitGroupUpdateInput, CircuitGroupOutput](
			metrics,
			"update_circuit_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CircuitGroupUpdateInput,
			) (*mcp.CallToolResult, CircuitGroupOutput, error) {
				return NewUpdateCircuitGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_circuit_group"],
		metrics,
		"delete_circuit_group",
		WrapToolHandler[CircuitGroupDeleteInput, struct{}](
			metrics,
			"delete_circuit_group",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CircuitGroupDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteCircuitGroupHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_circuit_group_assignments"],
		metrics,
		"get_circuit_group_assignments",
		WrapToolHandler[CircuitGroupAssignmentsInput, PaginatedOutput[domain.CircuitGroupAssignment]](
			metrics,
			"get_circuit_group_assignments",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CircuitGroupAssignmentsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.CircuitGroupAssignment], error) {
				return NewGetCircuitGroupAssignmentsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_circuit_group_assignment"],
		metrics,
		"create_circuit_group_assignment",
		WrapToolHandler[CircuitGroupAssignmentCreateInput, CircuitGroupAssignmentOutput](
			metrics,
			"create_circuit_group_assignment",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CircuitGroupAssignmentCreateInput,
			) (*mcp.CallToolResult, CircuitGroupAssignmentOutput, error) {
				return NewCreateCircuitGroupAssignmentHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_circuit_group_assignment"],
		metrics,
		"update_circuit_group_assignment",
		WrapToolHandler[CircuitGroupAssignmentUpdateInput, CircuitGroupAssignmentOutput](
			metrics,
			"update_circuit_group_assignment",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CircuitGroupAssignmentUpdateInput,
			) (*mcp.CallToolResult, CircuitGroupAssignmentOutput, error) {
				return NewUpdateCircuitGroupAssignmentHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_circuit_group_assignment"],
		metrics,
		"delete_circuit_group_assignment",
		WrapToolHandler[CircuitGroupAssignmentDeleteInput, struct{}](
			metrics,
			"delete_circuit_group_assignment",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CircuitGroupAssignmentDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteCircuitGroupAssignmentHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_virtual_circuits"],
		metrics,
		"get_virtual_circuits",
		WrapToolHandler[VirtualCircuitsInput, PaginatedOutput[domain.VirtualCircuit]](
			metrics,
			"get_virtual_circuits",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualCircuitsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.VirtualCircuit], error) {
				return NewGetVirtualCircuitsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_virtual_circuit"],
		metrics,
		"create_virtual_circuit",
		WrapToolHandler[VirtualCircuitCreateInput, VirtualCircuitOutput](
			metrics,
			"create_virtual_circuit",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualCircuitCreateInput,
			) (*mcp.CallToolResult, VirtualCircuitOutput, error) {
				return NewCreateVirtualCircuitHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_virtual_circuit"],
		metrics,
		"update_virtual_circuit",
		WrapToolHandler[VirtualCircuitUpdateInput, VirtualCircuitOutput](
			metrics,
			"update_virtual_circuit",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualCircuitUpdateInput,
			) (*mcp.CallToolResult, VirtualCircuitOutput, error) {
				return NewUpdateVirtualCircuitHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_virtual_circuit"],
		metrics,
		"delete_virtual_circuit",
		WrapToolHandler[VirtualCircuitDeleteInput, struct{}](
			metrics,
			"delete_virtual_circuit",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualCircuitDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteVirtualCircuitHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_virtual_circuit_types"],
		metrics,
		"get_virtual_circuit_types",
		WrapToolHandler[VirtualCircuitTypesInput, PaginatedOutput[domain.VirtualCircuitType]](
			metrics,
			"get_virtual_circuit_types",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualCircuitTypesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.VirtualCircuitType], error) {
				return NewGetVirtualCircuitTypesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_virtual_circuit_type"],
		metrics,
		"create_virtual_circuit_type",
		WrapToolHandler[VirtualCircuitTypeCreateInput, VirtualCircuitTypeOutput](
			metrics,
			"create_virtual_circuit_type",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualCircuitTypeCreateInput,
			) (*mcp.CallToolResult, VirtualCircuitTypeOutput, error) {
				return NewCreateVirtualCircuitTypeHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_virtual_circuit_type"],
		metrics,
		"update_virtual_circuit_type",
		WrapToolHandler[VirtualCircuitTypeUpdateInput, VirtualCircuitTypeOutput](
			metrics,
			"update_virtual_circuit_type",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualCircuitTypeUpdateInput,
			) (*mcp.CallToolResult, VirtualCircuitTypeOutput, error) {
				return NewUpdateVirtualCircuitTypeHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_virtual_circuit_type"],
		metrics,
		"delete_virtual_circuit_type",
		WrapToolHandler[VirtualCircuitTypeDeleteInput, struct{}](
			metrics,
			"delete_virtual_circuit_type",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualCircuitTypeDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteVirtualCircuitTypeHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_virtual_circuit_terminations"],
		metrics,
		"get_virtual_circuit_terminations",
		WrapToolHandler[VirtualCircuitTerminationsInput, PaginatedOutput[domain.VirtualCircuitTermination]](
			metrics,
			"get_virtual_circuit_terminations",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualCircuitTerminationsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.VirtualCircuitTermination], error) {
				return NewGetVirtualCircuitTerminationsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_virtual_circuit_termination"],
		metrics,
		"create_virtual_circuit_termination",
		WrapToolHandler[VirtualCircuitTerminationCreateInput, VirtualCircuitTerminationOutput](
			metrics,
			"create_virtual_circuit_termination",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualCircuitTerminationCreateInput,
			) (*mcp.CallToolResult, VirtualCircuitTerminationOutput, error) {
				return NewCreateVirtualCircuitTerminationHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_virtual_circuit_termination"],
		metrics,
		"update_virtual_circuit_termination",
		WrapToolHandler[VirtualCircuitTerminationUpdateInput, VirtualCircuitTerminationOutput](
			metrics,
			"update_virtual_circuit_termination",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualCircuitTerminationUpdateInput,
			) (*mcp.CallToolResult, VirtualCircuitTerminationOutput, error) {
				return NewUpdateVirtualCircuitTerminationHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_virtual_circuit_termination"],
		metrics,
		"delete_virtual_circuit_termination",
		WrapToolHandler[VirtualCircuitTerminationDeleteInput, struct{}](
			metrics,
			"delete_virtual_circuit_termination",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in VirtualCircuitTerminationDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteVirtualCircuitTerminationHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_roles"],
		metrics,
		"get_roles",
		WrapToolHandler[RolesInput, PaginatedOutput[domain.Role]](
			metrics,
			"get_roles",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in RolesInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.Role], error) {
				return NewGetRolesHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_tenants"],
		metrics,
		"get_tenants",
		WrapToolHandler[TenantsInput, PaginatedOutput[domain.Tenant]](
			metrics,
			"get_tenants",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in TenantsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.Tenant], error) {
				return NewGetTenantsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_tags"],
		metrics,
		"get_tags",
		WrapToolHandler[TagsInput, PaginatedOutput[domain.Tag]](
			metrics,
			"get_tags",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in TagsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.Tag], error) {
				return NewGetTagsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_tag"],
		metrics,
		"create_tag",
		WrapToolHandler[TagCreateInput, TagOutput](
			metrics,
			"create_tag",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in TagCreateInput,
			) (*mcp.CallToolResult, TagOutput, error) {
				return NewCreateTagHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_tag"],
		metrics,
		"update_tag",
		WrapToolHandler[TagUpdateInput, TagOutput](
			metrics,
			"update_tag",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in TagUpdateInput,
			) (*mcp.CallToolResult, TagOutput, error) {
				return NewUpdateTagHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_tag"],
		metrics,
		"delete_tag",
		WrapToolHandler[TagDeleteInput, struct{}](
			metrics,
			"delete_tag",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in TagDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteTagHandler(svc)(ctx, req, in)
			},
		),
	)

	registerTool(
		server,
		toolDefIndex()["get_custom_fields"],
		metrics,
		"get_custom_fields",
		WrapToolHandler[CustomFieldsInput, PaginatedOutput[domain.CustomField]](
			metrics,
			"get_custom_fields",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CustomFieldsInput,
			) (*mcp.CallToolResult, PaginatedOutput[domain.CustomField], error) {
				return NewGetCustomFieldsHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["create_custom_field"],
		metrics,
		"create_custom_field",
		WrapToolHandler[CustomFieldCreateInput, CustomFieldOutput](
			metrics,
			"create_custom_field",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CustomFieldCreateInput,
			) (*mcp.CallToolResult, CustomFieldOutput, error) {
				return NewCreateCustomFieldHandler(svc)(ctx, req, in)
			},
		),
	)

	registerWriteTool(
		server,
		toolDefIndex()["update_custom_field"],
		metrics,
		"update_custom_field",
		WrapToolHandler[CustomFieldUpdateInput, CustomFieldOutput](
			metrics,
			"update_custom_field",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CustomFieldUpdateInput,
			) (*mcp.CallToolResult, CustomFieldOutput, error) {
				return NewUpdateCustomFieldHandler(svc)(ctx, req, in)
			},
		),
	)

	registerDeleteTool(
		server,
		toolDefIndex()["delete_custom_field"],
		metrics,
		"delete_custom_field",
		WrapToolHandler[CustomFieldDeleteInput, struct{}](
			metrics,
			"delete_custom_field",
			func(
				ctx context.Context,
				req *mcp.CallToolRequest,
				in CustomFieldDeleteInput,
			) (*mcp.CallToolResult, struct{}, error) {
				return NewDeleteCustomFieldHandler(svc)(ctx, req, in)
			},
		),
	)
}
