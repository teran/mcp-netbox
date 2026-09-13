// Package netbox implements the domain.NetworkRepository interface as an HTTP
// client for the NetBox REST API.
//
// It provides:
//   - Typed list methods for all NetBox entities (sites, devices, IPs, etc.)
//   - A generic GetObject method that dynamically routes to the correct endpoint
//     based on a configurable object type-to-endpoint mapping
//   - Wire-to-domain model conversion functions for each entity type
//   - HTTP error handling with specific errors for 401, 403, 404, and 429 responses
package netbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"resty.dev/v3"

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

// maxBodySize is the maximum response body size accepted from NetBox. Larger
// responses are rejected to bound memory usage.
const maxBodySize = 10 * 1024 * 1024

// Client is an HTTP client for the NetBox API.
type Client struct {
	baseURL string
	client  *resty.Client
}

// NewClient creates a new NetBox API client backed by resty. The supplied
// *http.Client is wrapped via resty.NewWithClient so that its transport
// (DNS-rebinding dialer + circuit breaker), timeout and redirect policy are
// preserved.
func NewClient(baseURL string, httpClient *http.Client) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		client:  resty.NewWithClient(httpClient),
	}
}

func (c *Client) doRequest(ctx context.Context, token, method, path string, params map[string]string, body []byte) ([]byte, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	start := time.Now()

	req := c.client.R().
		SetContext(ctx).
		SetResponseBodyLimit(maxBodySize).
		SetResponseBodyUnlimitedReads(true).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+token).
		SetQueryParams(params)

	if body != nil {
		req = req.SetBody(body).SetHeader("Content-Type", "application/json")
	}

	// Propagate the correlation ID to the upstream NetBox call so records on
	// both sides of the wire can be matched (L9/G11).
	if requestID := domain.RequestIDFromContext(ctx); requestID != "" {
		req = req.SetHeader(requestIDHeader, requestID)
	}

	resp, err := req.Execute(method, u.String())
	duration := time.Since(start)

	// Emit an outbound per-request log record tagged with the same request_id.
	logger := getLogger()
	entry := logger.WithFields(logrus.Fields{
		"method":   method,
		"path":     u.Path,
		"duration": duration,
		"status":   statusCodeOf(resp),
	})
	if requestID := domain.RequestIDFromContext(ctx); requestID != "" {
		entry = entry.WithField("request_id", requestID)
	}
	entry.Debug("netbox request")

	if err != nil {
		if errors.Is(err, resty.ErrReadExceedsThresholdLimit) {
			return nil, fmt.Errorf("response body exceeds %d bytes (got at least %d bytes), truncation detected", maxBodySize, maxBodySize)
		}
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized: invalid token or insufficient permissions")
	}
	if resp.StatusCode() == http.StatusNotFound {
		return nil, fmt.Errorf("not found")
	}
	if resp.StatusCode() == http.StatusForbidden {
		return nil, fmt.Errorf("forbidden: token lacks required permissions")
	}

	if resp.StatusCode() == http.StatusTooManyRequests {
		retryAfter := resp.Header().Get("Retry-After")
		if retryAfter != "" {
			return nil, fmt.Errorf("rate limited by NetBox: retry after %ss", retryAfter)
		}
		return nil, fmt.Errorf("rate limited by NetBox")
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		// Only read the body for 400 so we can surface the NetBox validation
		// errors; the body is bounded by maxBodySize. Other unexpected statuses
		// keep the generic message (matching read behaviour).
		if resp.StatusCode() == http.StatusBadRequest {
			return nil, &netboxError{StatusCode: resp.StatusCode(), Body: resp.Bytes()}
		}
		return nil, fmt.Errorf("unexpected status %d from NetBox", resp.StatusCode())
	}

	return resp.Bytes(), nil
}

// netboxError carries the HTTP status and bounded response body for an
// unexpected non-2xx status. It lets write methods recover a 400 body (to build
// a domain.ValidationError) while preserving the generic message for reads.
type netboxError struct {
	StatusCode int
	Body       []byte
}

func (e *netboxError) Error() string {
	return fmt.Sprintf("unexpected status %d from NetBox", e.StatusCode)
}

