package models

import (
	"encoding/json"
	"strings"
	"testing"

	"scorp-agent/registry"
)

func TestBuildCommandCodePayload(t *testing.T) {
	model := &ModelConfig{
		Provider: "command-code",
		Model:    "deepseek/deepseek-v4-flash",
	}

	messages := []ChatMessage{
		{Role: "system", Content: "You are a helpful coding agent."},
		{Role: "user", Content: "Hello scorp!"},
		{
			Role:    "assistant",
			Content: "I will check the directory.",
			ToolCalls: []ToolCallResp{
				{
					ID:   "call_123",
					Type: "function",
					Function: struct {
						Name      string `json:"name"`
						Arguments string `json:"arguments"`
					}{
						Name:      "list_dir",
						Arguments: `{"path":"/tmp"}`,
					},
				},
			},
		},
		{
			Role:    "tool",
			Content: "a.txt, b.txt",
		},
	}

	tools := []registry.ToolSchema{
		{
			Type: "function",
			Function: registry.ToolSchemaFunc{
				Name:        "list_dir",
				Description: "List directory contents",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{"type": "string"},
					},
				},
			},
		},
	}

	payload, err := buildCommandCodePayload(model, messages, tools)
	if err != nil {
		t.Fatalf("unexpected error building payload: %v", err)
	}

	// 1. Verify Memory contains system prompt
	if payload.Memory != "You are a helpful coding agent." {
		t.Errorf("expected Memory to be system prompt, got: %q", payload.Memory)
	}

	// 2. Verify messages array does NOT contain role 'system'
	for _, m := range payload.Params.Messages {
		if m.Role == "system" {
			t.Errorf("messages array must never contain role 'system' (Command Code Zod violation)")
		}
	}

	// 3. Verify assistant tool call was converted
	if len(payload.Params.Messages) < 3 {
		t.Fatalf("expected at least 3 messages, got %d", len(payload.Params.Messages))
	}
	asstMsg := payload.Params.Messages[1]
	if asstMsg.Role != "assistant" {
		t.Errorf("expected role assistant, got %s", asstMsg.Role)
	}

	// 4. Verify tool schema was mapped
	if len(payload.Params.Tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(payload.Params.Tools))
	}
	if payload.Params.Tools[0].Name != "list_dir" {
		t.Errorf("expected tool name list_dir, got %s", payload.Params.Tools[0].Name)
	}

	// 5. Verify payload serializes cleanly to JSON
	bytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}
	if !strings.Contains(string(bytes), "deepseek/deepseek-v4-flash") {
		t.Errorf("expected json to contain model id")
	}
}

func TestExtractFallbackToolCalls(t *testing.T) {
	// Scenario 1: DSML tags
	dsmlInput := `<｜｜DSML｜｜tool_calls>
<｜｜DSML｜｜invoke name="read_file">
<｜｜DSML｜｜parameter name="path" string="true">/tmp/demo.txt</｜｜DSML｜｜parameter>
</｜｜DSML｜｜invoke>
</｜｜DSML｜｜tool_calls>
Teks penjelasan lain.`

	calls, clean := extractFallbackToolCalls(dsmlInput)
	if len(calls) != 1 {
		t.Fatalf("expected 1 DSML tool call, got %d", len(calls))
	}
	if calls[0].Name != "read_file" {
		t.Errorf("expected tool name read_file, got %s", calls[0].Name)
	}
	if calls[0].Args["path"] != "/tmp/demo.txt" {
		t.Errorf("expected path /tmp/demo.txt, got %v", calls[0].Args["path"])
	}
	if strings.Contains(clean, "DSML") {
		t.Errorf("expected clean text without DSML, got: %s", clean)
	}

	// Scenario 2: Generic XML <tool_call>
	xmlInput := `I'll run the command:
<tool_call>
{"name": "shell", "arguments": {"command": "ls -la"}}
</tool_call>`

	calls2, clean2 := extractFallbackToolCalls(xmlInput)
	if len(calls2) != 1 {
		t.Fatalf("expected 1 XML tool call, got %d", len(calls2))
	}
	if calls2[0].Name != "shell" {
		t.Errorf("expected tool name shell, got %s", calls2[0].Name)
	}
	if calls2[0].Args["command"] != "ls -la" {
		t.Errorf("expected command 'ls -la', got %v", calls2[0].Args["command"])
	}
	if strings.Contains(clean2, "<tool_call>") {
		t.Errorf("expected clean text without XML tags, got: %s", clean2)
	}

	// Scenario 3: Pseudo bracket format [Tool Call: ...]
	bracketInput := `Let me check.
[Tool Call: list_dir({"path": "/home"})]`

	calls3, clean3 := extractFallbackToolCalls(bracketInput)
	if len(calls3) != 1 {
		t.Fatalf("expected 1 bracket tool call, got %d", len(calls3))
	}
	if calls3[0].Name != "list_dir" {
		t.Errorf("expected tool name list_dir, got %s", calls3[0].Name)
	}
	if calls3[0].Args["path"] != "/home" {
		t.Errorf("expected path '/home', got %v", calls3[0].Args["path"])
	}
	if strings.Contains(clean3, "[Tool Call:") {
		t.Errorf("expected clean text without bracket call, got: %s", clean3)
	}
}

func TestSwitchActiveModel(t *testing.T) {
	ModelCfg = defaultModelConfig()

	// Switch to pro
	err := SwitchActiveModel("deepseek/deepseek-v4-pro")
	if err != nil {
		t.Fatalf("failed to switch model: %v", err)
	}

	if ModelCfg.AgentModel != "deepseek/deepseek-v4-pro" {
		t.Errorf("expected agent model deepseek/deepseek-v4-pro, got %s", ModelCfg.AgentModel)
	}

	// Switch to invalid
	err = SwitchActiveModel("nonexistent/model")
	if err == nil {
		t.Errorf("expected error for nonexistent model, got nil")
	}

	// Switch back to flash
	_ = SwitchActiveModel("deepseek/deepseek-v4-flash")
	if ModelCfg.AgentModel != "deepseek/deepseek-v4-flash" {
		t.Errorf("expected agent model deepseek/deepseek-v4-flash, got %s", ModelCfg.AgentModel)
	}
}
