# client.go

> 23 nodes

## Key Concepts

- **client.go** (35 connections) — `mcp/client.go`
- **encoding/json.RawMessage** (10 connections)
- **handleMCPRequest()** (8 connections) — `mcp/client.go`
- **jsonRPCResponse** (6 connections) — `mcp/client.go`
- **StartMCPServerMode()** (5 connections) — `mcp/client.go`
- **ExecuteToolByName()** (5 connections) — `registry/registry.go`
- **getExposedTools()** (4 connections) — `mcp/client.go`
- **GetMCPTools()** (4 connections) — `mcp/client.go`
- **sendMCPError()** (4 connections) — `mcp/client.go`
- **startMCPServerMode()** (4 connections) — `mcp/client.go`
- **MCPConfig** (4 connections) — `mcp/client.go`
- **executeMCPServerTool()** (3 connections) — `mcp/client.go`
- **sendMCPResult()** (3 connections) — `mcp/client.go`
- **mcpRequest** (3 connections) — `mcp/client.go`
- **mcpResponse** (3 connections) — `mcp/client.go`
- **MCPServerModeConfig** (3 connections) — `mcp/client.go`
- **ACPRequest** (2 connections) — `delegate/acp.go`
- **MCPToolsForPrompt()** (2 connections) — `mcp/client.go`
- **StopMCPServerMode()** (2 connections) — `mcp/client.go`
- **jsonRPCError** (2 connections) — `mcp/client.go`
- **mcpError** (2 connections) — `mcp/client.go`
- **MCPToolsSummary()** (1 connections) — `mcp/client.go`
- **jsonRPCRequest** (1 connections) — `mcp/client.go`

## Relationships

- [MCPServer](MCPServer.md) (11 shared connections)
- [LoadMCPConfig](LoadMCPConfig.md) (8 shared connections)
- [registry/registry.go](registry-registry.go.md) (4 shared connections)
- [registerMCPToolsAsNative](registerMCPToolsAsNative.md) (3 shared connections)
- [runSubagent](runSubagent.md) (2 shared connections)
- [StartDaemon](StartDaemon.md) (2 shared connections)
- [ChatMessage](ChatMessage.md) (1 shared connections)
- [collector_coolify.go](collector_coolify.go.md) (1 shared connections)
- [StopMCPServers](StopMCPServers.md) (1 shared connections)
- [runPlanningTurns](runPlanningTurns.md) (1 shared connections)
- [startCLI](startCLI.md) (1 shared connections)
- [.listenSSEStream](listenSSEStream.md) (1 shared connections)

## Source Files

- `delegate/acp.go`
- `mcp/client.go`
- `registry/registry.go`

## Audit Trail

- EXTRACTED: 76 (99%)
- INFERRED: 1 (1%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*