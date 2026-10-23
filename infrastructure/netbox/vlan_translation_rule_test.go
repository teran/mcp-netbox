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

//nolint:funlen // Table of subtests covering list success/error branches.
func TestClient_ListVLANTranslationRules(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/ipam/vlan-translation-rules/" {
				writer.WriteHeader(http.StatusNotFound)

				return
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{
				"count": 1,
				"results": [{
					"id": 1,
					"url": "http://nb/api/ipam/vlan-translation-rules/1/",
					"display": "100 -> 200 (datacenter)",
					"policy": 5,
					"local_vid": 100,
					"remote_vid": 200,
					"description": "edge",
					"created": "2024-01-01",
					"last_updated": "2024-01-01T00:00:00Z"
				}]
			}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		resp, err := client.ListVLANTranslationRules(context.Background(), "test-token", nil)
		if err != nil {
			t.Fatalf("ListVLANTranslationRules() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 {
			t.Errorf("resp = %+v, want count 1", resp)
		}
		if resp.Results[0].Policy != 5 {
			t.Errorf("resp.Results[0].Policy = %d, want 5", resp.Results[0].Policy)
		}
		if resp.Results[0].LocalVID != 100 {
			t.Errorf("resp.Results[0].LocalVID = %d, want 100", resp.Results[0].LocalVID)
		}
		if resp.Results[0].RemoteVID != 200 {
			t.Errorf("resp.Results[0].RemoteVID = %d, want 200", resp.Results[0].RemoteVID)
		}
		if resp.Results[0].Description != "edge" {
			t.Errorf("resp.Results[0].Description = %q, want edge", resp.Results[0].Description)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.ListVLANTranslationRules(context.Background(), "token", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

//nolint:funlen // Large table-driven/declarative function; splitting is not sensible.
func TestClient_CreateVLANTranslationRule(t *testing.T) {
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
			_, _ = writer.Write([]byte(`{"id":1,"policy":5,"local_vid":100,"remote_vid":200,` +
				`"created":"2024-01-01","last_updated":"2024-01-01T00:00:00Z"}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		localVID := 100
		remoteVID := 200
		rule, err := client.CreateVLANTranslationRule(
			context.Background(), "tok",
			domain.VLANTranslationRuleWrite{Policy: 5, LocalVID: &localVID, RemoteVID: &remoteVID},
		)
		if err != nil {
			t.Fatalf("CreateVLANTranslationRule() returned error: %v", err)
		}
		if gotMethod != http.MethodPost || gotPath != "/api/ipam/vlan-translation-rules/" {
			t.Errorf("method/path = %q %q, want POST /api/ipam/vlan-translation-rules/", gotMethod, gotPath)
		}
		if gotCT != "application/json" || gotAuth != "Token tok" {
			t.Errorf("CT/Auth = %q %q, want application/json / Token tok", gotCT, gotAuth)
		}
		var reqBody map[string]any
		_ = json.Unmarshal(gotBody, &reqBody)
		if reqBody["policy"] != float64(5) {
			t.Errorf("request body = %s, want policy 5", gotBody)
		}
		if reqBody["local_vid"] != float64(100) || reqBody["remote_vid"] != float64(200) {
			t.Errorf("request body local/remote = %v/%v, want 100/200", reqBody["local_vid"], reqBody["remote_vid"])
		}
		if rule.ID != 1 {
			t.Errorf("rule = %+v, want id 1", rule)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"policy":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.CreateVLANTranslationRule(context.Background(), "tok", domain.VLANTranslationRuleWrite{})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_UpdateVLANTranslationRule(t *testing.T) {
	t.Parallel()

	t.Run("success patch", func(t *testing.T) {
		t.Parallel()
		var gotMethod, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"id":7,"policy":5,"local_vid":100,"remote_vid":200,` +
				`"created":"","last_updated":""}`))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		rule, err := client.UpdateVLANTranslationRule(
			context.Background(), "tok", 7, domain.VLANTranslationRuleWrite{Policy: 5},
		)
		if err != nil {
			t.Fatalf("UpdateVLANTranslationRule() returned error: %v", err)
		}
		if gotMethod != http.MethodPatch || gotPath != "/api/ipam/vlan-translation-rules/7/" {
			t.Errorf("method/path = %q %q, want PATCH /api/ipam/vlan-translation-rules/7/", gotMethod, gotPath)
		}
		if rule.ID != 7 {
			t.Errorf("rule = %+v, want id 7", rule)
		}
	})

	t.Run("validation error 400", func(t *testing.T) {
		t.Parallel()
		body := `{"local_vid":["This field is required."]}`
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(body))
		}))
		defer srv.Close()

		client := NewClient(srv.URL, http.DefaultClient)
		_, err := client.UpdateVLANTranslationRule(context.Background(), "tok", 7, domain.VLANTranslationRuleWrite{Policy: 5})
		var validationErr *domain.ValidationError
		if !errors.As(err, &validationErr) || string(validationErr.Body) != body {
			t.Errorf("err = %v, want ValidationError body %s", err, body)
		}
	})
}

