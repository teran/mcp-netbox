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
func TestCreateConsolePortHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateConsolePortFunc: func(_ context.Context, _ string, in domain.ConsolePortWrite) (*domain.ConsolePort, error) {
				return &domain.ConsolePort{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateConsolePortHandler(svc)
		dev := 5
		result, out, err := handler(context.Background(), nil, handlers.ConsolePortCreateInput{Name: "CP1", Device: &dev})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "CP1" {
			t.Errorf("out.Data = %+v, want id 1 name CP1", out.Data)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateConsolePortHandler(svc)
		dev := 5
		result, _, _ := handler(context.Background(), nil, handlers.ConsolePortCreateInput{Name: "", Device: &dev})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateConsolePortFunc: func(_ context.Context, _ string, _ domain.ConsolePortWrite) (*domain.ConsolePort, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateConsolePortHandler(svc)
		dev := 5
		result, _, err := handler(context.Background(), nil, handlers.ConsolePortCreateInput{Name: "CP1", Device: &dev})
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
			CreateConsolePortFunc: func(_ context.Context, _ string, _ domain.ConsolePortWrite) (*domain.ConsolePort, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateConsolePortHandler(svc)
		dev := 5
		result, _, _ := handler(context.Background(), nil, handlers.ConsolePortCreateInput{Name: "CP1", Device: &dev})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateConsolePortHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateConsolePortFunc: func(
				_ context.Context, _ string, id int, _ domain.ConsolePortWrite,
			) (*domain.ConsolePort, error) {
				return &domain.ConsolePort{ID: id, Name: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateConsolePortHandler(svc)
		name := "Renamed"
		result, out, err := handler(context.Background(), nil, handlers.ConsolePortUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateConsolePortHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ConsolePortUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"type":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateConsolePortFunc: func(
				_ context.Context, _ string, _ int, _ domain.ConsolePortWrite,
			) (*domain.ConsolePort, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateConsolePortHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ConsolePortUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteConsolePortHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteConsolePortFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteConsolePortHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ConsolePortDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteConsolePortHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ConsolePortDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteConsolePortFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteConsolePortHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ConsolePortDeleteInput{ID: 3})
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

func TestCreateConsolePortHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateConsolePortHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ConsolePortCreateInput{Name: "CP1"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateConsolePortHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateConsolePortHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ConsolePortUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteConsolePortHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteConsolePortHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ConsolePortDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetConsolePortsHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListConsolePortsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ConsolePort], error) {
				return &domain.PaginatedResponse[domain.ConsolePort]{
					Count:   1,
					Results: []domain.ConsolePort{{ID: 1, Name: "CP1"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetConsolePortsHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.ConsolePortsInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 || len(out.Results) != 1 || out.Results[0].Name != "CP1" {
			t.Errorf("out = %+v, want count 1 with CP1", out)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListConsolePortsFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.ConsolePort], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.ConsolePort]{Count: 0, Results: []domain.ConsolePort{}}, nil
			},
		}, "token")

		handler := handlers.NewGetConsolePortsHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.ConsolePortsInput{
			Q:        "cp",
			Device:   "sw1",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "cp" || gotParams["device"] != "sw1" || gotParams["tag"] != "prod" {
			t.Errorf("params = %+v, want q cp, device sw1, tag prod", gotParams)
		}
		if gotParams["offset"] != "50" || gotParams["limit"] != "50" {
			t.Errorf("pagination = %+v, want offset 50 limit 50", gotParams)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListConsolePortsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ConsolePort], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetConsolePortsHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ConsolePortsInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
