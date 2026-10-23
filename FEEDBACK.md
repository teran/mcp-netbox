# MCP NetBox Server — Feedback

## Missing tool: Region management — ✅ RESOLVED

The server exposes full CRUD for **Sites** (`get_sites`, `create_site`, `update_site`,
`delete_site`), but there were **no tools for NetBox Regions**:

- ~~No `get_regions` (or region lookup)~~
- ~~No `create_region`~~
- ~~No `update_region`~~
- ~~No `delete_region`~~

A Site references a Region through a numeric `region` ID, so an agent could only
attach a Site to a Region that **already exists**; there was no way to enumerate or
create Regions from MCP.

**Status:** Addressed. Region tools have been added, mirroring the Site ones in the
same read → write/update → delete grouping used elsewhere:

- `get_regions` — list regions with optional filters (`parent` region ID, `q`
  free-text, `tag`, pagination), mirroring `get_sites`.
- `create_region` — `name` required, plus `slug`, `parent` (region ID), `description`,
  `tags`, `custom_fields`.
- `update_region` — partial PATCH merge by numeric ID (mirror `update_site`).
- `delete_region` — destructive, guarded like `delete_site`.

`region` was also added to the `get_object_by_id` supported object types
(`/api/dcim/regions/`). The entity set grew from 25 to 26 (93 tools total, 15 read +
78 typed write).

This closes the gap so agents can manage the full site/region hierarchy end-to-end via
MCP instead of falling back to the UI or raw API.

---

## Missing tool: Tenant Group management — ✅ RESOLVED

The server exposes `create_tenant` / `update_tenant` / `delete_tenant` for **Tenants**,
but there were **no tools for NetBox Tenant Groups**:

- ~~No `get_tenant_groups` (or tenant-group lookup)~~
- ~~No `create_tenant_group`~~
- ~~No `update_tenant_group`~~
- ~~No `delete_tenant_group`~~

In addition, the `create_tenant` tool did **not** expose a `group` field.

**Status:** Addressed. Tenant Group tools added in the same read → write/update → delete
grouping, and `group` (numeric tenant-group ID) wired into `create_tenant`/`update_tenant`:

- `get_tenant_groups` — list tenant groups (`parent` ID, `q`, `tag`, pagination).
- `create_tenant_group` — `name` required, plus `slug`, `parent`, `description`, `tags`, `custom_fields`.
- `update_tenant_group` — partial PATCH merge by numeric ID.
- `delete_tenant_group` — destructive, guarded like `delete_tenant`.
- `create_tenant` / `update_tenant` now accept and relay a `group` (numeric tenant-group ID).

`tenant_group` was also added to the `get_object_by_id` supported object types
(`/api/tenancy/tenant-groups/`). The entity set grew from 26 to 27 (97 tools total, 16 read +
81 typed write).

This closes the gap so agents can manage the full tenant/tenant-group hierarchy end-to-end via
MCP instead of falling back to the UI or raw API.

---

## Missing tool: Role & Tenant enumeration (read support)

We model independent service lines (compute, object storage, block storage, managed
databases, managed Kubernetes) as **Roles**, and install environments (incarnation ×
environment) as **Tenants**. Prefixes and VLANs are then tagged with a `role` and a
`tenant` alongside their `site`.

The server has write tools for these (`create_role`/`update_role`/`delete_role` and
`create_tenant`/`update_tenant`/`delete_tenant`), and `update_prefix`/`update_vlan`
expose `role` and `tenant` fields. However, there is **no enumeration/lookup** for
either entity:

- **No `get_roles`** — cannot list existing roles to obtain their numeric IDs.
- **No `get_tenants`** — cannot list existing tenants to obtain their numeric IDs.

Without them, an agent cannot resolve a role/tenant **by name** to an ID before
attaching it to a prefix or VLAN; it must fall back to the UI or raw API
(`/api/ipam/roles/`, `/api/tenancy/tenants/`) just to discover IDs.

**Suggestion:** add read tools mirroring the other list tools:

- `get_roles` — list roles (optional `q`, `tag`, pagination), mirroring `get_sites`.
- `get_tenants` — list tenants (optional `q`, `group`, `tag`, pagination).

Also, for convenience when an environment is expressed as a **Tag** rather than a
Tenant (so an incarnation can share fewer Tenants), consider read/write Tag support if
not already present (`get_tags`, `create_tag`, `update_tag`, `delete_tag`).

This rounds out the model so an agent can, end-to-end from MCP, enumerate roles and
tenants, create them, and attach them (plus site) to prefixes and VLANs — without
falling back to the UI or raw API.

---

## Missing tool: Custom Field management

We want to attach a reference to an external ticket to prefixes and VLANs (for example,
the work item under which the network was allocated). The natural NetBox mechanism is a
**Custom Field** on those objects, populated via `custom_fields` when creating/updating a
prefix or VLAN.

The create/update tools for many object types (prefix, VLAN, tenant, site, etc.) already
expose a `custom_fields` map, so **populating** existing custom field values works via
MCP. However, there are **no tools to manage the custom field definitions themselves**:

- **No `get_custom_fields`** — cannot enumerate existing custom fields.
- **No `create_custom_field`** — cannot define a new custom field.
- **No `update_custom_field`** — cannot alter an existing custom field definition.
- **No `delete_custom_field`** — cannot remove a custom field.

Without `create_custom_field`, an agent cannot set up a new custom field (such as a
ticket-reference field on prefixes/VLANs) purely from MCP; the definition has to be
created out-of-band via the UI or raw API (`/api/extras/custom-fields/`), and only then
can values be populated through the existing `custom_fields` parameter.

**Suggestion:** add Custom Field tools in the usual read → write/update → delete grouping:

- `get_custom_fields` — list custom fields (optional `q`, `content_type`/object-type
  filter, `tag`, pagination).
- `create_custom_field` — `name`/`label` required, plus `type`, `content_types`
  (which object types the field applies to), `group_name`, `description`, `required`,
  `default`, `filter_logic`, `weight`, `ui_visibility`, `choices` (for select/multi-select).
- `update_custom_field` — partial PATCH merge by numeric ID.
- `delete_custom_field` — destructive, guarded like the other delete tools.

This closes the gap so an agent can define, enumerate, and fill custom fields — including
ticket references on prefixes and VLANs — entirely from MCP.
