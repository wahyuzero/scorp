package main

import (
	"fmt"
	"strings"
	"sync"
)

// ──────────────────────────────────────────────
// Todo Tool — Task tracking for multi-step work
// ──────────────────────────────────────────────

type TodoItem struct {
	ID      string
	Content string
	Status  string // pending, in_progress, completed, cancelled
}

// In-memory todo list (single user = single list is fine)
var (
	todoList   []TodoItem
	todoMu     sync.Mutex
	todoIDSeq  int
)

// executeTodo handles the "todo" tool.
// No args: returns current list.
// With todos array: updates list (merge=false replaces, merge=true updates by id).
func executeTodo(args map[string]interface{}) (string, bool) {
	// No args → return formatted list
	if len(args) == 0 {
		return formatTodoList(), true
	}

	rawTodos, ok := args["todos"].([]interface{})
	if !ok {
		return "Error: 'todos' must be an array of {id, content, status}", false
	}

	merge := getBoolArg(args, "merge", false)

	todoMu.Lock()
	defer todoMu.Unlock()

	if !merge {
		todoList = nil
		todoIDSeq = 0
	}

	inProgressCount := 0
	for _, item := range todoList {
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
			todoIDSeq++
			id = fmt.Sprintf("t%d", todoIDSeq)
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
			for i := range todoList {
				if todoList[i].ID == id {
					todoList[i].Content = content
					todoList[i].Status = status
					updated = true
					break
				}
			}
		}

		if !updated {
			todoList = append(todoList, TodoItem{ID: id, Content: content, Status: status})
		}
	}

	// Enforce: only ONE in_progress at a time
	if inProgressCount > 1 {
		first := true
		for i := range todoList {
			if todoList[i].Status == "in_progress" {
				if first {
					first = false
				} else {
					todoList[i].Status = "pending"
				}
			}
		}
	}

	return formatTodoListLocked(), true
}

// formatTodoList returns the formatted todo list (thread-safe wrapper)
func formatTodoList() string {
	todoMu.Lock()
	defer todoMu.Unlock()
	return formatTodoListLocked()
}

// formatTodoListLocked returns formatted todo list WITHOUT locking.
// Caller MUST hold todoMu.
func formatTodoListLocked() string {
	if len(todoList) == 0 {
		return "📋 Todo list is empty.\nUse: todos=[{id, content, status}] to create items."
	}

	var sb strings.Builder
	sb.WriteString("📋 <b>Todo List</b>\n\n")

	statusIcon := map[string]string{
		"pending":      "🔲",
		"in_progress":  "🔄",
		"completed":    "✅",
		"cancelled":    "❌",
	}

	for i, item := range todoList {
		icon := statusIcon[item.Status]
		if icon == "" {
			icon = "🔲"
		}
		sb.WriteString(fmt.Sprintf("%d. %s <b>%s</b> <code>[%s]</code>\n", i+1, icon, item.Content, item.ID))
	}

	sb.WriteString("\n💡 <b>Usage:</b>\n")
	sb.WriteString(`<code>{"todos": [{"id": "t1", "content": "Task 1", "status": "in_progress"}]}</code> — replace list`)
	sb.WriteString("\n")
	sb.WriteString(`<code>{"todos": [{"id": "t1", "status": "completed"}], "merge": true}</code> — update item`)

	return sb.String()
}

// getStringArgFromMap gets string arg from map[string]interface{}
func getStringArgFromMap(m map[string]interface{}, key, defaultVal string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultVal
}
