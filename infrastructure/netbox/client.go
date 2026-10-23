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

// ErrNotFound is returned by the client when NetBox responds 404 (resource not
// found). Delete operations treat it as a successful no-op (idempotent delete).
var ErrNotFound = errors.New("not found")

var objectTypeToEndpoint = map[string]string{
	"site":                    "/api/dcim/sites/",
	"region":                  "/api/dcim/regions/",
	"device":                  "/api/dcim/devices/",
	"prefix":                  "/api/ipam/prefixes/",
	"ip_address":              "/api/ipam/ip-addresses/",
	"vlan":                    "/api/ipam/vlans/",
	"virtual_machine":         "/api/virtualization/virtual-machines/",
	"cluster":                 "/api/virtualization/clusters/",
	"circuit":                 "/api/circuits/circuits/",
	"provider":                "/api/circuits/providers/",
	"tenant":                  "/api/tenancy/tenants/",
	"tenant_group":            "/api/tenancy/tenant-groups/",
	"rack":                    "/api/dcim/racks/",
	"manufacturer":            "/api/dcim/manufacturers/",
	"device_type":             "/api/dcim/device-types/",
	"location":                "/api/dcim/locations/",
	"cluster_type":            "/api/virtualization/cluster-types/",
	"cluster_group":           "/api/virtualization/cluster-groups/",
	"circuit_type":            "/api/circuits/circuit-types/",
	"vrf":                     "/api/ipam/vrfs/",
	"vlan_group":              "/api/ipam/vlan-groups/",
	"role":                    "/api/ipam/roles/",
	"contact":                 "/api/tenancy/contacts/",
	"cable":                   "/api/dcim/cables/",
	"interface":               "/api/dcim/interfaces/",
	"vm_interface":            "/api/virtualization/interfaces/",
	"circuit_termination":     "/api/circuits/circuit-terminations/",
	"site_group":              "/api/dcim/site-groups/",
	"rack_group":              "/api/dcim/rack-groups/",
	"rack_role":               "/api/dcim/rack-roles/",
	"device_role":             "/api/dcim/device-roles/",
	"platform":                "/api/dcim/platforms/",
	"virtual_chassis":         "/api/dcim/virtual-chassis/",
	"power_panel":             "/api/dcim/power-panels/",
	"inventory_item_role":     "/api/dcim/inventory-item-roles/",
	"front_port":              "/api/dcim/front-ports/",
	"rear_port":               "/api/dcim/rear-ports/",
	"console_port":            "/api/dcim/console-ports/",
	"console_server_port":     "/api/dcim/console-server-ports/",
	"power_port":              "/api/dcim/power-ports/",
	"power_outlet":            "/api/dcim/power-outlets/",
	"power_feed":              "/api/dcim/power-feeds/",
	"device_bay":              "/api/dcim/device-bays/",
	"rack_reservation":        "/api/dcim/rack-reservations/",
	"inventory_item":          "/api/dcim/inventory-items/",
	"mac_address":             "/api/dcim/mac-addresses/",
	"rir":                     "/api/ipam/rirs/",
	"aggregate":               "/api/ipam/aggregates/",
	"ip_range":                "/api/ipam/ip-ranges/",
	"asn":                     "/api/ipam/asns/",
	"asn_range":               "/api/ipam/asn-ranges/",
	"route_target":            "/api/ipam/route-targets/",
	"service":                 "/api/ipam/services/",
	"service_template":        "/api/ipam/service-templates/",
	"fhrp_group":              "/api/ipam/fhrp-groups/",
	"fhrp_group_assignment":   "/api/ipam/fhrp-group-assignments/",
	"vlan_translation_policy": "/api/ipam/vlan-translation-policies/",
	"vlan_translation_rule":   "/api/ipam/vlan-translation-rules/",
	"contact_group":           "/api/tenancy/contact-groups/",
	"contact_role":            "/api/tenancy/contact-roles/",
	"virtual_disk":            "/api/virtualization/virtual-disks/",
	"virtual_machine_type":    "/api/virtualization/virtual-machine-types/",
}

// maxBodySize is the maximum response body size accepted from NetBox. Larger
// responses are rejected to bound memory usage.
const maxBodySize = 10 * 1024 * 1024

// Client is an HTTP client for the NetBox API.
type Client struct {
	baseURL  string
	client   *resty.Client
	upstream *UpstreamMetrics
}

// NewClient creates a new NetBox API client backed by resty. The supplied
// *http.Client is wrapped via resty.NewWithClient so that its transport
// (DNS-rebinding dialer + circuit breaker), timeout and redirect policy are
// preserved. Upstream metrics are disabled (nil).
func NewClient(baseURL string, httpClient *http.Client) *Client {
	return NewClientWithMetrics(baseURL, httpClient, nil)
}

// NewClientWithMetrics is NewClient with an optional upstream-metrics handle.
// The handle is nil-safe; passing nil disables upstream metrics recording.
func NewClientWithMetrics(baseURL string, httpClient *http.Client, upstream *UpstreamMetrics) *Client {
	// Wire resty's internal logger to the package logrus logger so resty emits
	// its own request/response detail records (redacting the Authorization
	// header), and enable debug mode so those records are produced. They only
	// reach the log when the logrus level is debug or trace; otherwise the
	// restyLogger adapter drops them. The custom access-log line in doRequest
	// is emitted independently and is unaffected.
	restyClient := resty.NewWithClient(httpClient).
		SetLogger(restyLogger{}).
		SetDebug(true).
		SetHeader("User-Agent", userAgent())

	return &Client{
		baseURL:  strings.TrimSuffix(baseURL, "/"),
		client:   restyClient,
		upstream: upstream,
	}
}

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func (c *Client) doRequest(ctx context.Context, token, method, path string, params map[string]string, body []byte) (
	[]byte, error,
) {
	parsedURL, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	start := time.Now()

	req := c.client.R().
		SetContext(ctx).
		SetResponseBodyLimit(maxBodySize).
		SetResponseBodyUnlimitedReads(true).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", authHeader(token)).
		SetQueryParams(params)

	if body != nil {
		req = req.SetBody(body).SetHeader("Content-Type", "application/json")
	}

	// Propagate the correlation ID to the upstream NetBox call so records on
	// both sides of the wire can be matched (L9/G11).
	if requestID := domain.RequestIDFromContext(ctx); requestID != "" {
		req = req.SetHeader(requestIDHeader, requestID)
	}

	resp, err := req.Execute(method, parsedURL.String())
	duration := time.Since(start)

	// O03: record upstream metrics (status counter, latency, response size).
	// The observe method is nil-safe when metrics are not configured.
	status := statusCodeOf(resp)
	var size int
	if resp != nil {
		size = len(resp.Bytes())
	}
	c.upstream.observe(status, duration, size)

	// Emit an outbound per-request log record tagged with the same request_id.
	logger := getLogger()
	entry := logger.WithFields(logrus.Fields{
		"method":   method,
		"path":     parsedURL.Path,
		"duration": duration,
		"status":   statusCodeOf(resp),
	})
	if requestID := domain.RequestIDFromContext(ctx); requestID != "" {
		entry = entry.WithField("request_id", requestID)
	}
	entry.Debug("netbox request")

	if err != nil {
		if errors.Is(err, resty.ErrReadExceedsThresholdLimit) {
			return nil, fmt.Errorf(
				"response body exceeds %d bytes (got at least %d bytes), truncation detected",
				maxBodySize, maxBodySize,
			)
		}

		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized: invalid token or insufficient permissions")
	}
	if resp.StatusCode() == http.StatusNotFound {
		return nil, ErrNotFound
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

// authHeader builds the Authorization header value for the outbound NetBox
// request. NetBox v2 personal access tokens (format nbt_<key>.<plaintext>) are
// accepted via the Bearer scheme, while legacy v1 plaintext tokens only work
// with the Token scheme (NetBox <=4.4 silently ignores Bearer and treats the
// request as anonymous). Auto-detection by the nbt_ prefix lets the same client
// authenticate against both NetBox <=4.4 and >=4.5. Empty and non-prefixed
// tokens are treated as v1 (Token).
func authHeader(token string) string {
	if strings.HasPrefix(token, "nbt_") {
		return "Bearer " + token
	}

	return "Token " + token
}

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
	body, err := c.doRequest(ctx, token, http.MethodDelete, path, nil, nil)
	// Idempotent delete (N25): a 404 means the object is already gone, which
	// from the caller's perspective is a successful no-op. Treating it as such
	// makes repeated delete calls safe to retry (matching IdempotentHint: true).
	if errors.Is(err, ErrNotFound) {
		return nil, nil
	}

	return body, err
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
		if nb, ok := errors.AsType[*netboxError](err); ok {
			return mapWriteStatus(nb.StatusCode, nb.Body)
		}

		return nil, err
	}

	return mapWriteStatus(successStatus, body)
}

// ListSites returns a paginated list of sites.
func (c *Client) ListSites(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.Site], error,
) {
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

// ListRegions returns a paginated list of regions.
func (c *Client) ListRegions(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.Region], error,
) {
	raw, err := c.get(ctx, token, "/api/dcim/regions/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireRegion]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal regions: %w", err)
	}

	return convertPaginated(&wireResp, wireRegionToDomain), nil
}

// ListTenantGroups returns a paginated list of tenant groups.
func (c *Client) ListTenantGroups(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.TenantGroup], error,
) {
	raw, err := c.get(ctx, token, "/api/tenancy/tenant-groups/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireTenantGroup]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tenant groups: %w", err)
	}

	return convertPaginated(&wireResp, wireTenantGroupToDomain), nil
}

// ListDevices returns a paginated list of devices.
func (c *Client) ListDevices(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.Device], error,
) {
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
func (c *Client) ListIPAddresses(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.IPAddress], error,
) {
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
func (c *Client) ListPrefixes(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.Prefix], error,
) {
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
func (c *Client) ListVLANs(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.VLAN], error,
) {
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
func (c *Client) ListVirtualMachines(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.VirtualMachine], error,
) {
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
func (c *Client) ListClusters(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.Cluster], error,
) {
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
func (c *Client) ListCircuits(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.Circuit], error,
) {
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
func (c *Client) ListCircuitTerminations(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.CircuitTermination], error,
) {
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
func (c *Client) ListCables(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.Cable], error,
) {
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
func (c *Client) ListRacks(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.Rack], error,
) {
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
func (c *Client) ListInterfaces(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.Interface], error,
) {
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
func (c *Client) ListVMInterfaces(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.VMInterface], error,
) {
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

// ListSiteGroups returns a paginated list of site groups.
func (c *Client) ListSiteGroups(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.SiteGroup], error,
) {
	raw, err := c.get(ctx, token, "/api/dcim/site-groups/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireSiteGroup]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal site groups: %w", err)
	}

	return convertPaginated(&wireResp, wireSiteGroupToDomain), nil
}

// ListRackGroups returns a paginated list of rack groups.
func (c *Client) ListRackGroups(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.RackGroup], error,
) {
	raw, err := c.get(ctx, token, "/api/dcim/rack-groups/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireRackGroup]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rack groups: %w", err)
	}

	return convertPaginated(&wireResp, wireRackGroupToDomain), nil
}

// ListRackRoles returns a paginated list of rack roles.
func (c *Client) ListRackRoles(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.RackRole], error,
) {
	raw, err := c.get(ctx, token, "/api/dcim/rack-roles/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireRackRole]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rack roles: %w", err)
	}

	return convertPaginated(&wireResp, wireRackRoleToDomain), nil
}

// ListDeviceRoles returns a paginated list of device roles.
func (c *Client) ListDeviceRoles(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.DeviceRole], error,
) {
	raw, err := c.get(ctx, token, "/api/dcim/device-roles/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireDeviceRole]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal device roles: %w", err)
	}

	return convertPaginated(&wireResp, wireDeviceRoleToDomain), nil
}

// ListPlatforms returns a paginated list of platforms.
func (c *Client) ListPlatforms(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.Platform], error,
) {
	raw, err := c.get(ctx, token, "/api/dcim/platforms/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WirePlatform]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal platforms: %w", err)
	}

	return convertPaginated(&wireResp, wirePlatformToDomain), nil
}

// ListVirtualChassis returns a paginated list of virtual chassis.
func (c *Client) ListVirtualChassis(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.VirtualChassis], error,
) {
	raw, err := c.get(ctx, token, "/api/dcim/virtual-chassis/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireVirtualChassis]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal virtual chassis: %w", err)
	}

	return convertPaginated(&wireResp, wireVirtualChassisToDomain), nil
}

// ListPowerPanels returns a paginated list of power panels.
func (c *Client) ListPowerPanels(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.PowerPanel], error,
) {
	raw, err := c.get(ctx, token, "/api/dcim/power-panels/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WirePowerPanel]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal power panels: %w", err)
	}

	return convertPaginated(&wireResp, wirePowerPanelToDomain), nil
}

// ListInventoryItemRoles returns a paginated list of inventory item roles.
func (c *Client) ListInventoryItemRoles(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.InventoryItemRole], error,
) {
	raw, err := c.get(ctx, token, "/api/dcim/inventory-item-roles/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireInventoryItemRole]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal inventory item roles: %w", err)
	}

	return convertPaginated(&wireResp, wireInventoryItemRoleToDomain), nil
}

// ListFrontPorts returns a paginated list of front ports.
func (c *Client) ListFrontPorts(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.FrontPort], error,
) {
	raw, err := c.get(ctx, token, "/api/dcim/front-ports/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireFrontPort]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal front ports: %w", err)
	}

	return convertPaginated(&wireResp, wireFrontPortToDomain), nil
}

// ListRearPorts returns a paginated list of rear ports.
func (c *Client) ListRearPorts(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.RearPort], error,
) {
	raw, err := c.get(ctx, token, "/api/dcim/rear-ports/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireRearPort]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rear ports: %w", err)
	}

	return convertPaginated(&wireResp, wireRearPortToDomain), nil
}

// ListConsolePorts returns a paginated list of console ports.
func (c *Client) ListConsolePorts(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.ConsolePort], error,
) {
	raw, err := c.get(ctx, token, "/api/dcim/console-ports/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireConsolePort]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal console ports: %w", err)
	}

	return convertPaginated(&wireResp, wireConsolePortToDomain), nil
}

// ListConsoleServerPorts returns a paginated list of console server ports.
func (c *Client) ListConsoleServerPorts(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.ConsoleServerPort], error,
) {
	raw, err := c.get(ctx, token, "/api/dcim/console-server-ports/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireConsoleServerPort]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal console server ports: %w", err)
	}

	return convertPaginated(&wireResp, wireConsoleServerPortToDomain), nil
}

// ListPowerPorts returns a paginated list of power ports.
func (c *Client) ListPowerPorts(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.PowerPort], error,
) {
	raw, err := c.get(ctx, token, "/api/dcim/power-ports/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WirePowerPort]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal power ports: %w", err)
	}

	return convertPaginated(&wireResp, wirePowerPortToDomain), nil
}

// ListPowerOutlets returns a paginated list of power outlets.
func (c *Client) ListPowerOutlets(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.PowerOutlet], error,
) {
	raw, err := c.get(ctx, token, "/api/dcim/power-outlets/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WirePowerOutlet]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal power outlets: %w", err)
	}

	return convertPaginated(&wireResp, wirePowerOutletToDomain), nil
}

// ListPowerFeeds returns a paginated list of power feeds.
func (c *Client) ListPowerFeeds(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.PowerFeed], error,
) {
	raw, err := c.get(ctx, token, "/api/dcim/power-feeds/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WirePowerFeed]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal power feeds: %w", err)
	}

	return convertPaginated(&wireResp, wirePowerFeedToDomain), nil
}

// ListDeviceBays returns a paginated list of device bays.
func (c *Client) ListDeviceBays(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.DeviceBay], error,
) {
	raw, err := c.get(ctx, token, "/api/dcim/device-bays/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireDeviceBay]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal device bays: %w", err)
	}

	return convertPaginated(&wireResp, wireDeviceBayToDomain), nil
}

// ListRackReservations returns a paginated list of rack reservations.
func (c *Client) ListRackReservations(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.RackReservation], error,
) {
	raw, err := c.get(ctx, token, "/api/dcim/rack-reservations/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireRackReservation]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rack reservations: %w", err)
	}

	return convertPaginated(&wireResp, wireRackReservationToDomain), nil
}

// ListInventoryItems returns a paginated list of inventory items.
func (c *Client) ListInventoryItems(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.InventoryItem], error,
) {
	raw, err := c.get(ctx, token, "/api/dcim/inventory-items/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireInventoryItem]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal inventory items: %w", err)
	}

	return convertPaginated(&wireResp, wireInventoryItemToDomain), nil
}

// ListMACAddresses returns a paginated list of MAC addresses.
func (c *Client) ListMACAddresses(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.MACAddress], error,
) {
	raw, err := c.get(ctx, token, "/api/dcim/mac-addresses/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireMACAddress]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal MAC addresses: %w", err)
	}

	return convertPaginated(&wireResp, wireMACAddressToDomain), nil
}

// ListRIRs returns a paginated list of RIRs.
func (c *Client) ListRIRs(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.RIR], error,
) {
	raw, err := c.get(ctx, token, "/api/ipam/rirs/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireRIR]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal RIRs: %w", err)
	}

	return convertPaginated(&wireResp, wireRIRToDomain), nil
}

// ListAggregates returns a paginated list of aggregates.
func (c *Client) ListAggregates(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.Aggregate], error,
) {
	raw, err := c.get(ctx, token, "/api/ipam/aggregates/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireAggregate]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal aggregates: %w", err)
	}

	return convertPaginated(&wireResp, wireAggregateToDomain), nil
}

// ListIPRanges returns a paginated list of IP ranges.
func (c *Client) ListIPRanges(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.IPRange], error,
) {
	raw, err := c.get(ctx, token, "/api/ipam/ip-ranges/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireIPRange]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal IP ranges: %w", err)
	}

	return convertPaginated(&wireResp, wireIPRangeToDomain), nil
}

// ListASNs returns a paginated list of ASNs.
func (c *Client) ListASNs(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.ASN], error,
) {
	raw, err := c.get(ctx, token, "/api/ipam/asns/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireASN]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ASNs: %w", err)
	}

	return convertPaginated(&wireResp, wireASNToDomain), nil
}

// ListASNRanges returns a paginated list of ASN ranges.
func (c *Client) ListASNRanges(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.ASNRange], error,
) {
	raw, err := c.get(ctx, token, "/api/ipam/asn-ranges/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireASNRange]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ASN ranges: %w", err)
	}

	return convertPaginated(&wireResp, wireASNRangeToDomain), nil
}

// ListRouteTargets returns a paginated list of route targets.
func (c *Client) ListRouteTargets(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.RouteTarget], error,
) {
	raw, err := c.get(ctx, token, "/api/ipam/route-targets/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireRouteTarget]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal route targets: %w", err)
	}

	return convertPaginated(&wireResp, wireRouteTargetToDomain), nil
}

// ListServices returns a paginated list of services.
func (c *Client) ListServices(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.Service], error,
) {
	raw, err := c.get(ctx, token, "/api/ipam/services/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireService]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal services: %w", err)
	}

	return convertPaginated(&wireResp, wireServiceToDomain), nil
}

