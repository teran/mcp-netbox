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

func TestCreateCircuitTerminationHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateCircuitTerminationFunc: func(_ context.Context, _ string, in domain.CircuitTerminationWrite) (*domain.CircuitTermination, error) {
				return &domain.CircuitTermination{ID: 1, TermSide: in.TermSide}, nil
			},
		}, "token")

		handler := handlers.NewCreateCircuitTerminationHandler(svc)
		circ := 2
		site := 3
		result, out, err := handler(context.Background(), nil, handlers.CircuitTerminationCreateInput{TermSide: "A", Circuit: &circ, Site: &site})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.TermSide != "A" {
			t.Errorf("out.Data = %+v, want id 1 term_side A", out.Data)
		}
	})

	t.Run("term_side required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateCircuitTerminationHandler(svc)
		circ := 2
		site := 3
		result, _, _ := handler(context.Background(), nil, handlers.CircuitTerminationCreateInput{TermSide: "", Circuit: &circ, Site: &site})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing term_side")
		}
	})

	t.Run("circuit required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateCircuitTerminationHandler(svc)
		site := 3
		result, _, _ := handler(context.Background(), nil, handlers.CircuitTerminationCreateInput{TermSide: "A", Site: &site})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing circuit")
		}
	})

	t.Run("site required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateCircuitTerminationHandler(svc)
		circ := 2
		result, _, _ := handler(context.Background(), nil, handlers.CircuitTerminationCreateInput{TermSide: "A", Circuit: &circ})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing site")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"term_side":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateCircuitTerminationFunc: func(_ context.Context, _ string, _ domain.CircuitTerminationWrite) (*domain.CircuitTermination, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateCircuitTerminationHandler(svc)
		circ := 2
		site := 3
		result, _, err := handler(context.Background(), nil, handlers.CircuitTerminationCreateInput{TermSide: "A", Circuit: &circ, Site: &site})
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
			CreateCircuitTerminationFunc: func(_ context.Context, _ string, _ domain.CircuitTerminationWrite) (*domain.CircuitTermination, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateCircuitTerminationHandler(svc)
		circ := 2
		site := 3
		result, _, _ := handler(context.Background(), nil, handlers.CircuitTerminationCreateInput{TermSide: "A", Circuit: &circ, Site: &site})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateCircuitTerminationHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateCircuitTerminationFunc: func(_ context.Context, _ string, id int, _ domain.CircuitTerminationWrite) (*domain.CircuitTermination, error) {
				return &domain.CircuitTermination{ID: id, TermSide: "Z"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateCircuitTerminationHandler(svc)
		side := "Z"
		result, out, err := handler(context.Background(), nil, handlers.CircuitTerminationUpdateInput{ID: 7, TermSide: &side})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateCircuitTerminationHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CircuitTerminationUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"site":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateCircuitTerminationFunc: func(_ context.Context, _ string, _ int, _ domain.CircuitTerminationWrite) (*domain.CircuitTermination, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateCircuitTerminationHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CircuitTerminationUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteCircuitTerminationHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteCircuitTerminationFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		}, "token")
		handler := handlers.NewDeleteCircuitTerminationHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CircuitTerminationDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteCircuitTerminationHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CircuitTerminationDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteCircuitTerminationFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteCircuitTerminationHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CircuitTerminationDeleteInput{ID: 3})
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

func TestCreateCircuitTerminationHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateCircuitTerminationHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.CircuitTerminationCreateInput{TermSide: "A"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateCircuitTerminationHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateCircuitTerminationHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.CircuitTerminationUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteCircuitTerminationHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteCircuitTerminationHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.CircuitTerminationDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}
