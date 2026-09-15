//go:build e2e

// Package e2e contains an end-to-end integration test that runs a real NetBox
// (via go-docker-testsuite) and drives the MCP server's tools over the full MCP
// JSON-RPC protocol, in-process and without OS IPC.
//
// This test is gated behind the `e2e` build tag because it requires a running
// Docker daemon and pulls a multi-container NetBox stack. The default unit run
// (`go test ./...`) builds without the tag, so this test is excluded; it is
// exercised by `make test-e2e` / `go test -tags e2e ./...`. The runtime guard
// (MCP_NETBOX_E2E=1) is kept as an additional safety check so the test is a
// no-op skip unless the operator explicitly opts in.
package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/teran/go-docker-testsuite/applications/netbox"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/handlers"
	infra "github.com/teran/mcp-netbox/infrastructure/netbox"
)

const netboxImage = "index.docker.io/netboxcommunity/netbox:v4.6-5.0.1"

// TestNetBoxE2E spins up a real NetBox and exercises CRUD for representative
// entities across NetBox applications (dcim, ipam, virtualization, circuits,
// tenancy) through the MCP protocol.
func TestNetBoxE2E(t *testing.T) {
	if os.Getenv("MCP_NETBOX_E2E") == "" {
		t.Skip("set MCP_NETBOX_E2E=1 to run the NetBox e2e test")
	}

	// First start runs ~200 DB migrations; allow up to 12 minutes for readiness.
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Minute)
	defer cancel()

	app, err := netbox.NewWithT(t, ctx, netboxImage)
	if err != nil {
		t.Fatalf("start netbox: %v", err)
	}
	defer func() { _ = app.Close(context.Background()) }()

	// Plain http.Client is required: the production server's DNS-rebinding dialer
	// rejects loopback, but NetBox publishes on 127.0.0.1:<random>.
	netboxClient := infra.NewClient(app.MustURL(), &http.Client{Timeout: 30 * time.Second})
	svc := application.NewNetworkService(netboxClient, app.SuperuserAPIToken())

	client := newE2E(t, svc)
	defer client.close()

	// tools/list must expose the full tool inventory, including the write tools.
	tools := client.listToolNames()
	for _, want := range []string{"get_sites", "create_site", "update_site", "delete_site", "create_circuit", "delete_tenant"} {
		if !contains(tools, want) {
			t.Fatalf("tools/list missing %q; got %d tools", want, len(tools))
		}
	}
	t.Logf("tools/list returned %d tools", len(tools))

	runID := strconv.FormatInt(time.Now().UnixNano(), 36)
	t.Run("site", func(t *testing.T) { siteCRUD(t, client, runID) })
	t.Run("prefix", func(t *testing.T) { prefixCRUD(t, client, runID) })
	t.Run("ip_address", func(t *testing.T) { ipAddressCRUD(t, client, runID) })
	t.Run("virtual_machine", func(t *testing.T) { virtualMachineCRUD(t, client, runID) })
	t.Run("tenant", func(t *testing.T) { tenantCRUD(t, client, runID) })
	t.Run("circuit", func(t *testing.T) { circuitCRUD(t, client, runID) })
}

// e2eClient is a thin typed wrapper over an MCP ClientSession that the CRUD
// helpers share. It forwards every call over the in-memory pipes.
type e2eClient struct {
	t  *testing.T
	cs *mcp.ClientSession
	// sessCtx lives for the duration of the MCP session; derived from
	// context.Background() so it is independent of the NetBox readiness ctx.
	sessCtx    context.Context
	sessCancel context.CancelFunc
}

// newE2E wires an in-process MCP server (with all tools registered on a real
// NetBox-backed service) to an SDK client over a pair of in-memory pipes.
func newE2E(t *testing.T, svc *application.NetworkService) *e2eClient {
	t.Helper()

	// Two pairs of pipes: requests flow client->server, responses server->client.
	reqR, reqW := io.Pipe()   // client writes requests here
	respR, respW := io.Pipe() // server writes responses here

	// Server reads requests from reqR and writes responses to respW.
	serverTransport := &mcp.IOTransport{Reader: reqR, Writer: respW}
	// Client writes requests to reqW and reads responses from respR.
	clientTransport := &mcp.IOTransport{Reader: respR, Writer: reqW}

	srv := mcp.NewServer(
		&mcp.Implementation{Name: "mcp-netbox", Version: "dev"},
		&mcp.ServerOptions{
			Capabilities: &mcp.ServerCapabilities{Tools: &mcp.ToolCapabilities{ListChanged: false}},
		},
	)
	handlers.RegisterTools(srv, nil, svc)

	sessCtx, sessCancel := context.WithCancel(context.Background())

	conn, err := srv.Connect(sessCtx, serverTransport, nil)
	if err != nil {
		sessCancel()
		t.Fatalf("server connect: %v", err)
	}
	go func() { _ = conn.Wait() }()

	client := mcp.NewClient(&mcp.Implementation{Name: "e2e", Version: "1"}, nil)
	cs, err := client.Connect(sessCtx, clientTransport, &mcp.ClientSessionOptions{})
	if err != nil {
		sessCancel()
		t.Fatalf("client connect: %v", err)
	}

	return &e2eClient{t: t, cs: cs, sessCtx: sessCtx, sessCancel: sessCancel}
}

