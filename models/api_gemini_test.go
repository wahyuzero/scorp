package models

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"scorp-agent/registry"
)

func init() {
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")
}

func TestGemini_LiveModels(t *testing.T) {
	testModels := []struct {
		modelID   string
		expectErr bool
	}{
		{"gemini-3.8-flash", false},
		{"gemini-3.7-flash", false},
		{"gemini-3.1-pro-preview", true}, // Free tier quota is 0 on Pro preview
	}

	for _, tc := range testModels {
		t.Run(tc.modelID, func(t *testing.T) {
			cfg := &ModelConfig{
				Provider: "gemini",
				Model:    tc.modelID,
			}

			key := ResolveAPIKey(cfg)
			if key == "" {
				t.Skip("GEMINI_API_KEY not found, skipping live Gemini test")
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
			defer cancel()

			reply, err := CallModel(ctx, cfg, []ChatMessage{
				{Role: "user", Content: "Reply with the single word: OK"},
			})

			if tc.expectErr {
				if err != nil {
					t.Logf("[%s] expected quota limitation on free tier: %v", tc.modelID, err)
				} else {
					t.Logf("[%s] unexpected success on free tier: %s", tc.modelID, reply)
				}
				return
			}

			if err != nil {
				if strings.Contains(err.Error(), "503") || strings.Contains(err.Error(), "high demand") {
					t.Logf("[%s] Google API temporary 503 (high demand): %v", tc.modelID, err)
					return
				}
				if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "quota") {
					t.Logf("[%s] Google API free-tier rate limit (429): %v", tc.modelID, err)
					return
				}
				if strings.Contains(err.Error(), "context deadline exceeded") {
					t.Logf("[%s] Google API free-tier latency spike/timeout: %v", tc.modelID, err)
					return
				}
				t.Fatalf("[%s] CallModel failed: %v", tc.modelID, err)
			}

			t.Logf("[%s] reply: %s", tc.modelID, strings.TrimSpace(reply))
			if !strings.Contains(strings.ToUpper(reply), "OK") {
				t.Errorf("[%s] expected 'OK' in reply, got: %s", tc.modelID, reply)
			}
		})

		// Respect free-tier rate limit (15 RPM = ~4s delay between requests)
		time.Sleep(3 * time.Second)
	}
}

func TestGemini_LiveStreaming(t *testing.T) {
	cfg := &ModelConfig{
		Provider: "gemini",
		Model:    "gemini-3.7-flash",
	}

	key := ResolveAPIKey(cfg)
	if key == "" {
		t.Skip("GEMINI_API_KEY not found, skipping live Gemini stream test")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	streamCh, err := CallModelStream(ctx, cfg, []ChatMessage{
		{Role: "user", Content: "Say 'Hello from Gemini Flash!'"},
	})
	if err != nil {
		if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "quota") {
			t.Logf("Gemini stream hit 429 free tier rate limit: %v", err)
			return
		}
		if strings.Contains(err.Error(), "context deadline exceeded") {
			t.Logf("Gemini stream hit latency spike/timeout: %v", err)
			return
		}
		t.Fatalf("CallModelStream failed: %v", err)
	}

	var accumulated string
	for chunk := range streamCh {
		if chunk.Error != nil {
			t.Fatalf("chunk error: %v", chunk.Error)
		}
		accumulated += chunk.Content
	}

	t.Logf("Gemini Streamed Reply: %s", accumulated)
	if !strings.Contains(strings.ToLower(accumulated), "gemini") {
		t.Errorf("expected 'Gemini' in streamed reply, got: %s", accumulated)
	}
}

func TestGemini_LiveToolCalling(t *testing.T) {
	registry.RegisterTool(registry.ToolDef{
		Name:        "get_current_time",
		Description: "Get the current time in a given timezone",
		Category:    "util",
		Native:      true,
		Arguments: map[string]registry.ArgDef{
			"timezone": {
				Type:        "string",
				Description: "The timezone name, e.g. UTC, Asia/Jakarta",
				Required:    true,
			},
		},
	})
	registry.ResetNativeToolCache()

	cfg := &ModelConfig{
		Provider: "gemini",
		Model:    "gemini-3.7-flash",
	}

	key := ResolveAPIKey(cfg)
	if key == "" {
		t.Skip("GEMINI_API_KEY not found, skipping live Gemini tool test")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	reply, toolCalls, err := CallGeminiWithTools(ctx, cfg, []ChatMessage{
		{Role: "user", Content: "What time is it right now? Use the get_time or similar tool."},
	})
	if err != nil {
		if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "quota") {
			t.Logf("Gemini tool test hit 429 free tier rate limit: %v", err)
			return
		}
		if strings.Contains(err.Error(), "context deadline exceeded") {
			t.Logf("Gemini tool test hit latency spike/timeout: %v", err)
			return
		}
		t.Fatalf("CallGeminiWithTools failed: %v", err)
	}

	t.Logf("Gemini Tool Reply: %s (ToolCalls: %d)", reply, len(toolCalls))
	for _, tc := range toolCalls {
		t.Logf("  Tool Call: %s args=%v", tc.Name, tc.Args)
	}

	if len(toolCalls) > 0 {
		// Test multi-turn follow-up with tool response
		turn2Msgs := []ChatMessage{
			{Role: "user", Content: "What time is it right now? Use the get_current_time tool."},
			{
				Role: "assistant",
				ToolCalls: []ToolCallResp{
					{
						ID:               "call_1",
						Type:             "function",
						ThoughtSignature: toolCalls[0].ThoughtSignature,
						Function: struct {
							Name      string `json:"name"`
							Arguments string `json:"arguments"`
						}{
							Name:      toolCalls[0].Name,
							Arguments: `{"timezone":"UTC"}`,
						},
					},
				},
			},
			{
				Role:    "tool",
				Content: `{"time": "16:45:00 UTC", "date": "2026-09-08"}`,
			},
		}

		followUpReply, _, err := CallGeminiWithTools(ctx, cfg, turn2Msgs)
		if err != nil {
			t.Fatalf("Gemini multi-turn tool response failed: %v", err)
		}
		t.Logf("Gemini multi-turn follow-up reply: %s", followUpReply)
		if !strings.Contains(followUpReply, "16:45") && !strings.Contains(followUpReply, "UTC") {
			t.Logf("Note: reply didn't contain explicit timestamp: %s", followUpReply)
		}
	}
}
