package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/teran/mcp-netbox/domain"
)

func TestWriteErrorResult_PlainError(t *testing.T) {
	t.Parallel()

	// A non-ValidationError must produce an IsError result and pass the error
	// through unchanged.
	err := errors.New("boom")
	res, e := writeErrorResult(err)
	if !res.IsError {
		t.Error("res.IsError = false, want true")
	}
	if !errors.Is(e, err) {
		t.Errorf("error = %v, want %v", e, err)
	}
	if res.StructuredContent != nil {
		t.Errorf("StructuredContent = %v, want nil for non-validation error", res.StructuredContent)
	}
}

func TestWriteErrorResult_ValidationError(t *testing.T) {
	t.Parallel()

	ve := &domain.ValidationError{StatusCode: 400, Body: []byte(`{"slug":["required"]}`)}
	res, e := writeErrorResult(ve)
	if e != nil {
		t.Errorf("error = %v, want nil", e)
	}
	if !res.IsError {
		t.Error("res.IsError = false, want true")
	}
	if res.StructuredContent == nil {
		t.Error("StructuredContent = nil, want validation_errors")
	}
}

func TestStrPtr(t *testing.T) {
	t.Parallel()

	if got := strPtr(""); got != nil {
		t.Errorf("strPtr(\"\") = %v, want nil", got)
	}
	got := strPtr("value")
	if got == nil || *got != "value" {
		t.Errorf("strPtr(\"value\") = %v, want &\"value\"", got)
	}
}

func TestVlanWriteFromUpdate_WithVID(t *testing.T) {
	t.Parallel()

	vid := 42
	in := VLANUpdateInput{VID: &vid, Name: nil}
	out := vlanWriteFromUpdate(in)
	if out.VID != 42 {
		t.Errorf("out.VID = %d, want 42", out.VID)
	}
}

func TestVrfWriteFromUpdate_WithRd(t *testing.T) {
	t.Parallel()

	rd := "65000:1"
	in := VrfUpdateInput{Rd: &rd}
	out := vrfWriteFromUpdate(in)
	if out.Rd != "65000:1" {
		t.Errorf("out.Rd = %q, want %q", out.Rd, "65000:1")
	}
}

