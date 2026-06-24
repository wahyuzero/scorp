package agent

// ── Self-Improvement Review (lightweight post-turn reflection) ──
//
// After each agent turn completes, maybeRunSelfReview() checks cadence and
// optionally spawns a background goroutine. The goroutine asks the LLM to
// extract 0-3 durable facts about the user from the conversation, then merges
// them into memory.json via existing setMemory().
//
// Safety:
//   - Runs every 5 turns (not every turn)
//   - Min 10 minutes between reviews per chat
//   - Max 3 new entries per review
//   - Memory capped at 50 entries total
//   - Silent: does not message user, does not block agent loop
//   - recover() on panic — never crashes the process

import (
	"scorp-agent/models"
	"scorp-agent/tools"
	"scorp-agent/internal/helpers"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

var (
	reviewTurnCount = make(map[int64]int)
	reviewLastTime  = make(map[int64]time.Time)
	reviewMu        sync.Mutex
)

const (
	reviewCadenceTurns = 5
	reviewMinInterval  = 10 * time.Minute
	reviewMaxFacts     = 3
	reviewMemoryCap    = 50
	reviewMaxMessages  = 20
	reviewMaxMsgChars  = 1500
)

type memoryFact struct {
	Key     string `json:"key"`
	Content string `json:"content"`
}

// maybeRunSelfReview is called after agent loop finishes (final answer sent).
// Non-blocking: either spawns a goroutine or returns immediately.
func maybeRunSelfReview(chatID int64, chatIDStr string) {
	reviewMu.Lock()
	reviewTurnCount[chatID]++
	turns := reviewTurnCount[chatID]
	lastRun := reviewLastTime[chatID]
	reviewMu.Unlock()

	// Cadence: every N turns
	if turns%reviewCadenceTurns != 0 {
		return
	}
	// Rate limit
	if time.Since(lastRun) < reviewMinInterval {
		return
	}

	reviewMu.Lock()
	reviewLastTime[chatID] = time.Now()
	reviewMu.Unlock()

	// Snapshot conversation
	history := getSessionHistory(chatIDStr)
	if len(history) < 4 { // need at least sys + user + assistant + user
		return
	}

	go runSelfReview(chatID, history)
}

func runSelfReview(chatID int64, history []AgentMessage) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[self-review] panic recovered: %v", r)
		}
	}()

	// Build conversation text (skip system prompt, take last N messages)
	start := len(history) - reviewMaxMessages
	if start < 1 {
		start = 1 // skip system message at index 0
	}

	var conv strings.Builder
	for _, m := range history[start:] {
		text, ok := m.Content.(string)
		if !ok {
			continue
		}
		if len(text) > reviewMaxMsgChars {
			text = text[:reviewMaxMsgChars] + "..."
		}
		switch m.Role {
		case "user":
			fmt.Fprintf(&conv, "User: %s\n", text)
		case "assistant":
			if len(text) < 1000 {
				fmt.Fprintf(&conv, "Assistant: %s\n", text)
			}
		}
	}

	if conv.Len() < 50 {
		return
	}

	prompt := fmt.Sprintf(`Analyze the conversation below. Extract 0-%d durable facts about the USER — their preferences, habits, role, work environment, communication style, or corrections to how the assistant should behave.

Only extract facts that will still matter in future sessions. Skip:
- Temporary task details or one-off questions
- Environment-specific errors (missing packages, wrong paths)
- System status reports

Return ONLY a JSON array. Each item: {"key":"short_descriptive_key","content":"the fact"}
If nothing notable, return: []

Conversation:
%s`, reviewMaxFacts, conv.String())

	model := models.RouteModel("chat")
	if model == nil {
		log.Printf("[self-review] no model configured, skipping")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := models.CallModel(ctx, model, []models.ChatMessage{
		{Role: "system", Content: "You are a memory extraction assistant. Return ONLY a valid JSON array, no markdown, no explanation."},
		{Role: "user", Content: prompt},
	})
	if err != nil {
		log.Printf("[self-review] LLM call failed: %v", err)
		return
	}

	// Strip markdown code fences if present
	resp = strings.TrimSpace(resp)
	resp = strings.TrimPrefix(resp, "```json")
	resp = strings.TrimPrefix(resp, "```")
	resp = strings.TrimSuffix(resp, "```")
	resp = strings.TrimSpace(resp)

	var facts []memoryFact
	if err := json.Unmarshal([]byte(resp), &facts); err != nil {
		log.Printf("[self-review] JSON parse failed: %v (raw: %.200s)", err, resp)
		return
	}

	if len(facts) == 0 {
		return
	}

	existing := tools.ListMemory()
	if len(existing) >= reviewMemoryCap {
		log.Printf("[self-review] memory at cap (%d), skipping", len(existing))
		return
	}

	added := 0
	for _, f := range facts {
		if f.Key == "" || f.Content == "" || added >= reviewMaxFacts {
			break
		}
		if _, exists := existing[f.Key]; exists {
			continue
		}
		// Simple content dedup
		dup := false
		for _, v := range existing {
			if strings.Contains(strings.ToLower(v), strings.ToLower(f.Content)) ||
				strings.Contains(strings.ToLower(f.Content), strings.ToLower(v)) {
				dup = true
				break
			}
		}
		if dup {
			continue
		}

		tools.SetMemory(f.Key, f.Content)
		added++
		log.Printf("[self-review] saved: %s = %s", f.Key, helpers.TruncateStr(f.Content, 80))
	}

	if added > 0 {
		log.Printf("[self-review] extracted %d fact(s) to memory", added)
	}
}
