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
func TestCreatePowerPanelHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreatePowerPanelFunc: func(_ context.Context, _ string, in domain.PowerPanelWrite) (*domain.PowerPanel, error) {
				return &domain.PowerPanel{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreatePowerPanelHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.PowerPanelCreateInput{Name: "PP-1"})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "PP-1" {
			t.Errorf("out.Data = %+v, want id 1 name PP-1", out.Data)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreatePowerPanelHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.PowerPanelCreateInput{Name: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreatePowerPanelFunc: func(_ context.Context, _ string, _ domain.PowerPanelWrite) (*domain.PowerPanel, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreatePowerPanelHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.PowerPanelCreateInput{Name: "A"})
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
			CreatePowerPanelFunc: func(_ context.Context, _ string, _ domain.PowerPanelWrite) (*domain.PowerPanel, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreatePowerPanelHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.PowerPanelCreateInput{Name: "A"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdatePowerPanelHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdatePowerPanelFunc: func(
				_ context.Context, _ string, id int, _ domain.PowerPanelWrite,
			) (*domain.PowerPanel, error) {
				return &domain.PowerPanel{ID: id, Name: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdatePowerPanelHandler(svc)
		name := "Renamed"
		result, out, err := handler(context.Background(), nil, handlers.PowerPanelUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdatePowerPanelHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.PowerPanelUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"site":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdatePowerPanelFunc: func(
				_ context.Context, _ string, _ int, _ domain.PowerPanelWrite,
			) (*domain.PowerPanel, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdatePowerPanelHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.PowerPanelUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeletePowerPanelHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeletePowerPanelFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeletePowerPanelHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.PowerPanelDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeletePowerPanelHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.PowerPanelDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeletePowerPanelFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeletePowerPanelHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.PowerPanelDeleteInput{ID: 3})
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

func TestCreatePowerPanelHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreatePowerPanelHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.PowerPanelCreateInput{Name: "A"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdatePowerPanelHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdatePowerPanelHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.PowerPanelUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeletePowerPanelHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeletePowerPanelHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.PowerPanelDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetPowerPanelsHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListPowerPanelsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.PowerPanel], error) {
				return &domain.PaginatedResponse[domain.PowerPanel]{
					Count:   1,
					Results: []domain.PowerPanel{{ID: 1, Name: "PP-1"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetPowerPanelsHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.PowerPanelsInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].Name != "PP-1" {
			t.Errorf("Results = %+v, want single power panel PP-1", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListPowerPanelsFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.PowerPanel], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.PowerPanel]{Count: 0, Results: []domain.PowerPanel{}}, nil
			},
		}, "token")

		handler := handlers.NewGetPowerPanelsHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.PowerPanelsInput{
			Q:        "pp",
			Site:     "dc1",
			Location: "row-a",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "pp" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "pp")
		}
		if gotParams["site"] != "dc1" {
			t.Errorf("site param = %q, want %q", gotParams["site"], "dc1")
		}
		if gotParams["location"] != "row-a" {
			t.Errorf("location param = %q, want %q", gotParams["location"], "row-a")
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
			ListPowerPanelsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.PowerPanel], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetPowerPanelsHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.PowerPanelsInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
