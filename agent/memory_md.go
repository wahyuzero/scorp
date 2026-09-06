package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"scorp-agent/config"
	"scorp-agent/internal/helpers"
	"scorp-agent/models"
)

// ──────────────────────────────────────────────
// Durable Memory — MEMORY.md protocol (P1.7)
//
// "Agent forgets across sessions" was a permanent community complaint; the
// auto-memory file is Anthropic's official answer and the 2026 consensus
// (quota ≈ 200 lines). Unlike the KV memory.json + self-review heuristic,
// this protocol writes DECISIONS & STATE via a dedicated extraction prompt
// at task completion and injects the file after the system prompt at session
// start. Silent, async, panic-safe — it must never disturb the task itself.
// ──────────────────────────────────────────────

var memoryMDFile = config.ScorpPath("MEMORY.md")

const (
	memoryMDMaxLines    = 200 // file quota (2026 consensus ≈ 200 lines)
	memoryMDMaxInject   = 4000 // injection cap (chars) for session-start context
	memoryExtractLimit  = 5    // max entries per extraction
	memoryExtractMsgs   = 20   // transcript window for extraction
	memoryExtractChars  = 1200 // per-message transcript cap
)

var memoryMDMu sync.Mutex

// ReadMemoryMD returns the durable memory file contents, bounded to the
// injection cap. Missing file → "".
func ReadMemoryMD() string {
	memoryMDMu.Lock()
	defer memoryMDMu.Unlock()
	return readMemoryMDLocked()
}

func readMemoryMDLocked() string {
	data, err := os.ReadFile(memoryMDFile)
	if err != nil {
		return ""
	}
	out := string(data)
	if len(out) > memoryMDMaxInject {
		out = out[:memoryMDMaxInject] + "\n... (truncated)"
	}
	return strings.TrimSpace(out)
}

// AppendMemoryMD adds entries as bullets, skipping duplicates
// (case-insensitive) and enforcing the line quota by dropping the OLDEST
// bullets. Returns how many entries were actually added. Header lines
// (starting with '#') always survive the quota trim.
func AppendMemoryMD(entries []string) int {
	memoryMDMu.Lock()
	defer memoryMDMu.Unlock()

	normalized := make([]string, 0, len(entries))
	seen := map[string]bool{}
	for _, e := range entries {
		e = strings.TrimSpace(e)
		if e == "" {
			continue
		}
		e = strings.TrimPrefix(e, "- ")
		if len(e) > 300 {
			e = e[:300]
		}
		key := strings.ToLower(e)
		if seen[key] {
			continue
		}
		seen[key] = true
		normalized = append(normalized, "- "+e)
	}
	if len(normalized) == 0 {
		return 0
	}

	existing := readMemoryMDLocked()
	var headers, bullets []string
	seenExisting := map[string]bool{}
	for _, line := range strings.Split(existing, "\n") {
		line = strings.TrimRight(line, " \t")
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			headers = append(headers, line)
			continue
		}
		key := strings.ToLower(strings.TrimPrefix(line, "- "))
		if seenExisting[key] {
			continue // drop pre-existing duplicates while we're here
		}
		seenExisting[key] = true
		bullets = append(bullets, line)
	}

	added := 0
	for _, e := range normalized {
		key := strings.ToLower(strings.TrimPrefix(e, "- "))
		if seenExisting[key] {
			continue // already known
		}
		seenExisting[key] = true
		bullets = append(bullets, e)
		added++
	}
	if added == 0 {
		return 0
	}

	// Quota: newest bullets win, oldest are dropped (headers always kept)
	quota := memoryMDMaxLines - len(headers)
	if quota < 10 {
		quota = 10
	}
	if len(bullets) > quota {
		bullets = bullets[len(bullets)-quota:]
	}

	var sb strings.Builder
	for _, h := range headers {
		sb.WriteString(h + "\n")
	}
	if len(headers) == 0 {
		sb.WriteString("# Durable Memory — decisions & state that outlive tasks\n\n")
	}
	for _, b := range bullets {
		sb.WriteString(b + "\n")
	}

	path := memoryMDFile
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(sb.String()), 0644); err != nil {
		log.Printf("[memory-md] write error: %v", err)
		return 0
	}
	if err := os.Rename(tmp, path); err != nil {
		log.Printf("[memory-md] rename error: %v", err)
		return 0
	}
	return added
}

