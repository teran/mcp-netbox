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
		if reqBody["circuit_type"] != float64(3) {
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
}

func ptrStr(s string) *string { return &s }
