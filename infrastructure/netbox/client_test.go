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

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_ListSites(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "Bearer test-token" {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}
			if req.URL.Path != "/api/dcim/sites/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
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
			t.Errorf("Count = %domainResult, want %domainResult", resp.Count, 1)
		}
		if len(resp.Results) != 1 {
			t.Fatalf("len(Results) = %domainResult, want %domainResult", len(resp.Results), 1)
		}
		if resp.Results[0].Name != "Test Site" {
			t.Errorf("Name = %q, want %q", resp.Results[0].Name, "Test Site")
		}
	})

	t.Run("unauthorized", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusUnauthorized)
			_, _ = writer.Write([]byte(`{"detail": "Invalid token"}`))
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
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
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
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{invalid json`))
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
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.URL.Query().Get("site") != "dc1" {
				writer.WriteHeader(http.StatusBadRequest)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
				"count": 2,
				"results": [
					{
						"id": 1,
						"name": "router1",
						"site": {"id": 1, "name": "DC1", "url": "", "slug": "dc1"},
						"created": "",
						"last_updated": ""
					},
					{
						"id": 2,
						"name": "switch1",
						"site": {"id": 1, "name": "DC1", "url": "", "slug": "dc1"},
						"created": "",
						"last_updated": ""
					}
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
			t.Errorf("Count = %domainResult, want %domainResult", resp.Count, 2)
		}
	})
}

func TestClient_ListIPAddresses(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
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
			t.Fatalf("len(Results) = %domainResult, want %domainResult", len(resp.Results), 1)
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
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
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
			t.Fatalf("len(Results) = %domainResult, want %domainResult", len(resp.Results), 1)
		}
		if resp.Results[0].Prefix != "10.0.0.0/8" {
			t.Errorf("Prefix = %q, want %q", resp.Results[0].Prefix, "10.0.0.0/8")
		}
	})
}

func TestClient_ListVLANs(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
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
			t.Fatalf("len(Results) = %domainResult, want %domainResult", len(resp.Results), 1)
		}
		if resp.Results[0].Name != "Users" {
			t.Errorf("Name = %q, want %q", resp.Results[0].Name, "Users")
		}
	})
}

func TestClient_ListVirtualMachines(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
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
			t.Fatalf("len(Results) = %domainResult, want %domainResult", len(resp.Results), 1)
		}
		if resp.Results[0].Name != "web-01" {
			t.Errorf("Name = %q, want %q", resp.Results[0].Name, "web-01")
		}
	})
}

func TestClient_ListClusters(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
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
			t.Fatalf("len(Results) = %domainResult, want %domainResult", len(resp.Results), 1)
		}
		if resp.Results[0].Name != "Prod Cluster" {
			t.Errorf("Name = %q, want %q", resp.Results[0].Name, "Prod Cluster")
		}
	})
}

func TestClient_ListCircuits(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"cid": "CIR-001",
					"provider": {"id": 1, "name": "ATT", "url": "", "slug": "att"},
					"type": {"id": 1, "name": "Dark Fiber", "url": "", "slug": "dark-fiber"},
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
			t.Fatalf("len(Results) = %domainResult, want %domainResult", len(resp.Results), 1)
		}
		if resp.Results[0].CID != "CIR-001" {
			t.Errorf("CID = %q, want %q", resp.Results[0].CID, "CIR-001")
		}
	})
}

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_ListInterfaces(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/dcim/interfaces/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
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
			t.Errorf("Count = %domainResult, want %domainResult", resp.Count, 2)
		}
		if len(resp.Results) != 2 {
			t.Fatalf("len(Results) = %domainResult, want %domainResult", len(resp.Results), 2)
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

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_ListVMInterfaces(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/virtualization/interfaces/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
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
			t.Errorf("Count = %domainResult, want %domainResult", resp.Count, 2)
		}
		if len(resp.Results) != 2 {
			t.Fatalf("len(Results) = %domainResult, want %domainResult", len(resp.Results), 2)
		}
		if resp.Results[0].Name != "eth0" {
			t.Errorf("Results[0].Name = %q, want %q", resp.Results[0].Name, "eth0")
		}
		if resp.Results[0].VirtualMachine.Name != "web-01" {
			t.Errorf("Results[0].VirtualMachine.Name = %q, want %q", resp.Results[0].VirtualMachine.Name, "web-01")
		}
		if resp.Results[0].MTU != 1500 {
			t.Errorf("Results[0].MTU = %domainResult, want %domainResult", resp.Results[0].MTU, 1500)
		}
	})
}

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_ListCircuitTerminations(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/circuits/circuit-terminations/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
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
			t.Errorf("Count = %domainResult, want %domainResult", resp.Count, 2)
		}
		if len(resp.Results) != 2 {
			t.Fatalf("len(Results) = %domainResult, want %domainResult", len(resp.Results), 2)
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

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_ListCables(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/dcim/cables/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
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
			t.Errorf("Count = %domainResult, want %domainResult", resp.Count, 2)
		}
		if len(resp.Results) != 2 {
			t.Fatalf("len(Results) = %domainResult, want %domainResult", len(resp.Results), 2)
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
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
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
			t.Fatalf("len(Results) = %domainResult, want %domainResult", len(resp.Results), 1)
		}
		if resp.Results[0].Name != "Rack-01" {
			t.Errorf("Name = %q, want %q", resp.Results[0].Name, "Rack-01")
		}
	})
}

func TestClient_GetObject(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{"id": 1, "name": "Test", "slug": "test"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.GetObject(context.Background(), "token", "site", 1, nil)
		if err != nil {
			t.Fatalf("GetObject() returned error: %v", err)
		}
		var result map[string]any
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
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
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
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusForbidden)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListSites(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("server error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListSites(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("rate limited (429) with Retry-After", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Retry-After", "30")
			writer.WriteHeader(http.StatusTooManyRequests)
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
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		sawAccept = req.Header.Get("Accept") == "application/json"
		sawAuth = req.Header.Get("Authorization") == "Bearer tok"
		sawQuery = req.URL.Query().Get("name") == "site-a" && req.URL.Query().Get("status") == "active"
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"count":0,"results":[]}`))
	}))
	defer srv.Close()

	rc := NewClient(srv.URL, http.DefaultClient)
	if _, ok := any(rc.client).(*resty.Client); !ok {
		t.Fatalf("client is %T, want *resty.Client", rc.client)
	}

	if _, err := rc.ListSites(context.Background(), "tok", map[string]string{
		"name":   "site-a",
		"status": "active",
	}); err != nil {
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

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"count": 0, "results": []}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	resp, err := client.ListSites(context.Background(), "token", nil)
	if err != nil {
		t.Fatalf("ListSites() returned error: %v", err)
	}
	if len(resp.Results) != 0 {
		t.Errorf("len(Results) = %domainResult, want 0", len(resp.Results))
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

	for _, testCase := range cases {
		t.Run(testCase.objectType, func(t *testing.T) {
			endpoint, ok := objectTypeToEndpoint[testCase.objectType]
			if !ok {
				t.Fatalf("object type %q not found in map", testCase.objectType)
			}
			if endpoint != testCase.expected {
				t.Errorf("endpoint = %q, want %q", endpoint, testCase.expected)
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

	wire := WireInterface{
		ID:      1,
		Name:    "ge-0/0/0",
		Device:  &WireNested{ID: 1, Name: "router1", Slug: "router1", URL: ""},
		Type:    &WireLabel{Value: "1000base-t", Label: "1000BASE-T"},
		Enabled: true,
		MTU:     1500,
	}
	domainResult := wireInterfaceToDomain(wire)
	if domainResult.Name != "ge-0/0/0" {
		t.Errorf("Name = %q, want %q", domainResult.Name, "ge-0/0/0")
	}
	if !domainResult.Enabled {
		t.Error("Enabled = false, want true")
	}
	if domainResult.Device.Name != "router1" {
		t.Errorf("Device.Name = %q, want %q", domainResult.Device.Name, "router1")
	}
}

func TestWireVMInterfaceToDomain(t *testing.T) {
	t.Parallel()

	wire := WireVMInterface{
		ID:             1,
		Name:           "eth0",
		VirtualMachine: &WireNested{ID: 1, Name: "web-01", Slug: "web-01", URL: ""},
		Enabled:        true,
	}
	domainResult := wireVMInterfaceToDomain(wire)
	if domainResult.Name != "eth0" {
		t.Errorf("Name = %q, want %q", domainResult.Name, "eth0")
	}
	if domainResult.VirtualMachine.Name != "web-01" {
		t.Errorf("VirtualMachine.Name = %q, want %q", domainResult.VirtualMachine.Name, "web-01")
	}
}

func TestWireCircuitTerminationToDomain(t *testing.T) {
	t.Parallel()

	wire := WireCircuitTermination{
		ID:       1,
		Circuit:  &WireNested{ID: 1, Name: "CIR-001", Slug: "cir-001", URL: ""},
		TermSide: "A",
		Site:     &WireNested{ID: 1, Name: "DC1", Slug: "dc1", URL: ""},
		Speed:    10000000,
	}
	domainResult := wireCircuitTerminationToDomain(wire)
	if domainResult.TermSide != "A" {
		t.Errorf("TermSide = %q, want %q", domainResult.TermSide, "A")
	}
	if domainResult.Circuit.Name != "CIR-001" {
		t.Errorf("Circuit.Name = %q, want %q", domainResult.Circuit.Name, "CIR-001")
	}
}

func TestWireCableToDomain(t *testing.T) {
	t.Parallel()

	wire := WireCable{
		ID:     1,
		Type:   &WireLabel{Value: "cat6a", Label: "CAT6a"},
		Status: &WireLabel{Value: "connected", Label: "Connected"},
		Label:  "link-01",
		Color:  "blue",
		Length: 10.5,
	}
	domainResult := wireCableToDomain(wire)
	if domainResult.Label != "link-01" {
		t.Errorf("Label = %q, want %q", domainResult.Label, "link-01")
	}
	if domainResult.Type.Value != "cat6a" {
		t.Errorf("Type.Value = %q, want %q", domainResult.Type.Value, "cat6a")
	}
	if domainResult.Status.Value != "connected" {
		t.Errorf("Status.Value = %q, want %q", domainResult.Status.Value, "connected")
	}
}

func TestWireCableToDomain_NilTypeAndStatus(t *testing.T) {
	t.Parallel()

	wire := WireCable{
		ID:     1,
		Type:   nil,
		Status: nil,
		Tags:   nil,
	}
	domainResult := wireCableToDomain(wire)
	if domainResult.Type != nil {
		t.Errorf("Type = %+v, want nil", domainResult.Type)
	}
	if domainResult.Status != nil {
		t.Errorf("Status = %+v, want nil", domainResult.Status)
	}
}

func TestWireVLANToDomain(t *testing.T) {
	t.Parallel()

	wire := WireVLAN{
		ID:   1,
		VID:  100,
		Name: "Prod-VLAN",
		Status: &WireLabel{
			Value: "active",
			Label: "Active",
		},
		Site: &WireNested{ID: 1, Name: "DC-1", Slug: "dc-1", URL: ""},
	}
	domainResult := wireVLANToDomain(wire)
	if domainResult.VID != 100 {
		t.Errorf("VID = %domainResult, want %domainResult", domainResult.VID, 100)
	}
	if domainResult.Name != "Prod-VLAN" {
		t.Errorf("Name = %q, want %q", domainResult.Name, "Prod-VLAN")
	}
	if domainResult.Status.Value != "active" {
		t.Errorf("Status.Value = %q, want %q", domainResult.Status.Value, "active")
	}
	if domainResult.Site.Name != "DC-1" {
		t.Errorf("Site.Name = %q, want %q", domainResult.Site.Name, "DC-1")
	}
}

func TestWireVirtualMachineToDomain(t *testing.T) {
	t.Parallel()

	wire := WireVirtualMachine{
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
	domainResult := wireVirtualMachineToDomain(wire)
	if domainResult.Name != "web-01" {
		t.Errorf("Name = %q, want %q", domainResult.Name, "web-01")
	}
	if domainResult.VCPUs != 4 {
		t.Errorf("VCPUs = %f, want %f", domainResult.VCPUs, 4.0)
	}
	if domainResult.Memory != 8192 {
		t.Errorf("Memory = %domainResult, want %domainResult", domainResult.Memory, 8192)
	}
	if domainResult.Disk != 100 {
		t.Errorf("Disk = %domainResult, want %domainResult", domainResult.Disk, 100)
	}
	if domainResult.Status.Value != "active" {
		t.Errorf("Status.Value = %q, want %q", domainResult.Status.Value, "active")
	}
}

func TestWireClusterToDomain(t *testing.T) {
	t.Parallel()

	wire := WireCluster{
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
	domainResult := wireClusterToDomain(wire)
	if domainResult.Name != "K8s-Prod" {
		t.Errorf("Name = %q, want %q", domainResult.Name, "K8s-Prod")
	}
	if domainResult.ClusterType.Name != "Kubernetes" {
		t.Errorf("ClusterType.Name = %q, want %q", domainResult.ClusterType.Name, "Kubernetes")
	}
	if domainResult.ClusterGroup.Name != "Production" {
		t.Errorf("ClusterGroup.Name = %q, want %q", domainResult.ClusterGroup.Name, "Production")
	}
	if domainResult.Site.Name != "DC-1" {
		t.Errorf("Site.Name = %q, want %q", domainResult.Site.Name, "DC-1")
	}
}

func TestWireCircuitToDomain(t *testing.T) {
	t.Parallel()

	wire := WireCircuit{
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
	domainResult := wireCircuitToDomain(wire)
	if domainResult.CID != "CIR-001" {
		t.Errorf("CID = %q, want %q", domainResult.CID, "CIR-001")
	}
	if domainResult.Provider.Name != "Equinix" {
		t.Errorf("Provider.Name = %q, want %q", domainResult.Provider.Name, "Equinix")
	}
	if domainResult.CircuitType.Name != "Dark Fiber" {
		t.Errorf("CircuitType.Name = %q, want %q", domainResult.CircuitType.Name, "Dark Fiber")
	}
	if domainResult.Status.Value != "active" {
		t.Errorf("Status.Value = %q, want %q", domainResult.Status.Value, "active")
	}
}

func TestWireRackToDomain(t *testing.T) {
	t.Parallel()

	wire := WireRack{
		ID:         1,
		Name:       "Rack-101",
		FacilityID: "A1-01",
		Site:       &WireNested{ID: 1, Name: "DC-1", Slug: "dc-1", URL: ""},
		Width:      19,
		UHeight:    42,
	}
	domainResult := wireRackToDomain(wire)
	if domainResult.Name != "Rack-101" {
		t.Errorf("Name = %q, want %q", domainResult.Name, "Rack-101")
	}
	if domainResult.FacilityID != "A1-01" {
		t.Errorf("FacilityID = %q, want %q", domainResult.FacilityID, "A1-01")
	}
	if domainResult.Site.Name != "DC-1" {
		t.Errorf("Site.Name = %q, want %q", domainResult.Site.Name, "DC-1")
	}
	if domainResult.Width != 19 {
		t.Errorf("Width = %domainResult, want %domainResult", domainResult.Width, 19)
	}
	if domainResult.UHeight != 42 {
		t.Errorf("UHeight = %domainResult, want %domainResult", domainResult.UHeight, 42)
	}
}

func TestConversionFunctions_Populated(t *testing.T) {
	t.Parallel()

	t.Run("wireNestedToDomain", func(t *testing.T) {
		wire := &WireNested{ID: 1, Name: "Test", Slug: "test", URL: "http://example.com"}
		domainResult := wireNestedToDomain(wire)
		if domainResult.ID != 1 || domainResult.Name != "Test" || domainResult.Slug != "test" {
			t.Errorf("unexpected result: %+v", domainResult)
		}
	})

	t.Run("wireLabelToDomain", func(t *testing.T) {
		wire := &WireLabel{Value: "active", Label: "Active"}
		domainResult := wireLabelToDomain(wire)
		if domainResult.Value != "active" || domainResult.Label != "Active" {
			t.Errorf("unexpected result: %+v", domainResult)
		}
	})

	t.Run("wireTagsToDomain", func(t *testing.T) {
		wire := []WireTag{{ID: 1, Name: "prod", Slug: "prod"}}
		domainResult := wireTagsToDomain(wire)
		if len(domainResult) != 1 || domainResult[0].Name != "prod" {
			t.Errorf("unexpected result: %+v", domainResult)
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
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{invalid`))
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

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"count": 0, "results": [{"id": 1, "prefix": "10.0.0.0/8", "family": null, "created": "` +
			`", "last_updated": ""}]}`))
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

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{
			"count": 1,
			"results": [{
				"id": 1,
				"address": "10.0.0.1/24",
				"status": {"value": "active", "label": "Active"},
				"assigned_object_type": "dcim.interface",
				"assigned_object_id": 100,
				"assigned_object": {
					"id": 100,
					"name": "ge-0/0/0",
					"device": {"id": 1, "name": "router1", "url": "", "slug": "router1"},
					"url": ""
				},
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
		t.Fatalf("len(Results) = %domainResult, want %domainResult", len(resp.Results), 1)
	}
	ipAddress := resp.Results[0]
	if ipAddress.AssignedObject == nil {
		t.Fatal("AssignedObject is nil")
	}
	if ipAddress.AssignedObject.Name != "ge-0/0/0" {
		t.Errorf("AssignedObject.Name = %q, want %q", ipAddress.AssignedObject.Name, "ge-0/0/0")
	}
	if ipAddress.AssignedObject.Device.Name != "router1" {
		t.Errorf("Device.Name = %q, want %q", ipAddress.AssignedObject.Device.Name, "router1")
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

	var wire WireSite
	if err := json.Unmarshal([]byte(raw), &wire); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	domainResult := wireSiteToDomain(wire)

	if domainResult.ID != 1 || domainResult.Name != "Test" || domainResult.Slug != "test" {
		t.Errorf("basic fields mismatch: %+v", domainResult)
	}
	if domainResult.Region.Name != "US" {
		t.Errorf("Region.Name = %q, want %q", domainResult.Region.Name, "US")
	}
	if domainResult.Tenant.Name != "Acme" {
		t.Errorf("Tenant.Name = %q, want %q", domainResult.Tenant.Name, "Acme")
	}
	if domainResult.Tenant.Name != "Acme" {
		t.Errorf("Tenant.Name = %q, want %q", domainResult.Tenant.Name, "Acme")
	}
	if len(domainResult.Tags) != 1 || domainResult.Tags[0].Name != "prod" {
		t.Errorf("tags mismatch: %+v", domainResult.Tags)
	}
}

// --- doRequest direct tests ---

func TestDoRequest_Success(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Accept") != "application/json" {
			writer.WriteHeader(http.StatusBadRequest)

			return
		}
		if req.Header.Get("Authorization") != "Bearer test-token" {
			writer.WriteHeader(http.StatusUnauthorized)

			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"name":"test"}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	body, err := client.doRequest(context.Background(), "test-token", http.MethodGet, "/api/dcim/sites/", nil, nil)
	if err != nil {
		t.Fatalf("doRequest() returned error: %v", err)
	}
	if string(body) != `{"name":"test"}` {
		t.Errorf("body = %q, want %q", string(body), `{"name":"test"}`)
	}
}

func TestDoRequest_Unauthorized(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	_, err := client.doRequest(context.Background(), "bad-token", http.MethodGet, "/api/dcim/sites/", nil, nil)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !contains(err.Error(), "unauthorized") {
		t.Errorf("error = %q, want it to contain 'unauthorized'", err.Error())
	}
}

func TestDoRequest_Forbidden(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	_, err := client.doRequest(context.Background(), "token", http.MethodGet, "/api/dcim/sites/", nil, nil)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !contains(err.Error(), "forbidden") {
		t.Errorf("error = %q, want it to contain 'forbidden'", err.Error())
	}
}

func TestDoRequest_NotFound(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	_, err := client.doRequest(context.Background(), "token", http.MethodGet, "/api/dcim/sites/", nil, nil)
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
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusTooManyRequests)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.doRequest(context.Background(), "token", http.MethodGet, "/api/dcim/sites/", nil, nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !contains(err.Error(), "rate limited") {
			t.Errorf("error = %q, want it to contain 'rate limited'", err.Error())
		}
	})

	t.Run("with Retry-After", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Retry-After", "5")
			writer.WriteHeader(http.StatusTooManyRequests)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.doRequest(context.Background(), "token", http.MethodGet, "/api/dcim/sites/", nil, nil)
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

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(`{"detail":"internal error"}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	_, err := client.doRequest(context.Background(), "token", http.MethodGet, "/api/dcim/sites/", nil, nil)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !contains(err.Error(), "unexpected status 500") {
		t.Errorf("error = %q, want it to contain 'unexpected status 500'", err.Error())
	}
}

func TestDoRequest_Status300IsError(t *testing.T) {
	t.Parallel()

	// A 300 status is outside the 2xx range and must be treated as an error.
	// Kills the CONDITIONALS_BOUNDARY on `statusCode >= 300` (whose mutant would
	// let exactly 300 through as a success).
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusMultipleChoices)
		_, _ = writer.Write([]byte(`{"detail":"multiple choices"}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	_, err := client.doRequest(context.Background(), "token", http.MethodGet, "/api/dcim/sites/", nil, nil)
	if err == nil {
		t.Fatal("Expected error for 300, got nil")
	}
	if !contains(err.Error(), "unexpected status 300") {
		t.Errorf("error = %q, want it to contain 'unexpected status 300'", err.Error())
	}
}

func TestDoRequest_QueryParams(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if req.URL.Query().Get("name") != "test-site" {
			writer.WriteHeader(http.StatusBadRequest)

			return
		}
		if req.URL.Query().Get("status") != "active" {
			writer.WriteHeader(http.StatusBadRequest)

			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"result":"ok"}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	body, err := client.doRequest(context.Background(), "token", http.MethodGet, "/api/dcim/sites/", map[string]string{
		"name":   "test-site",
		"status": "active",
	}, nil)
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

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write(largeBody)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	_, err := client.doRequest(context.Background(), "token", http.MethodGet, "/api/dcim/sites/", nil, nil)
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
	convert := func(wire WireSite) string {
		return wire.Name
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
		t.Errorf("Count = %domainResult, want %domainResult", converted.Count, 2)
	}
	if converted.Next != resp.Next {
		t.Errorf("Next = %q, want %q", converted.Next, resp.Next)
	}
	if converted.Previous != resp.Previous {
		t.Errorf("Previous = %q, want %q", converted.Previous, resp.Previous)
	}
	if len(converted.Results) != 2 {
		t.Fatalf("len(Results) = %domainResult, want %domainResult", len(converted.Results), 2)
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

	convert := func(wire WireSite) string {
		return wire.Name
	}

	resp := &domain.PaginatedResponse[WireSite]{
		Count:    0,
		Next:     "",
		Previous: "",
		Results:  []WireSite{},
	}

	converted := convertPaginated(resp, convert)

	if converted.Count != 0 {
		t.Errorf("Count = %domainResult, want %domainResult", converted.Count, 0)
	}
	if len(converted.Results) != 0 {
		t.Errorf("len(Results) = %domainResult, want %domainResult", len(converted.Results), 0)
	}
}

func TestClient_ListInterfaces_Errors(t *testing.T) {
	t.Parallel()

	t.Run("not_found", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListInterfaces(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("invalid_json", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{invalid`))
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
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListVMInterfaces(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("invalid_json", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{invalid`))
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
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListCircuitTerminations(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("invalid_json", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{invalid`))
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
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListCables(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("invalid_json", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{invalid`))
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
//
//nolint:funlen // table-driven test over remaining list methods; splitting is not sensible.
func TestClient_ListMethods_Errors(t *testing.T) {
	t.Parallel()

	type methodCase struct {
		name string
		call func(*Client) error
		path string
	}

	cases := []methodCase{
		{
			"devices",
			func(cable *Client) error {
				_, err := cable.ListDevices(context.Background(), "t", nil)

				return err
			},
			"/api/dcim/devices/",
		},
		{
			"ip_addresses",
			func(cable *Client) error {
				_, err := cable.ListIPAddresses(context.Background(), "t", nil)

				return err
			},
			"/api/ipam/ip-addresses/",
		},
		{
			"prefixes",
			func(cable *Client) error {
				_, err := cable.ListPrefixes(context.Background(), "t", nil)

				return err
			},
			"/api/ipam/prefixes/",
		},
		{
			"vlans",
			func(cable *Client) error {
				_, err := cable.ListVLANs(context.Background(), "t", nil)

				return err
			},
			"/api/ipam/vlans/",
		},
		{"virtual_machines", func(cable *Client) error {
			_, err := cable.ListVirtualMachines(context.Background(), "t", nil)

			return err
		}, "/api/virtualization/virtual-machines/"},
		{
			"clusters",
			func(cable *Client) error {
				_, err := cable.ListClusters(context.Background(), "t", nil)

				return err
			},
			"/api/virtualization/clusters/",
		},
		{
			"circuits",
			func(cable *Client) error {
				_, err := cable.ListCircuits(context.Background(), "t", nil)

				return err
			},
			"/api/circuits/circuits/",
		},
		{
			"racks",
			func(cable *Client) error {
				_, err := cable.ListRacks(context.Background(), "t", nil)

				return err
			},
			"/api/dcim/racks/",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name+"/not_found", func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
				writer.WriteHeader(http.StatusNotFound)
			}))
			defer srv.Close()
			if err := testCase.call(NewClient(srv.URL, http.DefaultClient)); err == nil {
				t.Error("expected error, got nil")
			}
		})
		t.Run(testCase.name+"/invalid_json", func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
				writer.Header().Set("Content-Type", "application/json")
				_, _ = writer.Write([]byte(`{invalid`))
			}))
			defer srv.Close()
			if err := testCase.call(NewClient(srv.URL, http.DefaultClient)); err == nil {
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

	wire := &WireCableTermination{ID: 7, URL: "http://x/7/", Name: "A-end", Type: "dcim.interface"}
	got := wireCableTerminationToDomain(wire)
	if got == nil {
		t.Fatal("wireCableTerminationToDomain(non-nil) = nil, want populated")
	}
	if got.ID != 7 || got.Name != "A-end" || got.Type != "dcim.interface" {
		t.Errorf("wireCableTerminationToDomain = %+v", got)
	}
}

func TestClient_ForwardsRequestID(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-Request-ID") != "corr-123" {
			writer.WriteHeader(http.StatusBadRequest)

			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"count":0,"results":[]}`))
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

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if got := req.Header.Get("X-Request-ID"); got != "" {
			t.Errorf("X-Request-ID = %q, want empty", got)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"count":0,"results":[]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	_, err := client.ListSites(context.Background(), "token", nil)
	if err != nil {
		t.Fatalf("ListSites() returned error: %v", err)
	}
}
