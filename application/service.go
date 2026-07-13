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

func (s *NetworkService) GetObject(ctx context.Context, objectType string, id int) (domain.RawObject, error) {
	resp, err := s.repo.GetObject(ctx, s.token.Value(), objectType, id, nil)
	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}
	return resp, nil
}
