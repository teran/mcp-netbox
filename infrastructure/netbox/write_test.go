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

func TestClient_CreateSite(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"name":"Site A","slug":"site-a","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		site, err := client.CreateSite(context.Background(), "tok", domain.SiteWrite{Name: "Site A"})
		if err != nil {
			t.Fatalf("CreateSite() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/dcim/sites/" {
			t.Errorf("path = %q, want /api/dcim/sites/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]string
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "Site A" {
			t.Errorf("request body = %s, want name Site A", gotBody)
		}
		if site.ID != 1 || site.Name != "Site A" {
			t.Errorf("site = %+v, want id 1 name Site A", site)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateSite(context.Background(), "tok", domain.SiteWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if ve.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode = %d, want 400", ve.StatusCode)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateSite(context.Background(), "tok", domain.SiteWrite{Name: "X"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			t.Errorf("err = %+v, want NOT a ValidationError for 404", ve)
		}
	})

	t.Run("unauthorized 401", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateSite(context.Background(), "tok", domain.SiteWrite{Name: "X"})
		if err == nil || !contains(err.Error(), "unauthorized") {
			t.Errorf("err = %v, want unauthorized", err)
		}
	})
}

func TestClient_UpdateSite(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"name":"Renamed","slug":"renamed","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		desc := "updated"
		site, err := client.UpdateSite(context.Background(), "tok", 7, domain.SiteWrite{Name: "Renamed", Description: &desc})
		if err != nil {
			t.Fatalf("UpdateSite() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/dcim/sites/7/" {
			t.Errorf("path = %q, want /api/dcim/sites/7/", gotPath)
		}
		if site.ID != 7 || site.Name != "Renamed" {
			t.Errorf("site = %+v, want id 7 name Renamed", site)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["description"] != "updated" {
			t.Errorf("request body = %s, want description present", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"slug":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateSite(context.Background(), "tok", 7, domain.SiteWrite{Name: "X"})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteSite(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteSite(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteSite() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/dcim/sites/3/" {
			t.Errorf("path = %q, want /api/dcim/sites/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteSite(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteSite(context.Background(), "tok", 999); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestMapWriteStatus(t *testing.T) {
	t.Parallel()

	t.Run("created 201 returns body", func(t *testing.T) {
		raw, err := mapWriteStatus(http.StatusCreated, []byte(`{"id":1}`))
		if err != nil {
			t.Fatalf("mapWriteStatus() error: %v", err)
		}
		if string(raw) != `{"id":1}` {
			t.Errorf("raw = %s, want body", raw)
		}
	})

	t.Run("ok 200 returns body", func(t *testing.T) {
		raw, err := mapWriteStatus(http.StatusOK, []byte(`{"id":2}`))
		if err != nil || string(raw) != `{"id":2}` {
			t.Errorf("raw = %s, err = %v", raw, err)
		}
	})

	t.Run("no content 204 returns nil", func(t *testing.T) {
		raw, err := mapWriteStatus(http.StatusNoContent, nil)
		if err != nil {
			t.Fatalf("mapWriteStatus() error: %v", err)
		}
		if raw != nil {
			t.Errorf("raw = %v, want nil", raw)
		}
	})

	t.Run("bad request 400 returns validation error", func(t *testing.T) {
		_, err := mapWriteStatus(http.StatusBadRequest, []byte(`{"e":["x"]}`))
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != `{"e":["x"]}` {
			t.Errorf("Body = %s, want body", ve.Body)
		}
	})

	t.Run("unexpected 500 returns error", func(t *testing.T) {
		_, err := mapWriteStatus(http.StatusInternalServerError, []byte(`{}`))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestClient_CreateDevice(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"name":"router1","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		dt := 2
		device, err := client.CreateDevice(context.Background(), "tok", domain.DeviceWrite{Name: "router1", DeviceType: &dt})
		if err != nil {
			t.Fatalf("CreateDevice() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/dcim/devices/" {
			t.Errorf("path = %q, want /api/dcim/devices/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "router1" {
			t.Errorf("request body = %s, want name router1", gotBody)
		}
		if reqBody["device_type"] != float64(2) {
			t.Errorf("request body = %s, want device_type 2", gotBody)
		}
		if device.ID != 1 || device.Name != "router1" {
			t.Errorf("device = %+v, want id 1 name router1", device)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateDevice(context.Background(), "tok", domain.DeviceWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if ve.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode = %d, want 400", ve.StatusCode)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateDevice(context.Background(), "tok", domain.DeviceWrite{Name: "X"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			t.Errorf("err = %+v, want NOT a ValidationError for 404", ve)
		}
	})
}

func TestClient_UpdateDevice(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"name":"router1","serial":"SN-X","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		serial := "SN-X"
		device, err := client.UpdateDevice(context.Background(), "tok", 7, domain.DeviceWrite{Serial: &serial})
		if err != nil {
			t.Fatalf("UpdateDevice() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/dcim/devices/7/" {
			t.Errorf("path = %q, want /api/dcim/devices/7/", gotPath)
		}
		if device.ID != 7 || device.Serial != "SN-X" {
			t.Errorf("device = %+v, want id 7 serial SN-X", device)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["serial"] != "SN-X" {
			t.Errorf("request body = %s, want serial present", gotBody)
		}
		if _, ok := reqBody["name"]; ok {
			t.Errorf("request body = %s, want name omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"role":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateDevice(context.Background(), "tok", 7, domain.DeviceWrite{Name: "X"})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteDevice(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteDevice(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteDevice() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/dcim/devices/3/" {
			t.Errorf("path = %q, want /api/dcim/devices/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteDevice(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteDevice(context.Background(), "tok", 999); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestSiteWriteToWire(t *testing.T) {
	t.Parallel()

	in := domain.SiteWrite{Name: "A", Slug: ptrStr("a")}
	w := siteWriteToWire(in)
	if w.Name != "A" || w.Slug == nil || *w.Slug != "a" {
		t.Errorf("wire = %+v, want name A slug a", w)
	}
}

func TestDeviceWriteToWire(t *testing.T) {
	t.Parallel()

	dt := 2
	serial := "SN"
	in := domain.DeviceWrite{Name: "A", DeviceType: &dt, Serial: &serial}
	w := deviceWriteToWire(in)
	if w.Name != "A" || w.DeviceType == nil || *w.DeviceType != 2 || w.Serial == nil || *w.Serial != "SN" {
		t.Errorf("wire = %+v, want name A device_type 2 serial SN", w)
	}
}

func TestClient_CreateIPAddress(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"address":"192.168.1.1/24","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		vrf := 2
		ip, err := client.CreateIPAddress(context.Background(), "tok", domain.IPAddressWrite{Address: "192.168.1.1/24", VRF: &vrf})
		if err != nil {
			t.Fatalf("CreateIPAddress() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/ipam/ip-addresses/" {
			t.Errorf("path = %q, want /api/ipam/ip-addresses/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["address"] != "192.168.1.1/24" {
			t.Errorf("request body = %s, want address 192.168.1.1/24", gotBody)
		}
		if reqBody["vrf"] != float64(2) {
			t.Errorf("request body = %s, want vrf 2", gotBody)
		}
		if ip.ID != 1 || ip.Address != "192.168.1.1/24" {
			t.Errorf("ip = %+v, want id 1 address 192.168.1.1/24", ip)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"address":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateIPAddress(context.Background(), "tok", domain.IPAddressWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if ve.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode = %d, want 400", ve.StatusCode)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateIPAddress(context.Background(), "tok", domain.IPAddressWrite{Address: "X"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			t.Errorf("err = %+v, want NOT a ValidationError for 404", ve)
		}
	})
}

func TestClient_UpdateIPAddress(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"address":"10.0.0.1/32","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		desc := "updated"
		ip, err := client.UpdateIPAddress(context.Background(), "tok", 7, domain.IPAddressWrite{Description: &desc})
		if err != nil {
			t.Fatalf("UpdateIPAddress() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/ipam/ip-addresses/7/" {
			t.Errorf("path = %q, want /api/ipam/ip-addresses/7/", gotPath)
		}
		if ip.ID != 7 {
			t.Errorf("ip = %+v, want id 7", ip)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["description"] != "updated" {
			t.Errorf("request body = %s, want description present", gotBody)
		}
		if _, ok := reqBody["address"]; ok {
			t.Errorf("request body = %s, want address omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"dns_name":["Enter a valid hostname."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateIPAddress(context.Background(), "tok", 7, domain.IPAddressWrite{Address: "X"})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteIPAddress(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteIPAddress(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteIPAddress() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/ipam/ip-addresses/3/" {
			t.Errorf("path = %q, want /api/ipam/ip-addresses/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteIPAddress(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteIPAddress(context.Background(), "tok", 999); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestIPAddressWriteToWire(t *testing.T) {
	t.Parallel()

	vrf := 2
	dns := "host.example.com"
	in := domain.IPAddressWrite{Address: "10.0.0.1/32", VRF: &vrf, DNSName: &dns}
	w := ipAddressWriteToWire(in)
	if w.Address != "10.0.0.1/32" || w.VRF == nil || *w.VRF != 2 || w.DNSName == nil || *w.DNSName != dns {
		t.Errorf("wire = %+v, want address 10.0.0.1/32 vrf 2 dns %s", w, dns)
	}
}

func TestClient_CreatePrefix(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"prefix":"10.0.0.0/24","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		site := 2
		prefix, err := client.CreatePrefix(context.Background(), "tok", domain.PrefixWrite{Prefix: "10.0.0.0/24", Site: &site})
		if err != nil {
			t.Fatalf("CreatePrefix() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/ipam/prefixes/" {
			t.Errorf("path = %q, want /api/ipam/prefixes/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["prefix"] != "10.0.0.0/24" {
			t.Errorf("request body = %s, want prefix 10.0.0.0/24", gotBody)
		}
		if reqBody["site"] != float64(2) {
			t.Errorf("request body = %s, want site 2", gotBody)
		}
		if prefix.ID != 1 || prefix.Prefix != "10.0.0.0/24" {
			t.Errorf("prefix = %+v, want id 1 prefix 10.0.0.0/24", prefix)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"prefix":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreatePrefix(context.Background(), "tok", domain.PrefixWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if ve.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode = %d, want 400", ve.StatusCode)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_UpdatePrefix(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"prefix":"10.0.0.0/24","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		desc := "updated"
		prefix, err := client.UpdatePrefix(context.Background(), "tok", 7, domain.PrefixWrite{Description: &desc})
		if err != nil {
			t.Fatalf("UpdatePrefix() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/ipam/prefixes/7/" {
			t.Errorf("path = %q, want /api/ipam/prefixes/7/", gotPath)
		}
		if prefix.ID != 7 {
			t.Errorf("prefix = %+v, want id 7", prefix)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["description"] != "updated" {
			t.Errorf("request body = %s, want description present", gotBody)
		}
		if _, ok := reqBody["prefix"]; ok {
			t.Errorf("request body = %s, want prefix omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"status":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdatePrefix(context.Background(), "tok", 7, domain.PrefixWrite{Prefix: "X"})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeletePrefix(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeletePrefix(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeletePrefix() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/ipam/prefixes/3/" {
			t.Errorf("path = %q, want /api/ipam/prefixes/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeletePrefix(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestPrefixWriteToWire(t *testing.T) {
	t.Parallel()

	site := 2
	pool := true
	in := domain.PrefixWrite{Prefix: "10.0.0.0/24", Site: &site, IsPool: &pool}
	w := prefixWriteToWire(in)
	if w.Prefix != "10.0.0.0/24" || w.Site == nil || *w.Site != 2 || w.IsPool == nil || !*w.IsPool {
		t.Errorf("wire = %+v, want prefix 10.0.0.0/24 site 2 is_pool true", w)
	}
}

func TestClient_CreateVLAN(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"vid":100,"name":"mgmt","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		site := 2
		vlan, err := client.CreateVLAN(context.Background(), "tok", domain.VLANWrite{VID: 100, Name: "mgmt", Site: &site})
		if err != nil {
			t.Fatalf("CreateVLAN() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/ipam/vlans/" {
			t.Errorf("path = %q, want /api/ipam/vlans/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "mgmt" {
			t.Errorf("request body = %s, want name mgmt", gotBody)
		}
		if reqBody["vid"] != float64(100) {
			t.Errorf("request body = %s, want vid 100", gotBody)
		}
		if reqBody["site"] != float64(2) {
			t.Errorf("request body = %s, want site 2", gotBody)
		}
		if vlan.ID != 1 || vlan.VID != 100 || vlan.Name != "mgmt" {
			t.Errorf("vlan = %+v, want id 1 vid 100 name mgmt", vlan)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateVLAN(context.Background(), "tok", domain.VLANWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if ve.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode = %d, want 400", ve.StatusCode)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_UpdateVLAN(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"vid":200,"name":"Renamed","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		desc := "updated"
		vlan, err := client.UpdateVLAN(context.Background(), "tok", 7, domain.VLANWrite{Description: &desc})
		if err != nil {
			t.Fatalf("UpdateVLAN() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/ipam/vlans/7/" {
			t.Errorf("path = %q, want /api/ipam/vlans/7/", gotPath)
		}
		if vlan.ID != 7 {
			t.Errorf("vlan = %+v, want id 7", vlan)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["description"] != "updated" {
			t.Errorf("request body = %s, want description present", gotBody)
		}
		if _, ok := reqBody["name"]; ok {
			t.Errorf("request body = %s, want name omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"status":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateVLAN(context.Background(), "tok", 7, domain.VLANWrite{Name: "X"})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteVLAN(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteVLAN(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteVLAN() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/ipam/vlans/3/" {
			t.Errorf("path = %q, want /api/ipam/vlans/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteVLAN(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestVLANWriteToWire(t *testing.T) {
	t.Parallel()

	site := 2
	in := domain.VLANWrite{VID: 100, Name: "mgmt", Site: &site}
	w := vlanWriteToWire(in)
	if w.VID != 100 || w.Name != "mgmt" || w.Site == nil || *w.Site != 2 {
		t.Errorf("wire = %+v, want vid 100 name mgmt site 2", w)
	}
}

func TestClient_CreateVirtualMachine(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"name":"web-01","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		cluster := 2
		vm, err := client.CreateVirtualMachine(context.Background(), "tok", domain.VirtualMachineWrite{Name: "web-01", Cluster: &cluster})
		if err != nil {
			t.Fatalf("CreateVirtualMachine() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/virtualization/virtual-machines/" {
			t.Errorf("path = %q, want /api/virtualization/virtual-machines/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "web-01" {
			t.Errorf("request body = %s, want name web-01", gotBody)
		}
		if reqBody["cluster"] != float64(2) {
			t.Errorf("request body = %s, want cluster 2", gotBody)
		}
		if vm.ID != 1 || vm.Name != "web-01" {
			t.Errorf("vm = %+v, want id 1 name web-01", vm)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateVirtualMachine(context.Background(), "tok", domain.VirtualMachineWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if ve.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode = %d, want 400", ve.StatusCode)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_UpdateVirtualMachine(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"name":"Renamed","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		mem := 8192
		vm, err := client.UpdateVirtualMachine(context.Background(), "tok", 7, domain.VirtualMachineWrite{Memory: &mem})
		if err != nil {
			t.Fatalf("UpdateVirtualMachine() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/virtualization/virtual-machines/7/" {
			t.Errorf("path = %q, want /api/virtualization/virtual-machines/7/", gotPath)
		}
		if vm.ID != 7 {
			t.Errorf("vm = %+v, want id 7", vm)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["memory"] != float64(8192) {
			t.Errorf("request body = %s, want memory present", gotBody)
		}
		if _, ok := reqBody["name"]; ok {
			t.Errorf("request body = %s, want name omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"status":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateVirtualMachine(context.Background(), "tok", 7, domain.VirtualMachineWrite{Name: "X"})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteVirtualMachine(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteVirtualMachine(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteVirtualMachine() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/virtualization/virtual-machines/3/" {
			t.Errorf("path = %q, want /api/virtualization/virtual-machines/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteVirtualMachine(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestVirtualMachineWriteToWire(t *testing.T) {
	t.Parallel()

	cluster := 2
	in := domain.VirtualMachineWrite{Name: "web-01", Cluster: &cluster}
	w := virtualMachineWriteToWire(in)
	if w.Name != "web-01" || w.Cluster == nil || *w.Cluster != 2 {
		t.Errorf("wire = %+v, want name web-01 cluster 2", w)
	}
}

func TestClient_CreateCluster(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"name":"prod","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		ct := 5
		cluster, err := client.CreateCluster(context.Background(), "tok", domain.ClusterWrite{Name: "prod", ClusterType: &ct})
		if err != nil {
			t.Fatalf("CreateCluster() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/virtualization/clusters/" {
			t.Errorf("path = %q, want /api/virtualization/clusters/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "prod" {
			t.Errorf("request body = %s, want name prod", gotBody)
		}
		if reqBody["type"] != float64(5) {
			t.Errorf("request body = %s, want type 5", gotBody)
		}
		if cluster.ID != 1 || cluster.Name != "prod" {
			t.Errorf("cluster = %+v, want id 1 name prod", cluster)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateCluster(context.Background(), "tok", domain.ClusterWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if ve.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode = %d, want 400", ve.StatusCode)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_UpdateCluster(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"name":"Renamed","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		desc := "updated"
		cluster, err := client.UpdateCluster(context.Background(), "tok", 7, domain.ClusterWrite{Description: &desc})
		if err != nil {
			t.Fatalf("UpdateCluster() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/virtualization/clusters/7/" {
			t.Errorf("path = %q, want /api/virtualization/clusters/7/", gotPath)
		}
		if cluster.ID != 7 {
			t.Errorf("cluster = %+v, want id 7", cluster)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["description"] != "updated" {
			t.Errorf("request body = %s, want description present", gotBody)
		}
		if _, ok := reqBody["name"]; ok {
			t.Errorf("request body = %s, want name omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"type":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateCluster(context.Background(), "tok", 7, domain.ClusterWrite{Name: "X"})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteCluster(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteCluster(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteCluster() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/virtualization/clusters/3/" {
			t.Errorf("path = %q, want /api/virtualization/clusters/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteCluster(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClusterWriteToWire(t *testing.T) {
	t.Parallel()

	ct := 5
	in := domain.ClusterWrite{Name: "prod", ClusterType: &ct}
	w := clusterWriteToWire(in)
	if w.Name != "prod" || w.ClusterType == nil || *w.ClusterType != 5 {
		t.Errorf("wire = %+v, want name prod type 5", w)
	}
}

func TestClient_CreateCircuit(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"cid":"CIR-001","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		prov := 2
		ct := 3
		circuit, err := client.CreateCircuit(context.Background(), "tok", domain.CircuitWrite{CID: "CIR-001", Provider: &prov, CircuitType: &ct})
		if err != nil {
			t.Fatalf("CreateCircuit() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/circuits/circuits/" {
			t.Errorf("path = %q, want /api/circuits/circuits/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["cid"] != "CIR-001" {
			t.Errorf("request body = %s, want cid CIR-001", gotBody)
		}
		if reqBody["provider"] != float64(2) {
			t.Errorf("request body = %s, want provider 2", gotBody)
		}
		if reqBody["type"] != float64(3) {
			t.Errorf("request body = %s, want circuit_type 3", gotBody)
		}
		if circuit.ID != 1 || circuit.CID != "CIR-001" {
			t.Errorf("circuit = %+v, want id 1 cid CIR-001", circuit)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"cid":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateCircuit(context.Background(), "tok", domain.CircuitWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if ve.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode = %d, want 400", ve.StatusCode)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_UpdateCircuit(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"cid":"Renamed","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		desc := "updated"
		circuit, err := client.UpdateCircuit(context.Background(), "tok", 7, domain.CircuitWrite{Description: &desc})
		if err != nil {
			t.Fatalf("UpdateCircuit() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/circuits/circuits/7/" {
			t.Errorf("path = %q, want /api/circuits/circuits/7/", gotPath)
		}
		if circuit.ID != 7 {
			t.Errorf("circuit = %+v, want id 7", circuit)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["description"] != "updated" {
			t.Errorf("request body = %s, want description present", gotBody)
		}
		if _, ok := reqBody["cid"]; ok {
			t.Errorf("request body = %s, want cid omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"provider":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateCircuit(context.Background(), "tok", 7, domain.CircuitWrite{CID: "X"})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteCircuit(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteCircuit(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteCircuit() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/circuits/circuits/3/" {
			t.Errorf("path = %q, want /api/circuits/circuits/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteCircuit(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestCircuitWriteToWire(t *testing.T) {
	t.Parallel()

	prov := 2
	ct := 3
	in := domain.CircuitWrite{CID: "CIR-001", Provider: &prov, CircuitType: &ct}
	w := circuitWriteToWire(in)
	if w.CID != "CIR-001" || w.Provider == nil || *w.Provider != 2 || w.CircuitType == nil || *w.CircuitType != 3 {
		t.Errorf("wire = %+v, want cid CIR-001 provider 2 circuit_type 3", w)
	}
}

// unmarshalableCustomFields returns a CustomFields value that json.Marshal
// cannot serialize, forcing the write methods' marshal-error branch.
func unmarshalableCustomFields() map[string]any {
	return map[string]any{"bad": make(chan int)}
}

func TestClient_WriteMarshalErrors(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	ctx := context.Background()

	badIP := domain.IPAddressWrite{Address: "10.0.0.1/32", CustomFields: unmarshalableCustomFields()}
	if _, err := client.CreateIPAddress(ctx, "tok", badIP); err == nil {
		t.Error("CreateIPAddress: expected marshal error, got nil")
	}
	if _, err := client.UpdateIPAddress(ctx, "tok", 7, badIP); err == nil {
		t.Error("UpdateIPAddress: expected marshal error, got nil")
	}

	badPrefix := domain.PrefixWrite{Prefix: "10.0.0.0/24", CustomFields: unmarshalableCustomFields()}
	if _, err := client.CreatePrefix(ctx, "tok", badPrefix); err == nil {
		t.Error("CreatePrefix: expected marshal error, got nil")
	}
	if _, err := client.UpdatePrefix(ctx, "tok", 7, badPrefix); err == nil {
		t.Error("UpdatePrefix: expected marshal error, got nil")
	}

	badVLAN := domain.VLANWrite{VID: 100, Name: "mgmt", CustomFields: unmarshalableCustomFields()}
	if _, err := client.CreateVLAN(ctx, "tok", badVLAN); err == nil {
		t.Error("CreateVLAN: expected marshal error, got nil")
	}
	if _, err := client.UpdateVLAN(ctx, "tok", 7, badVLAN); err == nil {
		t.Error("UpdateVLAN: expected marshal error, got nil")
	}

	badVM := domain.VirtualMachineWrite{Name: "web-01", CustomFields: unmarshalableCustomFields()}
	if _, err := client.CreateVirtualMachine(ctx, "tok", badVM); err == nil {
		t.Error("CreateVirtualMachine: expected marshal error, got nil")
	}
	if _, err := client.UpdateVirtualMachine(ctx, "tok", 7, badVM); err == nil {
		t.Error("UpdateVirtualMachine: expected marshal error, got nil")
	}

	badCluster := domain.ClusterWrite{Name: "prod", CustomFields: unmarshalableCustomFields()}
	if _, err := client.CreateCluster(ctx, "tok", badCluster); err == nil {
		t.Error("CreateCluster: expected marshal error, got nil")
	}
	if _, err := client.UpdateCluster(ctx, "tok", 7, badCluster); err == nil {
		t.Error("UpdateCluster: expected marshal error, got nil")
	}

	badCircuit := domain.CircuitWrite{CID: "CIR-001", CustomFields: unmarshalableCustomFields()}
	if _, err := client.CreateCircuit(ctx, "tok", badCircuit); err == nil {
		t.Error("CreateCircuit: expected marshal error, got nil")
	}
	if _, err := client.UpdateCircuit(ctx, "tok", 7, badCircuit); err == nil {
		t.Error("UpdateCircuit: expected marshal error, got nil")
	}

	badRack := domain.RackWrite{Name: "R1", CustomFields: unmarshalableCustomFields()}
	if _, err := client.CreateRack(ctx, "tok", badRack); err == nil {
		t.Error("CreateRack: expected marshal error, got nil")
	}
	if _, err := client.UpdateRack(ctx, "tok", 7, badRack); err == nil {
		t.Error("UpdateRack: expected marshal error, got nil")
	}

	badInterface := domain.InterfaceWrite{Name: "eth0", CustomFields: unmarshalableCustomFields()}
	if _, err := client.CreateInterface(ctx, "tok", badInterface); err == nil {
		t.Error("CreateInterface: expected marshal error, got nil")
	}
	if _, err := client.UpdateInterface(ctx, "tok", 7, badInterface); err == nil {
		t.Error("UpdateInterface: expected marshal error, got nil")
	}

	badCT := domain.CircuitTerminationWrite{TermSide: "A", CustomFields: unmarshalableCustomFields()}
	if _, err := client.CreateCircuitTermination(ctx, "tok", badCT); err == nil {
		t.Error("CreateCircuitTermination: expected marshal error, got nil")
	}
	if _, err := client.UpdateCircuitTermination(ctx, "tok", 7, badCT); err == nil {
		t.Error("UpdateCircuitTermination: expected marshal error, got nil")
	}
}

//nolint:gocognit // exhaustive per-path write test
func TestClient_CreateCircuitTermination(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"term_side":"A","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		circ := 2
		site := 3
		ct, err := client.CreateCircuitTermination(context.Background(), "tok", domain.CircuitTerminationWrite{TermSide: "A", Circuit: &circ, Site: &site})
		if err != nil {
			t.Fatalf("CreateCircuitTermination() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/circuits/circuit-terminations/" {
			t.Errorf("path = %q, want /api/circuits/circuit-terminations/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["term_side"] != "A" {
			t.Errorf("request body = %s, want term_side A", gotBody)
		}
		if reqBody["circuit"] != float64(2) {
			t.Errorf("request body = %s, want circuit 2", gotBody)
		}
		if reqBody["site"] != float64(3) {
			t.Errorf("request body = %s, want site 3", gotBody)
		}
		if ct.ID != 1 || ct.TermSide != "A" {
			t.Errorf("ct = %+v, want id 1 term_side A", ct)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"term_side":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateCircuitTermination(context.Background(), "tok", domain.CircuitTerminationWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if ve.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode = %d, want 400", ve.StatusCode)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateCircuitTermination(context.Background(), "tok", domain.CircuitTerminationWrite{TermSide: "A"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			t.Errorf("err = %+v, want NOT a ValidationError for 404", ve)
		}
	})
}

func TestClient_UpdateCircuitTermination(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"term_side":"A","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		speed := 1000000
		ct, err := client.UpdateCircuitTermination(context.Background(), "tok", 7, domain.CircuitTerminationWrite{Speed: &speed})
		if err != nil {
			t.Fatalf("UpdateCircuitTermination() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/circuits/circuit-terminations/7/" {
			t.Errorf("path = %q, want /api/circuits/circuit-terminations/7/", gotPath)
		}
		if ct.ID != 7 {
			t.Errorf("ct = %+v, want id 7", ct)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["speed"] != float64(1000000) {
			t.Errorf("request body = %s, want speed present", gotBody)
		}
		if _, ok := reqBody["term_side"]; ok {
			t.Errorf("request body = %s, want term_side omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"site":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateCircuitTermination(context.Background(), "tok", 7, domain.CircuitTerminationWrite{TermSide: "A"})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteCircuitTermination(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteCircuitTermination(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteCircuitTermination() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/circuits/circuit-terminations/3/" {
			t.Errorf("path = %q, want /api/circuits/circuit-terminations/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteCircuitTermination(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteCircuitTermination(context.Background(), "tok", 999); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestCircuitTerminationWriteToWire(t *testing.T) {
	t.Parallel()

	circ := 2
	site := 3
	speed := 1000000
	in := domain.CircuitTerminationWrite{TermSide: "A", Circuit: &circ, Site: &site, Speed: &speed}
	w := circuitTerminationWriteToWire(in)
	if w.TermSide != "A" || w.Circuit == nil || *w.Circuit != 2 || w.Site == nil || *w.Site != 3 || w.Speed == nil || *w.Speed != 1000000 {
		t.Errorf("wire = %+v, want term_side A circuit 2 site 3 speed 1000000", w)
	}
}

func TestCableWriteToWire(t *testing.T) {
	t.Parallel()

	ta := &domain.CableTerminationWrite{ObjectType: "dcim.interface", ObjectID: 101}
	tb := &domain.CableTerminationWrite{ObjectType: "dcim.interface", ObjectID: 102}
	label := "link-01"
	length := 10.5
	in := domain.CableWrite{TerminationA: ta, TerminationB: tb, Label: &label, Length: &length}
	w := cableWriteToWire(in)
	if w.TerminationA == nil || w.TerminationA.ObjectType != "dcim.interface" || w.TerminationA.ObjectID != 101 {
		t.Errorf("TerminationA = %+v, want dcim.interface 101", w.TerminationA)
	}
	if w.TerminationB == nil || w.TerminationB.ObjectID != 102 {
		t.Errorf("TerminationB = %+v, want object_id 102", w.TerminationB)
	}
	if w.Label == nil || *w.Label != "link-01" || w.Length == nil || *w.Length != 10.5 {
		t.Errorf("wire = %+v, want label link-01 length 10.5", w)
	}
}

func TestClient_CreateCable(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"label":"link-01","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		ta := &domain.CableTerminationWrite{ObjectType: "dcim.interface", ObjectID: 101}
		tb := &domain.CableTerminationWrite{ObjectType: "dcim.interface", ObjectID: 102}
		label := "link-01"
		c, err := client.CreateCable(context.Background(), "tok", domain.CableWrite{TerminationA: ta, TerminationB: tb, Label: &label})
		if err != nil {
			t.Fatalf("CreateCable() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/dcim/cables/" {
			t.Errorf("path = %q, want /api/dcim/cables/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		taMap, ok := reqBody["termination_a"].(map[string]any)
		if !ok || taMap["object_id"] != float64(101) {
			t.Errorf("request body = %s, want termination_a object_id 101", gotBody)
		}
		if c.ID != 1 || c.Label != "link-01" {
			t.Errorf("c = %+v, want id 1 label link-01", c)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"termination_a":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateCable(context.Background(), "tok", domain.CableWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if ve.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode = %d, want 400", ve.StatusCode)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateCable(context.Background(), "tok", domain.CableWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			t.Errorf("err = %+v, want NOT a ValidationError for 404", ve)
		}
	})
}

func TestClient_UpdateCable(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"label":"link-02","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		label := "link-02"
		c, err := client.UpdateCable(context.Background(), "tok", 7, domain.CableWrite{Label: &label})
		if err != nil {
			t.Fatalf("UpdateCable() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/dcim/cables/7/" {
			t.Errorf("path = %q, want /api/dcim/cables/7/", gotPath)
		}
		if c.ID != 7 {
			t.Errorf("c = %+v, want id 7", c)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["label"] != "link-02" {
			t.Errorf("request body = %s, want label present", gotBody)
		}
		if _, ok := reqBody["termination_a"]; ok {
			t.Errorf("request body = %s, want termination_a omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"status":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateCable(context.Background(), "tok", 7, domain.CableWrite{})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteCable(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteCable(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteCable() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/dcim/cables/3/" {
			t.Errorf("path = %q, want /api/dcim/cables/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteCable(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestVMInterfaceWriteToWire(t *testing.T) {
	t.Parallel()

	vm := 5
	mtu := 1500
	in := domain.VMInterfaceWrite{Name: "eth0", VirtualMachine: &vm, MTU: &mtu}
	w := vmInterfaceWriteToWire(in)
	if w.Name != "eth0" || w.VirtualMachine == nil || *w.VirtualMachine != 5 || w.MTU == nil || *w.MTU != 1500 {
		t.Errorf("wire = %+v, want name eth0 virtual_machine 5 mtu 1500", w)
	}
}

func TestClient_CreateVMInterface(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"name":"eth0","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		vm := 5
		vi, err := client.CreateVMInterface(context.Background(), "tok", domain.VMInterfaceWrite{Name: "eth0", VirtualMachine: &vm})
		if err != nil {
			t.Fatalf("CreateVMInterface() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/virtualization/interfaces/" {
			t.Errorf("path = %q, want /api/virtualization/interfaces/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "eth0" || reqBody["virtual_machine"] != float64(5) {
			t.Errorf("request body = %s, want name eth0 virtual_machine 5", gotBody)
		}
		if vi.ID != 1 || vi.Name != "eth0" {
			t.Errorf("vi = %+v, want id 1 name eth0", vi)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateVMInterface(context.Background(), "tok", domain.VMInterfaceWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if ve.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode = %d, want 400", ve.StatusCode)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateVMInterface(context.Background(), "tok", domain.VMInterfaceWrite{Name: "eth0"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			t.Errorf("err = %+v, want NOT a ValidationError for 404", ve)
		}
	})
}

func TestClient_UpdateVMInterface(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"name":"eth1","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		mtu := 9000
		vi, err := client.UpdateVMInterface(context.Background(), "tok", 7, domain.VMInterfaceWrite{MTU: &mtu})
		if err != nil {
			t.Fatalf("UpdateVMInterface() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/virtualization/interfaces/7/" {
			t.Errorf("path = %q, want /api/virtualization/interfaces/7/", gotPath)
		}
		if vi.ID != 7 {
			t.Errorf("vi = %+v, want id 7", vi)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["mtu"] != float64(9000) {
			t.Errorf("request body = %s, want mtu present", gotBody)
		}
		if _, ok := reqBody["name"]; ok {
			t.Errorf("request body = %s, want name omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"mac_address":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateVMInterface(context.Background(), "tok", 7, domain.VMInterfaceWrite{})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteVMInterface(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteVMInterface(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteVMInterface() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/virtualization/interfaces/3/" {
			t.Errorf("path = %q, want /api/virtualization/interfaces/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteVMInterface(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestProviderWriteToWire(t *testing.T) {
	t.Parallel()

	asn := 64512
	in := domain.ProviderWrite{Name: "ACME", Asn: &asn}
	w := providerWriteToWire(in)
	if w.Name != "ACME" || w.Asn == nil || *w.Asn != 64512 {
		t.Errorf("wire = %+v, want name ACME asn 64512", w)
	}
}

func TestUnmarshalWriteEntities_InvalidJSON(t *testing.T) {
	t.Parallel()

	bad := domain.RawObject(`{invalid`)
	if _, err := unmarshalCable(bad); err == nil {
		t.Error("unmarshalCable() = nil, want error for invalid json")
	}
	if _, err := unmarshalVMInterface(bad); err == nil {
		t.Error("unmarshalVMInterface() = nil, want error for invalid json")
	}
	if _, err := unmarshalProvider(bad); err == nil {
		t.Error("unmarshalProvider() = nil, want error for invalid json")
	}
	if _, err := unmarshalTenant(bad); err == nil {
		t.Error("unmarshalTenant() = nil, want error for invalid json")
	}
	if _, err := unmarshalManufacturer(bad); err == nil {
		t.Error("unmarshalManufacturer() = nil, want error for invalid json")
	}
	if _, err := unmarshalDeviceType(bad); err == nil {
		t.Error("unmarshalDeviceType() = nil, want error for invalid json")
	}
}

func TestDeviceTypeWriteToWire(t *testing.T) {
	t.Parallel()

	mf := 3
	in := domain.DeviceTypeWrite{Manufacturer: &mf, Model: "C9300"}
	w := deviceTypeWriteToWire(in)
	if w.Model != "C9300" || w.Manufacturer == nil || *w.Manufacturer != 3 {
		t.Errorf("wire = %+v, want model C9300 manufacturer 3", w)
	}
}

func TestManufacturerWriteToWire(t *testing.T) {
	t.Parallel()

	slug := "cisco"
	in := domain.ManufacturerWrite{Name: "Cisco", Slug: &slug}
	w := manufacturerWriteToWire(in)
	if w.Name != "Cisco" || w.Slug == nil || *w.Slug != "cisco" {
		t.Errorf("wire = %+v, want name Cisco slug cisco", w)
	}
}

func TestTenantWriteToWire(t *testing.T) {
	t.Parallel()

	slug := "acme"
	in := domain.TenantWrite{Name: "ACME", Slug: &slug}
	w := tenantWriteToWire(in)
	if w.Name != "ACME" || w.Slug == nil || *w.Slug != "acme" {
		t.Errorf("wire = %+v, want name ACME slug acme", w)
	}
}

func TestClient_CreateProvider(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"name":"ACME","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		p, err := client.CreateProvider(context.Background(), "tok", domain.ProviderWrite{Name: "ACME"})
		if err != nil {
			t.Fatalf("CreateProvider() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/circuits/providers/" {
			t.Errorf("path = %q, want /api/circuits/providers/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "ACME" {
			t.Errorf("request body = %s, want name ACME", gotBody)
		}
		if p.ID != 1 || p.Name != "ACME" {
			t.Errorf("p = %+v, want id 1 name ACME", p)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateProvider(context.Background(), "tok", domain.ProviderWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if ve.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode = %d, want 400", ve.StatusCode)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateProvider(context.Background(), "tok", domain.ProviderWrite{Name: "ACME"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			t.Errorf("err = %+v, want NOT a ValidationError for 404", ve)
		}
	})
}

func TestClient_UpdateProvider(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"name":"ACME2","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		asn := 64513
		p, err := client.UpdateProvider(context.Background(), "tok", 7, domain.ProviderWrite{Asn: &asn})
		if err != nil {
			t.Fatalf("UpdateProvider() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/circuits/providers/7/" {
			t.Errorf("path = %q, want /api/circuits/providers/7/", gotPath)
		}
		if p.ID != 7 {
			t.Errorf("p = %+v, want id 7", p)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["asn"] != float64(64513) {
			t.Errorf("request body = %s, want asn present", gotBody)
		}
		if _, ok := reqBody["name"]; ok {
			t.Errorf("request body = %s, want name omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"slug":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateProvider(context.Background(), "tok", 7, domain.ProviderWrite{})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteProvider(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteProvider(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteProvider() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/circuits/providers/3/" {
			t.Errorf("path = %q, want /api/circuits/providers/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteProvider(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_CreateTenant(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"name":"ACME","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		p, err := client.CreateTenant(context.Background(), "tok", domain.TenantWrite{Name: "ACME"})
		if err != nil {
			t.Fatalf("CreateTenant() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/tenancy/tenants/" {
			t.Errorf("path = %q, want /api/tenancy/tenants/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "ACME" {
			t.Errorf("request body = %s, want name ACME", gotBody)
		}
		if p.ID != 1 || p.Name != "ACME" {
			t.Errorf("p = %+v, want id 1 name ACME", p)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateTenant(context.Background(), "tok", domain.TenantWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if ve.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode = %d, want 400", ve.StatusCode)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateTenant(context.Background(), "tok", domain.TenantWrite{Name: "ACME"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			t.Errorf("err = %+v, want NOT a ValidationError for 404", ve)
		}
	})
}

func TestClient_UpdateTenant(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"name":"ACME2","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		desc := "new desc"
		p, err := client.UpdateTenant(context.Background(), "tok", 7, domain.TenantWrite{Description: &desc})
		if err != nil {
			t.Fatalf("UpdateTenant() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/tenancy/tenants/7/" {
			t.Errorf("path = %q, want /api/tenancy/tenants/7/", gotPath)
		}
		if p.ID != 7 {
			t.Errorf("p = %+v, want id 7", p)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["description"] != "new desc" {
			t.Errorf("request body = %s, want description present", gotBody)
		}
		if _, ok := reqBody["name"]; ok {
			t.Errorf("request body = %s, want name omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"slug":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateTenant(context.Background(), "tok", 7, domain.TenantWrite{})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteTenant(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteTenant(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteTenant() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/tenancy/tenants/3/" {
			t.Errorf("path = %q, want /api/tenancy/tenants/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteTenant(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_CreateManufacturer(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"name":"Cisco","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		p, err := client.CreateManufacturer(context.Background(), "tok", domain.ManufacturerWrite{Name: "Cisco"})
		if err != nil {
			t.Fatalf("CreateManufacturer() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/dcim/manufacturers/" {
			t.Errorf("path = %q, want /api/dcim/manufacturers/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "Cisco" {
			t.Errorf("request body = %s, want name Cisco", gotBody)
		}
		if p.ID != 1 || p.Name != "Cisco" {
			t.Errorf("p = %+v, want id 1 name Cisco", p)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateManufacturer(context.Background(), "tok", domain.ManufacturerWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if ve.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode = %d, want 400", ve.StatusCode)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateManufacturer(context.Background(), "tok", domain.ManufacturerWrite{Name: "Cisco"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			t.Errorf("err = %+v, want NOT a ValidationError for 404", ve)
		}
	})
}

func TestClient_UpdateManufacturer(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"name":"Cisco2","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		desc := "new desc"
		p, err := client.UpdateManufacturer(context.Background(), "tok", 7, domain.ManufacturerWrite{Description: &desc})
		if err != nil {
			t.Fatalf("UpdateManufacturer() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/dcim/manufacturers/7/" {
			t.Errorf("path = %q, want /api/dcim/manufacturers/7/", gotPath)
		}
		if p.ID != 7 {
			t.Errorf("p = %+v, want id 7", p)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["description"] != "new desc" {
			t.Errorf("request body = %s, want description present", gotBody)
		}
		if _, ok := reqBody["name"]; ok {
			t.Errorf("request body = %s, want name omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"slug":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateManufacturer(context.Background(), "tok", 7, domain.ManufacturerWrite{})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteManufacturer(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteManufacturer(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteManufacturer() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/dcim/manufacturers/3/" {
			t.Errorf("path = %q, want /api/dcim/manufacturers/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteManufacturer(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_CreateDeviceType(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"model":"C9300","manufacturer":{"id":3,"name":"Cisco"},"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		mf := 3
		p, err := client.CreateDeviceType(context.Background(), "tok", domain.DeviceTypeWrite{Manufacturer: &mf, Model: "C9300"})
		if err != nil {
			t.Fatalf("CreateDeviceType() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/dcim/device-types/" {
			t.Errorf("path = %q, want /api/dcim/device-types/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["model"] != "C9300" {
			t.Errorf("request body = %s, want model C9300", gotBody)
		}
		if p.ID != 1 || p.Model != "C9300" || p.Manufacturer == nil || p.Manufacturer.ID != 3 {
			t.Errorf("p = %+v, want id 1 model C9300 manufacturer 3", p)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"model":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateDeviceType(context.Background(), "tok", domain.DeviceTypeWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if ve.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode = %d, want 400", ve.StatusCode)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		mf := 3
		_, err := client.CreateDeviceType(context.Background(), "tok", domain.DeviceTypeWrite{Manufacturer: &mf, Model: "C9300"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			t.Errorf("err = %+v, want NOT a ValidationError for 404", ve)
		}
	})
}

func TestClient_UpdateDeviceType(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"model":"C9300-2","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		part := "PN-1"
		p, err := client.UpdateDeviceType(context.Background(), "tok", 7, domain.DeviceTypeWrite{PartNumber: &part})
		if err != nil {
			t.Fatalf("UpdateDeviceType() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/dcim/device-types/7/" {
			t.Errorf("path = %q, want /api/dcim/device-types/7/", gotPath)
		}
		if p.ID != 7 {
			t.Errorf("p = %+v, want id 7", p)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["part_number"] != "PN-1" {
			t.Errorf("request body = %s, want part_number present", gotBody)
		}
		if _, ok := reqBody["model"]; ok {
			t.Errorf("request body = %s, want model omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"slug":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateDeviceType(context.Background(), "tok", 7, domain.DeviceTypeWrite{})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteDeviceType(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteDeviceType(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteDeviceType() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/dcim/device-types/3/" {
			t.Errorf("path = %q, want /api/dcim/device-types/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteDeviceType(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

//nolint:gocognit // exhaustive per-path write test
func TestClient_CreateInterface(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"name":"eth0","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		dev := 5
		typ := "1000base-t"
		iface, err := client.CreateInterface(context.Background(), "tok", domain.InterfaceWrite{Name: "eth0", Device: &dev, Type: &typ})
		if err != nil {
			t.Fatalf("CreateInterface() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/dcim/interfaces/" {
			t.Errorf("path = %q, want /api/dcim/interfaces/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "eth0" {
			t.Errorf("request body = %s, want name eth0", gotBody)
		}
		if reqBody["device"] != float64(5) {
			t.Errorf("request body = %s, want device 5", gotBody)
		}
		if reqBody["type"] != "1000base-t" {
			t.Errorf("request body = %s, want type 1000base-t", gotBody)
		}
		if iface.ID != 1 || iface.Name != "eth0" {
			t.Errorf("iface = %+v, want id 1 name eth0", iface)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateInterface(context.Background(), "tok", domain.InterfaceWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if ve.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode = %d, want 400", ve.StatusCode)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateInterface(context.Background(), "tok", domain.InterfaceWrite{Name: "X"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			t.Errorf("err = %+v, want NOT a ValidationError for 404", ve)
		}
	})
}

func TestClient_UpdateInterface(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"name":"eth0","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		mtu := 9000
		iface, err := client.UpdateInterface(context.Background(), "tok", 7, domain.InterfaceWrite{MTU: &mtu})
		if err != nil {
			t.Fatalf("UpdateInterface() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/dcim/interfaces/7/" {
			t.Errorf("path = %q, want /api/dcim/interfaces/7/", gotPath)
		}
		if iface.ID != 7 {
			t.Errorf("iface = %+v, want id 7", iface)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["mtu"] != float64(9000) {
			t.Errorf("request body = %s, want mtu present", gotBody)
		}
		if _, ok := reqBody["name"]; ok {
			t.Errorf("request body = %s, want name omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"type":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateInterface(context.Background(), "tok", 7, domain.InterfaceWrite{Name: "X"})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteInterface(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteInterface(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteInterface() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/dcim/interfaces/3/" {
			t.Errorf("path = %q, want /api/dcim/interfaces/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteInterface(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteInterface(context.Background(), "tok", 999); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestInterfaceWriteToWire(t *testing.T) {
	t.Parallel()

	dev := 5
	typ := "1000base-t"
	enabled := true
	in := domain.InterfaceWrite{Name: "eth0", Device: &dev, Type: &typ, Enabled: &enabled}
	w := interfaceWriteToWire(in)
	if w.Name != "eth0" || w.Device == nil || *w.Device != 5 || w.Type == nil || *w.Type != typ || w.Enabled == nil || !*w.Enabled {
		t.Errorf("wire = %+v, want name eth0 device 5 type %s enabled true", w, typ)
	}
}

func TestClient_CreateRack(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"name":"R1","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		site := 2
		rack, err := client.CreateRack(context.Background(), "tok", domain.RackWrite{Name: "R1", Site: &site})
		if err != nil {
			t.Fatalf("CreateRack() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/dcim/racks/" {
			t.Errorf("path = %q, want /api/dcim/racks/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "R1" {
			t.Errorf("request body = %s, want name R1", gotBody)
		}
		if reqBody["site"] != float64(2) {
			t.Errorf("request body = %s, want site 2", gotBody)
		}
		if rack.ID != 1 || rack.Name != "R1" {
			t.Errorf("rack = %+v, want id 1 name R1", rack)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateRack(context.Background(), "tok", domain.RackWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if ve.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode = %d, want 400", ve.StatusCode)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateRack(context.Background(), "tok", domain.RackWrite{Name: "X"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			t.Errorf("err = %+v, want NOT a ValidationError for 404", ve)
		}
	})
}

func TestClient_UpdateRack(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"name":"R1","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		u := 42
		rack, err := client.UpdateRack(context.Background(), "tok", 7, domain.RackWrite{UHeight: &u})
		if err != nil {
			t.Fatalf("UpdateRack() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/dcim/racks/7/" {
			t.Errorf("path = %q, want /api/dcim/racks/7/", gotPath)
		}
		if rack.ID != 7 {
			t.Errorf("rack = %+v, want id 7", rack)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["u_height"] != float64(42) {
			t.Errorf("request body = %s, want u_height present", gotBody)
		}
		if _, ok := reqBody["name"]; ok {
			t.Errorf("request body = %s, want name omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"status":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateRack(context.Background(), "tok", 7, domain.RackWrite{Name: "X"})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteRack(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteRack(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteRack() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/dcim/racks/3/" {
			t.Errorf("path = %q, want /api/dcim/racks/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteRack(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteRack(context.Background(), "tok", 999); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestRackWriteToWire(t *testing.T) {
	t.Parallel()

	site := 2
	status := "active"
	in := domain.RackWrite{Name: "R1", Site: &site, Status: &status}
	w := rackWriteToWire(in)
	if w.Name != "R1" || w.Site == nil || *w.Site != 2 || w.Status == nil || *w.Status != "active" {
		t.Errorf("wire = %+v, want name R1 site 2 status active", w)
	}
}

func ptrStr(s string) *string { return &s }

func TestLocationWriteToWire(t *testing.T) {
	t.Parallel()

	site := 5
	in := domain.LocationWrite{Name: "Row A", Site: &site}
	w := locationWriteToWire(in)
	if w.Name != "Row A" || w.Site == nil || *w.Site != 5 {
		t.Errorf("wire = %+v, want name 'Row A' site 5", w)
	}
}

func TestClient_CreateLocation(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"name":"Row A","site":{"id":5,"name":"DC1"},"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		site := 5
		p, err := client.CreateLocation(context.Background(), "tok", domain.LocationWrite{Name: "Row A", Site: &site})
		if err != nil {
			t.Fatalf("CreateLocation() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/dcim/locations/" {
			t.Errorf("path = %q, want /api/dcim/locations/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "Row A" || reqBody["site"] != float64(5) {
			t.Errorf("request body = %s, want name 'Row A' site 5", gotBody)
		}
		if p.ID != 1 || p.Name != "Row A" || p.Site == nil || p.Site.ID != 5 {
			t.Errorf("p = %+v, want id 1 name 'Row A' site 5", p)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"site":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateLocation(context.Background(), "tok", domain.LocationWrite{Name: "Row A"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if ve.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode = %d, want 400", ve.StatusCode)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		site := 5
		_, err := client.CreateLocation(context.Background(), "tok", domain.LocationWrite{Name: "Row A", Site: &site})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			t.Errorf("err = %+v, want NOT a ValidationError for 404", ve)
		}
	})
}

func TestClient_UpdateLocation(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"name":"Row A","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		desc := "updated"
		p, err := client.UpdateLocation(context.Background(), "tok", 7, domain.LocationWrite{Description: &desc})
		if err != nil {
			t.Fatalf("UpdateLocation() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/dcim/locations/7/" {
			t.Errorf("path = %q, want /api/dcim/locations/7/", gotPath)
		}
		if p.ID != 7 {
			t.Errorf("p = %+v, want id 7", p)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["description"] != "updated" {
			t.Errorf("request body = %s, want description present", gotBody)
		}
		if _, ok := reqBody["name"]; ok {
			t.Errorf("request body = %s, want name omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"slug":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateLocation(context.Background(), "tok", 7, domain.LocationWrite{})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteLocation(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteLocation(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteLocation() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/dcim/locations/3/" {
			t.Errorf("path = %q, want /api/dcim/locations/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteLocation(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClusterTypeWriteToWire(t *testing.T) {
	t.Parallel()

	in := domain.ClusterTypeWrite{Name: "KVM"}
	w := clusterTypeWriteToWire(in)
	if w.Name != "KVM" {
		t.Errorf("wire = %+v, want name 'KVM'", w)
	}
}

func TestClient_CreateClusterType(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"name":"KVM","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		p, err := client.CreateClusterType(context.Background(), "tok", domain.ClusterTypeWrite{Name: "KVM"})
		if err != nil {
			t.Fatalf("CreateClusterType() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/virtualization/cluster-types/" {
			t.Errorf("path = %q, want /api/virtualization/cluster-types/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "KVM" {
			t.Errorf("request body = %s, want name 'KVM'", gotBody)
		}
		if p.ID != 1 || p.Name != "KVM" {
			t.Errorf("p = %+v, want id 1 name 'KVM'", p)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateClusterType(context.Background(), "tok", domain.ClusterTypeWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateClusterType(context.Background(), "tok", domain.ClusterTypeWrite{Name: "KVM"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			t.Errorf("err = %+v, want NOT a ValidationError for 404", ve)
		}
	})
}

func TestClient_UpdateClusterType(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"name":"KVM","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		desc := "updated"
		p, err := client.UpdateClusterType(context.Background(), "tok", 7, domain.ClusterTypeWrite{Description: &desc})
		if err != nil {
			t.Fatalf("UpdateClusterType() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/virtualization/cluster-types/7/" {
			t.Errorf("path = %q, want /api/virtualization/cluster-types/7/", gotPath)
		}
		if p.ID != 7 {
			t.Errorf("p = %+v, want id 7", p)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["description"] != "updated" {
			t.Errorf("request body = %s, want description present", gotBody)
		}
		if _, ok := reqBody["name"]; ok {
			t.Errorf("request body = %s, want name omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"slug":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateClusterType(context.Background(), "tok", 7, domain.ClusterTypeWrite{})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteClusterType(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteClusterType(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteClusterType() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/virtualization/cluster-types/3/" {
			t.Errorf("path = %q, want /api/virtualization/cluster-types/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteClusterType(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClusterGroupWriteToWire(t *testing.T) {
	t.Parallel()

	in := domain.ClusterGroupWrite{Name: "DC Clusters"}
	w := clusterGroupWriteToWire(in)
	if w.Name != "DC Clusters" {
		t.Errorf("wire = %+v, want name 'DC Clusters'", w)
	}
}

func TestClient_CreateClusterGroup(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"name":"DC Clusters","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		p, err := client.CreateClusterGroup(context.Background(), "tok", domain.ClusterGroupWrite{Name: "DC Clusters"})
		if err != nil {
			t.Fatalf("CreateClusterGroup() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/virtualization/cluster-groups/" {
			t.Errorf("path = %q, want /api/virtualization/cluster-groups/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "DC Clusters" {
			t.Errorf("request body = %s, want name 'DC Clusters'", gotBody)
		}
		if p.ID != 1 || p.Name != "DC Clusters" {
			t.Errorf("p = %+v, want id 1 name 'DC Clusters'", p)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateClusterGroup(context.Background(), "tok", domain.ClusterGroupWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateClusterGroup(context.Background(), "tok", domain.ClusterGroupWrite{Name: "DC Clusters"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			t.Errorf("err = %+v, want NOT a ValidationError for 404", ve)
		}
	})
}

func TestClient_UpdateClusterGroup(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"name":"DC Clusters","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		desc := "updated"
		p, err := client.UpdateClusterGroup(context.Background(), "tok", 7, domain.ClusterGroupWrite{Description: &desc})
		if err != nil {
			t.Fatalf("UpdateClusterGroup() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/virtualization/cluster-groups/7/" {
			t.Errorf("path = %q, want /api/virtualization/cluster-groups/7/", gotPath)
		}
		if p.ID != 7 {
			t.Errorf("p = %+v, want id 7", p)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["description"] != "updated" {
			t.Errorf("request body = %s, want description present", gotBody)
		}
		if _, ok := reqBody["name"]; ok {
			t.Errorf("request body = %s, want name omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"slug":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateClusterGroup(context.Background(), "tok", 7, domain.ClusterGroupWrite{})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteClusterGroup(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteClusterGroup(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteClusterGroup() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/virtualization/cluster-groups/3/" {
			t.Errorf("path = %q, want /api/virtualization/cluster-groups/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteClusterGroup(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestCircuitTypeWriteToWire(t *testing.T) {
	t.Parallel()

	in := domain.CircuitTypeWrite{Name: "Fiber"}
	w := circuitTypeWriteToWire(in)
	if w.Name != "Fiber" {
		t.Errorf("wire = %+v, want name 'Fiber'", w)
	}
}

func TestClient_CreateCircuitType(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"name":"Fiber","circuit_count":0,"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		p, err := client.CreateCircuitType(context.Background(), "tok", domain.CircuitTypeWrite{Name: "Fiber"})
		if err != nil {
			t.Fatalf("CreateCircuitType() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/circuits/circuit-types/" {
			t.Errorf("path = %q, want /api/circuits/circuit-types/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "Fiber" {
			t.Errorf("request body = %s, want name 'Fiber'", gotBody)
		}
		if p.ID != 1 || p.Name != "Fiber" {
			t.Errorf("p = %+v, want id 1 name 'Fiber'", p)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateCircuitType(context.Background(), "tok", domain.CircuitTypeWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateCircuitType(context.Background(), "tok", domain.CircuitTypeWrite{Name: "Fiber"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			t.Errorf("err = %+v, want NOT a ValidationError for 404", ve)
		}
	})
}

func TestClient_UpdateCircuitType(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"name":"Fiber","circuit_count":0,"created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		desc := "updated"
		p, err := client.UpdateCircuitType(context.Background(), "tok", 7, domain.CircuitTypeWrite{Description: &desc})
		if err != nil {
			t.Fatalf("UpdateCircuitType() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/circuits/circuit-types/7/" {
			t.Errorf("path = %q, want /api/circuits/circuit-types/7/", gotPath)
		}
		if p.ID != 7 {
			t.Errorf("p = %+v, want id 7", p)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["description"] != "updated" {
			t.Errorf("request body = %s, want description present", gotBody)
		}
		if _, ok := reqBody["name"]; ok {
			t.Errorf("request body = %s, want name omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"slug":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateCircuitType(context.Background(), "tok", 7, domain.CircuitTypeWrite{})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteCircuitType(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteCircuitType(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteCircuitType() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/circuits/circuit-types/3/" {
			t.Errorf("path = %q, want /api/circuits/circuit-types/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteCircuitType(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestVrfWriteToWire(t *testing.T) {
	t.Parallel()

	in := domain.VrfWrite{Name: "prod", Rd: "65000:1"}
	w := vrfWriteToWire(in)
	if w.Name != "prod" || w.Rd != "65000:1" {
		t.Errorf("wire = %+v, want name 'prod' rd '65000:1'", w)
	}
}

func TestClient_CreateVrf(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"name":"prod","rd":"65000:1","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		p, err := client.CreateVrf(context.Background(), "tok", domain.VrfWrite{Name: "prod", Rd: "65000:1"})
		if err != nil {
			t.Fatalf("CreateVrf() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/ipam/vrfs/" {
			t.Errorf("path = %q, want /api/ipam/vrfs/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "prod" || reqBody["rd"] != "65000:1" {
			t.Errorf("request body = %s, want name 'prod' rd '65000:1'", gotBody)
		}
		if p.ID != 1 || p.Name != "prod" || p.Rd != "65000:1" {
			t.Errorf("p = %+v, want id 1 name 'prod' rd '65000:1'", p)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateVrf(context.Background(), "tok", domain.VrfWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateVrf(context.Background(), "tok", domain.VrfWrite{Name: "prod", Rd: "65000:1"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			t.Errorf("err = %+v, want NOT a ValidationError for 404", ve)
		}
	})
}

func TestClient_UpdateVrf(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"name":"prod","rd":"65000:1","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		desc := "updated"
		p, err := client.UpdateVrf(context.Background(), "tok", 7, domain.VrfWrite{Description: &desc})
		if err != nil {
			t.Fatalf("UpdateVrf() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/ipam/vrfs/7/" {
			t.Errorf("path = %q, want /api/ipam/vrfs/7/", gotPath)
		}
		if p.ID != 7 {
			t.Errorf("p = %+v, want id 7", p)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["description"] != "updated" {
			t.Errorf("request body = %s, want description present", gotBody)
		}
		if _, ok := reqBody["name"]; ok {
			t.Errorf("request body = %s, want name omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"rd":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateVrf(context.Background(), "tok", 7, domain.VrfWrite{})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteVrf(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteVrf(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteVrf() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/ipam/vrfs/3/" {
			t.Errorf("path = %q, want /api/ipam/vrfs/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteVrf(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestVlanGroupWriteToWire(t *testing.T) {
	t.Parallel()

	in := domain.VlanGroupWrite{Name: "DC VLANs"}
	w := vlanGroupWriteToWire(in)
	if w.Name != "DC VLANs" {
		t.Errorf("wire = %+v, want name 'DC VLANs'", w)
	}
}

func TestClient_CreateVlanGroup(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"name":"DC VLANs","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		p, err := client.CreateVlanGroup(context.Background(), "tok", domain.VlanGroupWrite{Name: "DC VLANs"})
		if err != nil {
			t.Fatalf("CreateVlanGroup() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/ipam/vlan-groups/" {
			t.Errorf("path = %q, want /api/ipam/vlan-groups/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "DC VLANs" {
			t.Errorf("request body = %s, want name 'DC VLANs'", gotBody)
		}
		if p.ID != 1 || p.Name != "DC VLANs" {
			t.Errorf("p = %+v, want id 1 name 'DC VLANs'", p)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateVlanGroup(context.Background(), "tok", domain.VlanGroupWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateVlanGroup(context.Background(), "tok", domain.VlanGroupWrite{Name: "DC VLANs"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			t.Errorf("err = %+v, want NOT a ValidationError for 404", ve)
		}
	})
}

func TestClient_UpdateVlanGroup(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"name":"DC VLANs","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		desc := "updated"
		p, err := client.UpdateVlanGroup(context.Background(), "tok", 7, domain.VlanGroupWrite{Description: &desc})
		if err != nil {
			t.Fatalf("UpdateVlanGroup() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/ipam/vlan-groups/7/" {
			t.Errorf("path = %q, want /api/ipam/vlan-groups/7/", gotPath)
		}
		if p.ID != 7 {
			t.Errorf("p = %+v, want id 7", p)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["description"] != "updated" {
			t.Errorf("request body = %s, want description present", gotBody)
		}
		if _, ok := reqBody["name"]; ok {
			t.Errorf("request body = %s, want name omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"slug":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateVlanGroup(context.Background(), "tok", 7, domain.VlanGroupWrite{})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteVlanGroup(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteVlanGroup(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteVlanGroup() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/ipam/vlan-groups/3/" {
			t.Errorf("path = %q, want /api/ipam/vlan-groups/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteVlanGroup(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestWireVlanGroupScopeToDomain(t *testing.T) {
	t.Parallel()

	if got := wireVlanGroupScopeToDomain(nil); got != nil {
		t.Errorf("wireVlanGroupScopeToDomain(nil) = %+v, want nil", got)
	}

	got := wireVlanGroupScopeToDomain(&WireVlanGroupScope{ObjectType: "dcim.site", ObjectID: 5})
	if got == nil || got.ObjectType != "dcim.site" || got.ObjectID != 5 {
		t.Errorf("wireVlanGroupScopeToDomain = %+v, want scope_type dcim.site object_id 5", got)
	}
}

func TestWireVlanGroupToDomain(t *testing.T) {
	t.Parallel()

	minVID, maxVID := 100, 200
	w := WireVlanGroup{
		ID:           1,
		Name:         "DC VLANs",
		Scope:        &WireVlanGroupScope{ObjectType: "dcim.site", ObjectID: 5},
		MinVID:       &minVID,
		MaxVID:       &maxVID,
		CustomFields: map[string]any{"env": "prod"},
	}
	d := wireVlanGroupToDomain(w)
	if d.ID != 1 || d.Name != "DC VLANs" {
		t.Errorf("d = %+v, want id 1 name 'DC VLANs'", d)
	}
	if d.Scope == nil || d.Scope.ObjectID != 5 || d.Scope.ObjectType != "dcim.site" {
		t.Errorf("d.Scope = %+v, want dcim.site/5", d.Scope)
	}
	if d.MinVID == nil || *d.MinVID != 100 || d.MaxVID == nil || *d.MaxVID != 200 {
		t.Errorf("d.MinVID/MaxVID = %v/%v, want 100/200", d.MinVID, d.MaxVID)
	}
	if d.CustomFields["env"] != "prod" {
		t.Errorf("d.CustomFields = %+v, want env=prod", d.CustomFields)
	}
}

func TestUnmarshalNewEntitiesError(t *testing.T) {
	t.Parallel()

	bad := domain.RawObject([]byte("{invalid"))
	if _, err := unmarshalCircuitType(bad); err == nil {
		t.Error("unmarshalCircuitType(invalid) = nil error, want error")
	}
	if _, err := unmarshalVrf(bad); err == nil {
		t.Error("unmarshalVrf(invalid) = nil error, want error")
	}
	if _, err := unmarshalVlanGroup(bad); err == nil {
		t.Error("unmarshalVlanGroup(invalid) = nil error, want error")
	}
}

func TestRoleWriteToWire(t *testing.T) {
	t.Parallel()

	in := domain.RoleWrite{Name: "prod"}
	w := roleWriteToWire(in)
	if w.Name != "prod" {
		t.Errorf("wire = %+v, want name 'prod'", w)
	}
}

func TestClient_CreateRole(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"name":"prod","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		p, err := client.CreateRole(context.Background(), "tok", domain.RoleWrite{Name: "prod"})
		if err != nil {
			t.Fatalf("CreateRole() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/ipam/roles/" {
			t.Errorf("path = %q, want /api/ipam/roles/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "prod" {
			t.Errorf("request body = %s, want name 'prod'", gotBody)
		}
		if p.ID != 1 || p.Name != "prod" {
			t.Errorf("p = %+v, want id 1 name 'prod'", p)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateRole(context.Background(), "tok", domain.RoleWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateRole(context.Background(), "tok", domain.RoleWrite{Name: "prod"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			t.Errorf("err = %+v, want NOT a ValidationError for 404", ve)
		}
	})
}

func TestClient_UpdateRole(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"name":"prod","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		desc := "updated"
		p, err := client.UpdateRole(context.Background(), "tok", 7, domain.RoleWrite{Description: &desc})
		if err != nil {
			t.Fatalf("UpdateRole() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/ipam/roles/7/" {
			t.Errorf("path = %q, want /api/ipam/roles/7/", gotPath)
		}
		if p.ID != 7 {
			t.Errorf("p = %+v, want id 7", p)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["description"] != "updated" {
			t.Errorf("request body = %s, want description present", gotBody)
		}
		if _, ok := reqBody["name"]; ok {
			t.Errorf("request body = %s, want name omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateRole(context.Background(), "tok", 7, domain.RoleWrite{})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteRole(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteRole(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteRole() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/ipam/roles/3/" {
			t.Errorf("path = %q, want /api/ipam/roles/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteRole(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestContactWriteToWire(t *testing.T) {
	t.Parallel()

	in := domain.ContactWrite{Name: "ops"}
	w := contactWriteToWire(in)
	if w.Name != "ops" {
		t.Errorf("wire = %+v, want name 'ops'", w)
	}
}

func TestClient_CreateContact(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotCT = r.Header.Get("Content-Type")
			gotAuth = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1,"name":"ops","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		p, err := client.CreateContact(context.Background(), "tok", domain.ContactWrite{Name: "ops"})
		if err != nil {
			t.Fatalf("CreateContact() returned error: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/api/tenancy/contacts/" {
			t.Errorf("path = %q, want /api/tenancy/contacts/", gotPath)
		}
		if gotCT != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotCT)
		}
		if gotAuth != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "ops" {
			t.Errorf("request body = %s, want name 'ops'", gotBody)
		}
		if p.ID != 1 || p.Name != "ops" {
			t.Errorf("p = %+v, want id 1 name 'ops'", p)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateContact(context.Background(), "tok", domain.ContactWrite{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})

	t.Run("not found 404", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateContact(context.Background(), "tok", domain.ContactWrite{Name: "ops"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			t.Errorf("err = %+v, want NOT a ValidationError for 404", ve)
		}
	})
}

func TestClient_UpdateContact(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			gotBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":7,"name":"ops","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		phone := "+1000"
		p, err := client.UpdateContact(context.Background(), "tok", 7, domain.ContactWrite{Phone: &phone})
		if err != nil {
			t.Fatalf("UpdateContact() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", gotMethod)
		}
		if gotPath != "/api/tenancy/contacts/7/" {
			t.Errorf("path = %q, want /api/tenancy/contacts/7/", gotPath)
		}
		if p.ID != 7 {
			t.Errorf("p = %+v, want id 7", p)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["phone"] != "+1000" {
			t.Errorf("request body = %s, want phone present", gotBody)
		}
		if _, ok := reqBody["name"]; ok {
			t.Errorf("request body = %s, want name omitted (not provided)", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"name":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateContact(context.Background(), "tok", 7, domain.ContactWrite{})
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}

func TestClient_DeleteContact(t *testing.T) {
	t.Parallel()

	t.Run("success 204", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		if err := client.DeleteContact(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteContact() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", gotMethod)
		}
		if gotPath != "/api/tenancy/contacts/3/" {
			t.Errorf("path = %q, want /api/tenancy/contacts/3/", gotPath)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"detail":"Cannot delete object with dependent objects."}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		err := client.DeleteContact(context.Background(), "tok", 3)
		var ve *domain.ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("err = %v, want *domain.ValidationError", err)
		}
		if string(ve.Body) != body {
			t.Errorf("Body = %s, want %s", ve.Body, body)
		}
	})
}
