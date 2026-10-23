// Package application provides the business logic / use case layer.
package application

import (
	"context"
	"fmt"

	"github.com/teran/mcp-netbox/domain"
)

// NetworkService is the service layer that translates tool requests into
// repository calls with the user's token.
type NetworkService struct {
	repo  domain.NetworkRepository
	token *Token
}

// NewNetworkService creates a new NetworkService backed by the given repository,
// using the given API token for all requests.
func NewNetworkService(repo domain.NetworkRepository, token string) *NetworkService {
	return &NetworkService{
		repo:  repo,
		token: NewToken(token),
	}
}

// ListSites returns the page of sites matching params.
func (s *NetworkService) ListSites(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.Site], error) {
	resp, err := s.repo.ListSites(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list sites: %w", err)
	}

	return resp, nil
}

// ListRegions returns the page of regions matching params.
func (s *NetworkService) ListRegions(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.Region], error) {
	resp, err := s.repo.ListRegions(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list regions: %w", err)
	}

	return resp, nil
}

// ListTenantGroups returns the page of tenant groups matching params.
func (s *NetworkService) ListTenantGroups(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.TenantGroup], error) {
	resp, err := s.repo.ListTenantGroups(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list tenant groups: %w", err)
	}

	return resp, nil
}

// ListDevices returns the page of devices matching params.
func (s *NetworkService) ListDevices(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.Device], error) {
	resp, err := s.repo.ListDevices(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list devices: %w", err)
	}

	return resp, nil
}

// ListIPAddresses returns the page of IP addresses matching params.
func (s *NetworkService) ListIPAddresses(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.IPAddress], error) {
	resp, err := s.repo.ListIPAddresses(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list IP addresses: %w", err)
	}

	return resp, nil
}

// ListPrefixes returns the page of prefixes matching params.
func (s *NetworkService) ListPrefixes(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.Prefix], error) {
	resp, err := s.repo.ListPrefixes(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list prefixes: %w", err)
	}

	return resp, nil
}

// ListVLANs returns the page of VLANs matching params.
func (s *NetworkService) ListVLANs(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.VLAN], error) {
	resp, err := s.repo.ListVLANs(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list VLANs: %w", err)
	}

	return resp, nil
}

// ListVirtualMachines returns the page of virtual machines matching params.
func (s *NetworkService) ListVirtualMachines(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.VirtualMachine], error) {
	resp, err := s.repo.ListVirtualMachines(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list virtual machines: %w", err)
	}

	return resp, nil
}

// ListClusters returns the page of clusters matching params.
func (s *NetworkService) ListClusters(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.Cluster], error) {
	resp, err := s.repo.ListClusters(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list clusters: %w", err)
	}

	return resp, nil
}

// ListCircuits returns the page of circuits matching params.
func (s *NetworkService) ListCircuits(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.Circuit], error) {
	resp, err := s.repo.ListCircuits(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list circuits: %w", err)
	}

	return resp, nil
}

// ListInterfaces returns the page of device interfaces matching params.
func (s *NetworkService) ListInterfaces(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.Interface], error) {
	resp, err := s.repo.ListInterfaces(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list interfaces: %w", err)
	}

	return resp, nil
}

// ListVMInterfaces returns the page of virtual-machine interfaces matching params.
func (s *NetworkService) ListVMInterfaces(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.VMInterface], error) {
	resp, err := s.repo.ListVMInterfaces(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list VM interfaces: %w", err)
	}

	return resp, nil
}

// ListRacks returns the page of racks matching params.
func (s *NetworkService) ListRacks(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.Rack], error) {
	resp, err := s.repo.ListRacks(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list racks: %w", err)
	}

	return resp, nil
}

// ListCircuitTerminations returns the page of circuit terminations matching params.
func (s *NetworkService) ListCircuitTerminations(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.CircuitTermination], error) {
	resp, err := s.repo.ListCircuitTerminations(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list circuit terminations: %w", err)
	}

	return resp, nil
}

// ListCables returns the page of cables matching params.
func (s *NetworkService) ListCables(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.Cable], error) {
	resp, err := s.repo.ListCables(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list cables: %w", err)
	}

	return resp, nil
}

// GetObject returns the raw object of the given type and id, preserving the
// exact JSON structure returned by NetBox.
func (s *NetworkService) GetObject(
	ctx context.Context, objectType string, id int, params map[string]string,
) (domain.RawObject, error) {
	resp, err := s.repo.GetObject(ctx, s.token.Value(), objectType, id, params)
	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}

	return resp, nil
}

// CreateSite relays a site create to the repository. ValidationError is wrapped
// (via %w) and left intact so callers can recover it with errors.As.
func (s *NetworkService) CreateSite(ctx context.Context, in domain.SiteWrite) (*domain.Site, error) {
	resp, err := s.repo.CreateSite(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create site: %w", err)
	}

	return resp, nil
}

// UpdateSite relays a site update to the repository.
func (s *NetworkService) UpdateSite(ctx context.Context, id int, in domain.SiteWrite) (*domain.Site, error) {
	resp, err := s.repo.UpdateSite(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update site: %w", err)
	}

	return resp, nil
}

// DeleteSite relays a site delete to the repository.
func (s *NetworkService) DeleteSite(ctx context.Context, id int) error {
	if err := s.repo.DeleteSite(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete site: %w", err)
	}

	return nil
}

// CreateRegion relays a region create to the repository. ValidationError is
// wrapped (via %w) and left intact so callers can recover it with errors.As.
func (s *NetworkService) CreateRegion(ctx context.Context, in domain.RegionWrite) (*domain.Region, error) {
	resp, err := s.repo.CreateRegion(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create region: %w", err)
	}

	return resp, nil
}

// UpdateRegion relays a region update to the repository.
func (s *NetworkService) UpdateRegion(ctx context.Context, id int, in domain.RegionWrite) (*domain.Region, error) {
	resp, err := s.repo.UpdateRegion(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update region: %w", err)
	}

	return resp, nil
}

// DeleteRegion relays a region delete to the repository.
func (s *NetworkService) DeleteRegion(ctx context.Context, id int) error {
	if err := s.repo.DeleteRegion(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete region: %w", err)
	}

	return nil
}

// CreateTenantGroup relays a tenant group create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover it
// with errors.As.
func (s *NetworkService) CreateTenantGroup(
	ctx context.Context, in domain.TenantGroupWrite,
) (*domain.TenantGroup, error) {
	resp, err := s.repo.CreateTenantGroup(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create tenant group: %w", err)
	}

	return resp, nil
}

// UpdateTenantGroup relays a tenant group update to the repository.
func (s *NetworkService) UpdateTenantGroup(
	ctx context.Context, id int, in domain.TenantGroupWrite,
) (*domain.TenantGroup, error) {
	resp, err := s.repo.UpdateTenantGroup(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update tenant group: %w", err)
	}

	return resp, nil
}

// DeleteTenantGroup relays a tenant group delete to the repository.
func (s *NetworkService) DeleteTenantGroup(ctx context.Context, id int) error {
	if err := s.repo.DeleteTenantGroup(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete tenant group: %w", err)
	}

	return nil
}

// CreateDevice relays a device create to the repository. ValidationError is
// wrapped (via %w) and left intact so callers can recover it with errors.As.
func (s *NetworkService) CreateDevice(ctx context.Context, in domain.DeviceWrite) (*domain.Device, error) {
	resp, err := s.repo.CreateDevice(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create device: %w", err)
	}

	return resp, nil
}

// UpdateDevice relays a device update to the repository.
func (s *NetworkService) UpdateDevice(ctx context.Context, id int, in domain.DeviceWrite) (*domain.Device, error) {
	resp, err := s.repo.UpdateDevice(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update device: %w", err)
	}

	return resp, nil
}

// DeleteDevice relays a device delete to the repository.
func (s *NetworkService) DeleteDevice(ctx context.Context, id int) error {
	if err := s.repo.DeleteDevice(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete device: %w", err)
	}

	return nil
}

// CreateIPAddress relays an IP address create to the repository. ValidationError
// is wrapped (via %w) and left intact so callers can recover it with errors.As.
func (s *NetworkService) CreateIPAddress(ctx context.Context, in domain.IPAddressWrite) (*domain.IPAddress, error) {
	resp, err := s.repo.CreateIPAddress(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create IP address: %w", err)
	}

	return resp, nil
}

