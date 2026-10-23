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

// TestClient_ListRIRs exercises the RIR list endpoint.
func TestClient_ListRIRs(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "Token test-token" {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}
			if req.URL.Path != "/api/ipam/rirs/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"url": "http://nb/api/ipam/rirs/1/",
					"name": "ARIN",
					"slug": "arin",
					"is_private": false,
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListRIRs(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListRIRs() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "ARIN" {
			t.Errorf("resp = %+v, want count 1 with ARIN", resp)
		}
		if resp.Results[0].IsPrivate {
			t.Errorf("resp.Results[0].IsPrivate = true, want false")
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListRIRs(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_CreateRIR(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"name":"ARIN","slug":"arin","is_private":false,
				"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		isPrivate := true
		rir, err := client.CreateRIR(
			context.Background(), "tok", domain.RIRWrite{Name: "ARIN", IsPrivate: &isPrivate},
		)
		if err != nil {
			t.Fatalf("CreateRIR() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/ipam/rirs/" {
			t.Errorf("method/path = %q %q, want POST /api/ipam/rirs/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "ARIN" {
			t.Errorf("request body = %s, want name ARIN", gotBody)
		}
		if reqBody["is_private"] != true {
			t.Errorf("request body is_private = %v, want true", reqBody["is_private"])
		}
		if rir.ID != 1 || rir.Name != "ARIN" {
			t.Errorf("rir = %+v, want id 1 name ARIN", rir)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateRIR(context.Background(), "tok", domain.RIRWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateRIR(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"name":"Renamed","slug":"renamed","is_private":false,
				"created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		rir, err := client.UpdateRIR(context.Background(), "tok", 7, domain.RIRWrite{Name: "Renamed"})
		if err != nil {
			t.Fatalf("UpdateRIR() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/ipam/rirs/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/ipam/rirs/7/", gotMethod, gotPath)
		}
		if rir.ID != 7 || rir.Name != "Renamed" {
			t.Errorf("rir = %+v, want id 7 name Renamed", rir)
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
		_, err := client.UpdateRIR(context.Background(), "tok", 7, domain.RIRWrite{Name: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteRIR(t *testing.T) {
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
		if err := client.DeleteRIR(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteRIR() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/ipam/rirs/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/ipam/rirs/3/", gotMethod, gotPath)
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
		err := client.DeleteRIR(context.Background(), "tok", 3)
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
		if err := client.DeleteRIR(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteRIR() 404 returned error: %v", err)
		}
	})
}

func TestRIRWriteToWire(t *testing.T) {
	t.Parallel()

	isPrivate := true
	slug := "arin"
	desc := "registry"
	in := domain.RIRWrite{Name: "ARIN", Slug: &slug, IsPrivate: &isPrivate, Description: &desc, Tags: []string{"prod"}}
	wire := rirWriteToWire(in)
	if wire.Name != "ARIN" || wire.IsPrivate == nil || !*wire.IsPrivate {
		t.Errorf("wire = %+v, want name ARIN is_private true", wire)
	}
	if wire.Description == nil || *wire.Description != "registry" {
		t.Errorf("wire.Description = %v, want registry", wire.Description)
	}

	// The JSON keys must match NetBox's wire naming exactly.
	raw, err := json.Marshal(wire)
	if err != nil {
		t.Fatalf("marshal wire: %v", err)
	}
	for _, key := range []string{"name", "slug", "is_private", "description", "tags"} {
		if !jsonHasKey(raw, key) {
			t.Errorf("wire JSON missing key %q: %s", key, raw)
		}
	}
}

func TestWireRIRToDomain(t *testing.T) {
	t.Parallel()

	wire := WireRIR{
		ID:        1,
		URL:       "http://nb/api/ipam/rirs/1/",
		Name:      "ARIN",
		Slug:      "arin",
		IsPrivate: true,
		Created:   "2024-01-01",
	}
	got := wireRIRToDomain(wire)
	if got.ID != 1 || got.Name != "ARIN" || got.Slug != "arin" {
		t.Errorf("domain = %+v, want id 1 name ARIN slug arin", got)
	}
	if !got.IsPrivate {
		t.Errorf("domain.IsPrivate = false, want true")
	}
}

// jsonHasKey reports whether raw JSON contains the given top-level key. It is a
// package-level test helper shared by the IPAM entity wire tests to pin the
// exact NetBox wire JSON key names.
func jsonHasKey(raw []byte, key string) bool {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return false
	}
	_, ok := m[key]

	return ok
}
