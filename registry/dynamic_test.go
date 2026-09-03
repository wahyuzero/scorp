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
