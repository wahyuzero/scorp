package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// ──────────────────────────────────────────────
// Telegram Inline Query (Phase 4.5)
// ──────────────────────────────────────────────

// TGInlineQuery represents an incoming inline query
type TGInlineQuery struct {
	ID     string
	Query  string
	UserID int64
	Offset string
}

// answerInlineQuery sends results for an inline query
func answerInlineQuery(queryID string, results []map[string]interface{}, cacheTime int) bool {
	payload := map[string]interface{}{
		"inline_query_id": queryID,
		"results":         results,
		"cache_time":      cacheTime,
		"is_personal":     true,
	}
	resp, err := tgPost("/answerInlineQuery", payload)
	if err != nil {
		log.Printf("[inline] answer error: %v", err)
		return false
	}
	return resp.OK
}

// handleInlineQuery processes an inline query and returns results
func handleInlineQuery(query TGInlineQuery) {
	q := strings.TrimSpace(query.Query)
	results := buildInlineResults(q)
	if len(results) == 0 {
		results = buildInlineResults("help")
	}
	answerInlineQuery(query.ID, results, 30)
}

// buildInlineResults creates inline query results based on the query text
func buildInlineResults(query string) []map[string]interface{} {
	lower := strings.ToLower(query)
	var results []map[string]interface{}

	switch {
	case strings.HasPrefix(lower, "status") || strings.HasPrefix(lower, "st"):
		text := quickStatus()
		results = append(results, map[string]interface{}{
			"type":        "article",
			"id":          "status",
			"title":       "⚡ System Status",
			"description": text,
			"input_message_content": map[string]interface{}{
				"message_text":             text,
				"parse_mode":               "HTML",
				"disable_web_page_preview": true,
			},
		})

	case strings.HasPrefix(lower, "docker") || strings.HasPrefix(lower, "dc"):
		text := quickDocker()
		results = append(results, map[string]interface{}{
			"type":        "article",
			"id":          "docker",
			"title":       "🐳 Docker Status",
			"description": text,
			"input_message_content": map[string]interface{}{
				"message_text": text,
				"parse_mode":   "HTML",
			},
		})

	case strings.HasPrefix(lower, "disk") || strings.HasPrefix(lower, "storage"):
		text := quickStorage()
		results = append(results, map[string]interface{}{
			"type":        "article",
			"id":          "storage",
			"title":       "📁 Storage",
			"description": text,
			"input_message_content": map[string]interface{}{
				"message_text": text,
				"parse_mode":   "HTML",
			},
		})

	case strings.HasPrefix(lower, "help") || query == "":
		suggestions := []struct {
			id, title, desc, cmd string
		}{
			{"h1", "⚡ System Status", "Quick CPU/RAM/Disk overview", "status"},
			{"h2", "🐳 Docker Status", "View all containers", "docker"},
			{"h3", "📁 Storage", "Disk usage details", "storage"},
			{"h4", "🌐 Network", "Network interfaces & speed", "network"},
		}
		for _, s := range suggestions {
			results = append(results, map[string]interface{}{
				"type":        "article",
				"id":          s.id,
				"title":       s.title,
				"description": s.desc,
				"input_message_content": map[string]interface{}{
					"message_text": fmt.Sprintf("📡 Inline: %s\nType one of: status, docker, storage, network, help", s.cmd),
				},
			})
		}

	default:
		safeCmds := map[string]string{
			"df":     "💾 Disk Free",
			"free":   "🧠 Memory",
			"uptime": "⏱ Uptime",
			"date":   "📅 Date/Time",
		}
		if title, ok := safeCmds[query]; ok {
			output, errOut := runInlineSafeCmd(query)
			text := fmt.Sprintf("<b>%s</b>\n<code>%s</code>", title, escapeHTML(output))
			if errOut != "" {
				text = fmt.Sprintf("❌ <b>Error</b>\n<code>%s</code>", escapeHTML(errOut))
			}
			results = append(results, map[string]interface{}{
				"type":        "article",
				"id":          "cmd_" + query,
				"title":       title,
				"description": firstN(output, 80),
				"input_message_content": map[string]interface{}{
					"message_text": text,
					"parse_mode":   "HTML",
				},
			})
		}
	}

	return results
}

func quickStatus() string {
	sys := collectSystem()
	return fmt.Sprintf("⚡ <b>Status</b>\nCPU: %.0f%%\nRAM: %.0f%%\nDisk: %.0f%%",
		sys.CPUPercent, sys.RAMPercent, sys.DiskPercent)
}

func quickDocker() string {
	docker := collectDocker()
	var icons []string
	for _, c := range docker.Containers {
		if c.Status == "running" {
			icons = append(icons, fmt.Sprintf("🟢 %s", c.Name))
		} else {
			icons = append(icons, fmt.Sprintf("🔴 %s", c.Name))
		}
	}
	return fmt.Sprintf("🐳 <b>Docker</b>\n%s", strings.Join(icons, "\n"))
}

func quickStorage() string {
	// Shell-based quick check
	out, err := exec.Command("df", "-h", "/").CombinedOutput()
	if err != nil {
		return "📁 Storage: error"
	}
	return "📁 <b>Storage</b>\n<code>" + escapeHTML(strings.TrimSpace(string(out))) + "</code>"
}

// runInlineSafeCmd runs a safe read-only shell command
func runInlineSafeCmd(cmd string) (string, string) {
	var args []string
	switch cmd {
	case "df":
		args = []string{"df", "-h", "/"}
	case "free":
		args = []string{"free", "-h"}
	case "uptime":
		args = []string{"uptime"}
	case "date":
		args = []string{"date"}
	default:
		return "", "Command not allowed"
	}
	out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
	if err != nil {
		return "", err.Error()
	}
	return strings.TrimSpace(string(out)), ""
}

func firstN(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

// setupInlineMode enables inline queries for this bot
func setupInlineMode() {
	tgPost("/setMyDescription", map[string]interface{}{
		"description": "VPS Monitor Bot — Type @botname in any chat for quick status",
	})
	tgPost("/setMyName", map[string]interface{}{
		"name": "VPS Monitor",
	})
	tgPost("/setMyShortDescription", map[string]interface{}{
		"short_description": "Quick VPS status: @botname status|docker|disk|help",
	})
	log.Println("[inline] Inline mode configured")
}