// UpdateIPAddress relays an IP address update to the repository.
func (s *NetworkService) UpdateIPAddress(
	ctx context.Context, id int, in domain.IPAddressWrite,
) (*domain.IPAddress, error) {
	resp, err := s.repo.UpdateIPAddress(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update IP address: %w", err)
	}

	return resp, nil
}

// DeleteIPAddress relays an IP address delete to the repository.
func (s *NetworkService) DeleteIPAddress(ctx context.Context, id int) error {
	if err := s.repo.DeleteIPAddress(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete IP address: %w", err)
	}

	return nil
}

// CreatePrefix relays a prefix create to the repository. ValidationError is
// wrapped (via %w) and left intact so callers can recover it with errors.As.
func (s *NetworkService) CreatePrefix(ctx context.Context, in domain.PrefixWrite) (*domain.Prefix, error) {
	resp, err := s.repo.CreatePrefix(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create prefix: %w", err)
	}

	return resp, nil
}

// UpdatePrefix relays a prefix update to the repository.
func (s *NetworkService) UpdatePrefix(ctx context.Context, id int, in domain.PrefixWrite) (*domain.Prefix, error) {
	resp, err := s.repo.UpdatePrefix(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update prefix: %w", err)
	}

	return resp, nil
}

// DeletePrefix relays a prefix delete to the repository.
func (s *NetworkService) DeletePrefix(ctx context.Context, id int) error {
	if err := s.repo.DeletePrefix(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete prefix: %w", err)
	}

	return nil
}

// CreateVLAN relays a VLAN create to the repository. ValidationError is wrapped
// (via %w) and left intact so callers can recover it with errors.As.
func (s *NetworkService) CreateVLAN(ctx context.Context, in domain.VLANWrite) (*domain.VLAN, error) {
	resp, err := s.repo.CreateVLAN(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create VLAN: %w", err)
	}

	return resp, nil
}

// UpdateVLAN relays a VLAN update to the repository.
func (s *NetworkService) UpdateVLAN(ctx context.Context, id int, in domain.VLANWrite) (*domain.VLAN, error) {
	resp, err := s.repo.UpdateVLAN(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update VLAN: %w", err)
	}

	return resp, nil
}

// DeleteVLAN relays a VLAN delete to the repository.
func (s *NetworkService) DeleteVLAN(ctx context.Context, id int) error {
	if err := s.repo.DeleteVLAN(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete VLAN: %w", err)
	}

	return nil
}

// CreateVirtualMachine relays a virtual machine create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover it
// with errors.As.
func (s *NetworkService) CreateVirtualMachine(
	ctx context.Context, in domain.VirtualMachineWrite,
) (*domain.VirtualMachine, error) {
	resp, err := s.repo.CreateVirtualMachine(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create virtual machine: %w", err)
	}

	return resp, nil
}

// UpdateVirtualMachine relays a virtual machine update to the repository.
func (s *NetworkService) UpdateVirtualMachine(
	ctx context.Context, id int, in domain.VirtualMachineWrite,
) (*domain.VirtualMachine, error) {
	resp, err := s.repo.UpdateVirtualMachine(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update virtual machine: %w", err)
	}

	return resp, nil
}

// DeleteVirtualMachine relays a virtual machine delete to the repository.
func (s *NetworkService) DeleteVirtualMachine(ctx context.Context, id int) error {
	if err := s.repo.DeleteVirtualMachine(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete virtual machine: %w", err)
	}

	return nil
}

// CreateCluster relays a cluster create to the repository. ValidationError is
// wrapped (via %w) and left intact so callers can recover it with errors.As.
func (s *NetworkService) CreateCluster(ctx context.Context, in domain.ClusterWrite) (*domain.Cluster, error) {
	resp, err := s.repo.CreateCluster(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create cluster: %w", err)
	}

	return resp, nil
}

// UpdateCluster relays a cluster update to the repository.
func (s *NetworkService) UpdateCluster(ctx context.Context, id int, in domain.ClusterWrite) (*domain.Cluster, error) {
	resp, err := s.repo.UpdateCluster(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update cluster: %w", err)
	}

	return resp, nil
}

// DeleteCluster relays a cluster delete to the repository.
func (s *NetworkService) DeleteCluster(ctx context.Context, id int) error {
	if err := s.repo.DeleteCluster(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete cluster: %w", err)
	}

	return nil
}

// CreateCircuit relays a circuit create to the repository. ValidationError is
// wrapped (via %w) and left intact so callers can recover it with errors.As.
func (s *NetworkService) CreateCircuit(ctx context.Context, in domain.CircuitWrite) (*domain.Circuit, error) {
	resp, err := s.repo.CreateCircuit(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create circuit: %w", err)
	}

	return resp, nil
}

// UpdateCircuit relays a circuit update to the repository.
func (s *NetworkService) UpdateCircuit(ctx context.Context, id int, in domain.CircuitWrite) (*domain.Circuit, error) {
	resp, err := s.repo.UpdateCircuit(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update circuit: %w", err)
	}

	return resp, nil
}

// DeleteCircuit relays a circuit delete to the repository.
func (s *NetworkService) DeleteCircuit(ctx context.Context, id int) error {
	if err := s.repo.DeleteCircuit(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete circuit: %w", err)
	}

	return nil
}

// CreateRack relays a rack create to the repository. ValidationError is wrapped
// (via %w) and left intact so callers can recover it with errors.As.
func (s *NetworkService) CreateRack(ctx context.Context, in domain.RackWrite) (*domain.Rack, error) {
	resp, err := s.repo.CreateRack(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create rack: %w", err)
	}

	return resp, nil
}

// UpdateRack relays a rack update to the repository.
func (s *NetworkService) UpdateRack(ctx context.Context, id int, in domain.RackWrite) (*domain.Rack, error) {
	resp, err := s.repo.UpdateRack(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update rack: %w", err)
	}

	return resp, nil
}

// DeleteRack relays a rack delete to the repository.
func (s *NetworkService) DeleteRack(ctx context.Context, id int) error {
	if err := s.repo.DeleteRack(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete rack: %w", err)
	}

	return nil
}

// CreateInterface relays an interface create to the repository. ValidationError
// is wrapped (via %w) and left intact so callers can recover it with errors.As.
func (s *NetworkService) CreateInterface(ctx context.Context, in domain.InterfaceWrite) (*domain.Interface, error) {
	resp, err := s.repo.CreateInterface(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create interface: %w", err)
	}

	return resp, nil
}

// UpdateInterface relays an interface update to the repository.
func (s *NetworkService) UpdateInterface(
	ctx context.Context, id int, in domain.InterfaceWrite,
) (*domain.Interface, error) {
	resp, err := s.repo.UpdateInterface(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update interface: %w", err)
	}

	return resp, nil
}

// DeleteInterface relays an interface delete to the repository.
func (s *NetworkService) DeleteInterface(ctx context.Context, id int) error {
	if err := s.repo.DeleteInterface(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete interface: %w", err)
	}

	return nil
}

// CreateCircuitTermination relays a circuit termination create to the
// repository. ValidationError is wrapped (via %w) and left intact so callers
// can recover it with errors.As.
func (s *NetworkService) CreateCircuitTermination(
	ctx context.Context, in domain.CircuitTerminationWrite,
) (*domain.CircuitTermination, error) {
	resp, err := s.repo.CreateCircuitTermination(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create circuit termination: %w", err)
	}

	return resp, nil
}

// UpdateCircuitTermination relays a circuit termination update to the
// repository.
func (s *NetworkService) UpdateCircuitTermination(
	ctx context.Context, id int, in domain.CircuitTerminationWrite,
) (*domain.CircuitTermination, error) {
	resp, err := s.repo.UpdateCircuitTermination(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update circuit termination: %w", err)
	}

	return resp, nil
}

// DeleteCircuitTermination relays a circuit termination delete to the
// repository.
func (s *NetworkService) DeleteCircuitTermination(ctx context.Context, id int) error {
	if err := s.repo.DeleteCircuitTermination(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete circuit termination: %w", err)
	}

	return nil
}

// CreateCable relays a cable create to the repository. ValidationError is
// wrapped (via %w) and left intact so callers can recover it with errors.As.
func (s *NetworkService) CreateCable(ctx context.Context, in domain.CableWrite) (*domain.Cable, error) {
	resp, err := s.repo.CreateCable(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create cable: %w", err)
	}

	return resp, nil
}

