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

// TestClient_ListPowerPorts exercises the power port list endpoint.
func TestClient_ListPowerPorts(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "Token test-token" {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}
			if req.URL.Path != "/api/dcim/power-ports/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{"count":1,"results":[{"id":1,"url":"http://nb/api/dcim/power-ports/1/",
				"name":"PP1","device":{"id":5,"name":"sw1","url":"http://nb/api/dcim/devices/5/"},
				"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}]}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListPowerPorts(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListPowerPorts() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "PP1" {
			t.Errorf("resp = %+v, want count 1 with PP1", resp)
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
		_, err := client.ListPowerPorts(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestClient_CreatePowerPort(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"name":"PP1","url":"http://nb/api/dcim/power-ports/1/",
				"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		dev := 5
		draw := 100
		port, err := client.CreatePowerPort(
			context.Background(), "tok", domain.PowerPortWrite{Name: "PP1", Device: &dev, MaximumDraw: &draw},
		)
		if err != nil {
			t.Fatalf("CreatePowerPort() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/dcim/power-ports/" {
			t.Errorf("method/path = %q %q, want POST /api/dcim/power-ports/", gotMethod, gotPath)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "PP1" {
			t.Errorf("request body = %s, want name PP1", gotBody)
		}
		if port.ID != 1 || port.Name != "PP1" {
			t.Errorf("power port = %+v, want id 1 name PP1", port)
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
		_, err := client.CreatePowerPort(context.Background(), "tok", domain.PowerPortWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdatePowerPort(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"name":"Renamed","url":"http://nb/api/dcim/power-ports/7/",
				"created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		port, err := client.UpdatePowerPort(context.Background(), "tok", 7, domain.PowerPortWrite{Name: "Renamed"})
		if err != nil {
			t.Fatalf("UpdatePowerPort() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/dcim/power-ports/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/dcim/power-ports/7/", gotMethod, gotPath)
		}
		if port.ID != 7 || port.Name != "Renamed" {
			t.Errorf("power port = %+v, want id 7 name Renamed", port)
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
		_, err := client.UpdatePowerPort(context.Background(), "tok", 7, domain.PowerPortWrite{Name: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeletePowerPort(t *testing.T) {
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
		if err := client.DeletePowerPort(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeletePowerPort() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/dcim/power-ports/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/dcim/power-ports/3/", gotMethod, gotPath)
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
		err := client.DeletePowerPort(context.Background(), "tok", 3)
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
		if err := client.DeletePowerPort(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeletePowerPort() 404 returned error: %v", err)
		}
	})
}

func TestPowerPortWriteToWire(t *testing.T) {
	t.Parallel()

	dev := 5
	maxDraw := 100
	alloc := 60
	in := domain.PowerPortWrite{Name: "A", Device: &dev, MaximumDraw: &maxDraw, AllocatedDraw: &alloc}
	wire := powerPortWriteToWire(in)
	if wire.Name != "A" || wire.Device == nil || *wire.Device != 5 {
		t.Errorf("wire = %+v, want name A device 5", wire)
	}
	if wire.MaximumDraw == nil || *wire.MaximumDraw != 100 || wire.AllocatedDraw == nil || *wire.AllocatedDraw != 60 {
		t.Errorf("wire = %+v, want maximum_draw 100 allocated_draw 60", wire)
	}
}

func TestWirePowerPortToDomain(t *testing.T) {
	t.Parallel()

	wire := WirePowerPort{
		ID:      1,
		URL:     "http://nb/api/dcim/power-ports/1/",
		Name:    "PP1",
		Device:  &WireNested{ID: 5, Name: "sw1", Slug: "sw1"},
		Created: "2024-01-01",
	}
	d := wirePowerPortToDomain(wire)
	if d.ID != 1 || d.Name != "PP1" {
		t.Errorf("domain = %+v, want id 1 name PP1", d)
	}
	if d.Device == nil || d.Device.ID != 5 {
		t.Errorf("domain.Device = %+v, want nested id 5", d.Device)
	}

	if got := wirePowerPortToDomain(WirePowerPort{}).Device; got != nil {
		t.Errorf("empty wire Device = %+v, want nil", got)
	}
}
