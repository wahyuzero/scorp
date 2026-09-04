package tools

import (
	"strings"
	"testing"
)

func TestExecuteTermuxAPI_Simulation(t *testing.T) {
	// Test simulation behavior when running outside Android Termux
	res, ok := ExecuteTermuxAPI(map[string]interface{}{
		"action":  "notification",
		"title":   "Test",
		"content": "Hello",
	})
	if !ok {
		t.Fatalf("Expected simulation to succeed, got: %s", res)
	}
	if !strings.Contains(res, "Simulated Termux notification") {
		t.Errorf("Unexpected simulation response: %s", res)
	}

	res2, ok2 := ExecuteTermuxAPI(map[string]interface{}{
		"action": "wake_lock",
	})
	if !ok2 {
		t.Fatalf("Expected wake_lock simulation to succeed, got: %s", res2)
	}
	if !strings.Contains(res2, "Simulated Termux wake_lock") {
		t.Errorf("Unexpected simulation response: %s", res2)
	}
}
