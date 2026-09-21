package handlers_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/domain"
	"github.com/teran/mcp-netbox/handlers"
	"github.com/teran/mcp-netbox/mockrepo"
)

func TestCreateVMInterfaceHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVMInterfaceFunc: func(_ context.Context, _ string, in domain.VMInterfaceWrite) (*domain.VMInterface, error) {
				return &domain.VMInterface{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateVMInterfaceHandler(svc)
		vm := 5
		result, out, err := handler(context.Background(), nil, handlers.VMInterfaceCreateInput{Name: "eth0", VirtualMachine: &vm})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "eth0" {
			t.Errorf("out.Data = %+v, want id 1 name eth0", out.Data)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateVMInterfaceHandler(svc)
		vm := 5
		result, _, _ := handler(context.Background(), nil, handlers.VMInterfaceCreateInput{VirtualMachine: &vm})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("virtual_machine required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateVMInterfaceHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VMInterfaceCreateInput{Name: "eth0"})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing virtual_machine")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVMInterfaceFunc: func(_ context.Context, _ string, _ domain.VMInterfaceWrite) (*domain.VMInterface, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateVMInterfaceHandler(svc)
		vm := 5
		result, _, err := handler(context.Background(), nil, handlers.VMInterfaceCreateInput{Name: "eth0", VirtualMachine: &vm})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
		if result.StructuredContent == nil {
			t.Fatal("StructuredContent is nil, want validation_errors")
		}
		sc, ok := result.StructuredContent.(map[string]any)
		if !ok {
			t.Fatalf("StructuredContent type = %T, want map", result.StructuredContent)
		}
		raw, ok := sc["validation_errors"].(json.RawMessage)
		if !ok {
			t.Fatalf("validation_errors type = %T, want json.RawMessage", sc["validation_errors"])
		}
		if string(raw) != string(body) {
			t.Errorf("validation_errors = %s, want %s", raw, body)
		}
	})

	t.Run("generic error", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVMInterfaceFunc: func(_ context.Context, _ string, _ domain.VMInterfaceWrite) (*domain.VMInterface, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateVMInterfaceHandler(svc)
		vm := 5
		result, _, _ := handler(context.Background(), nil, handlers.VMInterfaceCreateInput{Name: "eth0", VirtualMachine: &vm})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateVMInterfaceHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateVMInterfaceFunc: func(_ context.Context, _ string, id int, _ domain.VMInterfaceWrite) (*domain.VMInterface, error) {
				return &domain.VMInterface{ID: id, Name: "eth1"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateVMInterfaceHandler(svc)
		name := "eth1"
		result, out, err := handler(context.Background(), nil, handlers.VMInterfaceUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateVMInterfaceHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VMInterfaceUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"mac_address":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateVMInterfaceFunc: func(_ context.Context, _ string, _ int, _ domain.VMInterfaceWrite) (*domain.VMInterface, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateVMInterfaceHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VMInterfaceUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteVMInterfaceHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteVMInterfaceFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		}, "token")
		handler := handlers.NewDeleteVMInterfaceHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VMInterfaceDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteVMInterfaceHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VMInterfaceDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteVMInterfaceFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteVMInterfaceHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VMInterfaceDeleteInput{ID: 3})
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

func TestCreateVMInterfaceHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateVMInterfaceHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VMInterfaceCreateInput{})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateVMInterfaceHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateVMInterfaceHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VMInterfaceUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteVMInterfaceHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteVMInterfaceHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VMInterfaceDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}
