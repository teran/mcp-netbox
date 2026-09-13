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
//
// It implements json.Marshaler and json.Unmarshaler by delegating to the
// underlying json.RawMessage, ensuring correct JSON-in-JSON serialization
// (rather than base64-encoded output) when used as a field in a struct.
type RawObject json.RawMessage

// MarshalJSON implements json.Marshaler so RawObject serialises as inline JSON
// instead of a base64-encoded string.
func (r RawObject) MarshalJSON() ([]byte, error) {
	return json.RawMessage(r).MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler.
func (r *RawObject) UnmarshalJSON(b []byte) error {
	return (*json.RawMessage)(r).UnmarshalJSON(b)
}

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

	CreateSite(ctx context.Context, token string, in SiteWrite) (*Site, error)
	UpdateSite(ctx context.Context, token string, id int, in SiteWrite) (*Site, error)
	DeleteSite(ctx context.Context, token string, id int) error

	CreateDevice(ctx context.Context, token string, in DeviceWrite) (*Device, error)
	UpdateDevice(ctx context.Context, token string, id int, in DeviceWrite) (*Device, error)
	DeleteDevice(ctx context.Context, token string, id int) error

	CreateIPAddress(ctx context.Context, token string, in IPAddressWrite) (*IPAddress, error)
	UpdateIPAddress(ctx context.Context, token string, id int, in IPAddressWrite) (*IPAddress, error)
	DeleteIPAddress(ctx context.Context, token string, id int) error

	CreatePrefix(ctx context.Context, token string, in PrefixWrite) (*Prefix, error)
	UpdatePrefix(ctx context.Context, token string, id int, in PrefixWrite) (*Prefix, error)
	DeletePrefix(ctx context.Context, token string, id int) error

	CreateVLAN(ctx context.Context, token string, in VLANWrite) (*VLAN, error)
	UpdateVLAN(ctx context.Context, token string, id int, in VLANWrite) (*VLAN, error)
	DeleteVLAN(ctx context.Context, token string, id int) error

	CreateVirtualMachine(ctx context.Context, token string, in VirtualMachineWrite) (*VirtualMachine, error)
	UpdateVirtualMachine(ctx context.Context, token string, id int, in VirtualMachineWrite) (*VirtualMachine, error)
	DeleteVirtualMachine(ctx context.Context, token string, id int) error

	CreateCluster(ctx context.Context, token string, in ClusterWrite) (*Cluster, error)
	UpdateCluster(ctx context.Context, token string, id int, in ClusterWrite) (*Cluster, error)
	DeleteCluster(ctx context.Context, token string, id int) error
}
