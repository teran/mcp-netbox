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
func TestCreateTagHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateTagFunc: func(_ context.Context, _ string, in domain.TagWrite) (*domain.Tag, error) {
				return &domain.Tag{ID: 1, Name: in.Name, Slug: "prod"}, nil
			},
		}, "token")

		handler := handlers.NewCreateTagHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.TagCreateInput{Name: "prod"})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "prod" {
			t.Errorf("out.Data = %+v, want id 1 name prod", out.Data)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateTagHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.TagCreateInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateTagFunc: func(_ context.Context, _ string, _ domain.TagWrite) (*domain.Tag, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateTagHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.TagCreateInput{Name: "A"})
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
			CreateTagFunc: func(_ context.Context, _ string, _ domain.TagWrite) (*domain.Tag, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateTagHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.TagCreateInput{Name: "A"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateTagHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateTagFunc: func(
				_ context.Context, _ string, id int, _ domain.TagWrite,
			) (*domain.Tag, error) {
				return &domain.Tag{ID: id, Name: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateTagHandler(svc)
		name := "Renamed"
		result, out, err := handler(context.Background(), nil, handlers.TagUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateTagHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.TagUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"slug":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateTagFunc: func(
				_ context.Context, _ string, _ int, _ domain.TagWrite,
			) (*domain.Tag, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateTagHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.TagUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteTagHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteTagFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteTagHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.TagDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteTagHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.TagDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteTagFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteTagHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.TagDeleteInput{ID: 3})
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

func TestCreateTagHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateTagHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.TagCreateInput{Name: "A"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateTagHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateTagHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.TagUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteTagHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteTagHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.TagDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetTagsHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListTagsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.Tag], error) {
				return &domain.PaginatedResponse[domain.Tag]{
					Count:   1,
					Results: []domain.Tag{{ID: 1, Name: "prod", Slug: "prod"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetTagsHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.TagsInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].Name != "prod" {
			t.Errorf("Results = %+v, want single tag prod", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListTagsFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.Tag], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.Tag]{Count: 0, Results: []domain.Tag{}}, nil
			},
		}, "token")

		handler := handlers.NewGetTagsHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.TagsInput{
			Q:        "prod",
			Tag:      "infra",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "prod" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "prod")
		}
		if gotParams["tag"] != "infra" {
			t.Errorf("tag param = %q, want %q", gotParams["tag"], "infra")
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
			ListTagsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.Tag], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetTagsHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.TagsInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
