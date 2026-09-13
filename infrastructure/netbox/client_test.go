package netbox

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"resty.dev/v3"

	"github.com/teran/mcp-netbox/domain"
)

func TestClient_ListSites(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer test-token" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if r.URL.Path != "/api/dcim/sites/" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"count": 1,
				"next": null,
				"previous": null,
				"results": [{
					"id": 1,
					"name": "Test Site",
					"slug": "test-site",
					"status": {"value": "active", "label": "Active"},
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListSites(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListSites() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("Count = %d, want %d", resp.Count, 1)
		}
		if len(resp.Results) != 1 {
			t.Fatalf("len(Results) = %d, want %d", len(resp.Results), 1)
		}
		if resp.Results[0].Name != "Test Site" {
			t.Errorf("Name = %q, want %q", resp.Results[0].Name, "Test Site")
		}
	})

	t.Run("unauthorized", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"detail": "Invalid token"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListSites(context.Background(), "bad-token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListSites(context.Background(), "token", map[string]string{"site": "nonexistent"})
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("invalid_json", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{invalid json`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListSites(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestClient_ListDevices(t *testing.T) {
	t.Parallel()

	t.Run("success with filters", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("site") != "dc1" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"count": 2,
				"results": [
					{"id": 1, "name": "router1", "site": {"id": 1, "name": "DC1", "url": "", "slug": "dc1"}, "created": "", "last_updated": ""},
					{"id": 2, "name": "switch1", "site": {"id": 1, "name": "DC1", "url": "", "slug": "dc1"}, "created": "", "last_updated": ""}
				]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListDevices(context.Background(), "token", map[string]string{"site": "dc1"})
		if err != nil {
			t.Fatalf("ListDevices() returned error: %v", err)
		}
		if resp.Count != 2 {
			t.Errorf("Count = %d, want %d", resp.Count, 2)
		}
	})
}

func TestClient_ListIPAddresses(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"address": "10.0.0.1/24",
					"status": {"value": "active", "label": "Active"},
					"dns_name": "server.example.com",
					"created": "",
					"last_updated": ""
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListIPAddresses(context.Background(), "token", nil)
		if err != nil {
			t.Fatalf("ListIPAddresses() returned error: %v", err)
		}
		if len(resp.Results) != 1 {
			t.Fatalf("len(Results) = %d, want %d", len(resp.Results), 1)
		}
		if resp.Results[0].Address != "10.0.0.1/24" {
			t.Errorf("Address = %q, want %q", resp.Results[0].Address, "10.0.0.1/24")
		}
		if resp.Results[0].DNSName != "server.example.com" {
			t.Errorf("DNSName = %q, want %q", resp.Results[0].DNSName, "server.example.com")
		}
	})
}

func TestClient_ListPrefixes(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"prefix": "10.0.0.0/8",
					"family": {"value": 4, "label": "IPv4"},
					"status": {"value": "active", "label": "Active"},
					"is_pool": false,
					"children": 0,
					"_depth": 0,
					"created": "",
					"last_updated": ""
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListPrefixes(context.Background(), "token", nil)
		if err != nil {
			t.Fatalf("ListPrefixes() returned error: %v", err)
		}
		if len(resp.Results) != 1 {
			t.Fatalf("len(Results) = %d, want %d", len(resp.Results), 1)
		}
		if resp.Results[0].Prefix != "10.0.0.0/8" {
			t.Errorf("Prefix = %q, want %q", resp.Results[0].Prefix, "10.0.0.0/8")
		}
	})
}

func TestClient_ListVLANs(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"vid": 100,
					"name": "Users",
					"created": "",
					"last_updated": ""
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListVLANs(context.Background(), "token", nil)
		if err != nil {
			t.Fatalf("ListVLANs() returned error: %v", err)
		}
		if len(resp.Results) != 1 {
			t.Fatalf("len(Results) = %d, want %d", len(resp.Results), 1)
		}
		if resp.Results[0].Name != "Users" {
			t.Errorf("Name = %q, want %q", resp.Results[0].Name, "Users")
		}
	})
}

func TestClient_ListVirtualMachines(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"name": "web-01",
					"vcpus": 4,
					"memory": 16384,
					"disk": 100,
					"created": "",
					"last_updated": ""
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListVirtualMachines(context.Background(), "token", nil)
		if err != nil {
			t.Fatalf("ListVirtualMachines() returned error: %v", err)
		}
		if len(resp.Results) != 1 {
			t.Fatalf("len(Results) = %d, want %d", len(resp.Results), 1)
		}
		if resp.Results[0].Name != "web-01" {
			t.Errorf("Name = %q, want %q", resp.Results[0].Name, "web-01")
		}
	})
}

func TestClient_ListClusters(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"name": "Prod Cluster",
					"type": {"id": 1, "name": "VMware", "url": "", "slug": "vmware"},
					"created": "",
					"last_updated": ""
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListClusters(context.Background(), "token", nil)
		if err != nil {
			t.Fatalf("ListClusters() returned error: %v", err)
		}
		if len(resp.Results) != 1 {
			t.Fatalf("len(Results) = %d, want %d", len(resp.Results), 1)
		}
		if resp.Results[0].Name != "Prod Cluster" {
			t.Errorf("Name = %q, want %q", resp.Results[0].Name, "Prod Cluster")
		}
	})
}

func TestClient_ListCircuits(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"cid": "CIR-001",
					"provider": {"id": 1, "name": "ATT", "url": "", "slug": "att"},
					"circuit_type": {"id": 1, "name": "Dark Fiber", "url": "", "slug": "dark-fiber"},
					"status": {"value": "active", "label": "Active"},
					"created": "",
					"last_updated": ""
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListCircuits(context.Background(), "token", nil)
		if err != nil {
			t.Fatalf("ListCircuits() returned error: %v", err)
		}
		if len(resp.Results) != 1 {
			t.Fatalf("len(Results) = %d, want %d", len(resp.Results), 1)
		}
		if resp.Results[0].CID != "CIR-001" {
			t.Errorf("CID = %q, want %q", resp.Results[0].CID, "CIR-001")
		}
	})
}

func TestClient_ListInterfaces(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/dcim/interfaces/" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"count": 2,
				"next": null,
				"previous": null,
				"results": [
					{
						"id": 1,
						"name": "ge-0/0/0",
						"device": {"id": 1, "name": "router1", "url": "", "slug": "router1"},
						"type": {"value": "1000base-t", "label": "1000BASE-T"},
						"enabled": true,
						"mtu": 1500,
						"mac_address": "aa:bb:cc:dd:ee:ff",
						"speed": 1000000,
						"description": "uplink",
						"tags": [],
						"created": "",
						"last_updated": ""
					},
					{
						"id": 2,
						"name": "ge-0/0/1",
						"device": {"id": 1, "name": "router1", "url": "", "slug": "router1"},
						"type": {"value": "1000base-t", "label": "1000BASE-T"},
						"enabled": false,
						"description": "",
						"tags": [],
						"created": "",
						"last_updated": ""
					}
				]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListInterfaces(context.Background(), "token", nil)
		if err != nil {
			t.Fatalf("ListInterfaces() returned error: %v", err)
		}
		if resp.Count != 2 {
			t.Errorf("Count = %d, want %d", resp.Count, 2)
		}
		if len(resp.Results) != 2 {
			t.Fatalf("len(Results) = %d, want %d", len(resp.Results), 2)
		}
		if resp.Results[0].Name != "ge-0/0/0" {
			t.Errorf("Results[0].Name = %q, want %q", resp.Results[0].Name, "ge-0/0/0")
		}
		if !resp.Results[0].Enabled {
			t.Error("Results[0].Enabled = false, want true")
		}
		if resp.Results[0].Device.Name != "router1" {
			t.Errorf("Results[0].Device.Name = %q, want %q", resp.Results[0].Device.Name, "router1")
		}
		if resp.Results[1].Enabled {
			t.Error("Results[1].Enabled = true, want false")
		}
	})
}

func TestClient_ListVMInterfaces(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/virtualization/interfaces/" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"count": 2,
				"next": null,
				"previous": null,
				"results": [
					{
						"id": 1,
						"name": "eth0",
						"virtual_machine": {"id": 1, "name": "web-01", "url": "", "slug": "web-01"},
						"enabled": true,
						"mtu": 1500,
						"mac_address": "11:22:33:44:55:66",
						"description": "primary",
						"tags": [],
						"created": "",
						"last_updated": ""
					},
					{
						"id": 2,
						"name": "eth1",
						"virtual_machine": {"id": 1, "name": "web-01", "url": "", "slug": "web-01"},
						"enabled": false,
						"description": "",
						"tags": [],
						"created": "",
						"last_updated": ""
					}
				]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListVMInterfaces(context.Background(), "token", nil)
		if err != nil {
			t.Fatalf("ListVMInterfaces() returned error: %v", err)
		}
		if resp.Count != 2 {
			t.Errorf("Count = %d, want %d", resp.Count, 2)
		}
		if len(resp.Results) != 2 {
			t.Fatalf("len(Results) = %d, want %d", len(resp.Results), 2)
		}
		if resp.Results[0].Name != "eth0" {
			t.Errorf("Results[0].Name = %q, want %q", resp.Results[0].Name, "eth0")
		}
		if resp.Results[0].VirtualMachine.Name != "web-01" {
			t.Errorf("Results[0].VirtualMachine.Name = %q, want %q", resp.Results[0].VirtualMachine.Name, "web-01")
		}
		if resp.Results[0].MTU != 1500 {
			t.Errorf("Results[0].MTU = %d, want %d", resp.Results[0].MTU, 1500)
		}
	})
}

func TestClient_ListCircuitTerminations(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/circuits/circuit-terminations/" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"count": 2,
				"next": null,
				"previous": null,
				"results": [
					{
						"id": 1,
						"circuit": {"id": 1, "name": "CIR-001", "url": "", "slug": "cir-001"},
						"term_side": "A",
						"site": {"id": 1, "name": "DC1", "url": "", "slug": "dc1"},
						"speed": 10000000,
						"upstream_speed": 10000000,
						"description": "primary termination",
						"tags": [],
						"created": "",
						"last_updated": ""
					},
					{
						"id": 2,
						"circuit": {"id": 1, "name": "CIR-001", "url": "", "slug": "cir-001"},
						"term_side": "Z",
						"site": {"id": 2, "name": "DC2", "url": "", "slug": "dc2"},
						"description": "",
						"tags": [],
						"created": "",
						"last_updated": ""
					}
				]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListCircuitTerminations(context.Background(), "token", nil)
		if err != nil {
			t.Fatalf("ListCircuitTerminations() returned error: %v", err)
		}
		if resp.Count != 2 {
			t.Errorf("Count = %d, want %d", resp.Count, 2)
		}
		if len(resp.Results) != 2 {
			t.Fatalf("len(Results) = %d, want %d", len(resp.Results), 2)
		}
		if resp.Results[0].TermSide != "A" {
			t.Errorf("Results[0].TermSide = %q, want %q", resp.Results[0].TermSide, "A")
		}
		if resp.Results[1].TermSide != "Z" {
			t.Errorf("Results[1].TermSide = %q, want %q", resp.Results[1].TermSide, "Z")
		}
		if resp.Results[0].Circuit.Name != "CIR-001" {
			t.Errorf("Results[0].Circuit.Name = %q, want %q", resp.Results[0].Circuit.Name, "CIR-001")
		}
	})
}

func TestClient_ListCables(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/dcim/cables/" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"count": 2,
				"next": null,
				"previous": null,
				"results": [
					{
						"id": 1,
						"type": {"value": "cat6a", "label": "CAT6a"},
						"status": {"value": "connected", "label": "Connected"},
						"label": "link-01",
						"color": "blue",
						"length": 10.5,
						"length_unit": {"value": "m", "label": "Meters"},
						"description": "server uplink",
						"tags": [{"id": 1, "name": "prod", "slug": "prod", "url": ""}],
						"created": "",
						"last_updated": ""
					},
					{
						"id": 2,
						"type": null,
						"status": {"value": "planned", "label": "Planned"},
						"label": "",
						"color": "",
						"description": "",
						"tags": [],
						"created": "",
						"last_updated": ""
					}
				]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListCables(context.Background(), "token", nil)
		if err != nil {
			t.Fatalf("ListCables() returned error: %v", err)
		}
		if resp.Count != 2 {
			t.Errorf("Count = %d, want %d", resp.Count, 2)
		}
		if len(resp.Results) != 2 {
			t.Fatalf("len(Results) = %d, want %d", len(resp.Results), 2)
		}
		if resp.Results[0].Label != "link-01" {
			t.Errorf("Results[0].Label = %q, want %q", resp.Results[0].Label, "link-01")
		}
		if resp.Results[0].Type.Value != "cat6a" {
			t.Errorf("Results[0].Type.Value = %q, want %q", resp.Results[0].Type.Value, "cat6a")
		}
		if resp.Results[1].Type != nil {
			t.Errorf("Results[1].Type = %+v, want nil", resp.Results[1].Type)
		}
	})
}

func TestClient_ListRacks(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"name": "Rack-01",
					"u_height": 42,
					"width": 19,
					"created": "",
					"last_updated": ""
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListRacks(context.Background(), "token", nil)
		if err != nil {
			t.Fatalf("ListRacks() returned error: %v", err)
		}
		if len(resp.Results) != 1 {
			t.Fatalf("len(Results) = %d, want %d", len(resp.Results), 1)
		}
		if resp.Results[0].Name != "Rack-01" {
			t.Errorf("Name = %q, want %q", resp.Results[0].Name, "Rack-01")
		}
	})
}

func TestClient_GetObject(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id": 1, "name": "Test", "slug": "test"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.GetObject(context.Background(), "token", "site", 1, nil)
		if err != nil {
			t.Fatalf("GetObject() returned error: %v", err)
		}
		var result map[string]interface{}
		if err := json.Unmarshal(resp, &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if result["name"] != "Test" {
			t.Errorf("name = %v, want %v", result["name"], "Test")
		}
	})

	t.Run("unknown type", func(t *testing.T) {
		client := NewClient("http://example.com", http.DefaultClient)
		_, err := client.GetObject(context.Background(), "token", "nonexistent", 1, nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("not found", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.GetObject(context.Background(), "token", "site", 999, nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestClient_HTTPErrors(t *testing.T) {
	t.Parallel()

	t.Run("forbidden", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListSites(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("server error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListSites(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("rate limited (429) with Retry-After", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Retry-After", "30")
			w.WriteHeader(http.StatusTooManyRequests)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListSites(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("request failure", func(t *testing.T) {
		client := NewClient("http://127.0.0.1:1", http.DefaultClient)
		_, err := client.ListSites(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

// TestNewClientUsesResty verifies the client is backed by a resty.Client and
// that outbound requests forward the Accept/Authorization headers and query
// parameters correctly through resty.
func TestNewClientUsesResty(t *testing.T) {
	t.Parallel()

	client := NewClient("http://example.com", http.DefaultClient)
	if client.client == nil {
		t.Fatal("NewClient did not initialise the resty client")
	}

	var sawAccept, sawAuth, sawQuery bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawAccept = r.Header.Get("Accept") == "application/json"
		sawAuth = r.Header.Get("Authorization") == "Bearer tok"
		sawQuery = r.URL.Query().Get("name") == "site-a" && r.URL.Query().Get("status") == "active"
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"count":0,"results":[]}`))
	}))
	defer srv.Close()

	rc := NewClient(srv.URL, http.DefaultClient)
	if _, ok := interface{}(rc.client).(*resty.Client); !ok {
		t.Fatalf("client is %T, want *resty.Client", rc.client)
	}

	if _, err := rc.ListSites(context.Background(), "tok", map[string]string{"name": "site-a", "status": "active"}); err != nil {
		t.Fatalf("ListSites() returned error: %v", err)
	}
	if !sawAccept {
		t.Error("Accept header was not forwarded")
	}
	if !sawAuth {
		t.Error("Authorization header was not forwarded")
	}
	if !sawQuery {
		t.Error("query params were not forwarded")
	}
}

