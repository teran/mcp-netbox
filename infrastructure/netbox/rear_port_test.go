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

// TestClient_ListRearPorts exercises the rear port list endpoint.
func TestClient_ListRearPorts(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "Token test-token" {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}
			if req.URL.Path != "/api/dcim/rear-ports/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{"count":1,"results":[{"id":1,"url":"http://nb/api/dcim/rear-ports/1/",
				"name":"RP1","type":{"value":"lc","label":"LC"},
				"device":{"id":5,"name":"sw1","url":"http://nb/api/dcim/devices/5/"},
				"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}]}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListRearPorts(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListRearPorts() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "RP1" {
			t.Errorf("resp = %+v, want count 1 with RP1", resp)
		}
		if resp.Results[0].Device == nil || resp.Results[0].Device.ID != 5 {
			t.Errorf("resp.Results[0].Device = %+v, want nested id 5", resp.Results[0].Device)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListRearPorts(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestClient_CreateRearPort(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			gotBody, _ = io.ReadAll(req.Body)
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusCreated)
			_, _ = writer.Write([]byte(`{"id":1,"name":"RP1","url":"http://nb/api/dcim/rear-ports/1/",
				"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		dev := 5
		rearPort, err := client.CreateRearPort(context.Background(), "tok", domain.RearPortWrite{Name: "RP1", Device: &dev})
		if err != nil {
			t.Fatalf("CreateRearPort() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/dcim/rear-ports/" {
			t.Errorf("method/path = %q %q, want POST /api/dcim/rear-ports/", gotMethod, gotPath)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "RP1" {
			t.Errorf("request body = %s, want name RP1", gotBody)
		}
		if rearPort.ID != 1 || rearPort.Name != "RP1" {
			t.Errorf("rear port = %+v, want id 1 name RP1", rearPort)
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
		_, err := client.CreateRearPort(context.Background(), "tok", domain.RearPortWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateRearPort(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"name":"Renamed","url":"http://nb/api/dcim/rear-ports/7/",
				"created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		rearPort, err := client.UpdateRearPort(context.Background(), "tok", 7, domain.RearPortWrite{Name: "Renamed"})
		if err != nil {
			t.Fatalf("UpdateRearPort() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/dcim/rear-ports/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/dcim/rear-ports/7/", gotMethod, gotPath)
		}
		if rearPort.ID != 7 || rearPort.Name != "Renamed" {
			t.Errorf("rear port = %+v, want id 7 name Renamed", rearPort)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"positions":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateRearPort(context.Background(), "tok", 7, domain.RearPortWrite{Name: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteRearPort(t *testing.T) {
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
		if err := client.DeleteRearPort(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteRearPort() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/dcim/rear-ports/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/dcim/rear-ports/3/", gotMethod, gotPath)
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
		err := client.DeleteRearPort(context.Background(), "tok", 3)
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
		if err := client.DeleteRearPort(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteRearPort() 404 returned error: %v", err)
		}
	})
}

func TestRearPortWriteToWire(t *testing.T) {
	t.Parallel()

	dev := 5
	pos := 2
	in := domain.RearPortWrite{Name: "A", Device: &dev, Positions: &pos}
	wire := rearPortWriteToWire(in)
	if wire.Name != "A" || wire.Device == nil || *wire.Device != 5 {
		t.Errorf("wire = %+v, want name A device 5", wire)
	}
	if wire.Positions == nil || *wire.Positions != 2 {
		t.Errorf("wire.Positions = %v, want 2", wire.Positions)
	}
}

func TestWireRearPortToDomain(t *testing.T) {
	t.Parallel()

	wire := WireRearPort{
		ID:      1,
		URL:     "http://nb/api/dcim/rear-ports/1/",
		Name:    "RP1",
		Type:    &WireLabel{Value: "lc", Label: "LC"},
		Device:  &WireNested{ID: 5, Name: "sw1", Slug: "sw1"},
		Created: "2024-01-01",
	}
	result := wireRearPortToDomain(wire)
	if result.ID != 1 || result.Name != "RP1" {
		t.Errorf("domain = %+v, want id 1 name RP1", result)
	}
	if result.Type == nil || result.Type.Value != "lc" || result.Type.Label != "LC" {
		t.Errorf("domain.Type = %+v, want label lc/LC", result.Type)
	}
	if result.Device == nil || result.Device.ID != 5 {
		t.Errorf("domain.Device = %+v, want nested id 5", result.Device)
	}

	if got := wireRearPortToDomain(WireRearPort{}).Device; got != nil {
		t.Errorf("empty wire Device = %+v, want nil", got)
	}
}
