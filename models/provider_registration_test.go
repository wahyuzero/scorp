package models

import (
	"context"
	"strings"
	"testing"

	"scorp-agent/registry"
)

// MockCustomProvider demonstrates how a 3rd-party or custom provider
// can be added as a standalone .go file with zero core modifications.
type MockCustomProvider struct {
	CalledCall       bool
	CalledCallTools  bool
	CalledStream     bool
	CalledResolveKey bool
}

func (m *MockCustomProvider) Format() string {
	return "mock-vendor"
}

func (m *MockCustomProvider) Call(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	m.CalledCall = true
	return "mock response from " + model.Model, nil
}

func (m *MockCustomProvider) CallWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	m.CalledCallTools = true
	return "mock tool response", []ToolCall{{Name: "test_tool", Args: map[string]interface{}{"q": "ok"}}}, nil
}

func (m *MockCustomProvider) CallStream(ctx context.Context, model *ModelConfig, messages []ChatMessage) (<-chan StreamChunk, error) {
	m.CalledStream = true
	ch := make(chan StreamChunk, 2)
	go func() {
		defer close(ch)
		ch <- StreamChunk{Content: "mock stream token"}
		ch <- StreamChunk{Finish: true}
	}()
	return ch, nil
}

func (m *MockCustomProvider) ResolveKey(cfg *ModelConfig) string {
	m.CalledResolveKey = true
	return "mock-resolved-secret-key-123"
}

func TestDynamicProviderRegistration(t *testing.T) {
	mock := &MockCustomProvider{}

	// Register dynamic provider with spec
	RegisterProvider(ProviderSpec{
		Name:             "mock-vendor",
		Aliases:          []string{"mockv", "mock-ai"},
		DisplayName:      "Mock Vendor AI",
		DefaultBaseURL:   "https://api.mockvendor.ai/v1",
		DefaultAPIFormat: "mock-vendor",
		KeyEnvs:          []string{"MOCK_VENDOR_API_KEY"},
	}, mock)

	// 1. Check provider retrieval
	p := GetProvider("mock-vendor")
	if p == nil || p.Format() != "mock-vendor" {
		t.Fatalf("expected mock-vendor provider, got %v", p)
	}

	// 2. Check alias retrieval
	pAlias := GetProvider("mockv")
	if pAlias == nil || pAlias.Format() != "mock-vendor" {
		t.Fatalf("expected alias mockv to resolve to mock-vendor, got %v", pAlias)
	}

	// 3. Check format resolution
	cfg := &ModelConfig{
		Provider: "mock-vendor",
		Model:    "mock-v1",
	}
	format := ResolveAPIFormat(cfg)
	if format != "mock-vendor" {
		t.Fatalf("expected format mock-vendor, got %s", format)
	}

	// 4. Check base URL resolution from spec
	baseURL := ResolveBaseURL(cfg)
	if baseURL != "https://api.mockvendor.ai/v1" {
		t.Fatalf("expected baseURL https://api.mockvendor.ai/v1, got %s", baseURL)
	}

	// 5. Check dynamic KeyResolver delegation
	apiKey := ResolveAPIKey(cfg)
	if apiKey != "mock-resolved-secret-key-123" {
		t.Fatalf("expected key 'mock-resolved-secret-key-123', got %s", apiKey)
	}
	if !mock.CalledResolveKey {
		t.Fatalf("expected KeyResolver.ResolveKey to be called")
	}

	// 6. Check CallModel dispatch
	reply, err := CallModel(context.Background(), cfg, []ChatMessage{{Role: "user", Content: "hello"}})
	if err != nil || !strings.Contains(reply, "mock response") {
		t.Fatalf("unexpected CallModel reply: %s, err: %v", reply, err)
	}
	if !mock.CalledCall {
		t.Fatalf("expected Call to be invoked")
	}

	// 7. Check CallModelWithTools dispatch
	replyWithTools, tools, err := CallModelWithTools(context.Background(), cfg, []ChatMessage{{Role: "user", Content: "call tool"}})
	if err != nil || len(tools) != 1 || tools[0].Name != "test_tool" {
		t.Fatalf("unexpected CallModelWithTools result: %s, tools: %v, err: %v", replyWithTools, tools, err)
	}
	if !mock.CalledCallTools {
		t.Fatalf("expected CallWithTools to be invoked")
	}

	// 8. Check CallModelStream dynamic StreamingProvider dispatch
	streamCh, err := CallModelStream(context.Background(), cfg, []ChatMessage{{Role: "user", Content: "stream"}})
	if err != nil {
		t.Fatalf("CallModelStream error: %v", err)
	}
	var received []string
	for chunk := range streamCh {
		if chunk.Content != "" {
			received = append(received, chunk.Content)
		}
	}
	if len(received) == 0 || received[0] != "mock stream token" {
		t.Fatalf("unexpected stream chunks: %v", received)
	}
	if !mock.CalledStream {
		t.Fatalf("expected StreamingProvider.CallStream to be invoked")
	}
}