// UpdateCable relays a cable update to the repository.
func (s *NetworkService) UpdateCable(ctx context.Context, id int, in domain.CableWrite) (*domain.Cable, error) {
	resp, err := s.repo.UpdateCable(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update cable: %w", err)
	}

	return resp, nil
}

// DeleteCable relays a cable delete to the repository.
func (s *NetworkService) DeleteCable(ctx context.Context, id int) error {
	if err := s.repo.DeleteCable(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete cable: %w", err)
	}

	return nil
}

// CreateVMInterface relays a VM interface create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover
// it with errors.As.
func (s *NetworkService) CreateVMInterface(
	ctx context.Context, in domain.VMInterfaceWrite,
) (*domain.VMInterface, error) {
	resp, err := s.repo.CreateVMInterface(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create vm interface: %w", err)
	}

	return resp, nil
}

// UpdateVMInterface relays a VM interface update to the repository.
func (s *NetworkService) UpdateVMInterface(
	ctx context.Context, id int, in domain.VMInterfaceWrite,
) (*domain.VMInterface, error) {
	resp, err := s.repo.UpdateVMInterface(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update vm interface: %w", err)
	}

	return resp, nil
}

// DeleteVMInterface relays a VM interface delete to the repository.
func (s *NetworkService) DeleteVMInterface(ctx context.Context, id int) error {
	if err := s.repo.DeleteVMInterface(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete vm interface: %w", err)
	}

	return nil
}

// CreateProvider relays a provider create to the repository. ValidationError is
// wrapped (via %w) and left intact so callers can recover it with errors.As.
func (s *NetworkService) CreateProvider(ctx context.Context, in domain.ProviderWrite) (*domain.Provider, error) {
	resp, err := s.repo.CreateProvider(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create provider: %w", err)
	}

	return resp, nil
}

// UpdateProvider relays a provider update to the repository.
func (s *NetworkService) UpdateProvider(
	ctx context.Context, id int, in domain.ProviderWrite,
) (*domain.Provider, error) {
	resp, err := s.repo.UpdateProvider(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update provider: %w", err)
	}

	return resp, nil
}

// DeleteProvider relays a provider delete to the repository.
func (s *NetworkService) DeleteProvider(ctx context.Context, id int) error {
	if err := s.repo.DeleteProvider(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete provider: %w", err)
	}

	return nil
}

// CreateTenant relays a tenant create to the repository. ValidationError is
// wrapped (via %w) and left intact so callers can recover it with errors.As.
func (s *NetworkService) CreateTenant(ctx context.Context, in domain.TenantWrite) (*domain.Tenant, error) {
	resp, err := s.repo.CreateTenant(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create tenant: %w", err)
	}

	return resp, nil
}

// UpdateTenant relays a tenant update to the repository.
func (s *NetworkService) UpdateTenant(ctx context.Context, id int, in domain.TenantWrite) (*domain.Tenant, error) {
	resp, err := s.repo.UpdateTenant(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update tenant: %w", err)
	}

	return resp, nil
}

// DeleteTenant relays a tenant delete to the repository.
func (s *NetworkService) DeleteTenant(ctx context.Context, id int) error {
	if err := s.repo.DeleteTenant(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete tenant: %w", err)
	}

	return nil
}

// CreateManufacturer relays a manufacturer create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover it
// with errors.As.
func (s *NetworkService) CreateManufacturer(
	ctx context.Context, in domain.ManufacturerWrite,
) (*domain.Manufacturer, error) {
	resp, err := s.repo.CreateManufacturer(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create manufacturer: %w", err)
	}

	return resp, nil
}

// UpdateManufacturer relays a manufacturer update to the repository.
func (s *NetworkService) UpdateManufacturer(
	ctx context.Context, id int, in domain.ManufacturerWrite,
) (*domain.Manufacturer, error) {
	resp, err := s.repo.UpdateManufacturer(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update manufacturer: %w", err)
	}

	return resp, nil
}

// DeleteManufacturer relays a manufacturer delete to the repository.
func (s *NetworkService) DeleteManufacturer(ctx context.Context, id int) error {
	if err := s.repo.DeleteManufacturer(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete manufacturer: %w", err)
	}

	return nil
}

// CreateDeviceType relays a device type create to the repository. ValidationError
// is wrapped (via %w) and left intact so callers can recover it with errors.As.
func (s *NetworkService) CreateDeviceType(ctx context.Context, in domain.DeviceTypeWrite) (*domain.DeviceType, error) {
	resp, err := s.repo.CreateDeviceType(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create device type: %w", err)
	}

	return resp, nil
}

// UpdateDeviceType relays a device type update to the repository.
func (s *NetworkService) UpdateDeviceType(
	ctx context.Context, id int, in domain.DeviceTypeWrite,
) (*domain.DeviceType, error) {
	resp, err := s.repo.UpdateDeviceType(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update device type: %w", err)
	}

	return resp, nil
}

// DeleteDeviceType relays a device type delete to the repository.
func (s *NetworkService) DeleteDeviceType(ctx context.Context, id int) error {
	if err := s.repo.DeleteDeviceType(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete device type: %w", err)
	}

	return nil
}

// CreateLocation relays a location create to the repository. ValidationError is
// wrapped (via %w) and left intact so callers can recover it with errors.As.
func (s *NetworkService) CreateLocation(ctx context.Context, in domain.LocationWrite) (*domain.Location, error) {
	resp, err := s.repo.CreateLocation(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create location: %w", err)
	}

	return resp, nil
}

// UpdateLocation relays a location update to the repository.
func (s *NetworkService) UpdateLocation(
	ctx context.Context, id int, in domain.LocationWrite,
) (*domain.Location, error) {
	resp, err := s.repo.UpdateLocation(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update location: %w", err)
	}

	return resp, nil
}

// DeleteLocation relays a location delete to the repository.
func (s *NetworkService) DeleteLocation(ctx context.Context, id int) error {
	if err := s.repo.DeleteLocation(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete location: %w", err)
	}

	return nil
}

// CreateClusterType relays a cluster type create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover
// it with errors.As.
func (s *NetworkService) CreateClusterType(
	ctx context.Context, in domain.ClusterTypeWrite,
) (*domain.ClusterType, error) {
	resp, err := s.repo.CreateClusterType(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create cluster type: %w", err)
	}

	return resp, nil
}

// UpdateClusterType relays a cluster type update to the repository.
func (s *NetworkService) UpdateClusterType(
	ctx context.Context, id int, in domain.ClusterTypeWrite,
) (*domain.ClusterType, error) {
	resp, err := s.repo.UpdateClusterType(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update cluster type: %w", err)
	}

	return resp, nil
}

// DeleteClusterType relays a cluster type delete to the repository.
func (s *NetworkService) DeleteClusterType(ctx context.Context, id int) error {
	if err := s.repo.DeleteClusterType(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete cluster type: %w", err)
	}

	return nil
}

// CreateClusterGroup relays a cluster group create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover
// it with errors.As.
func (s *NetworkService) CreateClusterGroup(
	ctx context.Context, in domain.ClusterGroupWrite,
) (*domain.ClusterGroup, error) {
	resp, err := s.repo.CreateClusterGroup(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create cluster group: %w", err)
	}

	return resp, nil
}

// UpdateClusterGroup relays a cluster group update to the repository.
func (s *NetworkService) UpdateClusterGroup(
	ctx context.Context, id int, in domain.ClusterGroupWrite,
) (*domain.ClusterGroup, error) {
	resp, err := s.repo.UpdateClusterGroup(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update cluster group: %w", err)
	}

	return resp, nil
}

// DeleteClusterGroup relays a cluster group delete to the repository.
func (s *NetworkService) DeleteClusterGroup(ctx context.Context, id int) error {
	if err := s.repo.DeleteClusterGroup(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete cluster group: %w", err)
	}

	return nil
}

// CreateCircuitType relays a circuit type create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover
// it with errors.As.
func (s *NetworkService) CreateCircuitType(
	ctx context.Context, in domain.CircuitTypeWrite,
) (*domain.CircuitType, error) {
	resp, err := s.repo.CreateCircuitType(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create circuit type: %w", err)
	}

	return resp, nil
}

// UpdateCircuitType relays a circuit type update to the repository.
func (s *NetworkService) UpdateCircuitType(
	ctx context.Context, id int, in domain.CircuitTypeWrite,
) (*domain.CircuitType, error) {
	resp, err := s.repo.UpdateCircuitType(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update circuit type: %w", err)
	}

	return resp, nil
}

