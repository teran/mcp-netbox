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
func TestCreateServiceTemplateHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.ServiceTemplateWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateServiceTemplateFunc: func(
				_ context.Context, _ string, in domain.ServiceTemplateWrite,
			) (*domain.ServiceTemplate, error) {
				gotWrite = in

				return &domain.ServiceTemplate{ID: 1, Name: in.Name, Ports: in.Ports}, nil
			},
		}, "token")

		handler := handlers.NewCreateServiceTemplateHandler(svc)
		protocol := "tcp"
		result, out, err := handler(context.Background(), nil, handlers.ServiceTemplateCreateInput{
			Name: "web", Protocol: &protocol, Ports: []int{80, 443},
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "web" {
			t.Errorf("out.Data = %+v, want id 1 name web", out.Data)
		}
		if gotWrite.Protocol == nil || *gotWrite.Protocol != "tcp" {
			t.Errorf("create mapped Protocol = %v, want tcp", gotWrite.Protocol)
		}
		if len(gotWrite.Ports) != 2 || gotWrite.Ports[0] != 80 || gotWrite.Ports[1] != 443 {
			t.Errorf("create mapped Ports = %v, want [80 443]", gotWrite.Ports)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateServiceTemplateHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ServiceTemplateCreateInput{Name: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"ports":["Enter a valid value."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateServiceTemplateFunc: func(
				_ context.Context, _ string, _ domain.ServiceTemplateWrite,
			) (*domain.ServiceTemplate, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateServiceTemplateHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ServiceTemplateCreateInput{
			Name: "web", Ports: []int{70000},
		})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
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
			CreateServiceTemplateFunc: func(
				_ context.Context, _ string, _ domain.ServiceTemplateWrite,
			) (*domain.ServiceTemplate, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateServiceTemplateHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ServiceTemplateCreateInput{Name: "web"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateServiceTemplateHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateServiceTemplateFunc: func(
				_ context.Context, _ string, id int, _ domain.ServiceTemplateWrite,
			) (*domain.ServiceTemplate, error) {
				return &domain.ServiceTemplate{ID: id, Name: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateServiceTemplateHandler(svc)
		desc := "https"
		result, out, err := handler(context.Background(), nil, handlers.ServiceTemplateUpdateInput{ID: 7, Description: &desc})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateServiceTemplateHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ServiceTemplateUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateServiceTemplateFunc: func(
				_ context.Context, _ string, _ int, _ domain.ServiceTemplateWrite,
			) (*domain.ServiceTemplate, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateServiceTemplateHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ServiceTemplateUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteServiceTemplateHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteServiceTemplateFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteServiceTemplateHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ServiceTemplateDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteServiceTemplateHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ServiceTemplateDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteServiceTemplateFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteServiceTemplateHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ServiceTemplateDeleteInput{ID: 3})
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

func TestCreateServiceTemplateHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateServiceTemplateHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ServiceTemplateCreateInput{Name: "web"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateServiceTemplateHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateServiceTemplateHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ServiceTemplateUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteServiceTemplateHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteServiceTemplateHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ServiceTemplateDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetServiceTemplatesHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListServiceTemplatesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ServiceTemplate], error) {
				return &domain.PaginatedResponse[domain.ServiceTemplate]{
					Count:   1,
					Results: []domain.ServiceTemplate{{ID: 1, Name: "web"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetServiceTemplatesHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.ServiceTemplatesInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].Name != "web" {
			t.Errorf("Results = %+v, want single service template web", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListServiceTemplatesFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.ServiceTemplate], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.ServiceTemplate]{Count: 0, Results: []domain.ServiceTemplate{}}, nil
			},
		}, "token")

		handler := handlers.NewGetServiceTemplatesHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.ServiceTemplatesInput{
			Q:        "web",
			Protocol: "tcp",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "web" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "web")
		}
		if gotParams["protocol"] != "tcp" {
			t.Errorf("protocol param = %q, want %q", gotParams["protocol"], "tcp")
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
			ListServiceTemplatesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ServiceTemplate], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetServiceTemplatesHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ServiceTemplatesInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
