package domain

import "context"

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
	GetObject(ctx context.Context, token string, objectType string, id int, params map[string]string) (interface{}, error)
}
