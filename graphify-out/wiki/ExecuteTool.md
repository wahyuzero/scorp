# ExecuteTool

> 11 nodes

## Key Concepts

- **ExecuteTool()** (31 connections) — `agent/prompt.go`
- **registerHookProbeTool()** (7 connections) — `agent/hooks_agent_test.go`
- **setHookEnvAgent()** (6 connections) — `agent/hooks_agent_test.go`
- **hooks_agent_test.go** (5 connections) — `agent/hooks_agent_test.go`
- **TestExecuteToolHookContextAppended()** (5 connections) — `agent/hooks_agent_test.go`
- **TestExecuteToolHookScopedToOtherToolDoesNotFire()** (5 connections) — `agent/hooks_agent_test.go`
- **TestExecuteToolPreHookBlocks()** (5 connections) — `agent/hooks_agent_test.go`
- **TestExecuteToolDenyRulesFirst()** (4 connections) — `agent/deny_test.go`
- **agent/prompt.go** (3 connections) — `agent/prompt.go`
- **FormatToolResult()** (2 connections) — `agent/prompt.go`
- **agent/deny_test.go** (1 connections) — `agent/deny_test.go`

## Relationships

- [testing.T](testing.T.md) (6 shared connections)
- [RegisterTool](RegisterTool.md) (4 shared connections)
- [GetAutonomyLevel](GetAutonomyLevel.md) (4 shared connections)
- [PermissionDecision](PermissionDecision.md) (4 shared connections)
- [CheckDenyRules](CheckDenyRules.md) (2 shared connections)
- [runPlanningTurns](runPlanningTurns.md) (2 shared connections)
- [RunPreToolUseHooks](RunPreToolUseHooks.md) (2 shared connections)
- [RecordToolReceipt](RecordToolReceipt.md) (2 shared connections)
- [startCLI](startCLI.md) (2 shared connections)
- [ToolCall](ToolCall.md) (2 shared connections)
- [RunAgentSessionLoop](RunAgentSessionLoop.md) (2 shared connections)
- [config/hooks.go](config-hooks.go.md) (1 shared connections)

## Source Files

- `agent/deny_test.go`
- `agent/hooks_agent_test.go`
- `agent/prompt.go`

## Audit Trail

- EXTRACTED: 42 (75%)
- INFERRED: 14 (25%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*