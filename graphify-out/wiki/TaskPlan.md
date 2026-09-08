# TaskPlan

> 10 nodes

## Key Concepts

- **TaskPlan** (21 connections) — `agent/taskplan.go`
- **PlanItem** (5 connections) — `agent/taskplan.go`
- **renderPlanItems()** (5 connections) — `agent/taskplan.go`
- **.snapshot()** (4 connections) — `agent/taskplan.go`
- **loadPlanFromDisk()** (3 connections) — `agent/taskplan.go`
- **.Unfinished()** (3 connections) — `agent/taskplan.go`
- **.UpdateItemStatus()** (3 connections) — `agent/taskplan.go`
- **.Progress()** (2 connections) — `agent/taskplan.go`
- **.Render()** (2 connections) — `agent/taskplan.go`
- **.Total()** (1 connections) — `agent/taskplan.go`

## Relationships

- [ClearTaskPlan](ClearTaskPlan.md) (7 shared connections)
- [SetTaskPlan](SetTaskPlan.md) (5 shared connections)
- [time.Time](time.Time.md) (2 shared connections)
- [RunAgentSessionLoop](RunAgentSessionLoop.md) (2 shared connections)
- [metasearch_engines.go](metasearch_engines.go.md) (1 shared connections)

## Source Files

- `agent/taskplan.go`

## Audit Trail

- EXTRACTED: 31 (94%)
- INFERRED: 2 (6%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*