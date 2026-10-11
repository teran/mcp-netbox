package netbox

import (
	"strings"
	"testing"

	"github.com/teran/mcp-netbox/domain"
)

func TestNetboxError_Error(t *testing.T) {
	t.Parallel()

	e := &netboxError{StatusCode: 502, Body: []byte("bad gateway")}
	msg := e.Error()
	if !strings.Contains(msg, "502") {
		t.Errorf("Error() = %q, want it to contain the status code", msg)
	}
}

func TestStatusCodeOf_NilResponse(t *testing.T) {
	t.Parallel()

	if got := statusCodeOf(nil); got != 0 {
		t.Errorf("statusCodeOf(nil) = %d, want 0", got)
	}
}

// TestUnmarshalFunctions_InvalidJSON exercises the json.Unmarshal error branch
// of each unmarshal<Type> helper with malformed input.
//
//nolint:funlen // table-driven test over all unmarshal helpers; splitting is not sensible.
//nolint:funlen // table-driven test over all unmarshal helpers; splitting is not sensible.
func TestUnmarshalFunctions_InvalidJSON(t *testing.T) {
	t.Parallel()

	bad := domain.RawObject([]byte(`{"id":`)) // truncated -> unmarshal error

	cases := []struct {
		name string
		call func() error
	}{
		{"site", func() error {
			_, err := unmarshalSite(bad)

			return err
		}},
		{"device", func() error {
			_, err := unmarshalDevice(bad)

			return err
		}},
		{"ip_address", func() error {
			_, err := unmarshalIPAddress(bad)

			return err
		}},
		{"prefix", func() error {
			_, err := unmarshalPrefix(bad)

			return err
		}},
		{"vlan", func() error {
			_, err := unmarshalVLAN(bad)

			return err
		}},
		{"virtual_machine", func() error {
			_, err := unmarshalVirtualMachine(bad)

			return err
		}},
		{"cluster", func() error {
			_, err := unmarshalCluster(bad)

			return err
		}},
		{"circuit", func() error {
			_, err := unmarshalCircuit(bad)

			return err
		}},
		{"rack", func() error {
			_, err := unmarshalRack(bad)

			return err
		}},
		{"interface", func() error {
			_, err := unmarshalInterface(bad)

			return err
		}},
		{"circuit_termination", func() error {
			_, err := unmarshalCircuitTermination(bad)

			return err
		}},
		{"location", func() error {
			_, err := unmarshalLocation(bad)

			return err
		}},
		{"cluster_type", func() error {
			_, err := unmarshalClusterType(bad)

			return err
		}},
		{"cluster_group", func() error {
			_, err := unmarshalClusterGroup(bad)

			return err
		}},
		{"role", func() error {
			_, err := unmarshalRole(bad)

			return err
		}},
		{"contact", func() error {
			_, err := unmarshalContact(bad)

			return err
		}},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if err := testCase.call(); err == nil {
				t.Errorf("%s: unmarshal = nil, want error for invalid json", testCase.name)
			}
		})
	}
}
