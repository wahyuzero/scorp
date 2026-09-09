package wizard

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"scorp-agent/config"
	"scorp-agent/models"
)

// ──────────────────────────────────────────────
// CLI Interactive Setup Wizard (Full English & Multi-Device)
// ──────────────────────────────────────────────

// RunQuickstart runs the interactive terminal setup wizard.
func RunQuickstart() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("\033[1;36m╔══════════════════════════════════════════════════════════════════════╗\033[0m")
	fmt.Println("\033[1;36m║                  🦂 SCORP AGENT — SETUP WIZARD                       ║\033[0m")
	fmt.Println("\033[1;36m║          Autonomous, Resilient & Lightweight AI Assistant            ║\033[0m")
	fmt.Println("\033[1;36m╚══════════════════════════════════════════════════════════════════════╝\033[0m")
	fmt.Println()
	fmt.Println("Welcome! This wizard will configure Scorp for your device in under 2 minutes.")
	fmt.Println()

	// ── Step 1: Operating Mode ──
	fmt.Println("\033[1m┌─ [Step 1/4] Select Operating Mode:\033[0m")
	fmt.Println("│  1) \033[1mCLI Terminal Mode\033[0m   (Interactive terminal REPL & one-shot CLI commands)")
	fmt.Println("│  2) \033[1mTelegram Bot Daemon\033[0m (24/7 background agent accessible via Telegram)")
	fmt.Println("│  3) \033[1mDual Mode\033[0m           (CLI terminal access + Telegram bot daemon)")
	fmt.Print("└─ Select [1-3] (default 1): ")

	scanner.Scan()
	modeChoice := strings.TrimSpace(scanner.Text())
	if modeChoice == "" {
		modeChoice = "1"
	}

	var tgToken, tgAllowedUsers string
	needsTelegram := modeChoice == "2" || modeChoice == "3"
	if needsTelegram {
		fmt.Println()
		fmt.Println("\033[1m┌─ Telegram Bot Configuration:\033[0m")
		existingToken := os.Getenv("TELEGRAM_BOT_TOKEN")
		if existingToken != "" {
			masked := maskString(existingToken)
			fmt.Printf("│  Found existing TELEGRAM_BOT_TOKEN (%s)\n", masked)
			fmt.Print("│  Press Enter to keep, or paste new token: ")
			scanner.Scan()
			input := strings.TrimSpace(scanner.Text())
			if input != "" {
				tgToken = input
			} else {
				tgToken = existingToken
			}
		} else {
			fmt.Println("│  Get a bot token from https://t.me/BotFather (e.g. 123456:ABC-DEF...)")
			fmt.Print("│  Telegram Bot Token: ")
			scanner.Scan()
			tgToken = strings.TrimSpace(scanner.Text())
		}

		existingUsers := os.Getenv("TELEGRAM_ALLOWED_USERS")
		fmt.Println("│")
		fmt.Println("│  Telegram Allowed Users (Security: comma-separated User IDs permitted to issue commands)")
		if existingUsers != "" {
			fmt.Printf("│  Found existing: %s\n", existingUsers)
			fmt.Print("└─ Press Enter to keep, or enter numeric IDs: ")
			scanner.Scan()
			input := strings.TrimSpace(scanner.Text())
			if input != "" {
				tgAllowedUsers = input
			} else {
				tgAllowedUsers = existingUsers
			}
		} else {
			fmt.Print("└─ Allowed User IDs (leave blank to allow all who message the bot): ")
			scanner.Scan()
			tgAllowedUsers = strings.TrimSpace(scanner.Text())
		}
	}

	// ── Step 2: Primary AI Provider Selection ──
	fmt.Println()
	fmt.Println("\033[1m┌─ [Step 2/4] Choose your Primary AI Provider:\033[0m")
	fmt.Println("│")
	fmt.Println("│  \033[1;32m── Verified & Tested in Production ──────────────────────────────\033[0m")
	fmt.Println("│  1) \033[1mGoogle Gemini\033[0m     \033[1;32m[TESTED]\033[0m   (Gemini 3.7 Flash, 3.8 Flash, 3.1 Pro)")
	fmt.Println("│  2) \033[1mOpenCode Zen\033[0m      \033[1;32m[TESTED]\033[0m   (mimo-v2.5-free, big-pickle) \033[32m[Free Tier]\033[0m")
	fmt.Println("│  3) \033[1mCommand Code\033[0m      \033[1;32m[TESTED]\033[0m   (DeepSeek Flash, Laguna, GLM) \033[32m[Free Tier]\033[0m")
	fmt.Println("│")
	fmt.Println("│  \033[1;36m── Local & Offline (No Cloud API Key) ───────────────────────────\033[0m")
	fmt.Println("│  4) \033[1mOllama Local\033[0m      \033[1;36m[OFFLINE]\033[0m  (qwen2.5-coder, llama3, deepseek-r1)")
	fmt.Println("│  5) \033[1mClaude CLI Bridge\033[0m \033[1;36m[LOCAL]\033[0m    (Uses local Claude Code auth & tools)")
	fmt.Println("│")
	fmt.Println("│  \033[33m── Cloud & Compatible Backends ──────────────────────────────────\033[0m")
	fmt.Println("│  6) \033[1mOpenAI\033[0m            \033[33m[UNTESTED]\033[0m (gpt-4o, gpt-4o-mini, o3-mini)")
	fmt.Println("│  7) \033[1mAnthropic Claude\033[0m  \033[33m[UNTESTED]\033[0m (claude-3-7-sonnet, claude-3-5-sonnet)")
	fmt.Println("│  8) \033[1mDeepSeek Official\033[0m \033[33m[UNTESTED]\033[0m (deepseek-chat, deepseek-reasoner)")
	fmt.Println("│  9) \033[1mMistral AI\033[0m        \033[33m[UNTESTED]\033[0m (codestral, mistral-large)")
	fmt.Println("│  10) \033[1mxAI Grok\033[0m         \033[33m[UNTESTED]\033[0m (grok-2, grok-beta)")
	fmt.Println("│  11) \033[1mPerplexity Sonar\033[0m \033[33m[UNTESTED]\033[0m (sonar-pro, sonar)")
	fmt.Println("│  12) \033[1mCustom Endpoint\033[0m  \033[35m[CUSTOM]\033[0m   (OpenRouter, Together, or any OpenAI-compatible URL)")
	fmt.Print("└─ Select [1-12] (default 1): ")

	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())
	if choice == "" {
		choice = "1"
	}

	var envKeyName, defaultModelKey, defaultModelID, providerName, apiFormat, baseURL string
	var maxTokens int

	switch choice {
	case "2":
		providerName = "opencode"
		envKeyName = "OPENCODE_API_KEY"
		defaultModelKey = "opencode/big-pickle"
		defaultModelID = "big-pickle"
		apiFormat = "opencode"
		baseURL = "https://opencode.ai/zen/v1"
		maxTokens = 32000
	case "3":
		providerName = "command-code"
		envKeyName = "COMMAND_CODE_API_KEY"
		defaultModelKey = "deepseek/deepseek-v4-flash"
		defaultModelID = "deepseek/deepseek-v4-flash"
		apiFormat = "command-code"
		baseURL = "https://api.commandcode.ai"
		maxTokens = 16384
	case "4":
		providerName = "ollama"
		envKeyName = ""
		defaultModelKey = "ollama/qwen2.5-coder"
		defaultModelID = "qwen2.5-coder"
		apiFormat = "openai"
		baseURL = "http://localhost:11434/v1"
		maxTokens = 16384
	case "5":
		providerName = "claude-cli"
		envKeyName = ""
		defaultModelKey = "claude-cli/claude-3-7-sonnet"
		defaultModelID = "claude-3-7-sonnet"
		apiFormat = "claude-cli"
		maxTokens = 32000
	case "6":
		providerName = "openai"
		envKeyName = "OPENAI_API_KEY"
		defaultModelKey = "openai/gpt-4o-mini"
		defaultModelID = "gpt-4o-mini"
		apiFormat = "openai"
		baseURL = "https://api.openai.com/v1"
		maxTokens = 16384
	case "7":
		providerName = "anthropic"
		envKeyName = "ANTHROPIC_API_KEY"
		defaultModelKey = "anthropic/claude-3-7-sonnet"
		defaultModelID = "claude-3-7-sonnet-20250219"
		apiFormat = "anthropic"
		baseURL = "https://api.anthropic.com/v1"
		maxTokens = 32000
	case "8":
		providerName = "deepseek"
		envKeyName = "DEEPSEEK_API_KEY"
		defaultModelKey = "deepseek/deepseek-chat"
		defaultModelID = "deepseek-chat"
		apiFormat = "openai"
		baseURL = "https://api.deepseek.com/v1"
		maxTokens = 8192
	case "9":
		providerName = "mistral"
		envKeyName = "MISTRAL_API_KEY"
		defaultModelKey = "mistral/codestral"
		defaultModelID = "codestral-latest"
		apiFormat = "mistral"
		baseURL = "https://api.mistral.ai/v1"
		maxTokens = 16384
	case "10":
		providerName = "xai"
		envKeyName = "XAI_API_KEY"
		defaultModelKey = "xai/grok-2"
		defaultModelID = "grok-2-latest"
		apiFormat = "xai"
		baseURL = "https://api.x.ai/v1"
		maxTokens = 32768
	case "11":
		providerName = "perplexity"
		envKeyName = "PERPLEXITY_API_KEY"
		defaultModelKey = "perplexity/sonar-pro"
		defaultModelID = "sonar-pro"
		apiFormat = "perplexity"
		baseURL = "https://api.perplexity.ai"
		maxTokens = 16384
	case "12":
		providerName = "custom"
		envKeyName = "CUSTOM_API_KEY"
		fmt.Print("\n  Custom Base URL (e.g. https://openrouter.ai/api/v1): ")
		scanner.Scan()
		baseURL = strings.TrimSpace(scanner.Text())
		fmt.Print("  Model Identifier (e.g. meta-llama/llama-3.3-70b-instruct): ")
		scanner.Scan()
		defaultModelID = strings.TrimSpace(scanner.Text())
		defaultModelKey = defaultModelID
		apiFormat = "openai"
		maxTokens = 16384
	default: // 1) Gemini
		providerName = "gemini"
		envKeyName = "GEMINI_API_KEY"
		defaultModelKey = "gemini/gemini-3.7-flash"
		defaultModelID = "gemini-3.7-flash"
		apiFormat = "gemini"
		baseURL = "https://generativelanguage.googleapis.com/v1beta"
		maxTokens = 65536
	}

	// ── Step 3: API Key Configuration ──
	var apiKey string
	if envKeyName != "" {
		existingKey := os.Getenv(envKeyName)
		fmt.Println()
		fmt.Println("\033[1m┌─ [Step 3/4] API Key Setup:\033[0m")
		if existingKey != "" {
			masked := maskString(existingKey)
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
			fmt.Printf("│  Enter your %s for %s\n", envKeyName, providerName)
			fmt.Print("└─ API Key: ")
			scanner.Scan()
			apiKey = strings.TrimSpace(scanner.Text())
		}
	} else {
		fmt.Println()
		fmt.Println("\033[1m┌─ [Step 3/4] Provider Setup:\033[0m")
		if providerName == "ollama" {
			fmt.Println("│  Local Ollama backend selected (no API key required).")
			fmt.Println("└─ Ensure Ollama service is running (`ollama serve`).")
		} else if providerName == "claude-cli" {
			fmt.Println("│  Local Claude Code CLI bridge selected.")
			fmt.Println("└─ Ensure Claude CLI is installed and authenticated (`claude`).")
		}
	}

	// ── Step 4: Autonomy & Safety Level ──
	fmt.Println()
	fmt.Println("\033[1m┌─ [Step 4/5] Autonomy & Safety Level:\033[0m")
	fmt.Println("│  1) \033[1;32mSupervised\033[0m [Default] (Read-only tools auto-run, destructive commands ask confirmation)")
	fmt.Println("│  2) \033[1;34mAuto\033[0m       (Smart automated risk classification with allowlist rules)")
	fmt.Println("│  3) \033[1;36mReadOnly\033[0m   (Audit mode — file writes and command executions blocked)")
	fmt.Println("│  4) \033[1;31mYOLO\033[0m       (Full autonomy — all commands run without human confirmation)")
	fmt.Print("└─ Select [1-4] (default 1): ")

	scanner.Scan()
	autoChoice := strings.TrimSpace(scanner.Text())
	autonomyLevel := "supervised"
	switch autoChoice {
	case "2":
		autonomyLevel = "auto"
	case "3":
		autonomyLevel = "readonly"
	case "4":
		autonomyLevel = "yolo"
	default:
		autonomyLevel = "supervised"
	}

	// ── Step 5: Shell Sandbox Isolation ──
	fmt.Println()
	fmt.Println("\033[1m┌─ [Step 5/5] Shell Execution Sandbox (Bubblewrap Isolation):\033[0m")
	fmt.Println("│  1) \033[1;32mEnabled (Safe Sandbox)\033[0m [Default] (Protects host root filesystem; system trees read-only)")
	fmt.Println("│  2) \033[1;33mDisabled (Host Ops)\033[0m     (Direct host execution — needed for kernel netns, ASan, raw sockets)")
	fmt.Print("└─ Select [1-2] (default 1): ")

	scanner.Scan()
	sandboxChoice := strings.TrimSpace(scanner.Text())
	sandboxMode := "on"
	if sandboxChoice == "2" {
		sandboxMode = "off"
	}

	// ── Systemd Service (Linux only) ──
	var setupService bool
	if runtime.GOOS == "linux" && needsTelegram {
		if _, err := exec.LookPath("systemctl"); err == nil {
			fmt.Println()
			fmt.Print("Would you like to install and enable Scorp as a systemd background service? [y/N]: ")
			scanner.Scan()
			svcResp := strings.ToLower(strings.TrimSpace(scanner.Text()))
			if svcResp == "y" || svcResp == "yes" {
				setupService = true
			}
		}
	}

	// Save all configurations
	saveConfig(envKeyName, apiKey, tgToken, tgAllowedUsers, autonomyLevel, sandboxMode, providerName, defaultModelKey, defaultModelID, apiFormat, baseURL, maxTokens)

	if setupService {
		installSystemdService()
	}

	fmt.Println()
	fmt.Println("\033[1;32m══════════════════════════════════════════════════════════════════════\033[0m")
	fmt.Println("\033[1;32m  ✅ Setup Complete! Scorp is configured and ready.\033[0m")
	fmt.Println("\033[1;32m══════════════════════════════════════════════════════════════════════\033[0m")
	fmt.Println()
	fmt.Println("Quick Start Commands:")
	if modeChoice == "1" || modeChoice == "3" {
		fmt.Println("  • Interactive Chat REPL:  scorp --cli")
		fmt.Println("  • One-shot Execution:     scorp \"check disk usage and memory\"")
		fmt.Println("  • Run System SOP:         scorp sop run health_audit")
	}
	if needsTelegram {
		fmt.Println("  • Run Daemon in Foregnd:  scorp daemon")
		if setupService {
			fmt.Println("  • Check Service Status:   systemctl status scorp")
		}
	}
	fmt.Println()
}

