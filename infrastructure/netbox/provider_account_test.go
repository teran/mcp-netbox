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

// TestClient_ListProviderAccounts exercises the provider account list endpoint.
func TestClient_ListProviderAccounts(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "Token test-token" {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}
			if req.URL.Path != "/api/circuits/provider-accounts/" {
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
					"url": "http://nb/api/circuits/provider-accounts/1/",
					"name": "Acc A",
					"provider": {"id": 2, "name": "ACME", "slug": "acme", "url": "http://nb/api/circuits/providers/2/"},
					"account": "acc-1",
					"description": "prod account",
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListProviderAccounts(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListProviderAccounts() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "Acc A" {
			t.Errorf("resp = %+v, want count 1 with Acc A", resp)
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
		_, err := client.ListProviderAccounts(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_CreateProviderAccount(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"name":"Acc A","account":"acc-1","provider":{"id":2},` +
				`"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		provider := 2
		account := "acc-1"
		acc, err := client.CreateProviderAccount(
			context.Background(), "tok", domain.ProviderAccountWrite{Name: "Acc A", Provider: &provider, Account: &account},
		)
		if err != nil {
			t.Fatalf("CreateProviderAccount() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/circuits/provider-accounts/" {
			t.Errorf("method/path = %q %q, want POST /api/circuits/provider-accounts/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "Acc A" {
			t.Errorf("request body = %s, want name Acc A", gotBody)
		}
		if reqBody["provider"] != float64(2) {
			t.Errorf("request body = %s, want provider 2", gotBody)
		}
		if reqBody["account"] != "acc-1" {
			t.Errorf("request body = %s, want account acc-1", gotBody)
		}
		if acc.ID != 1 || acc.Name != "Acc A" {
			t.Errorf("acc = %+v, want id 1 name Acc A", acc)
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
		_, err := client.CreateProviderAccount(context.Background(), "tok", domain.ProviderAccountWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateProviderAccount(t *testing.T) {
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
		account := "acc-2"
		acc, err := client.UpdateProviderAccount(
			context.Background(), "tok", 7, domain.ProviderAccountWrite{Name: "Renamed", Account: &account},
		)
		if err != nil {
			t.Fatalf("UpdateProviderAccount() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/circuits/provider-accounts/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/circuits/provider-accounts/7/", gotMethod, gotPath)
		}
		if acc.ID != 7 || acc.Name != "Renamed" {
			t.Errorf("acc = %+v, want id 7 name Renamed", acc)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["account"] != "acc-2" {
			t.Errorf("request body = %s, want account present", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"account":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateProviderAccount(context.Background(), "tok", 7, domain.ProviderAccountWrite{Name: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteProviderAccount(t *testing.T) {
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
		if err := client.DeleteProviderAccount(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteProviderAccount() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/circuits/provider-accounts/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/circuits/provider-accounts/3/", gotMethod, gotPath)
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
		err := client.DeleteProviderAccount(context.Background(), "tok", 3)
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
		if err := client.DeleteProviderAccount(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteProviderAccount() 404 returned error: %v", err)
		}
	})
}

func TestProviderAccountWriteToWire(t *testing.T) {
	t.Parallel()

	provider := 2
	account := "acc-1"
	desc := "d"
	in := domain.ProviderAccountWrite{Name: "A", Provider: &provider, Account: &account, Description: &desc}
	wire := providerAccountWriteToWire(in)
	if wire.Name != "A" || wire.Provider == nil || *wire.Provider != 2 {
		t.Errorf("wire = %+v, want name A provider 2", wire)
	}
	if wire.Account == nil || *wire.Account != "acc-1" {
		t.Errorf("wire.Account = %v, want acc-1", wire.Account)
	}

	raw, err := json.Marshal(wire)
	if err != nil {
		t.Fatalf("marshal wire: %v", err)
	}
	for _, key := range []string{"name", "provider", "account", "description"} {
		if !jsonHasKey(raw, key) {
			t.Errorf("wire JSON missing key %q: %s", key, raw)
		}
	}
}

func TestWireProviderAccountToDomain(t *testing.T) {
	t.Parallel()

	wire := WireProviderAccount{
		ID:       1,
		URL:      "http://nb/api/circuits/provider-accounts/1/",
		Name:     "Acc A",
		Provider: &WireNested{ID: 2, Name: "ACME", Slug: "acme"},
		Account:  "acc-1",
	}
	got := wireProviderAccountToDomain(wire)
	if got.ID != 1 || got.Name != "Acc A" {
		t.Errorf("domain = %+v, want id 1 name Acc A", got)
	}
	if got.Provider == nil || got.Provider.ID != 2 {
		t.Errorf("domain.Provider = %+v, want nested id 2", got.Provider)
	}
	if got.Account != "acc-1" {
		t.Errorf("domain.Account = %q, want acc-1", got.Account)
	}

	if got := wireProviderAccountToDomain(WireProviderAccount{}).Provider; got != nil {
		t.Errorf("empty wire Provider = %+v, want nil", got)
	}
}
