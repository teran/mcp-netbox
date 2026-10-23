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

// TestClient_ListContactRoles exercises the contact role list endpoint.
func TestClient_ListContactRoles(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "Token test-token" {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}
			if req.URL.Path != "/api/tenancy/contact-roles/" {
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
					"name": "Owner",
					"slug": "owner",
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListContactRoles(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListContactRoles() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "Owner" {
			t.Errorf("resp = %+v, want count 1 with Owner", resp)
		}
		if resp.Results[0].Slug != "owner" {
			t.Errorf("Slug = %q, want %q", resp.Results[0].Slug, "owner")
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListContactRoles(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestClient_CreateContactRole(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"name":"Owner","slug":"owner","created":"2024-01-01","last_update` +
				`d":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		role, err := client.CreateContactRole(context.Background(), "tok", domain.ContactRoleWrite{Name: "Owner"})
		if err != nil {
			t.Fatalf("CreateContactRole() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/tenancy/contact-roles/" {
			t.Errorf("method/path = %q %q, want POST /api/tenancy/contact-roles/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "Owner" {
			t.Errorf("request body = %s, want name Owner", gotBody)
		}
		if role.ID != 1 || role.Name != "Owner" {
			t.Errorf("role = %+v, want id 1 name Owner", role)
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
		_, err := client.CreateContactRole(context.Background(), "tok", domain.ContactRoleWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateContactRole(t *testing.T) {
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
		role, err := client.UpdateContactRole(
			context.Background(), "tok", 7, domain.ContactRoleWrite{Name: "Renamed"},
		)
		if err != nil {
			t.Fatalf("UpdateContactRole() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/tenancy/contact-roles/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/tenancy/contact-roles/7/", gotMethod, gotPath)
		}
		if role.ID != 7 || role.Name != "Renamed" {
			t.Errorf("role = %+v, want id 7 name Renamed", role)
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
		_, err := client.UpdateContactRole(context.Background(), "tok", 7, domain.ContactRoleWrite{Name: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteContactRole(t *testing.T) {
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
		if err := client.DeleteContactRole(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteContactRole() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/tenancy/contact-roles/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/tenancy/contact-roles/3/", gotMethod, gotPath)
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
		err := client.DeleteContactRole(context.Background(), "tok", 3)
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
		if err := client.DeleteContactRole(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteContactRole() 404 returned error: %v", err)
		}
	})
}

func TestContactRoleWriteToWire(t *testing.T) {
	t.Parallel()

	desc := "d"
	in := domain.ContactRoleWrite{Name: "A", Description: &desc}
	wire := contactRoleWriteToWire(in)
	if wire.Name != "A" {
		t.Errorf("wire = %+v, want name A", wire)
	}
	if wire.Description == nil || *wire.Description != "d" {
		t.Errorf("wire.Description = %v, want d", wire.Description)
	}
}

func TestWireContactRoleToDomain(t *testing.T) {
	t.Parallel()

	wire := WireContactRole{
		ID:      1,
		URL:     "http://nb/api/tenancy/contact-roles/1/",
		Name:    "Owner",
		Slug:    "owner",
		Created: "2024-01-01",
	}
	d := wireContactRoleToDomain(wire)
	if d.ID != 1 || d.Name != "Owner" || d.Slug != "owner" {
		t.Errorf("domain = %+v, want id 1 name Owner slug owner", d)
	}

	if got := wireContactRoleToDomain(WireContactRole{}).Name; got != "" {
		t.Errorf("empty wire name = %q, want empty", got)
	}
}
