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
func TestCreateSiteGroupHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateSiteGroupFunc: func(_ context.Context, _ string, in domain.SiteGroupWrite) (*domain.SiteGroup, error) {
				return &domain.SiteGroup{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateSiteGroupHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.SiteGroupCreateInput{Name: "Group A"})
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
		handler := handlers.NewCreateSiteGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.SiteGroupCreateInput{Name: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateSiteGroupFunc: func(_ context.Context, _ string, _ domain.SiteGroupWrite) (*domain.SiteGroup, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateSiteGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.SiteGroupCreateInput{Name: "A"})
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
			CreateSiteGroupFunc: func(_ context.Context, _ string, _ domain.SiteGroupWrite) (*domain.SiteGroup, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateSiteGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.SiteGroupCreateInput{Name: "A"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateSiteGroupHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateSiteGroupFunc: func(
				_ context.Context, _ string, id int, _ domain.SiteGroupWrite,
			) (*domain.SiteGroup, error) {
				return &domain.SiteGroup{ID: id, Name: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateSiteGroupHandler(svc)
		name := "Renamed"
		result, out, err := handler(context.Background(), nil, handlers.SiteGroupUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateSiteGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.SiteGroupUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"slug":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateSiteGroupFunc: func(
				_ context.Context, _ string, _ int, _ domain.SiteGroupWrite,
			) (*domain.SiteGroup, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateSiteGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.SiteGroupUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteSiteGroupHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteSiteGroupFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteSiteGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.SiteGroupDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteSiteGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.SiteGroupDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteSiteGroupFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteSiteGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.SiteGroupDeleteInput{ID: 3})
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

func TestCreateSiteGroupHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateSiteGroupHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.SiteGroupCreateInput{Name: "A"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateSiteGroupHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateSiteGroupHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.SiteGroupUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteSiteGroupHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteSiteGroupHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.SiteGroupDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetSiteGroupsHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListSiteGroupsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.SiteGroup], error) {
				return &domain.PaginatedResponse[domain.SiteGroup]{
					Count:   1,
					Results: []domain.SiteGroup{{ID: 1, Name: "EU-West"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetSiteGroupsHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.SiteGroupsInput{})
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
			t.Errorf("Results = %+v, want single site group EU-West", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListSiteGroupsFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.SiteGroup], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.SiteGroup]{Count: 0, Results: []domain.SiteGroup{}}, nil
			},
		}, "token")

		handler := handlers.NewGetSiteGroupsHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.SiteGroupsInput{
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
		if gotParams["offset"] != "50" {
			t.Errorf("offset param = %q, want %q", gotParams["offset"], "50")
		}
		if gotParams["limit"] != "50" {
			t.Errorf("limit param = %q, want %q", gotParams["limit"], "50")
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListSiteGroupsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.SiteGroup], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetSiteGroupsHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.SiteGroupsInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
