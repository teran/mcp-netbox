package main

import (
	"bufio"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/teran/mcp-netbox/config"
	"github.com/teran/mcp-netbox/logging"
)

// stdioSession drives a runStdio server over an in-memory pair of pipes.
type stdioSession struct {
	inR    *io.PipeReader // server reads requests from here
	inW    *io.PipeWriter // test writes requests here
	outR   *io.PipeReader // test reads responses from here
	outW   *io.PipeWriter // server writes responses here
	reader *bufio.Reader
}

func newStdioSession() *stdioSession {
	inR, inW := io.Pipe()
	outR, outW := io.Pipe()
	return &stdioSession{
		inR: inR, inW: inW,
		outR: outR, outW: outW,
		reader: bufio.NewReader(outR),
	}
}

func (s *stdioSession) transport() mcp.Transport {
	return &mcp.IOTransport{Reader: s.inR, Writer: s.outW}
}

func (s *stdioSession) send(obj map[string]interface{}) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	if _, err := s.inW.Write(append(b, '\n')); err != nil {
		return err
	}
	return nil
}

type rpcResp struct {
	ID     json.RawMessage `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  json.RawMessage `json:"error"`
}

func (s *stdioSession) readResp(t *testing.T) *rpcResp {
	t.Helper()

	type res struct {
		line []byte
		err  error
	}
	ch := make(chan res, 1)
	go func() {
		line, err := s.reader.ReadBytes('\n')
		ch <- res{line, err}
	}()
	select {
	case r := <-ch:
		if r.err != nil {
			t.Fatalf("reading response: %v", r.err)
		}
		var resp rpcResp
		if err := json.Unmarshal(r.line, &resp); err != nil {
			t.Fatalf("unmarshal response %q: %v", string(r.line), err)
		}
		return &resp
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for response")
		return nil
	}
}

func (s *stdioSession) close() {
	_ = s.inW.Close()
}

// testStdioConfig builds a stdio config whose NetBox URL is a loopback port
// that the DNS-rebinding dialer rejects immediately, so tools/call returns
// fast with an error instead of hanging on a real network connection.
func testStdioConfig() config.Config {
	return config.Config{
		NetBoxURL:          "http://127.0.0.1:1",
		AllowPrivateNetBox: true,
		Transport:          config.TransportStdio,
		NetBoxToken:        "test-token",
		RateLimitGlobal:    100,
		RateLimitPerClient: 10,
		WriteTimeout:       300 * time.Second,
	}
}

func TestRunStdio_EndToEnd(t *testing.T) {
	sess := newStdioSession()
	defer sess.close()

	cfg := testStdioConfig()
	logger, err := logging.Setup(cfg)
	if err != nil {
		t.Fatalf("logging.Setup: %v", err)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- runStdioWithTransport(cfg, logger, sess.transport())
	}()

	// 1. initialize
	if err := sess.send(map[string]interface{}{
		"jsonrpc": "2.0", "id": 1, "method": "initialize",
		"params": map[string]interface{}{
			"protocolVersion": "2025-06-18",
			"capabilities":    map[string]interface{}{},
			"clientInfo":      map[string]interface{}{"name": "test", "version": "1"},
		},
	}); err != nil {
		t.Fatalf("send initialize: %v", err)
	}
	if initResp := sess.readResp(t); len(initResp.Error) > 0 {
		t.Fatalf("initialize error: %s", initResp.Error)
	}

	// 2. initialized notification
	if err := sess.send(map[string]interface{}{
		"jsonrpc": "2.0", "method": "notifications/initialized",
	}); err != nil {
		t.Fatalf("send initialized: %v", err)
	}

	// 3. tools/list — must list all registered read-only tools
	if err := sess.send(map[string]interface{}{
		"jsonrpc": "2.0", "id": 2, "method": "tools/list",
	}); err != nil {
		t.Fatalf("send tools/list: %v", err)
	}
	listResp := sess.readResp(t)
	if len(listResp.Error) > 0 {
		t.Fatalf("tools/list error: %s", listResp.Error)
	}
	var tools struct {
		Tools []struct {
			Name string `json:"name"`
		} `json:"tools"`
	}
	if err := json.Unmarshal(listResp.Result, &tools); err != nil {
		t.Fatalf("unmarshal tools: %v", err)
	}
	if len(tools.Tools) == 0 {
		t.Fatal("tools/list returned no tools")
	}
	found := false
	for _, tl := range tools.Tools {
		if tl.Name == "get_sites" {
			found = true
		}
	}
	if !found {
		t.Errorf("tools/list missing get_sites; got %d tools", len(tools.Tools))
	}

	// 4. tools/call — attempts a real NetBox request that fails fast (the
	//    dialer rejects loopback), but must still return a JSON-RPC response.
	if err := sess.send(map[string]interface{}{
		"jsonrpc": "2.0", "id": 3, "method": "tools/call",
		"params": map[string]interface{}{
			"name":      "get_sites",
			"arguments": map[string]interface{}{"page": 1, "page_size": 5},
		},
	}); err != nil {
		t.Fatalf("send tools/call: %v", err)
	}
	callResp := sess.readResp(t)
	if len(callResp.Result) == 0 && len(callResp.Error) == 0 {
		t.Fatal("tools/call produced neither result nor error")
	}

	// Close the input to unblock conn.Wait() in the server goroutine.
	sess.close()

	select {
	case err := <-errCh:
		// runStdio may return nil or an EOF-derived error once the client
		// disconnects; either is acceptable for a clean shutdown.
		_ = err
	case <-time.After(5 * time.Second):
		t.Fatal("server did not stop after input closed")
	}
}

func TestRun_UnsupportedTransport(t *testing.T) {
	cfg := testStdioConfig()
	cfg.Transport = "websocket"
	if err := Run(cfg); err == nil {
		t.Fatal("Run() = nil, want error for unsupported transport")
	}
}
