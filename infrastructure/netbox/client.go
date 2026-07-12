package netbox

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/teran/mcp-netbox/domain"
)

var objectTypeToEndpoint = map[string]string{
	"site":            "/api/dcim/sites/",
	"device":          "/api/dcim/devices/",
	"prefix":          "/api/ipam/prefixes/",
	"ip_address":      "/api/ipam/ip-addresses/",
	"vlan":            "/api/ipam/vlans/",
	"virtual_machine": "/api/virtualization/virtual-machines/",
	"cluster":         "/api/virtualization/clusters/",
	"circuit":         "/api/circuits/circuits/",
	"provider":        "/api/circuits/providers/",
	"tenant":          "/api/tenancy/tenants/",
	"rack":            "/api/dcim/racks/",
	"manufacturer":    "/api/dcim/manufacturers/",
	"device_type":     "/api/dcim/device-types/",
	"location":        "/api/dcim/locations/",
	"cluster_type":    "/api/virtualization/cluster-types/",
	"cluster_group":   "/api/virtualization/cluster-groups/",
	"circuit_type":    "/api/circuits/circuit-types/",
	"vrf":             "/api/ipam/vrfs/",
	"vlan_group":      "/api/ipam/vlan-groups/",
	"role":            "/api/ipam/roles/",
	"contact":         "/api/tenancy/contacts/",
	"cable":           "/api/dcim/cables/",
}

// Client is an HTTP client for the NetBox API.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new NetBox API client.
func NewClient(baseURL string, httpClient *http.Client) *Client {
	return &Client{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		httpClient: httpClient,
	}
}

func (c *Client) doRequest(ctx context.Context, token, method, path string, params map[string]string) ([]byte, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, method, u.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized: invalid token or insufficient permissions")
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("not found")
	}
	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("forbidden: token lacks required permissions")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 100*1024*1024)) // 100 MB limit
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func (c *Client) get(ctx context.Context, token, path string, params map[string]string) ([]byte, error) {
	return c.doRequest(ctx, token, http.MethodGet, path, params)
}

// ListSites returns a paginated list of sites.
func (c *Client) ListSites(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Site], error) {
	raw, err := c.get(ctx, token, "/api/dcim/sites/", params)
	if err != nil {
		return nil, err
	}

	var wireResp WireJSONPaginated[WireSite]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sites: %w", err)
	}

	sites := make([]domain.Site, len(wireResp.Results))
	for i := range wireResp.Results {
		sites[i] = wireSiteToDomain(wireResp.Results[i])
	}

	return &domain.PaginatedResponse[domain.Site]{
		Count:    wireResp.Count,
		Next:     wireResp.Next,
		Previous: wireResp.Previous,
		Results:  sites,
	}, nil
}

// ListDevices returns a paginated list of devices.
func (c *Client) ListDevices(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Device], error) {
	raw, err := c.get(ctx, token, "/api/dcim/devices/", params)
	if err != nil {
		return nil, err
	}

	var wireResp WireJSONPaginated[WireDevice]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal devices: %w", err)
	}

	devices := make([]domain.Device, len(wireResp.Results))
	for i := range wireResp.Results {
		devices[i] = wireDeviceToDomain(wireResp.Results[i])
	}

	return &domain.PaginatedResponse[domain.Device]{
		Count:    wireResp.Count,
		Next:     wireResp.Next,
		Previous: wireResp.Previous,
		Results:  devices,
	}, nil
}

// ListIPAddresses returns a paginated list of IP addresses.
func (c *Client) ListIPAddresses(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.IPAddress], error) {
	raw, err := c.get(ctx, token, "/api/ipam/ip-addresses/", params)
	if err != nil {
		return nil, err
	}

	var wireResp WireJSONPaginated[WireIPAddress]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal IP addresses: %w", err)
	}

	ips := make([]domain.IPAddress, len(wireResp.Results))
	for i := range wireResp.Results {
		ips[i] = wireIPAddressToDomain(wireResp.Results[i])
	}

	return &domain.PaginatedResponse[domain.IPAddress]{
		Count:    wireResp.Count,
		Next:     wireResp.Next,
		Previous: wireResp.Previous,
		Results:  ips,
	}, nil
}

