package agent

import (
	"strings"
	"testing"
)

// ──────────────────────────────────────────────
// Helper: generate a tool result of given size
// ──────────────────────────────────────────────

func makeToolResult(size int) string {
	base := "[Tool Result: docker ps]\n"
	if size <= len(base) {
		return base
	}
	padding := strings.Repeat("x", size-len(base))
	return base + padding
}

// makeHistory builds a history slice: system + pairs of (user-tool-result, assistant-reply)
// toolResultSizes[i] = size of the i-th tool result; 0 = no tool result (just a normal user msg)
func makeHistory(toolResultSizes []int) []AgentMessage {
	var h []AgentMessage
	h = append(h, AgentMessage{Role: "system", Content: "You are a helpful assistant."})
	for _, size := range toolResultSizes {
		if size > 0 {
			h = append(h, AgentMessage{Role: "user", Content: makeToolResult(size)})
		} else {
			h = append(h, AgentMessage{Role: "user", Content: "What's the status?"})
		}
		h = append(h, AgentMessage{Role: "assistant", Content: "OK, got it."})
	}
	return h
}

// ──────────────────────────────────────────────
// Test 1: Recent tool result (age < 6) — should be kept at recentToolResultMax
// ──────────────────────────────────────────────

func TestPrune_RecentToolResult_KeptFull(t *testing.T) {
	// 1 tool result at the end (age=1 since assistant msg is last)
	// size = 2000, well under recentToolResultMax(3000) — should not be touched
	history := makeHistory([]int{2000})
	result, n := truncateToolResultsInHistory(history)

	if n != 0 {
		t.Errorf("expected 0 truncations, got %d", n)
	}
	content := result[1].Content.(string)
	if len(content) != 2000 {
		t.Errorf("recent tool result modified: expected 2000 chars, got %d", len(content))
	}
}

// ──────────────────────────────────────────────
// Test 2: Recent but oversized tool result — trimmed to recentToolResultMax
// ──────────────────────────────────────────────

func TestPrune_RecentOversized_TrimmedTo3000(t *testing.T) {
	// 1 tool result at end, size=5000, age=1 → should trim to recentToolResultMax=3000
	history := makeHistory([]int{5000})
	result, n := truncateToolResultsInHistory(history)

	if n != 1 {
		t.Fatalf("expected 1 truncation, got %d", n)
	}
	content := result[1].Content.(string)
	// Should be approximately 3000 (head 2250 + marker + tail 750)
	if len(content) > 3100 {
		t.Errorf("recent oversized not trimmed properly: got %d chars (expected ~3000)", len(content))
	}
	if !strings.Contains(content, "[trimmed]") {
		t.Error("trimmed content should contain [trimmed] marker")
	}
}

// ──────────────────────────────────────────────
// Test 3: Old tool result (age 6-11) — trimmed to oldToolResultMax (500)
// ──────────────────────────────────────────────

func TestPrune_OldToolResult_TrimmedTo500(t *testing.T) {
	// Build: 8 tool results of 2000 chars each.
	// Each pair = user(tool) + assistant. So 16 msgs + system = 17 total.
	// Last pair: tool at index 15, age=1 (recent)
	// Index 13: age=3 (recent)
	// Index 11: age=5 (recent, last in recent band since age < 6)
	// Index 9: age=7 (old band: 6-11)
	// Index 1: age=15 (very old band: >= 12)
	sizes := []int{2000, 2000, 2000, 2000, 2000, 2000, 2000, 2000}
	history := makeHistory(sizes)
	result, n := truncateToolResultsInHistory(history)

	// Tool results at indices 1,3,5,7,9,11,13,15 (system at 0)
	// age of each: 15,13,11,9,7,5,3,1
	// age < 6 → recent: indices 13(age3), 15(age1) — 2000 chars, under 3000 → NOT trimmed
	// age < 12 → old: indices 9(age7), 11(age5→wait, age=5 < 6 so recent)
	// Let me recalculate:
	// history has 17 elements (index 0-16)
	// n-1 = 16
	// tool at index 1: age=15 → very old (≥12) → stub
	// tool at index 3: age=13 → very old → stub
	// tool at index 5: age=11 → old (6-11) → 500
	// tool at index 7: age=9 → old → 500
	// tool at index 9: age=7 → old → 500
	// tool at index 11: age=5 → recent (<6) → 2000 < 3000 → no trim
	// tool at index 13: age=3 → recent → no trim
	// tool at index 15: age=1 → recent → no trim
	// So: 2 very old + 3 old = 5 trims expected
	if n != 5 {
		t.Fatalf("expected 5 truncations, got %d", n)
	}

	// Verify old band (age 6-11) trimmed to ~500
	oldContent := result[9].Content.(string) // age=7
	if len(oldContent) > 600 {
		t.Errorf("old tool result (age 7) not trimmed to ~500: got %d chars", len(oldContent))
	}

	// Verify recent (age < 6) NOT trimmed
	recentContent := result[11].Content.(string) // age=5
	if len(recentContent) != 2000 {
		t.Errorf("recent tool result (age 5) should not be trimmed: expected 2000, got %d", len(recentContent))
	}
}

