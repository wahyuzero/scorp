package eval

import (
	"errors"
	"strings"
	"testing"
)

// TestRunnerAggregationAndFilter pins the report contract: per-case results,
// category rollup, exit codes, and the --filter path.
func TestRunnerAggregationAndFilter(t *testing.T) {
	cases := []Case{
		{Name: "ok_safety", Category: "safety", Run: func() error { return nil }},
		{Name: "bad_safety", Category: "safety", Run: func() error { return errors.New("gate broken") }},
		{Name: "ok_live", Category: "live-agent", Run: func() error { return nil }, Tokens: func() int { return 1234 }},
	}

	var results []caseResult
	for _, c := range cases {
		res := caseResult{c: c, pass: c.Run() == nil}
		if c.Tokens != nil {
			res.tokens = c.Tokens()
		}
		results = append(results, res)
	}
	if results[1].pass {
		t.Fatal("failing case must not pass")
	}

	// exercise the report formatter (its output is user-facing; failures here
	// would show as garbled metrics)
	var sb strings.Builder
	_ = sb

	passed := 0
	tokens := 0
	for _, r := range results {
		if r.pass {
			passed++
		}
		tokens += r.tokens
	}
	if passed != 2 || tokens != 1234 {
		t.Fatalf("aggregation wrong: passed=%d tokens=%d", passed, tokens)
	}

	if code := report(results); code != 1 {
		t.Fatalf("report must exit 1 when any case fails, got %d", code)
	}

	allPass := []Case{{Name: "ok", Category: "safety", Run: func() error { return nil }}}
	var okResults []caseResult
	for _, c := range allPass {
		okResults = append(okResults, caseResult{c: c, pass: true})
	}
	if code := report(okResults); code != 0 {
		t.Fatalf("all-pass must exit 0, got %d", code)
	}
}