// statusCodeOf returns the HTTP status code of a resty response, or 0 when the
// response is nil (e.g. a transport-level error).
func statusCodeOf(resp *resty.Response) int {
	if resp == nil {
		return 0
	}
	return resp.StatusCode()
}

// requestIDHeader is the standard correlation header forwarded to upstream
// HTTP calls.
const requestIDHeader = "X-Request-ID"

func (c *Client) get(ctx context.Context, token, path string, params map[string]string) ([]byte, error) {
	return c.doRequest(ctx, token, http.MethodGet, path, params, nil)
}

func (c *Client) post(ctx context.Context, token, path string, body []byte) ([]byte, error) {
	return c.doRequest(ctx, token, http.MethodPost, path, nil, body)
}

func (c *Client) patch(ctx context.Context, token, path string, body []byte) ([]byte, error) {
	return c.doRequest(ctx, token, http.MethodPatch, path, nil, body)
}

func (c *Client) delete(ctx context.Context, token, path string) ([]byte, error) {
	return c.doRequest(ctx, token, http.MethodDelete, path, nil, nil)
}

// mapWriteStatus maps a write HTTP status to a result or error. 201/200 yield
// the created/updated object body, 204 (delete) yields nil, and 400 yields a
// domain.ValidationError carrying the NetBox body. 401/403/404/429 are handled
// upstream by doRequest and never reach this function.
func mapWriteStatus(code int, body []byte) (domain.RawObject, error) {
	switch code {
	case http.StatusCreated, http.StatusOK:
		return domain.RawObject(body), nil
	case http.StatusNoContent:
		return nil, nil
	case http.StatusBadRequest:
		return nil, &domain.ValidationError{StatusCode: code, Body: body}
	default:
		return nil, fmt.Errorf("unexpected status %d from NetBox", code)
	}
}

// doWrite performs a write request and maps the response status via
// mapWriteStatus. Transport errors and 401/403/404/429 are returned as-is from
// doRequest (not duplicated here); only a 400 body is recovered from the
// netboxError returned by doRequest.
func writeResult(successStatus int, body []byte, err error) (domain.RawObject, error) {
	if err != nil {
		var nb *netboxError
		if errors.As(err, &nb) {
			return mapWriteStatus(nb.StatusCode, nb.Body)
		}
		return nil, err
	}
	return mapWriteStatus(successStatus, body)
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

// CreateSite creates a new site via POST /api/dcim/sites/.
func (c *Client) CreateSite(ctx context.Context, token string, in domain.SiteWrite) (*domain.Site, error) {
	payload, err := json.Marshal(siteWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal site create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/sites/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}
	return unmarshalSite(raw)
}

// UpdateSite partially updates a site via PATCH /api/dcim/sites/<id>/.
func (c *Client) UpdateSite(ctx context.Context, token string, id int, in domain.SiteWrite) (*domain.Site, error) {
	payload, err := json.Marshal(siteWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal site update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/sites/%d/", id)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}
	return unmarshalSite(raw)
}

// DeleteSite deletes a site via DELETE /api/dcim/sites/<id>/.
func (c *Client) DeleteSite(ctx context.Context, token string, id int) error {
	path := fmt.Sprintf("/api/dcim/sites/%d/", id)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)
	return mErr
}

// CreateDevice creates a new device via POST /api/dcim/devices/.
func (c *Client) CreateDevice(ctx context.Context, token string, in domain.DeviceWrite) (*domain.Device, error) {
	payload, err := json.Marshal(deviceWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal device create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/devices/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}
	return unmarshalDevice(raw)
}

// UpdateDevice partially updates a device via PATCH /api/dcim/devices/<id>/.
func (c *Client) UpdateDevice(ctx context.Context, token string, id int, in domain.DeviceWrite) (*domain.Device, error) {
	payload, err := json.Marshal(deviceWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal device update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/devices/%d/", id)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}
	return unmarshalDevice(raw)
}

// DeleteDevice deletes a device via DELETE /api/dcim/devices/<id>/.
func (c *Client) DeleteDevice(ctx context.Context, token string, id int) error {
	path := fmt.Sprintf("/api/dcim/devices/%d/", id)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)
	return mErr
}

// CreateIPAddress creates a new IP address via POST /api/ipam/ip-addresses/.
func (c *Client) CreateIPAddress(ctx context.Context, token string, in domain.IPAddressWrite) (*domain.IPAddress, error) {
	payload, err := json.Marshal(ipAddressWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal IP address create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/ipam/ip-addresses/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}
	return unmarshalIPAddress(raw)
}

// UpdateIPAddress partially updates an IP address via
// PATCH /api/ipam/ip-addresses/<id>/.
func (c *Client) UpdateIPAddress(ctx context.Context, token string, id int, in domain.IPAddressWrite) (*domain.IPAddress, error) {
	payload, err := json.Marshal(ipAddressWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal IP address update request: %w", err)
	}
	path := fmt.Sprintf("/api/ipam/ip-addresses/%d/", id)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}
	return unmarshalIPAddress(raw)
}

