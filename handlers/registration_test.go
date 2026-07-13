package handlers

import (
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/teran/mcp-netbox/application"
	"github.com/teran/mcp-netbox/internal/mockrepo"
)

func TestRegisterTools(t *testing.T) {
	t.Parallel()

	t.Run("without metrics", func(t *testing.T) {
		srv := mcp.NewServer(&mcp.Implementation{
			Name:    "test",
			Version: "1.0.0",
		}, nil)

		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		RegisterTools(srv, nil, svc)
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
	})

	t.Run("handler with explicit service", func(t *testing.T) {
		srv := mcp.NewServer(&mcp.Implementation{
			Name:    "test",
			Version: "1.0.0",
		}, nil)

		svc := application.NewNetworkService(&mockrepo.NilRepo{}, "token")
		RegisterTools(srv, nil, svc)
	})
}
