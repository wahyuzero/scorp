# registerMCPToolsAsNative

> 13 nodes

## Key Concepts

- **registerMCPToolsAsNative()** (12 connections) — `mcp/client.go`
- **RestartServer()** (7 connections) — `mcp/watchdog.go`
- **.monitor()** (6 connections) — `mcp/watchdog.go`
- **RegisterWatchdog()** (6 connections) — `mcp/watchdog.go`
- **ServerWatchdog** (5 connections) — `mcp/watchdog.go`
- **watchdog.go** (5 connections) — `mcp/watchdog.go`
- **GetServerHealthStatus()** (4 connections) — `mcp/watchdog.go`
- **TruncOutputTool()** (3 connections) — `internal/helpers/helpers.go`
- **buildArgDefsFromInputSchema()** (3 connections) — `mcp/client.go`
- **MCPToolsDeferred()** (3 connections) — `mcp/client.go`
- **TestMCPToolsDeferredEnvParsing()** (3 connections) — `mcp/deferred_test.go`
- **MCPServer** (3 connections)
- **deferred_test.go** (1 connections) — `mcp/deferred_test.go`

## Relationships

- [LoadMCPConfig](LoadMCPConfig.md) (6 shared connections)
- [client.go](client.go.md) (3 shared connections)
- [MCPServer](MCPServer.md) (3 shared connections)
- [registry/registry.go](registry-registry.go.md) (2 shared connections)
- [testing.T](testing.T.md) (2 shared connections)
- [startCLI](startCLI.md) (2 shared connections)
- [HandleTelegramAction](HandleTelegramAction.md) (2 shared connections)
- [TruncateStr](TruncateStr.md) (1 shared connections)
- [GetIntArg](GetIntArg.md) (1 shared connections)
- [RegisterTool](RegisterTool.md) (1 shared connections)
- [RecordToolReceipt](RecordToolReceipt.md) (1 shared connections)
- [time.Time](time.Time.md) (1 shared connections)

## Source Files

- `internal/helpers/helpers.go`
- `mcp/client.go`
- `mcp/deferred_test.go`
- `mcp/watchdog.go`

## Audit Trail

- EXTRACTED: 34 (77%)
- INFERRED: 10 (23%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*