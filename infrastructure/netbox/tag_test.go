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

// TestClient_ListTags exercises the extras tag list endpoint.
func TestClient_ListTags(t *testing.T) { //nolint:funlen
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "Token test-token" {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}
			if req.URL.Path != "/api/extras/tags/" {
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
					"url": "http://nb/api/extras/tags/1/",
					"name": "prod",
					"slug": "prod",
					"color": "ff0000",
					"description": "production",
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListTags(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListTags() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "prod" {
			t.Errorf("resp = %+v, want count 1 with prod", resp)
		}
		if resp.Results[0].Slug != "prod" {
			t.Errorf("Slug = %q, want %q", resp.Results[0].Slug, "prod")
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListTags(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("bad_json", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{not-json`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListTags(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected unmarshal error, got nil")
		}
	})
}

func TestClient_CreateTag(t *testing.T) {
	t.Parallel()

	testClientCreateTagSuccess(t)
	testClientCreateTagValidationError(t)
}

func testClientCreateTagSuccess(t *testing.T) {
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
		_, _ = writer.Write([]byte(`{"id":1,"url":"http://nb/api/extras/tags/1/","name":"prod","slug":"prod",` +
			`"color":"ff0000","description":"production","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
	}))
	defer srv.Close()

	color := "ff0000"
	client := NewClient(srv.URL, http.DefaultClient)
	tag, err := client.CreateTag(context.Background(), "tok", domain.TagWrite{
		Name: "prod", Slug: new("prod"), Color: &color, Description: new("production"),
	})
	if err != nil {
		t.Fatalf("CreateTag() returned error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/api/extras/tags/" {
		t.Errorf("method/path = %q %q, want POST /api/extras/tags/", gotMethod, gotPath)
	}
	if gotCT != "application/json" || gotAuth != "Token tok" {
		t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
	}
	var reqBody map[string]any
	_ = json.Unmarshal(gotBody, &reqBody)
	if reqBody["name"] != "prod" {
		t.Errorf("request body = %s, want name prod", gotBody)
	}
	if reqBody["color"] != "ff0000" {
		t.Errorf("request body color = %v, want ff0000", reqBody["color"])
	}
	if tag.ID != 1 || tag.Name != "prod" {
		t.Errorf("tag = %+v, want id 1 name prod", tag)
	}
}

func testClientCreateTagValidationError(t *testing.T) {
	body := `{"name":["This field is required."]}`
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte(body))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	_, err := client.CreateTag(context.Background(), "tok", domain.TagWrite{})
	var validationErr *domain.ValidationError
	if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
		t.Errorf("err = %v, want ValidationError body %s", err, body)
	}
}

func TestClient_UpdateTag(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":7,"url":"","name":"Renamed","slug":"renamed","color":"00ff00",` +
				`"description":"","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		tag, err := client.UpdateTag(context.Background(), "tok", 7, domain.TagWrite{Name: "Renamed"})
		if err != nil {
			t.Fatalf("UpdateTag() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/extras/tags/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/extras/tags/7/", gotMethod, gotPath)
		}
		if tag.ID != 7 || tag.Name != "Renamed" {
			t.Errorf("tag = %+v, want id 7 name Renamed", tag)
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
		_, err := client.UpdateTag(context.Background(), "tok", 7, domain.TagWrite{Name: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteTag(t *testing.T) {
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
		if err := client.DeleteTag(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteTag() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/extras/tags/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/extras/tags/3/", gotMethod, gotPath)
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
		err := client.DeleteTag(context.Background(), "tok", 3)
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
		if err := client.DeleteTag(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteTag() 404 returned error: %v", err)
		}
	})
}

func TestTagWriteToWire(t *testing.T) {
	t.Parallel()

	color := "ff0000"
	in := domain.TagWrite{Name: "A", Slug: new("a"), Color: &color, Description: new("d")}
	wire := tagWriteToWire(in)
	if wire.Name != "A" {
		t.Errorf("wire.Name = %q, want A", wire.Name)
	}
	if wire.Color == nil || *wire.Color != "ff0000" {
		t.Errorf("wire.Color = %v, want ff0000", wire.Color)
	}
	if wire.Description == nil || *wire.Description != "d" {
		t.Errorf("wire.Description = %v, want d", wire.Description)
	}
}

func TestWireTagToDomain(t *testing.T) {
	t.Parallel()

	wire := WireTag{
		ID:   1,
		URL:  "http://nb/api/extras/tags/1/",
		Name: "prod",
		Slug: "prod",
	}
	d := wireTagToDomain(wire)
	if d.ID != 1 || d.Name != "prod" || d.Slug != "prod" {
		t.Errorf("domain = %+v, want id 1 name prod slug prod", d)
	}

	if got := wireTagToDomain(WireTag{}).Name; got != "" {
		t.Errorf("empty wire name = %q, want empty", got)
	}
}
