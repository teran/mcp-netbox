package handlers

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// guidanceScheme is the URI scheme under which the read-only usage-guidance
// resources are served (M08). Clients read them via resources/read.
const guidanceScheme = "netbox://guidance/"

// guidanceResource is a single read-only markdown resource served under the
// resources capability. It carries the audited server-instructions text split
// per topic so a client can fetch just the guidance it needs.
type guidanceResource struct {
	uri, name, title string
	markdown         string
}

// guidanceResources is the fixed inventory of read-only guidance resources.
// The ordering is invariant: overview, then read → write → delete tool groups,
// then validation, pagination and authentication.
var guidanceResources = []guidanceResource{
	{
		uri:   guidanceScheme + "overview",
		name:  "Overview",
		title: "Overview",
		markdown: "mcp-netbox provides CRUD access to the NetBox inventory. " +
			"Read tools (`get_*`) are safe and never modify NetBox; " +
			"write tools create/update objects; delete tools permanently remove objects.",
	},
	{
		uri:   guidanceScheme + "read-tools",
		name:  "Read tools",
		title: "Read tools",
		markdown: "Read tools (`get_*`) are safe and never modify NetBox. " +
			"All list tools are paginated: use `page` (1-based) and " +
			"`page_size` (max 1000) to navigate results. An empty results array " +
			"means no objects matched your filters — it is not an error. " +
			"Filters are additive: provide only the fields relevant to your query.",
	},
	{
		uri:   guidanceScheme + "write-tools",
		name:  "Write tools",
		title: "Write tools",
		markdown: "Write tools create and update objects (`create_*` / `update_*`); " +
			"they are not idempotent. Update tools perform a partial PATCH merge: " +
			"only the fields you provide are changed.",
	},
	{
		uri:   guidanceScheme + "delete-tools",
		name:  "Delete tools",
		title: "Delete tools",
		markdown: "Delete tools (`delete_*`) permanently and irreversibly remove " +
			"an object and are marked with a destructive hint. Never call a delete " +
			"tool unless the user has explicitly asked to delete that specific object.",
	},
	{
		uri:   guidanceScheme + "validation",
		name:  "Validation",
		title: "Validation",
		markdown: "Validation is enforced by NetBox: a rejected create/update " +
			"returns the DRF validation body to the model so it can correct the payload.",
	},
	{
		uri:   guidanceScheme + "pagination",
		name:  "Pagination",
		title: "Pagination",
		markdown: "All list tools are paginated: use `page` (1-based) and " +
			"`page_size` (max 1000) to navigate results. An empty results array " +
			"means no objects matched your filters — it is not an error.",
	},
	{
		uri:   guidanceScheme + "authentication",
		name:  "Authentication",
		title: "Authentication",
		markdown: "In HTTP/SSE mode the NetBox token is passed in the " +
			"`Authorization: Bearer` header. In STDIO mode it is provided via the " +
			"`NETBOX_TOKEN` environment variable at server startup. Authorization " +
			"is enforced entirely by NetBox; this server does not validate tokens.",
	},
}

// readGuidanceHandler returns the ResourceHandler that serves the read-only
// guidance resources. It resolves a request's URI against the fixed inventory
// and returns the matching markdown body, or ResourceNotFoundError when the URI
// is unknown.
func readGuidanceHandler() mcp.ResourceHandler {
	byURI := make(map[string]guidanceResource, len(guidanceResources))
	for _, g := range guidanceResources {
		byURI[g.uri] = g
	}

	return func(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		resource, ok := byURI[req.Params.URI]
		if !ok {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{URI: resource.uri, MIMEType: "text/markdown", Text: resource.markdown},
			},
		}, nil
	}
}

// RegisterResources registers the read-only guidance resources on the MCP
// server under the resources capability. All resources share a single handler
// that serves their fixed markdown bodies.
func RegisterResources(s *mcp.Server) {
	handler := readGuidanceHandler()
	for _, g := range guidanceResources {
		s.AddResource(&mcp.Resource{
			URI:         g.uri,
			Name:        g.name,
			Title:       g.title,
			Description: "Read-only usage guidance for mcp-netbox (served under the resources capability).",
			MIMEType:    "text/markdown",
		}, handler)
	}
}
