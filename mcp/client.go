package mcp

import (
	"scorp-agent/registry"
	"scorp-agent/internal/helpers"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)
// MCPConfig is the top-level mcp.json structure
type MCPConfig struct {
	MCPServers map[string]MCPServerConfig `json:"mcpServers"`
	// ServerMode config for when scorp-agent acts as an MCP server
	MCPServerMode *MCPServerModeConfig `json:"mcpServerMode,omitempty"`
}

// MCPServerModeConfig configures scorp-agent as an MCP server
type MCPServerModeConfig struct {
	Enabled     bool     `json:"enabled"`
	ExposedTools []string `json:"exposed_tools"` // list of tool names to expose
}

// MCPServerConfig is the config for a single MCP server in mcp.json
type MCPServerConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
}

// MCPTool represents a tool discovered from an MCP server
type MCPTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
	ServerName  string                 `json:"-"` // which server this belongs to
}

// MCPServer is a running MCP server process
type MCPServer struct {
	Name    string
	Config  MCPServerConfig
	cmd     *exec.Cmd
	stdin   *json.Encoder
	scanner *bufio.Scanner
	tools   []MCPTool
	mu      sync.Mutex
	reqID   int64
	alive   bool
}

// JSON-RPC 2.0 types
type jsonRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int64       `json:"id,omitempty"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type jsonRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int64           `json:"id,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *jsonRPCError   `json:"error,omitempty"`
}

type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Global MCP state
var (
	mcpServers   = make(map[string]*MCPServer)
	mcpServersMu sync.RWMutex
	mcpAllTools  []MCPTool // aggregated tools from all servers

	// Shutdown coordination
	mcpShutdownCtx    context.Context
	mcpShutdownCancel context.CancelFunc
	mcpShutdownWG     sync.WaitGroup
)

// ──────────────────────────────────────────────
// Lifecycle
// ──────────────────────────────────────────────

// LoadMCPConfig reads ~/.scorp/mcp.json
func LoadMCPConfig() (*MCPConfig, error) {
	path := os.ExpandEnv("$HOME") + "/.scorp/mcp.json"
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg MCPConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse mcp.json: %w", err)
	}
	return &cfg, nil
}

// StartMCPServers loads config and starts all configured MCP servers
func StartMCPServers() {
	// Create shutdown context with 10s timeout for graceful shutdown
	mcpShutdownCtx, mcpShutdownCancel = context.WithTimeout(context.Background(), 10*time.Second)

	cfg, err := LoadMCPConfig()
	if err != nil {
		log.Printf("[mcp] No MCP config: %v", err)
		return
	}

	for name, serverCfg := range cfg.MCPServers {
		srv, err := startMCPServer(name, serverCfg)
		if err != nil {
			log.Printf("[mcp] Failed to start %s: %v", name, err)
			continue
		}
		mcpServersMu.Lock()
		mcpServers[name] = srv
		mcpServersMu.Unlock()

		// Track each server in wait group for graceful shutdown
		mcpShutdownWG.Add(1)
		go func(s *MCPServer) {
			defer mcpShutdownWG.Done()
			<-mcpShutdownCtx.Done()
			s.Close()
		}(srv)

		log.Printf("[mcp] Started server %s with %d tools", name, len(srv.tools))
	}

	// Rebuild aggregated tool list
	rebuildMCPToolList()

	// Register all discovered MCP tools as first-class native tools
	// (like Hermes Agent: mcp_{server}_{tool})
	registerMCPToolsAsNative()
}

// StopMCPServers shuts down all running MCP servers gracefully
func StopMCPServers() {
	log.Println("[mcp] Initiating graceful shutdown of MCP servers...")

	// Cancel shutdown context to signal all server goroutines to close
	if mcpShutdownCancel != nil {
		mcpShutdownCancel()
	}

	// Wait for all server goroutines to complete (with timeout)
	done := make(chan struct{})
	go func() {
		mcpShutdownWG.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("[mcp] All MCP server goroutines stopped gracefully")
	case <-time.After(12 * time.Second):
		log.Println("[mcp] Timeout waiting for graceful shutdown, forcing kill")
	}

	mcpServersMu.Lock()
	defer mcpServersMu.Unlock()
	for name, srv := range mcpServers {
		srv.Close() // ensure cleanup
		log.Printf("[mcp] Stopped server %s", name)
	}
	mcpServers = make(map[string]*MCPServer)
	mcpAllTools = nil

	// Reset shutdown coordination
	mcpShutdownCtx = nil
	mcpShutdownCancel = nil
	log.Println("[mcp] MCP servers stopped")
}

