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

func TestCreateLocationHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateLocationFunc: func(_ context.Context, _ string, in domain.LocationWrite) (*domain.Location, error) {
				return &domain.Location{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateLocationHandler(svc)
		site := 5
		result, out, err := handler(context.Background(), nil, handlers.LocationCreateInput{Name: "Row A", Site: &site})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "Row A" {
			t.Errorf("out.Data = %+v, want id 1 name 'Row A'", out.Data)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateLocationHandler(svc)
		site := 5
		result, _, _ := handler(context.Background(), nil, handlers.LocationCreateInput{Site: &site})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("site required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateLocationHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.LocationCreateInput{Name: "Row A"})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing site")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"site":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateLocationFunc: func(_ context.Context, _ string, _ domain.LocationWrite) (*domain.Location, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateLocationHandler(svc)
		site := 5
		result, _, err := handler(context.Background(), nil, handlers.LocationCreateInput{Name: "Row A", Site: &site})
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
			CreateLocationFunc: func(_ context.Context, _ string, _ domain.LocationWrite) (*domain.Location, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateLocationHandler(svc)
		site := 5
		result, _, _ := handler(context.Background(), nil, handlers.LocationCreateInput{Name: "Row A", Site: &site})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateLocationHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateLocationFunc: func(_ context.Context, _ string, id int, _ domain.LocationWrite) (*domain.Location, error) {
				return &domain.Location{ID: id, Name: "Row B"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateLocationHandler(svc)
		name := "Row B"
		result, out, err := handler(context.Background(), nil, handlers.LocationUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateLocationHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.LocationUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"slug":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateLocationFunc: func(_ context.Context, _ string, _ int, _ domain.LocationWrite) (*domain.Location, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateLocationHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.LocationUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteLocationHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteLocationFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		}, "token")
		handler := handlers.NewDeleteLocationHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.LocationDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteLocationHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.LocationDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteLocationFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteLocationHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.LocationDeleteInput{ID: 3})
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

func TestCreateLocationHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateLocationHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.LocationCreateInput{})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateLocationHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateLocationHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.LocationUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteLocationHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteLocationHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.LocationDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}