func TestDefaultLimit(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"count": 0, "results": []}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	resp, err := client.ListSites(context.Background(), "token", nil)
	if err != nil {
		t.Fatalf("ListSites() returned error: %v", err)
	}
	if len(resp.Results) != 0 {
		t.Errorf("len(Results) = %d, want 0", len(resp.Results))
	}
}

func TestWireObjectTypeToEndpoint(t *testing.T) {
	t.Parallel()

	cases := []struct {
		objectType string
		expected   string
	}{
		{"site", "/api/dcim/sites/"},
		{"device", "/api/dcim/devices/"},
		{"prefix", "/api/ipam/prefixes/"},
		{"ip_address", "/api/ipam/ip-addresses/"},
		{"vlan", "/api/ipam/vlans/"},
		{"virtual_machine", "/api/virtualization/virtual-machines/"},
		{"cluster", "/api/virtualization/clusters/"},
		{"circuit", "/api/circuits/circuits/"},
		{"provider", "/api/circuits/providers/"},
		{"tenant", "/api/tenancy/tenants/"},
		{"rack", "/api/dcim/racks/"},
		{"manufacturer", "/api/dcim/manufacturers/"},
		{"device_type", "/api/dcim/device-types/"},
		{"location", "/api/dcim/locations/"},
		{"cluster_type", "/api/virtualization/cluster-types/"},
		{"cluster_group", "/api/virtualization/cluster-groups/"},
		{"circuit_type", "/api/circuits/circuit-types/"},
		{"vrf", "/api/ipam/vrfs/"},
		{"vlan_group", "/api/ipam/vlan-groups/"},
		{"role", "/api/ipam/roles/"},
		{"contact", "/api/tenancy/contacts/"},
		{"cable", "/api/dcim/cables/"},
	}

	for _, tc := range cases {
		t.Run(tc.objectType, func(t *testing.T) {
			endpoint, ok := objectTypeToEndpoint[tc.objectType]
			if !ok {
				t.Fatalf("object type %q not found in map", tc.objectType)
			}
			if endpoint != tc.expected {
				t.Errorf("endpoint = %q, want %q", endpoint, tc.expected)
			}
		})
	}
}