// parseMemoryEntries parses the extractor's reply into clean one-line
// entries: tolerates markdown fences and JSON arrays of strings or objects.
func parseMemoryEntries(resp string) []string {
	resp = strings.TrimSpace(resp)
	resp = strings.TrimPrefix(resp, "```json")
	resp = strings.TrimPrefix(resp, "```")
	resp = strings.TrimSuffix(resp, "```")
	resp = strings.TrimSpace(resp)

	var raw []interface{}
	if err := json.Unmarshal([]byte(resp), &raw); err != nil {
		return nil
	}
	var out []string
	for _, r := range raw {
		switch v := r.(type) {
		case string:
			if s := strings.TrimSpace(v); s != "" {
				out = append(out, s)
			}
		case map[string]interface{}:
			for _, k := range []string{"content", "memory", "fact", "text"} {
				if s, ok := v[k].(string); ok && strings.TrimSpace(s) != "" {
					out = append(out, strings.TrimSpace(s))
					break
				}
			}
		}
	}
	if len(out) > memoryExtractLimit {
		out = out[:memoryExtractLimit]
	}
	return out
}

// extractTaskMemory runs the dedicated memory-protocol extraction after a
// task concludes: a cheap model call turns the tail of the conversation into
// 0-5 durable entries appended to MEMORY.md. Async, silent, panic-safe.
func extractTaskMemory(chatIDStr string, history []AgentMessage) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[memory-md] panic recovered: %v", r)
		}
	}()

	start := len(history) - memoryExtractMsgs
	if start < 1 {
		start = 1 // skip system prompt
	}

	var conv strings.Builder
	for _, m := range history[start:] {
		text, ok := m.Content.(string)
		if !ok {
			continue
		}
		if len(text) > memoryExtractChars {
			text = text[:memoryExtractChars] + "..."
		}
		switch m.Role {
		case "user":
			fmt.Fprintf(&conv, "User: %s\n", text)
		case "assistant":
			fmt.Fprintf(&conv, "Assistant: %s\n", text)
		}
	}
	if conv.Len() < 80 {
		return
	}

	prompt := fmt.Sprintf(`A task just completed. Extract up to %d durable memory entries worth keeping for FUTURE sessions — decisions made, project/system state that changed, environment facts, user preferences, and where things live. One line each, self-contained (mention file paths/names explicitly).

Skip: transient task chatter, intermediate outputs, anything that will not matter next week, and anything already obvious.

Return ONLY a JSON array of strings. If nothing is worth remembering, return [].

Task conversation tail:
%s`, memoryExtractLimit, conv.String())

	model := models.RouteModel("chat")
	if model == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := models.CallModel(ctx, model, []models.ChatMessage{
		{Role: "system", Content: "You are a memory extraction assistant. Return ONLY a valid JSON array of strings, no markdown fences, no explanation."},
		{Role: "user", Content: prompt},
	})
	if err != nil {
		log.Printf("[memory-md] extraction LLM call failed: %v", err)
		return
	}

	entries := parseMemoryEntries(resp)
	if len(entries) == 0 {
		return
	}
	if added := AppendMemoryMD(entries); added > 0 {
		log.Printf("[memory-md] Task %s: added %d durable entr(ies) to MEMORY.md", helpers.TruncateStr(chatIDStr, 40), added)
	}
}

// SetMemoryMDFile overrides the durable-memory file location (eval/embedding
// use); empty string restores the default ~/.scorp/MEMORY.md.
func SetMemoryMDFile(path string) {
	memoryMDMu.Lock()
	defer memoryMDMu.Unlock()
	if path == "" {
		memoryMDFile = config.ScorpPath("MEMORY.md")
		return
	}
	memoryMDFile = path
}