func (c *e2eClient) close() {
	c.sessCancel()
	_ = c.cs.Close()
}

// listToolNames returns the names of all tools the server advertises.
func (c *e2eClient) listToolNames() []string {
	c.t.Helper()
	res, err := c.cs.ListTools(c.sessCtx, &mcp.ListToolsParams{})
	if err != nil {
		c.t.Fatalf("tools/list: %v", err)
	}
	names := make([]string, 0, len(res.Tools))
	for _, tl := range res.Tools {
		names = append(names, tl.Name)
	}
	return names
}

// call invokes a tool and fatals if the call fails or returns an error result.
func (c *e2eClient) call(name string, args map[string]interface{}) *mcp.CallToolResult {
	c.t.Helper()
	res, err := c.cs.CallTool(c.sessCtx, &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		c.t.Fatalf("call %s: %v", name, err)
	}
	if res.IsError {
		c.t.Fatalf("call %s returned error: %s", name, c.resultText(res))
	}
	return res
}

// callExpectError invokes a tool and asserts it returns an error result (used
// to confirm a delete actually removed the object: a get on a deleted object
// must fail with not-found).
func (c *e2eClient) callExpectError(name string, args map[string]interface{}) *mcp.CallToolResult {
	c.t.Helper()
	res, err := c.cs.CallTool(c.sessCtx, &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		c.t.Fatalf("call %s: %v", name, err)
	}
	if !res.IsError {
		c.t.Fatalf("call %s: expected error, got success: %s", name, c.resultText(res))
	}
	return res
}

// getByID fetches an object by its typed ID via get_object_by_id.
func (c *e2eClient) getByID(objectType string, id int) *mcp.CallToolResult {
	return c.call("get_object_by_id", map[string]interface{}{
		"object_type": objectType,
		"id":          id,
	})
}

// data extracts the "data" object from a tool result's structured content. The
// typed outputs are wrapped as {"data": {...}} on the wire.
func (c *e2eClient) data(res *mcp.CallToolResult) map[string]interface{} {
	c.t.Helper()
	m, ok := res.StructuredContent.(map[string]interface{})
	if !ok {
		c.t.Fatalf("structured content is not an object: %T (%v)", res.StructuredContent, res.StructuredContent)
	}
	data, ok := m["data"].(map[string]interface{})
	if !ok {
		c.t.Fatalf("result has no data object: %v", m)
	}
	return data
}

// mustID returns the numeric id of a created/updated object.
func (c *e2eClient) mustID(res *mcp.CallToolResult) int {
	c.t.Helper()
	data := c.data(res)
	id, ok := num(data["id"])
	if !ok {
		c.t.Fatalf("result has no numeric id: %v", data)
	}
	return id
}

// str returns a string field from an object data map.
func (c *e2eClient) str(data map[string]interface{}, key string) string {
	c.t.Helper()
	s, _ := data[key].(string)
	return s
}

// resultText renders the error text of an error result for diagnostics.
func (c *e2eClient) resultText(res *mcp.CallToolResult) string {
	for _, ct := range res.Content {
		if tc, ok := ct.(*mcp.TextContent); ok {
			return tc.Text
		}
	}
	if b, err := json.Marshal(res.StructuredContent); err == nil {
		return string(b)
	}
	return "(no detail)"
}

// num coerces a JSON number value (float64 or json.Number) to an int.
func num(v interface{}) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case json.Number:
		i, err := n.Int64()
		return int(i), err == nil
	case int:
		return n, true
	case int64:
		return int(n), true
	}
	return 0, false
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// --- CRUD helpers -----------------------------------------------------------

func siteCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-site-" + runID

	// NetBox 4.6 requires slug explicitly on site creation (it is not auto-derived).
	created := c.call("create_site", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created site name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("site", id)); c.str(data, "name") != name {
		t.Fatalf("get site name = %q, want %q", c.str(data, "name"), name)
	}

	updated := c.call("update_site", map[string]interface{}{"id": id, "description": "updated"})
	if got := c.str(c.data(updated), "description"); got != "updated" {
		t.Fatalf("updated site description = %q, want %q", got, "updated")
	}
	if data := c.data(c.getByID("site", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get after update description = %q, want %q", c.str(data, "description"), "updated")
	}

	c.call("delete_site", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "site", "id": id})
}

func prefixCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	cidr := "10." + fmt.Sprintf("%d", 100+hashSuffix(runID)%150) + ".0.0/24"

	created := c.call("create_prefix", map[string]interface{}{"prefix": cidr})
	id := c.mustID(created)

	if data := c.data(c.getByID("prefix", id)); c.str(data, "prefix") != cidr {
		t.Fatalf("get prefix = %q, want %q", c.str(data, "prefix"), cidr)
	}

	c.call("update_prefix", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("prefix", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get prefix after update description = %q", c.str(data, "description"))
	}

	c.call("delete_prefix", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "prefix", "id": id})
}

func ipAddressCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	addr := "192.168." + fmt.Sprintf("%d", hashSuffix(runID)%200) + ".1/24"

	created := c.call("create_ip_address", map[string]interface{}{"address": addr})
	id := c.mustID(created)

	if data := c.data(c.getByID("ip_address", id)); c.str(data, "address") != addr {
		t.Fatalf("get ip_address = %q, want %q", c.str(data, "address"), addr)
	}

	c.call("update_ip_address", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("ip_address", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get ip_address after update description = %q", c.str(data, "description"))
	}

	c.call("delete_ip_address", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "ip_address", "id": id})
}

func virtualMachineCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-vm-" + runID

	// A virtual machine must be assigned to a site, cluster, or device, so create
	// a site to host it (and clean it up afterwards).
	siteName := "e2e-vmsite-" + runID
	site := c.call("create_site", map[string]interface{}{"name": siteName, "slug": siteName})
	siteID := c.mustID(site)
	t.Cleanup(func() { c.call("delete_site", map[string]interface{}{"id": siteID}) })

	created := c.call("create_virtual_machine", map[string]interface{}{"name": name, "site": siteID})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created vm name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("virtual_machine", id)); c.str(data, "name") != name {
		t.Fatalf("get vm name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_virtual_machine", map[string]interface{}{"id": id, "vcpus": float64(4)})
	data := c.data(c.getByID("virtual_machine", id))
	if vc, ok := num(data["vcpus"]); !ok || vc != 4 {
		t.Fatalf("get vm vcpus after update = %v, want 4", data["vcpus"])
	}

	c.call("delete_virtual_machine", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "virtual_machine", "id": id})
}

func tenantCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()
	name := "e2e-tenant-" + runID

	created := c.call("create_tenant", map[string]interface{}{"name": name, "slug": name})
	id := c.mustID(created)
	if got := c.str(c.data(created), "name"); got != name {
		t.Fatalf("created tenant name = %q, want %q", got, name)
	}

	if data := c.data(c.getByID("tenant", id)); c.str(data, "name") != name {
		t.Fatalf("get tenant name = %q, want %q", c.str(data, "name"), name)
	}

	c.call("update_tenant", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("tenant", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get tenant after update description = %q", c.str(data, "description"))
	}

	c.call("delete_tenant", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "tenant", "id": id})
}

func circuitCRUD(t *testing.T, c *e2eClient, runID string) {
	t.Helper()

	// circuit requires a provider and a circuit_type, so create those first.
	// NetBox 4.6 requires an explicit slug for provider and circuit_type.
	provider := c.call("create_provider", map[string]interface{}{"name": "e2e-prov-" + runID, "slug": "e2e-prov-" + runID})
	providerID := c.mustID(provider)
	circuitType := c.call("create_circuit_type", map[string]interface{}{"name": "e2e-ct-" + runID, "slug": "e2e-ct-" + runID})
	circuitTypeID := c.mustID(circuitType)

	// Ensure the supporting objects are cleaned up too.
	t.Cleanup(func() {
		c.call("delete_provider", map[string]interface{}{"id": providerID})
		c.call("delete_circuit_type", map[string]interface{}{"id": circuitTypeID})
	})

	cid := "e2e-cid-" + runID
	created := c.call("create_circuit", map[string]interface{}{
		"cid":          cid,
		"provider":     providerID,
		"circuit_type": circuitTypeID,
	})
	id := c.mustID(created)
	if got := c.str(c.data(created), "cid"); got != cid {
		t.Fatalf("created circuit cid = %q, want %q", got, cid)
	}

	if data := c.data(c.getByID("circuit", id)); c.str(data, "cid") != cid {
		t.Fatalf("get circuit cid = %q, want %q", c.str(data, "cid"), cid)
	}

	c.call("update_circuit", map[string]interface{}{"id": id, "description": "updated"})
	if data := c.data(c.getByID("circuit", id)); c.str(data, "description") != "updated" {
		t.Fatalf("get circuit after update description = %q", c.str(data, "description"))
	}

	c.call("delete_circuit", map[string]interface{}{"id": id})
	c.callExpectError("get_object_by_id", map[string]interface{}{"object_type": "circuit", "id": id})
}

// hashSuffix returns a small deterministic-ish number derived from runID so
// each invocation uses a distinct CIDR/octet within NetBox's valid ranges.
func hashSuffix(runID string) int {
	sum := 0
	for _, r := range runID {
		sum = (sum*31 + int(r)) % 997
	}
	return sum
}