// ListPrefixes returns a paginated list of prefixes.
func (c *Client) ListPrefixes(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Prefix], error) {
	raw, err := c.get(ctx, token, "/api/ipam/prefixes/", params)
	if err != nil {
		return nil, err
	}

	var wireResp WireJSONPaginated[WirePrefix]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal prefixes: %w", err)
	}

	prefixes := make([]domain.Prefix, len(wireResp.Results))
	for i := range wireResp.Results {
		prefixes[i] = wirePrefixToDomain(wireResp.Results[i])
	}

	return &domain.PaginatedResponse[domain.Prefix]{
		Count:    wireResp.Count,
		Next:     wireResp.Next,
		Previous: wireResp.Previous,
		Results:  prefixes,
	}, nil
}

// ListVLANs returns a paginated list of VLANs.
func (c *Client) ListVLANs(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VLAN], error) {
	raw, err := c.get(ctx, token, "/api/ipam/vlans/", params)
	if err != nil {
		return nil, err
	}

	var wireResp WireJSONPaginated[WireVLAN]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal VLANs: %w", err)
	}

	vlans := make([]domain.VLAN, len(wireResp.Results))
	for i := range wireResp.Results {
		vlans[i] = wireVLANToDomain(wireResp.Results[i])
	}

	return &domain.PaginatedResponse[domain.VLAN]{
		Count:    wireResp.Count,
		Next:     wireResp.Next,
		Previous: wireResp.Previous,
		Results:  vlans,
	}, nil
}

// ListVirtualMachines returns a paginated list of virtual machines.
func (c *Client) ListVirtualMachines(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VirtualMachine], error) {
	raw, err := c.get(ctx, token, "/api/virtualization/virtual-machines/", params)
	if err != nil {
		return nil, err
	}

	var wireResp WireJSONPaginated[WireVirtualMachine]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal VMs: %w", err)
	}

	vms := make([]domain.VirtualMachine, len(wireResp.Results))
	for i := range wireResp.Results {
		vms[i] = wireVirtualMachineToDomain(wireResp.Results[i])
	}

	return &domain.PaginatedResponse[domain.VirtualMachine]{
		Count:    wireResp.Count,
		Next:     wireResp.Next,
		Previous: wireResp.Previous,
		Results:  vms,
	}, nil
}

// ListClusters returns a paginated list of clusters.
func (c *Client) ListClusters(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Cluster], error) {
	raw, err := c.get(ctx, token, "/api/virtualization/clusters/", params)
	if err != nil {
		return nil, err
	}

	var wireResp WireJSONPaginated[WireCluster]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal clusters: %w", err)
	}

	clusters := make([]domain.Cluster, len(wireResp.Results))
	for i := range wireResp.Results {
		clusters[i] = wireClusterToDomain(wireResp.Results[i])
	}

	return &domain.PaginatedResponse[domain.Cluster]{
		Count:    wireResp.Count,
		Next:     wireResp.Next,
		Previous: wireResp.Previous,
		Results:  clusters,
	}, nil
}

// ListCircuits returns a paginated list of circuits.
func (c *Client) ListCircuits(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Circuit], error) {
	raw, err := c.get(ctx, token, "/api/circuits/circuits/", params)
	if err != nil {
		return nil, err
	}

	var wireResp WireJSONPaginated[WireCircuit]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal circuits: %w", err)
	}

	circuits := make([]domain.Circuit, len(wireResp.Results))
	for i := range wireResp.Results {
		circuits[i] = wireCircuitToDomain(wireResp.Results[i])
	}

	return &domain.PaginatedResponse[domain.Circuit]{
		Count:    wireResp.Count,
		Next:     wireResp.Next,
		Previous: wireResp.Previous,
		Results:  circuits,
	}, nil
}

// ListRacks returns a paginated list of racks.
func (c *Client) ListRacks(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Rack], error) {
	raw, err := c.get(ctx, token, "/api/dcim/racks/", params)
	if err != nil {
		return nil, err
	}

	var wireResp WireJSONPaginated[WireRack]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal racks: %w", err)
	}

	racks := make([]domain.Rack, len(wireResp.Results))
	for i := range wireResp.Results {
		racks[i] = wireRackToDomain(wireResp.Results[i])
	}

	return &domain.PaginatedResponse[domain.Rack]{
		Count:    wireResp.Count,
		Next:     wireResp.Next,
		Previous: wireResp.Previous,
		Results:  racks,
	}, nil
}