// DeleteCircuitType relays a circuit type delete to the repository.
func (s *NetworkService) DeleteCircuitType(ctx context.Context, id int) error {
	if err := s.repo.DeleteCircuitType(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete circuit type: %w", err)
	}

	return nil
}

// CreateVrf relays a VRF create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover
// it with errors.As.
func (s *NetworkService) CreateVrf(ctx context.Context, in domain.VrfWrite) (*domain.Vrf, error) {
	resp, err := s.repo.CreateVrf(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create vrf: %w", err)
	}

	return resp, nil
}

// UpdateVrf relays a VRF update to the repository.
func (s *NetworkService) UpdateVrf(ctx context.Context, id int, in domain.VrfWrite) (*domain.Vrf, error) {
	resp, err := s.repo.UpdateVrf(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update vrf: %w", err)
	}

	return resp, nil
}

// DeleteVrf relays a VRF delete to the repository.
func (s *NetworkService) DeleteVrf(ctx context.Context, id int) error {
	if err := s.repo.DeleteVrf(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete vrf: %w", err)
	}

	return nil
}

// CreateVlanGroup relays a VLAN group create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover
// it with errors.As.
func (s *NetworkService) CreateVlanGroup(ctx context.Context, in domain.VlanGroupWrite) (*domain.VlanGroup, error) {
	resp, err := s.repo.CreateVlanGroup(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create vlan group: %w", err)
	}

	return resp, nil
}

// UpdateVlanGroup relays a VLAN group update to the repository.
func (s *NetworkService) UpdateVlanGroup(
	ctx context.Context, id int, in domain.VlanGroupWrite,
) (*domain.VlanGroup, error) {
	resp, err := s.repo.UpdateVlanGroup(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update vlan group: %w", err)
	}

	return resp, nil
}

// DeleteVlanGroup relays a VLAN group delete to the repository.
func (s *NetworkService) DeleteVlanGroup(ctx context.Context, id int) error {
	if err := s.repo.DeleteVlanGroup(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete vlan group: %w", err)
	}

	return nil
}

// CreateRole relays an IPAM role create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover
// it with errors.As.
func (s *NetworkService) CreateRole(ctx context.Context, in domain.RoleWrite) (*domain.Role, error) {
	resp, err := s.repo.CreateRole(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create role: %w", err)
	}

	return resp, nil
}

// UpdateRole relays an IPAM role update to the repository.
func (s *NetworkService) UpdateRole(ctx context.Context, id int, in domain.RoleWrite) (*domain.Role, error) {
	resp, err := s.repo.UpdateRole(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update role: %w", err)
	}

	return resp, nil
}

// DeleteRole relays an IPAM role delete to the repository.
func (s *NetworkService) DeleteRole(ctx context.Context, id int) error {
	if err := s.repo.DeleteRole(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete role: %w", err)
	}

	return nil
}

// CreateContact relays a tenancy contact create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover
// it with errors.As.
func (s *NetworkService) CreateContact(ctx context.Context, in domain.ContactWrite) (*domain.Contact, error) {
	resp, err := s.repo.CreateContact(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create contact: %w", err)
	}

	return resp, nil
}

// UpdateContact relays a tenancy contact update to the repository.
func (s *NetworkService) UpdateContact(ctx context.Context, id int, in domain.ContactWrite) (*domain.Contact, error) {
	resp, err := s.repo.UpdateContact(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update contact: %w", err)
	}

	return resp, nil
}

// DeleteContact relays a tenancy contact delete to the repository.
func (s *NetworkService) DeleteContact(ctx context.Context, id int) error {
	if err := s.repo.DeleteContact(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete contact: %w", err)
	}

	return nil
}

// ListSiteGroups returns the page of site groups matching params.
func (s *NetworkService) ListSiteGroups(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.SiteGroup], error) {
	resp, err := s.repo.ListSiteGroups(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list site groups: %w", err)
	}

	return resp, nil
}

// CreateSiteGroup relays a site group create to the repository. ValidationError
// is wrapped (via %w) and left intact so callers can recover it with errors.As.
func (s *NetworkService) CreateSiteGroup(ctx context.Context, in domain.SiteGroupWrite) (*domain.SiteGroup, error) {
	resp, err := s.repo.CreateSiteGroup(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create site group: %w", err)
	}

	return resp, nil
}

// UpdateSiteGroup relays a site group update to the repository.
func (s *NetworkService) UpdateSiteGroup(
	ctx context.Context, id int, in domain.SiteGroupWrite,
) (*domain.SiteGroup, error) {
	resp, err := s.repo.UpdateSiteGroup(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update site group: %w", err)
	}

	return resp, nil
}

// DeleteSiteGroup relays a site group delete to the repository.
func (s *NetworkService) DeleteSiteGroup(ctx context.Context, id int) error {
	if err := s.repo.DeleteSiteGroup(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete site group: %w", err)
	}

	return nil
}

// ListRackGroups returns the page of rack groups matching params.
func (s *NetworkService) ListRackGroups(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.RackGroup], error) {
	resp, err := s.repo.ListRackGroups(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list rack groups: %w", err)
	}

	return resp, nil
}

// CreateRackGroup relays a rack group create to the repository. ValidationError
// is wrapped (via %w) and left intact so callers can recover it with errors.As.
func (s *NetworkService) CreateRackGroup(ctx context.Context, in domain.RackGroupWrite) (*domain.RackGroup, error) {
	resp, err := s.repo.CreateRackGroup(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create rack group: %w", err)
	}

	return resp, nil
}

// UpdateRackGroup relays a rack group update to the repository.
func (s *NetworkService) UpdateRackGroup(
	ctx context.Context, id int, in domain.RackGroupWrite,
) (*domain.RackGroup, error) {
	resp, err := s.repo.UpdateRackGroup(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update rack group: %w", err)
	}

	return resp, nil
}

// DeleteRackGroup relays a rack group delete to the repository.
func (s *NetworkService) DeleteRackGroup(ctx context.Context, id int) error {
	if err := s.repo.DeleteRackGroup(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete rack group: %w", err)
	}

	return nil
}

// ListRackRoles returns the page of rack roles matching params.
func (s *NetworkService) ListRackRoles(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.RackRole], error) {
	resp, err := s.repo.ListRackRoles(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list rack roles: %w", err)
	}

	return resp, nil
}

// CreateRackRole relays a rack role create to the repository. ValidationError
// is wrapped (via %w) and left intact so callers can recover it with errors.As.
func (s *NetworkService) CreateRackRole(ctx context.Context, in domain.RackRoleWrite) (*domain.RackRole, error) {
	resp, err := s.repo.CreateRackRole(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create rack role: %w", err)
	}

	return resp, nil
}

// UpdateRackRole relays a rack role update to the repository.
func (s *NetworkService) UpdateRackRole(
	ctx context.Context, id int, in domain.RackRoleWrite,
) (*domain.RackRole, error) {
	resp, err := s.repo.UpdateRackRole(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update rack role: %w", err)
	}

	return resp, nil
}

// DeleteRackRole relays a rack role delete to the repository.
func (s *NetworkService) DeleteRackRole(ctx context.Context, id int) error {
	if err := s.repo.DeleteRackRole(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete rack role: %w", err)
	}

	return nil
}

// ListDeviceRoles returns the page of device roles matching params.
func (s *NetworkService) ListDeviceRoles(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.DeviceRole], error) {
	resp, err := s.repo.ListDeviceRoles(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list device roles: %w", err)
	}

	return resp, nil
}

// CreateDeviceRole relays a device role create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover it
// with errors.As.
func (s *NetworkService) CreateDeviceRole(
	ctx context.Context, in domain.DeviceRoleWrite,
) (*domain.DeviceRole, error) {
	resp, err := s.repo.CreateDeviceRole(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create device role: %w", err)
	}

	return resp, nil
}

// UpdateDeviceRole relays a device role update to the repository.
func (s *NetworkService) UpdateDeviceRole(
	ctx context.Context, id int, in domain.DeviceRoleWrite,
) (*domain.DeviceRole, error) {
	resp, err := s.repo.UpdateDeviceRole(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update device role: %w", err)
	}

	return resp, nil
}

// DeleteDeviceRole relays a device role delete to the repository.
func (s *NetworkService) DeleteDeviceRole(ctx context.Context, id int) error {
	if err := s.repo.DeleteDeviceRole(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete device role: %w", err)
	}

	return nil
}

