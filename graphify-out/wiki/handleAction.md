# handleAction

> 80 nodes · cohesion 0.06

## Key Concepts

- **handleAction()** (59 connections) — `main.go`
- **main()** (54 connections) — `main.go`
- **telegram.go** (28 connections) — `telegram/telegram.go`
- **files.go** (13 connections) — `telegram/files.go`
- **TgPost()** (13 connections) — `telegram/telegram.go`
- **sendHourlyReport()** (12 connections) — `main.go`
- **commandLoop()** (10 connections) — `main.go`
- **clarify.go** (10 connections) — `tools/clarify.go`
- **inline.go** (10 connections) — `tools/inline.go`
- **main.go** (9 connections) — `main.go`
- **hourlyLoop()** (8 connections) — `main.go`
- **HumanSize()** (8 connections) — `telegram/files.go`
- **SendFolderAsZip()** (8 connections) — `telegram/files.go`
- **SendMessage()** (8 connections) — `telegram/telegram.go`
- **buildInlineResults()** (8 connections) — `tools/inline.go`
- **SendFile()** (7 connections) — `telegram/files.go`
- **securityEventLoop()** (6 connections) — `main.go`
- **DirKeyboard()** (6 connections) — `telegram/files.go`
- **FileDetailKeyboard()** (6 connections) — `telegram/files.go`
- **AnswerCallback()** (6 connections) — `telegram/telegram.go`
- **PollUpdates()** (6 connections) — `telegram/telegram.go`
- **SendMessageGetID()** (6 connections) — `telegram/telegram.go`
- **WebhookHandler()** (6 connections) — `telegram/telegram.go`
- **HandleInlineQuery()** (6 connections) — `tools/inline.go`
- **FolderZipInfo()** (5 connections) — `telegram/files.go`
- *... and 55 more nodes in this community*

## Relationships

- [collector_system.go](collector_system.go.md) (18 shared connections)
- [chat.go](chat.go.md) (13 shared connections)
- [time.Time](time.Time.md) (13 shared connections)
- [config_paths.go](config_paths.go.md) (11 shared connections)
- [client.go](client.go.md) (7 shared connections)
- [wizard.go](wizard.go.md) (6 shared connections)
- [collector_security.go](collector_security.go.md) (5 shared connections)
- [collector_system_native.go](collector_system_native.go.md) (5 shared connections)
- [startCLI](startCLI.md) (5 shared connections)
- [ConfigMgr](ConfigMgr.md) (4 shared connections)
- [checker.go](checker.go.md) (3 shared connections)
- [TruncateStr](TruncateStr.md) (2 shared connections)

## Source Files

- `agent/chat.go`
- `main.go`
- `metrics/metrics.go`
- `telegram/files.go`
- `telegram/telegram.go`
- `tools/clarify.go`
- `tools/inline.go`

## Audit Trail

- EXTRACTED: 270 (95%)
- INFERRED: 13 (5%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*