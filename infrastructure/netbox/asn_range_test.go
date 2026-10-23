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

// TestClient_ListASNRanges exercises the ASN range list endpoint.
func TestClient_ListASNRanges(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/ipam/asn-ranges/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"url": "http://nb/api/ipam/asn-ranges/1/",
					"name": "PRIVATE",
					"start": 64512,
					"end": 64520,
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListASNRanges(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListASNRanges() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "PRIVATE" {
			t.Errorf("resp = %+v, want count 1 with PRIVATE", resp)
		}
		if resp.Results[0].Start != 64512 || resp.Results[0].End != 64520 {
			t.Errorf("resp.Results[0] = %+v, want start 64512 end 64520", resp.Results[0])
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListASNRanges(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_CreateASNRange(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"name":"PRIVATE","start":64512,"end":64520,` +
				`"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		start := 64512
		end := 64520
		asnRange, err := client.CreateASNRange(
			context.Background(), "tok", domain.ASNRangeWrite{Name: "PRIVATE", Start: &start, End: &end},
		)
		if err != nil {
			t.Fatalf("CreateASNRange() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/ipam/asn-ranges/" {
			t.Errorf("method/path = %q %q, want POST /api/ipam/asn-ranges/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "PRIVATE" {
			t.Errorf("request body = %s, want name PRIVATE", gotBody)
		}
		if reqBody["start"] != float64(64512) || reqBody["end"] != float64(64520) {
			t.Errorf("request body start/end = %v/%v, want 64512/64520", reqBody["start"], reqBody["end"])
		}
		if asnRange.ID != 1 || asnRange.Name != "PRIVATE" {
			t.Errorf("asn range = %+v, want id 1 name PRIVATE", asnRange)
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
		_, err := client.CreateASNRange(context.Background(), "tok", domain.ASNRangeWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateASNRange(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"name":"Renamed","start":1,"end":2,"created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		asnRange, err := client.UpdateASNRange(context.Background(), "tok", 7, domain.ASNRangeWrite{Name: "Renamed"})
		if err != nil {
			t.Fatalf("UpdateASNRange() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/ipam/asn-ranges/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/ipam/asn-ranges/7/", gotMethod, gotPath)
		}
		if asnRange.ID != 7 || asnRange.Name != "Renamed" {
			t.Errorf("asn range = %+v, want id 7 name Renamed", asnRange)
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
		_, err := client.UpdateASNRange(context.Background(), "tok", 7, domain.ASNRangeWrite{Name: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteASNRange(t *testing.T) {
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
		if err := client.DeleteASNRange(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteASNRange() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/ipam/asn-ranges/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/ipam/asn-ranges/3/", gotMethod, gotPath)
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
		err := client.DeleteASNRange(context.Background(), "tok", 3)
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
		if err := client.DeleteASNRange(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteASNRange() 404 returned error: %v", err)
		}
	})
}

func TestASNRangeWriteToWire(t *testing.T) {
	t.Parallel()

	start := 64512
	end := 64520
	rir := 3
	tenant := 4
	desc := "range"
	slug := "private"
	in := domain.ASNRangeWrite{
		Name: "PRIVATE", Slug: &slug, Start: &start, End: &end, RIR: &rir,
		Tenant: &tenant, Description: &desc, Tags: []string{"prod"},
	}
	wire := asnRangeWriteToWire(in)
	if wire.Name != "PRIVATE" || wire.Slug == nil || *wire.Slug != "private" ||
		wire.Start == nil || *wire.Start != 64512 || wire.End == nil || *wire.End != 64520 {
		t.Errorf("wire = %+v, want name PRIVATE slug private start 64512 end 64520", wire)
	}
	if wire.RIR == nil || *wire.RIR != 3 {
		t.Errorf("wire.RIR = %v, want 3", wire.RIR)
	}

	raw, err := json.Marshal(wire)
	if err != nil {
		t.Fatalf("marshal wire: %v", err)
	}
	for _, key := range []string{"name", "slug", "start", "end", "rir", "tenant", "description", "tags"} {
		if !jsonHasKey(raw, key) {
			t.Errorf("wire JSON missing key %q: %s", key, raw)
		}
	}
}

func TestWireASNRangeToDomain(t *testing.T) {
	t.Parallel()

	wire := WireASNRange{
		ID:      1,
		URL:     "http://nb/api/ipam/asn-ranges/1/",
		Name:    "PRIVATE",
		Start:   64512,
		End:     64520,
		RIR:     &WireNested{ID: 3, Name: "ARIN", Slug: "arin"},
		Tenant:  &WireNested{ID: 4, Name: "ACME", Slug: "acme"},
		Created: "2024-01-01",
	}
	got := wireASNRangeToDomain(wire)
	if got.ID != 1 || got.Name != "PRIVATE" || got.Start != 64512 || got.End != 64520 {
		t.Errorf("domain = %+v, want id 1 PRIVATE 64512-64520", got)
	}
	if got.RIR == nil || got.RIR.ID != 3 {
		t.Errorf("domain.RIR = %+v, want nested id 3", got.RIR)
	}
	if got.Tenant == nil || got.Tenant.ID != 4 {
		t.Errorf("domain.Tenant = %+v, want nested id 4", got.Tenant)
	}
	if got := wireASNRangeToDomain(WireASNRange{}).RIR; got != nil {
		t.Errorf("empty wire RIR = %+v, want nil", got)
	}
}
