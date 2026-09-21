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

func TestCreateCableHandler(t *testing.T) {
	t.Parallel()

	ta := &handlers.CableTerminationInput{ObjectType: "dcim.interface", ObjectID: 101}
	tb := &handlers.CableTerminationInput{ObjectType: "dcim.interface", ObjectID: 102}

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateCableFunc: func(_ context.Context, _ string, in domain.CableWrite) (*domain.Cable, error) {
				return &domain.Cable{ID: 1, Label: *in.Label}, nil
			},
		}, "token")

		handler := handlers.NewCreateCableHandler(svc)
		label := "link-01"
		result, out, err := handler(context.Background(), nil, handlers.CableCreateInput{TerminationA: ta, TerminationB: tb, Label: &label})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Label != "link-01" {
			t.Errorf("out.Data = %+v, want id 1 label link-01", out.Data)
		}
	})

	t.Run("termination_a required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateCableHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CableCreateInput{TerminationB: tb})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing termination_a")
		}
	})

	t.Run("termination_b required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateCableHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CableCreateInput{TerminationA: ta})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing termination_b")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"termination_a":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateCableFunc: func(_ context.Context, _ string, _ domain.CableWrite) (*domain.Cable, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateCableHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CableCreateInput{TerminationA: ta, TerminationB: tb})
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
			CreateCableFunc: func(_ context.Context, _ string, _ domain.CableWrite) (*domain.Cable, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateCableHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CableCreateInput{TerminationA: ta, TerminationB: tb})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateCableHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateCableFunc: func(_ context.Context, _ string, id int, _ domain.CableWrite) (*domain.Cable, error) {
				return &domain.Cable{ID: id, Label: "link-02"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateCableHandler(svc)
		label := "link-02"
		result, out, err := handler(context.Background(), nil, handlers.CableUpdateInput{ID: 7, Label: &label})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateCableHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CableUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"status":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateCableFunc: func(_ context.Context, _ string, _ int, _ domain.CableWrite) (*domain.Cable, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateCableHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CableUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteCableHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteCableFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		}, "token")
		handler := handlers.NewDeleteCableHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CableDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteCableHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CableDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteCableFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteCableHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CableDeleteInput{ID: 3})
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

func TestCreateCableHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateCableHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.CableCreateInput{})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateCableHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateCableHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.CableUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteCableHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteCableHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.CableDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}