// ──────────────────────────────────────────────
// Test 4: Very old tool result (age >= 12) — reduced to stub
// ──────────────────────────────────────────────

func TestPrune_VeryOldToolResult_StubOnly(t *testing.T) {
	// 7 tool results → indices 1,3,5,7,9,11,13 in 15-element history (n-1=14)
	// tool at index 1: age=13 → very old (≥12)
	// Others: age=11,9,7,5,3,1
	history := makeHistory([]int{5000, 2000, 2000, 2000, 2000, 2000, 2000})
	result, n := truncateToolResultsInHistory(history)

	// tool at index 1: age=13 → very old, size=5000 → stub
	// tool at index 3: age=11 → old, size=2000 > 500 → trimmed
	// tool at index 5: age=9 → old, size=2000 > 500 → trimmed
	// tool at index 7: age=7 → old, size=2000 > 500 → trimmed
	// tool at index 9: age=5 → recent, 2000 < 3000 → no trim
	// tool at index 11: age=3 → recent → no trim
	// tool at index 13: age=1 → recent → no trim
	if n != 4 {
		t.Fatalf("expected 4 truncations (1 very-old + 3 old), got %d", n)
	}

	stub := result[1].Content.(string)
	if !strings.Contains(stub, "[Tool Result: docker ps]") {
		t.Error("very old stub should preserve first line")
	}
	if !strings.Contains(stub, "chars trimmed]") {
		t.Error("very old stub should contain chars-trimmed indicator")
	}
	// Stub should be short (< 200 chars)
	if len(stub) > 200 {
		t.Errorf("very old stub too long: %d chars (expected < 200)", len(stub))
	}
}

// ──────────────────────────────────────────────
// Test 5: Non-tool-result messages are never touched
// ──────────────────────────────────────────────

func TestPrune_NonToolMessages_Preserved(t *testing.T) {
	history := []AgentMessage{
		{Role: "system", Content: "You are helpful."},
		{Role: "user", Content: "Check docker status"},
		{Role: "assistant", Content: "Let me check that for you."},
		{Role: "user", Content: makeToolResult(5000)}, // age=0 → recent but oversized
		{Role: "assistant", Content: "Docker is running fine."},
		{Role: "user", Content: "Great, thanks!"},
	}

	result, n := truncateToolResultsInHistory(history)

	// Only 1 tool result, age=0 (most recent user msg is not a tool result,
	// the tool result is at index 3, age = 5-1-3 = 1 → recent)
	if n != 1 {
		t.Fatalf("expected 1 truncation, got %d", n)
	}

	// All non-tool messages should be unchanged
	if result[0].Content.(string) != "You are helpful." {
		t.Error("system message modified")
	}
	if result[1].Content.(string) != "Check docker status" {
		t.Error("user message modified")
	}
	if result[2].Content.(string) != "Let me check that for you." {
		t.Error("assistant message modified")
	}
	if result[4].Content.(string) != "Docker is running fine." {
		t.Error("assistant reply modified")
	}
	if result[5].Content.(string) != "Great, thanks!" {
		t.Error("final user message modified")
	}
}

// ──────────────────────────────────────────────
// Test 6: Short tool results are never trimmed regardless of age
// ──────────────────────────────────────────────

