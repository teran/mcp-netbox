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
- **Health endpoint** — `GET /healthz` returns `{"status":"ok"}`.
- **Pagination** — all list tools support `page` and `page_size` parameters with `next`/`previous` navigation URLs.

## Tools

| Tool | Description |
|------|-------------|
| `get_sites` | List sites with optional filters and free-text search |
| `get_devices` | List devices with filters and free-text search |
| `get_ip_addresses` | Search IP addresses with free-text search |
| `get_prefixes` | Search IP prefixes with free-text search |
| `get_vlans` | List VLANs with free-text search |
| `get_virtual_machines` | List virtual machines with free-text search |
| `get_clusters` | List clusters with free-text search |
| `get_circuits` | List circuits with free-text search |
| `get_racks` | List racks with free-text search |
| `get_object_by_id` | Get any object by type and ID |

> All list tools support a `q` parameter for free-text search across all fields.

## Configuration

All configuration is via environment variables:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `NETBOX_URL` | Yes | — | Base URL of the NetBox instance |
| `LISTEN_ADDR` | No | `:8080` | TCP address to listen on |
| `PROMETHEUS_METRICS_ADDR` | No | `:8081` | Prometheus `/metrics` endpoint |
| `RATE_LIMIT_GLOBAL` | No | `100` | Global rate limit (requests/second) |
| `RATE_LIMIT_PER_CLIENT` | No | `10` | Per-client IP rate limit |
| `WRITE_TIMEOUT` | No | `300s` (5 minutes) | HTTP write timeout (Go duration, minimum 1s) |

The NetBox API token is supplied per-request in the `Authorization` header as `Bearer <token>`. Both v1 (`Token <token>`) and v2 (`Bearer nbt_<key>.<token>`) tokens are supported.

## MCP Client Configuration

The server uses the **Streamable HTTP** MCP transport (remote mode).
It does not support stdio/local mode. All clients must connect over HTTP.

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

## Security

- **Token handling**: The NetBox API token is passed per-request in the `Authorization` header. It is never stored on the server, written to logs, or persisted between requests.
- **Read-only**: The server only exposes GET operations. No write access to NetBox.
- **TLS**: Terminate TLS at a reverse proxy (nginx, Envoy) placed in front of the server.
- **Rate limiting**: Built-in rate limiting prevents abuse (configurable via environment variables).

## Troubleshooting

### 401 Unauthorized
The Authorization header is missing or malformed. Ensure you pass `Bearer <token>`.

### 429 Too Many Requests
Rate limit exceeded. Increase `RATE_LIMIT_GLOBAL` or `RATE_LIMIT_PER_CLIENT`, or wait before retrying.

### Connection issues
Verify `NETBOX_URL` is reachable from the server. Check `GET /healthz` endpoint.

### No results returned
Verify the API token has sufficient permissions in NetBox. Some filters may not match any objects.

## License

Apache 2.0
