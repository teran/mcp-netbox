package netbox

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/teran/mcp-netbox/domain"
)

// TestClient_CreateCircuitsNewEntities_InvalidJSON exercises the unmarshal
// error branch of the Create write path for every Batch-5 circuits entity
// (provider_account, provider_network, circuit_group, circuit_group_assignment,
// virtual_circuit, virtual_circuit_type, virtual_circuit_termination): the
// endpoint returns HTTP 201 (create success) but with a body that is not valid
// JSON, so the write succeeds at the transport layer but unmarshalling the
// created object fails and the method must surface an error.
//
//nolint:funlen // Large table-driven case list; splitting is not sensible.
func TestClient_CreateCircuitsNewEntities_InvalidJSON(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		path string
		call func(c *Client, ctx context.Context, token string) error
	}{
		{
			name: "ProviderAccounts",
			path: "/api/circuits/provider-accounts/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.CreateProviderAccount(ctx, token, domain.ProviderAccountWrite{Name: "ACC-1"})

				return err
			},
		},
		{
			name: "ProviderNetworks",
			path: "/api/circuits/provider-networks/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.CreateProviderNetwork(ctx, token, domain.ProviderNetworkWrite{Name: "Net A"})

				return err
			},
		},
		{
			name: "CircuitGroups",
			path: "/api/circuits/circuit-groups/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.CreateCircuitGroup(ctx, token, domain.CircuitGroupWrite{Name: "Grp"})

				return err
			},
		},
		{
			name: "CircuitGroupAssignments",
			path: "/api/circuits/circuit-group-assignments/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.CreateCircuitGroupAssignment(
					ctx, token, domain.CircuitGroupAssignmentWrite{Group: 5, MemberType: "circuits.circuit", MemberID: 6},
				)

				return err
			},
		},
		{
			name: "VirtualCircuits",
			path: "/api/circuits/virtual-circuits/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.CreateVirtualCircuit(ctx, token, domain.VirtualCircuitWrite{CID: "VC-1"})

				return err
			},
		},
		{
			name: "VirtualCircuitTypes",
			path: "/api/circuits/virtual-circuit-types/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.CreateVirtualCircuitType(ctx, token, domain.VirtualCircuitTypeWrite{Name: "VPLS"})

				return err
			},
		},
		{
			name: "VirtualCircuitTerminations",
			path: "/api/circuits/virtual-circuit-terminations/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.CreateVirtualCircuitTermination(
					ctx, token, domain.VirtualCircuitTerminationWrite{VirtualCircuit: 5},
				)

				return err
			},
		},
	}

	for _, tcCase := range cases {
		t.Run(tcCase.name, func(t *testing.T) {
			t.Parallel()
			srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
				if req.URL.Path != tcCase.path {
					writer.WriteHeader(http.StatusNotFound)

					return
				}
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusCreated)
				// 201 with a body that is not valid JSON.
				_, _ = writer.Write([]byte(`{"id": 1, "name": [`))
			}))
			defer srv.Close()

			client := NewClient(srv.URL, http.DefaultClient)
			if err := tcCase.call(client, context.Background(), "test-token"); err == nil {
				t.Fatalf("%s: expected unmarshal error, got nil", tcCase.name)
			}
		})
	}
}
