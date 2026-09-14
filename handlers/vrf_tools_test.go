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

func TestCreateVrfHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVrfFunc: func(_ context.Context, _ string, in domain.VrfWrite) (*domain.Vrf, error) {
				return &domain.Vrf{ID: 1, Name: in.Name, Rd: in.Rd}, nil
			},
		}, "token")

		handler := handlers.NewCreateVrfHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.VrfCreateInput{Name: "prod", Rd: "65000:1"})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "prod" || out.Data.Rd != "65000:1" {
			t.Errorf("out.Data = %+v, want id 1 name 'prod' rd '65000:1'", out.Data)
		}
	})

	t.Run("name and rd required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateVrfHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VrfCreateInput{Name: "prod"})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing rd")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVrfFunc: func(_ context.Context, _ string, _ domain.VrfWrite) (*domain.Vrf, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateVrfHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VrfCreateInput{Name: "prod", Rd: "65000:1"})
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
			CreateVrfFunc: func(_ context.Context, _ string, _ domain.VrfWrite) (*domain.Vrf, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateVrfHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VrfCreateInput{Name: "prod", Rd: "65000:1"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateVrfHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateVrfFunc: func(_ context.Context, _ string, id int, _ domain.VrfWrite) (*domain.Vrf, error) {
				return &domain.Vrf{ID: id, Name: "prod", Rd: "65000:1"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateVrfHandler(svc)
		name := "prod"
		result, out, err := handler(context.Background(), nil, handlers.VrfUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateVrfHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VrfUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"rd":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateVrfFunc: func(_ context.Context, _ string, _ int, _ domain.VrfWrite) (*domain.Vrf, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateVrfHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VrfUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteVrfHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteVrfFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		}, "token")
		handler := handlers.NewDeleteVrfHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VrfDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteVrfHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VrfDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteVrfFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteVrfHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VrfDeleteInput{ID: 3})
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

func TestCreateVrfHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateVrfHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VrfCreateInput{Name: "prod", Rd: "65000:1"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateVrfHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateVrfHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VrfUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteVrfHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteVrfHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VrfDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}
