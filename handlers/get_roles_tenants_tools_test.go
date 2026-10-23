package handlers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/domain"
	"github.com/teran/mcp-netbox/handlers"
	"github.com/teran/mcp-netbox/mockrepo"
)

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetRolesHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListRolesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.Role], error) {
				return &domain.PaginatedResponse[domain.Role]{
					Count:   1,
					Results: []domain.Role{{ID: 1, Name: "prod"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetRolesHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.RolesInput{})
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
			t.Errorf("Results = %+v, want single role prod", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListRolesFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.Role], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.Role]{Count: 0, Results: []domain.Role{}}, nil
			},
		}, "token")

		handler := handlers.NewGetRolesHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.RolesInput{
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
			ListRolesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.Role], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetRolesHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.RolesInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetTenantsHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListTenantsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.Tenant], error) {
				return &domain.PaginatedResponse[domain.Tenant]{
					Count:   1,
					Results: []domain.Tenant{{ID: 1, Name: "ACME"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetTenantsHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.TenantsInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].Name != "ACME" {
			t.Errorf("Results = %+v, want single tenant ACME", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListTenantsFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.Tenant], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.Tenant]{Count: 0, Results: []domain.Tenant{}}, nil
			},
		}, "token")

		handler := handlers.NewGetTenantsHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.TenantsInput{
			Q:        "acme",
			Group:    "7",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "acme" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "acme")
		}
		if gotParams["group"] != "7" {
			t.Errorf("group param = %q, want %q", gotParams["group"], "7")
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
			ListTenantsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.Tenant], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetTenantsHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.TenantsInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