// unregisterMCPNativeTools removes all MCP-native tools from the registry.
// Called before reload to prevent duplicate registrations.
func unregisterMCPNativeTools() {
	count := 0
	tools := registry.GetAllTools()
	for _, def := range tools {
		if def.Category == "mcp" {
			registry.UnregisterTool(def.Name)
			count++
		}
	}
	if count > 0 {
		log.Printf("[mcp] Unregistered %d old MCP native tools", count)
	}
}

// ReloadMCPServers restarts all MCP servers and re-registers native tools
func ReloadMCPServers() {
	StopMCPServers()
	unregisterMCPNativeTools()
	StartMCPServers()
}

// rebuildMCPToolList aggregates tools from all running servers
func rebuildMCPToolList() {
	mcpServersMu.RLock()
	defer mcpServersMu.RUnlock()
	var all []MCPTool
	for _, srv := range mcpServers {
		all = append(all, srv.tools...)
	}
	mcpAllTools = all
}

// GetMCPTools returns all discovered MCP tools
func GetMCPTools() []MCPTool {
	mcpServersMu.RLock()
	defer mcpServersMu.RUnlock()
	return mcpAllTools
}

// ──────────────────────────────────────────────
// First-class MCP tool registration (Hermes-style)
// ──────────────────────────────────────────────

// sanitizeMCPName replaces characters that are invalid in function names
func sanitizeMCPName(s string) string {
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, ".", "_")
	return s
}

// registerMCPToolsAsNative registers each discovered MCP tool as an individual
// first-class native tool in the registry, with the naming convention:
//   mcp_{server_name}_{tool_name}
// Hyphens and dots in names are replaced with underscores.
// The LLM can then call these tools directly without a generic dispatcher.
func registerMCPToolsAsNative() {
	mcpServersMu.RLock()
	defer mcpServersMu.RUnlock()

	registered := 0
	for serverName, srv := range mcpServers {
		for _, tool := range srv.tools {
			fullName := "mcp_" + sanitizeMCPName(serverName) + "_" + sanitizeMCPName(tool.Name)

			// Build registry.ArgDef from inputSchema for system prompt descriptions
			argDefs := buildArgDefsFromInputSchema(tool.InputSchema)

			// Capture for closure
			sn := serverName
			tn := tool.Name
			desc := tool.Description
			if len(desc) > 200 {
				desc = desc[:200] + "..."
			}

			registry.RegisterTool(registry.ToolDef{
				Name:           fullName,
				Description:    fmt.Sprintf("[MCP:%s] %s", serverName, desc),
				Category:       "mcp",
				Native:         true,
				Arguments:      argDefs,
				RawInputSchema: tool.InputSchema,
				Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
					mcpServersMu.RLock()
					s, ok := mcpServers[sn]
					mcpServersMu.RUnlock()
					if !ok {
						return fmt.Sprintf("Error: MCP server '%s' not connected", sn), false
					}
					result, err := s.CallTool(tn, args)
					if err != nil {
						return fmt.Sprintf("MCP tool error: %v", err), false
					}
					return helpers.TruncOutputTool(result), true
				},
			})
			registered++
		}
	}

	if registered > 0 {
		log.Printf("[mcp] Registered %d MCP tools as first-class native tools", registered)
		// Reset native tool cache so new tools are included
		registry.ResetNativeToolCache()
	}
}

// buildArgDefsFromInputSchema converts an MCP tool's JSON Schema inputSchema
// to a map of registry.ArgDef entries for system prompt descriptions.
func buildArgDefsFromInputSchema(schema map[string]interface{}) map[string]registry.ArgDef {
	args := make(map[string]registry.ArgDef)

	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		return args
	}

	// Get required fields
	requiredMap := make(map[string]bool)
	if reqList, ok := schema["required"].([]interface{}); ok {
		for _, r := range reqList {
			if s, ok := r.(string); ok {
				requiredMap[s] = true
			}
		}
	}

	for name, raw := range props {
		prop, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		typeStr, _ := prop["type"].(string)
		if typeStr == "" {
			typeStr = "string"
		}

		descStr, _ := prop["description"].(string)

		argDef := registry.ArgDef{
			Type:        typeStr,
			Description: descStr,
			Required:    requiredMap[name],
		}

		// Handle default
		if def, ok := prop["default"]; ok {
			argDef.Default = def
		}

		// Handle enum
		if enumList, ok := prop["enum"].([]interface{}); ok {
			for _, e := range enumList {
				if s, ok := e.(string); ok {
					argDef.Enum = append(argDef.Enum, s)
				}
			}
		}

		args[name] = argDef
	}

	return args
}

