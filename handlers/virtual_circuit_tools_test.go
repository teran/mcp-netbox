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
func TestCreateVirtualCircuitHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVirtualCircuitFunc: func(
				_ context.Context, _ string, in domain.VirtualCircuitWrite,
			) (*domain.VirtualCircuit, error) {
				return &domain.VirtualCircuit{ID: 1, CID: in.CID}, nil
			},
		}, "token")

		handler := handlers.NewCreateVirtualCircuitHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.VirtualCircuitCreateInput{CID: "VC-1"})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.CID != "VC-1" {
			t.Errorf("out.Data = %+v, want id 1 cid VC-1", out.Data)
		}
	})

	t.Run("cid required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateVirtualCircuitHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualCircuitCreateInput{CID: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing cid")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVirtualCircuitFunc: func(
				_ context.Context, _ string, _ domain.VirtualCircuitWrite,
			) (*domain.VirtualCircuit, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateVirtualCircuitHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitCreateInput{CID: "A"})
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
			CreateVirtualCircuitFunc: func(
				_ context.Context, _ string, _ domain.VirtualCircuitWrite,
			) (*domain.VirtualCircuit, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateVirtualCircuitHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualCircuitCreateInput{CID: "A"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateVirtualCircuitHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateVirtualCircuitFunc: func(
				_ context.Context, _ string, id int, _ domain.VirtualCircuitWrite,
			) (*domain.VirtualCircuit, error) {
				return &domain.VirtualCircuit{ID: id, CID: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateVirtualCircuitHandler(svc)
		cid := "Renamed"
		result, out, err := handler(context.Background(), nil, handlers.VirtualCircuitUpdateInput{ID: 7, CID: &cid})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateVirtualCircuitHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualCircuitUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"type":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateVirtualCircuitFunc: func(
				_ context.Context, _ string, _ int, _ domain.VirtualCircuitWrite,
			) (*domain.VirtualCircuit, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateVirtualCircuitHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteVirtualCircuitHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteVirtualCircuitFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteVirtualCircuitHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteVirtualCircuitHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualCircuitDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteVirtualCircuitFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteVirtualCircuitHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitDeleteInput{ID: 3})
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

func TestCreateVirtualCircuitHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateVirtualCircuitHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitCreateInput{CID: "A"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateVirtualCircuitHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateVirtualCircuitHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteVirtualCircuitHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteVirtualCircuitHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetVirtualCircuitsHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListVirtualCircuitsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualCircuit], error) {
				return &domain.PaginatedResponse[domain.VirtualCircuit]{
					Count:   1,
					Results: []domain.VirtualCircuit{{ID: 1, CID: "VC-1"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetVirtualCircuitsHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.VirtualCircuitsInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].CID != "VC-1" {
			t.Errorf("Results = %+v, want single virtual circuit VC-1", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListVirtualCircuitsFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualCircuit], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.VirtualCircuit]{Count: 0, Results: []domain.VirtualCircuit{}}, nil
			},
		}, "token")

		handler := handlers.NewGetVirtualCircuitsHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.VirtualCircuitsInput{
			Q:               "vc",
			Type:            "7",
			ProviderNetwork: "3",
			Status:          "planned",
			Tag:             "prod",
			Page:            2,
			PageSize:        50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "vc" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "vc")
		}
		if gotParams["type"] != "7" {
			t.Errorf("type param = %q, want %q", gotParams["type"], "7")
		}
		if gotParams["provider_network"] != "3" {
			t.Errorf("provider_network param = %q, want %q", gotParams["provider_network"], "3")
		}
		if gotParams["status"] != "planned" {
			t.Errorf("status param = %q, want %q", gotParams["status"], "planned")
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
			ListVirtualCircuitsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualCircuit], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetVirtualCircuitsHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualCircuitsInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

// TestVirtualCircuitMapsWrite verifies the optional writable fields are mapped
// through the create handler to the domain write DTO.
func TestVirtualCircuitMapsWrite(t *testing.T) {
	t.Parallel()

	var gotWrite domain.VirtualCircuitWrite
	svc := application.NewNetworkService(&mockrepo.MockRepo{
		CreateVirtualCircuitFunc: func(
			_ context.Context, _ string, in domain.VirtualCircuitWrite,
		) (*domain.VirtualCircuit, error) {
			gotWrite = in

			return &domain.VirtualCircuit{ID: 1, CID: in.CID}, nil
		},
	}, "token")

	vcType := 2
	providerNetwork := 3
	status := "planned"
	handler := handlers.NewCreateVirtualCircuitHandler(svc)
	result, _, err := handler(context.Background(), nil, handlers.VirtualCircuitCreateInput{
		CID: "A", Type: &vcType, ProviderNetwork: &providerNetwork, Status: &status,
	})
	if err != nil || result.IsError {
		t.Fatalf("handler failed: result.IsError=%v err=%v", result.IsError, err)
	}
	if gotWrite.Type == nil || *gotWrite.Type != 2 {
		t.Errorf("create mapped Type = %v, want 2", gotWrite.Type)
	}
	if gotWrite.ProviderNetwork == nil || *gotWrite.ProviderNetwork != 3 {
		t.Errorf("create mapped ProviderNetwork = %v, want 3", gotWrite.ProviderNetwork)
	}
	if gotWrite.Status == nil || *gotWrite.Status != "planned" {
		t.Errorf("create mapped Status = %v, want planned", gotWrite.Status)
	}
}
