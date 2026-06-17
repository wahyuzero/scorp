package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// ──────────────────────────────────────────────
// ACP (Agent Communication Protocol) Delegation
//
// Subprocess transport: delegate tasks to external CLI agents
// (Claude Code, OpenAI Codex, OpenCode) via JSON-RPC over stdio.
// ──────────────────────────────────────────────

// ACPRequest is a JSON-RPC 2.0 request sent to the subprocess.
type ACPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// ACPResponse is a JSON-RPC 2.0 response from the subprocess.
type ACPResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *ACPError       `json:"error,omitempty"`
}

// ACPError is a JSON-RPC error object.
type ACPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ACPInitializeParams for the "initialize" method.
type ACPInitializeParams struct {
	ProtocolVersion int `json:"protocolVersion"`
	Capabilities    struct {
		Streaming bool `json:"streaming"`
	} `json:"capabilities"`
}

// ACPMessageNewParams for the "message/new" method.
type ACPMessageNewParams struct {
	SessionID string        `json:"sessionId,omitempty"`
	Message   ACPUserMessage `json:"message"`
}

type ACPUserMessage struct {
	Role    string        `json:"role"`
	Parts   []ACPMessagePart `json:"parts"`
}

type ACPMessagePart struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// ACPSession holds a connected subprocess session.
type ACPSession struct {
	cmd      *exec.Cmd
	stdin    io.WriteCloser
	stdout   io.ReadCloser
	stderr   io.ReadCloser
	mu       sync.Mutex
	nextID   int
	pending  map[int]chan *ACPResponse
	sessionID string
	done     chan struct{}
}

// launchACP starts the external CLI subprocess.
func launchACP(ctx context.Context, command string, args []string) (*ACPSession, error) {
	cmd := exec.CommandContext(ctx, command, args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start %s: %w", command, err)
	}

	session := &ACPSession{
		cmd:     cmd,
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
		nextID:  1,
		pending: make(map[int]chan *ACPResponse),
		done:    make(chan struct{}),
	}

	// Read responses from stdout
	go session.readLoop()

	// Log stderr
	go session.logStderr()

	// Initialize protocol
	if err := session.initialize(); err != nil {
		session.Close()
		return nil, fmt.Errorf("ACP initialize failed: %w", err)
	}

	log.Printf("[acp] Subprocess %s started and initialized", command)
	return session, nil
}

func (s *ACPSession) readLoop() {
	scanner := bufio.NewScanner(s.stdout)
	buf := make([]byte, 0, 256*1024)
	scanner.Buffer(buf, 4*1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var resp ACPResponse
		if err := json.Unmarshal([]byte(line), &resp); err != nil {
			continue
		}

		s.mu.Lock()
		ch, ok := s.pending[resp.ID]
		if ok {
			delete(s.pending, resp.ID)
		}
		s.mu.Unlock()

		if ok {
			ch <- &resp
		}
	}
}

func (s *ACPSession) logStderr() {
	scanner := bufio.NewScanner(s.stderr)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			log.Printf("[acp:stderr] %s", line)
		}
	}
}

func (s *ACPSession) send(method string, params interface{}) (*ACPResponse, error) {
	s.mu.Lock()
	id := s.nextID
	s.nextID++
	ch := make(chan *ACPResponse, 1)
	s.pending[id] = ch
	s.mu.Unlock()

	var paramsRaw json.RawMessage
	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("marshal params: %w", err)
		}
		paramsRaw = data
	}

	req := ACPRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  paramsRaw,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	data = append(data, '\n')

	s.mu.Lock()
	_, err = s.stdin.Write(data)
	s.mu.Unlock()

	if err != nil {
		s.mu.Lock()
		delete(s.pending, id)
		s.mu.Unlock()
		return nil, fmt.Errorf("write to subprocess: %w", err)
	}

	select {
	case resp := <-ch:
		return resp, nil
	case <-time.After(5 * time.Minute):
		s.mu.Lock()
		delete(s.pending, id)
		s.mu.Unlock()
		return nil, fmt.Errorf("ACP request timeout (method=%s id=%d)", method, id)
	}
}

func (s *ACPSession) initialize() error {
	params := ACPInitializeParams{
		ProtocolVersion: 1,
	}
	params.Capabilities.Streaming = false

	resp, err := s.send("initialize", params)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return fmt.Errorf("init error: %s", resp.Error.Message)
	}

	// Send initialized notification (no response expected)
	s.mu.Lock()
