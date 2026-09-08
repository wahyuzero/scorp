# ClearTaskPlan

> 15 nodes

## Key Concepts

- **ClearTaskPlan()** (20 connections) — `agent/taskplan.go`
- **GetTaskPlan()** (19 connections) — `agent/taskplan.go`
- **execTaskPlanTool()** (15 connections) — `agent/taskplan.go`
- **TestRollUpSkipsActivePlan()** (6 connections) — `agent/history_window_test.go`
- **TestCancelPlanDropsLedger()** (6 connections) — `agent/planmode_test.go`
- **ApprovePlan()** (5 connections) — `agent/planmode.go`
- **CancelPlan()** (5 connections) — `agent/planmode.go`
- **taskplan_test.go** (5 connections) — `agent/taskplan_test.go`
- **TestTaskPlanConcurrentAccess()** (5 connections) — `agent/taskplan_test.go`
- **TestTaskPlanCreateUpdateLifecycle()** (5 connections) — `agent/taskplan_test.go`
- **TestTaskPlanRender()** (5 connections) — `agent/taskplan_test.go`
- **TestTaskPlanSessionsIsolated()** (5 connections) — `agent/taskplan_test.go`
- **TestApprovePlanWithoutLedger()** (4 connections) — `agent/planmode_test.go`
- **TestTaskPlanValidation()** (4 connections) — `agent/taskplan_test.go`
- **agent/planmode_test.go** (3 connections) — `agent/planmode_test.go`

## Relationships

- [SetTaskPlan](SetTaskPlan.md) (11 shared connections)
- [testing.T](testing.T.md) (8 shared connections)
- [runPlanningTurns](runPlanningTurns.md) (7 shared connections)
- [TaskPlan](TaskPlan.md) (7 shared connections)
- [RunAgentSessionLoop](RunAgentSessionLoop.md) (6 shared connections)
- [prepareNewTurnHistory](prepareNewTurnHistory.md) (4 shared connections)
- [HandleTelegramAction](HandleTelegramAction.md) (2 shared connections)
- [chat.go](chat.go.md) (2 shared connections)
- [startCLI](startCLI.md) (1 shared connections)
- [compaction_test.go](compaction_test.go.md) (1 shared connections)
- [GetStringArg](GetStringArg.md) (1 shared connections)

## Source Files

- `agent/history_window_test.go`
- `agent/planmode.go`
- `agent/planmode_test.go`
- `agent/taskplan.go`
- `agent/taskplan_test.go`

## Audit Trail

- EXTRACTED: 36 (44%)
- INFERRED: 45 (56%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*