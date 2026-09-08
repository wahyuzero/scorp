# SetTaskPlan

> 14 nodes

## Key Concepts

- **SetTaskPlan()** (12 connections) — `agent/taskplan.go`
- **planFilePath()** (11 connections) — `agent/taskplan.go`
- **TestClearTaskPlanRemovesPersistedFile()** (9 connections) — `agent/taskplan_persist_test.go`
- **taskplan_persist_test.go** (8 connections) — `agent/taskplan_persist_test.go`
- **TestPlanPersistenceRoundtripAcrossRestart()** (8 connections) — `agent/taskplan_persist_test.go`
- **TestStalePlanExpiredOnLoad()** (8 connections) — `agent/taskplan_persist_test.go`
- **TestUpdateItemStatusPersistsThroughRestart()** (7 connections) — `agent/taskplan_persist_test.go`
- **uniqueSess()** (7 connections) — `agent/taskplan_persist_test.go`
- **mkTestPlan()** (6 connections) — `agent/taskplan_persist_test.go`
- **simulateRestart()** (5 connections) — `agent/taskplan_persist_test.go`
- **savePlanToDisk()** (4 connections) — `agent/taskplan.go`
- **caseLedgerClear()** (4 connections) — `eval/core.go`
- **caseLedgerPersisted()** (4 connections) — `eval/core.go`
- **TestPlanFilePathSanitizesSessionID()** (3 connections) — `agent/taskplan_persist_test.go`

## Relationships

- [ClearTaskPlan](ClearTaskPlan.md) (11 shared connections)
- [testing.T](testing.T.md) (6 shared connections)
- [TaskPlan](TaskPlan.md) (5 shared connections)
- [ScorpPath](ScorpPath.md) (3 shared connections)
- [chat.go](chat.go.md) (2 shared connections)
- [compaction_test.go](compaction_test.go.md) (2 shared connections)
- [GetAutonomyLevel](GetAutonomyLevel.md) (2 shared connections)
- [runPlanningTurns](runPlanningTurns.md) (1 shared connections)

## Source Files

- `agent/taskplan.go`
- `agent/taskplan_persist_test.go`
- `eval/core.go`

## Audit Trail

- EXTRACTED: 45 (70%)
- INFERRED: 19 (30%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*