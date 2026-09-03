package mcp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRemoteMCPServer_HTTP(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req jsonRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		switch req.Method {
		case "initialize":
			resp := jsonRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: json.RawMessage(`{
					"protocolVersion": "2024-11-05",
					"capabilities": {"tools": {}},
					"serverInfo": {"name": "mock-mcp", "version": "1.0.0"}
				}`),
			}
			json.NewEncoder(w).Encode(resp)
		case "tools/list":
			resp := jsonRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: json.RawMessage(`{
					"tools": [
						{
							"name": "remote_ping",
							"description": "Ping remote service",
							"inputSchema": {
								"type": "object",
								"properties": {"msg": {"type": "string"}}
							}
						}
					]
				}`),
			}
			json.NewEncoder(w).Encode(resp)
		case "tools/call":
			resp := jsonRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: json.RawMessage(`{
					"content": [{"type": "text", "text": "pong from remote"}]
				}`),
			}
			json.NewEncoder(w).Encode(resp)
		default:
			fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"result":{}}`, req.ID)
		}
	}))
	defer ts.Close()

	cfg := MCPServerConfig{
		URL: ts.URL,
	}

	srv, err := startSSEServer("test_mock", cfg)
	if err != nil {
		t.Fatalf("startSSEServer failed: %v", err)
	}
	defer srv.Close()

	if len(srv.tools) != 1 {
		t.Fatalf("Expected 1 tool discovered, got %d", len(srv.tools))
	}
	if srv.tools[0].Name != "remote_ping" {
		t.Errorf("Expected tool name 'remote_ping', got '%s'", srv.tools[0].Name)
	}

	// Test calling the remote tool
	res, err := srv.CallTool("remote_ping", map[string]interface{}{"msg": "hello"})
	if err != nil {
		t.Fatalf("callTool failed: %v", err)
	}
	if res != "pong from remote" {
		t.Errorf("Expected 'pong from remote', got '%s'", res)
	}
}
