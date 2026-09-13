package handlers

// toolDef describes one MCP tool registration, including the per-tool
// Instructions metadata (M4) that guides the model on how and when to use it.
//
// The current official Go SDK (go-sdk v1.7.0) does not yet expose a per-tool
// Instructions field on mcp.Tool, so these instructions are surfaced to the
// model through the Description channel (the SDK's only per-tool guidance
// mechanism) and are kept as first-class, testable data here.
type toolDef struct {
	Name         string
	Title        string
	Description  string
	Instructions string
}

// toolDefs returns the ordered set of tool definitions registered by
// RegisterTools. The order follows the read -> write/update -> delete grouping
// (S3); every tool here is read-only.
func toolDefs() []toolDef {
	return []toolDef{
		{
			Name:         "get_sites",
			Title:        "List Sites",
			Description:  "List sites in NetBox with optional filters. Use q for free-text search, or filter by region, status, tenant, or tag. Results are paginated: set page and page_size (max 1000). An empty results array means no sites match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox sites. All filters are additive. Prefer this read-only tool over get_object_by_id when you need multiple sites or don't yet know the exact ID.",
		},
		{
			Name:         "get_devices",
			Title:        "List Devices",
			Description:  "List devices in NetBox with optional filters. Filter by site, role, manufacturer, device_type, status, name, tenant, rack, cluster, or tag, or use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no devices match.",
			Instructions: "Use to enumerate or search NetBox devices. Combine filters to narrow results; the tool never modifies NetBox.",
		},
		{
			Name:         "get_ip_addresses",
			Title:        "List IP Addresses",
			Description:  "Search IP addresses in NetBox. Filter by address (e.g. 192.168.1.0/24), assigned device, status, VRF, role, tenant, or tag, or use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no addresses match.",
			Instructions: "Use to look up IP addresses by prefix, device, or free-text. The address filter accepts a specific IP or a CIDR prefix.",
		},
		{
			Name:         "get_prefixes",
			Title:        "List Prefixes",
			Description:  "Search IP prefixes in NetBox. Filter by prefix (e.g. 10.0.0.0/8), site, VRF, status, role, tenant, family, or tag, or use within to find prefixes nested inside a given prefix. Results are paginated: set page and page_size (max 1000). An empty results array means no prefixes match.",
			Instructions: "Use to explore IP address space. Use 'within' to find child prefixes nested in a parent prefix.",
		},
		{
			Name:         "get_vlans",
			Title:        "List VLANs",
			Description:  "List VLANs in NetBox. Filter by site, group, status, tenant, or tag, or by a specific VLAN ID (vid). Use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no VLANs match.",
			Instructions: "Use to look up VLANs by site, group, or VLAN ID.",
		},
		{
			Name:         "get_virtual_machines",
			Title:        "List Virtual Machines",
			Description:  "List virtual machines in NetBox. Filter by cluster, cluster_group, role, status, tenant, name, or site, or use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no VMs match.",
			Instructions: "Use to enumerate or search NetBox virtual machines.",
		},
		{
			Name:         "get_clusters",
			Title:        "List Clusters",
			Description:  "List clusters in NetBox. Filter by cluster_type, cluster_group, site, tenant, or name, or use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no clusters match.",
			Instructions: "Use to look up virtualization clusters.",
		},
		{
			Name:         "get_circuits",
			Title:        "List Circuits",
			Description:  "List circuits in NetBox. Filter by provider, circuit_type, site, status, or tenant, or use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no circuits match.",
			Instructions: "Use to enumerate or search NetBox circuits.",
		},
		{
			Name:         "get_object_by_id",
			Title:        "Get Object by ID",
			Description:  "Retrieve any single NetBox object by its type and numeric ID. object_type is required (e.g. site, device, prefix, ip_address, vlan, virtual_machine, cluster, circuit, provider, tenant, rack, cable, interface). Returns a 404-style error if the object does not exist or the token lacks permission.",
			Instructions: "Use when you already know the object type and numeric ID. For searching without an ID, prefer the dedicated list tools.",
		},
		{
			Name:         "get_racks",
			Title:        "List Racks",
			Description:  "List racks in NetBox. Filter by site, location, status, or tenant, or use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no racks match.",
			Instructions: "Use to enumerate or search NetBox racks.",
		},
		{
			Name:         "get_interfaces",
			Title:        "List Interfaces",
			Description:  "List device interfaces in NetBox. Filter by device, type, enabled, or name, or use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no interfaces match.",
			Instructions: "Use to look up device interfaces, optionally filtering by enabled state.",
		},
		{
			Name:         "get_circuit_terminations",
			Title:        "List Circuit Terminations",
			Description:  "List circuit terminations in NetBox. Filter by circuit, site, term_side (A or Z), or tag, or use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no terminations match.",
			Instructions: "Use to look up circuit terminations, optionally filtering by termination side (A or Z).",
		},
		{
			Name:         "get_cables",
			Title:        "List Cables",
			Description:  "List cables in NetBox. Filter by type, status, site, color, or label, or use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no cables match.",
			Instructions: "Use to enumerate or search NetBox cables.",
		},
		{
			Name:         "get_vm_interfaces",
			Title:        "List VM Interfaces",
			Description:  "List VM interfaces in NetBox. Filter by virtual_machine, name, or tag, or use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty results array means no VM interfaces match.",
			Instructions: "Use to look up virtual machine interfaces.",
		},
	}
}
