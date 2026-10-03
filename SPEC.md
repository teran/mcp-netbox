# MCP NetBox — Specification

## Overview

An MCP (Model Context Protocol) server for [NetBox](https://netboxlabs.com/).
This server exposes NetBox infrastructure data through the MCP protocol using **Streamable HTTP** transport (remote mode), allowing AI assistants to query and manage DCIM, IPAM, virtualization, tenancy, and circuits data.

## Repository Location (R5)

The repository is hosted **publicly on GitHub.com** at
`github.com/teran/mcp-netbox` (module path `github.com/teran/mcp-netbox`,
container image `ghcr.io/teran/mcp-netbox`). Because the repository is **public**
(not an internal Forgejo host), the public module path is used and the S6
local-only-naming restriction does **not** apply. The default branch is `master`
(R2). This location was fixed before development and drives the module path, the
README badge URLs, the CI workflow triggers, and the image/release tagging scheme.

## Key Differentiators

- **Hybrid transport** — serves MCP over **Streamable HTTP** (remote) or **STDIO** (local), selected at startup via the `-mode http|stdio` flag (default `stdio`); the `TRANSPORT` env var remains a backward-compatible override.
- **Token from request headers (HTTP)** — the NetBox API token is read from the `Authorization` header of each MCP request, not from an environment variable. This enables per-user authentication in multi-tenant setups.
- **Token from environment (STDIO)** — in STDIO mode the NetBox token is read once from `NETBOX_TOKEN` at startup (there is no HTTP header to carry it per request).
- **Full CRUD** — 14 read-only `get_*` tools plus typed `create_*`/`update_*`/`delete_*` tools for all 25 entities (89 tools total). `create_*`/`update_*` modify NetBox (POST / partial PATCH); `delete_*` is **irreversible** and flagged `DestructiveHint: true`.
- **Transparent token relay** — the MCP server never inspects or validates the token; it passes it through to NetBox, which handles all authentication and authorization.

## Transport Decision

MCP servers choose between **STDIO** and **HTTP/SSE** transports based on the deployment task. `mcp-netbox` is a **Hybrid** server because the same tool surface serves two distinct use cases:

- **HTTP/SSE (remote)** — the default. Streamable HTTP lets many AI assistants / MCP clients connect to a shared NetBox instance over the network, each authenticating with its **own** NetBox token via the `Authorization: Bearer` header (per-request, per-user auth in multi-tenant setups). This is the deployment shape that also requires a container image.
- **STDIO (local)** — runs as a local child process of an MCP client (e.g. a desktop assistant or CLI). The client launches the server with `TRANSPORT=stdio` and passes the NetBox token via the `NETBOX_TOKEN` environment variable. Because there is no HTTP layer, per-request headers do not exist; a single shared token is used for the process lifetime.

The choice is driven by the task: NetBox is an infrastructure source of truth that is frequently shared and queried by many assistants, which favors the remote HTTP transport; but a local, single-tenant, zero-network-footprint mode (STDIO) is valuable for development and desktop use. Supporting both with one codebase (Hybrid) covers both without duplicating the tool surface. The transport (launch mode) is selected at startup via the `-mode http|stdio` flag (default `stdio`); the `TRANSPORT` environment variable is honoured as a backward-compatible override (`-mode` wins).

## Auth Decision (OAuth2)

**OAuth2 is NOT used.**

The MCP spec's OAuth2 flow is designed for servers that own their own resource
server and need to mint access tokens for end users. That does not apply here:
`mcp-netbox` is a **stateless relay** in front of NetBox, which is itself the
resource server and already implements its own authentication and authorization
model (NetBox personal access tokens).

Instead of OAuth2, the server relays an existing NetBox personal token to the
NetBox REST API:

- **HTTP transport**: the client supplies the token per request in the
  `Authorization: Bearer <token>` header. The server extracts it, passes it
  through to NetBox, and never validates or stores it.
- **STDIO transport**: there is no HTTP header, so the token is provided once via
  the `NETBOX_TOKEN` environment variable at startup and used by a single shared
  NetBox client for the process lifetime.

Authorization (which objects a token may read or modify, and which write
operations it may perform) is delegated entirely to NetBox. This server only
relays the credential and forwards the requested operation; it does not
implement an authorization decision layer. The server likewise does not
validate write payloads — NetBox is the source of truth for both authorization
and data validation (a NetBox 400 is surfaced as a `ValidationError`). This
keeps the server simple (no OAuth2 machinery, token endpoints, or consent flow)
and matches the existing NetBox security model.

## Deployment Type: Hybrid

The server is classified as **Hybrid**: it can run as a **Remote (HTTP/SSE)**
server and as a **Local (STDIO)** server. All deployment attributes follow from
the transport it is launched with:

| Aspect            | Remote (HTTP/SSE)                          | Local (STDIO)                         |
|-------------------|---------------------------------------------|----------------------------------------|
| **Auth**          | `Authorization: Bearer` header, per request | `NETBOX_TOKEN` env var, at startup     |
| **Logging**       | stdout (12-factor)                          | log file (default `/tmp/mcp-netbox.log`, chmod 600) |
| **Image**         | required and published by CI/CD             | not used                               |

Because the HTTP transport is supported, a **container image is built and
published** by CI/CD on every commit to `master` and on every release tag (see
[Release & Container Images](#release--container-images)). Logging follows the
transport (see [Logging](#logging)).

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

> **TLS termination** is expected to be handled by a reverse proxy (e.g., nginx, Envoy) placed in front of the MCP server. The server itself does not serve HTTPS directly. This also applies to the HTTP/SSE transport: TLS is never implemented inside the server (N1).

### Package Layout (Clean / DDD)

The application follows Clean architecture with a strict dependency direction —
inner layers depend on nothing internal, outer layers may depend inward:

| Layer             | Package(s)                              | Responsibility                                      |
|-------------------|-----------------------------------------|-----------------------------------------------------|
| Composition root  | `cmd/server`                            | Entrypoint; builds config, wiring, HTTP/STDIO transport |
| Infrastructure    | `infrastructure/netbox`, `infrastructure/circuitbreaker` | Adapters: NetBox REST client, circuit breaker transport |
| Application       | `application`                           | Use cases / `NetworkService` business logic          |
| Domain            | `domain`                                | Domain models + repository interfaces (ports)        |
| Config            | `config`                                | Env loading + validation (envconfig + ozzo-validation) |
| Logging           | `logging`                          | logrus setup, channel-by-transport (L1–L4), SDK `slog`→logrus handler |
| HTTP transport    | `handlers`                              | MCP tool handlers, middleware chain, metrics, registration |

> **Layout note (A1):** this repository deliberately uses the Clean/DDD layered
> layout above (composition root + layered packages) rather than the skill's
> default simple single-package layout. This is an explicit, documented choice
> for this project; dependency boundaries are enforced by `.go-arch-lint.yml`
> (C6).

### Tool Registry

All MCP tools are registered in **one place**: `handlers/registration.go`,
via the `RegisterTools()` function. Each tool declaration pairs a `mcp.Tool`
(name, description, annotations, instructions) with a wrapped handler
(`WrapToolHandler`) that decodes the input, calls the appropriate
`NetworkService` method, and formats the paginated output. In HTTP mode the
service is injected per request by `injectClientMiddleware`; in STDIO mode a
single shared service is passed directly to `RegisterTools`.

### Dependency Boundaries

Dependency rules are authored in `.go-arch-lint.yml` and enforced in CI by
`go-arch-lint check` (C6). The direction is:

```
domain  ←  application  ←  infrastructure, config, logging  ←  handlers  ←  cmd/server
```

`domain` never imports `application`, `handlers`, `infrastructure`, or
`config`. Third-party (vendor) dependencies may be used by any component.
Test files are excluded from architecture analysis.

### Transport Wiring

`cmd/server` reads the launch mode (`-mode` flag, else the `TRANSPORT` env var) and dispatches at startup:

- `http` → `runHTTP`: builds the Streamable HTTP handler (`/mcp`), the full
  HTTP middleware chain, and a separate internal observability server
  (metrics, pprof, probes); tokens are injected per request.
- `stdio` → `runStdio`: builds a single shared `NetworkService` from
  `NETBOX_TOKEN`, registers it directly, and serves over stdin/stdout via the
  SDK's `StdioTransport` (newline-delimited JSON).

The MCP `Server` is shared; only the transport wiring differs.

## Technology Stack

| Component         | Choice                                                          |
|-------------------|-----------------------------------------------------------------|
| Language          | Go 1.27+                                                        |
| MCP SDK           | `github.com/modelcontextprotocol/go-sdk` v1.8.0                 |
| Transport         | Hybrid — Streamable HTTP (MCP spec 2025-03-26+) and STDIO, selected via `-mode http|stdio` (default `stdio`; `TRANSPORT` override) |
| HTTP Router       | `net/http` standard library + middleware pattern                |
| Tool Registration | `handlers/registration.go` — `RegisterTools()` function         |
| Logging           | `github.com/sirupsen/logrus` — channel by transport, gated by `LOG_LEVEL`; SDK `slog` wired into logrus |
| Outbound HTTP     | `resty.dev/v3` — NetBox client (DNS-rebinding dialer + circuit breaker transport preserved) |
| Metrics           | Prometheus (Go runtime + custom MCP metrics), pprof, and health/readiness/startup probes on the internal address port 8081 |

## Configuration (Environment Variables)

| Variable               | Required | Default | Description                          |
|------------------------|----------|---------|--------------------------------------|
| `NETBOX_URL`           | Yes      | —       | Base URL of the NetBox instance (e.g. `http://netbox:8000`). Must be a valid HTTP(S) URL. Loopback, private, and link-local IP addresses are rejected for SSRF protection. |
| `NETBOX_TOKEN`         | STDIO only | `""` | NetBox API token for **STDIO** transport. Required when `TRANSPORT=stdio`. Ignored for HTTP. |
| `TRANSPORT`            | No       | `stdio` | MCP transport override: `http` (Streamable HTTP, remote) or `stdio` (stdin/stdout, local). Overridden by the `-mode` flag; defaults to `stdio`. |
| `LISTEN_ADDR`          | No       | `:8080` | TCP address for the MCP server (HTTP transport) |
| `INTERNAL_ADDR`        | No       | `:8081` | Internal observability address (HTTP transport): Prometheus `/metrics`, pprof `/debug/pprof/*`, and `healthz`/`readyz`/`startup` probes |
| `RATE_LIMIT_GLOBAL`    | No       | `100`   | Global rate limit (requests/second)  |
| `RATE_LIMIT_PER_CLIENT`| No       | `10`    | Per-client IP rate limit (requests/second) |
| `ALLOW_PRIVATE_NETBOX` | No       | `false` | When `true`, bypasses SSRF protection and allows `NETBOX_URL` to point to private/reserved IP addresses. Only enable if NetBox is on a private network without a public DNS name. |
| `TRUSTED_PROXY`        | No       | `""`    | CIDR prefix for the trusted reverse proxy (e.g. `10.0.0.0/8`). When set, the server uses the first IP from `X-Forwarded-For` for rate limiting instead of `RemoteAddr`. |
| `WRITE_TIMEOUT`        | No       | `300s`  | HTTP write timeout (Go duration format, e.g. `300s`). Minimum 1s. Note: 0 will fail validation; use a reverse proxy for no timeout. |
| `LOG_LEVEL`            | No       | (unset) | Logrus level (`trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`). **Unset ⇒ logging disabled.** |
| `LOG_FORMAT`           | No       | `text`  | Log format: `text` (logrus text, full absolute timestamp) or `json`. |
| `LOG_FILENAME`         | No       | `/tmp/mcp-netbox.log` | Log file path for **STDIO** transport (chmod 600). Ignored for HTTP (logs go to stdout). |

The NetBox API token is supplied per-request in the `Authorization` header as
`Bearer <token>` in HTTP mode. In STDIO mode it is read from `NETBOX_TOKEN` at
startup (there is no HTTP header to carry it per request).

The MCP server listens on the `/mcp` HTTP path via the Streamable HTTP handler
when the launch mode is `http` (`-mode http` or `TRANSPORT=http`).

## Logging

Logging uses **logrus**. Enablement follows the launch mode (L2):

- **HTTP** mode — logging is **always enabled**, defaulting to level **`info`**
  (overridable via `LOG_LEVEL`).
- **STDIO** mode — logging is **enabled only when `LOG_LEVEL` is set** — unset
  ⇒ **disabled**.

- **Channel by transport (L1):**
  - **HTTP** → **stdout** (12-factor style).
  - **STDIO** → a log file, because stdout carries the MCP protocol itself and
    must not be polluted with log lines. Default path `/tmp/mcp-netbox.log`,
    created with mode `0600`.
- **Path override (L3):** `LOG_FILENAME` changes the STDIO log file path.
- **Format (L4):** `LOG_FORMAT` — default `text` (logrus text, full absolute
  timestamp) or `json`.
- **Secrets (L5):** tokens, passwords and credentials are **never** logged. The
  token is redacted (see [Security](#security--secrets)) and URLs are logged via
  `url.Redacted()`.
- **Startup banner (L6/B5):** when logging is enabled (always in HTTP mode; in
  STDIO only when `LOG_LEVEL` is set), the very first log line at startup is
  the banner
  `Starting {appName}/{appVersion} (commit: {appCommitHash}; built at {appTimestamp}) ...`,
  built from ldflags-embedded metadata. No banner is emitted when logging is
  disabled.
- **SDK logger (L7/G10):** the MCP go-sdk's internal `slog` logger is wired into
  the server's logrus logger via `ServerOptions.Logger` (an `slog.Handler`
  forwarding to logrus), so SDK-level MCP events — session connect/end,
  tool-call results and errors, protocol warnings — are visible in the server
  logs at the configured level.
- **Per-request trace logging (L8):** when logging is enabled, each MCP tool
  call emits an access-log line at info level (not gated behind debug — it is
  visible whenever logging is on, L2) with structured fields: `tool`,
  `args` (only non-sensitive request params, never the Authorization token),
  `source` (derived for HTTP from `X-Real-IP` → `X-Forwarded-For` → `RemoteAddr`,
  comma-joined so all proxy hops are visible; `"STDIO"` for stdio), `duration`,
  `in_bytes`, `out_bytes`, and `outcome` (success/error).
- **Correlation via request_id (L9/G11):** every incoming request is assigned a
  unique `request_id` (reusing an inbound `X-Request-ID` when present), injected
  into the request context via `domain.WithRequestID`. A context-aware entry
  builder (`WithSession`) decorates all request-scoped log lines with
  `request_id`/`session_id`. The outbound NetBox client forwards the id as the
  `X-Request-ID` header and emits its own per-request log record (method, path,
  `in_bytes`, `out_bytes`, status, duration) tagged with the same `request_id`, so
  client- and server-side records can be matched across the wire.

## MCP Tools

Every tool is registered with **Annotations** and per-tool **Instructions** metadata.

The server exposes **89 tools**: 14 read-only queries and, for each of the 25
entities, a `create_*`, an `update_*`, and a `delete_*` (75 write tools).

The server-wide instructions (`ServerOptions.Instructions`) tell clients that
all **list** tools are **paginated** (use `page` 1-based and `page_size` max
1000, an empty array is not an error), that filters are **additive**, and that
authentication is enforced entirely by NetBox (Bearer in HTTP, `NETBOX_TOKEN`
env in STDIO).

The tools split into three annotation groups:

- **Read tools** (14): `readOnlyHint: true`, `idempotentHint: true`, `destructiveHint: false`, `openWorldHint: false` — they never mutate state.
- **Write tools — create/update** (50): `readOnlyHint: false`, `idempotentHint: false` (a create produces a new object each call), `destructiveHint: false`, `openWorldHint: false`.
- **Write tools — delete** (25): `readOnlyHint: false`, `idempotentHint: true` (deleting a non-existent object is a no-op), `destructiveHint: true` (irreversible), `openWorldHint: false`.

The per-tool `title` and `instructions` are listed under each read tool below.

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

**Annotations**: `title` — "List Sites"; `readOnlyHint: true`; `destructiveHint: false`; `idempotentHint: true`; `openWorldHint: false`.

**Instructions**: Use `q` for free-text search, or filter by `region`, `status`, `tenant`, or `tag`. Results are paginated (`page`, `page_size` max 1000); an empty array means no sites match, not an error.

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

**Annotations**: `title` — "List Devices"; `readOnlyHint: true`; `destructiveHint: false`; `idempotentHint: true`; `openWorldHint: false`.

**Instructions**: Filter by `site`, `role`, `manufacturer`, `device_type`, `status`, `name`, `tenant`, `rack`, `cluster`, or `tag`, or use `q` for free-text search. Paginated; an empty array means no devices match.

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

**Annotations**: `title` — "List IP Addresses"; `readOnlyHint: true`; `destructiveHint: false`; `idempotentHint: true`; `openWorldHint: false`.

**Instructions**: Filter by `address` (e.g. `192.168.1.0/24`), assigned `device`, `status`, `vrf`, `role`, `tenant`, or `tag`, or use `q`. Paginated; an empty array means no addresses match.

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

**Annotations**: `title` — "List Prefixes"; `readOnlyHint: true`; `destructiveHint: false`; `idempotentHint: true`; `openWorldHint: false`.

**Instructions**: Filter by `prefix` (e.g. `10.0.0.0/8`), `site`, `vrf`, `status`, `role`, `tenant`, `family`, or `tag`, or use `within` to find nested prefixes, or `q`. Paginated; an empty array means no prefixes match.

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

**Annotations**: `title` — "List VLANs"; `readOnlyHint: true`; `destructiveHint: false`; `idempotentHint: true`; `openWorldHint: false`.

**Instructions**: Filter by `site`, `group`, `status`, `tenant`, `tag`, or a specific VLAN ID (`vid`), or use `q`. Paginated; an empty array means no VLANs match.

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

**Annotations**: `title` — "List Virtual Machines"; `readOnlyHint: true`; `destructiveHint: false`; `idempotentHint: true`; `openWorldHint: false`.

**Instructions**: Filter by `cluster`, `cluster_group`, `role`, `status`, `tenant`, `name`, or `site`, or use `q`. Paginated; an empty array means no VMs match.

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

**Annotations**: `title` — "List Clusters"; `readOnlyHint: true`; `destructiveHint: false`; `idempotentHint: true`; `openWorldHint: false`.

**Instructions**: Filter by `cluster_type`, `cluster_group`, `site`, `tenant`, or `name`, or use `q`. Paginated; an empty array means no clusters match.

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

**Annotations**: `title` — "List Circuits"; `readOnlyHint: true`; `destructiveHint: false`; `idempotentHint: true`; `openWorldHint: false`.

**Instructions**: Filter by `provider`, `circuit_type`, `site`, `status`, or `tenant`, or use `q`. Paginated; an empty array means no circuits match.

---

### 9. `get_object_by_id`

Retrieve any NetBox object by its type and numeric ID.

**Input**:

| Parameter    | Type   | Required | Description                                    |
|--------------|--------|----------|------------------------------------------------|
| `object_type`| string | yes      | Object type. One of: `site`, `device`, `prefix`, `ip_address`, `vlan`, `virtual_machine`, `cluster`, `circuit`, `provider`, `tenant`, `rack`, `manufacturer`, `device_type`, `location`, `cluster_type`, `cluster_group`, `circuit_type`, `vrf`, `vlan_group`, `role`, `contact`, `cable`, `interface`, `vm_interface`, `circuit_termination` |
| `id`         | int    | yes      | Numeric ID of the object                       |
| `params`     | map    | no       | Additional query parameters to pass to NetBox (optional) |

**Output**: Full object detail. The `data` field is of arbitrary type (`any`) because it mirrors the raw NetBox object of the requested type.

**Annotations**: `title` — "Get Object by ID"; `readOnlyHint: true`; `destructiveHint: false`; `idempotentHint: true`; `openWorldHint: false`.

**Instructions**: `object_type` is required (e.g. `site`, `device`, `prefix`, `ip_address`, `vlan`, `virtual_machine`, `cluster`, `circuit`, `provider`, `tenant`, `rack`, `cable`, `interface`). Returns an error if the object does not exist or the token lacks permission.

**Rationale (M07/N29):** this is the **single** generic read passthrough in the
tool set, and it exists deliberately despite the "no generic CRUD/passthrough
mega-tool" rule:
- It maps to a real NetBox need: fetching an arbitrary object by numeric ID
  across the generic `GET /api/*/<type>/<id>/` endpoint, where the object shape
  is only known at runtime.
- It is **read-only** (`readOnlyHint: true`), not a write or CRUD mega-tool, so
  it does not weaken the write path.
- All 75 write tools (`create_*`/`update_*`/`delete_*`) remain specific and
  strongly typed; this generic tool never performs mutations.

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

**Annotations**: `title` — "List Racks"; `readOnlyHint: true`; `destructiveHint: false`; `idempotentHint: true`; `openWorldHint: false`.

**Instructions**: Filter by `site`, `location`, `status`, or `tenant`, or use `q`. Paginated; an empty array means no racks match.

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

**Annotations**: `title` — "List Interfaces"; `readOnlyHint: true`; `destructiveHint: false`; `idempotentHint: true`; `openWorldHint: false`.

**Instructions**: Filter by `device`, `type`, `enabled`, or `name`, or use `q`. Paginated; an empty array means no interfaces match.

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

**Annotations**: `title` — "List VM Interfaces"; `readOnlyHint: true`; `destructiveHint: false`; `idempotentHint: true`; `openWorldHint: false`.

**Instructions**: Filter by `virtual_machine`, `name`, or `tag`, or use `q`. Paginated; an empty array means no VM interfaces match.

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

**Annotations**: `title` — "List Circuit Terminations"; `readOnlyHint: true`; `destructiveHint: false`; `idempotentHint: true`; `openWorldHint: false`.

**Instructions**: Filter by `circuit`, `site`, `term_side` (A or Z), or `tag`, or use `q`. Paginated; an empty array means no terminations match.

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

**Annotations**: `title` — "List Cables"; `readOnlyHint: true`; `destructiveHint: false`; `idempotentHint: true`; `openWorldHint: false`.

**Instructions**: Filter by `type`, `status`, `site`, `color`, or `label`, or use `q`. Paginated; an empty array means no cables match.

## Middleware Chain (HTTP transport)

The server applies ten middleware layers to every HTTP request, executed in this order (outermost first):

1. **RequestIDMiddleware** — injects a per-request `request_id` into the context (reusing an inbound `X-Request-ID` when present) for correlation (L9/G11)
2. **RecoveryMiddleware** — catches panics, returns 500
3. **SecurityHeadersMiddleware** — sets security headers (X-Content-Type-Options, X-Frame-Options, Referrer-Policy)
4. **HostValidationMiddleware** — rejects requests with empty or malformed `Host` headers
5. **MetricsMiddleware** — tracks in-flight requests via gauge
6. **RateLimitMiddleware** — global (100 rps) + per-client (10 rps) token bucket
7. **BodyLimitMiddleware** — 1 MB request body limit
8. **LoggingMiddleware** — logs MCP tool, source, duration, byte sizes, outcome, and `request_id` at info level (never logs token) (L8)
9. **TokenMiddleware** — extracts token from `Authorization` header, stores as `*application.Token` in context (safe redaction via String/GoString/MarshalJSON)
10. **injectClientMiddleware** — creates NetBox API client with per-request token, stores service in context

The middleware chain applies only to the **HTTP** transport. In **STDIO** mode there is no HTTP layer, so the middleware is skipped entirely: a single shared `NetworkService` is built once with the `NETBOX_TOKEN` from the environment and registered directly on the MCP server.

## Authentication Flow

### HTTP transport

1. MCP Client sends POST to the Streamable HTTP endpoint with `Authorization: Bearer <token>`
2. Middleware chain extracts, validates format, and relays the token
3. Token is stored in request context, passed to the NetBox API client
4. Tool handlers retrieve the per-request service from context and call the NetBox API
5. Token is never stored on server — exists only for request lifetime

### STDIO transport

1. The server reads the token once from `NETBOX_TOKEN` at startup (required; validation fails if unset)
2. A single shared `NetworkService` is created with that token
3. Every tool call uses the shared service; no per-request token handling
4. The token is never logged

## Security

- **TLS (S1/N1):** TLS is **never** implemented inside the server. For the
  HTTP/SSE transport it is always the reverse proxy's job (nginx, Caddy,
  ingress). The server serves plain HTTP on its listener.
- **Secret hygiene (S2/N2):** tools never leak tokens, passwords, or secrets in
  outputs where avoidable. The NetBox token is never logged, and its type
  implements safe redaction (`String`/`GoString`/`MarshalJSON`). URLs are logged
  via `url.Redacted()` to strip any embedded credentials.
- **Annotation-based redaction (S02):** secrets are masked with a single
  central redaction helper driven by a `secret:"true"` struct-tag annotation,
  backed by structural safeguards:
  - `redact` performs reflection-based deep redaction of arbitrary
    values: any struct field tagged `secret:"true"` is replaced with
    `***redacted***` in a deep copy. `Redact` never mutates its input and
    handles nested structs, pointers, slices/arrays, maps, and interfaces;
    unexported fields are never rewritten. It exposes `MarshalJSON`/`String`
    and a logrus hook (`NewLogrusHook`) that is wired into the logger in
    `logging`.
  - The tool output path (`handlers.WrapToolHandler`) redacts the successful
    typed output before it is rendered to the model, so annotated secrets never
    appear in tool output.
  - The NetBox auth token is never written to logs. The `Authorization` header
    is excluded from the outbound request log, and the token type implements
    safe `String`/`GoString`/`MarshalJSON` redaction so it cannot leak even if
    logged indirectly. `config.Config.NetBoxToken` carries the `secret:"true"`
    tag so it is redacted if the whole config is ever logged.
  - `domain.ValidationError.Error()` carries **only the HTTP status code** and
    never the response body. The NetBox validation body is delivered to the
    model exclusively via `structuredContent` (the `validation_errors` field),
    so it never flows through the error string or any log line.
  - Consequence: no secrets, tokens, or passwords appear in logs or in tool
    output. Any secret-like value is either annotated (`secret:"true"`) or
    excluded at the source (never logged, never placed in the error string)
    rather than being scrubbed after the fact.
- **CRUD ordering (S3):** the read tools are grouped as **read** operations,
  and the write tools as typed create/update/delete operations. Read tools carry
  `readOnlyHint: true` and `destructiveHint: false`; create/update carry
  `destructiveHint: false`; delete tools carry `destructiveHint: true` because
  deletion is **irreversible**. Clients should treat delete tools with explicit
  confirmation.
- **Security-scan findings (S5/N8):** findings from security scanners
  (**gosec**, **govulncheck**) are **fixed rather than suppressed**. There are no
  blanket suppressions or default excludes; any `//nolint` is narrowly scoped
  and justified (best effort).
- **SSRF / DNS-rebinding protection:** `NETBOX_URL` is validated to reject
  loopback, private, and link-local IPs, and the HTTP client's dialer rejects
  connections to reserved IPs at dial time to prevent DNS rebinding
  (`ALLOW_PRIVATE_NETBOX` overrides for private networks).
- **Rate limiting & body limit:** global/per-client rate limits and a 1 MB body
  limit guard the HTTP listener against abuse.

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
| NetBox 400 validation on write (create/update/delete) | `ValidationError` — returned as structured content (field-level errors from the NetBox response body); the `Error()` string carries only the status code and never the body |

## Health Check

The probes below are served on the internal observability address
(`INTERNAL_ADDR`, default `:8081`), separate from the MCP address (`:8080`).

- `GET /healthz` — Returns `{"status":"ok"}` with HTTP 200. Used for liveness probes.
- `GET /readyz` — Returns `{"status":"ok"}` with HTTP 200, or `{"status":"degraded","circuit_breaker":"open"}` with HTTP 503 when the NetBox circuit breaker is open. Used for readiness probes.
- `GET /startup` — Returns `{"status":"ok"}` with HTTP 200. Used for startup probes.
- `GET /debug/pprof/*` — Standard `net/http/pprof` profiling endpoints (cmdline, profile, symbol, trace, and the named profiles). Registered only on the internal mux.

## Custom Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `mcp_tool_requests_total` | Counter | `{tool, status_class}` | Per-tool request count |
| `mcp_tool_duration_seconds` | Histogram | `{tool}` | Per-tool request duration |
| `mcp_active_requests` | Gauge | — | Current in-flight requests |

## Development Process / TDD

Changes follow a **test-first (TDD)** workflow:

- **`@qa` writes the tests first** using an **isolated context** (no access to
  the shared implementation or production state).
- **`@developer` writes the implementation** using an **isolated context**.

This keeps tests independent from the implementation and prevents a shared
context from biasing either role. The same workflow is documented for
contributors (see `README.md` → Contributing / TDD).

## Development

### Prerequisites
- Go **1.27+**
- golangci-lint
- goreleaser

### Building
```bash
goreleaser build --snapshot --clean
```

### Testing
```bash
go test -race -coverprofile=coverage.out -count=1 ./...
go tool cover -func=coverage.out   # total must be >= 95%

# End-to-end against a real NetBox (requires Docker; gated):
MCP_NETBOX_E2E=1 go test -race -count=1 ./e2e/...
```

The e2e test (`e2e/netbox_e2e_test.go`) boots a real NetBox via
go-docker-testsuite and drives CRUD over the full MCP protocol in-process. It is
gated behind `MCP_NETBOX_E2E=1` because it needs a running Docker daemon and
pulls a multi-container NetBox stack; CI runs `go test ./...` without Docker, so
it skips by default.

### Linting
```bash
golangci-lint run ./...   # includes gosec, gofumpt, gci
go-arch-lint check        # dependency rules (.go-arch-lint.yml)
govulncheck ./...         # findings must be fixed
```

### CI Gates
CI enforces, and **fails the build** when any gate is not met:

- **Coverage >= 95%** (C1/N6) — `go test -race` with a coverage threshold that
  fails the build below 95%.
- **Mutation testing (gremlins)** as a **hard gate** (C8/N19) — survivors fail
  the build; it is not informational/continue-on-error.
- **go-arch-lint** (C6) — dependency architecture is enforced.
- **golangci-lint** (C2) with gofumpt/gci formatting.
- **gosec** (C4) and **govulncheck** (C5) — findings are **fixed**, not
  suppressed (S5).

## Release & Container Images

Because the server supports the **HTTP (remote)** transport it is classified as
Hybrid, so a **container image is required** (R1/N17). CI/CD publishes a
multi-arch image (`linux/amd64`, `linux/arm64`) to `ghcr.io/teran/mcp-netbox`
on every commit to `master` and on every release tag.

- **On each git tag `X`** (release): image tags `X`, `X-{ts}`, `X-{commit}`, `X-{commit}-{ts}` (R3).
- **On each commit to `master`**: image tags `master-{commit}`, `master-{ts}`, `master-{commit}-{ts}` (R4).

There is **no `latest` tag**; deployments pin an immutable commit/timestamp tag.
The default branch is `master` (R2).
