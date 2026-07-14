package netbox

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/teran/mcp-netbox/domain"
)

var objectTypeToEndpoint = map[string]string{
	"site":                "/api/dcim/sites/",
	"device":              "/api/dcim/devices/",
	"prefix":              "/api/ipam/prefixes/",
	"ip_address":          "/api/ipam/ip-addresses/",
	"vlan":                "/api/ipam/vlans/",
	"virtual_machine":     "/api/virtualization/virtual-machines/",
	"cluster":             "/api/virtualization/clusters/",
	"circuit":             "/api/circuits/circuits/",
	"provider":            "/api/circuits/providers/",
	"tenant":              "/api/tenancy/tenants/",
	"rack":                "/api/dcim/racks/",
	"manufacturer":        "/api/dcim/manufacturers/",
	"device_type":         "/api/dcim/device-types/",
	"location":            "/api/dcim/locations/",
	"cluster_type":        "/api/virtualization/cluster-types/",
	"cluster_group":       "/api/virtualization/cluster-groups/",
	"circuit_type":        "/api/circuits/circuit-types/",
	"vrf":                 "/api/ipam/vrfs/",
	"vlan_group":          "/api/ipam/vlan-groups/",
	"role":                "/api/ipam/roles/",
	"contact":             "/api/tenancy/contacts/",
	"cable":               "/api/dcim/cables/",
	"interface":           "/api/dcim/interfaces/",
	"vm_interface":        "/api/virtualization/interfaces/",
	"circuit_termination": "/api/circuits/circuit-terminations/",
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
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Warn("failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized: invalid token or insufficient permissions")
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("not found")
	}
	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("forbidden: token lacks required permissions")
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		if retryAfter != "" {
			return nil, fmt.Errorf("rate limited by NetBox: retry after %ss", retryAfter)
		}
		return nil, fmt.Errorf("rate limited by NetBox")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status %d from NetBox", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024)) // 10 MB limit
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
	var wireResp domain.PaginatedResponse[WireSite]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sites: %w", err)
	}
	return convertPaginated(&wireResp, wireSiteToDomain), nil
}

// ListDevices returns a paginated list of devices.
func (c *Client) ListDevices(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Device], error) {
	raw, err := c.get(ctx, token, "/api/dcim/devices/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireDevice]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal devices: %w", err)
	}
	return convertPaginated(&wireResp, wireDeviceToDomain), nil
}

// ListIPAddresses returns a paginated list of IP addresses.
func (c *Client) ListIPAddresses(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.IPAddress], error) {
	raw, err := c.get(ctx, token, "/api/ipam/ip-addresses/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireIPAddress]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal IP addresses: %w", err)
	}
	return convertPaginated(&wireResp, wireIPAddressToDomain), nil
}

// ListPrefixes returns a paginated list of prefixes.
func (c *Client) ListPrefixes(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Prefix], error) {
	raw, err := c.get(ctx, token, "/api/ipam/prefixes/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WirePrefix]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal prefixes: %w", err)
	}
	return convertPaginated(&wireResp, wirePrefixToDomain), nil
}

// ListVLANs returns a paginated list of VLANs.
func (c *Client) ListVLANs(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VLAN], error) {
	raw, err := c.get(ctx, token, "/api/ipam/vlans/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireVLAN]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal VLANs: %w", err)
	}
	return convertPaginated(&wireResp, wireVLANToDomain), nil
}

// ListVirtualMachines returns a paginated list of virtual machines.
func (c *Client) ListVirtualMachines(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VirtualMachine], error) {
	raw, err := c.get(ctx, token, "/api/virtualization/virtual-machines/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireVirtualMachine]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal VMs: %w", err)
	}
	return convertPaginated(&wireResp, wireVirtualMachineToDomain), nil
}

// ListClusters returns a paginated list of clusters.
func (c *Client) ListClusters(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Cluster], error) {
	raw, err := c.get(ctx, token, "/api/virtualization/clusters/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireCluster]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal clusters: %w", err)
	}
	return convertPaginated(&wireResp, wireClusterToDomain), nil
}

// ListCircuits returns a paginated list of circuits.
func (c *Client) ListCircuits(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Circuit], error) {
	raw, err := c.get(ctx, token, "/api/circuits/circuits/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireCircuit]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal circuits: %w", err)
	}
	return convertPaginated(&wireResp, wireCircuitToDomain), nil
}

// ListCircuitTerminations returns a paginated list of circuit terminations.
func (c *Client) ListCircuitTerminations(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.CircuitTermination], error) {
	raw, err := c.get(ctx, token, "/api/circuits/circuit-terminations/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireCircuitTermination]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal circuit terminations: %w", err)
	}
	return convertPaginated(&wireResp, wireCircuitTerminationToDomain), nil
}

// ListCables returns a paginated list of cables.
func (c *Client) ListCables(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Cable], error) {
	raw, err := c.get(ctx, token, "/api/dcim/cables/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireCable]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cables: %w", err)
	}
	return convertPaginated(&wireResp, wireCableToDomain), nil
}