func TestConversionFunctions_NilSafety(t *testing.T) {
	t.Parallel()

	t.Run("wireNestedToDomain nil", func(t *testing.T) {
		if wireNestedToDomain(nil) != nil {
			t.Error("expected nil")
		}
	})

	t.Run("wireLabelToDomain nil", func(t *testing.T) {
		if wireLabelToDomain(nil) != nil {
			t.Error("expected nil")
		}
	})

	t.Run("wireTagsToDomain nil", func(t *testing.T) {
		if wireTagsToDomain(nil) != nil {
			t.Error("expected nil")
		}
	})

	t.Run("wireAssignedObjToDomain nil", func(t *testing.T) {
		if wireAssignedObjToDomain(nil) != nil {
			t.Error("expected nil")
		}
	})
}

func TestWireInterfaceToDomain(t *testing.T) {
	t.Parallel()

	w := WireInterface{
		ID:      1,
		Name:    "ge-0/0/0",
		Device:  &WireNested{ID: 1, Name: "router1", Slug: "router1", URL: ""},
		Type:    &WireLabel{Value: "1000base-t", Label: "1000BASE-T"},
		Enabled: true,
		MTU:     1500,
	}
	d := wireInterfaceToDomain(w)
	if d.Name != "ge-0/0/0" {
		t.Errorf("Name = %q, want %q", d.Name, "ge-0/0/0")
	}
	if !d.Enabled {
		t.Error("Enabled = false, want true")
	}
	if d.Device.Name != "router1" {
		t.Errorf("Device.Name = %q, want %q", d.Device.Name, "router1")
	}
}