func maskString(s string) string {
	if len(s) <= 8 {
		return "••••••••"
	}
	return s[:4] + "..." + s[len(s)-4:]
}

func saveConfig(envKeyName, apiKey, tgToken, tgUsers, autonomy, sandboxMode, provider, modelKey, modelID, apiFormat, baseURL string, maxTokens int) {
	// 1. Save to .env (both current directory and ~/.scorp/.env)
	envTargets := []string{".env", filepath.Join(config.ScorpDir(), ".env")}
	for _, envPath := range envTargets {
		_ = os.MkdirAll(filepath.Dir(envPath), 0755)
		lines := []string{}
		if data, err := os.ReadFile(envPath); err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				trimmed := strings.TrimSpace(line)
				if trimmed == "" || strings.HasPrefix(trimmed, "#") {
					lines = append(lines, line)
					continue
				}
				parts := strings.SplitN(trimmed, "=", 2)
				if len(parts) == 2 {
					k := parts[0]
					if k == envKeyName || k == "SCORP_AUTONOMY" || k == "SCORP_SANDBOX" ||
						(tgToken != "" && k == "TELEGRAM_BOT_TOKEN") ||
						(tgUsers != "" && k == "TELEGRAM_ALLOWED_USERS") {
						continue // skip to overwrite below
					}
				}
				lines = append(lines, line)
			}
		}

		if envKeyName != "" && apiKey != "" {
			lines = append(lines, fmt.Sprintf("%s=%s", envKeyName, apiKey))
			_ = os.Setenv(envKeyName, apiKey)
		}
		if tgToken != "" {
			lines = append(lines, fmt.Sprintf("TELEGRAM_BOT_TOKEN=%s", tgToken))
			_ = os.Setenv("TELEGRAM_BOT_TOKEN", tgToken)
		}
		if tgUsers != "" {
			lines = append(lines, fmt.Sprintf("TELEGRAM_ALLOWED_USERS=%s", tgUsers))
			_ = os.Setenv("TELEGRAM_ALLOWED_USERS", tgUsers)
		}
		lines = append(lines, fmt.Sprintf("SCORP_AUTONOMY=%s", autonomy))
		_ = os.Setenv("SCORP_AUTONOMY", autonomy)
		lines = append(lines, fmt.Sprintf("SCORP_SANDBOX=%s", sandboxMode))
		_ = os.Setenv("SCORP_SANDBOX", sandboxMode)

		_ = os.WriteFile(envPath, []byte(strings.Join(lines, "\n")+"\n"), 0600)
	}

	config.SetAutonomyLevel(autonomy)

	// 2. Update ~/.scorp/models.json
	models.LoadModelConfig()
	if models.ModelCfg == nil {
		models.ModelCfg = &models.ModelRouterConfig{
			Models: make(map[string]models.ModelConfig),
		}
	}
	if models.ModelCfg.Models == nil {
		models.ModelCfg.Models = make(map[string]models.ModelConfig)
	}

	models.ModelCfg.DefaultModel = modelKey
	models.ModelCfg.AgentModel = modelKey

	// Add model definition if not already present
	models.ModelCfg.Models[modelKey] = models.ModelConfig{
		Provider:  provider,
		Model:     modelID,
		KeyEnv:    envKeyName,
		BaseURL:   baseURL,
		API:       apiFormat,
		MaxTokens: maxTokens,
	}

	models.SaveModelConfig()

	// 3. Ensure SOPs dir initialized
	_ = os.MkdirAll(filepath.Join(config.ScorpDir(), "sops"), 0755)
}

