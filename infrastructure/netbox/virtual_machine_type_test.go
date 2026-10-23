package netbox

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/teran/mcp-netbox/domain"
)

// TestClient_ListVirtualMachineTypes exercises the virtual machine type list
// endpoint.
func TestClient_ListVirtualMachineTypes(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "Token test-token" {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}
			if req.URL.Path != "/api/virtualization/virtual-machine-types/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
				"count": 1,
				"next": null,
				"previous": null,
				"results": [{
					"id": 1,
					"name": "Standard",
					"slug": "standard",
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListVirtualMachineTypes(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListVirtualMachineTypes() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "Standard" {
			t.Errorf("resp = %+v, want count 1 with Standard", resp)
		}
		if resp.Results[0].Slug != "standard" {
			t.Errorf("Slug = %q, want %q", resp.Results[0].Slug, "standard")
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListVirtualMachineTypes(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestClient_CreateVirtualMachineType(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			gotCT = req.Header.Get("Content-Type")
			gotAuth = req.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(req.Body)
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusCreated)
			_, _ = writer.Write([]byte(`{"id":1,"name":"Standard","slug":"standard","created":"2024-01-01",
				"last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		vmt, err := client.CreateVirtualMachineType(
			context.Background(), "tok", domain.VirtualMachineTypeWrite{Name: "Standard"},
		)
		if err != nil {
			t.Fatalf("CreateVirtualMachineType() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/virtualization/virtual-machine-types/" {
			t.Errorf("method/path = %q %q, want POST /api/virtualization/virtual-machine-types/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "Standard" {
			t.Errorf("request body = %s, want name Standard", gotBody)
		}
		if vmt.ID != 1 || vmt.Name != "Standard" {
			t.Errorf("vmt = %+v, want id 1 name Standard", vmt)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateVirtualMachineType(context.Background(), "tok", domain.VirtualMachineTypeWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateVirtualMachineType(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			gotBody, _ = io.ReadAll(req.Body)
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"name":"Renamed","slug":"renamed","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		vmt, err := client.UpdateVirtualMachineType(
			context.Background(), "tok", 7, domain.VirtualMachineTypeWrite{Name: "Renamed"},
		)
		if err != nil {
			t.Fatalf("UpdateVirtualMachineType() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/virtualization/virtual-machine-types/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/virtualization/virtual-machine-types/7/", gotMethod, gotPath)
		}
		if vmt.ID != 7 || vmt.Name != "Renamed" {
			t.Errorf("vmt = %+v, want id 7 name Renamed", vmt)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "Renamed" {
			t.Errorf("request body = %s, want name present", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"slug":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateVirtualMachineType(context.Background(), "tok", 7, domain.VirtualMachineTypeWrite{Name: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteVirtualMachineType(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			writer.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteVirtualMachineType(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteVirtualMachineType() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/virtualization/virtual-machine-types/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/virtualization/virtual-machine-types/3/", gotMethod, gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteVirtualMachineType(context.Background(), "tok", 3)
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})

	t.Run("not found 404 (idempotent no-op)", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteVirtualMachineType(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteVirtualMachineType() 404 returned error: %v", err)
		}
	})
}

func TestVirtualMachineTypeWriteToWire(t *testing.T) {
	t.Parallel()

	desc := "d"
	in := domain.VirtualMachineTypeWrite{Name: "A", Description: &desc}
	wire := virtualMachineTypeWriteToWire(in)
	if wire.Name != "A" {
		t.Errorf("wire = %+v, want name A", wire)
	}
	if wire.Description == nil || *wire.Description != "d" {
		t.Errorf("wire.Description = %v, want d", wire.Description)
	}
}

func TestWireVirtualMachineTypeToDomain(t *testing.T) {
	t.Parallel()

	wire := WireVirtualMachineType{
		ID:      1,
		URL:     "http://nb/api/virtualization/virtual-machine-types/1/",
		Name:    "Standard",
		Slug:    "standard",
		Created: "2024-01-01",
	}
	d := wireVirtualMachineTypeToDomain(wire)
	if d.ID != 1 || d.Name != "Standard" || d.Slug != "standard" {
		t.Errorf("domain = %+v, want id 1 name Standard slug standard", d)
	}

	if got := wireVirtualMachineTypeToDomain(WireVirtualMachineType{}).Name; got != "" {
		t.Errorf("empty wire name = %q, want empty", got)
	}
}
