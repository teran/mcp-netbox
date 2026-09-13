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

func TestClient_CreateSite(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"name":"Site A","slug":"site-a","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		site, err := client.CreateSite(context.Background(), "tok", domain.SiteWrite{Name: "Site A"})
		if err != nil {
			t.Fatalf("CreateSite() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/dcim/sites/" {
			t.Errorf("path = %q, want /api/dcim/sites/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]string
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "Site A" {
			t.Errorf("request body = %s, want name Site A", gotBody)
		}
		if site.ID != 1 || site.Name != "Site A" {
			t.Errorf("site = %+v, want id 1 name Site A", site)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateSite(context.Background(), "tok", domain.SiteWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if ve.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode = %d, want 400", ve.StatusCode)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateSite(context.Background(), "tok", domain.SiteWrite{Name: "X"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			t.Errorf("err = %+v, want NOT a ValidationError for 404", ve)
		}
	})

	t.Run("unauthorized 401", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateSite(context.Background(), "tok", domain.SiteWrite{Name: "X"})
		if err == nil || !contains(err.Error(), "unauthorized") {
			t.Errorf("err = %v, want unauthorized", err)
		}
	})
}

func TestClient_UpdateSite(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"name":"Renamed","slug":"renamed","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		desc := "updated"
		site, err := client.UpdateSite(context.Background(), "tok", 7, domain.SiteWrite{Name: "Renamed", Description: &desc})
		if err != nil {
			t.Fatalf("UpdateSite() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/dcim/sites/7/" {
			t.Errorf("path = %q, want /api/dcim/sites/7/", gotPath)
		}
		if site.ID != 7 || site.Name != "Renamed" {
			t.Errorf("site = %+v, want id 7 name Renamed", site)
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
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateSite(context.Background(), "tok", 7, domain.SiteWrite{Name: "X"})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteSite(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteSite(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteSite() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/dcim/sites/3/" {
			t.Errorf("path = %q, want /api/dcim/sites/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteSite(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteSite(context.Background(), "tok", 999); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestMapWriteStatus(t *testing.T) {
	t.Parallel()

	t.Run("created 201 returns body", func(t *testing.T) {
		raw, err := mapWriteStatus(http.StatusCreated, []byte(`{"id":1}`))
		if err != nil {
			t.Fatalf("mapWriteStatus() error: %v", err)
		}
		if string(raw) != `{"id":1}` {
			t.Errorf("raw = %s, want body", raw)
		}
	})

	t.Run("ok 200 returns body", func(t *testing.T) {
		raw, err := mapWriteStatus(http.StatusOK, []byte(`{"id":2}`))
		if err != nil || string(raw) != `{"id":2}` {
			t.Errorf("raw = %s, err = %v", raw, err)
		}
	})

	t.Run("no content 204 returns nil", func(t *testing.T) {
		raw, err := mapWriteStatus(http.StatusNoContent, nil)
		if err != nil {
			t.Fatalf("mapWriteStatus() error: %v", err)
		}
		if raw != nil {
			t.Errorf("raw = %v, want nil", raw)
		}
	})

	t.Run("bad request 400 returns validation error", func(t *testing.T) {
		_, err := mapWriteStatus(http.StatusBadRequest, []byte(`{"e":["x"]}`))
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != `{"e":["x"]}` {
			t.Errorf("Body = %s, want body", ve.Body)
		}
	})

	t.Run("unexpected 500 returns error", func(t *testing.T) {
		_, err := mapWriteStatus(http.StatusInternalServerError, []byte(`{}`))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestSiteWriteToWire(t *testing.T) {
	t.Parallel()

	in := domain.SiteWrite{Name: "A", Slug: ptrStr("a")}
	w := siteWriteToWire(in)
	if w.Name != "A" || w.Slug == nil || *w.Slug != "a" {
		t.Errorf("wire = %+v, want name A slug a", w)
	}
}

func ptrStr(s string) *string { return &s }
