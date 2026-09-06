package tools

import (
	"testing"
	"time"
)

// Receipt state is manipulated directly (same package) so tests never touch
// the real ~/.scorp/receipts.json that RecordToolReceipt persists to.
func setTestReceipts(t *testing.T, receipts []ToolReceipt) {
	t.Helper()
	recentReceiptsMu.Lock()
	recentReceipts = receipts
	receiptsLoaded = true
	recentReceiptsMu.Unlock()
	t.Cleanup(func() {
		recentReceiptsMu.Lock()
		recentReceipts = nil
		recentReceipts = recentReceipts[:0]
		receiptsLoaded = false
		recentReceiptsMu.Unlock()
		testGateBoundary = time.Time{}
	})
}

func mkReceipt(tool string, minutesAgo int, success bool, meta map[string]string) ToolReceipt {
	return ToolReceipt{
		ReceiptID: tool,
		Tool:      tool,
		Success:   success,
		Timestamp: time.Now().Add(-time.Duration(minutesAgo) * time.Minute),
		Meta:      meta,
	}
}

func TestIsTestRelatedPath(t *testing.T) {
	yes := []string{
		"agent/loop_test.go", "internal/foo_test.go",
		"tests/test_login.py", "tests/conftest.py", "conftest.py",
		"src/app.component.spec.ts", "src/util.test.js",
		"src/__tests__/store.tsx", "ci/.github/workflows/ci.yml",
		".gitlab-ci.yml", "Jenkinsfile", "phpunit.xml",
		"crate/src/lib/tests/mod.rs", "app/tests/FooTest.java",
	}
	no := []string{
		"agent/loop.go", "README.md", "main.py", "utils.js",
		"latest_test_report.pdf", "attest.json", "scripts/testing_helper.sh",
		"", "/tmp/notes.txt",
	}
	for _, p := range yes {
		if !IsTestRelatedPath(p) {
			t.Errorf("IsTestRelatedPath(%q) = false, want true", p)
		}
	}
	for _, p := range no {
		if IsTestRelatedPath(p) {
			t.Errorf("IsTestRelatedPath(%q) = true, want false", p)
		}
	}
}

func TestIsTestRunCommand(t *testing.T) {
	yes := []string{"go test ./...", "npm test -- --watchAll=false", "python -m pytest -q", "cargo test --release", "make test"}
	no := []string{"go build ./...", "ls -la", "cat Makefile", "echo done"}
	for _, c := range yes {
		if !IsTestRunCommand(c) {
			t.Errorf("IsTestRunCommand(%q) = false, want true", c)
		}
	}
	for _, c := range no {
		if IsTestRunCommand(c) {
			t.Errorf("IsTestRunCommand(%q) = true, want false", c)
		}
	}
}

func TestTestIntegrityStatus_NoTouches(t *testing.T) {
	setTestReceipts(t, []ToolReceipt{
		mkReceipt("shell", 5, true, map[string]string{"cmd": "go build ./..."}),
		mkReceipt("read_file", 4, true, map[string]string{"path": "agent/loop_test.go"}),
	})
	touched, green := TestIntegrityStatus()
	if touched {
		t.Fatal("reading test files must not count as a touch")
	}
	if green {
		t.Fatal("no test run recorded — green must be false")
	}
}

func TestTestIntegrityStatus_TouchBlocksUntilGreenRun(t *testing.T) {
	// write_file touched a test file 10 minutes ago; no test run after → closed
	setTestReceipts(t, []ToolReceipt{
		mkReceipt("shell", 12, true, map[string]string{"cmd": "go test ./..."}),
		mkReceipt("write_file", 10, true, map[string]string{"path": "tools/exec_gate_test.go"}),
	})
	touched, green := TestIntegrityStatus()
	if !touched || green {
		t.Fatalf("gate must be closed (touched=%v green=%v), want touched=true green=false", touched, green)
	}

	// a green run AFTER the touch reopens it
	setTestReceipts(t, []ToolReceipt{
		mkReceipt("write_file", 10, true, map[string]string{"path": "tools/exec_gate_test.go"}),
		mkReceipt("shell", 2, true, map[string]string{"cmd": "go test ./tools/..."}),
	})
	touched, green = TestIntegrityStatus()
	if !touched || !green {
		t.Fatalf("gate must be open after a later green run (touched=%v green=%v)", touched, green)
	}
}

func TestTestIntegrityStatus_GreenRunMustBeAfterEdit(t *testing.T) {
	// test run happened BEFORE the edit — that green says nothing about the
	// modified suite → gate stays closed
	setTestReceipts(t, []ToolReceipt{
		mkReceipt("shell", 12, true, map[string]string{"cmd": "go test ./..."}),
		mkReceipt("shell", 10, true, map[string]string{"cmd": "sed -i 's/t.Skip(t)//' agent/loop_test.go"}),
	})
	touched, green := TestIntegrityStatus()
	if !touched || green {
		t.Fatalf("edit after green run must keep the gate closed (touched=%v green=%v)", touched, green)
	}
}

func TestTestIntegrityStatus_FailingSuiteDoesNotCount(t *testing.T) {
	setTestReceipts(t, []ToolReceipt{
		mkReceipt("shell", 10, true, map[string]string{"path": "tools/x_test.go", "cmd": "echo rewrite > tools/x_test.go"}),
		mkReceipt("shell", 2, false, map[string]string{"cmd": "go test ./..."}),
	})
	touched, green := TestIntegrityStatus()
	if !touched || green {
		t.Fatalf("a FAILED test run must not open the gate (touched=%v green=%v)", touched, green)
	}
}

func TestTestIntegrityStatus_TestRunWithRedirectIsNotATouch(t *testing.T) {
	// `go test ./tests/... -v > /tmp/log` mentions a test dir and a redirect,
	// but it IS the green run — it must not block itself.
	setTestReceipts(t, []ToolReceipt{
		mkReceipt("shell", 3, true, map[string]string{"cmd": "go test ./tests/... -v > /tmp/log 2>&1"}),
	})
	touched, green := TestIntegrityStatus()
	if touched {
		t.Fatal("a test-run command must not register as a test-file touch")
	}
	if !green {
		t.Fatal("successful test run must register as green")
	}
}

func TestTestIntegrityStatus_ShellEditDetected(t *testing.T) {
	setTestReceipts(t, []ToolReceipt{
		mkReceipt("shell", 4, true, map[string]string{"cmd": "sed -i 's/x/y/' agent/loop_test.go && echo done"}),
	})
	touched, _ := TestIntegrityStatus()
	if !touched {
		t.Fatal("shell-based test-file edit must register as a touch")
	}
}

func TestTestIntegrityStatus_TaskBoundaryExcludesOldReceipts(t *testing.T) {
	setTestReceipts(t, []ToolReceipt{
		mkReceipt("write_file", 30, true, map[string]string{"path": "old_task_test.go"}),
	})
	MarkTaskBoundary()
	time.Sleep(10 * time.Millisecond)
	touched, _ := TestIntegrityStatus()
	if touched {
		t.Fatal("receipts from before the task boundary must be ignored")
	}
}
