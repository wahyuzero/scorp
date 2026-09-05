package main

import (
	"fmt"
	"strings"

	"scorp-agent/mcp"
	"scorp-agent/mcp/marketplace"
	"scorp-agent/mcp/transpiler"
)

// handleMCPCommand dispatches the /mcp command family:
//
//	/mcp                     — server health status (existing)
//	/mcp restart <server>    — restart one server (existing)
//	/mcp search [term]       — search the Scorp Marketplace
//	/mcp install <name> [1|2|3] — Tri-Option install (interactive when [opt] omitted)
//	/mcp install upstream <runtime>:<package> — register an unlisted upstream server
//	/mcp share <name>        — package a local rebuild for a marketplace PR
func handleMCPCommand(parts []string) {
	if len(parts) >= 3 && parts[1] == "restart" {
		targetServer := parts[2]
		fmt.Printf("🔄 Restarting MCP server '%s'...\n", targetServer)
		if err := mcp.RestartServer(targetServer); err != nil {
			fmt.Printf("❌ %v\n", err)
		} else {
			fmt.Printf("✓ MCP server '%s' successfully restarted and healthy!\n", targetServer)
		}
		return
	}

	if len(parts) >= 2 {
		switch parts[1] {
		case "search":
			term := strings.Join(parts[2:], " ")
			out, err := marketplace.CLISearch(term)
			if err != nil {
				fmt.Printf("❌ %v\n", err)
				return
			}
			fmt.Println(formatTerminalText(out))
			return

		case "share":
			if len(parts) < 3 {
				fmt.Println("Usage: /mcp share <name>  (after a successful local rebuild)")
				return
			}
			out, ok := transpiler.ShareToMarketplace(parts[2])
			fmt.Println(formatTerminalText(out))
			if !ok {
				return
			}
			return

		case "install":
			if len(parts) < 3 {
				fmt.Println("Usage: /mcp install <name> [1|2|3]")
				fmt.Println("       /mcp install upstream <runtime>:<package>")
				return
			}
			name := parts[2]
			opt := strings.Join(parts[3:], " ") // keep multi-token specs intact
			for {
				out, ok, needsSelection := marketplace.CLIInstall(name, opt)
				fmt.Println(formatTerminalText(out))
				if !needsSelection || !ok {
					return
				}
				selection, err := readInteractiveInput("")
				if err != nil {
					return
				}
				selection = strings.TrimSpace(selection)
				if selection == "" || selection == "q" || selection == "quit" {
					fmt.Println("Install cancelled.")
					return
				}
				opt = selection
			}
		}
	}

	fmt.Println(formatTerminalText(mcp.GetServerHealthStatus()))
}
