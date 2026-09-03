package wizard

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"scorp-agent/config"
	"scorp-agent/models"
)

// ──────────────────────────────────────────────
// CLI Quickstart Wizard (ZeroClaw Parity)
// Super simple, interactive 60-second setup for terminal users
// ──────────────────────────────────────────────

// RunQuickstart runs the interactive terminal setup wizard
func RunQuickstart() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("\033[1;36m╔════════════════════════════════════════════════════════════╗\033[0m")
	fmt.Println("\033[1;36m║             🦂 SCORP AGENT — QUICKSTART WIZARD             ║\033[0m")
	fmt.Println("\033[1;36m║     Ultra-Lightweight & Fast AI Agent (<15MB RAM)          ║\033[0m")
	fmt.Println("\033[1;36m╚════════════════════════════════════════════════════════════╝\033[0m")
	fmt.Println()
	fmt.Println("Welcome! Let's get Scorp configured in under 60 seconds.")
	fmt.Println()

	// ── Step 1: Provider selection ──
	fmt.Println("\033[1m┌─ [Step 1/3] Choose your primary AI Provider:\033[0m")
	fmt.Println("│  1) \033[32mCommand Code (cmd)\033[0m [Free Tier, DeepSeek Flash, Laguna, GLM] \033[33m(Recommended)\033[0m")
	fmt.Println("│  2) DeepSeek Official (deepseek-chat / deepseek-reasoner)")
	fmt.Println("│  3) Google Gemini (Gemini 2.0 Flash / Pro)")
	fmt.Println("│  4) OpenAI (GPT-4o / GPT-4o-mini / o3-mini)")
	fmt.Println("│  5) Anthropic Claude (Claude 3.5 Sonnet / Haiku)")
	fmt.Println("│  6) Ollama (Local offline models, e.g. qwen2.5-coder)")
	fmt.Println("│  7) OpenRouter / Custom OpenAI-Compatible Endpoint")
	fmt.Print("└─ Select [1-7] (default 1): ")

	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())
	if choice == "" {
		choice = "1"
	}

	var envKeyName, defaultModel, providerName, apiFormat string
	switch choice {
	case "2":
		providerName = "deepseek"
		envKeyName = "DEEPSEEK_API_KEY"
		defaultModel = "deepseek-chat"
		apiFormat = "openai"
	case "3":
		providerName = "gemini"
		envKeyName = "GEMINI_API_KEY"
		defaultModel = "gemini-2.0-flash"
		apiFormat = "gemini"
	case "4":
		providerName = "openai"
		envKeyName = "OPENAI_API_KEY"
		defaultModel = "gpt-4o-mini"
		apiFormat = "openai"
	case "5":
		providerName = "anthropic"
		envKeyName = "ANTHROPIC_API_KEY"
		defaultModel = "claude-3-5-sonnet-latest"
		apiFormat = "anthropic"
	case "6":
		providerName = "ollama"
		envKeyName = ""
		defaultModel = "qwen2.5-coder"
		apiFormat = "openai"
	case "7":
		providerName = "openrouter"
		envKeyName = "OPENROUTER_API_KEY"
		defaultModel = "deepseek/deepseek-chat"
		apiFormat = "openai"
	default:
		providerName = "command-code"
		envKeyName = "COMMAND_CODE_API_KEY"
		defaultModel = "deepseek/deepseek-v4-flash"
		apiFormat = "command-code"
	}

	// ── Step 2: API Key input ──
	var apiKey string
	if envKeyName != "" {
		existingKey := os.Getenv(envKeyName)
		if existingKey != "" {
			masked := existingKey
			if len(masked) > 8 {
				masked = masked[:4] + "..." + masked[len(masked)-4:]
			}
			fmt.Printf("\n┌─ [Step 2/3] API Key Configuration:\n")
			fmt.Printf("│  Found existing %s (%s)\n", envKeyName, masked)
			fmt.Print("└─ Press Enter to keep, or paste new key: ")
			scanner.Scan()
			input := strings.TrimSpace(scanner.Text())
			if input != "" {
				apiKey = input
			} else {
				apiKey = existingKey
			}
		} else {
			fmt.Printf("\n┌─ [Step 2/3] Enter your %s:\n", envKeyName)
			fmt.Print("└─ Key: ")
			scanner.Scan()
			apiKey = strings.TrimSpace(scanner.Text())
		}
	} else {
		fmt.Println("\n┌─ [Step 2/3] Local Ollama Provider (no API key required)")
		fmt.Println("└─ Ensure Ollama is running (`ollama serve`).")
	}

	// ── Step 3: Autonomy Mode ──
	fmt.Println("\n┌─ [Step 3/3] Autonomy & Safety Level:")
	fmt.Println("│  1) Supervised (Default — safe tools run auto, dangerous ask confirmation)")
	fmt.Println("│  2) ReadOnly   (Audit mode — only read, search, and inspect tools allowed)")
	fmt.Println("│  3) YOLO       (Full autonomy — runs all commands without confirmation)")
	fmt.Print("└─ Select [1-3] (default 1): ")

	scanner.Scan()
	autoChoice := strings.TrimSpace(scanner.Text())
	autonomyLevel := "supervised"
	switch autoChoice {
	case "2":
		autonomyLevel = "readonly"
	case "3":
		autonomyLevel = "yolo"
	default:
		autonomyLevel = "supervised"
	}

	// Save configuration
	saveQuickstartConfig(envKeyName, apiKey, providerName, defaultModel, apiFormat, autonomyLevel)

	fmt.Println()
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println("  ✅ Setup Complete! Scorp is ready to run.")
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("Try these commands to get started:")
	fmt.Println("  • Interactive Chat REPL:  ./scorp --cli")
	fmt.Println("  • One-shot Execution:     ./scorp \"inspect disk and ram\"")
	fmt.Println("  • Run Built-in SOP:       ./scorp sop run health_audit")
	fmt.Println("  • Launch Web Gateway:     ./scorp gateway --port 8080")
	fmt.Println()
}

func saveQuickstartConfig(envKeyName, apiKey, providerName, defaultModel, apiFormat, autonomyLevel string) {
	// 1. Update .env if key provided
	if envKeyName != "" && apiKey != "" {
		envPath := ".env"
		lines := []string{}
		if data, err := os.ReadFile(envPath); err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				if strings.HasPrefix(strings.TrimSpace(line), envKeyName+"=") {
					continue
				}
				lines = append(lines, line)
			}
		}
		lines = append(lines, fmt.Sprintf("%s=%s", envKeyName, apiKey))
		lines = append(lines, fmt.Sprintf("SCORP_AUTONOMY=%s", autonomyLevel))
		_ = os.WriteFile(envPath, []byte(strings.Join(lines, "\n")), 0600)
		_ = os.Setenv(envKeyName, apiKey)
	}

	_ = os.Setenv("SCORP_AUTONOMY", autonomyLevel)
	config.SetAutonomyLevel(autonomyLevel)

	// 2. Ensure models config is set
	models.LoadModelConfig()
	if models.ModelCfg != nil {
		models.ModelCfg.DefaultModel = defaultModel
		models.ModelCfg.AgentModel = defaultModel
		models.SaveModelConfig()
	}

	// 3. Ensure SOPs dir initialized
	_ = os.MkdirAll(filepath.Join(config.ScorpDir(), "sops"), 0755)
}
