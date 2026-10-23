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
func TestCreateVirtualCircuitTypeHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVirtualCircuitTypeFunc: func(
				_ context.Context, _ string, in domain.VirtualCircuitTypeWrite,
			) (*domain.VirtualCircuitType, error) {
				return &domain.VirtualCircuitType{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateVirtualCircuitTypeHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.VirtualCircuitTypeCreateInput{Name: "VPLS"})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "VPLS" {
			t.Errorf("out.Data = %+v, want id 1 name VPLS", out.Data)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateVirtualCircuitTypeHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualCircuitTypeCreateInput{Name: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVirtualCircuitTypeFunc: func(
				_ context.Context, _ string, _ domain.VirtualCircuitTypeWrite,
			) (*domain.VirtualCircuitType, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateVirtualCircuitTypeHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitTypeCreateInput{Name: "A"})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
		if result.StructuredContent == nil {
			t.Fatal("StructuredContent is nil, want validation_errors")
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
			CreateVirtualCircuitTypeFunc: func(
				_ context.Context, _ string, _ domain.VirtualCircuitTypeWrite,
			) (*domain.VirtualCircuitType, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateVirtualCircuitTypeHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualCircuitTypeCreateInput{Name: "A"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateVirtualCircuitTypeHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateVirtualCircuitTypeFunc: func(
				_ context.Context, _ string, id int, _ domain.VirtualCircuitTypeWrite,
			) (*domain.VirtualCircuitType, error) {
				return &domain.VirtualCircuitType{ID: id, Name: "VPLS"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateVirtualCircuitTypeHandler(svc)
		name := "VPLS"
		result, out, err := handler(context.Background(), nil, handlers.VirtualCircuitTypeUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateVirtualCircuitTypeHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualCircuitTypeUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"slug":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateVirtualCircuitTypeFunc: func(
				_ context.Context, _ string, _ int, _ domain.VirtualCircuitTypeWrite,
			) (*domain.VirtualCircuitType, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateVirtualCircuitTypeHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitTypeUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteVirtualCircuitTypeHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteVirtualCircuitTypeFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteVirtualCircuitTypeHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitTypeDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteVirtualCircuitTypeHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualCircuitTypeDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteVirtualCircuitTypeFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteVirtualCircuitTypeHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitTypeDeleteInput{ID: 3})
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

func TestCreateVirtualCircuitTypeHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateVirtualCircuitTypeHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitTypeCreateInput{Name: "A"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateVirtualCircuitTypeHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateVirtualCircuitTypeHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitTypeUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteVirtualCircuitTypeHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteVirtualCircuitTypeHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitTypeDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetVirtualCircuitTypesHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListVirtualCircuitTypesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualCircuitType], error) {
				return &domain.PaginatedResponse[domain.VirtualCircuitType]{
					Count:   1,
					Results: []domain.VirtualCircuitType{{ID: 1, Name: "VPLS"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetVirtualCircuitTypesHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.VirtualCircuitTypesInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].Name != "VPLS" {
			t.Errorf("Results = %+v, want single virtual circuit type VPLS", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListVirtualCircuitTypesFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualCircuitType], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.VirtualCircuitType]{Count: 0, Results: []domain.VirtualCircuitType{}}, nil
			},
		}, "token")

		handler := handlers.NewGetVirtualCircuitTypesHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.VirtualCircuitTypesInput{
			Q:        "vpls",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "vpls" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "vpls")
		}
		if gotParams["tag"] != "prod" {
			t.Errorf("tag param = %q, want %q", gotParams["tag"], "prod")
		}
		if gotParams["offset"] != "50" || gotParams["limit"] != "50" {
			t.Errorf("offset/limit = %q/%q, want 50/50", gotParams["offset"], gotParams["limit"])
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListVirtualCircuitTypesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualCircuitType], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetVirtualCircuitTypesHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualCircuitTypesInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
