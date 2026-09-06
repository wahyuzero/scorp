package registry

import (
	"os"
	"testing"
)

func TestDynamicToolTTL(t *testing.T) {
	os.Setenv("SCORP_DYNAMIC_TOOLS", "true")
	defer os.Unsetenv("SCORP_DYNAMIC_TOOLS")

	ResetDynamicTools()

	// Core tool is always active
	coreDef := ToolDef{Name: "read_file", Native: true}
	if !IsToolActive(coreDef) {
		t.Errorf("Expected core tool read_file to be active")
	}

	// Specialized tool is inactive by default in dynamic mode
	specDef := ToolDef{Name: "sql", Native: true}
	if IsToolActive(specDef) {
		t.Errorf("Expected specialized tool sql to be inactive initially in dynamic mode")
	}

	// Activate with TTL=2
	ActivateToolWithTTL("sql", 2)
	if !IsToolActive(specDef) {
		t.Errorf("Expected specialized tool sql to be active after TTL activation")
	}

	// Turn 1 tick -> TTL becomes 1
	TickToolTTL()
	if !IsToolActive(specDef) {
		t.Errorf("Expected specialized tool sql to still be active after 1 tick (TTL=1)")
	}

	// Turn 2 tick -> TTL expires (0)
	TickToolTTL()
	if IsToolActive(specDef) {
		t.Errorf("Expected specialized tool sql to be expired and inactive after 2 ticks")
	}
}

func clearDynamicEnv(t *testing.T, dynamic string) {
	t.Helper()
	orig, had := os.LookupEnv("SCORP_DYNAMIC_TOOLS")
	if dynamic == "" {
		os.Unsetenv("SCORP_DYNAMIC_TOOLS")
	} else {
		os.Setenv("SCORP_DYNAMIC_TOOLS", dynamic)
	}
	ResetDynamicTools()
	ResetNativeToolCache()
	t.Cleanup(func() {
		if had {
			os.Setenv("SCORP_DYNAMIC_TOOLS", orig)
		} else {
			os.Unsetenv("SCORP_DYNAMIC_TOOLS")
		}
		ResetDynamicTools()
		ResetNativeToolCache()
	})
}

func registerTempTool(t *testing.T, name string, deferred bool) {
	t.Helper()
	RegisterTool(ToolDef{
		Name:        name,
		Description: "temp test tool",
		Category:    "test",
		Native:      true,
		Deferred:    deferred,
		Execute:     func(map[string]interface{}, int64) (string, bool) { return "ok", true },
	})
	t.Cleanup(func() {
		UnregisterTool(name)
		ResetNativeToolCache()
	})
}

// TestIsToolActiveStaticModeWithTTL pins the P2.9 contract: in the default
// (static) mode a deferred tool is withheld from the schema until discovered
// via tool_search / invoked via tool_call, then lives for its TTL window and
// returns to deferred.
func TestIsToolActiveStaticModeWithTTL(t *testing.T) {
	clearDynamicEnv(t, "")
	registerTempTool(t, "zz_temp_deferred", true)

	def, ok := GetTool("zz_temp_deferred")
	if !ok {
		t.Fatal("temp tool missing")
	}
	if IsToolActive(def) {
		t.Fatal("deferred tool must be inactive in static mode before discovery")
	}

	ActivateToolWithTTL("zz_temp_deferred", 2)
	if !IsToolActive(def) {
		t.Fatal("activated deferred tool must be active")
	}

	// schema must include it while active
	schema := GenerateNativeToolsSchema()
	found := false
	for _, s := range schema {
		if s.Function.Name == "zz_temp_deferred" {
			found = true
		}
	}
	if !found {
		t.Fatal("activated tool must appear in the native schema")
	}

	TickToolTTL()
	TickToolTTL()
	if IsToolActive(def) {
		t.Fatal("tool must return to deferred after TTL expiry")
	}
}

func TestIsToolActiveNonDeferredAlwaysActiveInStaticMode(t *testing.T) {
	clearDynamicEnv(t, "")
	registerTempTool(t, "zz_temp_native", false)

	def, _ := GetTool("zz_temp_native")
	if !IsToolActive(def) {
		t.Fatal("non-deferred native tool must stay active in static mode")
	}
}

func TestGenerateNativeToolsSchemaExcludesDeferred(t *testing.T) {
	clearDynamicEnv(t, "")
	registerTempTool(t, "zz_temp_deferred", true)
	registerTempTool(t, "zz_temp_native", false)

	schema := GenerateNativeToolsSchema()
	var deferredSeen, nativeSeen bool
	for _, s := range schema {
		switch s.Function.Name {
		case "zz_temp_deferred":
			deferredSeen = true
		case "zz_temp_native":
			nativeSeen = true
		}
	}
	if deferredSeen {
		t.Fatal("deferred tool must NOT be injected into the native schema")
	}
	if !nativeSeen {
		t.Fatal("non-deferred tool must be injected into the native schema")
	}
}