// DeleteIPAddress deletes an IP address via DELETE /api/ipam/ip-addresses/<id>/.
func (c *Client) DeleteIPAddress(ctx context.Context, token string, id int) error {
	path := fmt.Sprintf("/api/ipam/ip-addresses/%d/", id)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)
	return mErr
}

// CreatePrefix creates a new prefix via POST /api/ipam/prefixes/.
func (c *Client) CreatePrefix(ctx context.Context, token string, in domain.PrefixWrite) (*domain.Prefix, error) {
	payload, err := json.Marshal(prefixWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal prefix create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/ipam/prefixes/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}
	return unmarshalPrefix(raw)
}

// UpdatePrefix partially updates a prefix via PATCH /api/ipam/prefixes/<id>/.
func (c *Client) UpdatePrefix(ctx context.Context, token string, id int, in domain.PrefixWrite) (*domain.Prefix, error) {
	payload, err := json.Marshal(prefixWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal prefix update request: %w", err)
	}
	path := fmt.Sprintf("/api/ipam/prefixes/%d/", id)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}
	return unmarshalPrefix(raw)
}

// DeletePrefix deletes a prefix via DELETE /api/ipam/prefixes/<id>/.
func (c *Client) DeletePrefix(ctx context.Context, token string, id int) error {
	path := fmt.Sprintf("/api/ipam/prefixes/%d/", id)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)
	return mErr
}

// CreateVLAN creates a new VLAN via POST /api/ipam/vlans/.
func (c *Client) CreateVLAN(ctx context.Context, token string, in domain.VLANWrite) (*domain.VLAN, error) {
	payload, err := json.Marshal(vlanWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal VLAN create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/ipam/vlans/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}
	return unmarshalVLAN(raw)
}

// UpdateVLAN partially updates a VLAN via PATCH /api/ipam/vlans/<id>/.
func (c *Client) UpdateVLAN(ctx context.Context, token string, id int, in domain.VLANWrite) (*domain.VLAN, error) {
	payload, err := json.Marshal(vlanWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal VLAN update request: %w", err)
	}
	path := fmt.Sprintf("/api/ipam/vlans/%d/", id)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}
	return unmarshalVLAN(raw)
}

// DeleteVLAN deletes a VLAN via DELETE /api/ipam/vlans/<id>/.
func (c *Client) DeleteVLAN(ctx context.Context, token string, id int) error {
	path := fmt.Sprintf("/api/ipam/vlans/%d/", id)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)
	return mErr
}

// CreateVirtualMachine creates a new virtual machine via POST
// /api/virtualization/virtual-machines/.
func (c *Client) CreateVirtualMachine(ctx context.Context, token string, in domain.VirtualMachineWrite) (*domain.VirtualMachine, error) {
	payload, err := json.Marshal(virtualMachineWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal virtual machine create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/virtualization/virtual-machines/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}
	return unmarshalVirtualMachine(raw)
}

