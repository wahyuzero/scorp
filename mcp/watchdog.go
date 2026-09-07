package mcp

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// ──────────────────────────────────────────────
// MCP Server Watchdog & Health Monitor
// Continuously monitors running MCP servers, detects unexpected crashes/EOF,
// and performs automated recovery with exponential backoff.
// ──────────────────────────────────────────────

type ServerWatchdog struct {
	serverName   string
	restartCount int
	maxRestarts  int
	lastCrash    time.Time
	// current is the MCPServer instance this monitor run is responsible for.
	// A monitor whose process was superseded (restart) or whose server was
	// intentionally closed must not count the exit as a crash.
	current *MCPServer
	stopped bool
	mu      sync.Mutex
	stopCh  chan struct{}
}

var (
	watchdogsMu sync.Mutex
	watchdogs   = make(map[string]*ServerWatchdog)
)

// RegisterWatchdog starts monitoring an MCP server process
func RegisterWatchdog(serverName string, srv *MCPServer) {
	watchdogsMu.Lock()
	defer watchdogsMu.Unlock()

	if wd, exists := watchdogs[serverName]; exists {
		// Reuse the existing watchdog so restartCount persists across
		// restarts — recreating it reset the counter to zero and a
		// fast-crashing server restarted forever (always "Attempt 1/5").
		// The previous monitor run exits on its own: it waits on the old
		// process, and once reaped it sees current != old and bails without
		// counting a crash.
		wd.mu.Lock()
		wd.current = srv
		wd.mu.Unlock()
		go wd.monitor(srv)
		return
	}

	wd := &ServerWatchdog{
		serverName:  serverName,
		maxRestarts: 5,
		current:     srv,
		stopCh:      make(chan struct{}),
	}
	watchdogs[serverName] = wd

	go wd.monitor(srv)
}

// StopWatchdogs marks every watchdog stopped so intentional shutdowns are
// not counted as crashes. Called from StopMCPServers before closing servers.
func StopWatchdogs() {
	watchdogsMu.Lock()
	defer watchdogsMu.Unlock()
	for _, wd := range watchdogs {
		wd.mu.Lock()
		wd.stopped = true
		wd.mu.Unlock()
	}
}

func (wd *ServerWatchdog) monitor(srv *MCPServer) {
	if srv.cmd == nil || srv.cmd.Process == nil {
		return
	}

	// Single-flight reap: shares cmd.Wait() with Close() instead of racing a
	// second Wait on the same process.
	err := srv.reapWait()

	if srv.wasClosed() {
		return
	}

	wd.mu.Lock()
	superseded := wd.current != srv
	stopped := wd.stopped
	wd.mu.Unlock()
	if superseded || stopped {
		return
	}

	wd.mu.Lock()
	wd.lastCrash = time.Now()
	wd.restartCount++
	count := wd.restartCount
	wd.mu.Unlock()

	log.Printf("[mcp-watchdog] Server '%s' crashed (exit: %v). Attempt %d/%d",
		wd.serverName, err, count, wd.maxRestarts)

	// Mark server dead
	mcpServersMu.Lock()
	if s, ok := mcpServers[wd.serverName]; ok && s == srv {
		s.alive = false
	}
	mcpServersMu.Unlock()

	if count > wd.maxRestarts {
		log.Printf("[mcp-watchdog] Server '%s' exceeded max restarts. Disabling watchdog.", wd.serverName)
		wd.mu.Lock()
		wd.stopped = true
		wd.mu.Unlock()
		return
	}

	// Exponential backoff: 1s, 2s, 4s, 8s... cap 30s
	backoff := time.Duration(1<<uint(count-1)) * time.Second
	if backoff > 30*time.Second {
		backoff = 30 * time.Second
	}
	log.Printf("[mcp-watchdog] Restarting '%s' in %s...", wd.serverName, backoff)

	select {
	case <-wd.stopCh:
		return
	case <-time.After(backoff):
	}

	wd.mu.Lock()
	stopped = wd.stopped
	wd.mu.Unlock()
	if stopped || srv.wasClosed() {
		return
	}

	// Attempt restart
	cfg, err := LoadMCPConfig()
	if err != nil {
		log.Printf("[mcp-watchdog] Config reload error: %v", err)
		return
	}

	serverCfg, exists := cfg.MCPServers[wd.serverName]
	if !exists {
		log.Printf("[mcp-watchdog] Server '%s' no longer in config.", wd.serverName)
		wd.mu.Lock()
		wd.stopped = true
		wd.mu.Unlock()
		return
	}

	newSrv, err := startMCPServer(wd.serverName, serverCfg)
	if err != nil {
		log.Printf("[mcp-watchdog] Failed to restart '%s': %v", wd.serverName, err)
		// Trigger next retry — stay the current monitor so the counter
		// keeps incrementing toward maxRestarts.
		go wd.monitor(srv)
		return
	}

	mcpServersMu.Lock()
	mcpServers[wd.serverName] = newSrv
	mcpServersMu.Unlock()

	// Re-register tools
	registerMCPToolsAsNative()
	log.Printf("[mcp-watchdog] Server '%s' successfully recovered and running!", wd.serverName)

	// Hand monitoring to a fresh run bound to the new process. This reuses
	// this watchdog (counter persists) and supersedes THIS run.
	RegisterWatchdog(wd.serverName, newSrv)
}

// RestartServer manually forces an MCP server restart
func RestartServer(serverName string) error {
	cfg, err := LoadMCPConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	serverCfg, ok := cfg.MCPServers[serverName]
	if !ok {
		return fmt.Errorf("server '%s' not configured", serverName)
	}

	// Stop existing
	mcpServersMu.Lock()
	if existing, exists := mcpServers[serverName]; exists {
		existing.Close()
		delete(mcpServers, serverName)
	}
	mcpServersMu.Unlock()

	// Start new
	newSrv, err := startMCPServer(serverName, serverCfg)
	if err != nil {
		return fmt.Errorf("start server '%s': %w", serverName, err)
	}

	mcpServersMu.Lock()
	mcpServers[serverName] = newSrv
	mcpServersMu.Unlock()

	registerMCPToolsAsNative()
	RegisterWatchdog(serverName, newSrv)

	return nil
}

// GetServerHealthStatus returns formatted health indicators for all MCP servers
func GetServerHealthStatus() string {
	cfg, err := LoadMCPConfig()
	if err != nil {
		return fmt.Sprintf("Error loading config: %v", err)
	}

	if len(cfg.MCPServers) == 0 {
		return "📋 No MCP servers configured.\nAdd servers in <code>~/.scorp/mcp.json</code>"
	}

	mcpServersMu.RLock()
	defer mcpServersMu.RUnlock()

	watchdogsMu.Lock()
	defer watchdogsMu.Unlock()

	var sb strings.Builder
	sb.WriteString("🔌 <b>MCP Servers Health & Watchdog</b>\n\n")

	for name := range cfg.MCPServers {
		srv, isRunning := mcpServers[name]
		status := "🔴 Stopped"
		toolsCount := 0
		if isRunning && srv.alive {
			status = "🟢 Healthy"
			toolsCount = len(srv.tools)
		}

		restarts := 0
		if wd, hasWD := watchdogs[name]; hasWD {
			restarts = wd.restartCount
		}

		sb.WriteString(fmt.Sprintf("%s <b>%s</b> — %d tools (restarts: %d)\n", status, name, toolsCount, restarts))
	}

	sb.WriteString("\nCommands:\n• <code>/mcp list</code>\n• <code>/mcp restart &lt;server&gt;</code>")
	return sb.String()
}
