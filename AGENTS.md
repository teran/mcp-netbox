# AGENTS.md — Agent Documentation

## Overview

This document describes the agents/assistants involved in the development and operation of `mcp-netbox`.

## Agent Roles

### User Agent
- **Role**: End-user interacting via an MCP-compatible AI assistant (e.g., Claude, Copilot).
- **Scope**: Sends natural-language queries that get translated into MCP tool calls.
- **No direct access** to NetBox API.

### MCP Server (`mcp-netbox`)
- **Role**: Mediator between the AI assistant and NetBox.
- **Scope**: Translates MCP tool invocations into NetBox REST API calls.
- **Transport**: **Hybrid** — Streamable HTTP (remote) or STDIO (local), selected via `TRANSPORT` (default `http`).
- **Responsible for**: Transparent token relay, request routing, response formatting.
- **Token handling**: HTTP → per-request `Authorization: Bearer` header; STDIO → `NETBOX_TOKEN` env var at startup.
- **Does not validate tokens** — authentication and authorization are delegated entirely to the NetBox backend.
- **OAuth2** is not used; NetBox personal access tokens are relayed directly.

### NetBox
- **Role**: Infrastructure source of truth backend.
- **Scope**: Stores and serves DCIM, IPAM, virtualization, tenancy, and circuits data.
- **API**: REST JSON API under `/api/`.

## Package Layout

| Package / File                              | Purpose                                         |
|---------------------------------------------|-------------------------------------------------|
| `cmd/server/main.go`                        | Entrypoint; HTTP/STDIO transport wiring, MCP server + instructions, SDK slog→logrus logger, B5 startup banner, build metadata (ldflags `appName/appVersion/appCommitHash/appTimestamp`) |
| `config/config.go`                          | Configuration loading (`envconfig` + ozzo-validation); includes `TRANSPORT`, `NETBOX_TOKEN`, `LOG_*` |
| `handlers/middleware.go`                    | Token extraction, body limit, request_id injection, logging/trace, batch validation middleware |
| `handlers/ratelimit.go`                     | Rate limiting middleware (global + per-client)  |
| `handlers/metrics.go`                       | Prometheus metrics collectors + middleware + `WrapToolHandler` |
| `handlers/tools.go`                         | MCP tool handler factories + I/O types          |
| `handlers/registration.go`                  | Tool registration via `RegisterTools()` (annotations + instructions); tool definitions in `handlers/tooldefs.go` |
| `handlers/server.go`                        | HTTP mux builder, middleware chain assembly, service-per-request injection |
| `application/service.go`                    | Business logic / use case layer                 |
| `domain/`                                   | Domain models + repository interfaces (ports), request_id context helpers (`WithRequestID`/`RequestIDFromContext`) |
| `domain/repository.go`                      | NetworkRepository interface (port), RawObject type for generic object retrieval |
| `infrastructure/netbox/client.go`           | NetBox REST client via `resty.dev/v3` (DNS-rebinding dialer + circuit breaker transport), `X-Request-ID` forwarding + per-request outbound log |
| `infrastructure/netbox/models.go`           | JSON wire models + `toDomain()` conversion      |
| `internal/logging/logging.go`               | logrus setup: channel-by-transport, `LOG_LEVEL` gating, `LOG_FORMAT`/`LOG_FILENAME` |
| `internal/logging/slog.go`                  | slog→logrus handler, `NewSlogLogger` (SDK logger wiring, L7/G10) |

## Tool-to-Agent Mapping

| MCP Tool                   | Agent Role | NetBox Endpoint                                   |
|----------------------------|------------|---------------------------------------------------|
| `get_sites`                | MCP Server | `GET /api/dcim/sites/`                            |
| `get_devices`              | MCP Server | `GET /api/dcim/devices/`                          |
| `get_ip_addresses`         | MCP Server | `GET /api/ipam/ip-addresses/`                     |
| `get_prefixes`             | MCP Server | `GET /api/ipam/prefixes/`                         |
| `get_vlans`                | MCP Server | `GET /api/ipam/vlans/`                            |
| `get_virtual_machines`     | MCP Server | `GET /api/virtualization/virtual-machines/`       |
| `get_clusters`             | MCP Server | `GET /api/virtualization/clusters/`               |
| `get_circuits`             | MCP Server | `GET /api/circuits/circuits/`                     |
| `get_racks`                | MCP Server | `GET /api/dcim/racks/`                            |
| `get_interfaces`           | MCP Server | `GET /api/dcim/interfaces/`                       |
| `get_vm_interfaces`        | MCP Server | `GET /api/virtualization/interfaces/`             |
| `get_circuit_terminations` | MCP Server | `GET /api/circuits/circuit-terminations/`         |
| `get_cables`               | MCP Server | `GET /api/dcim/cables/`                           |
| `get_object_by_id`         | MCP Server | `GET /api/*/<type>/<id>/`                         |

## Metrics

The server exposes Prometheus metrics on a separate HTTP server (default port `:8081`):

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `mcp_tool_requests_total` | Counter | `{tool, status_class}` | Per-tool request count |
| `mcp_tool_duration_seconds` | Histogram | `{tool}` | Per-tool request duration |
| `mcp_active_requests` | Gauge | — | Current in-flight requests |
| `go_*` | Various | — | Go runtime metrics |

## CI Pipeline

Every commit on any branch is checked by the following CI gates (a failed gate
blocks the build):

1. **golangci-lint** — static analysis with `gosec`, gofumpt and gci formatting.
2. **go test -race** — race-enabled unit tests with a **coverage gate of 95%**
   (the build fails below it).
3. **govulncheck** — dependency vulnerability audit; findings are **fixed**.
4. **go-arch-lint** — dependency architecture rules (`.go-arch-lint.yml`).
5. **gremlins** — mutation testing as a **hard gate** (survivors fail the
   build; not informational).

Additional workflows: `master.yml` and `release.yml` build & publish the Docker
image with the documented tag scheme; `ci.yml` also runs a gitleaks secret scan.

## Development Agents

| Agent       | Responsible For                                    |
|-------------|----------------------------------------------------|
| `architect` | High-level design decisions, system boundaries     |
| `developer` | Writing Go code, implementing tools and client     |
| `qa`        | Writing tests, verifying correctness               |
| `security`  | Reviewing auth flow, token handling                |
| `code-review` | Reviewing merge requests before deployment      |
| `devops`    | CI/CD pipelines, Docker image, deployment          |

## Conflict Resolution

1. **Security first** — any recommendation that weakens the auth boundary is rejected.
2. **SPEC compliance** — the choice that best matches SPEC.md wins.
3. **Simplicity** — prefer the solution with fewer moving parts.
4. **Go idioms** — prefer standard library over external dependencies.