func installSystemdService() {
	exePath, err := os.Executable()
	if err != nil {
		exePath = "/usr/local/bin/scorp"
	}
	exePath, _ = filepath.EvalSymlinks(exePath)

	cwd, _ := os.Getwd()
	serviceContent := fmt.Sprintf(`[Unit]
Description=Scorp - Autonomous Personal AI Agent
After=network.target

[Service]
Type=simple
WorkingDirectory=%s
EnvironmentFile=-%s/.env
EnvironmentFile=-%s/.env
ExecStart=%s daemon
Restart=always
RestartSec=5
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
`, cwd, cwd, config.ScorpDir(), exePath)

	tmpService := "/tmp/scorp.service"
	if err := os.WriteFile(tmpService, []byte(serviceContent), 0644); err != nil {
		fmt.Printf("⚠️ Failed to write service file: %v\n", err)
		return
	}

	fmt.Println("Installing /etc/systemd/system/scorp.service (may prompt for sudo password)...")
	cmd := exec.Command("sudo", "cp", tmpService, "/etc/systemd/system/scorp.service")
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("⚠️ Could not copy service file: %v (%s)\n", err, string(out))
		return
	}

	_ = exec.Command("sudo", "systemctl", "daemon-reload").Run()
	_ = exec.Command("sudo", "systemctl", "enable", "--now", "scorp").Run()
	fmt.Println("✅ Systemd service installed and enabled.")
}
