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

// TestClient_ListCircuitGroupAssignments exercises the circuit group assignment list endpoint.
func TestClient_ListCircuitGroupAssignments(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/circuits/circuit-group-assignments/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"url": "http://nb/api/circuits/circuit-group-assignments/1/",
					"group": {"id": 5, "name": "Prod", "url": "http://nb/api/circuits/circuit-groups/5/"},
					"member": {"id": 9, "name": "CKT-001", "url": "http://nb/api/circuits/circuits/9/"},
					"priority": {"value": "primary", "label": "Primary"},
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListCircuitGroupAssignments(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListCircuitGroupAssignments() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 {
			t.Errorf("resp = %+v, want count 1", resp)
		}
		if resp.Results[0].Group == nil || resp.Results[0].Group.ID != 5 {
			t.Errorf("resp.Results[0].Group = %+v, want nested id 5", resp.Results[0].Group)
		}
		if resp.Results[0].Member == nil || resp.Results[0].Member.ID != 9 {
			t.Errorf("resp.Results[0].Member = %+v, want nested id 9", resp.Results[0].Member)
		}
		if resp.Results[0].Priority == nil || resp.Results[0].Priority.Value != "primary" {
			t.Errorf("resp.Results[0].Priority = %+v, want primary", resp.Results[0].Priority)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListCircuitGroupAssignments(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_CreateCircuitGroupAssignment(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"group":{"id":5},"member":{"id":9},` +
				`"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		priority := "primary"
		assn, err := client.CreateCircuitGroupAssignment(
			context.Background(), "tok", domain.CircuitGroupAssignmentWrite{
				Group: 5, MemberType: "circuits.circuit", MemberID: 9, Priority: &priority,
			},
		)
		if err != nil {
			t.Fatalf("CreateCircuitGroupAssignment() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/circuits/circuit-group-assignments/" {
			t.Errorf("method/path = %q %q, want POST /api/circuits/circuit-group-assignments/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["group"] != float64(5) ||
			reqBody["member_type"] != "circuits.circuit" || reqBody["member_id"] != float64(9) {
			t.Errorf("request body = %s, want group 5 member circuits.circuit/9", gotBody)
		}
		if reqBody["priority"] != "primary" {
			t.Errorf("request body priority = %v, want primary", reqBody["priority"])
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
		_, err := client.CreateCircuitGroupAssignment(context.Background(), "tok", domain.CircuitGroupAssignmentWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateCircuitGroupAssignment(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"group":{"id":5},"member":{"id":9},` +
				`"created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		assn, err := client.UpdateCircuitGroupAssignment(
			context.Background(), "tok", 7, domain.CircuitGroupAssignmentWrite{Group: 5},
		)
		if err != nil {
			t.Fatalf("UpdateCircuitGroupAssignment() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/circuits/circuit-group-assignments/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/circuits/circuit-group-assignments/7/", gotMethod, gotPath)
		}
		if assn.ID != 7 {
			t.Errorf("assignment = %+v, want id 7", assn)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"circuit":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateCircuitGroupAssignment(
			context.Background(), "tok", 7, domain.CircuitGroupAssignmentWrite{Group: 5},
		)
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteCircuitGroupAssignment(t *testing.T) {
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
		if err := client.DeleteCircuitGroupAssignment(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteCircuitGroupAssignment() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/circuits/circuit-group-assignments/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/circuits/circuit-group-assignments/3/", gotMethod, gotPath)
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
		err := client.DeleteCircuitGroupAssignment(context.Background(), "tok", 3)
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
		if err := client.DeleteCircuitGroupAssignment(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteCircuitGroupAssignment() 404 returned error: %v", err)
		}
	})
}

func TestCircuitGroupAssignmentWriteToWire(t *testing.T) {
	t.Parallel()

	priority := "primary"
	in := domain.CircuitGroupAssignmentWrite{Group: 5, MemberType: "circuits.circuit", MemberID: 9, Priority: &priority}
	wire := circuitGroupAssignmentWriteToWire(in)
	if wire.Group != 5 || wire.MemberType != "circuits.circuit" || wire.MemberID != 9 {
		t.Errorf("wire = %+v, want group 5 member circuits.circuit/9", wire)
	}
	if wire.Priority == nil || *wire.Priority != "primary" {
		t.Errorf("wire.Priority = %v, want primary", wire.Priority)
	}

	raw, err := json.Marshal(wire)
	if err != nil {
		t.Fatalf("marshal wire: %v", err)
	}
	for _, key := range []string{"group", "member_type", "member_id", "priority"} {
		if !jsonHasKey(raw, key) {
			t.Errorf("wire JSON missing key %q: %s", key, raw)
		}
	}
}

func TestWireCircuitGroupAssignmentToDomain(t *testing.T) {
	t.Parallel()

	wire := WireCircuitGroupAssignment{
		ID:       1,
		URL:      "http://nb/api/circuits/circuit-group-assignments/1/",
		Group:    &WireNested{ID: 5, Name: "Prod", Slug: "prod"},
		Member:   &WireNested{ID: 9, Name: "CKT-001", Slug: "ckt-001"},
		Priority: &WireLabel{Value: "primary", Label: "Primary"},
		Created:  "2024-01-01",
	}
	got := wireCircuitGroupAssignmentToDomain(wire)
	if got.ID != 1 || got.Priority == nil || got.Priority.Value != "primary" {
		t.Errorf("domain = %+v, want id 1 priority primary", got)
	}
	if got.Group == nil || got.Group.ID != 5 {
		t.Errorf("domain.Group = %+v, want nested id 5", got.Group)
	}
	if got.Member == nil || got.Member.ID != 9 {
		t.Errorf("domain.Member = %+v, want nested id 9", got.Member)
	}
	if got := wireCircuitGroupAssignmentToDomain(WireCircuitGroupAssignment{}).Group; got != nil {
		t.Errorf("empty wire Group = %+v, want nil", got)
	}
}
