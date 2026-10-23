package netbox

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/teran/mcp-netbox/domain"
)

// TestListEntities_UnmarshalError verifies that each list client method
// surfaces an error when the response body is not valid JSON (covering the
// json.Unmarshal error branch of every list method).
//
//nolint:funlen // Table-driven list of every entity list method.
func TestListEntities_UnmarshalError(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		list func(c *Client) error
	}{
		{
			"ListRegions", func(c *Client) error {
				_, err := c.ListRegions(context.Background(), "t", nil)

				return err
			},
		},
		{
			"ListTenantGroups", func(c *Client) error {
				_, err := c.ListTenantGroups(context.Background(), "t", nil)

				return err
			},
		},
		{
			"ListFrontPorts", func(c *Client) error {
				_, err := c.ListFrontPorts(context.Background(), "t", nil)

				return err
			},
		},
		{
			"ListRearPorts", func(c *Client) error {
				_, err := c.ListRearPorts(context.Background(), "t", nil)

				return err
			},
		},
		{
			"ListConsolePorts", func(c *Client) error {
				_, err := c.ListConsolePorts(context.Background(), "t", nil)

				return err
			},
		},
		{
			"ListConsoleServerPorts", func(c *Client) error {
				_, err := c.ListConsoleServerPorts(context.Background(), "t", nil)

				return err
			},
		},
		{
			"ListPowerPorts", func(c *Client) error {
				_, err := c.ListPowerPorts(context.Background(), "t", nil)

				return err
			},
		},
		{
			"ListPowerOutlets", func(c *Client) error {
				_, err := c.ListPowerOutlets(context.Background(), "t", nil)

				return err
			},
		},
		{
			"ListPowerFeeds", func(c *Client) error {
				_, err := c.ListPowerFeeds(context.Background(), "t", nil)

				return err
			},
		},
		{
			"ListDeviceBays", func(c *Client) error {
				_, err := c.ListDeviceBays(context.Background(), "t", nil)

				return err
			},
		},
		{
			"ListRackReservations", func(c *Client) error {
				_, err := c.ListRackReservations(context.Background(), "t", nil)

				return err
			},
		},
		{
			"ListInventoryItems", func(c *Client) error {
				_, err := c.ListInventoryItems(context.Background(), "t", nil)

				return err
			},
		},
		{
			"ListMACAddresses", func(c *Client) error {
				_, err := c.ListMACAddresses(context.Background(), "t", nil)

				return err
			},
		},
		{
			"ListContactGroups", func(c *Client) error {
				_, err := c.ListContactGroups(context.Background(), "t", nil)

				return err
			},
		},
		{
			"ListContactRoles", func(c *Client) error {
				_, err := c.ListContactRoles(context.Background(), "t", nil)

				return err
			},
		},
		{
			"ListVirtualDisks", func(c *Client) error {
				_, err := c.ListVirtualDisks(context.Background(), "t", nil)

				return err
			},
		},
		{
			"ListVirtualMachineTypes", func(c *Client) error {
				_, err := c.ListVirtualMachineTypes(context.Background(), "t", nil)

				return err
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`not-json`))
			}))
			defer srv.Close()

			client := NewClient(srv.URL, http.DefaultClient)
			if err := testCase.list(client); err == nil {
				t.Errorf("%s: expected unmarshal error, got nil", testCase.name)
			}
		})
	}
}

// TestUnmarshalEntities_Error verifies that each unmarshal helper returns an
// error when handed invalid JSON (covering their json.Unmarshal error branch).
//
//nolint:funlen // Table-driven list of every unmarshal helper.
func TestUnmarshalEntities_Error(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		fn   func(domain.RawObject) error
	}{
		{
			"unmarshalRegion", func(raw domain.RawObject) error {
				_, err := unmarshalRegion(raw)

				return err
			},
		},
		{
			"unmarshalTenantGroup", func(raw domain.RawObject) error {
				_, err := unmarshalTenantGroup(raw)

				return err
			},
		},
		{
			"unmarshalFrontPort", func(raw domain.RawObject) error {
				_, err := unmarshalFrontPort(raw)

				return err
			},
		},
		{
			"unmarshalRearPort", func(raw domain.RawObject) error {
				_, err := unmarshalRearPort(raw)

				return err
			},
		},
		{
			"unmarshalConsolePort", func(raw domain.RawObject) error {
				_, err := unmarshalConsolePort(raw)

				return err
			},
		},
		{
			"unmarshalConsoleServerPort", func(raw domain.RawObject) error {
				_, err := unmarshalConsoleServerPort(raw)

				return err
			},
		},
		{
			"unmarshalPowerPort", func(raw domain.RawObject) error {
				_, err := unmarshalPowerPort(raw)

				return err
			},
		},
		{
			"unmarshalPowerOutlet", func(raw domain.RawObject) error {
				_, err := unmarshalPowerOutlet(raw)

				return err
			},
		},
		{
			"unmarshalPowerFeed", func(raw domain.RawObject) error {
				_, err := unmarshalPowerFeed(raw)

				return err
			},
		},
		{
			"unmarshalDeviceBay", func(raw domain.RawObject) error {
				_, err := unmarshalDeviceBay(raw)

				return err
			},
		},
		{
			"unmarshalRackReservation", func(raw domain.RawObject) error {
				_, err := unmarshalRackReservation(raw)

				return err
			},
		},
		{
			"unmarshalInventoryItem", func(raw domain.RawObject) error {
				_, err := unmarshalInventoryItem(raw)

				return err
			},
		},
		{
			"unmarshalMACAddress", func(raw domain.RawObject) error {
				_, err := unmarshalMACAddress(raw)

				return err
			},
		},
		{
			"unmarshalContactGroup", func(raw domain.RawObject) error {
				_, err := unmarshalContactGroup(raw)

				return err
			},
		},
		{
			"unmarshalContactRole", func(raw domain.RawObject) error {
				_, err := unmarshalContactRole(raw)

				return err
			},
		},
		{
			"unmarshalVirtualDisk", func(raw domain.RawObject) error {
				_, err := unmarshalVirtualDisk(raw)

				return err
			},
		},
		{
			"unmarshalVirtualMachineType", func(raw domain.RawObject) error {
				_, err := unmarshalVirtualMachineType(raw)

				return err
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			if err := testCase.fn(domain.RawObject(`{`)); err == nil {
				t.Errorf("%s: expected unmarshal error, got nil", testCase.name)
			}
		})
	}
}
