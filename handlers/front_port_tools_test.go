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
func TestCreateFrontPortHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateFrontPortFunc: func(_ context.Context, _ string, in domain.FrontPortWrite) (*domain.FrontPort, error) {
				return &domain.FrontPort{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateFrontPortHandler(svc)
		dev := 5
		typ := "lc"
		result, out, err := handler(context.Background(), nil, handlers.FrontPortCreateInput{
			Name: "FP1", Device: &dev, Type: &typ,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "FP1" {
			t.Errorf("out.Data = %+v, want id 1 name FP1", out.Data)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateFrontPortHandler(svc)
		dev := 5
		typ := "lc"
		result, _, _ := handler(context.Background(), nil, handlers.FrontPortCreateInput{Name: "", Device: &dev, Type: &typ})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("device required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateFrontPortHandler(svc)
		typ := "lc"
		result, _, _ := handler(context.Background(), nil, handlers.FrontPortCreateInput{Name: "FP1", Type: &typ})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing device")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateFrontPortFunc: func(_ context.Context, _ string, _ domain.FrontPortWrite) (*domain.FrontPort, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateFrontPortHandler(svc)
		dev := 5
		typ := "lc"
		result, _, err := handler(context.Background(), nil, handlers.FrontPortCreateInput{
			Name: "FP1", Device: &dev, Type: &typ,
		})
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
			CreateFrontPortFunc: func(_ context.Context, _ string, _ domain.FrontPortWrite) (*domain.FrontPort, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateFrontPortHandler(svc)
		dev := 5
		typ := "lc"
		result, _, _ := handler(context.Background(), nil, handlers.FrontPortCreateInput{
			Name: "FP1", Device: &dev, Type: &typ,
		})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateFrontPortHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateFrontPortFunc: func(_ context.Context, _ string, id int, _ domain.FrontPortWrite) (*domain.FrontPort, error) {
				return &domain.FrontPort{ID: id, Name: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateFrontPortHandler(svc)
		name := "Renamed"
		result, out, err := handler(context.Background(), nil, handlers.FrontPortUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateFrontPortHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.FrontPortUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"rear_port":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateFrontPortFunc: func(_ context.Context, _ string, _ int, _ domain.FrontPortWrite) (*domain.FrontPort, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateFrontPortHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.FrontPortUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteFrontPortHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteFrontPortFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteFrontPortHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.FrontPortDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteFrontPortHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.FrontPortDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteFrontPortFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteFrontPortHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.FrontPortDeleteInput{ID: 3})
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

func TestCreateFrontPortHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateFrontPortHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.FrontPortCreateInput{Name: "FP1"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateFrontPortHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateFrontPortHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.FrontPortUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteFrontPortHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteFrontPortHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.FrontPortDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetFrontPortsHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListFrontPortsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.FrontPort], error) {
				return &domain.PaginatedResponse[domain.FrontPort]{
					Count:   1,
					Results: []domain.FrontPort{{ID: 1, Name: "FP1"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetFrontPortsHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.FrontPortsInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].Name != "FP1" {
			t.Errorf("Results = %+v, want single front port FP1", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListFrontPortsFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.FrontPort], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.FrontPort]{Count: 0, Results: []domain.FrontPort{}}, nil
			},
		}, "token")

		handler := handlers.NewGetFrontPortsHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.FrontPortsInput{
			Q:        "fp",
			Device:   "sw1",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "fp" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "fp")
		}
		if gotParams["device"] != "sw1" {
			t.Errorf("device param = %q, want %q", gotParams["device"], "sw1")
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
			ListFrontPortsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.FrontPort], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetFrontPortsHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.FrontPortsInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
