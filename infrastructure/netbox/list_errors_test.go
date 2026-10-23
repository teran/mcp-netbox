package netbox

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestClient_ListIPAMEntities_InvalidJSON exercises the json.Unmarshal error
// branch of every Batch-3 IPAM List method: the endpoint returns HTTP 200 with
// a body that is not valid JSON, so c.get succeeds but unmarshalling into the
// typed page fails and the method must surface an error.
//
//nolint:funlen // Large table-driven case list; splitting is not sensible.
func TestClient_ListIPAMEntities_InvalidJSON(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		path string
		call func(c *Client, ctx context.Context, token string) error
	}{
		{
			name: "RIRs",
			path: "/api/ipam/rirs/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.ListRIRs(ctx, token, nil)

				return err
			},
		},
		{
			name: "Aggregates",
			path: "/api/ipam/aggregates/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.ListAggregates(ctx, token, nil)

				return err
			},
		},
		{
			name: "IPRanges",
			path: "/api/ipam/ip-ranges/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.ListIPRanges(ctx, token, nil)

				return err
			},
		},
		{
			name: "ASNs",
			path: "/api/ipam/asns/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.ListASNs(ctx, token, nil)

				return err
			},
		},
		{
			name: "ASNRanges",
			path: "/api/ipam/asn-ranges/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.ListASNRanges(ctx, token, nil)

				return err
			},
		},
		{
			name: "RouteTargets",
			path: "/api/ipam/route-targets/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.ListRouteTargets(ctx, token, nil)

				return err
			},
		},
		{
			name: "Services",
			path: "/api/ipam/services/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.ListServices(ctx, token, nil)

				return err
			},
		},
		{
			name: "ServiceTemplates",
			path: "/api/ipam/service-templates/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.ListServiceTemplates(ctx, token, nil)

				return err
			},
		},
		{
			name: "FHRPGroups",
			path: "/api/ipam/fhrp-groups/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.ListFHRPGroups(ctx, token, nil)

				return err
			},
		},
		{
			name: "FHRPGroupAssignments",
			path: "/api/ipam/fhrp-group-assignments/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.ListFHRPGroupAssignments(ctx, token, nil)

				return err
			},
		},
		{
			name: "VLANTranslationPolicies",
			path: "/api/ipam/vlan-translation-policies/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.ListVLANTranslationPolicies(ctx, token, nil)

				return err
			},
		},
		{
			name: "VLANTranslationRules",
			path: "/api/ipam/vlan-translation-rules/",
			call: func(c *Client, ctx context.Context, token string) error {
				_, err := c.ListVLANTranslationRules(ctx, token, nil)

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
