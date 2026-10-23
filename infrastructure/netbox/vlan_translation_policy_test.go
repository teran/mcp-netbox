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

// TestClient_ListVLANTranslationPolicies exercises the VLAN translation policy list endpoint.
func TestClient_ListVLANTranslationPolicies(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/ipam/vlan-translation-policies/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"url": "http://nb/api/ipam/vlan-translation-policies/1/",
					"name": "datacenter",
					"description": "core",
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListVLANTranslationPolicies(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListVLANTranslationPolicies() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "datacenter" {
			t.Errorf("resp = %+v, want count 1 with datacenter", resp)
		}
		if resp.Results[0].Description != "core" {
			t.Errorf("resp.Results[0].Description = %q, want core", resp.Results[0].Description)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListVLANTranslationPolicies(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestClient_CreateVLANTranslationPolicy(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"name":"datacenter","created":"2024-01-01",` +
				`"last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		policy, err := client.CreateVLANTranslationPolicy(
			context.Background(), "tok", domain.VLANTranslationPolicyWrite{Name: "datacenter", Tags: []string{"prod"}},
		)
		if err != nil {
			t.Fatalf("CreateVLANTranslationPolicy() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/ipam/vlan-translation-policies/" {
			t.Errorf("method/path = %q %q, want POST /api/ipam/vlan-translation-policies/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["name"] != "datacenter" {
			t.Errorf("request body = %s, want name datacenter", gotBody)
		}
		if policy.ID != 1 || policy.Name != "datacenter" {
			t.Errorf("policy = %+v, want id 1 name datacenter", policy)
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
		_, err := client.CreateVLANTranslationPolicy(context.Background(), "tok", domain.VLANTranslationPolicyWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateVLANTranslationPolicy(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
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
		policy, err := client.UpdateVLANTranslationPolicy(
			context.Background(), "tok", 7, domain.VLANTranslationPolicyWrite{Name: "Renamed"},
		)
		if err != nil {
			t.Fatalf("UpdateVLANTranslationPolicy() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/ipam/vlan-translation-policies/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/ipam/vlan-translation-policies/7/", gotMethod, gotPath)
		}
		if policy.ID != 7 || policy.Name != "Renamed" {
			t.Errorf("policy = %+v, want id 7 name Renamed", policy)
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
		_, err := client.UpdateVLANTranslationPolicy(
			context.Background(), "tok", 7, domain.VLANTranslationPolicyWrite{Name: "X"},
		)
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteVLANTranslationPolicy(t *testing.T) {
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
		if err := client.DeleteVLANTranslationPolicy(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteVLANTranslationPolicy() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/ipam/vlan-translation-policies/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/ipam/vlan-translation-policies/3/", gotMethod, gotPath)
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
		err := client.DeleteVLANTranslationPolicy(context.Background(), "tok", 3)
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
		if err := client.DeleteVLANTranslationPolicy(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteVLANTranslationPolicy() 404 returned error: %v", err)
		}
	})
}

func TestVLANTranslationPolicyWriteToWire(t *testing.T) {
	t.Parallel()

	desc := "core"
	in := domain.VLANTranslationPolicyWrite{Name: "datacenter", Description: &desc, Tags: []string{"prod"}}
	wire := vlanTranslationPolicyWriteToWire(in)
	if wire.Name != "datacenter" || wire.Description == nil || *wire.Description != "core" {
		t.Errorf("wire = %+v, want name datacenter description core", wire)
	}

	raw, err := json.Marshal(wire)
	if err != nil {
		t.Fatalf("marshal wire: %v", err)
	}
	for _, key := range []string{"name", "description", "tags"} {
		if !jsonHasKey(raw, key) {
			t.Errorf("wire JSON missing key %q: %s", key, raw)
		}
	}
}

func TestWireVLANTranslationPolicyToDomain(t *testing.T) {
	t.Parallel()

	wire := WireVLANTranslationPolicy{
		ID:          1,
		URL:         "http://nb/api/ipam/vlan-translation-policies/1/",
		Name:        "datacenter",
		Description: "core",
		Created:     "2024-01-01",
	}
	got := wireVLANTranslationPolicyToDomain(wire)
	if got.ID != 1 || got.Name != "datacenter" || got.Description != "core" {
		t.Errorf("domain = %+v, want id 1 name datacenter description core", got)
	}
}
