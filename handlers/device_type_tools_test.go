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

func TestCreateDeviceTypeHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateDeviceTypeFunc: func(_ context.Context, _ string, in domain.DeviceTypeWrite) (*domain.DeviceType, error) {
				return &domain.DeviceType{ID: 1, Model: in.Model}, nil
			},
		}, "token")

		handler := handlers.NewCreateDeviceTypeHandler(svc)
		mf := 3
		result, out, err := handler(context.Background(), nil, handlers.DeviceTypeCreateInput{Manufacturer: &mf, Model: "C9300"})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Model != "C9300" {
			t.Errorf("out.Data = %+v, want id 1 model C9300", out.Data)
		}
	})

	t.Run("manufacturer required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateDeviceTypeHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.DeviceTypeCreateInput{Model: "C9300"})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing manufacturer")
		}
	})

	t.Run("model required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateDeviceTypeHandler(svc)
		mf := 3
		result, _, _ := handler(context.Background(), nil, handlers.DeviceTypeCreateInput{Manufacturer: &mf})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing model")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"model":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateDeviceTypeFunc: func(_ context.Context, _ string, _ domain.DeviceTypeWrite) (*domain.DeviceType, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateDeviceTypeHandler(svc)
		mf := 3
		result, _, err := handler(context.Background(), nil, handlers.DeviceTypeCreateInput{Manufacturer: &mf, Model: "C9300"})
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
			CreateDeviceTypeFunc: func(_ context.Context, _ string, _ domain.DeviceTypeWrite) (*domain.DeviceType, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateDeviceTypeHandler(svc)
		mf := 3
		result, _, _ := handler(context.Background(), nil, handlers.DeviceTypeCreateInput{Manufacturer: &mf, Model: "C9300"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateDeviceTypeHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateDeviceTypeFunc: func(_ context.Context, _ string, id int, _ domain.DeviceTypeWrite) (*domain.DeviceType, error) {
				return &domain.DeviceType{ID: id, Model: "C9300-2"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateDeviceTypeHandler(svc)
		name := "C9300-2"
		result, out, err := handler(context.Background(), nil, handlers.DeviceTypeUpdateInput{ID: 7, Model: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateDeviceTypeHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.DeviceTypeUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"slug":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateDeviceTypeFunc: func(_ context.Context, _ string, _ int, _ domain.DeviceTypeWrite) (*domain.DeviceType, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateDeviceTypeHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.DeviceTypeUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteDeviceTypeHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteDeviceTypeFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		}, "token")
		handler := handlers.NewDeleteDeviceTypeHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.DeviceTypeDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteDeviceTypeHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.DeviceTypeDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteDeviceTypeFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteDeviceTypeHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.DeviceTypeDeleteInput{ID: 3})
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

func TestCreateDeviceTypeHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateDeviceTypeHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.DeviceTypeCreateInput{})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateDeviceTypeHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateDeviceTypeHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.DeviceTypeUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteDeviceTypeHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteDeviceTypeHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.DeviceTypeDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}
