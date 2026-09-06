package tools

import (
	"strings"
	"testing"

	"scorp-agent/registry"
)

// registerSearchTarget registers one deferred tool for discovery tests.
func registerSearchTarget(t *testing.T) {
	t.Helper()
	registry.RegisterTool(registry.ToolDef{
		Name:        "zz_mcp_testwidget_send",
		Description: "[MCP:testwidgets] Send a widget to a test recipient",
		Category:    "mcp",
		Native:      true,
		Deferred:    true,
		Execute:     func(map[string]interface{}, int64) (string, bool) { return "widget sent", true },
		Arguments: map[string]registry.ArgDef{
			"recipient": {Type: "string", Description: "who receives it", Required: true},
		},
	})
	t.Cleanup(func() {
		registry.UnregisterTool("zz_mcp_testwidget_send")
		registry.ResetNativeToolCache()
	})
}

func deferredToolActive() bool {
	def, ok := registry.GetTool("zz_mcp_testwidget_send")
	if !ok {
		return false
	}
	return registry.IsToolActive(def)
}

// TestToolSearchActivatesDeferredToolInStaticMode pins the P2.9 discovery
// flow: tool_search must surface deferred MCP tools and inject them into the
// native schema for a TTL window — WITHOUT requiring SCORP_DYNAMIC_TOOLS.
func TestToolSearchActivatesDeferredToolInStaticMode(t *testing.T) {
	registerSearchTarget(t)
	t.Setenv("SCORP_DYNAMIC_TOOLS", "")
	osUnsetDynamic(t)

	if deferredToolActive() {
		t.Fatal("deferred tool must start inactive")
	}

	out, ok := ExecuteToolSearch(map[string]interface{}{"query": "widget"}, 0)
	if !ok {
		t.Fatalf("tool_search failed: %q", out)
	}
	if !strings.Contains(out, "zz_mcp_testwidget_send") {
		t.Fatalf("search must surface the deferred tool, got %q", out)
	}
	if !deferredToolActive() {
		t.Fatal("discovered deferred tool must be activated into the native schema")
	}

	// tool_call also activates (and executes through the deferred tool)
	out, ok = ExecuteToolCall(map[string]interface{}{
		"name":      "zz_mcp_testwidget_send",
		"arguments": map[string]interface{}{"recipient": "tester"},
	}, 0)
	if !ok || !strings.Contains(out, "widget sent") {
		t.Fatalf("tool_call through deferred tool failed: ok=%v out=%q", ok, out)
	}
	if !deferredToolActive() {
		t.Fatal("tool_call must keep the deferred tool activated")
	}
}

func osUnsetDynamic(t *testing.T) {
	t.Helper()
	t.Setenv("SCORP_DYNAMIC_TOOLS", "0") // not "true"/"1" → dynamic mode OFF
}
