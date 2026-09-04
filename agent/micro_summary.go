package agent

import (
	"context"
	"fmt"
	"strings"
	"time"

	"scorp-agent/models"
)

// ──────────────────────────────────────────────
// Micro-Summarizer for Deep / Old Tool Results
// Instead of blind string truncation, invokes free/fast model
// to extract essential facts and keys into a concise technical bullet.
// ──────────────────────────────────────────────

// SummarizeOldToolResult uses the fast/free model to extract technical keypoints
func SummarizeOldToolResult(toolName, content string) string {
	// Stub fallback format
	firstLine := content
	if idx := strings.Index(content, "\n"); idx > 0 {
		firstLine = content[:idx]
	}
	if len(firstLine) > 80 {
		firstLine = firstLine[:80]
	}
	stubFallback := fmt.Sprintf("%s ...[%d chars trimmed]", firstLine, len(content))

	if len(content) < 500 {
		return stubFallback
	}

	// Truncate input sample if extraordinarily huge to save latency
	sample := content
	if len(sample) > 2000 {
		sample = sample[:2000]
	}

	prompt := fmt.Sprintf("Compress the following tool output (%s) into 1-2 ultra-concise factual bullet points. Under 150 chars total:\n\n%s", toolName, sample)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	chatMsgs := []models.ChatMessage{
		{Role: "system", Content: "You are a concise technical summarizer. Output facts only."},
		{Role: "user", Content: prompt},
	}

	reply, _, err := models.CallModelWithFallback(ctx, "fast", chatMsgs)
	if err != nil || strings.TrimSpace(reply) == "" {
		return stubFallback
	}

	summary := strings.TrimSpace(reply)
	if len(summary) > 180 {
		summary = summary[:180] + "..."
	}
	return fmt.Sprintf("%s ...[%d chars trimmed]\n%s", firstLine, len(content), summary)
}
