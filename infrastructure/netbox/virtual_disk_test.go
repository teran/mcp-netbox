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

// TestClient_ListVirtualDisks exercises the virtual disk list endpoint.
//
//nolint:funlen // Large declarative httptest function; splitting is not sensible.
func TestClient_ListVirtualDisks(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "Token test-token" {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}
			if req.URL.Path != "/api/virtualization/virtual-disks/" {
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
					"name": "Data",
					"virtual_machine": {"id": 5, "name": "vm1", "url": "http://nb/api/virtualization/virtual-machines/5/"},
					"size": 100,
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListVirtualDisks(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListVirtualDisks() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "Data" {
			t.Errorf("resp = %+v, want count 1 with Data", resp)
		}
		if resp.Results[0].VirtualMachine == nil || resp.Results[0].VirtualMachine.ID != 5 {
			t.Errorf("resp.Results[0].VirtualMachine = %+v, want nested id 5", resp.Results[0].VirtualMachine)
		}
		if resp.Results[0].Size != 100 {
			t.Errorf("Size = %d, want 100", resp.Results[0].Size)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListVirtualDisks(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_CreateVirtualDisk(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"name":"Data","virtual_machine":{"id":5},"size":100,
				"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		vm := 5
		size := 100
		disk, err := client.CreateVirtualDisk(
			context.Background(), "tok", domain.VirtualDiskWrite{Name: "Data", VirtualMachine: &vm, Size: &size},
		)
		if err != nil {
			t.Fatalf("CreateVirtualDisk() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/virtualization/virtual-disks/" {
			t.Errorf("method/path = %q %q, want POST /api/virtualization/virtual-disks/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		// The wire JSON keys for virtual disk are pinned: virtual_machine and size.
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "Data" {
			t.Errorf("request body = %s, want name Data", gotBody)
		}
		if reqBody["virtual_machine"] != float64(5) {
			t.Errorf("request body = %s, want virtual_machine 5", gotBody)
		}
		if reqBody["size"] != float64(100) {
			t.Errorf("request body = %s, want size 100", gotBody)
		}
		if disk.ID != 1 || disk.Name != "Data" {
			t.Errorf("disk = %+v, want id 1 name Data", disk)
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
		_, err := client.CreateVirtualDisk(context.Background(), "tok", domain.VirtualDiskWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateVirtualDisk(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":7,"name":"Renamed","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		size := 200
		disk, err := client.UpdateVirtualDisk(
			context.Background(), "tok", 7, domain.VirtualDiskWrite{Name: "Renamed", Size: &size},
		)
		if err != nil {
			t.Fatalf("UpdateVirtualDisk() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/virtualization/virtual-disks/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/virtualization/virtual-disks/7/", gotMethod, gotPath)
		}
		if disk.ID != 7 || disk.Name != "Renamed" {
			t.Errorf("disk = %+v, want id 7 name Renamed", disk)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["size"] != float64(200) {
			t.Errorf("request body = %s, want size present", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"virtual_machine":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateVirtualDisk(context.Background(), "tok", 7, domain.VirtualDiskWrite{Name: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteVirtualDisk(t *testing.T) {
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
		if err := client.DeleteVirtualDisk(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteVirtualDisk() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/virtualization/virtual-disks/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/virtualization/virtual-disks/3/", gotMethod, gotPath)
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
		err := client.DeleteVirtualDisk(context.Background(), "tok", 3)
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
		if err := client.DeleteVirtualDisk(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteVirtualDisk() 404 returned error: %v", err)
		}
	})
}

func TestVirtualDiskWriteToWire(t *testing.T) {
	t.Parallel()

	vm := 5
	size := 100
	in := domain.VirtualDiskWrite{Name: "A", VirtualMachine: &vm, Size: &size}
	wire := virtualDiskWriteToWire(in)
	if wire.Name != "A" || wire.VirtualMachine == nil || *wire.VirtualMachine != 5 {
		t.Errorf("wire = %+v, want name A virtual_machine 5", wire)
	}
	if wire.Size == nil || *wire.Size != 100 {
		t.Errorf("wire.Size = %v, want 100", wire.Size)
	}
}

func TestWireVirtualDiskToDomain(t *testing.T) {
	t.Parallel()

	wire := WireVirtualDisk{
		ID:             1,
		URL:            "http://nb/api/virtualization/virtual-disks/1/",
		Name:           "Data",
		VirtualMachine: &WireNested{ID: 5, Name: "vm1", Slug: "vm1"},
		Size:           100,
		Created:        "2024-01-01",
	}
	got := wireVirtualDiskToDomain(wire)
	if got.ID != 1 || got.Name != "Data" {
		t.Errorf("domain = %+v, want id 1 name Data", got)
	}
	if got.VirtualMachine == nil || got.VirtualMachine.ID != 5 {
		t.Errorf("domain.VirtualMachine = %+v, want nested id 5", got.VirtualMachine)
	}
	if got.Size != 100 {
		t.Errorf("domain.Size = %d, want 100", got.Size)
	}

	if got := wireVirtualDiskToDomain(WireVirtualDisk{}).VirtualMachine; got != nil {
		t.Errorf("empty wire VirtualMachine = %+v, want nil", got)
	}
}
