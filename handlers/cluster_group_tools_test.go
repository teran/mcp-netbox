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

func TestCreateClusterGroupHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateClusterGroupFunc: func(_ context.Context, _ string, in domain.ClusterGroupWrite) (*domain.ClusterGroup, error) {
				return &domain.ClusterGroup{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateClusterGroupHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.ClusterGroupCreateInput{Name: "DC Clusters"})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "DC Clusters" {
			t.Errorf("out.Data = %+v, want id 1 name 'DC Clusters'", out.Data)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateClusterGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ClusterGroupCreateInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateClusterGroupFunc: func(_ context.Context, _ string, _ domain.ClusterGroupWrite) (*domain.ClusterGroup, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateClusterGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ClusterGroupCreateInput{Name: "DC Clusters"})
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
			CreateClusterGroupFunc: func(_ context.Context, _ string, _ domain.ClusterGroupWrite) (*domain.ClusterGroup, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateClusterGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ClusterGroupCreateInput{Name: "DC Clusters"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateClusterGroupHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateClusterGroupFunc: func(_ context.Context, _ string, id int, _ domain.ClusterGroupWrite) (*domain.ClusterGroup, error) {
				return &domain.ClusterGroup{ID: id, Name: "DC Clusters"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateClusterGroupHandler(svc)
		name := "DC Clusters"
		result, out, err := handler(context.Background(), nil, handlers.ClusterGroupUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateClusterGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ClusterGroupUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"slug":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateClusterGroupFunc: func(_ context.Context, _ string, _ int, _ domain.ClusterGroupWrite) (*domain.ClusterGroup, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateClusterGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ClusterGroupUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteClusterGroupHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteClusterGroupFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}
				return nil
			},
		}, "token")
		handler := handlers.NewDeleteClusterGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ClusterGroupDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteClusterGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ClusterGroupDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteClusterGroupFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteClusterGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ClusterGroupDeleteInput{ID: 3})
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

func TestCreateClusterGroupHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateClusterGroupHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ClusterGroupCreateInput{})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateClusterGroupHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateClusterGroupHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ClusterGroupUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteClusterGroupHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteClusterGroupHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ClusterGroupDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}
