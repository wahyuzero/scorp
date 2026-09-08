# startCLI

> 77 nodes

## Key Concepts

- **startCLI()** (51 connections) — `cli.go`
- **main()** (21 connections) — `main.go`
- **cli.go** (17 connections) — `cli.go`
- **executeOneShot()** (17 connections) — `cli.go`
- **HandleConfirmation()** (13 connections) — `agent/confirmation.go`
- **StorePendingConfirmation()** (11 connections) — `agent/confirmation.go`
- **wireCLICallbacks()** (11 connections) — `cli_callbacks.go`
- **v2_skills.go** (11 connections) — `skills/v2_skills.go`
- **handleMCPCommand()** (10 connections) — `cli_mcp.go`
- **confirmation.go** (9 connections) — `agent/confirmation.go`
- **printStatus()** (8 connections) — `cli.go`
- **InitDefaultSOPs()** (8 connections) — `sop/sop.go`
- **formatTerminalText()** (7 connections) — `cli_format.go`
- **handleCLISession()** (7 connections) — `cli.go`
- **SOP** (7 connections) — `sop/sop.go`
- **ListSOPs()** (7 connections) — `sop/sop.go`
- **ExecuteSOP()** (7 connections) — `tools/sop.go`
- **getPendingConfirmation()** (6 connections) — `agent/confirmation.go`
- **StorePendingConfirmationArgs()** (6 connections) — `agent/confirmation.go`
- **executeTurn()** (6 connections) — `cli.go`
- **handleCLISOP()** (6 connections) — `cli.go`
- **printCostUsage()** (6 connections) — `cli.go`
- **printCurrentModel()** (6 connections) — `cli.go`
- **LoadAllSkills()** (6 connections) — `skills/v2_skills.go`
- **Dir()** (6 connections) — `sop/sop.go`
- *... and 52 more nodes in this community*

## Relationships

- [chat.go](chat.go.md) (11 shared connections)
- [RunAgentSessionLoop](RunAgentSessionLoop.md) (9 shared connections)
- [GetAutonomyLevel](GetAutonomyLevel.md) (9 shared connections)
- [ScorpPath](ScorpPath.md) (7 shared connections)
- [testing.T](testing.T.md) (6 shared connections)
- [HandleTelegramAction](HandleTelegramAction.md) (6 shared connections)
- [time.Time](time.Time.md) (6 shared connections)
- [TruncateStr](TruncateStr.md) (6 shared connections)
- [StartDaemon](StartDaemon.md) (5 shared connections)
- [runPlanningTurns](runPlanningTurns.md) (4 shared connections)
- [context.Context](context.Context.md) (4 shared connections)
- [CreateCheckpoint](CreateCheckpoint.md) (4 shared connections)

## Source Files

- `agent/auto_test.go`
- `agent/chat.go`
- `agent/confirmation.go`
- `bootstrap/autonomous.go`
- `cli.go`
- `cli_callbacks.go`
- `cli_format.go`
- `cli_lock.go`
- `cli_lock_test.go`
- `cli_mcp.go`
- `cli_test.go`
- `gateway/gateway.go`
- `main.go`
- `models/model_router.go`
- `skills/v2_skills.go`
- `sop/sop.go`
- `sop/sop_test.go`
- `tools/skill_activate.go`
- `tools/sop.go`

## Audit Trail

- EXTRACTED: 236 (86%)
- INFERRED: 40 (14%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*