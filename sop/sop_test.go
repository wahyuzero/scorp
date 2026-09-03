package sop

import (
	"testing"
)

func TestSOPLifecycle(t *testing.T) {
	InitDefaultSOPs()

	sops := ListSOPs()
	if len(sops) == 0 {
		t.Fatalf("Expected default SOPs to be initialized")
	}

	foundHealth := false
	for _, s := range sops {
		if s.Name == "health_audit" {
			foundHealth = true
			break
		}
	}
	if !foundHealth {
		t.Errorf("Expected 'health_audit' SOP in list")
	}

	// Test GetSOP
	s, err := GetSOP("health_audit")
	if err != nil {
		t.Fatalf("GetSOP failed: %v", err)
	}
	if len(s.Steps) == 0 {
		t.Errorf("Expected steps in health_audit SOP")
	}

	// Test Save custom SOP
	custom := SOP{
		Name:        "custom_test",
		Description: "Custom test procedure",
		Steps:       []string{"step 1", "step 2"},
		Prompt:      "Execute custom test",
	}
	if err := SaveSOP(custom); err != nil {
		t.Fatalf("SaveSOP failed: %v", err)
	}

	retrieved, err := GetSOP("custom_test")
	if err != nil {
		t.Fatalf("GetSOP custom failed: %v", err)
	}
	if retrieved.Description != "Custom test procedure" {
		t.Errorf("Unexpected description: %s", retrieved.Description)
	}
}
