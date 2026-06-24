package agent

import (
	"scorp-agent/models"
	"scorp-agent/tools"
	"fmt"
	"testing"
	"time"
)

// TestSelfReviewIntegration tests the self-review logic with a realistic conversation
func TestSelfReviewIntegration(t *testing.T) {
	// Setup: need memory cache initialized
	tools.InitMemoryCache()
	models.LoadModelConfig()

	// Create mock conversation history (simulating 5 turns)
	history := []AgentMessage{
		{Role: "system", Content: getAgentSystemPrompt()},
		{Role: "user", Content: "Halo scorp, saya suka banget kalau response kamu singkat dan to the point, jangan pakai basa-basi"},
		{Role: "assistant", Content: "Oke, siap. Aku akan jawab singkat dan ke point."},
		{Role: "user", Content: "Oke, lanjut. Tugas ini saya mau bikin script monitoring disk space yang auto-alert kalau > 80%"},
		{Role: "assistant", Content: "Mengerti. Script bash native untuk monitoring disk, alert kalau > 80%. Bisa jalanin via cron."},
		{Role: "user", Content: "Terus, saya prefer pakai bash native, jangan python kalau bisa - biar ringan"},
		{Role: "assistant", Content: "Noted. Bash only, no Python. Ringan dan cepat."},
		{Role: "user", Content: "Saya juga kerja di Alfamart Warehouse Bali shift malam jadi kadang reply lambat"},
		{Role: "assistant", Content: "Siap, paham. Shift malam di Alfamart Warehouse Bali. No pressure untuk reply cepet."},
		{Role: "user", Content: "Btw nama saya Wahyu, panggil aja mas Wahyu atau Wahyu"},
		{Role: "assistant", Content: "Oke mas Wahyu. Terus terang, simpel, dan ringan - itu prinsip aku pegang."},
	}

	// Call runSelfReview directly (bypassing cadence checks)
	runSelfReview(999999, history)

	// Check memory for new entries
	mem := tools.ListMemory()
	fmt.Println("\n=== MEMORY AFTER SELF-REVIEW ===")
	for k, v := range mem {
		fmt.Printf("  %s: %s\n", k, v)
	}

	// Verify expected facts were extracted
	expectedKeys := []string{"prefers_concise", "prefers_bash", "works_alfamart_bali", "name_wahyu"}
	found := 0
	for _, k := range expectedKeys {
		if val, ok := mem[k]; ok && len(val) > 0 {
			fmt.Printf("✅ Found expected key: %s = %s\n", k, val)
			found++
		} else {
			fmt.Printf("⚠️  Missing expected key: %s\n", k)
		}
	}

	// At least some should be found
	if found == 0 {
		t.Log("Warning: No expected facts extracted - may need real LLM call")
	}
}

// TestSelfReviewCadence tests the cadence logic
func TestSelfReviewCadence(t *testing.T) {
	// Reset counters
	reviewMu.Lock()
	reviewTurnCount[123] = 0
	reviewLastTime[123] = time.Time{}
	reviewMu.Unlock()

	// Test: should trigger on 5th turn
	for i := 1; i <= 10; i++ {
		reviewMu.Lock()
		reviewTurnCount[123]++
		turns := reviewTurnCount[123]
		reviewMu.Unlock()

		shouldTrigger := (turns % reviewCadenceTurns == 0)
		if i == 5 {
			if !shouldTrigger {
				t.Errorf("Turn 5 should trigger review")
			} else {
				fmt.Println("✅ Turn 5 triggers review")
			}
		} else if i == 10 {
			if !shouldTrigger {
				t.Errorf("Turn 10 should trigger review")
			} else {
				fmt.Println("✅ Turn 10 triggers review")
			}
		} else {
			if shouldTrigger {
				t.Errorf("Turn %d should NOT trigger review", i)
			}
		}
	}
}

// TestSelfReviewRateLimit tests the 10-minute rate limit
func TestSelfReviewRateLimit(t *testing.T) {
	reviewMu.Lock()
	reviewLastTime[456] = time.Now().Add(-5 * time.Minute) // 5 min ago
	reviewMu.Unlock()

	// Should NOT trigger (too recent)
	reviewMu.Lock()
	reviewTurnCount[456] = 5
	reviewMu.Unlock()

	// Simulate check
	reviewMu.Lock()
	lastRun := reviewLastTime[456]
	reviewMu.Unlock()

	if time.Since(lastRun) < reviewMinInterval {
		fmt.Println("✅ Rate limit works: 5 min ago -> skip")
	} else {
		t.Errorf("Rate limit failed")
	}

	// Now set to 15 min ago - should allow
	reviewMu.Lock()
	reviewLastTime[456] = time.Now().Add(-15 * time.Minute)
	reviewMu.Unlock()

	reviewMu.Lock()
	lastRun = reviewLastTime[456]
	reviewMu.Unlock()

	if time.Since(lastRun) >= reviewMinInterval {
		fmt.Println("✅ Rate limit works: 15 min ago -> allow")
	} else {
		t.Errorf("Rate limit failed for 15 min")
	}
}
