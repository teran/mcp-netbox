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
func TestCreateASNHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.ASNWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateASNFunc: func(_ context.Context, _ string, in domain.ASNWrite) (*domain.ASN, error) {
				gotWrite = in

				return &domain.ASN{ID: 1, ASN: in.ASN}, nil
			},
		}, "token")

		handler := handlers.NewCreateASNHandler(svc)
		rir := 3
		result, out, err := handler(context.Background(), nil, handlers.ASNCreateInput{ASN: 64512, RIR: &rir})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.ASN != 64512 {
			t.Errorf("out.Data = %+v, want id 1 asn 64512", out.Data)
		}
		if gotWrite.ASN != 64512 || gotWrite.RIR == nil || *gotWrite.RIR != 3 {
			t.Errorf("create mapped ASN = %d rir %v, want 64512 and 3", gotWrite.ASN, gotWrite.RIR)
		}
	})

	t.Run("asn required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateASNHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ASNCreateInput{ASN: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing asn")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"asn":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateASNFunc: func(_ context.Context, _ string, _ domain.ASNWrite) (*domain.ASN, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateASNHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ASNCreateInput{ASN: 64512})
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
			CreateASNFunc: func(_ context.Context, _ string, _ domain.ASNWrite) (*domain.ASN, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateASNHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ASNCreateInput{ASN: 64512})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateASNHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.ASNWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateASNFunc: func(
				_ context.Context, _ string, id int, in domain.ASNWrite,
			) (*domain.ASN, error) {
				gotWrite = in

				return &domain.ASN{ID: id, ASN: 64513}, nil
			},
		}, "token")
		handler := handlers.NewUpdateASNHandler(svc)
		asn := 64513
		result, out, err := handler(context.Background(), nil, handlers.ASNUpdateInput{ID: 7, ASN: &asn})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
		if gotWrite.ASN != 64513 {
			t.Errorf("update mapped ASN = %d, want 64513", gotWrite.ASN)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateASNHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ASNUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"rir":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateASNFunc: func(
				_ context.Context, _ string, _ int, _ domain.ASNWrite,
			) (*domain.ASN, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateASNHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ASNUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteASNHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteASNFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteASNHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ASNDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteASNHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ASNDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteASNFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteASNHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ASNDeleteInput{ID: 3})
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

func TestCreateASNHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateASNHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ASNCreateInput{ASN: 64512})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateASNHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateASNHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ASNUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteASNHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteASNHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ASNDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetASNsHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListASNsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ASN], error) {
				return &domain.PaginatedResponse[domain.ASN]{
					Count:   1,
					Results: []domain.ASN{{ID: 1, ASN: 64512}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetASNsHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.ASNsInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].ASN != 64512 {
			t.Errorf("Results = %+v, want single ASN 64512", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListASNsFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.ASN], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.ASN]{Count: 0, Results: []domain.ASN{}}, nil
			},
		}, "token")

		handler := handlers.NewGetASNsHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.ASNsInput{
			Q:        "64512",
			RIR:      "arin",
			Tenant:   "acme",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "64512" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "64512")
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
			ListASNsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ASN], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetASNsHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ASNsInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
