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

// TestClient_ListVirtualCircuitTerminations exercises the virtual circuit termination list endpoint.
func TestClient_ListVirtualCircuitTerminations(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/circuits/virtual-circuit-terminations/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"url": "http://nb/api/circuits/virtual-circuit-terminations/1/",
					"virtual_circuit": {"id": 5, "name": "VC-1", "url": "http://nb/api/circuits/virtual-circuits/5/"},
					"role": {"value": "peer", "label": "Peer"},
					"interface": {"id": 9, "name": "eth0", "url": "http://nb/api/dcim/interfaces/9/"},
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListVirtualCircuitTerminations(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListVirtualCircuitTerminations() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 {
			t.Errorf("resp = %+v, want count 1", resp)
		}
		if resp.Results[0].VirtualCircuit == nil || resp.Results[0].VirtualCircuit.ID != 5 {
			t.Errorf("resp.Results[0].VirtualCircuit = %+v, want nested id 5", resp.Results[0].VirtualCircuit)
		}
		if resp.Results[0].Interface == nil || resp.Results[0].Interface.ID != 9 {
			t.Errorf(
				"resp.Results[0].Interface = %+v, want nested id 9",
				resp.Results[0].Interface,
			)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListVirtualCircuitTerminations(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_CreateVirtualCircuitTermination(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"virtual_circuit":{"id":5},` +
				`"role":{"value":"peer","label":"Peer"},` +
				`"interface":{"id":9,"name":"eth0"},` +
				`"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		role := "peer"
		iface := 9
		term, err := client.CreateVirtualCircuitTermination(
			context.Background(), "tok", domain.VirtualCircuitTerminationWrite{
				VirtualCircuit: 5, Role: &role, Interface: iface,
			},
		)
		if err != nil {
			t.Fatalf("CreateVirtualCircuitTermination() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/circuits/virtual-circuit-terminations/" {
			t.Errorf("method/path = %q %q, want POST /api/circuits/virtual-circuit-terminations/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["virtual_circuit"] != float64(5) {
			t.Errorf("request body = %s, want virtual_circuit 5", gotBody)
		}
		if reqBody["role"] != "peer" {
			t.Errorf("request body = %s, want role peer", gotBody)
		}
		if reqBody["interface"] != float64(9) {
			t.Errorf("request body interface = %v, want 9", reqBody["interface"])
		}
		if term.ID != 1 {
			t.Errorf("term = %+v, want id 1", term)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"virtual_circuit":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateVirtualCircuitTermination(context.Background(), "tok", domain.VirtualCircuitTerminationWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateVirtualCircuitTermination(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"virtual_circuit":{"id":5},` +
				`"role":{"value":"peer","label":"Peer"},` +
				`"created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		term, err := client.UpdateVirtualCircuitTermination(
			context.Background(), "tok", 7, domain.VirtualCircuitTerminationWrite{VirtualCircuit: 5},
		)
		if err != nil {
			t.Fatalf("UpdateVirtualCircuitTermination() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/circuits/virtual-circuit-terminations/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/circuits/virtual-circuit-terminations/7/", gotMethod, gotPath)
		}
		if term.ID != 7 {
			t.Errorf("term = %+v, want id 7", term)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"virtual_circuit":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateVirtualCircuitTermination(
			context.Background(), "tok", 7, domain.VirtualCircuitTerminationWrite{VirtualCircuit: 5},
		)
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteVirtualCircuitTermination(t *testing.T) {
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
		if err := client.DeleteVirtualCircuitTermination(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteVirtualCircuitTermination() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/circuits/virtual-circuit-terminations/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/circuits/virtual-circuit-terminations/3/", gotMethod, gotPath)
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
		err := client.DeleteVirtualCircuitTermination(context.Background(), "tok", 3)
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
		if err := client.DeleteVirtualCircuitTermination(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteVirtualCircuitTermination() 404 returned error: %v", err)
		}
	})
}

func TestVirtualCircuitTerminationWriteToWire(t *testing.T) {
	t.Parallel()

	role := "peer"
	iface := 9
	in := domain.VirtualCircuitTerminationWrite{
		VirtualCircuit: 5, Role: &role, Interface: iface,
	}
	wire := virtualCircuitTerminationWriteToWire(in)
	if wire.VirtualCircuit != 5 {
		t.Errorf("wire = %+v, want virtual_circuit 5", wire)
	}
	if wire.Role == nil || *wire.Role != "peer" {
		t.Errorf("wire.Role = %v, want peer", wire.Role)
	}
	if wire.Interface != 9 {
		t.Errorf("wire.Interface = %v, want 9", wire.Interface)
	}

	raw, err := json.Marshal(wire)
	if err != nil {
		t.Fatalf("marshal wire: %v", err)
	}
	for _, key := range []string{"virtual_circuit", "role", "interface"} {
		if !jsonHasKey(raw, key) {
			t.Errorf("wire JSON missing key %q: %s", key, raw)
		}
	}
}

func TestWireVirtualCircuitTerminationToDomain(t *testing.T) {
	t.Parallel()

	wire := WireVirtualCircuitTermination{
		ID:             1,
		URL:            "http://nb/api/circuits/virtual-circuit-terminations/1/",
		VirtualCircuit: &WireNested{ID: 5, Name: "VC-1", Slug: "vc-1"},
		Role:           &WireLabel{Value: "peer", Label: "Peer"},
		Interface:      &WireNested{ID: 9, Name: "eth0", Slug: "eth0"},
		Created:        "2024-01-01",
	}
	got := wireVirtualCircuitTerminationToDomain(wire)
	if got.ID != 1 {
		t.Errorf("domain = %+v, want id 1", got)
	}
	if got.VirtualCircuit == nil || got.VirtualCircuit.ID != 5 {
		t.Errorf("domain.VirtualCircuit = %+v, want nested id 5", got.VirtualCircuit)
	}
	if got.Interface == nil || got.Interface.ID != 9 {
		t.Errorf("domain.Interface = %+v, want nested id 9", got.Interface)
	}
	if got.Role == nil || got.Role.Value != "peer" {
		t.Errorf("domain.Role = %+v, want peer", got.Role)
	}
	if got := wireVirtualCircuitTerminationToDomain(WireVirtualCircuitTermination{}).VirtualCircuit; got != nil {
		t.Errorf("empty wire VirtualCircuit = %+v, want nil", got)
	}
}
