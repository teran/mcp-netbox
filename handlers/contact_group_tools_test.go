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
func TestCreateContactGroupHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateContactGroupFunc: func(
				_ context.Context, _ string, in domain.ContactGroupWrite,
			) (*domain.ContactGroup, error) {
				return &domain.ContactGroup{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateContactGroupHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.ContactGroupCreateInput{Name: "Admins"})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "Admins" {
			t.Errorf("out.Data = %+v, want id 1 name Admins", out.Data)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateContactGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ContactGroupCreateInput{Name: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateContactGroupFunc: func(_ context.Context, _ string, _ domain.ContactGroupWrite) (*domain.ContactGroup, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateContactGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ContactGroupCreateInput{Name: "A"})
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
			CreateContactGroupFunc: func(_ context.Context, _ string, _ domain.ContactGroupWrite) (*domain.ContactGroup, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateContactGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ContactGroupCreateInput{Name: "A"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateContactGroupHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateContactGroupFunc: func(
				_ context.Context, _ string, id int, _ domain.ContactGroupWrite,
			) (*domain.ContactGroup, error) {
				return &domain.ContactGroup{ID: id, Name: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateContactGroupHandler(svc)
		name := "Renamed"
		result, out, err := handler(context.Background(), nil, handlers.ContactGroupUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateContactGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ContactGroupUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"slug":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateContactGroupFunc: func(
				_ context.Context, _ string, _ int, _ domain.ContactGroupWrite,
			) (*domain.ContactGroup, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateContactGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ContactGroupUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteContactGroupHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteContactGroupFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteContactGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ContactGroupDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteContactGroupHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ContactGroupDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteContactGroupFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteContactGroupHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.ContactGroupDeleteInput{ID: 3})
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

func TestCreateContactGroupHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateContactGroupHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ContactGroupCreateInput{Name: "A"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateContactGroupHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateContactGroupHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ContactGroupUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteContactGroupHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteContactGroupHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.ContactGroupDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetContactGroupsHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListContactGroupsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ContactGroup], error) {
				return &domain.PaginatedResponse[domain.ContactGroup]{
					Count:   1,
					Results: []domain.ContactGroup{{ID: 1, Name: "Admins"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetContactGroupsHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.ContactGroupsInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].Name != "Admins" {
			t.Errorf("Results = %+v, want single contact group Admins", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListContactGroupsFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.ContactGroup], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.ContactGroup]{Count: 0, Results: []domain.ContactGroup{}}, nil
			},
		}, "token")

		handler := handlers.NewGetContactGroupsHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.ContactGroupsInput{
			Q:        "adm",
			Parent:   "7",
			Tag:      "prod",
			Page:     2,
			PageSize: 50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "adm" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "adm")
		}
		if gotParams["parent"] != "7" {
			t.Errorf("parent param = %q, want %q", gotParams["parent"], "7")
		}
		if gotParams["tag"] != "prod" {
			t.Errorf("tag param = %q, want %q", gotParams["tag"], "prod")
		}
		if gotParams["offset"] != "50" {
			t.Errorf("offset param = %q, want %q", gotParams["offset"], "50")
		}
		if gotParams["limit"] != "50" {
			t.Errorf("limit param = %q, want %q", gotParams["limit"], "50")
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListContactGroupsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ContactGroup], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetContactGroupsHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.ContactGroupsInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

// TestContactGroupMapsParent verifies the numeric parent contact-group ID is
// forwarded through the create handler to the write DTO.
func TestContactGroupMapsParent(t *testing.T) {
	t.Parallel()

	var gotWrite domain.ContactGroupWrite
	svc := application.NewNetworkService(&mockrepo.MockRepo{
		CreateContactGroupFunc: func(_ context.Context, _ string, in domain.ContactGroupWrite) (*domain.ContactGroup, error) {
			gotWrite = in

			return &domain.ContactGroup{ID: 1, Name: in.Name}, nil
		},
	}, "token")

	parent := 42
	handler := handlers.NewCreateContactGroupHandler(svc)
	result, _, err := handler(context.Background(), nil, handlers.ContactGroupCreateInput{Name: "Admins", Parent: &parent})
	if err != nil || result.IsError {
		t.Fatalf("handler failed: result.IsError=%v err=%v", result.IsError, err)
	}
	if gotWrite.Parent == nil || *gotWrite.Parent != 42 {
		t.Errorf("create mapped Parent = %v, want 42", gotWrite.Parent)
	}
}