// UpdateVirtualMachine partially updates a virtual machine via PATCH
// /api/virtualization/virtual-machines/<id>/.
func (c *Client) UpdateVirtualMachine(ctx context.Context, token string, id int, in domain.VirtualMachineWrite) (*domain.VirtualMachine, error) {
	payload, err := json.Marshal(virtualMachineWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal virtual machine update request: %w", err)
	}
	path := fmt.Sprintf("/api/virtualization/virtual-machines/%d/", id)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}
	return unmarshalVirtualMachine(raw)
}

// DeleteVirtualMachine deletes a virtual machine via DELETE
// /api/virtualization/virtual-machines/<id>/.
func (c *Client) DeleteVirtualMachine(ctx context.Context, token string, id int) error {
	path := fmt.Sprintf("/api/virtualization/virtual-machines/%d/", id)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)
	return mErr
}

// CreateCluster creates a new cluster via POST /api/virtualization/clusters/.
func (c *Client) CreateCluster(ctx context.Context, token string, in domain.ClusterWrite) (*domain.Cluster, error) {
	payload, err := json.Marshal(clusterWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal cluster create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/virtualization/clusters/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}
	return unmarshalCluster(raw)
}

// UpdateCluster partially updates a cluster via PATCH
// /api/virtualization/clusters/<id>/.
func (c *Client) UpdateCluster(ctx context.Context, token string, id int, in domain.ClusterWrite) (*domain.Cluster, error) {
	payload, err := json.Marshal(clusterWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal cluster update request: %w", err)
	}
	path := fmt.Sprintf("/api/virtualization/clusters/%d/", id)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}
	return unmarshalCluster(raw)
}

// DeleteCluster deletes a cluster via DELETE /api/virtualization/clusters/<id>/.
func (c *Client) DeleteCluster(ctx context.Context, token string, id int) error {
	path := fmt.Sprintf("/api/virtualization/clusters/%d/", id)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)
	return mErr
}

// CreateCircuit creates a new circuit via POST /api/circuits/circuits/.
func (c *Client) CreateCircuit(ctx context.Context, token string, in domain.CircuitWrite) (*domain.Circuit, error) {
	payload, err := json.Marshal(circuitWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal circuit create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/circuits/circuits/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}
	return unmarshalCircuit(raw)
}

// UpdateCircuit partially updates a circuit via PATCH /api/circuits/circuits/<id>/.
func (c *Client) UpdateCircuit(ctx context.Context, token string, id int, in domain.CircuitWrite) (*domain.Circuit, error) {
	payload, err := json.Marshal(circuitWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal circuit update request: %w", err)
	}
	path := fmt.Sprintf("/api/circuits/circuits/%d/", id)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}
	return unmarshalCircuit(raw)
}

// DeleteCircuit deletes a circuit via DELETE /api/circuits/circuits/<id>/.
func (c *Client) DeleteCircuit(ctx context.Context, token string, id int) error {
	path := fmt.Sprintf("/api/circuits/circuits/%d/", id)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)
	return mErr
}

// convertPaginated converts a paginated response from wire type W to domain type D.
func convertPaginated[W, D any](resp *domain.PaginatedResponse[W], convert func(W) D) *domain.PaginatedResponse[D] {
	if resp == nil {
		return nil
	}
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

// siteWriteToWire converts a domain write DTO to the wire request model. The
// fields map one-to-one; nil pointers on the domain side remain nil on the wire
// side so json.Marshal omits them (required for a partial PATCH body).
func siteWriteToWire(in domain.SiteWrite) WireSiteWrite {
	return WireSiteWrite{
		Name:            in.Name,
		Slug:            in.Slug,
		Status:          in.Status,
		Region:          in.Region,
		Tenant:          in.Tenant,
		Facility:        in.Facility,
		TimeZone:        in.TimeZone,
		Description:     in.Description,
		PhysicalAddress: in.PhysicalAddress,
		ShippingAddress: in.ShippingAddress,
		Comments:        in.Comments,
		Tags:            in.Tags,
		CustomFields:    in.CustomFields,
	}
}

// unmarshalSite decodes a site response body into a domain.Site.
func unmarshalSite(raw domain.RawObject) (*domain.Site, error) {
	var w WireSite
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal site: %w", err)
	}
	d := wireSiteToDomain(w)
	return &d, nil
}

