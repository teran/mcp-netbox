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
func TestCreateFHRPGroupHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.FHRPGroupWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateFHRPGroupFunc: func(_ context.Context, _ string, in domain.FHRPGroupWrite) (*domain.FHRPGroup, error) {
				gotWrite = in

				return &domain.FHRPGroup{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateFHRPGroupHandler(svc)
		protocol := "vrrp"
		groupID := 10
		result, out, err := handler(context.Background(), nil, handlers.FHRPGroupCreateInput{
			Name: "gw", Protocol: &protocol, GroupID: &groupID,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "gw" {
			t.Errorf("out.Data = %+v, want id 1 name gw", out.Data)
		}
		if gotWrite.GroupID == nil || *gotWrite.GroupID != 10 {
			t.Errorf("create mapped GroupID = %v, want 10", gotWrite.GroupID)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateFHRPGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.FHRPGroupCreateInput{Name: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateFHRPGroupFunc: func(_ context.Context, _ string, _ domain.FHRPGroupWrite) (*domain.FHRPGroup, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateFHRPGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.FHRPGroupCreateInput{Name: "gw"})
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
			CreateFHRPGroupFunc: func(_ context.Context, _ string, _ domain.FHRPGroupWrite) (*domain.FHRPGroup, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateFHRPGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.FHRPGroupCreateInput{Name: "gw"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateFHRPGroupHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.FHRPGroupWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateFHRPGroupFunc: func(
				_ context.Context, _ string, id int, in domain.FHRPGroupWrite,
			) (*domain.FHRPGroup, error) {
				gotWrite = in

				return &domain.FHRPGroup{ID: id, Name: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateFHRPGroupHandler(svc)
		groupID := 20
		result, out, err := handler(context.Background(), nil, handlers.FHRPGroupUpdateInput{ID: 7, GroupID: &groupID})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
		if gotWrite.GroupID == nil || *gotWrite.GroupID != 20 {
			t.Errorf("update mapped GroupID = %v, want 20", gotWrite.GroupID)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateFHRPGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.FHRPGroupUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"protocol":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateFHRPGroupFunc: func(
				_ context.Context, _ string, _ int, _ domain.FHRPGroupWrite,
			) (*domain.FHRPGroup, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateFHRPGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.FHRPGroupUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteFHRPGroupHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteFHRPGroupFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteFHRPGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.FHRPGroupDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteFHRPGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.FHRPGroupDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteFHRPGroupFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteFHRPGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.FHRPGroupDeleteInput{ID: 3})
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

func TestCreateFHRPGroupHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateFHRPGroupHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.FHRPGroupCreateInput{Name: "gw"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateFHRPGroupHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateFHRPGroupHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.FHRPGroupUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteFHRPGroupHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteFHRPGroupHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.FHRPGroupDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetFHRPGroupsHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListFHRPGroupsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.FHRPGroup], error) {
				return &domain.PaginatedResponse[domain.FHRPGroup]{
					Count:   1,
					Results: []domain.FHRPGroup{{ID: 1, Name: "gw"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetFHRPGroupsHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.FHRPGroupsInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].Name != "gw" {
			t.Errorf("Results = %+v, want single FHRP group gw", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListFHRPGroupsFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.FHRPGroup], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.FHRPGroup]{Count: 0, Results: []domain.FHRPGroup{}}, nil
			},
		}, "token")

		handler := handlers.NewGetFHRPGroupsHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.FHRPGroupsInput{
			Q:        "gw",
			Protocol: "vrrp",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "gw" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "gw")
		}
		if gotParams["protocol"] != "vrrp" {
			t.Errorf("protocol param = %q, want %q", gotParams["protocol"], "vrrp")
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
			ListFHRPGroupsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.FHRPGroup], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetFHRPGroupsHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.FHRPGroupsInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