notifData, _ := json.Marshal(ACPRequest{
		JSONRPC: "2.0",
		Method:  "notifications/initialized",
	})
	notifData = append(notifData, '\n')
	s.stdin.Write(notifData)
	s.mu.Unlock()

	return nil
}

// SendMessage sends a user message and returns the text response.
func (s *ACPSession) SendMessage(prompt string) (string, error) {
	params := ACPMessageNewParams{
		SessionID: s.sessionID,
		Message: ACPUserMessage{
			Role: "user",
			Parts: []ACPMessagePart{
				{Type: "text", Text: prompt},
			},
		},
	}

	resp, err := s.send("message/new", params)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", fmt.Errorf("message error: %s", resp.Error.Message)
	}

	// Parse result — extract text from parts
	var result struct {
		SessionID string `json:"sessionId"`
		Message   struct {
			Role  string `json:"role"`
			Parts []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"message"`
	}
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		// Fallback: return raw result
		return string(resp.Result), nil
	}

	s.sessionID = result.SessionID

	var sb strings.Builder
	for _, part := range result.Message.Parts {
		if part.Type == "text" {
			sb.WriteString(part.Text)
		}
	}
	return sb.String(), nil
}

func (s *ACPSession) Close() {
	select {
	case <-s.done:
		return
	default:
	}
	close(s.done)

	if s.stdin != nil {
		s.stdin.Close()
	}
	if s.cmd != nil && s.cmd.Process != nil {
		s.cmd.Process.Kill()
		s.cmd.Wait()
	}
}

// ──────────────────────────────────────────────
// ACP Delegate Execution
// ──────────────────────────────────────────────

// runSubagentACP executes a task via external CLI subprocess.
func runSubagentACP(params delegateTaskParams) delegateResult {
	start := time.Now()
	subagentID := fmt.Sprintf("acp_%d", start.UnixNano())

	if params.ACPCommand == "" {
		return delegateResult{
			SubagentID: subagentID,
			Task:       params.Task,
			Status:     "failed",
			Error:      "acp_command is empty",
			Duration:   time.Since(start),
		}
	}

	// Build the full prompt with context
	fullPrompt := params.Task
	if params.Context != "" {
		fullPrompt = fmt.Sprintf("Context:\n%s\n\nTask:\n%s", params.Context, params.Task)
	}
	if len(params.Tools) > 0 {
		fullPrompt += fmt.Sprintf("\n\nAvailable tools: %s", strings.Join(params.Tools, ", "))
	}

	// Launch subprocess
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	args := params.ACPArgs
	if len(args) == 0 {
		args = []string{"--acp", "--stdio"}
	}

	session, err := launchACP(ctx, params.ACPCommand, args)
	if err != nil {
		return delegateResult{
			SubagentID: subagentID,
			Task:       params.Task,
			Role:       params.Role,
			Status:     "failed",
			Error:      fmt.Sprintf("launch failed: %v", err),
			Duration:   time.Since(start),
		}
	}
	defer session.Close()

	// Send task
	result, err := session.SendMessage(fullPrompt)
	if err != nil {
		return delegateResult{
			SubagentID: subagentID,
			Task:       params.Task,
			Role:       params.Role,
			Status:     "failed",
			Error:      fmt.Sprintf("message failed: %v", err),
			Duration:   time.Since(start),
		}
	}

	modelLabel := params.ACPCommand
	return delegateResult{
		SubagentID: subagentID,
		Task:       params.Task,
		Role:       params.Role + " (ACP)",
		ModelUsed:  modelLabel,
		Status:     "completed",
		Result:     result,
		ToolsUsed:  []string{params.ACPCommand},
		Iterations: 1,
		Duration:   time.Since(start),
	}
}

// knownACPCommands maps CLI names to their binary paths.
var knownACPCommands = map[string]string{
	"claude":  "claude",
	"codex":   "codex",
	"opencode": "opencode",
}

// checkACPAvailable checks if a CLI is installed and available.
func checkACPAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// listAvailableACP returns list of available ACP agents.
func listAvailableACP() []string {
	var available []string
	for name := range knownACPCommands {
		if checkACPAvailable(name) {
			available = append(available, name)
		}
	}
	return available
}