// GetObject retrieves a single object by its type and ID.
func (c *Client) GetObject(ctx context.Context, token, objectType string, id int, params map[string]string) (interface{}, error) {
	endpoint, ok := objectTypeToEndpoint[objectType]
	if !ok {
		return nil, fmt.Errorf("unknown object type: %s", objectType)
	}

	path := fmt.Sprintf("%s%d/", endpoint, id)
	body, err := c.get(ctx, token, path, params)
	if err != nil {
		return nil, err
	}

	// Return raw JSON for generic object retrieval
	var raw interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return raw, nil
}

// Ensure Client implements domain.NetworkRepository.
var _ domain.NetworkRepository = (*Client)(nil)

// — conversion functions —

func wireNestedToDomain(w *WireNested) *domain.Nested {
	if w == nil {
		return nil
	}
	return &domain.Nested{
		ID:   w.ID,
		URL:  w.URL,
		Name: w.Name,
		Slug: w.Slug,
	}
}

func wireLabelToDomain(w *WireLabel) *domain.Label {
	if w == nil {
		return nil
	}
	return &domain.Label{
		Value: w.Value,
		Label: w.Label,
	}
}

func wireTagsToDomain(tags []WireTag) []domain.Tag {
	if tags == nil {
		return nil
	}
	result := make([]domain.Tag, len(tags))
	for i := range tags {
		result[i] = domain.Tag{
			ID:   tags[i].ID,
			Name: tags[i].Name,
			Slug: tags[i].Slug,
			URL:  tags[i].URL,
		}
	}
	return result
}

func wireSiteToDomain(w WireSite) domain.Site {
	return domain.Site{
		ID:              w.ID,
		Name:            w.Name,
		Slug:            w.Slug,
		Status:          wireLabelToDomain(w.Status),
		Region:          wireNestedToDomain(w.Region),
		Tenant:          wireNestedToDomain(w.Tenant),
		Facility:        w.Facility,
		TimeZone:        w.TimeZone,
		Description:     w.Description,
		PhysicalAddress: w.PhysicalAddress,
		ShippingAddress: w.ShippingAddress,
		Comments:        w.Comments,
		Tags:            wireTagsToDomain(w.Tags),
		Created:         w.Created,
		LastUpdated:     w.LastUpdated,
	}
}

func wireDeviceToDomain(w WireDevice) domain.Device {
	return domain.Device{
		ID:          w.ID,
		Name:        w.Name,
		DeviceType:  wireNestedToDomain(w.DeviceType),
		Role:        wireNestedToDomain(w.Role),
		Tenant:      wireNestedToDomain(w.Tenant),
		Platform:    wireNestedToDomain(w.Platform),
		Serial:      w.Serial,
		AssetTag:    w.AssetTag,
		Site:        wireNestedToDomain(w.Site),
		Rack:        wireNestedToDomain(w.Rack),
		Position:    w.Position,
		Face:        wireLabelToDomain(w.Face),
		Status:      wireLabelToDomain(w.Status),
		Cluster:     wireNestedToDomain(w.Cluster),
		Comments:    w.Comments,
		Tags:        wireTagsToDomain(w.Tags),
		Created:     w.Created,
		LastUpdated: w.LastUpdated,
	}
}

func wireAssignedObjToDomain(w *WireAssignedObj) *domain.AssignedObject {
	if w == nil {
		return nil
	}
	return &domain.AssignedObject{
		ID:     w.ID,
		URL:    w.URL,
		Device: wireNestedToDomain(w.Device),
		Name:   w.Name,
	}
}

func wireIPAddressToDomain(w WireIPAddress) domain.IPAddress {
	return domain.IPAddress{
		ID:                 w.ID,
		Address:            w.Address,
		VRF:                wireNestedToDomain(w.VRF),
		Tenant:             wireNestedToDomain(w.Tenant),
		Status:             wireLabelToDomain(w.Status),
		Role:               wireLabelToDomain(w.Role),
		AssignedObjectType: w.AssignedObjectType,
		AssignedObjectID:   w.AssignedObjectID,
		AssignedObject:     wireAssignedObjToDomain(w.AssignedObject),
		DNSName:            w.DNSName,
		Description:        w.Description,
		Comments:           w.Comments,
		Tags:               wireTagsToDomain(w.Tags),
		Created:            w.Created,
		LastUpdated:        w.LastUpdated,
	}
}

