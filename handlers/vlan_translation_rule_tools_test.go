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
func TestCreateVLANTranslationRuleHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.VLANTranslationRuleWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVLANTranslationRuleFunc: func(
				_ context.Context, _ string, in domain.VLANTranslationRuleWrite,
			) (*domain.VLANTranslationRule, error) {
				gotWrite = in

				return &domain.VLANTranslationRule{ID: 1}, nil
			},
		}, "token")

		handler := handlers.NewCreateVLANTranslationRuleHandler(svc)
		localVID := 100
		remoteVID := 200
		result, out, err := handler(context.Background(), nil, handlers.VLANTranslationRuleCreateInput{
			Policy: 5, LocalVID: &localVID, RemoteVID: &remoteVID,
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
		if gotWrite.Policy != 5 || gotWrite.LocalVID == nil || *gotWrite.LocalVID != 100 {
			t.Errorf("create mapped Policy/LocalVID = %d/%v, want 5/100", gotWrite.Policy, gotWrite.LocalVID)
		}
		if gotWrite.RemoteVID == nil || *gotWrite.RemoteVID != 200 {
			t.Errorf("create mapped RemoteVID = %v, want 200", gotWrite.RemoteVID)
		}
	})

	t.Run("policy required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateVLANTranslationRuleHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VLANTranslationRuleCreateInput{Policy: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing policy")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"policy":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVLANTranslationRuleFunc: func(
				_ context.Context, _ string, _ domain.VLANTranslationRuleWrite,
			) (*domain.VLANTranslationRule, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateVLANTranslationRuleHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VLANTranslationRuleCreateInput{Policy: 5})
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
			CreateVLANTranslationRuleFunc: func(
				_ context.Context, _ string, _ domain.VLANTranslationRuleWrite,
			) (*domain.VLANTranslationRule, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateVLANTranslationRuleHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VLANTranslationRuleCreateInput{Policy: 5})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateVLANTranslationRuleHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var gotWrite domain.VLANTranslationRuleWrite
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateVLANTranslationRuleFunc: func(
				_ context.Context, _ string, id int, in domain.VLANTranslationRuleWrite,
			) (*domain.VLANTranslationRule, error) {
				gotWrite = in

				return &domain.VLANTranslationRule{ID: id}, nil
			},
		}, "token")
		handler := handlers.NewUpdateVLANTranslationRuleHandler(svc)
		remoteVID := 300
		result, out, err := handler(context.Background(), nil, handlers.VLANTranslationRuleUpdateInput{
			ID: 7, RemoteVID: &remoteVID,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
		if gotWrite.RemoteVID == nil || *gotWrite.RemoteVID != 300 {
			t.Errorf("update mapped RemoteVID = %v, want 300", gotWrite.RemoteVID)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateVLANTranslationRuleHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VLANTranslationRuleUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"local_vid":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateVLANTranslationRuleFunc: func(
				_ context.Context, _ string, _ int, _ domain.VLANTranslationRuleWrite,
			) (*domain.VLANTranslationRule, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateVLANTranslationRuleHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VLANTranslationRuleUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteVLANTranslationRuleHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteVLANTranslationRuleFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteVLANTranslationRuleHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VLANTranslationRuleDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteVLANTranslationRuleHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VLANTranslationRuleDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteVLANTranslationRuleFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteVLANTranslationRuleHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VLANTranslationRuleDeleteInput{ID: 3})
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

func TestCreateVLANTranslationRuleHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateVLANTranslationRuleHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VLANTranslationRuleCreateInput{Policy: 5})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateVLANTranslationRuleHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateVLANTranslationRuleHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VLANTranslationRuleUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteVLANTranslationRuleHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteVLANTranslationRuleHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VLANTranslationRuleDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetVLANTranslationRulesHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListVLANTranslationRulesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VLANTranslationRule], error) {
				return &domain.PaginatedResponse[domain.VLANTranslationRule]{
					Count:   1,
					Results: []domain.VLANTranslationRule{{ID: 1}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetVLANTranslationRulesHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.VLANTranslationRulesInput{})
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
			t.Errorf("Results = %+v, want single translation rule", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListVLANTranslationRulesFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.VLANTranslationRule], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.VLANTranslationRule]{Count: 0, Results: []domain.VLANTranslationRule{}}, nil
			},
		}, "token")

		handler := handlers.NewGetVLANTranslationRulesHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.VLANTranslationRulesInput{
			Q:        "100",
			Policy:   "5",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "100" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "100")
		}
		if gotParams["policy"] != "5" {
			t.Errorf("policy param = %q, want %q", gotParams["policy"], "5")
		}
		if gotParams["offset"] != "50" || gotParams["limit"] != "50" {
			t.Errorf("offset/limit = %q/%q, want 50/50", gotParams["offset"], gotParams["limit"])
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListVLANTranslationRulesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VLANTranslationRule], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetVLANTranslationRulesHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VLANTranslationRulesInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}
