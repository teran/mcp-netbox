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

// TestClient_ListCustomFields exercises the extras custom-fields list endpoint.
func TestClient_ListCustomFields(t *testing.T) { //nolint:funlen
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Authorization") != "Token test-token" {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}
			if req.URL.Path != "/api/extras/custom-fields/" {
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
					"url": "http://nb/api/extras/custom-fields/1/",
					"name": "env",
					"label": "Environment",
					"type": {"value": "text", "label": "Text"},
					"object_types": ["dcim.device", "dcim.site"],
					"group_name": "core",
					"description": "deployment environment",
					"required": true,
					"default": "prod",
					"filter_logic": {"value": "exact", "label": "Exact"},
					"weight": 10,
					"ui_visible": {"value": "read-write", "label": "Read-write"},
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListCustomFields(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListCustomFields() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "env" {
			t.Errorf("resp = %+v, want count 1 with env", resp)
		}
		if resp.Results[0].Label != "Environment" {
			t.Errorf("Label = %q, want %q", resp.Results[0].Label, "Environment")
		}
		if resp.Results[0].Type != "text" {
			t.Errorf("Type = %q, want %q", resp.Results[0].Type, "text")
		}
		if len(resp.Results[0].ContentTypes) != 2 || resp.Results[0].ContentTypes[0] != "dcim.device" {
			t.Errorf("ContentTypes = %v, want [dcim.device dcim.site]", resp.Results[0].ContentTypes)
		}
		if resp.Results[0].Weight != 10 {
			t.Errorf("Weight = %d, want %d", resp.Results[0].Weight, 10)
		}
		if !resp.Results[0].Required {
			t.Error("Required = false, want true")
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListCustomFields(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("bad_json", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{not-json`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListCustomFields(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected unmarshal error, got nil")
		}
	})
}

func TestClient_CreateCustomField(t *testing.T) {
	t.Parallel()

	testClientCreateCustomFieldSuccess(t)
	testClientCreateCustomFieldValidationError(t)
}

func testClientCreateCustomFieldSuccess(t *testing.T) {
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
		_, _ = writer.Write([]byte(`{"id":1,"url":"http://nb/api/extras/custom-fields/1/","name":"env",` +
			`"label":"Environment","type":{"value":"text","label":"Text"},"object_types":["dcim.device"],"group_name":"core",` +
			`"description":"","required":true,"default":"","filter_logic":{"value":"exact","label":"Exact"},"weight":10,` +
			`"ui_visible":{"value":"read-write","label":"Read-write"},"created":"2024-01-01",` +
			`"last_updated":"2024-01-01T00:00:00Z"}`))
	}))
	defer srv.Close()

	typ := "text"
	required := true
	weight := 10
	client := NewClient(srv.URL, http.DefaultClient)
	created, err := client.CreateCustomField(context.Background(), "tok", domain.CustomFieldWrite{
		Name: "env", Label: "Environment", Type: &typ, ContentTypes: []string{"dcim.device"},
		Required: &required, Weight: &weight,
	})
	if err != nil {
		t.Fatalf("CreateCustomField() returned error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/api/extras/custom-fields/" {
		t.Errorf("method/path = %q %q, want POST /api/extras/custom-fields/", gotMethod, gotPath)
	}
	if gotCT != "application/json" || gotAuth != "Token tok" {
		t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
	}
	var reqBody map[string]any
	_ = json.Unmarshal(gotBody, &reqBody)
	if reqBody["name"] != "env" {
		t.Errorf("request body = %s, want name env", gotBody)
	}
	if reqBody["label"] != "Environment" {
		t.Errorf("request body label = %v, want Environment", reqBody["label"])
	}
	contentTypes, ok := reqBody["object_types"].([]any)
	if !ok || len(contentTypes) != 1 || contentTypes[0] != "dcim.device" {
		t.Errorf("request body object_types = %v, want [dcim.device]", reqBody["object_types"])
	}
	if created.ID != 1 || created.Name != "env" || created.Label != "Environment" {
		t.Errorf("created = %+v, want id 1 name env label Environment", created)
	}
	if created.Type != "text" || created.Weight != 10 {
		t.Errorf("created = %+v, want type text weight 10", created)
	}
}

func testClientCreateCustomFieldValidationError(t *testing.T) {
	body := `{"name":["This field is required."]}`
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte(body))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, http.DefaultClient)
	_, err := client.CreateCustomField(context.Background(), "tok", domain.CustomFieldWrite{})
	var validationErr *domain.ValidationError
	if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
		t.Errorf("err = %v, want ValidationError body %s", err, body)
	}
}

func TestClient_UpdateCustomField(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		var gotBody []byte
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			gotBody, _ = io.ReadAll(req.Body)
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"url":"","name":"renamed","label":"Renamed",` +
				`"type":{"value":"text","label":"Text"},"object_types":[],"group_name":"","description":"",` +
				`"required":false,"default":"","filter_logic":{"value":"loose","label":"Loose"},"weight":100,` +
				`"ui_visible":{"value":"always","label":"Always"},"created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		updated, err := client.UpdateCustomField(context.Background(), "tok", 7, domain.CustomFieldWrite{Name: "renamed"})
		if err != nil {
			t.Fatalf("UpdateCustomField() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/extras/custom-fields/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/extras/custom-fields/7/", gotMethod, gotPath)
		}
		if updated.ID != 7 || updated.Name != "renamed" {
			t.Errorf("updated = %+v, want id 7 name renamed", updated)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "renamed" {
			t.Errorf("request body = %s, want name present", gotBody)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"label":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateCustomField(context.Background(), "tok", 7, domain.CustomFieldWrite{Name: "X"})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteCustomField(t *testing.T) {
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
		if err := client.DeleteCustomField(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteCustomField() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/extras/custom-fields/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/extras/custom-fields/3/", gotMethod, gotPath)
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
		err := client.DeleteCustomField(context.Background(), "tok", 3)
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
		if err := client.DeleteCustomField(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteCustomField() 404 returned error: %v", err)
		}
	})
}

func TestCustomFieldWriteToWire(t *testing.T) {
	t.Parallel()

	typ := "text"
	groupName := "core"
	required := true
	weight := 10
	input := domain.CustomFieldWrite{
		Name:         "env",
		Label:        "Environment",
		Type:         &typ,
		ContentTypes: []string{"dcim.device"},
		GroupName:    &groupName,
		Required:     &required,
		Weight:       &weight,
		Choices:      []string{"prod", "dev"},
	}
	wire := customFieldWriteToWire(input)
	if wire.Name != "env" {
		t.Errorf("wire.Name = %q, want env", wire.Name)
	}
	if wire.Label != "Environment" {
		t.Errorf("wire.Label = %q, want Environment", wire.Label)
	}
	if wire.Type == nil || *wire.Type != "text" {
		t.Errorf("wire.Type = %v, want text", wire.Type)
	}
	if len(wire.ContentTypes) != 1 || wire.ContentTypes[0] != "dcim.device" {
		t.Errorf("wire.ContentTypes = %v, want [dcim.device]", wire.ContentTypes)
	}
	if wire.Weight == nil || *wire.Weight != 10 {
		t.Errorf("wire.Weight = %v, want 10", wire.Weight)
	}
	if len(wire.Choices) != 2 {
		t.Errorf("wire.Choices = %v, want 2 entries", wire.Choices)
	}
}

func TestWireCustomFieldToDomain(t *testing.T) {
	t.Parallel()

	wire := WireCustomField{
		ID:           1,
		URL:          "http://nb/api/extras/custom-fields/1/",
		Name:         "env",
		Label:        "Environment",
		Type:         WireChoice{Value: "text", Label: "Text"},
		ContentTypes: []string{"dcim.device"},
		GroupName:    "core",
		Description:  "desc",
		Required:     true,
		Default:      "prod",
		FilterLogic:  WireChoice{Value: "exact", Label: "Exact"},
		Weight:       10,
		UIVisible:    WireChoice{Value: "read-write", Label: "Read-write"},
	}
	got := wireCustomFieldToDomain(wire)
	if got.ID != 1 || got.Name != "env" || got.Label != "Environment" {
		t.Errorf("domain = %+v, want id 1 name env label Environment", got)
	}
	if got.Type != "text" || got.Weight != 10 || !got.Required {
		t.Errorf("domain = %+v, want type text weight 10 required true", got)
	}
	if got.FilterLogic != "exact" {
		t.Errorf("domain.FilterLogic = %q, want exact", got.FilterLogic)
	}
	if got.UIVisibility != "read-write" {
		t.Errorf("domain.UIVisibility = %q, want read-write", got.UIVisibility)
	}
	if len(got.ContentTypes) != 1 || got.ContentTypes[0] != "dcim.device" {
		t.Errorf("domain.ContentTypes = %v, want [dcim.device]", got.ContentTypes)
	}
	if len(got.Choices) != 0 {
		t.Errorf("domain.Choices = %v, want empty (NetBox no longer returns inline choices)", got.Choices)
	}

	if got := wireCustomFieldToDomain(WireCustomField{}).Name; got != "" {
		t.Errorf("empty wire name = %q, want empty", got)
	}
}

func TestUnmarshalCustomField_BadJSON(t *testing.T) {
	t.Parallel()

	if _, err := unmarshalCustomField(domain.RawObject(`{not-json`)); err == nil {
		t.Fatal("unmarshalCustomField() expected error on invalid JSON, got nil")
	}
}