// deviceWriteToWire converts a domain device write DTO to the wire request
// model. The fields map one-to-one; nil pointers on the domain side remain nil
// on the wire side so json.Marshal omits them (required for a partial PATCH
// body).
func deviceWriteToWire(in domain.DeviceWrite) WireDeviceWrite {
	return WireDeviceWrite{
		Name:         in.Name,
		DeviceType:   in.DeviceType,
		Role:         in.Role,
		Tenant:       in.Tenant,
		Platform:     in.Platform,
		Serial:       in.Serial,
		AssetTag:     in.AssetTag,
		Site:         in.Site,
		Rack:         in.Rack,
		Position:     in.Position,
		Face:         in.Face,
		Status:       in.Status,
		Cluster:      in.Cluster,
		Comments:     in.Comments,
		Tags:         in.Tags,
		CustomFields: in.CustomFields,
	}
}

// unmarshalDevice decodes a device response body into a domain.Device.
func unmarshalDevice(raw domain.RawObject) (*domain.Device, error) {
	var w WireDevice
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal device: %w", err)
	}
	d := wireDeviceToDomain(w)
	return &d, nil
}

// ipAddressWriteToWire converts a domain IP address write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func ipAddressWriteToWire(in domain.IPAddressWrite) WireIPAddressWrite {
	return WireIPAddressWrite{
		Address:            in.Address,
		Status:             in.Status,
		Role:               in.Role,
		VRF:                in.VRF,
		Tenant:             in.Tenant,
		DNSName:            in.DNSName,
		Description:        in.Description,
		AssignedObjectType: in.AssignedObjectType,
		AssignedObjectID:   in.AssignedObjectID,
		Comments:           in.Comments,
		Tags:               in.Tags,
		CustomFields:       in.CustomFields,
	}
}

// unmarshalIPAddress decodes an IP address response body into a
// domain.IPAddress.
func unmarshalIPAddress(raw domain.RawObject) (*domain.IPAddress, error) {
	var w WireIPAddress
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal IP address: %w", err)
	}
	d := wireIPAddressToDomain(w)
	return &d, nil
}

// prefixWriteToWire converts a domain prefix write DTO to the wire request
// model. The fields map one-to-one; nil pointers on the domain side remain nil
// on the wire side so json.Marshal omits them (required for a partial PATCH
// body).
func prefixWriteToWire(in domain.PrefixWrite) WirePrefixWrite {
	return WirePrefixWrite{
		Prefix:       in.Prefix,
		Site:         in.Site,
		VRF:          in.VRF,
		Tenant:       in.Tenant,
		VLAN:         in.VLAN,
		Status:       in.Status,
		Role:         in.Role,
		IsPool:       in.IsPool,
		Description:  in.Description,
		Comments:     in.Comments,
		Tags:         in.Tags,
		CustomFields: in.CustomFields,
	}
}

// unmarshalPrefix decodes a prefix response body into a domain.Prefix.
func unmarshalPrefix(raw domain.RawObject) (*domain.Prefix, error) {
	var w WirePrefix
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal prefix: %w", err)
	}
	d := wirePrefixToDomain(w)
	return &d, nil
}

// vlanWriteToWire converts a domain VLAN write DTO to the wire request model.
// The fields map one-to-one; nil pointers on the domain side remain nil on the
// wire side so json.Marshal omits them (required for a partial PATCH body).
func vlanWriteToWire(in domain.VLANWrite) WireVLANWrite {
	return WireVLANWrite{
		VID:          in.VID,
		Name:         in.Name,
		Site:         in.Site,
		Group:        in.Group,
		Tenant:       in.Tenant,
		Status:       in.Status,
		Role:         in.Role,
		Description:  in.Description,
		Comments:     in.Comments,
		Tags:         in.Tags,
		CustomFields: in.CustomFields,
	}
}

// unmarshalVLAN decodes a VLAN response body into a domain.VLAN.
func unmarshalVLAN(raw domain.RawObject) (*domain.VLAN, error) {
	var w WireVLAN
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal VLAN: %w", err)
	}
	d := wireVLANToDomain(w)
	return &d, nil
}

