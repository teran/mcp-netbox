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
	if e != err {
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
