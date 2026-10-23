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

// TestClient_ListFrontPorts exercises the front port list endpoint.
func TestClient_ListFrontPorts(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "Token test-token" {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}
			if req.URL.Path != "/api/dcim/front-ports/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"url": "http://nb/api/dcim/front-ports/1/",
					"name": "FP1",
					"type": {"value": "lc", "label": "LC"},
					"device": {"id": 5, "name": "sw1", "url": "http://nb/api/dcim/devices/5/"},
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListFrontPorts(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListFrontPorts() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "FP1" {
			t.Errorf("resp = %+v, want count 1 with FP1", resp)
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
		_, err := client.ListFrontPorts(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_CreateFrontPort(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"name":"FP1","url":"http://nb/api/dcim/front-ports/1/",
				"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		dev := 5
		typ := "lc"
		frontPort, err := client.CreateFrontPort(
			context.Background(), "tok", domain.FrontPortWrite{Name: "FP1", Device: &dev, Type: &typ},
		)
		if err != nil {
			t.Fatalf("CreateFrontPort() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/dcim/front-ports/" {
			t.Errorf("method/path = %q %q, want POST /api/dcim/front-ports/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "FP1" {
			t.Errorf("request body = %s, want name FP1", gotBody)
		}
		if frontPort.ID != 1 || frontPort.Name != "FP1" {
			t.Errorf("front port = %+v, want id 1 name FP1", frontPort)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateFrontPort(context.Background(), "tok", domain.FrontPortWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateFrontPort(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"name":"Renamed","url":"http://nb/api/dcim/front-ports/7/",
				"created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		name := "Renamed"
		frontPort, err := client.UpdateFrontPort(
			context.Background(), "tok", 7, domain.FrontPortWrite{Name: name},
		)
		if err != nil {
			t.Fatalf("UpdateFrontPort() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/dcim/front-ports/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/dcim/front-ports/7/", gotMethod, gotPath)
		}
		if frontPort.ID != 7 || frontPort.Name != "Renamed" {
			t.Errorf("front port = %+v, want id 7 name Renamed", frontPort)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"type":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateFrontPort(context.Background(), "tok", 7, domain.FrontPortWrite{Name: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteFrontPort(t *testing.T) {
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
		if err := client.DeleteFrontPort(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteFrontPort() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/dcim/front-ports/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/dcim/front-ports/3/", gotMethod, gotPath)
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
		err := client.DeleteFrontPort(context.Background(), "tok", 3)
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
		if err := client.DeleteFrontPort(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteFrontPort() 404 returned error: %v", err)
		}
	})
}

func TestFrontPortWriteToWire(t *testing.T) {
	t.Parallel()

	dev := 5
	rear := 9
	pos := 2
	typ := "lc"
	in := domain.FrontPortWrite{Name: "A", Device: &dev, RearPort: &rear, RearPortPosition: &pos, Type: &typ}
	wire := frontPortWriteToWire(in)
	if wire.Name != "A" || wire.Device == nil || *wire.Device != 5 {
		t.Errorf("wire = %+v, want name A device 5", wire)
	}
	if wire.RearPort == nil || *wire.RearPort != 9 || wire.RearPortPosition == nil || *wire.RearPortPosition != 2 {
		t.Errorf("wire = %+v, want rear_port 9 rear_port_position 2", wire)
	}
}

func TestWireFrontPortToDomain(t *testing.T) {
	t.Parallel()

	wire := WireFrontPort{
		ID:       1,
		URL:      "http://nb/api/dcim/front-ports/1/",
		Name:     "FP1",
		Type:     &WireLabel{Value: "lc", Label: "LC"},
		Device:   &WireNested{ID: 5, Name: "sw1", Slug: "sw1"},
		RearPort: &WireNested{ID: 9, Name: "RP1", Slug: "rp1"},
		Created:  "2024-01-01",
	}
	result := wireFrontPortToDomain(wire)
	if result.ID != 1 || result.Name != "FP1" {
		t.Errorf("domain = %+v, want id 1 name FP1", result)
	}
	if result.Type == nil || result.Type.Value != "lc" || result.Type.Label != "LC" {
		t.Errorf("domain.Type = %+v, want label lc/LC", result.Type)
	}
	if result.Device == nil || result.Device.ID != 5 {
		t.Errorf("domain.Device = %+v, want nested id 5", result.Device)
	}
	if result.RearPort == nil || result.RearPort.ID != 9 {
		t.Errorf("domain.RearPort = %+v, want nested id 9", result.RearPort)
	}

	if got := wireFrontPortToDomain(WireFrontPort{}).Device; got != nil {
		t.Errorf("empty wire Device = %+v, want nil", got)
	}
}
