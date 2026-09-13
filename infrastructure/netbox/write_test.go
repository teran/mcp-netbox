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

func ptrStr(s string) *string { return &s }
