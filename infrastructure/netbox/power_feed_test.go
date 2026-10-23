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

// TestClient_ListPowerFeeds exercises the power feed list endpoint.
func TestClient_ListPowerFeeds(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "Token test-token" {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}
			if req.URL.Path != "/api/dcim/power-feeds/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{"count":1,"results":[{"id":1,"url":"http://nb/api/dcim/power-feeds/1/",
				"name":"PF1","power_panel":{"id":5,"name":"PP1","url":"http://nb/api/dcim/power-panels/5/"},
				"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}]}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListPowerFeeds(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListPowerFeeds() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "PF1" {
			t.Errorf("resp = %+v, want count 1 with PF1", resp)
		}
		if resp.Results[0].PowerPanel == nil || resp.Results[0].PowerPanel.ID != 5 {
			t.Errorf("resp.Results[0].PowerPanel = %+v, want nested id 5", resp.Results[0].PowerPanel)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListPowerFeeds(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestClient_CreatePowerFeed(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"name":"PF1","url":"http://nb/api/dcim/power-feeds/1/",
				"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		panel := 5
		feed, err := client.CreatePowerFeed(
			context.Background(), "tok", domain.PowerFeedWrite{Name: "PF1", PowerPanel: &panel},
		)
		if err != nil {
			t.Fatalf("CreatePowerFeed() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/dcim/power-feeds/" {
			t.Errorf("method/path = %q %q, want POST /api/dcim/power-feeds/", gotMethod, gotPath)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "PF1" {
			t.Errorf("request body = %s, want name PF1", gotBody)
		}
		if feed.ID != 1 || feed.Name != "PF1" {
			t.Errorf("power feed = %+v, want id 1 name PF1", feed)
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
		_, err := client.CreatePowerFeed(context.Background(), "tok", domain.PowerFeedWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdatePowerFeed(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"name":"Renamed","url":"http://nb/api/dcim/power-feeds/7/",
				"created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		feed, err := client.UpdatePowerFeed(context.Background(), "tok", 7, domain.PowerFeedWrite{Name: "Renamed"})
		if err != nil {
			t.Fatalf("UpdatePowerFeed() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/dcim/power-feeds/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/dcim/power-feeds/7/", gotMethod, gotPath)
		}
		if feed.ID != 7 || feed.Name != "Renamed" {
			t.Errorf("power feed = %+v, want id 7 name Renamed", feed)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"power_panel":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdatePowerFeed(context.Background(), "tok", 7, domain.PowerFeedWrite{Name: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeletePowerFeed(t *testing.T) {
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
		if err := client.DeletePowerFeed(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeletePowerFeed() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/dcim/power-feeds/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/dcim/power-feeds/3/", gotMethod, gotPath)
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
		err := client.DeletePowerFeed(context.Background(), "tok", 3)
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
		if err := client.DeletePowerFeed(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeletePowerFeed() 404 returned error: %v", err)
		}
	})
}

func TestPowerFeedWriteToWire(t *testing.T) {
	t.Parallel()

	panel := 5
	volt := 230
	in := domain.PowerFeedWrite{Name: "A", PowerPanel: &panel, Voltage: &volt}
	wire := powerFeedWriteToWire(in)
	if wire.Name != "A" || wire.PowerPanel == nil || *wire.PowerPanel != 5 {
		t.Errorf("wire = %+v, want name A power_panel 5", wire)
	}
	if wire.Voltage == nil || *wire.Voltage != 230 {
		t.Errorf("wire.Voltage = %v, want 230", wire.Voltage)
	}
}

func TestWirePowerFeedToDomain(t *testing.T) {
	t.Parallel()

	wire := WirePowerFeed{
		ID:         1,
		URL:        "http://nb/api/dcim/power-feeds/1/",
		Name:       "PF1",
		PowerPanel: &WireNested{ID: 5, Name: "PP1", Slug: "pp1"},
		Created:    "2024-01-01",
	}
	d := wirePowerFeedToDomain(wire)
	if d.ID != 1 || d.Name != "PF1" {
		t.Errorf("domain = %+v, want id 1 name PF1", d)
	}
	if d.PowerPanel == nil || d.PowerPanel.ID != 5 {
		t.Errorf("domain.PowerPanel = %+v, want nested id 5", d.PowerPanel)
	}

	if got := wirePowerFeedToDomain(WirePowerFeed{}).PowerPanel; got != nil {
		t.Errorf("empty wire PowerPanel = %+v, want nil", got)
	}
}
