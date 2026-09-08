# runSubagent

> 74 nodes

## Key Concepts

- **runSubagent()** (19 connections) — `delegate/run.go`
- **acp.go** (13 connections) — `delegate/acp.go`
- **ACPSession** (13 connections) — `delegate/acp.go`
- **execute.go** (12 connections) — `delegate/execute.go`
- **bg.go** (12 connections) — `tools/bg.go`
- **ParseDelegateParams()** (11 connections) — `delegate/execute.go`
- **ExecuteDelegate()** (10 connections) — `delegate/execute.go`
- **ExecuteDelegateBatch()** (10 connections) — `delegate/execute.go`
- **sync.Mutex** (10 connections)
- **runSubagentACP()** (9 connections) — `delegate/acp.go`
- **delegateTaskParams** (9 connections) — `delegate/execute.go`
- **isolation.go** (9 connections) — `delegate/isolation.go`
- **launchACP()** (8 connections) — `delegate/acp.go`
- **SubagentRole** (8 connections) — `delegate/execute.go`
- **ExecuteBgProcess()** (8 connections) — `tools/bg.go`
- **SubagentIsolation** (7 connections) — `delegate/isolation.go`
- **getBGProcess()** (7 connections) — `tools/bg.go`
- **BGProcess** (7 connections) — `tools/bg.go`
- **runOpenCodeCLI()** (6 connections) — `delegate/acp.go`
- **delegateResult** (6 connections) — `delegate/execute.go`
- **ACPResponse** (5 connections) — `delegate/acp.go`
- **FormatDelegateResult()** (5 connections) — `delegate/execute.go`
- **ValidateSubagentTools()** (5 connections) — `delegate/execute.go`
- **buildSubagentPrompt()** (5 connections) — `delegate/prompt.go`
- **StartTestEndpoint()** (5 connections) — `testutil/endpoint.go`
- *... and 49 more nodes in this community*

## Relationships

- [testing.T](testing.T.md) (3 shared connections)
- [init](init.md) (3 shared connections)
- [TruncateStr](TruncateStr.md) (3 shared connections)
- [GetIntArg](GetIntArg.md) (3 shared connections)
- [client.go](client.go.md) (2 shared connections)
- [context.Context](context.Context.md) (2 shared connections)
- [time.Time](time.Time.md) (2 shared connections)
- [RegisterTool](RegisterTool.md) (2 shared connections)
- [registry/registry.go](registry-registry.go.md) (2 shared connections)
- [cost_router.go](cost_router.go.md) (2 shared connections)
- [ToolCall](ToolCall.md) (2 shared connections)
- [MCPServer](MCPServer.md) (2 shared connections)

## Source Files

- `agent/auto.go`
- `delegate/acp.go`
- `delegate/delegate_test.go`
- `delegate/execute.go`
- `delegate/isolation.go`
- `delegate/prompt.go`
- `delegate/run.go`
- `docs/IMPLEMENTATION_PLAN_SCORP.md`
- `testutil/endpoint.go`
- `tools/bg.go`

## Audit Trail

- EXTRACTED: 171 (91%)
- INFERRED: 16 (9%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*