func TestToolSchemaTransform(t *testing.T) {
	rawSchema := map[string]interface{}{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"title":   "TestParams",
		"type":    "object",
		"properties": map[string]interface{}{
			"prompt": map[string]interface{}{
				"type":        "string",
				"description": "text prompt",
				"anyOf": []interface{}{
					map[string]interface{}{"type": "string"},
					map[string]interface{}{"type": "null"},
				},
			},
		},
		"additionalProperties": false,
	}

	sanitized := SanitizeToolSchema(rawSchema, "simple")

	// Verify $schema and title stripped
	if _, hasSchema := sanitized["$schema"]; hasSchema {
		t.Errorf("expected $schema to be stripped")
	}
	if _, hasTitle := sanitized["title"]; hasTitle {
		t.Errorf("expected title to be stripped")
	}

	// Verify properties preserved and sanitized
	props, ok := sanitized["properties"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected properties map")
	}
	promptProp, ok := props["prompt"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected prompt property map")
	}
	if _, hasAnyOf := promptProp["anyOf"]; hasAnyOf {
		t.Errorf("expected anyOf to be flattened in simple mode")
	}

	// Test TransformToolDefinitions helper
	tools := []registry.ToolSchema{
		{
			Type: "function",
			Function: registry.ToolSchemaFunc{
				Name:        "sample",
				Description: "sample desc",
				Parameters:  rawSchema,
			},
		},
	}
	transformed := TransformToolDefinitions(tools, "simple")
	if len(transformed) != 1 {
		t.Fatalf("expected 1 transformed tool")
	}
	if _, hasSchema := transformed[0].Function.Parameters["$schema"]; hasSchema {
		t.Errorf("expected $schema to be stripped in transformed tool parameters")
	}
}

func TestCoreAndExtendedProvidersRegistered(t *testing.T) {
	expectedProviders := []string{
		"openai",
		"anthropic",
		"gemini",
		"command-code",
		"opencode",
		"azure",
		"claude-cli",
		"mistral",
		"xai",
		"grok",
		"perplexity",
		"cohere",
		"qwen",
		"dashscope",
		"together",
		"fireworks",
		"siliconflow",
		"cerebras",
		"novita",
		"deepinfra",
		"deepseek",
		"groq",
		"openrouter",
		"kimi",
		"minimax",
		"nvidia",
		"ollama",
		"lmstudio",
		"zai",
	}

	for _, pName := range expectedProviders {
		cfg := &ModelConfig{Provider: pName, Model: "test-model"}
		format := ResolveAPIFormat(cfg)
		if format == "" {
			t.Errorf("provider %s failed to resolve format", pName)
		}
		p := GetProvider(format)
		if p == nil {
			t.Errorf("provider %s (format %s) returned nil adapter", pName, format)
		}
	}
}