// ListRacks returns a paginated list of racks.
func (c *Client) ListRacks(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Rack], error) {
	raw, err := c.get(ctx, token, "/api/dcim/racks/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireRack]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal racks: %w", err)
	}
	return convertPaginated(&wireResp, wireRackToDomain), nil
}

// ListInterfaces returns a paginated list of device interfaces.
func (c *Client) ListInterfaces(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.Interface], error) {
	raw, err := c.get(ctx, token, "/api/dcim/interfaces/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireInterface]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal interfaces: %w", err)
	}
	return convertPaginated(&wireResp, wireInterfaceToDomain), nil
}

// ListVMInterfaces returns a paginated list of VM interfaces.
func (c *Client) ListVMInterfaces(ctx context.Context, token string, params map[string]string) (*domain.PaginatedResponse[domain.VMInterface], error) {
	raw, err := c.get(ctx, token, "/api/virtualization/interfaces/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireVMInterface]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal VM interfaces: %w", err)
	}
	return convertPaginated(&wireResp, wireVMInterfaceToDomain), nil
}

// GetObject retrieves a single object by its type and ID.
func (c *Client) GetObject(ctx context.Context, token, objectType string, id int, params map[string]string) (domain.RawObject, error) {
	endpoint, ok := objectTypeToEndpoint[objectType]
	if !ok {
		return nil, fmt.Errorf("unknown object type: %s", objectType)
	}

	path := fmt.Sprintf("%s%d/", endpoint, id)
	body, err := c.get(ctx, token, path, params)
	if err != nil {
		return nil, fmt.Errorf("%s #%d: %w", objectType, id, err)
	}

	return domain.RawObject(body), nil
}

// convertPaginated converts a paginated response from wire type W to domain type D.
func convertPaginated[W, D any](resp *domain.PaginatedResponse[W], convert func(W) D) *domain.PaginatedResponse[D] {
	results := make([]D, len(resp.Results))
	for i, w := range resp.Results {
		results[i] = convert(w)
	}
	return &domain.PaginatedResponse[D]{
		Count:    resp.Count,
		Next:     resp.Next,
		Previous: resp.Previous,
		Results:  results,
	}
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
		Display:         w.Display,
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
		Display:     w.Display,
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
		Display:            w.Display,
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

func wireCircuitTerminationToDomain(w WireCircuitTermination) domain.CircuitTermination {
	return domain.CircuitTermination{
		ID:            w.ID,
		Display:       w.Display,
		Circuit:       wireNestedToDomain(w.Circuit),
		TermSide:      w.TermSide,
		Site:          wireNestedToDomain(w.Site),
		Speed:         w.Speed,
		UpstreamSpeed: w.UpstreamSpeed,
		Description:   w.Description,
		Tags:          wireTagsToDomain(w.Tags),
		Created:       w.Created,
		LastUpdated:   w.LastUpdated,
	}
}

func wireCableToDomain(w WireCable) domain.Cable {
	return domain.Cable{
		ID:          w.ID,
		Display:     w.Display,
		Type:        wireLabelToDomain(w.Type),
		Status:      wireLabelToDomain(w.Status),
		Label:       w.Label,
		Color:       w.Color,
		Length:      w.Length,
		LengthUnit:  wireLabelToDomain(w.LengthUnit),
		Description: w.Description,
		Tags:        wireTagsToDomain(w.Tags),
		Created:     w.Created,
		LastUpdated: w.LastUpdated,
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
		Display:     w.Display,
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
		Display:     w.Display,
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
		Display:     w.Display,
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
		Display:      w.Display,
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
		Display:     w.Display,
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

func wireInterfaceToDomain(w WireInterface) domain.Interface {
	return domain.Interface{
		ID:          w.ID,
		Name:        w.Name,
		Display:     w.Display,
		Device:      wireNestedToDomain(w.Device),
		Type:        wireLabelToDomain(w.Type),
		Enabled:     w.Enabled,
		MTU:         w.MTU,
		MACAddress:  w.MACAddress,
		Speed:       w.Speed,
		Description: w.Description,
		Tags:        wireTagsToDomain(w.Tags),
		Created:     w.Created,
		LastUpdated: w.LastUpdated,
	}
}

func wireVMInterfaceToDomain(w WireVMInterface) domain.VMInterface {
	return domain.VMInterface{
		ID:             w.ID,
		Name:           w.Name,
		Display:        w.Display,
		VirtualMachine: wireNestedToDomain(w.VirtualMachine),
		Enabled:        w.Enabled,
		MTU:            w.MTU,
		MACAddress:     w.MACAddress,
		Description:    w.Description,
		Tags:           wireTagsToDomain(w.Tags),
		Created:        w.Created,
		LastUpdated:    w.LastUpdated,
	}
}

func wireRackToDomain(w WireRack) domain.Rack {
	return domain.Rack{
		ID:          w.ID,
		Name:        w.Name,
		Display:     w.Display,
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
