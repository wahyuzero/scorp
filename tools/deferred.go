package tools

import (
	"scorp-agent/internal/helpers"
	"scorp-agent/registry"
	"fmt"
	"sort"
	"strings"
)

// ──────────────────────────────────────────────
// Deferred Tool Loading — tool_search + tool_call
// Allows the LLM to discover and invoke tools that
// are NOT sent in the initial schema (marked Deferred=true).
// This reduces context overhead for rarely-used tools.
// ──────────────────────────────────────────────

// executeToolSearch searches available (including deferred) tools by keyword.
// Returns name + description + argument schema so the LLM knows how to call them.
func ExecuteToolSearch(args map[string]interface{}, chatID int64) (string, bool) {
	query := strings.ToLower(helpers.GetStringArg(args, "query", ""))
	limit := helpers.GetIntArg(args, "limit", 10)
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	type scoredTool struct {
		def   registry.ToolDef
		score int
	}

	var matches []scoredTool

	for _, def := range registry.GetAllTools() {
		// Search all tools (both native and deferred)
		score := 0
		name := strings.ToLower(def.Name)
		desc := strings.ToLower(def.Description)
		cat := strings.ToLower(def.Category)

		if query == "" {
			score = 1
		} else {
			// Name exact match → highest score
			if name == query {
				score += 100
			}
			// Name contains query
			if strings.Contains(name, query) {
				score += 50
			}
			// Description contains query
			if strings.Contains(desc, query) {
				score += 30
			}
			// Category match
			if strings.Contains(cat, query) {
				score += 20
			}
			// Individual word matches in description
			for _, word := range strings.Fields(query) {
				if strings.Contains(desc, word) {
					score += 5
				}
			}
		}

		if score > 0 {
			matches = append(matches, scoredTool{def: def, score: score})
		}
	}

	if len(matches) == 0 {
		return fmt.Sprintf("No tools found matching \"%s\"", query), true
	}

	// Sort by score descending
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].score > matches[j].score
	})

	// Limit results
	if len(matches) > limit {
		matches = matches[:limit]
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🔧 Found %d tool(s)", len(matches)))
	if query != "" {
		sb.WriteString(fmt.Sprintf(" matching \"%s\"", query))
	}
	sb.WriteString(":\n\n")

	for i, m := range matches {
		tags := ""
		if m.def.Deferred {
			tags = " ⏸️deferred"
		}
		if !m.def.Native {
			tags += " (non-native)"
		}

		sb.WriteString(fmt.Sprintf("**%d. %s**%s\n   %s\n", i+1, m.def.Name, tags, m.def.Description))

		// Show arguments
		if len(m.def.Arguments) > 0 {
			sb.WriteString("   Args:\n")
			for argName, argDef := range m.def.Arguments {
				req := ""
				if argDef.Required {
					req = " (required)"
				}
				sb.WriteString(fmt.Sprintf("     - %s (%s)%s: %s\n", argName, argDef.Type, req, argDef.Description))
			}
		}

		// Show raw schema for MCP tools
		if m.def.RawInputSchema != nil {
			sb.WriteString("   Schema: available (MCP tool)\n")
		}

		sb.WriteString("\n")
	}

	return sb.String(), true
}

// executeToolCall invokes a deferred tool by name.
// The LLM uses this after discovering a tool via tool_search.
func ExecuteToolCall(args map[string]interface{}, chatID int64) (string, bool) {
	toolName := helpers.GetStringArg(args, "name", "")
	if toolName == "" {
		return "Error: 'name' is required (the tool to call)", false
	}

	def, ok := registry.GetTool(toolName)
	if !ok {
		return fmt.Sprintf("Error: tool '%s' not found. Use tool_search to find available tools.", toolName), false
	}

	// Extract arguments object
	argsObj, ok := args["arguments"].(map[string]interface{})
	if !ok {
		// If no arguments provided, use empty map
		argsObj = make(map[string]interface{})
	}

	// Execute the tool
	if def.Execute == nil {
		return fmt.Sprintf("Error: tool '%s' has no executor", toolName), false
	}

	result, ok := def.Execute(argsObj, chatID)
	if !ok {
		return fmt.Sprintf("Tool '%s' failed: %s", toolName, result), false
	}

	return result, true
}

// executeToolList shows a compact list of all available tools
func ExecuteToolList(args map[string]interface{}, chatID int64) (string, bool) {
	category := helpers.GetStringArg(args, "category", "")

	// Group by category
	byCategory := make(map[string][]registry.ToolDef)
	for _, def := range registry.GetAllTools() {
		if category != "" && def.Category != category {
			continue
		}
		byCategory[def.Category] = append(byCategory[def.Category], def)
	}

	// Sort categories
	var cats []string
	for c := range byCategory {
		cats = append(cats, c)
	}
	sort.Strings(cats)

	var sb strings.Builder
	total := 0
	for _, c := range cats {
		tools := byCategory[c]
		total += len(tools)
		sb.WriteString(fmt.Sprintf("📂 **%s** (%d)\n", c, len(tools)))
		for _, t := range tools {
			marker := ""
			if t.Deferred {
				marker = " ⏸️"
			}
			sb.WriteString(fmt.Sprintf("  • %s%s\n", t.Name, marker))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("Total: %d tools (%d active, %d deferred)", total, countActiveTools(), countDeferredTools()))
	return sb.String(), true
}

func countActiveTools() int {
	n := 0
	for _, def := range registry.GetAllTools() {
		if def.Native && !def.Deferred {
			n++
		}
	}
	return n
}

func countDeferredTools() int {
	n := 0
	for _, def := range registry.GetAllTools() {
		if def.Deferred {
			n++
		}
	}
	return n
}
