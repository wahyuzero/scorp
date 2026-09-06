package eval

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// ──────────────────────────────────────────────
// scorp eval — private evaluation arena (P4.15)
//
// SWE-bench was retired; the 2026 consensus is "build your own private
// arena from real work". This codifies the verify26 method that guarded
// scorp's own releases:
//
//   - core cases: deterministic in-process checks of every safety/persistence
//     layer — no model, no network, fast enough for every deploy
//   - live cases (opt-in --live): real agent tasks run as subprocesses
//     (`scorp -p <task> -s scorp-eval-<case>`), verified by INDEPENDENT
//     artifact checkers — the agent's claims are never trusted
//   - metrics: pass rate per category + tokens per live task (from
//     model_usage.json), and exit code 1 on any failure so CI/deploy
//     scripts can gate on it
// ──────────────────────────────────────────────

// Case is one evaluation unit. Run returns an error on failure; the message
// should state WHAT was independently verified, not just what broke.
type Case struct {
	Name     string
	Category string
	Run      func() error
	Tokens   func() int // optional: model tokens consumed (live cases)
}

type caseResult struct {
	c        Case
	pass     bool
	err      error
	duration time.Duration
	tokens   int // live cases only: model tokens consumed by the task
}

// Run executes the suite and returns the process exit code.
func Run(args []string) int {
	live := false
	filter := ""
	for _, a := range args {
		switch {
		case a == "--live":
			live = true
		case strings.HasPrefix(a, "--filter="):
			filter = strings.TrimPrefix(a, "--filter=")
		}
	}

	cases := CoreCases()
	if live {
		cases = append(cases, LiveCases()...)
	}

	var selected []Case
	for _, c := range cases {
		if filter == "" || strings.EqualFold(c.Category, filter) || strings.EqualFold(c.Name, filter) {
			selected = append(selected, c)
		}
	}
	if len(selected) == 0 {
		fmt.Println("🧪 scorp eval: no cases match filter", filter)
		return 1
	}

	fmt.Printf("🧪 scorp eval — %d case(s)%s\n\n", len(selected), liveLabel(live))

	results := make([]caseResult, 0, len(selected))
	for i, c := range selected {
		fmt.Printf("[%2d/%d] %-12s %-42s ", i+1, len(selected), c.Category, c.Name)
		start := time.Now()
		res := caseResult{c: c}
		func() {
			defer func() {
				if r := recover(); r != nil {
					res.err = fmt.Errorf("panic: %v", r)
				}
			}()
			res.err = c.Run()
		}()
		res.duration = time.Since(start)
		res.pass = res.err == nil
		if c.Tokens != nil {
			res.tokens = c.Tokens()
		}
		if res.pass {
			fmt.Printf("✅ PASS (%s)\n", res.duration.Round(time.Millisecond))
		} else {
			fmt.Printf("❌ FAIL — %v\n", res.err)
		}
		results = append(results, res)
	}

	fmt.Println()
	return report(results)
}

func report(results []caseResult) int {
	passed := 0
	byCat := map[string][2]int{} // cat -> {passed, total}
	catOrder := []string{}
	var liveTokens, liveCount int

	for _, r := range results {
		cat := r.c.Category
		if _, ok := byCat[cat]; !ok {
			catOrder = append(catOrder, cat)
		}
		slot := byCat[cat]
		if r.pass {
			passed++
			slot[0]++
		}
		slot[1]++
		byCat[cat] = slot
		if strings.HasPrefix(cat, "live") {
			liveTokens += r.tokens
			liveCount++
		}
	}

	for _, cat := range catOrder {
		s := byCat[cat]
		fmt.Printf("  %-14s %d/%d passed\n", cat, s[0], s[1])
	}
	pct := 0
	if len(results) > 0 {
		pct = passed * 100 / len(results)
	}
	fmt.Printf("\nTOTAL: %d/%d passed (%d%%)", passed, len(results), pct)
	if liveCount > 0 {
		avg := 0
		if liveCount > 0 {
			avg = liveTokens / liveCount
		}
		fmt.Printf(" · tokens/task ≈ %d (%d live tasks)", avg, liveCount)
	}
	fmt.Println()

	if passed != len(results) {
		fmt.Println("❌ eval gate FAILED — fix failures before deploying.")
		return 1
	}
	fmt.Println("✅ eval gate PASSED.")
	return 0
}

func liveLabel(live bool) string {
	if live {
		return " (incl. LIVE model tasks)"
	}
	return " (core only — add --live for model tasks)"
}

// evalSandboxDir creates a unique working directory for one live case.
func evalSandboxDir(tag string) (string, error) {
	return os.MkdirTemp("", "scorp-eval-"+tag+"-")
}
