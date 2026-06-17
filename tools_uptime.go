package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ──────────────────────────────────────────────
// Uptime/Health Monitor
// ──────────────────────────────────────────────

type UptimeTarget struct {
	Name    string
	URL     string
	Status  int  // expected HTTP status code (0 = any 2xx)
	Timeout int  // seconds
}

type UptimeResult struct {
	Name      string
	URL       string
	Up        bool
	Status    int
	Latency   time.Duration
	Error     string
	Timestamp time.Time
}

type UptimeMonitor struct {
	mu       sync.RWMutex
	targets  []UptimeTarget
	history  map[string][]UptimeResult
	interval time.Duration
}

var uptimeMon = &UptimeMonitor{
	history:  make(map[string][]UptimeResult),
	interval: 5 * time.Minute,
}

// AddUptimeTarget adds a target to monitor
func AddUptimeTarget(name, url string, expectedStatus int, timeout int) {
	uptimeMon.mu.Lock()
	defer uptimeMon.mu.Unlock()
	for i, t := range uptimeMon.targets {
		if t.Name == name {
			uptimeMon.targets[i] = UptimeTarget{Name: name, URL: url, Status: expectedStatus, Timeout: timeout}
			return
		}
	}
	uptimeMon.targets = append(uptimeMon.targets, UptimeTarget{Name: name, URL: url, Status: expectedStatus, Timeout: timeout})
}

// RemoveUptimeTarget removes a target
func RemoveUptimeTarget(name string) {
	uptimeMon.mu.Lock()
	defer uptimeMon.mu.Unlock()
	var keep []UptimeTarget
	for _, t := range uptimeMon.targets {
		if t.Name != name {
			keep = append(keep, t)
		}
	}
	uptimeMon.targets = keep
	delete(uptimeMon.history, name)
}

// ListUptimeTargets returns all targets
func ListUptimeTargets() []UptimeTarget {
	uptimeMon.mu.RLock()
	defer uptimeMon.mu.RUnlock()
	out := make([]UptimeTarget, len(uptimeMon.targets))
	copy(out, uptimeMon.targets)
	return out
}

// runUptimeCheck performs a single check on all targets
func runUptimeCheck() {
	uptimeMon.mu.RLock()
	targets := make([]UptimeTarget, len(uptimeMon.targets))
	copy(targets, uptimeMon.targets)
	uptimeMon.mu.RUnlock()

	if len(targets) == 0 {
		return
	}

	var wg sync.WaitGroup
	results := make(chan UptimeResult, len(targets))

	for _, target := range targets {
		wg.Add(1)
		go func(t UptimeTarget) {
			defer wg.Done()
			results <- checkTarget(t)
		}(target)
	}

	wg.Wait()
	close(results)

	var downAlerts []UptimeResult
	uptimeMon.mu.Lock()
	for r := range results {
		name := r.Name
		uptimeMon.history[name] = append(uptimeMon.history[name], r)
		if len(uptimeMon.history[name]) > 20 {
			uptimeMon.history[name] = uptimeMon.history[name][len(uptimeMon.history[name])-20:]
		}
		if !r.Up {
			downAlerts = append(downAlerts, r)
		}
	}
	uptimeMon.mu.Unlock()

	for _, r := range downAlerts {
		log.Printf("[uptime] DOWN: %s (%s) - %s (%dms)", r.Name, r.URL, r.Error, r.Latency.Milliseconds())
		msg := fmt.Sprintf("🔴 <b>Uptime Alert: %s</b>\n🔗 <code>%s</code>\n📡 HTTP %d / %dms\n❌ %s\n🕐 %s",
			r.Name, r.URL, r.Status, r.Latency.Milliseconds(), r.Error, r.Timestamp.Format("15:04:05"))
		sendMessage(msg, nil)
	}
}

func checkTarget(target UptimeTarget) UptimeResult {
	timeout := target.Timeout
	if timeout <= 0 {
		timeout = 10
	}
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	start := time.Now()

	resp, err := client.Get(target.URL)
	latency := time.Since(start)

	if err != nil {
		return UptimeResult{
			Name: target.Name, URL: target.URL, Up: false,
			Status: 0, Latency: latency, Error: err.Error(),
			Timestamp: time.Now(),
		}
	}
	defer resp.Body.Close()

	up := true
	if target.Status > 0 && resp.StatusCode != target.Status {
		up = false
	}
	if target.Status == 0 && resp.StatusCode >= 300 {
		up = false
	}

	errMsg := ""
	if !up {
		errMsg = fmt.Sprintf("HTTP %d (expected %d)", resp.StatusCode, target.Status)
	}

	return UptimeResult{
		Name: target.Name, URL: target.URL, Up: up,
		Status: resp.StatusCode, Latency: latency,
		Error: errMsg, Timestamp: time.Now(),
	}
}

// uptimeLoop runs periodic health checks
func uptimeLoop(done chan struct{}) {
	time.Sleep(10 * time.Second)
	for {
		select {
		case <-done:
			return
		default:
		}
		runUptimeCheck()
		select {
		case <-done:
			return
		case <-time.After(uptimeMon.interval):
		}
	}
}

// executeUptime is the agent tool for managing uptime targets
func executeUptime(args map[string]interface{}) (string, bool) {
	action := getStringArg(args, "action", "list")
	switch action {
	case "list":
		targets := ListUptimeTargets()
		if len(targets) == 0 {
			return "📡 <b>Uptime Targets</b>\n\nNo targets configured. Use `action=add` to add one.", true
		}
		var b strings.Builder
		b.WriteString("📡 <b>Uptime Targets</b>\n\n")
		for _, t := range targets {
			b.WriteString(fmt.Sprintf("⬜ <b>%s</b>\n  🔗 <code>%s</code>\n", t.Name, t.URL))
		}
		return b.String(), true
	case "add":
		name := getStringArg(args, "name", "")
		url := getStringArg(args, "url", "")
		if name == "" || url == "" {
			return "❌ `name` and `url` are required", false
		}
		expectedStatus := getIntArg(args, "expected_status", 200)
		timeout := getIntArg(args, "timeout", 10)
		AddUptimeTarget(name, url, expectedStatus, timeout)
		return fmt.Sprintf("✅ Added uptime target: <b>%s</b>\n🔗 <code>%s</code>", name, url), true
	case "remove":
		name := getStringArg(args, "name", "")
		if name == "" {
			return "❌ `name` is required", false
		}
		RemoveUptimeTarget(name)
		return fmt.Sprintf("✅ Removed uptime target: <b>%s</b>", name), true
	case "check":
		if len(ListUptimeTargets()) == 0 {
			return "📡 No uptime targets configured. Use `action=add` first.", true
		}
		runUptimeCheck()
		return "✅ Uptime check completed", true
	default:
		return fmt.Sprintf("Unknown action: %s (use: list, add, remove, check)", action), false
	}
}
