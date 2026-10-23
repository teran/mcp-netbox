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

// TestClient_ListFHRPGroupAssignments exercises the FHRP group assignment list endpoint.
func TestClient_ListFHRPGroupAssignments(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/ipam/fhrp-group-assignments/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"url": "http://nb/api/ipam/fhrp-group-assignments/1/",
					"group": {"id": 5, "name": "gw", "url": "http://nb/api/ipam/fhrp-groups/5/"},
					"interface_type": "dcim.interface",
					"interface_id": 9,
					"interface": {"id": 9, "name": "eth0", "url": "http://nb/api/dcim/interfaces/9/"},
					"priority": 100,
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListFHRPGroupAssignments(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListFHRPGroupAssignments() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 {
			t.Errorf("resp = %+v, want count 1", resp)
		}
		if resp.Results[0].Group == nil || resp.Results[0].Group.ID != 5 {
			t.Errorf("resp.Results[0].Group = %+v, want nested id 5", resp.Results[0].Group)
		}
		if resp.Results[0].InterfaceID != 9 {
			t.Errorf("resp.Results[0].InterfaceID = %d, want 9", resp.Results[0].InterfaceID)
		}
		if resp.Results[0].Priority != 100 {
			t.Errorf("resp.Results[0].Priority = %d, want 100", resp.Results[0].Priority)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListFHRPGroupAssignments(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_CreateFHRPGroupAssignment(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"group":{"id":5},"interface_type":"dcim.interface",` +
				`"interface_id":9,"interface":{"id":9,"name":"eth0"},` +
				`"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		priority := 100
		assn, err := client.CreateFHRPGroupAssignment(
			context.Background(), "tok", domain.FHRPGroupAssignmentWrite{
				Group: 5, InterfaceType: "dcim.interface", InterfaceID: 9, Priority: &priority,
			},
		)
		if err != nil {
			t.Fatalf("CreateFHRPGroupAssignment() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/ipam/fhrp-group-assignments/" {
			t.Errorf("method/path = %q %q, want POST /api/ipam/fhrp-group-assignments/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["group"] != float64(5) || reqBody["interface_id"] != float64(9) {
			t.Errorf("request body = %s, want group 5 interface_id 9", gotBody)
		}
		if reqBody["interface_type"] != "dcim.interface" {
			t.Errorf("request body interface_type = %v, want dcim.interface", reqBody["interface_type"])
		}
		if reqBody["priority"] != float64(100) {
			t.Errorf("request body priority = %v, want 100", reqBody["priority"])
		}
		if assn.ID != 1 {
			t.Errorf("assignment = %+v, want id 1", assn)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"group":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateFHRPGroupAssignment(context.Background(), "tok", domain.FHRPGroupAssignmentWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateFHRPGroupAssignment(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"group":{"id":5},"interface_type":"dcim.interface",` +
				`"interface_id":9,"interface":{"id":9,"name":"eth0"},` +
				`"created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		assn, err := client.UpdateFHRPGroupAssignment(
			context.Background(), "tok", 7, domain.FHRPGroupAssignmentWrite{Group: 5},
		)
		if err != nil {
			t.Fatalf("UpdateFHRPGroupAssignment() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/ipam/fhrp-group-assignments/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/ipam/fhrp-group-assignments/7/", gotMethod, gotPath)
		}
		if assn.ID != 7 {
			t.Errorf("assignment = %+v, want id 7", assn)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"interface_id":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateFHRPGroupAssignment(context.Background(), "tok", 7, domain.FHRPGroupAssignmentWrite{Group: 5})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteFHRPGroupAssignment(t *testing.T) {
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
		if err := client.DeleteFHRPGroupAssignment(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteFHRPGroupAssignment() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/ipam/fhrp-group-assignments/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/ipam/fhrp-group-assignments/3/", gotMethod, gotPath)
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
		err := client.DeleteFHRPGroupAssignment(context.Background(), "tok", 3)
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
		if err := client.DeleteFHRPGroupAssignment(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteFHRPGroupAssignment() 404 returned error: %v", err)
		}
	})
}

func TestFHRPGroupAssignmentWriteToWire(t *testing.T) {
	t.Parallel()

	priority := 100
	in := domain.FHRPGroupAssignmentWrite{Group: 5, InterfaceType: "dcim.interface", InterfaceID: 9, Priority: &priority}
	wire := fhrpGroupAssignmentWriteToWire(in)
	if wire.Group != 5 || wire.InterfaceID != 9 {
		t.Errorf("wire = %+v, want group 5 interface_id 9", wire)
	}
	if wire.InterfaceType != "dcim.interface" {
		t.Errorf("wire.InterfaceType = %q, want dcim.interface", wire.InterfaceType)
	}
	if wire.Priority == nil || *wire.Priority != 100 {
		t.Errorf("wire.Priority = %v, want 100", wire.Priority)
	}

	raw, err := json.Marshal(wire)
	if err != nil {
		t.Fatalf("marshal wire: %v", err)
	}
	for _, key := range []string{"group", "interface_type", "interface_id", "priority"} {
		if !jsonHasKey(raw, key) {
			t.Errorf("wire JSON missing key %q: %s", key, raw)
		}
	}
}

func TestWireFHRPGroupAssignmentToDomain(t *testing.T) {
	t.Parallel()

	wire := WireFHRPGroupAssignment{
		ID:            1,
		URL:           "http://nb/api/ipam/fhrp-group-assignments/1/",
		Group:         &WireNested{ID: 5, Name: "gw", Slug: "gw"},
		InterfaceType: "dcim.interface",
		InterfaceID:   9,
		Interface:     &WireNested{ID: 9, Name: "eth0", Slug: "eth0"},
		Priority:      100,
		Created:       "2024-01-01",
	}
	got := wireFHRPGroupAssignmentToDomain(wire)
	if got.ID != 1 || got.Priority != 100 {
		t.Errorf("domain = %+v, want id 1 priority 100", got)
	}
	if got.Group == nil || got.Group.ID != 5 {
		t.Errorf("domain.Group = %+v, want nested id 5", got.Group)
	}
	if got.InterfaceID != 9 || got.InterfaceType != "dcim.interface" {
		t.Errorf("domain.Interface = %s/%d, want dcim.interface/9", got.InterfaceType, got.InterfaceID)
	}
	if got.Interface == nil || got.Interface.ID != 9 {
		t.Errorf("domain.Interface = %+v, want nested id 9", got.Interface)
	}
	if got := wireFHRPGroupAssignmentToDomain(WireFHRPGroupAssignment{}).Group; got != nil {
		t.Errorf("empty wire Group = %+v, want nil", got)
	}
}