func wirePrefixToDomain(w WirePrefix) domain.Prefix {
	family := (*domain.Family)(nil)
	if w.Family != nil {
		family = &domain.Family{
			Value: w.Family.Value,
			Label: w.Family.Label,
		}
	}
	return domain.Prefix{
		ID:          w.ID,
		Prefix:      w.Prefix,
		Site:        wireNestedToDomain(w.Site),
		VRF:         wireNestedToDomain(w.VRF),
		Tenant:      wireNestedToDomain(w.Tenant),
		VLAN:        wireNestedToDomain(w.VLAN),
		Status:      wireLabelToDomain(w.Status),
		Role:        wireNestedToDomain(w.Role),
		IsPool:      w.IsPool,
		Description: w.Description,
		Comments:    w.Comments,
		Tags:        wireTagsToDomain(w.Tags),
		Created:     w.Created,
		LastUpdated: w.LastUpdated,
		Children:    w.Children,
		Depth:       w.Depth,
		Family:      family,
	}
}

func wireVLANToDomain(w WireVLAN) domain.VLAN {
	return domain.VLAN{
		ID:          w.ID,
		Site:        wireNestedToDomain(w.Site),
		Group:       wireNestedToDomain(w.Group),
		VID:         w.VID,
		Name:        w.Name,
		Tenant:      wireNestedToDomain(w.Tenant),
		Status:      wireLabelToDomain(w.Status),
		Role:        wireNestedToDomain(w.Role),
		Description: w.Description,
		Comments:    w.Comments,
		Tags:        wireTagsToDomain(w.Tags),
		Created:     w.Created,
		LastUpdated: w.LastUpdated,
	}
}

func wireVirtualMachineToDomain(w WireVirtualMachine) domain.VirtualMachine {
	return domain.VirtualMachine{
		ID:          w.ID,
		Name:        w.Name,
		Cluster:     wireNestedToDomain(w.Cluster),
		Role:        wireNestedToDomain(w.Role),
		Tenant:      wireNestedToDomain(w.Tenant),
		Platform:    wireNestedToDomain(w.Platform),
		Status:      wireLabelToDomain(w.Status),
		Site:        wireNestedToDomain(w.Site),
		VCPUs:       w.VCPUs,
		Memory:      w.Memory,
		Disk:        w.Disk,
		Comments:    w.Comments,
		Tags:        wireTagsToDomain(w.Tags),
		Created:     w.Created,
		LastUpdated: w.LastUpdated,
	}
}

func wireClusterToDomain(w WireCluster) domain.Cluster {
	return domain.Cluster{
		ID:           w.ID,
		Name:         w.Name,
		ClusterType:  wireNestedToDomain(w.ClusterType),
		ClusterGroup: wireNestedToDomain(w.ClusterGroup),
		Site:         wireNestedToDomain(w.Site),
		Tenant:       wireNestedToDomain(w.Tenant),
		Description:  w.Description,
		Comments:     w.Comments,
		Tags:         wireTagsToDomain(w.Tags),
		Created:      w.Created,
		LastUpdated:  w.LastUpdated,
	}
}

func wireCircuitToDomain(w WireCircuit) domain.Circuit {
	return domain.Circuit{
		ID:          w.ID,
		CID:         w.CID,
		Provider:    wireNestedToDomain(w.Provider),
		CircuitType: wireNestedToDomain(w.CircuitType),
		Tenant:      wireNestedToDomain(w.Tenant),
		Status:      wireLabelToDomain(w.Status),
		Description: w.Description,
		Comments:    w.Comments,
		Tags:        wireTagsToDomain(w.Tags),
		Created:     w.Created,
		LastUpdated: w.LastUpdated,
		InstallDate: w.InstallDate,
		CommitRate:  w.CommitRate,
	}
}

func wireRackToDomain(w WireRack) domain.Rack {
	return domain.Rack{
		ID:          w.ID,
		Name:        w.Name,
		FacilityID:  w.FacilityID,
		Site:        wireNestedToDomain(w.Site),
		Location:    wireNestedToDomain(w.Location),
		Tenant:      wireNestedToDomain(w.Tenant),
		Status:      wireLabelToDomain(w.Status),
		Role:        wireNestedToDomain(w.Role),
		Serial:      w.Serial,
		AssetTag:    w.AssetTag,
		Type:        wireLabelToDomain(w.Type),
		Width:       w.Width,
		UHeight:     w.UHeight,
		Comments:    w.Comments,
		Tags:        wireTagsToDomain(w.Tags),
		Created:     w.Created,
		LastUpdated: w.LastUpdated,
	}
}