// virtualMachineWriteToWire converts a domain virtual machine write DTO to the
// wire request model. The fields map one-to-one; nil pointers on the domain
// side remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func virtualMachineWriteToWire(in domain.VirtualMachineWrite) WireVirtualMachineWrite {
	return WireVirtualMachineWrite{
		Name:         in.Name,
		Cluster:      in.Cluster,
		Role:         in.Role,
		Tenant:       in.Tenant,
		Platform:     in.Platform,
		Status:       in.Status,
		Site:         in.Site,
		VCPUs:        in.VCPUs,
		Memory:       in.Memory,
		Disk:         in.Disk,
		Comments:     in.Comments,
		Tags:         in.Tags,
		CustomFields: in.CustomFields,
	}
}

// unmarshalVirtualMachine decodes a virtual machine response body into a
// domain.VirtualMachine.
func unmarshalVirtualMachine(raw domain.RawObject) (*domain.VirtualMachine, error) {
	var w WireVirtualMachine
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal virtual machine: %w", err)
	}
	d := wireVirtualMachineToDomain(w)
	return &d, nil
}

// clusterWriteToWire converts a domain cluster write DTO to the wire request
// model. The fields map one-to-one; nil pointers on the domain side remain nil
// on the wire side so json.Marshal omits them (required for a partial PATCH
// body).
func clusterWriteToWire(in domain.ClusterWrite) WireClusterWrite {
	return WireClusterWrite{
		Name:         in.Name,
		ClusterType:  in.ClusterType,
		ClusterGroup: in.ClusterGroup,
		Site:         in.Site,
		Tenant:       in.Tenant,
		Description:  in.Description,
		Comments:     in.Comments,
		Tags:         in.Tags,
		CustomFields: in.CustomFields,
	}
}

// unmarshalCluster decodes a cluster response body into a domain.Cluster.
func unmarshalCluster(raw domain.RawObject) (*domain.Cluster, error) {
	var w WireCluster
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cluster: %w", err)
	}
	d := wireClusterToDomain(w)
	return &d, nil
}

// circuitWriteToWire converts a domain circuit write DTO to the wire request
// model. The fields map one-to-one; nil pointers on the domain side remain nil
// on the wire side so json.Marshal omits them (required for a partial PATCH
// body).
func circuitWriteToWire(in domain.CircuitWrite) WireCircuitWrite {
	return WireCircuitWrite{
		CID:          in.CID,
		Provider:     in.Provider,
		CircuitType:  in.CircuitType,
		Tenant:       in.Tenant,
		Status:       in.Status,
		Description:  in.Description,
		Comments:     in.Comments,
		Tags:         in.Tags,
		CustomFields: in.CustomFields,
		InstallDate:  in.InstallDate,
		CommitRate:   in.CommitRate,
	}
}

// unmarshalCircuit decodes a circuit response body into a domain.Circuit.
func unmarshalCircuit(raw domain.RawObject) (*domain.Circuit, error) {
	var w WireCircuit
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal circuit: %w", err)
	}
	d := wireCircuitToDomain(w)
	return &d, nil
}

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
		URL:             w.URL,
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
		CustomFields:    w.CustomFields,
		Created:         w.Created,
		LastUpdated:     w.LastUpdated,
	}
}

func wireDeviceToDomain(w WireDevice) domain.Device {
	return domain.Device{
		ID:           w.ID,
		URL:          w.URL,
		Name:         w.Name,
		Display:      w.Display,
		DeviceType:   wireNestedToDomain(w.DeviceType),
		Role:         wireNestedToDomain(w.Role),
		Tenant:       wireNestedToDomain(w.Tenant),
		Platform:     wireNestedToDomain(w.Platform),
		Serial:       w.Serial,
		AssetTag:     w.AssetTag,
		Site:         wireNestedToDomain(w.Site),
		Rack:         wireNestedToDomain(w.Rack),
		Position:     w.Position,
		Face:         wireLabelToDomain(w.Face),
		Status:       wireLabelToDomain(w.Status),
		Cluster:      wireNestedToDomain(w.Cluster),
		Comments:     w.Comments,
		Tags:         wireTagsToDomain(w.Tags),
		CustomFields: w.CustomFields,
		Created:      w.Created,
		LastUpdated:  w.LastUpdated,
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
		URL:                w.URL,
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
		CustomFields:       w.CustomFields,
		Created:            w.Created,
		LastUpdated:        w.LastUpdated,
	}
}

