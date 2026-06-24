package testutil

import (
	"scorp-agent/agent"
	"scorp-agent/registry"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
)

// ──────────────────────────────────────────────
// Test Endpoint — localhost-only HTTP injection
// Allows testing agent commands without Telegram client
// POST {"command": "...", "chat_id": 123456} → triggers runAgentLoop
// ──────────────────────────────────────────────

// Per-chatID locks so only one agent loop runs at a time per chat.
// Prevents race conditions on shared session history.
var (
	chatLocksMu sync.Mutex
	chatLocks   = make(map[int64]*sync.Mutex)
)

func getChatLock(chatID int64) *sync.Mutex {
	chatLocksMu.Lock()
	defer chatLocksMu.Unlock()
	if _, ok := chatLocks[chatID]; !ok {
		chatLocks[chatID] = &sync.Mutex{}
	}
	return chatLocks[chatID]
}

func StartTestEndpoint() {
	port := os.Getenv("SCORP_TEST_PORT")
	if port == "" {
		port = "8765"
	}

	mux := http.NewServeMux()

	// POST /inject — inject a command into the agent loop
	mux.HandleFunc("/inject", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Command string `json:"command"`
			ChatID  int64  `json:"chat_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Bad JSON: %v", err), http.StatusBadRequest)
			return
		}

		if req.Command == "" {
			http.Error(w, "Missing 'command' field", http.StatusBadRequest)
			return
		}

		if req.ChatID == 0 {
			req.ChatID = 1 // default for CLI/testing
		}

		log.Printf("[test-endpoint] Injecting command for chat %d: %s", req.ChatID, req.Command)

		// Run agent loop in goroutine, but serialize per chatID
		// to prevent concurrent agent loops corrupting shared session history.
		go func() {
			lock := getChatLock(req.ChatID)
			lock.Lock()
			defer lock.Unlock()
			agent.RunAgentLoop(req.ChatID, req.Command, 0)
		}()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "injected",
			"chat_id": req.ChatID,
			"command": req.Command,
		})
	})

	// GET /health — simple health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "ok",
			"tools":  len(registry.GetAllTools()),
		})
	})

	// Listen on localhost only (security)
	ln, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		log.Printf("[test-endpoint] Failed to listen on port %s: %v", port, err)
		return
	}

	log.Printf("[test-endpoint] Listening on http://127.0.0.1:%s", port)
	go func() {
		if err := http.Serve(ln, mux); err != nil {
			log.Printf("[test-endpoint] Server error: %v", err)
		}
	}()
}
