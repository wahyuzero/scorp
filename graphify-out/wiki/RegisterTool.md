# RegisterTool

> 17 nodes

## Key Concepts

- **RegisterTool()** (26 connections) — `registry/registry.go`
- **UnregisterTool()** (14 connections) — `registry/registry.go`
- **TestExecuteToolAutoPresetTrustedAndRecorded()** (7 connections) — `agent/auto_test.go`
- **caseHooksBlockAndContext()** (6 connections) — `eval/core.go`
- **TestValidateSubagentTools()** (5 connections) — `delegate/delegate_test.go`
- **init()** (2 connections) — `bootstrap/browser.go`
- **init()** (2 connections) — `bootstrap/monitor.go`
- **init()** (2 connections) — `bootstrap/provider.go`
- **init()** (2 connections) — `bootstrap/script.go`
- **init()** (2 connections) — `tools/complete_task.go`
- **init()** (2 connections) — `tools/task_plan.go`
- **bootstrap/browser.go** (1 connections) — `bootstrap/browser.go`
- **bootstrap/monitor.go** (1 connections) — `bootstrap/monitor.go`
- **bootstrap/provider.go** (1 connections) — `bootstrap/provider.go`
- **bootstrap/script.go** (1 connections) — `bootstrap/script.go`
- **complete_task.go** (1 connections) — `tools/complete_task.go`
- **task_plan.go** (1 connections) — `tools/task_plan.go`

## Relationships

- [registry/registry.go](registry-registry.go.md) (13 shared connections)
- [PermissionDecision](PermissionDecision.md) (6 shared connections)
- [ExecuteTool](ExecuteTool.md) (4 shared connections)
- [GetAutonomyLevel](GetAutonomyLevel.md) (4 shared connections)
- [testing.T](testing.T.md) (2 shared connections)
- [runSubagent](runSubagent.md) (2 shared connections)
- [startCLI](startCLI.md) (2 shared connections)
- [RecordToolReceipt](RecordToolReceipt.md) (1 shared connections)
- [config/hooks.go](config-hooks.go.md) (1 shared connections)
- [GetStringArg](GetStringArg.md) (1 shared connections)
- [init](init.md) (1 shared connections)
- [registerMCPToolsAsNative](registerMCPToolsAsNative.md) (1 shared connections)

## Source Files

- `agent/auto_test.go`
- `bootstrap/browser.go`
- `bootstrap/monitor.go`
- `bootstrap/provider.go`
- `bootstrap/script.go`
- `delegate/delegate_test.go`
- `eval/core.go`
- `registry/registry.go`
- `tools/complete_task.go`
- `tools/task_plan.go`

## Audit Trail

- EXTRACTED: 51 (88%)
- INFERRED: 7 (12%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*