// ListServiceTemplates returns a paginated list of service templates.
func (c *Client) ListServiceTemplates(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.ServiceTemplate], error,
) {
	raw, err := c.get(ctx, token, "/api/ipam/service-templates/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireServiceTemplate]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal service templates: %w", err)
	}

	return convertPaginated(&wireResp, wireServiceTemplateToDomain), nil
}

// ListFHRPGroups returns a paginated list of FHRP groups.
func (c *Client) ListFHRPGroups(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.FHRPGroup], error,
) {
	raw, err := c.get(ctx, token, "/api/ipam/fhrp-groups/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireFHRPGroup]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal FHRP groups: %w", err)
	}

	return convertPaginated(&wireResp, wireFHRPGroupToDomain), nil
}

// ListFHRPGroupAssignments returns a paginated list of FHRP group assignments.
func (c *Client) ListFHRPGroupAssignments(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.FHRPGroupAssignment], error,
) {
	raw, err := c.get(ctx, token, "/api/ipam/fhrp-group-assignments/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireFHRPGroupAssignment]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal FHRP group assignments: %w", err)
	}

	return convertPaginated(&wireResp, wireFHRPGroupAssignmentToDomain), nil
}

// ListVLANTranslationPolicies returns a paginated list of VLAN translation
// policies.
func (c *Client) ListVLANTranslationPolicies(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.VLANTranslationPolicy], error,
) {
	raw, err := c.get(ctx, token, "/api/ipam/vlan-translation-policies/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireVLANTranslationPolicy]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal VLAN translation policies: %w", err)
	}

	return convertPaginated(&wireResp, wireVLANTranslationPolicyToDomain), nil
}

// ListVLANTranslationRules returns a paginated list of VLAN translation rules.
func (c *Client) ListVLANTranslationRules(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.VLANTranslationRule], error,
) {
	raw, err := c.get(ctx, token, "/api/ipam/vlan-translation-rules/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireVLANTranslationRule]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal VLAN translation rules: %w", err)
	}

	return convertPaginated(&wireResp, wireVLANTranslationRuleToDomain), nil
}

// ListContactGroups returns a paginated list of contact groups.
func (c *Client) ListContactGroups(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.ContactGroup], error,
) {
	raw, err := c.get(ctx, token, "/api/tenancy/contact-groups/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireContactGroup]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal contact groups: %w", err)
	}

	return convertPaginated(&wireResp, wireContactGroupToDomain), nil
}

// ListContactRoles returns a paginated list of contact roles.
func (c *Client) ListContactRoles(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.ContactRole], error,
) {
	raw, err := c.get(ctx, token, "/api/tenancy/contact-roles/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireContactRole]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal contact roles: %w", err)
	}

	return convertPaginated(&wireResp, wireContactRoleToDomain), nil
}

// ListVirtualDisks returns a paginated list of virtual disks.
func (c *Client) ListVirtualDisks(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.VirtualDisk], error,
) {
	raw, err := c.get(ctx, token, "/api/virtualization/virtual-disks/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireVirtualDisk]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal virtual disks: %w", err)
	}

	return convertPaginated(&wireResp, wireVirtualDiskToDomain), nil
}

// ListVirtualMachineTypes returns a paginated list of virtual machine types.
func (c *Client) ListVirtualMachineTypes(ctx context.Context, token string, params map[string]string) (
	*domain.PaginatedResponse[domain.VirtualMachineType], error,
) {
	raw, err := c.get(ctx, token, "/api/virtualization/virtual-machine-types/", params)
	if err != nil {
		return nil, err
	}
	var wireResp domain.PaginatedResponse[WireVirtualMachineType]
	if err := json.Unmarshal(raw, &wireResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal virtual machine types: %w", err)
	}

	return convertPaginated(&wireResp, wireVirtualMachineTypeToDomain), nil
}

