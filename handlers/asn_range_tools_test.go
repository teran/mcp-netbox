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
func TestCreateASNRangeHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.ASNRangeWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateASNRangeFunc: func(_ context.Context, _ string, in domain.ASNRangeWrite) (*domain.ASNRange, error) {
				gotWrite = in

				return &domain.ASNRange{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateASNRangeHandler(svc)
		start := 64512
		end := 64520
		slug := "private"
		result, out, err := handler(context.Background(), nil, handlers.ASNRangeCreateInput{
			Name: "PRIVATE", Slug: &slug, Start: &start, End: &end,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "PRIVATE" {
			t.Errorf("out.Data = %+v, want id 1 name PRIVATE", out.Data)
		}
		if gotWrite.Start == nil || *gotWrite.Start != 64512 || gotWrite.End == nil || *gotWrite.End != 64520 {
			t.Errorf("create mapped Start/End = %v/%v, want 64512/64520", gotWrite.Start, gotWrite.End)
		}
		if gotWrite.Slug == nil || *gotWrite.Slug != "private" {
			t.Errorf("create mapped Slug = %v, want private", gotWrite.Slug)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateASNRangeHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ASNRangeCreateInput{Name: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateASNRangeFunc: func(_ context.Context, _ string, _ domain.ASNRangeWrite) (*domain.ASNRange, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateASNRangeHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ASNRangeCreateInput{Name: "PRIVATE"})
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
			CreateASNRangeFunc: func(_ context.Context, _ string, _ domain.ASNRangeWrite) (*domain.ASNRange, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateASNRangeHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ASNRangeCreateInput{Name: "PRIVATE"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateASNRangeHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.ASNRangeWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateASNRangeFunc: func(
				_ context.Context, _ string, id int, in domain.ASNRangeWrite,
			) (*domain.ASNRange, error) {
				gotWrite = in

				return &domain.ASNRange{ID: id, Name: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateASNRangeHandler(svc)
		rir := 3
		result, out, err := handler(context.Background(), nil, handlers.ASNRangeUpdateInput{ID: 7, RIR: &rir})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
		if gotWrite.RIR == nil || *gotWrite.RIR != 3 {
			t.Errorf("update mapped RIR = %v, want 3", gotWrite.RIR)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateASNRangeHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ASNRangeUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"rir":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateASNRangeFunc: func(
				_ context.Context, _ string, _ int, _ domain.ASNRangeWrite,
			) (*domain.ASNRange, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateASNRangeHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ASNRangeUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteASNRangeHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteASNRangeFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteASNRangeHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ASNRangeDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteASNRangeHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ASNRangeDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteASNRangeFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteASNRangeHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ASNRangeDeleteInput{ID: 3})
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

func TestCreateASNRangeHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateASNRangeHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ASNRangeCreateInput{Name: "PRIVATE"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateASNRangeHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateASNRangeHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ASNRangeUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteASNRangeHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteASNRangeHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ASNRangeDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetASNRangesHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListASNRangesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ASNRange], error) {
				return &domain.PaginatedResponse[domain.ASNRange]{
					Count:   1,
					Results: []domain.ASNRange{{ID: 1, Name: "PRIVATE"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetASNRangesHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.ASNRangesInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].Name != "PRIVATE" {
			t.Errorf("Results = %+v, want single ASN range PRIVATE", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListASNRangesFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.ASNRange], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.ASNRange]{Count: 0, Results: []domain.ASNRange{}}, nil
			},
		}, "token")

		handler := handlers.NewGetASNRangesHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.ASNRangesInput{
			Q:        "PRIVATE",
			RIR:      "arin",
			Tenant:   "acme",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "PRIVATE" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "PRIVATE")
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
			ListASNRangesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ASNRange], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetASNRangesHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ASNRangesInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
