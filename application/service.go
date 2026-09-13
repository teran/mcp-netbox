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

func (s *NetworkService) ListSites(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Site], error) {
	resp, err := s.repo.ListSites(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list sites: %w", err)
	}
	return resp, nil
}

func (s *NetworkService) ListDevices(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Device], error) {
	resp, err := s.repo.ListDevices(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list devices: %w", err)
	}
	return resp, nil
}

func (s *NetworkService) ListIPAddresses(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.IPAddress], error) {
	resp, err := s.repo.ListIPAddresses(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list IP addresses: %w", err)
	}
	return resp, nil
}

func (s *NetworkService) ListPrefixes(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Prefix], error) {
	resp, err := s.repo.ListPrefixes(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list prefixes: %w", err)
	}
	return resp, nil
}

func (s *NetworkService) ListVLANs(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.VLAN], error) {
	resp, err := s.repo.ListVLANs(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list VLANs: %w", err)
	}
	return resp, nil
}

func (s *NetworkService) ListVirtualMachines(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.VirtualMachine], error) {
	resp, err := s.repo.ListVirtualMachines(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list virtual machines: %w", err)
	}
	return resp, nil
}

func (s *NetworkService) ListClusters(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Cluster], error) {
	resp, err := s.repo.ListClusters(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list clusters: %w", err)
	}
	return resp, nil
}

func (s *NetworkService) ListCircuits(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Circuit], error) {
	resp, err := s.repo.ListCircuits(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list circuits: %w", err)
	}
	return resp, nil
}

func (s *NetworkService) ListInterfaces(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Interface], error) {
	resp, err := s.repo.ListInterfaces(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list interfaces: %w", err)
	}
	return resp, nil
}

func (s *NetworkService) ListVMInterfaces(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.VMInterface], error) {
	resp, err := s.repo.ListVMInterfaces(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list VM interfaces: %w", err)
	}
	return resp, nil
}

func (s *NetworkService) ListRacks(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Rack], error) {
	resp, err := s.repo.ListRacks(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list racks: %w", err)
	}
	return resp, nil
}

func (s *NetworkService) ListCircuitTerminations(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.CircuitTermination], error) {
	resp, err := s.repo.ListCircuitTerminations(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list circuit terminations: %w", err)
	}
	return resp, nil
}

func (s *NetworkService) ListCables(ctx context.Context, params map[string]string) (*domain.PaginatedResponse[domain.Cable], error) {
	resp, err := s.repo.ListCables(ctx, s.token.Value(), params)
	if err != nil {
		return nil, fmt.Errorf("list cables: %w", err)
	}
	return resp, nil
}

func (s *NetworkService) GetObject(ctx context.Context, objectType string, id int, params map[string]string) (domain.RawObject, error) {
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
func (s *NetworkService) UpdateIPAddress(ctx context.Context, id int, in domain.IPAddressWrite) (*domain.IPAddress, error) {
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
func (s *NetworkService) CreateVirtualMachine(ctx context.Context, in domain.VirtualMachineWrite) (*domain.VirtualMachine, error) {
	resp, err := s.repo.CreateVirtualMachine(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create virtual machine: %w", err)
	}
	return resp, nil
}

// UpdateVirtualMachine relays a virtual machine update to the repository.
func (s *NetworkService) UpdateVirtualMachine(ctx context.Context, id int, in domain.VirtualMachineWrite) (*domain.VirtualMachine, error) {
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
func (s *NetworkService) UpdateInterface(ctx context.Context, id int, in domain.InterfaceWrite) (*domain.Interface, error) {
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
func (s *NetworkService) CreateCircuitTermination(ctx context.Context, in domain.CircuitTerminationWrite) (*domain.CircuitTermination, error) {
	resp, err := s.repo.CreateCircuitTermination(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create circuit termination: %w", err)
	}
	return resp, nil
}

// UpdateCircuitTermination relays a circuit termination update to the
// repository.
func (s *NetworkService) UpdateCircuitTermination(ctx context.Context, id int, in domain.CircuitTerminationWrite) (*domain.CircuitTermination, error) {
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
func (s *NetworkService) CreateVMInterface(ctx context.Context, in domain.VMInterfaceWrite) (*domain.VMInterface, error) {
	resp, err := s.repo.CreateVMInterface(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create vm interface: %w", err)
	}
	return resp, nil
}

// UpdateVMInterface relays a VM interface update to the repository.
func (s *NetworkService) UpdateVMInterface(ctx context.Context, id int, in domain.VMInterfaceWrite) (*domain.VMInterface, error) {
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
func (s *NetworkService) UpdateProvider(ctx context.Context, id int, in domain.ProviderWrite) (*domain.Provider, error) {
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
func (s *NetworkService) CreateManufacturer(ctx context.Context, in domain.ManufacturerWrite) (*domain.Manufacturer, error) {
	resp, err := s.repo.CreateManufacturer(ctx, s.token.Value(), in)
	if err != nil {
		return nil, fmt.Errorf("create manufacturer: %w", err)
	}
	return resp, nil
}

// UpdateManufacturer relays a manufacturer update to the repository.
func (s *NetworkService) UpdateManufacturer(ctx context.Context, id int, in domain.ManufacturerWrite) (*domain.Manufacturer, error) {
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
func (s *NetworkService) UpdateDeviceType(ctx context.Context, id int, in domain.DeviceTypeWrite) (*domain.DeviceType, error) {
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
func (s *NetworkService) UpdateLocation(ctx context.Context, id int, in domain.LocationWrite) (*domain.Location, error) {
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
