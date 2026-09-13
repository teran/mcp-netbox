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

func TestCreateInterfaceHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateInterfaceFunc: func(_ context.Context, _ string, in domain.InterfaceWrite) (*domain.Interface, error) {
				return &domain.Interface{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateInterfaceHandler(svc)
		dev := 5
		typ := "1000base-t"
		result, out, err := handler(context.Background(), nil, handlers.InterfaceCreateInput{Name: "eth0", Device: &dev, Type: &typ})
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
		handler := handlers.NewCreateInterfaceHandler(svc)
		dev := 5
		typ := "1000base-t"
		result, _, _ := handler(context.Background(), nil, handlers.InterfaceCreateInput{Name: "", Device: &dev, Type: &typ})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("device required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateInterfaceHandler(svc)
		typ := "1000base-t"
		result, _, _ := handler(context.Background(), nil, handlers.InterfaceCreateInput{Name: "eth0", Type: &typ})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing device")
		}
	})

	t.Run("type required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateInterfaceHandler(svc)
		dev := 5
		result, _, _ := handler(context.Background(), nil, handlers.InterfaceCreateInput{Name: "eth0", Device: &dev})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing type")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateInterfaceFunc: func(_ context.Context, _ string, _ domain.InterfaceWrite) (*domain.Interface, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateInterfaceHandler(svc)
		dev := 5
		typ := "1000base-t"
		result, _, err := handler(context.Background(), nil, handlers.InterfaceCreateInput{Name: "eth0", Device: &dev, Type: &typ})
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
			CreateInterfaceFunc: func(_ context.Context, _ string, _ domain.InterfaceWrite) (*domain.Interface, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateInterfaceHandler(svc)
		dev := 5
		typ := "1000base-t"
		result, _, _ := handler(context.Background(), nil, handlers.InterfaceCreateInput{Name: "eth0", Device: &dev, Type: &typ})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateInterfaceHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateInterfaceFunc: func(_ context.Context, _ string, id int, _ domain.InterfaceWrite) (*domain.Interface, error) {
				return &domain.Interface{ID: id, Name: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateInterfaceHandler(svc)
		name := "Renamed"
		result, out, err := handler(context.Background(), nil, handlers.InterfaceUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateInterfaceHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.InterfaceUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"type":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateInterfaceFunc: func(_ context.Context, _ string, _ int, _ domain.InterfaceWrite) (*domain.Interface, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateInterfaceHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.InterfaceUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteInterfaceHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteInterfaceFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		}, "token")
		handler := handlers.NewDeleteInterfaceHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.InterfaceDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteInterfaceHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.InterfaceDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteInterfaceFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteInterfaceHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.InterfaceDeleteInput{ID: 3})
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

func TestCreateInterfaceHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateInterfaceHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.InterfaceCreateInput{Name: "eth0"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateInterfaceHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateInterfaceHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.InterfaceUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteInterfaceHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteInterfaceHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.InterfaceDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}
