# compaction_test.go

> 28 nodes

## Key Concepts

- **compaction_test.go** (22 connections) — `agent/compaction_test.go`
- **truncateToolResultsInHistory()** (18 connections) — `agent/compaction.go`
- **makeHistory()** (10 connections) — `agent/compaction_test.go`
- **compactionTestSetup()** (8 connections) — `agent/compaction_test.go`
- **TestActiveLoopCompactionPreservesContext()** (8 connections) — `agent/compaction_test.go`
- **makeToolResult()** (5 connections) — `agent/compaction_test.go`
- **TestPreservationNoteContents()** (5 connections) — `agent/compaction_test.go`
- **TestPrune_DockerScenario_41Messages()** (5 connections) — `agent/compaction_test.go`
- **TestPrune_TokenSavings()** (5 connections) — `agent/compaction_test.go`
- **TestCompactionNoopUnderThreshold()** (4 connections) — `agent/compaction_test.go`
- **TestEstimateHistoryTokens()** (4 connections) — `agent/compaction_test.go`
- **TestPrune_BoundaryAges()** (4 connections) — `agent/compaction_test.go`
- **TestPrune_NonToolMessages_Preserved()** (4 connections) — `agent/compaction_test.go`
- **TestPrune_OldToolResult_TrimmedTo500()** (4 connections) — `agent/compaction_test.go`
- **TestPrune_RecentOversized_TrimmedTo3000()** (4 connections) — `agent/compaction_test.go`
- **TestPrune_RecentToolResult_KeptFull()** (4 connections) — `agent/compaction_test.go`
- **TestPrune_ShortToolResult_NeverTrimmed()** (4 connections) — `agent/compaction_test.go`
- **TestPrune_VeryOldToolResult_StubOnly()** (4 connections) — `agent/compaction_test.go`
- **mkCompactionHistory()** (3 connections) — `agent/compaction_test.go`
- **TestEstimateTokens()** (3 connections) — `agent/compaction_test.go`
- **TestPreservationNoteSkipsMachinery()** (3 connections) — `agent/compaction_test.go`
- **TestPrune_EmptyHistory()** (3 connections) — `agent/compaction_test.go`
- **TestPrune_HeadTailFormat()** (3 connections) — `agent/compaction_test.go`
- **TestPrune_SingleMessage()** (3 connections) — `agent/compaction_test.go`
- **SummarizeOldToolResult()** (3 connections) — `agent/micro_summary.go`
- *... and 3 more nodes in this community*

## Relationships

- [testing.T](testing.T.md) (19 shared connections)
- [chat.go](chat.go.md) (16 shared connections)
- [SetTaskPlan](SetTaskPlan.md) (2 shared connections)
- [ClearTaskPlan](ClearTaskPlan.md) (1 shared connections)
- [context.Context](context.Context.md) (1 shared connections)

## Source Files

- `agent/compaction.go`
- `agent/compaction_test.go`
- `agent/micro_summary.go`
- `docs/IMPLEMENTATION_PLAN_SCORP.md`

## Audit Trail

- EXTRACTED: 64 (69%)
- INFERRED: 29 (31%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*