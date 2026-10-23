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
func TestCreateTenantGroupHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateTenantGroupFunc: func(_ context.Context, _ string, in domain.TenantGroupWrite) (*domain.TenantGroup, error) {
				return &domain.TenantGroup{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateTenantGroupHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.TenantGroupCreateInput{Name: "Group A"})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "Group A" {
			t.Errorf("out.Data = %+v, want id 1 name Group A", out.Data)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateTenantGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.TenantGroupCreateInput{Name: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateTenantGroupFunc: func(_ context.Context, _ string, _ domain.TenantGroupWrite) (*domain.TenantGroup, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateTenantGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.TenantGroupCreateInput{Name: "A"})
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
			CreateTenantGroupFunc: func(_ context.Context, _ string, _ domain.TenantGroupWrite) (*domain.TenantGroup, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateTenantGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.TenantGroupCreateInput{Name: "A"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateTenantGroupHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateTenantGroupFunc: func(
				_ context.Context, _ string, id int, _ domain.TenantGroupWrite,
			) (*domain.TenantGroup, error) {
				return &domain.TenantGroup{ID: id, Name: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateTenantGroupHandler(svc)
		name := "Renamed"
		result, out, err := handler(context.Background(), nil, handlers.TenantGroupUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateTenantGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.TenantGroupUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"slug":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateTenantGroupFunc: func(
				_ context.Context, _ string, _ int, _ domain.TenantGroupWrite,
			) (*domain.TenantGroup, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateTenantGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.TenantGroupUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteTenantGroupHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteTenantGroupFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteTenantGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.TenantGroupDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteTenantGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.TenantGroupDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteTenantGroupFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteTenantGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.TenantGroupDeleteInput{ID: 3})
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

func TestCreateTenantGroupHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateTenantGroupHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.TenantGroupCreateInput{Name: "A"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateTenantGroupHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateTenantGroupHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.TenantGroupUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteTenantGroupHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteTenantGroupHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.TenantGroupDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetTenantGroupsHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListTenantGroupsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.TenantGroup], error) {
				return &domain.PaginatedResponse[domain.TenantGroup]{
					Count:   1,
					Results: []domain.TenantGroup{{ID: 1, Name: "EU-West"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetTenantGroupsHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.TenantGroupsInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].Name != "EU-West" {
			t.Errorf("Results = %+v, want single tenant group EU-West", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListTenantGroupsFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.TenantGroup], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.TenantGroup]{Count: 0, Results: []domain.TenantGroup{}}, nil
			},
		}, "token")

		handler := handlers.NewGetTenantGroupsHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.TenantGroupsInput{
			Q:        "west",
			Parent:   "7",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "west" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "west")
		}
		if gotParams["parent"] != "7" {
			t.Errorf("parent param = %q, want %q", gotParams["parent"], "7")
		}
		if gotParams["tag"] != "prod" {
			t.Errorf("tag param = %q, want %q", gotParams["tag"], "prod")
		}
		// page/page_size are translated into NetBox offset/limit pagination.
		// page=2, page_size=50 => offset=50, limit=50.
		if gotParams["offset"] != "50" {
			t.Errorf("offset param = %q, want %q", gotParams["offset"], "50")
		}
		if gotParams["limit"] != "50" {
			t.Errorf("limit param = %q, want %q", gotParams["limit"], "50")
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListTenantGroupsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.TenantGroup], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetTenantGroupsHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.TenantGroupsInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

//nolint:funlen // Table of subtests covering the tenant group mapping; splitting is not sensible.
func TestTenantGroupMapsTenantGroup(t *testing.T) {
	t.Parallel()

	t.Run("create maps group", func(t *testing.T) {
		var gotWrite domain.TenantWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateTenantFunc: func(_ context.Context, _ string, in domain.TenantWrite) (*domain.Tenant, error) {
				gotWrite = in

				return &domain.Tenant{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		group := 42
		handler := handlers.NewCreateTenantHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.TenantCreateInput{Name: "ACME", Group: &group})
		if err != nil || result.IsError {
			t.Fatalf("handler failed: result.IsError=%v err=%v", result.IsError, err)
		}
		if gotWrite.Group == nil || *gotWrite.Group != 42 {
			t.Errorf("create mapped Group = %v, want 42", gotWrite.Group)
		}
	})

	t.Run("update maps group", func(t *testing.T) {
		var gotWrite domain.TenantWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateTenantFunc: func(_ context.Context, _ string, _ int, in domain.TenantWrite) (*domain.Tenant, error) {
				gotWrite = in

				return &domain.Tenant{ID: 1, Name: "ACME"}, nil
			},
		}, "token")

		group := 7
		handler := handlers.NewUpdateTenantHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.TenantUpdateInput{ID: 1, Group: &group})
		if err != nil || result.IsError {
			t.Fatalf("handler failed: result.IsError=%v err=%v", result.IsError, err)
		}
		if gotWrite.Group == nil || *gotWrite.Group != 7 {
			t.Errorf("update mapped Group = %v, want 7", gotWrite.Group)
		}
	})

	t.Run("update without group leaves it unset", func(t *testing.T) {
		var gotWrite domain.TenantWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateTenantFunc: func(_ context.Context, _ string, _ int, in domain.TenantWrite) (*domain.Tenant, error) {
				gotWrite = in

				return &domain.Tenant{ID: 1, Name: "ACME"}, nil
			},
		}, "token")

		handler := handlers.NewUpdateTenantHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.TenantUpdateInput{ID: 1})
		if err != nil || result.IsError {
			t.Fatalf("handler failed: result.IsError=%v err=%v", result.IsError, err)
		}
		if gotWrite.Group != nil {
			t.Errorf("update left Group = %v, want nil when not provided", gotWrite.Group)
		}
	})
}
