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
func TestCreateVirtualDiskHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVirtualDiskFunc: func(_ context.Context, _ string, in domain.VirtualDiskWrite) (*domain.VirtualDisk, error) {
				return &domain.VirtualDisk{ID: 1, Name: in.Name}, nil
			},
		}, "token")

		handler := handlers.NewCreateVirtualDiskHandler(svc)
		vm := 5
		result, out, err := handler(
			context.Background(), nil,
			handlers.VirtualDiskCreateInput{Name: "Data", VirtualMachine: &vm},
		)
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Data.ID != 1 || out.Data.Name != "Data" {
			t.Errorf("out.Data = %+v, want id 1 name Data", out.Data)
		}
	})

	t.Run("name required", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewCreateVirtualDiskHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualDiskCreateInput{Name: ""})
		if !result.IsError {
			t.Error("result.IsError = false, want true for missing name")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"name":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			CreateVirtualDiskFunc: func(_ context.Context, _ string, _ domain.VirtualDiskWrite) (*domain.VirtualDisk, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")

		handler := handlers.NewCreateVirtualDiskHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VirtualDiskCreateInput{Name: "A"})
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
			CreateVirtualDiskFunc: func(_ context.Context, _ string, _ domain.VirtualDiskWrite) (*domain.VirtualDisk, error) {
				return nil, &domain.ValidationError{StatusCode: 500, Body: []byte(`{}`)}
			},
		}, "token")
		handler := handlers.NewCreateVirtualDiskHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualDiskCreateInput{Name: "A"})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestUpdateVirtualDiskHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateVirtualDiskFunc: func(
				_ context.Context, _ string, id int, _ domain.VirtualDiskWrite,
			) (*domain.VirtualDisk, error) {
				return &domain.VirtualDisk{ID: id, Name: "Renamed"}, nil
			},
		}, "token")
		handler := handlers.NewUpdateVirtualDiskHandler(svc)
		name := "Renamed"
		result, out, err := handler(context.Background(), nil, handlers.VirtualDiskUpdateInput{ID: 7, Name: &name})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError || out.Data.ID != 7 {
			t.Errorf("result = %+v, out.Data = %+v", result, out.Data)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewUpdateVirtualDiskHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualDiskUpdateInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error in structured content", func(t *testing.T) {
		body := []byte(`{"virtual_machine":["This field is required."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			UpdateVirtualDiskFunc: func(
				_ context.Context, _ string, _ int, _ domain.VirtualDiskWrite,
			) (*domain.VirtualDisk, error) {
				return nil, &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewUpdateVirtualDiskHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VirtualDiskUpdateInput{ID: 7})
		if err != nil {
			t.Fatalf("handler returned non-nil error, want structured result: %v", err)
		}
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

func TestDeleteVirtualDiskHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteVirtualDiskFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}, "token")
		handler := handlers.NewDeleteVirtualDiskHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VirtualDiskDeleteInput{ID: 3})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		handler := handlers.NewDeleteVirtualDiskHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualDiskDeleteInput{ID: 0})
		if !result.IsError {
			t.Error("result.IsError = false, want true for invalid id")
		}
	})

	t.Run("validation error propagated", func(t *testing.T) {
		body := []byte(`{"detail":["Cannot delete object with dependent objects."]}`)
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			DeleteVirtualDiskFunc: func(_ context.Context, _ string, _ int) error {
				return &domain.ValidationError{StatusCode: 400, Body: body}
			},
		}, "token")
		handler := handlers.NewDeleteVirtualDiskHandler(svc)
		result, _, err := handler(context.Background(), nil, handlers.VirtualDiskDeleteInput{ID: 3})
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

func TestCreateVirtualDiskHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewCreateVirtualDiskHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VirtualDiskCreateInput{Name: "A"})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestUpdateVirtualDiskHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewUpdateVirtualDiskHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VirtualDiskUpdateInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

func TestDeleteVirtualDiskHandler_NoService(t *testing.T) {
	t.Parallel()
	handler := handlers.NewDeleteVirtualDiskHandler(nil)
	result, _, err := handler(context.Background(), nil, handlers.VirtualDiskDeleteInput{ID: 1})
	if err == nil || !result.IsError {
		t.Errorf("result.IsError = %v, err = %v; want error when no service", result.IsError, err)
	}
}

//nolint:funlen // Table of subtests covering the list handler's branches; splitting is not sensible.
func TestGetVirtualDisksHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListVirtualDisksFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualDisk], error) {
				return &domain.PaginatedResponse[domain.VirtualDisk]{
					Count:   1,
					Results: []domain.VirtualDisk{{ID: 1, Name: "Data"}},
				}, nil
			},
		}, "token")

		handler := handlers.NewGetVirtualDisksHandler(svc)
		result, out, err := handler(context.Background(), nil, handlers.VirtualDisksInput{})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if result.IsError {
			t.Error("result.IsError = true, want false")
		}
		if out.Count != 1 {
			t.Errorf("Count = %d, want %d", out.Count, 1)
		}
		if len(out.Results) != 1 || out.Results[0].Name != "Data" {
			t.Errorf("Results = %+v, want single virtual disk Data", out.Results)
		}
	})

	t.Run("filter params propagated", func(t *testing.T) {
		var gotParams map[string]string
		svc := application.NewNetworkService(&mockrepo.MockRepo{
			ListVirtualDisksFunc: func(
				_ context.Context, _ string, params map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualDisk], error) {
				gotParams = params

				return &domain.PaginatedResponse[domain.VirtualDisk]{Count: 0, Results: []domain.VirtualDisk{}}, nil
			},
		}, "token")

		handler := handlers.NewGetVirtualDisksHandler(svc)
		_, _, err := handler(context.Background(), nil, handlers.VirtualDisksInput{
			Q:              "data",
			VirtualMachine: "7",
			Tag:            "prod",
			Page:           2,
			PageSize:       50,
		})
		if err != nil {
			t.Fatalf("handler returned error: %v", err)
		}
		if gotParams["q"] != "data" {
			t.Errorf("q param = %q, want %q", gotParams["q"], "data")
		}
		if gotParams["virtual_machine"] != "7" {
			t.Errorf("virtual_machine param = %q, want %q", gotParams["virtual_machine"], "7")
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
			ListVirtualDisksFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualDisk], error) {
				return nil, errors.New("repo error")
			},
		}, "token")

		handler := handlers.NewGetVirtualDisksHandler(svc)
		result, _, _ := handler(context.Background(), nil, handlers.VirtualDisksInput{})
		if !result.IsError {
			t.Error("result.IsError = false, want true")
		}
	})
}

// TestVirtualDiskMapsReferences verifies the numeric virtual_machine and size
// fields are forwarded through the create handler to the write DTO.
func TestVirtualDiskMapsReferences(t *testing.T) {
	t.Parallel()

	var gotWrite domain.VirtualDiskWrite
	svc := application.NewNetworkService(&mockrepo.MockRepo{
		CreateVirtualDiskFunc: func(_ context.Context, _ string, in domain.VirtualDiskWrite) (*domain.VirtualDisk, error) {
			gotWrite = in

			return &domain.VirtualDisk{ID: 1, Name: in.Name}, nil
		},
	}, "token")

	vm := 5
	size := 100
	handler := handlers.NewCreateVirtualDiskHandler(svc)
	result, _, err := handler(context.Background(), nil, handlers.VirtualDiskCreateInput{
		Name: "Data", VirtualMachine: &vm, Size: &size,
	})
	if err != nil || result.IsError {
		t.Fatalf("handler failed: result.IsError=%v err=%v", result.IsError, err)
	}
	if gotWrite.VirtualMachine == nil || *gotWrite.VirtualMachine != 5 {
		t.Errorf("create mapped VirtualMachine = %v, want 5", gotWrite.VirtualMachine)
	}
	if gotWrite.Size == nil || *gotWrite.Size != 100 {
		t.Errorf("create mapped Size = %v, want 100", gotWrite.Size)
	}
}
