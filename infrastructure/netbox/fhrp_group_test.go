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

// TestClient_ListFHRPGroups exercises the FHRP group list endpoint.
func TestClient_ListFHRPGroups(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/ipam/fhrp-groups/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"url": "http://nb/api/ipam/fhrp-groups/1/",
					"name": "gw",
					"protocol": "vrrp",
					"group_id": 10,
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListFHRPGroups(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListFHRPGroups() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "gw" {
			t.Errorf("resp = %+v, want count 1 with gw", resp)
		}
		if resp.Results[0].GroupID != 10 {
			t.Errorf("resp.Results[0].GroupID = %d, want 10", resp.Results[0].GroupID)
		}
		if resp.Results[0].Protocol == nil || resp.Results[0].Protocol.Value != "vrrp" {
			t.Errorf("resp.Results[0].Protocol = %+v, want vrrp", resp.Results[0].Protocol)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListFHRPGroups(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_CreateFHRPGroup(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"name":"gw","group_id":10,` +
				`"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		protocol := "vrrp"
		groupID := 10
		group, err := client.CreateFHRPGroup(
			context.Background(), "tok", domain.FHRPGroupWrite{Name: "gw", Protocol: &protocol, GroupID: &groupID},
		)
		if err != nil {
			t.Fatalf("CreateFHRPGroup() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/ipam/fhrp-groups/" {
			t.Errorf("method/path = %q %q, want POST /api/ipam/fhrp-groups/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "gw" {
			t.Errorf("request body = %s, want name gw", gotBody)
		}
		if reqBody["group_id"] != float64(10) {
			t.Errorf("request body group_id = %v, want 10", reqBody["group_id"])
		}
		if group.ID != 1 || group.Name != "gw" {
			t.Errorf("fhrp group = %+v, want id 1 name gw", group)
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
		_, err := client.CreateFHRPGroup(context.Background(), "tok", domain.FHRPGroupWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateFHRPGroup(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"name":"Renamed","group_id":20,"created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		group, err := client.UpdateFHRPGroup(context.Background(), "tok", 7, domain.FHRPGroupWrite{Name: "Renamed"})
		if err != nil {
			t.Fatalf("UpdateFHRPGroup() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/ipam/fhrp-groups/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/ipam/fhrp-groups/7/", gotMethod, gotPath)
		}
		if group.ID != 7 || group.Name != "Renamed" {
			t.Errorf("fhrp group = %+v, want id 7 name Renamed", group)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"protocol":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateFHRPGroup(context.Background(), "tok", 7, domain.FHRPGroupWrite{Name: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteFHRPGroup(t *testing.T) {
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
		if err := client.DeleteFHRPGroup(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteFHRPGroup() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/ipam/fhrp-groups/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/ipam/fhrp-groups/3/", gotMethod, gotPath)
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
		err := client.DeleteFHRPGroup(context.Background(), "tok", 3)
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
		if err := client.DeleteFHRPGroup(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteFHRPGroup() 404 returned error: %v", err)
		}
	})
}

func TestFHRPGroupWriteToWire(t *testing.T) {
	t.Parallel()

	protocol := "vrrp"
	groupID := 10
	desc := "gateway"
	in := domain.FHRPGroupWrite{
		Name: "gw", Protocol: &protocol, GroupID: &groupID,
		Description: &desc, Tags: []string{"prod"},
	}
	wire := fhrpGroupWriteToWire(in)
	if wire.Name != "gw" || wire.Protocol == nil || *wire.Protocol != "vrrp" {
		t.Errorf("wire = %+v, want name gw protocol vrrp", wire)
	}
	if wire.GroupID == nil || *wire.GroupID != 10 {
		t.Errorf("wire.GroupID = %v, want 10", wire.GroupID)
	}

	raw, err := json.Marshal(wire)
	if err != nil {
		t.Fatalf("marshal wire: %v", err)
	}
	for _, key := range []string{"name", "protocol", "group_id", "description", "tags"} {
		if !jsonHasKey(raw, key) {
			t.Errorf("wire JSON missing key %q: %s", key, raw)
		}
	}
}

func TestWireFHRPGroupToDomain(t *testing.T) {
	t.Parallel()

	wire := WireFHRPGroup{
		ID:       1,
		URL:      "http://nb/api/ipam/fhrp-groups/1/",
		Name:     "gw",
		Protocol: "vrrp",
		GroupID:  10,
		Created:  "2024-01-01",
	}
	got := wireFHRPGroupToDomain(wire)
	if got.ID != 1 || got.Name != "gw" || got.GroupID != 10 {
		t.Errorf("domain = %+v, want id 1 name gw group_id 10", got)
	}
	if got.Protocol == nil || got.Protocol.Value != "vrrp" {
		t.Errorf("domain.Protocol = %+v, want vrrp", got.Protocol)
	}
}
