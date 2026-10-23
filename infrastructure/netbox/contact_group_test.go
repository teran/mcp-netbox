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

// TestClient_ListContactGroups exercises the contact group list endpoint.
func TestClient_ListContactGroups(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "Token test-token" {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}
			if req.URL.Path != "/api/tenancy/contact-groups/" {
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
					"name": "Admins",
					"slug": "admins",
					"parent": {"id": 2, "name": "Root", "slug": "root", "url": "http://nb/api/tenancy/contact-groups/2/"},
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListContactGroups(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListContactGroups() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "Admins" {
			t.Errorf("resp = %+v, want count 1 with Admins", resp)
		}
		if resp.Results[0].Parent == nil || resp.Results[0].Parent.ID != 2 {
			t.Errorf("resp.Results[0].Parent = %+v, want nested id 2", resp.Results[0].Parent)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListContactGroups(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_CreateContactGroup(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"name":"Admins","slug":"admins","created":"2024-01-01","last_updat` +
				`ed":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		parent := 2
		group, err := client.CreateContactGroup(
			context.Background(), "tok", domain.ContactGroupWrite{Name: "Admins", Parent: &parent},
		)
		if err != nil {
			t.Fatalf("CreateContactGroup() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/tenancy/contact-groups/" {
			t.Errorf("method/path = %q %q, want POST /api/tenancy/contact-groups/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "Admins" {
			t.Errorf("request body = %s, want name Admins", gotBody)
		}
		if reqBody["parent"] != float64(2) {
			t.Errorf("request body = %s, want parent 2", gotBody)
		}
		if group.ID != 1 || group.Name != "Admins" {
			t.Errorf("group = %+v, want id 1 name Admins", group)
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
		_, err := client.CreateContactGroup(context.Background(), "tok", domain.ContactGroupWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateContactGroup(t *testing.T) {
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
		group, err := client.UpdateContactGroup(
			context.Background(), "tok", 7, domain.ContactGroupWrite{Name: "Renamed", Description: &desc},
		)
		if err != nil {
			t.Fatalf("UpdateContactGroup() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/tenancy/contact-groups/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/tenancy/contact-groups/7/", gotMethod, gotPath)
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
		_, err := client.UpdateContactGroup(context.Background(), "tok", 7, domain.ContactGroupWrite{Name: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteContactGroup(t *testing.T) {
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
		if err := client.DeleteContactGroup(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteContactGroup() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/tenancy/contact-groups/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/tenancy/contact-groups/3/", gotMethod, gotPath)
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
		err := client.DeleteContactGroup(context.Background(), "tok", 3)
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
		if err := client.DeleteContactGroup(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteContactGroup() 404 returned error: %v", err)
		}
	})
}

func TestContactGroupWriteToWire(t *testing.T) {
	t.Parallel()

	parent := 2
	desc := "d"
	in := domain.ContactGroupWrite{Name: "A", Parent: &parent, Description: &desc}
	wire := contactGroupWriteToWire(in)
	if wire.Name != "A" || wire.Parent == nil || *wire.Parent != 2 {
		t.Errorf("wire = %+v, want name A parent 2", wire)
	}
	if wire.Description == nil || *wire.Description != "d" {
		t.Errorf("wire.Description = %v, want d", wire.Description)
	}
}

func TestWireContactGroupToDomain(t *testing.T) {
	t.Parallel()

	wire := WireContactGroup{
		ID:      1,
		URL:     "http://nb/api/tenancy/contact-groups/1/",
		Name:    "Admins",
		Slug:    "admins",
		Parent:  &WireNested{ID: 2, Name: "Root", Slug: "root"},
		Created: "2024-01-01",
	}
	d := wireContactGroupToDomain(wire)
	if d.ID != 1 || d.Name != "Admins" || d.Slug != "admins" {
		t.Errorf("domain = %+v, want id 1 name Admins slug admins", d)
	}
	if d.Parent == nil || d.Parent.ID != 2 {
		t.Errorf("domain.Parent = %+v, want nested id 2", d.Parent)
	}

	if got := wireContactGroupToDomain(WireContactGroup{}).Parent; got != nil {
		t.Errorf("empty wire Parent = %+v, want nil", got)
	}
}
