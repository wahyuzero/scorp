package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"scorp-agent/agent"
	"scorp-agent/config"
	"scorp-agent/models"
	"scorp-agent/registry"
	"scorp-agent/sop"
	"scorp-agent/tools"
)

// ──────────────────────────────────────────────
// Micro HTTP Gateway & Embedded Web Dashboard (ZeroClaw Parity)
// Uses < 2MB RAM (pure net/http standard library, zero framework bloat)
// ──────────────────────────────────────────────

var (
	gatewayStartTime = time.Now()
	gatewayMu        sync.Mutex
)

// StartGateway launches the HTTP Gateway daemon and Web Dashboard
func StartGateway(port int) error {
	if port <= 0 {
		port = 8080
	}

	mux := http.NewServeMux()

	// ── Web Dashboard (Single-file embedded UI) ──
	mux.HandleFunc("/", handleDashboard)

	// ── REST API ──
	mux.HandleFunc("/api/status", handleStatus)
	mux.HandleFunc("/api/chat", handleChat)
	mux.HandleFunc("/api/tools", handleTools)
	mux.HandleFunc("/api/sops", handleSOPs)
	mux.HandleFunc("/api/receipts", handleReceipts)

	addr := fmt.Sprintf(":%d", port)
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Printf("║  🌐 SCORP MICRO GATEWAY & WEB DASHBOARD                     ║\n")
	fmt.Printf("║  URL: http://localhost:%d                                ║\n", port)
	fmt.Println("║  Footprint: < 2MB RAM | Single Binary                       ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	log.Printf("[gateway] Listening on http://0.0.0.0:%d", port)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	return server.ListenAndServe()
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	status := map[string]interface{}{
		"status":      "online",
		"version":     "v2.0-modern",
		"uptime":      time.Since(gatewayStartTime).Round(time.Second).String(),
		"autonomy":    config.GetAutonomyLevel(),
		"ram_alloc_mb": float64(m.Alloc) / (1024 * 1024),
		"goroutines":  runtime.NumGoroutine(),
		"models_count": len(models.ModelCfg.Models),
		"active_model": models.ModelCfg.AgentModel,
	}
	json.NewEncoder(w).Encode(status)
}

func handleTools(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	all := registry.GetAllTools()
	json.NewEncoder(w).Encode(all)
}

func handleSOPs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sops := sop.ListSOPs()
	json.NewEncoder(w).Encode(sops)
}

func handleReceipts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	receipts := tools.GetRecentReceipts()
	json.NewEncoder(w).Encode(receipts)
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Message) == "" {
		http.Error(w, "Missing 'message' field", http.StatusBadRequest)
		return
	}

	gatewayMu.Lock()
	defer gatewayMu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	// Call model with tools and fallback
	ctx, cancel := contextWithTimeout(90 * time.Second)
	defer cancel()

	messages := []models.ChatMessage{
		{Role: "system", Content: "You are Scorp Agent running via Micro Gateway."},
		{Role: "user", Content: req.Message},
	}

	reply, toolCalls, modelUsed, err := models.CallModelWithToolsAndFallback(ctx, "agent", messages)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	executedTools := []map[string]interface{}{}
	for _, tc := range toolCalls {
		out, ok := agent.ExecuteTool(tc, 0)
		executedTools = append(executedTools, map[string]interface{}{
			"tool":    tc.Name,
			"args":    tc.Args,
			"output":  out,
			"success": ok,
		})
	}

	resp := map[string]interface{}{
		"reply":      reply,
		"model_used": modelUsed,
		"tools":      executedTools,
		"timestamp":  time.Now(),
	}
	json.NewEncoder(w).Encode(resp)
}

func contextWithTimeout(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d)
}

