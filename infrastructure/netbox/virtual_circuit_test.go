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

// TestClient_ListVirtualCircuits exercises the virtual circuit list endpoint.
//
//nolint:funlen // Exercises the full list payload; splitting is not sensible.
func TestClient_ListVirtualCircuits(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "Token test-token" {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}
			if req.URL.Path != "/api/circuits/virtual-circuits/" {
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
					"url": "http://nb/api/circuits/virtual-circuits/1/",
					"cid": "VC-1",
					"type": {"id": 2, "name": "VPLS", "slug": "vpls", "url": "http://nb/api/circuits/virtual-circuit-types/2/"},
					"provider_network": {"id": 3, "name": "Net A", "slug": "net-a",
						"url": "http://nb/api/circuits/provider-networks/3/"},
					"status": {"value": "planned", "label": "Planned"},
					"description": "a virtual circuit",
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListVirtualCircuits(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListVirtualCircuits() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].CID != "VC-1" {
			t.Errorf("resp = %+v, want count 1 with VC-1", resp)
		}
		if resp.Results[0].Type == nil || resp.Results[0].Type.ID != 2 {
			t.Errorf("resp.Results[0].Type = %+v, want nested id 2", resp.Results[0].Type)
		}
		if resp.Results[0].ProviderNetwork == nil || resp.Results[0].ProviderNetwork.ID != 3 {
			t.Errorf("resp.Results[0].ProviderNetwork = %+v, want nested id 3", resp.Results[0].ProviderNetwork)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListVirtualCircuits(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_CreateVirtualCircuit(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"cid":"VC-1","type":{"id":2},` +
				`"provider_network":{"id":3},"status":{"value":"planned","label":"Planned"},` +
				`"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		vcType := 2
		providerNetwork := 3
		status := "planned"
		got, err := client.CreateVirtualCircuit(
			context.Background(), "tok", domain.VirtualCircuitWrite{
				CID: "VC-1", Type: &vcType, ProviderNetwork: &providerNetwork, Status: &status,
			},
		)
		if err != nil {
			t.Fatalf("CreateVirtualCircuit() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/circuits/virtual-circuits/" {
			t.Errorf("method/path = %q %q, want POST /api/circuits/virtual-circuits/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["cid"] != "VC-1" {
			t.Errorf("request body = %s, want cid VC-1", gotBody)
		}
		if reqBody["type"] != float64(2) {
			t.Errorf("request body = %s, want type 2", gotBody)
		}
		if reqBody["provider_network"] != float64(3) {
			t.Errorf("request body = %s, want provider_network 3", gotBody)
		}
		if reqBody["status"] != "planned" {
			t.Errorf("request body = %s, want status planned", gotBody)
		}
		if got.ID != 1 || got.CID != "VC-1" {
			t.Errorf("vc = %+v, want id 1 cid VC-1", got)
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
		_, err := client.CreateVirtualCircuit(context.Background(), "tok", domain.VirtualCircuitWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateVirtualCircuit(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":7,"cid":"Renamed","type":{"id":2},` +
				`"provider_network":{"id":3},"status":{"value":"planned","label":"Planned"},"created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		status := "active"
		got, err := client.UpdateVirtualCircuit(
			context.Background(), "tok", 7, domain.VirtualCircuitWrite{CID: "Renamed", Status: &status},
		)
		if err != nil {
			t.Fatalf("UpdateVirtualCircuit() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/circuits/virtual-circuits/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/circuits/virtual-circuits/7/", gotMethod, gotPath)
		}
		if got.ID != 7 || got.CID != "Renamed" {
			t.Errorf("vc = %+v, want id 7 cid Renamed", got)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["status"] != "active" {
			t.Errorf("request body = %s, want status present", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"type":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateVirtualCircuit(context.Background(), "tok", 7, domain.VirtualCircuitWrite{CID: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteVirtualCircuit(t *testing.T) {
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
		if err := client.DeleteVirtualCircuit(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteVirtualCircuit() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/circuits/virtual-circuits/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/circuits/virtual-circuits/3/", gotMethod, gotPath)
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
		err := client.DeleteVirtualCircuit(context.Background(), "tok", 3)
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
		if err := client.DeleteVirtualCircuit(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteVirtualCircuit() 404 returned error: %v", err)
		}
	})
}

func TestVirtualCircuitWriteToWire(t *testing.T) {
	t.Parallel()

	vcType := 2
	providerNetwork := 3
	status := "planned"
	description := "d"
	in := domain.VirtualCircuitWrite{
		CID: "A", Type: &vcType, ProviderNetwork: &providerNetwork, Status: &status, Description: &description,
	}
	wire := virtualCircuitWriteToWire(in)
	if wire.CID != "A" || wire.Type == nil || *wire.Type != 2 {
		t.Errorf("wire = %+v, want cid A type 2", wire)
	}
	if wire.ProviderNetwork == nil || *wire.ProviderNetwork != 3 {
		t.Errorf("wire.ProviderNetwork = %v, want 3", wire.ProviderNetwork)
	}
	if wire.Status == nil || *wire.Status != "planned" {
		t.Errorf("wire.Status = %v, want planned", wire.Status)
	}

	raw, err := json.Marshal(wire)
	if err != nil {
		t.Fatalf("marshal wire: %v", err)
	}
	for _, key := range []string{"cid", "type", "provider_network", "status", "description"} {
		if !jsonHasKey(raw, key) {
			t.Errorf("wire JSON missing key %q: %s", key, raw)
		}
	}
}

func TestWireVirtualCircuitToDomain(t *testing.T) {
	t.Parallel()

	wire := WireVirtualCircuit{
		ID:              1,
		URL:             "http://nb/api/circuits/virtual-circuits/1/",
		CID:             "VC-1",
		Type:            &WireNested{ID: 2, Name: "VPLS", Slug: "vpls"},
		ProviderNetwork: &WireNested{ID: 3, Name: "Net A", Slug: "net-a"},
		Status:          &WireLabel{Value: "planned", Label: "Planned"},
	}
	got := wireVirtualCircuitToDomain(wire)
	if got.ID != 1 || got.CID != "VC-1" {
		t.Errorf("domain = %+v, want id 1 cid VC-1", got)
	}
	if got.Type == nil || got.Type.ID != 2 {
		t.Errorf("domain.Type = %+v, want nested id 2", got.Type)
	}
	if got.ProviderNetwork == nil || got.ProviderNetwork.ID != 3 {
		t.Errorf("domain.ProviderNetwork = %+v, want nested id 3", got.ProviderNetwork)
	}
	if got.Status == nil || got.Status.Value != "planned" {
		t.Errorf("domain.Status = %+v, want planned", got.Status)
	}

	if got := wireVirtualCircuitToDomain(WireVirtualCircuit{}).Type; got != nil {
		t.Errorf("empty wire Type = %+v, want nil", got)
	}
}
