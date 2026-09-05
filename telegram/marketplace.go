package telegram

import (
	"context"
	"fmt"

	"scorp-agent/mcp/marketplace"
)

// handleTelegramMarketplaceInstall resolves a marketplace install target and
// either renders the Tri-Option inline keyboard (opt == "") or executes the
// chosen option and reports the outcome. Mirrors the CLI dialog flow so both
// interfaces stay behaviorally identical.
func handleTelegramMarketplaceInstall(chatID, messageID int64, name, opt string, edit bool) {
	ctx := context.Background()

	// Case B: unlisted upstream runtime from a raw spec.
	if name == "upstream" {
		out, ok := marketplace.InstallUpstreamSpec(opt)
		reportInstallOutcome(chatID, messageID, out, ok, edit)
		return
	}

	entry, err := marketplace.FindEntry(ctx, name)
	if err != nil {
		SendMessage(fmt.Sprintf(
			"💡 No Go port found for <code>%s</code> in the Scorp Marketplace.\n\n"+
				"Install the upstream runtime directly:\n"+
				"<code>/mcp install upstream npx:@scope/some-mcp-server</code>\n"+
				"<code>/mcp install upstream uvx:some-mcp-server</code>", name), BackButtonKeyboard())
		return
	}

	m, err := marketplace.GetManifest(ctx, *entry)
	if err != nil {
		SendMessage(fmt.Sprintf("❌ %v", err), BackButtonKeyboard())
		return
	}

	if opt == "" {
		text := marketplace.DisclosureCard(entry, m) + "\nSelect Installation Method:"
		SendMessage(text, MarketplaceTriOptionKeyboard(name))
		return
	}

	choice := marketplace.ParseInstallOption(opt)
	if choice == marketplace.OptionUnknown {
		SendMessage(fmt.Sprintf("❌ Invalid option %q — use 1, 2, or 3.", opt), BackButtonKeyboard())
		return
	}

	out, ok := marketplace.Install(ctx, entry, choice)
	reportInstallOutcome(chatID, messageID, out, ok, edit)
}

func reportInstallOutcome(chatID, messageID int64, out string, ok bool, edit bool) {
	kb := BackButtonKeyboard()
	if edit && messageID != 0 {
		EditMessage(chatID, messageID, out, kb)
		return
	}
	SendMessage(out, kb)
}

// MarketplaceTriOptionKeyboard renders the three installation choices.
// callback_data stays well under Telegram's 64-byte limit.
func MarketplaceTriOptionKeyboard(name string) map[string]interface{} {
	return map[string]interface{}{
		"inline_keyboard": []interface{}{
			[]interface{}{
				map[string]string{"text": "⚡ Prebuilt Binary (SHA-256 verified)", "callback_data": "mcpm:" + name + ":1"},
			},
			[]interface{}{
				map[string]string{"text": "🛠️ Rebuild Native Go from Source", "callback_data": "mcpm:" + name + ":2"},
			},
			[]interface{}{
				map[string]string{"text": "📦 Upstream Original Runtime", "callback_data": "mcpm:" + name + ":3"},
				map[string]string{"text": "🛒 Marketplace", "callback_data": "mcpm:menu"},
			},
		},
	}
}
