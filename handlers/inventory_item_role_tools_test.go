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
func TestCreateInventoryItemRoleHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateInventoryItemRoleFunc: func(
				_ context.Context, _ string, in domain.InventoryItemRoleWrite,
			) (*domain.InventoryItemRole, error) {
				return &domain.InventoryItemRole{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateInventoryItemRoleHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.InventoryItemRoleCreateInput{Name: "Role A"})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "Role A" {
			t.Errorf("out.Data = %+v, want id 1 name Role A", out.Data)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateInventoryItemRoleHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.InventoryItemRoleCreateInput{Name: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateInventoryItemRoleFunc: func(
				_ context.Context, _ string, _ domain.InventoryItemRoleWrite,
			) (*domain.InventoryItemRole, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateInventoryItemRoleHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.InventoryItemRoleCreateInput{Name: "A"})
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
			CreateInventoryItemRoleFunc: func(
				_ context.Context, _ string, _ domain.InventoryItemRoleWrite,
			) (*domain.InventoryItemRole, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateInventoryItemRoleHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.InventoryItemRoleCreateInput{Name: "A"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateInventoryItemRoleHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateInventoryItemRoleFunc: func(
				_ context.Context, _ string, id int, _ domain.InventoryItemRoleWrite,
			) (*domain.InventoryItemRole, error) {
				return &domain.InventoryItemRole{ID: id, Name: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateInventoryItemRoleHandler(svc)
		name := "Renamed"
		result, out, err := handler(context.Background(), nil, handlers.InventoryItemRoleUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateInventoryItemRoleHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.InventoryItemRoleUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"color":["Enter a valid color."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateInventoryItemRoleFunc: func(
				_ context.Context, _ string, _ int, _ domain.InventoryItemRoleWrite,
			) (*domain.InventoryItemRole, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateInventoryItemRoleHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.InventoryItemRoleUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteInventoryItemRoleHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteInventoryItemRoleFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteInventoryItemRoleHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.InventoryItemRoleDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteInventoryItemRoleHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.InventoryItemRoleDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteInventoryItemRoleFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteInventoryItemRoleHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.InventoryItemRoleDeleteInput{ID: 3})
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

func TestCreateInventoryItemRoleHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateInventoryItemRoleHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.InventoryItemRoleCreateInput{Name: "A"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateInventoryItemRoleHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateInventoryItemRoleHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.InventoryItemRoleUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteInventoryItemRoleHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteInventoryItemRoleHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.InventoryItemRoleDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetInventoryItemRolesHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListInventoryItemRolesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.InventoryItemRole], error) {
				return &domain.PaginatedResponse[domain.InventoryItemRole]{
					Count:   1,
					Results: []domain.InventoryItemRole{{ID: 1, Name: "Power"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetInventoryItemRolesHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.InventoryItemRolesInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].Name != "Power" {
			t.Errorf("Results = %+v, want single inventory item role Power", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListInventoryItemRolesFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.InventoryItemRole], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.InventoryItemRole]{Count: 0, Results: []domain.InventoryItemRole{}}, nil
			},
		}, "token")

		handler := handlers.NewGetInventoryItemRolesHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.InventoryItemRolesInput{
			Q:        "power",
			Color:    "00ff00",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "power" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "power")
		}
		if gotParams["color"] != "00ff00" {
			t.Errorf("color param = %q, want %q", gotParams["color"], "00ff00")
		}
		if gotParams["tag"] != "prod" {
			t.Errorf("tag param = %q, want %q", gotParams["tag"], "prod")
		}
		if gotParams["offset"] != "50" {
			t.Errorf("offset param = %q, want %q", gotParams["offset"], "50")
		}
		if gotParams["limit"] != "50" {
			t.Errorf("limit param = %q, want %q", gotParams["limit"], "50")
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListInventoryItemRolesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.InventoryItemRole], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetInventoryItemRolesHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.InventoryItemRolesInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
