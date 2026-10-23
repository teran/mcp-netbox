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
func TestCreateAggregateHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.AggregateWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateAggregateFunc: func(_ context.Context, _ string, in domain.AggregateWrite) (*domain.Aggregate, error) {
				gotWrite = in

				return &domain.Aggregate{ID: 1, Prefix: in.Prefix}, nil
			},
		}, "token")

		handler := handlers.NewCreateAggregateHandler(svc)
		rir := 3
		result, out, err := handler(context.Background(), nil, handlers.AggregateCreateInput{
			Prefix: "10.0.0.0/8", RIR: &rir,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Prefix != "10.0.0.0/8" {
			t.Errorf("out.Data = %+v, want id 1 prefix 10.0.0.0/8", out.Data)
		}
		if gotWrite.RIR == nil || *gotWrite.RIR != 3 {
			t.Errorf("create mapped RIR = %v, want 3", gotWrite.RIR)
		}
	})

	t.Run("prefix required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateAggregateHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.AggregateCreateInput{Prefix: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing prefix")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"prefix":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateAggregateFunc: func(_ context.Context, _ string, _ domain.AggregateWrite) (*domain.Aggregate, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateAggregateHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.AggregateCreateInput{Prefix: "10.0.0.0/8"})
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
			CreateAggregateFunc: func(_ context.Context, _ string, _ domain.AggregateWrite) (*domain.Aggregate, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateAggregateHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.AggregateCreateInput{Prefix: "10.0.0.0/8"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateAggregateHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.AggregateWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateAggregateFunc: func(
				_ context.Context, _ string, id int, in domain.AggregateWrite,
			) (*domain.Aggregate, error) {
				gotWrite = in

				return &domain.Aggregate{ID: id, Prefix: "10.0.0.0/8"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateAggregateHandler(svc)
		dateAdded := "2024-01-01"
		result, out, err := handler(context.Background(), nil, handlers.AggregateUpdateInput{ID: 7, DateAdded: &dateAdded})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
		if gotWrite.DateAdded == nil || *gotWrite.DateAdded != "2024-01-01" {
			t.Errorf("update mapped DateAdded = %v, want 2024-01-01", gotWrite.DateAdded)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateAggregateHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.AggregateUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"rir":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateAggregateFunc: func(
				_ context.Context, _ string, _ int, _ domain.AggregateWrite,
			) (*domain.Aggregate, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateAggregateHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.AggregateUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteAggregateHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteAggregateFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteAggregateHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.AggregateDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteAggregateHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.AggregateDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteAggregateFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteAggregateHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.AggregateDeleteInput{ID: 3})
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

func TestCreateAggregateHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateAggregateHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.AggregateCreateInput{Prefix: "10.0.0.0/8"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateAggregateHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateAggregateHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.AggregateUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteAggregateHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteAggregateHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.AggregateDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetAggregatesHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListAggregatesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.Aggregate], error) {
				return &domain.PaginatedResponse[domain.Aggregate]{
					Count:   1,
					Results: []domain.Aggregate{{ID: 1, Prefix: "10.0.0.0/8"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetAggregatesHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.AggregatesInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].Prefix != "10.0.0.0/8" {
			t.Errorf("Results = %+v, want single aggregate 10.0.0.0/8", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListAggregatesFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.Aggregate], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.Aggregate]{Count: 0, Results: []domain.Aggregate{}}, nil
			},
		}, "token")

		handler := handlers.NewGetAggregatesHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.AggregatesInput{
			Q:        "10.0.0.0",
			RIR:      "arin",
			Tenant:   "acme",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "10.0.0.0" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "10.0.0.0")
		}
		if gotParams["rir"] != "arin" {
			t.Errorf("rir param = %q, want %q", gotParams["rir"], "arin")
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
			ListAggregatesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.Aggregate], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetAggregatesHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.AggregatesInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
