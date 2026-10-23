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
func TestCreateVirtualChassisHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVirtualChassisFunc: func(
				_ context.Context, _ string, in domain.VirtualChassisWrite,
			) (*domain.VirtualChassis, error) {
				return &domain.VirtualChassis{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateVirtualChassisHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.VirtualChassisCreateInput{Name: "VC-1"})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "VC-1" {
			t.Errorf("out.Data = %+v, want id 1 name VC-1", out.Data)
		}
	})

	t.Run("name optional", func(t *testing.T) {
		// Name is optional on create for virtual chassis; an empty input must
		// still be accepted.
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVirtualChassisFunc: func(
				_ context.Context, _ string, _ domain.VirtualChassisWrite,
			) (*domain.VirtualChassis, error) {
				return &domain.VirtualChassis{ID: 1}, nil
			},
		}, "token")
		handler := handlers.NewCreateVirtualChassisHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VirtualChassisCreateInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false for empty name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVirtualChassisFunc: func(
				_ context.Context, _ string, _ domain.VirtualChassisWrite,
			) (*domain.VirtualChassis, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateVirtualChassisHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VirtualChassisCreateInput{Name: "A"})
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
			CreateVirtualChassisFunc: func(
				_ context.Context, _ string, _ domain.VirtualChassisWrite,
			) (*domain.VirtualChassis, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateVirtualChassisHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualChassisCreateInput{Name: "A"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateVirtualChassisHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateVirtualChassisFunc: func(
				_ context.Context, _ string, id int, _ domain.VirtualChassisWrite,
			) (*domain.VirtualChassis, error) {
				return &domain.VirtualChassis{ID: id, Name: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateVirtualChassisHandler(svc)
		name := "Renamed"
		result, out, err := handler(context.Background(), nil, handlers.VirtualChassisUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateVirtualChassisHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualChassisUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"master":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateVirtualChassisFunc: func(
				_ context.Context, _ string, _ int, _ domain.VirtualChassisWrite,
			) (*domain.VirtualChassis, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateVirtualChassisHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VirtualChassisUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteVirtualChassisHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteVirtualChassisFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteVirtualChassisHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VirtualChassisDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteVirtualChassisHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualChassisDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteVirtualChassisFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteVirtualChassisHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VirtualChassisDeleteInput{ID: 3})
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

func TestCreateVirtualChassisHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateVirtualChassisHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VirtualChassisCreateInput{Name: "A"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateVirtualChassisHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateVirtualChassisHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VirtualChassisUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteVirtualChassisHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteVirtualChassisHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VirtualChassisDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetVirtualChassisHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListVirtualChassisFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualChassis], error) {
				return &domain.PaginatedResponse[domain.VirtualChassis]{
					Count:   1,
					Results: []domain.VirtualChassis{{ID: 1, Name: "VC-1"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetVirtualChassisHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.VirtualChassisInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].Name != "VC-1" {
			t.Errorf("Results = %+v, want single virtual chassis VC-1", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListVirtualChassisFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualChassis], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.VirtualChassis]{Count: 0, Results: []domain.VirtualChassis{}}, nil
			},
		}, "token")

		handler := handlers.NewGetVirtualChassisHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.VirtualChassisInput{
			Q:        "vc",
			Name:     "VC-1",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "vc" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "vc")
		}
		if gotParams["name"] != "VC-1" {
			t.Errorf("name param = %q, want %q", gotParams["name"], "VC-1")
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
			ListVirtualChassisFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualChassis], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetVirtualChassisHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualChassisInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
