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

// Client-level CRUD tests for the 8 DCIM entities added in the
// feature/expand-entities branch: site_group, rack_group, rack_role,
// device_role, platform, virtual_chassis, power_panel, inventory_item_role.
//
// Each entity is exercised over its full REST lifecycle (List / Create /
// Update / Delete) with the shared write helpers (WriteToWire + WireToDomain),
// following the tenant_group_test.go pattern. These tests cover the
// previously-untested client methods so the package stays above the 95%
// coverage gate.

// --- site_group -----------------------------------------------------------

func TestClient_ListSiteGroups(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/dcim/site-groups/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{"count":1,"results":[{"id":1,"name":"EU-West","slug":"eu-west",
				"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}]}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListSiteGroups(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListSiteGroups() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "EU-West" {
			t.Errorf("resp = %+v, want count 1 with EU-West", resp)
		}
	})
}

func TestClient_CreateSiteGroup(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath, gotCT, gotAuth string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			gotCT = req.Header.Get("Content-Type")
			gotAuth = req.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(req.Body)
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusCreated)
			_, _ = writer.Write([]byte(`{"id":1,"name":"Group A","slug":"group-a","created":"2024-01-01",` +
				`"last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		siteGroup, err := client.CreateSiteGroup(context.Background(), "tok", domain.SiteGroupWrite{Name: "Group A"})
		if err != nil {
			t.Fatalf("CreateSiteGroup() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/dcim/site-groups/" {
			t.Errorf("method/path = %q %q, want POST /api/dcim/site-groups/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]string
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "Group A" {
			t.Errorf("request body = %s, want name Group A", gotBody)
		}
		if siteGroup.ID != 1 || siteGroup.Name != "Group A" {
			t.Errorf("site group = %+v, want id 1 name Group A", siteGroup)
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
		_, err := client.CreateSiteGroup(context.Background(), "tok", domain.SiteGroupWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateSiteGroup(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"name":"Renamed","slug":"renamed","created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		siteGroup, err := client.UpdateSiteGroup(context.Background(), "tok", 7, domain.SiteGroupWrite{Name: "Renamed"})
		if err != nil {
			t.Fatalf("UpdateSiteGroup() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/dcim/site-groups/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/dcim/site-groups/7/", gotMethod, gotPath)
		}
		if siteGroup.ID != 7 || siteGroup.Name != "Renamed" {
			t.Errorf("site group = %+v, want id 7 name Renamed", siteGroup)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"slug":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateSiteGroup(context.Background(), "tok", 7, domain.SiteGroupWrite{Name: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteSiteGroup(t *testing.T) {
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
		if err := client.DeleteSiteGroup(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteSiteGroup() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/dcim/site-groups/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/dcim/site-groups/3/", gotMethod, gotPath)
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
		err := client.DeleteSiteGroup(context.Background(), "tok", 3)
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
		if err := client.DeleteSiteGroup(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteSiteGroup() 404 returned error: %v", err)
		}
	})
}

func TestSiteGroupWriteToWire(t *testing.T) {
	t.Parallel()

	parent := 2
	desc := "core"
	in := domain.SiteGroupWrite{Name: "A", Parent: &parent, Description: &desc, Tags: []string{"prod"}}
	wire := siteGroupWriteToWire(in)
	if wire.Name != "A" || wire.Parent == nil || *wire.Parent != 2 {
		t.Errorf("wire = %+v, want name A parent 2", wire)
	}
	if wire.Description == nil || *wire.Description != "core" {
		t.Errorf("wire.Description = %v, want core", wire.Description)
	}
}

func TestWireSiteGroupToDomain(t *testing.T) {
	t.Parallel()

	wire := WireSiteGroup{
		ID:      1,
		URL:     "http://nb/api/dcim/site-groups/1/",
		Name:    "EU-West",
		Slug:    "eu-west",
		Parent:  &WireNested{ID: 9, Name: "EMEA", Slug: "emea"},
		Created: "2024-01-01",
	}
	got := wireSiteGroupToDomain(wire)
	if got.ID != 1 || got.Name != "EU-West" || got.Slug != "eu-west" {
		t.Errorf("domain = %+v, want id 1 name EU-West slug eu-west", got)
	}
	if got.Parent == nil || got.Parent.ID != 9 || got.Parent.Name != "EMEA" {
		t.Errorf("domain.Parent = %+v, want nested id 9 name EMEA", got.Parent)
	}

	if got := wireSiteGroupToDomain(WireSiteGroup{}).Parent; got != nil {
		t.Errorf("empty wire Parent = %+v, want nil", got)
	}
}

// --- rack_group -----------------------------------------------------------

func TestClient_ListRackGroups(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/dcim/rack-groups/" {
			writer.WriteHeader(http.StatusNotFound)

			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"count":1,"results":[{"id":1,"name":"RowA","slug":"row-a",
			"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	resp, err := client.ListRackGroups(context.Background(), "test-token", nil)
	if err != nil {
		t.Fatalf("ListRackGroups() returned error: %v", err)
	}
	if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "RowA" {
		t.Errorf("resp = %+v, want count 1 with RowA", resp)
	}
}

func TestClient_CreateRackGroup(t *testing.T) {
	t.Parallel()
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		_, _ = writer.Write([]byte(`{"id":1,"name":"RowA","slug":"row-a","created":"2024-01-01",` +
			`"last_updated":"2024-01-01T00:00:00Z"}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	rackGroup, err := client.CreateRackGroup(context.Background(), "tok", domain.RackGroupWrite{Name: "RowA"})
	if err != nil {
		t.Fatalf("CreateRackGroup() returned error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/api/dcim/rack-groups/" {
		t.Errorf("method/path = %q %q, want POST /api/dcim/rack-groups/", gotMethod, gotPath)
	}
	if rackGroup.ID != 1 || rackGroup.Name != "RowA" {
		t.Errorf("rack group = %+v, want id 1 name RowA", rackGroup)
	}
}

func TestClient_UpdateRackGroup(t *testing.T) {
	t.Parallel()
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(`{"id":7,"name":"Renamed","slug":"renamed","created":"","last_updated":""}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	rackGroup, err := client.UpdateRackGroup(context.Background(), "tok", 7, domain.RackGroupWrite{Name: "Renamed"})
	if err != nil {
		t.Fatalf("UpdateRackGroup() returned error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/api/dcim/rack-groups/7/" {
		t.Errorf("method/path = %q %q, want PATCH /api/dcim/rack-groups/7/", gotMethod, gotPath)
	}
	if rackGroup.ID != 7 || rackGroup.Name != "Renamed" {
		t.Errorf("rack group = %+v, want id 7 name Renamed", rackGroup)
	}
}

func TestClient_DeleteRackGroup(t *testing.T) {
	t.Parallel()
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		writer.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	if err := client.DeleteRackGroup(context.Background(), "tok", 3); err != nil {
		t.Fatalf("DeleteRackGroup() returned error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/api/dcim/rack-groups/3/" {
		t.Errorf("method/path = %q %q, want DELETE /api/dcim/rack-groups/3/", gotMethod, gotPath)
	}
}

func TestRackGroupWriteToWire(t *testing.T) {
	t.Parallel()

	site := 5
	parent := 2
	in := domain.RackGroupWrite{Name: "A", Site: &site, Parent: &parent}
	wire := rackGroupWriteToWire(in)
	if wire.Name != "A" || wire.Site == nil || *wire.Site != 5 || wire.Parent == nil || *wire.Parent != 2 {
		t.Errorf("wire = %+v, want name A site 5 parent 2", wire)
	}
}

func TestWireRackGroupToDomain(t *testing.T) {
	t.Parallel()

	wire := WireRackGroup{
		ID:      1,
		URL:     "http://nb/api/dcim/rack-groups/1/",
		Name:    "RowA",
		Slug:    "row-a",
		Site:    &WireNested{ID: 5, Name: "DC1", Slug: "dc1"},
		Parent:  &WireNested{ID: 2, Name: "EMEA", Slug: "emea"},
		Created: "2024-01-01",
	}
	got := wireRackGroupToDomain(wire)
	if got.ID != 1 || got.Name != "RowA" {
		t.Errorf("domain = %+v, want id 1 name RowA", got)
	}
	if got.Site == nil || got.Site.ID != 5 {
		t.Errorf("domain.Site = %+v, want nested id 5", got.Site)
	}
	if got.Parent == nil || got.Parent.ID != 2 {
		t.Errorf("domain.Parent = %+v, want nested id 2", got.Parent)
	}

	if got := wireRackGroupToDomain(WireRackGroup{}).Site; got != nil {
		t.Errorf("empty wire Site = %+v, want nil", got)
	}
}

// --- rack_role ------------------------------------------------------------

func TestClient_ListRackRoles(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/dcim/rack-roles/" {
			writer.WriteHeader(http.StatusNotFound)

			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"count":1,"results":[{"id":1,"name":"Core","slug":"core","color":"00ff00",
			"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	resp, err := client.ListRackRoles(context.Background(), "test-token", nil)
	if err != nil {
		t.Fatalf("ListRackRoles() returned error: %v", err)
	}
	if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "Core" || resp.Results[0].Color != "00ff00" {
		t.Errorf("resp = %+v, want count 1 with Core/00ff00", resp)
	}
}

func TestClient_CreateRackRole(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		_, _ = writer.Write([]byte(`{"id":1,"name":"Core","slug":"core","color":"00ff00",` +
			`"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	rackRole, err := client.CreateRackRole(context.Background(), "tok", domain.RackRoleWrite{Name: "Core"})
	if err != nil {
		t.Fatalf("CreateRackRole() returned error: %v", err)
	}
	if rackRole.ID != 1 || rackRole.Name != "Core" || rackRole.Color != "00ff00" {
		t.Errorf("rack role = %+v, want id 1 name Core color 00ff00", rackRole)
	}
}

func TestClient_UpdateRackRole(t *testing.T) {
	t.Parallel()
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(`{"id":7,"name":"Renamed","slug":"renamed","color":"ff0000",` +
			`"created":"","last_updated":""}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	rackRole, err := client.UpdateRackRole(context.Background(), "tok", 7, domain.RackRoleWrite{Name: "Renamed"})
	if err != nil {
		t.Fatalf("UpdateRackRole() returned error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/api/dcim/rack-roles/7/" {
		t.Errorf("method/path = %q %q, want PATCH /api/dcim/rack-roles/7/", gotMethod, gotPath)
	}
	if rackRole.ID != 7 || rackRole.Name != "Renamed" {
		t.Errorf("rack role = %+v, want id 7 name Renamed", rackRole)
	}
}

func TestClient_DeleteRackRole(t *testing.T) {
	t.Parallel()
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		writer.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	if err := client.DeleteRackRole(context.Background(), "tok", 3); err != nil {
		t.Fatalf("DeleteRackRole() returned error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/api/dcim/rack-roles/3/" {
		t.Errorf("method/path = %q %q, want DELETE /api/dcim/rack-roles/3/", gotMethod, gotPath)
	}
}

func TestRackRoleWriteToWire(t *testing.T) {
	t.Parallel()

	color := "00ff00"
	in := domain.RackRoleWrite{Name: "A", Color: &color}
	wire := rackRoleWriteToWire(in)
	if wire.Name != "A" || wire.Color == nil || *wire.Color != "00ff00" {
		t.Errorf("wire = %+v, want name A color 00ff00", wire)
	}
}

func TestWireRackRoleToDomain(t *testing.T) {
	t.Parallel()

	wire := WireRackRole{
		ID:      1,
		URL:     "http://nb/api/dcim/rack-roles/1/",
		Name:    "Core",
		Slug:    "core",
		Color:   "00ff00",
		Created: "2024-01-01",
	}
	got := wireRackRoleToDomain(wire)
	if got.ID != 1 || got.Name != "Core" || got.Color != "00ff00" {
		t.Errorf("domain = %+v, want id 1 name Core color 00ff00", got)
	}
	if got := wireRackRoleToDomain(WireRackRole{}).Color; got != "" {
		t.Errorf("empty wire Color = %q, want empty", got)
	}
}

// --- device_role ----------------------------------------------------------

func TestClient_ListDeviceRoles(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/dcim/device-roles/" {
			writer.WriteHeader(http.StatusNotFound)

			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"count":1,"results":[{"id":1,"name":"Spine","slug":"spine","color":"ff0000",
			"weight":100,"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	resp, err := client.ListDeviceRoles(context.Background(), "test-token", nil)
	if err != nil {
		t.Fatalf("ListDeviceRoles() returned error: %v", err)
	}
	if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "Spine" || resp.Results[0].Weight != 100 {
		t.Errorf("resp = %+v, want count 1 with Spine/weight 100", resp)
	}
}

func TestClient_CreateDeviceRole(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		_, _ = writer.Write([]byte(`{"id":1,"name":"Spine","slug":"spine","color":"ff0000",` +
			`"weight":100,"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	deviceRole, err := client.CreateDeviceRole(context.Background(), "tok", domain.DeviceRoleWrite{Name: "Spine"})
	if err != nil {
		t.Fatalf("CreateDeviceRole() returned error: %v", err)
	}
	if deviceRole.ID != 1 || deviceRole.Name != "Spine" {
		t.Errorf("device role = %+v, want id 1 name Spine", deviceRole)
	}
}

func TestClient_UpdateDeviceRole(t *testing.T) {
	t.Parallel()
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(`{"id":7,"name":"Renamed","slug":"renamed","created":"","last_updated":""}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	deviceRole, err := client.UpdateDeviceRole(context.Background(), "tok", 7, domain.DeviceRoleWrite{Name: "Renamed"})
	if err != nil {
		t.Fatalf("UpdateDeviceRole() returned error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/api/dcim/device-roles/7/" {
		t.Errorf("method/path = %q %q, want PATCH /api/dcim/device-roles/7/", gotMethod, gotPath)
	}
	if deviceRole.ID != 7 || deviceRole.Name != "Renamed" {
		t.Errorf("device role = %+v, want id 7 name Renamed", deviceRole)
	}
}

func TestClient_DeleteDeviceRole(t *testing.T) {
	t.Parallel()
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		writer.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	if err := client.DeleteDeviceRole(context.Background(), "tok", 3); err != nil {
		t.Fatalf("DeleteDeviceRole() returned error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/api/dcim/device-roles/3/" {
		t.Errorf("method/path = %q %q, want DELETE /api/dcim/device-roles/3/", gotMethod, gotPath)
	}
}

func TestDeviceRoleWriteToWire(t *testing.T) {
	t.Parallel()

	vmRole := true
	weight := 50
	ct := 9
	in := domain.DeviceRoleWrite{Name: "A", VMRole: &vmRole, Weight: &weight, ConfigTemplate: &ct}
	wire := deviceRoleWriteToWire(in)
	if wire.Name != "A" || wire.VMRole == nil || !*wire.VMRole || wire.Weight == nil || *wire.Weight != 50 {
		t.Errorf("wire = %+v, want name A vm_role true weight 50", wire)
	}
	if wire.ConfigTemplate == nil || *wire.ConfigTemplate != 9 {
		t.Errorf("wire.ConfigTemplate = %v, want 9", wire.ConfigTemplate)
	}
}

func TestWireDeviceRoleToDomain(t *testing.T) {
	t.Parallel()

	wire := WireDeviceRole{
		ID:             1,
		URL:            "http://nb/api/dcim/device-roles/1/",
		Name:           "Spine",
		Slug:           "spine",
		VMRole:         true,
		ConfigTemplate: &WireNested{ID: 9, Name: "cfg", Slug: "cfg"},
		Weight:         100,
		Created:        "2024-01-01",
	}
	got := wireDeviceRoleToDomain(wire)
	if got.ID != 1 || got.Name != "Spine" || got.Weight != 100 {
		t.Errorf("domain = %+v, want id 1 name Spine weight 100", got)
	}
	if got.ConfigTemplate == nil || got.ConfigTemplate.ID != 9 {
		t.Errorf("domain.ConfigTemplate = %+v, want nested id 9", got.ConfigTemplate)
	}
	if got := wireDeviceRoleToDomain(WireDeviceRole{}).ConfigTemplate; got != nil {
		t.Errorf("empty wire ConfigTemplate = %+v, want nil", got)
	}
}

// --- platform -------------------------------------------------------------

func TestClient_ListPlatforms(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/dcim/platforms/" {
			writer.WriteHeader(http.StatusNotFound)

			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"count":1,"results":[{"id":1,"name":"Linux","slug":"linux",
			"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	resp, err := client.ListPlatforms(context.Background(), "test-token", nil)
	if err != nil {
		t.Fatalf("ListPlatforms() returned error: %v", err)
	}
	if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "Linux" {
		t.Errorf("resp = %+v, want count 1 with Linux", resp)
	}
}

func TestClient_CreatePlatform(t *testing.T) {
	t.Parallel()
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		_, _ = writer.Write([]byte(`{"id":1,"name":"Linux","slug":"linux","created":"2024-01-01",` +
			`"last_updated":"2024-01-01T00:00:00Z"}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	platform, err := client.CreatePlatform(context.Background(), "tok", domain.PlatformWrite{Name: "Linux"})
	if err != nil {
		t.Fatalf("CreatePlatform() returned error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/api/dcim/platforms/" {
		t.Errorf("method/path = %q %q, want POST /api/dcim/platforms/", gotMethod, gotPath)
	}
	if platform.ID != 1 || platform.Name != "Linux" {
		t.Errorf("platform = %+v, want id 1 name Linux", platform)
	}
}

func TestClient_UpdatePlatform(t *testing.T) {
	t.Parallel()
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(`{"id":7,"name":"Renamed","slug":"renamed","created":"","last_updated":""}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	platform, err := client.UpdatePlatform(context.Background(), "tok", 7, domain.PlatformWrite{Name: "Renamed"})
	if err != nil {
		t.Fatalf("UpdatePlatform() returned error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/api/dcim/platforms/7/" {
		t.Errorf("method/path = %q %q, want PATCH /api/dcim/platforms/7/", gotMethod, gotPath)
	}
	if platform.ID != 7 || platform.Name != "Renamed" {
		t.Errorf("platform = %+v, want id 7 name Renamed", platform)
	}
}

func TestClient_DeletePlatform(t *testing.T) {
	t.Parallel()
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		writer.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	if err := client.DeletePlatform(context.Background(), "tok", 3); err != nil {
		t.Fatalf("DeletePlatform() returned error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/api/dcim/platforms/3/" {
		t.Errorf("method/path = %q %q, want DELETE /api/dcim/platforms/3/", gotMethod, gotPath)
	}
}

func TestPlatformWriteToWire(t *testing.T) {
	t.Parallel()

	mfr := 3
	in := domain.PlatformWrite{Name: "A", Manufacturer: &mfr}
	wire := platformWriteToWire(in)
	if wire.Name != "A" || wire.Manufacturer == nil || *wire.Manufacturer != 3 {
		t.Errorf("wire = %+v, want name A manufacturer 3", wire)
	}
}

func TestWirePlatformToDomain(t *testing.T) {
	t.Parallel()

	wire := WirePlatform{
		ID:           1,
		URL:          "http://nb/api/dcim/platforms/1/",
		Name:         "Linux",
		Slug:         "linux",
		Manufacturer: &WireNested{ID: 3, Name: "Cisco", Slug: "cisco"},
		Created:      "2024-01-01",
	}
	got := wirePlatformToDomain(wire)
	if got.ID != 1 || got.Name != "Linux" {
		t.Errorf("domain = %+v, want id 1 name Linux", got)
	}
	if got.Manufacturer == nil || got.Manufacturer.ID != 3 {
		t.Errorf("domain.Manufacturer = %+v, want nested id 3", got.Manufacturer)
	}
	if got := wirePlatformToDomain(WirePlatform{}).Manufacturer; got != nil {
		t.Errorf("empty wire Manufacturer = %+v, want nil", got)
	}
}

// --- virtual_chassis ------------------------------------------------------

func TestClient_ListVirtualChassis(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/dcim/virtual-chassis/" {
			writer.WriteHeader(http.StatusNotFound)

			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"count":1,"results":[{"id":1,"name":"VC-1","domain":"example.com",
			"member_count":2,"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	resp, err := client.ListVirtualChassis(context.Background(), "test-token", nil)
	if err != nil {
		t.Fatalf("ListVirtualChassis() returned error: %v", err)
	}
	if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "VC-1" || resp.Results[0].MemberCount != 2 {
		t.Errorf("resp = %+v, want count 1 with VC-1/member_count 2", resp)
	}
}

func TestClient_CreateVirtualChassis(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		_, _ = writer.Write([]byte(`{"id":1,"name":"VC-1","domain":"example.com","member_count":1,` +
			`"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	chassis, err := client.CreateVirtualChassis(context.Background(), "tok", domain.VirtualChassisWrite{Name: "VC-1"})
	if err != nil {
		t.Fatalf("CreateVirtualChassis() returned error: %v", err)
	}
	if chassis.ID != 1 || chassis.Name != "VC-1" || chassis.Domain != "example.com" {
		t.Errorf("virtual chassis = %+v, want id 1 name VC-1 domain example.com", chassis)
	}
}

func TestClient_UpdateVirtualChassis(t *testing.T) {
	t.Parallel()
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(`{"id":7,"name":"Renamed","domain":"new.example.com","created":"","last_updated":""}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	chassis, err := client.UpdateVirtualChassis(
		context.Background(), "tok", 7, domain.VirtualChassisWrite{Name: "Renamed"},
	)
	if err != nil {
		t.Fatalf("UpdateVirtualChassis() returned error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/api/dcim/virtual-chassis/7/" {
		t.Errorf("method/path = %q %q, want PATCH /api/dcim/virtual-chassis/7/", gotMethod, gotPath)
	}
	if chassis.ID != 7 || chassis.Name != "Renamed" {
		t.Errorf("virtual chassis = %+v, want id 7 name Renamed", chassis)
	}
}

func TestClient_DeleteVirtualChassis(t *testing.T) {
	t.Parallel()
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		writer.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	if err := client.DeleteVirtualChassis(context.Background(), "tok", 3); err != nil {
		t.Fatalf("DeleteVirtualChassis() returned error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/api/dcim/virtual-chassis/3/" {
		t.Errorf("method/path = %q %q, want DELETE /api/dcim/virtual-chassis/3/", gotMethod, gotPath)
	}
}

func TestVirtualChassisWriteToWire(t *testing.T) {
	t.Parallel()

	domainStr := "example.com"
	in := domain.VirtualChassisWrite{Name: "A", Domain: &domainStr}
	wire := virtualChassisWriteToWire(in)
	if wire.Name != "A" || wire.Domain == nil || *wire.Domain != "example.com" {
		t.Errorf("wire = %+v, want name A domain example.com", wire)
	}
}

func TestWireVirtualChassisToDomain(t *testing.T) {
	t.Parallel()

	wire := WireVirtualChassis{
		ID:          1,
		URL:         "http://nb/api/dcim/virtual-chassis/1/",
		Name:        "VC-1",
		Domain:      "example.com",
		Master:      &WireNested{ID: 4, Name: "sw1", Slug: "sw1"},
		MemberCount: 2,
		Created:     "2024-01-01",
	}
	got := wireVirtualChassisToDomain(wire)
	if got.ID != 1 || got.Name != "VC-1" || got.Domain != "example.com" || got.MemberCount != 2 {
		t.Errorf("domain = %+v, want id 1 name VC-1 domain example.com member_count 2", got)
	}
	if got.Master == nil || got.Master.ID != 4 {
		t.Errorf("domain.Master = %+v, want nested id 4", got.Master)
	}
	if got := wireVirtualChassisToDomain(WireVirtualChassis{}).Master; got != nil {
		t.Errorf("empty wire Master = %+v, want nil", got)
	}
}

// --- power_panel ----------------------------------------------------------

func TestClient_ListPowerPanels(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/dcim/power-panels/" {
			writer.WriteHeader(http.StatusNotFound)

			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"count":1,"results":[{"id":1,"name":"PP-1",
			"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	resp, err := client.ListPowerPanels(context.Background(), "test-token", nil)
	if err != nil {
		t.Fatalf("ListPowerPanels() returned error: %v", err)
	}
	if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "PP-1" {
		t.Errorf("resp = %+v, want count 1 with PP-1", resp)
	}
}

func TestClient_CreatePowerPanel(t *testing.T) {
	t.Parallel()
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		_, _ = writer.Write([]byte(`{"id":1,"name":"PP-1","created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	panel, err := client.CreatePowerPanel(context.Background(), "tok", domain.PowerPanelWrite{Name: "PP-1"})
	if err != nil {
		t.Fatalf("CreatePowerPanel() returned error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/api/dcim/power-panels/" {
		t.Errorf("method/path = %q %q, want POST /api/dcim/power-panels/", gotMethod, gotPath)
	}
	if panel.ID != 1 || panel.Name != "PP-1" {
		t.Errorf("power panel = %+v, want id 1 name PP-1", panel)
	}
}

func TestClient_UpdatePowerPanel(t *testing.T) {
	t.Parallel()
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(`{"id":7,"name":"Renamed","created":"","last_updated":""}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	panel, err := client.UpdatePowerPanel(context.Background(), "tok", 7, domain.PowerPanelWrite{Name: "Renamed"})
	if err != nil {
		t.Fatalf("UpdatePowerPanel() returned error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/api/dcim/power-panels/7/" {
		t.Errorf("method/path = %q %q, want PATCH /api/dcim/power-panels/7/", gotMethod, gotPath)
	}
	if panel.ID != 7 || panel.Name != "Renamed" {
		t.Errorf("power panel = %+v, want id 7 name Renamed", panel)
	}
}

func TestClient_DeletePowerPanel(t *testing.T) {
	t.Parallel()
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		writer.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	if err := client.DeletePowerPanel(context.Background(), "tok", 3); err != nil {
		t.Fatalf("DeletePowerPanel() returned error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/api/dcim/power-panels/3/" {
		t.Errorf("method/path = %q %q, want DELETE /api/dcim/power-panels/3/", gotMethod, gotPath)
	}
}

func TestPowerPanelWriteToWire(t *testing.T) {
	t.Parallel()

	site := 5
	loc := 6
	in := domain.PowerPanelWrite{Name: "A", Site: &site, Location: &loc}
	wire := powerPanelWriteToWire(in)
	if wire.Name != "A" || wire.Site == nil || *wire.Site != 5 || wire.Location == nil || *wire.Location != 6 {
		t.Errorf("wire = %+v, want name A site 5 location 6", wire)
	}
}

func TestWirePowerPanelToDomain(t *testing.T) {
	t.Parallel()

	wire := WirePowerPanel{
		ID:       1,
		URL:      "http://nb/api/dcim/power-panels/1/",
		Name:     "PP-1",
		Display:  "PP-1",
		Site:     &WireNested{ID: 5, Name: "DC1", Slug: "dc1"},
		Location: &WireNested{ID: 6, Name: "RowA", Slug: "row-a"},
		Created:  "2024-01-01",
	}
	got := wirePowerPanelToDomain(wire)
	if got.ID != 1 || got.Name != "PP-1" {
		t.Errorf("domain = %+v, want id 1 name PP-1", got)
	}
	if got.Site == nil || got.Site.ID != 5 {
		t.Errorf("domain.Site = %+v, want nested id 5", got.Site)
	}
	if got.Location == nil || got.Location.ID != 6 {
		t.Errorf("domain.Location = %+v, want nested id 6", got.Location)
	}
	if got := wirePowerPanelToDomain(WirePowerPanel{}).Site; got != nil {
		t.Errorf("empty wire Site = %+v, want nil", got)
	}
}

// --- inventory_item_role --------------------------------------------------

func TestClient_ListInventoryItemRoles(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/dcim/inventory-item-roles/" {
			writer.WriteHeader(http.StatusNotFound)

			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"count":1,"results":[{"id":1,"name":"PSU","slug":"psu","color":"0000ff",
			"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	resp, err := client.ListInventoryItemRoles(context.Background(), "test-token", nil)
	if err != nil {
		t.Fatalf("ListInventoryItemRoles() returned error: %v", err)
	}
	if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "PSU" || resp.Results[0].Color != "0000ff" {
		t.Errorf("resp = %+v, want count 1 with PSU/0000ff", resp)
	}
}

func TestClient_CreateInventoryItemRole(t *testing.T) {
	t.Parallel()
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		_, _ = writer.Write([]byte(`{"id":1,"name":"PSU","slug":"psu","color":"0000ff",` +
			`"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	invRole, err := client.CreateInventoryItemRole(context.Background(), "tok", domain.InventoryItemRoleWrite{Name: "PSU"})
	if err != nil {
		t.Fatalf("CreateInventoryItemRole() returned error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/api/dcim/inventory-item-roles/" {
		t.Errorf("method/path = %q %q, want POST /api/dcim/inventory-item-roles/", gotMethod, gotPath)
	}
	if invRole.ID != 1 || invRole.Name != "PSU" || invRole.Color != "0000ff" {
		t.Errorf("inventory item role = %+v, want id 1 name PSU color 0000ff", invRole)
	}
}

func TestClient_UpdateInventoryItemRole(t *testing.T) {
	t.Parallel()
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(`{"id":7,"name":"Renamed","slug":"renamed","color":"ff0000",` +
			`"created":"","last_updated":""}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	invRole, err := client.UpdateInventoryItemRole(
		context.Background(), "tok", 7, domain.InventoryItemRoleWrite{Name: "Renamed"},
	)
	if err != nil {
		t.Fatalf("UpdateInventoryItemRole() returned error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/api/dcim/inventory-item-roles/7/" {
		t.Errorf("method/path = %q %q, want PATCH /api/dcim/inventory-item-roles/7/", gotMethod, gotPath)
	}
	if invRole.ID != 7 || invRole.Name != "Renamed" {
		t.Errorf("inventory item role = %+v, want id 7 name Renamed", invRole)
	}
}

func TestClient_DeleteInventoryItemRole(t *testing.T) {
	t.Parallel()
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		writer.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	if err := client.DeleteInventoryItemRole(context.Background(), "tok", 3); err != nil {
		t.Fatalf("DeleteInventoryItemRole() returned error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/api/dcim/inventory-item-roles/3/" {
		t.Errorf("method/path = %q %q, want DELETE /api/dcim/inventory-item-roles/3/", gotMethod, gotPath)
	}
}

func TestInventoryItemRoleWriteToWire(t *testing.T) {
	t.Parallel()

	color := "0000ff"
	in := domain.InventoryItemRoleWrite{Name: "A", Color: &color}
	wire := inventoryItemRoleWriteToWire(in)
	if wire.Name != "A" || wire.Color == nil || *wire.Color != "0000ff" {
		t.Errorf("wire = %+v, want name A color 0000ff", wire)
	}
}

func TestWireInventoryItemRoleToDomain(t *testing.T) {
	t.Parallel()

	wire := WireInventoryItemRole{
		ID:      1,
		URL:     "http://nb/api/dcim/inventory-item-roles/1/",
		Name:    "PSU",
		Slug:    "psu",
		Color:   "0000ff",
		Created: "2024-01-01",
	}
	got := wireInventoryItemRoleToDomain(wire)
	if got.ID != 1 || got.Name != "PSU" || got.Color != "0000ff" {
		t.Errorf("domain = %+v, want id 1 name PSU color 0000ff", got)
	}
	if got := wireInventoryItemRoleToDomain(WireInventoryItemRole{}).Color; got != "" {
		t.Errorf("empty wire Color = %q, want empty", got)
	}
}

// TestClient_ListNewEntities_Errors exercises the error branches of the 8 new
// DCIM list methods: a transport-level error (server returns 500) and a
// malformed JSON response body (unmarshal failure). These branches are
// otherwise uncovered and would otherwise fail the coverage gate.
//
//nolint:funlen // Large table-driven test; splitting is not sensible.
func TestClient_ListNewEntities_Errors(t *testing.T) {
	t.Parallel()

	list := map[string]func(*Client, context.Context) (int, error){
		"ListSiteGroups": func(c *Client, ctx context.Context) (int, error) {
			resp, err := c.ListSiteGroups(ctx, "t", nil)
			if resp == nil {
				return 0, err
			}

			return resp.Count, err
		},
		"ListRackGroups": func(c *Client, ctx context.Context) (int, error) {
			resp, err := c.ListRackGroups(ctx, "t", nil)
			if resp == nil {
				return 0, err
			}

			return resp.Count, err
		},
		"ListRackRoles": func(c *Client, ctx context.Context) (int, error) {
			resp, err := c.ListRackRoles(ctx, "t", nil)
			if resp == nil {
				return 0, err
			}

			return resp.Count, err
		},
		"ListDeviceRoles": func(c *Client, ctx context.Context) (int, error) {
			resp, err := c.ListDeviceRoles(ctx, "t", nil)
			if resp == nil {
				return 0, err
			}

			return resp.Count, err
		},
		"ListPlatforms": func(c *Client, ctx context.Context) (int, error) {
			resp, err := c.ListPlatforms(ctx, "t", nil)
			if resp == nil {
				return 0, err
			}

			return resp.Count, err
		},
		"ListVirtualChassis": func(c *Client, ctx context.Context) (int, error) {
			resp, err := c.ListVirtualChassis(ctx, "t", nil)
			if resp == nil {
				return 0, err
			}

			return resp.Count, err
		},
		"ListPowerPanels": func(c *Client, ctx context.Context) (int, error) {
			resp, err := c.ListPowerPanels(ctx, "t", nil)
			if resp == nil {
				return 0, err
			}

			return resp.Count, err
		},
		"ListInventoryItemRoles": func(c *Client, ctx context.Context) (int, error) {
			resp, err := c.ListInventoryItemRoles(ctx, "t", nil)
			if resp == nil {
				return 0, err
			}

			return resp.Count, err
		},
	}

	t.Run("transport error", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		for name, call := range list {
			t.Run(name, func(t *testing.T) {
				if _, err := call(client, context.Background()); err == nil {
					t.Errorf("%s: expected error on 500, got nil", name)
				}
			})
		}
	})

	t.Run("unmarshal error", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{"count":`)) // truncated JSON
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		for name, call := range list {
			t.Run(name, func(t *testing.T) {
				if _, err := call(client, context.Background()); err == nil {
					t.Errorf("%s: expected unmarshal error, got nil", name)
				}
			})
		}
	})
}

// TestClient_NewEntities_ValidationErrors exercises the 400-validation-error
// branch of the create/update/delete methods for all 8 new DCIM entities.
// Each write method maps a 400 response to a domain.ValidationError; without
// these tests that branch is left uncovered.
//
//nolint:funlen,lll,nlreturn // Large table-driven test; long closure lines are inherent.
func TestClient_NewEntities_ValidationErrors(t *testing.T) {
	t.Parallel()

	body := `{"name":["This field is required."]}`
	detailBody := `{"detail":"Cannot delete object with dependent objects."}`

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodDelete {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(detailBody))

			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte(body))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	ctx := context.Background()

	t.Run("create validation error", func(t *testing.T) {
		cases := []struct {
			name string
			call func() error
		}{
			{"site_group", func() error { _, err := client.CreateSiteGroup(ctx, "t", domain.SiteGroupWrite{}); return err }},
			{"rack_group", func() error { _, err := client.CreateRackGroup(ctx, "t", domain.RackGroupWrite{}); return err }},
			{"rack_role", func() error { _, err := client.CreateRackRole(ctx, "t", domain.RackRoleWrite{}); return err }},
			{"device_role", func() error { _, err := client.CreateDeviceRole(ctx, "t", domain.DeviceRoleWrite{}); return err }},
			{"platform", func() error { _, err := client.CreatePlatform(ctx, "t", domain.PlatformWrite{}); return err }},
			{"virtual_chassis", func() error {
				_, err := client.CreateVirtualChassis(ctx, "t", domain.VirtualChassisWrite{})
				return err
			}},
			{"power_panel", func() error { _, err := client.CreatePowerPanel(ctx, "t", domain.PowerPanelWrite{}); return err }},
			{"inventory_item_role", func() error {
				_, err := client.CreateInventoryItemRole(ctx, "t", domain.InventoryItemRoleWrite{})
				return err
			}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				err := tc.call()
				var validationErr *domain.ValidationError
				if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
					t.Errorf("%s: err = %v, want ValidationError body %s", tc.name, err, body)
				}
			})
		}
	})

	t.Run("update validation error", func(t *testing.T) {
		cases := []struct {
			name string
			call func() error
		}{
			{"site_group", func() error { _, err := client.UpdateSiteGroup(ctx, "t", 7, domain.SiteGroupWrite{}); return err }},
			{"rack_group", func() error { _, err := client.UpdateRackGroup(ctx, "t", 7, domain.RackGroupWrite{}); return err }},
			{"rack_role", func() error { _, err := client.UpdateRackRole(ctx, "t", 7, domain.RackRoleWrite{}); return err }},
			{"device_role", func() error { _, err := client.UpdateDeviceRole(ctx, "t", 7, domain.DeviceRoleWrite{}); return err }},
			{"platform", func() error { _, err := client.UpdatePlatform(ctx, "t", 7, domain.PlatformWrite{}); return err }},
			{"virtual_chassis", func() error {
				_, err := client.UpdateVirtualChassis(ctx, "t", 7, domain.VirtualChassisWrite{})
				return err
			}},
			{"power_panel", func() error { _, err := client.UpdatePowerPanel(ctx, "t", 7, domain.PowerPanelWrite{}); return err }},
			{"inventory_item_role", func() error {
				_, err := client.UpdateInventoryItemRole(ctx, "t", 7, domain.InventoryItemRoleWrite{})
				return err
			}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				err := tc.call()
				var validationErr *domain.ValidationError
				if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
					t.Errorf("%s: err = %v, want ValidationError body %s", tc.name, err, body)
				}
			})
		}
	})

	t.Run("delete validation error", func(t *testing.T) {
		cases := []struct {
			name string
			call func() error
		}{
			{"site_group", func() error { return client.DeleteSiteGroup(ctx, "t", 3) }},
			{"rack_group", func() error { return client.DeleteRackGroup(ctx, "t", 3) }},
			{"rack_role", func() error { return client.DeleteRackRole(ctx, "t", 3) }},
			{"device_role", func() error { return client.DeleteDeviceRole(ctx, "t", 3) }},
			{"platform", func() error { return client.DeletePlatform(ctx, "t", 3) }},
			{"virtual_chassis", func() error { return client.DeleteVirtualChassis(ctx, "t", 3) }},
			{"power_panel", func() error { return client.DeletePowerPanel(ctx, "t", 3) }},
			{"inventory_item_role", func() error { return client.DeleteInventoryItemRole(ctx, "t", 3) }},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				err := tc.call()
				var validationErr *domain.ValidationError
				if !errors.As(err, &validationErr) || string(validationErr.Body) != detailBody {
					t.Errorf("%s: err = %v, want ValidationError body %s", tc.name, err, detailBody)
				}
			})
		}
	})
}
