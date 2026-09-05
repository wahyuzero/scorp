package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"scorp-agent/bootstrap"
	"scorp-agent/config"
	"scorp-agent/gateway"
	"scorp-agent/mcp"
	"scorp-agent/mcp/marketplace"
	"scorp-agent/mcp/transpiler"
	"scorp-agent/models"
	"scorp-agent/sop"
	"scorp-agent/telegram"
	"scorp-agent/updater"
	"scorp-agent/wizard"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime)

	// Marketplace Option 2 (local rebuild) delegates to the AI transpiler.
	marketplace.RebuildHook = transpiler.Rebuild

	// In CLI mode, redirect diagnostic logs to ~/.scorp/scorp.log unless SCORP_DEBUG or --debug
	if isCLIMode() && os.Getenv("SCORP_DEBUG") == "" && !hasDebugFlag() {
		logDir := config.ScorpDir()
		_ = os.MkdirAll(logDir, 0755)
		if logFile, err := os.OpenFile(config.ScorpPath("scorp.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err == nil {
			log.SetOutput(logFile)
		} else {
			log.SetOutput(io.Discard)
		}
	}

	// Load environment configuration
	_ = config.LoadConfig()

	// 1. Subcommand: update
	if len(os.Args) > 1 && os.Args[1] == "update" {
		msg, err := updater.SelfUpdate()
		if err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}
		fmt.Println(msg)
		return
	}

	// 2. Subcommand: version
	if len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("scorp %s (%s/%s)\n", updater.Version, runtime.GOOS, runtime.GOARCH)
		return
	}

	// 3. Subcommand: --mcp-server (Stdio MCP server mode)
	if len(os.Args) > 1 && os.Args[1] == "--mcp-server" {
		log.Println("[mcp] MCP-only mode starting...")
		mcp.StartMCPServerMode()
		return
	}

	// 4. Subcommand: quickstart / setup
	if len(os.Args) > 1 && (os.Args[1] == "quickstart" || os.Args[1] == "setup") {
		wizard.RunQuickstart()
		return
	}

	// 5. Subcommand: gateway [--port 8080]
	if len(os.Args) > 1 && os.Args[1] == "gateway" {
		port := 8080
		for i := 2; i < len(os.Args); i++ {
			if (os.Args[i] == "--port" || os.Args[i] == "-p") && i+1 < len(os.Args) {
				if p, err := strconv.Atoi(os.Args[i+1]); err == nil && p > 0 {
					port = p
				}
			} else if strings.HasPrefix(os.Args[i], "--port=") {
				if p, err := strconv.Atoi(strings.TrimPrefix(os.Args[i], "--port=")); err == nil && p > 0 {
					port = p
				}
			}
		}
		config.InitConfigManager()
		models.LoadModelConfig()
		models.InitModelUsage()
		sop.InitDefaultSOPs()
		bootstrap.RegisterAutonomous()
		if err := gateway.StartGateway(port); err != nil {
			log.Fatalf("Gateway error: %v", err)
		}
		return
	}

	// 6. Subcommand: sop [list|run <name>]
	if len(os.Args) > 1 && os.Args[1] == "sop" {
		sop.InitDefaultSOPs()
		if len(os.Args) > 2 && os.Args[2] == "run" && len(os.Args) > 3 {
			sopName := os.Args[3]
			s, err := sop.GetSOP(sopName)
			if err != nil {
				fmt.Printf("❌ %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("🚀 Executing SOP '%s'...\n", s.Name)
			startCLI(s.Prompt)
			return
		}

		sops := sop.ListSOPs()
		fmt.Println("\n📋 Available Standard Operating Procedures (SOPs):")
		for _, s := range sops {
			fmt.Printf("  • %-18s — %s\n", s.Name, s.Description)
		}
		fmt.Println("To run an SOP: ./scorp sop run <name>")
		return
	}

	// Global Autonomy Mode flag parsing
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "--mode=") {
			config.SetAutonomyLevel(strings.TrimPrefix(arg, "--mode="))
		}
	}
	if envMode := os.Getenv("SCORP_AUTONOMY"); envMode != "" {
		config.SetAutonomyLevel(envMode)
	}

	// 7. CLI mode (interactive or one-shot execution)
	if isCLIMode() {
		var promptArgs []string
		for i := 1; i < len(os.Args); i++ {
			arg := os.Args[i]
			if arg == "--cli" || arg == "-c" || arg == "-p" || arg == "--debug" {
				continue
			}
			if strings.HasPrefix(arg, "--mode=") {
				continue
			}
			if (arg == "--session" || arg == "-s") && i+1 < len(os.Args) {
				i++ // skip value
				continue
			}
			if strings.HasPrefix(arg, "--session=") {
				continue
			}
			promptArgs = append(promptArgs, arg)
		}

		startCLI(promptArgs...)
		return
	}

	// 8. 24/7 Telegram Daemon mode
	if config.Cfg.TelegramBotToken == "" || config.Cfg.TelegramChatID == "" {
		log.Fatalf("Config error: TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID are required for Telegram mode. Run './scorp --cli' or './scorp setup'.")
	}

	telegram.StartDaemon()
}

// isCLIMode returns true if running in CLI mode (--cli flag, commands, or no Telegram token)
func isCLIMode() bool {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "--cli" || arg == "-c" || arg == "-p" {
			return true
		}
		if arg != "update" && arg != "version" && arg != "--version" && arg != "-v" && arg != "--mcp-server" && !strings.HasPrefix(arg, "-") {
			return true
		}
	}
	// If no token configured, default to CLI mode
	if config.Cfg.TelegramBotToken == "" && os.Getenv("TELEGRAM_BOT_TOKEN") == "" {
		return true
	}
	return false
}