func TestWireVMInterfaceToDomain(t *testing.T) {
	t.Parallel()

	w := WireVMInterface{
		ID:             1,
		Name:           "eth0",
		VirtualMachine: &WireNested{ID: 1, Name: "web-01", Slug: "web-01", URL: ""},
		Enabled:        true,
	}
	d := wireVMInterfaceToDomain(w)
	if d.Name != "eth0" {
		t.Errorf("Name = %q, want %q", d.Name, "eth0")
	}
	if d.VirtualMachine.Name != "web-01" {
		t.Errorf("VirtualMachine.Name = %q, want %q", d.VirtualMachine.Name, "web-01")
	}
}

func TestWireCircuitTerminationToDomain(t *testing.T) {
	t.Parallel()

	w := WireCircuitTermination{
		ID:       1,
		Circuit:  &WireNested{ID: 1, Name: "CIR-001", Slug: "cir-001", URL: ""},
		TermSide: "A",
		Site:     &WireNested{ID: 1, Name: "DC1", Slug: "dc1", URL: ""},
		Speed:    10000000,
	}
	d := wireCircuitTerminationToDomain(w)
	if d.TermSide != "A" {
		t.Errorf("TermSide = %q, want %q", d.TermSide, "A")
	}
	if d.Circuit.Name != "CIR-001" {
		t.Errorf("Circuit.Name = %q, want %q", d.Circuit.Name, "CIR-001")
	}
}

func TestWireCableToDomain(t *testing.T) {
	t.Parallel()

	w := WireCable{
		ID:     1,
		Type:   &WireLabel{Value: "cat6a", Label: "CAT6a"},
		Status: &WireLabel{Value: "connected", Label: "Connected"},
		Label:  "link-01",
		Color:  "blue",
		Length: 10.5,
	}
	d := wireCableToDomain(w)
	if d.Label != "link-01" {
		t.Errorf("Label = %q, want %q", d.Label, "link-01")
	}
	if d.Type.Value != "cat6a" {
		t.Errorf("Type.Value = %q, want %q", d.Type.Value, "cat6a")
	}
	if d.Status.Value != "connected" {
		t.Errorf("Status.Value = %q, want %q", d.Status.Value, "connected")
	}
}