func TestPrune_ShortToolResult_NeverTrimmed(t *testing.T) {
	// All tool results are 100 chars — well under any threshold
	sizes := []int{100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	history := makeHistory(sizes)
	_, n := truncateToolResultsInHistory(history)

	if n != 0 {
		t.Errorf("short tool results should never be trimmed, got %d truncations", n)
	}
}

// ──────────────────────────────────────────────
// Test 7: Boundary ages — age exactly 5, 6, 11, 12
// ──────────────────────────────────────────────

func TestPrune_BoundaryAges(t *testing.T) {
	// Build history with tool results at specific ages.
	// We want: age=5 (last recent), age=6 (first old), age=11 (last old), age=12 (first very-old)
	// History = system + [user(tool), assistant] * 7 + [user(tool), assistant] at the end
	// Indices: 0=system, 1=t0, 2=a0, 3=t1, 4=a1, 5=t2, 6=a2, 7=t3, 8=a3,
	//           9=t4, 10=a4, 11=t5, 12=a5, 13=t6, 14=a6
	// n=15, n-1=14
	// t0 at idx 1: age=13 (very-old)
	// t1 at idx 3: age=11 (old, last old)
	// t2 at idx 5: age=9 (old)
	// t3 at idx 7: age=7 (old)
	// t4 at idx 9: age=5 (recent, last recent)
	// t5 at idx 11: age=3 (recent)
	// t6 at idx 13: age=1 (recent)

	// To get exact boundary ages 5, 6, 11, 12 we need more control.
	// Let's place tool results at precise positions:
	history := []AgentMessage{
		{Role: "system", Content: "sys"}, // idx 0
	}
	// We'll add filler to create specific ages
	// Target: age=12 → very-old boundary, age=11 → old boundary,
	//         age=6 → old boundary, age=5 → recent boundary
	// Total history size will be 15 (indices 0-14, n-1=14)
	// age=12 → idx 2, age=11 → idx 3, age=6 → idx 8, age=5 → idx 9

	// idx 1: filler user
	history = append(history, AgentMessage{Role: "user", Content: "filler"})
	// idx 2: age=12 tool result (first very-old)
	history = append(history, AgentMessage{Role: "user", Content: makeToolResult(1000)})
	// idx 3: age=11 tool result (last old)
	history = append(history, AgentMessage{Role: "user", Content: makeToolResult(1000)})
	// idx 4-7: fillers
	for i := 0; i < 4; i++ {
		history = append(history, AgentMessage{Role: "assistant", Content: "filler"})
	}
	// idx 8: age=6 tool result (first old)
	history = append(history, AgentMessage{Role: "user", Content: makeToolResult(1000)})
	// idx 9: age=5 tool result (last recent)
	history = append(history, AgentMessage{Role: "user", Content: makeToolResult(4000)}) // oversized recent
	// idx 10-14: fillers
	for i := 0; i < 5; i++ {
		history = append(history, AgentMessage{Role: "assistant", Content: "filler"})
	}

	// n = 15, n-1 = 14
	result, n := truncateToolResultsInHistory(history)

	// idx 2 (age=12): very-old, 1000 > 120 → stub
	stubContent := result[2].Content.(string)
	if len(stubContent) > 200 {
		t.Errorf("age=12 should be stub (<200 chars), got %d", len(stubContent))
	}

	// idx 3 (age=11): old, 1000 > 500 → trimmed to ~500
	oldContent := result[3].Content.(string)
	if len(oldContent) > 600 {
		t.Errorf("age=11 should be trimmed to ~500, got %d", len(oldContent))
	}

	// idx 8 (age=6): old, 1000 > 500 → trimmed to ~500
	oldContent2 := result[8].Content.(string)
	if len(oldContent2) > 600 {
		t.Errorf("age=6 should be trimmed to ~500, got %d", len(oldContent2))
	}

	// idx 9 (age=5): recent, 4000 > 3000 → trimmed to ~3000
	recentContent := result[9].Content.(string)
	if len(recentContent) > 3100 {
		t.Errorf("age=5 oversized should be trimmed to ~3000, got %d", len(recentContent))
	}

	_ = n // truncation count varies, we checked individual results
}

// ──────────────────────────────────────────────
// Test 8: Empty and minimal history edge cases
// ──────────────────────────────────────────────

func TestPrune_EmptyHistory(t *testing.T) {
	result, n := truncateToolResultsInHistory([]AgentMessage{})
	if n != 0 {
		t.Errorf("empty history should have 0 truncations, got %d", n)
	}
	if len(result) != 0 {
		t.Errorf("empty history should return empty, got %d elements", len(result))
	}
}

func TestPrune_SingleMessage(t *testing.T) {
	history := []AgentMessage{
		{Role: "system", Content: "sys"},
	}
	result, n := truncateToolResultsInHistory(history)
	if n != 0 {
		t.Errorf("single system message should have 0 truncations, got %d", n)
	}
	if len(result) != 1 {
		t.Errorf("should return 1 element, got %d", len(result))
	}
}

// ──────────────────────────────────────────────
// Test 9: Head+tail format preserves beginning and end
// ──────────────────────────────────────────────

func TestPrune_HeadTailFormat(t *testing.T) {
	// Tool result with identifiable head and tail
	content := "[Tool Result: test]\nHEAD_MARKER_" + strings.Repeat("x", 4000) + "_TAIL_MARKER"
	history := []AgentMessage{
		{Role: "system", Content: "sys"},
		{Role: "user", Content: content}, // age=1, recent, size ~4100 → trim to 3000
		{Role: "assistant", Content: "ok"},
	}

	result, n := truncateToolResultsInHistory(history)
	if n != 1 {
		t.Fatalf("expected 1 truncation, got %d", n)
	}

	trimmed := result[1].Content.(string)
	if !strings.HasPrefix(trimmed, "[Tool Result: test]") {
		t.Error("trimmed content should preserve head (prefix)")
	}
	if !strings.Contains(trimmed, "HEAD_MARKER") {
		t.Error("trimmed content should contain HEAD_MARKER")
	}
	if !strings.Contains(trimmed, "TAIL_MARKER") {
		t.Error("trimmed content should contain TAIL_MARKER")
	}
}

// ──────────────────────────────────────────────
// Test 10: estimateTokens basic sanity
// ──────────────────────────────────────────────

func TestEstimateTokens(t *testing.T) {
	tokens := estimateTokens("hello world")
	// 11 chars * 0.25 = 2.75 → 2
	if tokens != 2 {
		t.Errorf("expected 2 tokens for 'hello world', got %d", tokens)
	}

	// Empty string
	if estimateTokens("") != 0 {
		t.Error("empty string should be 0 tokens")
	}

	// Large string
	big := strings.Repeat("a", 1000)
	if estimateTokens(big) != 250 {
		t.Errorf("expected 250 tokens for 1000 chars, got %d", estimateTokens(big))
	}
}

// ──────────────────────────────────────────────
// Test 11: estimateHistoryTokens with mixed content
// ──────────────────────────────────────────────

func TestEstimateHistoryTokens(t *testing.T) {
	history := []AgentMessage{
		{Role: "system", Content: "You are helpful."},     // 16 chars → 4 tokens + 4 overhead = 8
		{Role: "user", Content: makeToolResult(100)},      // 100 chars → 25 tokens + 4 = 29
		{Role: "assistant", Content: "OK got it."},        // 10 chars → 2 tokens + 4 = 6
	}

	tokens := estimateHistoryTokens(history)
	// 8 + 29 + 6 = 43
	if tokens != 43 {
		t.Errorf("expected 43 tokens, got %d", tokens)
	}
}

// ──────────────────────────────────────────────
// Test 12: Token savings after pruning (the whole point)
// ──────────────────────────────────────────────

func TestPrune_TokenSavings(t *testing.T) {
	// 8 tool results of 3000 chars each
	// Without pruning: 8 * 3000 * 0.25 = 6000 tokens just from tool results
	// After pruning: 3 recent (3000 each) + 3 old (500 each) + 2 very-old (120 each)
	//               = 9000 + 1500 + 240 = 10740 chars → ~2685 tokens
	// Savings: ~3315 tokens
	sizes := []int{3000, 3000, 3000, 3000, 3000, 3000, 3000, 3000}
	history := makeHistory(sizes)

	beforeTokens := estimateHistoryTokens(history)
	pruned, _ := truncateToolResultsInHistory(history)
	afterTokens := estimateHistoryTokens(pruned)

	savings := beforeTokens - afterTokens
	if savings < 2000 {
		t.Errorf("expected significant token savings (>2000), got only %d (before=%d, after=%d)",
			savings, beforeTokens, afterTokens)
	}

	t.Logf("Token savings: %d → %d (saved %d tokens, %.1f%% reduction)",
		beforeTokens, afterTokens, savings,
		float64(savings)/float64(beforeTokens)*100)
}

// ──────────────────────────────────────────────
// Test 13: Simulates the bug scenario — 41 messages with Docker data
// ──────────────────────────────────────────────

func TestPrune_DockerScenario_41Messages(t *testing.T) {
	// Simulate: system + 20 tool results (docker ps, docker stats, etc.) of 2800 chars each
	// Total = 1 + 40 = 41 messages
	sizes := make([]int, 20)
	for i := range sizes {
		sizes[i] = 2800
	}
	history := makeHistory(sizes)

	// This should be 41 messages
	if len(history) != 41 {
		t.Fatalf("expected 41 messages, got %d", len(history))
	}

	beforeTokens := estimateHistoryTokens(history)
	pruned, truncCount := truncateToolResultsInHistory(history)
	afterTokens := estimateHistoryTokens(pruned)

	// Should have trimmed most old tool results
	// Recent (age < 6): ~3 tool results not trimmed (they're under 3000)
	// Old (age 6-11): ~6 trimmed to 500
	// Very old (age >= 12): ~11 trimmed to stub
	if truncCount < 15 {
		t.Errorf("expected at least 15 truncations for 41-message scenario, got %d", truncCount)
	}

	savings := beforeTokens - afterTokens
	if savings < 10000 {
		t.Errorf("expected >10k token savings, got %d (before=%d, after=%d)",
			savings, beforeTokens, afterTokens)
	}

	t.Logf("Docker 41-msg scenario: %d → %d tokens (saved %d, %.1f%% reduction, %d truncations)",
		beforeTokens, afterTokens, savings,
		float64(savings)/float64(beforeTokens)*100, truncCount)
}
