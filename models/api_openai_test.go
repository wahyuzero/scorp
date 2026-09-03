package models

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestLLMProviderFactory(t *testing.T) {
	// 1. Check default OpenAI provider
	pOpenAI := GetProvider("openai")
	if pOpenAI == nil || pOpenAI.Format() != "openai" {
		t.Errorf("Expected openai provider, got %v", pOpenAI)
	}

	// 2. Check Anthropic provider
	pAnthropic := GetProvider("anthropic")
	if pAnthropic == nil || pAnthropic.Format() != "anthropic" {
		t.Errorf("Expected anthropic provider, got %v", pAnthropic)
	}

	// 3. Check Gemini provider
	pGemini := GetProvider("gemini")
	if pGemini == nil || pGemini.Format() != "gemini" {
		t.Errorf("Expected gemini provider, got %v", pGemini)
	}

	// 4. Check Command Code provider
	pCC := GetProvider("command-code")
	if pCC == nil || pCC.Format() != "command-code" {
		t.Errorf("Expected command-code provider, got %v", pCC)
	}

	// 5. Unknown defaults to OpenAI
	pDefault := GetProvider("unknown_custom_engine")
	if pDefault == nil || pDefault.Format() != "openai" {
		t.Errorf("Expected fallback to openai for unknown format, got %v", pDefault)
	}
}

func TestCallOpenAI_MockServer(t *testing.T) {
	// Mock OpenAI API server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/chat/completions") {
			http.Error(w, "Not found", 404)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"choices": [
				{
					"message": {
						"role": "assistant",
						"content": "Hello from Mock OpenAI"
					},
					"finish_reason": "stop"
				}
			],
			"usage": {
				"prompt_tokens": 10,
				"completion_tokens": 5,
				"total_tokens": 15
			}
		}`))
	}))
	defer ts.Close()

	os.Setenv("TEST_OPENAI_KEY", "sk-mocktest12345678901234567890")
	defer os.Unsetenv("TEST_OPENAI_KEY")

	model := &ModelConfig{
		Provider:  "openai",
		Model:     "gpt-4o-mini",
		BaseURL:   ts.URL,
		KeyEnv:    "TEST_OPENAI_KEY",
		MaxTokens: 1000,
	}

	reply, err := CallModel(context.Background(), model, []ChatMessage{
		{Role: "user", Content: "hi"},
	})
	if err != nil {
		t.Fatalf("CallModel failed: %v", err)
	}
	if reply != "Hello from Mock OpenAI" {
		t.Errorf("Expected 'Hello from Mock OpenAI', got '%s'", reply)
	}
}
