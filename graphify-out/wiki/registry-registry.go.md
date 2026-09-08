# registry/registry.go

> 50 nodes

## Key Concepts

- **registry/registry.go** (13 connections) — `registry/registry.go`
- **ResetNativeToolCache()** (11 connections) — `registry/registry.go`
- **GenerateNativeToolsSchema()** (10 connections) — `registry/registry.go`
- **GetAllTools()** (10 connections) — `registry/registry.go`
- **IsToolActive()** (9 connections) — `registry/dynamic.go`
- **TestIsToolActiveStaticModeWithTTL()** (9 connections) — `registry/dynamic_test.go`
- **registerTempTool()** (8 connections) — `registry/dynamic_test.go`
- **GetTool()** (8 connections) — `registry/registry.go`
- **ExecuteToolSearch()** (8 connections) — `tools/deferred.go`
- **dynamic.go** (7 connections) — `registry/dynamic.go`
- **clearDynamicEnv()** (7 connections) — `registry/dynamic_test.go`
- **ToolDef** (7 connections) — `registry/registry.go`
- **TestToolSearchActivatesDeferredToolInStaticMode()** (7 connections) — `tools/deferred_activation_test.go`
- **TestGenerateNativeToolsSchema()** (6 connections) — `models/tools_test.go`
- **ActivateToolWithTTL()** (6 connections) — `registry/dynamic.go`
- **dynamic_test.go** (6 connections) — `registry/dynamic_test.go`
- **TestDynamicToolTTL()** (6 connections) — `registry/dynamic_test.go`
- **TestIsToolActiveNonDeferredAlwaysActiveInStaticMode()** (6 connections) — `registry/dynamic_test.go`
- **TickToolTTL()** (6 connections) — `registry/dynamic.go`
- **TestRegisterPlugin()** (6 connections) — `registry/plugin_test.go`
- **registerSearchTarget()** (6 connections) — `tools/deferred_activation_test.go`
- **ExecuteToolCall()** (6 connections) — `tools/deferred.go`
- **ExecuteToolList()** (6 connections) — `tools/deferred.go`
- **TestGenerateNativeToolsSchemaExcludesDeferred()** (5 connections) — `registry/dynamic_test.go`
- **RegisterPlugin()** (5 connections) — `registry/plugin.go`
- *... and 25 more nodes in this community*

## Relationships

- [RegisterTool](RegisterTool.md) (13 shared connections)
- [testing.T](testing.T.md) (11 shared connections)
- [client.go](client.go.md) (4 shared connections)
- [init](init.md) (3 shared connections)
- [GetStringArg](GetStringArg.md) (3 shared connections)
- [ChatMessage](ChatMessage.md) (2 shared connections)
- [context.Context](context.Context.md) (2 shared connections)
- [RunAgentSessionLoop](RunAgentSessionLoop.md) (2 shared connections)
- [CallCommandCodeWithTools](CallCommandCodeWithTools.md) (2 shared connections)
- [runSubagent](runSubagent.md) (2 shared connections)
- [registerMCPToolsAsNative](registerMCPToolsAsNative.md) (2 shared connections)
- [LoadMCPConfig](LoadMCPConfig.md) (1 shared connections)

## Source Files

- `mcp/client.go`
- `models/model_router.go`
- `models/tools_test.go`
- `registry/dynamic.go`
- `registry/dynamic_test.go`
- `registry/plugin.go`
- `registry/plugin_test.go`
- `registry/registry.go`
- `tools/deferred.go`
- `tools/deferred_activation_test.go`

## Audit Trail

- EXTRACTED: 124 (80%)
- INFERRED: 31 (20%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*