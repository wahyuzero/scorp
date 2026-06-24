package registry

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

// ──────────────────────────────────────────────
// Tool Registry — Centralized tool definitions
// ──────────────────────────────────────────────

// ArgDef describes a tool argument
type ArgDef struct {
	Type        string   `json:"type"`          // string, integer, boolean, object, array
	Description string   `json:"description"`   // human-readable description
	Required    bool     `json:"required"`      // whether this arg is required
	Default     any      `json:"default,omitempty"`
	Enum        []string `json:"enum,omitempty"` // allowed values (for string)
}

// ToolDef describes a tool for the registry
type ToolDef struct {
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Arguments      map[string]ArgDef      `json:"arguments"`
	Execute        func(map[string]interface{}, int64) (string, bool)
	Native         bool                   `json:"native"`       // available via native function calling
	Category       string                 `json:"category"`     // shell, system, browser, mcp, vision, etc.
	Deferred       bool                   `json:"deferred"`     // if true, not sent to LLM schema (use tool_search + tool_call)
	RawInputSchema map[string]interface{} `json:"-"`            // raw JSON Schema (for MCP tools); overrides Arguments in native schema
}

// ToolSchema represents a tool definition for LLM API function calling
type ToolSchema struct {
	Type     string         `json:"type"` // always "function"
	Function ToolSchemaFunc `json:"function"`
}

// ToolSchemaFunc represents the function part of a tool schema
type ToolSchemaFunc struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// Global tool registry
var toolRegistry = make(map[string]ToolDef)

// RegisterTool adds a tool to the registry
func RegisterTool(def ToolDef) {
	if def.Name == "" {
		log.Printf("[registry] Warning: attempted to register tool with empty name")
		return
	}
	toolRegistry[def.Name] = def
	log.Printf("[registry] Registered tool: %s (%s)", def.Name, def.Category)
}

// GetTool retrieves a tool definition by name
func GetTool(name string) (ToolDef, bool) {
	def, ok := toolRegistry[name]
	return def, ok
}

// GetAllTools returns all registered tools
func GetAllTools() []ToolDef {
	tools := make([]ToolDef, 0, len(toolRegistry))
	for _, def := range toolRegistry {
		tools = append(tools, def)
	}
	return tools
}

// GetToolsByCategory returns tools filtered by category
func GetToolsByCategory(category string) []ToolDef {
	var tools []ToolDef
	for _, def := range toolRegistry {
		if def.Category == category {
			tools = append(tools, def)
		}
	}
	return tools
}

// ExecuteToolByName calls a tool from the registry
func ExecuteToolByName(name string, args map[string]interface{}, chatID int64) (string, bool) {
	def, ok := toolRegistry[name]
	if !ok {
		return "Unknown tool: " + name, false
	}
	return def.Execute(args, chatID)
}

// UnregisterTool removes a tool from the registry by name
func UnregisterTool(name string) bool {
	if _, ok := toolRegistry[name]; ok {
		delete(toolRegistry, name)
		return true
	}
	return false
}

// cached native tool schema
var (
	cachedNativeTools []ToolSchema
	cachedNativeOnce  sync.Once
)

// ResetNativeToolCache resets the sync.Once so the native tool schema is rebuilt
// on next call. Called after dynamically registering tools (e.g. MCP tools).
func ResetNativeToolCache() {
	cachedNativeOnce = sync.Once{}
	log.Printf("[registry] Native tool cache reset — will rebuild on next use")
}

// GenerateNativeToolsSchema returns tools formatted for native function calling API
// Uses sync.Once to cache the schema after first build.
func GenerateNativeToolsSchema() []ToolSchema {
	cachedNativeOnce.Do(func() {
		var tools []ToolSchema
		for _, def := range toolRegistry {
			if !def.Native || def.Deferred {
				continue
			}
			// If the tool has a raw JSON Schema (MCP tools), use it directly
			if def.RawInputSchema != nil {
				schema := def.RawInputSchema
				if _, ok := schema["type"]; !ok {
					schema = make(map[string]interface{})
					for k, v := range def.RawInputSchema {
						schema[k] = v
					}
					schema["type"] = "object"
				}
				tools = append(tools, ToolSchema{
					Type: "function",
					Function: ToolSchemaFunc{
						Name:        def.Name,
						Description: def.Description,
						Parameters:  schema,
					},
				})
				continue
			}
			// Build schema from ArgDef
			props := make(map[string]interface{})
			required := []string{}
			for argName, argDef := range def.Arguments {
				prop := map[string]interface{}{
					"type":        argDef.Type,
					"description": argDef.Description,
				}
				if len(argDef.Enum) > 0 {
					prop["enum"] = argDef.Enum
				}
				if argDef.Default != nil {
					prop["default"] = argDef.Default
				}
				props[argName] = prop
				if argDef.Required {
					required = append(required, argName)
				}
			}
			tools = append(tools, ToolSchema{
				Type: "function",
				Function: ToolSchemaFunc{
					Name:        def.Name,
					Description: def.Description,
					Parameters: map[string]interface{}{
						"type":       "object",
						"properties": props,
						"required":   required,
					},
				},
			})
		}
		cachedNativeTools = tools
		log.Printf("[registry] Cached native tool schema: %d tools", len(cachedNativeTools))
	})
	return cachedNativeTools
}

// GenerateSystemPromptDescriptions returns tool descriptions for system prompt
func GenerateSystemPromptDescriptions() string {
	var sb strings.Builder

	categories := []string{"shell", "system", "browser", "vision", "mcp", "code", "docker", "process", "network", "database", "http", "other"}

	for _, cat := range categories {
		tools := GetToolsByCategory(cat)
		if len(tools) == 0 {
			continue
		}

		sb.WriteString("\n### " + strings.ToUpper(cat) + " Tools\n")
		for _, t := range tools {
			sb.WriteString(fmt.Sprintf("- %s: %s\n", t.Name, t.Description))
			for argName, argDef := range t.Arguments {
				req := ""
				if argDef.Required {
					req = " (required)"
				}
				defStr := ""
				if argDef.Default != nil {
					defStr = fmt.Sprintf(" [default: %v]", argDef.Default)
				}
				enumStr := ""
				if len(argDef.Enum) > 0 {
					enumStr = fmt.Sprintf(" [enum: %s]", strings.Join(argDef.Enum, ", "))
				}
				sb.WriteString(fmt.Sprintf("  - %s (%s)%s%s%s: %s\n", argName, argDef.Type, req, defStr, enumStr, argDef.Description))
			}
		}
	}

	return sb.String()
}
