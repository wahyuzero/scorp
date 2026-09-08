# RunAgentSessionLoop

> 20 nodes

## Key Concepts

- **RunAgentSessionLoop()** (63 connections) — `agent/loop.go`
- **resumeAgentLoop()** (40 connections) — `agent/loop.go`
- **loop.go** (8 connections) — `agent/loop.go`
- **confirmationDisplay()** (6 connections) — `agent/auto.go`
- **sendScorpReply()** (6 connections) — `agent/chat.go`
- **buildThinkingMessage()** (6 connections) — `agent/thinking.go`
- **toolDescription()** (6 connections) — `agent/thinking.go`
- **maxTurnTimeout()** (5 connections) — `agent/loop.go`
- **RunAgentLoop()** (5 connections) — `agent/loop.go`
- **ConsumeStopRequest()** (4 connections) — `agent/chat.go`
- **cleanToolCallTags()** (4 connections) — `agent/loop.go`
- **maxIterations()** (4 connections) — `agent/loop.go`
- **thinking.go** (4 connections) — `agent/thinking.go`
- **shouldUpdateThinking()** (4 connections) — `agent/thinking.go`
- **toolCallSignature()** (4 connections) — `agent/thinking.go`
- **confirmKeyboard()** (3 connections) — `agent/confirmation.go`
- **TickActiveSkills()** (3 connections) — `skills/v2_skills.go`
- **AgentMessage** (2 connections) — `agent/loop.go`
- **ClearStopRequest()** (2 connections) — `agent/chat.go`
- **getSessionSearchContext()** (2 connections) — `agent/loop.go`

## Relationships

- [chat.go](chat.go.md) (13 shared connections)
- [startCLI](startCLI.md) (9 shared connections)
- [ToolCall](ToolCall.md) (6 shared connections)
- [ClearTaskPlan](ClearTaskPlan.md) (6 shared connections)
- [time.Time](time.Time.md) (5 shared connections)
- [HasGreenTestRun](HasGreenTestRun.md) (5 shared connections)
- [GetStringArg](GetStringArg.md) (4 shared connections)
- [PermissionDecision](PermissionDecision.md) (4 shared connections)
- [runPlanningTurns](runPlanningTurns.md) (4 shared connections)
- [GetAutonomyLevel](GetAutonomyLevel.md) (4 shared connections)
- [ExecuteTermuxAPI](ExecuteTermuxAPI.md) (4 shared connections)
- [TruncateStr](TruncateStr.md) (3 shared connections)

## Source Files

- `agent/auto.go`
- `agent/chat.go`
- `agent/confirmation.go`
- `agent/loop.go`
- `agent/thinking.go`
- `skills/v2_skills.go`

## Audit Trail

- EXTRACTED: 76 (54%)
- INFERRED: 65 (46%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*