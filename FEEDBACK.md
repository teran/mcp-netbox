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
