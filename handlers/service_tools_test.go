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
func TestCreateServiceHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.ServiceWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateServiceFunc: func(_ context.Context, _ string, in domain.ServiceWrite) (*domain.Service, error) {
				gotWrite = in

				return &domain.Service{ID: 1, Name: in.Name, Ports: in.Ports}, nil
			},
		}, "token")

		handler := handlers.NewCreateServiceHandler(svc)
		parentType := "dcim.device"
		parentID := 5
		protocol := "tcp"
		result, out, err := handler(context.Background(), nil, handlers.ServiceCreateInput{
			Name: "web", ParentObjectType: &parentType, ParentObjectID: &parentID,
			Protocol: &protocol, Ports: []int{80, 443},
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "web" {
			t.Errorf("out.Data = %+v, want id 1 name web", out.Data)
		}
		if gotWrite.ParentObjectID == nil || *gotWrite.ParentObjectID != 5 {
			t.Errorf("create mapped ParentObjectID = %v, want 5", gotWrite.ParentObjectID)
		}
		if len(gotWrite.Ports) != 2 || gotWrite.Ports[0] != 80 || gotWrite.Ports[1] != 443 {
			t.Errorf("create mapped Ports = %v, want [80 443]", gotWrite.Ports)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateServiceHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ServiceCreateInput{Name: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"ports":["Enter a valid value."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateServiceFunc: func(_ context.Context, _ string, _ domain.ServiceWrite) (*domain.Service, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateServiceHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ServiceCreateInput{Name: "web", Ports: []int{70000}})
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
			CreateServiceFunc: func(_ context.Context, _ string, _ domain.ServiceWrite) (*domain.Service, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateServiceHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ServiceCreateInput{Name: "web"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateServiceHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.ServiceWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateServiceFunc: func(
				_ context.Context, _ string, id int, in domain.ServiceWrite,
			) (*domain.Service, error) {
				gotWrite = in

				return &domain.Service{ID: id, Name: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateServiceHandler(svc)
		parentID := 8
		parentType := "dcim.device"
		result, out, err := handler(context.Background(), nil, handlers.ServiceUpdateInput{
			ID: 7, ParentObjectID: &parentID, ParentObjectType: &parentType,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
		if gotWrite.ParentObjectID == nil || *gotWrite.ParentObjectID != 8 ||
			gotWrite.ParentObjectType == nil || *gotWrite.ParentObjectType != "dcim.device" {
			t.Errorf("update mapped ParentObject = %v/%v, want 8/dcim.device",
				gotWrite.ParentObjectID, gotWrite.ParentObjectType)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateServiceHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ServiceUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateServiceFunc: func(
				_ context.Context, _ string, _ int, _ domain.ServiceWrite,
			) (*domain.Service, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateServiceHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ServiceUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteServiceHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteServiceFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteServiceHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ServiceDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteServiceHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ServiceDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteServiceFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteServiceHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ServiceDeleteInput{ID: 3})
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

func TestCreateServiceHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateServiceHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ServiceCreateInput{Name: "web"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateServiceHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateServiceHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ServiceUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteServiceHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteServiceHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ServiceDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetServicesHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListServicesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.Service], error) {
				return &domain.PaginatedResponse[domain.Service]{
					Count:   1,
					Results: []domain.Service{{ID: 1, Name: "web"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetServicesHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.ServicesInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].Name != "web" {
			t.Errorf("Results = %+v, want single service web", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListServicesFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.Service], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.Service]{Count: 0, Results: []domain.Service{}}, nil
			},
		}, "token")

		handler := handlers.NewGetServicesHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.ServicesInput{
			Q:              "web",
			Device:         "sw1",
			VirtualMachine: "vm1",
			Tag:            "prod",
			Page:           2,
			PageSize:       50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "web" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "web")
		}
		if gotParams["device"] != "sw1" {
			t.Errorf("device param = %q, want %q", gotParams["device"], "sw1")
		}
		if gotParams["virtual_machine"] != "vm1" {
			t.Errorf("virtual_machine param = %q, want %q", gotParams["virtual_machine"], "vm1")
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
			ListServicesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.Service], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetServicesHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ServicesInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
