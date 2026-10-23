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

// TestClient_ListServiceTemplates exercises the service template list endpoint.
func TestClient_ListServiceTemplates(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/ipam/service-templates/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"url": "http://nb/api/ipam/service-templates/1/",
					"name": "web",
					"protocol": {"value": "tcp", "label": "TCP"},
					"ports": [80, 443],
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListServiceTemplates(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListServiceTemplates() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "web" {
			t.Errorf("resp = %+v, want count 1 with web", resp)
		}
		if resp.Results[0].Protocol == nil || resp.Results[0].Protocol.Value != "tcp" {
			t.Errorf("resp.Results[0].Protocol = %+v, want tcp", resp.Results[0].Protocol)
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
		_, err := client.ListServiceTemplates(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_CreateServiceTemplate(t *testing.T) {
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
		protocol := "tcp"
		tmpl, err := client.CreateServiceTemplate(
			context.Background(), "tok", domain.ServiceTemplateWrite{Name: "web", Protocol: &protocol, Ports: []int{80, 443}},
		)
		if err != nil {
			t.Fatalf("CreateServiceTemplate() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/ipam/service-templates/" {
			t.Errorf("method/path = %q %q, want POST /api/ipam/service-templates/", gotMethod, gotPath)
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
		if !ok || len(ports) != 2 || ports[0] != float64(80) {
			t.Errorf("request body ports = %v, want [80 443] array", reqBody["ports"])
		}
		if tmpl.ID != 1 || tmpl.Name != "web" {
			t.Errorf("service template = %+v, want id 1 name web", tmpl)
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
		_, err := client.CreateServiceTemplate(context.Background(), "tok", domain.ServiceTemplateWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateServiceTemplate(t *testing.T) {
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
		tmpl, err := client.UpdateServiceTemplate(
			context.Background(), "tok", 7, domain.ServiceTemplateWrite{Name: "Renamed"},
		)
		if err != nil {
			t.Fatalf("UpdateServiceTemplate() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/ipam/service-templates/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/ipam/service-templates/7/", gotMethod, gotPath)
		}
		if tmpl.ID != 7 || tmpl.Name != "Renamed" {
			t.Errorf("service template = %+v, want id 7 name Renamed", tmpl)
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
		_, err := client.UpdateServiceTemplate(context.Background(), "tok", 7, domain.ServiceTemplateWrite{Name: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteServiceTemplate(t *testing.T) {
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
		if err := client.DeleteServiceTemplate(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteServiceTemplate() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/ipam/service-templates/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/ipam/service-templates/3/", gotMethod, gotPath)
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
		err := client.DeleteServiceTemplate(context.Background(), "tok", 3)
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
		if err := client.DeleteServiceTemplate(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteServiceTemplate() 404 returned error: %v", err)
		}
	})
}

func TestServiceTemplateWriteToWire(t *testing.T) {
	t.Parallel()

	protocol := "tcp"
	desc := "https"
	in := domain.ServiceTemplateWrite{Name: "web", Protocol: &protocol, Ports: []int{80, 443}, Description: &desc}
	wire := serviceTemplateWriteToWire(in)
	if wire.Name != "web" || wire.Protocol == nil || *wire.Protocol != "tcp" {
		t.Errorf("wire = %+v, want name web protocol tcp", wire)
	}
	if len(wire.Ports) != 2 || wire.Ports[0] != 80 || wire.Ports[1] != 443 {
		t.Errorf("wire.Ports = %v, want [80 443]", wire.Ports)
	}

	raw, err := json.Marshal(wire)
	if err != nil {
		t.Fatalf("marshal wire: %v", err)
	}
	for _, key := range []string{"name", "protocol", "ports", "description"} {
		if !jsonHasKey(raw, key) {
			t.Errorf("wire JSON missing key %q: %s", key, raw)
		}
	}
}

func TestWireServiceTemplateToDomain(t *testing.T) {
	t.Parallel()

	wire := WireServiceTemplate{
		ID:       1,
		URL:      "http://nb/api/ipam/service-templates/1/",
		Name:     "web",
		Protocol: &WireLabel{Value: "tcp", Label: "TCP"},
		Ports:    []int{80, 443},
		Created:  "2024-01-01",
	}
	got := wireServiceTemplateToDomain(wire)
	if got.ID != 1 || got.Name != "web" {
		t.Errorf("domain = %+v, want id 1 name web", got)
	}
	if got.Protocol == nil || got.Protocol.Value != "tcp" {
		t.Errorf("domain.Protocol = %+v, want tcp", got.Protocol)
	}
	if len(got.Ports) != 2 || got.Ports[1] != 443 {
		t.Errorf("domain.Ports = %v, want [80 443]", got.Ports)
	}
}