func wireCircuitTerminationToDomain(w WireCircuitTermination) domain.CircuitTermination {
	return domain.CircuitTermination{
		ID:            w.ID,
		URL:           w.URL,
		Display:       w.Display,
		Circuit:       wireNestedToDomain(w.Circuit),
		TermSide:      w.TermSide,
		Site:          wireNestedToDomain(w.Site),
		Speed:         w.Speed,
		UpstreamSpeed: w.UpstreamSpeed,
		Description:   w.Description,
		Tags:          wireTagsToDomain(w.Tags),
		CustomFields:  w.CustomFields,
		Created:       w.Created,
		LastUpdated:   w.LastUpdated,
	}
}

func wireCableToDomain(w WireCable) domain.Cable {
	return domain.Cable{
		ID:           w.ID,
		URL:          w.URL,
		Display:      w.Display,
		Type:         wireLabelToDomain(w.Type),
		Status:       wireLabelToDomain(w.Status),
		Label:        w.Label,
		Color:        w.Color,
		Length:       w.Length,
		LengthUnit:   wireLabelToDomain(w.LengthUnit),
		TerminationA: wireCableTerminationToDomain(w.TerminationA),
		TerminationB: wireCableTerminationToDomain(w.TerminationB),
		Description:  w.Description,
		Tags:         wireTagsToDomain(w.Tags),
		CustomFields: w.CustomFields,
		Created:      w.Created,
		LastUpdated:  w.LastUpdated,
	}
}

func wireCableTerminationToDomain(w *WireCableTermination) *domain.CableTermination {
	if w == nil {
		return nil
	}
	return &domain.CableTermination{
		ID:   w.ID,
		URL:  w.URL,
		Name: w.Name,
		Type: w.Type,
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
		ID:           w.ID,
		URL:          w.URL,
		Prefix:       w.Prefix,
		Display:      w.Display,
		Site:         wireNestedToDomain(w.Site),
		VRF:          wireNestedToDomain(w.VRF),
		Tenant:       wireNestedToDomain(w.Tenant),
		VLAN:         wireNestedToDomain(w.VLAN),
		Status:       wireLabelToDomain(w.Status),
		Role:         wireNestedToDomain(w.Role),
		IsPool:       w.IsPool,
		Description:  w.Description,
		Comments:     w.Comments,
		Tags:         wireTagsToDomain(w.Tags),
		CustomFields: w.CustomFields,
		Created:      w.Created,
		LastUpdated:  w.LastUpdated,
		Children:     w.Children,
		Depth:        w.Depth,
		Family:       family,
	}
}

func wireVLANToDomain(w WireVLAN) domain.VLAN {
	return domain.VLAN{
		ID:           w.ID,
		URL:          w.URL,
		Site:         wireNestedToDomain(w.Site),
		Group:        wireNestedToDomain(w.Group),
		VID:          w.VID,
		Name:         w.Name,
		Display:      w.Display,
		Tenant:       wireNestedToDomain(w.Tenant),
		Status:       wireLabelToDomain(w.Status),
		Role:         wireNestedToDomain(w.Role),
		Description:  w.Description,
		Comments:     w.Comments,
		Tags:         wireTagsToDomain(w.Tags),
		CustomFields: w.CustomFields,
		Created:      w.Created,
		LastUpdated:  w.LastUpdated,
	}
}

func wireVirtualMachineToDomain(w WireVirtualMachine) domain.VirtualMachine {
	return domain.VirtualMachine{
		ID:           w.ID,
		URL:          w.URL,
		Name:         w.Name,
		Display:      w.Display,
		Cluster:      wireNestedToDomain(w.Cluster),
		Role:         wireNestedToDomain(w.Role),
		Tenant:       wireNestedToDomain(w.Tenant),
		Platform:     wireNestedToDomain(w.Platform),
		Status:       wireLabelToDomain(w.Status),
		Site:         wireNestedToDomain(w.Site),
		VCPUs:        w.VCPUs,
		Memory:       w.Memory,
		Disk:         w.Disk,
		Comments:     w.Comments,
		Tags:         wireTagsToDomain(w.Tags),
		CustomFields: w.CustomFields,
		Created:      w.Created,
		LastUpdated:  w.LastUpdated,
	}
}

