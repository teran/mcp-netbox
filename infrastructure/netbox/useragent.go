package netbox

import (
	"fmt"
	"sync/atomic"
)

// Outbound User-Agent identity (M10). The NetBox client sends an identifying
// User-Agent in the form "app-name/version (reference)". For this public repo
// the reference is the repository URL. It is a public identifier and must never
// contain the NetBox token.
const (
	appName      = "mcp-netbox"
	appReference = "https://github.com/teran/mcp-netbox"
)

// defaultUserAgent is the non-empty identifying value used when no version is
// set (e.g. a local/dev build or unit tests), so the header is never empty or
// malformed (N44).
const defaultUserAgent = appName + " (" + appReference + ")"

// userAgentVersion is the ldflags-embedded app version, injected at startup via
// SetUserAgent. It is empty until then, in which case userAgent falls back to
// defaultUserAgent (N44).
var userAgentVersion atomic.Pointer[string]

// SetUserAgent configures the app version used in the outbound User-Agent
// header (M10). When version is empty, userAgent returns the non-empty
// identifying default so the header is never empty or malformed (N44).
func SetUserAgent(version string) {
	userAgentVersion.Store(&version)
}

// userAgent returns the User-Agent value sent on outbound NetBox requests.
// It is "mcp-netbox/<version> (https://github.com/teran/mcp-netbox)" when a
// version is set, or the non-empty default otherwise (N44).
func userAgent() string {
	if v := userAgentVersion.Load(); v != nil && *v != "" {
		return fmt.Sprintf("%s/%s (%s)", appName, *v, appReference)
	}

	return defaultUserAgent
}
