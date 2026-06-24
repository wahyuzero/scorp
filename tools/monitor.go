// +build !nobrowser

package tools

import (
	"scorp-agent/internal/helpers"
	"scorp-agent/rag"
	"scorp-agent/config"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)
// ──────────────────────────────────────────────
// Browser Monitor — Scheduled scraping with change detection
// ──────────────────────────────────────────────

// MonitorTarget defines a page to watch.
type MonitorTarget struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	Selector    string `json:"selector,omitempty"`    // CSS selector to extract specific content
	Screenshot  bool   `json:"screenshot"`            // Capture screenshot on each check
	IngestRAG   bool   `json:"ingest_rag"`            // Auto-ingest extracted content to RAG
	IntervalMin int    `json:"interval_min"`          // Check interval in minutes
	LastHash    string `json:"last_hash,omitempty"`   // SHA-256 of last content
	LastCheck   string `json:"last_check,omitempty"`  // ISO timestamp
	LastChange  string `json:"last_change,omitempty"` // ISO timestamp of last detected change
	CheckCount  int    `json:"check_count"`
	ChangeCount int    `json:"change_count"`
	Created     string `json:"created"`
}

var (
	monitorTargets   []MonitorTarget
	monitorTargetsMu sync.Mutex
	monitorRunning   bool
)

// InitMonitor loads saved monitor targets.
func InitMonitor() {
	loadMonitorTargets()
	// Start background checker
	go monitorLoop()
}

func loadMonitorTargets() {
	monitorTargetsMu.Lock()
	defer monitorTargetsMu.Unlock()

	if err := config.ConfigMgr().Load("monitor_targets.json", &monitorTargets); err != nil {
		log.Printf("[monitor] Load error: %v", err)
	}
	if monitorTargets == nil {
		monitorTargets = []MonitorTarget{}
	}
	log.Printf("[monitor] Loaded %d targets", len(monitorTargets))
}

func saveMonitorTargets() {
	monitorTargetsMu.Lock()
	defer monitorTargetsMu.Unlock()
	if err := config.ConfigMgr().Save("monitor_targets.json", monitorTargets); err != nil {
		log.Printf("[monitor] Save error: %v", err)
	}
}

// monitorLoop runs in background, checking targets on their intervals.
func monitorLoop() {
	if monitorRunning {
		return
	}
	monitorRunning = true

	ticker := time.NewTicker(60 * time.Second) // Check every minute
	defer ticker.Stop()

	for range ticker.C {
		monitorTargetsMu.Lock()
		targets := make([]MonitorTarget, len(monitorTargets))
		copy(targets, monitorTargets)
		monitorTargetsMu.Unlock()

		for i := range targets {
			t := &targets[i]
			if t.IntervalMin <= 0 {
				continue
			}

			// Parse last check time
			var lastCheck time.Time
			if t.LastCheck != "" {
				lastCheck, _ = time.Parse(time.RFC3339, t.LastCheck)
			}

			if time.Since(lastCheck) < time.Duration(t.IntervalMin)*time.Minute {
				continue
			}

			// Run check
			log.Printf("[monitor] Checking: %s (%s)", t.Name, t.URL)
			changed, content, err := monitorCheckOne(t)
			if err != nil {
				log.Printf("[monitor] Error checking %s: %v", t.Name, err)
			}

			if changed {
				log.Printf("[monitor] 🔄 Change detected: %s", t.Name)
				t.ChangeCount++
				t.LastChange = time.Now().Format(time.RFC3339)

				// Alert via Telegram
				alertMsg := fmt.Sprintf("🔄 <b>Page Change Detected</b>\n\n<b>%s</b>\n%s\n\nChanged content (%d chars):\n<code>%s</code>",
					t.Name, t.URL, len(content), truncateForAlert(content))
				SendMessage(alertMsg, nil)

				// Auto-ingest to RAG if enabled
				if t.IngestRAG && content != "" {
					go func(name, url, text string) {
						source := fmt.Sprintf("monitor:%s", name)
						ragIngestText(text, source, url)
						log.Printf("[monitor] Ingested to RAG: %s (%d chars)", name, len(text))
					}(t.Name, t.URL, content)
				}
			}

			t.CheckCount++
			t.LastCheck = time.Now().Format(time.RFC3339)

			// Update in global slice
			monitorTargetsMu.Lock()
			for j := range monitorTargets {
				if monitorTargets[j].ID == t.ID {
					monitorTargets[j] = *t
					break
				}
			}
			monitorTargetsMu.Unlock()
		}

		saveMonitorTargets()
	}
}

