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
func TestCreateIPRangeHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.IPRangeWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateIPRangeFunc: func(_ context.Context, _ string, in domain.IPRangeWrite) (*domain.IPRange, error) {
				gotWrite = in

				return &domain.IPRange{ID: 1, StartAddress: in.StartAddress, EndAddress: in.EndAddress}, nil
			},
		}, "token")

		handler := handlers.NewCreateIPRangeHandler(svc)
		vrf := 2
		result, out, err := handler(context.Background(), nil, handlers.IPRangeCreateInput{
			StartAddress: "10.0.0.1", EndAddress: "10.0.0.10", VRF: &vrf,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.StartAddress != "10.0.0.1" || out.Data.EndAddress != "10.0.0.10" {
			t.Errorf("out.Data = %+v, want id 1 range 10.0.0.1-10.0.0.10", out.Data)
		}
		if gotWrite.VRF == nil || *gotWrite.VRF != 2 {
			t.Errorf("create mapped VRF = %v, want 2", gotWrite.VRF)
		}
	})

	t.Run("start_address required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateIPRangeHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.IPRangeCreateInput{
			StartAddress: "", EndAddress: "10.0.0.10",
		})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing start_address")
		}
	})

	t.Run("end_address required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateIPRangeHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.IPRangeCreateInput{
			StartAddress: "10.0.0.1", EndAddress: "",
		})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing end_address")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"start_address":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateIPRangeFunc: func(_ context.Context, _ string, _ domain.IPRangeWrite) (*domain.IPRange, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateIPRangeHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.IPRangeCreateInput{
			StartAddress: "10.0.0.1", EndAddress: "10.0.0.10",
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
			CreateIPRangeFunc: func(_ context.Context, _ string, _ domain.IPRangeWrite) (*domain.IPRange, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateIPRangeHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.IPRangeCreateInput{
			StartAddress: "10.0.0.1", EndAddress: "10.0.0.10",
		})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateIPRangeHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.IPRangeWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateIPRangeFunc: func(
				_ context.Context, _ string, id int, in domain.IPRangeWrite,
			) (*domain.IPRange, error) {
				gotWrite = in

				return &domain.IPRange{ID: id, StartAddress: "10.0.0.1"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateIPRangeHandler(svc)
		status := "active"
		result, out, err := handler(context.Background(), nil, handlers.IPRangeUpdateInput{ID: 7, Status: &status})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
		if gotWrite.Status == nil || *gotWrite.Status != "active" {
			t.Errorf("update mapped Status = %v, want active", gotWrite.Status)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateIPRangeHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.IPRangeUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"end_address":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateIPRangeFunc: func(
				_ context.Context, _ string, _ int, _ domain.IPRangeWrite,
			) (*domain.IPRange, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateIPRangeHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.IPRangeUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteIPRangeHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteIPRangeFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteIPRangeHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.IPRangeDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteIPRangeHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.IPRangeDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteIPRangeFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteIPRangeHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.IPRangeDeleteInput{ID: 3})
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

func TestCreateIPRangeHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateIPRangeHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.IPRangeCreateInput{
		StartAddress: "10.0.0.1", EndAddress: "10.0.0.10",
	})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateIPRangeHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateIPRangeHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.IPRangeUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteIPRangeHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteIPRangeHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.IPRangeDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetIPRangesHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListIPRangesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.IPRange], error) {
				return &domain.PaginatedResponse[domain.IPRange]{
					Count:   1,
					Results: []domain.IPRange{{ID: 1, StartAddress: "10.0.0.1", EndAddress: "10.0.0.10"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetIPRangesHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.IPRangesInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].StartAddress != "10.0.0.1" {
			t.Errorf("Results = %+v, want single IP range", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListIPRangesFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.IPRange], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.IPRange]{Count: 0, Results: []domain.IPRange{}}, nil
			},
		}, "token")

		handler := handlers.NewGetIPRangesHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.IPRangesInput{
			Q:        "10.0.0.1",
			VRF:      "vrf1",
			Tenant:   "acme",
			Status:   "active",
			Role:     "service",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "10.0.0.1" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "10.0.0.1")
		}
		if gotParams["vrf"] != "vrf1" {
			t.Errorf("vrf param = %q, want %q", gotParams["vrf"], "vrf1")
		}
		if gotParams["tenant"] != "acme" {
			t.Errorf("tenant param = %q, want %q", gotParams["tenant"], "acme")
		}
		if gotParams["status"] != "active" {
			t.Errorf("status param = %q, want %q", gotParams["status"], "active")
		}
		if gotParams["role"] != "service" {
			t.Errorf("role param = %q, want %q", gotParams["role"], "service")
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
			ListIPRangesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.IPRange], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetIPRangesHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.IPRangesInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
