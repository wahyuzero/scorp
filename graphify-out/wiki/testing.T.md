# testing.T

> 68 nodes · cohesion 0.06

## Key Concepts

- **testing.T** (83 connections)
- **compaction_test.go** (16 connections) — `agent/compaction_test.go`
- **truncateToolResultsInHistory()** (15 connections) — `agent/compaction.go`
- **prompt_test.go** (13 connections) — `agent/prompt_test.go`
- **makeHistory()** (10 connections) — `agent/compaction_test.go`
- **collector_native_test.go** (10 connections) — `collectors/collector_native_test.go`
- **estimateHistoryTokens()** (8 connections) — `agent/compaction.go`
- **SortByCPUDesc()** (7 connections) — `collectors/collector_system_native.go`
- **compaction.go** (5 connections) — `agent/compaction.go`
- **makeToolResult()** (5 connections) — `agent/compaction_test.go`
- **TestPrune_DockerScenario_41Messages()** (5 connections) — `agent/compaction_test.go`
- **TestPrune_TokenSavings()** (5 connections) — `agent/compaction_test.go`
- **AgentMessage** (4 connections)
- **TestEstimateHistoryTokens()** (4 connections) — `agent/compaction_test.go`
- **TestPrune_BoundaryAges()** (4 connections) — `agent/compaction_test.go`
- **TestPrune_NonToolMessages_Preserved()** (4 connections) — `agent/compaction_test.go`
- **TestPrune_OldToolResult_TrimmedTo500()** (4 connections) — `agent/compaction_test.go`
- **TestPrune_RecentOversized_TrimmedTo3000()** (4 connections) — `agent/compaction_test.go`
- **TestPrune_RecentToolResult_KeptFull()** (4 connections) — `agent/compaction_test.go`
- **TestPrune_ShortToolResult_NeverTrimmed()** (4 connections) — `agent/compaction_test.go`
- **TestPrune_VeryOldToolResult_StubOnly()** (4 connections) — `agent/compaction_test.go`
- **GetStringSliceArg()** (4 connections) — `internal/helpers/helpers.go`
- **TestSwitchActiveModel()** (4 connections) — `models/api_commandcode_test.go`
- **session_fallback_test.go** (4 connections) — `session/session_fallback_test.go`
- **estimateTokens()** (3 connections) — `agent/compaction.go`
- *... and 43 more nodes in this community*

## Relationships

- [ConfigMgr](ConfigMgr.md) (10 shared connections)
- [chat.go](chat.go.md) (9 shared connections)
- [collector_system_native.go](collector_system_native.go.md) (9 shared connections)
- [collector_system.go](collector_system.go.md) (9 shared connections)
- [TruncateStr](TruncateStr.md) (7 shared connections)
- [rag_vector.go](rag_vector.go.md) (7 shared connections)
- [config_paths.go](config_paths.go.md) (6 shared connections)
- [GetStringArg](GetStringArg.md) (4 shared connections)
- [startCLI](startCLI.md) (4 shared connections)
- [time.Time](time.Time.md) (3 shared connections)
- [runSelfReview](runSelfReview.md) (2 shared connections)
- [runSubagent](runSubagent.md) (1 shared connections)

## Source Files

- `agent/compaction.go`
- `agent/compaction_test.go`
- `agent/prompt_test.go`
- `agent/self_improve_test.go`
- `collectors/collector_native_test.go`
- `collectors/collector_system_native.go`
- `internal/helpers/helpers.go`
- `models/api_commandcode_test.go`
- `session/session_fallback_test.go`
- `tools_native_test.go`

## Audit Trail

- EXTRACTED: 180 (85%)
- INFERRED: 32 (15%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*