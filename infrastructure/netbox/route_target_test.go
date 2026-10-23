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

// TestClient_ListRouteTargets exercises the route target list endpoint.
func TestClient_ListRouteTargets(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/ipam/route-targets/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"url": "http://nb/api/ipam/route-targets/1/",
					"name": "65000:100",
					"tenant": {"id": 4, "name": "ACME", "url": "http://nb/api/tenancy/tenants/4/"},
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListRouteTargets(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListRouteTargets() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "65000:100" {
			t.Errorf("resp = %+v, want count 1 with 65000:100", resp)
		}
		if resp.Results[0].Tenant == nil || resp.Results[0].Tenant.ID != 4 {
			t.Errorf("resp.Results[0].Tenant = %+v, want nested id 4", resp.Results[0].Tenant)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListRouteTargets(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestClient_CreateRouteTarget(t *testing.T) { //nolint:funlen
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
			_, _ = writer.Write([]byte(`{"id":1,"name":"65000:100","created":"2024-01-01",` +
				`"last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		tenant := 4
		routeTarget, err := client.CreateRouteTarget(
			context.Background(), "tok", domain.RouteTargetWrite{Name: "65000:100", Tenant: &tenant},
		)
		if err != nil {
			t.Fatalf("CreateRouteTarget() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/ipam/route-targets/" {
			t.Errorf("method/path = %q %q, want POST /api/ipam/route-targets/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "65000:100" {
			t.Errorf("request body = %s, want name 65000:100", gotBody)
		}
		if reqBody["tenant"] != float64(4) {
			t.Errorf("request body tenant = %v, want 4", reqBody["tenant"])
		}
		if routeTarget.ID != 1 || routeTarget.Name != "65000:100" {
			t.Errorf("route target = %+v, want id 1 name 65000:100", routeTarget)
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
		_, err := client.CreateRouteTarget(context.Background(), "tok", domain.RouteTargetWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateRouteTarget(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"name":"65000:200","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		routeTarget, err := client.UpdateRouteTarget(
			context.Background(), "tok", 7, domain.RouteTargetWrite{Name: "65000:200"},
		)
		if err != nil {
			t.Fatalf("UpdateRouteTarget() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/ipam/route-targets/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/ipam/route-targets/7/", gotMethod, gotPath)
		}
		if routeTarget.ID != 7 || routeTarget.Name != "65000:200" {
			t.Errorf("route target = %+v, want id 7 name 65000:200", routeTarget)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"tenant":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateRouteTarget(context.Background(), "tok", 7, domain.RouteTargetWrite{Name: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteRouteTarget(t *testing.T) {
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
		if err := client.DeleteRouteTarget(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteRouteTarget() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/ipam/route-targets/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/ipam/route-targets/3/", gotMethod, gotPath)
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
		err := client.DeleteRouteTarget(context.Background(), "tok", 3)
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
		if err := client.DeleteRouteTarget(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteRouteTarget() 404 returned error: %v", err)
		}
	})
}

func TestRouteTargetWriteToWire(t *testing.T) {
	t.Parallel()

	tenant := 4
	desc := "vpn"
	in := domain.RouteTargetWrite{Name: "65000:100", Tenant: &tenant, Description: &desc, Tags: []string{"prod"}}
	wire := routeTargetWriteToWire(in)
	if wire.Name != "65000:100" || wire.Tenant == nil || *wire.Tenant != 4 {
		t.Errorf("wire = %+v, want name 65000:100 tenant 4", wire)
	}

	raw, err := json.Marshal(wire)
	if err != nil {
		t.Fatalf("marshal wire: %v", err)
	}
	for _, key := range []string{"name", "tenant", "description", "tags"} {
		if !jsonHasKey(raw, key) {
			t.Errorf("wire JSON missing key %q: %s", key, raw)
		}
	}
}

func TestWireRouteTargetToDomain(t *testing.T) {
	t.Parallel()

	wire := WireRouteTarget{
		ID:      1,
		URL:     "http://nb/api/ipam/route-targets/1/",
		Name:    "65000:100",
		Tenant:  &WireNested{ID: 4, Name: "ACME", Slug: "acme"},
		Created: "2024-01-01",
	}
	got := wireRouteTargetToDomain(wire)
	if got.ID != 1 || got.Name != "65000:100" {
		t.Errorf("domain = %+v, want id 1 name 65000:100", got)
	}
	if got.Tenant == nil || got.Tenant.ID != 4 {
		t.Errorf("domain.Tenant = %+v, want nested id 4", got.Tenant)
	}
	if got := wireRouteTargetToDomain(WireRouteTarget{}).Tenant; got != nil {
		t.Errorf("empty wire Tenant = %+v, want nil", got)
	}
}
