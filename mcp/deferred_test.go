package mcp

import (
	"os"
	"testing"
)

func TestMCPToolsDeferredEnvParsing(t *testing.T) {
	orig, had := os.LookupEnv("SCORP_MCP_DEFERRED")
	t.Cleanup(func() {
		if had {
			os.Setenv("SCORP_MCP_DEFERRED", orig)
		} else {
			os.Unsetenv("SCORP_MCP_DEFERRED")
		}
	})

	for _, tc := range []struct {
		val  string
		want bool
	}{
		{"", true},      // default: deferred (P2.9)
		{"on", true},    // explicit on
		{"1", true},     // truthy
		{"anything", true},
		{"off", false},  // escape hatch
		{"false", false},
		{"0", false},
		{"no", false},
		{" OFF ", false},
	} {
		if tc.val == "" {
			os.Unsetenv("SCORP_MCP_DEFERRED")
		} else {
			os.Setenv("SCORP_MCP_DEFERRED", tc.val)
		}
		if got := MCPToolsDeferred(); got != tc.want {
			t.Errorf("SCORP_MCP_DEFERRED=%q: MCPToolsDeferred()=%v, want %v", tc.val, got, tc.want)
		}
	}
}