func TestClient_DeleteVLANTranslationRule(t *testing.T) {
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
		if err := client.DeleteVLANTranslationRule(context.Background(), "tok", 3); err != nil {
			t.Fatalf("DeleteVLANTranslationRule() returned error: %v", err)
		}
		if gotMethod != http.MethodDelete || gotPath != "/api/ipam/vlan-translation-rules/3/" {
			t.Errorf("method/path = %q %q, want DELETE /api/ipam/vlan-translation-rules/3/", gotMethod, gotPath)
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
		err := client.DeleteVLANTranslationRule(context.Background(), "tok", 3)
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
		if err := client.DeleteVLANTranslationRule(context.Background(), "tok", 999); err != nil {
			t.Fatalf("DeleteVLANTranslationRule() 404 returned error: %v", err)
		}
	})
}

func TestVLANTranslationRuleWriteToWire(t *testing.T) {
	t.Parallel()

	localVID := 100
	remoteVID := 200
	in := domain.VLANTranslationRuleWrite{Policy: 5, LocalVID: &localVID, RemoteVID: &remoteVID}
	wire := vlanTranslationRuleWriteToWire(in)
	if wire.Policy != 5 || wire.LocalVID == nil || *wire.LocalVID != 100 {
		t.Errorf("wire = %+v, want policy 5 local_vid 100", wire)
	}
	if wire.RemoteVID == nil || *wire.RemoteVID != 200 {
		t.Errorf("wire.RemoteVID = %v, want 200", wire.RemoteVID)
	}

	raw, err := json.Marshal(wire)
	if err != nil {
		t.Fatalf("marshal wire: %v", err)
	}
	for _, key := range []string{"policy", "local_vid", "remote_vid"} {
		if !jsonHasKey(raw, key) {
			t.Errorf("wire JSON missing key %q: %s", key, raw)
		}
	}
}

func TestWireVLANTranslationRuleToDomain(t *testing.T) {
	t.Parallel()

	wire := WireVLANTranslationRule{
		ID:          1,
		URL:         "http://nb/api/ipam/vlan-translation-rules/1/",
		Display:     "100 -> 200 (datacenter)",
		Policy:      5,
		LocalVID:    100,
		RemoteVID:   200,
		Description: "edge",
		Created:     "2024-01-01",
	}
	got := wireVLANTranslationRuleToDomain(wire)
	if got.ID != 1 {
		t.Errorf("domain = %+v, want id 1", got)
	}
	if got.Policy != 5 {
		t.Errorf("domain.Policy = %d, want 5", got.Policy)
	}
	if got.LocalVID != 100 {
		t.Errorf("domain.LocalVID = %d, want 100", got.LocalVID)
	}
	if got.RemoteVID != 200 {
		t.Errorf("domain.RemoteVID = %d, want 200", got.RemoteVID)
	}
	if got.Description != "edge" {
		t.Errorf("domain.Description = %q, want edge", got.Description)
	}
	if got := wireVLANTranslationRuleToDomain(WireVLANTranslationRule{}).Policy; got != 0 {
		t.Errorf("empty wire Policy = %d, want 0", got)
	}
}
