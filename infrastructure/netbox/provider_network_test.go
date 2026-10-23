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

// TestClient_ListProviderNetworks exercises the provider network list endpoint.
func TestClient_ListProviderNetworks(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "Token test-token" {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}
			if req.URL.Path != "/api/circuits/provider-networks/" {
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
					"url": "http://nb/api/circuits/provider-networks/1/",
					"name": "Net A",
					"provider": {"id": 2, "name": "ACME", "slug": "acme", "url": "http://nb/api/circuits/providers/2/"},
					"description": "prod network",
					"comments": "note",
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListProviderNetworks(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListProviderNetworks() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "Net A" {
			t.Errorf("resp = %+v, want count 1 with Net A", resp)
		}
		if resp.Results[0].Provider == nil || resp.Results[0].Provider.ID != 2 {
			t.Errorf("resp.Results[0].Provider = %+v, want nested id 2", resp.Results[0].Provider)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListProviderNetworks(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_CreateProviderNetwork(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"name":"Net A","provider":{"id":2},` +
				`"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		provider := 2
		net, err := client.CreateProviderNetwork(
			context.Background(), "tok", domain.ProviderNetworkWrite{Name: "Net A", Provider: &provider},
		)
		if err != nil {
			t.Fatalf("CreateProviderNetwork() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/circuits/provider-networks/" {
			t.Errorf("method/path = %q %q, want POST /api/circuits/provider-networks/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "Net A" {
			t.Errorf("request body = %s, want name Net A", gotBody)
		}
		if reqBody["provider"] != float64(2) {
			t.Errorf("request body = %s, want provider 2", gotBody)
		}
		if net.ID != 1 || net.Name != "Net A" {
			t.Errorf("net = %+v, want id 1 name Net A", net)
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
		_, err := client.CreateProviderNetwork(context.Background(), "tok", domain.ProviderNetworkWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateProviderNetwork(t *testing.T) {
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
		comments := "note"
		net, err := client.UpdateProviderNetwork(
			context.Background(), "tok", 7, domain.ProviderNetworkWrite{Name: "Renamed", Comments: &comments},
		)
		if err != nil {
			t.Fatalf("UpdateProviderNetwork() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/circuits/provider-networks/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/circuits/provider-networks/7/", gotMethod, gotPath)
		}
		if net.ID != 7 || net.Name != "Renamed" {
			t.Errorf("net = %+v, want id 7 name Renamed", net)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["comments"] != "note" {
			t.Errorf("request body = %s, want comments present", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"provider":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateProviderNetwork(context.Background(), "tok", 7, domain.ProviderNetworkWrite{Name: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteProviderNetwork(t *testing.T) {
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
		if err := client.DeleteProviderNetwork(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteProviderNetwork() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/circuits/provider-networks/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/circuits/provider-networks/3/", gotMethod, gotPath)
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
		err := client.DeleteProviderNetwork(context.Background(), "tok", 3)
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
		if err := client.DeleteProviderNetwork(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteProviderNetwork() 404 returned error: %v", err)
		}
	})
}

func TestProviderNetworkWriteToWire(t *testing.T) {
	t.Parallel()

	provider := 2
	comments := "c"
	description := "d"
	in := domain.ProviderNetworkWrite{Name: "A", Provider: &provider, Comments: &comments, Description: &description}
	wire := providerNetworkWriteToWire(in)
	if wire.Name != "A" || wire.Provider == nil || *wire.Provider != 2 {
		t.Errorf("wire = %+v, want name A provider 2", wire)
	}
	if wire.Comments == nil || *wire.Comments != "c" {
		t.Errorf("wire.Comments = %v, want c", wire.Comments)
	}

	raw, err := json.Marshal(wire)
	if err != nil {
		t.Fatalf("marshal wire: %v", err)
	}
	for _, key := range []string{"name", "provider", "description", "comments"} {
		if !jsonHasKey(raw, key) {
			t.Errorf("wire JSON missing key %q: %s", key, raw)
		}
	}
}

func TestWireProviderNetworkToDomain(t *testing.T) {
	t.Parallel()

	wire := WireProviderNetwork{
		ID:       1,
		URL:      "http://nb/api/circuits/provider-networks/1/",
		Name:     "Net A",
		Provider: &WireNested{ID: 2, Name: "ACME", Slug: "acme"},
		Comments: "note",
	}
	got := wireProviderNetworkToDomain(wire)
	if got.ID != 1 || got.Name != "Net A" {
		t.Errorf("domain = %+v, want id 1 name Net A", got)
	}
	if got.Provider == nil || got.Provider.ID != 2 {
		t.Errorf("domain.Provider = %+v, want nested id 2", got.Provider)
	}
	if got.Comments != "note" {
		t.Errorf("domain.Comments = %q, want note", got.Comments)
	}

	if got := wireProviderNetworkToDomain(WireProviderNetwork{}).Provider; got != nil {
		t.Errorf("empty wire Provider = %+v, want nil", got)
	}
}
