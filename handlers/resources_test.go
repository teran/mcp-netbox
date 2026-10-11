package handlers

import (
	"context"
	"strings"
	"testing"
	"unicode"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestReadGuidanceHandler_AllURIsResolvable(t *testing.T) {
	t.Parallel()

	h := readGuidanceHandler()
	for _, resource := range guidanceResources {
		t.Run(resource.uri, func(t *testing.T) {
			res, err := h(context.Background(), &mcp.ReadResourceRequest{
				Params: &mcp.ReadResourceParams{URI: resource.uri},
			})
			if err != nil {
				t.Fatalf("read %s: %v", resource.uri, err)
			}
			if len(res.Contents) != 1 {
				t.Fatalf("read %s: len(Contents) = %d, want 1", resource.uri, len(res.Contents))
			}
			c := res.Contents[0]
			if c.URI != resource.uri {
				t.Errorf("content URI = %q, want %q", c.URI, resource.uri)
			}
			if c.Text != resource.markdown {
				t.Errorf("content text does not match the resource markdown for %s", resource.uri)
			}
		})
	}
}

func TestReadGuidanceHandler_UnknownURIReturnsNotFound(t *testing.T) {
	t.Parallel()

	h := readGuidanceHandler()
	_, err := h(context.Background(), &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{URI: guidanceScheme + "does-not-exist"},
	})
	if err == nil {
		t.Fatal("handler returned nil error for an unknown URI, want ResourceNotFoundError")
	}
	if !strings.Contains(err.Error(), "Resource not found") {
		t.Errorf("error = %q, want it to carry 'Resource not found'", err.Error())
	}
}

func TestReadGuidanceHandler_MIMETypeIsTextMarkdown(t *testing.T) {
	t.Parallel()

	h := readGuidanceHandler()
	for _, resource := range guidanceResources {
		res, err := h(context.Background(), &mcp.ReadResourceRequest{
			Params: &mcp.ReadResourceParams{URI: resource.uri},
		})
		if err != nil {
			t.Fatalf("read %s: %v", resource.uri, err)
		}
		if got := res.Contents[0].MIMEType; got != "text/markdown" {
			t.Errorf("read %s: MIMEType = %q, want %q", resource.uri, got, "text/markdown")
		}
	}
}

func TestReadGuidanceHandler_MarkdownControlClean(t *testing.T) {
	t.Parallel()

	// Hygiene: no served markdown body may contain ANSI escapes or control
	// characters (S09/N23), so guidance can never inject terminal escapes.
	for _, guidance := range guidanceResources {
		for _, r := range guidance.markdown {
			if r == '\n' || r == '\t' {
				continue
			}
			if unicode.IsControl(r) {
				t.Errorf("resource %s markdown contains control character %q", guidance.uri, r)
			}
		}
	}
}

func TestRegisterResources_RegistersAllGuidanceResources(t *testing.T) {
	t.Parallel()

	// A freshly built server must expose every guidance resource so the
	// resources/list capability advertises the full invariant set.
	s := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "test"}, &mcp.ServerOptions{})
	RegisterResources(s)

	for _, resource := range guidanceResources {
		// The handler must resolve every registered URI end to end.
		res, err := readGuidanceHandler()(context.Background(), &mcp.ReadResourceRequest{
			Params: &mcp.ReadResourceParams{URI: resource.uri},
		})
		if err != nil {
			t.Errorf("registered resource %s: %v", resource.uri, err)

			continue
		}
		if res == nil || len(res.Contents) == 0 {
			t.Errorf("registered resource %s returned no contents", resource.uri)
		}
	}
}
