package handlers

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sirupsen/logrus"
)

// secretArgs is a tool-input-like struct with a secret:"true" field used to
// verify the stdio access log redacts tagged secrets from the args field.
type secretArgs struct {
	Secret string `json:"secret" secret:"true"`
	Name   string `json:"name"`
}

// captureAccessLog runs action with a capture-enabled logger and the per-tool
// access log enabled, then returns everything written to the capture buffer.
func captureAccessLog(t *testing.T, action func()) string {
	t.Helper()

	var buf bytes.Buffer
	l := logrus.New()
	l.SetOutput(&buf)
	l.SetLevel(logrus.InfoLevel)
	SetLogger(l)
	EnableToolAccessLog()

	t.Cleanup(func() {
		SetLogger(nil)
		toolAccessLogEnabled.Store(false)
	})

	action()

	return buf.String()
}

func callToolRequest(name string) *mcp.CallToolRequest {
	return &mcp.CallToolRequest{Params: &mcp.CallToolParamsRaw{Name: name}}
}

func TestWrapToolHandler_StdioAccessLog_Success(t *testing.T) {
	out := captureAccessLog(t, func() {
		handler := WrapToolHandler[testInput, testOutput](nil, "get_sites", func(
			_ context.Context, _ *mcp.CallToolRequest, _ testInput,
		) (*mcp.CallToolResult, testOutput, error) {
			return &mcp.CallToolResult{}, testOutput{Result: "ok"}, nil
		})
		_, _, _ = handler(context.Background(), callToolRequest("get_sites"), testInput{Name: "dc1"})
	})

	// L08: the per-tool access log is at info level, carrying tool, source,
	// duration and outcome, and never a debug gate.
	for _, want := range []string{"mcp_request", "get_sites", "source=stdio", "outcome=success"} {
		if !strings.Contains(out, want) {
			t.Errorf("access log missing %q; got: %s", want, out)
		}
	}
}

func TestWrapToolHandler_StdioAccessLog_Error(t *testing.T) {
	out := captureAccessLog(t, func() {
		handler := WrapToolHandler[testInput, testOutput](nil, "get_sites", func(
			_ context.Context, _ *mcp.CallToolRequest, _ testInput,
		) (*mcp.CallToolResult, testOutput, error) {
			return &mcp.CallToolResult{IsError: true}, testOutput{}, nil
		})
		_, _, _ = handler(context.Background(), callToolRequest("get_sites"), testInput{})
	})

	if !strings.Contains(out, "outcome=error") {
		t.Errorf("access log = %q, want it to contain outcome=error", out)
	}
}

func TestWrapToolHandler_StdioAccessLog_RedactsSecrets(t *testing.T) {
	const secret = "super-secret-value"
	out := captureAccessLog(t, func() {
		handler := WrapToolHandler[secretArgs, testOutput](nil, "update_site", func(
			_ context.Context, _ *mcp.CallToolRequest, _ secretArgs,
		) (*mcp.CallToolResult, testOutput, error) {
			return &mcp.CallToolResult{}, testOutput{Result: "ok"}, nil
		})
		_, _, _ = handler(context.Background(), callToolRequest("update_site"), secretArgs{Secret: secret, Name: "site-a"})
	})

	if strings.Contains(out, secret) {
		t.Errorf("access log leaks a secret value; got: %s", out)
	}
	if !strings.Contains(out, "site-a") {
		t.Errorf("access log should still carry non-secret args; got: %s", out)
	}
}

func TestWrapToolHandler_AccessLogDisabledByDefault(t *testing.T) {
	// The flag is off by default: no access log is emitted until the stdio
	// transport enables it (so the HTTP transport is not double-logged).
	var buf bytes.Buffer
	l := logrus.New()
	l.SetOutput(&buf)
	l.SetLevel(logrus.InfoLevel)
	SetLogger(l)
	t.Cleanup(func() { SetLogger(nil) })

	handler := WrapToolHandler[testInput, testOutput](nil, "get_sites", func(
		_ context.Context, _ *mcp.CallToolRequest, _ testInput,
	) (*mcp.CallToolResult, testOutput, error) {
		return &mcp.CallToolResult{}, testOutput{Result: "ok"}, nil
	})
	_, _, _ = handler(context.Background(), callToolRequest("get_sites"), testInput{})

	if strings.Contains(buf.String(), "mcp_request") {
		t.Errorf("access log emitted while disabled; got: %s", buf.String())
	}
}
