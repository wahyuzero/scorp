// +build !nobrowser

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// ──────────────────────────────────────────────
// Browser Automation Scripting Engine
// ──────────────────────────────────────────────

// ScriptStep defines a single action in a browser automation script.
type ScriptStep struct {
	Action   string `json:"action"`             // goto, click, type, extract, wait, screenshot, evaluate, scroll, submit
	URL      string `json:"url,omitempty"`      // URL for goto/screenshot
	Selector string `json:"selector,omitempty"` // CSS selector for click/type/extract/wait
	Value    string `json:"value,omitempty"`    // Text for type or JS code for evaluate
	WaitMs   int    `json:"wait_ms,omitempty"`  // Wait after action (ms)
	Timeout  int    `json:"timeout,omitempty"`  // Wait for element timeout (seconds)
	Name     string `json:"name,omitempty"`     // Variable name to store extracted data
	Retry    int    `json:"retry,omitempty"`    // Retry count on failure
}

// ScriptResult captures the result of a script execution.
type ScriptResult struct {
	Success    bool              `json:"success"`
	Error      string            `json:"error,omitempty"`
	StepCount  int               `json:"step_count"`
	FailedStep int               `json:"failed_step,omitempty"`
	Duration   string            `json:"duration"`
	Captures   map[string]string `json:"captures,omitempty"`
}

// executeScript runs a browser automation script from a file or inline JSON.
func executeScript(args map[string]interface{}, chatID int64) (string, bool) {
	source := getStringArg(args, "source", "")
	inline := getStringArg(args, "inline", "")

	var script string
	if inline != "" {
		script = inline
	} else if source != "" {
		data, err := os.ReadFile(source)
		if err != nil {
			return fmt.Sprintf("Error reading script file: %v", err), false
		}
		script = string(data)
	} else {
		return "Error: provide either 'source' (file path) or 'inline' (JSON string)", false
	}

	// Parse the script
	var steps []ScriptStep
	if err := json.Unmarshal([]byte(script), &steps); err != nil {
		return fmt.Sprintf("Error parsing script JSON: %v", err), false
	}

	if len(steps) == 0 {
		return "Error: script must contain at least one step", false
	}

	start := time.Now()
	result := &ScriptResult{
		Captures: make(map[string]string),
	}

	// Get or create browser session
	sess := getOrCreateBrowserSession(chatID)

	// Execute each step
	for i, step := range steps {
		log.Printf("[script] Step %d/%d: %s", i+1, len(steps), step.Action)

		// Retry loop
		maxRetries := step.Retry
		if maxRetries <= 0 {
			maxRetries = 1
		}

		var lastErr error
	stepLoop:
		for retry := 0; retry < maxRetries; retry++ {
			lastErr = executeStep(sess.Ctx, &step, result)
			if lastErr == nil {
				break stepLoop
			}
			if retry < maxRetries-1 {
				log.Printf("[script] Step %d failed (retry %d/%d): %v", i+1, retry+1, maxRetries, lastErr)
				time.Sleep(1 * time.Second)
			}
		}

		if lastErr != nil {
			result.Success = false
			result.Error = lastErr.Error()
			result.FailedStep = i + 1
			result.StepCount = i + 1
			result.Duration = time.Since(start).Round(time.Millisecond).String()
			return formatScriptResult(result), false
		}

		// Post-step wait
		if step.WaitMs > 0 {
			time.Sleep(time.Duration(step.WaitMs) * time.Millisecond)
		}
	}

	result.Success = true
	result.StepCount = len(steps)
	result.Duration = time.Since(start).Round(time.Millisecond).String()

	return formatScriptResult(result), true
}

