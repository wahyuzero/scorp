package models

import (
	"context"
	"testing"
	"time"
)

func TestOpenCodeProvider_RegistrationAndKey(t *testing.T) {
	p := GetProvider("opencode")
	if p == nil || p.Format() != "opencode" {
		t.Fatalf("expected opencode provider, got %v", p)
	}

	pZen := GetProvider("opencode-zen")
	if pZen == nil || pZen.Format() != "opencode" {
		t.Fatalf("expected opencode-zen alias to resolve to opencode, got %v", pZen)
	}

	cfg := &ModelConfig{
		Provider: "opencode",
		Model:    "mimo-v2.5-free",
	}

	format := ResolveAPIFormat(cfg)
	if format != "opencode" {
		t.Errorf("expected format 'opencode', got '%s'", format)
	}

	baseURL := ResolveBaseURL(cfg)
	if baseURL != OpenCodeZenBaseURL {
		t.Errorf("expected base_url '%s', got '%s'", OpenCodeZenBaseURL, baseURL)
	}

	key := ResolveAPIKey(cfg)
	if key == "" {
		t.Logf("Note: OpenCode key not found on disk or env (skipping live network test)")
		return
	}

	t.Logf("OpenCode key successfully resolved via dynamic KeyResolver (len=%d)", len(key))

	// Live network test with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	reply, err := CallModel(ctx, cfg, []ChatMessage{
		{Role: "user", Content: "Respond with the single word 'PONG'."},
	})
	if err != nil {
		t.Logf("Live OpenCode CallModel returned error (might be rate-limited or offline): %v", err)
		return
	}

	t.Logf("Live OpenCode response: %s", reply)

	// Live streaming test
	streamCtx, streamCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer streamCancel()

	streamCh, err := CallModelStream(streamCtx, cfg, []ChatMessage{
		{Role: "user", Content: "Say 'Hello OpenCode'"},
	})
	if err != nil {
		t.Logf("CallModelStream error: %v", err)
		return
	}

	var streamedText string
	for chunk := range streamCh {
		if chunk.Error != nil {
			t.Logf("chunk error: %v", chunk.Error)
			break
		}
		streamedText += chunk.Content
	}
	t.Logf("Streamed OpenCode response: %s", streamedText)
}
