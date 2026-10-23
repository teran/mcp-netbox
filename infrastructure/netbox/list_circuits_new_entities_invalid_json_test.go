package netbox

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestClient_ListCircuitsNewEntities_InvalidJSON exercises the json.Unmarshal
// error branch of every Batch-5 circuits List method (provider_account,
// provider_network, circuit_group, circuit_group_assignment, virtual_circuit,
// virtual_circuit_type, virtual_circuit_termination): the endpoint returns
// HTTP 200 with a body that is not valid JSON, so c.get succeeds but
// unmarshalling into the typed page fails and the method must surface an error.
//
//nolint:funlen // Large table-driven case list; splitting is not sensible.
func TestClient_ListCircuitsNewEntities_InvalidJSON(t *testing.T) {
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
				_, err := c.ListProviderAccounts(ctx, token, nil)

				return err
			},
		},
		{
			name: "ProviderNetworks",
			path: "/api/circuits/provider-networks/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.ListProviderNetworks(ctx, token, nil)

				return err
			},
		},
		{
			name: "CircuitGroups",
			path: "/api/circuits/circuit-groups/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.ListCircuitGroups(ctx, token, nil)

				return err
			},
		},
		{
			name: "CircuitGroupAssignments",
			path: "/api/circuits/circuit-group-assignments/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.ListCircuitGroupAssignments(ctx, token, nil)

				return err
			},
		},
		{
			name: "VirtualCircuits",
			path: "/api/circuits/virtual-circuits/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.ListVirtualCircuits(ctx, token, nil)

				return err
			},
		},
		{
			name: "VirtualCircuitTypes",
			path: "/api/circuits/virtual-circuit-types/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.ListVirtualCircuitTypes(ctx, token, nil)

				return err
			},
		},
		{
			name: "VirtualCircuitTerminations",
			path: "/api/circuits/virtual-circuit-terminations/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.ListVirtualCircuitTerminations(ctx, token, nil)

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
				// 200 with a body that is not valid JSON.
				_, _ = writer.Write([]byte(`{"count": "not-a-number", "results": [`))
			}))
			defer srv.Close()

			client := NewClient(srv.URL, http.DefaultClient)
			if err := tcCase.call(client, context.Background(), "test-token"); err == nil {
				t.Fatalf("%s: expected unmarshal error, got nil", tcCase.name)
			}
		})
	}
}
