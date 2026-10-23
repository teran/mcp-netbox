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
func TestCreateRIRHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.RIRWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateRIRFunc: func(_ context.Context, _ string, in domain.RIRWrite) (*domain.RIR, error) {
				gotWrite = in

				return &domain.RIR{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateRIRHandler(svc)
		isPrivate := true
		result, out, err := handler(context.Background(), nil, handlers.RIRCreateInput{
			Name: "ARIN", IsPrivate: &isPrivate,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "ARIN" {
			t.Errorf("out.Data = %+v, want id 1 name ARIN", out.Data)
		}
		if gotWrite.IsPrivate == nil || !*gotWrite.IsPrivate {
			t.Errorf("create mapped IsPrivate = %v, want true", gotWrite.IsPrivate)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateRIRHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.RIRCreateInput{Name: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateRIRFunc: func(_ context.Context, _ string, _ domain.RIRWrite) (*domain.RIR, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateRIRHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.RIRCreateInput{Name: "ARIN"})
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
			CreateRIRFunc: func(_ context.Context, _ string, _ domain.RIRWrite) (*domain.RIR, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateRIRHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.RIRCreateInput{Name: "ARIN"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateRIRHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateRIRFunc: func(
				_ context.Context, _ string, id int, _ domain.RIRWrite,
			) (*domain.RIR, error) {
				return &domain.RIR{ID: id, Name: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateRIRHandler(svc)
		name := "Renamed"
		result, out, err := handler(context.Background(), nil, handlers.RIRUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateRIRHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.RIRUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"slug":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateRIRFunc: func(
				_ context.Context, _ string, _ int, _ domain.RIRWrite,
			) (*domain.RIR, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateRIRHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.RIRUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteRIRHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteRIRFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteRIRHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.RIRDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteRIRHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.RIRDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteRIRFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteRIRHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.RIRDeleteInput{ID: 3})
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

func TestCreateRIRHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateRIRHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.RIRCreateInput{Name: "ARIN"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateRIRHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateRIRHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.RIRUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteRIRHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteRIRHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.RIRDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetRIRsHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListRIRsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.RIR], error) {
				return &domain.PaginatedResponse[domain.RIR]{
					Count:   1,
					Results: []domain.RIR{{ID: 1, Name: "ARIN"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetRIRsHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.RIRsInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].Name != "ARIN" {
			t.Errorf("Results = %+v, want single RIR ARIN", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListRIRsFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.RIR], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.RIR]{Count: 0, Results: []domain.RIR{}}, nil
			},
		}, "token")

		handler := handlers.NewGetRIRsHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.RIRsInput{
			Q:        "arin",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "arin" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "arin")
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
			ListRIRsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.RIR], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetRIRsHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.RIRsInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
