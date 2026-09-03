package agent

import (
	"strings"
	"testing"
)

func TestGetRepoMap(t *testing.T) {
	InvalidateRepoMap()
	rm := GetRepoMap()
	if rm == "" {
		t.Fatalf("Expected non-empty repo map")
	}

	if !strings.HasPrefix(rm, "```") || !strings.HasSuffix(rm, "```") {
		t.Errorf("Expected repo map to be wrapped in code block")
	}

	// Should contain files in the working directory
	if !strings.Contains(rm, "prompt.go") && !strings.Contains(rm, "loop.go") {
		t.Errorf("Expected repo map to contain files in dir, got: %s", rm)
	}

	// Should not contain .git or graphify-out
	if strings.Contains(rm, "graphify-out") {
		t.Errorf("Repo map should ignore graphify-out")
	}
}
