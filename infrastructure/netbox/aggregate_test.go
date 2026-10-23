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

// TestClient_ListAggregates exercises the aggregate list endpoint.
func TestClient_ListAggregates(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/ipam/aggregates/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"url": "http://nb/api/ipam/aggregates/1/",
					"prefix": "10.0.0.0/8",
					"rir": {"id": 3, "name": "ARIN", "url": "http://nb/api/ipam/rirs/3/"},
					"date_added": "2024-01-01",
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListAggregates(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListAggregates() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Prefix != "10.0.0.0/8" {
			t.Errorf("resp = %+v, want count 1 with 10.0.0.0/8", resp)
		}
		if resp.Results[0].RIR == nil || resp.Results[0].RIR.ID != 3 {
			t.Errorf("resp.Results[0].RIR = %+v, want nested id 3", resp.Results[0].RIR)
		}
		if resp.Results[0].DateAdded != "2024-01-01" {
			t.Errorf("resp.Results[0].DateAdded = %q, want 2024-01-01", resp.Results[0].DateAdded)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListAggregates(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_CreateAggregate(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"prefix":"10.0.0.0/8","created":"2024-01-01",` +
				`"last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		rir := 3
		aggregate, err := client.CreateAggregate(
			context.Background(), "tok", domain.AggregateWrite{Prefix: "10.0.0.0/8", RIR: &rir},
		)
		if err != nil {
			t.Fatalf("CreateAggregate() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/ipam/aggregates/" {
			t.Errorf("method/path = %q %q, want POST /api/ipam/aggregates/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["prefix"] != "10.0.0.0/8" {
			t.Errorf("request body = %s, want prefix 10.0.0.0/8", gotBody)
		}
		if reqBody["rir"] != float64(3) {
			t.Errorf("request body rir = %v, want 3", reqBody["rir"])
		}
		if aggregate.ID != 1 || aggregate.Prefix != "10.0.0.0/8" {
			t.Errorf("aggregate = %+v, want id 1 prefix 10.0.0.0/8", aggregate)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"prefix":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateAggregate(context.Background(), "tok", domain.AggregateWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateAggregate(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"prefix":"10.0.0.0/8","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		aggregate, err := client.UpdateAggregate(context.Background(), "tok", 7, domain.AggregateWrite{Prefix: "10.0.0.0/8"})
		if err != nil {
			t.Fatalf("UpdateAggregate() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/ipam/aggregates/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/ipam/aggregates/7/", gotMethod, gotPath)
		}
		if aggregate.ID != 7 || aggregate.Prefix != "10.0.0.0/8" {
			t.Errorf("aggregate = %+v, want id 7 prefix 10.0.0.0/8", aggregate)
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
		_, err := client.UpdateAggregate(context.Background(), "tok", 7, domain.AggregateWrite{Prefix: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteAggregate(t *testing.T) {
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
		if err := client.DeleteAggregate(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteAggregate() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/ipam/aggregates/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/ipam/aggregates/3/", gotMethod, gotPath)
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
		err := client.DeleteAggregate(context.Background(), "tok", 3)
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
		if err := client.DeleteAggregate(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteAggregate() 404 returned error: %v", err)
		}
	})
}

func TestAggregateWriteToWire(t *testing.T) {
	t.Parallel()

	rir := 3
	tenant := 4
	dateAdded := "2024-01-01"
	desc := "core"
	in := domain.AggregateWrite{
		Prefix: "10.0.0.0/8", RIR: &rir, Tenant: &tenant, DateAdded: &dateAdded, Description: &desc, Tags: []string{"prod"},
	}
	wire := aggregateWriteToWire(in)
	if wire.Prefix != "10.0.0.0/8" || wire.RIR == nil || *wire.RIR != 3 || wire.Tenant == nil || *wire.Tenant != 4 {
		t.Errorf("wire = %+v, want prefix 10.0.0.0/8 rir 3 tenant 4", wire)
	}
	if wire.DateAdded == nil || *wire.DateAdded != "2024-01-01" {
		t.Errorf("wire.DateAdded = %v, want 2024-01-01", wire.DateAdded)
	}

	raw, err := json.Marshal(wire)
	if err != nil {
		t.Fatalf("marshal wire: %v", err)
	}
	for _, key := range []string{"prefix", "rir", "tenant", "date_added", "description", "tags"} {
		if !jsonHasKey(raw, key) {
			t.Errorf("wire JSON missing key %q: %s", key, raw)
		}
	}
}

func TestWireAggregateToDomain(t *testing.T) {
	t.Parallel()

	wire := WireAggregate{
		ID:        1,
		URL:       "http://nb/api/ipam/aggregates/1/",
		Prefix:    "10.0.0.0/8",
		RIR:       &WireNested{ID: 3, Name: "ARIN", Slug: "arin"},
		Tenant:    &WireNested{ID: 4, Name: "ACME", Slug: "acme"},
		DateAdded: "2024-01-01",
		Created:   "2024-01-01",
	}
	got := wireAggregateToDomain(wire)
	if got.ID != 1 || got.Prefix != "10.0.0.0/8" || got.DateAdded != "2024-01-01" {
		t.Errorf("domain = %+v, want id 1 prefix 10.0.0.0/8 date_added 2024-01-01", got)
	}
	if got.RIR == nil || got.RIR.ID != 3 {
		t.Errorf("domain.RIR = %+v, want nested id 3", got.RIR)
	}
	if got.Tenant == nil || got.Tenant.ID != 4 {
		t.Errorf("domain.Tenant = %+v, want nested id 4", got.Tenant)
	}
	if got := wireAggregateToDomain(WireAggregate{}).RIR; got != nil {
		t.Errorf("empty wire RIR = %+v, want nil", got)
	}
}
