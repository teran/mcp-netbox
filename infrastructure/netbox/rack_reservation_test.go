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

// TestClient_ListRackReservations exercises the rack reservation list endpoint.
func TestClient_ListRackReservations(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "Token test-token" {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}
			if req.URL.Path != "/api/dcim/rack-reservations/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{"count":1,"results":[{"id":1,"url":"http://nb/api/dcim/rack-reservations/1/",
				"description":"reserved","rack":{"id":3,"name":"R1","url":"http://nb/api/dcim/racks/3/"},
				"units":[1,2],"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}]}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListRackReservations(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListRackReservations() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Description != "reserved" {
			t.Errorf("resp = %+v, want count 1 with reserved", resp)
		}
		if resp.Results[0].Rack == nil || resp.Results[0].Rack.ID != 3 {
			t.Errorf("resp.Results[0].Rack = %+v, want nested id 3", resp.Results[0].Rack)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListRackReservations(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestClient_CreateRackReservation(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"description":"reserved","url":"http://nb/api/dcim/rack-reservations/1/",
				"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		rack := 3
		units := []int{1, 2}
		reservation, err := client.CreateRackReservation(
			context.Background(), "tok", domain.RackReservationWrite{Description: "reserved", Rack: &rack, Units: units},
		)
		if err != nil {
			t.Fatalf("CreateRackReservation() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/dcim/rack-reservations/" {
			t.Errorf("method/path = %q %q, want POST /api/dcim/rack-reservations/", gotMethod, gotPath)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["description"] != "reserved" {
			t.Errorf("request body = %s, want description reserved", gotBody)
		}
		if reservation.ID != 1 || reservation.Description != "reserved" {
			t.Errorf("rack reservation = %+v, want id 1 description reserved", reservation)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"description":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateRackReservation(context.Background(), "tok", domain.RackReservationWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateRackReservation(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"description":"Updated","url":"http://nb/api/dcim/rack-reservations/7/",
				"created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		reservation, err := client.UpdateRackReservation(
			context.Background(), "tok", 7, domain.RackReservationWrite{Description: "Updated"},
		)
		if err != nil {
			t.Fatalf("UpdateRackReservation() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/dcim/rack-reservations/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/dcim/rack-reservations/7/", gotMethod, gotPath)
		}
		if reservation.ID != 7 || reservation.Description != "Updated" {
			t.Errorf("rack reservation = %+v, want id 7 description Updated", reservation)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"rack":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateRackReservation(context.Background(), "tok", 7, domain.RackReservationWrite{Description: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteRackReservation(t *testing.T) {
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
		if err := client.DeleteRackReservation(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteRackReservation() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/dcim/rack-reservations/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/dcim/rack-reservations/3/", gotMethod, gotPath)
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
		err := client.DeleteRackReservation(context.Background(), "tok", 3)
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
		if err := client.DeleteRackReservation(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteRackReservation() 404 returned error: %v", err)
		}
	})
}

func TestRackReservationWriteToWire(t *testing.T) {
	t.Parallel()

	rack := 3
	units := []int{1, 2}
	in := domain.RackReservationWrite{Description: "A", Rack: &rack, Units: units}
	wire := rackReservationWriteToWire(in)
	if wire.Description != "A" || wire.Rack == nil || *wire.Rack != 3 {
		t.Errorf("wire = %+v, want description A rack 3", wire)
	}
	if len(wire.Units) != 2 || wire.Units[0] != 1 || wire.Units[1] != 2 {
		t.Errorf("wire.Units = %v, want [1 2]", wire.Units)
	}
}

func TestWireRackReservationToDomain(t *testing.T) {
	t.Parallel()

	wire := WireRackReservation{
		ID:          1,
		URL:         "http://nb/api/dcim/rack-reservations/1/",
		Description: "reserved",
		Rack:        &WireNested{ID: 3, Name: "R1", Slug: "r1"},
		Created:     "2024-01-01",
	}
	d := wireRackReservationToDomain(wire)
	if d.ID != 1 || d.Description != "reserved" {
		t.Errorf("domain = %+v, want id 1 description reserved", d)
	}
	if d.Rack == nil || d.Rack.ID != 3 {
		t.Errorf("domain.Rack = %+v, want nested id 3", d.Rack)
	}

	if got := wireRackReservationToDomain(WireRackReservation{}).Rack; got != nil {
		t.Errorf("empty wire Rack = %+v, want nil", got)
	}
}
