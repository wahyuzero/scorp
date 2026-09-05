package mcp

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"scorp-agent/config"
	"scorp-agent/internal/helpers"
	"strings"
)

// executeMCPManage handles: list, add, remove, reload of MCP servers
func ExecuteMCPManage(args map[string]interface{}) (string, bool) {
	action := helpers.GetStringArg(args, "action", "list")

	switch action {
	case "list":
		return mcpManageList()
	case "add":
		return mcpManageAdd(args)
	case "remove":
		return mcpManageRemove(args)
	case "reload":
		return mcpManageReload()
	default:
		return fmt.Sprintf("Unknown action '%s'. Use: list, add, remove, reload", action), false
	}
}

// AddServerEntry writes a new server entry to ~/.scorp/mcp.json and hot-reloads
// so the new server (and its native tools) come online immediately. This is the
// terminal state of every installation path — mcp_manage add, the marketplace
// Tri-Option installer, and transpiler rebuilds all converge here.
func AddServerEntry(name string, cfg MCPServerConfig) (string, bool) {
	if name == "" {
		return "Error: server name is required", false
	}
	if cfg.Command == "" && cfg.URL == "" {
		return "Error: command or url is required", false
	}

	configPath := config.MCPConfigFilePath()
	mcpcfg, err := LoadMCPConfig()
	if err != nil {
		return fmt.Sprintf("Error loading config: %v", err), false
	}

	if mcpcfg.MCPServers == nil {
		mcpcfg.MCPServers = make(map[string]MCPServerConfig)
	}

	if _, exists := mcpcfg.MCPServers[name]; exists {
		return fmt.Sprintf("Error: MCP server '%s' already exists. Remove it first or use reload.", name), false
	}

	// Verify command exists (stdio servers only; remote URLs have no command)
	if cfg.Command != "" {
		if _, err := exec.LookPath(cfg.Command); err != nil {
			log.Printf("[mcp] Warning: command '%s' not in PATH: %v", cfg.Command, err)
		}
	}

	mcpcfg.MCPServers[name] = cfg

	data, err := json.MarshalIndent(mcpcfg, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error encoding config: %v", err), false
	}
	if err := os.MkdirAll(config.ScorpDir(), 0o755); err != nil {
		return fmt.Sprintf("Error creating config dir: %v", err), false
	}
	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		return fmt.Sprintf("Error writing config: %v", err), false
	}

	sb := fmt.Sprintf("✅ Added MCP server '%s' to config.\n", name)
	sb += fmt.Sprintf("   Command: %s %s\n", cfg.Command, strings.Join(cfg.Args, " "))

	sb += "\n🔄 Reloading MCP servers...\n"
	ReloadMCPServers()

	mcpServersMu.RLock()
	srv, ok := mcpServers[name]
	mcpServersMu.RUnlock()

	if ok {
		sb += fmt.Sprintf("\n✅ Server '%s' is running with %d tools:\n", name, len(srv.tools))
		for _, t := range srv.tools {
			nativeName := "mcp_" + sanitizeMCPName(name) + "_" + sanitizeMCPName(t.Name)
			sb += fmt.Sprintf("  • %s\n", nativeName)
		}
	} else {
		sb += fmt.Sprintf("\n⚠️ Server '%s' failed to start. Check logs: journalctl -u scorp-agent -n 20\n", name)
	}

	return sb, true
}

// RemoveServerEntry deletes a server entry from ~/.scorp/mcp.json and reloads.
func RemoveServerEntry(name string) (string, bool) {
	configPath := config.MCPConfigFilePath()
	mcpcfg, err := LoadMCPConfig()
	if err != nil {
		return fmt.Sprintf("Error loading config: %v", err), false
	}

	if _, exists := mcpcfg.MCPServers[name]; !exists {
		return fmt.Sprintf("Error: MCP server '%s' not found in config", name), false
	}

	delete(mcpcfg.MCPServers, name)

	data, err := json.MarshalIndent(mcpcfg, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error encoding config: %v", err), false
	}
	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		return fmt.Sprintf("Error writing config: %v", err), false
	}

	sb := fmt.Sprintf("✅ Removed MCP server '%s' from config.\n", name)
	sb += "\n🔄 Reloading MCP servers...\n"
	ReloadMCPServers()
	sb += "\n✅ Done. Remaining servers reloaded.\n"
	return sb, true
}

