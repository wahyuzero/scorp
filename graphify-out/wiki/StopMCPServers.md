# StopMCPServers

> 4 nodes

## Key Concepts

- **StopMCPServers()** (5 connections) — `mcp/client.go`
- **.Close()** (4 connections) — `mcp/client.go`
- **.reapWait()** (2 connections) — `mcp/client.go`
- **StopWatchdogs()** (2 connections) — `mcp/watchdog.go`

## Relationships

- [LoadMCPConfig](LoadMCPConfig.md) (2 shared connections)
- [MCPServer](MCPServer.md) (2 shared connections)
- [StartDaemon](StartDaemon.md) (1 shared connections)
- [client.go](client.go.md) (1 shared connections)
- [registerMCPToolsAsNative](registerMCPToolsAsNative.md) (1 shared connections)

## Source Files

- `mcp/client.go`
- `mcp/watchdog.go`

## Audit Trail

- EXTRACTED: 9 (90%)
- INFERRED: 1 (10%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*