func TestWireCableToDomain_NilTypeAndStatus(t *testing.T) {
	t.Parallel()

	w := WireCable{
		ID:     1,
		Type:   nil,
		Status: nil,
		Tags:   nil,
	}
	d := wireCableToDomain(w)
	if d.Type != nil {
		t.Errorf("Type = %+v, want nil", d.Type)
	}
	if d.Status != nil {
		t.Errorf("Status = %+v, want nil", d.Status)
	}
}

func TestWireVLANToDomain(t *testing.T) {
	t.Parallel()

	w := WireVLAN{
		ID:   1,
		VID:  100,
		Name: "Prod-VLAN",
		Status: &WireLabel{
			Value: "active",
			Label: "Active",
		},
		Site: &WireNested{ID: 1, Name: "DC-1", Slug: "dc-1", URL: ""},
	}
	d := wireVLANToDomain(w)
	if d.VID != 100 {
		t.Errorf("VID = %d, want %d", d.VID, 100)
	}
	if d.Name != "Prod-VLAN" {
		t.Errorf("Name = %q, want %q", d.Name, "Prod-VLAN")
	}
	if d.Status.Value != "active" {
		t.Errorf("Status.Value = %q, want %q", d.Status.Value, "active")
	}
	if d.Site.Name != "DC-1" {
		t.Errorf("Site.Name = %q, want %q", d.Site.Name, "DC-1")
	}
}

func TestWireVirtualMachineToDomain(t *testing.T) {
	t.Parallel()

	w := WireVirtualMachine{
		ID:   1,
		Name: "web-01",
		Status: &WireLabel{
			Value: "active",
			Label: "Active",
		},
		VCPUs:  4,
		Memory: 8192,
		Disk:   100,
	}
	d := wireVirtualMachineToDomain(w)
	if d.Name != "web-01" {
		t.Errorf("Name = %q, want %q", d.Name, "web-01")
	}
	if d.VCPUs != 4 {
		t.Errorf("VCPUs = %f, want %f", d.VCPUs, 4.0)
	}
	if d.Memory != 8192 {
		t.Errorf("Memory = %d, want %d", d.Memory, 8192)
	}
	if d.Disk != 100 {
		t.Errorf("Disk = %d, want %d", d.Disk, 100)
	}
	if d.Status.Value != "active" {
		t.Errorf("Status.Value = %q, want %q", d.Status.Value, "active")
	}
}

func TestWireClusterToDomain(t *testing.T) {
	t.Parallel()

	w := WireCluster{
		ID:   1,
		Name: "K8s-Prod",
		ClusterType: &WireNested{
			ID: 1, Name: "Kubernetes", Slug: "kubernetes", URL: "",
		},
		ClusterGroup: &WireNested{
			ID: 2, Name: "Production", Slug: "production", URL: "",
		},
		Site: &WireNested{ID: 1, Name: "DC-1", Slug: "dc-1", URL: ""},
	}
	d := wireClusterToDomain(w)
	if d.Name != "K8s-Prod" {
		t.Errorf("Name = %q, want %q", d.Name, "K8s-Prod")
	}
	if d.ClusterType.Name != "Kubernetes" {
		t.Errorf("ClusterType.Name = %q, want %q", d.ClusterType.Name, "Kubernetes")
	}
	if d.ClusterGroup.Name != "Production" {
		t.Errorf("ClusterGroup.Name = %q, want %q", d.ClusterGroup.Name, "Production")
	}
	if d.Site.Name != "DC-1" {
		t.Errorf("Site.Name = %q, want %q", d.Site.Name, "DC-1")
	}
}

func TestWireCircuitToDomain(t *testing.T) {
	t.Parallel()

	w := WireCircuit{
		ID:  1,
		CID: "CIR-001",
		Provider: &WireNested{
			ID: 1, Name: "Equinix", Slug: "equinix", URL: "",
		},
		CircuitType: &WireNested{
			ID: 1, Name: "Dark Fiber", Slug: "dark-fiber", URL: "",
		},
		Status: &WireLabel{
			Value: "active",
			Label: "Active",
		},
	}
	d := wireCircuitToDomain(w)
	if d.CID != "CIR-001" {
		t.Errorf("CID = %q, want %q", d.CID, "CIR-001")
	}
	if d.Provider.Name != "Equinix" {
		t.Errorf("Provider.Name = %q, want %q", d.Provider.Name, "Equinix")
	}
	if d.CircuitType.Name != "Dark Fiber" {
		t.Errorf("CircuitType.Name = %q, want %q", d.CircuitType.Name, "Dark Fiber")
	}
	if d.Status.Value != "active" {
		t.Errorf("Status.Value = %q, want %q", d.Status.Value, "active")
	}
}

