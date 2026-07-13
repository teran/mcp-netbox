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
- **Responsible for**: Transparent token relay, request routing, response formatting.
- **Does not validate tokens** — authentication and authorization are delegated entirely to the NetBox backend.

### NetBox
- **Role**: Infrastructure source of truth backend.
- **Scope**: Stores and serves DCIM, IPAM, virtualization, tenancy, and circuits data.
- **API**: REST JSON API under `/api/`.

## Package Layout

| Package / File                              | Purpose                                         |
|---------------------------------------------|-------------------------------------------------|
| `cmd/server/main.go`                        | Entrypoint, HTTP server, middleware wiring      |
| `config/config.go`                          | Configuration loading (`envconfig` + ozzo-validation) |
| `handlers/middleware.go`                    | Token extraction, body limit, logging, batch validation middleware |
| `handlers/ratelimit.go`                     | Rate limiting middleware (global + per-client)  |
| `handlers/metrics.go`                       | Prometheus metrics collectors + middleware + `WrapToolHandler` |
| `handlers/tools.go`                         | MCP tool handler factories + I/O types          |
| `handlers/registration.go`                  | Tool registration via `RegisterTools()`         |
| `handlers/server.go`                        | HTTP mux builder, middleware chain assembly, service-per-request injection |
| `application/service.go`                    | Business logic / use case layer                 |
| `domain/`                                   | Domain models + repository interfaces (ports)   |
| `domain/repository.go`                      | NetworkRepository interface (port), RawObject type for generic object retrieval |
| `infrastructure/netbox/client.go`           | NetBox HTTP API client (adapters)               |
| `infrastructure/netbox/models.go`           | JSON wire models + `toDomain()` conversion      |

## Tool-to-Agent Mapping

| MCP Tool               | Agent Role | NetBox Endpoint                           |
|------------------------|------------|-------------------------------------------|
| `get_sites`            | MCP Server | `GET /api/dcim/sites/`                    |
| `get_devices`          | MCP Server | `GET /api/dcim/devices/`                  |
| `get_ip_addresses`     | MCP Server | `GET /api/ipam/ip-addresses/`             |
| `get_prefixes`         | MCP Server | `GET /api/ipam/prefixes/`                 |
| `get_vlans`            | MCP Server | `GET /api/ipam/vlans/`                    |
| `get_virtual_machines` | MCP Server | `GET /api/virtualization/virtual-machines/` |
| `get_clusters`         | MCP Server | `GET /api/virtualization/clusters/`       |
| `get_circuits`         | MCP Server | `GET /api/circuits/circuits/`             |
| `get_object_by_id`     | MCP Server | `GET /api/*/<type>/<id>/`                 |
| `get_racks`            | MCP Server | `GET /api/dcim/racks/`                    |

## Metrics

The server exposes Prometheus metrics on a separate HTTP server (default port `:8081`):

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `mcp_tool_requests_total` | Counter | `{tool, status_class}` | Per-tool request count |
| `mcp_tool_duration_seconds` | Histogram | `{tool}` | Per-tool request duration |
| `mcp_active_requests` | Gauge | — | Current in-flight requests |
| `go_*` | Various | — | Go runtime metrics |

## CI Pipeline

Every commit on any branch is checked by three workflows:

1. **golangci-lint** — static analysis with `gosec` enabled.
2. **go test** — unit tests with coverage profile (threshold: 85%).
3. **gremlins unleash** — mutation testing (informational, does not block).

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