// monitorCheckOne visits a URL, extracts content, and checks for changes.
func monitorCheckOne(t *MonitorTarget) (changed bool, content string, err error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var text string
	var buf []byte

	actions := []chromedp.Action{
		chromedp.Navigate(t.URL),
		chromedp.WaitReady("body"),
		chromedp.Sleep(2 * time.Second), // Let JS render
	}

	if t.Selector != "" {
		actions = append(actions, chromedp.Text(t.Selector, &text, chromedp.ByQuery))
	} else {
		actions = append(actions, chromedp.Text("body", &text, chromedp.ByQuery))
	}

	if t.Screenshot {
		actions = append(actions, chromedp.FullScreenshot(&buf, 90))
	}

	err = chromedp.Run(ctx, actions...)
	if err != nil {
		return false, "", err
	}

	// Normalize content
	text = strings.TrimSpace(text)
	if len(text) > 50000 {
		text = text[:50000]
	}

	// Hash for change detection
	hasher := sha256.New()
	hasher.Write([]byte(text))
	newHash := hex.EncodeToString(hasher.Sum(nil))

	changed = t.LastHash != "" && newHash != t.LastHash
	t.LastHash = newHash

	// Save screenshot
	if t.Screenshot && len(buf) > 0 {
		ssDir := config.ScreenshotsDir()
		os.MkdirAll(ssDir, 0755)
		filename := filepath.Join(ssDir, fmt.Sprintf("%s_%d.png", sanitizeFilename(t.Name), time.Now().Unix()))
		os.WriteFile(filename, buf, 0644)
	}

	return changed, text, nil
}

// ──────────────────────────────────────────────
// Monitor Tools (exposed to agent)
// ──────────────────────────────────────────────

