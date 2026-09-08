# chat.go

> 81 nodes

## Key Concepts

- **chat.go** (46 connections) — `agent/chat.go`
- **maybeCompactHistory()** (16 connections) — `agent/compaction.go`
- **getSession()** (15 connections) — `agent/chat.go`
- **setSession()** (14 connections) — `agent/chat.go`
- **estimateHistoryTokens()** (13 connections) — `agent/compaction.go`
- **CompactSessionHistory()** (12 connections) — `agent/compact_manual.go`
- **RenameSession()** (12 connections) — `agent/session_mgr.go`
- **HandleUploadInAgentMode()** (11 connections) — `agent/upload.go`
- **appendSessionHistory()** (10 connections) — `agent/chat.go`
- **ClearChatSession()** (10 connections) — `agent/chat.go`
- **getOrCreateSession()** (10 connections) — `agent/chat.go`
- **getSessionHistory()** (10 connections) — `agent/chat.go`
- **compaction.go** (10 connections) — `agent/compaction.go`
- **setLoopActive()** (9 connections) — `agent/chat.go`
- **summarizeHistory()** (9 connections) — `agent/chat.go`
- **AgentMessage** (9 connections)
- **preservationNote()** (9 connections) — `agent/compaction.go`
- **TestSessionManager()** (9 connections) — `agent/session_mgr_test.go`
- **session_ui.go** (9 connections) — `telegram/session_ui.go`
- **HandleSessionCallback()** (9 connections) — `telegram/session_ui.go`
- **getSessionMap()** (8 connections) — `agent/chat.go`
- **historyFilePath()** (8 connections) — `agent/chat.go`
- **DeleteSession()** (8 connections) — `agent/session_mgr.go`
- **ListSessions()** (8 connections) — `agent/session_mgr.go`
- **saveHistoryToDisk()** (7 connections) — `agent/chat.go`
- *... and 56 more nodes in this community*

## Relationships

- [compaction_test.go](compaction_test.go.md) (16 shared connections)
- [RunAgentSessionLoop](RunAgentSessionLoop.md) (13 shared connections)
- [time.Time](time.Time.md) (12 shared connections)
- [HandleTelegramAction](HandleTelegramAction.md) (12 shared connections)
- [startCLI](startCLI.md) (11 shared connections)
- [runSelfReview](runSelfReview.md) (3 shared connections)
- [StartDaemon](StartDaemon.md) (3 shared connections)
- [TruncateStr](TruncateStr.md) (3 shared connections)
- [context.Context](context.Context.md) (3 shared connections)
- [testing.T](testing.T.md) (3 shared connections)
- [ScorpPath](ScorpPath.md) (3 shared connections)
- [runPlanningTurns](runPlanningTurns.md) (2 shared connections)

## Source Files

- `agent/chat.go`
- `agent/compact_manual.go`
- `agent/compaction.go`
- `agent/prompt_test.go`
- `agent/session_mgr.go`
- `agent/session_mgr_test.go`
- `agent/sessions.go`
- `agent/upload.go`
- `config/config_paths.go`
- `telegram/session_ui.go`

## Audit Trail

- EXTRACTED: 223 (76%)
- INFERRED: 71 (24%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*