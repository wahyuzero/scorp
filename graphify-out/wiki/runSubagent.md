# runSubagent

> 51 nodes · cohesion 0.08

## Key Concepts

- **runSubagent()** (19 connections) — `delegate/run.go`
- **acp.go** (13 connections) — `delegate/acp.go`
- **ACPSession** (13 connections) — `delegate/acp.go`
- **execute.go** (12 connections) — `delegate/execute.go`
- **ExecuteDelegate()** (10 connections) — `delegate/execute.go`
- **ExecuteDelegateBatch()** (10 connections) — `delegate/execute.go`
- **runSubagentACP()** (9 connections) — `delegate/acp.go`
- **ParseDelegateParams()** (9 connections) — `delegate/execute.go`
- **isolation.go** (9 connections) — `delegate/isolation.go`
- **launchACP()** (8 connections) — `delegate/acp.go`
- **delegateTaskParams** (8 connections) — `delegate/execute.go`
- **SubagentRole** (8 connections) — `delegate/execute.go`
- **SubagentIsolation** (7 connections) — `delegate/isolation.go`
- **runOpenCodeCLI()** (6 connections) — `delegate/acp.go`
- **delegateResult** (6 connections) — `delegate/execute.go`
- **ACPResponse** (5 connections) — `delegate/acp.go`
- **FormatDelegateResult()** (5 connections) — `delegate/execute.go`
- **buildSubagentPrompt()** (5 connections) — `delegate/prompt.go`
- **.send()** (4 connections) — `delegate/acp.go`
- **ValidateSubagentTools()** (4 connections) — `delegate/execute.go`
- **.Close()** (3 connections) — `delegate/acp.go`
- **.initialize()** (3 connections) — `delegate/acp.go`
- **.SendMessage()** (3 connections) — `delegate/acp.go`
- **ACPUserMessage** (3 connections) — `delegate/acp.go`
- **DefaultSubagentTools()** (3 connections) — `delegate/execute.go`
- *... and 26 more nodes in this community*

## Relationships

- [TruncateStr](TruncateStr.md) (8 shared connections)
- [GetStringArg](GetStringArg.md) (6 shared connections)
- [client.go](client.go.md) (3 shared connections)
- [RunAgentLoop](RunAgentLoop.md) (3 shared connections)
- [config_paths.go](config_paths.go.md) (1 shared connections)
- [time.Time](time.Time.md) (1 shared connections)
- [collector_system_native.go](collector_system_native.go.md) (1 shared connections)
- [testing.T](testing.T.md) (1 shared connections)
- [ConfigMgr](ConfigMgr.md) (1 shared connections)

## Source Files

- `delegate/acp.go`
- `delegate/execute.go`
- `delegate/isolation.go`
- `delegate/prompt.go`
- `delegate/run.go`

## Audit Trail

- EXTRACTED: 117 (90%)
- INFERRED: 13 (10%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*