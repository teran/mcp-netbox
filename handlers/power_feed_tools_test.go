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
func TestCreatePowerFeedHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreatePowerFeedFunc: func(_ context.Context, _ string, in domain.PowerFeedWrite) (*domain.PowerFeed, error) {
				return &domain.PowerFeed{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreatePowerFeedHandler(svc)
		panel := 5
		result, out, err := handler(context.Background(), nil, handlers.PowerFeedCreateInput{
			Name: "PF1", PowerPanel: &panel,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "PF1" {
			t.Errorf("out.Data = %+v, want id 1 name PF1", out.Data)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreatePowerFeedHandler(svc)
		panel := 5
		result, _, _ := handler(context.Background(), nil, handlers.PowerFeedCreateInput{Name: "", PowerPanel: &panel})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreatePowerFeedFunc: func(_ context.Context, _ string, _ domain.PowerFeedWrite) (*domain.PowerFeed, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreatePowerFeedHandler(svc)
		panel := 5
		result, _, err := handler(context.Background(), nil, handlers.PowerFeedCreateInput{
			Name: "PF1", PowerPanel: &panel,
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
			CreatePowerFeedFunc: func(_ context.Context, _ string, _ domain.PowerFeedWrite) (*domain.PowerFeed, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreatePowerFeedHandler(svc)
		panel := 5
		result, _, _ := handler(context.Background(), nil, handlers.PowerFeedCreateInput{
			Name: "PF1", PowerPanel: &panel,
		})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdatePowerFeedHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdatePowerFeedFunc: func(_ context.Context, _ string, id int, _ domain.PowerFeedWrite) (*domain.PowerFeed, error) {
				return &domain.PowerFeed{ID: id, Name: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdatePowerFeedHandler(svc)
		name := "Renamed"
		result, out, err := handler(context.Background(), nil, handlers.PowerFeedUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdatePowerFeedHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.PowerFeedUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"power_panel":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdatePowerFeedFunc: func(_ context.Context, _ string, _ int, _ domain.PowerFeedWrite) (*domain.PowerFeed, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdatePowerFeedHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.PowerFeedUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeletePowerFeedHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeletePowerFeedFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeletePowerFeedHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.PowerFeedDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeletePowerFeedHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.PowerFeedDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeletePowerFeedFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeletePowerFeedHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.PowerFeedDeleteInput{ID: 3})
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

func TestCreatePowerFeedHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreatePowerFeedHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.PowerFeedCreateInput{Name: "PF1"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdatePowerFeedHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdatePowerFeedHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.PowerFeedUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeletePowerFeedHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeletePowerFeedHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.PowerFeedDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetPowerFeedsHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListPowerFeedsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.PowerFeed], error) {
				return &domain.PaginatedResponse[domain.PowerFeed]{
					Count:   1,
					Results: []domain.PowerFeed{{ID: 1, Name: "PF1"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetPowerFeedsHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.PowerFeedsInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 || len(out.Results) != 1 || out.Results[0].Name != "PF1" {
			t.Errorf("out = %+v, want count 1 with PF1", out)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListPowerFeedsFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.PowerFeed], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.PowerFeed]{Count: 0, Results: []domain.PowerFeed{}}, nil
			},
		}, "token")

		handler := handlers.NewGetPowerFeedsHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.PowerFeedsInput{
			Q:          "pf",
			PowerPanel: "pp1",
			Tag:        "prod",
			Page:       2,
			PageSize:   50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "pf" || gotParams["power_panel"] != "pp1" || gotParams["tag"] != "prod" {
			t.Errorf("params = %+v, want q pf, power_panel pp1, tag prod", gotParams)
		}
		if gotParams["offset"] != "50" || gotParams["limit"] != "50" {
			t.Errorf("pagination = %+v, want offset 50 limit 50", gotParams)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListPowerFeedsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.PowerFeed], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetPowerFeedsHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.PowerFeedsInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
