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
func TestCreateCustomFieldHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateCustomFieldFunc: func(_ context.Context, _ string, in domain.CustomFieldWrite) (*domain.CustomField, error) {
				return &domain.CustomField{ID: 1, Name: in.Name, Label: in.Label}, nil
			},
		}, "token")

		handler := handlers.NewCreateCustomFieldHandler(svc)
		result, out, err := handler(
			context.Background(), nil,
			handlers.CustomFieldCreateInput{Name: "env", Label: "Environment"},
		)
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "env" || out.Data.Label != "Environment" {
			t.Errorf("out.Data = %+v, want id 1 name env label Environment", out.Data)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateCustomFieldHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CustomFieldCreateInput{Name: "", Label: "Environment"})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("label required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateCustomFieldHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CustomFieldCreateInput{Name: "env", Label: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing label")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateCustomFieldFunc: func(_ context.Context, _ string, _ domain.CustomFieldWrite) (*domain.CustomField, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateCustomFieldHandler(svc)
		result, _, err := handler(
			context.Background(), nil,
			handlers.CustomFieldCreateInput{Name: "env", Label: "Environment"},
		)
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
			CreateCustomFieldFunc: func(_ context.Context, _ string, _ domain.CustomFieldWrite) (*domain.CustomField, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateCustomFieldHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CustomFieldCreateInput{Name: "env", Label: "Environment"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateCustomFieldHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateCustomFieldFunc: func(
				_ context.Context, _ string, id int, _ domain.CustomFieldWrite,
			) (*domain.CustomField, error) {
				return &domain.CustomField{ID: id, Name: "renamed", Label: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateCustomFieldHandler(svc)
		name := "renamed"
		result, out, err := handler(context.Background(), nil, handlers.CustomFieldUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateCustomFieldHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CustomFieldUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"label":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateCustomFieldFunc: func(
				_ context.Context, _ string, _ int, _ domain.CustomFieldWrite,
			) (*domain.CustomField, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateCustomFieldHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CustomFieldUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteCustomFieldHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteCustomFieldFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteCustomFieldHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CustomFieldDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteCustomFieldHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CustomFieldDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteCustomFieldFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteCustomFieldHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CustomFieldDeleteInput{ID: 3})
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

func TestCreateCustomFieldHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateCustomFieldHandler(nil)
	result, _, err := handler(
		context.Background(), nil,
		handlers.CustomFieldCreateInput{Name: "env", Label: "Environment"},
	)
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateCustomFieldHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateCustomFieldHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.CustomFieldUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteCustomFieldHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteCustomFieldHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.CustomFieldDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetCustomFieldsHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListCustomFieldsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.CustomField], error) {
				return &domain.PaginatedResponse[domain.CustomField]{
					Count:   1,
					Results: []domain.CustomField{{ID: 1, Name: "env", Label: "Environment"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetCustomFieldsHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.CustomFieldsInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].Name != "env" {
			t.Errorf("Results = %+v, want single custom field env", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListCustomFieldsFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.CustomField], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.CustomField]{Count: 0, Results: []domain.CustomField{}}, nil
			},
		}, "token")

		handler := handlers.NewGetCustomFieldsHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.CustomFieldsInput{
			Q:           "env",
			ContentType: "dcim.device",
			Tag:         "infra",
			Page:        2,
			PageSize:    50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "env" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "env")
		}
		if gotParams["content_type"] != "dcim.device" {
			t.Errorf("content_type param = %q, want %q", gotParams["content_type"], "dcim.device")
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
			ListCustomFieldsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.CustomField], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetCustomFieldsHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CustomFieldsInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