func TestWireRackToDomain(t *testing.T) {
	t.Parallel()

	w := WireRack{
		ID:         1,
		Name:       "Rack-101",
		FacilityID: "A1-01",
		Site:       &WireNested{ID: 1, Name: "DC-1", Slug: "dc-1", URL: ""},
		Width:      19,
		UHeight:    42,
	}
	d := wireRackToDomain(w)
	if d.Name != "Rack-101" {
		t.Errorf("Name = %q, want %q", d.Name, "Rack-101")
	}
	if d.FacilityID != "A1-01" {
		t.Errorf("FacilityID = %q, want %q", d.FacilityID, "A1-01")
	}
	if d.Site.Name != "DC-1" {
		t.Errorf("Site.Name = %q, want %q", d.Site.Name, "DC-1")
	}
	if d.Width != 19 {
		t.Errorf("Width = %d, want %d", d.Width, 19)
	}
	if d.UHeight != 42 {
		t.Errorf("UHeight = %d, want %d", d.UHeight, 42)
	}
}

func TestConversionFunctions_Populated(t *testing.T) {
	t.Parallel()

	t.Run("wireNestedToDomain", func(t *testing.T) {
		w := &WireNested{ID: 1, Name: "Test", Slug: "test", URL: "http://example.com"}
		d := wireNestedToDomain(w)
		if d.ID != 1 || d.Name != "Test" || d.Slug != "test" {
			t.Errorf("unexpected result: %+v", d)
		}
	})

	t.Run("wireLabelToDomain", func(t *testing.T) {
		w := &WireLabel{Value: "active", Label: "Active"}
		d := wireLabelToDomain(w)
		if d.Value != "active" || d.Label != "Active" {
			t.Errorf("unexpected result: %+v", d)
		}
	})

	t.Run("wireTagsToDomain", func(t *testing.T) {
		w := []WireTag{{ID: 1, Name: "prod", Slug: "prod"}}
		d := wireTagsToDomain(w)
		if len(d) != 1 || d[0].Name != "prod" {
			t.Errorf("unexpected result: %+v", d)
		}
	})
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	client := NewClient("http://netbox:8000/", http.DefaultClient)
	if client.baseURL != "http://netbox:8000" {
		t.Errorf("baseURL = %q, want %q", client.baseURL, "http://netbox:8000")
	}

	client2 := NewClient("http://netbox:8000", http.DefaultClient)
	if client2.baseURL != "http://netbox:8000" {
		t.Errorf("baseURL = %q, want %q", client2.baseURL, "http://netbox:8000")
	}
}

