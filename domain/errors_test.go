package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestValidationError_Error(t *testing.T) {
	t.Parallel()

	ve := &ValidationError{StatusCode: 400, Body: []byte(`{"name":["required"]}`)}
	msg := ve.Error()
	if !strings.Contains(msg, "400") {
		t.Errorf("Error() = %q, want it to contain the status code 400", msg)
	}
	if strings.Contains(msg, "required") || strings.Contains(msg, "name") {
		t.Errorf("Error() = %q, must NOT contain the response body", msg)
	}
	if strings.Contains(msg, `{`) {
		t.Errorf("Error() = %q, must not contain JSON body", msg)
	}
}

func TestValidationError_Error_NoBody(t *testing.T) {
	t.Parallel()

	ve := &ValidationError{StatusCode: 422, Body: []byte(`{"detail":"something secret"}`)}
	if got := ve.Error(); got != "netbox validation error: status 422" {
		t.Errorf("Error() = %q, want %q", got, "netbox validation error: status 422")
	}
}

func TestSiteWrite_PartialUpdateOmitsUnsetFields(t *testing.T) {
	t.Parallel()

	// A partially-populated write (e.g. update of description only) must marshal
	// to JSON containing only the set fields — nil pointers must be omitted so
	// the PATCH body is minimal and does not clear unrelated fields.
	w := SiteWrite{Description: ptr("new desc")}

	b, err := json.Marshal(w)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	s := string(b)
	if !strings.Contains(s, `"description":"new desc"`) {
		t.Errorf("marshal = %s, want description present", s)
	}
	for _, forbidden := range []string{"name", "slug", "status", "region", "tenant", "facility", "time_zone", "comments"} {
		if strings.Contains(s, `"`+forbidden+`"`) {
			t.Errorf("marshal = %s, field %q should be omitted when unset", s, forbidden)
		}
	}
}

func TestSiteWrite_CreateIncludesName(t *testing.T) {
	t.Parallel()

	w := SiteWrite{Name: "Site A", Slug: ptr("site-a")}
	b, err := json.Marshal(w)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	s := string(b)
	if !strings.Contains(s, `"name":"Site A"`) {
		t.Errorf("marshal = %s, want name present", s)
	}
	if !strings.Contains(s, `"slug":"site-a"`) {
		t.Errorf("marshal = %s, want slug present", s)
	}
}

func TestSiteWrite_ExplicitEmptyString(t *testing.T) {
	t.Parallel()

	// A pointer to an empty string must marshal as an explicit empty value so a
	// PATCH can clear a field (distinct from "not provided").
	w := SiteWrite{Description: ptr("")}
	b, err := json.Marshal(w)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	if !strings.Contains(string(b), `"description":""`) {
		t.Errorf("marshal = %s, want explicit empty description", string(b))
	}
}

func TestDeviceWrite_PartialUpdateOmitsUnsetFields(t *testing.T) {
	t.Parallel()

	// A partially-populated write (e.g. update of serial only) must marshal to
	// JSON containing only the set fields — nil pointers must be omitted so the
	// PATCH body is minimal and does not clear unrelated fields.
	w := DeviceWrite{Serial: ptr("SN-123")}

	b, err := json.Marshal(w)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	s := string(b)
	if !strings.Contains(s, `"serial":"SN-123"`) {
		t.Errorf("marshal = %s, want serial present", s)
	}
	for _, forbidden := range []string{"name", "device_type", "role", "tenant", "platform", "asset_tag", "site", "rack", "position", "face", "status", "cluster", "comments"} {
		if strings.Contains(s, `"`+forbidden+`"`) {
			t.Errorf("marshal = %s, field %q should be omitted when unset", s, forbidden)
		}
	}
}

func TestDeviceWrite_CreateIncludesName(t *testing.T) {
	t.Parallel()

	w := DeviceWrite{Name: "router1", DeviceType: ptrInt(2), Site: ptrInt(3)}
	b, err := json.Marshal(w)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	s := string(b)
	if !strings.Contains(s, `"name":"router1"`) {
		t.Errorf("marshal = %s, want name present", s)
	}
	if !strings.Contains(s, `"device_type":2`) {
		t.Errorf("marshal = %s, want device_type present", s)
	}
	if !strings.Contains(s, `"site":3`) {
		t.Errorf("marshal = %s, want site present", s)
	}
}

func TestDeviceWrite_ExplicitEmptyString(t *testing.T) {
	t.Parallel()

	// A pointer to an empty string must marshal as an explicit empty value so a
	// PATCH can clear a field (distinct from "not provided").
	w := DeviceWrite{Serial: ptr("")}
	b, err := json.Marshal(w)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	if !strings.Contains(string(b), `"serial":""`) {
		t.Errorf("marshal = %s, want explicit empty serial", string(b))
	}
}

func ptr(s string) *string { return &s }

func ptrInt(i int) *int { return &i }
