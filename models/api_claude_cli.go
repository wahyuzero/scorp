package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// ──────────────────────────────────────────────
// Claude CLI Provider
// Bridges to local 'claude' CLI binary via subprocess
// ──────────────────────────────────────────────

type ClaudeCliProvider struct {
	Command string
}

func init() {
	p := &ClaudeCliProvider{Command: "claude"}
	RegisterProvider(ProviderSpec{
		Name:             "claude-cli",
		Aliases:          []string{"claude-code"},
		DisplayName:      "Claude Code CLI",
		DefaultBaseURL:   "local://claude-cli",
		DefaultAPIFormat: "claude-cli",
		NoAuth:           true,
	}, p)

	RegisterCatalog("claude-cli", []CatalogEntry{
		{"claude-3-7-sonnet", 16384, true, "claude-sonnet-cli"},
		{"claude-3-5-sonnet", 16384, true, "claude-35-sonnet-cli"},
		{"claude-3-5-haiku", 8192, false, "claude-haiku-cli"},
	})
}

func (p *ClaudeCliProvider) Format() string {
	return "claude-cli"
}

func (p *ClaudeCliProvider) Call(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	text, _, err := p.CallWithTools(ctx, model, messages)
	return text, err
}

func (p *ClaudeCliProvider) CallWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	cmdName := p.Command
	if cmdName == "" {
		cmdName = "claude"
	}

	// Check if CLI is available
	if _, err := exec.LookPath(cmdName); err != nil {
		return "", nil, fmt.Errorf("claude CLI not found in PATH — install via 'npm install -g @anthropic-ai/claude-code'")
	}

	prompt := formatMessagesForCLI(messages)

	args := []string{"-p", "--output-format", "json", "--dangerously-skip-permissions"}
	if model.Model != "" && model.Model != "claude-cli" && model.Model != "claude-code" {
		args = append(args, "--model", model.Model)
	}
	args = append(args, "-") // read prompt from stdin

	cmd := exec.CommandContext(ctx, cmdName, args...)
	cmd.Stdin = bytes.NewReader([]byte(prompt))

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		stderrStr := strings.TrimSpace(stderr.String())
		if stderrStr != "" {
			return "", nil, fmt.Errorf("claude CLI error: %s", stderrStr)
		}
		return "", nil, fmt.Errorf("claude CLI failed: %w", err)
	}

	output := stdout.String()

	// Parse JSON response from claude CLI
	var cliResp struct {
		Result  string `json:"result"`
		IsError bool   `json:"is_error"`
		Usage   struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal([]byte(output), &cliResp); err == nil {
		if cliResp.IsError {
			return "", nil, fmt.Errorf("claude CLI returned error: %s", cliResp.Result)
		}
		if cliResp.Usage.InputTokens > 0 || cliResp.Usage.OutputTokens > 0 {
			TrackModelUsage(model.Model, cliResp.Usage.InputTokens, cliResp.Usage.OutputTokens)
			RecordCost(model.Model, cliResp.Usage.InputTokens, cliResp.Usage.OutputTokens)
		}
		return cliResp.Result, nil, nil
	}

	// Fallback to raw stdout
	return strings.TrimSpace(output), nil, nil
}

func formatMessagesForCLI(messages []ChatMessage) string {
	var sb strings.Builder
	for _, m := range messages {
		switch m.Role {
		case "system":
			sb.WriteString("[System Instructions]:\n")
			sb.WriteString(m.Content)
			sb.WriteString("\n\n")
		case "user":
			sb.WriteString("User: ")
			sb.WriteString(m.Content)
			sb.WriteString("\n\n")
		case "assistant":
			sb.WriteString("Assistant: ")
			sb.WriteString(m.Content)
			sb.WriteString("\n\n")
		case "tool":
			sb.WriteString("[Tool Result]: ")
			sb.WriteString(m.Content)
			sb.WriteString("\n\n")
		}
	}
	return sb.String()
}
