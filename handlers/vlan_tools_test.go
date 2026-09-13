package handlers_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/domain"
	"github.com/teran/mcp-netbox/handlers"
	"github.com/teran/mcp-netbox/internal/mockrepo"
)

func TestCreateVLANHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVLANFunc: func(_ context.Context, _ string, in domain.VLANWrite) (*domain.VLAN, error) {
				return &domain.VLAN{ID: 1, VID: in.VID, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateVLANHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.VLANCreateInput{VID: 100, Name: "mgmt"})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.VID != 100 || out.Data.Name != "mgmt" {
			t.Errorf("out.Data = %+v, want id 1 vid 100 name mgmt", out.Data)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateVLANHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VLANCreateInput{VID: 100, Name: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("vid required positive", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateVLANHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VLANCreateInput{VID: 0, Name: "mgmt"})
		if !result.IsError {
			t.Error("result.IsError = false, want true for non-positive vid")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVLANFunc: func(_ context.Context, _ string, _ domain.VLANWrite) (*domain.VLAN, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateVLANHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VLANCreateInput{VID: 100, Name: "mgmt"})
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
			CreateVLANFunc: func(_ context.Context, _ string, _ domain.VLANWrite) (*domain.VLAN, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateVLANHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VLANCreateInput{VID: 100, Name: "mgmt"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateVLANHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateVLANFunc: func(_ context.Context, _ string, id int, _ domain.VLANWrite) (*domain.VLAN, error) {
				return &domain.VLAN{ID: id, VID: 200, Name: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateVLANHandler(svc)
		name := "Renamed"
		result, out, err := handler(context.Background(), nil, handlers.VLANUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateVLANHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VLANUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"status":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateVLANFunc: func(_ context.Context, _ string, _ int, _ domain.VLANWrite) (*domain.VLAN, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateVLANHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VLANUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteVLANHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteVLANFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		}, "token")
		handler := handlers.NewDeleteVLANHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VLANDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteVLANHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VLANDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteVLANFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteVLANHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VLANDeleteInput{ID: 3})
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

func TestCreateVLANHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateVLANHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VLANCreateInput{VID: 100, Name: "mgmt"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateVLANHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateVLANHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VLANUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteVLANHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteVLANHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VLANDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}
