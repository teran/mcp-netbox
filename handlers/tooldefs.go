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
// (S3): the read tools appear first, followed by create/update and delete tools.
func toolDefs() []toolDef { //nolint:funlen,maintidx,lll // Data-declaration function; size/cyclomatic count is inherent to listing every tool.
	return []toolDef{
		{
			Name:  "get_sites",
			Title: "List Sites",
			Description: "List sites in NetBox with optional filters. Use q for free-text search, or filter by " +
				"region, status, tenant, or tag. Results are paginated: set page and page_size (max 1000). " +
				"An empty results array means no sites match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox sites. All filters are additive. Prefer this " +
				"read-only tool over get_object_by_id when you need multiple sites or don't yet know the exact ID.",
		},
		{
			Name:  "get_regions",
			Title: "List Regions",
			Description: "List regions in NetBox with optional filters. Use q for free-text search, or filter by " +
				"parent (region ID), or tag. Results are paginated: set page and page_size (max 1000). " +
				"An empty results array means no regions match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox regions. All filters are additive. Prefer this " +
				"read-only tool over get_object_by_id when you need multiple regions or don't yet know the exact ID.",
		},
		{
			Name:  "get_tenant_groups",
			Title: "List Tenant Groups",
			Description: "List tenant groups in NetBox with optional filters. Use q for free-text search, or filter " +
				"by parent (tenant group ID), or tag. Results are paginated: set page and page_size (max 1000). " +
				"An empty results array means no tenant groups match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox tenant groups. All filters are additive. Prefer this " +
				"read-only tool over get_object_by_id when you need multiple tenant groups or don't yet know the exact ID.",
		},
		{
			Name:  "get_devices",
			Title: "List Devices",
			Description: "List devices in NetBox with optional filters. Filter by site, role, manufacturer, " +
				"device_type, status, name, tenant, rack, cluster, or tag, or use q for free-text search. " +
				"Results are paginated: set page and page_size (max 1000). An empty results array means no devices match.",
			Instructions: "Use to enumerate or search NetBox devices. Combine filters to narrow results; " +
				"the tool never modifies NetBox.",
		},
		{
			Name:  "get_ip_addresses",
			Title: "List IP Addresses",
			Description: "Search IP addresses in NetBox. Filter by address (e.g. 192.168.1.0/24), assigned " +
				"device, status, VRF, role, tenant, or tag, or use q for free-text search. Results are " +
				"paginated: set page and page_size (max 1000). An empty results array means no addresses match.",
			Instructions: "Use to look up IP addresses by prefix, device, or free-text. The address filter " +
				"accepts a specific IP or a CIDR prefix.",
		},
		{
			Name:  "get_prefixes",
			Title: "List Prefixes",
			Description: "Search IP prefixes in NetBox. Filter by prefix (e.g. 10.0.0.0/8), site, VRF, " +
				"status, role, tenant, family, or tag, or use within to find prefixes nested inside a given " +
				"prefix. Results are paginated: set page and page_size (max 1000). An empty results array means no prefixes match.",
			Instructions: "Use to explore IP address space. Use 'within' to find child prefixes nested in a parent prefix.",
		},
		{
			Name:  "get_vlans",
			Title: "List VLANs",
			Description: "List VLANs in NetBox. Filter by site, group, status, tenant, or tag, or by a " +
				"specific VLAN ID (vid). Use q for free-text search. Results are paginated: set page and " +
				"page_size (max 1000). An empty results array means no VLANs match.",
			Instructions: "Use to look up VLANs by site, group, or VLAN ID.",
		},
		{
			Name:  "get_virtual_machines",
			Title: "List Virtual Machines",
			Description: "List virtual machines in NetBox. Filter by cluster, cluster_group, role, status, " +
				"tenant, name, or site, or use q for free-text search. Results are paginated: set page and " +
				"page_size (max 1000). An empty results array means no VMs match.",
			Instructions: "Use to enumerate or search NetBox virtual machines.",
		},
		{
			Name:  "get_clusters",
			Title: "List Clusters",
			Description: "List clusters in NetBox. Filter by cluster_type, cluster_group, site, tenant, or " +
				"name, or use q for free-text search. Results are paginated: set page and page_size (max " +
				"1000). An empty results array means no clusters match.",
			Instructions: "Use to look up virtualization clusters.",
		},
		{
			Name:  "get_circuits",
			Title: "List Circuits",
			Description: "List circuits in NetBox. Filter by provider, circuit_type, site, status, or tenant, " +
				"or use q for free-text search. Results are paginated: set page and page_size (max 1000). " +
				"An empty results array means no circuits match.",
			Instructions: "Use to enumerate or search NetBox circuits.",
		},
		{
			Name:  "get_object_by_id",
			Title: "Get Object by ID",
			Description: "Retrieve any single NetBox object by its type and numeric ID. object_type is " +
				"required (e.g. site, region, device, prefix, ip_address, vlan, virtual_machine, cluster, circuit, " +
				"provider, tenant, rack, cable, interface, site_group, rack_group, rack_role, device_role, platform, " +
				"virtual_chassis, power_panel, inventory_item_role, front_port, rear_port, console_port, " +
				"console_server_port, power_port, power_outlet, power_feed, device_bay, rack_reservation, " +
				"inventory_item, mac_address, rir, aggregate, ip_range, asn, asn_range, route_target, " +
				"service, service_template, fhrp_group, fhrp_group_assignment, vlan_translation_policy, " +
				"vlan_translation_rule, contact_group, contact_role, virtual_disk, virtual_machine_type, " +
				"provider_account, provider_network, circuit_group, circuit_group_assignment, virtual_circuit, " +
				"virtual_circuit_type, virtual_circuit_termination). " +
				"Returns a 404-style error if the object does " +
				"not exist or the token lacks permission.",
			Instructions: "Use when you already know the object type and numeric ID. For searching without " +
				"an ID, prefer the dedicated list tools.",
		},
		{
			Name:  "get_racks",
			Title: "List Racks",
			Description: "List racks in NetBox. Filter by site, location, status, or tenant, or use q for " +
				"free-text search. Results are paginated: set page and page_size (max 1000). An empty " +
				"results array means no racks match.",
			Instructions: "Use to enumerate or search NetBox racks.",
		},
		{
			Name:  "get_interfaces",
			Title: "List Interfaces",
			Description: "List device interfaces in NetBox. Filter by device, type, enabled, or name, or " +
				"use q for free-text search. Results are paginated: set page and page_size (max 1000). " +
				"An empty results array means no interfaces match.",
			Instructions: "Use to look up device interfaces, optionally filtering by enabled state.",
		},
		{
			Name:  "get_circuit_terminations",
			Title: "List Circuit Terminations",
			Description: "List circuit terminations in NetBox. Filter by circuit, site, term_side (A or Z), " +
				"or tag, or use q for free-text search. Results are paginated: set page and page_size " +
				"(max 1000). An empty results array means no terminations match.",
			Instructions: "Use to look up circuit terminations, optionally filtering by termination side (A or Z).",
		},
		{
			Name:  "get_cables",
			Title: "List Cables",
			Description: "List cables in NetBox. Filter by type, status, site, color, or label, or use q for " +
				"free-text search. Results are paginated: set page and page_size (max 1000). An empty " +
				"results array means no cables match.",
			Instructions: "Use to enumerate or search NetBox cables.",
		},
		{
			Name:  "get_vm_interfaces",
			Title: "List VM Interfaces",
			Description: "List VM interfaces in NetBox. Filter by virtual_machine, name, or tag, or use q " +
				"for free-text search. Results are paginated: set page and page_size (max 1000). An empty " +
				"results array means no VM interfaces match.",
			Instructions: "Use to look up virtual machine interfaces.",
		},
		{
			Name:  "get_site_groups",
			Title: "List Site Groups",
			Description: "List site groups in NetBox with optional filters. Use q for free-text search, or " +
				"filter by parent (site group ID) or tag. Results are paginated: set page and page_size " +
				"(max 1000). An empty results array means no site groups match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox site groups. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple site groups or don't yet know the exact ID.",
		},
		{
			Name:  "get_rack_groups",
			Title: "List Rack Groups",
			Description: "List rack groups in NetBox with optional filters. Use q for free-text search, or " +
				"filter by site, parent (rack group ID) or tag. Results are paginated: set page and " +
				"page_size (max 1000). An empty results array means no rack groups match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox rack groups. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple rack groups or don't yet know the exact ID.",
		},
		{
			Name:  "get_rack_roles",
			Title: "List Rack Roles",
			Description: "List rack roles in NetBox with optional filters. Use q for free-text search, or " +
				"filter by color or tag. Results are paginated: set page and page_size (max 1000). An empty " +
				"results array means no rack roles match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox rack roles. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple rack roles or don't yet know the exact ID.",
		},
		{
			Name:  "get_device_roles",
			Title: "List Device Roles",
			Description: "List device roles in NetBox with optional filters. Use q for free-text search, or " +
				"filter by color or tag. Results are paginated: set page and page_size (max 1000). An empty " +
				"results array means no device roles match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox device roles. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple device roles or don't yet know the exact ID.",
		},
		{
			Name:  "get_platforms",
			Title: "List Platforms",
			Description: "List platforms in NetBox with optional filters. Use q for free-text search, or " +
				"filter by manufacturer or tag. Results are paginated: set page and page_size (max 1000). " +
				"An empty results array means no platforms match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox platforms. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple platforms or don't yet know the exact ID.",
		},
		{
			Name:  "get_virtual_chassis",
			Title: "List Virtual Chassis",
			Description: "List virtual chassis in NetBox with optional filters. Use q for free-text search, " +
				"or filter by name or tag. Results are paginated: set page and page_size (max 1000). An " +
				"empty results array means no virtual chassis match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox virtual chassis. Prefer this read-only tool " +
				"over get_object_by_id when you need multiple virtual chassis or don't yet know the exact ID.",
		},
		{
			Name:  "get_power_panels",
			Title: "List Power Panels",
			Description: "List power panels in NetBox with optional filters. Use q for free-text search, or " +
				"filter by site, location or tag. Results are paginated: set page and page_size (max 1000). " +
				"An empty results array means no power panels match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox power panels. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple power panels or don't yet know the exact ID.",
		},
		{
			Name:  "get_inventory_item_roles",
			Title: "List Inventory Item Roles",
			Description: "List inventory item roles in NetBox with optional filters. Use q for free-text " +
				"search, or filter by color or tag. Results are paginated: set page and page_size (max " +
				"1000). An empty results array means no inventory item roles match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox inventory item roles. Prefer this read-only " +
				"tool over get_object_by_id when you need multiple inventory item roles or don't yet know " +
				"the exact ID.",
		},
		{
			Name:  "create_site",
			Title: "Create Site",
			Description: "Create a new site in NetBox. name is required; all other fields are optional. " +
				"Provide region/tenant as numeric IDs (look them up first), status as one of " +
				"planned/staged/active/decommissioning/retired, and tags as a list of tag names. Returns the created site.",
			Instructions: "Use only when the caller explicitly wants to create a site and provides a name. " +
				"This tool mutates NetBox; confirm the name and key attributes before calling.",
		},
		{
			Name:  "update_site",
			Title: "Update Site",
			Description: "Partially update a site in NetBox by its numeric ID. Only the fields you " +
				"explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated site.",
			Instructions: "Use to modify one or more fields of an existing site. Provide id plus only the " +
				"fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_site",
			Title: "Delete Site",
			Description: "WARNING: This operation is irreversible. Permanently deletes a site (and its " +
				"dependent objects) from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the site, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_region",
			Title: "Create Region",
			Description: "Create a new region in NetBox. name is required; all other fields are optional. " +
				"Provide parent as a numeric region ID (look it up first), and tags as a list of tag names. " +
				"Returns the created region.",
			Instructions: "Use only when the caller explicitly wants to create a region and provides a name. " +
				"This tool mutates NetBox; confirm the name and key attributes before calling.",
		},
		{
			Name:  "update_region",
			Title: "Update Region",
			Description: "Partially update a region in NetBox by its numeric ID. Only the fields you " +
				"explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated region.",
			Instructions: "Use to modify one or more fields of an existing region. Provide id plus only the " +
				"fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_region",
			Title: "Delete Region",
			Description: "WARNING: This operation is irreversible. Permanently deletes a region (and its " +
				"dependent objects) from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the region, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_tenant_group",
			Title: "Create Tenant Group",
			Description: "Create a new tenancy tenant group in NetBox. name is required; all other fields are " +
				"optional. Provide parent as a numeric tenant group ID (look it up first), and tags as a list of " +
				"tag names. Returns the created tenant group.",
			Instructions: "Use only when the caller explicitly wants to create a tenant group and provides a name. " +
				"This tool mutates NetBox; confirm the name and key attributes before calling.",
		},
		{
			Name:  "update_tenant_group",
			Title: "Update Tenant Group",
			Description: "Partially update a tenancy tenant group in NetBox by its numeric ID. Only the fields " +
				"you explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated tenant group.",
			Instructions: "Use to modify one or more fields of an existing tenant group. Provide id plus only the " +
				"fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_tenant_group",
			Title: "Delete Tenant Group",
			Description: "WARNING: This operation is irreversible. Permanently deletes a tenancy tenant group " +
				"from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the tenant group, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_device",
			Title: "Create Device",
			Description: "Create a new device in NetBox. name is required; all other fields are optional. " +
				"Provide device_type/role/site/rack/cluster/tenant/platform as numeric IDs (look them up " +
				"first), status as one of offline/active/planned/staged/failed/inventory/decommissioning, " +
				"face as front/rear, and tags as a list of tag names. Returns the created device.",
			Instructions: "Use only when the caller explicitly wants to create a device and provides a " +
				"name. This tool mutates NetBox; confirm the name and key attributes (device_type, role, site) before calling.",
		},
		{
			Name:  "update_device",
			Title: "Update Device",
			Description: "Partially update a device in NetBox by its numeric ID. Only the fields you " +
				"explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated device.",
			Instructions: "Use to modify one or more fields of an existing device. Provide id plus only " +
				"the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_device",
			Title: "Delete Device",
			Description: "WARNING: This operation is irreversible. Permanently deletes a device from NetBox " +
				"by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the device, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_ip_address",
			Title: "Create IP Address",
			Description: "Create a new IP address in NetBox. address is required (e.g. 192.168.1.1/24); " +
				"all other fields are optional. Provide vrf/tenant/assigned_object_id as numeric IDs " +
				"(look them up first), status as one of active/reserved/deprecated/dhcp/slaac, role as one " +
				"of loopback/secondary/anycast/vip/vrrp/hsrp/glbp/carp, and tags as a list " +
				"of tag names. Returns the created IP address.",
			Instructions: "Use only when the caller explicitly wants to create an IP address and provides " +
				"an address. This tool mutates NetBox; confirm the address and key attributes before calling.",
		},
		{
			Name:  "update_ip_address",
			Title: "Update IP Address",
			Description: "Partially update an IP address in NetBox by its numeric ID. Only the fields you " +
				"explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display, assigned_object) are ignored. " +
				"Returns the updated IP address.",
			Instructions: "Use to modify one or more fields of an existing IP address. Provide id plus " +
				"only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_ip_address",
			Title: "Delete IP Address",
			Description: "WARNING: This operation is irreversible. Permanently deletes an IP address from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the IP address, ideally by ID, and understands it is permanent and cannot " +
				"be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_prefix",
			Title: "Create Prefix",
			Description: "Create a new prefix in NetBox. prefix is required (CIDR, e.g. 10.0.0.0/24); all " +
				"other fields are optional. Provide site/vrf/tenant/vlan/role as numeric IDs (look them up " +
				"first), status as one of container/active/reserved/deprecated, is_pool as a boolean, and " +
				"tags as a list of tag names. Returns the created prefix.",
			Instructions: "Use only when the caller explicitly wants to create a prefix and provides a " +
				"CIDR. This tool mutates NetBox; confirm the prefix and key attributes before calling.",
		},
		{
			Name:  "update_prefix",
			Title: "Update Prefix",
			Description: "Partially update a prefix in NetBox by its numeric ID. Only the fields you " +
				"explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display, children, _depth, family) are " +
				"ignored. Returns the updated prefix.",
			Instructions: "Use to modify one or more fields of an existing prefix. Provide id plus only " +
				"the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_prefix",
			Title: "Delete Prefix",
			Description: "WARNING: This operation is irreversible. Permanently deletes a prefix from NetBox " +
				"by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the prefix, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_vlan",
			Title: "Create VLAN",
			Description: "Create a new VLAN in NetBox. vid (1-4094) and name are required; all other fields " +
				"are optional. Provide site/group/tenant/role as numeric IDs (look them up first), status " +
				"as one of active/reserved/deprecated, and tags as a list of tag names. Returns the created VLAN.",
			Instructions: "Use only when the caller explicitly wants to create a VLAN and provides both a " +
				"VLAN ID (vid) and a name. This tool mutates NetBox; confirm the vid and key attributes before calling.",
		},
		{
			Name:  "update_vlan",
			Title: "Update VLAN",
			Description: "Partially update a VLAN in NetBox by its numeric ID. Only the fields you " +
				"explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated VLAN.",
			Instructions: "Use to modify one or more fields of an existing VLAN. Provide id plus only the " +
				"fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_vlan",
			Title: "Delete VLAN",
			Description: "WARNING: This operation is irreversible. Permanently deletes a VLAN from NetBox " +
				"by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the VLAN, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_virtual_machine",
			Title: "Create Virtual Machine",
			Description: "Create a new virtual machine in NetBox. name is required; all other fields are " +
				"optional. Provide cluster/role/tenant/platform/site as numeric IDs (look them up first), " +
				"status as one of offline/active/planned/staged/failed/decommissioning, vcpus as a number, " +
				"memory (MB) and disk (GB) as integers, and tags as a list of tag names. Returns the created virtual machine.",
			Instructions: "Use only when the caller explicitly wants to create a virtual machine and " +
				"provides a name. This tool mutates NetBox; confirm the name and key attributes before calling.",
		},
		{
			Name:  "update_virtual_machine",
			Title: "Update Virtual Machine",
			Description: "Partially update a virtual machine in NetBox by its numeric ID. Only the fields " +
				"you explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display) are ignored. Returns the " +
				"updated virtual machine.",
			Instructions: "Use to modify one or more fields of an existing virtual machine. Provide id " +
				"plus only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_virtual_machine",
			Title: "Delete Virtual Machine",
			Description: "WARNING: This operation is irreversible. Permanently deletes a virtual machine " +
				"from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the virtual machine, ideally by ID, and understands it is permanent and " +
				"cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming " +
				"deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "create_cluster",
			Title: "Create Cluster",
			Description: "Create a new cluster in NetBox. name and type (cluster type) are required; all " +
				"other fields are optional. Provide type/group/site/tenant as numeric IDs (look them up " +
				"first), and tags as a list of tag names. Returns the created cluster.",
			Instructions: "Use only when the caller explicitly wants to create a cluster and provides a " +
				"name and a cluster type. This tool mutates NetBox; confirm the name and key attributes before calling.",
		},
		{
			Name:  "update_cluster",
			Title: "Update Cluster",
			Description: "Partially update a cluster in NetBox by its numeric ID. Only the fields you " +
				"explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated cluster.",
			Instructions: "Use to modify one or more fields of an existing cluster. Provide id plus only " +
				"the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_cluster",
			Title: "Delete Cluster",
			Description: "WARNING: This operation is irreversible. Permanently deletes a cluster from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the cluster, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_circuit",
			Title: "Create Circuit",
			Description: "Create a new circuit in NetBox. cid and provider and circuit_type are required; " +
				"all other fields are optional. Provide provider/circuit_type/tenant as numeric IDs (look " +
				"them up first), status as one of planned/provisioning/active/offline/decommissioning, " +
				"install_date as YYYY-MM-DD, commit_rate in kbps, and tags as a list of tag names. Returns the created circuit.",
			Instructions: "Use only when the caller explicitly wants to create a circuit and provides a " +
				"cid, a provider and a circuit_type. This tool mutates NetBox; confirm the cid and key " +
				"attributes before calling.",
		},
		{
			Name:  "update_circuit",
			Title: "Update Circuit",
			Description: "Partially update a circuit in NetBox by its numeric ID. Only the fields you " +
				"explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated circuit.",
			Instructions: "Use to modify one or more fields of an existing circuit. Provide id plus only " +
				"the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_circuit",
			Title: "Delete Circuit",
			Description: "WARNING: This operation is irreversible. Permanently deletes a circuit from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the circuit, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_rack",
			Title: "Create Rack",
			Description: "Create a new rack in NetBox. name is required; all other fields are optional. " +
				"Provide site/location/tenant/role as numeric IDs (look them up first), status as one of " +
				"reserved/available/planned/active/decommissioning, type as one of " +
				"2-post-frame/4-post-frame/4-post-cabinet/wall-frame/wall-cabinet, width in inches, " +
				"u_height in rack units, and tags as a list of tag names. Returns the created rack.",
			Instructions: "Use only when the caller explicitly wants to create a rack and provides a name. " +
				"This tool mutates NetBox; confirm the name and key attributes before calling.",
		},
		{
			Name:  "update_rack",
			Title: "Update Rack",
			Description: "Partially update a rack in NetBox by its numeric ID. Only the fields you " +
				"explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated rack.",
			Instructions: "Use to modify one or more fields of an existing rack. Provide id plus only the " +
				"fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_rack",
			Title: "Delete Rack",
			Description: "WARNING: This operation is irreversible. Permanently deletes a rack from NetBox " +
				"by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the rack, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_interface",
			Title: "Create Interface",
			Description: "Create a new device interface in NetBox. name, device and type are required; all " +
				"other fields are optional. Provide device as a numeric ID (look it up first), type as one " +
				"of the NetBox interface type values (e.g. 1000base-t, 10gbase-x-sfpp), and tags as a list " +
				"of tag names. Returns the created interface.",
			Instructions: "Use only when the caller explicitly wants to create an interface and provides " +
				"a name, a device and a type. This tool mutates NetBox; confirm the name and key " +
				"attributes before calling.",
		},
		{
			Name:  "update_interface",
			Title: "Update Interface",
			Description: "Partially update a device interface in NetBox by its numeric ID. Only the fields " +
				"you explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated interface.",
			Instructions: "Use to modify one or more fields of an existing interface. Provide id plus only " +
				"the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_interface",
			Title: "Delete Interface",
			Description: "WARNING: This operation is irreversible. Permanently deletes a device interface " +
				"from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the interface, ideally by ID, and understands it is permanent and cannot " +
				"be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_circuit_termination",
			Title: "Create Circuit Termination",
			Description: "Create a new circuit termination in NetBox. circuit, term_side (A or Z) and site " +
				"are required; all other fields are optional. Provide circuit and site as numeric IDs " +
				"(look them up first), speed and upstream_speed in kbps, and tags as a list of tag names. " +
				"Returns the created circuit termination.",
			Instructions: "Use only when the caller explicitly wants to create a circuit termination and " +
				"provides a circuit, a term_side and a site. This tool mutates NetBox; confirm the key " +
				"attributes before calling.",
		},
		{
			Name:  "update_circuit_termination",
			Title: "Update Circuit Termination",
			Description: "Partially update a circuit termination in NetBox by its numeric ID. Only the " +
				"fields you explicitly provide are changed (PATCH merge); omitted fields are left " +
				"untouched. Read-only fields (id, url, created, last_updated, display) are ignored. " +
				"Returns the updated circuit termination.",
			Instructions: "Use to modify one or more fields of an existing circuit termination. Provide id " +
				"plus only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_circuit_termination",
			Title: "Delete Circuit Termination",
			Description: "WARNING: This operation is irreversible. Permanently deletes a circuit " +
				"termination from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the circuit termination, ideally by ID, and understands it is permanent " +
				"and cannot be undone. If the user merely asks to 'remove' or 'clean up' without " +
				"confirming deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "create_cable",
			Title: "Create Cable",
			Description: "Create a new cable in NetBox. termination_a and termination_b are required; " +
				"each is an object with object_type (e.g. dcim.interface) and object_id. All other fields " +
				"are optional. Provide type/status/length_unit as NetBox choice strings, length as a " +
				"number, and tags as a list of tag names. Returns the created cable.",
			Instructions: "Use only when the caller explicitly wants to create a cable and provides both " +
				"terminations. This tool mutates NetBox; confirm the key attributes before calling.",
		},
		{
			Name:  "update_cable",
			Title: "Update Cable",
			Description: "Partially update a cable in NetBox by its numeric ID. Only the fields you " +
				"explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated cable.",
			Instructions: "Use to modify one or more fields of an existing cable. Provide id plus only " +
				"the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_cable",
			Title: "Delete Cable",
			Description: "WARNING: This operation is irreversible. Permanently deletes a cable from NetBox " +
				"by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the cable, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_vm_interface",
			Title: "Create VM Interface",
			Description: "Create a new virtual-machine interface in NetBox. virtual_machine (the parent " +
				"VM's numeric ID) and name are required; all other fields are optional. Provide mtu as an " +
				"integer, mac_address as a string, and tags as a list of tag names. Returns the created VM interface.",
			Instructions: "Use only when the caller explicitly wants to create a VM interface and provides " +
				"a virtual_machine and a name. This tool mutates NetBox; confirm the key attributes before calling.",
		},
		{
			Name:  "update_vm_interface",
			Title: "Update VM Interface",
			Description: "Partially update a virtual-machine interface in NetBox by its numeric ID. Only " +
				"the fields you explicitly provide are changed (PATCH merge); omitted fields are left " +
				"untouched. Read-only fields (id, url, created, last_updated, display) are ignored. " +
				"Returns the updated VM interface.",
			Instructions: "Use to modify one or more fields of an existing VM interface. Provide id plus " +
				"only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_vm_interface",
			Title: "Delete VM Interface",
			Description: "WARNING: This operation is irreversible. Permanently deletes a virtual-machine " +
				"interface from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the VM interface, ideally by ID, and understands it is permanent and " +
				"cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming " +
				"deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "create_provider",
			Title: "Create Provider",
			Description: "Create a new circuits provider in NetBox. name is required; all other fields are " +
				"optional. Provide asn as an integer, slug/account/portal_url/noc_contact/" +
				"admin_contact/comments as strings, and tags as a list of tag names. Returns the created provider.",
			Instructions: "Use only when the caller explicitly wants to create a provider and provides a " +
				"name. This tool mutates NetBox; confirm the key attributes before calling.",
		},
		{
			Name:  "update_provider",
			Title: "Update Provider",
			Description: "Partially update a circuits provider in NetBox by its numeric ID. Only the " +
				"fields you explicitly provide are changed (PATCH merge); omitted fields are left " +
				"untouched. Read-only fields (id, url, created, last_updated, display) are ignored. " +
				"Returns the updated provider.",
			Instructions: "Use to modify one or more fields of an existing provider. Provide id plus only " +
				"the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_provider",
			Title: "Delete Provider",
			Description: "WARNING: This operation is irreversible. Permanently deletes a circuits provider " +
				"from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the provider, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_tenant",
			Title: "Create Tenant",
			Description: "Create a new tenancy tenant in NetBox. name is required; all other fields are " +
				"optional. Provide slug/description/comments as strings, group as a numeric tenant group ID, " +
				"and tags as a list of tag names. " +
				"Returns the created tenant.",
			Instructions: "Use only when the caller explicitly wants to create a tenant and provides a " +
				"name. This tool mutates NetBox; confirm the key attributes before calling.",
		},
		{
			Name:  "update_tenant",
			Title: "Update Tenant",
			Description: "Partially update a tenancy tenant in NetBox by its numeric ID. Only the fields " +
				"you explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated tenant.",
			Instructions: "Use to modify one or more fields of an existing tenant. Provide id plus only " +
				"the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_tenant",
			Title: "Delete Tenant",
			Description: "WARNING: This operation is irreversible. Permanently deletes a tenancy tenant " +
				"from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the tenant, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_manufacturer",
			Title: "Create Manufacturer",
			Description: "Create a new device manufacturer in NetBox. name is required; all other fields " +
				"are optional. Provide slug/description as strings, and tags as a list of tag names. " +
				"Returns the created manufacturer.",
			Instructions: "Use only when the caller explicitly wants to create a manufacturer and provides " +
				"a name. This tool mutates NetBox; confirm the key attributes before calling.",
		},
		{
			Name:  "update_manufacturer",
			Title: "Update Manufacturer",
			Description: "Partially update a device manufacturer in NetBox by its numeric ID. Only the " +
				"fields you explicitly provide are changed (PATCH merge); omitted fields are left " +
				"untouched. Read-only fields (id, url, created, last_updated, display) are ignored. " +
				"Returns the updated manufacturer.",
			Instructions: "Use to modify one or more fields of an existing manufacturer. Provide id plus " +
				"only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_manufacturer",
			Title: "Delete Manufacturer",
			Description: "WARNING: This operation is irreversible. Permanently deletes a device " +
				"manufacturer from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the manufacturer, ideally by ID, and understands it is permanent and " +
				"cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming " +
				"deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "create_device_type",
			Title: "Create Device Type",
			Description: "Create a new device type in NetBox. manufacturer (numeric ID) and model are " +
				"required; all other fields are optional. Provide part_number/comments as strings, u_height " +
				"as a number, is_full_depth/subdevice_role as booleans, slug as a string, and tags as a " +
				"list of tag names. Returns the created device type.",
			Instructions: "Use only when the caller explicitly wants to create a device type and provides " +
				"a manufacturer and model. This tool mutates NetBox; confirm the key attributes before calling.",
		},
		{
			Name:  "update_device_type",
			Title: "Update Device Type",
			Description: "Partially update a device type in NetBox by its numeric ID. Only the fields you " +
				"explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display) are ignored. Returns the " +
				"updated device type.",
			Instructions: "Use to modify one or more fields of an existing device type. Provide id plus " +
				"only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_device_type",
			Title: "Delete Device Type",
			Description: "WARNING: This operation is irreversible. Permanently deletes a device type from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the device type, ideally by ID, and understands it is permanent and cannot " +
				"be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_location",
			Title: "Create Location",
			Description: "Create a new site location in NetBox. name and site (numeric ID) are required; " +
				"all other fields are optional. Provide slug/parent (numeric ID)/description/status as " +
				"strings or integers, and tags as a list of tag names. Returns the created location.",
			Instructions: "Use only when the caller explicitly wants to create a location and provides a " +
				"name and a site. This tool mutates NetBox; confirm the key attributes before calling.",
		},
		{
			Name:  "update_location",
			Title: "Update Location",
			Description: "Partially update a site location in NetBox by its numeric ID. Only the fields " +
				"you explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated location.",
			Instructions: "Use to modify one or more fields of an existing location. Provide id plus only " +
				"the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_location",
			Title: "Delete Location",
			Description: "WARNING: This operation is irreversible. Permanently deletes a site location " +
				"from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the location, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_cluster_type",
			Title: "Create Cluster Type",
			Description: "Create a new cluster type in NetBox. name is required; all other fields are " +
				"optional. Provide slug/description as strings, and tags as a list of tag names. Returns " +
				"the created cluster type.",
			Instructions: "Use only when the caller explicitly wants to create a cluster type and provides " +
				"a name. This tool mutates NetBox; confirm the key attributes before calling.",
		},
		{
			Name:  "update_cluster_type",
			Title: "Update Cluster Type",
			Description: "Partially update a cluster type in NetBox by its numeric ID. Only the fields " +
				"you explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display) are ignored. Returns the " +
				"updated cluster type.",
			Instructions: "Use to modify one or more fields of an existing cluster type. Provide id plus " +
				"only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_cluster_type",
			Title: "Delete Cluster Type",
			Description: "WARNING: This operation is irreversible. Permanently deletes a cluster type from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the cluster type, ideally by ID, and understands it is permanent and " +
				"cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming " +
				"deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "create_cluster_group",
			Title: "Create Cluster Group",
			Description: "Create a new cluster group in NetBox. name is required; all other fields are " +
				"optional. Provide slug/description as strings, and tags as a list of tag names. Returns " +
				"the created cluster group.",
			Instructions: "Use only when the caller explicitly wants to create a cluster group and provides " +
				"a name. This tool mutates NetBox; confirm the key attributes before calling.",
		},
		{
			Name:  "update_cluster_group",
			Title: "Update Cluster Group",
			Description: "Partially update a cluster group in NetBox by its numeric ID. Only the fields " +
				"you explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display) are ignored. Returns the " +
				"updated cluster group.",
			Instructions: "Use to modify one or more fields of an existing cluster group. Provide id plus " +
				"only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_cluster_group",
			Title: "Delete Cluster Group",
			Description: "WARNING: This operation is irreversible. Permanently deletes a cluster group " +
				"from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the cluster group, ideally by ID, and understands it is permanent and " +
				"cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming " +
				"deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "create_circuit_type",
			Title: "Create Circuit Type",
			Description: "Create a new circuit type in NetBox. name is required; all other fields are " +
				"optional. Provide slug/description as strings, and tags as a list of tag names. Returns " +
				"the created circuit type.",
			Instructions: "Use only when the caller explicitly wants to create a circuit type and provides " +
				"a name. This tool mutates NetBox; confirm the key attributes before calling.",
		},
		{
			Name:  "update_circuit_type",
			Title: "Update Circuit Type",
			Description: "Partially update a circuit type in NetBox by its numeric ID. Only the fields " +
				"you explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, circuit_count, created, last_updated, display) are ignored. " +
				"Returns the updated circuit type.",
			Instructions: "Use to modify one or more fields of an existing circuit type. Provide id plus " +
				"only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_circuit_type",
			Title: "Delete Circuit Type",
			Description: "WARNING: This operation is irreversible. Permanently deletes a circuit type from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the circuit type, ideally by ID, and understands it is permanent and " +
				"cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming " +
				"deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "create_vrf",
			Title: "Create VRF",
			Description: "Create a new VRF in NetBox. name and rd are required; all other fields are " +
				"optional. Provide tenant as a numeric tenant ID, description as a string, and tags as a " +
				"list of tag names. Returns the created VRF.",
			Instructions: "Use only when the caller explicitly wants to create a VRF and provides a name " +
				"and route distinguisher (rd). This tool mutates NetBox; confirm the key attributes before calling.",
		},
		{
			Name:  "update_vrf",
			Title: "Update VRF",
			Description: "Partially update a VRF in NetBox by its numeric ID. Only the fields you " +
				"explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated VRF.",
			Instructions: "Use to modify one or more fields of an existing VRF. Provide id plus only the " +
				"fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_vrf",
			Title: "Delete VRF",
			Description: "WARNING: This operation is irreversible. Permanently deletes a VRF from NetBox " +
				"by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the VRF, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_vlan_group",
			Title: "Create VLAN Group",
			Description: "Create a new VLAN group in NetBox. name is required; all other fields are " +
				"optional. Provide slug/description as strings, scope_type (e.g. dcim.site) and scope_id " +
				"for the scope, min_vid/max_vid as integers, and tags as a list of tag names. Returns the created VLAN group.",
			Instructions: "Use only when the caller explicitly wants to create a VLAN group and provides a " +
				"name. This tool mutates NetBox; confirm the key attributes before calling.",
		},
		{
			Name:  "update_vlan_group",
			Title: "Update VLAN Group",
			Description: "Partially update a VLAN group in NetBox by its numeric ID. Only the fields you " +
				"explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display) are ignored. Returns the " +
				"updated VLAN group.",
			Instructions: "Use to modify one or more fields of an existing VLAN group. Provide id plus " +
				"only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_vlan_group",
			Title: "Delete VLAN Group",
			Description: "WARNING: This operation is irreversible. Permanently deletes a VLAN group from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the VLAN group, ideally by ID, and understands it is permanent and cannot " +
				"be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_role",
			Title: "Create Role",
			Description: "Create a new IPAM role in NetBox. name is required; all other fields are " +
				"optional. Provide slug and description as strings, weight as an integer, and tags as a " +
				"list of tag names. Returns the created role.",
			Instructions: "Use only when the caller explicitly wants to create an IPAM role and provides " +
				"a name. This tool mutates NetBox; confirm the key attributes before calling.",
		},
		{
			Name:  "update_role",
			Title: "Update Role",
			Description: "Partially update an IPAM role in NetBox by its numeric ID. Only the fields you " +
				"explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated) are ignored. Returns the updated role.",
			Instructions: "Use to modify one or more fields of an existing IPAM role. Provide id plus only " +
				"the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_role",
			Title: "Delete Role",
			Description: "WARNING: This operation is irreversible. Permanently deletes an IPAM role from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the role, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_contact",
			Title: "Create Contact",
			Description: "Create a new tenancy contact in NetBox. name is required; all other fields are " +
				"optional. Provide group as a numeric contact group ID, and " +
				"title/phone/email/address/link/description/comments as strings. Returns the created contact.",
			Instructions: "Use only when the caller explicitly wants to create a contact and provides a " +
				"name. This tool mutates NetBox; confirm the key attributes before calling.",
		},
		{
			Name:  "update_contact",
			Title: "Update Contact",
			Description: "Partially update a tenancy contact in NetBox by its numeric ID. Only the fields " +
				"you explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated) are ignored. Returns the updated contact.",
			Instructions: "Use to modify one or more fields of an existing contact. Provide id plus only " +
				"the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_contact",
			Title: "Delete Contact",
			Description: "WARNING: This operation is irreversible. Permanently deletes a tenancy contact " +
				"from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the contact, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_site_group",
			Title: "Create Site Group",
			Description: "Create a new site group in NetBox. name is required; all other fields are " +
				"optional. Provide parent as a numeric site group ID (look it up first), and tags as a " +
				"list of tag names. Returns the created site group.",
			Instructions: "Use only when the caller explicitly wants to create a site group and provides " +
				"a name. This tool mutates NetBox; confirm the name and key attributes before calling.",
		},
		{
			Name:  "update_site_group",
			Title: "Update Site Group",
			Description: "Partially update a site group in NetBox by its numeric ID. Only the fields you " +
				"explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated site group.",
			Instructions: "Use to modify one or more fields of an existing site group. Provide id plus " +
				"only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_site_group",
			Title: "Delete Site Group",
			Description: "WARNING: This operation is irreversible. Permanently deletes a site group (and " +
				"its dependent objects) from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the site group, ideally by ID, and understands it is permanent and cannot " +
				"be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_rack_group",
			Title: "Create Rack Group",
			Description: "Create a new rack group in NetBox. name is required; all other fields are " +
				"optional. Provide site and parent as numeric IDs (look them up first), and tags as a " +
				"list of tag names. Returns the created rack group.",
			Instructions: "Use only when the caller explicitly wants to create a rack group and provides " +
				"a name. This tool mutates NetBox; confirm the name and key attributes before calling.",
		},
		{
			Name:  "update_rack_group",
			Title: "Update Rack Group",
			Description: "Partially update a rack group in NetBox by its numeric ID. Only the fields you " +
				"explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated rack group.",
			Instructions: "Use to modify one or more fields of an existing rack group. Provide id plus " +
				"only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_rack_group",
			Title: "Delete Rack Group",
			Description: "WARNING: This operation is irreversible. Permanently deletes a rack group (and " +
				"its dependent objects) from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the rack group, ideally by ID, and understands it is permanent and cannot " +
				"be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_rack_role",
			Title: "Create Rack Role",
			Description: "Create a new rack role in NetBox. name is required; all other fields are " +
				"optional. Provide slug/color/description as strings, and tags as a list of tag names. " +
				"Returns the created rack role.",
			Instructions: "Use only when the caller explicitly wants to create a rack role and provides a " +
				"name. This tool mutates NetBox; confirm the key attributes before calling.",
		},
		{
			Name:  "update_rack_role",
			Title: "Update Rack Role",
			Description: "Partially update a rack role in NetBox by its numeric ID. Only the fields you " +
				"explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated) are ignored. Returns the updated rack role.",
			Instructions: "Use to modify one or more fields of an existing rack role. Provide id plus " +
				"only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_rack_role",
			Title: "Delete Rack Role",
			Description: "WARNING: This operation is irreversible. Permanently deletes a rack role from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the rack role, ideally by ID, and understands it is permanent and cannot " +
				"be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_device_role",
			Title: "Create Device Role",
			Description: "Create a new device role in NetBox. name is required; all other fields are " +
				"optional. Provide slug/color/description as strings, vm_role as a boolean, weight as an " +
				"integer, config_template as a numeric ID, and tags as a list of tag names. Returns the " +
				"created device role.",
			Instructions: "Use only when the caller explicitly wants to create a device role and provides " +
				"a name. This tool mutates NetBox; confirm the key attributes before calling.",
		},
		{
			Name:  "update_device_role",
			Title: "Update Device Role",
			Description: "Partially update a device role in NetBox by its numeric ID. Only the fields you " +
				"explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated) are ignored. Returns the updated device role.",
			Instructions: "Use to modify one or more fields of an existing device role. Provide id plus " +
				"only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_device_role",
			Title: "Delete Device Role",
			Description: "WARNING: This operation is irreversible. Permanently deletes a device role from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the device role, ideally by ID, and understands it is permanent and cannot " +
				"be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_platform",
			Title: "Create Platform",
			Description: "Create a new platform in NetBox. name is required; all other fields are " +
				"optional. Provide manufacturer and config_template as numeric IDs (look them up first), " +
				"slug/description as strings, and tags as a list of tag names. Returns the created platform.",
			Instructions: "Use only when the caller explicitly wants to create a platform and provides a " +
				"name. This tool mutates NetBox; confirm the key attributes before calling.",
		},
		{
			Name:  "update_platform",
			Title: "Update Platform",
			Description: "Partially update a platform in NetBox by its numeric ID. Only the fields you " +
				"explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated) are ignored. Returns the updated platform.",
			Instructions: "Use to modify one or more fields of an existing platform. Provide id plus only " +
				"the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_platform",
			Title: "Delete Platform",
			Description: "WARNING: This operation is irreversible. Permanently deletes a platform from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the platform, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_virtual_chassis",
			Title: "Create Virtual Chassis",
			Description: "Create a new virtual chassis in NetBox. name is optional here (unlike most " +
				"entities); all other fields are optional. Provide master as a numeric device ID (look it " +
				"up first), domain/description as strings, and tags as a list of tag names. Returns the " +
				"created virtual chassis.",
			Instructions: "Use only when the caller explicitly wants to create a virtual chassis. This " +
				"tool mutates NetBox; confirm the key attributes before calling.",
		},
		{
			Name:  "update_virtual_chassis",
			Title: "Update Virtual Chassis",
			Description: "Partially update a virtual chassis in NetBox by its numeric ID. Only the fields " +
				"you explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, member_count) are ignored. Returns the " +
				"updated virtual chassis.",
			Instructions: "Use to modify one or more fields of an existing virtual chassis. Provide id " +
				"plus only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_virtual_chassis",
			Title: "Delete Virtual Chassis",
			Description: "WARNING: This operation is irreversible. Permanently deletes a virtual chassis " +
				"from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the virtual chassis, ideally by ID, and understands it is permanent and " +
				"cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming " +
				"deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "create_power_panel",
			Title: "Create Power Panel",
			Description: "Create a new power panel in NetBox. name is required; all other fields are " +
				"optional. Provide site and location as numeric IDs (look them up first), description as " +
				"a string, and tags as a list of tag names. Returns the created power panel.",
			Instructions: "Use only when the caller explicitly wants to create a power panel and provides " +
				"a name. This tool mutates NetBox; confirm the key attributes before calling.",
		},
		{
			Name:  "update_power_panel",
			Title: "Update Power Panel",
			Description: "Partially update a power panel in NetBox by its numeric ID. Only the fields you " +
				"explicitly provide are changed (PATCH merge); omitted fields are left untouched. " +
				"Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated power panel.",
			Instructions: "Use to modify one or more fields of an existing power panel. Provide id plus " +
				"only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_power_panel",
			Title: "Delete Power Panel",
			Description: "WARNING: This operation is irreversible. Permanently deletes a power panel from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the power panel, ideally by ID, and understands it is permanent and cannot " +
				"be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "create_inventory_item_role",
			Title: "Create Inventory Item Role",
			Description: "Create a new inventory item role in NetBox. name is required; all other fields " +
				"are optional. Provide slug/color/description as strings, and tags as a list of tag names. " +
				"Returns the created inventory item role.",
			Instructions: "Use only when the caller explicitly wants to create an inventory item role and " +
				"provides a name. This tool mutates NetBox; confirm the key attributes before calling.",
		},
		{
			Name:  "update_inventory_item_role",
			Title: "Update Inventory Item Role",
			Description: "Partially update an inventory item role in NetBox by its numeric ID. Only the " +
				"fields you explicitly provide are changed (PATCH merge); omitted fields are left " +
				"untouched. Read-only fields (id, url, created, last_updated) are ignored. Returns the " +
				"updated inventory item role.",
			Instructions: "Use to modify one or more fields of an existing inventory item role. Provide id " +
				"plus only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:  "delete_inventory_item_role",
			Title: "Delete Inventory Item Role",
			Description: "WARNING: This operation is irreversible. Permanently deletes an inventory item " +
				"role from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the inventory item role, ideally by ID, and understands it is permanent " +
				"and cannot be undone. If the user merely asks to 'remove' or 'clean up' without " +
				"confirming deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "get_front_ports",
			Title: "List Front Ports",
			Description: "List front ports in NetBox. Filter by device or tag, or use q for free-text search. " +
				"Results are paginated: set page and page_size (max 1000). An empty results array means no " +
				"front ports match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox front ports. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple front ports or don't yet know the exact ID.",
		},
		{
			Name:  "create_front_port",
			Title: "Create Front Port",
			Description: "Create a front port in NetBox. Name and device are required; type is a choice " +
				"string (e.g. lc, sc). Provide rear_port and rear_port_position to link a rear port.",
			Instructions: "Use to add a new front port to a device.",
		},
		{
			Name:  "update_front_port",
			Title: "Update Front Port",
			Description: "Partially update a front port in NetBox by its numeric ID. Only the provided " +
				"fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing front port.",
		},
		{
			Name:  "delete_front_port",
			Title: "Delete Front Port",
			Description: "WARNING: This operation is irreversible. Permanently deletes a front port from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the front port, ideally by ID, and understands it is permanent and cannot " +
				"be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "get_rear_ports",
			Title: "List Rear Ports",
			Description: "List rear ports in NetBox. Filter by device or tag, or use q for free-text search. " +
				"Results are paginated: set page and page_size (max 1000). An empty results array means no " +
				"rear ports match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox rear ports. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple rear ports or don't yet know the exact ID.",
		},
		{
			Name:  "create_rear_port",
			Title: "Create Rear Port",
			Description: "Create a rear port in NetBox. Name and device are required; type is a choice " +
				"string (e.g. lc, sc) and positions sets the number of positions.",
			Instructions: "Use to add a new rear port to a device.",
		},
		{
			Name:  "update_rear_port",
			Title: "Update Rear Port",
			Description: "Partially update a rear port in NetBox by its numeric ID. Only the provided " +
				"fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing rear port.",
		},
		{
			Name:  "delete_rear_port",
			Title: "Delete Rear Port",
			Description: "WARNING: This operation is irreversible. Permanently deletes a rear port from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the rear port, ideally by ID, and understands it is permanent and cannot " +
				"be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "get_console_ports",
			Title: "List Console Ports",
			Description: "List console ports in NetBox. Filter by device or tag, or use q for free-text " +
				"search. Results are paginated: set page and page_size (max 1000). An empty results array " +
				"means no console ports match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox console ports. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple console ports or don't yet know the exact ID.",
		},
		{
			Name:  "create_console_port",
			Title: "Create Console Port",
			Description: "Create a console port in NetBox. Name and device are required; type is a choice " +
				"string (e.g. de-9, usb-c).",
			Instructions: "Use to add a new console port to a device.",
		},
		{
			Name:  "update_console_port",
			Title: "Update Console Port",
			Description: "Partially update a console port in NetBox by its numeric ID. Only the provided " +
				"fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing console port.",
		},
		{
			Name:  "delete_console_port",
			Title: "Delete Console Port",
			Description: "WARNING: This operation is irreversible. Permanently deletes a console port from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the console port, ideally by ID, and understands it is permanent and cannot " +
				"be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "get_console_server_ports",
			Title: "List Console Server Ports",
			Description: "List console server ports in NetBox. Filter by device or tag, or use q for " +
				"free-text search. Results are paginated: set page and page_size (max 1000). An empty " +
				"results array means no console server ports match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox console server ports. Prefer this read-only " +
				"tool over get_object_by_id when you need multiple console server ports or don't yet know " +
				"the exact ID.",
		},
		{
			Name:  "create_console_server_port",
			Title: "Create Console Server Port",
			Description: "Create a console server port in NetBox. Name and device are required; type is a " +
				"choice string.",
			Instructions: "Use to add a new console server port to a device.",
		},
		{
			Name:  "update_console_server_port",
			Title: "Update Console Server Port",
			Description: "Partially update a console server port in NetBox by its numeric ID. Only the " +
				"provided fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing console server port.",
		},
		{
			Name:  "delete_console_server_port",
			Title: "Delete Console Server Port",
			Description: "WARNING: This operation is irreversible. Permanently deletes a console server " +
				"port from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the console server port, ideally by ID, and understands it is permanent " +
				"and cannot be undone. If the user merely asks to 'remove' or 'clean up' without " +
				"confirming deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "get_power_ports",
			Title: "List Power Ports",
			Description: "List power ports in NetBox. Filter by device or tag, or use q for free-text " +
				"search. Results are paginated: set page and page_size (max 1000). An empty results array " +
				"means no power ports match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox power ports. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple power ports or don't yet know the exact ID.",
		},
		{
			Name:  "create_power_port",
			Title: "Create Power Port",
			Description: "Create a power port in NetBox. Name and device are required; type is a choice " +
				"string. Provide maximum_draw and allocated_draw (in VA).",
			Instructions: "Use to add a new power port to a device.",
		},
		{
			Name:  "update_power_port",
			Title: "Update Power Port",
			Description: "Partially update a power port in NetBox by its numeric ID. Only the provided " +
				"fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing power port.",
		},
		{
			Name:  "delete_power_port",
			Title: "Delete Power Port",
			Description: "WARNING: This operation is irreversible. Permanently deletes a power port from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the power port, ideally by ID, and understands it is permanent and cannot " +
				"be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "get_power_outlets",
			Title: "List Power Outlets",
			Description: "List power outlets in NetBox. Filter by device or tag, or use q for free-text " +
				"search. Results are paginated: set page and page_size (max 1000). An empty results array " +
				"means no power outlets match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox power outlets. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple power outlets or don't yet know the exact ID.",
		},
		{
			Name:  "create_power_outlet",
			Title: "Create Power Outlet",
			Description: "Create a power outlet in NetBox. Name and device are required; type is a choice " +
				"string. Provide power_port to link a power port and feed_leg for the feed leg.",
			Instructions: "Use to add a new power outlet to a device.",
		},
		{
			Name:  "update_power_outlet",
			Title: "Update Power Outlet",
			Description: "Partially update a power outlet in NetBox by its numeric ID. Only the provided " +
				"fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing power outlet.",
		},
		{
			Name:  "delete_power_outlet",
			Title: "Delete Power Outlet",
			Description: "WARNING: This operation is irreversible. Permanently deletes a power outlet from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the power outlet, ideally by ID, and understands it is permanent and cannot " +
				"be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "get_power_feeds",
			Title: "List Power Feeds",
			Description: "List power feeds in NetBox. Filter by power panel or tag, or use q for free-text " +
				"search. Results are paginated: set page and page_size (max 1000). An empty results array " +
				"means no power feeds match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox power feeds. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple power feeds or don't yet know the exact ID.",
		},
		{
			Name:  "create_power_feed",
			Title: "Create Power Feed",
			Description: "Create a power feed in NetBox. Name is required; provide power_panel to attach " +
				"the feed to a power panel, plus voltage, amperage and other electrical properties.",
			Instructions: "Use to add a new power feed.",
		},
		{
			Name:  "update_power_feed",
			Title: "Update Power Feed",
			Description: "Partially update a power feed in NetBox by its numeric ID. Only the provided " +
				"fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing power feed.",
		},
		{
			Name:  "delete_power_feed",
			Title: "Delete Power Feed",
			Description: "WARNING: This operation is irreversible. Permanently deletes a power feed from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the power feed, ideally by ID, and understands it is permanent and cannot " +
				"be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "get_device_bays",
			Title: "List Device Bays",
			Description: "List device bays in NetBox. Filter by device or tag, or use q for free-text " +
				"search. Results are paginated: set page and page_size (max 1000). An empty results array " +
				"means no device bays match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox device bays. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple device bays or don't yet know the exact ID.",
		},
		{
			Name:         "create_device_bay",
			Title:        "Create Device Bay",
			Description:  "Create a device bay in NetBox. Name and device are required.",
			Instructions: "Use to add a new device bay to a device.",
		},
		{
			Name:  "update_device_bay",
			Title: "Update Device Bay",
			Description: "Partially update a device bay in NetBox by its numeric ID. Only the provided " +
				"fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing device bay.",
		},
		{
			Name:  "delete_device_bay",
			Title: "Delete Device Bay",
			Description: "WARNING: This operation is irreversible. Permanently deletes a device bay from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the device bay, ideally by ID, and understands it is permanent and cannot " +
				"be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "get_rack_reservations",
			Title: "List Rack Reservations",
			Description: "List rack reservations in NetBox. Filter by rack or tag, or use q for free-text " +
				"search. Results are paginated: set page and page_size (max 1000). An empty results array " +
				"means no rack reservations match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox rack reservations. Prefer this read-only tool " +
				"over get_object_by_id when you need multiple rack reservations or don't yet know the exact ID.",
		},
		{
			Name:  "create_rack_reservation",
			Title: "Create Rack Reservation",
			Description: "Create a rack reservation in NetBox. Description is required; provide rack and " +
				"units to reserve specific rack units.",
			Instructions: "Use to reserve space in a rack.",
		},
		{
			Name:  "update_rack_reservation",
			Title: "Update Rack Reservation",
			Description: "Partially update a rack reservation in NetBox by its numeric ID. Only the " +
				"provided fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing rack reservation.",
		},
		{
			Name:  "delete_rack_reservation",
			Title: "Delete Rack Reservation",
			Description: "WARNING: This operation is irreversible. Permanently deletes a rack reservation " +
				"from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the rack reservation, ideally by ID, and understands it is permanent and " +
				"cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming " +
				"deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "get_inventory_items",
			Title: "List Inventory Items",
			Description: "List inventory items in NetBox. Filter by device or tag, or use q for free-text " +
				"search. Results are paginated: set page and page_size (max 1000). An empty results array " +
				"means no inventory items match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox inventory items. Prefer this read-only tool " +
				"over get_object_by_id when you need multiple inventory items or don't yet know the exact ID.",
		},
		{
			Name:  "create_inventory_item",
			Title: "Create Inventory Item",
			Description: "Create an inventory item in NetBox. Name and device are required; provide " +
				"manufacturer, serial and part_id as applicable.",
			Instructions: "Use to add a new inventory item to a device.",
		},
		{
			Name:  "update_inventory_item",
			Title: "Update Inventory Item",
			Description: "Partially update an inventory item in NetBox by its numeric ID. Only the " +
				"provided fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing inventory item.",
		},
		{
			Name:  "delete_inventory_item",
			Title: "Delete Inventory Item",
			Description: "WARNING: This operation is irreversible. Permanently deletes an inventory item " +
				"from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the inventory item, ideally by ID, and understands it is permanent and " +
				"cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming " +
				"deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "get_mac_addresses",
			Title: "List MAC Addresses",
			Description: "List MAC addresses in NetBox. Filter by tag, or use q for free-text search. " +
				"Results are paginated: set page and page_size (max 1000). An empty results array means no " +
				"MAC addresses match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox MAC addresses. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple MAC addresses or don't yet know the exact ID.",
		},
		{
			Name:  "create_mac_address",
			Title: "Create MAC Address",
			Description: "Create a MAC address in NetBox. Address is required; provide " +
				"assigned_object_type (e.g. dcim.interface) and assigned_object_id to assign the MAC to an " +
				"object.",
			Instructions: "Use to add a new MAC address.",
		},
		{
			Name:  "update_mac_address",
			Title: "Update MAC Address",
			Description: "Partially update a MAC address in NetBox by its numeric ID. Only the provided " +
				"fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing MAC address.",
		},
		{
			Name:  "delete_mac_address",
			Title: "Delete MAC Address",
			Description: "WARNING: This operation is irreversible. Permanently deletes a MAC address from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the MAC address, ideally by ID, and understands it is permanent and cannot " +
				"be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "get_rirs",
			Title: "List RIRs",
			Description: "List regional internet registries (RIRs) in NetBox. Filter by tag, or use q for " +
				"free-text search. Results are paginated: set page and page_size (max 1000). An empty results " +
				"array means no RIRs match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox RIRs. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple RIRs or don't yet know the exact ID.",
		},
		{
			Name:  "create_rir",
			Title: "Create RIR",
			Description: "Create a regional internet registry (RIR) in NetBox. Name is required; provide " +
				"slug, is_private, description, tags and custom_fields as needed.",
			Instructions: "Use to add a new RIR.",
		},
		{
			Name:  "update_rir",
			Title: "Update RIR",
			Description: "Partially update an RIR in NetBox by its numeric ID. Only the provided fields are " +
				"patched (PATCH merge).",
			Instructions: "Use to modify an existing RIR.",
		},
		{
			Name:  "delete_rir",
			Title: "Delete RIR",
			Description: "WARNING: This operation is irreversible. Permanently deletes an RIR from NetBox " +
				"by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the RIR, ideally by ID, and understands it is permanent and cannot be undone. " +
				"If the user merely asks to 'remove' or 'clean up' without confirming deletion, ask for " +
				"explicit confirmation first.",
		},
		{
			Name:  "get_contact_groups",
			Title: "List Contact Groups",
			Description: "List tenancy contact groups in NetBox. Filter by parent or tag, or use q for " +
				"free-text search. Results are paginated: set page and page_size (max 1000). An empty results " +
				"array means no contact groups match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox contact groups. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple contact groups or don't yet know the exact ID.",
		},
		{
			Name:  "create_contact_group",
			Title: "Create Contact Group",
			Description: "Create a tenancy contact group in NetBox. Name is required; provide slug, parent, " +
				"description, tags and custom_fields as needed.",
			Instructions: "Use to add a new contact group.",
		},
		{
			Name:  "update_contact_group",
			Title: "Update Contact Group",
			Description: "Partially update a tenancy contact group in NetBox by its numeric ID. Only the " +
				"provided fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing contact group.",
		},
		{
			Name:  "delete_contact_group",
			Title: "Delete Contact Group",
			Description: "WARNING: This operation is irreversible. Permanently deletes a contact group from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the contact group, ideally by ID, and understands it is permanent and cannot " +
				"be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, ask " +
				"for explicit confirmation first.",
		},
		{
			Name:  "get_contact_roles",
			Title: "List Contact Roles",
			Description: "List tenancy contact roles in NetBox. Filter by tag, or use q for free-text search. " +
				"Results are paginated: set page and page_size (max 1000). An empty results array means no " +
				"contact roles match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox contact roles. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple contact roles or don't yet know the exact ID.",
		},
		{
			Name:  "create_contact_role",
			Title: "Create Contact Role",
			Description: "Create a tenancy contact role in NetBox. Name is required; provide slug, " +
				"description, tags and custom_fields as needed.",
			Instructions: "Use to add a new contact role.",
		},
		{
			Name:  "update_contact_role",
			Title: "Update Contact Role",
			Description: "Partially update a tenancy contact role in NetBox by its numeric ID. Only the " +
				"provided fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing contact role.",
		},
		{
			Name:  "delete_contact_role",
			Title: "Delete Contact Role",
			Description: "WARNING: This operation is irreversible. Permanently deletes a contact role from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the contact role, ideally by ID, and understands it is permanent and cannot " +
				"be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, ask " +
				"for explicit confirmation first.",
		},
		{
			Name:  "get_virtual_disks",
			Title: "List Virtual Disks",
			Description: "List virtualization virtual disks in NetBox. Filter by virtual_machine or tag, or " +
				"use q for free-text search. Results are paginated: set page and page_size (max 1000). An empty " +
				"results array means no virtual disks match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox virtual disks. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple virtual disks or don't yet know the exact ID.",
		},
		{
			Name:  "create_virtual_disk",
			Title: "Create Virtual Disk",
			Description: "Create a virtualization virtual disk in NetBox. Name is required; provide " +
				"virtual_machine, size, description, tags and custom_fields as needed.",
			Instructions: "Use to add a new virtual disk.",
		},
		{
			Name:  "update_virtual_disk",
			Title: "Update Virtual Disk",
			Description: "Partially update a virtualization virtual disk in NetBox by its numeric ID. Only " +
				"the provided fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing virtual disk.",
		},
		{
			Name:  "delete_virtual_disk",
			Title: "Delete Virtual Disk",
			Description: "WARNING: This operation is irreversible. Permanently deletes a virtual disk from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the virtual disk, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, ask for " +
				"explicit confirmation first.",
		},
		{
			Name:  "get_virtual_machine_types",
			Title: "List Virtual Machine Types",
			Description: "List virtualization virtual machine types in NetBox. Filter by tag, or use q for " +
				"free-text search. Results are paginated: set page and page_size (max 1000). An empty results " +
				"array means no virtual machine types match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox virtual machine types. Prefer this read-only " +
				"tool over get_object_by_id when you need multiple virtual machine types or don't yet know the " +
				"exact ID.",
		},
		{
			Name:  "create_virtual_machine_type",
			Title: "Create Virtual Machine Type",
			Description: "Create a virtualization virtual machine type in NetBox. Name is required; provide " +
				"slug, description, tags and custom_fields as needed.",
			Instructions: "Use to add a new virtual machine type.",
		},
		{
			Name:  "update_virtual_machine_type",
			Title: "Update Virtual Machine Type",
			Description: "Partially update a virtualization virtual machine type in NetBox by its numeric " +
				"ID. Only the provided fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing virtual machine type.",
		},
		{
			Name:  "delete_virtual_machine_type",
			Title: "Delete Virtual Machine Type",
			Description: "WARNING: This operation is irreversible. Permanently deletes a virtual machine " +
				"type from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the virtual machine type, ideally by ID, and understands it is permanent and " +
				"cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming " +
				"deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "get_aggregates",
			Title: "List Aggregates",
			Description: "List IPAM aggregates in NetBox. Filter by rir, tenant or tag, or use q for " +
				"free-text search. Results are paginated: set page and page_size (max 1000). An empty results " +
				"array means no aggregates match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox aggregates. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple aggregates or don't yet know the exact ID.",
		},
		{
			Name:  "create_aggregate",
			Title: "Create Aggregate",
			Description: "Create an IPAM aggregate in NetBox. Prefix is required; provide rir, tenant, " +
				"date_added, description, tags and custom_fields as needed.",
			Instructions: "Use to add a new aggregate.",
		},
		{
			Name:  "update_aggregate",
			Title: "Update Aggregate",
			Description: "Partially update an aggregate in NetBox by its numeric ID. Only the provided " +
				"fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing aggregate.",
		},
		{
			Name:  "delete_aggregate",
			Title: "Delete Aggregate",
			Description: "WARNING: This operation is irreversible. Permanently deletes an aggregate from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the aggregate, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, ask " +
				"for explicit confirmation first.",
		},
		{
			Name:  "get_ip_ranges",
			Title: "List IP Ranges",
			Description: "List IPAM IP ranges in NetBox. Filter by vrf, tenant, status, role or tag, or use " +
				"q for free-text search. Results are paginated: set page and page_size (max 1000). An empty " +
				"results array means no IP ranges match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox IP ranges. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple IP ranges or don't yet know the exact ID.",
		},
		{
			Name:  "create_ip_range",
			Title: "Create IP Range",
			Description: "Create an IPAM IP range in NetBox. Start address and end address are required; " +
				"provide vrf, tenant, status, role, description, tags and custom_fields as needed.",
			Instructions: "Use to add a new IP range.",
		},
		{
			Name:  "update_ip_range",
			Title: "Update IP Range",
			Description: "Partially update an IP range in NetBox by its numeric ID. Only the provided " +
				"fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing IP range.",
		},
		{
			Name:  "delete_ip_range",
			Title: "Delete IP Range",
			Description: "WARNING: This operation is irreversible. Permanently deletes an IP range from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the IP range, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, ask " +
				"for explicit confirmation first.",
		},
		{
			Name:  "get_asns",
			Title: "List ASNs",
			Description: "List autonomous system numbers (ASNs) in NetBox. Filter by rir, tenant or tag, or " +
				"use q for free-text search. Results are paginated: set page and page_size (max 1000). An " +
				"empty results array means no ASNs match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox ASNs. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple ASNs or don't yet know the exact ID.",
		},
		{
			Name:  "create_asn",
			Title: "Create ASN",
			Description: "Create an autonomous system number (ASN) in NetBox. Asn is required; provide " +
				"rir, tenant, description, tags and custom_fields as needed.",
			Instructions: "Use to add a new ASN.",
		},
		{
			Name:  "update_asn",
			Title: "Update ASN",
			Description: "Partially update an ASN in NetBox by its numeric ID. Only the provided fields are " +
				"patched (PATCH merge).",
			Instructions: "Use to modify an existing ASN.",
		},
		{
			Name:  "delete_asn",
			Title: "Delete ASN",
			Description: "WARNING: This operation is irreversible. Permanently deletes an ASN from NetBox " +
				"by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the ASN, ideally by ID, and understands it is permanent and cannot be undone. " +
				"If the user merely asks to 'remove' or 'clean up' without confirming deletion, ask for " +
				"explicit confirmation first.",
		},
		{
			Name:  "get_asn_ranges",
			Title: "List ASN Ranges",
			Description: "List IPAM ASN ranges in NetBox. Filter by rir, tenant or tag, or use q for " +
				"free-text search. Results are paginated: set page and page_size (max 1000). An empty results " +
				"array means no ASN ranges match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox ASN ranges. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple ASN ranges or don't yet know the exact ID.",
		},
		{
			Name:  "create_asn_range",
			Title: "Create ASN Range",
			Description: "Create an IPAM ASN range in NetBox. Name is required; provide start, end, rir, " +
				"tenant, description, tags and custom_fields as needed.",
			Instructions: "Use to add a new ASN range.",
		},
		{
			Name:  "update_asn_range",
			Title: "Update ASN Range",
			Description: "Partially update an ASN range in NetBox by its numeric ID. Only the provided " +
				"fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing ASN range.",
		},
		{
			Name:  "delete_asn_range",
			Title: "Delete ASN Range",
			Description: "WARNING: This operation is irreversible. Permanently deletes an ASN range from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the ASN range, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, ask " +
				"for explicit confirmation first.",
		},
		{
			Name:  "get_route_targets",
			Title: "List Route Targets",
			Description: "List IPAM route targets in NetBox. Filter by tenant or tag, or use q for " +
				"free-text search. Results are paginated: set page and page_size (max 1000). An empty results " +
				"array means no route targets match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox route targets. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple route targets or don't yet know the exact ID.",
		},
		{
			Name:  "create_route_target",
			Title: "Create Route Target",
			Description: "Create an IPAM route target in NetBox. Name is required; provide tenant, " +
				"description, tags and custom_fields as needed.",
			Instructions: "Use to add a new route target.",
		},
		{
			Name:  "update_route_target",
			Title: "Update Route Target",
			Description: "Partially update a route target in NetBox by its numeric ID. Only the provided " +
				"fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing route target.",
		},
		{
			Name:  "delete_route_target",
			Title: "Delete Route Target",
			Description: "WARNING: This operation is irreversible. Permanently deletes a route target from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the route target, ideally by ID, and understands it is permanent and cannot " +
				"be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "get_services",
			Title: "List Services",
			Description: "List IPAM services in NetBox. Filter by device, virtual_machine or tag, or use q " +
				"for free-text search. Results are paginated: set page and page_size (max 1000). An empty " +
				"results array means no services match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox services. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple services or don't yet know the exact ID.",
		},
		{
			Name:  "create_service",
			Title: "Create Service",
			Description: "Create an IPAM service in NetBox. Name is required; provide device, " +
				"virtual_machine, protocol, ports, description, tags and custom_fields as needed.",
			Instructions: "Use to add a new service.",
		},
		{
			Name:  "update_service",
			Title: "Update Service",
			Description: "Partially update a service in NetBox by its numeric ID. Only the provided fields " +
				"are patched (PATCH merge).",
			Instructions: "Use to modify an existing service.",
		},
		{
			Name:  "delete_service",
			Title: "Delete Service",
			Description: "WARNING: This operation is irreversible. Permanently deletes a service from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the service, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, ask " +
				"for explicit confirmation first.",
		},
		{
			Name:  "get_service_templates",
			Title: "List Service Templates",
			Description: "List IPAM service templates in NetBox. Filter by protocol or tag, or use q for " +
				"free-text search. Results are paginated: set page and page_size (max 1000). An empty results " +
				"array means no service templates match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox service templates. Prefer this read-only tool " +
				"over get_object_by_id when you need multiple service templates or don't yet know the exact ID.",
		},
		{
			Name:  "create_service_template",
			Title: "Create Service Template",
			Description: "Create an IPAM service template in NetBox. Name is required; provide protocol, " +
				"ports, description, tags and custom_fields as needed.",
			Instructions: "Use to add a new service template.",
		},
		{
			Name:  "update_service_template",
			Title: "Update Service Template",
			Description: "Partially update a service template in NetBox by its numeric ID. Only the provided " +
				"fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing service template.",
		},
		{
			Name:  "delete_service_template",
			Title: "Delete Service Template",
			Description: "WARNING: This operation is irreversible. Permanently deletes a service template " +
				"from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the service template, ideally by ID, and understands it is permanent and " +
				"cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming " +
				"deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "get_fhrp_groups",
			Title: "List FHRP Groups",
			Description: "List IPAM FHRP (first-hop redundancy protocol) groups in NetBox. Filter by " +
				"protocol or tag, or use q for free-text search. Results are paginated: set page and page_size " +
				"(max 1000). An empty results array means no FHRP groups match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox FHRP groups. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple FHRP groups or don't yet know the exact ID.",
		},
		{
			Name:  "create_fhrp_group",
			Title: "Create FHRP Group",
			Description: "Create an IPAM FHRP group in NetBox. Name is required; provide protocol, group_id, " +
				"description, tags and custom_fields as needed.",
			Instructions: "Use to add a new FHRP group.",
		},
		{
			Name:  "update_fhrp_group",
			Title: "Update FHRP Group",
			Description: "Partially update an FHRP group in NetBox by its numeric ID. Only the provided " +
				"fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing FHRP group.",
		},
		{
			Name:  "delete_fhrp_group",
			Title: "Delete FHRP Group",
			Description: "WARNING: This operation is irreversible. Permanently deletes an FHRP group from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the FHRP group, ideally by ID, and understands it is permanent and cannot be " +
				"undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, ask " +
				"for explicit confirmation first.",
		},
		{
			Name:  "get_fhrp_group_assignments",
			Title: "List FHRP Group Assignments",
			Description: "List IPAM FHRP group assignments in NetBox. Filter by group or ip_address, or use " +
				"q for free-text search. Results are paginated: set page and page_size (max 1000). An empty " +
				"results array means no FHRP group assignments match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox FHRP group assignments. Prefer this read-only " +
				"tool over get_object_by_id when you need multiple FHRP group assignments or don't yet know " +
				"the exact ID.",
		},
		{
			Name:  "create_fhrp_group_assignment",
			Title: "Create FHRP Group Assignment",
			Description: "Create an IPAM FHRP group assignment in NetBox. Group and ip_address are " +
				"required; provide priority as needed.",
			Instructions: "Use to map an FHRP group to an IP address.",
		},
		{
			Name:  "update_fhrp_group_assignment",
			Title: "Update FHRP Group Assignment",
			Description: "Partially update an FHRP group assignment in NetBox by its numeric ID. Only the " +
				"provided fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing FHRP group assignment.",
		},
		{
			Name:  "delete_fhrp_group_assignment",
			Title: "Delete FHRP Group Assignment",
			Description: "WARNING: This operation is irreversible. Permanently deletes an FHRP group " +
				"assignment from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the FHRP group assignment, ideally by ID, and understands it is permanent " +
				"and cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming " +
				"deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "get_vlan_translation_policies",
			Title: "List VLAN Translation Policies",
			Description: "List IPAM VLAN translation policies in NetBox. Filter by tag, or use q for " +
				"free-text search. Results are paginated: set page and page_size (max 1000). An empty results " +
				"array means no VLAN translation policies match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox VLAN translation policies. Prefer this " +
				"read-only tool over get_object_by_id when you need multiple policies or don't yet know the " +
				"exact ID.",
		},
		{
			Name:  "create_vlan_translation_policy",
			Title: "Create VLAN Translation Policy",
			Description: "Create an IPAM VLAN translation policy in NetBox. Name is required; provide " +
				"description, tags and custom_fields as needed.",
			Instructions: "Use to add a new VLAN translation policy.",
		},
		{
			Name:  "update_vlan_translation_policy",
			Title: "Update VLAN Translation Policy",
			Description: "Partially update a VLAN translation policy in NetBox by its numeric ID. Only the " +
				"provided fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing VLAN translation policy.",
		},
		{
			Name:  "delete_vlan_translation_policy",
			Title: "Delete VLAN Translation Policy",
			Description: "WARNING: This operation is irreversible. Permanently deletes a VLAN translation " +
				"policy from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the VLAN translation policy, ideally by ID, and understands it is permanent " +
				"and cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming " +
				"deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "get_vlan_translation_rules",
			Title: "List VLAN Translation Rules",
			Description: "List IPAM VLAN translation rules in NetBox. Filter by policy, or use q for " +
				"free-text search. Results are paginated: set page and page_size (max 1000). An empty results " +
				"array means no VLAN translation rules match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox VLAN translation rules. Prefer this read-only " +
				"tool over get_object_by_id when you need multiple rules or don't yet know the exact ID.",
		},
		{
			Name:  "create_vlan_translation_rule",
			Title: "Create VLAN Translation Rule",
			Description: "Create an IPAM VLAN translation rule in NetBox. Policy is required; provide " +
				"local_vid and remote_vid as needed.",
			Instructions: "Use to add a new VLAN translation rule.",
		},
		{
			Name:  "update_vlan_translation_rule",
			Title: "Update VLAN Translation Rule",
			Description: "Partially update a VLAN translation rule in NetBox by its numeric ID. Only the " +
				"provided fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing VLAN translation rule.",
		},
		{
			Name:  "delete_vlan_translation_rule",
			Title: "Delete VLAN Translation Rule",
			Description: "WARNING: This operation is irreversible. Permanently deletes a VLAN translation " +
				"rule from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the VLAN translation rule, ideally by ID, and understands it is permanent " +
				"and cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming " +
				"deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "get_provider_accounts",
			Title: "List Provider Accounts",
			Description: "List circuits provider accounts in NetBox. Filter by provider or tag, or use q " +
				"for free-text search. Results are paginated: set page and page_size (max 1000). An empty " +
				"results array means no provider accounts match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox provider accounts. Prefer this read-only tool " +
				"over get_object_by_id when you need multiple provider accounts or don't yet know the exact ID.",
		},
		{
			Name:  "create_provider_account",
			Title: "Create Provider Account",
			Description: "Create a circuits provider account in NetBox. Name is required; provide provider, " +
				"account, description, tags and custom_fields as needed.",
			Instructions: "Use to add a new provider account.",
		},
		{
			Name:  "update_provider_account",
			Title: "Update Provider Account",
			Description: "Partially update a circuits provider account in NetBox by its numeric ID. Only " +
				"the provided fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing provider account.",
		},
		{
			Name:  "delete_provider_account",
			Title: "Delete Provider Account",
			Description: "WARNING: This operation is irreversible. Permanently deletes a provider account " +
				"from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the provider account, ideally by ID, and understands it is permanent and " +
				"cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming " +
				"deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "get_provider_networks",
			Title: "List Provider Networks",
			Description: "List circuits provider networks in NetBox. Filter by provider or tag, or use q " +
				"for free-text search. Results are paginated: set page and page_size (max 1000). An empty " +
				"results array means no provider networks match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox provider networks. Prefer this read-only tool " +
				"over get_object_by_id when you need multiple provider networks or don't yet know the exact ID.",
		},
		{
			Name:  "create_provider_network",
			Title: "Create Provider Network",
			Description: "Create a circuits provider network in NetBox. Name is required; provide provider, " +
				"description, comments, tags and custom_fields as needed.",
			Instructions: "Use to add a new provider network.",
		},
		{
			Name:  "update_provider_network",
			Title: "Update Provider Network",
			Description: "Partially update a circuits provider network in NetBox by its numeric ID. Only " +
				"the provided fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing provider network.",
		},
		{
			Name:  "delete_provider_network",
			Title: "Delete Provider Network",
			Description: "WARNING: This operation is irreversible. Permanently deletes a provider network " +
				"from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the provider network, ideally by ID, and understands it is permanent and " +
				"cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming " +
				"deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "get_circuit_groups",
			Title: "List Circuit Groups",
			Description: "List circuits circuit groups in NetBox. Filter by tag, or use q for free-text " +
				"search. Results are paginated: set page and page_size (max 1000). An empty results array " +
				"means no circuit groups match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox circuit groups. Prefer this read-only tool over " +
				"get_object_by_id when you need multiple circuit groups or don't yet know the exact ID.",
		},
		{
			Name:  "create_circuit_group",
			Title: "Create Circuit Group",
			Description: "Create a circuits circuit group in NetBox. Name is required; provide description, " +
				"tags and custom_fields as needed.",
			Instructions: "Use to add a new circuit group.",
		},
		{
			Name:  "update_circuit_group",
			Title: "Update Circuit Group",
			Description: "Partially update a circuits circuit group in NetBox by its numeric ID. Only the " +
				"provided fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing circuit group.",
		},
		{
			Name:  "delete_circuit_group",
			Title: "Delete Circuit Group",
			Description: "WARNING: This operation is irreversible. Permanently deletes a circuit group from " +
				"NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the circuit group, ideally by ID, and understands it is permanent and cannot " +
				"be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, " +
				"ask for explicit confirmation first.",
		},
		{
			Name:  "get_circuit_group_assignments",
			Title: "List Circuit Group Assignments",
			Description: "List circuits circuit group assignments in NetBox. Filter by group or circuit, or " +
				"use q for free-text search. Results are paginated: set page and page_size (max 1000). An " +
				"empty results array means no circuit group assignments match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox circuit group assignments. Prefer this " +
				"read-only tool over get_object_by_id when you need multiple assignments or don't yet know " +
				"the exact ID.",
		},
		{
			Name:  "create_circuit_group_assignment",
			Title: "Create Circuit Group Assignment",
			Description: "Create a circuits circuit group assignment in NetBox. Group and circuit are " +
				"required; provide priority as needed.",
			Instructions: "Use to map a circuit to a circuit group.",
		},
		{
			Name:  "update_circuit_group_assignment",
			Title: "Update Circuit Group Assignment",
			Description: "Partially update a circuits circuit group assignment in NetBox by its numeric " +
				"ID. Only the provided fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing circuit group assignment.",
		},
		{
			Name:  "delete_circuit_group_assignment",
			Title: "Delete Circuit Group Assignment",
			Description: "WARNING: This operation is irreversible. Permanently deletes a circuit group " +
				"assignment from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the circuit group assignment, ideally by ID, and understands it is permanent " +
				"and cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming " +
				"deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "get_virtual_circuits",
			Title: "List Virtual Circuits",
			Description: "List circuits virtual circuits in NetBox. Filter by type, provider_network, " +
				"status or tag, or use q for free-text search. Results are paginated: set page and page_size " +
				"(max 1000). An empty results array means no virtual circuits match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox virtual circuits. Prefer this read-only tool " +
				"over get_object_by_id when you need multiple virtual circuits or don't yet know the exact ID.",
		},
		{
			Name:  "create_virtual_circuit",
			Title: "Create Virtual Circuit",
			Description: "Create a circuits virtual circuit in NetBox. Name is required; provide type, " +
				"provider_network, status, description, tags and custom_fields as needed.",
			Instructions: "Use to add a new virtual circuit.",
		},
		{
			Name:  "update_virtual_circuit",
			Title: "Update Virtual Circuit",
			Description: "Partially update a circuits virtual circuit in NetBox by its numeric ID. Only " +
				"the provided fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing virtual circuit.",
		},
		{
			Name:  "delete_virtual_circuit",
			Title: "Delete Virtual Circuit",
			Description: "WARNING: This operation is irreversible. Permanently deletes a virtual circuit " +
				"from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the virtual circuit, ideally by ID, and understands it is permanent and " +
				"cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming " +
				"deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "get_virtual_circuit_types",
			Title: "List Virtual Circuit Types",
			Description: "List circuits virtual circuit types in NetBox. Filter by tag, or use q for " +
				"free-text search. Results are paginated: set page and page_size (max 1000). An empty results " +
				"array means no virtual circuit types match the filters, not an error.",
			Instructions: "Use to enumerate or search NetBox virtual circuit types. Prefer this read-only " +
				"tool over get_object_by_id when you need multiple virtual circuit types or don't yet know " +
				"the exact ID.",
		},
		{
			Name:  "create_virtual_circuit_type",
			Title: "Create Virtual Circuit Type",
			Description: "Create a circuits virtual circuit type in NetBox. Name is required; provide slug, " +
				"description, tags and custom_fields as needed.",
			Instructions: "Use to add a new virtual circuit type.",
		},
		{
			Name:  "update_virtual_circuit_type",
			Title: "Update Virtual Circuit Type",
			Description: "Partially update a circuits virtual circuit type in NetBox by its numeric ID. " +
				"Only the provided fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing virtual circuit type.",
		},
		{
			Name:  "delete_virtual_circuit_type",
			Title: "Delete Virtual Circuit Type",
			Description: "WARNING: This operation is irreversible. Permanently deletes a virtual circuit " +
				"type from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the virtual circuit type, ideally by ID, and understands it is permanent " +
				"and cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming " +
				"deletion, ask for explicit confirmation first.",
		},
		{
			Name:  "get_virtual_circuit_terminations",
			Title: "List Virtual Circuit Terminations",
			Description: "List circuits virtual circuit terminations in NetBox. Filter by virtual_circuit, " +
				"interface_type or interface_id, or use q for free-text search. Results are paginated: set " +
				"page and page_size (max 1000). An empty results array means no terminations match the " +
				"filters, not an error.",
			Instructions: "Use to enumerate or search NetBox virtual circuit terminations. Prefer this " +
				"read-only tool over get_object_by_id when you need multiple terminations or don't yet know " +
				"the exact ID.",
		},
		{
			Name:  "create_virtual_circuit_termination",
			Title: "Create Virtual Circuit Termination",
			Description: "Create a circuits virtual circuit termination in NetBox. Virtual circuit is " +
				"required; provide role, interface_type and interface_id as needed.",
			Instructions: "Use to attach a virtual circuit to an interface.",
		},
		{
			Name:  "update_virtual_circuit_termination",
			Title: "Update Virtual Circuit Termination",
			Description: "Partially update a circuits virtual circuit termination in NetBox by its numeric " +
				"ID. Only the provided fields are patched (PATCH merge).",
			Instructions: "Use to modify an existing virtual circuit termination.",
		},
		{
			Name:  "delete_virtual_circuit_termination",
			Title: "Delete Virtual Circuit Termination",
			Description: "WARNING: This operation is irreversible. Permanently deletes a virtual circuit " +
				"termination from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they " +
				"want to delete the virtual circuit termination, ideally by ID, and understands it is " +
				"permanent and cannot be undone. If the user merely asks to 'remove' or 'clean up' without " +
				"confirming deletion, ask for explicit confirmation first.",
		},
	}
}
