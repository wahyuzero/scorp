package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ──────────────────────────────────────────────
// execute_code tool — Python orchestration
// Runs Python scripts that can call scorp-agent tools
// ──────────────────────────────────────────────

// scorpToolsModule is injected into the Python subprocess via PYTHONPATH.
// It provides wrapper functions that call back into the agent's tool system.
const scorpToolsPy = `
import json, os, subprocess, sys, shlex, re, time
from pathlib import Path

_tool_calls = 0
_MAX_TOOL_CALLS = 50

def _check_limit():
    global _tool_calls
    _tool_calls += 1
    if _tool_calls > _MAX_TOOL_CALLS:
        raise RuntimeError(f"Max tool calls ({_MAX_TOOL_CALLS}) exceeded")

def _call_tool(name, **kwargs):
    """Internal: call a tool via the bridge file protocol."""
    _check_limit()
    bridge_dir = os.environ.get('SCORP_BRIDGE_DIR', '/tmp/scorp_bridge')
    req_id = str(os.getpid()) + '_' + str(_tool_calls)
    req_path = os.path.join(bridge_dir, f'req_{req_id}.json')
    resp_path = os.path.join(bridge_dir, f'resp_{req_id}.json')
    
    req = {'tool': name, 'args': kwargs}
    with open(req_path, 'w') as f:
        json.dump(req, f)
    
    # Wait for response (timeout 30s)
    for _ in range(300):
        if os.path.exists(resp_path):
            with open(resp_path) as f:
                resp = json.load(f)
            os.remove(resp_path)
            os.remove(req_path)
            if not resp.get('success', True):
                raise RuntimeError(resp.get('error', 'Tool failed'))
            return resp.get('result', '')
        time.sleep(0.1)
    
    raise TimeoutError(f"Tool '{name}' timed out")

# ── File Operations ──

def read_file(path, offset=1, limit=500):
    """Read file contents. Lines are 1-indexed."""
    result = _call_tool('read_file', path=path, offset=offset, limit=limit)
    return result

def write_file(path, content):
    """Write content to file (overwrites)."""
    return _call_tool('write_file', path=path, content=content)

def search_files(pattern, target="content", path=".", file_glob=None, limit=50):
    """Search files by content (regex) or by name (glob)."""
    kwargs = {'pattern': pattern, 'target': target, 'path': path, 'limit': limit}
    if file_glob:
        kwargs['file_glob'] = file_glob
    return _call_tool('search_code', **kwargs)

def patch(path, old_string, new_string, replace_all=False):
    """Find and replace text in a file."""
    return _call_tool('patch', path=path, old_string=old_string, new_string=new_string, replace_all=replace_all)

# ── Shell ──

def terminal(command, timeout=180):
    """Execute a shell command and return stdout."""
    result = _call_tool('shell', command=command, timeout=timeout)
    return result

# ── HTTP ──

def http_request(method, url, headers=None, body=None):
    """Make HTTP request."""
    kwargs = {'method': method, 'url': url}
    if headers:
        kwargs['headers'] = headers
    if body:
        kwargs['body'] = body
    return _call_tool('http', **kwargs)

# ── Memory ──

def memory(action, key=None, value=None):
    """KV memory store."""
    kwargs = {'action': action}
    if key:
        kwargs['key'] = key
    if value:
        kwargs['value'] = value
    return _call_tool('memory', **kwargs)

# ── Helpers ──

def json_parse(text):
    """Parse JSON with lenient settings."""
    return json.loads(text, strict=False)

def shell_quote(s):
    """Quote a string for safe shell usage."""
    return shlex.quote(str(s))

def retry(fn, max_attempts=3, delay=2):
    """Retry a function with exponential backoff."""
    last_err = None
    for i in range(max_attempts):
        try:
            return fn()
        except Exception as e:
            last_err = e
            if i < max_attempts - 1:
                time.sleep(delay * (2 ** i))
    raise last_err
`

