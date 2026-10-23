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
func TestCreateCircuitGroupHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateCircuitGroupFunc: func(
				_ context.Context, _ string, in domain.CircuitGroupWrite,
			) (*domain.CircuitGroup, error) {
				return &domain.CircuitGroup{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateCircuitGroupHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.CircuitGroupCreateInput{Name: "Prod"})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "Prod" {
			t.Errorf("out.Data = %+v, want id 1 name Prod", out.Data)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateCircuitGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CircuitGroupCreateInput{Name: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateCircuitGroupFunc: func(
				_ context.Context, _ string, _ domain.CircuitGroupWrite,
			) (*domain.CircuitGroup, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateCircuitGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CircuitGroupCreateInput{Name: "A"})
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
			CreateCircuitGroupFunc: func(
				_ context.Context, _ string, _ domain.CircuitGroupWrite,
			) (*domain.CircuitGroup, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateCircuitGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CircuitGroupCreateInput{Name: "A"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateCircuitGroupHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateCircuitGroupFunc: func(
				_ context.Context, _ string, id int, _ domain.CircuitGroupWrite,
			) (*domain.CircuitGroup, error) {
				return &domain.CircuitGroup{ID: id, Name: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateCircuitGroupHandler(svc)
		name := "Renamed"
		result, out, err := handler(context.Background(), nil, handlers.CircuitGroupUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateCircuitGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CircuitGroupUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"slug":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateCircuitGroupFunc: func(
				_ context.Context, _ string, _ int, _ domain.CircuitGroupWrite,
			) (*domain.CircuitGroup, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateCircuitGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CircuitGroupUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteCircuitGroupHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteCircuitGroupFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteCircuitGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CircuitGroupDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteCircuitGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CircuitGroupDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteCircuitGroupFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteCircuitGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CircuitGroupDeleteInput{ID: 3})
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

func TestCreateCircuitGroupHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateCircuitGroupHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.CircuitGroupCreateInput{Name: "A"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateCircuitGroupHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateCircuitGroupHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.CircuitGroupUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteCircuitGroupHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteCircuitGroupHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.CircuitGroupDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetCircuitGroupsHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListCircuitGroupsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.CircuitGroup], error) {
				return &domain.PaginatedResponse[domain.CircuitGroup]{
					Count:   1,
					Results: []domain.CircuitGroup{{ID: 1, Name: "Prod"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetCircuitGroupsHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.CircuitGroupsInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].Name != "Prod" {
			t.Errorf("Results = %+v, want single circuit group Prod", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListCircuitGroupsFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.CircuitGroup], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.CircuitGroup]{Count: 0, Results: []domain.CircuitGroup{}}, nil
			},
		}, "token")

		handler := handlers.NewGetCircuitGroupsHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.CircuitGroupsInput{
			Q:        "prod",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "prod" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "prod")
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
			ListCircuitGroupsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.CircuitGroup], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetCircuitGroupsHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CircuitGroupsInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
