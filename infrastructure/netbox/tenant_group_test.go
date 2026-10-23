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

// TestClient_ListTenantGroups exercises the tenant group list endpoint.
func TestClient_ListTenantGroups(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "Token test-token" {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}
			if req.URL.Path != "/api/tenancy/tenant-groups/" {
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
					"name": "Clouds",
					"slug": "clouds",
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListTenantGroups(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListTenantGroups() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("Count = %d, want %d", resp.Count, 1)
		}
		if len(resp.Results) != 1 {
			t.Fatalf("len(Results) = %d, want %d", len(resp.Results), 1)
		}
		if resp.Results[0].Name != "Clouds" {
			t.Errorf("Name = %q, want %q", resp.Results[0].Name, "Clouds")
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListTenantGroups(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_CreateTenantGroup(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"name":"Clouds","slug":"clouds","created":"2024-01-01","last_update` +
				`d":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		group, err := client.CreateTenantGroup(context.Background(), "tok", domain.TenantGroupWrite{Name: "Clouds"})
		if err != nil {
			t.Fatalf("CreateTenantGroup() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/tenancy/tenant-groups/" {
			t.Errorf("path = %q, want /api/tenancy/tenant-groups/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Token tok" {
			t.Errorf("Authorization = %q, want Token tok", gotAuth)
		}
		var reqBody map[string]string
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "Clouds" {
			t.Errorf("request body = %s, want name Clouds", gotBody)
		}
		if group.ID != 1 || group.Name != "Clouds" {
			t.Errorf("group = %+v, want id 1 name Clouds", group)
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
		_, err := client.CreateTenantGroup(context.Background(), "tok", domain.TenantGroupWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if validationErr.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode = %d, want 400", validationErr.StatusCode)
		}
		if string(validationErr.Body) != body {
			t.Errorf("Body = %s, want %s", validationErr.Body, body)
		}
	})
}

func TestClient_UpdateTenantGroup(t *testing.T) {
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
		desc := "updated"
		group, err := client.UpdateTenantGroup(
			context.Background(), "tok", 7, domain.TenantGroupWrite{Name: "Renamed", Description: &desc},
		)
		if err != nil {
			t.Fatalf("UpdateTenantGroup() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/tenancy/tenant-groups/7/" {
			t.Errorf("path = %q, want /api/tenancy/tenant-groups/7/", gotPath)
		}
		if group.ID != 7 || group.Name != "Renamed" {
			t.Errorf("group = %+v, want id 7 name Renamed", group)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["description"] != "updated" {
			t.Errorf("request body = %s, want description present", gotBody)
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
		_, err := client.UpdateTenantGroup(context.Background(), "tok", 7, domain.TenantGroupWrite{Name: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(validationErr.Body) != body {
			t.Errorf("Body = %s, want %s", validationErr.Body, body)
		}
	})
}

func TestClient_DeleteTenantGroup(t *testing.T) {
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
		if err := client.DeleteTenantGroup(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteTenantGroup() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/tenancy/tenant-groups/3/" {
			t.Errorf("path = %q, want /api/tenancy/tenant-groups/3/", gotPath)
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
		err := client.DeleteTenantGroup(context.Background(), "tok", 3)
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(validationErr.Body) != body {
			t.Errorf("Body = %s, want %s", validationErr.Body, body)
		}
	})

	t.Run("not found 404 (idempotent no-op)", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteTenantGroup(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteTenantGroup() returned error: %v", err)
		}
	})
}

func TestTenantGroupWriteToWire(t *testing.T) {
	t.Parallel()

	parent := 2
	in := domain.TenantGroupWrite{Name: "A", Parent: &parent}
	wire := tenantGroupWriteToWire(in)
	if wire.Name != "A" || wire.Parent == nil || *wire.Parent != 2 {
		t.Errorf("wire = %+v, want name A parent 2", wire)
	}
}

func TestWireTenantGroupToDomain(t *testing.T) {
	t.Parallel()

	wire := WireTenantGroup{
		ID:      1,
		URL:     "http://nb/api/tenancy/tenant-groups/1/",
		Name:    "Clouds",
		Slug:    "clouds",
		Created: "2024-01-01",
	}
	d := wireTenantGroupToDomain(wire)
	if d.ID != 1 || d.Name != "Clouds" || d.Slug != "clouds" {
		t.Errorf("domain = %+v, want id 1 name Clouds slug clouds", d)
	}

	if got := wireTenantGroupToDomain(WireTenantGroup{}).Name; got != "" {
		t.Errorf("empty wire name = %q, want empty", got)
	}
}

// TestTenantWriteToWire_WithGroup verifies the numeric tenant-group ID is
// forwarded to the wire request model.
func TestTenantWriteToWire_WithGroup(t *testing.T) {
	t.Parallel()

	group := 42
	in := domain.TenantWrite{Name: "ACME", Group: &group}
	wire := tenantWriteToWire(in)
	if wire.Group == nil || *wire.Group != 42 {
		t.Errorf("wire.Group = %v, want 42", wire.Group)
	}
}

// TestWireTenantToDomain_WithGroup verifies the tenant group reference is
// returned on the domain model.
func TestWireTenantToDomain_WithGroup(t *testing.T) {
	t.Parallel()

	wire := WireTenant{
		ID:    1,
		URL:   "http://nb/api/tenancy/tenants/1/",
		Name:  "ACME",
		Slug:  "acme",
		Group: &WireNested{ID: 42, Name: "Clouds", Slug: "clouds"},
	}
	d := wireTenantToDomain(wire)
	if d.Group == nil || d.Group.ID != 42 || d.Group.Name != "Clouds" {
		t.Errorf("domain.Group = %+v, want nested id 42 name Clouds", d.Group)
	}

	if got := wireTenantToDomain(WireTenant{}).Group; got != nil {
		t.Errorf("empty wire Group = %+v, want nil", got)
	}
}
