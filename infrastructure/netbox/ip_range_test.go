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

// TestClient_ListIPRanges exercises the IP range list endpoint.
func TestClient_ListIPRanges(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/ipam/ip-ranges/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"url": "http://nb/api/ipam/ip-ranges/1/",
					"start_address": "10.0.0.1",
					"end_address": "10.0.0.10",
					"status": {"value": "active", "label": "Active"},
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListIPRanges(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListIPRanges() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 {
			t.Errorf("resp = %+v, want count 1", resp)
		}
		if resp.Results[0].StartAddress != "10.0.0.1" || resp.Results[0].EndAddress != "10.0.0.10" {
			t.Errorf("resp.Results[0] = %+v, want range 10.0.0.1-10.0.0.10", resp.Results[0])
		}
		if resp.Results[0].Status == nil || resp.Results[0].Status.Value != "active" {
			t.Errorf("resp.Results[0].Status = %+v, want active", resp.Results[0].Status)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListIPRanges(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_CreateIPRange(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"start_address":"10.0.0.1","end_address":"10.0.0.10",` +
				`"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		vrf := 2
		ipRange, err := client.CreateIPRange(
			context.Background(), "tok", domain.IPRangeWrite{StartAddress: "10.0.0.1", EndAddress: "10.0.0.10", VRF: &vrf},
		)
		if err != nil {
			t.Fatalf("CreateIPRange() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/ipam/ip-ranges/" {
			t.Errorf("method/path = %q %q, want POST /api/ipam/ip-ranges/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["start_address"] != "10.0.0.1" || reqBody["end_address"] != "10.0.0.10" {
			t.Errorf("request body = %s, want start/end addresses", gotBody)
		}
		if reqBody["vrf"] != float64(2) {
			t.Errorf("request body vrf = %v, want 2", reqBody["vrf"])
		}
		if ipRange.ID != 1 || ipRange.StartAddress != "10.0.0.1" || ipRange.EndAddress != "10.0.0.10" {
			t.Errorf("ip range = %+v, want id 1 range", ipRange)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"start_address":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateIPRange(context.Background(), "tok", domain.IPRangeWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateIPRange(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"start_address":"10.0.0.1","end_address":"10.0.0.10",` +
				`"created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		ipRange, err := client.UpdateIPRange(
			context.Background(), "tok", 7, domain.IPRangeWrite{StartAddress: "10.0.0.1"},
		)
		if err != nil {
			t.Fatalf("UpdateIPRange() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/ipam/ip-ranges/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/ipam/ip-ranges/7/", gotMethod, gotPath)
		}
		if ipRange.ID != 7 || ipRange.StartAddress != "10.0.0.1" {
			t.Errorf("ip range = %+v, want id 7", ipRange)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"end_address":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateIPRange(context.Background(), "tok", 7, domain.IPRangeWrite{StartAddress: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteIPRange(t *testing.T) {
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
		if err := client.DeleteIPRange(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteIPRange() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/ipam/ip-ranges/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/ipam/ip-ranges/3/", gotMethod, gotPath)
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
		err := client.DeleteIPRange(context.Background(), "tok", 3)
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
		if err := client.DeleteIPRange(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteIPRange() 404 returned error: %v", err)
		}
	})
}

func TestIPRangeWriteToWire(t *testing.T) {
	t.Parallel()

	vrf := 2
	tenant := 4
	status := "active"
	role := 6
	desc := "range"
	in := domain.IPRangeWrite{
		StartAddress: "10.0.0.1", EndAddress: "10.0.0.10", VRF: &vrf, Tenant: &tenant,
		Status: &status, Role: &role, Description: &desc, Tags: []string{"prod"},
	}
	wire := ipRangeWriteToWire(in)
	if wire.StartAddress != "10.0.0.1" || wire.EndAddress != "10.0.0.10" {
		t.Errorf("wire = %+v, want start/end addresses", wire)
	}
	if wire.VRF == nil || *wire.VRF != 2 ||
		wire.Status == nil || *wire.Status != "active" ||
		wire.Role == nil || *wire.Role != 6 {
		t.Errorf("wire = %+v, want vrf 2 status active role 6", wire)
	}

	raw, err := json.Marshal(wire)
	if err != nil {
		t.Fatalf("marshal wire: %v", err)
	}
	for _, key := range []string{
		"start_address", "end_address", "vrf", "tenant",
		"status", "role", "description", "tags",
	} {
		if !jsonHasKey(raw, key) {
			t.Errorf("wire JSON missing key %q: %s", key, raw)
		}
	}
}

func TestWireIPRangeToDomain(t *testing.T) {
	t.Parallel()

	wire := WireIPRange{
		ID:           1,
		URL:          "http://nb/api/ipam/ip-ranges/1/",
		StartAddress: "10.0.0.1",
		EndAddress:   "10.0.0.10",
		VRF:          &WireNested{ID: 2, Name: "vrf1", Slug: "vrf1"},
		Status:       &WireLabel{Value: "active", Label: "Active"},
		Role:         &WireNested{ID: 6, Name: "service", Slug: "service"},
		Created:      "2024-01-01",
	}
	got := wireIPRangeToDomain(wire)
	if got.ID != 1 || got.StartAddress != "10.0.0.1" || got.EndAddress != "10.0.0.10" {
		t.Errorf("domain = %+v, want id 1 range", got)
	}
	if got.VRF == nil || got.VRF.ID != 2 {
		t.Errorf("domain.VRF = %+v, want nested id 2", got.VRF)
	}
	if got.Status == nil || got.Status.Value != "active" {
		t.Errorf("domain.Status = %+v, want active", got.Status)
	}
	if got.Role == nil || got.Role.ID != 6 {
		t.Errorf("domain.Role = %+v, want nested id 6", got.Role)
	}
	if got := wireIPRangeToDomain(WireIPRange{}).VRF; got != nil {
		t.Errorf("empty wire VRF = %+v, want nil", got)
	}
}
