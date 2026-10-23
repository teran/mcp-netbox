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

// TestClient_ListASNs exercises the ASN list endpoint.
func TestClient_ListASNs(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/ipam/asns/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"url": "http://nb/api/ipam/asns/1/",
					"asn": 64512,
					"rir": {"id": 3, "name": "ARIN", "url": "http://nb/api/ipam/rirs/3/"},
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListASNs(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListASNs() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].ASN != 64512 {
			t.Errorf("resp = %+v, want count 1 with asn 64512", resp)
		}
		if resp.Results[0].RIR == nil || resp.Results[0].RIR.ID != 3 {
			t.Errorf("resp.Results[0].RIR = %+v, want nested id 3", resp.Results[0].RIR)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListASNs(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestClient_CreateASN(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"asn":64512,"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		rir := 3
		asn, err := client.CreateASN(context.Background(), "tok", domain.ASNWrite{ASN: 64512, RIR: &rir})
		if err != nil {
			t.Fatalf("CreateASN() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/ipam/asns/" {
			t.Errorf("method/path = %q %q, want POST /api/ipam/asns/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["asn"] != float64(64512) {
			t.Errorf("request body asn = %v, want 64512 (as a number)", reqBody["asn"])
		}
		if reqBody["rir"] != float64(3) {
			t.Errorf("request body rir = %v, want 3", reqBody["rir"])
		}
		if asn.ID != 1 || asn.ASN != 64512 {
			t.Errorf("asn = %+v, want id 1 asn 64512", asn)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"asn":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateASN(context.Background(), "tok", domain.ASNWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateASN(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"asn":64513,"created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		asn, err := client.UpdateASN(context.Background(), "tok", 7, domain.ASNWrite{ASN: 64513})
		if err != nil {
			t.Fatalf("UpdateASN() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/ipam/asns/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/ipam/asns/7/", gotMethod, gotPath)
		}
		if asn.ID != 7 || asn.ASN != 64513 {
			t.Errorf("asn = %+v, want id 7 asn 64513", asn)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"rir":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateASN(context.Background(), "tok", 7, domain.ASNWrite{ASN: 64513})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteASN(t *testing.T) {
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
		if err := client.DeleteASN(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteASN() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/ipam/asns/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/ipam/asns/3/", gotMethod, gotPath)
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
		err := client.DeleteASN(context.Background(), "tok", 3)
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
		if err := client.DeleteASN(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteASN() 404 returned error: %v", err)
		}
	})
}

func TestASNWriteToWire(t *testing.T) {
	t.Parallel()

	rir := 3
	tenant := 4
	desc := "edge"
	in := domain.ASNWrite{ASN: 64512, RIR: &rir, Tenant: &tenant, Description: &desc, Tags: []string{"prod"}}
	wire := asnWriteToWire(in)
	if wire.ASN != 64512 || wire.RIR == nil || *wire.RIR != 3 || wire.Tenant == nil || *wire.Tenant != 4 {
		t.Errorf("wire = %+v, want asn 64512 rir 3 tenant 4", wire)
	}

	raw, err := json.Marshal(wire)
	if err != nil {
		t.Fatalf("marshal wire: %v", err)
	}
	for _, key := range []string{"asn", "rir", "tenant", "description", "tags"} {
		if !jsonHasKey(raw, key) {
			t.Errorf("wire JSON missing key %q: %s", key, raw)
		}
	}
}

func TestWireASNToDomain(t *testing.T) {
	t.Parallel()

	wire := WireASN{
		ID:      1,
		URL:     "http://nb/api/ipam/asns/1/",
		ASN:     64512,
		RIR:     &WireNested{ID: 3, Name: "ARIN", Slug: "arin"},
		Tenant:  &WireNested{ID: 4, Name: "ACME", Slug: "acme"},
		Created: "2024-01-01",
	}
	got := wireASNToDomain(wire)
	if got.ID != 1 || got.ASN != 64512 {
		t.Errorf("domain = %+v, want id 1 asn 64512", got)
	}
	if got.RIR == nil || got.RIR.ID != 3 {
		t.Errorf("domain.RIR = %+v, want nested id 3", got.RIR)
	}
	if got.Tenant == nil || got.Tenant.ID != 4 {
		t.Errorf("domain.Tenant = %+v, want nested id 4", got.Tenant)
	}
	if got := wireASNToDomain(WireASN{}).RIR; got != nil {
		t.Errorf("empty wire RIR = %+v, want nil", got)
	}
}
