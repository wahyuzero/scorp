# chat.go

> 78 nodes · cohesion 0.06

## Key Concepts

- **chat.go** (41 connections) — `agent/chat.go`
- **RunAgentLoop()** (35 connections) — `agent/loop.go`
- **loop.go** (25 connections) — `agent/loop.go`
- **resumeAgentLoop()** (23 connections) — `agent/loop.go`
- **setSession()** (12 connections) — `agent/chat.go`
- **getSession()** (11 connections) — `agent/chat.go`
- **HandleUploadInAgentMode()** (11 connections) — `agent/loop.go`
- **appendSessionHistory()** (10 connections) — `agent/chat.go`
- **HandleConfirmation()** (10 connections) — `agent/loop.go`
- **StorePendingConfirmation()** (10 connections) — `agent/loop.go`
- **getOrCreateSession()** (9 connections) — `agent/chat.go`
- **getSessionHistory()** (9 connections) — `agent/chat.go`
- **maybeCompactHistory()** (9 connections) — `agent/compaction.go`
- **chatSession** (8 connections) — `agent/chat.go`
- **IsDangerousCommand()** (8 connections) — `agent/prompt.go`
- **ClearChatSession()** (7 connections) — `agent/chat.go`
- **ExitAgentMode()** (7 connections) — `agent/chat.go`
- **getSessionMap()** (7 connections) — `agent/chat.go`
- **ExecuteTool()** (7 connections) — `agent/prompt.go`
- **AgentMessage** (6 connections) — `agent/loop.go`
- **convertTableToList()** (6 connections) — `agent/chat.go`
- **EnterAgentMode()** (6 connections) — `agent/chat.go`
- **SendMessageSmart()** (6 connections) — `agent/chat.go`
- **sendScorpReply()** (6 connections) — `agent/chat.go`
- **summarizeHistory()** (6 connections) — `agent/chat.go`
- *... and 53 more nodes in this community*

## Relationships

- [TruncateStr](TruncateStr.md) (17 shared connections)
- [handleAction](handleAction.md) (13 shared connections)
- [startCLI](startCLI.md) (13 shared connections)
- [GetStringArg](GetStringArg.md) (10 shared connections)
- [time.Time](time.Time.md) (9 shared connections)
- [testing.T](testing.T.md) (9 shared connections)
- [runSelfReview](runSelfReview.md) (6 shared connections)
- [ConfigMgr](ConfigMgr.md) (3 shared connections)
- [RunAgentLoop](RunAgentLoop.md) (2 shared connections)
- [config_paths.go](config_paths.go.md) (2 shared connections)
- [client.go](client.go.md) (2 shared connections)
- [collector_system.go](collector_system.go.md) (1 shared connections)

## Source Files

- `agent/chat.go`
- `agent/compaction.go`
- `agent/loop.go`
- `agent/prompt.go`
- `agent/self_improve.go`
- `agent/sessions.go`
- `tools/memory.go`

## Audit Trail

- EXTRACTED: 239 (85%)
- INFERRED: 43 (15%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*