// GetObject retrieves a single object by its type and ID.
func (c *Client) GetObject(ctx context.Context, token, objectType string, objectID int, params map[string]string) (
	domain.RawObject, error,
) {
	endpoint, ok := objectTypeToEndpoint[objectType]
	if !ok {
		return nil, fmt.Errorf("unknown object type: %s", objectType)
	}

	path := fmt.Sprintf("%s%d/", endpoint, objectID)
	body, err := c.get(ctx, token, path, params)
	if err != nil {
		return nil, fmt.Errorf("%s #%d: %w", objectType, objectID, err)
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
func (c *Client) UpdateSite(ctx context.Context, token string, objectID int, in domain.SiteWrite) (
	*domain.Site, error,
) {
	payload, err := json.Marshal(siteWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal site update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/sites/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalSite(raw)
}

// DeleteSite deletes a site via DELETE /api/dcim/sites/<id>/.
func (c *Client) DeleteSite(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/sites/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateRegion creates a new region via POST /api/dcim/regions/.
func (c *Client) CreateRegion(ctx context.Context, token string, in domain.RegionWrite) (*domain.Region, error) {
	payload, err := json.Marshal(regionWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal region create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/regions/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalRegion(raw)
}

// UpdateRegion partially updates a region via PATCH /api/dcim/regions/<id>/.
func (c *Client) UpdateRegion(ctx context.Context, token string, objectID int, in domain.RegionWrite) (
	*domain.Region, error,
) {
	payload, err := json.Marshal(regionWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal region update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/regions/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalRegion(raw)
}

// DeleteRegion deletes a region via DELETE /api/dcim/regions/<id>/.
func (c *Client) DeleteRegion(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/regions/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateTenantGroup creates a new tenant group via POST /api/tenancy/tenant-groups/.
func (c *Client) CreateTenantGroup(ctx context.Context, token string, in domain.TenantGroupWrite) (
	*domain.TenantGroup, error,
) {
	payload, err := json.Marshal(tenantGroupWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal tenant group create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/tenancy/tenant-groups/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalTenantGroup(raw)
}

// UpdateTenantGroup partially updates a tenant group via PATCH
// /api/tenancy/tenant-groups/<id>/.
func (c *Client) UpdateTenantGroup(ctx context.Context, token string, objectID int, in domain.TenantGroupWrite) (
	*domain.TenantGroup, error,
) {
	payload, err := json.Marshal(tenantGroupWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal tenant group update request: %w", err)
	}
	path := fmt.Sprintf("/api/tenancy/tenant-groups/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalTenantGroup(raw)
}

// DeleteTenantGroup deletes a tenant group via DELETE
// /api/tenancy/tenant-groups/<id>/.
func (c *Client) DeleteTenantGroup(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/tenancy/tenant-groups/%d/", objectID)
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
func (c *Client) UpdateDevice(ctx context.Context, token string, objectID int, in domain.DeviceWrite) (
	*domain.Device, error,
) {
	payload, err := json.Marshal(deviceWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal device update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/devices/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalDevice(raw)
}

// DeleteDevice deletes a device via DELETE /api/dcim/devices/<id>/.
func (c *Client) DeleteDevice(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/devices/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateIPAddress creates a new IP address via POST /api/ipam/ip-addresses/.
func (c *Client) CreateIPAddress(ctx context.Context, token string, in domain.IPAddressWrite) (
	*domain.IPAddress, error,
) {
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
func (c *Client) UpdateIPAddress(ctx context.Context, token string, objectID int, in domain.IPAddressWrite) (
	*domain.IPAddress, error,
) {
	payload, err := json.Marshal(ipAddressWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal IP address update request: %w", err)
	}
	path := fmt.Sprintf("/api/ipam/ip-addresses/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalIPAddress(raw)
}

// DeleteIPAddress deletes an IP address via DELETE /api/ipam/ip-addresses/<id>/.
func (c *Client) DeleteIPAddress(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/ipam/ip-addresses/%d/", objectID)
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
func (c *Client) UpdatePrefix(ctx context.Context, token string, objectID int, in domain.PrefixWrite) (
	*domain.Prefix, error,
) {
	payload, err := json.Marshal(prefixWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal prefix update request: %w", err)
	}
	path := fmt.Sprintf("/api/ipam/prefixes/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalPrefix(raw)
}

// DeletePrefix deletes a prefix via DELETE /api/ipam/prefixes/<id>/.
func (c *Client) DeletePrefix(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/ipam/prefixes/%d/", objectID)
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
func (c *Client) UpdateVLAN(ctx context.Context, token string, objectID int, in domain.VLANWrite) (
	*domain.VLAN, error,
) {
	payload, err := json.Marshal(vlanWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal VLAN update request: %w", err)
	}
	path := fmt.Sprintf("/api/ipam/vlans/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalVLAN(raw)
}

// DeleteVLAN deletes a VLAN via DELETE /api/ipam/vlans/<id>/.
func (c *Client) DeleteVLAN(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/ipam/vlans/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateVirtualMachine creates a new virtual machine via POST
// /api/virtualization/virtual-machines/.
func (c *Client) CreateVirtualMachine(ctx context.Context, token string, in domain.VirtualMachineWrite) (
	*domain.VirtualMachine, error,
) {
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
func (c *Client) UpdateVirtualMachine(ctx context.Context, token string, objectID int, in domain.VirtualMachineWrite) (
	*domain.VirtualMachine, error,
) {
	payload, err := json.Marshal(virtualMachineWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal virtual machine update request: %w", err)
	}
	path := fmt.Sprintf("/api/virtualization/virtual-machines/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalVirtualMachine(raw)
}

// DeleteVirtualMachine deletes a virtual machine via DELETE
// /api/virtualization/virtual-machines/<id>/.
func (c *Client) DeleteVirtualMachine(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/virtualization/virtual-machines/%d/", objectID)
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
func (c *Client) UpdateCluster(ctx context.Context, token string, objectID int, in domain.ClusterWrite) (
	*domain.Cluster, error,
) {
	payload, err := json.Marshal(clusterWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal cluster update request: %w", err)
	}
	path := fmt.Sprintf("/api/virtualization/clusters/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalCluster(raw)
}

// DeleteCluster deletes a cluster via DELETE /api/virtualization/clusters/<id>/.
func (c *Client) DeleteCluster(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/virtualization/clusters/%d/", objectID)
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
func (c *Client) UpdateCircuit(ctx context.Context, token string, objectID int, in domain.CircuitWrite) (
	*domain.Circuit, error,
) {
	payload, err := json.Marshal(circuitWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal circuit update request: %w", err)
	}
	path := fmt.Sprintf("/api/circuits/circuits/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalCircuit(raw)
}

// DeleteCircuit deletes a circuit via DELETE /api/circuits/circuits/<id>/.
func (c *Client) DeleteCircuit(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/circuits/circuits/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateRack creates a new rack via POST /api/dcim/racks/.
func (c *Client) CreateRack(ctx context.Context, token string, in domain.RackWrite) (*domain.Rack, error) {
	payload, err := json.Marshal(rackWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal rack create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/racks/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalRack(raw)
}

// UpdateRack partially updates a rack via PATCH /api/dcim/racks/<id>/.
func (c *Client) UpdateRack(ctx context.Context, token string, objectID int, in domain.RackWrite) (
	*domain.Rack, error,
) {
	payload, err := json.Marshal(rackWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal rack update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/racks/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalRack(raw)
}

// DeleteRack deletes a rack via DELETE /api/dcim/racks/<id>/.
func (c *Client) DeleteRack(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/racks/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateInterface creates a new interface via POST /api/dcim/interfaces/.
func (c *Client) CreateInterface(ctx context.Context, token string, in domain.InterfaceWrite) (
	*domain.Interface, error,
) {
	payload, err := json.Marshal(interfaceWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal interface create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/interfaces/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalInterface(raw)
}

// UpdateInterface partially updates an interface via PATCH /api/dcim/interfaces/<id>/.
func (c *Client) UpdateInterface(ctx context.Context, token string, objectID int, in domain.InterfaceWrite) (
	*domain.Interface, error,
) {
	payload, err := json.Marshal(interfaceWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal interface update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/interfaces/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalInterface(raw)
}

// DeleteInterface deletes an interface via DELETE /api/dcim/interfaces/<id>/.
func (c *Client) DeleteInterface(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/interfaces/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateCircuitTermination creates a new circuit termination via POST
// /api/circuits/circuit-terminations/.
func (c *Client) CreateCircuitTermination(ctx context.Context, token string, in domain.CircuitTerminationWrite) (
	*domain.CircuitTermination, error,
) {
	payload, err := json.Marshal(circuitTerminationWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal circuit termination create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/circuits/circuit-terminations/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalCircuitTermination(raw)
}

// UpdateCircuitTermination partially updates a circuit termination via PATCH
// /api/circuits/circuit-terminations/<id>/.
func (c *Client) UpdateCircuitTermination(
	ctx context.Context, token string, objectID int, in domain.CircuitTerminationWrite,
) (
	*domain.CircuitTermination, error,
) {
	payload, err := json.Marshal(circuitTerminationWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal circuit termination update request: %w", err)
	}
	path := fmt.Sprintf("/api/circuits/circuit-terminations/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalCircuitTermination(raw)
}

// DeleteCircuitTermination deletes a circuit termination via DELETE
// /api/circuits/circuit-terminations/<id>/.
func (c *Client) DeleteCircuitTermination(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/circuits/circuit-terminations/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateCable creates a new cable via POST /api/dcim/cables/.
func (c *Client) CreateCable(ctx context.Context, token string, in domain.CableWrite) (*domain.Cable, error) {
	payload, err := json.Marshal(cableWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal cable create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/cables/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalCable(raw)
}

// UpdateCable partially updates a cable via PATCH /api/dcim/cables/<id>/.
func (c *Client) UpdateCable(ctx context.Context, token string, objectID int, in domain.CableWrite) (
	*domain.Cable, error,
) {
	payload, err := json.Marshal(cableWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal cable update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/cables/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalCable(raw)
}

// DeleteCable deletes a cable via DELETE /api/dcim/cables/<id>/.
func (c *Client) DeleteCable(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/cables/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateVMInterface creates a new VM interface via POST
// /api/virtualization/interfaces/.
func (c *Client) CreateVMInterface(ctx context.Context, token string, in domain.VMInterfaceWrite) (
	*domain.VMInterface, error,
) {
	payload, err := json.Marshal(vmInterfaceWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal vm interface create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/virtualization/interfaces/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalVMInterface(raw)
}

// UpdateVMInterface partially updates a VM interface via PATCH
// /api/virtualization/interfaces/<id>/.
func (c *Client) UpdateVMInterface(ctx context.Context, token string, objectID int, in domain.VMInterfaceWrite) (
	*domain.VMInterface, error,
) {
	payload, err := json.Marshal(vmInterfaceWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal vm interface update request: %w", err)
	}
	path := fmt.Sprintf("/api/virtualization/interfaces/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalVMInterface(raw)
}

// DeleteVMInterface deletes a VM interface via DELETE
// /api/virtualization/interfaces/<id>/.
func (c *Client) DeleteVMInterface(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/virtualization/interfaces/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateProvider creates a new provider via POST /api/circuits/providers/.
func (c *Client) CreateProvider(ctx context.Context, token string, in domain.ProviderWrite) (*domain.Provider, error) {
	payload, err := json.Marshal(providerWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal provider create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/circuits/providers/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalProvider(raw)
}

// UpdateProvider partially updates a provider via PATCH
// /api/circuits/providers/<id>/.
func (c *Client) UpdateProvider(ctx context.Context, token string, objectID int, in domain.ProviderWrite) (
	*domain.Provider, error,
) {
	payload, err := json.Marshal(providerWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal provider update request: %w", err)
	}
	path := fmt.Sprintf("/api/circuits/providers/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalProvider(raw)
}

// DeleteProvider deletes a provider via DELETE /api/circuits/providers/<id>/.
func (c *Client) DeleteProvider(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/circuits/providers/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateTenant creates a new tenant via POST /api/tenancy/tenants/.
func (c *Client) CreateTenant(ctx context.Context, token string, in domain.TenantWrite) (*domain.Tenant, error) {
	payload, err := json.Marshal(tenantWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal tenant create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/tenancy/tenants/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalTenant(raw)
}

// UpdateTenant partially updates a tenant via PATCH /api/tenancy/tenants/<id>/.
func (c *Client) UpdateTenant(ctx context.Context, token string, objectID int, in domain.TenantWrite) (
	*domain.Tenant, error,
) {
	payload, err := json.Marshal(tenantWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal tenant update request: %w", err)
	}
	path := fmt.Sprintf("/api/tenancy/tenants/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalTenant(raw)
}

// DeleteTenant deletes a tenant via DELETE /api/tenancy/tenants/<id>/.
func (c *Client) DeleteTenant(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/tenancy/tenants/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateManufacturer creates a new manufacturer via POST /api/dcim/manufacturers/.
func (c *Client) CreateManufacturer(ctx context.Context, token string, in domain.ManufacturerWrite) (
	*domain.Manufacturer, error,
) {
	payload, err := json.Marshal(manufacturerWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal manufacturer create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/manufacturers/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalManufacturer(raw)
}

// UpdateManufacturer partially updates a manufacturer via PATCH
// /api/dcim/manufacturers/<id>/.
func (c *Client) UpdateManufacturer(ctx context.Context, token string, objectID int, in domain.ManufacturerWrite) (
	*domain.Manufacturer, error,
) {
	payload, err := json.Marshal(manufacturerWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal manufacturer update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/manufacturers/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalManufacturer(raw)
}

// DeleteManufacturer deletes a manufacturer via DELETE /api/dcim/manufacturers/<id>/.
func (c *Client) DeleteManufacturer(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/manufacturers/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateDeviceType creates a new device type via POST /api/dcim/device-types/.
func (c *Client) CreateDeviceType(ctx context.Context, token string, in domain.DeviceTypeWrite) (
	*domain.DeviceType, error,
) {
	payload, err := json.Marshal(deviceTypeWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal device type create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/device-types/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalDeviceType(raw)
}

// UpdateDeviceType partially updates a device type via PATCH
// /api/dcim/device-types/<id>/.
func (c *Client) UpdateDeviceType(ctx context.Context, token string, objectID int, in domain.DeviceTypeWrite) (
	*domain.DeviceType, error,
) {
	payload, err := json.Marshal(deviceTypeWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal device type update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/device-types/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalDeviceType(raw)
}

// DeleteDeviceType deletes a device type via DELETE /api/dcim/device-types/<id>/.
func (c *Client) DeleteDeviceType(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/device-types/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateLocation creates a new site location via POST /api/dcim/locations/.
func (c *Client) CreateLocation(ctx context.Context, token string, in domain.LocationWrite) (*domain.Location, error) {
	payload, err := json.Marshal(locationWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal location create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/locations/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalLocation(raw)
}

// UpdateLocation partially updates a site location via PATCH
// /api/dcim/locations/<id>/.
func (c *Client) UpdateLocation(ctx context.Context, token string, objectID int, in domain.LocationWrite) (
	*domain.Location, error,
) {
	payload, err := json.Marshal(locationWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal location update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/locations/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalLocation(raw)
}

// DeleteLocation deletes a site location via DELETE /api/dcim/locations/<id>/.
func (c *Client) DeleteLocation(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/locations/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateClusterType creates a new cluster type via POST
// /api/virtualization/cluster-types/.
func (c *Client) CreateClusterType(ctx context.Context, token string, in domain.ClusterTypeWrite) (
	*domain.ClusterType, error,
) {
	payload, err := json.Marshal(clusterTypeWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal cluster type create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/virtualization/cluster-types/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalClusterType(raw)
}

// UpdateClusterType partially updates a cluster type via PATCH
// /api/virtualization/cluster-types/<id>/.
func (c *Client) UpdateClusterType(ctx context.Context, token string, objectID int, in domain.ClusterTypeWrite) (
	*domain.ClusterType, error,
) {
	payload, err := json.Marshal(clusterTypeWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal cluster type update request: %w", err)
	}
	path := fmt.Sprintf("/api/virtualization/cluster-types/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalClusterType(raw)
}

// DeleteClusterType deletes a cluster type via DELETE
// /api/virtualization/cluster-types/<id>/.
func (c *Client) DeleteClusterType(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/virtualization/cluster-types/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateClusterGroup creates a new cluster group via POST
// /api/virtualization/cluster-groups/.
func (c *Client) CreateClusterGroup(ctx context.Context, token string, in domain.ClusterGroupWrite) (
	*domain.ClusterGroup, error,
) {
	payload, err := json.Marshal(clusterGroupWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal cluster group create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/virtualization/cluster-groups/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalClusterGroup(raw)
}

// UpdateClusterGroup partially updates a cluster group via PATCH
// /api/virtualization/cluster-groups/<id>/.
func (c *Client) UpdateClusterGroup(ctx context.Context, token string, objectID int, in domain.ClusterGroupWrite) (
	*domain.ClusterGroup, error,
) {
	payload, err := json.Marshal(clusterGroupWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal cluster group update request: %w", err)
	}
	path := fmt.Sprintf("/api/virtualization/cluster-groups/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalClusterGroup(raw)
}

// DeleteClusterGroup deletes a cluster group via DELETE
// /api/virtualization/cluster-groups/<id>/.
func (c *Client) DeleteClusterGroup(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/virtualization/cluster-groups/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateCircuitType creates a new circuit type via POST
// /api/circuits/circuit-types/.
func (c *Client) CreateCircuitType(ctx context.Context, token string, in domain.CircuitTypeWrite) (
	*domain.CircuitType, error,
) {
	payload, err := json.Marshal(circuitTypeWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal circuit type create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/circuits/circuit-types/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalCircuitType(raw)
}

// UpdateCircuitType partially updates a circuit type via PATCH
// /api/circuits/circuit-types/<id>/.
func (c *Client) UpdateCircuitType(ctx context.Context, token string, objectID int, in domain.CircuitTypeWrite) (
	*domain.CircuitType, error,
) {
	payload, err := json.Marshal(circuitTypeWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal circuit type update request: %w", err)
	}
	path := fmt.Sprintf("/api/circuits/circuit-types/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalCircuitType(raw)
}

// DeleteCircuitType deletes a circuit type via DELETE
// /api/circuits/circuit-types/<id>/.
func (c *Client) DeleteCircuitType(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/circuits/circuit-types/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateVrf creates a new VRF via POST /api/ipam/vrfs/.
func (c *Client) CreateVrf(ctx context.Context, token string, in domain.VrfWrite) (*domain.Vrf, error) {
	payload, err := json.Marshal(vrfWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal vrf create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/ipam/vrfs/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalVrf(raw)
}

// UpdateVrf partially updates a VRF via PATCH /api/ipam/vrfs/<id>/.
func (c *Client) UpdateVrf(ctx context.Context, token string, objectID int, in domain.VrfWrite) (*domain.Vrf, error) {
	payload, err := json.Marshal(vrfWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal vrf update request: %w", err)
	}
	path := fmt.Sprintf("/api/ipam/vrfs/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalVrf(raw)
}

// DeleteVrf deletes a VRF via DELETE /api/ipam/vrfs/<id>/.
func (c *Client) DeleteVrf(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/ipam/vrfs/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateVlanGroup creates a new VLAN group via POST /api/ipam/vlan-groups/.
func (c *Client) CreateVlanGroup(ctx context.Context, token string, in domain.VlanGroupWrite) (
	*domain.VlanGroup, error,
) {
	payload, err := json.Marshal(vlanGroupWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal vlan group create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/ipam/vlan-groups/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalVlanGroup(raw)
}

// UpdateVlanGroup partially updates a VLAN group via PATCH
// /api/ipam/vlan-groups/<id>/.
func (c *Client) UpdateVlanGroup(ctx context.Context, token string, objectID int, in domain.VlanGroupWrite) (
	*domain.VlanGroup, error,
) {
	payload, err := json.Marshal(vlanGroupWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal vlan group update request: %w", err)
	}
	path := fmt.Sprintf("/api/ipam/vlan-groups/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalVlanGroup(raw)
}

// DeleteVlanGroup deletes a VLAN group via DELETE /api/ipam/vlan-groups/<id>/.
func (c *Client) DeleteVlanGroup(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/ipam/vlan-groups/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateRole creates a new IPAM role via POST /api/ipam/roles/.
func (c *Client) CreateRole(ctx context.Context, token string, in domain.RoleWrite) (*domain.Role, error) {
	payload, err := json.Marshal(roleWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal role create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/ipam/roles/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalRole(raw)
}

// UpdateRole partially updates an IPAM role via PATCH /api/ipam/roles/<id>/.
func (c *Client) UpdateRole(ctx context.Context, token string, objectID int, in domain.RoleWrite) (
	*domain.Role, error,
) {
	payload, err := json.Marshal(roleWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal role update request: %w", err)
	}
	path := fmt.Sprintf("/api/ipam/roles/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalRole(raw)
}

// DeleteRole deletes an IPAM role via DELETE /api/ipam/roles/<id>/.
func (c *Client) DeleteRole(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/ipam/roles/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateContact creates a new tenancy contact via POST /api/tenancy/contacts/.
func (c *Client) CreateContact(ctx context.Context, token string, in domain.ContactWrite) (*domain.Contact, error) {
	payload, err := json.Marshal(contactWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal contact create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/tenancy/contacts/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalContact(raw)
}

// UpdateContact partially updates a tenancy contact via
// PATCH /api/tenancy/contacts/<id>/.
func (c *Client) UpdateContact(ctx context.Context, token string, objectID int, in domain.ContactWrite) (
	*domain.Contact, error,
) {
	payload, err := json.Marshal(contactWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal contact update request: %w", err)
	}
	path := fmt.Sprintf("/api/tenancy/contacts/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalContact(raw)
}

// DeleteContact deletes a tenancy contact via DELETE /api/tenancy/contacts/<id>/.
func (c *Client) DeleteContact(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/tenancy/contacts/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateSiteGroup creates a new site group via POST /api/dcim/site-groups/.
func (c *Client) CreateSiteGroup(ctx context.Context, token string, in domain.SiteGroupWrite) (
	*domain.SiteGroup, error,
) {
	payload, err := json.Marshal(siteGroupWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal site group create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/site-groups/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalSiteGroup(raw)
}

// UpdateSiteGroup partially updates a site group via PATCH
// /api/dcim/site-groups/<id>/.
func (c *Client) UpdateSiteGroup(ctx context.Context, token string, objectID int, in domain.SiteGroupWrite) (
	*domain.SiteGroup, error,
) {
	payload, err := json.Marshal(siteGroupWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal site group update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/site-groups/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalSiteGroup(raw)
}

// DeleteSiteGroup deletes a site group via DELETE /api/dcim/site-groups/<id>/.
func (c *Client) DeleteSiteGroup(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/site-groups/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateRackGroup creates a new rack group via POST /api/dcim/rack-groups/.
func (c *Client) CreateRackGroup(ctx context.Context, token string, in domain.RackGroupWrite) (
	*domain.RackGroup, error,
) {
	payload, err := json.Marshal(rackGroupWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal rack group create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/rack-groups/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalRackGroup(raw)
}

// UpdateRackGroup partially updates a rack group via PATCH
// /api/dcim/rack-groups/<id>/.
func (c *Client) UpdateRackGroup(ctx context.Context, token string, objectID int, in domain.RackGroupWrite) (
	*domain.RackGroup, error,
) {
	payload, err := json.Marshal(rackGroupWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal rack group update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/rack-groups/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalRackGroup(raw)
}

// DeleteRackGroup deletes a rack group via DELETE /api/dcim/rack-groups/<id>/.
func (c *Client) DeleteRackGroup(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/rack-groups/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateRackRole creates a new rack role via POST /api/dcim/rack-roles/.
func (c *Client) CreateRackRole(ctx context.Context, token string, in domain.RackRoleWrite) (*domain.RackRole, error) {
	payload, err := json.Marshal(rackRoleWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal rack role create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/rack-roles/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalRackRole(raw)
}

// UpdateRackRole partially updates a rack role via PATCH
// /api/dcim/rack-roles/<id>/.
func (c *Client) UpdateRackRole(ctx context.Context, token string, objectID int, in domain.RackRoleWrite) (
	*domain.RackRole, error,
) {
	payload, err := json.Marshal(rackRoleWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal rack role update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/rack-roles/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalRackRole(raw)
}

// DeleteRackRole deletes a rack role via DELETE /api/dcim/rack-roles/<id>/.
func (c *Client) DeleteRackRole(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/rack-roles/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateDeviceRole creates a new device role via POST /api/dcim/device-roles/.
func (c *Client) CreateDeviceRole(ctx context.Context, token string, in domain.DeviceRoleWrite) (
	*domain.DeviceRole, error,
) {
	payload, err := json.Marshal(deviceRoleWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal device role create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/device-roles/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalDeviceRole(raw)
}

// UpdateDeviceRole partially updates a device role via PATCH
// /api/dcim/device-roles/<id>/.
func (c *Client) UpdateDeviceRole(ctx context.Context, token string, objectID int, in domain.DeviceRoleWrite) (
	*domain.DeviceRole, error,
) {
	payload, err := json.Marshal(deviceRoleWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal device role update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/device-roles/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalDeviceRole(raw)
}

// DeleteDeviceRole deletes a device role via DELETE /api/dcim/device-roles/<id>/.
func (c *Client) DeleteDeviceRole(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/device-roles/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreatePlatform creates a new platform via POST /api/dcim/platforms/.
func (c *Client) CreatePlatform(ctx context.Context, token string, in domain.PlatformWrite) (*domain.Platform, error) {
	payload, err := json.Marshal(platformWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal platform create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/platforms/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalPlatform(raw)
}

// UpdatePlatform partially updates a platform via PATCH /api/dcim/platforms/<id>/.
func (c *Client) UpdatePlatform(ctx context.Context, token string, objectID int, in domain.PlatformWrite) (
	*domain.Platform, error,
) {
	payload, err := json.Marshal(platformWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal platform update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/platforms/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalPlatform(raw)
}

// DeletePlatform deletes a platform via DELETE /api/dcim/platforms/<id>/.
func (c *Client) DeletePlatform(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/platforms/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateVirtualChassis creates a new virtual chassis via POST
// /api/dcim/virtual-chassis/.
func (c *Client) CreateVirtualChassis(
	ctx context.Context, token string, in domain.VirtualChassisWrite,
) (*domain.VirtualChassis, error) {
	payload, err := json.Marshal(virtualChassisWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal virtual chassis create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/virtual-chassis/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalVirtualChassis(raw)
}

// UpdateVirtualChassis partially updates a virtual chassis via PATCH
// /api/dcim/virtual-chassis/<id>/.
func (c *Client) UpdateVirtualChassis(
	ctx context.Context, token string, objectID int, in domain.VirtualChassisWrite,
) (*domain.VirtualChassis, error) {
	payload, err := json.Marshal(virtualChassisWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal virtual chassis update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/virtual-chassis/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalVirtualChassis(raw)
}

// DeleteVirtualChassis deletes a virtual chassis via DELETE
// /api/dcim/virtual-chassis/<id>/.
func (c *Client) DeleteVirtualChassis(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/virtual-chassis/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreatePowerPanel creates a new power panel via POST /api/dcim/power-panels/.
func (c *Client) CreatePowerPanel(ctx context.Context, token string, in domain.PowerPanelWrite) (
	*domain.PowerPanel, error,
) {
	payload, err := json.Marshal(powerPanelWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal power panel create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/power-panels/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalPowerPanel(raw)
}

// UpdatePowerPanel partially updates a power panel via PATCH
// /api/dcim/power-panels/<id>/.
func (c *Client) UpdatePowerPanel(ctx context.Context, token string, objectID int, in domain.PowerPanelWrite) (
	*domain.PowerPanel, error,
) {
	payload, err := json.Marshal(powerPanelWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal power panel update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/power-panels/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalPowerPanel(raw)
}

// DeletePowerPanel deletes a power panel via DELETE /api/dcim/power-panels/<id>/.
func (c *Client) DeletePowerPanel(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/power-panels/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateInventoryItemRole creates a new inventory item role via POST
// /api/dcim/inventory-item-roles/.
func (c *Client) CreateInventoryItemRole(
	ctx context.Context, token string, in domain.InventoryItemRoleWrite,
) (*domain.InventoryItemRole, error) {
	payload, err := json.Marshal(inventoryItemRoleWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal inventory item role create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/inventory-item-roles/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalInventoryItemRole(raw)
}

// UpdateInventoryItemRole partially updates an inventory item role via PATCH
// /api/dcim/inventory-item-roles/<id>/.
func (c *Client) UpdateInventoryItemRole(
	ctx context.Context, token string, objectID int, in domain.InventoryItemRoleWrite,
) (*domain.InventoryItemRole, error) {
	payload, err := json.Marshal(inventoryItemRoleWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal inventory item role update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/inventory-item-roles/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalInventoryItemRole(raw)
}

// DeleteInventoryItemRole deletes an inventory item role via DELETE
// /api/dcim/inventory-item-roles/<id>/.
func (c *Client) DeleteInventoryItemRole(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/inventory-item-roles/%d/", objectID)
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

// CreateFrontPort creates a new front port via POST /api/dcim/front-ports/.
func (c *Client) CreateFrontPort(ctx context.Context, token string, in domain.FrontPortWrite) (
	*domain.FrontPort, error,
) {
	payload, err := json.Marshal(frontPortWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal front port create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/front-ports/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalFrontPort(raw)
}

// UpdateFrontPort partially updates a front port via PATCH
// /api/dcim/front-ports/<id>/.
func (c *Client) UpdateFrontPort(ctx context.Context, token string, objectID int, in domain.FrontPortWrite) (
	*domain.FrontPort, error,
) {
	payload, err := json.Marshal(frontPortWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal front port update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/front-ports/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalFrontPort(raw)
}

// DeleteFrontPort deletes a front port via DELETE /api/dcim/front-ports/<id>/.
func (c *Client) DeleteFrontPort(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/front-ports/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateRearPort creates a new rear port via POST /api/dcim/rear-ports/.
func (c *Client) CreateRearPort(ctx context.Context, token string, in domain.RearPortWrite) (*domain.RearPort, error) {
	payload, err := json.Marshal(rearPortWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal rear port create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/rear-ports/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalRearPort(raw)
}

// UpdateRearPort partially updates a rear port via PATCH
// /api/dcim/rear-ports/<id>/.
func (c *Client) UpdateRearPort(ctx context.Context, token string, objectID int, in domain.RearPortWrite) (
	*domain.RearPort, error,
) {
	payload, err := json.Marshal(rearPortWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal rear port update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/rear-ports/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalRearPort(raw)
}

// DeleteRearPort deletes a rear port via DELETE /api/dcim/rear-ports/<id>/.
func (c *Client) DeleteRearPort(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/rear-ports/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateConsolePort creates a new console port via POST /api/dcim/console-ports/.
func (c *Client) CreateConsolePort(ctx context.Context, token string, in domain.ConsolePortWrite) (
	*domain.ConsolePort, error,
) {
	payload, err := json.Marshal(consolePortWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal console port create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/console-ports/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalConsolePort(raw)
}

// UpdateConsolePort partially updates a console port via PATCH
// /api/dcim/console-ports/<id>/.
func (c *Client) UpdateConsolePort(ctx context.Context, token string, objectID int, in domain.ConsolePortWrite) (
	*domain.ConsolePort, error,
) {
	payload, err := json.Marshal(consolePortWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal console port update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/console-ports/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalConsolePort(raw)
}

// DeleteConsolePort deletes a console port via DELETE /api/dcim/console-ports/<id>/.
func (c *Client) DeleteConsolePort(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/console-ports/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateConsoleServerPort creates a new console server port via POST
// /api/dcim/console-server-ports/.
func (c *Client) CreateConsoleServerPort(ctx context.Context, token string, in domain.ConsoleServerPortWrite) (
	*domain.ConsoleServerPort, error,
) {
	payload, err := json.Marshal(consoleServerPortWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal console server port create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/console-server-ports/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalConsoleServerPort(raw)
}

// UpdateConsoleServerPort partially updates a console server port via PATCH
// /api/dcim/console-server-ports/<id>/.
func (c *Client) UpdateConsoleServerPort(
	ctx context.Context, token string, objectID int, in domain.ConsoleServerPortWrite,
) (
	*domain.ConsoleServerPort, error,
) {
	payload, err := json.Marshal(consoleServerPortWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal console server port update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/console-server-ports/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalConsoleServerPort(raw)
}

// DeleteConsoleServerPort deletes a console server port via DELETE
// /api/dcim/console-server-ports/<id>/.
func (c *Client) DeleteConsoleServerPort(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/console-server-ports/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreatePowerPort creates a new power port via POST /api/dcim/power-ports/.
func (c *Client) CreatePowerPort(ctx context.Context, token string, in domain.PowerPortWrite) (
	*domain.PowerPort, error,
) {
	payload, err := json.Marshal(powerPortWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal power port create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/power-ports/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalPowerPort(raw)
}

// UpdatePowerPort partially updates a power port via PATCH
// /api/dcim/power-ports/<id>/.
func (c *Client) UpdatePowerPort(ctx context.Context, token string, objectID int, in domain.PowerPortWrite) (
	*domain.PowerPort, error,
) {
	payload, err := json.Marshal(powerPortWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal power port update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/power-ports/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalPowerPort(raw)
}

// DeletePowerPort deletes a power port via DELETE /api/dcim/power-ports/<id>/.
func (c *Client) DeletePowerPort(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/power-ports/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreatePowerOutlet creates a new power outlet via POST /api/dcim/power-outlets/.
func (c *Client) CreatePowerOutlet(ctx context.Context, token string, in domain.PowerOutletWrite) (
	*domain.PowerOutlet, error,
) {
	payload, err := json.Marshal(powerOutletWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal power outlet create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/power-outlets/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalPowerOutlet(raw)
}

// UpdatePowerOutlet partially updates a power outlet via PATCH
// /api/dcim/power-outlets/<id>/.
func (c *Client) UpdatePowerOutlet(ctx context.Context, token string, objectID int, in domain.PowerOutletWrite) (
	*domain.PowerOutlet, error,
) {
	payload, err := json.Marshal(powerOutletWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal power outlet update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/power-outlets/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalPowerOutlet(raw)
}

// DeletePowerOutlet deletes a power outlet via DELETE /api/dcim/power-outlets/<id>/.
func (c *Client) DeletePowerOutlet(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/power-outlets/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreatePowerFeed creates a new power feed via POST /api/dcim/power-feeds/.
func (c *Client) CreatePowerFeed(ctx context.Context, token string, in domain.PowerFeedWrite) (
	*domain.PowerFeed, error,
) {
	payload, err := json.Marshal(powerFeedWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal power feed create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/power-feeds/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalPowerFeed(raw)
}

// UpdatePowerFeed partially updates a power feed via PATCH
// /api/dcim/power-feeds/<id>/.
func (c *Client) UpdatePowerFeed(ctx context.Context, token string, objectID int, in domain.PowerFeedWrite) (
	*domain.PowerFeed, error,
) {
	payload, err := json.Marshal(powerFeedWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal power feed update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/power-feeds/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalPowerFeed(raw)
}

// DeletePowerFeed deletes a power feed via DELETE /api/dcim/power-feeds/<id>/.
func (c *Client) DeletePowerFeed(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/power-feeds/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateDeviceBay creates a new device bay via POST /api/dcim/device-bays/.
func (c *Client) CreateDeviceBay(ctx context.Context, token string, in domain.DeviceBayWrite) (
	*domain.DeviceBay, error,
) {
	payload, err := json.Marshal(deviceBayWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal device bay create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/device-bays/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalDeviceBay(raw)
}

// UpdateDeviceBay partially updates a device bay via PATCH
// /api/dcim/device-bays/<id>/.
func (c *Client) UpdateDeviceBay(ctx context.Context, token string, objectID int, in domain.DeviceBayWrite) (
	*domain.DeviceBay, error,
) {
	payload, err := json.Marshal(deviceBayWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal device bay update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/device-bays/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalDeviceBay(raw)
}

// DeleteDeviceBay deletes a device bay via DELETE /api/dcim/device-bays/<id>/.
func (c *Client) DeleteDeviceBay(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/device-bays/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateRackReservation creates a new rack reservation via POST
// /api/dcim/rack-reservations/.
func (c *Client) CreateRackReservation(ctx context.Context, token string, in domain.RackReservationWrite) (
	*domain.RackReservation, error,
) {
	payload, err := json.Marshal(rackReservationWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal rack reservation create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/rack-reservations/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalRackReservation(raw)
}

// UpdateRackReservation partially updates a rack reservation via PATCH
// /api/dcim/rack-reservations/<id>/.
func (c *Client) UpdateRackReservation(
	ctx context.Context, token string, objectID int, in domain.RackReservationWrite,
) (
	*domain.RackReservation, error,
) {
	payload, err := json.Marshal(rackReservationWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal rack reservation update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/rack-reservations/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalRackReservation(raw)
}

// DeleteRackReservation deletes a rack reservation via DELETE
// /api/dcim/rack-reservations/<id>/.
func (c *Client) DeleteRackReservation(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/rack-reservations/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateInventoryItem creates a new inventory item via POST
// /api/dcim/inventory-items/.
func (c *Client) CreateInventoryItem(ctx context.Context, token string, in domain.InventoryItemWrite) (
	*domain.InventoryItem, error,
) {
	payload, err := json.Marshal(inventoryItemWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal inventory item create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/inventory-items/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalInventoryItem(raw)
}

// UpdateInventoryItem partially updates an inventory item via PATCH
// /api/dcim/inventory-items/<id>/.
func (c *Client) UpdateInventoryItem(ctx context.Context, token string, objectID int, in domain.InventoryItemWrite) (
	*domain.InventoryItem, error,
) {
	payload, err := json.Marshal(inventoryItemWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal inventory item update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/inventory-items/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalInventoryItem(raw)
}

// DeleteInventoryItem deletes an inventory item via DELETE
// /api/dcim/inventory-items/<id>/.
func (c *Client) DeleteInventoryItem(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/inventory-items/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateMACAddress creates a new MAC address via POST /api/dcim/mac-addresses/.
func (c *Client) CreateMACAddress(ctx context.Context, token string, in domain.MACAddressWrite) (
	*domain.MACAddress, error,
) {
	payload, err := json.Marshal(macAddressWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal MAC address create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/dcim/mac-addresses/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalMACAddress(raw)
}

// UpdateMACAddress partially updates a MAC address via PATCH
// /api/dcim/mac-addresses/<id>/.
func (c *Client) UpdateMACAddress(ctx context.Context, token string, objectID int, in domain.MACAddressWrite) (
	*domain.MACAddress, error,
) {
	payload, err := json.Marshal(macAddressWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal MAC address update request: %w", err)
	}
	path := fmt.Sprintf("/api/dcim/mac-addresses/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalMACAddress(raw)
}

// DeleteMACAddress deletes a MAC address via DELETE /api/dcim/mac-addresses/<id>/.
func (c *Client) DeleteMACAddress(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/dcim/mac-addresses/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateRIR creates a new RIR via POST /api/ipam/rirs/.
func (c *Client) CreateRIR(ctx context.Context, token string, in domain.RIRWrite) (*domain.RIR, error) {
	payload, err := json.Marshal(rirWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal RIR create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/ipam/rirs/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalRIR(raw)
}

// UpdateRIR partially updates an RIR via PATCH /api/ipam/rirs/<id>/.
func (c *Client) UpdateRIR(ctx context.Context, token string, objectID int, in domain.RIRWrite) (*domain.RIR, error) {
	payload, err := json.Marshal(rirWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal RIR update request: %w", err)
	}
	path := fmt.Sprintf("/api/ipam/rirs/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalRIR(raw)
}

// DeleteRIR deletes an RIR via DELETE /api/ipam/rirs/<id>/.
func (c *Client) DeleteRIR(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/ipam/rirs/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateAggregate creates a new aggregate via POST /api/ipam/aggregates/.
func (c *Client) CreateAggregate(
	ctx context.Context, token string, in domain.AggregateWrite,
) (*domain.Aggregate, error) {
	payload, err := json.Marshal(aggregateWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal aggregate create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/ipam/aggregates/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalAggregate(raw)
}

// UpdateAggregate partially updates an aggregate via PATCH /api/ipam/aggregates/<id>/.
func (c *Client) UpdateAggregate(
	ctx context.Context, token string, objectID int, in domain.AggregateWrite,
) (*domain.Aggregate, error) {
	payload, err := json.Marshal(aggregateWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal aggregate update request: %w", err)
	}
	path := fmt.Sprintf("/api/ipam/aggregates/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalAggregate(raw)
}

// DeleteAggregate deletes an aggregate via DELETE /api/ipam/aggregates/<id>/.
func (c *Client) DeleteAggregate(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/ipam/aggregates/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateIPRange creates a new IP range via POST /api/ipam/ip-ranges/.
func (c *Client) CreateIPRange(ctx context.Context, token string, in domain.IPRangeWrite) (*domain.IPRange, error) {
	payload, err := json.Marshal(ipRangeWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal IP range create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/ipam/ip-ranges/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalIPRange(raw)
}

// UpdateIPRange partially updates an IP range via PATCH /api/ipam/ip-ranges/<id>/.
func (c *Client) UpdateIPRange(
	ctx context.Context, token string, objectID int, in domain.IPRangeWrite,
) (*domain.IPRange, error) {
	payload, err := json.Marshal(ipRangeWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal IP range update request: %w", err)
	}
	path := fmt.Sprintf("/api/ipam/ip-ranges/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalIPRange(raw)
}

// DeleteIPRange deletes an IP range via DELETE /api/ipam/ip-ranges/<id>/.
func (c *Client) DeleteIPRange(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/ipam/ip-ranges/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateASN creates a new ASN via POST /api/ipam/asns/.
func (c *Client) CreateASN(ctx context.Context, token string, in domain.ASNWrite) (*domain.ASN, error) {
	payload, err := json.Marshal(asnWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal ASN create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/ipam/asns/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalASN(raw)
}

// UpdateASN partially updates an ASN via PATCH /api/ipam/asns/<id>/.
func (c *Client) UpdateASN(ctx context.Context, token string, objectID int, in domain.ASNWrite) (*domain.ASN, error) {
	payload, err := json.Marshal(asnWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal ASN update request: %w", err)
	}
	path := fmt.Sprintf("/api/ipam/asns/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalASN(raw)
}

// DeleteASN deletes an ASN via DELETE /api/ipam/asns/<id>/.
func (c *Client) DeleteASN(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/ipam/asns/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateASNRange creates a new ASN range via POST /api/ipam/asn-ranges/.
func (c *Client) CreateASNRange(ctx context.Context, token string, in domain.ASNRangeWrite) (*domain.ASNRange, error) {
	payload, err := json.Marshal(asnRangeWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal ASN range create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/ipam/asn-ranges/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalASNRange(raw)
}

// UpdateASNRange partially updates an ASN range via PATCH /api/ipam/asn-ranges/<id>/.
func (c *Client) UpdateASNRange(
	ctx context.Context, token string, objectID int, in domain.ASNRangeWrite,
) (*domain.ASNRange, error) {
	payload, err := json.Marshal(asnRangeWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal ASN range update request: %w", err)
	}
	path := fmt.Sprintf("/api/ipam/asn-ranges/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalASNRange(raw)
}

// DeleteASNRange deletes an ASN range via DELETE /api/ipam/asn-ranges/<id>/.
func (c *Client) DeleteASNRange(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/ipam/asn-ranges/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateRouteTarget creates a new route target via POST /api/ipam/route-targets/.
func (c *Client) CreateRouteTarget(
	ctx context.Context, token string, in domain.RouteTargetWrite,
) (*domain.RouteTarget, error) {
	payload, err := json.Marshal(routeTargetWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal route target create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/ipam/route-targets/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalRouteTarget(raw)
}

// UpdateRouteTarget partially updates a route target via PATCH /api/ipam/route-targets/<id>/.
func (c *Client) UpdateRouteTarget(
	ctx context.Context, token string, objectID int, in domain.RouteTargetWrite,
) (*domain.RouteTarget, error) {
	payload, err := json.Marshal(routeTargetWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal route target update request: %w", err)
	}
	path := fmt.Sprintf("/api/ipam/route-targets/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalRouteTarget(raw)
}

// DeleteRouteTarget deletes a route target via DELETE /api/ipam/route-targets/<id>/.
func (c *Client) DeleteRouteTarget(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/ipam/route-targets/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateService creates a new service via POST /api/ipam/services/.
func (c *Client) CreateService(ctx context.Context, token string, in domain.ServiceWrite) (*domain.Service, error) {
	payload, err := json.Marshal(serviceWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal service create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/ipam/services/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalService(raw)
}

// UpdateService partially updates a service via PATCH /api/ipam/services/<id>/.
func (c *Client) UpdateService(
	ctx context.Context, token string, objectID int, in domain.ServiceWrite,
) (*domain.Service, error) {
	payload, err := json.Marshal(serviceWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal service update request: %w", err)
	}
	path := fmt.Sprintf("/api/ipam/services/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalService(raw)
}

// DeleteService deletes a service via DELETE /api/ipam/services/<id>/.
func (c *Client) DeleteService(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/ipam/services/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateServiceTemplate creates a new service template via POST
// /api/ipam/service-templates/.
func (c *Client) CreateServiceTemplate(
	ctx context.Context, token string, in domain.ServiceTemplateWrite,
) (*domain.ServiceTemplate, error) {
	payload, err := json.Marshal(serviceTemplateWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal service template create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/ipam/service-templates/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalServiceTemplate(raw)
}

// UpdateServiceTemplate partially updates a service template via PATCH
// /api/ipam/service-templates/<id>/.
func (c *Client) UpdateServiceTemplate(
	ctx context.Context, token string, objectID int, in domain.ServiceTemplateWrite,
) (*domain.ServiceTemplate, error) {
	payload, err := json.Marshal(serviceTemplateWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal service template update request: %w", err)
	}
	path := fmt.Sprintf("/api/ipam/service-templates/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalServiceTemplate(raw)
}

// DeleteServiceTemplate deletes a service template via DELETE
// /api/ipam/service-templates/<id>/.
func (c *Client) DeleteServiceTemplate(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/ipam/service-templates/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateFHRPGroup creates a new FHRP group via POST /api/ipam/fhrp-groups/.
func (c *Client) CreateFHRPGroup(
	ctx context.Context, token string, in domain.FHRPGroupWrite,
) (*domain.FHRPGroup, error) {
	payload, err := json.Marshal(fhrpGroupWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal FHRP group create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/ipam/fhrp-groups/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalFHRPGroup(raw)
}

// UpdateFHRPGroup partially updates an FHRP group via PATCH /api/ipam/fhrp-groups/<id>/.
func (c *Client) UpdateFHRPGroup(
	ctx context.Context, token string, objectID int, in domain.FHRPGroupWrite,
) (*domain.FHRPGroup, error) {
	payload, err := json.Marshal(fhrpGroupWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal FHRP group update request: %w", err)
	}
	path := fmt.Sprintf("/api/ipam/fhrp-groups/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalFHRPGroup(raw)
}

// DeleteFHRPGroup deletes an FHRP group via DELETE /api/ipam/fhrp-groups/<id>/.
func (c *Client) DeleteFHRPGroup(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/ipam/fhrp-groups/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateFHRPGroupAssignment creates a new FHRP group assignment via POST
// /api/ipam/fhrp-group-assignments/.
func (c *Client) CreateFHRPGroupAssignment(
	ctx context.Context, token string, in domain.FHRPGroupAssignmentWrite,
) (*domain.FHRPGroupAssignment, error) {
	payload, err := json.Marshal(fhrpGroupAssignmentWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal FHRP group assignment create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/ipam/fhrp-group-assignments/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalFHRPGroupAssignment(raw)
}

// UpdateFHRPGroupAssignment partially updates an FHRP group assignment via
// PATCH /api/ipam/fhrp-group-assignments/<id>/.
func (c *Client) UpdateFHRPGroupAssignment(
	ctx context.Context, token string, objectID int, in domain.FHRPGroupAssignmentWrite,
) (*domain.FHRPGroupAssignment, error) {
	payload, err := json.Marshal(fhrpGroupAssignmentWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal FHRP group assignment update request: %w", err)
	}
	path := fmt.Sprintf("/api/ipam/fhrp-group-assignments/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalFHRPGroupAssignment(raw)
}

// DeleteFHRPGroupAssignment deletes an FHRP group assignment via DELETE
// /api/ipam/fhrp-group-assignments/<id>/.
func (c *Client) DeleteFHRPGroupAssignment(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/ipam/fhrp-group-assignments/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateVLANTranslationPolicy creates a new VLAN translation policy via POST
// /api/ipam/vlan-translation-policies/.
func (c *Client) CreateVLANTranslationPolicy(
	ctx context.Context, token string, in domain.VLANTranslationPolicyWrite,
) (*domain.VLANTranslationPolicy, error) {
	payload, err := json.Marshal(vlanTranslationPolicyWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal VLAN translation policy create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/ipam/vlan-translation-policies/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalVLANTranslationPolicy(raw)
}

// UpdateVLANTranslationPolicy partially updates a VLAN translation policy via
// PATCH /api/ipam/vlan-translation-policies/<id>/.
func (c *Client) UpdateVLANTranslationPolicy(
	ctx context.Context, token string, objectID int, in domain.VLANTranslationPolicyWrite,
) (*domain.VLANTranslationPolicy, error) {
	payload, err := json.Marshal(vlanTranslationPolicyWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal VLAN translation policy update request: %w", err)
	}
	path := fmt.Sprintf("/api/ipam/vlan-translation-policies/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalVLANTranslationPolicy(raw)
}

// DeleteVLANTranslationPolicy deletes a VLAN translation policy via DELETE
// /api/ipam/vlan-translation-policies/<id>/.
func (c *Client) DeleteVLANTranslationPolicy(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/ipam/vlan-translation-policies/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateVLANTranslationRule creates a new VLAN translation rule via POST
// /api/ipam/vlan-translation-rules/.
func (c *Client) CreateVLANTranslationRule(
	ctx context.Context, token string, in domain.VLANTranslationRuleWrite,
) (*domain.VLANTranslationRule, error) {
	payload, err := json.Marshal(vlanTranslationRuleWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal VLAN translation rule create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/ipam/vlan-translation-rules/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalVLANTranslationRule(raw)
}

// UpdateVLANTranslationRule partially updates a VLAN translation rule via PATCH
// /api/ipam/vlan-translation-rules/<id>/.
func (c *Client) UpdateVLANTranslationRule(
	ctx context.Context, token string, objectID int, in domain.VLANTranslationRuleWrite,
) (*domain.VLANTranslationRule, error) {
	payload, err := json.Marshal(vlanTranslationRuleWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal VLAN translation rule update request: %w", err)
	}
	path := fmt.Sprintf("/api/ipam/vlan-translation-rules/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalVLANTranslationRule(raw)
}

// DeleteVLANTranslationRule deletes a VLAN translation rule via DELETE
// /api/ipam/vlan-translation-rules/<id>/.
func (c *Client) DeleteVLANTranslationRule(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/ipam/vlan-translation-rules/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateContactGroup creates a new contact group via POST
// /api/tenancy/contact-groups/.
func (c *Client) CreateContactGroup(
	ctx context.Context, token string, in domain.ContactGroupWrite,
) (*domain.ContactGroup, error) {
	payload, err := json.Marshal(contactGroupWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal contact group create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/tenancy/contact-groups/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalContactGroup(raw)
}

// UpdateContactGroup partially updates a contact group via PATCH
// /api/tenancy/contact-groups/<id>/.
func (c *Client) UpdateContactGroup(
	ctx context.Context, token string, objectID int, in domain.ContactGroupWrite,
) (*domain.ContactGroup, error) {
	payload, err := json.Marshal(contactGroupWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal contact group update request: %w", err)
	}
	path := fmt.Sprintf("/api/tenancy/contact-groups/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalContactGroup(raw)
}

// DeleteContactGroup deletes a contact group via DELETE
// /api/tenancy/contact-groups/<id>/.
func (c *Client) DeleteContactGroup(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/tenancy/contact-groups/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateContactRole creates a new contact role via POST
// /api/tenancy/contact-roles/.
func (c *Client) CreateContactRole(
	ctx context.Context, token string, in domain.ContactRoleWrite,
) (*domain.ContactRole, error) {
	payload, err := json.Marshal(contactRoleWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal contact role create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/tenancy/contact-roles/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalContactRole(raw)
}

// UpdateContactRole partially updates a contact role via PATCH
// /api/tenancy/contact-roles/<id>/.
func (c *Client) UpdateContactRole(
	ctx context.Context, token string, objectID int, in domain.ContactRoleWrite,
) (*domain.ContactRole, error) {
	payload, err := json.Marshal(contactRoleWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal contact role update request: %w", err)
	}
	path := fmt.Sprintf("/api/tenancy/contact-roles/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalContactRole(raw)
}

// DeleteContactRole deletes a contact role via DELETE
// /api/tenancy/contact-roles/<id>/.
func (c *Client) DeleteContactRole(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/tenancy/contact-roles/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateVirtualDisk creates a new virtual disk via POST
// /api/virtualization/virtual-disks/.
func (c *Client) CreateVirtualDisk(
	ctx context.Context, token string, in domain.VirtualDiskWrite,
) (*domain.VirtualDisk, error) {
	payload, err := json.Marshal(virtualDiskWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal virtual disk create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/virtualization/virtual-disks/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalVirtualDisk(raw)
}

// UpdateVirtualDisk partially updates a virtual disk via PATCH
// /api/virtualization/virtual-disks/<id>/.
func (c *Client) UpdateVirtualDisk(
	ctx context.Context, token string, objectID int, in domain.VirtualDiskWrite,
) (*domain.VirtualDisk, error) {
	payload, err := json.Marshal(virtualDiskWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal virtual disk update request: %w", err)
	}
	path := fmt.Sprintf("/api/virtualization/virtual-disks/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalVirtualDisk(raw)
}

// DeleteVirtualDisk deletes a virtual disk via DELETE
// /api/virtualization/virtual-disks/<id>/.
func (c *Client) DeleteVirtualDisk(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/virtualization/virtual-disks/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// CreateVirtualMachineType creates a new virtual machine type via POST
// /api/virtualization/virtual-machine-types/.
func (c *Client) CreateVirtualMachineType(
	ctx context.Context, token string, in domain.VirtualMachineTypeWrite,
) (*domain.VirtualMachineType, error) {
	payload, err := json.Marshal(virtualMachineTypeWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal virtual machine type create request: %w", err)
	}
	body, err := c.post(ctx, token, "/api/virtualization/virtual-machine-types/", payload)
	raw, mErr := writeResult(http.StatusCreated, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalVirtualMachineType(raw)
}

// UpdateVirtualMachineType partially updates a virtual machine type via PATCH
// /api/virtualization/virtual-machine-types/<id>/.
func (c *Client) UpdateVirtualMachineType(
	ctx context.Context, token string, objectID int, in domain.VirtualMachineTypeWrite,
) (*domain.VirtualMachineType, error) {
	payload, err := json.Marshal(virtualMachineTypeWriteToWire(in))
	if err != nil {
		return nil, fmt.Errorf("marshal virtual machine type update request: %w", err)
	}
	path := fmt.Sprintf("/api/virtualization/virtual-machine-types/%d/", objectID)
	body, err := c.patch(ctx, token, path, payload)
	raw, mErr := writeResult(http.StatusOK, body, err)
	if mErr != nil {
		return nil, mErr
	}

	return unmarshalVirtualMachineType(raw)
}

// DeleteVirtualMachineType deletes a virtual machine type via DELETE
// /api/virtualization/virtual-machine-types/<id>/.
func (c *Client) DeleteVirtualMachineType(ctx context.Context, token string, objectID int) error {
	path := fmt.Sprintf("/api/virtualization/virtual-machine-types/%d/", objectID)
	body, err := c.delete(ctx, token, path)
	_, mErr := writeResult(http.StatusNoContent, body, err)

	return mErr
}

// — conversion functions —

// siteWriteToWire converts a domain write DTO to the wire request model. The
// fields map one-to-one; nil pointers on the domain side remain nil on the wire
// side so json.Marshal omits them (required for a partial PATCH body).
func siteWriteToWire(input domain.SiteWrite) WireSiteWrite {
	return WireSiteWrite{
		Name:            input.Name,
		Slug:            input.Slug,
		Status:          input.Status,
		Region:          input.Region,
		Tenant:          input.Tenant,
		Facility:        input.Facility,
		TimeZone:        input.TimeZone,
		Description:     input.Description,
		PhysicalAddress: input.PhysicalAddress,
		ShippingAddress: input.ShippingAddress,
		Comments:        input.Comments,
		Tags:            input.Tags,
		CustomFields:    input.CustomFields,
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

// regionWriteToWire converts a domain region write DTO to the wire request
// model. The fields map one-to-one; nil pointers on the domain side remain nil
// on the wire side so json.Marshal omits them (required for a partial PATCH
// body).
func regionWriteToWire(input domain.RegionWrite) WireRegionWrite {
	return WireRegionWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Parent:       input.Parent,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalRegion decodes a region response body into a domain.Region.
func unmarshalRegion(raw domain.RawObject) (*domain.Region, error) {
	var w WireRegion
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal region: %w", err)
	}
	d := wireRegionToDomain(w)

	return &d, nil
}

// tenantGroupWriteToWire converts a domain tenant group write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func tenantGroupWriteToWire(input domain.TenantGroupWrite) WireTenantGroupWrite {
	return WireTenantGroupWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Parent:       input.Parent,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalTenantGroup decodes a tenant group response body into a
// domain.TenantGroup.
func unmarshalTenantGroup(raw domain.RawObject) (*domain.TenantGroup, error) {
	var w WireTenantGroup
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tenant group: %w", err)
	}
	d := wireTenantGroupToDomain(w)

	return &d, nil
}

// deviceWriteToWire converts a domain device write DTO to the wire request
// model. The fields map one-to-one; nil pointers on the domain side remain nil
// on the wire side so json.Marshal omits them (required for a partial PATCH
// body).
func deviceWriteToWire(input domain.DeviceWrite) WireDeviceWrite {
	return WireDeviceWrite{
		Name:         input.Name,
		DeviceType:   input.DeviceType,
		Role:         input.Role,
		Tenant:       input.Tenant,
		Platform:     input.Platform,
		Serial:       input.Serial,
		AssetTag:     input.AssetTag,
		Site:         input.Site,
		Rack:         input.Rack,
		Position:     input.Position,
		Face:         input.Face,
		Status:       input.Status,
		Cluster:      input.Cluster,
		Comments:     input.Comments,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
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
func ipAddressWriteToWire(input domain.IPAddressWrite) WireIPAddressWrite {
	return WireIPAddressWrite{
		Address:            input.Address,
		Status:             input.Status,
		Role:               input.Role,
		VRF:                input.VRF,
		Tenant:             input.Tenant,
		DNSName:            input.DNSName,
		Description:        input.Description,
		AssignedObjectType: input.AssignedObjectType,
		AssignedObjectID:   input.AssignedObjectID,
		Comments:           input.Comments,
		Tags:               input.Tags,
		CustomFields:       input.CustomFields,
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
func prefixWriteToWire(input domain.PrefixWrite) WirePrefixWrite {
	return WirePrefixWrite{
		Prefix:       input.Prefix,
		Site:         input.Site,
		VRF:          input.VRF,
		Tenant:       input.Tenant,
		VLAN:         input.VLAN,
		Status:       input.Status,
		Role:         input.Role,
		IsPool:       input.IsPool,
		Description:  input.Description,
		Comments:     input.Comments,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
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
func vlanWriteToWire(input domain.VLANWrite) WireVLANWrite {
	return WireVLANWrite{
		VID:          input.VID,
		Name:         input.Name,
		Site:         input.Site,
		Group:        input.Group,
		Tenant:       input.Tenant,
		Status:       input.Status,
		Role:         input.Role,
		Description:  input.Description,
		Comments:     input.Comments,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
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
func virtualMachineWriteToWire(input domain.VirtualMachineWrite) WireVirtualMachineWrite {
	return WireVirtualMachineWrite{
		Name:         input.Name,
		Cluster:      input.Cluster,
		Role:         input.Role,
		Tenant:       input.Tenant,
		Platform:     input.Platform,
		Status:       input.Status,
		Site:         input.Site,
		VCPUs:        input.VCPUs,
		Memory:       input.Memory,
		Disk:         input.Disk,
		Comments:     input.Comments,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
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
func clusterWriteToWire(input domain.ClusterWrite) WireClusterWrite {
	return WireClusterWrite{
		Name:         input.Name,
		ClusterType:  input.ClusterType,
		ClusterGroup: input.ClusterGroup,
		Site:         input.Site,
		Tenant:       input.Tenant,
		Description:  input.Description,
		Comments:     input.Comments,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
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
func circuitWriteToWire(input domain.CircuitWrite) WireCircuitWrite {
	return WireCircuitWrite{
		CID:          input.CID,
		Provider:     input.Provider,
		CircuitType:  input.CircuitType,
		Tenant:       input.Tenant,
		Status:       input.Status,
		Description:  input.Description,
		Comments:     input.Comments,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
		InstallDate:  input.InstallDate,
		CommitRate:   input.CommitRate,
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

// rackWriteToWire converts a domain rack write DTO to the wire request model.
// The fields map one-to-one; nil pointers on the domain side remain nil on the
// wire side so json.Marshal omits them (required for a partial PATCH body).
func rackWriteToWire(input domain.RackWrite) WireRackWrite {
	return WireRackWrite{
		Name:         input.Name,
		FacilityID:   input.FacilityID,
		Site:         input.Site,
		Location:     input.Location,
		Tenant:       input.Tenant,
		Status:       input.Status,
		Role:         input.Role,
		Serial:       input.Serial,
		AssetTag:     input.AssetTag,
		Type:         input.Type,
		Width:        input.Width,
		UHeight:      input.UHeight,
		Comments:     input.Comments,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalRack decodes a rack response body into a domain.Rack.
func unmarshalRack(raw domain.RawObject) (*domain.Rack, error) {
	var w WireRack
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rack: %w", err)
	}
	d := wireRackToDomain(w)

	return &d, nil
}

// interfaceWriteToWire converts a domain interface write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func interfaceWriteToWire(input domain.InterfaceWrite) WireInterfaceWrite {
	return WireInterfaceWrite{
		Device:       input.Device,
		Name:         input.Name,
		Type:         input.Type,
		Enabled:      input.Enabled,
		MTU:          input.MTU,
		MACAddress:   input.MACAddress,
		Speed:        input.Speed,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalInterface decodes an interface response body into a domain.Interface.
func unmarshalInterface(raw domain.RawObject) (*domain.Interface, error) {
	var w WireInterface
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal interface: %w", err)
	}
	d := wireInterfaceToDomain(w)

	return &d, nil
}

// circuitTerminationWriteToWire converts a domain circuit termination write DTO
// to the wire request model. The fields map one-to-one; nil pointers on the
// domain side remain nil on the wire side so json.Marshal omits them (required
// for a partial PATCH body).
func circuitTerminationWriteToWire(input domain.CircuitTerminationWrite) WireCircuitTerminationWrite {
	// NetBox 4.6 models a circuit termination's terminating object generically
	// (termination_type/termination_id). The domain's Site maps to a dcim.site
	// termination.
	terminationType := ""
	if input.Site != nil {
		terminationType = "dcim.site"
	}

	return WireCircuitTerminationWrite{
		Circuit:         input.Circuit,
		TermSide:        input.TermSide,
		TerminationType: terminationType,
		TerminationID:   input.Site,
		Speed:           input.Speed,
		UpstreamSpeed:   input.UpstreamSpeed,
		Description:     input.Description,
		Tags:            input.Tags,
		CustomFields:    input.CustomFields,
	}
}

// unmarshalCircuitTermination decodes a circuit termination response body into a
// domain.CircuitTermination.
func unmarshalCircuitTermination(raw domain.RawObject) (*domain.CircuitTermination, error) {
	var w WireCircuitTermination
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal circuit termination: %w", err)
	}
	d := wireCircuitTerminationToDomain(w)

	return &d, nil
}

// cableTerminationWriteToWire converts a domain cable termination DTO to the
// wire request model.
func cableTerminationWriteToWire(input *domain.CableTerminationWrite) *WireCableWriteTermination {
	if input == nil {
		return nil
	}

	return &WireCableWriteTermination{
		ObjectType: input.ObjectType,
		ObjectID:   input.ObjectID,
	}
}

// cableWriteToWire converts a domain cable write DTO to the wire request model.
// The fields map one-to-one; nil pointers on the domain side remain nil on the
// wire side so json.Marshal omits them (required for a partial PATCH body).
func cableWriteToWire(input domain.CableWrite) WireCableWrite {
	return WireCableWrite{
		TerminationA: cableTerminationWriteToWire(input.TerminationA),
		TerminationB: cableTerminationWriteToWire(input.TerminationB),
		Type:         input.Type,
		Status:       input.Status,
		Label:        input.Label,
		Color:        input.Color,
		Length:       input.Length,
		LengthUnit:   input.LengthUnit,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalCable decodes a cable response body into a domain.Cable.
func unmarshalCable(raw domain.RawObject) (*domain.Cable, error) {
	var w WireCable
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cable: %w", err)
	}
	d := wireCableToDomain(w)

	return &d, nil
}

// vmInterfaceWriteToWire converts a domain VM interface write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func vmInterfaceWriteToWire(input domain.VMInterfaceWrite) WireVMInterfaceWrite {
	return WireVMInterfaceWrite{
		VirtualMachine: input.VirtualMachine,
		Name:           input.Name,
		Enabled:        input.Enabled,
		MTU:            input.MTU,
		MACAddress:     input.MACAddress,
		Description:    input.Description,
		Tags:           input.Tags,
		CustomFields:   input.CustomFields,
	}
}

// unmarshalVMInterface decodes a VM interface response body into a
// domain.VMInterface.
func unmarshalVMInterface(raw domain.RawObject) (*domain.VMInterface, error) {
	var w WireVMInterface
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal vm interface: %w", err)
	}
	d := wireVMInterfaceToDomain(w)

	return &d, nil
}

// providerWriteToWire converts a domain provider write DTO to the wire request
// model. The fields map one-to-one; nil pointers on the domain side remain nil
// on the wire side so json.Marshal omits them (required for a partial PATCH
// body).
func providerWriteToWire(input domain.ProviderWrite) WireProviderWrite {
	return WireProviderWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Asn:          input.Asn,
		Account:      input.Account,
		PortalURL:    input.PortalURL,
		NocContact:   input.NocContact,
		AdminContact: input.AdminContact,
		Comments:     input.Comments,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalProvider decodes a provider response body into a domain.Provider.
func unmarshalProvider(raw domain.RawObject) (*domain.Provider, error) {
	var w WireProvider
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal provider: %w", err)
	}
	d := wireProviderToDomain(w)

	return &d, nil
}

// wireProviderToDomain converts a wire provider to the domain model.
func wireProviderToDomain(wire WireProvider) domain.Provider {
	return domain.Provider{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Slug:         wire.Slug,
		Display:      wire.Display,
		Asn:          wire.Asn,
		Account:      wire.Account,
		PortalURL:    wire.PortalURL,
		NocContact:   wire.NocContact,
		AdminContact: wire.AdminContact,
		Comments:     wire.Comments,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// tenantWriteToWire converts a domain tenant write DTO to the wire request
// model. The fields map one-to-one; nil pointers on the domain side remain nil
// on the wire side so json.Marshal omits them (required for a partial PATCH
// body).
func tenantWriteToWire(input domain.TenantWrite) WireTenantWrite {
	return WireTenantWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Group:        input.Group,
		Description:  input.Description,
		Comments:     input.Comments,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalTenant decodes a tenant response body into a domain.Tenant.
func unmarshalTenant(raw domain.RawObject) (*domain.Tenant, error) {
	var w WireTenant
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tenant: %w", err)
	}
	d := wireTenantToDomain(w)

	return &d, nil
}

// wireTenantToDomain converts a wire tenant to the domain model.
func wireTenantToDomain(wire WireTenant) domain.Tenant {
	return domain.Tenant{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Slug:         wire.Slug,
		Display:      wire.Display,
		Group:        wireNestedToDomain(wire.Group),
		Description:  wire.Description,
		Comments:     wire.Comments,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// manufacturerWriteToWire converts a domain manufacturer write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func manufacturerWriteToWire(input domain.ManufacturerWrite) WireManufacturerWrite {
	return WireManufacturerWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalManufacturer decodes a manufacturer response body into a
// domain.Manufacturer.
func unmarshalManufacturer(raw domain.RawObject) (*domain.Manufacturer, error) {
	var w WireManufacturer
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal manufacturer: %w", err)
	}
	d := wireManufacturerToDomain(w)

	return &d, nil
}

// wireManufacturerToDomain converts a wire manufacturer to the domain model.
func wireManufacturerToDomain(wire WireManufacturer) domain.Manufacturer {
	return domain.Manufacturer{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Slug:         wire.Slug,
		Display:      wire.Display,
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// deviceTypeWriteToWire converts a domain device type write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func deviceTypeWriteToWire(input domain.DeviceTypeWrite) WireDeviceTypeWrite {
	return WireDeviceTypeWrite{
		Manufacturer:  input.Manufacturer,
		Model:         input.Model,
		Slug:          input.Slug,
		PartNumber:    input.PartNumber,
		UHeight:       input.UHeight,
		IsFullDepth:   input.IsFullDepth,
		SubdeviceRole: input.SubdeviceRole,
		Comments:      input.Comments,
		Tags:          input.Tags,
		CustomFields:  input.CustomFields,
	}
}

// unmarshalDeviceType decodes a device type response body into a
// domain.DeviceType.
func unmarshalDeviceType(raw domain.RawObject) (*domain.DeviceType, error) {
	var w WireDeviceType
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal device type: %w", err)
	}
	d := wireDeviceTypeToDomain(w)

	return &d, nil
}

// wireDeviceTypeToDomain converts a wire device type to the domain model.
func wireDeviceTypeToDomain(wire WireDeviceType) domain.DeviceType {
	return domain.DeviceType{
		ID:            wire.ID,
		URL:           wire.URL,
		Manufacturer:  wireNestedToDomain(wire.Manufacturer),
		Model:         wire.Model,
		Slug:          wire.Slug,
		Display:       wire.Display,
		PartNumber:    wire.PartNumber,
		UHeight:       wire.UHeight,
		IsFullDepth:   wire.IsFullDepth,
		SubdeviceRole: wireLabelToDomain(wire.SubdeviceRole),
		Comments:      wire.Comments,
		Tags:          wireTagsToDomain(wire.Tags),
		CustomFields:  wire.CustomFields,
		Created:       wire.Created,
		LastUpdated:   wire.LastUpdated,
	}
}

func wireNestedToDomain(wire *WireNested) *domain.Nested {
	if wire == nil {
		return nil
	}

	return &domain.Nested{
		ID:   wire.ID,
		URL:  wire.URL,
		Name: wire.Name,
		Slug: wire.Slug,
	}
}

func wireLabelToDomain(wire *WireLabel) *domain.Label {
	if wire == nil {
		return nil
	}

	return &domain.Label{
		Value: wire.Value,
		Label: wire.Label,
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

func wireSiteToDomain(wire WireSite) domain.Site {
	return domain.Site{
		ID:              wire.ID,
		URL:             wire.URL,
		Name:            wire.Name,
		Slug:            wire.Slug,
		Display:         wire.Display,
		Status:          wireLabelToDomain(wire.Status),
		Region:          wireNestedToDomain(wire.Region),
		Tenant:          wireNestedToDomain(wire.Tenant),
		Facility:        wire.Facility,
		TimeZone:        wire.TimeZone,
		Description:     wire.Description,
		PhysicalAddress: wire.PhysicalAddress,
		ShippingAddress: wire.ShippingAddress,
		Comments:        wire.Comments,
		Tags:            wireTagsToDomain(wire.Tags),
		CustomFields:    wire.CustomFields,
		Created:         wire.Created,
		LastUpdated:     wire.LastUpdated,
	}
}

func wireRegionToDomain(wire WireRegion) domain.Region {
	return domain.Region{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Slug:         wire.Slug,
		Display:      wire.Display,
		Description:  wire.Description,
		Parent:       wireNestedToDomain(wire.Parent),
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// wireTenantGroupToDomain converts a wire tenant group to the domain model.
func wireTenantGroupToDomain(wire WireTenantGroup) domain.TenantGroup {
	return domain.TenantGroup{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Slug:         wire.Slug,
		Display:      wire.Display,
		Description:  wire.Description,
		Parent:       wireNestedToDomain(wire.Parent),
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

func wireDeviceToDomain(wire WireDevice) domain.Device {
	return domain.Device{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Display:      wire.Display,
		DeviceType:   wireNestedToDomain(wire.DeviceType),
		Role:         wireNestedToDomain(wire.Role),
		Tenant:       wireNestedToDomain(wire.Tenant),
		Platform:     wireNestedToDomain(wire.Platform),
		Serial:       wire.Serial,
		AssetTag:     wire.AssetTag,
		Site:         wireNestedToDomain(wire.Site),
		Rack:         wireNestedToDomain(wire.Rack),
		Position:     wire.Position,
		Face:         wireLabelToDomain(wire.Face),
		Status:       wireLabelToDomain(wire.Status),
		Cluster:      wireNestedToDomain(wire.Cluster),
		Comments:     wire.Comments,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

func wireAssignedObjToDomain(wire *WireAssignedObj) *domain.AssignedObject {
	if wire == nil {
		return nil
	}

	return &domain.AssignedObject{
		ID:     wire.ID,
		URL:    wire.URL,
		Device: wireNestedToDomain(wire.Device),
		Name:   wire.Name,
	}
}

func wireIPAddressToDomain(wire WireIPAddress) domain.IPAddress {
	return domain.IPAddress{
		ID:                 wire.ID,
		URL:                wire.URL,
		Address:            wire.Address,
		Display:            wire.Display,
		VRF:                wireNestedToDomain(wire.VRF),
		Tenant:             wireNestedToDomain(wire.Tenant),
		Status:             wireLabelToDomain(wire.Status),
		Role:               wireLabelToDomain(wire.Role),
		AssignedObjectType: wire.AssignedObjectType,
		AssignedObjectID:   wire.AssignedObjectID,
		AssignedObject:     wireAssignedObjToDomain(wire.AssignedObject),
		DNSName:            wire.DNSName,
		Description:        wire.Description,
		Comments:           wire.Comments,
		Tags:               wireTagsToDomain(wire.Tags),
		CustomFields:       wire.CustomFields,
		Created:            wire.Created,
		LastUpdated:        wire.LastUpdated,
	}
}

func wireCircuitTerminationToDomain(wire WireCircuitTermination) domain.CircuitTermination {
	return domain.CircuitTermination{
		ID:            wire.ID,
		URL:           wire.URL,
		Display:       wire.Display,
		Circuit:       wireNestedToDomain(wire.Circuit),
		TermSide:      wire.TermSide,
		Site:          wireNestedToDomain(wire.Site),
		Speed:         wire.Speed,
		UpstreamSpeed: wire.UpstreamSpeed,
		Description:   wire.Description,
		Tags:          wireTagsToDomain(wire.Tags),
		CustomFields:  wire.CustomFields,
		Created:       wire.Created,
		LastUpdated:   wire.LastUpdated,
	}
}

func wireCableToDomain(wire WireCable) domain.Cable {
	return domain.Cable{
		ID:           wire.ID,
		URL:          wire.URL,
		Display:      wire.Display,
		Type:         wireLabelToDomain(wire.Type),
		Status:       wireLabelToDomain(wire.Status),
		Label:        wire.Label,
		Color:        wire.Color,
		Length:       wire.Length,
		LengthUnit:   wireLabelToDomain(wire.LengthUnit),
		TerminationA: wireCableTerminationToDomain(wire.TerminationA),
		TerminationB: wireCableTerminationToDomain(wire.TerminationB),
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

func wireCableTerminationToDomain(wire *WireCableTermination) *domain.CableTermination {
	if wire == nil {
		return nil
	}

	return &domain.CableTermination{
		ID:   wire.ID,
		URL:  wire.URL,
		Name: wire.Name,
		Type: wire.Type,
	}
}

func wirePrefixToDomain(wire WirePrefix) domain.Prefix {
	family := (*domain.Family)(nil)
	if wire.Family != nil {
		family = &domain.Family{
			Value: wire.Family.Value,
			Label: wire.Family.Label,
		}
	}

	return domain.Prefix{
		ID:           wire.ID,
		URL:          wire.URL,
		Prefix:       wire.Prefix,
		Display:      wire.Display,
		Site:         wireNestedToDomain(wire.Site),
		VRF:          wireNestedToDomain(wire.VRF),
		Tenant:       wireNestedToDomain(wire.Tenant),
		VLAN:         wireNestedToDomain(wire.VLAN),
		Status:       wireLabelToDomain(wire.Status),
		Role:         wireNestedToDomain(wire.Role),
		IsPool:       wire.IsPool,
		Description:  wire.Description,
		Comments:     wire.Comments,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
		Children:     wire.Children,
		Depth:        wire.Depth,
		Family:       family,
	}
}

func wireVLANToDomain(wire WireVLAN) domain.VLAN {
	return domain.VLAN{
		ID:           wire.ID,
		URL:          wire.URL,
		Site:         wireNestedToDomain(wire.Site),
		Group:        wireNestedToDomain(wire.Group),
		VID:          wire.VID,
		Name:         wire.Name,
		Display:      wire.Display,
		Tenant:       wireNestedToDomain(wire.Tenant),
		Status:       wireLabelToDomain(wire.Status),
		Role:         wireNestedToDomain(wire.Role),
		Description:  wire.Description,
		Comments:     wire.Comments,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

func wireVirtualMachineToDomain(wire WireVirtualMachine) domain.VirtualMachine {
	return domain.VirtualMachine{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Display:      wire.Display,
		Cluster:      wireNestedToDomain(wire.Cluster),
		Role:         wireNestedToDomain(wire.Role),
		Tenant:       wireNestedToDomain(wire.Tenant),
		Platform:     wireNestedToDomain(wire.Platform),
		Status:       wireLabelToDomain(wire.Status),
		Site:         wireNestedToDomain(wire.Site),
		VCPUs:        wire.VCPUs,
		Memory:       wire.Memory,
		Disk:         wire.Disk,
		Comments:     wire.Comments,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

func wireClusterToDomain(wire WireCluster) domain.Cluster {
	return domain.Cluster{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Display:      wire.Display,
		ClusterType:  wireNestedToDomain(wire.ClusterType),
		ClusterGroup: wireNestedToDomain(wire.ClusterGroup),
		Site:         wireNestedToDomain(wire.Site),
		Tenant:       wireNestedToDomain(wire.Tenant),
		Description:  wire.Description,
		Comments:     wire.Comments,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

func wireCircuitToDomain(wire WireCircuit) domain.Circuit {
	return domain.Circuit{
		ID:           wire.ID,
		URL:          wire.URL,
		CID:          wire.CID,
		Display:      wire.Display,
		Provider:     wireNestedToDomain(wire.Provider),
		CircuitType:  wireNestedToDomain(wire.CircuitType),
		Tenant:       wireNestedToDomain(wire.Tenant),
		Status:       wireLabelToDomain(wire.Status),
		Description:  wire.Description,
		Comments:     wire.Comments,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
		InstallDate:  wire.InstallDate,
		CommitRate:   wire.CommitRate,
	}
}

func wireInterfaceToDomain(wire WireInterface) domain.Interface {
	return domain.Interface{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Display:      wire.Display,
		Device:       wireNestedToDomain(wire.Device),
		Type:         wireLabelToDomain(wire.Type),
		Enabled:      wire.Enabled,
		MTU:          wire.MTU,
		MACAddress:   wire.MACAddress,
		Speed:        wire.Speed,
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

func wireVMInterfaceToDomain(wire WireVMInterface) domain.VMInterface {
	return domain.VMInterface{
		ID:             wire.ID,
		URL:            wire.URL,
		Name:           wire.Name,
		Display:        wire.Display,
		VirtualMachine: wireNestedToDomain(wire.VirtualMachine),
		Enabled:        wire.Enabled,
		MTU:            wire.MTU,
		MACAddress:     wire.MACAddress,
		Description:    wire.Description,
		Tags:           wireTagsToDomain(wire.Tags),
		CustomFields:   wire.CustomFields,
		Created:        wire.Created,
		LastUpdated:    wire.LastUpdated,
	}
}

func wireRackToDomain(wire WireRack) domain.Rack {
	return domain.Rack{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Display:      wire.Display,
		FacilityID:   wire.FacilityID,
		Site:         wireNestedToDomain(wire.Site),
		Location:     wireNestedToDomain(wire.Location),
		Tenant:       wireNestedToDomain(wire.Tenant),
		Status:       wireLabelToDomain(wire.Status),
		Role:         wireNestedToDomain(wire.Role),
		Serial:       wire.Serial,
		AssetTag:     wire.AssetTag,
		Type:         wireLabelToDomain(wire.Type),
		Width:        wire.Width,
		UHeight:      wire.UHeight,
		Comments:     wire.Comments,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

func locationWriteToWire(input domain.LocationWrite) WireLocationWrite {
	return WireLocationWrite{
		Name:         input.Name,
		Site:         input.Site,
		Slug:         input.Slug,
		Parent:       input.Parent,
		Description:  input.Description,
		Status:       input.Status,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalLocation decodes a location response body into a domain.Location.
func unmarshalLocation(raw domain.RawObject) (*domain.Location, error) {
	var w WireLocation
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal location: %w", err)
	}
	d := wireLocationToDomain(w)

	return &d, nil
}

// wireLocationToDomain converts a wire location to the domain model.
func wireLocationToDomain(wire WireLocation) domain.Location {
	return domain.Location{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Slug:         wire.Slug,
		Display:      wire.Display,
		Site:         wireNestedToDomain(wire.Site),
		Parent:       wireNestedToDomain(wire.Parent),
		Status:       wireLabelToDomain(wire.Status),
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

func clusterTypeWriteToWire(input domain.ClusterTypeWrite) WireClusterTypeWrite {
	return WireClusterTypeWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalClusterType decodes a cluster type response body into a
// domain.ClusterType.
func unmarshalClusterType(raw domain.RawObject) (*domain.ClusterType, error) {
	var w WireClusterType
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cluster type: %w", err)
	}
	d := wireClusterTypeToDomain(w)

	return &d, nil
}

// wireClusterTypeToDomain converts a wire cluster type to the domain model.
func wireClusterTypeToDomain(wire WireClusterType) domain.ClusterType {
	return domain.ClusterType{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Slug:         wire.Slug,
		Display:      wire.Display,
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

func clusterGroupWriteToWire(input domain.ClusterGroupWrite) WireClusterGroupWrite {
	return WireClusterGroupWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalClusterGroup decodes a cluster group response body into a
// domain.ClusterGroup.
func unmarshalClusterGroup(raw domain.RawObject) (*domain.ClusterGroup, error) {
	var w WireClusterGroup
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cluster group: %w", err)
	}
	d := wireClusterGroupToDomain(w)

	return &d, nil
}

// wireClusterGroupToDomain converts a wire cluster group to the domain model.
func wireClusterGroupToDomain(wire WireClusterGroup) domain.ClusterGroup {
	return domain.ClusterGroup{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Slug:         wire.Slug,
		Display:      wire.Display,
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

func circuitTypeWriteToWire(input domain.CircuitTypeWrite) WireCircuitTypeWrite {
	return WireCircuitTypeWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalCircuitType decodes a circuit type response body into a
// domain.CircuitType.
func unmarshalCircuitType(raw domain.RawObject) (*domain.CircuitType, error) {
	var w WireCircuitType
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal circuit type: %w", err)
	}
	d := wireCircuitTypeToDomain(w)

	return &d, nil
}

// wireCircuitTypeToDomain converts a wire circuit type to the domain model.
func wireCircuitTypeToDomain(wire WireCircuitType) domain.CircuitType {
	return domain.CircuitType{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Slug:         wire.Slug,
		Display:      wire.Display,
		Description:  wire.Description,
		CircuitCount: wire.CircuitCount,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

func vrfWriteToWire(input domain.VrfWrite) WireVrfWrite {
	return WireVrfWrite{
		Name:         input.Name,
		Rd:           input.Rd,
		Tenant:       input.Tenant,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalVrf decodes a VRF response body into a domain.Vrf.
func unmarshalVrf(raw domain.RawObject) (*domain.Vrf, error) {
	var w WireVrf
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal vrf: %w", err)
	}
	d := wireVrfToDomain(w)

	return &d, nil
}

// wireVrfToDomain converts a wire VRF to the domain model.
func wireVrfToDomain(wire WireVrf) domain.Vrf {
	return domain.Vrf{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Rd:           wire.Rd,
		Tenant:       wireNestedToDomain(wire.Tenant),
		Display:      wire.Display,
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

func vlanGroupWriteToWire(input domain.VlanGroupWrite) WireVlanGroupWrite {
	return WireVlanGroupWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Description:  input.Description,
		ScopeType:    input.ScopeType,
		ScopeID:      input.ScopeID,
		MinVID:       input.MinVID,
		MaxVID:       input.MaxVID,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalVlanGroup decodes a VLAN group response body into a domain.VlanGroup.
func unmarshalVlanGroup(raw domain.RawObject) (*domain.VlanGroup, error) {
	var w WireVlanGroup
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal vlan group: %w", err)
	}
	d := wireVlanGroupToDomain(w)

	return &d, nil
}

// wireVlanGroupToDomain converts a wire VLAN group to the domain model.
func wireVlanGroupToDomain(wire WireVlanGroup) domain.VlanGroup {
	return domain.VlanGroup{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Slug:         wire.Slug,
		Display:      wire.Display,
		Description:  wire.Description,
		Scope:        wireVlanGroupScopeToDomain(wire.Scope),
		MinVID:       wire.MinVID,
		MaxVID:       wire.MaxVID,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

func wireVlanGroupScopeToDomain(scope *WireVlanGroupScope) *domain.VlanGroupScope {
	if scope == nil {
		return nil
	}

	return &domain.VlanGroupScope{
		ObjectType: scope.ObjectType,
		ObjectID:   scope.ObjectID,
	}
}

func roleWriteToWire(input domain.RoleWrite) WireRoleWrite {
	return WireRoleWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Weight:       input.Weight,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalRole decodes an IPAM role response body into a domain.Role.
func unmarshalRole(raw domain.RawObject) (*domain.Role, error) {
	var w WireRole
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal role: %w", err)
	}
	d := wireRoleToDomain(w)

	return &d, nil
}

// wireRoleToDomain converts a wire IPAM role to the domain model.
func wireRoleToDomain(wire WireRole) domain.Role {
	return domain.Role{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Slug:         wire.Slug,
		Weight:       wire.Weight,
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

func contactWriteToWire(input domain.ContactWrite) WireContactWrite {
	return WireContactWrite{
		Name:         input.Name,
		Group:        input.Group,
		Title:        input.Title,
		Phone:        input.Phone,
		Email:        input.Email,
		Address:      input.Address,
		Link:         input.Link,
		Description:  input.Description,
		Comments:     input.Comments,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalContact decodes a tenancy contact response body into a domain.Contact.
func unmarshalContact(raw domain.RawObject) (*domain.Contact, error) {
	var w WireContact
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal contact: %w", err)
	}
	d := wireContactToDomain(w)

	return &d, nil
}

// wireContactToDomain converts a wire tenancy contact to the domain model.
func wireContactToDomain(wire WireContact) domain.Contact {
	return domain.Contact{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Group:        wireNestedToDomain(wire.Group),
		Title:        wire.Title,
		Phone:        wire.Phone,
		Email:        wire.Email,
		Address:      wire.Address,
		Link:         wire.Link,
		Description:  wire.Description,
		Comments:     wire.Comments,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// siteGroupWriteToWire converts a domain site group write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func siteGroupWriteToWire(input domain.SiteGroupWrite) WireSiteGroupWrite {
	return WireSiteGroupWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Parent:       input.Parent,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalSiteGroup decodes a site group response body into a
// domain.SiteGroup.
func unmarshalSiteGroup(raw domain.RawObject) (*domain.SiteGroup, error) {
	var w WireSiteGroup
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal site group: %w", err)
	}
	d := wireSiteGroupToDomain(w)

	return &d, nil
}

// wireSiteGroupToDomain converts a wire site group to the domain model.
func wireSiteGroupToDomain(wire WireSiteGroup) domain.SiteGroup {
	return domain.SiteGroup{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Slug:         wire.Slug,
		Display:      wire.Display,
		Parent:       wireNestedToDomain(wire.Parent),
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// rackGroupWriteToWire converts a domain rack group write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func rackGroupWriteToWire(input domain.RackGroupWrite) WireRackGroupWrite {
	return WireRackGroupWrite{
		Name:        input.Name,
		Slug:        input.Slug,
		Site:        input.Site,
		Parent:      input.Parent,
		Description: input.Description,
		Tags:        input.Tags,
	}
}

// unmarshalRackGroup decodes a rack group response body into a
// domain.RackGroup.
func unmarshalRackGroup(raw domain.RawObject) (*domain.RackGroup, error) {
	var w WireRackGroup
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rack group: %w", err)
	}
	d := wireRackGroupToDomain(w)

	return &d, nil
}

// wireRackGroupToDomain converts a wire rack group to the domain model.
func wireRackGroupToDomain(wire WireRackGroup) domain.RackGroup {
	return domain.RackGroup{
		ID:          wire.ID,
		URL:         wire.URL,
		Name:        wire.Name,
		Slug:        wire.Slug,
		Display:     wire.Display,
		Site:        wireNestedToDomain(wire.Site),
		Parent:      wireNestedToDomain(wire.Parent),
		Description: wire.Description,
		Tags:        wireTagsToDomain(wire.Tags),
		Created:     wire.Created,
		LastUpdated: wire.LastUpdated,
	}
}

// rackRoleWriteToWire converts a domain rack role write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func rackRoleWriteToWire(input domain.RackRoleWrite) WireRackRoleWrite {
	return WireRackRoleWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Color:        input.Color,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalRackRole decodes a rack role response body into a domain.RackRole.
func unmarshalRackRole(raw domain.RawObject) (*domain.RackRole, error) {
	var w WireRackRole
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rack role: %w", err)
	}
	d := wireRackRoleToDomain(w)

	return &d, nil
}

// wireRackRoleToDomain converts a wire rack role to the domain model.
func wireRackRoleToDomain(wire WireRackRole) domain.RackRole {
	return domain.RackRole{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Slug:         wire.Slug,
		Color:        wire.Color,
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// deviceRoleWriteToWire converts a domain device role write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func deviceRoleWriteToWire(input domain.DeviceRoleWrite) WireDeviceRoleWrite {
	return WireDeviceRoleWrite{
		Name:           input.Name,
		Slug:           input.Slug,
		Color:          input.Color,
		VMRole:         input.VMRole,
		ConfigTemplate: input.ConfigTemplate,
		Description:    input.Description,
		Tags:           input.Tags,
		CustomFields:   input.CustomFields,
		Weight:         input.Weight,
	}
}

// unmarshalDeviceRole decodes a device role response body into a
// domain.DeviceRole.
func unmarshalDeviceRole(raw domain.RawObject) (*domain.DeviceRole, error) {
	var w WireDeviceRole
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal device role: %w", err)
	}
	d := wireDeviceRoleToDomain(w)

	return &d, nil
}

// wireDeviceRoleToDomain converts a wire device role to the domain model.
func wireDeviceRoleToDomain(wire WireDeviceRole) domain.DeviceRole {
	return domain.DeviceRole{
		ID:             wire.ID,
		URL:            wire.URL,
		Name:           wire.Name,
		Slug:           wire.Slug,
		Color:          wire.Color,
		VMRole:         wire.VMRole,
		ConfigTemplate: wireNestedToDomain(wire.ConfigTemplate),
		Description:    wire.Description,
		Tags:           wireTagsToDomain(wire.Tags),
		CustomFields:   wire.CustomFields,
		Created:        wire.Created,
		LastUpdated:    wire.LastUpdated,
		Weight:         wire.Weight,
	}
}

// platformWriteToWire converts a domain platform write DTO to the wire request
// model. The fields map one-to-one; nil pointers on the domain side remain nil
// on the wire side so json.Marshal omits them (required for a partial PATCH
// body).
func platformWriteToWire(input domain.PlatformWrite) WirePlatformWrite {
	return WirePlatformWrite{
		Name:           input.Name,
		Slug:           input.Slug,
		Manufacturer:   input.Manufacturer,
		ConfigTemplate: input.ConfigTemplate,
		Description:    input.Description,
		Tags:           input.Tags,
		CustomFields:   input.CustomFields,
	}
}

// unmarshalPlatform decodes a platform response body into a domain.Platform.
func unmarshalPlatform(raw domain.RawObject) (*domain.Platform, error) {
	var w WirePlatform
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal platform: %w", err)
	}
	d := wirePlatformToDomain(w)

	return &d, nil
}

// wirePlatformToDomain converts a wire platform to the domain model.
func wirePlatformToDomain(wire WirePlatform) domain.Platform {
	return domain.Platform{
		ID:             wire.ID,
		URL:            wire.URL,
		Name:           wire.Name,
		Slug:           wire.Slug,
		Manufacturer:   wireNestedToDomain(wire.Manufacturer),
		ConfigTemplate: wireNestedToDomain(wire.ConfigTemplate),
		Description:    wire.Description,
		Tags:           wireTagsToDomain(wire.Tags),
		CustomFields:   wire.CustomFields,
		Created:        wire.Created,
		LastUpdated:    wire.LastUpdated,
	}
}

// virtualChassisWriteToWire converts a domain virtual chassis write DTO to the
// wire request model. The fields map one-to-one; nil pointers on the domain
// side remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func virtualChassisWriteToWire(input domain.VirtualChassisWrite) WireVirtualChassisWrite {
	return WireVirtualChassisWrite{
		Name:         input.Name,
		Domain:       input.Domain,
		Master:       input.Master,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalVirtualChassis decodes a virtual chassis response body into a
// domain.VirtualChassis.
func unmarshalVirtualChassis(raw domain.RawObject) (*domain.VirtualChassis, error) {
	var w WireVirtualChassis
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal virtual chassis: %w", err)
	}
	d := wireVirtualChassisToDomain(w)

	return &d, nil
}

// wireVirtualChassisToDomain converts a wire virtual chassis to the domain
// model.
func wireVirtualChassisToDomain(wire WireVirtualChassis) domain.VirtualChassis {
	return domain.VirtualChassis{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Domain:       wire.Domain,
		Master:       wireNestedToDomain(wire.Master),
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
		MemberCount:  wire.MemberCount,
	}
}

// powerPanelWriteToWire converts a domain power panel write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func powerPanelWriteToWire(input domain.PowerPanelWrite) WirePowerPanelWrite {
	return WirePowerPanelWrite{
		Site:         input.Site,
		Location:     input.Location,
		Name:         input.Name,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalPowerPanel decodes a power panel response body into a
// domain.PowerPanel.
func unmarshalPowerPanel(raw domain.RawObject) (*domain.PowerPanel, error) {
	var w WirePowerPanel
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal power panel: %w", err)
	}
	d := wirePowerPanelToDomain(w)

	return &d, nil
}

// wirePowerPanelToDomain converts a wire power panel to the domain model.
func wirePowerPanelToDomain(wire WirePowerPanel) domain.PowerPanel {
	return domain.PowerPanel{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Display:      wire.Display,
		Site:         wireNestedToDomain(wire.Site),
		Location:     wireNestedToDomain(wire.Location),
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// inventoryItemRoleWriteToWire converts a domain inventory item role write DTO
// to the wire request model. The fields map one-to-one; nil pointers on the
// domain side remain nil on the wire side so json.Marshal omits them (required
// for a partial PATCH body).
func inventoryItemRoleWriteToWire(input domain.InventoryItemRoleWrite) WireInventoryItemRoleWrite {
	return WireInventoryItemRoleWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Color:        input.Color,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalInventoryItemRole decodes an inventory item role response body into
// a domain.InventoryItemRole.
func unmarshalInventoryItemRole(raw domain.RawObject) (*domain.InventoryItemRole, error) {
	var w WireInventoryItemRole
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal inventory item role: %w", err)
	}
	d := wireInventoryItemRoleToDomain(w)

	return &d, nil
}

// wireInventoryItemRoleToDomain converts a wire inventory item role to the
// domain model.
func wireInventoryItemRoleToDomain(wire WireInventoryItemRole) domain.InventoryItemRole {
	return domain.InventoryItemRole{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Slug:         wire.Slug,
		Color:        wire.Color,
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// frontPortWriteToWire converts a domain front port write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func frontPortWriteToWire(input domain.FrontPortWrite) WireFrontPortWrite {
	return WireFrontPortWrite{
		Name:             input.Name,
		Device:           input.Device,
		Type:             input.Type,
		RearPort:         input.RearPort,
		RearPortPosition: input.RearPortPosition,
		Description:      input.Description,
		Tags:             input.Tags,
		CustomFields:     input.CustomFields,
	}
}

// unmarshalFrontPort decodes a front port response body into a
// domain.FrontPort.
func unmarshalFrontPort(raw domain.RawObject) (*domain.FrontPort, error) {
	var w WireFrontPort
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal front port: %w", err)
	}
	d := wireFrontPortToDomain(w)

	return &d, nil
}

// wireFrontPortToDomain converts a wire front port to the domain model.
func wireFrontPortToDomain(wire WireFrontPort) domain.FrontPort {
	return domain.FrontPort{
		ID:               wire.ID,
		URL:              wire.URL,
		Name:             wire.Name,
		Display:          wire.Display,
		Device:           wireNestedToDomain(wire.Device),
		Type:             wireLabelToDomain(wire.Type),
		RearPort:         wireNestedToDomain(wire.RearPort),
		RearPortPosition: wire.RearPortPosition,
		Description:      wire.Description,
		Tags:             wireTagsToDomain(wire.Tags),
		CustomFields:     wire.CustomFields,
		Created:          wire.Created,
		LastUpdated:      wire.LastUpdated,
	}
}

// rearPortWriteToWire converts a domain rear port write DTO to the wire request
// model. The fields map one-to-one; nil pointers on the domain side remain nil
// on the wire side so json.Marshal omits them (required for a partial PATCH
// body).
func rearPortWriteToWire(input domain.RearPortWrite) WireRearPortWrite {
	return WireRearPortWrite{
		Name:         input.Name,
		Device:       input.Device,
		Type:         input.Type,
		Positions:    input.Positions,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalRearPort decodes a rear port response body into a domain.RearPort.
func unmarshalRearPort(raw domain.RawObject) (*domain.RearPort, error) {
	var w WireRearPort
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rear port: %w", err)
	}
	d := wireRearPortToDomain(w)

	return &d, nil
}

// wireRearPortToDomain converts a wire rear port to the domain model.
func wireRearPortToDomain(wire WireRearPort) domain.RearPort {
	return domain.RearPort{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Display:      wire.Display,
		Device:       wireNestedToDomain(wire.Device),
		Type:         wireLabelToDomain(wire.Type),
		Positions:    wire.Positions,
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// consolePortWriteToWire converts a domain console port write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func consolePortWriteToWire(input domain.ConsolePortWrite) WireConsolePortWrite {
	return WireConsolePortWrite{
		Name:         input.Name,
		Device:       input.Device,
		Type:         input.Type,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalConsolePort decodes a console port response body into a
// domain.ConsolePort.
func unmarshalConsolePort(raw domain.RawObject) (*domain.ConsolePort, error) {
	var w WireConsolePort
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal console port: %w", err)
	}
	d := wireConsolePortToDomain(w)

	return &d, nil
}

// wireConsolePortToDomain converts a wire console port to the domain model.
func wireConsolePortToDomain(wire WireConsolePort) domain.ConsolePort {
	return domain.ConsolePort{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Display:      wire.Display,
		Device:       wireNestedToDomain(wire.Device),
		Type:         wireLabelToDomain(wire.Type),
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// consoleServerPortWriteToWire converts a domain console server port write DTO
// to the wire request model. The fields map one-to-one; nil pointers on the
// domain side remain nil on the wire side so json.Marshal omits them (required
// for a partial PATCH body).
func consoleServerPortWriteToWire(input domain.ConsoleServerPortWrite) WireConsoleServerPortWrite {
	return WireConsoleServerPortWrite{
		Name:         input.Name,
		Device:       input.Device,
		Type:         input.Type,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalConsoleServerPort decodes a console server port response body into
// a domain.ConsoleServerPort.
func unmarshalConsoleServerPort(raw domain.RawObject) (*domain.ConsoleServerPort, error) {
	var w WireConsoleServerPort
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal console server port: %w", err)
	}
	d := wireConsoleServerPortToDomain(w)

	return &d, nil
}

// wireConsoleServerPortToDomain converts a wire console server port to the
// domain model.
func wireConsoleServerPortToDomain(wire WireConsoleServerPort) domain.ConsoleServerPort {
	return domain.ConsoleServerPort{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Display:      wire.Display,
		Device:       wireNestedToDomain(wire.Device),
		Type:         wireLabelToDomain(wire.Type),
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// powerPortWriteToWire converts a domain power port write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func powerPortWriteToWire(input domain.PowerPortWrite) WirePowerPortWrite {
	return WirePowerPortWrite{
		Name:          input.Name,
		Device:        input.Device,
		Type:          input.Type,
		MaximumDraw:   input.MaximumDraw,
		AllocatedDraw: input.AllocatedDraw,
		Description:   input.Description,
		Tags:          input.Tags,
		CustomFields:  input.CustomFields,
	}
}

// unmarshalPowerPort decodes a power port response body into a
// domain.PowerPort.
func unmarshalPowerPort(raw domain.RawObject) (*domain.PowerPort, error) {
	var w WirePowerPort
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal power port: %w", err)
	}
	d := wirePowerPortToDomain(w)

	return &d, nil
}

// wirePowerPortToDomain converts a wire power port to the domain model.
func wirePowerPortToDomain(wire WirePowerPort) domain.PowerPort {
	return domain.PowerPort{
		ID:            wire.ID,
		URL:           wire.URL,
		Name:          wire.Name,
		Display:       wire.Display,
		Device:        wireNestedToDomain(wire.Device),
		Type:          wireLabelToDomain(wire.Type),
		MaximumDraw:   wire.MaximumDraw,
		AllocatedDraw: wire.AllocatedDraw,
		Description:   wire.Description,
		Tags:          wireTagsToDomain(wire.Tags),
		CustomFields:  wire.CustomFields,
		Created:       wire.Created,
		LastUpdated:   wire.LastUpdated,
	}
}

// powerOutletWriteToWire converts a domain power outlet write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func powerOutletWriteToWire(input domain.PowerOutletWrite) WirePowerOutletWrite {
	return WirePowerOutletWrite{
		Name:         input.Name,
		Device:       input.Device,
		Type:         input.Type,
		PowerPort:    input.PowerPort,
		FeedLeg:      input.FeedLeg,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalPowerOutlet decodes a power outlet response body into a
// domain.PowerOutlet.
func unmarshalPowerOutlet(raw domain.RawObject) (*domain.PowerOutlet, error) {
	var w WirePowerOutlet
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal power outlet: %w", err)
	}
	d := wirePowerOutletToDomain(w)

	return &d, nil
}

// wirePowerOutletToDomain converts a wire power outlet to the domain model.
func wirePowerOutletToDomain(wire WirePowerOutlet) domain.PowerOutlet {
	return domain.PowerOutlet{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Display:      wire.Display,
		Device:       wireNestedToDomain(wire.Device),
		Type:         wireLabelToDomain(wire.Type),
		PowerPort:    wireNestedToDomain(wire.PowerPort),
		FeedLeg:      wireLabelToDomain(wire.FeedLeg),
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// powerFeedWriteToWire converts a domain power feed write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func powerFeedWriteToWire(input domain.PowerFeedWrite) WirePowerFeedWrite {
	return WirePowerFeedWrite{
		Name:           input.Name,
		PowerPanel:     input.PowerPanel,
		Status:         input.Status,
		Type:           input.Type,
		Supply:         input.Supply,
		Phase:          input.Phase,
		Voltage:        input.Voltage,
		Amperage:       input.Amperage,
		MaxUtilization: input.MaxUtilization,
		Description:    input.Description,
		Tags:           input.Tags,
		CustomFields:   input.CustomFields,
	}
}

// unmarshalPowerFeed decodes a power feed response body into a domain.PowerFeed.
func unmarshalPowerFeed(raw domain.RawObject) (*domain.PowerFeed, error) {
	var w WirePowerFeed
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal power feed: %w", err)
	}
	d := wirePowerFeedToDomain(w)

	return &d, nil
}

// wirePowerFeedToDomain converts a wire power feed to the domain model.
func wirePowerFeedToDomain(wire WirePowerFeed) domain.PowerFeed {
	return domain.PowerFeed{
		ID:             wire.ID,
		URL:            wire.URL,
		Name:           wire.Name,
		Display:        wire.Display,
		PowerPanel:     wireNestedToDomain(wire.PowerPanel),
		Status:         wireLabelToDomain(wire.Status),
		Type:           wireLabelToDomain(wire.Type),
		Supply:         wireLabelToDomain(wire.Supply),
		Phase:          wireLabelToDomain(wire.Phase),
		Voltage:        wire.Voltage,
		Amperage:       wire.Amperage,
		MaxUtilization: wire.MaxUtilization,
		Description:    wire.Description,
		Tags:           wireTagsToDomain(wire.Tags),
		CustomFields:   wire.CustomFields,
		Created:        wire.Created,
		LastUpdated:    wire.LastUpdated,
	}
}

// deviceBayWriteToWire converts a domain device bay write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func deviceBayWriteToWire(input domain.DeviceBayWrite) WireDeviceBayWrite {
	return WireDeviceBayWrite{
		Name:         input.Name,
		Device:       input.Device,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalDeviceBay decodes a device bay response body into a domain.DeviceBay.
func unmarshalDeviceBay(raw domain.RawObject) (*domain.DeviceBay, error) {
	var w WireDeviceBay
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal device bay: %w", err)
	}
	d := wireDeviceBayToDomain(w)

	return &d, nil
}

// wireDeviceBayToDomain converts a wire device bay to the domain model.
func wireDeviceBayToDomain(wire WireDeviceBay) domain.DeviceBay {
	return domain.DeviceBay{
		ID:              wire.ID,
		URL:             wire.URL,
		Name:            wire.Name,
		Display:         wire.Display,
		Device:          wireNestedToDomain(wire.Device),
		InstalledDevice: wireNestedToDomain(wire.InstalledDevice),
		Description:     wire.Description,
		Tags:            wireTagsToDomain(wire.Tags),
		CustomFields:    wire.CustomFields,
		Created:         wire.Created,
		LastUpdated:     wire.LastUpdated,
	}
}

// rackReservationWriteToWire converts a domain rack reservation write DTO to
// the wire request model. The fields map one-to-one; nil pointers on the domain
// side remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func rackReservationWriteToWire(input domain.RackReservationWrite) WireRackReservationWrite {
	return WireRackReservationWrite{
		Description:  input.Description,
		Rack:         input.Rack,
		Units:        input.Units,
		User:         input.User,
		Tenant:       input.Tenant,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalRackReservation decodes a rack reservation response body into a
// domain.RackReservation.
func unmarshalRackReservation(raw domain.RawObject) (*domain.RackReservation, error) {
	var w WireRackReservation
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rack reservation: %w", err)
	}
	d := wireRackReservationToDomain(w)

	return &d, nil
}

// wireRackReservationToDomain converts a wire rack reservation to the domain
// model.
func wireRackReservationToDomain(wire WireRackReservation) domain.RackReservation {
	return domain.RackReservation{
		ID:          wire.ID,
		URL:         wire.URL,
		Rack:        wireNestedToDomain(wire.Rack),
		Units:       wire.Units,
		Description: wire.Description,
		User:        wireNestedToDomain(wire.User),
		Created:     wire.Created,
		LastUpdated: wire.LastUpdated,
	}
}

// inventoryItemWriteToWire converts a domain inventory item write DTO to the
// wire request model. The fields map one-to-one; nil pointers on the domain
// side remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func inventoryItemWriteToWire(input domain.InventoryItemWrite) WireInventoryItemWrite {
	return WireInventoryItemWrite{
		Name:         input.Name,
		Device:       input.Device,
		Role:         input.Role,
		Manufacturer: input.Manufacturer,
		Serial:       input.Serial,
		PartID:       input.PartID,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalInventoryItem decodes an inventory item response body into a
// domain.InventoryItem.
func unmarshalInventoryItem(raw domain.RawObject) (*domain.InventoryItem, error) {
	var w WireInventoryItem
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal inventory item: %w", err)
	}
	d := wireInventoryItemToDomain(w)

	return &d, nil
}

// wireInventoryItemToDomain converts a wire inventory item to the domain model.
func wireInventoryItemToDomain(wire WireInventoryItem) domain.InventoryItem {
	return domain.InventoryItem{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Display:      wire.Display,
		Device:       wireNestedToDomain(wire.Device),
		Role:         wireNestedToDomain(wire.Role),
		Manufacturer: wireNestedToDomain(wire.Manufacturer),
		Serial:       wire.Serial,
		PartID:       wire.PartID,
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// macAddressWriteToWire converts a domain MAC address write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func macAddressWriteToWire(input domain.MACAddressWrite) WireMACAddressWrite {
	return WireMACAddressWrite{
		Address:            input.Address,
		AssignedObjectType: input.AssignedObjectType,
		AssignedObjectID:   input.AssignedObjectID,
		Description:        input.Description,
		Tags:               input.Tags,
		CustomFields:       input.CustomFields,
	}
}

// unmarshalMACAddress decodes a MAC address response body into a
// domain.MACAddress.
func unmarshalMACAddress(raw domain.RawObject) (*domain.MACAddress, error) {
	var w WireMACAddress
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal MAC address: %w", err)
	}
	d := wireMACAddressToDomain(w)

	return &d, nil
}

// wireMACAddressToDomain converts a wire MAC address to the domain model.
func wireMACAddressToDomain(wire WireMACAddress) domain.MACAddress {
	return domain.MACAddress{
		ID:           wire.ID,
		URL:          wire.URL,
		Address:      wire.Address,
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// rirWriteToWire converts a domain RIR write DTO to the wire request model. The
// fields map one-to-one; nil pointers on the domain side remain nil on the wire
// side so json.Marshal omits them (required for a partial PATCH body).
func rirWriteToWire(input domain.RIRWrite) WireRIRWrite {
	return WireRIRWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		IsPrivate:    input.IsPrivate,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalRIR decodes an RIR response body into a domain.RIR.
func unmarshalRIR(raw domain.RawObject) (*domain.RIR, error) {
	var w WireRIR
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal RIR: %w", err)
	}
	d := wireRIRToDomain(w)

	return &d, nil
}

// wireRIRToDomain converts a wire RIR to the domain model.
func wireRIRToDomain(wire WireRIR) domain.RIR {
	return domain.RIR{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Slug:         wire.Slug,
		Display:      wire.Display,
		IsPrivate:    wire.IsPrivate,
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// contactGroupWriteToWire converts a domain contact group write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func contactGroupWriteToWire(input domain.ContactGroupWrite) WireContactGroupWrite {
	return WireContactGroupWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Parent:       input.Parent,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalContactGroup decodes a contact group response body into a
// domain.ContactGroup.
func unmarshalContactGroup(raw domain.RawObject) (*domain.ContactGroup, error) {
	var w WireContactGroup
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal contact group: %w", err)
	}
	d := wireContactGroupToDomain(w)

	return &d, nil
}

// wireContactGroupToDomain converts a wire contact group to the domain model.
func wireContactGroupToDomain(wire WireContactGroup) domain.ContactGroup {
	return domain.ContactGroup{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Slug:         wire.Slug,
		Display:      wire.Display,
		Parent:       wireNestedToDomain(wire.Parent),
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// contactRoleWriteToWire converts a domain contact role write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func contactRoleWriteToWire(input domain.ContactRoleWrite) WireContactRoleWrite {
	return WireContactRoleWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalContactRole decodes a contact role response body into a
// domain.ContactRole.
func unmarshalContactRole(raw domain.RawObject) (*domain.ContactRole, error) {
	var w WireContactRole
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal contact role: %w", err)
	}
	d := wireContactRoleToDomain(w)

	return &d, nil
}

// wireContactRoleToDomain converts a wire contact role to the domain model.
func wireContactRoleToDomain(wire WireContactRole) domain.ContactRole {
	return domain.ContactRole{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Slug:         wire.Slug,
		Display:      wire.Display,
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// virtualDiskWriteToWire converts a domain virtual disk write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func virtualDiskWriteToWire(input domain.VirtualDiskWrite) WireVirtualDiskWrite {
	return WireVirtualDiskWrite{
		Name:           input.Name,
		VirtualMachine: input.VirtualMachine,
		Size:           input.Size,
		Description:    input.Description,
		Tags:           input.Tags,
		CustomFields:   input.CustomFields,
	}
}

// unmarshalVirtualDisk decodes a virtual disk response body into a
// domain.VirtualDisk.
func unmarshalVirtualDisk(raw domain.RawObject) (*domain.VirtualDisk, error) {
	var w WireVirtualDisk
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal virtual disk: %w", err)
	}
	d := wireVirtualDiskToDomain(w)

	return &d, nil
}

// wireVirtualDiskToDomain converts a wire virtual disk to the domain model.
func wireVirtualDiskToDomain(wire WireVirtualDisk) domain.VirtualDisk {
	return domain.VirtualDisk{
		ID:             wire.ID,
		URL:            wire.URL,
		Name:           wire.Name,
		Display:        wire.Display,
		VirtualMachine: wireNestedToDomain(wire.VirtualMachine),
		Size:           wire.Size,
		Description:    wire.Description,
		Tags:           wireTagsToDomain(wire.Tags),
		CustomFields:   wire.CustomFields,
		Created:        wire.Created,
		LastUpdated:    wire.LastUpdated,
	}
}

// virtualMachineTypeWriteToWire converts a domain virtual machine type write
// DTO to the wire request model. The fields map one-to-one; nil pointers on the
// domain side remain nil on the wire side so json.Marshal omits them (required
// for a partial PATCH body).
func virtualMachineTypeWriteToWire(input domain.VirtualMachineTypeWrite) WireVirtualMachineTypeWrite {
	return WireVirtualMachineTypeWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalVirtualMachineType decodes a virtual machine type response body into
// a domain.VirtualMachineType.
func unmarshalVirtualMachineType(raw domain.RawObject) (*domain.VirtualMachineType, error) {
	var w WireVirtualMachineType
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal virtual machine type: %w", err)
	}
	d := wireVirtualMachineTypeToDomain(w)

	return &d, nil
}

// wireVirtualMachineTypeToDomain converts a wire virtual machine type to the
// domain model.
func wireVirtualMachineTypeToDomain(wire WireVirtualMachineType) domain.VirtualMachineType {
	return domain.VirtualMachineType{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Slug:         wire.Slug,
		Display:      wire.Display,
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// aggregateWriteToWire converts a domain aggregate write DTO to the wire request
// model. The fields map one-to-one; nil pointers on the domain side remain nil
// on the wire side so json.Marshal omits them (required for a partial PATCH
// body).
func aggregateWriteToWire(input domain.AggregateWrite) WireAggregateWrite {
	return WireAggregateWrite{
		Prefix:       input.Prefix,
		RIR:          input.RIR,
		Tenant:       input.Tenant,
		DateAdded:    input.DateAdded,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalAggregate decodes an aggregate response body into a domain.Aggregate.
func unmarshalAggregate(raw domain.RawObject) (*domain.Aggregate, error) {
	var w WireAggregate
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal aggregate: %w", err)
	}
	d := wireAggregateToDomain(w)

	return &d, nil
}

// wireAggregateToDomain converts a wire aggregate to the domain model.
func wireAggregateToDomain(wire WireAggregate) domain.Aggregate {
	return domain.Aggregate{
		ID:           wire.ID,
		URL:          wire.URL,
		Prefix:       wire.Prefix,
		Display:      wire.Display,
		RIR:          wireNestedToDomain(wire.RIR),
		Tenant:       wireNestedToDomain(wire.Tenant),
		DateAdded:    wire.DateAdded,
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// ipRangeWriteToWire converts a domain IP range write DTO to the wire request
// model. The fields map one-to-one; nil pointers on the domain side remain nil
// on the wire side so json.Marshal omits them (required for a partial PATCH
// body).
func ipRangeWriteToWire(input domain.IPRangeWrite) WireIPRangeWrite {
	return WireIPRangeWrite{
		StartAddress: input.StartAddress,
		EndAddress:   input.EndAddress,
		VRF:          input.VRF,
		Tenant:       input.Tenant,
		Status:       input.Status,
		Role:         input.Role,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalIPRange decodes an IP range response body into a domain.IPRange.
func unmarshalIPRange(raw domain.RawObject) (*domain.IPRange, error) {
	var w WireIPRange
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal IP range: %w", err)
	}
	d := wireIPRangeToDomain(w)

	return &d, nil
}

// wireIPRangeToDomain converts a wire IP range to the domain model.
func wireIPRangeToDomain(wire WireIPRange) domain.IPRange {
	return domain.IPRange{
		ID:           wire.ID,
		URL:          wire.URL,
		StartAddress: wire.StartAddress,
		EndAddress:   wire.EndAddress,
		Display:      wire.Display,
		VRF:          wireNestedToDomain(wire.VRF),
		Tenant:       wireNestedToDomain(wire.Tenant),
		Status:       wireLabelToDomain(wire.Status),
		Role:         wireNestedToDomain(wire.Role),
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// asnWriteToWire converts a domain ASN write DTO to the wire request model. The
// fields map one-to-one; nil pointers on the domain side remain nil on the wire
// side so json.Marshal omits them (required for a partial PATCH body).
func asnWriteToWire(input domain.ASNWrite) WireASNWrite {
	return WireASNWrite{
		ASN:          input.ASN,
		RIR:          input.RIR,
		Tenant:       input.Tenant,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalASN decodes an ASN response body into a domain.ASN.
func unmarshalASN(raw domain.RawObject) (*domain.ASN, error) {
	var w WireASN
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ASN: %w", err)
	}
	d := wireASNToDomain(w)

	return &d, nil
}

// wireASNToDomain converts a wire ASN to the domain model.
func wireASNToDomain(wire WireASN) domain.ASN {
	return domain.ASN{
		ID:           wire.ID,
		URL:          wire.URL,
		ASN:          wire.ASN,
		Display:      wire.Display,
		RIR:          wireNestedToDomain(wire.RIR),
		Tenant:       wireNestedToDomain(wire.Tenant),
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// asnRangeWriteToWire converts a domain ASN range write DTO to the wire request
// model. The fields map one-to-one; nil pointers on the domain side remain nil
// on the wire side so json.Marshal omits them (required for a partial PATCH
// body).
func asnRangeWriteToWire(input domain.ASNRangeWrite) WireASNRangeWrite {
	return WireASNRangeWrite{
		Name:         input.Name,
		Slug:         input.Slug,
		Start:        input.Start,
		End:          input.End,
		RIR:          input.RIR,
		Tenant:       input.Tenant,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalASNRange decodes an ASN range response body into a domain.ASNRange.
func unmarshalASNRange(raw domain.RawObject) (*domain.ASNRange, error) {
	var w WireASNRange
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ASN range: %w", err)
	}
	d := wireASNRangeToDomain(w)

	return &d, nil
}

// wireASNRangeToDomain converts a wire ASN range to the domain model.
func wireASNRangeToDomain(wire WireASNRange) domain.ASNRange {
	return domain.ASNRange{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Display:      wire.Display,
		Start:        wire.Start,
		End:          wire.End,
		RIR:          wireNestedToDomain(wire.RIR),
		Tenant:       wireNestedToDomain(wire.Tenant),
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// routeTargetWriteToWire converts a domain route target write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a
// partial PATCH body).
func routeTargetWriteToWire(input domain.RouteTargetWrite) WireRouteTargetWrite {
	return WireRouteTargetWrite{
		Name:         input.Name,
		Tenant:       input.Tenant,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalRouteTarget decodes a route target response body into a
// domain.RouteTarget.
func unmarshalRouteTarget(raw domain.RawObject) (*domain.RouteTarget, error) {
	var w WireRouteTarget
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal route target: %w", err)
	}
	d := wireRouteTargetToDomain(w)

	return &d, nil
}

// wireRouteTargetToDomain converts a wire route target to the domain model.
func wireRouteTargetToDomain(wire WireRouteTarget) domain.RouteTarget {
	return domain.RouteTarget{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Display:      wire.Display,
		Tenant:       wireNestedToDomain(wire.Tenant),
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// serviceWriteToWire converts a domain service write DTO to the wire request
// model. The fields map one-to-one; nil pointers on the domain side remain nil
// on the wire side so json.Marshal omits them (required for a partial PATCH
// body).
func serviceWriteToWire(input domain.ServiceWrite) WireServiceWrite {
	return WireServiceWrite{
		Name:             input.Name,
		ParentObjectType: input.ParentObjectType,
		ParentObjectID:   input.ParentObjectID,
		Protocol:         input.Protocol,
		Ports:            input.Ports,
		Description:      input.Description,
		Tags:             input.Tags,
		CustomFields:     input.CustomFields,
	}
}

// unmarshalService decodes a service response body into a domain.Service.
func unmarshalService(raw domain.RawObject) (*domain.Service, error) {
	var w WireService
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal service: %w", err)
	}
	d := wireServiceToDomain(w)

	return &d, nil
}

// wireServiceToDomain converts a wire service to the domain model.
func wireServiceToDomain(wire WireService) domain.Service {
	return domain.Service{
		ID:               wire.ID,
		URL:              wire.URL,
		Name:             wire.Name,
		Display:          wire.Display,
		ParentObjectType: wire.ParentObjectType,
		ParentObjectID:   wire.ParentObjectID,
		Parent:           wireNestedToDomain(wire.Parent),
		Protocol:         wireLabelToDomain(wire.Protocol),
		Ports:            wire.Ports,
		Description:      wire.Description,
		Tags:             wireTagsToDomain(wire.Tags),
		CustomFields:     wire.CustomFields,
		Created:          wire.Created,
		LastUpdated:      wire.LastUpdated,
	}
}

// serviceTemplateWriteToWire converts a domain service template write DTO to the
// wire request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a partial
// PATCH body).
func serviceTemplateWriteToWire(input domain.ServiceTemplateWrite) WireServiceTemplateWrite {
	return WireServiceTemplateWrite{
		Name:         input.Name,
		Protocol:     input.Protocol,
		Ports:        input.Ports,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalServiceTemplate decodes a service template response body into a
// domain.ServiceTemplate.
func unmarshalServiceTemplate(raw domain.RawObject) (*domain.ServiceTemplate, error) {
	var w WireServiceTemplate
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal service template: %w", err)
	}
	d := wireServiceTemplateToDomain(w)

	return &d, nil
}

// wireServiceTemplateToDomain converts a wire service template to the domain
// model.
func wireServiceTemplateToDomain(wire WireServiceTemplate) domain.ServiceTemplate {
	return domain.ServiceTemplate{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Display:      wire.Display,
		Protocol:     wireLabelToDomain(wire.Protocol),
		Ports:        wire.Ports,
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// fhrpGroupWriteToWire converts a domain FHRP group write DTO to the wire
// request model. The fields map one-to-one; nil pointers on the domain side
// remain nil on the wire side so json.Marshal omits them (required for a partial
// PATCH body).
func fhrpGroupWriteToWire(input domain.FHRPGroupWrite) WireFHRPGroupWrite {
	return WireFHRPGroupWrite{
		Name:         input.Name,
		Protocol:     input.Protocol,
		GroupID:      input.GroupID,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalFHRPGroup decodes an FHRP group response body into a domain.FHRPGroup.
func unmarshalFHRPGroup(raw domain.RawObject) (*domain.FHRPGroup, error) {
	var w WireFHRPGroup
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal FHRP group: %w", err)
	}
	d := wireFHRPGroupToDomain(w)

	return &d, nil
}

// wireFHRPGroupToDomain converts a wire FHRP group to the domain model.
func wireFHRPGroupToDomain(wire WireFHRPGroup) domain.FHRPGroup {
	return domain.FHRPGroup{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Display:      wire.Display,
		Protocol:     fhrpProtocolToDomain(wire.Protocol),
		GroupID:      wire.GroupID,
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// fhrpProtocolToDomain converts the FHRP group protocol (a plain string in the
// NetBox JSON API, unlike the ChoiceField objects used elsewhere) into the
// domain Label representation.
func fhrpProtocolToDomain(protocol string) *domain.Label {
	if protocol == "" {
		return nil
	}

	return &domain.Label{Value: protocol}
}

// fhrpGroupAssignmentWriteToWire converts a domain FHRP group assignment write
// DTO to the wire request model. The fields map one-to-one; nil pointers on the
// domain side remain nil on the wire side so json.Marshal omits them (required
// for a partial PATCH body).
func fhrpGroupAssignmentWriteToWire(input domain.FHRPGroupAssignmentWrite) WireFHRPGroupAssignmentWrite {
	return WireFHRPGroupAssignmentWrite{
		Group:         input.Group,
		InterfaceType: input.InterfaceType,
		InterfaceID:   input.InterfaceID,
		Priority:      input.Priority,
	}
}

// unmarshalFHRPGroupAssignment decodes an FHRP group assignment response body
// into a domain.FHRPGroupAssignment.
func unmarshalFHRPGroupAssignment(raw domain.RawObject) (*domain.FHRPGroupAssignment, error) {
	var w WireFHRPGroupAssignment
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal FHRP group assignment: %w", err)
	}
	d := wireFHRPGroupAssignmentToDomain(w)

	return &d, nil
}

// wireFHRPGroupAssignmentToDomain converts a wire FHRP group assignment to the
// domain model.
func wireFHRPGroupAssignmentToDomain(wire WireFHRPGroupAssignment) domain.FHRPGroupAssignment {
	return domain.FHRPGroupAssignment{
		ID:            wire.ID,
		URL:           wire.URL,
		Group:         wireNestedToDomain(wire.Group),
		InterfaceType: wire.InterfaceType,
		InterfaceID:   wire.InterfaceID,
		Interface:     wireNestedToDomain(wire.Interface),
		Priority:      wire.Priority,
		Created:       wire.Created,
		LastUpdated:   wire.LastUpdated,
	}
}

// vlanTranslationPolicyWriteToWire converts a domain VLAN translation policy
// write DTO to the wire request model. The fields map one-to-one; nil pointers
// on the domain side remain nil on the wire side so json.Marshal omits them
// (required for a partial PATCH body).
func vlanTranslationPolicyWriteToWire(input domain.VLANTranslationPolicyWrite) WireVLANTranslationPolicyWrite {
	return WireVLANTranslationPolicyWrite{
		Name:         input.Name,
		Description:  input.Description,
		Tags:         input.Tags,
		CustomFields: input.CustomFields,
	}
}

// unmarshalVLANTranslationPolicy decodes a VLAN translation policy response
// body into a domain.VLANTranslationPolicy.
func unmarshalVLANTranslationPolicy(raw domain.RawObject) (*domain.VLANTranslationPolicy, error) {
	var w WireVLANTranslationPolicy
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal VLAN translation policy: %w", err)
	}
	d := wireVLANTranslationPolicyToDomain(w)

	return &d, nil
}

// wireVLANTranslationPolicyToDomain converts a wire VLAN translation policy to
// the domain model.
func wireVLANTranslationPolicyToDomain(wire WireVLANTranslationPolicy) domain.VLANTranslationPolicy {
	return domain.VLANTranslationPolicy{
		ID:           wire.ID,
		URL:          wire.URL,
		Name:         wire.Name,
		Display:      wire.Display,
		Description:  wire.Description,
		Tags:         wireTagsToDomain(wire.Tags),
		CustomFields: wire.CustomFields,
		Created:      wire.Created,
		LastUpdated:  wire.LastUpdated,
	}
}

// vlanTranslationRuleWriteToWire converts a domain VLAN translation rule write
// DTO to the wire request model. The fields map one-to-one; nil pointers on the
// domain side remain nil on the wire side so json.Marshal omits them (required
// for a partial PATCH body).
func vlanTranslationRuleWriteToWire(input domain.VLANTranslationRuleWrite) WireVLANTranslationRuleWrite {
	return WireVLANTranslationRuleWrite{
		Policy:      input.Policy,
		LocalVID:    input.LocalVID,
		RemoteVID:   input.RemoteVID,
		Description: input.Description,
	}
}

// unmarshalVLANTranslationRule decodes a VLAN translation rule response body
// into a domain.VLANTranslationRule.
func unmarshalVLANTranslationRule(raw domain.RawObject) (*domain.VLANTranslationRule, error) {
	var w WireVLANTranslationRule
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("failed to unmarshal VLAN translation rule: %w", err)
	}
	d := wireVLANTranslationRuleToDomain(w)

	return &d, nil
}

// wireVLANTranslationRuleToDomain converts a wire VLAN translation rule to the
// domain model.
func wireVLANTranslationRuleToDomain(wire WireVLANTranslationRule) domain.VLANTranslationRule {
	return domain.VLANTranslationRule{
		ID:          wire.ID,
		URL:         wire.URL,
		Display:     wire.Display,
		Policy:      wire.Policy,
		LocalVID:    wire.LocalVID,
		RemoteVID:   wire.RemoteVID,
		Description: wire.Description,
		Created:     wire.Created,
		LastUpdated: wire.LastUpdated,
	}
}