func TestJSONUnmarshalErrors(t *testing.T) {
	t.Parallel()

	t.Run("unmarshal prefix with bad data", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{invalid`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListPrefixes(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestWirePrefixFamilyNil(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"count": 0, "results": [{"id": 1, "prefix": "10.0.0.0/8", "family": null, "created": "", "last_updated": ""}]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	resp, err := client.ListPrefixes(context.Background(), "token", nil)
	if err != nil {
		t.Fatalf("ListPrefixes() returned error: %v", err)
	}
	if resp.Results[0].Family != nil {
		t.Errorf("Family = %+v, want nil", resp.Results[0].Family)
	}
}

func TestWireAssignedObjConversion(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"count": 1,
			"results": [{
				"id": 1,
				"address": "10.0.0.1/24",
				"status": {"value": "active", "label": "Active"},
				"assigned_object_type": "dcim.interface",
				"assigned_object_id": 100,
				"assigned_object": {"id": 100, "name": "ge-0/0/0", "device": {"id": 1, "name": "router1", "url": "", "slug": "router1"}, "url": ""},
				"created": "",
				"last_updated": ""
			}]
		}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	resp, err := client.ListIPAddresses(context.Background(), "token", nil)
	if err != nil {
		t.Fatalf("ListIPAddresses() returned error: %v", err)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("len(Results) = %d, want %d", len(resp.Results), 1)
	}
	ip := resp.Results[0]
	if ip.AssignedObject == nil {
		t.Fatal("AssignedObject is nil")
	}
	if ip.AssignedObject.Name != "ge-0/0/0" {
		t.Errorf("AssignedObject.Name = %q, want %q", ip.AssignedObject.Name, "ge-0/0/0")
	}
	if ip.AssignedObject.Device.Name != "router1" {
		t.Errorf("Device.Name = %q, want %q", ip.AssignedObject.Device.Name, "router1")
	}
}

func TestJSONRoundTripWireSites(t *testing.T) {
	t.Parallel()

	raw := `{
		"id": 1,
		"name": "Test",
		"slug": "test",
		"status": {"value": "active", "label": "Active"},
		"region": {"id": 1, "name": "US", "url": "", "slug": "us"},
		"tenant": {"id": 1, "name": "Acme", "url": "", "slug": "acme"},
		"facility": "Facility1",
		"time_zone": "America/New_York",
		"description": "desc",
		"physical_address": "123 Main St",
		"shipping_address": "456 Oak Ave",
		"comments": "comment",
		"tags": [{"id": 1, "name": "prod", "slug": "prod", "url": ""}],
		"created": "2024-01-01",
		"last_updated": "2024-01-01T00:00:00Z"
	}`

	var w WireSite
	if err := json.Unmarshal([]byte(raw), &w); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	d := wireSiteToDomain(w)

	if d.ID != 1 || d.Name != "Test" || d.Slug != "test" {
		t.Errorf("basic fields mismatch: %+v", d)
	}
	if d.Region.Name != "US" {
		t.Errorf("Region.Name = %q, want %q", d.Region.Name, "US")
	}
	if d.Tenant.Name != "Acme" {
		t.Errorf("Tenant.Name = %q, want %q", d.Tenant.Name, "Acme")
	}
	if d.Tenant.Name != "Acme" {
		t.Errorf("Tenant.Name = %q, want %q", d.Tenant.Name, "Acme")
	}
	if len(d.Tags) != 1 || d.Tags[0].Name != "prod" {
		t.Errorf("tags mismatch: %+v", d.Tags)
	}
}

// --- doRequest direct tests ---

func TestDoRequest_Success(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"test"}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	body, err := client.doRequest(context.Background(), "test-token", http.MethodGet, "/api/dcim/sites/", nil)
	if err != nil {
		t.Fatalf("doRequest() returned error: %v", err)
	}
	if string(body) != `{"name":"test"}` {
		t.Errorf("body = %q, want %q", string(body), `{"name":"test"}`)
	}
}

func TestDoRequest_Unauthorized(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	_, err := client.doRequest(context.Background(), "bad-token", http.MethodGet, "/api/dcim/sites/", nil)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !contains(err.Error(), "unauthorized") {
		t.Errorf("error = %q, want it to contain 'unauthorized'", err.Error())
	}
}

func TestDoRequest_Forbidden(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	_, err := client.doRequest(context.Background(), "token", http.MethodGet, "/api/dcim/sites/", nil)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !contains(err.Error(), "forbidden") {
		t.Errorf("error = %q, want it to contain 'forbidden'", err.Error())
	}
}

func TestDoRequest_NotFound(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	_, err := client.doRequest(context.Background(), "token", http.MethodGet, "/api/dcim/sites/", nil)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !contains(err.Error(), "not found") {
		t.Errorf("error = %q, want it to contain 'not found'", err.Error())
	}
}

func TestDoRequest_RateLimited(t *testing.T) {
	t.Parallel()

	t.Run("without Retry-After", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.doRequest(context.Background(), "token", http.MethodGet, "/api/dcim/sites/", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !contains(err.Error(), "rate limited") {
			t.Errorf("error = %q, want it to contain 'rate limited'", err.Error())
		}
	})

	t.Run("with Retry-After", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Retry-After", "5")
			w.WriteHeader(http.StatusTooManyRequests)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.doRequest(context.Background(), "token", http.MethodGet, "/api/dcim/sites/", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !contains(err.Error(), "retry after 5s") {
			t.Errorf("error = %q, want it to contain 'retry after 5s'", err.Error())
		}
	})
}

func TestDoRequest_UnexpectedStatus(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"detail":"internal error"}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	_, err := client.doRequest(context.Background(), "token", http.MethodGet, "/api/dcim/sites/", nil)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !contains(err.Error(), "unexpected status 500") {
		t.Errorf("error = %q, want it to contain 'unexpected status 500'", err.Error())
	}
}

func TestDoRequest_QueryParams(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("name") != "test-site" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.URL.Query().Get("status") != "active" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"ok"}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	body, err := client.doRequest(context.Background(), "token", http.MethodGet, "/api/dcim/sites/", map[string]string{
		"name":   "test-site",
		"status": "active",
	})
	if err != nil {
		t.Fatalf("doRequest() returned error: %v", err)
	}
	if string(body) != `{"result":"ok"}` {
		t.Errorf("body = %q, want %q", string(body), `{"result":"ok"}`)
	}
}

// TestDoRequest_WithBodyLimit verifies that doRequest returns an error when
// the response body exceeds 10 MB.
func TestDoRequest_WithBodyLimit(t *testing.T) {
	t.Parallel()

	// Generate a payload larger than 10 MB.
	size := 11 * 1024 * 1024
	largeBody := []byte(`{"data":"`)
	largeBody = append(largeBody, bytes.Repeat([]byte("x"), size)...)
	largeBody = append(largeBody, `"}`...)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(largeBody)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	_, err := client.doRequest(context.Background(), "token", http.MethodGet, "/api/dcim/sites/", nil)
	if err == nil {
		t.Fatal("Expected error due to body size limit, got nil")
	}
	if !contains(err.Error(), "exceeds") && !contains(err.Error(), "truncation") {
		t.Errorf("error = %q, want it to contain 'exceeds' or 'truncation'", err.Error())
	}
}

// TestDoRequest_GetObjectUnknownType tests the object type lookup error path.
func TestDoRequest_GetObjectUnknownType(t *testing.T) {
	t.Parallel()

	client := NewClient("http://example.com", http.DefaultClient)
	_, err := client.GetObject(context.Background(), "token", "nonexistent_type", 1, nil)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !contains(err.Error(), "unknown object type") {
		t.Errorf("error = %q, want it to contain 'unknown object type'", err.Error())
	}
}

// --- convertPaginated generic function test ---

func TestConvertPaginated(t *testing.T) {
	t.Parallel()

	// Dummy conversion function: WireSite -> string (just the name).
	convert := func(w WireSite) string {
		return w.Name
	}

	resp := &domain.PaginatedResponse[WireSite]{
		Count:    2,
		Next:     "http://example.com/api/dcim/sites/?page=2",
		Previous: "http://example.com/api/dcim/sites/?page=1",
		Results: []WireSite{
			{ID: 1, Name: "Site A", Slug: "site-a"},
			{ID: 2, Name: "Site B", Slug: "site-b"},
		},
	}

	converted := convertPaginated(resp, convert)

	if converted.Count != 2 {
		t.Errorf("Count = %d, want %d", converted.Count, 2)
	}
	if converted.Next != resp.Next {
		t.Errorf("Next = %q, want %q", converted.Next, resp.Next)
	}
	if converted.Previous != resp.Previous {
		t.Errorf("Previous = %q, want %q", converted.Previous, resp.Previous)
	}
	if len(converted.Results) != 2 {
		t.Fatalf("len(Results) = %d, want %d", len(converted.Results), 2)
	}
	if converted.Results[0] != "Site A" {
		t.Errorf("Results[0] = %q, want %q", converted.Results[0], "Site A")
	}
	if converted.Results[1] != "Site B" {
		t.Errorf("Results[1] = %q, want %q", converted.Results[1], "Site B")
	}
}

// TestConvertPaginated_Empty tests convertPaginated with zero results.
func TestConvertPaginated_Empty(t *testing.T) {
	t.Parallel()

	convert := func(w WireSite) string {
		return w.Name
	}

	resp := &domain.PaginatedResponse[WireSite]{
		Count:    0,
		Next:     "",
		Previous: "",
		Results:  []WireSite{},
	}

	converted := convertPaginated(resp, convert)

	if converted.Count != 0 {
		t.Errorf("Count = %d, want %d", converted.Count, 0)
	}
	if len(converted.Results) != 0 {
		t.Errorf("len(Results) = %d, want %d", len(converted.Results), 0)
	}
}

func TestClient_ListInterfaces_Errors(t *testing.T) {
	t.Parallel()

	t.Run("not_found", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListInterfaces(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("invalid_json", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{invalid`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListInterfaces(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestClient_ListVMInterfaces_Errors(t *testing.T) {
	t.Parallel()

	t.Run("not_found", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListVMInterfaces(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("invalid_json", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{invalid`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListVMInterfaces(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestClient_ListCircuitTerminations_Errors(t *testing.T) {
	t.Parallel()

	t.Run("not_found", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListCircuitTerminations(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("invalid_json", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{invalid`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListCircuitTerminations(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestClient_ListCables_Errors(t *testing.T) {
	t.Parallel()

	t.Run("not_found", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListCables(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("invalid_json", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{invalid`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListCables(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

// contains is a helper to check substring presence.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsStr(s, substr)
}

// containsStr is a simple substring check without importing strings.
func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestClient_ListMethods_Errors covers the get() and json.Unmarshal error paths
// of the remaining list methods (Devices, IPAddresses, Prefixes, VLANs,
// VirtualMachines, Clusters, Circuits, Racks) that lack dedicated error tests.
func TestClient_ListMethods_Errors(t *testing.T) {
	t.Parallel()

	type methodCase struct {
		name string
		call func(*Client) error
		path string
	}

	cases := []methodCase{
		{"devices", func(c *Client) error { _, err := c.ListDevices(context.Background(), "t", nil); return err }, "/api/dcim/devices/"},
		{"ip_addresses", func(c *Client) error { _, err := c.ListIPAddresses(context.Background(), "t", nil); return err }, "/api/ipam/ip-addresses/"},
		{"prefixes", func(c *Client) error { _, err := c.ListPrefixes(context.Background(), "t", nil); return err }, "/api/ipam/prefixes/"},
		{"vlans", func(c *Client) error { _, err := c.ListVLANs(context.Background(), "t", nil); return err }, "/api/ipam/vlans/"},
		{"virtual_machines", func(c *Client) error { _, err := c.ListVirtualMachines(context.Background(), "t", nil); return err }, "/api/virtualization/virtual-machines/"},
		{"clusters", func(c *Client) error { _, err := c.ListClusters(context.Background(), "t", nil); return err }, "/api/virtualization/clusters/"},
		{"circuits", func(c *Client) error { _, err := c.ListCircuits(context.Background(), "t", nil); return err }, "/api/circuits/circuits/"},
		{"racks", func(c *Client) error { _, err := c.ListRacks(context.Background(), "t", nil); return err }, "/api/dcim/racks/"},
	}

	for _, tc := range cases {

		t.Run(tc.name+"/not_found", func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}))
			defer srv.Close()
			if err := tc.call(NewClient(srv.URL, http.DefaultClient)); err == nil {
				t.Error("expected error, got nil")
			}
		})
		t.Run(tc.name+"/invalid_json", func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{invalid`))
			}))
			defer srv.Close()
			if err := tc.call(NewClient(srv.URL, http.DefaultClient)); err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestConvertPaginated_Nil(t *testing.T) {
	t.Parallel()

	if got := convertPaginated[WireSite, domain.Site](nil, wireSiteToDomain); got != nil {
		t.Errorf("convertPaginated(nil) = %v, want nil", got)
	}
}

func TestWireCableTerminationToDomain_Populated(t *testing.T) {
	t.Parallel()

	w := &WireCableTermination{ID: 7, URL: "http://x/7/", Name: "A-end", Type: "dcim.interface"}
	got := wireCableTerminationToDomain(w)
	if got == nil {
		t.Fatal("wireCableTerminationToDomain(non-nil) = nil, want populated")
	}
	if got.ID != 7 || got.Name != "A-end" || got.Type != "dcim.interface" {
		t.Errorf("wireCableTerminationToDomain = %+v", got)
	}
}

func TestClient_ForwardsRequestID(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Request-ID") != "corr-123" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"count":0,"results":[]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	ctx := domain.WithRequestID(context.Background(), "corr-123")
	_, err := client.ListSites(ctx, "token", nil)
	if err != nil {
		t.Fatalf("ListSites() returned error: %v", err)
	}
}

func TestClient_NoRequestIDHeaderWhenAbsent(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Request-ID"); got != "" {
			t.Errorf("X-Request-ID = %q, want empty", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"count":0,"results":[]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	_, err := client.ListSites(context.Background(), "token", nil)
	if err != nil {
		t.Fatalf("ListSites() returned error: %v", err)
	}
}
