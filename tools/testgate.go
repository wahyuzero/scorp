package tools

import (
	"path/filepath"
	"strings"
	"time"
)

// ──────────────────────────────────────────────
// Test-Integrity Gate (P0.3)
//
// Community evidence (docs/RESEARCH_AI_AGENT_FEATURES_2026.md): in a 2026
// audit, 19/45 "all tests pass" claims were fabricated; agents weaken tests
// to make suites green. Contract: if the current task touched test files or
// CI config, complete_task is rejected until a SUCCESSFUL test-suite shell
// receipt exists AFTER the last such touch.
//
// Evidence source: receipts.json (every model-executed tool call is recorded
// with meta facts), so the gate works across daemon restarts and survives
// confirm-resume flows (which record receipts too).
// ──────────────────────────────────────────────

// testGateBoundary is the wall-clock start of the current task. Set when a
// fresh RunAgentSessionLoop begins; resume loops (confirmation continuation)
// deliberately do NOT reset it, so edits made before a user confirmation still
// count toward the gate.
var testGateBoundary time.Time

// MarkTaskBoundary starts a fresh test-integrity window for a new task.
func MarkTaskBoundary() {
	testGateBoundary = time.Now()
}

// testPathMarkers match test files and CI configs across common stacks.
var testFileSuffixes = []string{
	"_test.go", "_test.py", "_test.rs", "_test.rb", "_test.jsx", "_test.ts",
	"_test.tsx", "_test.mjs", "_test.cjs",
	".test.js", ".test.ts", ".test.jsx", ".test.tsx", ".test.mjs", ".test.cjs",
	".spec.js", ".spec.ts", ".spec.jsx", ".spec.tsx",
	"test.java", "tests.java", "spec.rb", "_test.dart",
}

var testFileNames = []string{
	"conftest.py", "jenkinsfile", ".gitlab-ci.yml", "azure-pipelines.yml",
	"pytest.ini", "tox.ini", "jest.config.js", "jest.config.ts",
	"vitest.config.ts", "vitest.config.js", "karma.conf.js", ".phpunit.xml",
	"phpunit.xml", ".circleci",
}

var testPathSegments = []string{
	"/tests/", "/test/", "/__tests__/", "/spec/", "/.github/workflows/", "/.circleci/",
}

// IsTestRelatedPath reports whether a file path is a test source or CI config.
func IsTestRelatedPath(p string) bool {
	if p == "" {
		return false
	}
	norm := "/" + strings.ToLower(filepath.ToSlash(strings.TrimSpace(p))) + "/"
	for _, seg := range testPathSegments {
		if strings.Contains(norm, seg) {
			return true
		}
	}
	base := strings.ToLower(filepath.Base(p))
	for _, s := range testFileSuffixes {
		if strings.HasSuffix(base, s) {
			return true
		}
	}
	for _, n := range testFileNames {
		if base == n {
			return true
		}
	}
	return false
}

// testRunCommands detect a test-suite execution in a shell command. A green
// receipt is only recorded when the shell tool also reports success.
var testRunCommands = []string{
	"go test", "npm test", "npm run test", "yarn test", "pnpm test",
	"pytest", "cargo test", "make test", "mvn test", "gradle test",
	"rake test", "phpunit", "jest", "vitest", "dart test", "flutter test",
	"rspec", "dotnet test", "mix test", "go vet ./",
}

// IsTestRunCommand reports whether a shell command executes a test suite.
func IsTestRunCommand(cmd string) bool {
	lower := strings.ToLower(cmd)
	for _, marker := range testRunCommands {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

// shellWriteMarkers approximate "this command modifies files". Reading a test
// file (cat, grep) must not trigger the gate; rewriting one must.
var shellWriteMarkers = []string{
	">", "sed -i", "tee ", "rm ", "mv ", "cp ", "truncate", "patch ",
	"unlink ", "rsync ", "git checkout ", "git restore ",
}

// shellTouchesTestFile is a conservative heuristic for shell-based test-file
// edits: some token in the command must look like a test/CI path AND the
// command must contain a write marker. Reading a test file (cat, grep) does
// not trigger the gate. False positives only force an extra test run — the
// safe bias.
func shellTouchesTestFile(cmd string) bool {
	lower := strings.ToLower(cmd)
	hasTestToken := false
	for _, tok := range strings.FieldsFunc(lower, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\'' || r == '"' || r == ';' ||
			r == '|' || r == '&' || r == '(' || r == ')' || r == '<' || r == '>'
	}) {
		if IsTestRelatedPath(tok) {
			hasTestToken = true
			break
		}
	}
	if !hasTestToken {
		return false
	}
	for _, m := range shellWriteMarkers {
		if strings.Contains(lower, m) {
			return true
		}
	}
	return false
}

// isWriteToolName classifies structured tools by name — any tool whose name
// implies mutation counts when its meta path points at a test file.
func isWriteToolName(name string) bool {
	lower := strings.ToLower(name)
	for _, marker := range []string{
		"write", "edit", "patch", "creat", "updat", "append", "move",
		"rename", "copy", "delet", "remov", "truncat", "import", "generat", "apply",
	} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

// TestIntegrityStatus scans the receipt window of the current task:
// touched = test/CI files were modified; green = a successful test-suite run
// was recorded AFTER the last modification. Gate closes when touched && !green.
func TestIntegrityStatus() (touched, green bool) {
	var lastTouch, lastGreen time.Time
	for _, r := range GetRecentReceipts() {
		if !testGateBoundary.IsZero() && r.Timestamp.Before(testGateBoundary) {
			continue // receipt predates the current task
		}
		if r.Tool == "shell" {
			cmd := r.Meta["cmd"]
			if cmd == "" {
				continue
			}
			isTestRun := IsTestRunCommand(cmd)
			if isTestRun && r.Success {
				lastGreen = r.Timestamp
			}
			// A command that runs the suite counts as green, not as a touch —
			// otherwise `go test ./tests/... -v > /tmp/log` would block itself.
			if !isTestRun && shellTouchesTestFile(cmd) {
				lastTouch = r.Timestamp
			}
			continue
		}
		if p := r.Meta["path"]; p != "" && IsTestRelatedPath(p) && isWriteToolName(r.Tool) {
			lastTouch = r.Timestamp
		}
	}
	return !lastTouch.IsZero(), lastGreen.After(lastTouch)
}