func wireClusterToDomain(w WireCluster) domain.Cluster {
	return domain.Cluster{
		ID:           w.ID,
		URL:          w.URL,
		Name:         w.Name,
		Display:      w.Display,
		ClusterType:  wireNestedToDomain(w.ClusterType),
		ClusterGroup: wireNestedToDomain(w.ClusterGroup),
		Site:         wireNestedToDomain(w.Site),
		Tenant:       wireNestedToDomain(w.Tenant),
		Description:  w.Description,
		Comments:     w.Comments,
		Tags:         wireTagsToDomain(w.Tags),
		CustomFields: w.CustomFields,
		Created:      w.Created,
		LastUpdated:  w.LastUpdated,
	}
}

func wireCircuitToDomain(w WireCircuit) domain.Circuit {
	return domain.Circuit{
		ID:           w.ID,
		URL:          w.URL,
		CID:          w.CID,
		Display:      w.Display,
		Provider:     wireNestedToDomain(w.Provider),
		CircuitType:  wireNestedToDomain(w.CircuitType),
		Tenant:       wireNestedToDomain(w.Tenant),
		Status:       wireLabelToDomain(w.Status),
		Description:  w.Description,
		Comments:     w.Comments,
		Tags:         wireTagsToDomain(w.Tags),
		CustomFields: w.CustomFields,
		Created:      w.Created,
		LastUpdated:  w.LastUpdated,
		InstallDate:  w.InstallDate,
		CommitRate:   w.CommitRate,
	}
}

func wireInterfaceToDomain(w WireInterface) domain.Interface {
	return domain.Interface{
		ID:           w.ID,
		URL:          w.URL,
		Name:         w.Name,
		Display:      w.Display,
		Device:       wireNestedToDomain(w.Device),
		Type:         wireLabelToDomain(w.Type),
		Enabled:      w.Enabled,
		MTU:          w.MTU,
		MACAddress:   w.MACAddress,
		Speed:        w.Speed,
		Description:  w.Description,
		Tags:         wireTagsToDomain(w.Tags),
		CustomFields: w.CustomFields,
		Created:      w.Created,
		LastUpdated:  w.LastUpdated,
	}
}

func wireVMInterfaceToDomain(w WireVMInterface) domain.VMInterface {
	return domain.VMInterface{
		ID:             w.ID,
		URL:            w.URL,
		Name:           w.Name,
		Display:        w.Display,
		VirtualMachine: wireNestedToDomain(w.VirtualMachine),
		Enabled:        w.Enabled,
		MTU:            w.MTU,
		MACAddress:     w.MACAddress,
		Description:    w.Description,
		Tags:           wireTagsToDomain(w.Tags),
		CustomFields:   w.CustomFields,
		Created:        w.Created,
		LastUpdated:    w.LastUpdated,
	}
}

func wireRackToDomain(w WireRack) domain.Rack {
	return domain.Rack{
		ID:           w.ID,
		URL:          w.URL,
		Name:         w.Name,
		Display:      w.Display,
		FacilityID:   w.FacilityID,
		Site:         wireNestedToDomain(w.Site),
		Location:     wireNestedToDomain(w.Location),
		Tenant:       wireNestedToDomain(w.Tenant),
		Status:       wireLabelToDomain(w.Status),
		Role:         wireNestedToDomain(w.Role),
		Serial:       w.Serial,
		AssetTag:     w.AssetTag,
		Type:         wireLabelToDomain(w.Type),
		Width:        w.Width,
		UHeight:      w.UHeight,
		Comments:     w.Comments,
		Tags:         wireTagsToDomain(w.Tags),
		CustomFields: w.CustomFields,
		Created:      w.Created,
		LastUpdated:  w.LastUpdated,
	}
}
