package netbox

import (
	"testing"

	"github.com/teran/mcp-netbox/domain"
)

// TestUnmarshalIPAMEntities_InvalidJSON exercises the json.Unmarshal error
// branch of every Batch-3 unmarshal helper: a malformed raw object must surface
// an error rather than silently producing a zero-value domain object.
func TestUnmarshalIPAMEntities_InvalidJSON(t *testing.T) { //nolint:funlen
	t.Parallel()

	cases := []struct {
		name string
		fn   func(raw domain.RawObject) error
	}{
		{
			name: "RIR",
			fn: func(raw domain.RawObject) error {
				_, err := unmarshalRIR(raw)

				return err
			},
		},
		{
			name: "Aggregate",
			fn: func(raw domain.RawObject) error {
				_, err := unmarshalAggregate(raw)

				return err
			},
		},
		{
			name: "IPRange",
			fn: func(raw domain.RawObject) error {
				_, err := unmarshalIPRange(raw)

				return err
			},
		},
		{
			name: "ASN",
			fn: func(raw domain.RawObject) error {
				_, err := unmarshalASN(raw)

				return err
			},
		},
		{
			name: "ASNRange",
			fn: func(raw domain.RawObject) error {
				_, err := unmarshalASNRange(raw)

				return err
			},
		},
		{
			name: "RouteTarget",
			fn: func(raw domain.RawObject) error {
				_, err := unmarshalRouteTarget(raw)

				return err
			},
		},
		{
			name: "Service",
			fn: func(raw domain.RawObject) error {
				_, err := unmarshalService(raw)

				return err
			},
		},
		{
			name: "ServiceTemplate",
			fn: func(raw domain.RawObject) error {
				_, err := unmarshalServiceTemplate(raw)

				return err
			},
		},
		{
			name: "FHRPGroup",
			fn: func(raw domain.RawObject) error {
				_, err := unmarshalFHRPGroup(raw)

				return err
			},
		},
		{
			name: "FHRPGroupAssignment",
			fn: func(raw domain.RawObject) error {
				_, err := unmarshalFHRPGroupAssignment(raw)

				return err
			},
		},
		{
			name: "VLANTranslationPolicy",
			fn: func(raw domain.RawObject) error {
				_, err := unmarshalVLANTranslationPolicy(raw)

				return err
			},
		},
		{
			name: "VLANTranslationRule",
			fn: func(raw domain.RawObject) error {
				_, err := unmarshalVLANTranslationRule(raw)

				return err
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := tc.fn(domain.RawObject(`{"id":`)); err == nil {
				t.Fatalf("%s: expected unmarshal error, got nil", tc.name)
			}
		})
	}
}
