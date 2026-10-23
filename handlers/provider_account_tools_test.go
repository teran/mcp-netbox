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
func TestCreateProviderAccountHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateProviderAccountFunc: func(
				_ context.Context, _ string, in domain.ProviderAccountWrite,
			) (*domain.ProviderAccount, error) {
				return &domain.ProviderAccount{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateProviderAccountHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.ProviderAccountCreateInput{Name: "Acc A"})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "Acc A" {
			t.Errorf("out.Data = %+v, want id 1 name Acc A", out.Data)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateProviderAccountHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ProviderAccountCreateInput{Name: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateProviderAccountFunc: func(
				_ context.Context, _ string, _ domain.ProviderAccountWrite,
			) (*domain.ProviderAccount, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateProviderAccountHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ProviderAccountCreateInput{Name: "A"})
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
			CreateProviderAccountFunc: func(
				_ context.Context, _ string, _ domain.ProviderAccountWrite,
			) (*domain.ProviderAccount, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateProviderAccountHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ProviderAccountCreateInput{Name: "A"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateProviderAccountHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateProviderAccountFunc: func(
				_ context.Context, _ string, id int, _ domain.ProviderAccountWrite,
			) (*domain.ProviderAccount, error) {
				return &domain.ProviderAccount{ID: id, Name: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateProviderAccountHandler(svc)
		name := "Renamed"
		result, out, err := handler(context.Background(), nil, handlers.ProviderAccountUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateProviderAccountHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ProviderAccountUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"account":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateProviderAccountFunc: func(
				_ context.Context, _ string, _ int, _ domain.ProviderAccountWrite,
			) (*domain.ProviderAccount, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateProviderAccountHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ProviderAccountUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteProviderAccountHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteProviderAccountFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteProviderAccountHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ProviderAccountDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteProviderAccountHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ProviderAccountDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteProviderAccountFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteProviderAccountHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ProviderAccountDeleteInput{ID: 3})
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

func TestCreateProviderAccountHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateProviderAccountHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ProviderAccountCreateInput{Name: "A"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateProviderAccountHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateProviderAccountHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ProviderAccountUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteProviderAccountHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteProviderAccountHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ProviderAccountDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetProviderAccountsHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListProviderAccountsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ProviderAccount], error) {
				return &domain.PaginatedResponse[domain.ProviderAccount]{
					Count:   1,
					Results: []domain.ProviderAccount{{ID: 1, Name: "Acc A"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetProviderAccountsHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.ProviderAccountsInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].Name != "Acc A" {
			t.Errorf("Results = %+v, want single provider account Acc A", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListProviderAccountsFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.ProviderAccount], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.ProviderAccount]{Count: 0, Results: []domain.ProviderAccount{}}, nil
			},
		}, "token")

		handler := handlers.NewGetProviderAccountsHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.ProviderAccountsInput{
			Q:        "acc",
			Provider: "7",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "acc" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "acc")
		}
		if gotParams["provider"] != "7" {
			t.Errorf("provider param = %q, want %q", gotParams["provider"], "7")
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
			ListProviderAccountsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ProviderAccount], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetProviderAccountsHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ProviderAccountsInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

// TestProviderAccountMapsWrite verifies the optional writable fields are mapped
// through the create handler to the domain write DTO.
func TestProviderAccountMapsWrite(t *testing.T) {
	t.Parallel()

	var gotWrite domain.ProviderAccountWrite
	svc := application.NewNetworkService(&mockrepo.MockRepo{
		CreateProviderAccountFunc: func(
			_ context.Context, _ string, in domain.ProviderAccountWrite,
		) (*domain.ProviderAccount, error) {
			gotWrite = in

			return &domain.ProviderAccount{ID: 1, Name: in.Name}, nil
		},
	}, "token")

	provider := 42
	account := "acc-1"
	desc := "d"
	handler := handlers.NewCreateProviderAccountHandler(svc)
	result, _, err := handler(context.Background(), nil, handlers.ProviderAccountCreateInput{
		Name: "A", Provider: &provider, Account: &account, Description: &desc,
	})
	if err != nil || result.IsError {
		t.Fatalf("handler failed: result.IsError=%v err=%v", result.IsError, err)
	}
	if gotWrite.Provider == nil || *gotWrite.Provider != 42 {
		t.Errorf("create mapped Provider = %v, want 42", gotWrite.Provider)
	}
	if gotWrite.Account == nil || *gotWrite.Account != "acc-1" {
		t.Errorf("create mapped Account = %v, want acc-1", gotWrite.Account)
	}
	if gotWrite.Description == nil || *gotWrite.Description != "d" {
		t.Errorf("create mapped Description = %v, want d", gotWrite.Description)
	}
}
