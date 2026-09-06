package agent

import (
	"strings"
	"testing"
)

func heavyHistory() []AgentMessage {
	h := []AgentMessage{{Role: "system", Content: "system prompt"}}
	h = append(h, AgentMessage{Role: "user", Content: "old task: do DNS analysis with many steps"})
	for i := 0; i < 30; i++ {
		h = append(h, AgentMessage{Role: "assistant", Content: "running step"})
		h = append(h, AgentMessage{Role: "user", Content: "[Tool Result: shell] " + strings.Repeat("x", 1200)})
	}
	h = append(h, AgentMessage{Role: "assistant", Content: strings.Repeat("Final report of the DNS task. ", 20)})
	return h
}

func TestRollUpSkipsSmallHistory(t *testing.T) {
	small := []AgentMessage{
		{Role: "system", Content: "sys"},
		{Role: "user", Content: "hi"},
		{Role: "assistant", Content: "hello"},
	}
	got := prepareNewTurnHistory("sess", small, false)
	if len(got) != 3 {
		t.Fatalf("small history must pass through, got %d msgs", len(got))
	}
}

func TestRollUpSkipsContinuation(t *testing.T) {
	got := prepareNewTurnHistory("sess", heavyHistory(), true)
	if len(got) <= 34 {
		t.Fatal("continuation must keep full history")
	}
}

func TestRollUpSkipsActivePlan(t *testing.T) {
	sid := "rollup-plan"
	defer ClearTaskPlan(sid)
	execTaskPlanTool(sid, map[string]interface{}{"action": "create", "goal": "g", "items": []interface{}{"a", "b"}})

	got := prepareNewTurnHistory(sid, heavyHistory(), false)
	if len(got) <= 34 {
		t.Fatal("active plan must keep full history")
	}
}

func TestRollUpTrimsAndMarksBoundary(t *testing.T) {
	sid := "rollup-fresh" // no plan
	got := prepareNewTurnHistory(sid, heavyHistory(), false)

	if estimateHistoryTokens(got) >= estimateHistoryTokens(heavyHistory()) {
		t.Fatal("roll-up must shrink heavy history")
	}

	last := got[len(got)-1].Content.(string)
	if !strings.Contains(last, "TURN BOUNDARY") {
		t.Fatalf("boundary marker missing, last msg: %q", last[:min(80, len(last))])
	}

	joined := ""
	for _, m := range got {
		if c, ok := m.Content.(string); ok {
			joined += c + "\n"
		}
	}
	if !strings.Contains(joined, "Final report of the DNS task") {
		t.Fatal("previous task final report must be preserved")
	}
	if !strings.Contains(joined, "system prompt") {
		t.Fatal("system prompt must be preserved")
	}
	// Old tool results must be gone.
	if strings.Contains(joined, strings.Repeat("x", 1200)) {
		t.Fatal("stale tool results must not survive roll-up")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