// ──────────────────────────────────────────────
// Server process management
// ──────────────────────────────────────────────

func startMCPServer(name string, cfg MCPServerConfig) (*MCPServer, error) {
	log.Printf("[mcp] Starting server %s: %s %s", name, cfg.Command, strings.Join(cfg.Args, " "))

	cmd := exec.Command(cfg.Command, cfg.Args...)

	// Set environment
	cmd.Env = os.Environ()
	for k, v := range cfg.Env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	// Setup stdio pipes
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("stdin pipe: %w", err)
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}

	// Discard stderr to avoid blocking
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start process: %w", err)
	}

	srv := &MCPServer{
		Name:    name,
		Config:  cfg,
		cmd:     cmd,
		stdin:   json.NewEncoder(stdinPipe),
		scanner: bufio.NewScanner(stdoutPipe),
		alive:   true,
	}

	// Set larger buffer for scanner (some responses can be big)
	srv.scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)

	// Initialize the connection
	if err := srv.initialize(); err != nil {
		cmd.Process.Kill()
		return nil, fmt.Errorf("initialize: %w", err)
	}

	// Discover tools
	tools, err := srv.listTools()
	if err != nil {
		log.Printf("[mcp] Warning: couldn't list tools for %s: %v", name, err)
	}
	srv.tools = tools

	return srv, nil
}

// Close shuts down the MCP server process gracefully
func (s *MCPServer) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.alive = false

	// Try graceful shutdown: send MCP shutdown notification if process still alive
	if s.cmd != nil && s.cmd.Process != nil {
		// Send shutdown notification (JSON-RPC 2.0) to stdin
		shutdownReq := jsonRPCRequest{
			JSONRPC: "2.0",
			Method:  "shutdown",
			Params:  nil,
		}
		if s.stdin != nil {
			s.stdin.Encode(shutdownReq)
			// Give it a moment to process shutdown
			time.Sleep(200 * time.Millisecond)
		}

		// Try graceful exit first
		s.cmd.Process.Signal(os.Interrupt)

		// Wait for graceful exit with timeout
		done := make(chan error, 1)
		go func() {
			done <- s.cmd.Wait()
		}()

		select {
		case <-done:
			// Process exited gracefully
			return
		case <-time.After(3 * time.Second):
			// Force kill after timeout
			s.cmd.Process.Kill()
			<-done // wait for kill to complete
		}
	}
}

// ──────────────────────────────────────────────
// JSON-RPC communication
// ──────────────────────────────────────────────

func (s *MCPServer) sendRequest(method string, params interface{}) (*jsonRPCResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.alive {
		return nil, fmt.Errorf("server %s is not alive", s.Name)
	}

	id := atomic.AddInt64(&s.reqID, 1)
	req := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}

	if err := s.stdin.Encode(req); err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	// Read response (skip notifications — lines without matching id)
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		if !s.scanner.Scan() {
			if err := s.scanner.Err(); err != nil {
				return nil, fmt.Errorf("read response: %w", err)
			}
			return nil, fmt.Errorf("server %s closed connection", s.Name)
		}

		line := s.scanner.Text()
		if line == "" {
			continue
		}

		var resp jsonRPCResponse
		if err := json.Unmarshal([]byte(line), &resp); err != nil {
			log.Printf("[mcp] Skipping non-JSON line from %s: %s", s.Name, helpers.TruncateStr(line, 100))
			continue
		}

		// Skip notifications (no id)
		if resp.ID == 0 && resp.Result == nil && resp.Error == nil {
			continue
		}

		// Skip responses for other requests
		if resp.ID != id {
			continue
		}

		if resp.Error != nil {
			return nil, fmt.Errorf("MCP error %d: %s", resp.Error.Code, resp.Error.Message)
		}

		return &resp, nil
	}

	return nil, fmt.Errorf("timeout waiting for response from %s", s.Name)
}

func (s *MCPServer) sendNotification(method string, params interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.alive {
		return fmt.Errorf("server %s is not alive", s.Name)
	}

	req := jsonRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	return s.stdin.Encode(req)
}

// ──────────────────────────────────────────────
// MCP Protocol methods
// ──────────────────────────────────────────────