// ListPlatforms returns the page of platforms matching params.
func (s *NetworkService) ListPlatforms(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.Platform], error) {
	resp, err := s.repo.ListPlatforms(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list platforms: %w", err)
	}

	return resp, nil
}

// CreatePlatform relays a platform create to the repository. ValidationError is
// wrapped (via %w) and left intact so callers can recover it with errors.As.
func (s *NetworkService) CreatePlatform(ctx context.Context, in domain.PlatformWrite) (*domain.Platform, error) {
	resp, err := s.repo.CreatePlatform(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create platform: %w", err)
	}

	return resp, nil
}

// UpdatePlatform relays a platform update to the repository.
func (s *NetworkService) UpdatePlatform(
	ctx context.Context, id int, in domain.PlatformWrite,
) (*domain.Platform, error) {
	resp, err := s.repo.UpdatePlatform(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update platform: %w", err)
	}

	return resp, nil
}

// DeletePlatform relays a platform delete to the repository.
func (s *NetworkService) DeletePlatform(ctx context.Context, id int) error {
	if err := s.repo.DeletePlatform(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete platform: %w", err)
	}

	return nil
}

// ListVirtualChassis returns the page of virtual chassis matching params.
func (s *NetworkService) ListVirtualChassis(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.VirtualChassis], error) {
	resp, err := s.repo.ListVirtualChassis(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list virtual chassis: %w", err)
	}

	return resp, nil
}

// CreateVirtualChassis relays a virtual chassis create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover it
// with errors.As.
func (s *NetworkService) CreateVirtualChassis(
	ctx context.Context, in domain.VirtualChassisWrite,
) (*domain.VirtualChassis, error) {
	resp, err := s.repo.CreateVirtualChassis(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create virtual chassis: %w", err)
	}

	return resp, nil
}

// UpdateVirtualChassis relays a virtual chassis update to the repository.
func (s *NetworkService) UpdateVirtualChassis(
	ctx context.Context, id int, in domain.VirtualChassisWrite,
) (*domain.VirtualChassis, error) {
	resp, err := s.repo.UpdateVirtualChassis(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update virtual chassis: %w", err)
	}

	return resp, nil
}

// DeleteVirtualChassis relays a virtual chassis delete to the repository.
func (s *NetworkService) DeleteVirtualChassis(ctx context.Context, id int) error {
	if err := s.repo.DeleteVirtualChassis(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete virtual chassis: %w", err)
	}

	return nil
}

// ListPowerPanels returns the page of power panels matching params.
func (s *NetworkService) ListPowerPanels(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.PowerPanel], error) {
	resp, err := s.repo.ListPowerPanels(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list power panels: %w", err)
	}

	return resp, nil
}

// CreatePowerPanel relays a power panel create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover it
// with errors.As.
func (s *NetworkService) CreatePowerPanel(
	ctx context.Context, in domain.PowerPanelWrite,
) (*domain.PowerPanel, error) {
	resp, err := s.repo.CreatePowerPanel(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create power panel: %w", err)
	}

	return resp, nil
}

// UpdatePowerPanel relays a power panel update to the repository.
func (s *NetworkService) UpdatePowerPanel(
	ctx context.Context, id int, in domain.PowerPanelWrite,
) (*domain.PowerPanel, error) {
	resp, err := s.repo.UpdatePowerPanel(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update power panel: %w", err)
	}

	return resp, nil
}

// DeletePowerPanel relays a power panel delete to the repository.
func (s *NetworkService) DeletePowerPanel(ctx context.Context, id int) error {
	if err := s.repo.DeletePowerPanel(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete power panel: %w", err)
	}

	return nil
}

// ListInventoryItemRoles returns the page of inventory item roles matching
// params.
func (s *NetworkService) ListInventoryItemRoles(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.InventoryItemRole], error) {
	resp, err := s.repo.ListInventoryItemRoles(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list inventory item roles: %w", err)
	}

	return resp, nil
}

// CreateInventoryItemRole relays an inventory item role create to the
// repository. ValidationError is wrapped (via %w) and left intact so callers
// can recover it with errors.As.
func (s *NetworkService) CreateInventoryItemRole(
	ctx context.Context, in domain.InventoryItemRoleWrite,
) (*domain.InventoryItemRole, error) {
	resp, err := s.repo.CreateInventoryItemRole(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create inventory item role: %w", err)
	}

	return resp, nil
}

// UpdateInventoryItemRole relays an inventory item role update to the
// repository.
func (s *NetworkService) UpdateInventoryItemRole(
	ctx context.Context, id int, in domain.InventoryItemRoleWrite,
) (*domain.InventoryItemRole, error) {
	resp, err := s.repo.UpdateInventoryItemRole(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update inventory item role: %w", err)
	}

	return resp, nil
}

// DeleteInventoryItemRole relays an inventory item role delete to the
// repository.
func (s *NetworkService) DeleteInventoryItemRole(ctx context.Context, id int) error {
	if err := s.repo.DeleteInventoryItemRole(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete inventory item role: %w", err)
	}

	return nil
}

// ListFrontPorts returns the page of front ports matching params.
func (s *NetworkService) ListFrontPorts(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.FrontPort], error) {
	resp, err := s.repo.ListFrontPorts(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list front ports: %w", err)
	}

	return resp, nil
}

// CreateFrontPort relays a front port create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover
// it with errors.As.
func (s *NetworkService) CreateFrontPort(ctx context.Context, in domain.FrontPortWrite) (*domain.FrontPort, error) {
	resp, err := s.repo.CreateFrontPort(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create front port: %w", err)
	}

	return resp, nil
}

// UpdateFrontPort relays a front port update to the repository.
func (s *NetworkService) UpdateFrontPort(
	ctx context.Context, id int, in domain.FrontPortWrite,
) (*domain.FrontPort, error) {
	resp, err := s.repo.UpdateFrontPort(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update front port: %w", err)
	}

	return resp, nil
}

// DeleteFrontPort relays a front port delete to the repository.
func (s *NetworkService) DeleteFrontPort(ctx context.Context, id int) error {
	if err := s.repo.DeleteFrontPort(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete front port: %w", err)
	}

	return nil
}

// ListRearPorts returns the page of rear ports matching params.
func (s *NetworkService) ListRearPorts(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.RearPort], error) {
	resp, err := s.repo.ListRearPorts(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list rear ports: %w", err)
	}

	return resp, nil
}

// CreateRearPort relays a rear port create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover
// it with errors.As.
func (s *NetworkService) CreateRearPort(ctx context.Context, in domain.RearPortWrite) (*domain.RearPort, error) {
	resp, err := s.repo.CreateRearPort(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create rear port: %w", err)
	}

	return resp, nil
}

// UpdateRearPort relays a rear port update to the repository.
func (s *NetworkService) UpdateRearPort(
	ctx context.Context, id int, in domain.RearPortWrite,
) (*domain.RearPort, error) {
	resp, err := s.repo.UpdateRearPort(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update rear port: %w", err)
	}

	return resp, nil
}

// DeleteRearPort relays a rear port delete to the repository.
func (s *NetworkService) DeleteRearPort(ctx context.Context, id int) error {
	if err := s.repo.DeleteRearPort(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete rear port: %w", err)
	}

	return nil
}

// ListConsolePorts returns the page of console ports matching params.
func (s *NetworkService) ListConsolePorts(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.ConsolePort], error) {
	resp, err := s.repo.ListConsolePorts(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list console ports: %w", err)
	}

	return resp, nil
}

// CreateConsolePort relays a console port create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover
// it with errors.As.
func (s *NetworkService) CreateConsolePort(
	ctx context.Context, in domain.ConsolePortWrite,
) (*domain.ConsolePort, error) {
	resp, err := s.repo.CreateConsolePort(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create console port: %w", err)
	}

	return resp, nil
}

// UpdateConsolePort relays a console port update to the repository.
func (s *NetworkService) UpdateConsolePort(
	ctx context.Context, id int, in domain.ConsolePortWrite,
) (*domain.ConsolePort, error) {
	resp, err := s.repo.UpdateConsolePort(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update console port: %w", err)
	}

	return resp, nil
}

// DeleteConsolePort relays a console port delete to the repository.
func (s *NetworkService) DeleteConsolePort(ctx context.Context, id int) error {
	if err := s.repo.DeleteConsolePort(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete console port: %w", err)
	}

	return nil
}

