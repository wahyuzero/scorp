# StartDaemon

> 26 nodes

## Key Concepts

- **StartDaemon()** (43 connections) — `telegram/daemon.go`
- **telegram.go** (29 connections) — `telegram/telegram.go`
- **TgPost()** (13 connections) — `telegram/telegram.go`
- **runCommandLoop()** (10 connections) — `telegram/daemon.go`
- **SendMessage()** (8 connections) — `telegram/telegram.go`
- **PollUpdates()** (6 connections) — `telegram/telegram.go`
- **EditMessageByID()** (5 connections) — `telegram/telegram.go`
- **SendMessageGetID()** (5 connections) — `telegram/telegram.go`
- **EditMessage()** (4 connections) — `telegram/telegram.go`
- **UploadsDir()** (3 connections) — `config/config_paths.go`
- **daemon.go** (3 connections) — `telegram/daemon.go`
- **DeleteWebhook()** (3 connections) — `telegram/telegram.go`
- **SendChatAction()** (3 connections) — `telegram/telegram.go`
- **SendMessageGetIDWithKeyboard()** (3 connections) — `telegram/telegram.go`
- **SettingsMenuText()** (3 connections) — `telegram/telegram.go`
- **SetupBotCommands()** (3 connections) — `telegram/telegram.go`
- **SetWebhook()** (3 connections) — `telegram/telegram.go`
- **StartWebhookServer()** (3 connections) — `telegram/telegram.go`
- **StopWebhookServer()** (3 connections) — `telegram/telegram.go`
- **BackAndRefreshKeyboard()** (2 connections) — `telegram/telegram.go`
- **InitTelegram()** (2 connections) — `telegram/telegram.go`
- **ReplyMenuKeyboard()** (2 connections) — `telegram/telegram.go`
- **TGCallback** (2 connections) — `telegram/telegram.go`
- **TGCommand** (2 connections) — `telegram/telegram.go`
- **TgResponse** (2 connections) — `telegram/telegram.go`
- *... and 1 more nodes in this community*

## Relationships

- [HandleTelegramAction](HandleTelegramAction.md) (17 shared connections)
- [clarify.go](clarify.go.md) (5 shared connections)
- [startCLI](startCLI.md) (5 shared connections)
- [time.Time](time.Time.md) (5 shared connections)
- [inline.go](inline.go.md) (4 shared connections)
- [Manifest](Manifest.md) (4 shared connections)
- [chat.go](chat.go.md) (3 shared connections)
- [ScorpPath](ScorpPath.md) (2 shared connections)
- [client.go](client.go.md) (2 shared connections)
- [cost_router.go](cost_router.go.md) (2 shared connections)
- [init](init.md) (2 shared connections)
- [LoadConfig](LoadConfig.md) (2 shared connections)

## Source Files

- `config/config_paths.go`
- `telegram/daemon.go`
- `telegram/telegram.go`
- `tools/inline.go`

## Audit Trail

- EXTRACTED: 89 (75%)
- INFERRED: 29 (25%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*