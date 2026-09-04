package tools

import (
	"fmt"
	"scorp-agent/internal/helpers"
	"strings"
	"sync"
)

// ──────────────────────────────────────────────
// Todo Tool — Task tracking for multi-step work
// Encapsulated within TodoManager
// ──────────────────────────────────────────────

type TodoItem struct {
	ID      string
	Content string
	Status  string // pending, in_progress, completed, cancelled
}

type TodoManager struct {
	mu         sync.Mutex
	sessionMap map[string][]TodoItem
	seqMap     map[string]int
}

var defaultTodoManager = NewTodoManager()

func NewTodoManager() *TodoManager {
	return &TodoManager{
		sessionMap: make(map[string][]TodoItem),
		seqMap:     make(map[string]int),
	}
}

func GetDefaultTodoManager() *TodoManager {
	return defaultTodoManager
}

func (tm *TodoManager) ResolveKey(args map[string]interface{}) string {
	if s, ok := args["_session_id"].(string); ok && s != "" {
		return s
	}
	if s, ok := args["session_id"].(string); ok && s != "" {
		return s
	}
	return "default"
}

func (tm *TodoManager) Execute(args map[string]interface{}) (string, bool) {
	sessionKey := tm.ResolveKey(args)

	tm.mu.Lock()
	defer tm.mu.Unlock()

	rawTodos, ok := args["todos"].([]interface{})
	if !ok || len(rawTodos) == 0 {
		return tm.formatLocked(sessionKey), true
	}

	merge := helpers.GetBoolArg(args, "merge", false)

	if !merge {
		tm.sessionMap[sessionKey] = nil
		tm.seqMap[sessionKey] = 0
	}

	list := tm.sessionMap[sessionKey]
	seq := tm.seqMap[sessionKey]

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

		if status != "pending" && status != "in_progress" && status != "completed" && status != "cancelled" {
			status = "pending"
		}

		if status == "in_progress" {
			inProgressCount++
		}

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

	tm.sessionMap[sessionKey] = list
	tm.seqMap[sessionKey] = seq

	return tm.formatLocked(sessionKey), true
}

func (tm *TodoManager) formatLocked(sessionKey string) string {
	list := tm.sessionMap[sessionKey]
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

// ──────────────────────────────────────────────
// Compatibility Layer
// ──────────────────────────────────────────────

func ExecuteTodo(args map[string]interface{}) (string, bool) {
	return defaultTodoManager.Execute(args)
}

func formatTodoList() string {
	defaultTodoManager.mu.Lock()
	defer defaultTodoManager.mu.Unlock()
	return defaultTodoManager.formatLocked("default")
}

func getStringArgFromMap(m map[string]interface{}, key, defaultVal string) string {
	if val, ok := m[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return defaultVal
}