// mcpManageList shows all configured MCP servers (from config + running state)
func mcpManageList() (string, bool) {
	// Load config
	cfg, err := LoadMCPConfig()
	if err != nil {
		return fmt.Sprintf("Error loading config: %v", err), false
	}

	var sb strings.Builder
	sb.WriteString("📋 MCP Server Configuration\n\n")

	if len(cfg.MCPServers) == 0 {
		sb.WriteString("No MCP servers configured.\n")
		sb.WriteString("\nTo add one:\n")
		sb.WriteString("  mcp_manage(action=\"add\", name=\"myserver\", command=\"npx\", args=[\"-y\", \"@some/mcp-server\"])\n")
		return sb.String(), true
	}

	// Show running state
	mcpServersMu.RLock()
	defer mcpServersMu.RUnlock()

	for name, serverCfg := range cfg.MCPServers {
		srv, isRunning := mcpServers[name]
		status := "🔴 stopped"
		toolCount := 0
		if isRunning {
			status = "🟢 running"
			toolCount = len(srv.tools)
		}

		sb.WriteString(fmt.Sprintf("%s %s (%d tools)\n", status, name, toolCount))
		sb.WriteString(fmt.Sprintf("  Command: %s %s\n", serverCfg.Command, strings.Join(serverCfg.Args, " ")))

		// Show registered native tool names
		if isRunning {
			for _, t := range srv.tools {
				nativeName := "mcp_" + sanitizeMCPName(name) + "_" + sanitizeMCPName(t.Name)
				desc := t.Description
				if len(desc) > 50 {
					desc = desc[:50] + "..."
				}
				sb.WriteString(fmt.Sprintf("  • %s — %s\n", nativeName, desc))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String(), true
}

// mcpManageAdd adds a new MCP server to config, then hot-reloads
func mcpManageAdd(args map[string]interface{}) (string, bool) {
	name := helpers.GetStringArg(args, "name", "")
	if name == "" {
		return "Error: 'name' is required", false
	}
	command := helpers.GetStringArg(args, "command", "")
	if command == "" {
		return "Error: 'command' is required", false
	}

	// Get args array
	cmdArgs := helpers.GetStringSliceArg(args, "args")

	// Get env map
	envMap := make(map[string]string)
	if rawEnv, ok := args["env"]; ok {
		if envJSON, err := json.Marshal(rawEnv); err == nil {
			json.Unmarshal(envJSON, &envMap)
		}
	}

	// Also accept remote URL configs (SSE/HTTP transport)
	url := helpers.GetStringArg(args, "url", "")
	transport := helpers.GetStringArg(args, "transport", "")

	cfg := MCPServerConfig{
		Command:   command,
		Args:      cmdArgs,
		Env:       envMap,
		URL:       url,
		Transport: transport,
	}
	if url != "" {
		cfg.Command = ""
	}

	return AddServerEntry(name, cfg)
}

// mcpManageRemove removes an MCP server from config, then hot-reloads
func mcpManageRemove(args map[string]interface{}) (string, bool) {
	name := helpers.GetStringArg(args, "name", "")
	if name == "" {
		return "Error: 'name' is required", false
	}
	return RemoveServerEntry(name)
}

// mcpManageReload hot-reloads all MCP servers without restarting scorp-agent
func mcpManageReload() (string, bool) {
	sb := "🔄 Hot-reloading MCP servers...\n\n"
	ReloadMCPServers()

	// Show results
	mcpServersMu.RLock()
	defer mcpServersMu.RUnlock()

	if len(mcpServers) == 0 {
		sb += "No MCP servers running after reload."
		return sb, true
	}

	for name, srv := range mcpServers {
		status := "🟢"
		if !srv.alive {
			status = "🔴"
		}
		sb += fmt.Sprintf("%s %s (%d tools)\n", status, name, len(srv.tools))
		for _, t := range srv.tools {
			nativeName := "mcp_" + sanitizeMCPName(name) + "_" + sanitizeMCPName(t.Name)
			sb += fmt.Sprintf("  • %s\n", nativeName)
		}
	}

	return sb, true
}
