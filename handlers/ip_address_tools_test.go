package handlers_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/domain"
	"github.com/teran/mcp-netbox/handlers"
	"github.com/teran/mcp-netbox/internal/mockrepo"
)

func TestCreateIPAddressHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateIPAddressFunc: func(_ context.Context, _ string, in domain.IPAddressWrite) (*domain.IPAddress, error) {
				return &domain.IPAddress{ID: 1, Address: in.Address}, nil
			},
		}, "token")

		handler := handlers.NewCreateIPAddressHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.IPAddressCreateInput{Address: "192.168.1.1/24"})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Address != "192.168.1.1/24" {
			t.Errorf("out.Data = %+v, want id 1 address 192.168.1.1/24", out.Data)
		}
	})

	t.Run("address required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateIPAddressHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.IPAddressCreateInput{Address: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing address")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"address":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateIPAddressFunc: func(_ context.Context, _ string, _ domain.IPAddressWrite) (*domain.IPAddress, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateIPAddressHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.IPAddressCreateInput{Address: "A"})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
		if result.StructuredContent == nil {
			t.Fatal("StructuredContent is nil, want validation_errors")
		}
		sc, ok := result.StructuredContent.(map[string]any)
		if !ok {
			t.Fatalf("StructuredContent type = %T, want map", result.StructuredContent)
		}
		raw, ok := sc["validation_errors"].(json.RawMessage)
		if !ok {
			t.Fatalf("validation_errors type = %T, want json.RawMessage", sc["validation_errors"])
		}
		if string(raw) != string(body) {
			t.Errorf("validation_errors = %s, want %s", raw, body)
		}
	})

	t.Run("generic error", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateIPAddressFunc: func(_ context.Context, _ string, _ domain.IPAddressWrite) (*domain.IPAddress, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateIPAddressHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.IPAddressCreateInput{Address: "A"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateIPAddressHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateIPAddressFunc: func(_ context.Context, _ string, id int, _ domain.IPAddressWrite) (*domain.IPAddress, error) {
				return &domain.IPAddress{ID: id, Address: "10.0.0.1/32"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateIPAddressHandler(svc)
		addr := "10.0.0.1/32"
		result, out, err := handler(context.Background(), nil, handlers.IPAddressUpdateInput{ID: 7, Address: &addr})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateIPAddressHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.IPAddressUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"dns_name":["Enter a valid hostname."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateIPAddressFunc: func(_ context.Context, _ string, _ int, _ domain.IPAddressWrite) (*domain.IPAddress, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateIPAddressHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.IPAddressUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteIPAddressHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteIPAddressFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		}, "token")
		handler := handlers.NewDeleteIPAddressHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.IPAddressDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteIPAddressHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.IPAddressDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteIPAddressFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteIPAddressHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.IPAddressDeleteInput{ID: 3})
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

func TestCreateIPAddressHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateIPAddressHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.IPAddressCreateInput{Address: "A"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}