func (s *MCPServer) initialize() error {
	params := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo": map[string]interface{}{
			"name":    "scorp-agent",
			"version": "1.0.0",
		},
	}

	resp, err := s.sendRequest("initialize", params)
	if err != nil {
		return err
	}

	log.Printf("[mcp] Initialized %s: %s", s.Name, string(resp.Result))

	// Send initialized notification
	return s.sendNotification("notifications/initialized", nil)
}

func (s *MCPServer) listTools() ([]MCPTool, error) {
	resp, err := s.sendRequest("tools/list", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	var result struct {
		Tools []MCPTool `json:"tools"`
	}
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("parse tools: %w", err)
	}

	// Tag tools with server name
	for i := range result.Tools {
		result.Tools[i].ServerName = s.Name
	}

	return result.Tools, nil
}

// CallTool calls a tool on this MCP server
func (s *MCPServer) CallTool(toolName string, args map[string]interface{}) (string, error) {
	params := map[string]interface{}{
		"name":      toolName,
		"arguments": args,
	}

	resp, err := s.sendRequest("tools/call", params)
	if err != nil {
		return "", err
	}

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		IsError bool `json:"isError"`
	}
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return "", fmt.Errorf("parse tool result: %w", err)
	}

	// Collect all text content
	var texts []string
	for _, c := range result.Content {
		if c.Type == "text" {
			texts = append(texts, c.Text)
		}
	}

	output := strings.Join(texts, "\n")

	if result.IsError {
		return "", fmt.Errorf("tool error: %s", output)
	}

	return output, nil
}

// ──────────────────────────────────────────────
// Helper: find server for a tool
// ──────────────────────────────────────────────

// FindMCPTool finds which server has a given tool name
func FindMCPTool(toolName string) (*MCPServer, *MCPTool) {
	mcpServersMu.RLock()
	defer mcpServersMu.RUnlock()
	for _, srv := range mcpServers {
		for i, t := range srv.tools {
			if t.Name == toolName {
				return srv, &srv.tools[i]
			}
		}
	}
	return nil, nil
}

// MCPToolsSummary returns a human-readable summary of all MCP tools for display
func MCPToolsSummary() string {
	mcpServersMu.RLock()
	defer mcpServersMu.RUnlock()

	if len(mcpServers) == 0 {
		return "No MCP servers configured."
	}

	var sb strings.Builder
	for name, srv := range mcpServers {
		status := "🟢"
		if !srv.alive {
			status = "🔴"
		}
		sb.WriteString(fmt.Sprintf("%s <b>%s</b> (%d tools)\n", status, name, len(srv.tools)))
		for _, t := range srv.tools {
			desc := t.Description
			if len(desc) > 60 {
				desc = desc[:60] + "..."
			}
			sb.WriteString(fmt.Sprintf("  • <code>%s</code> — %s\n", t.Name, desc))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// MCPToolsForPrompt returns a prompt-friendly description of all MCP tools
func MCPToolsForPrompt() string {
	tools := GetMCPTools()
	if len(tools) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\n\nYou also have access to MCP tools. To use them, call the 'mcp_tool' function with:\n")
	sb.WriteString("- server: the server name\n")
	sb.WriteString("- tool: the tool name\n")
	sb.WriteString("- arguments: JSON object of arguments\n\n")
	sb.WriteString("Available MCP tools:\n")

	for _, t := range tools {
		sb.WriteString(fmt.Sprintf("- [%s] %s: %s\n", t.ServerName, t.Name, t.Description))
		// Include input schema hints for required params
		if schema, ok := t.InputSchema["properties"]; ok {
			if props, ok := schema.(map[string]interface{}); ok {
				var paramNames []string
				for pn := range props {
					paramNames = append(paramNames, pn)
				}
				sb.WriteString(fmt.Sprintf("  Parameters: %s\n", strings.Join(paramNames, ", ")))
			}
		}
	}

	return sb.String()
}

// ──────────────────────────────────────────────
// MCP Server Mode — scorp-agent as an MCP server
// ──────────────────────────────────────────────

var mcpServerModeRunning atomic.Bool

// StartMCPServerMode starts scorp-agent as an MCP server over stdio
// This should be called in a goroutine as it blocks on stdin
func StartMCPServerMode() {
	cfg, err := LoadMCPConfig()
	if err != nil {
		log.Printf("[mcp] No MCP config for server mode: %v", err)
		return
	}

	if cfg.MCPServerMode == nil || !cfg.MCPServerMode.Enabled {
		log.Printf("[mcp] Server mode not enabled in config")
		return
	}

	mcpServerModeRunning.Store(true)
	log.Println("[mcp] Starting MCP server mode on stdio...")
	startMCPServerMode()
}

func StopMCPServerMode() {
	mcpServerModeRunning.Store(false)
	log.Println("[mcp] MCP server mode stopped")
}

// startMCPServerMode runs the JSON-RPC 2.0 server on stdin/stdout
func startMCPServerMode() {
	scanner := bufio.NewScanner(os.Stdin)
	for mcpServerModeRunning.Load() && scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var req mcpRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			sendMCPError(nil, -32700, "Parse error")
			continue
		}

		handleMCPRequest(req)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("[mcp] Stdin error: %v", err)
	}
}

// mcpRequest represents an incoming JSON-RPC 2.0 request
type mcpRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type mcpResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *mcpError       `json:"error,omitempty"`
}

type mcpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type mcpTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// getExposedTools returns the tools to expose based on config
func getExposedTools(cfg *MCPServerModeConfig) []mcpTool {
	allTools := []mcpTool{
		{
			Name:        "shell",
			Description: "Execute a shell command on the VPS",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"command": map[string]interface{}{"type": "string", "description": "Command to execute"},
					"timeout": map[string]interface{}{"type": "integer", "description": "Timeout seconds", "default": 30},
				},
				"required": []string{"command"},
			},
		},
		{
			Name:        "system_info",
			Description: "Get system information (CPU, memory, disk, network, docker, services)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"type": map[string]interface{}{"type": "string", "description": "Type: full, cpu, memory, disk, network, docker, services", "default": "full"},
				},
			},
		},
		{
			Name:        "search_code",
			Description: "Search codebase using ripgrep",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"pattern": map[string]interface{}{"type": "string", "description": "Regex pattern"},
					"path":    map[string]interface{}{"type": "string", "description": "Search path", "default": "."},
				},
				"required": []string{"pattern"},
			},
		},
		{
			Name:        "log",
			Description: "Fetch logs from docker, journal, or file",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"source": map[string]interface{}{"type": "string", "description": "docker, journal, or file"},
					"target": map[string]interface{}{"type": "string", "description": "Container/unit/file path"},
					"lines":  map[string]interface{}{"type": "integer", "description": "Lines to fetch", "default": 50},
				},
				"required": []string{"source", "target"},
			},
		},
	}

	if cfg == nil || len(cfg.ExposedTools) == 0 {
		return allTools
	}

	// Filter by exposed tools list
	exposed := make(map[string]bool)
	for _, t := range cfg.ExposedTools {
		exposed[t] = true
	}

	var filtered []mcpTool
	for _, t := range allTools {
		if exposed[t.Name] {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

func handleMCPRequest(req mcpRequest) {
	cfg, _ := LoadMCPConfig()
	tools := getExposedTools(cfg.MCPServerMode)

	switch req.Method {
	case "initialize":
		result := map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]string{
				"name":    "scorp-agent",
				"version": "1.0.0",
			},
		}
		sendMCPResult(req.ID, result)

	case "tools/list":
		result := map[string]interface{}{
			"tools": tools,
		}
		sendMCPResult(req.ID, result)

	case "tools/call":
		var params struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			sendMCPError(req.ID, -32602, "Invalid params: "+err.Error())
			return
		}

		result, ok := executeMCPServerTool(params.Name, params.Arguments)
		if !ok {
			sendMCPError(req.ID, -32603, "Tool execution failed: "+result)
			return
		}

		sendMCPResult(req.ID, map[string]interface{}{
			"content": []map[string]string{
				{"type": "text", "text": result},
			},
		})

	default:
		sendMCPError(req.ID, -32601, "Method not found: "+req.Method)
	}
}

func executeMCPServerTool(name string, args map[string]interface{}) (string, bool) {
	switch name {
	case "shell":
		return registry.ExecuteToolByName("shell", args, 0)
	case "system_info":
		res, _ := registry.ExecuteToolByName("system_info", args, 0)
		return res, true
	case "search_code":
		res, _ := registry.ExecuteToolByName("search_code", args, 0)
		return res, true
	case "log":
		res, _ := registry.ExecuteToolByName("log", args, 0)
		return res, true
	default:
		return "Unknown tool: " + name, false
	}
}

func sendMCPResult(id json.RawMessage, result interface{}) {
	resp := mcpResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	data, _ := json.Marshal(resp)
	fmt.Println(string(data))
}

func sendMCPError(id json.RawMessage, code int, message string) {
	resp := mcpResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &mcpError{Code: code, Message: message},
	}
	data, _ := json.Marshal(resp)
	fmt.Println(string(data))
}
