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
func TestCreateCircuitGroupAssignmentHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.CircuitGroupAssignmentWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateCircuitGroupAssignmentFunc: func(
				_ context.Context, _ string, in domain.CircuitGroupAssignmentWrite,
			) (*domain.CircuitGroupAssignment, error) {
				gotWrite = in

				return &domain.CircuitGroupAssignment{ID: 1, Group: &domain.Nested{ID: in.Group}}, nil
			},
		}, "token")

		handler := handlers.NewCreateCircuitGroupAssignmentHandler(svc)
		priority := "primary"
		result, out, err := handler(context.Background(), nil, handlers.CircuitGroupAssignmentCreateInput{
			Group: 5, MemberType: "circuits.circuit", MemberID: 9, Priority: &priority,
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
		if gotWrite.Group != 5 || gotWrite.MemberID != 9 {
			t.Errorf("create mapped Group/MemberID = %d/%d, want 5/9", gotWrite.Group, gotWrite.MemberID)
		}
		if gotWrite.Priority == nil || *gotWrite.Priority != "primary" {
			t.Errorf("create mapped Priority = %v, want primary", gotWrite.Priority)
		}
	})

	t.Run("group required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateCircuitGroupAssignmentHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CircuitGroupAssignmentCreateInput{
			Group: 0, MemberType: "circuits.circuit", MemberID: 9,
		})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing group")
		}
	})

	t.Run("member required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateCircuitGroupAssignmentHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CircuitGroupAssignmentCreateInput{
			Group: 5, MemberType: "", MemberID: 0,
		})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing member")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"group":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateCircuitGroupAssignmentFunc: func(
				_ context.Context, _ string, _ domain.CircuitGroupAssignmentWrite,
			) (*domain.CircuitGroupAssignment, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateCircuitGroupAssignmentHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CircuitGroupAssignmentCreateInput{
			Group: 5, MemberType: "circuits.circuit", MemberID: 9,
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
			CreateCircuitGroupAssignmentFunc: func(
				_ context.Context, _ string, _ domain.CircuitGroupAssignmentWrite,
			) (*domain.CircuitGroupAssignment, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateCircuitGroupAssignmentHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CircuitGroupAssignmentCreateInput{
			Group: 5, MemberType: "circuits.circuit", MemberID: 9,
		})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

//nolint:funlen // Table of subtests covering update branches (pointer-set + validation); splitting is not sensible.
func TestUpdateCircuitGroupAssignmentHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.CircuitGroupAssignmentWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateCircuitGroupAssignmentFunc: func(
				_ context.Context, _ string, id int, in domain.CircuitGroupAssignmentWrite,
			) (*domain.CircuitGroupAssignment, error) {
				gotWrite = in

				return &domain.CircuitGroupAssignment{ID: id}, nil
			},
		}, "token")
		handler := handlers.NewUpdateCircuitGroupAssignmentHandler(svc)
		priority := "secondary"
		result, out, err := handler(context.Background(), nil, handlers.CircuitGroupAssignmentUpdateInput{
			ID: 7, Priority: &priority,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
		if gotWrite.Priority == nil || *gotWrite.Priority != "secondary" {
			t.Errorf("update mapped Priority = %v, want secondary", gotWrite.Priority)
		}
	})

	t.Run("success with group and member set", func(t *testing.T) {
		var gotWrite domain.CircuitGroupAssignmentWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateCircuitGroupAssignmentFunc: func(
				_ context.Context, _ string, id int, in domain.CircuitGroupAssignmentWrite,
			) (*domain.CircuitGroupAssignment, error) {
				gotWrite = in

				return &domain.CircuitGroupAssignment{ID: id}, nil
			},
		}, "token")
		handler := handlers.NewUpdateCircuitGroupAssignmentHandler(svc)
		group := 5
		memberType := "circuits.circuit"
		memberID := 6
		result, out, err := handler(context.Background(), nil, handlers.CircuitGroupAssignmentUpdateInput{
			ID: 7, Group: &group, MemberType: &memberType, MemberID: &memberID,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
		if gotWrite.Group != 5 || gotWrite.MemberID != 6 {
			t.Errorf("update mapped group/member = %d/%d, want 5/6", gotWrite.Group, gotWrite.MemberID)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateCircuitGroupAssignmentHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CircuitGroupAssignmentUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"circuit":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateCircuitGroupAssignmentFunc: func(
				_ context.Context, _ string, _ int, _ domain.CircuitGroupAssignmentWrite,
			) (*domain.CircuitGroupAssignment, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateCircuitGroupAssignmentHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CircuitGroupAssignmentUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteCircuitGroupAssignmentHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteCircuitGroupAssignmentFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteCircuitGroupAssignmentHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CircuitGroupAssignmentDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteCircuitGroupAssignmentHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CircuitGroupAssignmentDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteCircuitGroupAssignmentFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteCircuitGroupAssignmentHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.CircuitGroupAssignmentDeleteInput{ID: 3})
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

func TestCreateCircuitGroupAssignmentHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateCircuitGroupAssignmentHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.CircuitGroupAssignmentCreateInput{
		Group: 5, MemberType: "circuits.circuit", MemberID: 9,
	})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateCircuitGroupAssignmentHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateCircuitGroupAssignmentHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.CircuitGroupAssignmentUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteCircuitGroupAssignmentHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteCircuitGroupAssignmentHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.CircuitGroupAssignmentDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetCircuitGroupAssignmentsHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListCircuitGroupAssignmentsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.CircuitGroupAssignment], error) {
				return &domain.PaginatedResponse[domain.CircuitGroupAssignment]{
					Count:   1,
					Results: []domain.CircuitGroupAssignment{{ID: 1}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetCircuitGroupAssignmentsHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.CircuitGroupAssignmentsInput{})
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
			t.Errorf("Results = %+v, want single circuit group assignment", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListCircuitGroupAssignmentsFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.CircuitGroupAssignment], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.CircuitGroupAssignment]{
					Count:   0,
					Results: []domain.CircuitGroupAssignment{},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetCircuitGroupAssignmentsHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.CircuitGroupAssignmentsInput{
			Q:        "10.0.0.1",
			Group:    "5",
			Circuit:  "9",
			Page:     2,
			PageSize: 50,
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
		if gotParams["circuit"] != "9" {
			t.Errorf("circuit param = %q, want %q", gotParams["circuit"], "9")
		}
		if gotParams["offset"] != "50" || gotParams["limit"] != "50" {
			t.Errorf("offset/limit = %q/%q, want 50/50", gotParams["offset"], gotParams["limit"])
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListCircuitGroupAssignmentsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.CircuitGroupAssignment], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetCircuitGroupAssignmentsHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.CircuitGroupAssignmentsInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
