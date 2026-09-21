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

func TestCreateCircuitTypeHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateCircuitTypeFunc: func(_ context.Context, _ string, in domain.CircuitTypeWrite) (*domain.CircuitType, error) {
				return &domain.CircuitType{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateCircuitTypeHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.CircuitTypeCreateInput{Name: "Fiber"})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "Fiber" {
			t.Errorf("out.Data = %+v, want id 1 name 'Fiber'", out.Data)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateCircuitTypeHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CircuitTypeCreateInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateCircuitTypeFunc: func(_ context.Context, _ string, _ domain.CircuitTypeWrite) (*domain.CircuitType, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateCircuitTypeHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CircuitTypeCreateInput{Name: "Fiber"})
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
			CreateCircuitTypeFunc: func(_ context.Context, _ string, _ domain.CircuitTypeWrite) (*domain.CircuitType, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateCircuitTypeHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CircuitTypeCreateInput{Name: "Fiber"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateCircuitTypeHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateCircuitTypeFunc: func(_ context.Context, _ string, id int, _ domain.CircuitTypeWrite) (*domain.CircuitType, error) {
				return &domain.CircuitType{ID: id, Name: "Fiber"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateCircuitTypeHandler(svc)
		name := "Fiber"
		result, out, err := handler(context.Background(), nil, handlers.CircuitTypeUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateCircuitTypeHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CircuitTypeUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"slug":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateCircuitTypeFunc: func(_ context.Context, _ string, _ int, _ domain.CircuitTypeWrite) (*domain.CircuitType, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateCircuitTypeHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CircuitTypeUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteCircuitTypeHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteCircuitTypeFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		}, "token")
		handler := handlers.NewDeleteCircuitTypeHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CircuitTypeDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteCircuitTypeHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CircuitTypeDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteCircuitTypeFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteCircuitTypeHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CircuitTypeDeleteInput{ID: 3})
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

func TestCreateCircuitTypeHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateCircuitTypeHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.CircuitTypeCreateInput{})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateCircuitTypeHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateCircuitTypeHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.CircuitTypeUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteCircuitTypeHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteCircuitTypeHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.CircuitTypeDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}
