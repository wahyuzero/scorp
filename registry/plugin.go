package registry

import (
	"context"
	"encoding/json"
	"fmt"
)

// ──────────────────────────────────────────────
// ToolPlugin — Clean, standard Go interface for plugins
// Designed for extreme simplicity and easy open-source contributions
// (as outlined in COMPETITIVE_ANALYSIS_ZEROCLAW.md)
// ──────────────────────────────────────────────

// ToolPlugin is the standard interface for external/contributor tools.
type ToolPlugin interface {
	Name() string
	Description() string
	Execute(ctx context.Context, args json.RawMessage) (string, error)
}

// ToolPluginWithSchema is an optional interface for plugins that declare custom JSON schema parameters.
type ToolPluginWithSchema interface {
	ToolPlugin
	Parameters() map[string]interface{}
}

// RegisterPlugin adapts a ToolPlugin into a registry ToolDef and registers it.
func RegisterPlugin(plugin ToolPlugin) {
	if plugin == nil || plugin.Name() == "" {
		return
	}

	toolDef := ToolDef{
		Name:        plugin.Name(),
		Description: plugin.Description(),
		Category:    "plugin",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			rawJSON, err := json.Marshal(args)
			if err != nil {
				return fmt.Sprintf("Error serializing arguments: %v", err), false
			}

			res, err := plugin.Execute(context.Background(), rawJSON)
			if err != nil {
				return fmt.Sprintf("Plugin error: %v", err), false
			}
			return res, true
		},
	}

	if schemaPlugin, ok := plugin.(ToolPluginWithSchema); ok {
		toolDef.RawInputSchema = schemaPlugin.Parameters()
	}

	RegisterTool(toolDef)
	ResetNativeToolCache()
}
