//go:build e2e

package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/teran/go-docker-testsuite/applications/netbox"
)

func TestDiagService(t *testing.T) {
	if testing.Short() {
		t.Skip("short")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Minute)
	defer cancel()
	app, err := netbox.NewWithT(t, ctx, netboxImage)
	if err != nil {
		t.Fatalf("start netbox: %v", err)
	}
	defer app.Close(context.Background())

	token := app.SuperuserAPIToken()
	base := app.MustURL()

	post := func(path string, payload map[string]any) {
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, base+path, bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Logf("POST %s error: %v", path, err)
			return
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		t.Logf("POST %s -> status=%d body=%s", path, resp.StatusCode, string(b))
	}

	post("/api/dcim/sites/", map[string]any{"name": "s1", "slug": "s1"})
	post("/api/dcim/manufacturers/", map[string]any{"name": "m1", "slug": "m1"})
	post("/api/dcim/device-types/", map[string]any{"model": "dt1", "slug": "dt1", "manufacturer": 1})
	post("/api/dcim/device-roles/", map[string]any{"name": "r1", "slug": "r1"})
	post("/api/dcim/devices/", map[string]any{"name": "dev1", "device_type": 1, "role": 1, "site": 1})
	post("/api/virtualization/virtual-machines/", map[string]any{"name": "vm1", "site": 1})

	post("/api/ipam/services/", map[string]any{
		"name": "svc1", "parent_object_type": "dcim.device", "parent_object_id": 1,
		"protocol": "tcp", "ports": []int{80, 443},
	})
	post("/api/ipam/services/", map[string]any{
		"name": "svc2", "parent_object_type": "virtualization.virtualmachine", "parent_object_id": 1,
		"protocol": "udp", "ports": []int{53},
	})
}
