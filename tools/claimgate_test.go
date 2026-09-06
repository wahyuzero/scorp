package tools

import (
	"testing"
)

func TestLooksLikeTestPassClaim(t *testing.T) {
	claims := []string{
		"All tests pass.",
		"all tests passed",
		"The test suite is green.",
		"Tests are passing now.",
		"all checks passed",
		"Semua test lulus.",
		"tests lulus semua",
		"100% tests lulus",
		"Everything is green.",
		"the build is green",
	}
	for _, c := range claims {
		if !LooksLikeTestPassClaim(c) {
			t.Errorf("expected claim detected: %q", c)
		}
	}
	nonClaims := []string{
		"Running the tests now.",
		"The test suite is red — 3 failures.",
		"tests fail and I need help",
		"I will write the test file next.",
		"Fixed a typo in main.go.",
	}
	for _, c := range nonClaims {
		if LooksLikeTestPassClaim(c) {
			t.Errorf("false positive: %q", c)
		}
	}
}

func TestHasGreenTestRun(t *testing.T) {
	setTestReceipts(t, nil)
	MarkTaskBoundary()

	if HasGreenTestRun() {
		t.Fatal("no receipts yet — must be false")
	}

	// Failed test run does not count.
	RecordToolReceipt("shell", map[string]interface{}{"command": "go test ./..."}, "FAIL", false)
	if HasGreenTestRun() {
		t.Fatal("failed test run must not count as green")
	}

	// Successful test run counts.
	RecordToolReceipt("shell", map[string]interface{}{"command": "go test ./... "}, "ok", true)
	if !HasGreenTestRun() {
		t.Fatal("successful go test receipt must satisfy the claim gate")
	}

	// Receipts from BEFORE the task boundary do not count.
	MarkTaskBoundary()
	if HasGreenTestRun() {
		t.Fatal("green run predating the task boundary must not count")
	}
}
