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

func TestCreatePrefixHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreatePrefixFunc: func(_ context.Context, _ string, in domain.PrefixWrite) (*domain.Prefix, error) {
				return &domain.Prefix{ID: 1, Prefix: in.Prefix}, nil
			},
		}, "token")

		handler := handlers.NewCreatePrefixHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.PrefixCreateInput{Prefix: "10.0.0.0/24"})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Prefix != "10.0.0.0/24" {
			t.Errorf("out.Data = %+v, want id 1 prefix 10.0.0.0/24", out.Data)
		}
	})

	t.Run("prefix required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreatePrefixHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.PrefixCreateInput{Prefix: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing prefix")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"prefix":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreatePrefixFunc: func(_ context.Context, _ string, _ domain.PrefixWrite) (*domain.Prefix, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreatePrefixHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.PrefixCreateInput{Prefix: "A"})
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
			CreatePrefixFunc: func(_ context.Context, _ string, _ domain.PrefixWrite) (*domain.Prefix, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreatePrefixHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.PrefixCreateInput{Prefix: "A"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdatePrefixHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdatePrefixFunc: func(_ context.Context, _ string, id int, _ domain.PrefixWrite) (*domain.Prefix, error) {
				return &domain.Prefix{ID: id, Prefix: "10.0.0.0/24"}, nil
			},
		}, "token")
		handler := handlers.NewUpdatePrefixHandler(svc)
		prefix := "10.0.0.0/24"
		result, out, err := handler(context.Background(), nil, handlers.PrefixUpdateInput{ID: 7, Prefix: &prefix})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdatePrefixHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.PrefixUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"status":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdatePrefixFunc: func(_ context.Context, _ string, _ int, _ domain.PrefixWrite) (*domain.Prefix, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdatePrefixHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.PrefixUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeletePrefixHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeletePrefixFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		}, "token")
		handler := handlers.NewDeletePrefixHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.PrefixDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeletePrefixHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.PrefixDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeletePrefixFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeletePrefixHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.PrefixDeleteInput{ID: 3})
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

func TestCreatePrefixHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreatePrefixHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.PrefixCreateInput{Prefix: "A"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdatePrefixHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdatePrefixHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.PrefixUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeletePrefixHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeletePrefixHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.PrefixDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}
