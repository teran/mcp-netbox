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

func TestCreateCircuitHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateCircuitFunc: func(_ context.Context, _ string, in domain.CircuitWrite) (*domain.Circuit, error) {
				return &domain.Circuit{ID: 1, CID: in.CID}, nil
			},
		}, "token")

		handler := handlers.NewCreateCircuitHandler(svc)
		prov := 2
		ct := 3
		result, out, err := handler(context.Background(), nil, handlers.CircuitCreateInput{CID: "CIR-001", Provider: &prov, CircuitType: &ct})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.CID != "CIR-001" {
			t.Errorf("out.Data = %+v, want id 1 cid CIR-001", out.Data)
		}
	})

	t.Run("cid required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateCircuitHandler(svc)
		prov := 2
		ct := 3
		result, _, _ := handler(context.Background(), nil, handlers.CircuitCreateInput{CID: "", Provider: &prov, CircuitType: &ct})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing cid")
		}
	})

	t.Run("provider required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateCircuitHandler(svc)
		ct := 3
		result, _, _ := handler(context.Background(), nil, handlers.CircuitCreateInput{CID: "CIR-001", CircuitType: &ct})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing provider")
		}
	})

	t.Run("circuit_type required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateCircuitHandler(svc)
		prov := 2
		result, _, _ := handler(context.Background(), nil, handlers.CircuitCreateInput{CID: "CIR-001", Provider: &prov})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing circuit_type")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"cid":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateCircuitFunc: func(_ context.Context, _ string, _ domain.CircuitWrite) (*domain.Circuit, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateCircuitHandler(svc)
		prov := 2
		ct := 3
		result, _, err := handler(context.Background(), nil, handlers.CircuitCreateInput{CID: "CIR-001", Provider: &prov, CircuitType: &ct})
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
			CreateCircuitFunc: func(_ context.Context, _ string, _ domain.CircuitWrite) (*domain.Circuit, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateCircuitHandler(svc)
		prov := 2
		ct := 3
		result, _, _ := handler(context.Background(), nil, handlers.CircuitCreateInput{CID: "CIR-001", Provider: &prov, CircuitType: &ct})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateCircuitHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateCircuitFunc: func(_ context.Context, _ string, id int, _ domain.CircuitWrite) (*domain.Circuit, error) {
				return &domain.Circuit{ID: id, CID: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateCircuitHandler(svc)
		cid := "Renamed"
		result, out, err := handler(context.Background(), nil, handlers.CircuitUpdateInput{ID: 7, CID: &cid})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateCircuitHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CircuitUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"status":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateCircuitFunc: func(_ context.Context, _ string, _ int, _ domain.CircuitWrite) (*domain.Circuit, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateCircuitHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CircuitUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteCircuitHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteCircuitFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		}, "token")
		handler := handlers.NewDeleteCircuitHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CircuitDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteCircuitHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CircuitDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteCircuitFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteCircuitHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CircuitDeleteInput{ID: 3})
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

func TestCreateCircuitHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateCircuitHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.CircuitCreateInput{CID: "CIR-001"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateCircuitHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateCircuitHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.CircuitUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteCircuitHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteCircuitHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.CircuitDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}
