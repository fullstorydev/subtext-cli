package tunnelmcp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

// exchange sends a single JSON-RPC request to a Server and returns the decoded response.
func exchange(t *testing.T, srv *Server, req map[string]any) map[string]any {
	t.Helper()
	reqData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	in := bytes.NewReader(append(reqData, '\n'))
	var out bytes.Buffer

	ctx, cancel := context.WithCancel(context.Background())
	// Run Serve in the background; it exits once stdin hits EOF.
	done := make(chan error, 1)
	go func() {
		done <- srv.Serve(ctx, in, &out)
		cancel()
	}()
	<-done

	// Parse the single response line.
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	if !scanner.Scan() {
		t.Fatal("no response from server")
	}
	var resp map[string]any
	if err := json.Unmarshal(scanner.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v\nraw: %s", err, scanner.Bytes())
	}
	return resp
}

func TestInitialize(t *testing.T) {
	srv := New("test-version")
	resp := exchange(t, srv, map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]any{},
		},
	})

	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result, ok := resp["result"].(map[string]any)
	if !ok {
		t.Fatalf("result missing or wrong type: %v", resp["result"])
	}
	if result["protocolVersion"] != "2024-11-05" {
		t.Errorf("protocolVersion = %v, want 2024-11-05", result["protocolVersion"])
	}
	si, _ := result["serverInfo"].(map[string]any)
	if si["name"] != "subtext_tunnel" {
		t.Errorf("serverInfo.name = %v, want subtext_tunnel", si["name"])
	}
	if si["version"] != "test-version" {
		t.Errorf("serverInfo.version = %v, want test-version", si["version"])
	}
}

func TestToolsList(t *testing.T) {
	srv := New("test-version")
	resp := exchange(t, srv, map[string]any{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/list",
	})

	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result, ok := resp["result"].(map[string]any)
	if !ok {
		t.Fatalf("result missing: %v", resp)
	}
	tools, ok := result["tools"].([]any)
	if !ok {
		t.Fatalf("tools missing or wrong type: %v", result["tools"])
	}
	want := map[string]bool{"tunnel-connect": true, "tunnel-disconnect": true, "tunnel-status": true}
	for _, raw := range tools {
		tool := raw.(map[string]any)
		name := tool["name"].(string)
		delete(want, name)
	}
	if len(want) != 0 {
		t.Errorf("missing tools: %v", want)
	}
}

func TestTunnelStatusEmpty(t *testing.T) {
	srv := New("test-version")
	resp := exchange(t, srv, map[string]any{
		"jsonrpc": "2.0",
		"id":      3,
		"method":  "tools/call",
		"params": map[string]any{
			"name":      "tunnel-status",
			"arguments": map[string]any{},
		},
	})

	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	result := resp["result"].(map[string]any)
	content := result["content"].([]any)
	if len(content) == 0 {
		t.Fatal("empty content")
	}
	item := content[0].(map[string]any)
	var payload map[string]any
	if err := json.Unmarshal([]byte(item["text"].(string)), &payload); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if payload["count"].(float64) != 0 {
		t.Errorf("count = %v, want 0", payload["count"])
	}
}

func TestTunnelDisconnectUnknown(t *testing.T) {
	srv := New("test-version")
	resp := exchange(t, srv, map[string]any{
		"jsonrpc": "2.0",
		"id":      4,
		"method":  "tools/call",
		"params": map[string]any{
			"name":      "tunnel-disconnect",
			"arguments": map[string]any{"tunnelId": "nonexistent"},
		},
	})

	if resp["error"] != nil {
		t.Fatalf("unexpected rpc error: %v", resp["error"])
	}
	result := result(t, resp)
	if result["isError"] != true {
		t.Errorf("expected isError=true for unknown tunnelId")
	}
}

func TestTunnelConnectMissingRelayURL(t *testing.T) {
	srv := New("test-version")
	resp := exchange(t, srv, map[string]any{
		"jsonrpc": "2.0",
		"id":      5,
		"method":  "tools/call",
		"params": map[string]any{
			"name":      "tunnel-connect",
			"arguments": map[string]any{},
		},
	})

	if resp["error"] != nil {
		t.Fatalf("unexpected rpc error: %v", resp["error"])
	}
	r := result(t, resp)
	if r["isError"] != true {
		t.Errorf("expected isError=true when relayUrl is missing")
	}
}

func TestTunnelConnectBadOrigin(t *testing.T) {
	srv := New("test-version")
	resp := exchange(t, srv, map[string]any{
		"jsonrpc": "2.0",
		"id":      6,
		"method":  "tools/call",
		"params": map[string]any{
			"name": "tunnel-connect",
			"arguments": map[string]any{
				"relayUrl":       "ws://relay.example.com/tunnel",
				"allowedOrigins": []string{"evil.com:80"}, // non-loopback
			},
		},
	})

	if resp["error"] != nil {
		t.Fatalf("unexpected rpc error: %v", resp["error"])
	}
	r := result(t, resp)
	if r["isError"] != true {
		t.Errorf("expected isError=true for non-loopback origin")
	}
}

func TestUnknownMethod(t *testing.T) {
	srv := New("test-version")
	resp := exchange(t, srv, map[string]any{
		"jsonrpc": "2.0",
		"id":      7,
		"method":  "unknown/method",
	})
	if resp["error"] == nil {
		t.Errorf("expected rpc error for unknown method")
	}
}

func TestNotificationIgnored(t *testing.T) {
	// A notification (no id) followed by tools/list; only the tools/list response is emitted.
	notification := `{"jsonrpc":"2.0","method":"notifications/initialized"}` + "\n"
	request := `{"jsonrpc":"2.0","id":8,"method":"tools/list"}` + "\n"

	var out bytes.Buffer
	srv := New("test-version")
	if err := srv.Serve(context.Background(), strings.NewReader(notification+request), &out); err != nil {
		t.Fatalf("Serve: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 response line, got %d: %v", len(lines), lines)
	}
}

// result extracts the MCP tool-call result payload from an rpcResponse.
func result(t *testing.T, resp map[string]any) map[string]any {
	t.Helper()
	r, ok := resp["result"].(map[string]any)
	if !ok {
		t.Fatalf("result missing: %v", resp)
	}
	return r
}
