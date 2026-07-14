package netbox

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

// TestDoRequest_WithBodyLimit verifies that doRequest enforces a 10 MB body limit.
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
	body, err := client.doRequest(context.Background(), "token", http.MethodGet, "/api/dcim/sites/", nil)
	if err != nil {
		t.Fatalf("doRequest() returned error: %v", err)
	}
	if len(body) > 10*1024*1024 {
		t.Errorf("body length = %d, expected at most %d (10 MB limit)", len(body), 10*1024*1024)
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
