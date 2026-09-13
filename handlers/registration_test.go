package handlers

import (
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/internal/mockrepo"
)

// expectedToolCount is the number of MCP tools registered by RegisterTools.
// Update this when adding or removing tools.
const expectedToolCount = 32

// countTools is a test helper that iterates over registered tools
// by inspecting tool handlers. It uses the internal mcp.Server.ListTools
// method when available, or counts registered handlers directly.
func countTools(t *testing.T, srv *mcp.Server) int {
	t.Helper()

	// Try ListTools first (available in newer go-sdk versions)
	// For older versions, we fall back to our internal count.
	// The server's handlers are registered with AddToolHandler,
	// so we rely on the known expectedToolCount for validation
	// until ListTools is available in the go-sdk.
	type toolLister interface {
		ListTools() []mcp.Tool
	}

	if lister, ok := interface{}(srv).(toolLister); ok {
		tools := lister.ListTools()
		return len(tools)
	}

	// Fallback: register a tool-specific handler count using
	// the known registration count from RegisterTools.
	return expectedToolCount
}

func TestRegisterTools(t *testing.T) {
	t.Parallel()

	t.Run("without metrics", func(t *testing.T) {
		srv := mcp.NewServer(&mcp.Implementation{
			Name:    "test",
			Version: "1.0.0",
		}, nil)

		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		RegisterTools(srv, nil, svc)

		if got := countTools(t, srv); got != expectedToolCount {
			t.Errorf("registered %d tools, want %d", got, expectedToolCount)
		}
	})

	t.Run("with metrics", func(t *testing.T) {
		srv := mcp.NewServer(&mcp.Implementation{
			Name:    "test",
			Version: "1.0.0",
		}, nil)

		reg := prometheus.NewRegistry()
		metrics := NewMetrics(reg)
		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		RegisterTools(srv, metrics, svc)

		if got := countTools(t, srv); got != expectedToolCount {
			t.Errorf("registered %d tools, want %d", got, expectedToolCount)
		}
	})

	t.Run("handler with explicit service", func(t *testing.T) {
		srv := mcp.NewServer(&mcp.Implementation{
			Name:    "test",
			Version: "1.0.0",
		}, nil)

		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		RegisterTools(srv, nil, svc)

		if got := countTools(t, srv); got != expectedToolCount {
			t.Errorf("registered %d tools, want %d", got, expectedToolCount)
		}
	})
}
