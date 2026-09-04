package tools

import (
	"fmt"
	"scorp-agent/internal/helpers"
	"strings"
	"sync"
)

// ──────────────────────────────────────────────
// Todo Tool — Task tracking for multi-step work
// Supports session isolation so CLI and Telegram sessions don't clash.
// ──────────────────────────────────────────────

type TodoItem struct {
	ID      string
	Content string
	Status  string // pending, in_progress, completed, cancelled
}

// Session-isolated todo store
var (
	todoSessionMap = make(map[string][]TodoItem)
	todoIDSeqMap   = make(map[string]int)
	todoMu         sync.Mutex
)

// resolveSessionKey extracts a session identifier from args or defaults to "default"
func resolveSessionKey(args map[string]interface{}) string {
	if s, ok := args["_session_id"].(string); ok && s != "" {
		return s
	}
	if s, ok := args["session_id"].(string); ok && s != "" {
		return s
	}
	return "default"
}

// ExecuteTodo handles the "todo" tool.
// No args: returns current list.
// With todos array: updates list (merge=false replaces, merge=true updates by id).
func ExecuteTodo(args map[string]interface{}) (string, bool) {
	sessionKey := resolveSessionKey(args)

	todoMu.Lock()
	defer todoMu.Unlock()

	// No args or only session_id → return formatted list
	rawTodos, ok := args["todos"].([]interface{})
	if !ok || len(rawTodos) == 0 {
		return formatTodoListLocked(sessionKey), true
	}

	merge := helpers.GetBoolArg(args, "merge", false)

	if !merge {
		todoSessionMap[sessionKey] = nil
		todoIDSeqMap[sessionKey] = 0
	}

	list := todoSessionMap[sessionKey]
	seq := todoIDSeqMap[sessionKey]

	inProgressCount := 0
	for _, item := range list {
		if item.Status == "in_progress" {
			inProgressCount++
		}
	}

	for _, raw := range rawTodos {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		id := getStringArgFromMap(m, "id", "")
		content := getStringArgFromMap(m, "content", "")
		status := getStringArgFromMap(m, "status", "pending")

		if content == "" {
			continue
		}

		if id == "" {
			seq++
			id = fmt.Sprintf("t%d", seq)
		}

		// Validate status
		if status != "pending" && status != "in_progress" && status != "completed" && status != "cancelled" {
			status = "pending"
		}

		if status == "in_progress" {
			inProgressCount++
		}

		// If merge=true, update existing by id
		updated := false
		if merge {
			for i := range list {
				if list[i].ID == id {
					list[i].Content = content
					list[i].Status = status
					updated = true
					break
				}
			}
		}

		if !updated {
			list = append(list, TodoItem{ID: id, Content: content, Status: status})
		}
	}

	// Enforce: only ONE in_progress at a time
	if inProgressCount > 1 {
		first := true
		for i := range list {
			if list[i].Status == "in_progress" {
				if first {
					first = false
				} else {
					list[i].Status = "pending"
				}
			}
		}
	}

	todoSessionMap[sessionKey] = list
	todoIDSeqMap[sessionKey] = seq

	return formatTodoListLocked(sessionKey), true
}

// formatTodoList returns the formatted todo list (thread-safe wrapper)
func formatTodoList() string {
	todoMu.Lock()
	defer todoMu.Unlock()
	return formatTodoListLocked("default")
}

// formatTodoListLocked returns formatted todo list WITHOUT locking.
// Caller MUST hold todoMu.
func formatTodoListLocked(sessionKey string) string {
	list := todoSessionMap[sessionKey]
	if len(list) == 0 {
		return "📋 Todo list is empty.\nUse: todos=[{id, content, status}] to create items."
	}

	var sb strings.Builder
	sb.WriteString("📋 <b>Todo List</b>\n\n")

	statusIcon := map[string]string{
		"pending":     "🔲",
		"in_progress": "🔄",
		"completed":   "✅",
		"cancelled":   "❌",
	}

	for _, item := range list {
		icon := statusIcon[item.Status]
		if icon == "" {
			icon = "🔲"
		}

		if item.Status == "completed" {
			sb.WriteString(fmt.Sprintf("%s <s>%s. %s</s>\n", icon, item.ID, item.Content))
		} else if item.Status == "in_progress" {
			sb.WriteString(fmt.Sprintf("%s <b>%s. %s</b>\n", icon, item.ID, item.Content))
		} else {
			sb.WriteString(fmt.Sprintf("%s %s. %s\n", icon, item.ID, item.Content))
		}
	}

	return sb.String()
}

func getStringArgFromMap(m map[string]interface{}, key, defaultVal string) string {
	if val, ok := m[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return defaultVal
}
