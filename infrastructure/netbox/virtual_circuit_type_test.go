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

// TestClient_ListVirtualCircuitTypes exercises the virtual circuit type list endpoint.
func TestClient_ListVirtualCircuitTypes(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "Token test-token" {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}
			if req.URL.Path != "/api/circuits/virtual-circuit-types/" {
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
					"url": "http://nb/api/circuits/virtual-circuit-types/1/",
					"name": "VPLS",
					"slug": "vpls",
					"description": "VPLS service",
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListVirtualCircuitTypes(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListVirtualCircuitTypes() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "VPLS" {
			t.Errorf("resp = %+v, want count 1 with VPLS", resp)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListVirtualCircuitTypes(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestClient_CreateVirtualCircuitType(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"name":"VPLS","slug":"vpls",` +
				`"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		vcType, err := client.CreateVirtualCircuitType(
			context.Background(), "tok", domain.VirtualCircuitTypeWrite{Name: "VPLS"},
		)
		if err != nil {
			t.Fatalf("CreateVirtualCircuitType() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/circuits/virtual-circuit-types/" {
			t.Errorf("method/path = %q %q, want POST /api/circuits/virtual-circuit-types/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "VPLS" {
			t.Errorf("request body = %s, want name VPLS", gotBody)
		}
		if vcType.ID != 1 || vcType.Name != "VPLS" {
			t.Errorf("vcType = %+v, want id 1 name VPLS", vcType)
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
		_, err := client.CreateVirtualCircuitType(context.Background(), "tok", domain.VirtualCircuitTypeWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateVirtualCircuitType(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"name":"VPLS","slug":"vpls","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		vcType, err := client.UpdateVirtualCircuitType(
			context.Background(), "tok", 7, domain.VirtualCircuitTypeWrite{Name: "VPLS"},
		)
		if err != nil {
			t.Fatalf("UpdateVirtualCircuitType() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/circuits/virtual-circuit-types/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/circuits/virtual-circuit-types/7/", gotMethod, gotPath)
		}
		if vcType.ID != 7 || vcType.Name != "VPLS" {
			t.Errorf("vcType = %+v, want id 7 name VPLS", vcType)
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
		_, err := client.UpdateVirtualCircuitType(context.Background(), "tok", 7, domain.VirtualCircuitTypeWrite{Name: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteVirtualCircuitType(t *testing.T) {
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
		if err := client.DeleteVirtualCircuitType(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteVirtualCircuitType() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/circuits/virtual-circuit-types/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/circuits/virtual-circuit-types/3/", gotMethod, gotPath)
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
		err := client.DeleteVirtualCircuitType(context.Background(), "tok", 3)
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
		if err := client.DeleteVirtualCircuitType(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteVirtualCircuitType() 404 returned error: %v", err)
		}
	})
}

func TestVirtualCircuitTypeWriteToWire(t *testing.T) {
	t.Parallel()

	slug := "vpls"
	description := "d"
	in := domain.VirtualCircuitTypeWrite{Name: "A", Slug: &slug, Description: &description}
	wire := virtualCircuitTypeWriteToWire(in)
	if wire.Name != "A" {
		t.Errorf("wire = %+v, want name A", wire)
	}
	if wire.Slug == nil || *wire.Slug != "vpls" {
		t.Errorf("wire.Slug = %v, want vpls", wire.Slug)
	}

	raw, err := json.Marshal(wire)
	if err != nil {
		t.Fatalf("marshal wire: %v", err)
	}
	for _, key := range []string{"name", "slug", "description"} {
		if !jsonHasKey(raw, key) {
			t.Errorf("wire JSON missing key %q: %s", key, raw)
		}
	}
}

func TestWireVirtualCircuitTypeToDomain(t *testing.T) {
	t.Parallel()

	wire := WireVirtualCircuitType{
		ID:          1,
		URL:         "http://nb/api/circuits/virtual-circuit-types/1/",
		Name:        "VPLS",
		Slug:        "vpls",
		Description: "VPLS service",
	}
	d := wireVirtualCircuitTypeToDomain(wire)
	if d.ID != 1 || d.Name != "VPLS" || d.Slug != "vpls" {
		t.Errorf("domain = %+v, want id 1 name VPLS slug vpls", d)
	}
	if d.Description != "VPLS service" {
		t.Errorf("domain.Description = %q, want VPLS service", d.Description)
	}
}
