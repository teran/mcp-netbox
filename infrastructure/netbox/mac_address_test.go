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

// TestClient_ListMACAddresses exercises the MAC address list endpoint.
func TestClient_ListMACAddresses(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "Token test-token" {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}
			if req.URL.Path != "/api/dcim/mac-addresses/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{"count":1,"results":[{"id":1,"url":"http://nb/api/dcim/mac-addresses/1/",
				"mac_address":"00:11:22:33:44:55","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}]}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListMACAddresses(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListMACAddresses() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Address != "00:11:22:33:44:55" {
			t.Errorf("resp = %+v, want count 1 with MAC", resp)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListMACAddresses(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestClient_CreateMACAddress(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"mac_address":"00:11:22:33:44:55","url":"http://nb/api/dcim/mac-addresses/1/",
				"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		mac, err := client.CreateMACAddress(
			context.Background(), "tok", domain.MACAddressWrite{Address: "00:11:22:33:44:55"},
		)
		if err != nil {
			t.Fatalf("CreateMACAddress() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/dcim/mac-addresses/" {
			t.Errorf("method/path = %q %q, want POST /api/dcim/mac-addresses/", gotMethod, gotPath)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["mac_address"] != "00:11:22:33:44:55" {
			t.Errorf("request body = %s, want mac_address 00:11:22:33:44:55", gotBody)
		}
		if mac.ID != 1 || mac.Address != "00:11:22:33:44:55" {
			t.Errorf("MAC address = %+v, want id 1 with MAC", mac)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"address":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateMACAddress(context.Background(), "tok", domain.MACAddressWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateMACAddress(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"mac_address":"00:11:22:33:44:55","url":"http://nb/api/dcim/mac-addresses/7/",
				"created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		mac, err := client.UpdateMACAddress(
			context.Background(), "tok", 7, domain.MACAddressWrite{Address: "00:11:22:33:44:55"},
		)
		if err != nil {
			t.Fatalf("UpdateMACAddress() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/dcim/mac-addresses/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/dcim/mac-addresses/7/", gotMethod, gotPath)
		}
		if mac.ID != 7 || mac.Address != "00:11:22:33:44:55" {
			t.Errorf("MAC address = %+v, want id 7 with MAC", mac)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"assigned_object_type":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateMACAddress(context.Background(), "tok", 7, domain.MACAddressWrite{Address: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteMACAddress(t *testing.T) {
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
		if err := client.DeleteMACAddress(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteMACAddress() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/dcim/mac-addresses/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/dcim/mac-addresses/3/", gotMethod, gotPath)
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
		err := client.DeleteMACAddress(context.Background(), "tok", 3)
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
		if err := client.DeleteMACAddress(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteMACAddress() 404 returned error: %v", err)
		}
	})
}

func TestMACAddressWriteToWire(t *testing.T) {
	t.Parallel()

	objType := "dcim.interface"
	objID := 9
	in := domain.MACAddressWrite{Address: "00:11:22:33:44:55", AssignedObjectType: &objType, AssignedObjectID: &objID}
	wire := macAddressWriteToWire(in)
	if wire.Address != "00:11:22:33:44:55" {
		t.Errorf("wire.Address = %q, want 00:11:22:33:44:55", wire.Address)
	}
	if wire.AssignedObjectType == nil || *wire.AssignedObjectType != "dcim.interface" {
		t.Errorf("wire.AssignedObjectType = %v, want dcim.interface", wire.AssignedObjectType)
	}
	if wire.AssignedObjectID == nil || *wire.AssignedObjectID != 9 {
		t.Errorf("wire.AssignedObjectID = %v, want 9", wire.AssignedObjectID)
	}
}

func TestWireMACAddressToDomain(t *testing.T) {
	t.Parallel()

	wire := WireMACAddress{
		ID:      1,
		URL:     "http://nb/api/dcim/mac-addresses/1/",
		Address: "00:11:22:33:44:55",
		Created: "2024-01-01",
	}
	d := wireMACAddressToDomain(wire)
	if d.ID != 1 || d.Address != "00:11:22:33:44:55" {
		t.Errorf("domain = %+v, want id 1 with MAC", d)
	}

	if got := wireMACAddressToDomain(WireMACAddress{}).Address; got != "" {
		t.Errorf("empty wire Address = %q, want empty", got)
	}
}