func executeStep(ctx context.Context, step *ScriptStep, result *ScriptResult) error {
	timeoutCtx := ctx
	if step.Timeout > 0 {
		var cancel context.CancelFunc
		timeoutCtx, cancel = context.WithTimeout(ctx, time.Duration(step.Timeout)*time.Second)
		defer cancel()
	}

	switch step.Action {
	case "goto", "navigate":
		if step.URL == "" {
			return fmt.Errorf("url required for goto")
		}
		return chromedp.Run(timeoutCtx,
			chromedp.Navigate(step.URL),
			chromedp.WaitReady("body"),
		)

	case "click":
		if step.Selector == "" {
			return fmt.Errorf("selector required for click")
		}
		return chromedp.Run(timeoutCtx,
			chromedp.Click(step.Selector, chromedp.ByQuery),
			chromedp.Sleep(500*time.Millisecond),
		)

	case "type":
		if step.Selector == "" || step.Value == "" {
			return fmt.Errorf("selector and value required for type")
		}
		return chromedp.Run(timeoutCtx,
			chromedp.WaitReady(step.Selector),
			chromedp.Clear(step.Selector, chromedp.ByQuery),
			chromedp.SendKeys(step.Selector, step.Value, chromedp.ByQuery),
		)

	case "submit":
		if step.Selector == "" {
			return fmt.Errorf("selector required for submit")
		}
		return chromedp.Run(timeoutCtx,
			chromedp.Submit(step.Selector, chromedp.ByQuery),
			chromedp.Sleep(1*time.Second),
		)

	case "wait":
		if step.Selector != "" {
			return chromedp.Run(timeoutCtx,
				chromedp.WaitReady(step.Selector),
			)
		}
		time.Sleep(time.Duration(step.Timeout) * time.Second)
		return nil

	case "extract":
		if step.Selector == "" {
			return fmt.Errorf("selector required for extract")
		}
		var text string
		err := chromedp.Run(timeoutCtx,
			chromedp.Text(step.Selector, &text, chromedp.ByQuery),
		)
		if err != nil {
			return err
		}
		text = strings.TrimSpace(text)
		if step.Name != "" {
			result.Captures[step.Name] = text
		}
		return nil

	case "screenshot":
		var buf []byte
		err := chromedp.Run(timeoutCtx,
			chromedp.Sleep(500*time.Millisecond),
			chromedp.FullScreenshot(&buf, 90),
		)
		if err != nil {
			return err
		}
		ssDir := screenshotsDir()
		os.MkdirAll(ssDir, 0755)
		filename := fmt.Sprintf("%s/script_%d.png", ssDir, time.Now().Unix())
		if err := os.WriteFile(filename, buf, 0644); err != nil {
			return err
		}
		if step.Name != "" {
			result.Captures[step.Name] = filename
		}
		return nil

	case "evaluate":
		if step.Value == "" {
			return fmt.Errorf("value (JS code) required for evaluate")
		}
		var res interface{}
		err := chromedp.Run(timeoutCtx,
			chromedp.Evaluate(step.Value, &res),
		)
		if err != nil {
			return err
		}
		if step.Name != "" {
			result.Captures[step.Name] = fmt.Sprintf("%v", res)
		}
		return nil

	case "scroll":
		direction := step.Value
		if direction == "" {
			direction = "down"
		}
		pixels := "500"
		if direction == "up" {
			pixels = "-500"
		}
		return chromedp.Run(timeoutCtx,
			chromedp.Evaluate(fmt.Sprintf("window.scrollBy(0, %s)", pixels), nil),
		)

	default:
		return fmt.Errorf("unknown action: %s (use: goto, click, type, submit, wait, extract, screenshot, evaluate, scroll)", step.Action)
	}
}

func formatScriptResult(r *ScriptResult) string {
	var sb strings.Builder
	if r.Success {
		sb.WriteString(fmt.Sprintf("✅ Script completed: %d steps in %s", r.StepCount, r.Duration))
	} else {
		sb.WriteString(fmt.Sprintf("❌ Script failed at step %d/%d: %s", r.FailedStep, r.StepCount, r.Error))
	}
	if len(r.Captures) > 0 {
		sb.WriteString("\n\n📝 Captured data:\n")
		for name, val := range r.Captures {
			if len(val) > 100 {
				val = val[:100] + "..."
			}
			sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", name, val))
		}
	}
	return sb.String()
}

// ──────────────────────────────────────────────
// Script File Management Tools
// ──────────────────────────────────────────────

// executeScriptList lists saved automation scripts.
func executeScriptList(args map[string]interface{}) (string, bool) {
	scriptsDir := scorpPath("scripts")
	files, err := os.ReadDir(scriptsDir)
	if err != nil {
		return "No scripts directory. Create one at ~/.scorp-agent/scripts/", true
	}

	var sb strings.Builder
	sb.WriteString("📋 Saved Scripts:\n\n")
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
			info, _ := f.Info()
			sb.WriteString(fmt.Sprintf("- %s (%d B, %s)\n", f.Name(), info.Size(), info.ModTime().Format("02 Jan 15:04")))
		}
	}
	return sb.String(), true
}