// ListConsoleServerPorts returns the page of console server ports matching
// params.
func (s *NetworkService) ListConsoleServerPorts(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.ConsoleServerPort], error) {
	resp, err := s.repo.ListConsoleServerPorts(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list console server ports: %w", err)
	}

	return resp, nil
}

// CreateConsoleServerPort relays a console server port create to the
// repository. ValidationError is wrapped (via %w) and left intact so callers
// can recover it with errors.As.
func (s *NetworkService) CreateConsoleServerPort(
	ctx context.Context, in domain.ConsoleServerPortWrite,
) (*domain.ConsoleServerPort, error) {
	resp, err := s.repo.CreateConsoleServerPort(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create console server port: %w", err)
	}

	return resp, nil
}

// UpdateConsoleServerPort relays a console server port update to the
// repository.
func (s *NetworkService) UpdateConsoleServerPort(
	ctx context.Context, id int, in domain.ConsoleServerPortWrite,
) (*domain.ConsoleServerPort, error) {
	resp, err := s.repo.UpdateConsoleServerPort(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update console server port: %w", err)
	}

	return resp, nil
}

// DeleteConsoleServerPort relays a console server port delete to the
// repository.
func (s *NetworkService) DeleteConsoleServerPort(ctx context.Context, id int) error {
	if err := s.repo.DeleteConsoleServerPort(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete console server port: %w", err)
	}

	return nil
}

// ListPowerPorts returns the page of power ports matching params.
func (s *NetworkService) ListPowerPorts(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.PowerPort], error) {
	resp, err := s.repo.ListPowerPorts(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list power ports: %w", err)
	}

	return resp, nil
}

// CreatePowerPort relays a power port create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover
// it with errors.As.
func (s *NetworkService) CreatePowerPort(ctx context.Context, in domain.PowerPortWrite) (*domain.PowerPort, error) {
	resp, err := s.repo.CreatePowerPort(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create power port: %w", err)
	}

	return resp, nil
}

// UpdatePowerPort relays a power port update to the repository.
func (s *NetworkService) UpdatePowerPort(
	ctx context.Context, id int, in domain.PowerPortWrite,
) (*domain.PowerPort, error) {
	resp, err := s.repo.UpdatePowerPort(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update power port: %w", err)
	}

	return resp, nil
}

// DeletePowerPort relays a power port delete to the repository.
func (s *NetworkService) DeletePowerPort(ctx context.Context, id int) error {
	if err := s.repo.DeletePowerPort(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete power port: %w", err)
	}

	return nil
}

// ListPowerOutlets returns the page of power outlets matching params.
func (s *NetworkService) ListPowerOutlets(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.PowerOutlet], error) {
	resp, err := s.repo.ListPowerOutlets(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list power outlets: %w", err)
	}

	return resp, nil
}

// CreatePowerOutlet relays a power outlet create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover
// it with errors.As.
func (s *NetworkService) CreatePowerOutlet(
	ctx context.Context, in domain.PowerOutletWrite,
) (*domain.PowerOutlet, error) {
	resp, err := s.repo.CreatePowerOutlet(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create power outlet: %w", err)
	}

	return resp, nil
}

// UpdatePowerOutlet relays a power outlet update to the repository.
func (s *NetworkService) UpdatePowerOutlet(
	ctx context.Context, id int, in domain.PowerOutletWrite,
) (*domain.PowerOutlet, error) {
	resp, err := s.repo.UpdatePowerOutlet(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update power outlet: %w", err)
	}

	return resp, nil
}

// DeletePowerOutlet relays a power outlet delete to the repository.
func (s *NetworkService) DeletePowerOutlet(ctx context.Context, id int) error {
	if err := s.repo.DeletePowerOutlet(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete power outlet: %w", err)
	}

	return nil
}

// ListPowerFeeds returns the page of power feeds matching params.
func (s *NetworkService) ListPowerFeeds(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.PowerFeed], error) {
	resp, err := s.repo.ListPowerFeeds(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list power feeds: %w", err)
	}

	return resp, nil
}

// CreatePowerFeed relays a power feed create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover
// it with errors.As.
func (s *NetworkService) CreatePowerFeed(ctx context.Context, in domain.PowerFeedWrite) (*domain.PowerFeed, error) {
	resp, err := s.repo.CreatePowerFeed(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create power feed: %w", err)
	}

	return resp, nil
}

// UpdatePowerFeed relays a power feed update to the repository.
func (s *NetworkService) UpdatePowerFeed(
	ctx context.Context, id int, in domain.PowerFeedWrite,
) (*domain.PowerFeed, error) {
	resp, err := s.repo.UpdatePowerFeed(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update power feed: %w", err)
	}

	return resp, nil
}

// DeletePowerFeed relays a power feed delete to the repository.
func (s *NetworkService) DeletePowerFeed(ctx context.Context, id int) error {
	if err := s.repo.DeletePowerFeed(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete power feed: %w", err)
	}

	return nil
}

// ListDeviceBays returns the page of device bays matching params.
func (s *NetworkService) ListDeviceBays(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.DeviceBay], error) {
	resp, err := s.repo.ListDeviceBays(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list device bays: %w", err)
	}

	return resp, nil
}

// CreateDeviceBay relays a device bay create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover
// it with errors.As.
func (s *NetworkService) CreateDeviceBay(ctx context.Context, in domain.DeviceBayWrite) (*domain.DeviceBay, error) {
	resp, err := s.repo.CreateDeviceBay(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create device bay: %w", err)
	}

	return resp, nil
}

// UpdateDeviceBay relays a device bay update to the repository.
func (s *NetworkService) UpdateDeviceBay(
	ctx context.Context, id int, in domain.DeviceBayWrite,
) (*domain.DeviceBay, error) {
	resp, err := s.repo.UpdateDeviceBay(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update device bay: %w", err)
	}

	return resp, nil
}

// DeleteDeviceBay relays a device bay delete to the repository.
func (s *NetworkService) DeleteDeviceBay(ctx context.Context, id int) error {
	if err := s.repo.DeleteDeviceBay(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete device bay: %w", err)
	}

	return nil
}

// ListRackReservations returns the page of rack reservations matching params.
func (s *NetworkService) ListRackReservations(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.RackReservation], error) {
	resp, err := s.repo.ListRackReservations(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list rack reservations: %w", err)
	}

	return resp, nil
}

// CreateRackReservation relays a rack reservation create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover
// it with errors.As.
func (s *NetworkService) CreateRackReservation(
	ctx context.Context, in domain.RackReservationWrite,
) (*domain.RackReservation, error) {
	resp, err := s.repo.CreateRackReservation(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create rack reservation: %w", err)
	}

	return resp, nil
}

// UpdateRackReservation relays a rack reservation update to the repository.
func (s *NetworkService) UpdateRackReservation(
	ctx context.Context, id int, in domain.RackReservationWrite,
) (*domain.RackReservation, error) {
	resp, err := s.repo.UpdateRackReservation(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update rack reservation: %w", err)
	}

	return resp, nil
}

// DeleteRackReservation relays a rack reservation delete to the repository.
func (s *NetworkService) DeleteRackReservation(ctx context.Context, id int) error {
	if err := s.repo.DeleteRackReservation(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete rack reservation: %w", err)
	}

	return nil
}

// ListInventoryItems returns the page of inventory items matching params.
func (s *NetworkService) ListInventoryItems(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.InventoryItem], error) {
	resp, err := s.repo.ListInventoryItems(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list inventory items: %w", err)
	}

	return resp, nil
}

// CreateInventoryItem relays an inventory item create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover
// it with errors.As.
func (s *NetworkService) CreateInventoryItem(
	ctx context.Context, in domain.InventoryItemWrite,
) (*domain.InventoryItem, error) {
	resp, err := s.repo.CreateInventoryItem(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create inventory item: %w", err)
	}

	return resp, nil
}

// UpdateInventoryItem relays an inventory item update to the repository.
func (s *NetworkService) UpdateInventoryItem(
	ctx context.Context, id int, in domain.InventoryItemWrite,
) (*domain.InventoryItem, error) {
	resp, err := s.repo.UpdateInventoryItem(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update inventory item: %w", err)
	}

	return resp, nil
}

// DeleteInventoryItem relays an inventory item delete to the repository.
func (s *NetworkService) DeleteInventoryItem(ctx context.Context, id int) error {
	if err := s.repo.DeleteInventoryItem(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete inventory item: %w", err)
	}

	return nil
}

