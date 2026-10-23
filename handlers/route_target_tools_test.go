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
func TestCreateRouteTargetHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.RouteTargetWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateRouteTargetFunc: func(_ context.Context, _ string, in domain.RouteTargetWrite) (*domain.RouteTarget, error) {
				gotWrite = in

				return &domain.RouteTarget{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateRouteTargetHandler(svc)
		tenant := 4
		result, out, err := handler(context.Background(), nil, handlers.RouteTargetCreateInput{
			Name: "65000:100", Tenant: &tenant,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "65000:100" {
			t.Errorf("out.Data = %+v, want id 1 name 65000:100", out.Data)
		}
		if gotWrite.Tenant == nil || *gotWrite.Tenant != 4 {
			t.Errorf("create mapped Tenant = %v, want 4", gotWrite.Tenant)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateRouteTargetHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.RouteTargetCreateInput{Name: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateRouteTargetFunc: func(_ context.Context, _ string, _ domain.RouteTargetWrite) (*domain.RouteTarget, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateRouteTargetHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.RouteTargetCreateInput{Name: "65000:100"})
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
			CreateRouteTargetFunc: func(_ context.Context, _ string, _ domain.RouteTargetWrite) (*domain.RouteTarget, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateRouteTargetHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.RouteTargetCreateInput{Name: "65000:100"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateRouteTargetHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateRouteTargetFunc: func(
				_ context.Context, _ string, id int, _ domain.RouteTargetWrite,
			) (*domain.RouteTarget, error) {
				return &domain.RouteTarget{ID: id, Name: "65000:200"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateRouteTargetHandler(svc)
		name := "65000:200"
		result, out, err := handler(context.Background(), nil, handlers.RouteTargetUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateRouteTargetHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.RouteTargetUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"tenant":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateRouteTargetFunc: func(
				_ context.Context, _ string, _ int, _ domain.RouteTargetWrite,
			) (*domain.RouteTarget, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateRouteTargetHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.RouteTargetUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteRouteTargetHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteRouteTargetFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteRouteTargetHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.RouteTargetDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteRouteTargetHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.RouteTargetDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteRouteTargetFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteRouteTargetHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.RouteTargetDeleteInput{ID: 3})
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

func TestCreateRouteTargetHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateRouteTargetHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.RouteTargetCreateInput{Name: "65000:100"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateRouteTargetHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateRouteTargetHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.RouteTargetUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteRouteTargetHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteRouteTargetHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.RouteTargetDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetRouteTargetsHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListRouteTargetsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.RouteTarget], error) {
				return &domain.PaginatedResponse[domain.RouteTarget]{
					Count:   1,
					Results: []domain.RouteTarget{{ID: 1, Name: "65000:100"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetRouteTargetsHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.RouteTargetsInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].Name != "65000:100" {
			t.Errorf("Results = %+v, want single route target 65000:100", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListRouteTargetsFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.RouteTarget], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.RouteTarget]{Count: 0, Results: []domain.RouteTarget{}}, nil
			},
		}, "token")

		handler := handlers.NewGetRouteTargetsHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.RouteTargetsInput{
			Q:        "65000",
			Tenant:   "acme",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "65000" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "65000")
		}
		if gotParams["tenant"] != "acme" {
			t.Errorf("tenant param = %q, want %q", gotParams["tenant"], "acme")
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
			ListRouteTargetsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.RouteTarget], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetRouteTargetsHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.RouteTargetsInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
