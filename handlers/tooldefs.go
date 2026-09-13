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
		{
			Name:         "create_site",
			Title:        "Create Site",
			Description:  "Create a new site in NetBox. name is required; all other fields are optional. Provide region/tenant as numeric IDs (look them up first), status as one of planned/staged/active/decommissioning/retired, and tags as a list of tag names. Returns the created site.",
			Instructions: "Use only when the caller explicitly wants to create a site and provides a name. This tool mutates NetBox; confirm the name and key attributes before calling.",
		},
		{
			Name:         "update_site",
			Title:        "Update Site",
			Description:  "Partially update a site in NetBox by its numeric ID. Only the fields you explicitly provide are changed (PATCH merge); omitted fields are left untouched. Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated site.",
			Instructions: "Use to modify one or more fields of an existing site. Provide id plus only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:         "delete_site",
			Title:        "Delete Site",
			Description:  "WARNING: This operation is irreversible. Permanently deletes a site (and its dependent objects) from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they want to delete the site, ideally by ID, and understands it is permanent and cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, ask for explicit confirmation first.",
		},
		{
			Name:         "create_device",
			Title:        "Create Device",
			Description:  "Create a new device in NetBox. name is required; all other fields are optional. Provide device_type/role/site/rack/cluster/tenant/platform as numeric IDs (look them up first), status as one of offline/active/planned/staged/failed/inventory/decommissioning, face as front/rear, and tags as a list of tag names. Returns the created device.",
			Instructions: "Use only when the caller explicitly wants to create a device and provides a name. This tool mutates NetBox; confirm the name and key attributes (device_type, role, site) before calling.",
		},
		{
			Name:         "update_device",
			Title:        "Update Device",
			Description:  "Partially update a device in NetBox by its numeric ID. Only the fields you explicitly provide are changed (PATCH merge); omitted fields are left untouched. Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated device.",
			Instructions: "Use to modify one or more fields of an existing device. Provide id plus only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:         "delete_device",
			Title:        "Delete Device",
			Description:  "WARNING: This operation is irreversible. Permanently deletes a device from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they want to delete the device, ideally by ID, and understands it is permanent and cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, ask for explicit confirmation first.",
		},
		{
			Name:         "create_ip_address",
			Title:        "Create IP Address",
			Description:  "Create a new IP address in NetBox. address is required (e.g. 192.168.1.1/24); all other fields are optional. Provide vrf/tenant/assigned_object_id as numeric IDs (look them up first), status as one of active/reserved/deprecated/dhcp/slaac, role as one of loopback/secondary/anycast/vip/vrrp/hsrp/glbp/carp, and tags as a list of tag names. Returns the created IP address.",
			Instructions: "Use only when the caller explicitly wants to create an IP address and provides an address. This tool mutates NetBox; confirm the address and key attributes before calling.",
		},
		{
			Name:         "update_ip_address",
			Title:        "Update IP Address",
			Description:  "Partially update an IP address in NetBox by its numeric ID. Only the fields you explicitly provide are changed (PATCH merge); omitted fields are left untouched. Read-only fields (id, url, created, last_updated, display, assigned_object) are ignored. Returns the updated IP address.",
			Instructions: "Use to modify one or more fields of an existing IP address. Provide id plus only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:         "delete_ip_address",
			Title:        "Delete IP Address",
			Description:  "WARNING: This operation is irreversible. Permanently deletes an IP address from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they want to delete the IP address, ideally by ID, and understands it is permanent and cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, ask for explicit confirmation first.",
		},
		{
			Name:         "create_prefix",
			Title:        "Create Prefix",
			Description:  "Create a new prefix in NetBox. prefix is required (CIDR, e.g. 10.0.0.0/24); all other fields are optional. Provide site/vrf/tenant/vlan/role as numeric IDs (look them up first), status as one of container/active/reserved/deprecated, is_pool as a boolean, and tags as a list of tag names. Returns the created prefix.",
			Instructions: "Use only when the caller explicitly wants to create a prefix and provides a CIDR. This tool mutates NetBox; confirm the prefix and key attributes before calling.",
		},
		{
			Name:         "update_prefix",
			Title:        "Update Prefix",
			Description:  "Partially update a prefix in NetBox by its numeric ID. Only the fields you explicitly provide are changed (PATCH merge); omitted fields are left untouched. Read-only fields (id, url, created, last_updated, display, children, _depth, family) are ignored. Returns the updated prefix.",
			Instructions: "Use to modify one or more fields of an existing prefix. Provide id plus only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:         "delete_prefix",
			Title:        "Delete Prefix",
			Description:  "WARNING: This operation is irreversible. Permanently deletes a prefix from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they want to delete the prefix, ideally by ID, and understands it is permanent and cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, ask for explicit confirmation first.",
		},
		{
			Name:         "create_vlan",
			Title:        "Create VLAN",
			Description:  "Create a new VLAN in NetBox. vid (1-4094) and name are required; all other fields are optional. Provide site/group/tenant/role as numeric IDs (look them up first), status as one of active/reserved/deprecated, and tags as a list of tag names. Returns the created VLAN.",
			Instructions: "Use only when the caller explicitly wants to create a VLAN and provides both a VLAN ID (vid) and a name. This tool mutates NetBox; confirm the vid and key attributes before calling.",
		},
		{
			Name:         "update_vlan",
			Title:        "Update VLAN",
			Description:  "Partially update a VLAN in NetBox by its numeric ID. Only the fields you explicitly provide are changed (PATCH merge); omitted fields are left untouched. Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated VLAN.",
			Instructions: "Use to modify one or more fields of an existing VLAN. Provide id plus only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:         "delete_vlan",
			Title:        "Delete VLAN",
			Description:  "WARNING: This operation is irreversible. Permanently deletes a VLAN from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they want to delete the VLAN, ideally by ID, and understands it is permanent and cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, ask for explicit confirmation first.",
		},
		{
			Name:         "create_virtual_machine",
			Title:        "Create Virtual Machine",
			Description:  "Create a new virtual machine in NetBox. name is required; all other fields are optional. Provide cluster/role/tenant/platform/site as numeric IDs (look them up first), status as one of offline/active/planned/staged/failed/decommissioning, vcpus as a number, memory (MB) and disk (GB) as integers, and tags as a list of tag names. Returns the created virtual machine.",
			Instructions: "Use only when the caller explicitly wants to create a virtual machine and provides a name. This tool mutates NetBox; confirm the name and key attributes before calling.",
		},
		{
			Name:         "update_virtual_machine",
			Title:        "Update Virtual Machine",
			Description:  "Partially update a virtual machine in NetBox by its numeric ID. Only the fields you explicitly provide are changed (PATCH merge); omitted fields are left untouched. Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated virtual machine.",
			Instructions: "Use to modify one or more fields of an existing virtual machine. Provide id plus only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:         "delete_virtual_machine",
			Title:        "Delete Virtual Machine",
			Description:  "WARNING: This operation is irreversible. Permanently deletes a virtual machine from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they want to delete the virtual machine, ideally by ID, and understands it is permanent and cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, ask for explicit confirmation first.",
		},
		{
			Name:         "create_cluster",
			Title:        "Create Cluster",
			Description:  "Create a new cluster in NetBox. name and type (cluster type) are required; all other fields are optional. Provide type/group/site/tenant as numeric IDs (look them up first), and tags as a list of tag names. Returns the created cluster.",
			Instructions: "Use only when the caller explicitly wants to create a cluster and provides a name and a cluster type. This tool mutates NetBox; confirm the name and key attributes before calling.",
		},
		{
			Name:         "update_cluster",
			Title:        "Update Cluster",
			Description:  "Partially update a cluster in NetBox by its numeric ID. Only the fields you explicitly provide are changed (PATCH merge); omitted fields are left untouched. Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated cluster.",
			Instructions: "Use to modify one or more fields of an existing cluster. Provide id plus only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:         "delete_cluster",
			Title:        "Delete Cluster",
			Description:  "WARNING: This operation is irreversible. Permanently deletes a cluster from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they want to delete the cluster, ideally by ID, and understands it is permanent and cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, ask for explicit confirmation first.",
		},
		{
			Name:         "create_circuit",
			Title:        "Create Circuit",
			Description:  "Create a new circuit in NetBox. cid and provider and circuit_type are required; all other fields are optional. Provide provider/circuit_type/tenant as numeric IDs (look them up first), status as one of planned/provisioning/active/offline/decommissioning, install_date as YYYY-MM-DD, commit_rate in kbps, and tags as a list of tag names. Returns the created circuit.",
			Instructions: "Use only when the caller explicitly wants to create a circuit and provides a cid, a provider and a circuit_type. This tool mutates NetBox; confirm the cid and key attributes before calling.",
		},
		{
			Name:         "update_circuit",
			Title:        "Update Circuit",
			Description:  "Partially update a circuit in NetBox by its numeric ID. Only the fields you explicitly provide are changed (PATCH merge); omitted fields are left untouched. Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated circuit.",
			Instructions: "Use to modify one or more fields of an existing circuit. Provide id plus only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:         "delete_circuit",
			Title:        "Delete Circuit",
			Description:  "WARNING: This operation is irreversible. Permanently deletes a circuit from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they want to delete the circuit, ideally by ID, and understands it is permanent and cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, ask for explicit confirmation first.",
		},
		{
			Name:         "create_rack",
			Title:        "Create Rack",
			Description:  "Create a new rack in NetBox. name is required; all other fields are optional. Provide site/location/tenant/role as numeric IDs (look them up first), status as one of reserved/available/planned/active/decommissioning, type as one of 2-post-frame/4-post-frame/4-post-cabinet/wall-frame/wall-cabinet, width in inches, u_height in rack units, and tags as a list of tag names. Returns the created rack.",
			Instructions: "Use only when the caller explicitly wants to create a rack and provides a name. This tool mutates NetBox; confirm the name and key attributes before calling.",
		},
		{
			Name:         "update_rack",
			Title:        "Update Rack",
			Description:  "Partially update a rack in NetBox by its numeric ID. Only the fields you explicitly provide are changed (PATCH merge); omitted fields are left untouched. Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated rack.",
			Instructions: "Use to modify one or more fields of an existing rack. Provide id plus only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:         "delete_rack",
			Title:        "Delete Rack",
			Description:  "WARNING: This operation is irreversible. Permanently deletes a rack from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they want to delete the rack, ideally by ID, and understands it is permanent and cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, ask for explicit confirmation first.",
		},
		{
			Name:         "create_interface",
			Title:        "Create Interface",
			Description:  "Create a new device interface in NetBox. name, device and type are required; all other fields are optional. Provide device as a numeric ID (look it up first), type as one of the NetBox interface type values (e.g. 1000base-t, 10gbase-x-sfpp), and tags as a list of tag names. Returns the created interface.",
			Instructions: "Use only when the caller explicitly wants to create an interface and provides a name, a device and a type. This tool mutates NetBox; confirm the name and key attributes before calling.",
		},
		{
			Name:         "update_interface",
			Title:        "Update Interface",
			Description:  "Partially update a device interface in NetBox by its numeric ID. Only the fields you explicitly provide are changed (PATCH merge); omitted fields are left untouched. Read-only fields (id, url, created, last_updated, display) are ignored. Returns the updated interface.",
			Instructions: "Use to modify one or more fields of an existing interface. Provide id plus only the fields to change. This tool mutates NetBox.",
		},
		{
			Name:         "delete_interface",
			Title:        "Delete Interface",
			Description:  "WARNING: This operation is irreversible. Permanently deletes a device interface from NetBox by its numeric ID. There is no undo.",
			Instructions: "Only call this after the user has explicitly and unambiguously confirmed they want to delete the interface, ideally by ID, and understands it is permanent and cannot be undone. If the user merely asks to 'remove' or 'clean up' without confirming deletion, ask for explicit confirmation first.",
		},
	}
}
