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
