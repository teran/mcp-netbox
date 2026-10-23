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
func TestCreateMACAddressHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateMACAddressFunc: func(_ context.Context, _ string, in domain.MACAddressWrite) (*domain.MACAddress, error) {
				return &domain.MACAddress{ID: 1, Address: in.Address}, nil
			},
		}, "token")

		handler := handlers.NewCreateMACAddressHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.MACAddressCreateInput{
			Address: "00:11:22:33:44:55",
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Address != "00:11:22:33:44:55" {
			t.Errorf("out.Data = %+v, want id 1 address 00:11:22:33:44:55", out.Data)
		}
	})

	t.Run("address required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateMACAddressHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.MACAddressCreateInput{Address: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing address")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"address":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateMACAddressFunc: func(_ context.Context, _ string, _ domain.MACAddressWrite) (*domain.MACAddress, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateMACAddressHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.MACAddressCreateInput{
			Address: "00:11:22:33:44:55",
		})
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
			CreateMACAddressFunc: func(_ context.Context, _ string, _ domain.MACAddressWrite) (*domain.MACAddress, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateMACAddressHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.MACAddressCreateInput{
			Address: "00:11:22:33:44:55",
		})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateMACAddressHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateMACAddressFunc: func(
				_ context.Context, _ string, id int, _ domain.MACAddressWrite,
			) (*domain.MACAddress, error) {
				return &domain.MACAddress{ID: id, Address: "00:11:22:33:44:55"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateMACAddressHandler(svc)
		addr := "00:11:22:33:44:55"
		result, out, err := handler(context.Background(), nil, handlers.MACAddressUpdateInput{ID: 7, Address: &addr})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateMACAddressHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.MACAddressUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"assigned_object_type":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateMACAddressFunc: func(
				_ context.Context, _ string, _ int, _ domain.MACAddressWrite,
			) (*domain.MACAddress, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateMACAddressHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.MACAddressUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteMACAddressHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteMACAddressFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteMACAddressHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.MACAddressDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteMACAddressHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.MACAddressDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteMACAddressFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteMACAddressHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.MACAddressDeleteInput{ID: 3})
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

func TestCreateMACAddressHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateMACAddressHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.MACAddressCreateInput{Address: "00:11:22:33:44:55"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateMACAddressHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateMACAddressHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.MACAddressUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteMACAddressHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteMACAddressHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.MACAddressDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetMACAddressesHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListMACAddressesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.MACAddress], error) {
				return &domain.PaginatedResponse[domain.MACAddress]{
					Count:   1,
					Results: []domain.MACAddress{{ID: 1, Address: "00:11:22:33:44:55"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetMACAddressesHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.MACAddressesInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 || len(out.Results) != 1 || out.Results[0].Address != "00:11:22:33:44:55" {
			t.Errorf("out = %+v, want count 1 with MAC", out)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListMACAddressesFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.MACAddress], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.MACAddress]{Count: 0, Results: []domain.MACAddress{}}, nil
			},
		}, "token")

		handler := handlers.NewGetMACAddressesHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.MACAddressesInput{
			Q:        "00:11",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "00:11" || gotParams["tag"] != "prod" {
			t.Errorf("params = %+v, want q 00:11, tag prod", gotParams)
		}
		if gotParams["offset"] != "50" || gotParams["limit"] != "50" {
			t.Errorf("pagination = %+v, want offset 50 limit 50", gotParams)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListMACAddressesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.MACAddress], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetMACAddressesHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.MACAddressesInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
