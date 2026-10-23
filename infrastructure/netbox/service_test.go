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

// TestClient_ListServices exercises the service list endpoint.
func TestClient_ListServices(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/ipam/services/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"url": "http://nb/api/ipam/services/1/",
					"name": "web",
					"parent_object_type": "dcim.device",
					"parent_object_id": 5,
					"parent": {"id": 5, "name": "sw1", "url": "http://nb/api/dcim/devices/5/"},
					"protocol": {"value": "tcp", "label": "TCP"},
					"ports": [80, 443],
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListServices(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListServices() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "web" {
			t.Errorf("resp = %+v, want count 1 with web", resp)
		}
		if resp.Results[0].Parent == nil || resp.Results[0].Parent.ID != 5 {
			t.Errorf("resp.Results[0].Parent = %+v, want nested id 5", resp.Results[0].Parent)
		}
		if len(resp.Results[0].Ports) != 2 || resp.Results[0].Ports[0] != 80 {
			t.Errorf("resp.Results[0].Ports = %v, want [80 443]", resp.Results[0].Ports)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListServices(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_CreateService(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"name":"web","ports":[80,443],` +
				`"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		parentType := "dcim.device"
		parentID := 5
		protocol := "tcp"
		svc, err := client.CreateService(
			context.Background(), "tok", domain.ServiceWrite{
				Name: "web", ParentObjectType: &parentType, ParentObjectID: &parentID,
				Protocol: &protocol, Ports: []int{80, 443},
			},
		)
		if err != nil {
			t.Fatalf("CreateService() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/ipam/services/" {
			t.Errorf("method/path = %q %q, want POST /api/ipam/services/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "web" {
			t.Errorf("request body = %s, want name web", gotBody)
		}
		ports, ok := reqBody["ports"].([]any)
		if !ok || len(ports) != 2 || ports[0] != float64(80) || ports[1] != float64(443) {
			t.Errorf("request body ports = %v, want [80 443] array", reqBody["ports"])
		}
		if reqBody["parent_object_type"] != "dcim.device" || reqBody["parent_object_id"] != float64(5) {
			t.Errorf("request body parent = %v/%v, want dcim.device/5",
				reqBody["parent_object_type"], reqBody["parent_object_id"])
		}
		if svc.ID != 1 || svc.Name != "web" {
			t.Errorf("service = %+v, want id 1 name web", svc)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"ports":["Enter a valid value."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateService(context.Background(), "tok", domain.ServiceWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateService(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"name":"Renamed","ports":[80],"created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		svc, err := client.UpdateService(context.Background(), "tok", 7, domain.ServiceWrite{Name: "Renamed"})
		if err != nil {
			t.Fatalf("UpdateService() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/ipam/services/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/ipam/services/7/", gotMethod, gotPath)
		}
		if svc.ID != 7 || svc.Name != "Renamed" {
			t.Errorf("service = %+v, want id 7 name Renamed", svc)
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
		_, err := client.UpdateService(context.Background(), "tok", 7, domain.ServiceWrite{Name: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteService(t *testing.T) {
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
		if err := client.DeleteService(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteService() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/ipam/services/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/ipam/services/3/", gotMethod, gotPath)
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
		err := client.DeleteService(context.Background(), "tok", 3)
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
		if err := client.DeleteService(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteService() 404 returned error: %v", err)
		}
	})
}

func TestServiceWriteToWire(t *testing.T) {
	t.Parallel()

	parentType := "dcim.device"
	parentID := 5
	protocol := "tcp"
	desc := "https"
	in := domain.ServiceWrite{
		Name: "web", ParentObjectType: &parentType, ParentObjectID: &parentID,
		Protocol: &protocol, Ports: []int{80, 443}, Description: &desc, Tags: []string{"prod"},
	}
	wire := serviceWriteToWire(in)
	if wire.Name != "web" || wire.ParentObjectType == nil || *wire.ParentObjectType != "dcim.device" ||
		wire.ParentObjectID == nil || *wire.ParentObjectID != 5 {
		t.Errorf("wire = %+v, want name web parent dcim.device/5", wire)
	}
	if wire.Protocol == nil || *wire.Protocol != "tcp" {
		t.Errorf("wire.Protocol = %v, want tcp", wire.Protocol)
	}
	if len(wire.Ports) != 2 || wire.Ports[0] != 80 || wire.Ports[1] != 443 {
		t.Errorf("wire.Ports = %v, want [80 443]", wire.Ports)
	}

	raw, err := json.Marshal(wire)
	if err != nil {
		t.Fatalf("marshal wire: %v", err)
	}
	for _, key := range []string{
		"name", "parent_object_type", "parent_object_id",
		"protocol", "ports", "description", "tags",
	} {
		if !jsonHasKey(raw, key) {
			t.Errorf("wire JSON missing key %q: %s", key, raw)
		}
	}
}

func TestWireServiceToDomain(t *testing.T) {
	t.Parallel()

	wire := WireService{
		ID:               1,
		URL:              "http://nb/api/ipam/services/1/",
		Name:             "web",
		ParentObjectType: "dcim.device",
		ParentObjectID:   5,
		Parent:           &WireNested{ID: 5, Name: "sw1", Slug: "sw1"},
		Protocol:         &WireLabel{Value: "tcp", Label: "TCP"},
		Ports:            []int{80, 443},
		Created:          "2024-01-01",
	}
	got := wireServiceToDomain(wire)
	if got.ID != 1 || got.Name != "web" {
		t.Errorf("domain = %+v, want id 1 name web", got)
	}
	if got.Parent == nil || got.Parent.ID != 5 {
		t.Errorf("domain.Parent = %+v, want nested id 5", got.Parent)
	}
	if got.Protocol == nil || got.Protocol.Value != "tcp" {
		t.Errorf("domain.Protocol = %+v, want tcp", got.Protocol)
	}
	if len(got.Ports) != 2 || got.Ports[1] != 443 {
		t.Errorf("domain.Ports = %v, want [80 443]", got.Ports)
	}
	if got := wireServiceToDomain(WireService{}).Parent; got != nil {
		t.Errorf("empty wire Parent = %+v, want nil", got)
	}
}
