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
func TestCreateRackReservationHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateRackReservationFunc: func(
				_ context.Context, _ string, in domain.RackReservationWrite,
			) (*domain.RackReservation, error) {
				return &domain.RackReservation{ID: 1, Description: in.Description}, nil
			},
		}, "token")

		handler := handlers.NewCreateRackReservationHandler(svc)
		rack := 3
		result, out, err := handler(context.Background(), nil, handlers.RackReservationCreateInput{
			Description: "reserved", Rack: &rack,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Description != "reserved" {
			t.Errorf("out.Data = %+v, want id 1 description reserved", out.Data)
		}
	})

	t.Run("description required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateRackReservationHandler(svc)
		rack := 3
		result, _, _ := handler(context.Background(), nil, handlers.RackReservationCreateInput{
			Description: "", Rack: &rack,
		})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing description")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"description":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateRackReservationFunc: func(
				_ context.Context, _ string, _ domain.RackReservationWrite,
			) (*domain.RackReservation, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateRackReservationHandler(svc)
		rack := 3
		result, _, err := handler(context.Background(), nil, handlers.RackReservationCreateInput{
			Description: "reserved", Rack: &rack,
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
			CreateRackReservationFunc: func(
				_ context.Context, _ string, _ domain.RackReservationWrite,
			) (*domain.RackReservation, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateRackReservationHandler(svc)
		rack := 3
		result, _, _ := handler(context.Background(), nil, handlers.RackReservationCreateInput{
			Description: "reserved", Rack: &rack,
		})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateRackReservationHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateRackReservationFunc: func(
				_ context.Context, _ string, id int, _ domain.RackReservationWrite,
			) (*domain.RackReservation, error) {
				return &domain.RackReservation{ID: id, Description: "Updated"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateRackReservationHandler(svc)
		desc := "Updated"
		result, out, err := handler(context.Background(), nil, handlers.RackReservationUpdateInput{
			ID: 7, Description: &desc,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateRackReservationHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.RackReservationUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"rack":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateRackReservationFunc: func(
				_ context.Context, _ string, _ int, _ domain.RackReservationWrite,
			) (*domain.RackReservation, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateRackReservationHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.RackReservationUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteRackReservationHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteRackReservationFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteRackReservationHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.RackReservationDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteRackReservationHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.RackReservationDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteRackReservationFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteRackReservationHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.RackReservationDeleteInput{ID: 3})
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

func TestCreateRackReservationHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateRackReservationHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.RackReservationCreateInput{Description: "reserved"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateRackReservationHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateRackReservationHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.RackReservationUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteRackReservationHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteRackReservationHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.RackReservationDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetRackReservationsHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListRackReservationsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.RackReservation], error) {
				return &domain.PaginatedResponse[domain.RackReservation]{
					Count:   1,
					Results: []domain.RackReservation{{ID: 1, Description: "reserved"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetRackReservationsHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.RackReservationsInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 || len(out.Results) != 1 || out.Results[0].Description != "reserved" {
			t.Errorf("out = %+v, want count 1 with reserved", out)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListRackReservationsFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.RackReservation], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.RackReservation]{Count: 0, Results: []domain.RackReservation{}}, nil
			},
		}, "token")

		handler := handlers.NewGetRackReservationsHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.RackReservationsInput{
			Q:        "res",
			Rack:     "r1",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "res" || gotParams["rack"] != "r1" || gotParams["tag"] != "prod" {
			t.Errorf("params = %+v, want q res, rack r1, tag prod", gotParams)
		}
		if gotParams["offset"] != "50" || gotParams["limit"] != "50" {
			t.Errorf("pagination = %+v, want offset 50 limit 50", gotParams)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListRackReservationsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.RackReservation], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetRackReservationsHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.RackReservationsInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
