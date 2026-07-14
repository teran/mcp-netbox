// Package domain provides the core domain models, repository interfaces (ports),
// and generic types for the NetBox MCP server.
//
// This package is free of external dependencies and defines:
// - Domain models for all NetBox entity types (Site, Device, IPAddress, etc.)
// - The NetworkRepository interface (port) that infrastructure adapters implement
// - Generic PaginatedResponse[T] for typed paginated API results
// - RawObject type for generic object retrieval via get_object_by_id
package domain

import (
	"context"
	"encoding/json"
)

// RawObject is a JSON-preserving type for generic object retrieval via get_object_by_id.
// It preserves the exact JSON structure from NetBox, avoiding map[string]interface{} lossiness.
type RawObject json.RawMessage

// NetworkRepository defines the port (interface) for NetBox data access.
// Implementations provide CRUD operations for sites, devices, IP addresses,
// prefixes, VLANs, virtual machines, clusters, circuits, and racks.
type NetworkRepository interface {
	ListSites(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[Site], error)
	ListDevices(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[Device], error)
	ListIPAddresses(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[IPAddress], error)
	ListPrefixes(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[Prefix], error)
	ListVLANs(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[VLAN], error)
	ListVirtualMachines(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[VirtualMachine], error)
	ListClusters(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[Cluster], error)
	ListCircuits(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[Circuit], error)
	ListRacks(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[Rack], error)
	ListInterfaces(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[Interface], error)
	ListVMInterfaces(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[VMInterface], error)
	ListCircuitTerminations(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[CircuitTermination], error)
	ListCables(ctx context.Context, token string, params map[string]string) (*PaginatedResponse[Cable], error)
	GetObject(ctx context.Context, token string, objectType string, id int, params map[string]string) (RawObject, error)
}
