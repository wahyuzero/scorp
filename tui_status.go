package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"scorp-agent/agent"
	"scorp-agent/config"
	"scorp-agent/models"
)

// getShortCwd returns a shortened, clean directory path
func getShortCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "workspace"
	}
	home := os.Getenv("HOME")
	if home != "" && strings.HasPrefix(cwd, home) {
		cwd = "~" + strings.TrimPrefix(cwd, home)
	}
	// If path is too long, show last two segments
	parts := strings.Split(cwd, "/")
	if len(parts) > 3 && parts[0] != "~" {
		return ".../" + strings.Join(parts[len(parts)-2:], "/")
	}
	return cwd
}

// getGitStatus returns the active git branch name and dirty indicator
func getGitStatus() string {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return ""
	}
	branch := strings.TrimSpace(string(out))
	if branch == "" {
		return ""
	}

	// Check if repo is dirty
	statusOut, err := exec.Command("git", "status", "--porcelain").Output()
	dirty := ""
	if err == nil && len(strings.TrimSpace(string(statusOut))) > 0 {
		dirty = "*"
	}

	return fmt.Sprintf("git:(%s%s)", branch, dirty)
}

// getContextPill returns the active context window usage
func getContextPill(sessionID string) string {
	tokens := agent.GetHistoryTokenEstimate(sessionID)
	// Default context window threshold
	maxTokens := 64000
	if m := models.RouteModel("agent"); m != nil {
		if m.MaxTokens > 0 {
			maxTokens = m.MaxTokens * 4 // approx window
		}
	}
	if maxTokens <= 0 {
		maxTokens = 64000
	}

	pct := float64(tokens) / float64(maxTokens) * 100.0
	if pct > 100.0 {
		pct = 100.0
	}

	// Mini progress bar: 5 segments
	filled := int(pct / 20.0)
	if filled > 5 {
		filled = 5
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", 5-filled)

	// ANSI color based on usage
	color := "\033[32m" // green
	if pct > 50.0 {
		color = "\033[33m" // yellow
	}
	if pct > 80.0 {
		color = "\033[31m" // red
	}

	var tokStr string
	if tokens < 1000 {
		tokStr = fmt.Sprintf("%d", tokens)
	} else {
		tokStr = fmt.Sprintf("%.1fk", float64(tokens)/1000.0)
	}

	return fmt.Sprintf("🧠 %sctx: %s (%.0f%%) [%s]\033[0m", color, tokStr, pct, bar)
}

// renderStatusFooter returns the styled single-line persistent bottom bar
func renderStatusFooter(sessionID string) string {
	folder := getShortCwd()
	gitInfo := getGitStatus()
	locStr := "📂 " + folder
	if gitInfo != "" {
		locStr += " \033[35m" + gitInfo + "\033[0m"
	}

	// Active Model
	modelName := "default"
	if m := models.RouteModel("agent"); m != nil {
		modelName = m.Model
		if idx := strings.LastIndex(modelName, "/"); idx >= 0 {
			modelName = modelName[idx+1:]
		}
	}
	modelPill := "🤖 \033[1;36m" + modelName + "\033[0m"

	// Context Window
	ctxPill := getContextPill(sessionID)

	// Daily cost
	spend := models.GetDailyTotalUSD()
	costPill := fmt.Sprintf("💰 \033[32m$%.4f\033[0m", spend)

	// Autonomy Mode
	mode := config.GetAutonomyLevel()
	modeIcon := "🛡️"
	if mode == "yolo" {
		modeIcon = "⚡"
	} else if mode == "readonly" {
		modeIcon = "🔒"
	}
	modePill := fmt.Sprintf("%s \033[33m%s\033[0m", modeIcon, mode)

	// Combine into a sleek powerline capsule
	return fmt.Sprintf("\033[2m╰─\033[0m %s \033[2m─\033[0m %s \033[2m─\033[0m %s \033[2m─\033[0m %s \033[2m─\033[0m %s \033[2m─╯\033[0m",
		locStr, modelPill, ctxPill, costPill, modePill)
}