// handleDashboard serves an ultra-clean, modern single-page dashboard
func handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(dashboardHTML))
}

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Scorp Agent — Gateway Dashboard</title>
  <style>
    :root {
      --bg: #0d1117;
      --card-bg: #161b22;
      --border: #30363d;
      --text: #c9d1d9;
      --text-dim: #8b949e;
      --accent: #58a6ff;
      --success: #238636;
      --danger: #da3633;
    }
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
      background: var(--bg);
      color: var(--text);
      line-height: 1.5;
      padding: 20px;
    }
    .container { max-width: 900px; margin: 0 auto; }
    header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding-bottom: 16px;
      border-bottom: 1px solid var(--border);
      margin-bottom: 20px;
    }
    h1 { font-size: 1.4rem; display: flex; align-items: center; gap: 8px; }
    .badge {
      background: #1f6feb22;
      color: var(--accent);
      border: 1px solid #1f6feb;
      font-size: 0.75rem;
      padding: 2px 8px;
      border-radius: 12px;
      text-transform: uppercase;
      font-weight: 600;
    }
    .stats-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
      gap: 12px;
      margin-bottom: 20px;
    }
    .stat-card {
      background: var(--card-bg);
      border: 1px solid var(--border);
      padding: 14px;
      border-radius: 8px;
    }
    .stat-val { font-size: 1.3rem; font-weight: bold; color: #f0f6fc; margin-top: 4px; }
    .stat-label { font-size: 0.8rem; color: var(--text-dim); }
    .chat-box {
      background: var(--card-bg);
      border: 1px solid var(--border);
      border-radius: 8px;
      padding: 16px;
      margin-bottom: 20px;
    }
    .chat-history {
      max-height: 350px;
      overflow-y: auto;
      margin-bottom: 14px;
      display: flex;
      flex-direction: column;
      gap: 10px;
    }
    .msg { padding: 10px 14px; border-radius: 6px; font-size: 0.95rem; }
    .msg-user { background: #1f6feb22; border-left: 3px solid var(--accent); align-self: flex-end; max-width: 85%; }
    .msg-agent { background: #21262d; border-left: 3px solid var(--success); align-self: flex-start; max-width: 90%; white-space: pre-wrap; }
    .input-row { display: flex; gap: 8px; }
    input[type="text"] {
      flex: 1;
      background: #0d1117;
      border: 1px solid var(--border);
      color: #fff;
      padding: 10px 14px;
      border-radius: 6px;
      font-size: 0.95rem;
    }
    button {
      background: var(--success);
      color: #fff;
      border: none;
      padding: 10px 18px;
      border-radius: 6px;
      font-weight: 600;
      cursor: pointer;
    }
    button:hover { opacity: 0.9; }
    button:disabled { opacity: 0.5; cursor: not-allowed; }
  </style>
</head>
<body>
  <div class="container">
    <header>
      <h1>🦂 Scorp Gateway <span class="badge" id="autonomy-badge">supervised</span></h1>
      <span style="color:var(--text-dim); font-size:0.85rem;" id="uptime-text">Uptime: --</span>
    </header>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-label">Allocated RAM</div>
        <div class="stat-val" id="ram-val">-- MB</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Active Model</div>
        <div class="stat-val" id="model-val" style="font-size:0.95rem;">--</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Registered Tools</div>
        <div class="stat-val" id="tools-val">--</div>
      </div>
    </div>

    <div class="chat-box">
      <div class="chat-history" id="history">
        <div class="msg msg-agent">👋 Hello! Scorp Agent is ready. Send a task, command, or query below.</div>
      </div>
      <div class="input-row">
        <input type="text" id="prompt-input" placeholder="Type a task (e.g. 'check server health' or 'read url https://...')">
        <button id="send-btn" onclick="sendTask()">Send</button>
      </div>
    </div>
  </div>

  <script>
    async function loadStatus() {
      try {
        const res = await fetch('/api/status');
        const d = await res.json();
        document.getElementById('uptime-text').innerText = 'Uptime: ' + d.uptime;
        document.getElementById('ram-val').innerText = d.ram_alloc_mb.toFixed(2) + ' MB';
        document.getElementById('model-val').innerText = d.active_model;
        document.getElementById('autonomy-badge').innerText = d.autonomy;

        const toolsRes = await fetch('/api/tools');
        const tools = await toolsRes.json();
        document.getElementById('tools-val').innerText = tools.length;
      } catch (e) {
        console.error(e);
      }
    }

    async function sendTask() {
      const input = document.getElementById('prompt-input');
      const text = input.value.trim();
      if (!text) return;

      const history = document.getElementById('history');
      const userMsg = document.createElement('div');
      userMsg.className = 'msg msg-user';
      userMsg.innerText = text;
      history.appendChild(userMsg);
      input.value = '';

      const btn = document.getElementById('send-btn');
      btn.disabled = true;
      btn.innerText = 'Working...';

      const agentMsg = document.createElement('div');
      agentMsg.className = 'msg msg-agent';
      agentMsg.innerText = '⏳ Processing task...';
      history.appendChild(agentMsg);
      history.scrollTop = history.scrollHeight;

      try {
        const res = await fetch('/api/chat', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify({message: text})
        });
        const data = await res.json();
        if (data.error) {
          agentMsg.innerText = '❌ Error: ' + data.error;
        } else {
          agentMsg.innerText = data.reply || (data.tools && data.tools.length ? 'Task executed via tools.' : 'Done.');
        }
      } catch (err) {
        agentMsg.innerText = '❌ Request failed: ' + err;
      } finally {
        btn.disabled = false;
        btn.innerText = 'Send';
        history.scrollTop = history.scrollHeight;
        loadStatus();
      }
    }

    document.getElementById('prompt-input').addEventListener('keypress', (e) => {
      if (e.key === 'Enter') sendTask();
    });

    loadStatus();
    setInterval(loadStatus, 10000);
  </script>
</body>
</html>
`