// executeCodeTool runs a Python script with access to scorp tools.
func executeCodeTool(args map[string]interface{}, chatID int64) (string, bool) {
	code, _ := args["code"].(string)
	if code == "" {
		return "Missing 'code' parameter", false
	}

	// Create bridge directory for IPC
	bridgeDir := fmt.Sprintf("/tmp/scorp_bridge_%d_%d", chatID, time.Now().UnixNano())
	if err := os.MkdirAll(bridgeDir, 0755); err != nil {
		return fmt.Sprintf("Failed to create bridge dir: %v", err), false
	}
	defer os.RemoveAll(bridgeDir)

	// Write the scorp_tools module
	toolsPath := filepath.Join(bridgeDir, "scorp_tools.py")
	if err := os.WriteFile(toolsPath, []byte(scorpToolsPy), 0644); err != nil {
		return fmt.Sprintf("Failed to write tools module: %v", err), false
	}

	// Write user script
	scriptPath := filepath.Join(bridgeDir, "user_script.py")
	header := "import sys; sys.path.insert(0, '" + bridgeDir + "')\nfrom scorp_tools import *\n\n"
	if err := os.WriteFile(scriptPath, []byte(header+code), 0644); err != nil {
		return fmt.Sprintf("Failed to write script: %v", err), false
	}

	// Start Python process with env vars
	pythonBin := "/usr/bin/python3"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, pythonBin, scriptPath)
	cmd.Env = append(os.Environ(),
		"PYTHONPATH="+bridgeDir,
		"SCORP_BRIDGE_DIR="+bridgeDir,
	)

	// Start a goroutine to service tool requests from the Python script
	go serviceBridgeRequests(bridgeDir, chatID, ctx)

	// Set stdout/stderr capture with 50KB cap
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()

	// Cap output at 50KB
	stdout := stdoutBuf.String()
	if len(stdout) > 50000 {
		stdout = stdout[:50000] + "\n... (truncated at 50KB)"
	}
	stderr := stderrBuf.String()

	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Sprintf("⏰ Script timed out (5 min limit)\n\n**stdout:**\n%s\n**stderr:**\n%s", stdout, stderr), true
	}

	if err != nil {
		// Still return stdout if we have it
		result := ""
		if stdout != "" {
			result = stdout + "\n\n"
		}
		result += fmt.Sprintf("❌ Error: %v\n", err)
		if stderr != "" {
			result += fmt.Sprintf("**stderr:**\n%s", stderr)
		}
		return result, false
	}

	if stderr != "" {
		return fmt.Sprintf("%s\n\n⚠️ stderr:\n%s", stdout, stderr), true
	}

	return stdout, true
}

// serviceBridgeRequests polls the bridge directory for tool requests from the Python script
// and executes them using the agent's tool system.
func serviceBridgeRequests(bridgeDir string, chatID int64, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		entries, err := os.ReadDir(bridgeDir)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		for _, entry := range entries {
			name := entry.Name()
			if !strings.HasPrefix(name, "req_") {
				continue
			}

			reqPath := filepath.Join(bridgeDir, name)
			// Derive response path
			respName := strings.Replace(name, "req_", "resp_", 1)
			respPath := filepath.Join(bridgeDir, respName)

			// Skip if response already exists
			if _, err := os.Stat(respPath); err == nil {
				continue
			}

			// Read request
			reqData, err := os.ReadFile(reqPath)
			if err != nil {
				continue
			}

			var req struct {
				Tool string                 `json:"tool"`
				Args map[string]interface{} `json:"args"`
			}
			if err := json.Unmarshal(reqData, &req); err != nil {
				writeBridgeResponse(respPath, "", fmt.Sprintf("parse error: %v", err))
				continue
			}

			// Execute tool
			tc := ToolCall{Name: req.Tool, Args: req.Args}
			result, ok := executeTool(tc, chatID)

			writeBridgeResponse(respPath, result, func() string {
				if !ok {
					return "tool returned error"
				}
				return ""
			}())
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func writeBridgeResponse(path, result, errMsg string) {
	type bridgeResp struct {
		Result  string `json:"result"`
		Success bool   `json:"success"`
		Error   string `json:"error,omitempty"`
	}
	resp := bridgeResp{
		Result:  result,
		Success: errMsg == "",
		Error:   errMsg,
	}
	data, _ := json.Marshal(resp)
	os.WriteFile(path, data, 0644)
}

func init() {
	registerTool(ToolDef{
		Name:        "execute_code",
		Description: "Run a Python script with access to agent tools (read_file, write_file, terminal, search_files, patch, http_request, memory, json_parse, shell_quote, retry). 5-min timeout, 50KB stdout cap, max 50 tool calls. Use for batch operations, loops, data processing.",
		Category:    "code",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return executeCodeTool(args, chatID)
		},
		Arguments: map[string]ArgDef{
			"code": {Type: "string", Description: "Python code to execute. Tools are pre-imported: read_file(), write_file(), terminal(), search_files(), patch(), http_request(), memory(), json_parse(), shell_quote(), retry()", Required: true},
		},
	})
}
