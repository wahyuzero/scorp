package registry

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

// ──────────────────────────────────────────────
// Tool Registry — Centralized tool definitions
// ──────────────────────────────────────────────

// ArgDef describes a tool argument
type ArgDef struct {
	Type        string   `json:"type"`        // string, integer, boolean, object, array
	Description string   `json:"description"` // human-readable description
	Required    bool     `json:"required"`    // whether this arg is required
	Default     any      `json:"default,omitempty"`
	Enum        []string `json:"enum,omitempty"` // allowed values (for string)
}

// ToolDef describes a tool for the registry
type ToolDef struct {
	Name           string                                             `json:"name"`
	Description    string                                             `json:"description"`
	Arguments      map[string]ArgDef                                  `json:"arguments"`
	Execute        func(map[string]interface{}, int64) (string, bool) `json:"-"`
	Native         bool                                               `json:"native"`   // available via native function calling
	Category       string                                             `json:"category"` // shell, system, browser, mcp, vision, etc.
	Deferred       bool                                               `json:"deferred"` // if true, not sent to LLM schema (use tool_search + tool_call)
	RawInputSchema map[string]interface{}                             `json:"-"`        // raw JSON Schema (for MCP tools); overrides Arguments in native schema
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

// Global thread-safe tool registry
var (
	toolRegistry   = make(map[string]ToolDef)
	toolRegistryMu sync.RWMutex
	maxCustomTools = 200 // Cap dynamically registered tools to prevent memory exhaustion
)

// RegisterTool adds a tool to the registry in a thread-safe manner
func RegisterTool(def ToolDef) {
	if def.Name == "" {
		log.Printf("[registry] Warning: attempted to register tool with empty name")
		return
	}
	toolRegistryMu.Lock()
	defer toolRegistryMu.Unlock()

	// If registry exceeds max custom tools and tool is dynamic/MCP, prevent unbounded growth
	if len(toolRegistry) >= maxCustomTools && !def.Native {
		log.Printf("[registry] Warning: registry capacity (%d) reached, skipping tool %s", maxCustomTools, def.Name)
		return
	}

	toolRegistry[def.Name] = def
	ResetNativeToolCacheLocked()
	if os.Getenv("SCORP_DEBUG") != "" {
		log.Printf("[registry] Registered tool: %s (%s)", def.Name, def.Category)
	}
}

// GetTool retrieves a tool definition by name safely
func GetTool(name string) (ToolDef, bool) {
	toolRegistryMu.RLock()
	defer toolRegistryMu.RUnlock()
	def, ok := toolRegistry[name]
	return def, ok
}

// GetAllTools returns all registered tools safely
func GetAllTools() []ToolDef {
	toolRegistryMu.RLock()
	defer toolRegistryMu.RUnlock()
	tools := make([]ToolDef, 0, len(toolRegistry))
	for _, def := range toolRegistry {
		tools = append(tools, def)
	}
	return tools
}

// GetToolsByCategory returns tools filtered by category safely
func GetToolsByCategory(category string) []ToolDef {
	toolRegistryMu.RLock()
	defer toolRegistryMu.RUnlock()
	var tools []ToolDef
	for _, def := range toolRegistry {
		if def.Category == category {
			tools = append(tools, def)
		}
	}
	return tools
}

// ExecuteToolByName calls a tool from the registry safely
func ExecuteToolByName(name string, args map[string]interface{}, chatID int64) (string, bool) {
	toolRegistryMu.RLock()
	def, ok := toolRegistry[name]
	toolRegistryMu.RUnlock()
	if !ok {
		return "Unknown tool: " + name, false
	}
	return def.Execute(args, chatID)
}

// UnregisterTool removes a tool from the registry by name safely
func UnregisterTool(name string) bool {
	toolRegistryMu.Lock()
	defer toolRegistryMu.Unlock()
	if _, ok := toolRegistry[name]; ok {
		delete(toolRegistry, name)
		ResetNativeToolCacheLocked()
		ClearToolTTL(name)
		return true
	}
	return false
}

// cached native tool schema
var (
	cachedNativeTools []ToolSchema
	cachedNativeMu    sync.RWMutex
	cachedNativeValid bool
)

// ResetNativeToolCache marks the cached native tool schema as invalid so it is rebuilt on next call.
func ResetNativeToolCache() {
	cachedNativeMu.Lock()
	cachedNativeValid = false
	cachedNativeMu.Unlock()
	if os.Getenv("SCORP_DEBUG") != "" {
		log.Printf("[registry] Native tool cache invalidated — will rebuild on next use")
	}
}

// ResetNativeToolCacheLocked invalidates cache when caller already holds registry locks
func ResetNativeToolCacheLocked() {
	cachedNativeMu.Lock()
	cachedNativeValid = false
	cachedNativeMu.Unlock()
}

// GenerateNativeToolsSchema returns tools formatted for native function calling API with RWMutex protection
func GenerateNativeToolsSchema() []ToolSchema {
	cachedNativeMu.RLock()
	if cachedNativeValid {
		res := make([]ToolSchema, len(cachedNativeTools))
		copy(res, cachedNativeTools)
		cachedNativeMu.RUnlock()
		return res
	}
	cachedNativeMu.RUnlock()

	cachedNativeMu.Lock()
	defer cachedNativeMu.Unlock()

	if cachedNativeValid {
		res := make([]ToolSchema, len(cachedNativeTools))
		copy(res, cachedNativeTools)
		return res
	}

	toolRegistryMu.RLock()
	defs := make([]ToolDef, 0, len(toolRegistry))
	for _, def := range toolRegistry {
		defs = append(defs, def)
	}
	toolRegistryMu.RUnlock()

	var tools []ToolSchema
	for _, def := range defs {
		if !IsToolActive(def) {
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
	cachedNativeValid = true
	log.Printf("[registry] Cached native tool schema: %d tools", len(cachedNativeTools))

	res := make([]ToolSchema, len(cachedNativeTools))
	copy(res, cachedNativeTools)
	return res
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