func ExecuteMonitor(args map[string]interface{}, chatID int64) (string, bool) {
	action := helpers.GetStringArg(args, "action", "")

	switch action {
	case "add":
		url := helpers.GetStringArg(args, "url", "")
		name := helpers.GetStringArg(args, "name", "")
		if url == "" || name == "" {
			return "Error: url and name required for monitor add", false
		}
		selector := helpers.GetStringArg(args, "selector", "")
		interval := helpers.GetIntArg(args, "interval_min", 30)
		screenshot := helpers.GetBoolArg(args, "screenshot", false)
		ingestRAG := helpers.GetBoolArg(args, "ingest_rag", true)

		target := MonitorTarget{
			ID:          fmt.Sprintf("mon_%d", time.Now().UnixNano()),
			URL:         url,
			Name:        name,
			Selector:    selector,
			IntervalMin: interval,
			Screenshot:  screenshot,
			IngestRAG:   ingestRAG,
			Created:     time.Now().Format(time.RFC3339),
		}

		monitorTargetsMu.Lock()
		monitorTargets = append(monitorTargets, target)
		monitorTargetsMu.Unlock()
		saveMonitorTargets()

		return fmt.Sprintf("✅ Monitor target added:\n  Name: %s\n  URL: %s\n  Selector: %s\n  Interval: %dm\n  Screenshot: %v\n  RAG: %v\n  ID: %s",
			name, url, selector, interval, screenshot, ingestRAG, target.ID), true

	case "remove":
		id := helpers.GetStringArg(args, "id", "")
		name := helpers.GetStringArg(args, "name", "")
		if id == "" && name == "" {
			return "Error: id or name required for monitor remove", false
		}

		monitorTargetsMu.Lock()
		removed := false
		for i, t := range monitorTargets {
			if (id != "" && t.ID == id) || (name != "" && t.Name == name) {
				monitorTargets = append(monitorTargets[:i], monitorTargets[i+1:]...)
				removed = true
				break
			}
		}
		monitorTargetsMu.Unlock()
		if removed {
			saveMonitorTargets()
			return fmt.Sprintf("✅ Removed monitor target: %s", name), true
		}
		return "No matching monitor target found", false

	case "list":
		monitorTargetsMu.Lock()
		defer monitorTargetsMu.Unlock()
		if len(monitorTargets) == 0 {
			return "No monitor targets. Use monitor_add to create one.", true
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("📋 Monitor Targets (%d):\n\n", len(monitorTargets)))
		for _, t := range monitorTargets {
			sb.WriteString(fmt.Sprintf("┌─ %s\n", t.Name))
			sb.WriteString(fmt.Sprintf("│  URL: %s\n", t.URL))
			if t.Selector != "" {
				sb.WriteString(fmt.Sprintf("│  Selector: %s\n", t.Selector))
			}
			sb.WriteString(fmt.Sprintf("│  Interval: %dm | Checks: %d | Changes: %d\n", t.IntervalMin, t.CheckCount, t.ChangeCount))
			if t.LastCheck != "" {
				sb.WriteString(fmt.Sprintf("│  Last check: %s\n", t.LastCheck))
			}
			if t.LastChange != "" {
				sb.WriteString(fmt.Sprintf("│  Last change: %s\n", t.LastChange))
			}
			sb.WriteString(fmt.Sprintf("│  Screenshot: %v | RAG: %v\n", t.Screenshot, t.IngestRAG))
			sb.WriteString(fmt.Sprintf("│  ID: %s\n└──\n\n", t.ID))
		}
		return sb.String(), true

	case "check":
		// Force immediate check of a target
		name := helpers.GetStringArg(args, "name", "")
		id := helpers.GetStringArg(args, "id", "")

		monitorTargetsMu.Lock()
		var target *MonitorTarget
		for i := range monitorTargets {
			if (id != "" && monitorTargets[i].ID == id) || (name != "" && monitorTargets[i].Name == name) {
				target = &monitorTargets[i]
				break
			}
		}
		monitorTargetsMu.Unlock()

		if target == nil {
			return "No matching monitor target found", false
		}

		changed, content, err := monitorCheckOne(target)
		if err != nil {
			return fmt.Sprintf("Error checking %s: %v", target.Name, err), false
		}

		target.CheckCount++
		target.LastCheck = time.Now().Format(time.RFC3339)
		if changed {
			target.ChangeCount++
			target.LastChange = time.Now().Format(time.RFC3339)
		}

		monitorTargetsMu.Lock()
		for j := range monitorTargets {
			if monitorTargets[j].ID == target.ID {
				monitorTargets[j] = *target
				break
			}
		}
		monitorTargetsMu.Unlock()
		saveMonitorTargets()

		status := "no change"
		if changed {
			status = "🔄 CHANGED"
		}
		preview := truncateForAlert(content)
		return fmt.Sprintf("✅ Check complete: %s\nStatus: %s\nContent: %d chars\nPreview: %s",
			target.Name, status, len(content), preview), true

	default:
		return "Unknown monitor action (use: add, remove, list, check)", false
	}
}

// ──────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────

func truncateForAlert(s string) string {
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	if len(s) > 200 {
		return s[:200] + "..."
	}
	return s
}

func sanitizeFilename(s string) string {
	s = strings.ToLower(s)
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, s)
	return s
}

// ragIngestText feeds text into the RAG vector index using existing infrastructure.
func ragIngestText(text, source, url string) {
	if rag.VecIndex == nil {
		log.Printf("[monitor] RAG index not initialized, skipping ingest")
		return
	}
	// Use the same chunking logic as ragVecIngest
	chunks := rag.SmartChunk(text, 500)
	for _, chunk := range chunks {
		rag.VecIndex.AddVecChunk(source, chunk)
	}
	rag.VecIndex.Persist()
	log.Printf("[monitor] Ingested %d chunks to RAG from %s", len(chunks), source)
}
