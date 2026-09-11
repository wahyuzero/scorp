package main

import (
	"fmt"
	"strings"
	"testing"

	"scorp-agent/models"
)

func TestStripHTML(t *testing.T) {
	input := "<b>Hello</b> &amp; <i>World</i> &lt;test&gt; <pre><code>code block</code></pre>"
	output := stripHTML(input)

	if strings.Contains(output, "<b>") || strings.Contains(output, "</b>") {
		t.Errorf("expected no HTML tags, got: %s", output)
	}
	if !strings.Contains(output, "&") {
		t.Errorf("expected &amp; to be unescaped to &, got: %s", output)
	}
	if !strings.Contains(output, "<test>") {
		t.Errorf("expected &lt;test&gt; to be unescaped, got: %s", output)
	}
	if !strings.Contains(output, "```") {
		t.Errorf("expected code block backticks in output, got: %s", output)
	}
}

func TestFormatFinalResponse(t *testing.T) {
	input := "🤖 <b>Scorp:</b>\n\nThis is the final response."
	output := formatFinalResponse(input)

	if strings.HasPrefix(output, "🤖") || strings.HasPrefix(output, "Scorp:") {
		t.Errorf("expected Scorp header prefix to be stripped, got: %s", output)
	}
	if !strings.Contains(output, "This is the final response.") {
		t.Errorf("expected main text preserved, got: %s", output)
	}
}

func TestInteractiveModelSelectionRendering(t *testing.T) {
	name := "deepseek/deepseek-v4-flash"
	provider := "command-code"
	active := true
	keyStatus := "env:COMMAND_CODE_API_KEY"

	nameCol := name
	if len(nameCol) < 33 {
		nameCol += strings.Repeat(" ", 33-len(nameCol))
	}
	activeLabel := ""
	if active {
		activeLabel = "★ ACTIVE"
	}

	choice := SelectChoice{
		ID: name,
		Render: func(selected bool) string {
			if selected {
				return fmt.Sprintf("▶ \033[1;36m%-33s\033[0m  %-15s  %-10s  %s",
					name, provider, activeLabel, keyStatus)
			}
			return fmt.Sprintf("  %-33s  %-15s  %-10s  %s",
				name, provider, activeLabel, keyStatus)
		},
	}

	selectedStr := choice.Render(true)
	if !strings.Contains(selectedStr, "▶ \033[1;36mdeepseek/deepseek-v4-flash") {
		t.Errorf("selected string missing highlighted model name: %s", selectedStr)
	}
	if !strings.Contains(selectedStr, "★ ACTIVE") {
		t.Errorf("selected string missing active label: %s", selectedStr)
	}
	if !strings.Contains(selectedStr, "command-code") {
		t.Errorf("selected string missing provider: %s", selectedStr)
	}

	unselectedStr := choice.Render(false)
	if strings.Contains(unselectedStr, "▶") {
		t.Errorf("unselected string should not have cursor indicator: %s", unselectedStr)
	}
}

func TestInteractiveSessionSelectionRendering(t *testing.T) {
	sessionID := "chat-0911-test"
	active := true
	info := "(5 msgs, 11 Sep 21:00, hello world)"

	marker := " "
	if active {
		marker = "●"
	}

	choice := SelectChoice{
		ID: sessionID,
		Render: func(selected bool) string {
			if selected {
				return fmt.Sprintf("▶ \033[1;36m%-24s\033[0m %s %s", sessionID, marker, info)
			}
			return fmt.Sprintf("  %-24s %s %s", sessionID, marker, info)
		},
	}

	selectedStr := choice.Render(true)
	if !strings.Contains(selectedStr, "▶ \033[1;36mchat-0911-test") {
		t.Errorf("selected session string missing highlighted ID: %s", selectedStr)
	}
	if !strings.Contains(selectedStr, "●") {
		t.Errorf("selected session string missing active dot marker: %s", selectedStr)
	}
	if !strings.Contains(selectedStr, "5 msgs") {
		t.Errorf("selected session string missing msg count: %s", selectedStr)
	}
}

func TestInteractiveModelSwitchExecution(t *testing.T) {
	models.ModelCfgMu.Lock()
	oldCfg := models.ModelCfg
	models.ModelCfg = &models.ModelRouterConfig{
		DefaultModel: "model-one",
		AgentModel:   "model-one",
		Models: map[string]models.ModelConfig{
			"model-one": {
				Provider: "test",
				Model:    "model-one",
			},
			"model-two": {
				Provider: "test",
				Model:    "model-two",
			},
		},
	}
	models.ModelCfgMu.Unlock()
	defer func() {
		models.ModelCfgMu.Lock()
		models.ModelCfg = oldCfg
		models.ModelCfgMu.Unlock()
	}()

	// Switch model to model-two
	err := models.SwitchActiveModel("model-two")
	if err != nil {
		t.Fatalf("SwitchActiveModel failed: %v", err)
	}

	models.ModelCfgMu.RLock()
	current := models.ModelCfg.AgentModel
	models.ModelCfgMu.RUnlock()

	if current != "model-two" {
		t.Fatalf("expected active model 'model-two', got %q", current)
	}
}
