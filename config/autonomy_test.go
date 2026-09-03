package config

import (
	"testing"
)

func TestAutonomyLevels(t *testing.T) {
	// Default is supervised
	SetAutonomyLevel("supervised")
	if GetAutonomyLevel() != AutonomySupervised {
		t.Errorf("Expected default supervised, got %s", GetAutonomyLevel())
	}

	// ReadOnly mode blocks write tools
	SetAutonomyLevel("readonly")
	if GetAutonomyLevel() != AutonomyReadOnly {
		t.Errorf("Expected readonly, got %s", GetAutonomyLevel())
	}

	allowed, _ := IsToolAllowed("read_file")
	if !allowed {
		t.Errorf("Expected read_file to be allowed in ReadOnly mode")
	}

	writeAllowed, reason := IsToolAllowed("write_file")
	if writeAllowed {
		t.Errorf("Expected write_file to be blocked in ReadOnly mode")
	}
	if reason == "" {
		t.Errorf("Expected rejection reason for blocked tool")
	}

	// YOLO mode allows all tools
	SetAutonomyLevel("yolo")
	if GetAutonomyLevel() != AutonomyYOLO {
		t.Errorf("Expected yolo, got %s", GetAutonomyLevel())
	}
	yoloAllowed, _ := IsToolAllowed("write_file")
	if !yoloAllowed {
		t.Errorf("Expected write_file to be allowed in YOLO mode")
	}

	// Sandbox restricted path check
	restricted, _ := IsPathRestricted("/etc/shadow")
	if !restricted {
		t.Errorf("Expected /etc/shadow to be restricted")
	}
	normalPathRestricted, _ := IsPathRestricted("/home/user/project/main.go")
	if normalPathRestricted {
		t.Errorf("Expected normal path not to be restricted")
	}

	// Reset to supervised
	SetAutonomyLevel("supervised")
}