// ListMACAddresses returns the page of MAC addresses matching params.
func (s *NetworkService) ListMACAddresses(
	ctx context.Context, params map[string]string,
) (*domain.PaginatedResponse[domain.MACAddress], error) {
	resp, err := s.repo.ListMACAddresses(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list MAC addresses: %w", err)
	}

	return resp, nil
}

// CreateMACAddress relays a MAC address create to the repository.
// ValidationError is wrapped (via %w) and left intact so callers can recover
// it with errors.As.
func (s *NetworkService) CreateMACAddress(
	ctx context.Context, in domain.MACAddressWrite,
) (*domain.MACAddress, error) {
	resp, err := s.repo.CreateMACAddress(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create MAC address: %w", err)
	}

	return resp, nil
}

// UpdateMACAddress relays a MAC address update to the repository.
func (s *NetworkService) UpdateMACAddress(
	ctx context.Context, id int, in domain.MACAddressWrite,
) (*domain.MACAddress, error) {
	resp, err := s.repo.UpdateMACAddress(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update MAC address: %w", err)
	}

	return resp, nil
}

// DeleteMACAddress relays a MAC address delete to the repository.
func (s *NetworkService) DeleteMACAddress(ctx context.Context, id int) error {
	if err := s.repo.DeleteMACAddress(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete MAC address: %w", err)
	}

	return nil
}

// ListRIRs returns the page of RIRs matching params.
func (s *NetworkService) ListRIRs(ctx context.Context, params map[string]string) (
	*domain.PaginatedResponse[domain.RIR], error,
) {
	resp, err := s.repo.ListRIRs(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list RIRs: %w", err)
	}

	return resp, nil
}

// CreateRIR relays an RIR create to the repository.
func (s *NetworkService) CreateRIR(ctx context.Context, in domain.RIRWrite) (*domain.RIR, error) {
	resp, err := s.repo.CreateRIR(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create RIR: %w", err)
	}

	return resp, nil
}

// UpdateRIR relays an RIR update to the repository.
func (s *NetworkService) UpdateRIR(ctx context.Context, id int, in domain.RIRWrite) (*domain.RIR, error) {
	resp, err := s.repo.UpdateRIR(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update RIR: %w", err)
	}

	return resp, nil
}

// DeleteRIR relays an RIR delete to the repository.
func (s *NetworkService) DeleteRIR(ctx context.Context, id int) error {
	if err := s.repo.DeleteRIR(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete RIR: %w", err)
	}

	return nil
}

// ListAggregates returns the page of aggregates matching params.
func (s *NetworkService) ListAggregates(ctx context.Context, params map[string]string) (
	*domain.PaginatedResponse[domain.Aggregate], error,
) {
	resp, err := s.repo.ListAggregates(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list aggregates: %w", err)
	}

	return resp, nil
}

// CreateAggregate relays an aggregate create to the repository.
func (s *NetworkService) CreateAggregate(ctx context.Context, in domain.AggregateWrite) (*domain.Aggregate, error) {
	resp, err := s.repo.CreateAggregate(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create aggregate: %w", err)
	}

	return resp, nil
}

// UpdateAggregate relays an aggregate update to the repository.
func (s *NetworkService) UpdateAggregate(
	ctx context.Context, id int, in domain.AggregateWrite,
) (*domain.Aggregate, error) {
	resp, err := s.repo.UpdateAggregate(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update aggregate: %w", err)
	}

	return resp, nil
}

// DeleteAggregate relays an aggregate delete to the repository.
func (s *NetworkService) DeleteAggregate(ctx context.Context, id int) error {
	if err := s.repo.DeleteAggregate(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete aggregate: %w", err)
	}

	return nil
}

// ListIPRanges returns the page of IP ranges matching params.
func (s *NetworkService) ListIPRanges(ctx context.Context, params map[string]string) (
	*domain.PaginatedResponse[domain.IPRange], error,
) {
	resp, err := s.repo.ListIPRanges(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list IP ranges: %w", err)
	}

	return resp, nil
}

// CreateIPRange relays an IP range create to the repository.
func (s *NetworkService) CreateIPRange(ctx context.Context, in domain.IPRangeWrite) (*domain.IPRange, error) {
	resp, err := s.repo.CreateIPRange(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create IP range: %w", err)
	}

	return resp, nil
}

// UpdateIPRange relays an IP range update to the repository.
func (s *NetworkService) UpdateIPRange(
	ctx context.Context, id int, in domain.IPRangeWrite,
) (*domain.IPRange, error) {
	resp, err := s.repo.UpdateIPRange(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update IP range: %w", err)
	}

	return resp, nil
}

// DeleteIPRange relays an IP range delete to the repository.
func (s *NetworkService) DeleteIPRange(ctx context.Context, id int) error {
	if err := s.repo.DeleteIPRange(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete IP range: %w", err)
	}

	return nil
}

// ListASNs returns the page of ASNs matching params.
func (s *NetworkService) ListASNs(ctx context.Context, params map[string]string) (
	*domain.PaginatedResponse[domain.ASN], error,
) {
	resp, err := s.repo.ListASNs(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list ASNs: %w", err)
	}

	return resp, nil
}

// CreateASN relays an ASN create to the repository.
func (s *NetworkService) CreateASN(ctx context.Context, in domain.ASNWrite) (*domain.ASN, error) {
	resp, err := s.repo.CreateASN(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create ASN: %w", err)
	}

	return resp, nil
}

// UpdateASN relays an ASN update to the repository.
func (s *NetworkService) UpdateASN(ctx context.Context, id int, in domain.ASNWrite) (*domain.ASN, error) {
	resp, err := s.repo.UpdateASN(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update ASN: %w", err)
	}

	return resp, nil
}

// DeleteASN relays an ASN delete to the repository.
func (s *NetworkService) DeleteASN(ctx context.Context, id int) error {
	if err := s.repo.DeleteASN(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete ASN: %w", err)
	}

	return nil
}

// ListASNRanges returns the page of ASN ranges matching params.
func (s *NetworkService) ListASNRanges(ctx context.Context, params map[string]string) (
	*domain.PaginatedResponse[domain.ASNRange], error,
) {
	resp, err := s.repo.ListASNRanges(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list ASN ranges: %w", err)
	}

	return resp, nil
}

// CreateASNRange relays an ASN range create to the repository.
func (s *NetworkService) CreateASNRange(ctx context.Context, in domain.ASNRangeWrite) (*domain.ASNRange, error) {
	resp, err := s.repo.CreateASNRange(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create ASN range: %w", err)
	}

	return resp, nil
}

// UpdateASNRange relays an ASN range update to the repository.
func (s *NetworkService) UpdateASNRange(
	ctx context.Context, id int, in domain.ASNRangeWrite,
) (*domain.ASNRange, error) {
	resp, err := s.repo.UpdateASNRange(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update ASN range: %w", err)
	}

	return resp, nil
}

// DeleteASNRange relays an ASN range delete to the repository.
func (s *NetworkService) DeleteASNRange(ctx context.Context, id int) error {
	if err := s.repo.DeleteASNRange(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete ASN range: %w", err)
	}

	return nil
}

// ListRouteTargets returns the page of route targets matching params.
func (s *NetworkService) ListRouteTargets(ctx context.Context, params map[string]string) (
	*domain.PaginatedResponse[domain.RouteTarget], error,
) {
	resp, err := s.repo.ListRouteTargets(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list route targets: %w", err)
	}

	return resp, nil
}

// CreateRouteTarget relays a route target create to the repository.
func (s *NetworkService) CreateRouteTarget(
	ctx context.Context, in domain.RouteTargetWrite,
) (*domain.RouteTarget, error) {
	resp, err := s.repo.CreateRouteTarget(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create route target: %w", err)
	}

	return resp, nil
}

// UpdateRouteTarget relays a route target update to the repository.
func (s *NetworkService) UpdateRouteTarget(
	ctx context.Context, id int, in domain.RouteTargetWrite,
) (*domain.RouteTarget, error) {
	resp, err := s.repo.UpdateRouteTarget(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update route target: %w", err)
	}

	return resp, nil
}

// DeleteRouteTarget relays a route target delete to the repository.
func (s *NetworkService) DeleteRouteTarget(ctx context.Context, id int) error {
	if err := s.repo.DeleteRouteTarget(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete route target: %w", err)
	}

	return nil
}

