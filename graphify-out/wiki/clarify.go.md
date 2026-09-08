# clarify.go

> 12 nodes

## Key Concepts

- **clarify.go** (10 connections) — `tools/clarify.go`
- **AnswerCallback()** (6 connections) — `telegram/telegram.go`
- **executeClarify()** (4 connections) — `tools/clarify.go`
- **init()** (4 connections) — `tools/clarify.go`
- **ResolveClarify()** (4 connections) — `tools/clarify.go`
- **sendClarifyMessage()** (4 connections) — `tools/clarify.go`
- **GetClarifyChatID()** (2 connections) — `tools/clarify.go`
- **handleClarifyResponse()** (2 connections) — `tools/clarify.go`
- **HasPendingClarify()** (2 connections) — `tools/clarify.go`
- **SetClarifyChatID()** (2 connections) — `tools/clarify.go`
- **PendingClarify** (2 connections) — `tools/clarify.go`
- **cleanupStaleClarifies()** (1 connections) — `tools/clarify.go`

## Relationships

- [StartDaemon](StartDaemon.md) (5 shared connections)
- [HandleTelegramAction](HandleTelegramAction.md) (3 shared connections)
- [TestGatewayEndpoints](TestGatewayEndpoints.md) (1 shared connections)
- [RegisterTool](RegisterTool.md) (1 shared connections)
- [time.Time](time.Time.md) (1 shared connections)

## Source Files

- `telegram/telegram.go`
- `tools/clarify.go`

## Audit Trail

- EXTRACTED: 22 (81%)
- INFERRED: 5 (19%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*