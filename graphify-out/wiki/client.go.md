# client.go

> 75 nodes · cohesion 0.05

## Key Concepts

- **client.go** (33 connections) — `mcp/client.go`
- **MCPServer** (15 connections) — `mcp/client.go`
- **registry.go** (13 connections) — `registry/registry.go`
- **RegisterTool()** (12 connections) — `registry/registry.go`
- **encoding/json.RawMessage** (9 connections)
- **handleMCPRequest()** (8 connections) — `mcp/client.go`
- **LoadMCPConfig()** (8 connections) — `mcp/client.go`
- **registerMCPToolsAsNative()** (8 connections) — `mcp/client.go`
- **StartMCPServers()** (8 connections) — `mcp/client.go`
- **ReloadMCPServers()** (7 connections) — `mcp/client.go`
- **ExecuteMCPManage()** (7 connections) — `mcp/manage.go`
- **mcpManageAdd()** (7 connections) — `mcp/manage.go`
- **startMCPServer()** (6 connections) — `mcp/client.go`
- **.sendRequest()** (6 connections) — `mcp/client.go`
- **MCPTool** (6 connections) — `mcp/client.go`
- **ToolDef** (6 connections) — `registry/registry.go`
- **sanitizeMCPName()** (5 connections) — `mcp/client.go`
- **mcp/manage.go** (5 connections) — `mcp/manage.go`
- **mcpManageRemove()** (5 connections) — `mcp/manage.go`
- **serviceBridgeRequests()** (5 connections) — `tools/exec_code.go`
- **RegisterAutonomous()** (4 connections) — `bootstrap/autonomous.go`
- **getExposedTools()** (4 connections) — `mcp/client.go`
- **GetMCPTools()** (4 connections) — `mcp/client.go`
- **sendMCPError()** (4 connections) — `mcp/client.go`
- **StartMCPServerMode()** (4 connections) — `mcp/client.go`
- *... and 50 more nodes in this community*

## Relationships

- [GetStringArg](GetStringArg.md) (10 shared connections)
- [handleAction](handleAction.md) (7 shared connections)
- [TruncateStr](TruncateStr.md) (7 shared connections)
- [runSubagent](runSubagent.md) (3 shared connections)
- [ConfigMgr](ConfigMgr.md) (2 shared connections)
- [chat.go](chat.go.md) (2 shared connections)
- [RunAgentLoop](RunAgentLoop.md) (2 shared connections)
- [startCLI](startCLI.md) (1 shared connections)
- [collector_coolify.go](collector_coolify.go.md) (1 shared connections)
- [time.Time](time.Time.md) (1 shared connections)
- [testing.T](testing.T.md) (1 shared connections)

## Source Files

- `bootstrap/autonomous.go`
- `bootstrap/browser.go`
- `bootstrap/monitor.go`
- `bootstrap/provider.go`
- `bootstrap/script.go`
- `internal/helpers/helpers.go`
- `mcp/client.go`
- `mcp/manage.go`
- `registry/registry.go`
- `tools/exec_code.go`

## Audit Trail

- EXTRACTED: 173 (95%)
- INFERRED: 10 (5%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*