// ListServices returns the page of services matching params.
func (s *NetworkService) ListServices(ctx context.Context, params map[string]string) (
	*domain.PaginatedResponse[domain.Service], error,
) {
	resp, err := s.repo.ListServices(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list services: %w", err)
	}

	return resp, nil
}

// CreateService relays a service create to the repository.
func (s *NetworkService) CreateService(ctx context.Context, in domain.ServiceWrite) (*domain.Service, error) {
	resp, err := s.repo.CreateService(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create service: %w", err)
	}

	return resp, nil
}

// UpdateService relays a service update to the repository.
func (s *NetworkService) UpdateService(
	ctx context.Context, id int, in domain.ServiceWrite,
) (*domain.Service, error) {
	resp, err := s.repo.UpdateService(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update service: %w", err)
	}

	return resp, nil
}

// DeleteService relays a service delete to the repository.
func (s *NetworkService) DeleteService(ctx context.Context, id int) error {
	if err := s.repo.DeleteService(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete service: %w", err)
	}

	return nil
}

// ListServiceTemplates returns the page of service templates matching params.
func (s *NetworkService) ListServiceTemplates(ctx context.Context, params map[string]string) (
	*domain.PaginatedResponse[domain.ServiceTemplate], error,
) {
	resp, err := s.repo.ListServiceTemplates(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list service templates: %w", err)
	}

	return resp, nil
}

// CreateServiceTemplate relays a service template create to the repository.
func (s *NetworkService) CreateServiceTemplate(
	ctx context.Context, in domain.ServiceTemplateWrite,
) (*domain.ServiceTemplate, error) {
	resp, err := s.repo.CreateServiceTemplate(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create service template: %w", err)
	}

	return resp, nil
}

// UpdateServiceTemplate relays a service template update to the repository.
func (s *NetworkService) UpdateServiceTemplate(
	ctx context.Context, id int, in domain.ServiceTemplateWrite,
) (*domain.ServiceTemplate, error) {
	resp, err := s.repo.UpdateServiceTemplate(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update service template: %w", err)
	}

	return resp, nil
}

// DeleteServiceTemplate relays a service template delete to the repository.
func (s *NetworkService) DeleteServiceTemplate(ctx context.Context, id int) error {
	if err := s.repo.DeleteServiceTemplate(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete service template: %w", err)
	}

	return nil
}

// ListFHRPGroups returns the page of FHRP groups matching params.
func (s *NetworkService) ListFHRPGroups(ctx context.Context, params map[string]string) (
	*domain.PaginatedResponse[domain.FHRPGroup], error,
) {
	resp, err := s.repo.ListFHRPGroups(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list FHRP groups: %w", err)
	}

	return resp, nil
}

// CreateFHRPGroup relays an FHRP group create to the repository.
func (s *NetworkService) CreateFHRPGroup(ctx context.Context, in domain.FHRPGroupWrite) (*domain.FHRPGroup, error) {
	resp, err := s.repo.CreateFHRPGroup(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create FHRP group: %w", err)
	}

	return resp, nil
}

// UpdateFHRPGroup relays an FHRP group update to the repository.
func (s *NetworkService) UpdateFHRPGroup(
	ctx context.Context, id int, in domain.FHRPGroupWrite,
) (*domain.FHRPGroup, error) {
	resp, err := s.repo.UpdateFHRPGroup(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update FHRP group: %w", err)
	}

	return resp, nil
}

// DeleteFHRPGroup relays an FHRP group delete to the repository.
func (s *NetworkService) DeleteFHRPGroup(ctx context.Context, id int) error {
	if err := s.repo.DeleteFHRPGroup(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete FHRP group: %w", err)
	}

	return nil
}

// ListFHRPGroupAssignments returns the page of FHRP group assignments matching
// params.
func (s *NetworkService) ListFHRPGroupAssignments(ctx context.Context, params map[string]string) (
	*domain.PaginatedResponse[domain.FHRPGroupAssignment], error,
) {
	resp, err := s.repo.ListFHRPGroupAssignments(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list FHRP group assignments: %w", err)
	}

	return resp, nil
}

// CreateFHRPGroupAssignment relays an FHRP group assignment create to the
// repository.
func (s *NetworkService) CreateFHRPGroupAssignment(
	ctx context.Context, in domain.FHRPGroupAssignmentWrite,
) (*domain.FHRPGroupAssignment, error) {
	resp, err := s.repo.CreateFHRPGroupAssignment(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create FHRP group assignment: %w", err)
	}

	return resp, nil
}

// UpdateFHRPGroupAssignment relays an FHRP group assignment update to the
// repository.
func (s *NetworkService) UpdateFHRPGroupAssignment(
	ctx context.Context, id int, in domain.FHRPGroupAssignmentWrite,
) (*domain.FHRPGroupAssignment, error) {
	resp, err := s.repo.UpdateFHRPGroupAssignment(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update FHRP group assignment: %w", err)
	}

	return resp, nil
}

// DeleteFHRPGroupAssignment relays an FHRP group assignment delete to the
// repository.
func (s *NetworkService) DeleteFHRPGroupAssignment(ctx context.Context, id int) error {
	if err := s.repo.DeleteFHRPGroupAssignment(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete FHRP group assignment: %w", err)
	}

	return nil
}

// ListVLANTranslationPolicies returns the page of VLAN translation policies
// matching params.
func (s *NetworkService) ListVLANTranslationPolicies(ctx context.Context, params map[string]string) (
	*domain.PaginatedResponse[domain.VLANTranslationPolicy], error,
) {
	resp, err := s.repo.ListVLANTranslationPolicies(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list VLAN translation policies: %w", err)
	}

	return resp, nil
}

// CreateVLANTranslationPolicy relays a VLAN translation policy create to the
// repository.
func (s *NetworkService) CreateVLANTranslationPolicy(
	ctx context.Context, in domain.VLANTranslationPolicyWrite,
) (*domain.VLANTranslationPolicy, error) {
	resp, err := s.repo.CreateVLANTranslationPolicy(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create VLAN translation policy: %w", err)
	}

	return resp, nil
}

// UpdateVLANTranslationPolicy relays a VLAN translation policy update to the
// repository.
func (s *NetworkService) UpdateVLANTranslationPolicy(
	ctx context.Context, id int, in domain.VLANTranslationPolicyWrite,
) (*domain.VLANTranslationPolicy, error) {
	resp, err := s.repo.UpdateVLANTranslationPolicy(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update VLAN translation policy: %w", err)
	}

	return resp, nil
}

// DeleteVLANTranslationPolicy relays a VLAN translation policy delete to the
// repository.
func (s *NetworkService) DeleteVLANTranslationPolicy(ctx context.Context, id int) error {
	if err := s.repo.DeleteVLANTranslationPolicy(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete VLAN translation policy: %w", err)
	}

	return nil
}

// ListVLANTranslationRules returns the page of VLAN translation rules matching
// params.
func (s *NetworkService) ListVLANTranslationRules(ctx context.Context, params map[string]string) (
	*domain.PaginatedResponse[domain.VLANTranslationRule], error,
) {
	resp, err := s.repo.ListVLANTranslationRules(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list VLAN translation rules: %w", err)
	}

	return resp, nil
}

// CreateVLANTranslationRule relays a VLAN translation rule create to the
// repository.
func (s *NetworkService) CreateVLANTranslationRule(
	ctx context.Context, in domain.VLANTranslationRuleWrite,
) (*domain.VLANTranslationRule, error) {
	resp, err := s.repo.CreateVLANTranslationRule(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create VLAN translation rule: %w", err)
	}

	return resp, nil
}

// UpdateVLANTranslationRule relays a VLAN translation rule update to the
// repository.
func (s *NetworkService) UpdateVLANTranslationRule(
	ctx context.Context, id int, in domain.VLANTranslationRuleWrite,
) (*domain.VLANTranslationRule, error) {
	resp, err := s.repo.UpdateVLANTranslationRule(ctx, s.token.Value(), id, in)
	if err != nil {
		return nil, fmt.Errorf("update VLAN translation rule: %w", err)
	}

	return resp, nil
}

// DeleteVLANTranslationRule relays a VLAN translation rule delete to the
// repository.
func (s *NetworkService) DeleteVLANTranslationRule(ctx context.Context, id int) error {
	if err := s.repo.DeleteVLANTranslationRule(ctx, s.token.Value(), id); err != nil {
		return fmt.Errorf("delete VLAN translation rule: %w", err)
	}

	return nil
}