//nolint:funlen,maintidx // Table of write handlers hitting the errServiceNotAvailable branch.
func TestWriteHandlers_ServiceNotAvailable(t *testing.T) {
	t.Parallel()

	// Passing a nil service (and no service in context) must hit the
	// errServiceNotAvailable branch in each write handler.
	ctx := context.Background()

	cases := []struct {
		name    string
		handler func() (*mcp.CallToolResult, error)
	}{
		{
			name: "update_site",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdateSiteHandler(nil)(ctx, nil, SiteUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_site",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeleteSiteHandler(nil)(ctx, nil, SiteDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "update_region",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdateRegionHandler(nil)(ctx, nil, RegionUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_region",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeleteRegionHandler(nil)(ctx, nil, RegionDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "update_tenant_group",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdateTenantGroupHandler(nil)(ctx, nil, TenantGroupUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_tenant_group",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeleteTenantGroupHandler(nil)(ctx, nil, TenantGroupDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "update_device",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdateDeviceHandler(nil)(ctx, nil, DeviceUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_device",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeleteDeviceHandler(nil)(ctx, nil, DeviceDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "update_site_group",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdateSiteGroupHandler(nil)(ctx, nil, SiteGroupUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_site_group",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeleteSiteGroupHandler(nil)(ctx, nil, SiteGroupDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "update_rack_group",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdateRackGroupHandler(nil)(ctx, nil, RackGroupUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_rack_group",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeleteRackGroupHandler(nil)(ctx, nil, RackGroupDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "update_rack_role",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdateRackRoleHandler(nil)(ctx, nil, RackRoleUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_rack_role",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeleteRackRoleHandler(nil)(ctx, nil, RackRoleDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "update_device_role",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdateDeviceRoleHandler(nil)(ctx, nil, DeviceRoleUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_device_role",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeleteDeviceRoleHandler(nil)(ctx, nil, DeviceRoleDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "update_platform",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdatePlatformHandler(nil)(ctx, nil, PlatformUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_platform",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeletePlatformHandler(nil)(ctx, nil, PlatformDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "update_virtual_chassis",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdateVirtualChassisHandler(nil)(ctx, nil, VirtualChassisUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_virtual_chassis",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeleteVirtualChassisHandler(nil)(ctx, nil, VirtualChassisDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "update_power_panel",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdatePowerPanelHandler(nil)(ctx, nil, PowerPanelUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_power_panel",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeletePowerPanelHandler(nil)(ctx, nil, PowerPanelDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "update_inventory_item_role",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdateInventoryItemRoleHandler(nil)(ctx, nil, InventoryItemRoleUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_inventory_item_role",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeleteInventoryItemRoleHandler(nil)(ctx, nil, InventoryItemRoleDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "create_front_port",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewCreateFrontPortHandler(nil)(ctx, nil, FrontPortCreateInput{Name: "FP1"})

				return res, err
			},
		},
		{
			name: "update_front_port",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdateFrontPortHandler(nil)(ctx, nil, FrontPortUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_front_port",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeleteFrontPortHandler(nil)(ctx, nil, FrontPortDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "create_rear_port",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewCreateRearPortHandler(nil)(ctx, nil, RearPortCreateInput{Name: "RP1"})

				return res, err
			},
		},
		{
			name: "update_rear_port",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdateRearPortHandler(nil)(ctx, nil, RearPortUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_rear_port",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeleteRearPortHandler(nil)(ctx, nil, RearPortDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "create_console_port",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewCreateConsolePortHandler(nil)(ctx, nil, ConsolePortCreateInput{Name: "CP1"})

				return res, err
			},
		},
		{
			name: "update_console_port",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdateConsolePortHandler(nil)(ctx, nil, ConsolePortUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_console_port",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeleteConsolePortHandler(nil)(ctx, nil, ConsolePortDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "create_console_server_port",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewCreateConsoleServerPortHandler(nil)(ctx, nil, ConsoleServerPortCreateInput{Name: "CSP1"})

				return res, err
			},
		},
		{
			name: "update_console_server_port",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdateConsoleServerPortHandler(nil)(ctx, nil, ConsoleServerPortUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_console_server_port",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeleteConsoleServerPortHandler(nil)(ctx, nil, ConsoleServerPortDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "create_power_port",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewCreatePowerPortHandler(nil)(ctx, nil, PowerPortCreateInput{Name: "PP1"})

				return res, err
			},
		},
		{
			name: "update_power_port",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdatePowerPortHandler(nil)(ctx, nil, PowerPortUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_power_port",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeletePowerPortHandler(nil)(ctx, nil, PowerPortDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "create_power_outlet",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewCreatePowerOutletHandler(nil)(ctx, nil, PowerOutletCreateInput{Name: "PO1"})

				return res, err
			},
		},
		{
			name: "update_power_outlet",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdatePowerOutletHandler(nil)(ctx, nil, PowerOutletUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_power_outlet",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeletePowerOutletHandler(nil)(ctx, nil, PowerOutletDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "create_power_feed",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewCreatePowerFeedHandler(nil)(ctx, nil, PowerFeedCreateInput{Name: "PF1"})

				return res, err
			},
		},
		{
			name: "update_power_feed",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdatePowerFeedHandler(nil)(ctx, nil, PowerFeedUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_power_feed",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeletePowerFeedHandler(nil)(ctx, nil, PowerFeedDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "create_device_bay",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewCreateDeviceBayHandler(nil)(ctx, nil, DeviceBayCreateInput{Name: "Bay 1"})

				return res, err
			},
		},
		{
			name: "update_device_bay",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdateDeviceBayHandler(nil)(ctx, nil, DeviceBayUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_device_bay",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeleteDeviceBayHandler(nil)(ctx, nil, DeviceBayDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "create_rack_reservation",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewCreateRackReservationHandler(nil)(ctx, nil, RackReservationCreateInput{Description: "res"})

				return res, err
			},
		},
		{
			name: "update_rack_reservation",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdateRackReservationHandler(nil)(ctx, nil, RackReservationUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_rack_reservation",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeleteRackReservationHandler(nil)(ctx, nil, RackReservationDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "create_inventory_item",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewCreateInventoryItemHandler(nil)(ctx, nil, InventoryItemCreateInput{Name: "PSU-1"})

				return res, err
			},
		},
		{
			name: "update_inventory_item",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdateInventoryItemHandler(nil)(ctx, nil, InventoryItemUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_inventory_item",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeleteInventoryItemHandler(nil)(ctx, nil, InventoryItemDeleteInput{ID: 1})

				return res, err
			},
		},
		{
			name: "create_mac_address",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewCreateMACAddressHandler(nil)(ctx, nil, MACAddressCreateInput{Address: "00:11:22:33:44:55"})

				return res, err
			},
		},
		{
			name: "update_mac_address",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewUpdateMACAddressHandler(nil)(ctx, nil, MACAddressUpdateInput{ID: 1})

				return res, err
			},
		},
		{
			name: "delete_mac_address",
			handler: func() (*mcp.CallToolResult, error) {
				res, _, err := NewDeleteMACAddressHandler(nil)(ctx, nil, MACAddressDeleteInput{ID: 1})

				return res, err
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.handler()
			if !res.IsError {
				t.Error("res.IsError = false, want true")
			}
			if !errors.Is(err, errServiceNotAvailable) {
				t.Errorf("err = %v, want errServiceNotAvailable", err)
			}
		})
	}
}
