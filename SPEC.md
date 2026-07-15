# MCP NetBox — Specification

## Overview

An MCP (Model Context Protocol) server for [NetBox](https://netboxlabs.com/).
This server exposes NetBox infrastructure data through the MCP protocol using **Streamable HTTP** transport (remote mode), allowing AI assistants to query DCIM, IPAM, virtualization, tenancy, and circuits data.

## Key Differentiators

- **Remote (HTTP) transport** — uses MCP Streamable HTTP protocol; not stdio-bound.
- **Token from request headers** — the NetBox API token is read from the `Authorization` header of each MCP request, not from an environment variable. This enables per-user authentication in multi-tenant setups.
- **Read-only** — only exposes `GET` operations against NetBox. No create, update, or delete capabilities.
- **Transparent token relay** — the MCP server never inspects or validates the token; it passes it through to NetBox, which handles all authentication and authorization.

## Architecture

```
┌──────────────────┐      MCP (Streamable HTTP)      ┌───────────────────────┐
│  MCP Client      │  ◄──────────────────────────►   │  mcp-netbox           │
│  (AI Assistant)  │     Authorization: Bearer <tok>  │  (Go server)          │
└──────────────────┘                                   └──────┬────────────────
                                                               │ HTTP (Token auth)
                                                               ▼
                                                    ┌───────────────────────┐
                                                    │  NetBox              │
                                                    │  REST API            │
                                                    └───────────────────────┘
```

> **TLS termination** is expected to be handled by a reverse proxy (e.g., nginx, Envoy) placed in front of the MCP server. The server itself does not serve HTTPS directly.

## Technology Stack

| Component         | Choice                                                          |
|-------------------|-----------------------------------------------------------------|
| Language          | Go                                                              |
| MCP SDK           | `github.com/modelcontextprotocol/go-sdk`                        |
| Transport         | Streamable HTTP (MCP spec 2025-03-26+, remote-capable)          |
| HTTP Router       | `net/http` standard library + middleware pattern                |
| Tool Registration | `handlers/registration.go` — `RegisterTools()` function         |
| Metrics           | Prometheus (Go runtime + custom MCP metrics) on port 8081       |

## Configuration (Environment Variables)

| Variable               | Required | Default | Description                          |
|------------------------|----------|---------|--------------------------------------|
| `NETBOX_URL`           | Yes      | —       | Base URL of the NetBox instance (e.g. `http://netbox:8000`). Must be a valid HTTP(S) URL. Loopback, private, and link-local IP addresses are rejected for SSRF protection. |
| `LISTEN_ADDR`          | No       | `:8080` | TCP address to listen on             |
| `PROMETHEUS_METRICS_ADDR` | No    | `:8081` | TCP address for the Prometheus `/metrics` endpoint |
| `RATE_LIMIT_GLOBAL`    | No       | `100`   | Global rate limit (requests/second)  |
| `RATE_LIMIT_PER_CLIENT`| No       | `10`    | Per-client IP rate limit (requests/second) |
| `ALLOW_PRIVATE_NETBOX` | No       | `false` | When `true`, bypasses SSRF protection and allows `NETBOX_URL` to point to private/reserved IP addresses. Only enable if NetBox is on a private network without a public DNS name. |
| `TRUSTED_PROXY`        | No       | `""`    | CIDR prefix for the trusted reverse proxy (e.g. `10.0.0.0/8`). When set, the server uses the first IP from `X-Forwarded-For` for rate limiting instead of `RemoteAddr`. |
| `WRITE_TIMEOUT`        | No       | `300s`  | HTTP write timeout (Go duration format, e.g. `300s`). Minimum 1s. Note: 0 will fail validation; use a reverse proxy for no timeout. |

The NetBox API token is **not** set via environment variables. It is supplied per-request in the `Authorization` header as `Bearer <token>`.

The MCP server listens on the `/mcp` HTTP path via the Streamable HTTP handler.

## MCP Tools

### 1. `get_sites`

List sites in NetBox with optional filters.

**Input**:

| Parameter    | Type   | Required | Description                                  |
|--------------|--------|----------|----------------------------------------------|
| `region`     | string | no       | Filter by region (slug or name)              |
| `status`     | string | no       | Filter by status (active, planned, etc.)     |
| `tenant`     | string | no       | Filter by tenant (slug or name)              |
| `tag`        | string | no       | Filter by tag (slug)                         |
| `page`       | int    | no       | Page number (default: 1)                     |
| `page_size`  | int    | no       | Results per page (default: 25, max: 1000)     |
| `q`          | string | no       | Free-text search across all fields           |

**Output**: Paginated list of sites.

---

### 2. `get_devices`

List devices in NetBox with optional filters.

**Input**:

| Parameter      | Type   | Required | Description                                  |
|----------------|--------|----------|----------------------------------------------|
| `site`         | string | no       | Filter by site (slug)                        |
| `role`         | string | no       | Filter by device role (slug)                 |
| `manufacturer` | string | no       | Filter by manufacturer (slug)                |
| `device_type`  | string | no       | Filter by device type (model slug)           |
| `status`       | string | no       | Filter by status (active, planned, etc.)     |
| `name`         | string | no       | Filter by name (partial match with `__ic`)   |
| `tenant`       | string | no       | Filter by tenant (slug)                      |
| `rack`         | string | no       | Filter by rack (name)                        |
| `cluster`      | string | no       | Filter by cluster (name)                     |
| `tag`          | string | no       | Filter by tag (slug)                         |
| `page`         | int    | no       | Page number (default: 1)                     |
| `page_size`    | int    | no       | Results per page (default: 25, max: 1000)     |
| `q`            | string | no       | Free-text search across all fields           |

**Output**: Paginated list of devices.

---

### 3. `get_ip_addresses`

Search IP addresses in NetBox.

**Input**:

| Parameter      | Type   | Required | Description                                  |
|----------------|--------|----------|----------------------------------------------|
| `address`      | string | no       | Filter by address (e.g. `192.168.1.0/24`)    |
| `device`       | string | no       | Filter by assigned device name               |
| `status`       | string | no       | Filter by status (active, reserved, etc.)    |
| `vrf`          | string | no       | Filter by VRF (rd or name)                   |
| `role`         | string | no       | Filter by role (loopback, etc.)              |
| `tenant`       | string | no       | Filter by tenant (slug)                      |
| `tag`          | string | no       | Filter by tag (slug)                         |
| `page`         | int    | no       | Page number (default: 1)                     |
| `page_size`    | int    | no       | Results per page (default: 25, max: 1000)     |
| `q`            | string | no       | Free-text search across all fields           |

**Output**: Paginated list of IP addresses.

---

### 4. `get_prefixes`

Search IP prefixes in NetBox.

**Input**:

| Parameter    | Type   | Required | Description                                  |
|--------------|--------|----------|----------------------------------------------|
| `prefix`     | string | no       | Filter by prefix (e.g. `10.0.0.0/8`)        |
| `site`       | string | no       | Filter by site (slug)                        |
| `vrf`        | string | no       | Filter by VRF (rd or name)                   |
| `status`     | string | no       | Filter by status (active, container, etc.)   |
| `role`       | string | no       | Filter by role (slug)                        |
| `tenant`     | string | no       | Filter by tenant (slug)                      |
| `within`     | string | no       | Find prefixes within a given prefix          |
| `family`     | int    | no       | Address family: 4 or 6                       |
| `tag`        | string | no       | Filter by tag (slug)                         |
| `page`       | int    | no       | Page number (default: 1)                     |
| `page_size`  | int    | no       | Results per page (default: 25, max: 1000)     |
| `q`          | string | no       | Free-text search across all fields           |

**Output**: Paginated list of prefixes.

---

### 5. `get_vlans`

List VLANs in NetBox.

**Input**:

| Parameter    | Type   | Required | Description                                  |
|--------------|--------|----------|----------------------------------------------|
| `site`       | string | no       | Filter by site (slug)                        |
| `group`      | string | no       | Filter by VLAN group (slug)                  |
| `status`     | string | no       | Filter by status (active, reserved, etc.)    |
| `tenant`     | string | no       | Filter by tenant (slug)                      |
| `vid`        | int    | no       | Filter by VLAN ID                            |
| `tag`        | string | no       | Filter by tag (slug)                         |
| `page`       | int    | no       | Page number (default: 1)                     |
| `page_size`  | int    | no       | Results per page (default: 25, max: 1000)     |
| `q`          | string | no       | Free-text search across all fields           |

**Output**: Paginated list of VLANs.

---

### 6. `get_virtual_machines`

List virtual machines in NetBox.

**Input**:

| Parameter      | Type   | Required | Description                                  |
|----------------|--------|----------|----------------------------------------------|
| `cluster`      | string | no       | Filter by cluster (name)                     |
| `cluster_group`| string | no       | Filter by cluster group (slug)               |
| `role`         | string | no       | Filter by VM role (slug)                     |
| `status`       | string | no       | Filter by status (active, staged, etc.)      |
| `tenant`       | string | no       | Filter by tenant (slug)                      |
| `name`         | string | no       | Filter by name (partial match with `__ic`)   |
| `site`         | string | no       | Filter by site (slug)                        |
| `tag`          | string | no       | Filter by tag (slug)                         |
| `page`         | int    | no       | Page number (default: 1)                     |
| `page_size`    | int    | no       | Results per page (default: 25, max: 1000)     |
| `q`            | string | no       | Free-text search across all fields           |

**Output**: Paginated list of virtual machines.

---

### 7. `get_clusters`

List clusters in NetBox.

**Input**:

| Parameter      | Type   | Required | Description                                  |
|----------------|--------|----------|----------------------------------------------|
| `cluster_type` | string | no       | Filter by cluster type (slug)                |
| `cluster_group`| string | no       | Filter by cluster group (slug)               |
| `site`         | string | no       | Filter by site (slug)                        |
| `tenant`       | string | no       | Filter by tenant (slug)                      |
| `name`         | string | no       | Filter by name (partial match with `__ic`)   |
| `tag`          | string | no       | Filter by tag (slug)                         |
| `page`         | int    | no       | Page number (default: 1)                     |
| `page_size`    | int    | no       | Results per page (default: 25, max: 1000)     |
| `q`            | string | no       | Free-text search across all fields           |

**Output**: Paginated list of clusters.

---

### 8. `get_circuits`

List circuits in NetBox.

**Input**:

| Parameter      | Type   | Required | Description                                  |
|----------------|--------|----------|----------------------------------------------|
| `provider`     | string | no       | Filter by provider (slug)                    |
| `circuit_type` | string | no       | Filter by circuit type (slug)                |
| `site`         | string | no       | Filter by site (slug)                        |
| `status`       | string | no       | Filter by status (active, planned, etc.)     |
| `tenant`       | string | no       | Filter by tenant (slug)                      |
| `tag`          | string | no       | Filter by tag (slug)                         |
| `page`         | int    | no       | Page number (default: 1)                     |
| `page_size`    | int    | no       | Results per page (default: 25, max: 1000)     |
| `q`            | string | no       | Free-text search across all fields           |

**Output**: Paginated list of circuits.

---

### 9. `get_object_by_id`

Retrieve any NetBox object by its type and numeric ID.

**Input**:

| Parameter    | Type   | Required | Description                                    |
|--------------|--------|----------|------------------------------------------------|
| `object_type`| string | yes      | Object type. One of: `site`, `device`, `prefix`, `ip_address`, `vlan`, `virtual_machine`, `cluster`, `circuit`, `provider`, `tenant`, `rack`, `manufacturer`, `device_type`, `location`, `cluster_type`, `cluster_group`, `circuit_type`, `vrf`, `vlan_group`, `role`, `contact`, `cable`, `interface`, `vm_interface`, `circuit_termination` |
| `id`         | int    | yes      | Numeric ID of the object                       |
| `params`     | map    | no       | Additional query parameters to pass to NetBox (optional) |

**Output**: Full object detail.

---

### 10. `get_racks`

List racks in NetBox.

**Input**:

| Parameter    | Type   | Required | Description                                  |
|--------------|--------|----------|----------------------------------------------|
| `site`       | string | no       | Filter by site (slug)                        |
| `location`   | string | no       | Filter by location (slug)                    |
| `status`     | string | no       | Filter by status (active, planned, etc.)     |
| `tenant`     | string | no       | Filter by tenant (slug)                      |
| `tag`        | string | no       | Filter by tag (slug)                         |
| `page`       | int    | no       | Page number (default: 1)                     |
| `page_size`  | int    | no       | Results per page (default: 25, max: 1000)     |
| `q`          | string | no       | Free-text search across all fields           |

**Output**: Paginated list of racks.

---

### 11. `get_interfaces`

List device interfaces in NetBox. Includes physical interfaces of devices (switches, routers, servers).

**Input**:

| Parameter    | Type   | Required | Description                                  |
|--------------|--------|----------|----------------------------------------------|
| `device`     | string | no       | Filter by device (name)                      |
| `type`       | string | no       | Filter by interface type (slug)              |
| `enabled`    | bool   | no       | Filter by enabled status (true/false)        |
| `name`       | string | no       | Filter by name (case-insensitive partial match) |
| `tag`        | string | no       | Filter by tag (slug)                         |
| `page`       | int    | no       | Page number (default: 1)                     |
| `page_size`  | int    | no       | Results per page (default: 25, max: 1000)     |
| `q`          | string | no       | Free-text search across all fields           |

**Output**: Paginated list of interfaces.

---

### 12. `get_vm_interfaces`

List VM interfaces in NetBox. These are virtual NICs attached to virtual machines.

**Input**:

| Parameter        | Type   | Required | Description                                  |
|------------------|--------|----------|----------------------------------------------|
| `virtual_machine`| string | no       | Filter by virtual machine (name)             |
| `name`           | string | no       | Filter by name (case-insensitive partial match) |
| `tag`            | string | no       | Filter by tag (slug)                         |
| `page`           | int    | no       | Page number (default: 1)                     |
| `page_size`      | int    | no       | Results per page (default: 25, max: 1000)     |
| `q`              | string | no       | Free-text search across all fields           |

**Output**: Paginated list of VM interfaces.

---

### 13. `get_circuit_terminations`

List circuit terminations in NetBox. A termination represents one end of a circuit at a site.

**Input**:

| Parameter    | Type   | Required | Description                                  |
|--------------|--------|----------|----------------------------------------------|
| `circuit`    | string | no       | Filter by circuit (ID or CID)                |
| `site`       | string | no       | Filter by site (slug)                        |
| `term_side`  | string | no       | Filter by termination side (`A` or `Z`)      |
| `tag`        | string | no       | Filter by tag (slug)                         |
| `page`       | int    | no       | Page number (default: 1)                     |
| `page_size`  | int    | no       | Results per page (default: 25, max: 1000)     |
| `q`          | string | no       | Free-text search across all fields           |

**Output**: Paginated list of circuit terminations.

---

### 14. `get_cables`

List cables in NetBox with optional filters.

**Input**:

| Parameter    | Type   | Required | Description                                  |
|--------------|--------|----------|----------------------------------------------|
| `type`       | string | no       | Filter by cable type (slug)                  |
| `status`     | string | no       | Filter by status (connected, planned, decommissioning) |
| `site`       | string | no       | Filter by site (slug)                        |
| `color`      | string | no       | Filter by color (slug)                       |
| `label`      | string | no       | Filter by label (case-insensitive partial match) |
| `tag`        | string | no       | Filter by tag (slug)                         |
| `page`       | int    | no       | Page number (default: 1)                     |
| `page_size`  | int    | no       | Results per page (default: 25, max: 1000)     |
| `q`          | string | no       | Free-text search across all fields           |

**Output**: Paginated list of cables.

## Middleware Chain

The server applies eight middleware layers to every HTTP request, executed in this order (outermost first):

1. **RecoveryMiddleware** — catches panics, returns 500
2. **SecurityHeadersMiddleware** — sets security headers (X-Content-Type-Options, X-Frame-Options, Referrer-Policy)
3. **HostValidationMiddleware** — rejects requests with empty or malformed `Host` headers
4. **MetricsMiddleware** — tracks in-flight requests via gauge
5. **RateLimitMiddleware** — global (100 rps) + per-client (10 rps) token bucket
6. **BodyLimitMiddleware** — 1 MB request body limit
7. **LoggingMiddleware** — logs MCP method, duration, status, sizes (never logs token)
8. **TokenMiddleware** — extracts token from `Authorization` header, stores as `*application.Token` in context (safe redaction via String/GoString/MarshalJSON)
9. **injectClientMiddleware** — creates NetBox API client with per-request token, stores service in context

## Authentication Flow

1. MCP Client sends POST to Streamable HTTP endpoint with `Authorization: Bearer <token>`
2. Middleware chain extracts, validates format, and relays the token
3. Token is stored in request context, passed to NetBox API client
4. Tool handlers retrieve services from context and call NetBox API
5. Token is never stored on server — exists only for request lifetime

## Error Handling

### HTTP Level (Middleware)

| Status | Cause |
|--------|-------|
| 429 Too Many Requests | Rate limit exceeded |
| 401 Unauthorized | Missing or malformed Authorization header |
| 403 Forbidden | Token lacks permissions |
| 413 Request Entity Too Large | Body exceeds 1 MB limit |
| 500 Internal Server Error | Panic recovery |
| 400 Bad Request | Batch request exceeds max size (100) |

### MCP Level (Tool Handlers)

| Scenario | MCP Error |
|----------|-----------|
| Resource not found | `isError: true` |
| NetBox unavailable | `isError: true` |
| Invalid token | `isError: true` |

## Health Check

`GET /healthz` — Returns `{"status":"ok"}` with HTTP 200. Used for liveness probes.

## Custom Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `mcp_tool_requests_total` | Counter | `{tool, status_class}` | Per-tool request count |
| `mcp_tool_duration_seconds` | Histogram | `{tool}` | Per-tool request duration |
| `mcp_active_requests` | Gauge | — | Current in-flight requests |

## Development

### Prerequisites
- Go 1.26+
- golangci-lint
- goreleaser

### Building
```bash
goreleaser build --snapshot --clean
```

### Testing
```bash
go test -race -coverprofile=coverage.out -count=1 ./...
```

### Linting
```bash
golangci-lint run ./...
```
