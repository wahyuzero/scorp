# runPlanningTurns

> 14 nodes

## Key Concepts

- **runPlanningTurns()** (13 connections) — `agent/planmode.go`
- **getAgentSystemPrompt()** (11 connections) — `agent/prompt.go`
- **planmode.go** (9 connections) — `agent/planmode.go`
- **RunPlanningLoop()** (9 connections) — `agent/planmode.go`
- **RevisePlan()** (8 connections) — `agent/planmode.go`
- **BeginPlanning()** (6 connections) — `agent/planmode.go`
- **EndPlanning()** (6 connections) — `agent/planmode.go`
- **TestPlanningStateTransitions()** (6 connections) — `agent/planmode_test.go`
- **PlanningState()** (2 connections) — `agent/planmode.go`
- **FormatSkillsIndexForSystemPrompt()** (2 connections) — `skills/v2_skills.go`
- **GetActiveSkillsContext()** (2 connections) — `skills/v2_skills.go`
- **Plan Mode Workflow (P1.4)** (2 connections) — `docs/IMPLEMENTATION_PLAN_SCORP.md`
- **Persistent Task Ledger (P1.5)** (2 connections) — `docs/IMPLEMENTATION_PLAN_SCORP.md`
- **AgentMessage** (1 connections)

## Relationships

- [ClearTaskPlan](ClearTaskPlan.md) (7 shared connections)
- [startCLI](startCLI.md) (4 shared connections)
- [RunAgentSessionLoop](RunAgentSessionLoop.md) (4 shared connections)
- [GetAutonomyLevel](GetAutonomyLevel.md) (3 shared connections)
- [HandleTelegramAction](HandleTelegramAction.md) (2 shared connections)
- [chat.go](chat.go.md) (2 shared connections)
- [ExecuteTool](ExecuteTool.md) (2 shared connections)
- [runSelfReview](runSelfReview.md) (2 shared connections)
- [time.Time](time.Time.md) (1 shared connections)
- [TruncateStr](TruncateStr.md) (1 shared connections)
- [ToolCall](ToolCall.md) (1 shared connections)
- [TestSteeringQueue](TestSteeringQueue.md) (1 shared connections)

## Source Files

- `agent/planmode.go`
- `agent/planmode_test.go`
- `agent/prompt.go`
- `docs/IMPLEMENTATION_PLAN_SCORP.md`
- `skills/v2_skills.go`

## Audit Trail

- EXTRACTED: 36 (63%)
- INFERRED: 21 (37%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*