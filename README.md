# MCP NetBox

MCP (Model Context Protocol) server for [NetBox](https://netboxlabs.com/) — infrastructure source of truth.

This server exposes NetBox DCIM, IPAM, virtualization, tenancy, and circuits data through the MCP protocol using **Streamable HTTP** transport (remote mode), allowing AI assistants to query your NetBox instance.

## Features

- **Read-only** — only exposes `GET` operations. No create, update, or delete capabilities.
- **Remote (HTTP) transport** — uses MCP Streamable HTTP protocol.
- **Per-request token authentication** — the NetBox API token is passed in the `Authorization` header of each MCP request. No server-side token storage.
- **Comprehensive NetBox coverage** — sites, devices, IP addresses, prefixes, VLANs, VMs, clusters, circuits, racks, and generic `get_object_by_id`.
- **Prometheus metrics** — on a separate HTTP server (default `:8081`).
- **Rate limiting** — configurable global and per-client rate limits.

## Tools

| Tool | Description |
|------|-------------|
| `get_sites` | List sites with optional filters |
| `get_devices` | List devices with search/filters |
| `get_ip_addresses` | Search IP addresses |
| `get_prefixes` | Search IP prefixes |
| `get_vlans` | List VLANs |
| `get_virtual_machines` | List virtual machines |
| `get_clusters` | List clusters |
| `get_circuits` | List circuits |
| `get_racks` | List racks |
| `get_object_by_id` | Get any object by type and ID |

## Configuration

All configuration is via environment variables:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `NETBOX_URL` | Yes | — | Base URL of the NetBox instance |
| `LISTEN_ADDR` | No | `:8080` | TCP address to listen on |
| `PROMETHEUS_METRICS_ADDR` | No | `:8081` | Prometheus `/metrics` endpoint |
| `RATE_LIMIT_GLOBAL` | No | `100` | Global rate limit (requests/second) |
| `RATE_LIMIT_PER_CLIENT` | No | `10` | Per-client IP rate limit |
| `WRITE_TIMEOUT` | No | `300` | HTTP write timeout in seconds |

The NetBox API token is supplied per-request in the `Authorization` header as `Bearer <token>`. Both v1 (`Token <token>`) and v2 (`Bearer nbt_<key>.<token>`) tokens are supported.

## MCP Client Configuration

For use with OpenCode or any MCP client:

```json
{
  "mcp": {
    "netbox": {
      "type": "local",
      "enabled": true,
      "command": ["/path/to/mcp-netbox"],
      "environment": {
        "NETBOX_URL": "http://netbox:8000"
      }
    }
  }
}
```

For remote usage via SSE/Streamable HTTP:

```json
{
  "mcp": {
    "netbox": {
      "type": "remote",
      "enabled": true,
      "url": "http://localhost:8080/mcp",
      "headers": {
        "Authorization": "Bearer <netbox-api-token>"
      }
    }
  }
}
```

## Building

```bash
# Using goreleaser (recommended)
goreleaser build --snapshot --clean

# Using go build directly
go build -o mcp-netbox ./cmd/server
```

## Testing

```bash
go test -race -coverprofile=coverage.out -count=1 ./...
go tool cover -func=coverage.out
```

## Linting

```bash
golangci-lint run ./...
```

## Docker

```bash
docker buildx build --platform linux/amd64,linux/arm64 -t ghcr.io/teran/mcp-netbox:latest .
```

## License

Apache 2.0
