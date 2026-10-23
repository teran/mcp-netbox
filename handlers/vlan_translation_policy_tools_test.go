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
func TestCreateVLANTranslationPolicyHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.VLANTranslationPolicyWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVLANTranslationPolicyFunc: func(
				_ context.Context, _ string, in domain.VLANTranslationPolicyWrite,
			) (*domain.VLANTranslationPolicy, error) {
				gotWrite = in

				return &domain.VLANTranslationPolicy{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateVLANTranslationPolicyHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.VLANTranslationPolicyCreateInput{
			Name: "datacenter", Tags: []string{"prod"},
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "datacenter" {
			t.Errorf("out.Data = %+v, want id 1 name datacenter", out.Data)
		}
		if len(gotWrite.Tags) != 1 || gotWrite.Tags[0] != "prod" {
			t.Errorf("create mapped Tags = %v, want [prod]", gotWrite.Tags)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateVLANTranslationPolicyHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VLANTranslationPolicyCreateInput{Name: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVLANTranslationPolicyFunc: func(
				_ context.Context, _ string, _ domain.VLANTranslationPolicyWrite,
			) (*domain.VLANTranslationPolicy, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateVLANTranslationPolicyHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VLANTranslationPolicyCreateInput{Name: "datacenter"})
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
			CreateVLANTranslationPolicyFunc: func(
				_ context.Context, _ string, _ domain.VLANTranslationPolicyWrite,
			) (*domain.VLANTranslationPolicy, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateVLANTranslationPolicyHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VLANTranslationPolicyCreateInput{Name: "datacenter"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateVLANTranslationPolicyHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.VLANTranslationPolicyWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateVLANTranslationPolicyFunc: func(
				_ context.Context, _ string, id int, in domain.VLANTranslationPolicyWrite,
			) (*domain.VLANTranslationPolicy, error) {
				gotWrite = in

				return &domain.VLANTranslationPolicy{ID: id, Name: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateVLANTranslationPolicyHandler(svc)
		desc := "core"
		result, out, err := handler(context.Background(), nil, handlers.VLANTranslationPolicyUpdateInput{
			ID: 7, Description: &desc,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
		if gotWrite.Description == nil || *gotWrite.Description != "core" {
			t.Errorf("update mapped Description = %v, want core", gotWrite.Description)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateVLANTranslationPolicyHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VLANTranslationPolicyUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateVLANTranslationPolicyFunc: func(
				_ context.Context, _ string, _ int, _ domain.VLANTranslationPolicyWrite,
			) (*domain.VLANTranslationPolicy, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateVLANTranslationPolicyHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VLANTranslationPolicyUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteVLANTranslationPolicyHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteVLANTranslationPolicyFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteVLANTranslationPolicyHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VLANTranslationPolicyDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteVLANTranslationPolicyHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VLANTranslationPolicyDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteVLANTranslationPolicyFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteVLANTranslationPolicyHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VLANTranslationPolicyDeleteInput{ID: 3})
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

func TestCreateVLANTranslationPolicyHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateVLANTranslationPolicyHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VLANTranslationPolicyCreateInput{Name: "datacenter"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateVLANTranslationPolicyHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateVLANTranslationPolicyHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VLANTranslationPolicyUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteVLANTranslationPolicyHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteVLANTranslationPolicyHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VLANTranslationPolicyDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetVLANTranslationPoliciesHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListVLANTranslationPoliciesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VLANTranslationPolicy], error) {
				return &domain.PaginatedResponse[domain.VLANTranslationPolicy]{
					Count:   1,
					Results: []domain.VLANTranslationPolicy{{ID: 1, Name: "datacenter"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetVLANTranslationPoliciesHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.VLANTranslationPoliciesInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].Name != "datacenter" {
			t.Errorf("Results = %+v, want single policy datacenter", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListVLANTranslationPoliciesFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.VLANTranslationPolicy], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.VLANTranslationPolicy]{
					Count: 0, Results: []domain.VLANTranslationPolicy{},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetVLANTranslationPoliciesHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.VLANTranslationPoliciesInput{
			Q:        "datacenter",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "datacenter" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "datacenter")
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
			ListVLANTranslationPoliciesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VLANTranslationPolicy], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetVLANTranslationPoliciesHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VLANTranslationPoliciesInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
