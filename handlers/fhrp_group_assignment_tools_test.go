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
func TestCreateFHRPGroupAssignmentHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.FHRPGroupAssignmentWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateFHRPGroupAssignmentFunc: func(
				_ context.Context, _ string, in domain.FHRPGroupAssignmentWrite,
			) (*domain.FHRPGroupAssignment, error) {
				gotWrite = in

				return &domain.FHRPGroupAssignment{ID: 1, Group: &domain.Nested{ID: in.Group}}, nil
			},
		}, "token")

		handler := handlers.NewCreateFHRPGroupAssignmentHandler(svc)
		priority := 100
		result, out, err := handler(context.Background(), nil, handlers.FHRPGroupAssignmentCreateInput{
			Group: 5, InterfaceType: "dcim.interface", InterfaceID: 9, Priority: &priority,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 {
			t.Errorf("out.Data = %+v, want id 1", out.Data)
		}
		if gotWrite.Group != 5 || gotWrite.InterfaceID != 9 {
			t.Errorf("create mapped Group/InterfaceID = %d/%d, want 5/9", gotWrite.Group, gotWrite.InterfaceID)
		}
		if gotWrite.InterfaceType != "dcim.interface" {
			t.Errorf("create mapped InterfaceType = %q, want dcim.interface", gotWrite.InterfaceType)
		}
		if gotWrite.Priority == nil || *gotWrite.Priority != 100 {
			t.Errorf("create mapped Priority = %v, want 100", gotWrite.Priority)
		}
	})

	t.Run("group required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateFHRPGroupAssignmentHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.FHRPGroupAssignmentCreateInput{
			Group: 0, InterfaceType: "dcim.interface", InterfaceID: 9,
		})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing group")
		}
	})

	t.Run("interface_id required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateFHRPGroupAssignmentHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.FHRPGroupAssignmentCreateInput{
			Group: 5, InterfaceType: "dcim.interface", InterfaceID: 0,
		})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing interface_id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"group":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateFHRPGroupAssignmentFunc: func(
				_ context.Context, _ string, _ domain.FHRPGroupAssignmentWrite,
			) (*domain.FHRPGroupAssignment, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateFHRPGroupAssignmentHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.FHRPGroupAssignmentCreateInput{
			Group: 5, InterfaceType: "dcim.interface", InterfaceID: 9,
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
			CreateFHRPGroupAssignmentFunc: func(
				_ context.Context, _ string, _ domain.FHRPGroupAssignmentWrite,
			) (*domain.FHRPGroupAssignment, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateFHRPGroupAssignmentHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.FHRPGroupAssignmentCreateInput{
			Group: 5, InterfaceType: "dcim.interface", InterfaceID: 9,
		})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateFHRPGroupAssignmentHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.FHRPGroupAssignmentWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateFHRPGroupAssignmentFunc: func(
				_ context.Context, _ string, id int, in domain.FHRPGroupAssignmentWrite,
			) (*domain.FHRPGroupAssignment, error) {
				gotWrite = in

				return &domain.FHRPGroupAssignment{ID: id}, nil
			},
		}, "token")
		handler := handlers.NewUpdateFHRPGroupAssignmentHandler(svc)
		priority := 200
		result, out, err := handler(context.Background(), nil, handlers.FHRPGroupAssignmentUpdateInput{
			ID: 7, Priority: &priority,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
		if gotWrite.Priority == nil || *gotWrite.Priority != 200 {
			t.Errorf("update mapped Priority = %v, want 200", gotWrite.Priority)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateFHRPGroupAssignmentHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.FHRPGroupAssignmentUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"ip_address":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateFHRPGroupAssignmentFunc: func(
				_ context.Context, _ string, _ int, _ domain.FHRPGroupAssignmentWrite,
			) (*domain.FHRPGroupAssignment, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateFHRPGroupAssignmentHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.FHRPGroupAssignmentUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteFHRPGroupAssignmentHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteFHRPGroupAssignmentFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteFHRPGroupAssignmentHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.FHRPGroupAssignmentDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteFHRPGroupAssignmentHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.FHRPGroupAssignmentDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteFHRPGroupAssignmentFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteFHRPGroupAssignmentHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.FHRPGroupAssignmentDeleteInput{ID: 3})
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

func TestCreateFHRPGroupAssignmentHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateFHRPGroupAssignmentHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.FHRPGroupAssignmentCreateInput{
		Group: 5, InterfaceType: "dcim.interface", InterfaceID: 9,
	})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateFHRPGroupAssignmentHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateFHRPGroupAssignmentHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.FHRPGroupAssignmentUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteFHRPGroupAssignmentHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteFHRPGroupAssignmentHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.FHRPGroupAssignmentDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetFHRPGroupAssignmentsHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListFHRPGroupAssignmentsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.FHRPGroupAssignment], error) {
				return &domain.PaginatedResponse[domain.FHRPGroupAssignment]{
					Count:   1,
					Results: []domain.FHRPGroupAssignment{{ID: 1}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetFHRPGroupAssignmentsHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.FHRPGroupAssignmentsInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].ID != 1 {
			t.Errorf("Results = %+v, want single FHRP group assignment", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListFHRPGroupAssignmentsFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.FHRPGroupAssignment], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.FHRPGroupAssignment]{Count: 0, Results: []domain.FHRPGroupAssignment{}}, nil
			},
		}, "token")

		handler := handlers.NewGetFHRPGroupAssignmentsHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.FHRPGroupAssignmentsInput{
			Q:         "10.0.0.1",
			Group:     "5",
			IPAddress: "9",
			Page:      2,
			PageSize:  50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "10.0.0.1" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "10.0.0.1")
		}
		if gotParams["group"] != "5" {
			t.Errorf("group param = %q, want %q", gotParams["group"], "5")
		}
		if gotParams["ip_address"] != "9" {
			t.Errorf("ip_address param = %q, want %q", gotParams["ip_address"], "9")
		}
		if gotParams["offset"] != "50" || gotParams["limit"] != "50" {
			t.Errorf("offset/limit = %q/%q, want 50/50", gotParams["offset"], gotParams["limit"])
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListFHRPGroupAssignmentsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.FHRPGroupAssignment], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetFHRPGroupAssignmentsHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.FHRPGroupAssignmentsInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
