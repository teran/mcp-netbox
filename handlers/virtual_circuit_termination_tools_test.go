package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/domain"
	"github.com/teran/mcp-netbox/handlers"
	"github.com/teran/mcp-netbox/mockrepo"
)

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestCreateVirtualCircuitTerminationHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.VirtualCircuitTerminationWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVirtualCircuitTerminationFunc: func(
				_ context.Context, _ string, in domain.VirtualCircuitTerminationWrite,
			) (*domain.VirtualCircuitTermination, error) {
				gotWrite = in

				return &domain.VirtualCircuitTermination{ID: 1, VirtualCircuit: &domain.Nested{ID: in.VirtualCircuit}}, nil
			},
		}, "token")

		role := "peer"
		iface := 9
		handler := handlers.NewCreateVirtualCircuitTerminationHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.VirtualCircuitTerminationCreateInput{
			VirtualCircuit: 5, Role: &role, Interface: iface,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 {
			t.Errorf("out.Data = %+v, want id 1", out.Data)
		}
		if gotWrite.VirtualCircuit != 5 {
			t.Errorf("create mapped VirtualCircuit = %d, want 5", gotWrite.VirtualCircuit)
		}
		if gotWrite.Role == nil || *gotWrite.Role != "peer" {
			t.Errorf("create mapped Role = %v, want peer", gotWrite.Role)
		}
		if gotWrite.Interface != 9 {
			t.Errorf("create mapped Interface = %d, want 9", gotWrite.Interface)
		}
	})

	t.Run("virtual_circuit required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateVirtualCircuitTerminationHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualCircuitTerminationCreateInput{
			VirtualCircuit: 0,
		})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing virtual_circuit")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"virtual_circuit":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVirtualCircuitTerminationFunc: func(
				_ context.Context, _ string, _ domain.VirtualCircuitTerminationWrite,
			) (*domain.VirtualCircuitTermination, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateVirtualCircuitTerminationHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitTerminationCreateInput{
			VirtualCircuit: 5,
		})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
		structuredContent, ok := result.StructuredContent.(map[string]any)
		if !ok {
			t.Fatalf("StructuredContent type = %T, want map", result.StructuredContent)
		}
		raw, ok := structuredContent["validation_errors"].(json.RawMessage)
		if !ok {
			t.Fatalf("validation_errors type = %T, want json.RawMessage", structuredContent["validation_errors"])
		}
		if string(raw) != string(body) {
			t.Errorf("validation_errors = %s, want %s", raw, body)
		}
	})

	t.Run("generic error", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVirtualCircuitTerminationFunc: func(
				_ context.Context, _ string, _ domain.VirtualCircuitTerminationWrite,
			) (*domain.VirtualCircuitTermination, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateVirtualCircuitTerminationHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualCircuitTerminationCreateInput{
			VirtualCircuit: 5,
		})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

//nolint:funlen // Table of subtests covering update branches (pointer-set + validation); splitting is not sensible.
func TestUpdateVirtualCircuitTerminationHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.VirtualCircuitTerminationWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateVirtualCircuitTerminationFunc: func(
				_ context.Context, _ string, id int, in domain.VirtualCircuitTerminationWrite,
			) (*domain.VirtualCircuitTermination, error) {
				gotWrite = in

				return &domain.VirtualCircuitTermination{ID: id}, nil
			},
		}, "token")
		handler := handlers.NewUpdateVirtualCircuitTerminationHandler(svc)
		role := "peer"
		result, out, err := handler(context.Background(), nil, handlers.VirtualCircuitTerminationUpdateInput{
			ID: 7, Role: &role,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
		if gotWrite.Role == nil || *gotWrite.Role != "peer" {
			t.Errorf("update mapped Role = %v, want peer", gotWrite.Role)
		}
	})

	t.Run("success with virtual_circuit set", func(t *testing.T) {
		var gotWrite domain.VirtualCircuitTerminationWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateVirtualCircuitTerminationFunc: func(
				_ context.Context, _ string, id int, in domain.VirtualCircuitTerminationWrite,
			) (*domain.VirtualCircuitTermination, error) {
				gotWrite = in

				return &domain.VirtualCircuitTermination{ID: id}, nil
			},
		}, "token")
		handler := handlers.NewUpdateVirtualCircuitTerminationHandler(svc)
		vc := 9
		result, out, err := handler(context.Background(), nil, handlers.VirtualCircuitTerminationUpdateInput{
			ID: 7, VirtualCircuit: &vc,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
		if gotWrite.VirtualCircuit != 9 {
			t.Errorf("update mapped VirtualCircuit = %d, want 9", gotWrite.VirtualCircuit)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateVirtualCircuitTerminationHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualCircuitTerminationUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"virtual_circuit":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateVirtualCircuitTerminationFunc: func(
				_ context.Context, _ string, _ int, _ domain.VirtualCircuitTerminationWrite,
			) (*domain.VirtualCircuitTermination, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateVirtualCircuitTerminationHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitTerminationUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteVirtualCircuitTerminationHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteVirtualCircuitTerminationFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteVirtualCircuitTerminationHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitTerminationDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteVirtualCircuitTerminationHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualCircuitTerminationDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteVirtualCircuitTerminationFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteVirtualCircuitTerminationHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitTerminationDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
		if result.StructuredContent == nil {
			t.Error("StructuredContent = nil, want validation_errors")
		}
	})
}

func TestCreateVirtualCircuitTerminationHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateVirtualCircuitTerminationHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitTerminationCreateInput{
		VirtualCircuit: 5,
	})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateVirtualCircuitTerminationHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateVirtualCircuitTerminationHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitTerminationUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteVirtualCircuitTerminationHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteVirtualCircuitTerminationHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitTerminationDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetVirtualCircuitTerminationsHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListVirtualCircuitTerminationsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualCircuitTermination], error) {
				return &domain.PaginatedResponse[domain.VirtualCircuitTermination]{
					Count:   1,
					Results: []domain.VirtualCircuitTermination{{ID: 1}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetVirtualCircuitTerminationsHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.VirtualCircuitTerminationsInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].ID != 1 {
			t.Errorf("Results = %+v, want single virtual circuit termination", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListVirtualCircuitTerminationsFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualCircuitTermination], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.VirtualCircuitTermination]{
					Count:   0,
					Results: []domain.VirtualCircuitTermination{},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetVirtualCircuitTerminationsHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.VirtualCircuitTerminationsInput{
			Q:              "vc",
			VirtualCircuit: "5",
			InterfaceType:  "dcim.interface",
			InterfaceID:    "9",
			Page:           2,
			PageSize:       50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "vc" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "vc")
		}
		if gotParams["virtual_circuit"] != "5" {
			t.Errorf("virtual_circuit param = %q, want %q", gotParams["virtual_circuit"], "5")
		}
		if gotParams["interface_type"] != "dcim.interface" {
			t.Errorf("interface_type param = %q, want %q", gotParams["interface_type"], "dcim.interface")
		}
		if gotParams["interface_id"] != "9" {
			t.Errorf("interface_id param = %q, want %q", gotParams["interface_id"], "9")
		}
		if gotParams["offset"] != "50" || gotParams["limit"] != "50" {
			t.Errorf("offset/limit = %q/%q, want 50/50", gotParams["offset"], gotParams["limit"])
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListVirtualCircuitTerminationsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualCircuitTermination], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetVirtualCircuitTerminationsHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualCircuitTerminationsInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
